package util

import tftypes "github.com/hashicorp/terraform-plugin-framework/types"

// IsStringAttributeNonEmpty returns true if the string value is not null and not empty.
func IsStringAttributeNonEmpty(s tftypes.String) bool {
	return !s.IsNull() && s.ValueString() != ""
}
