package exchangepolicy

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	tftypes "github.com/hashicorp/terraform-plugin-framework/types"
)

func resourceSchema() schema.Schema {
	return schema.Schema{
		MarkdownDescription: "Manages a Cofide Connect exchange policy. Exchange policies govern Credex token exchanges within a trust zone by specifying match conditions on inbound tokens and determining whether exchanges are permitted or denied.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The ID of the exchange policy.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"org_id": schema.StringAttribute{
				Description: "The ID of the organization. Derived from the trust zone by Cofide Connect.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"trust_zone_id": schema.StringAttribute{
				Description: "The ID of the trust zone to which this policy applies. Cannot be changed after creation.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"name": schema.StringAttribute{
				Description: "The name of the exchange policy.",
				Required:    true,
			},
			"action": schema.StringAttribute{
				Description: "Action to take when all conditions match. One of `ALLOW`, or `DENY`. Defaults to ALLOW when unset.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf(
						"ALLOW",
						"DENY",
					),
				},
			},
			"subject_identity": stringSetResourceAttribute("Match conditions on the subject identity of the inbound token."),
			"subject_issuer":   stringSetResourceAttribute("Match conditions on the issuer of the inbound subject token."),
			"actor_identity":   stringSetResourceAttribute("Match conditions on the actor identity of the inbound token."),
			"actor_issuer":     stringSetResourceAttribute("Match conditions on the issuer of the inbound actor token."),
			"client_id":        stringSetResourceAttribute("Match conditions on the OAuth client_id presenting the exchange request."),
			"target_audience":  stringSetResourceAttribute("Match conditions on the requested target audience."),
			"outbound_scopes": schema.ListAttribute{
				Description: "Outbound scopes to grant. Only relevant when action is allow.",
				Optional:    true,
				Computed:    true,
				ElementType: tftypes.StringType,
				PlanModifiers: []planmodifier.List{
					listplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func stringSetResourceAttribute(description string) schema.Attribute {
	return schema.ListNestedAttribute{
		Description: description,
		Optional:    true,
		NestedObject: schema.NestedAttributeObject{
			Attributes: map[string]schema.Attribute{
				"exact": schema.StringAttribute{
					Description: "Exact string match.",
					Optional:    true,
					Validators: []validator.String{
						stringvalidator.ExactlyOneOf(
							path.MatchRelative().AtParent().AtName("exact"),
							path.MatchRelative().AtParent().AtName("glob"),
						),
					},
				},
				"glob": schema.StringAttribute{
					Description: "Glob pattern match (e.g. `spiffe://trust.domain/ns/*/sa/*`).",
					Optional:    true,
					Validators: []validator.String{
						stringvalidator.ExactlyOneOf(
							path.MatchRelative().AtParent().AtName("exact"),
							path.MatchRelative().AtParent().AtName("glob"),
						),
					},
				},
			},
		},
	}
}

func (r *ExchangePolicyResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = resourceSchema()
}
