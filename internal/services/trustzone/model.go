package trustzone

import "github.com/hashicorp/terraform-plugin-framework/types"

type TrustZoneModel struct {
	ID                    types.String `tfsdk:"id"`
	Name                  types.String `tfsdk:"name"`
	TrustDomain           types.String `tfsdk:"trust_domain"`
	OrgID                 types.String `tfsdk:"org_id"`
	IsManagementZone      types.Bool   `tfsdk:"is_management_zone"`
	BundleEndpointURL     types.String `tfsdk:"bundle_endpoint_url"`
	BundleEndpointProfile types.String `tfsdk:"bundle_endpoint_profile"`
	JWTIssuer             types.String `tfsdk:"jwt_issuer"`
}
