package cluster

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/structpb"

	trustproviderpb "github.com/cofide/cofide-api-sdk/gen/go/proto/trust_provider/v1alpha1"
)

func TestNewTrustProvider(t *testing.T) {
	tests := []struct {
		name          string
		kind          string
		want          *trustproviderpb.TrustProvider
		wantErr       bool
		wantErrString string
	}{
		{
			name: "kubernetes",
			kind: "kubernetes",
			want: &trustproviderpb.TrustProvider{Kind: ptrOf("kubernetes")},
		},
		{
			name:          "invalid",
			kind:          "invalid",
			want:          nil,
			wantErr:       true,
			wantErrString: "invalid trust provider kind: invalid",
		},
		{
			name:          "empty",
			kind:          "",
			want:          nil,
			wantErr:       true,
			wantErrString: "invalid trust provider kind: ",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := newTrustProvider(tt.kind)
			if !tt.wantErr {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
				require.ErrorContains(t, err, tt.wantErrString)
			}

			assert.Equal(t, tt.want, got)
		})
	}
}

func TestParseExtraHelmValues(t *testing.T) {
	tests := []struct {
		name          string
		valueStr      string
		want          map[string]interface{}
		wantErr       bool
		wantErrString string
	}{
		{
			name:     "valid",
			valueStr: `{"key": "value"}`,
			want: map[string]interface{}{
				"key": "value",
			},
		},
		{
			name:     "valid empty",
			valueStr: `{}`,
			want:     map[string]interface{}{},
		},
		{
			name:     "empty",
			valueStr: "",
			want:     nil,
		},
		{
			name:          "invalid",
			valueStr:      "invalid",
			want:          nil,
			wantErr:       true,
			wantErrString: "invalid YAML in extra_helm_values: yaml: unmarshal errors:\n  line 1: cannot unmarshal !!str `invalid` into map[string]interface {}",
		},
		{
			name:     "nested objects",
			valueStr: `{"config": {"nested": "value", "number": 42}}`,
			want: map[string]interface{}{
				"config": map[string]interface{}{
					"nested": "value",
					"number": float64(42),
				},
			},
		},
		{
			name:     "array values",
			valueStr: `{"list": ["item1", "item2"]}`,
			want: map[string]interface{}{
				"list": []interface{}{"item1", "item2"},
			},
		},
		{
			name:     "multiple types",
			valueStr: `{"string": "value", "number": 42, "bool": true, "null": null}`,
			want: map[string]interface{}{
				"string": "value",
				"number": float64(42),
				"bool":   true,
				"null":   nil,
			},
		},
		{
			name:          "unclosed object",
			valueStr:      `{"key": "value"`,
			want:          nil,
			wantErr:       true,
			wantErrString: "invalid YAML in extra_helm_values: yaml: line 1: did not find expected ',' or '}'",
		},
		{
			name:          "invalid json with multiple errors",
			valueStr:      `{"key": value, "unclosed": "string}`,
			want:          nil,
			wantErr:       true,
			wantErrString: "invalid YAML in extra_helm_values: yaml: found unexpected end of stream",
		},
		{
			name:     "whitespace",
			valueStr: `  {  "key"  :  "value"  }  `,
			want: map[string]interface{}{
				"key": "value",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseExtraHelmValues(types.StringValue(tt.valueStr))
			if !tt.wantErr {
				require.NoError(t, err)
				if tt.want == nil {
					assert.Nil(t, got)
				} else {
					require.NotNil(t, got)
					assert.Equal(t, tt.want, got.AsMap())
				}
			} else {
				require.Error(t, err)
				require.ErrorContains(t, err, tt.wantErrString)
			}
		})
	}
}

func TestValidateHelmValues(t *testing.T) {
	tests := []struct {
		name           string
		planValues     string
		responseValues *structpb.Struct
		wantErr        bool
		wantErrString  string
	}{
		{
			name:       "valid",
			planValues: `{"key": "value"}`,
			responseValues: &structpb.Struct{
				Fields: map[string]*structpb.Value{
					"key": structpb.NewStringValue("value"),
				},
			},
		},
		{
			name:       "valid, nested values",
			planValues: `{"config": {"nested": "value", "number": 42}}`,
			responseValues: &structpb.Struct{
				Fields: map[string]*structpb.Value{
					"config": structpb.NewStructValue(&structpb.Struct{
						Fields: map[string]*structpb.Value{
							"nested": structpb.NewStringValue("value"),
							"number": structpb.NewNumberValue(42),
						},
					}),
				},
			},
		},
		{
			name:       "empty plan values, populated response values",
			planValues: "",
			responseValues: &structpb.Struct{
				Fields: map[string]*structpb.Value{
					"key": structpb.NewStringValue("value"),
				},
			},
			wantErr:       true,
			wantErrString: "invalid Helm values: plan is empty while the response is not",
		},
		{
			name:           "populated plan values, empty response values",
			planValues:     `{"key": "value"}`,
			responseValues: nil,
			wantErr:        true,
			wantErrString:  "invalid Helm values: plan is not empty while the response is nil",
		},
		{
			name:           "empty plan values, empty response values",
			planValues:     "",
			responseValues: nil,
		},
		{
			name:       "case sensitive keys don't match",
			planValues: `{"extraEnv": {"name": "test"}}`,
			responseValues: &structpb.Struct{
				Fields: map[string]*structpb.Value{
					"extraenv": structpb.NewStructValue(&structpb.Struct{
						Fields: map[string]*structpb.Value{
							"name": structpb.NewStringValue("test"),
						},
					}),
				},
			},
			wantErr:       true,
			wantErrString: "Helm values mismatch: plan values don't match response values",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateHelmValues(types.StringValue(tt.planValues), tt.responseValues)
			if !tt.wantErr {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
				require.ErrorContains(t, err, tt.wantErrString)
			}
		})
	}
}

func ptrOf[T any](v T) *T {
	return &v
}
