package cloudaccount

import (
	"github.com/cofide/terraform-provider-cofide/internal/services/cloudprovider"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type CloudAccountModel struct {
	ID                  types.String     `tfsdk:"id"`
	OrgID               types.String     `tfsdk:"org_id"`
	CloudOrganizationID types.String     `tfsdk:"cloud_organization_id"`
	Name                types.String     `tfsdk:"name"`
	AWS                 *AWSAccountModel `tfsdk:"aws"`
	Suppressed          types.Bool       `tfsdk:"suppressed"`
	ManagedByDiscovery  types.Bool       `tfsdk:"managed_by_discovery"`
	CreatedAt           types.String     `tfsdk:"created_at"`
	LastUpdatedAt       types.String     `tfsdk:"last_updated_at"`
}

type AWSAccountModel struct {
	AccountID                types.String             `tfsdk:"account_id"`
	LambdaDiscoveryConfig    *AWSDiscoveryConfigModel `tfsdk:"lambda_discovery_config"`
	AgentCoreDiscoveryConfig *AWSDiscoveryConfigModel `tfsdk:"agent_core_discovery_config"`
}

// AWSDiscoveryConfigModel is shared by lambda_discovery_config and
// agent_core_discovery_config, which have identical shapes in the proto
// (AWSLambdaDiscoveryConfig and AWSAgentCoreDiscoveryConfig).
type AWSDiscoveryConfigModel struct {
	Audience                types.String                   `tfsdk:"audience"`
	Regions                 types.List                     `tfsdk:"regions"`
	DiscoveryEnabled        types.Bool                     `tfsdk:"discovery_enabled"`
	RoleChain               []cloudprovider.RoleChainModel `tfsdk:"role_chain"`
	DiscoveryInterval       types.String                   `tfsdk:"discovery_interval"`
	Status                  types.String                   `tfsdk:"status"`
	LastSuccessfulDiscovery types.String                   `tfsdk:"last_successful_discovery"`
	StatusLastUpdatedAt     types.String                   `tfsdk:"status_last_updated_at"`
}

type CloudAccountsDataSourceModel struct {
	OrgIDs               types.List          `tfsdk:"org_ids"`
	CloudOrganizationIDs types.List          `tfsdk:"cloud_organization_ids"`
	IncludeSuppressed    types.Bool          `tfsdk:"include_suppressed"`
	CloudAccounts        []CloudAccountModel `tfsdk:"cloud_accounts"`
}
