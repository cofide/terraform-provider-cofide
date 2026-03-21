package cluster

import (
	"context"

	"github.com/cofide/terraform-provider-cofide/internal/planmodifiers"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ resource.ResourceWithConfigValidators = (*ClusterResource)(nil)

func ResourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		MarkdownDescription: "Manages a Cofide Connect cluster. A cluster represents a Kubernetes cluster registered with a trust zone.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The ID of the cluster.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description: "The name of the cluster.",
				Required:    true,
			},
			"org_id": schema.StringAttribute{
				Description: "The ID of the organisation.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					planmodifiers.OptionalComputedModifier{},
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"trust_zone_id": schema.StringAttribute{
				Description: "The ID of the associated trust zone.",
				Required:    true,
			},
			"kubernetes_context": schema.StringAttribute{
				Description: "The Kubernetes context of the cluster.",
				Required:    true,
			},
			"trust_provider": schema.SingleNestedAttribute{
				Description: "The trust provider of the cluster.",
				Required:    true,
				Attributes: map[string]schema.Attribute{
					"kind": schema.StringAttribute{
						Description: "The kind of trust provider. Currently only `kubernetes` is supported.",
						Required:    true,
						Validators: []validator.String{
							stringvalidator.OneOf("kubernetes"),
						},
					},
					"k8s_psat_config": schema.SingleNestedAttribute{
						Description: "Configuration for the k8s PSAT node attestor plugin.",
						Optional:    true,
						Attributes: map[string]schema.Attribute{
							"enabled": schema.BoolAttribute{
								Description: "Whether to enable the k8s PSAT node attestor plugin with a Connect datasource.",
								Required:    true,
							},
							"allowed_service_accounts": schema.ListNestedAttribute{
								Description: "Service accounts whose tokens agents may use to attest nodes in this cluster.",
								Optional:    true,
								NestedObject: schema.NestedAttributeObject{
									Attributes: map[string]schema.Attribute{
										"namespace": schema.StringAttribute{
											Description: "The namespace of the service account.",
											Required:    true,
										},
										"service_account_name": schema.StringAttribute{
											Description: "The name of the service account.",
											Required:    true,
										},
									},
								},
							},
							"allowed_node_label_keys": schema.ListAttribute{
								Description: "Node label keys that may be used as selectors in this cluster.",
								Optional:    true,
								ElementType: types.StringType,
							},
							"allowed_pod_label_keys": schema.ListAttribute{
								Description: "Pod label keys that may be used as selectors in this cluster.",
								Optional:    true,
								ElementType: types.StringType,
							},
							"api_server_ca_cert": schema.StringAttribute{
								Description: "Base64-encoded CA certificate of the cluster's API server.",
								Optional:    true,
							},
							"api_server_url": schema.StringAttribute{
								Description: "URL of the cluster's API server.",
								Optional:    true,
							},
							"api_server_tls_server_name": schema.StringAttribute{
								Description: "Alternative TLS server name to verify the API server certificate against.",
								Optional:    true,
							},
							"api_server_proxy_url": schema.StringAttribute{
								Description: "Proxy URL for the cluster's API server.",
								Optional:    true,
							},
							"spire_server_audience": schema.StringAttribute{
								Description: "Audience the SPIRE server uses in the JWT presented to the cluster's API server.",
								Optional:    true,
							},
						},
					},
				},
			},
			"extra_helm_values": schema.StringAttribute{
				Description: "Additional Helm values for the Cofide agent Helm chart installation, in YAML format. Use `yamlencode()` to generate from a Terraform map.",
				Optional:    true,
			},
			"profile": schema.StringAttribute{
				Description: "The Cofide profile used by the cluster (e.g. `kubernetes`, `istio`).",
				Required:    true,
			},
			"external_server": schema.BoolAttribute{
				Description: "Whether the SPIRE server runs externally to this cluster. Set to `true` for clusters that delegate to a centralized SPIRE server.",
				Required:    true,
			},
			"oidc_issuer_url": schema.StringAttribute{
				Description: "The OIDC issuer URL of the cluster.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					planmodifiers.OptionalComputedModifier{},
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"oidc_issuer_ca_cert": schema.StringAttribute{
				Description: "The CA certificate (base64-encoded) to validate the cluster's OIDC issuer URL. Use `base64encode(file(...))` to supply a PEM certificate file.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					planmodifiers.OptionalComputedModifier{},
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (c *ClusterResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = ResourceSchema(ctx)
}

func (c *ClusterResource) ConfigValidators(_ context.Context) []resource.ConfigValidator {
	return []resource.ConfigValidator{}
}
