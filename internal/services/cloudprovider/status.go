package cloudprovider

import (
	"strings"

	cloudproviderpb "github.com/cofide/cofide-api-sdk/gen/go/proto/cloud_provider/v1alpha1"
	tftypes "github.com/hashicorp/terraform-plugin-framework/types"
)

const discoveryStatusPrefix = "DISCOVERY_STATUS_"

// DiscoveryStatusToString converts a DiscoveryStatus enum to its trimmed
// string representation (e.g. "DISCOVERING", "ERROR", "UNSPECIFIED").
func DiscoveryStatusToString(status cloudproviderpb.DiscoveryStatus) tftypes.String {
	return tftypes.StringValue(strings.TrimPrefix(status.String(), discoveryStatusPrefix))
}
