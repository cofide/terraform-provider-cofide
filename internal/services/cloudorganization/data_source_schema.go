package cloudorganization

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	tftypes "github.com/hashicorp/terraform-plugin-framework/types"
)

var _ datasource.DataSource = &CloudOrganizationDataSource{}
var _ datasource.DataSource = &CloudOrganizationsDataSource{}

func cloudOrganizationNestedAttributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Description: "The ID of the cloud organization.",
			Computed:    true,
		},
		"org_id": schema.StringAttribute{
			Description: "The ID of the Cofide organization this cloud organization belongs to.",
			Computed:    true,
		},
		"name": schema.StringAttribute{
			Description: "The name of the cloud organization.",
			Computed:    true,
		},
		"aws": schema.SingleNestedAttribute{
			Description: "AWS-specific configuration for the cloud organization.",
			Computed:    true,
			Attributes: map[string]schema.Attribute{
				"aws_org_id": schema.StringAttribute{
					Description: "The AWS Organization ID (e.g. `o-xxxxxxxxxx`).",
					Computed:    true,
				},
				"audience": schema.StringAttribute{
					Description: "Audience value for the initial SPIFFE JWT-based assume role call. Only used when `assume_through_oidc` is true.",
					Computed:    true,
				},
				"assume_through_oidc": schema.BoolAttribute{
					Description: "Whether the first role in `role_chain` is assumed via SPIFFE JWT-based AssumeRoleWithWebIdentity or via ambient credentials such as EKS Pod Identity.",
					Computed:    true,
				},
				"role_chain": schema.ListNestedAttribute{
					Description: "Ordered chain of IAM roles to assume when discovering resources in this cloud organization.",
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
			},
		},
		"discovery_enabled": schema.BoolAttribute{
			Description: "Whether discovery is enabled for this cloud organization.",
			Computed:    true,
		},
		"discovery_interval": schema.StringAttribute{
			Description: "How frequently discovery runs for this cloud organization (a Go duration string, e.g. `1m`, `1h`).",
			Computed:    true,
		},
		"status": schema.StringAttribute{
			Description: "The current status of cloud resource discovery (e.g. `DISCOVERING`, `ERROR`).",
			Computed:    true,
		},
		"last_discovered_at": schema.StringAttribute{
			Description: "The timestamp (RFC3339) of the last successful discovery.",
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
	attrs := cloudOrganizationNestedAttributes()
	attrs["name"] = schema.StringAttribute{
		Description: "The name of the cloud organization.",
		Required:    true,
	}
	attrs["org_id"] = schema.StringAttribute{
		Description: "The ID of the Cofide organization this cloud organization belongs to.",
		Optional:    true,
	}

	return schema.Schema{
		MarkdownDescription: "Provides information about a Cofide Connect cloud organization.",
		Attributes:          attrs,
	}
}

func (d *CloudOrganizationDataSource) Schema(ctx context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = DataSourceSchema(ctx)
}

func ListDataSourceSchema(_ context.Context) schema.Schema {
	return schema.Schema{
		MarkdownDescription: "Provides information about Cofide Connect cloud organizations.",
		Attributes: map[string]schema.Attribute{
			"org_ids": schema.ListAttribute{
				Description: "Filter by Cofide organization IDs.",
				Optional:    true,
				ElementType: tftypes.StringType,
			},
			"cloud_organizations": schema.ListNestedAttribute{
				Description: "The list of cloud organizations.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: cloudOrganizationNestedAttributes(),
				},
			},
		},
	}
}

func (d *CloudOrganizationsDataSource) Schema(ctx context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = ListDataSourceSchema(ctx)
}
