package cloudprovider

import (
	"testing"

	cloudproviderpb "github.com/cofide/cofide-api-sdk/gen/go/proto/cloud_provider/v1alpha1"
	tftypes "github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
)

func TestDiscoveryStatusToString(t *testing.T) {
	tests := []struct {
		name   string
		status cloudproviderpb.DiscoveryStatus
		want   tftypes.String
	}{
		{
			name:   "unspecified",
			status: cloudproviderpb.DiscoveryStatus_DISCOVERY_STATUS_UNSPECIFIED,
			want:   tftypes.StringValue("UNSPECIFIED"),
		},
		{
			name:   "discovering",
			status: cloudproviderpb.DiscoveryStatus_DISCOVERY_STATUS_DISCOVERING,
			want:   tftypes.StringValue("DISCOVERING"),
		},
		{
			name:   "error",
			status: cloudproviderpb.DiscoveryStatus_DISCOVERY_STATUS_ERROR,
			want:   tftypes.StringValue("ERROR"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, DiscoveryStatusToString(tt.status))
		})
	}
}
