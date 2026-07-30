package cloudaccount

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	tftypes "github.com/hashicorp/terraform-plugin-framework/types"
)

var _ datasource.DataSource = &CloudAccountDataSource{}
var _ datasource.DataSource = &CloudAccountsDataSource{}

func awsDiscoveryConfigNestedAttributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"audience": schema.StringAttribute{
			Description: "Audience value for the initial SPIFFE JWT-based assume role call.",
			Computed:    true,
		},
		"regions": schema.ListAttribute{
			Description: "AWS regions to discover resources in.",
			Computed:    true,
			ElementType: tftypes.StringType,
		},
		"discovery_enabled": schema.BoolAttribute{
			Description: "Whether discovery is enabled for this config.",
			Computed:    true,
		},
		"role_chain": schema.ListNestedAttribute{
			Description: "Ordered chain of IAM roles to assume.",
			Computed:    true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: map[string]schema.Attribute{
					"iam_role_arn": schema.StringAttribute{
						Description: "ARN of the IAM role to assume.",
						Computed:    true,
					},
					"external_id": schema.StringAttribute{
						Description: "Optional external ID for confused-deputy protection.",
						Computed:    true,
					},
				},
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

func cloudAccountNestedAttributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Description: "The ID of the cloud account.",
			Computed:    true,
		},
		"org_id": schema.StringAttribute{
			Description: "The ID of the Cofide organization this cloud account belongs to.",
			Computed:    true,
		},
		"cloud_organization_id": schema.StringAttribute{
			Description: "The ID of the cloud organization this cloud account belongs to, if any.",
			Computed:    true,
		},
		"name": schema.StringAttribute{
			Description: "The name of the cloud account.",
			Computed:    true,
		},
		"aws": schema.SingleNestedAttribute{
			Description: "AWS-specific configuration for the cloud account.",
			Computed:    true,
			Attributes: map[string]schema.Attribute{
				"account_id": schema.StringAttribute{
					Description: "The 12-digit AWS account ID.",
					Computed:    true,
				},
				"lambda_discovery_config": schema.SingleNestedAttribute{
					Description: "Configuration for Lambda function discovery.",
					Computed:    true,
					Attributes:  awsDiscoveryConfigNestedAttributes(),
				},
				"agent_core_discovery_config": schema.SingleNestedAttribute{
					Description: "Configuration for Amazon Bedrock AgentCore Runtime discovery.",
					Computed:    true,
					Attributes:  awsDiscoveryConfigNestedAttributes(),
				},
			},
		},
		"suppressed": schema.BoolAttribute{
			Description: "When true, discovery is suspended for this account and existing discovered resources are hidden from findings.",
			Computed:    true,
		},
		"managed_by_discovery": schema.BoolAttribute{
			Description: "When true, the existence of this cloud account in Connect is managed by the discovery mechanism.",
			Computed:    true,
		},
		"created_at": schema.StringAttribute{
			Description: "The timestamp (RFC3339) of resource creation.",
			Computed:    true,
		},
		"last_updated_at": schema.StringAttribute{
			Description: "The timestamp (RFC3339) of the last resource update.",
			Computed:    true,
		},
	}
}

func DataSourceSchema(_ context.Context) schema.Schema {
	attrs := cloudAccountNestedAttributes()
	attrs["name"] = schema.StringAttribute{
		Description: "The name of the cloud account.",
		Required:    true,
	}
	attrs["org_id"] = schema.StringAttribute{
		Description: "The ID of the Cofide organization this cloud account belongs to.",
		Optional:    true,
	}

	return schema.Schema{
		MarkdownDescription: "Provides information about a Cofide Connect cloud account.",
		Attributes:          attrs,
	}
}

func (d *CloudAccountDataSource) Schema(ctx context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = DataSourceSchema(ctx)
}

func ListDataSourceSchema(_ context.Context) schema.Schema {
	return schema.Schema{
		MarkdownDescription: "Provides information about Cofide Connect cloud accounts.",
		Attributes: map[string]schema.Attribute{
			"org_ids": schema.ListAttribute{
				Description: "Filter by Cofide organization IDs.",
				Optional:    true,
				ElementType: tftypes.StringType,
			},
			"cloud_organization_ids": schema.ListAttribute{
				Description: "Filter by cloud organization IDs.",
				Optional:    true,
				ElementType: tftypes.StringType,
			},
			"include_suppressed": schema.BoolAttribute{
				Description: "When true, suppressed accounts are included in the response. Defaults to false.",
				Optional:    true,
			},
			"cloud_accounts": schema.ListNestedAttribute{
				Description: "The list of cloud accounts.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: cloudAccountNestedAttributes(),
				},
			},
		},
	}
}

func (d *CloudAccountsDataSource) Schema(ctx context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = ListDataSourceSchema(ctx)
}
