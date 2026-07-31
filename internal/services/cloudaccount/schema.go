package cloudaccount

import (
	"context"

	"github.com/cofide/terraform-provider-cofide/internal/planmodifiers"
	"github.com/cofide/terraform-provider-cofide/internal/services/cloudprovider"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	tftypes "github.com/hashicorp/terraform-plugin-framework/types"
)

// awsDiscoveryConfigAttributes returns the schema shared by lambda_discovery_config
// and agent_core_discovery_config, which have identical shapes in the proto.
func awsDiscoveryConfigAttributes(roleChainDescription string) map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"audience": schema.StringAttribute{
			Description: "Audience value for the initial SPIFFE JWT-based assume role call. Only used when `assume_through_oidc` is true.",
			Required:    true,
		},
		"assume_through_oidc": schema.BoolAttribute{
			Description: "Whether the first role in `role_chain` is assumed via SPIFFE JWT-based AssumeRoleWithWebIdentity (the default) or via ambient credentials such as EKS Pod Identity (`false`).",
			Optional:    true,
			Computed:    true,
			Default:     booldefault.StaticBool(true),
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
		"role_chain": cloudprovider.RoleChainAttribute(roleChainDescription),
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
	}
}

func ResourceSchema(_ context.Context) schema.Schema {
	return schema.Schema{
		MarkdownDescription: "Manages a Cofide Connect cloud account. A cloud account represents a cloud provider account (e.g. an AWS account) linked to a Cofide organization and is the target for workload discovery.\n\n" +
			"The `aws.lambda_discovery_config`, `aws.agent_core_discovery_config`, and `suppressed` fields may either be managed inline here, or independently via the `cofide_connect_aws_lambda_discovery_config`, `cofide_connect_aws_agent_core_discovery_config`, and `cofide_connect_cloud_account_suppression_config` resources respectively. This is useful because a cloud account's existence may be controlled by automatic discovery rather than Terraform: omit these fields here to leave them under external management. Do not manage the same field both inline and via a standalone resource.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The ID of the cloud account.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"org_id": schema.StringAttribute{
				Description: "The ID of the Cofide organization this cloud account belongs to. Cannot be changed after creation.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"cloud_organization_id": schema.StringAttribute{
				Description: "The ID of the cloud organization this cloud account belongs to, if any. Cannot be changed after creation.",
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"name": schema.StringAttribute{
				Description: "The name of the cloud account.",
				Required:    true,
			},
			"aws": schema.SingleNestedAttribute{
				Description: "AWS-specific configuration for the cloud account.",
				Required:    true,
				Attributes: map[string]schema.Attribute{
					"account_id": schema.StringAttribute{
						Description: "The 12-digit AWS account ID. Cannot be changed after creation.",
						Required:    true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.RequiresReplace(),
						},
					},
					"lambda_discovery_config": schema.SingleNestedAttribute{
						Description: "Configuration for Lambda function discovery. May instead be managed with the standalone `cofide_connect_aws_lambda_discovery_config` resource.",
						Optional:    true,
						Computed:    true,
						PlanModifiers: []planmodifier.Object{
							planmodifiers.OptionalComputedModifier{},
							objectplanmodifier.UseStateForUnknown(),
						},
						Attributes: awsDiscoveryConfigAttributes("Ordered chain of IAM roles to assume when discovering Lambda functions."),
					},
					"agent_core_discovery_config": schema.SingleNestedAttribute{
						Description: "Configuration for Amazon Bedrock AgentCore Runtime discovery. May instead be managed with the standalone `cofide_connect_aws_agent_core_discovery_config` resource.",
						Optional:    true,
						Computed:    true,
						PlanModifiers: []planmodifier.Object{
							planmodifiers.OptionalComputedModifier{},
							objectplanmodifier.UseStateForUnknown(),
						},
						Attributes: awsDiscoveryConfigAttributes("Ordered chain of IAM roles to assume when discovering AgentCore runtimes."),
					},
				},
			},
			"suppressed": schema.BoolAttribute{
				Description: "When true, discovery is suspended for this account and existing discovered resources are hidden from findings. May instead be managed with the standalone `cofide_connect_cloud_account_suppression_config` resource.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					planmodifiers.OptionalComputedModifier{},
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"managed_by_discovery": schema.BoolAttribute{
				Description: "When true, the existence of this cloud account in Connect is managed by the discovery mechanism: if a subsequent discovery run for the parent cloud organization does not find this account, it is deleted. Set to false to manually manage the existence of this resource outside of cloud organization level discovery.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
			"created_at": schema.StringAttribute{
				Description: "The timestamp (RFC3339) of resource creation.",
				Computed:    true,
			},
			"last_updated_at": schema.StringAttribute{
				Description: "The timestamp (RFC3339) of the last resource update.",
				Computed:    true,
			},
		},
	}
}

func (c *CloudAccountResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = ResourceSchema(ctx)
}
