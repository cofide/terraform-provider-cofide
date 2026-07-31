package cloudorganization

import (
	"github.com/cofide/terraform-provider-cofide/internal/services/cloudprovider"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type CloudOrganizationModel struct {
	ID                types.String          `tfsdk:"id"`
	OrgID             types.String          `tfsdk:"org_id"`
	Name              types.String          `tfsdk:"name"`
	AWS               *AWSOrganizationModel `tfsdk:"aws"`
	DiscoveryEnabled  types.Bool            `tfsdk:"discovery_enabled"`
	DiscoveryInterval types.String          `tfsdk:"discovery_interval"`
	Status            types.String          `tfsdk:"status"`
	LastDiscoveredAt  types.String          `tfsdk:"last_discovered_at"`
	CreatedAt         types.String          `tfsdk:"created_at"`
	LastUpdatedAt     types.String          `tfsdk:"last_updated_at"`
}

type AWSOrganizationModel struct {
	AWSOrgID          types.String                   `tfsdk:"aws_org_id"`
	Audience          types.String                   `tfsdk:"audience"`
	AssumeThroughOidc types.Bool                     `tfsdk:"assume_through_oidc"`
	RoleChain         []cloudprovider.RoleChainModel `tfsdk:"role_chain"`
}

type CloudOrganizationsDataSourceModel struct {
	OrgIDs             types.List               `tfsdk:"org_ids"`
	CloudOrganizations []CloudOrganizationModel `tfsdk:"cloud_organizations"`
}
