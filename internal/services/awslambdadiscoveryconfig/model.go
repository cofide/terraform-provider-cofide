package awslambdadiscoveryconfig

import (
	"github.com/cofide/terraform-provider-cofide/internal/services/cloudprovider"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type Model struct {
	ID                      types.String                   `tfsdk:"id"`
	CloudAccountID          types.String                   `tfsdk:"cloud_account_id"`
	Audience                types.String                   `tfsdk:"audience"`
	Regions                 types.List                     `tfsdk:"regions"`
	DiscoveryEnabled        types.Bool                     `tfsdk:"discovery_enabled"`
	RoleChain               []cloudprovider.RoleChainModel `tfsdk:"role_chain"`
	DiscoveryInterval       types.String                   `tfsdk:"discovery_interval"`
	Status                  types.String                   `tfsdk:"status"`
	LastSuccessfulDiscovery types.String                   `tfsdk:"last_successful_discovery"`
	StatusLastUpdatedAt     types.String                   `tfsdk:"status_last_updated_at"`
}
