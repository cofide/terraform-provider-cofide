package awsagentcorediscoveryconfig

import (
	"context"

	"github.com/cofide/terraform-provider-cofide/internal/planmodifiers"
	"github.com/cofide/terraform-provider-cofide/internal/services/cloudprovider"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	tftypes "github.com/hashicorp/terraform-plugin-framework/types"
)

func ResourceSchema(_ context.Context) schema.Schema {
	return schema.Schema{
		MarkdownDescription: "Manages the Amazon Bedrock AgentCore Runtime discovery configuration of a Cofide Connect cloud account, independently of the `cofide_connect_cloud_account` resource. Useful when the cloud account's existence is managed by automatic discovery rather than Terraform: only this specific configuration is managed here, referencing the account by ID. Do not manage this configuration both here and inline on `cofide_connect_cloud_account.aws.agent_core_discovery_config`.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The ID of this resource. Matches `cloud_account_id`.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"cloud_account_id": schema.StringAttribute{
				Description: "The ID of the cloud account this configuration applies to. Cannot be changed after creation.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"audience": schema.StringAttribute{
				Description: "Audience value for the initial SPIFFE JWT-based assume role call.",
				Required:    true,
			},
			"regions": schema.ListAttribute{
				Description: "AWS regions to discover resources in.",
				Optional:    true,
				ElementType: tftypes.StringType,
			},
			"discovery_enabled": schema.BoolAttribute{
				Description: "Whether discovery is enabled for this config.",
				Optional:    true,
			},
			"role_chain": cloudprovider.RoleChainAttribute("Ordered chain of IAM roles to assume when discovering AgentCore runtimes."),
			"discovery_interval": schema.StringAttribute{
				Description: "How frequently discovery runs for this config (a Go duration string, e.g. `1m`, `1h`). Defaults to a server-assigned value when unset.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					planmodifiers.OptionalComputedModifier{},
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"status": schema.StringAttribute{
				Description: "The current status of cloud resource discovery (e.g. `DISCOVERING`, `ERROR`).",
				Computed:    true,
			},
			"last_successful_discovery": schema.StringAttribute{
				Description: "The timestamp (RFC3339) of the last successful discovery run.",
				Computed:    true,
			},
			"status_last_updated_at": schema.StringAttribute{
				Description: "The timestamp (RFC3339) of the last status update.",
				Computed:    true,
			},
		},
	}
}

func (r *Resource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = ResourceSchema(ctx)
}
