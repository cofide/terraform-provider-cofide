package cloudprovider

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	tftypes "github.com/hashicorp/terraform-plugin-framework/types"
	"google.golang.org/protobuf/types/known/durationpb"
)

// DurationToString converts a protobuf duration to a Go duration string
// (e.g. "1m0s"), returning a null string if nil.
func DurationToString(d *durationpb.Duration) tftypes.String {
	if d == nil {
		return tftypes.StringNull()
	}
	return tftypes.StringValue(d.AsDuration().String())
}

// DurationToProto parses a Go duration string (e.g. "1m", "30s") into a
// protobuf duration, returning nil if the string is null.
func DurationToProto(s tftypes.String) (*durationpb.Duration, diag.Diagnostics) {
	var diags diag.Diagnostics

	if s.IsNull() {
		return nil, diags
	}

	d, err := time.ParseDuration(s.ValueString())
	if err != nil {
		diags.AddError(
			"Invalid discovery_interval",
			fmt.Sprintf("Could not parse discovery_interval %q as a duration: %s", s.ValueString(), err),
		)
		return nil, diags
	}

	return durationpb.New(d), diags
}
