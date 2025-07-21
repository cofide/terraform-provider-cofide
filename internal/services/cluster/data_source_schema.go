package cluster

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

var _ datasource.DataSource = &ClusterDataSource{}

func DataSourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		MarkdownDescription: "Provides information about a cluster resource.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The ID of the cluster.",
				Computed:    true,
			},
			"name": schema.StringAttribute{
				Description: "The name of the cluster.",
				Required:    true,
			},
			"org_id": schema.StringAttribute{
				Description: "The ID of the organisation.",
				Required:    true,
			},
			"trust_zone_id": schema.StringAttribute{
				Description: "The ID of the associated trust zone.",
				Optional:    true,
			},
			"kubernetes_context": schema.StringAttribute{
				Description: "The Kubernetes context of the cluster.",
				Computed:    true,
			},
			"trust_provider": schema.SingleNestedAttribute{
				Description: "The trust provider of the cluster.",
				Computed:    true,
				Attributes: map[string]schema.Attribute{
					"kind": schema.StringAttribute{
						Description: "The kind of trust provider.",
						Computed:    true,
					},
				},
			},
			"extra_helm_values": schema.StringAttribute{
				Description: "The extra Helm values to provide to the cluster.",
				Computed:    true,
			},
			"profile": schema.StringAttribute{
				Description: "The Cofide profile used by the cluster.",
				Computed:    true,
			},
			"external_server": schema.BoolAttribute{
				Description: "Whether or not the SPIRE server runs externally.",
				Computed:    true,
			},
			"oidc_issuer_url": schema.StringAttribute{
				Description: "The OIDC issuer URL of the cluster.",
				Computed:    true,
			},
			"oidc_issuer_ca_cert": schema.StringAttribute{
				Description: "The OIDC issuer CA certificate of the cluster.",
				Computed:    true,
				Sensitive:   true,
			},
		},
	}
}

func (c *ClusterDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = DataSourceSchema(ctx)
}

func (c *ClusterDataSource) ConfigValidators(_ context.Context) []datasource.ConfigValidator {
	return []datasource.ConfigValidator{}
}
