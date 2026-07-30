package cloudprovider

import (
	"testing"
	"time"

	tftypes "github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/durationpb"
)

func TestDurationToString(t *testing.T) {
	tests := []struct {
		name string
		d    *durationpb.Duration
		want tftypes.String
	}{
		{
			name: "nil",
			d:    nil,
			want: tftypes.StringNull(),
		},
		{
			name: "one minute",
			d:    durationpb.New(time.Minute),
			want: tftypes.StringValue("1m0s"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, DurationToString(tt.d))
		})
	}
}

func TestDurationToProto(t *testing.T) {
	t.Run("null", func(t *testing.T) {
		got, diags := DurationToProto(tftypes.StringNull())
		require.False(t, diags.HasError())
		assert.Nil(t, got)
	})

	t.Run("valid duration", func(t *testing.T) {
		got, diags := DurationToProto(tftypes.StringValue("1m"))
		require.False(t, diags.HasError())
		assert.Equal(t, durationpb.New(time.Minute), got)
	})

	t.Run("invalid duration", func(t *testing.T) {
		got, diags := DurationToProto(tftypes.StringValue("not-a-duration"))
		assert.True(t, diags.HasError())
		assert.Nil(t, got)
	})
}
