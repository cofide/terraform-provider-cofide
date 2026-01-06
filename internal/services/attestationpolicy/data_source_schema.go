package attestationpolicy

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	tftypes "github.com/hashicorp/terraform-plugin-framework/types"
)

var _ datasource.DataSource = &AttestationPolicyDataSource{}

func DataSourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		MarkdownDescription: "Provides information about an attestation policy resource.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The ID of the attestation policy.",
				Computed:    true,
			},
			"name": schema.StringAttribute{
				Description: "The name of the attestation policy.",
				Required:    true,
			},
			"org_id": schema.StringAttribute{
				Description: "The ID of the organisation.",
				Required:    true,
			},
			"kubernetes": schema.SingleNestedAttribute{
				Description: "The configuration of the Kubernetes attestation policy.",
				Computed:    true,
				Attributes: map[string]schema.Attribute{
					"namespace_selector": schema.SingleNestedAttribute{
						Description: "The configuration of the namespace selector for the Kubernetes attestation policy.",
						Computed:    true,
						Attributes: map[string]schema.Attribute{
							"match_labels": schema.MapAttribute{
								Description: "The list of labels to match for the namespace selector.",
								Computed:    true,
								ElementType: tftypes.StringType,
							},
							"match_expressions": schema.ListNestedAttribute{
								Description: "The list of match expressions for the namespace selector.",
								Computed:    true,
								NestedObject: schema.NestedAttributeObject{
									Attributes: map[string]schema.Attribute{
										"key": schema.StringAttribute{
											Description: "The key of the match expression.",
											Computed:    true,
										},
										"operator": schema.StringAttribute{
											Description: "The operator of the match expression.",
											Computed:    true,
										},
										"values": schema.ListAttribute{
											Description: "The values of the match expression.",
											Computed:    true,
											ElementType: tftypes.StringType,
										},
									},
								},
							},
						},
					},
					"pod_selector": schema.SingleNestedAttribute{
						Description: "The configuration of the pod selector for the Kubernetes attestation policy.",
						Computed:    true,
						Attributes: map[string]schema.Attribute{
							"match_labels": schema.MapAttribute{
								Description: "The list of labels to match for the pod selector.",
								Computed:    true,
								ElementType: tftypes.StringType,
							},
							"match_expressions": schema.ListNestedAttribute{
								Description: "The list of match expressions for the pod selector.",
								Computed:    true,
								NestedObject: schema.NestedAttributeObject{
									Attributes: map[string]schema.Attribute{
										"key": schema.StringAttribute{
											Description: "The key of the match expression.",
											Computed:    true,
										},
										"operator": schema.StringAttribute{
											Description: "The operator of the match expression.",
											Computed:    true,
										},
										"values": schema.ListAttribute{
											Description: "The values of the match expression.",
											Computed:    true,
											ElementType: tftypes.StringType,
										},
									},
								},
							},
						},
					},
					"dns_name_templates": schema.ListAttribute{
						Description: "The list of DNS name templates for the Kubernetes attestation policy.",
						Computed:    true,
						ElementType: tftypes.StringType,
					},
					"spiffe_id_path_template": schema.StringAttribute{
						Description: "The SPIFFE ID path template for the Kubernetes attestation policy.",
						Computed:    true,
					},
				},
			},
			"static": schema.SingleNestedAttribute{
				Description: "The configuration of the static attestation policy.",
				Computed:    true,
				Attributes: map[string]schema.Attribute{
					"spiffe_id_path": schema.StringAttribute{
						Description: "The SPIFFE ID path for the static attestation policy.",
						Computed:    true,
					},
					"parent_id_path": schema.StringAttribute{
						Description: "The parent ID path for the static attestation policy.",
						Computed:    true,
					},
					"selectors": schema.ListNestedAttribute{
						Description: "The list of selectors for the static attestation policy.",
						Computed:    true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"type": schema.StringAttribute{
									Description: "The type of the selector.",
									Computed:    true,
								},
								"value": schema.StringAttribute{
									Description: "The value of the selector.",
									Computed:    true,
								},
							},
						},
					},
					"dns_names": schema.ListAttribute{
						Description: "The list of DNS names for the static attestation policy.",
						Computed:    true,
						ElementType: tftypes.StringType,
					},
				},
			},
		},
	}
}

func (a *AttestationPolicyDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = DataSourceSchema(ctx)
}

func (a *AttestationPolicyDataSource) ConfigValidators(_ context.Context) []datasource.ConfigValidator {
	return []datasource.ConfigValidator{}
}
