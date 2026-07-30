package cloudprovider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	tftypes "github.com/hashicorp/terraform-plugin-framework/types"
)

// StringListToProto converts a Terraform string list to a []string, returning
// nil (rather than an error) if the list is null or unknown.
func StringListToProto(ctx context.Context, list tftypes.List) ([]string, diag.Diagnostics) {
	if list.IsNull() || list.IsUnknown() {
		return nil, nil
	}

	var values []string
	diags := list.ElementsAs(ctx, &values, false)
	return values, diags
}

// StringListFromProto converts a []string to a Terraform string list, using
// a null list (rather than an empty one) when there are no values.
func StringListFromProto(values []string) tftypes.List {
	if len(values) == 0 {
		return tftypes.ListNull(tftypes.StringType)
	}

	elems := make([]attr.Value, 0, len(values))
	for _, v := range values {
		elems = append(elems, tftypes.StringValue(v))
	}
	return tftypes.ListValueMust(tftypes.StringType, elems)
}
