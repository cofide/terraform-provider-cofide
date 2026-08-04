package exchangepolicy

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	tftypes "github.com/hashicorp/terraform-plugin-framework/types"
)

// exactlyOneVariantValidator enforces that exactly one variant attribute of an
// object representing a protobuf oneof is non-null. It lives on the object
// itself so it fires even when no variant is set (unlike per-variant
// ExactlyOneOf validators, which only fire when the attribute they are attached
// to is non-null).
//
// When adding a new variant: add it as Optional inside the object's Attributes
// map. This validator requires no modification.
type exactlyOneVariantValidator struct {
	// name is the attribute name of the object, used in diagnostics.
	name string
}

func (v exactlyOneVariantValidator) Description(_ context.Context) string {
	return fmt.Sprintf("Exactly one %s variant must be set.", v.name)
}

func (v exactlyOneVariantValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v exactlyOneVariantValidator) ValidateObject(_ context.Context, req validator.ObjectRequest, resp *validator.ObjectResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}
	set := 0
	for _, attr := range req.ConfigValue.Attributes() {
		if !attr.IsNull() && !attr.IsUnknown() {
			set++
		}
	}
	if set != 1 {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			fmt.Sprintf("Invalid %s configuration", v.name),
			fmt.Sprintf("Exactly one %s variant must be set, but got %d.", v.name, set),
		)
	}
}

// oauthAsValidator enforces that at least one of issuer_url or token_url is
// configured. It lives on the oauth_as object rather than on the two URL
// attributes because a per-attribute AtLeastOneOf validator only fires when the
// attribute it is attached to is non-null, so it would not catch both being
// omitted. grant_type is not checked here: it is Required in the schema, so
// Terraform enforces it whenever oauth_as is set.
//
// Unknown values count as set, since their eventual value is not knowable at
// validation time.
type oauthAsValidator struct{}

func (oauthAsValidator) Description(_ context.Context) string {
	return "At least one of issuer_url or token_url must be set."
}

func (oauthAsValidator) MarkdownDescription(ctx context.Context) string {
	return oauthAsValidator{}.Description(ctx)
}

func (oauthAsValidator) ValidateObject(_ context.Context, req validator.ObjectRequest, resp *validator.ObjectResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}
	attrs := req.ConfigValue.Attributes()
	isSet := func(name string) bool {
		attr, ok := attrs[name]
		return ok && !attr.IsNull()
	}
	if !isSet("issuer_url") && !isSet("token_url") {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Invalid oauth_as configuration",
			"At least one of issuer_url or token_url is required when oauth_as is set.",
		)
	}
}

func ResourceSchema() schema.Schema {
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
			"subject_audience": stringSetResourceAttribute("Match conditions on the audience claim of the inbound subject token."),
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
			"outbound_issuer": schema.SingleNestedAttribute{
				Description: "Outbound token issuer configuration. When set, Credex will obtain an outbound token from this issuer rather than minting one itself. At most one outbound_issuer variant can be set.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.UseStateForUnknown(),
				},
				Validators: []validator.Object{
					// exactlyOneVariantValidator fires on the outbound_issuer object itself
					// so it catches both empty blocks and multiple variants being set.
					// When adding a new issuer variant, add it as Optional below — no change here.
					exactlyOneVariantValidator{name: "outbound_issuer"},
				},
				Attributes: map[string]schema.Attribute{
					"oauth_as": schema.SingleNestedAttribute{
						Description: "Use an external OAuth 2.0 authorisation server as the outbound issuer. At least one of `issuer_url` or `token_url` is required.",
						Optional:    true,
						Validators: []validator.Object{
							oauthAsValidator{},
						},
						Attributes: map[string]schema.Attribute{
							"grant_type": schema.StringAttribute{
								Description: "OAuth 2.0 grant type.",
								Required:    true,
							},
							"issuer_url": schema.StringAttribute{
								Description: "Issuer URL of the OAuth 2.0 authorisation server.",
								Optional:    true,
								Computed:    true,
								PlanModifiers: []planmodifier.String{
									stringplanmodifier.UseStateForUnknown(),
								},
							},
							"token_url": schema.StringAttribute{
								Description: "Token endpoint URL of the OAuth 2.0 authorisation server.",
								Optional:    true,
								Computed:    true,
								PlanModifiers: []planmodifier.String{
									stringplanmodifier.UseStateForUnknown(),
								},
							},
							"audiences": schema.ListAttribute{
								Description: "Audiences to request in the outbound token.",
								Optional:    true,
								Computed:    true,
								ElementType: tftypes.StringType,
								PlanModifiers: []planmodifier.List{
									listplanmodifier.UseStateForUnknown(),
								},
							},
							"timeout": schema.Int64Attribute{
								Description: "Timeout for token requests to the authorisation server, in seconds.",
								Optional:    true,
							},
						},
					},
				},
			},
			"outbound_identity": schema.StringAttribute{
				Description: "Outbound identity to assert in the exchanged token. When set, Credex will use this identity rather than the inbound subject identity.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"external_hooks": schema.ListNestedAttribute{
				Description: "Post-matching hooks that transform outbound token claims before Credex mints them.",
				Optional:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Description: "Name of the hook, unique within the policy.",
							Required:    true,
						},
						"description": schema.StringAttribute{
							Description: "Optional description of the hook.",
							Optional:    true,
						},
						"url": schema.StringAttribute{
							Description: "URL of the external hook endpoint.",
							Required:    true,
						},
						"auth": schema.SingleNestedAttribute{
							Description: "Authentication configuration for the hook endpoint. Exactly one auth variant must be set.",
							Required:    true,
							Validators: []validator.Object{
								// exactlyOneVariantValidator fires on the auth object itself so it
								// catches both empty auth blocks and multiple variants being set.
								// When adding a new auth variant, add it as Optional below — no change here.
								exactlyOneVariantValidator{name: "auth"},
							},
							Attributes: map[string]schema.Attribute{
								"spiffe_mtls": schema.SingleNestedAttribute{
									Description: "Authenticate to the hook using SPIFFE mTLS.",
									Optional:    true,
									Attributes: map[string]schema.Attribute{
										"spiffe_id": schema.StringAttribute{
											Description: "SPIFFE ID to present when connecting to the hook endpoint.",
											Required:    true,
										},
									},
								},
							},
						},
						"timeout": schema.Int64Attribute{
							Description: "Timeout for the hook request, in seconds.",
							Optional:    true,
						},
					},
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
	resp.Schema = ResourceSchema()
}
