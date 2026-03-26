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

func TestTrustProviderRoundTrip(t *testing.T) {
	tests := []struct {
		name  string
		model *TrustProviderModel
	}{
		{
			name: "kind only",
			model: &TrustProviderModel{
				Kind:          types.StringValue("kubernetes"),
				K8sPsatConfig: nil,
			},
		},
		{
			name: "with k8s_psat_config disabled",
			model: &TrustProviderModel{
				Kind: types.StringValue("kubernetes"),
				K8sPsatConfig: &K8sPsatConfigModel{
					Enabled:                types.BoolValue(false),
					AllowedServiceAccounts: nil,
					AllowedNodeLabelKeys:   nil,
					AllowedPodLabelKeys:    nil,
					ApiServerCaCert:        types.StringNull(),
					ApiServerURL:           types.StringNull(),
					ApiServerTLSServerName: types.StringNull(),
					ApiServerProxyURL:      types.StringNull(),
					SpireServerAudience:    types.StringNull(),
				},
			},
		},
		{
			name: "with minimal k8s_psat_config",
			model: &TrustProviderModel{
				Kind: types.StringValue("kubernetes"),
				K8sPsatConfig: &K8sPsatConfigModel{
					Enabled:                types.BoolValue(true),
					AllowedServiceAccounts: nil,
					AllowedNodeLabelKeys:   nil,
					AllowedPodLabelKeys:    nil,
					ApiServerCaCert:        types.StringNull(),
					ApiServerURL:           types.StringNull(),
					ApiServerTLSServerName: types.StringNull(),
					ApiServerProxyURL:      types.StringNull(),
					SpireServerAudience:    types.StringNull(),
				},
			},
		},
		{
			name: "with full k8s_psat_config",
			model: &TrustProviderModel{
				Kind: types.StringValue("kubernetes"),
				K8sPsatConfig: &K8sPsatConfigModel{
					Enabled: types.BoolValue(true),
					AllowedServiceAccounts: []ServiceAccountModel{
						{
							Namespace:          types.StringValue("spire"),
							ServiceAccountName: types.StringValue("spire-agent"),
						},
						{
							Namespace:          types.StringValue("default"),
							ServiceAccountName: types.StringValue("app-agent"),
						},
					},
					AllowedNodeLabelKeys:   []types.String{types.StringValue("kubernetes.io/hostname")},
					AllowedPodLabelKeys:    []types.String{types.StringValue("app"), types.StringValue("version")},
					ApiServerCaCert:        types.StringValue("dGVzdC1jYQ=="), // base64("test-ca")
					ApiServerURL:           types.StringValue("https://kubernetes.default.svc"),
					ApiServerTLSServerName: types.StringValue("kubernetes"),
					ApiServerProxyURL:      types.StringValue("http://proxy:3128"),
					SpireServerAudience:    types.StringValue("spire-server"),
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			proto, err := trustProviderToProto(tt.model)
			require.NoError(t, err)

			got := trustProviderFromProto(proto)
			assert.Equal(t, tt.model, got)
		})
	}
}

func TestTrustProviderFromProto_Nil(t *testing.T) {
	got := trustProviderFromProto(nil)
	assert.Equal(t, types.StringValue(""), got.Kind)
	assert.Nil(t, got.K8sPsatConfig)
}

func TestTrustProviderToProto_InvalidKind(t *testing.T) {
	model := &TrustProviderModel{
		Kind:          types.StringValue("invalid"),
		K8sPsatConfig: nil,
	}
	_, err := trustProviderToProto(model)
	require.ErrorContains(t, err, "invalid trust provider kind: invalid")
}

func TestK8sPsatConfigToProto_InvalidBase64CaCert(t *testing.T) {
	model := &K8sPsatConfigModel{
		Enabled:         types.BoolValue(true),
		ApiServerCaCert: types.StringValue("not-valid-base64!!!"),
	}
	_, err := k8sPsatConfigToProto(model)
	require.ErrorContains(t, err, "failed to decode api_server_ca_cert from base64")
}


func TestK8sPsatConfigForState(t *testing.T) {
	sa := ServiceAccountModel{
		Namespace:          types.StringValue("spire"),
		ServiceAccountName: types.StringValue("spire-agent"),
	}

	tests := []struct {
		name  string
		model *K8sPsatConfigModel
		prev  *K8sPsatConfigModel
		want  *K8sPsatConfigModel
	}{
		{
			name:  "nil API service accounts, nil prev → nil",
			model: &K8sPsatConfigModel{Enabled: types.BoolValue(true), AllowedServiceAccounts: nil},
			prev:  &K8sPsatConfigModel{Enabled: types.BoolValue(true), AllowedServiceAccounts: nil},
			want:  &K8sPsatConfigModel{Enabled: types.BoolValue(true), AllowedServiceAccounts: nil},
		},
		{
			name:  "nil API service accounts, empty prev → empty preserved",
			model: &K8sPsatConfigModel{Enabled: types.BoolValue(true), AllowedServiceAccounts: nil},
			prev:  &K8sPsatConfigModel{Enabled: types.BoolValue(true), AllowedServiceAccounts: []ServiceAccountModel{}},
			want:  &K8sPsatConfigModel{Enabled: types.BoolValue(true), AllowedServiceAccounts: []ServiceAccountModel{}},
		},
		{
			name:  "non-empty API service accounts override empty prev",
			model: &K8sPsatConfigModel{Enabled: types.BoolValue(true), AllowedServiceAccounts: []ServiceAccountModel{sa}},
			prev:  &K8sPsatConfigModel{Enabled: types.BoolValue(true), AllowedServiceAccounts: []ServiceAccountModel{}},
			want:  &K8sPsatConfigModel{Enabled: types.BoolValue(true), AllowedServiceAccounts: []ServiceAccountModel{sa}},
		},
		{
			name:  "nil API service accounts, non-empty prev → nil (API removal wins)",
			model: &K8sPsatConfigModel{Enabled: types.BoolValue(true), AllowedServiceAccounts: nil},
			prev:  &K8sPsatConfigModel{Enabled: types.BoolValue(true), AllowedServiceAccounts: []ServiceAccountModel{sa}},
			want:  &K8sPsatConfigModel{Enabled: types.BoolValue(true), AllowedServiceAccounts: nil},
		},
		{
			name:  "nil API node label keys, nil prev → nil",
			model: &K8sPsatConfigModel{Enabled: types.BoolValue(true), AllowedNodeLabelKeys: nil},
			prev:  &K8sPsatConfigModel{Enabled: types.BoolValue(true), AllowedNodeLabelKeys: nil},
			want:  &K8sPsatConfigModel{Enabled: types.BoolValue(true), AllowedNodeLabelKeys: nil},
		},
		{
			name:  "nil API node label keys, empty prev → empty preserved",
			model: &K8sPsatConfigModel{Enabled: types.BoolValue(true), AllowedNodeLabelKeys: nil},
			prev:  &K8sPsatConfigModel{Enabled: types.BoolValue(true), AllowedNodeLabelKeys: []types.String{}},
			want:  &K8sPsatConfigModel{Enabled: types.BoolValue(true), AllowedNodeLabelKeys: []types.String{}},
		},
		{
			name:  "nil API pod label keys, empty prev → empty preserved",
			model: &K8sPsatConfigModel{Enabled: types.BoolValue(true), AllowedPodLabelKeys: nil},
			prev:  &K8sPsatConfigModel{Enabled: types.BoolValue(true), AllowedPodLabelKeys: []types.String{}},
			want:  &K8sPsatConfigModel{Enabled: types.BoolValue(true), AllowedPodLabelKeys: []types.String{}},
		},
		{
			name:  "nil API node label keys, non-empty prev → nil (API removal wins)",
			model: &K8sPsatConfigModel{Enabled: types.BoolValue(true), AllowedNodeLabelKeys: nil},
			prev:  &K8sPsatConfigModel{Enabled: types.BoolValue(true), AllowedNodeLabelKeys: []types.String{types.StringValue("kubernetes.io/hostname")}},
			want:  &K8sPsatConfigModel{Enabled: types.BoolValue(true), AllowedNodeLabelKeys: nil},
		},
		{
			name:  "nil API pod label keys, non-empty prev → nil (API removal wins)",
			model: &K8sPsatConfigModel{Enabled: types.BoolValue(true), AllowedPodLabelKeys: nil},
			prev:  &K8sPsatConfigModel{Enabled: types.BoolValue(true), AllowedPodLabelKeys: []types.String{types.StringValue("app")}},
			want:  &K8sPsatConfigModel{Enabled: types.BoolValue(true), AllowedPodLabelKeys: nil},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := k8sPsatConfigForState(tt.model, tt.prev)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestTrustProviderForState(t *testing.T) {
	tests := []struct {
		name  string
		model *TrustProviderModel
		prev  *TrustProviderModel
		want  *TrustProviderModel
	}{
		{
			name: "nil prev → model returned as-is",
			model: &TrustProviderModel{
				Kind:          types.StringValue("kubernetes"),
				K8sPsatConfig: &K8sPsatConfigModel{Enabled: types.BoolValue(true), AllowedServiceAccounts: nil},
			},
			prev: nil,
			want: &TrustProviderModel{
				Kind:          types.StringValue("kubernetes"),
				K8sPsatConfig: &K8sPsatConfigModel{Enabled: types.BoolValue(true), AllowedServiceAccounts: nil},
			},
		},
		{
			name: "nil prev k8s_psat_config → model returned as-is",
			model: &TrustProviderModel{
				Kind:          types.StringValue("kubernetes"),
				K8sPsatConfig: &K8sPsatConfigModel{Enabled: types.BoolValue(true), AllowedServiceAccounts: nil},
			},
			prev: &TrustProviderModel{
				Kind:          types.StringValue("kubernetes"),
				K8sPsatConfig: nil,
			},
			want: &TrustProviderModel{
				Kind:          types.StringValue("kubernetes"),
				K8sPsatConfig: &K8sPsatConfigModel{Enabled: types.BoolValue(true), AllowedServiceAccounts: nil},
			},
		},
		{
			name: "nil model k8s_psat_config → model returned as-is",
			model: &TrustProviderModel{
				Kind:          types.StringValue("kubernetes"),
				K8sPsatConfig: nil,
			},
			prev: &TrustProviderModel{
				Kind:          types.StringValue("kubernetes"),
				K8sPsatConfig: &K8sPsatConfigModel{Enabled: types.BoolValue(true), AllowedServiceAccounts: []ServiceAccountModel{}},
			},
			want: &TrustProviderModel{
				Kind:          types.StringValue("kubernetes"),
				K8sPsatConfig: nil,
			},
		},
		{
			name: "both non-nil → list fields merged from prev",
			model: &TrustProviderModel{
				Kind: types.StringValue("kubernetes"),
				K8sPsatConfig: &K8sPsatConfigModel{
					Enabled:                types.BoolValue(true),
					AllowedServiceAccounts: nil,
					AllowedNodeLabelKeys:   nil,
					AllowedPodLabelKeys:    nil,
				},
			},
			prev: &TrustProviderModel{
				Kind: types.StringValue("kubernetes"),
				K8sPsatConfig: &K8sPsatConfigModel{
					Enabled:                types.BoolValue(true),
					AllowedServiceAccounts: []ServiceAccountModel{},
					AllowedNodeLabelKeys:   []types.String{},
					AllowedPodLabelKeys:    nil,
				},
			},
			want: &TrustProviderModel{
				Kind: types.StringValue("kubernetes"),
				K8sPsatConfig: &K8sPsatConfigModel{
					Enabled:                types.BoolValue(true),
					AllowedServiceAccounts: []ServiceAccountModel{},
					AllowedNodeLabelKeys:   []types.String{},
					AllowedPodLabelKeys:    nil,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Build a proto from the model, then call trustProviderForState with the prev.
			// This exercises the full path including trustProviderFromProto.
			proto, err := trustProviderToProto(tt.model)
			require.NoError(t, err)
			got := trustProviderForState(proto, tt.prev)
			assert.Equal(t, tt.want, got)
		})
	}
}

func ptrOf[T any](v T) *T {
	return &v
}
