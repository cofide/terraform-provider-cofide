package cluster

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ datasource.DataSource = &ClusterDataSource{}

func DataSourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		MarkdownDescription: "Provides information about a Cofide Connect cluster.",
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
				Description: "The ID of the organization.",
				Optional:    true,
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
					"k8s_psat_config": schema.SingleNestedAttribute{
						Description: "Configuration for the k8s PSAT node attestor plugin.",
						Computed:    true,
						Attributes: map[string]schema.Attribute{
							"enabled": schema.BoolAttribute{
								Description: "Whether to enable the k8s PSAT node attestor plugin with a Connect datasource.",
								Computed:    true,
							},
							"allowed_service_accounts": schema.ListNestedAttribute{
								Description: "Service accounts whose tokens agents may use to attest nodes in this cluster.",
								Computed:    true,
								NestedObject: schema.NestedAttributeObject{
									Attributes: map[string]schema.Attribute{
										"namespace": schema.StringAttribute{
											Description: "The namespace of the service account.",
											Computed:    true,
										},
										"service_account_name": schema.StringAttribute{
											Description: "The name of the service account.",
											Computed:    true,
										},
									},
								},
							},
							"allowed_node_label_keys": schema.ListAttribute{
								Description: "Node label keys that may be used as selectors in this cluster.",
								Computed:    true,
								ElementType: types.StringType,
							},
							"allowed_pod_label_keys": schema.ListAttribute{
								Description: "Pod label keys that may be used as selectors in this cluster.",
								Computed:    true,
								ElementType: types.StringType,
							},
							"api_server_ca_cert": schema.StringAttribute{
								Description: "Base64-encoded CA certificate of the cluster's API server.",
								Computed:    true,
							},
							"api_server_url": schema.StringAttribute{
								Description: "URL of the cluster's API server.",
								Computed:    true,
							},
							"api_server_tls_server_name": schema.StringAttribute{
								Description: "Alternative TLS server name to verify the API server certificate against.",
								Computed:    true,
							},
							"api_server_proxy_url": schema.StringAttribute{
								Description: "Proxy URL for the cluster's API server.",
								Computed:    true,
							},
							"spire_server_audience": schema.StringAttribute{
								Description: "Audience the SPIRE server uses in the JWT presented to the cluster's API server.",
								Computed:    true,
							},
						},
					},
				},
			},
			"extra_helm_values": schema.StringAttribute{
				Description: "Additional Helm values for the Cofide SPIRE Helm chart installation, in YAML format.",
				Computed:    true,
			},
			"profile": schema.StringAttribute{
				Description: "The Cofide profile used by the cluster (e.g. `kubernetes`, `istio`). Ensures Cofide SPIRE is configured correctly for the target environment.",
				Computed:    true,
			},
			"external_server": schema.BoolAttribute{
				Description: "Whether the SPIRE server runs externally to this cluster.",
				Computed:    true,
			},
			"oidc_issuer_url": schema.StringAttribute{
				Description: "The OIDC issuer URL of the cluster.",
				Computed:    true,
			},
			"oidc_issuer_ca_cert": schema.StringAttribute{
				Description: "The CA certificate (base64-encoded) to validate the cluster's OIDC issuer URL.",
				Computed:    true,
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
