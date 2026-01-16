package attestationpolicy

import (
	"context"
	"reflect"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	tftypes "github.com/hashicorp/terraform-plugin-framework/types"
)

var _ resource.ResourceWithConfigValidators = (*AttestationPolicyResource)(nil)

func ResourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		MarkdownDescription: "Provides an attestation policy resource.",
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
				Optional:    true,
				Computed:    true,
			},
			"kubernetes": schema.SingleNestedAttribute{
				Description: "The configuration of the Kubernetes attestation policy.",
				Optional:    true,
				Attributes: map[string]schema.Attribute{
					"namespace_selector": schema.SingleNestedAttribute{
						Description: "The configuration of the namespace selector for the Kubernetes attestation policy.",
						Optional:    true,
						Attributes: map[string]schema.Attribute{
							"match_labels": schema.MapAttribute{
								Description: "The list of labels to match for the namespace selector.",
								Optional:    true,
								ElementType: tftypes.StringType,
							},
							"match_expressions": schema.ListNestedAttribute{
								Description: "The list of match expressions for the namespace selector.",
								Optional:    true,
								NestedObject: schema.NestedAttributeObject{
									Attributes: map[string]schema.Attribute{
										"key": schema.StringAttribute{
											Description: "The key of the match expression.",
											Required:    true,
										},
										"operator": schema.StringAttribute{
											Description: "The operator of the match expression.",
											Required:    true,
										},
										"values": schema.ListAttribute{
											Description: "The values of the match expression.",
											Optional:    true,
											ElementType: tftypes.StringType,
										},
									},
								},
							},
						},
					},
					"pod_selector": schema.SingleNestedAttribute{
						Description: "The configuration of the pod selector for the Kubernetes attestation policy.",
						Optional:    true,
						Attributes: map[string]schema.Attribute{
							"match_labels": schema.MapAttribute{
								Description: "The list of labels to match for the pod selector.",
								Optional:    true,
								ElementType: tftypes.StringType,
							},
							"match_expressions": schema.ListNestedAttribute{
								Description: "The list of match expressions for the pod selector.",
								Optional:    true,
								NestedObject: schema.NestedAttributeObject{
									Attributes: map[string]schema.Attribute{
										"key": schema.StringAttribute{
											Description: "The key of the match expression.",
											Required:    true,
										},
										"operator": schema.StringAttribute{
											Description: "The operator of the match expression.",
											Required:    true,
										},
										"values": schema.ListAttribute{
											Description: "The values of the match expression.",
											Optional:    true,
											ElementType: tftypes.StringType,
										},
									},
								},
							},
						},
					},
					"dns_name_templates": schema.ListAttribute{
						Description: "The list of DNS name templates for the Kubernetes attestation policy.",
						Optional:    true,
						ElementType: tftypes.StringType,
					},
					"spiffe_id_path_template": schema.StringAttribute{
						Description: "The SPIFFE ID path template for the Kubernetes attestation policy.",
						Optional:    true,
					},
				},
			},
			"static": schema.SingleNestedAttribute{
				Description: "The configuration of the static attestation policy.",
				Optional:    true,
				Attributes: map[string]schema.Attribute{
					"spiffe_id_path": schema.StringAttribute{
						Description: "The SPIFFE ID path for the static attestation policy.",
						Required:    true,
					},
					"parent_id_path": schema.StringAttribute{
						Description: "The parent ID path for the static attestation policy.",
						Required:    true,
					},
					"selectors": schema.ListNestedAttribute{
						Description: "The list of selectors for the static attestation policy.",
						Required:    true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"type": schema.StringAttribute{
									Description: "The type of the selector.",
									Required:    true,
								},
								"value": schema.StringAttribute{
									Description: "The value of the selector.",
									Required:    true,
								},
							},
						},
					},
					"dns_names": schema.ListAttribute{
						Description: "The list of DNS names for the static attestation policy.",
						Optional:    true,
						ElementType: tftypes.StringType,
					},
				},
			},
			"tpm_node": schema.SingleNestedAttribute{
				Description: "The configuration of the TPM node attestation policy.",
				Optional:    true,
				Attributes: map[string]schema.Attribute{
					"attestation": schema.SingleNestedAttribute{
						Description: "The TPM attestation configuration.",
						Required:    true,
						Attributes: map[string]schema.Attribute{
							"ek_hash": schema.StringAttribute{
								Description: "SHA-256 hash of the Endorsement Key (EK) certificate of the TPM.",
								Required:    true,
							},
						},
					},
					"selector_values": schema.ListAttribute{
						Description: "The list of selector values for the TPM node attestation policy.",
						Optional:    true,
						ElementType: tftypes.StringType,
					},
				},
			},
		},
	}
}

func (a *AttestationPolicyResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = ResourceSchema(ctx)
}

type exactlyOneOfValidator struct{}

var _ resource.ConfigValidator = exactlyOneOfValidator{}

func (v exactlyOneOfValidator) Description(ctx context.Context) string {
	return "exactly one of 'kubernetes', 'static', or 'tpm_node' must be configured"
}

func (v exactlyOneOfValidator) MarkdownDescription(ctx context.Context) string {
	return "exactly one of `kubernetes`, `static`, or `tpm_node` must be configured"
}

func (v exactlyOneOfValidator) ValidateResource(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	var data AttestationPolicyModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	valid, reason := isExactlyOneNonNil(data.Kubernetes, data.Static, data.TPMNode)
	if !valid {
		resp.Diagnostics.AddError(
			"Invalid configuration",
			"Exactly one of kubernetes, static or tpm_node blocks must be configured, but "+reason+" were provided.",
		)
		return
	}
}

func (a *AttestationPolicyResource) ConfigValidators(_ context.Context) []resource.ConfigValidator {
	return []resource.ConfigValidator{
		exactlyOneOfValidator{},
	}
}

// isExactlyOneNonNil returns true if the value of exactly one of the arguments is non-nil.
// Otherwise, it returns false and a string reason of "none" or "multiple".
func isExactlyOneNonNil(input ...any) (bool, string) {
	count := 0
	for _, v := range input {
		// An interface value is not nil if it has a type, even if the underlying value is nil.
		// We must use reflection to check if the underlying value is nil.
		if v == nil {
			continue
		}
		if !reflect.ValueOf(v).IsNil() {
			count++
		}
	}
	if count == 0 {
		return false, "none"
	}
	if count > 1 {
		return false, "multiple"
	}
	return true, ""
}
