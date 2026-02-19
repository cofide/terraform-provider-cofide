package cluster

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
)

var _ resource.ResourceWithConfigValidators = (*ClusterResource)(nil)

func ResourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		MarkdownDescription: "Provides a cluster resource.",
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
					optionalComputedModifier{},
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
						Description: "The kind of trust provider.",
						Required:    true,
					},
				},
			},
			"extra_helm_values": schema.StringAttribute{
				Description: "The extra Helm values to provide to the cluster.",
				Optional:    true,
			},
			"profile": schema.StringAttribute{
				Description: "The Cofide profile used by the cluster.",
				Required:    true,
			},
			"external_server": schema.BoolAttribute{
				Description: "Whether or not the SPIRE server runs externally.",
				Required:    true,
			},
			"oidc_issuer_url": schema.StringAttribute{
				Description: "The OIDC issuer URL of the cluster.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					optionalComputedModifier{},
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"oidc_issuer_ca_cert": schema.StringAttribute{
				Description: "The CA certificate (base64-encoded) to validate the cluster's OIDC issuer URL.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					optionalComputedModifier{},
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
