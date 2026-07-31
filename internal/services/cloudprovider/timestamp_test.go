package cloudprovider

import (
	"testing"
	"time"

	tftypes "github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestTimestampToString(t *testing.T) {
	fixed := time.Date(2026, time.July, 30, 12, 0, 0, 0, time.UTC)

	tests := []struct {
		name string
		ts   *timestamppb.Timestamp
		want tftypes.String
	}{
		{
			name: "nil timestamp",
			ts:   nil,
			want: tftypes.StringNull(),
		},
		{
			name: "set timestamp",
			ts:   timestamppb.New(fixed),
			want: tftypes.StringValue("2026-07-30T12:00:00Z"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, TimestampToString(tt.ts))
		})
	}
}
