package trustzoneserver

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	tftypes "github.com/hashicorp/terraform-plugin-framework/types"
)

var _ datasource.DataSource = &TrustZoneServerDataSource{}
var _ datasource.DataSource = &TrustZoneServersDataSource{}

func DataSourceSchema(_ context.Context) schema.Schema {
	return schema.Schema{
		MarkdownDescription: "Provides information about a Cofide Connect trust zone server.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The ID of the trust zone server.",
				Optional:    true,
				Computed:    true,
			},
			"trust_zone_id": schema.StringAttribute{
				Description: "The ID of the trust zone managed by this server.",
				Optional:    true,
				Computed:    true,
			},
			"cluster_id": schema.StringAttribute{
				Description: "The ID of the cluster on which the server is deployed.",
				Optional:    true,
				Computed:    true,
			},
			"kubernetes_namespace": schema.StringAttribute{
				Description: "The Kubernetes namespace in which the server is deployed.",
				Computed:    true,
			},
			"kubernetes_service_account": schema.StringAttribute{
				Description: "The name of the Kubernetes service account deployed with the server.",
				Computed:    true,
			},
			"org_id": schema.StringAttribute{
				Description: "The ID of the organisation.",
				Optional:    true,
				Computed:    true,
			},
			"helm_values": schema.StringAttribute{
				Description: "Helm values configured for the server install (JSON).",
				Computed:    true,
			},
			"status": schema.SingleNestedAttribute{
				Description: "The current lifecycle status of the trust zone server.",
				Computed:    true,
				Attributes: map[string]schema.Attribute{
					"status": schema.StringAttribute{
						Description: "The status of the trust zone server (e.g. `TRUST_ZONE_SERVER_STATUS_PROVISIONED`).",
						Computed:    true,
					},
					"last_transition_time": schema.StringAttribute{
						Description: "The time of the last status transition (RFC3339).",
						Computed:    true,
					},
				},
			},
			"connect_k8s_psat_config": schema.SingleNestedAttribute{
				Description: "Configuration for the k8s PSAT node attestor plugin.",
				Computed:    true,
				Attributes: map[string]schema.Attribute{
					"audiences": schema.ListAttribute{
						Description: "Audiences that SPIRE agents in remote clusters can present for node attestation.",
						Computed:    true,
						ElementType: tftypes.StringType,
					},
					"spire_server_spiffe_id_path": schema.StringAttribute{
						Description: "SPIFFE ID path used in the JWT presented by the SPIRE server to the cluster's API server.",
						Computed:    true,
					},
				},
			},
		},
	}
}

func (d *TrustZoneServerDataSource) Schema(ctx context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = DataSourceSchema(ctx)
}

func (d *TrustZoneServerDataSource) ConfigValidators(_ context.Context) []datasource.ConfigValidator {
	return []datasource.ConfigValidator{}
}

func ListDataSourceSchema(_ context.Context) schema.Schema {
	return schema.Schema{
		MarkdownDescription: "Provides information about Cofide Connect trust zone servers.",
		Attributes: map[string]schema.Attribute{
			"trust_zone_id": schema.StringAttribute{
				Description: "Filter by trust zone ID.",
				Optional:    true,
			},
			"cluster_id": schema.StringAttribute{
				Description: "Filter by cluster ID.",
				Optional:    true,
			},
			"org_id": schema.StringAttribute{
				Description: "Filter by organisation ID.",
				Optional:    true,
			},
			"trust_zone_servers": schema.ListNestedAttribute{
				Description: "The list of trust zone servers.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Description: "The ID of the trust zone server.",
							Computed:    true,
						},
						"trust_zone_id": schema.StringAttribute{
							Description: "The ID of the trust zone managed by this server.",
							Computed:    true,
						},
						"cluster_id": schema.StringAttribute{
							Description: "The ID of the cluster on which the server is deployed.",
							Computed:    true,
						},
						"kubernetes_namespace": schema.StringAttribute{
							Description: "The Kubernetes namespace in which the server is deployed.",
							Computed:    true,
						},
						"kubernetes_service_account": schema.StringAttribute{
							Description: "The name of the Kubernetes service account deployed with the server.",
							Computed:    true,
						},
						"org_id": schema.StringAttribute{
							Description: "The ID of the organisation.",
							Computed:    true,
						},
						"helm_values": schema.StringAttribute{
							Description: "Helm values configured for the server install (JSON).",
							Computed:    true,
						},
						"status": schema.SingleNestedAttribute{
							Description: "The current lifecycle status of the trust zone server.",
							Computed:    true,
							Attributes: map[string]schema.Attribute{
								"status": schema.StringAttribute{
									Description: "The status of the trust zone server (e.g. `TRUST_ZONE_SERVER_STATUS_PROVISIONED`).",
									Computed:    true,
								},
								"last_transition_time": schema.StringAttribute{
									Description: "The time of the last status transition (RFC3339).",
									Computed:    true,
								},
							},
						},
						"connect_k8s_psat_config": schema.SingleNestedAttribute{
							Description: "Configuration for the k8s PSAT node attestor plugin.",
							Computed:    true,
							Attributes: map[string]schema.Attribute{
								"audiences": schema.ListAttribute{
									Description: "Audiences that SPIRE agents in remote clusters can present for node attestation.",
									Computed:    true,
									ElementType: tftypes.StringType,
								},
								"spire_server_spiffe_id_path": schema.StringAttribute{
									Description: "SPIFFE ID path used in the JWT presented by the SPIRE server to the cluster's API server.",
									Computed:    true,
								},
							},
						},
					},
				},
			},
		},
	}
}

func (d *TrustZoneServersDataSource) Schema(ctx context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = ListDataSourceSchema(ctx)
}

func (d *TrustZoneServersDataSource) ConfigValidators(_ context.Context) []datasource.ConfigValidator {
	return []datasource.ConfigValidator{}
}
