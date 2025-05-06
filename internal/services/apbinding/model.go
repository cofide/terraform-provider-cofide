package apbinding

import "github.com/hashicorp/terraform-plugin-framework/types"

type APBindingModel struct {
	ID          types.String               `tfsdk:"id"`
	OrgID       types.String               `tfsdk:"org_id"`
	TrustZoneID types.String               `tfsdk:"trust_zone_id"`
	PolicyID    types.String               `tfsdk:"policy_id"`
	Federations []APBindingFederationModel `tfsdk:"federations"`
}

type APBindingFederationModel struct {
	TrustZoneID types.String `tfsdk:"trust_zone_id"`
}
