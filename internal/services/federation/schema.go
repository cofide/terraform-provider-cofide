package federation

import (
	"context"

	"github.com/cofide/terraform-provider-cofide/internal/planmodifiers"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
)

var _ resource.ResourceWithConfigValidators = (*FederationResource)(nil)

func ResourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		MarkdownDescription: "Manages a Cofide Connect federation. Establishes a trust relationship between two trust zones so their workloads can mutually authenticate.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The ID of the federation.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
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
			"remote_trust_zone_id": schema.StringAttribute{
				Description: "The ID of the associated remote trust zone.",
				Required:    true,
			},
		},
	}
}

func (f *FederationResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = ResourceSchema(ctx)
}

func (f *FederationResource) ConfigValidators(_ context.Context) []resource.ConfigValidator {
	return []resource.ConfigValidator{}
}
