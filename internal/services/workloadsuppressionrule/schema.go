package workloadsuppressionrule

import (
	"context"
	"reflect"

	"github.com/cofide/terraform-provider-cofide/internal/planmodifiers"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	tftypes "github.com/hashicorp/terraform-plugin-framework/types"
)

var _ resource.ResourceWithConfigValidators = (*WorkloadSuppressionRuleResource)(nil)

func ResourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		MarkdownDescription: "Manages a Cofide Connect workload suppression rule. Suppression rules hide matching workloads from the Connect findings view. Exactly one matcher (currently only `kubernetes_pod`) must be configured.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The ID of the workload suppression rule.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"org_id": schema.StringAttribute{
				Description: "The ID of the organization. Cannot be changed after creation.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					planmodifiers.OptionalComputedModifier{},
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description: "The name of the workload suppression rule.",
				Required:    true,
			},
			"description": schema.StringAttribute{
				Description: "A human-readable explanation of why the rule exists.",
				Optional:    true,
			},
			"enabled": schema.BoolAttribute{
				Description: "Whether the rule is currently active. Disabled rules are retained but have no suppression effect.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					planmodifiers.OptionalComputedModifier{},
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"kubernetes_pod": schema.SingleNestedAttribute{
				Description: "Matches Kubernetes pod workloads.",
				Optional:    true,
				Attributes: map[string]schema.Attribute{
					"trust_zone_ids": schema.ListAttribute{
						Description: "Matches only pods within these trust zones.",
						Optional:    true,
						ElementType: tftypes.StringType,
					},
					"cluster_ids": schema.ListAttribute{
						Description: "Matches only pods within these clusters.",
						Optional:    true,
						ElementType: tftypes.StringType,
					},
					"namespaces": schema.ListAttribute{
						Description: "Matches pods in these Kubernetes namespaces.",
						Optional:    true,
						ElementType: tftypes.StringType,
					},
					"labels": schema.MapAttribute{
						Description: "Matches pods with these Kubernetes labels.",
						Optional:    true,
						ElementType: tftypes.StringType,
					},
				},
			},
			"created_at": schema.StringAttribute{
				Description: "The time the rule was created, in RFC3339 format.",
				Computed:    true,
			},
			"last_updated_at": schema.StringAttribute{
				Description: "The time the rule was last updated, in RFC3339 format.",
				Computed:    true,
			},
		},
	}
}

func (r *WorkloadSuppressionRuleResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = ResourceSchema(ctx)
}

type exactlyOneOfMatcherValidator struct{}

var _ resource.ConfigValidator = exactlyOneOfMatcherValidator{}

func (v exactlyOneOfMatcherValidator) Description(_ context.Context) string {
	return "exactly one matcher (currently only 'kubernetes_pod') must be configured"
}

func (v exactlyOneOfMatcherValidator) MarkdownDescription(_ context.Context) string {
	return "exactly one matcher (currently only `kubernetes_pod`) must be configured"
}

func (v exactlyOneOfMatcherValidator) ValidateResource(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	var data WorkloadSuppressionRuleModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	valid, reason := isExactlyOneNonNil(data.KubernetesPod)
	if !valid {
		resp.Diagnostics.AddError(
			"Invalid configuration",
			"Exactly one matcher block (kubernetes_pod) must be configured, but "+reason+" were provided.",
		)
		return
	}
}

func (r *WorkloadSuppressionRuleResource) ConfigValidators(_ context.Context) []resource.ConfigValidator {
	return []resource.ConfigValidator{
		exactlyOneOfMatcherValidator{},
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
