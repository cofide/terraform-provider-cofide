package trustzone

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
)

var _ resource.ResourceWithConfigValidators = (*TrustZoneResource)(nil)

func ResourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		MarkdownDescription: "Provides a trust zone resource.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The ID of the trust zone.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description: "The name of the trust zone.",
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
			"trust_domain": schema.StringAttribute{
				Description: "The trust domain of the trust zone.",
				Required:    true,
			},
			"is_management_zone": schema.BoolAttribute{
				Description: "Whether or not this is a management trust zone.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					optionalComputedModifier{},
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"bundle_endpoint_url": schema.StringAttribute{
				Description: "The bundle endpoint URL of the trust zone.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"bundle_endpoint_profile": schema.StringAttribute{
				Description: "The bundle endpoint profile of the trust zone.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"jwt_issuer": schema.StringAttribute{
				Description: "The JWT issuer of the trust zone.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (t *TrustZoneResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = ResourceSchema(ctx)
}

func (t *TrustZoneResource) ConfigValidators(_ context.Context) []resource.ConfigValidator {
	return []resource.ConfigValidator{}
}
