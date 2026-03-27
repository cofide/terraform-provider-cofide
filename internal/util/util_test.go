package util

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/structpb"
)

func TestHelmValuesForState(t *testing.T) {
	nonEmptyStruct, err := structpb.NewStruct(map[string]any{"key": "value"})
	require.NoError(t, err)

	tests := []struct {
		name       string
		apiValues  *structpb.Struct
		stateValue types.String
		wantNull   bool
		wantValue  string
	}{
		// Empty API response cases — the edge cases raised in code review.
		{
			name:       "nil API, null state returns null",
			apiValues:  nil,
			stateValue: types.StringNull(),
			wantNull:   true,
		},
		{
			name:       "nil API, empty string state is preserved",
			apiValues:  nil,
			stateValue: types.StringValue(""),
			wantValue:  "",
		},
		{
			name:       "nil API, empty JSON state is preserved",
			apiValues:  nil,
			stateValue: types.StringValue("{}"),
			wantValue:  "{}",
		},
		{
			name:       "empty struct API, empty JSON state is preserved",
			apiValues:  &structpb.Struct{},
			stateValue: types.StringValue("{}"),
			wantValue:  "{}",
		},
		{
			name:       "empty struct API, null state returns null",
			apiValues:  &structpb.Struct{},
			stateValue: types.StringNull(),
			wantNull:   true,
		},
		{
			// State has values but the API has none — externally removed, should show drift.
			name:       "nil API, non-empty state returns null",
			apiValues:  nil,
			stateValue: types.StringValue(`{"key": "value"}`),
			wantNull:   true,
		},
		// Non-empty API response cases.
		{
			name:       "non-empty API, null state returns API JSON",
			apiValues:  nonEmptyStruct,
			stateValue: types.StringNull(),
			wantValue:  `{"key":"value"}`,
		},
		{
			name:       "non-empty API, matching JSON state is preserved",
			apiValues:  nonEmptyStruct,
			stateValue: types.StringValue(`{"key": "value"}`),
			wantValue:  `{"key": "value"}`,
		},
		{
			name:       "non-empty API, matching YAML state is preserved",
			apiValues:  nonEmptyStruct,
			stateValue: types.StringValue("key: value\n"),
			wantValue:  "key: value\n",
		},
		{
			name:       "non-empty API, different state returns API JSON",
			apiValues:  nonEmptyStruct,
			stateValue: types.StringValue(`{"other": "value"}`),
			wantValue:  `{"key":"value"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := HelmValuesForState(tt.apiValues, tt.stateValue)
			require.NoError(t, err)
			if tt.wantNull {
				assert.True(t, got.IsNull(), "expected null, got %q", got.ValueString())
			} else {
				assert.False(t, got.IsNull())
				assert.Equal(t, tt.wantValue, got.ValueString())
			}
		})
	}
}
