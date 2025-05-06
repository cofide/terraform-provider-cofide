package federation

import "github.com/hashicorp/terraform-plugin-framework/types"

type FederationModel struct {
	ID                types.String `tfsdk:"id"`
	OrgID             types.String `tfsdk:"org_id"`
	TrustZoneID       types.String `tfsdk:"trust_zone_id"`
	RemoteTrustZoneID types.String `tfsdk:"remote_trust_zone_id"`
}
