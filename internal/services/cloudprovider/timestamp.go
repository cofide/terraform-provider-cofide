package cloudprovider

import (
	"time"

	tftypes "github.com/hashicorp/terraform-plugin-framework/types"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// TimestampToString converts a protobuf timestamp to an RFC3339-formatted
// string, returning a null string if the timestamp is nil.
func TimestampToString(ts *timestamppb.Timestamp) tftypes.String {
	if ts == nil {
		return tftypes.StringNull()
	}
	return tftypes.StringValue(ts.AsTime().Format(time.RFC3339))
}
