package trustzone

import (
	"context"

	"github.com/cofide/terraform-provider-cofide/internal/planmodifiers"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
)

var _ resource.ResourceWithConfigValidators = (*TrustZoneResource)(nil)

func ResourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		MarkdownDescription: "Manages a Cofide Connect trust zone. A trust zone contains a SPIFFE trust domain.",
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
					planmodifiers.OptionalComputedModifier{},
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"trust_domain": schema.StringAttribute{
				Description: "The SPIFFE trust domain for this trust zone (e.g. `example.cofide.dev`).",
				Required:    true,
			},
			"is_management_zone": schema.BoolAttribute{
				Description: "Whether this is a management trust zone. Cannot be changed after creation.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					planmodifiers.OptionalComputedModifier{},
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"bundle_endpoint_url": schema.StringAttribute{
				Description: "The URL of the SPIFFE bundle endpoint for this trust zone. Set by Cofide Connect.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"bundle_endpoint_profile": schema.StringAttribute{
				Description: "The SPIFFE bundle endpoint profile for this trust zone (`BUNDLE_ENDPOINT_PROFILE_HTTPS_SPIFFE` or `BUNDLE_ENDPOINT_PROFILE_HTTPS_WEB`). Set by Cofide Connect.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"jwt_issuer": schema.StringAttribute{
				Description: "The JWT issuer URL for this trust zone. Set by Cofide Connect.",
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
