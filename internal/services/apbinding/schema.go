package apbinding

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ resource.ResourceWithConfigValidators = (*APBindingResource)(nil)

func ResourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		MarkdownDescription: "Provides an attestation policy binding resource.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The ID of the attestation policy binding.",
				Computed:    true,
			},
			"org_id": schema.StringAttribute{
				Description: "The ID of the organisation.",
				Optional:    true,
				Computed:    true,
			},
			"trust_zone_id": schema.StringAttribute{
				Description: "The ID of the trust zone.",
				Required:    true,
			},
			"policy_id": schema.StringAttribute{
				Description: "The ID of the attestation policy.",
				Required:    true,
			},
			"federations": schema.ListAttribute{
				Description: "The list of associated federations.",
				Optional:    true,
				ElementType: types.ObjectType{
					AttrTypes: map[string]attr.Type{
						"trust_zone_id": types.StringType,
					},
				},
			},
		},
	}
}

func (a *APBindingResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = ResourceSchema(ctx)
}

func (a *APBindingResource) ConfigValidators(_ context.Context) []resource.ConfigValidator {
	return []resource.ConfigValidator{}
}
