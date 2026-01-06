package organization

import (
	tftypes "github.com/hashicorp/terraform-plugin-framework/types"
)

// OrganizationModel describes the data source data model.
type OrganizationModel struct {
	ID   tftypes.String `tfsdk:"id"`
	Name tftypes.String `tfsdk:"name"`
}
