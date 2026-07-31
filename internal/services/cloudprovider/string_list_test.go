package cloudprovider

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	tftypes "github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func regionsList() tftypes.List {
	return tftypes.ListValueMust(tftypes.StringType, []attr.Value{
		tftypes.StringValue("eu-west-1"),
		tftypes.StringValue("us-east-1"),
	})
}

func TestStringListToProto(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name string
		list tftypes.List
		want []string
	}{
		{
			name: "null",
			list: tftypes.ListNull(tftypes.StringType),
			want: nil,
		},
		{
			name: "unknown",
			list: tftypes.ListUnknown(tftypes.StringType),
			want: nil,
		},
		{
			name: "values",
			list: regionsList(),
			want: []string{"eu-west-1", "us-east-1"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, diags := StringListToProto(ctx, tt.list)
			require.False(t, diags.HasError())
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestStringListFromProto(t *testing.T) {
	tests := []struct {
		name   string
		values []string
		want   tftypes.List
	}{
		{
			name:   "empty",
			values: nil,
			want:   tftypes.ListNull(tftypes.StringType),
		},
		{
			name:   "values",
			values: []string{"eu-west-1", "us-east-1"},
			want:   regionsList(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, StringListFromProto(tt.values))
		})
	}
}
