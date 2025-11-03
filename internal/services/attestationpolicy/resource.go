package attestationpolicy

import (
	"context"
	"fmt"

	sdkclient "github.com/cofide/cofide-api-sdk/pkg/connect/client"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	tftypes "github.com/hashicorp/terraform-plugin-framework/types"
	spiretypes "github.com/spiffe/spire-api-sdk/proto/spire/api/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	attestationpolicypb "github.com/cofide/cofide-api-sdk/gen/go/proto/attestation_policy/v1alpha1"
)

var _ resource.Resource = &AttestationPolicyResource{}
var _ resource.ResourceWithImportState = &AttestationPolicyResource{}
var _ resource.ResourceWithValidateConfig = &AttestationPolicyResource{}

type AttestationPolicyResource struct {
	client sdkclient.ClientSet
}

func NewResource() resource.Resource {
	return &AttestationPolicyResource{}
}

func (r *AttestationPolicyResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_connect_attestation_policy"
}

func (r *AttestationPolicyResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(sdkclient.ClientSet)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected resource configure type",
			fmt.Sprintf("Expected sdkclient.ClientSet, got: %T", req.ProviderData),
		)

		return
	}

	r.client = client
}

func (r *AttestationPolicyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan AttestationPolicyModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	policy := &attestationpolicypb.AttestationPolicy{
		Name:  plan.Name.ValueString(),
		OrgId: plan.OrgID.ValueStringPointer(),
	}

	if plan.Kubernetes != nil {
		k8sPolicy := &attestationpolicypb.APKubernetes{}
		if plan.Kubernetes.NamespaceSelector != nil {
			k8sPolicy.NamespaceSelector = convertLabelSelector(plan.Kubernetes.NamespaceSelector)
		}
		if plan.Kubernetes.PodSelector != nil {
			k8sPolicy.PodSelector = convertLabelSelector(plan.Kubernetes.PodSelector)
		}
		policy.Policy = &attestationpolicypb.AttestationPolicy_Kubernetes{
			Kubernetes: k8sPolicy,
		}
	}

	if plan.Static != nil {
		staticPolicy := &attestationpolicypb.APStatic{
			SpiffeId: plan.Static.SpiffeID.ValueStringPointer(),
		}

		selectors, err := extractSelectors(ctx, plan.Static.Selectors)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error extracting selectors",
				err.Error(),
			)
			return
		}
		staticPolicy.Selectors = selectors
		policy.Policy = &attestationpolicypb.AttestationPolicy_Static{
			Static: staticPolicy,
		}
	}

	createResp, err := r.client.AttestationPolicyV1Alpha1().CreateAttestationPolicy(ctx, policy)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating attestation policy",
			fmt.Sprintf("Could not create attestation policy: %s", err.Error()),
		)
		return
	}

	plan.ID = tftypes.StringValue(createResp.GetId())
	plan.Name = tftypes.StringValue(createResp.GetName())
	plan.OrgID = tftypes.StringValue(createResp.GetOrgId())

	if k8s := createResp.GetKubernetes(); k8s != nil && plan.Kubernetes != nil {
		plan.Kubernetes = &APKubernetesModel{
			NamespaceSelector: plan.Kubernetes.NamespaceSelector,
			PodSelector:       plan.Kubernetes.PodSelector,
			DNSNameTemplates:  plan.Kubernetes.DNSNameTemplates,
		}
	}

	if static := createResp.GetStatic(); static != nil && plan.Static != nil {
		plan.Static = &APStaticModel{
			SpiffeID:  tftypes.StringValue(static.GetSpiffeId()),
			Selectors: plan.Static.Selectors,
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *AttestationPolicyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state AttestationPolicyModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *AttestationPolicyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var state AttestationPolicyModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var plan AttestationPolicyModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	policyID := state.ID.ValueString()
	if policyID == "" {
		resp.Diagnostics.AddError(
			"Error updating attestation policy",
			"Policy ID not found in state. The resource might not have been created properly.",
		)
		return
	}

	policy := &attestationpolicypb.AttestationPolicy{
		Id:    &policyID,
		Name:  plan.Name.ValueString(),
		OrgId: plan.OrgID.ValueStringPointer(),
	}

	if plan.Kubernetes != nil {
		k8sPolicy := &attestationpolicypb.APKubernetes{}
		if plan.Kubernetes.NamespaceSelector != nil {
			k8sPolicy.NamespaceSelector = convertLabelSelector(plan.Kubernetes.NamespaceSelector)
		}
		if plan.Kubernetes.PodSelector != nil {
			k8sPolicy.PodSelector = convertLabelSelector(plan.Kubernetes.PodSelector)
		}
		policy.Policy = &attestationpolicypb.AttestationPolicy_Kubernetes{
			Kubernetes: k8sPolicy,
		}
	}

	if plan.Static != nil {
		staticPolicy := &attestationpolicypb.APStatic{
			SpiffeId: plan.Static.SpiffeID.ValueStringPointer(),
		}
		selectors, err := extractSelectors(ctx, plan.Static.Selectors)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error extracting selectors",
				err.Error(),
			)
			return
		}
		staticPolicy.Selectors = selectors
		policy.Policy = &attestationpolicypb.AttestationPolicy_Static{
			Static: staticPolicy,
		}
	}

	updateResp, err := r.client.AttestationPolicyV1Alpha1().UpdateAttestationPolicy(ctx, policy)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating attestation policy",
			fmt.Sprintf("Could not update attestation policy: %s", err.Error()),
		)
		return
	}

	plan.ID = tftypes.StringValue(updateResp.GetId())
	plan.Name = tftypes.StringValue(updateResp.GetName())
	plan.OrgID = tftypes.StringValue(updateResp.GetOrgId())

	if k8s := updateResp.GetKubernetes(); k8s != nil && plan.Kubernetes != nil {
		plan.Kubernetes = &APKubernetesModel{
			NamespaceSelector: plan.Kubernetes.NamespaceSelector,
			PodSelector:       plan.Kubernetes.PodSelector,
		}
	}

	if static := updateResp.GetStatic(); static != nil && plan.Static != nil {
		plan.Static = &APStaticModel{
			SpiffeID:  tftypes.StringValue(static.GetSpiffeId()),
			Selectors: plan.Static.Selectors,
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *AttestationPolicyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state AttestationPolicyModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.AttestationPolicyV1Alpha1().DestroyAttestationPolicy(ctx, state.ID.ValueString())
	if err != nil {
		if status.Code(err) != codes.NotFound {
			resp.Diagnostics.AddError(
				"Error deleting attestation policy",
				fmt.Sprintf("Could not delete attestation policy: %s", err),
			)
			return
		}
	}
}

func (r *AttestationPolicyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *AttestationPolicyResource) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	var data AttestationPolicyModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	hasKubernetes := data.Kubernetes != nil
	hasStatic := data.Static != nil

	if !hasKubernetes && !hasStatic {
		resp.Diagnostics.AddError(
			"Invalid configuration",
			"Either kubernetes or static block must be configured, but neither was provided.",
		)
		return
	}

	if hasKubernetes && hasStatic {
		resp.Diagnostics.AddError(
			"Invalid configuration",
			"Only one of kubernetes or static block can be configured, but both were provided.",
		)
		return
	}
}

func extractSelectors(ctx context.Context, selectors []APStaticSelectorModel) ([]*spiretypes.Selector, error) {
	protoSelectors := []*spiretypes.Selector{}

	for _, s := range selectors {
		protoSelectors = append(protoSelectors, &spiretypes.Selector{
			Type:  s.Type.ValueString(),
			Value: s.Value.ValueString(),
		})
	}

	return protoSelectors, nil
}

func convertLabelSelector(selector *APLabelSelectorModel) *attestationpolicypb.APLabelSelector {
	if selector == nil {
		return nil
	}

	result := &attestationpolicypb.APLabelSelector{
		MatchLabels: make(map[string]string),
	}

	if !selector.MatchLabels.IsNull() {
		elements := selector.MatchLabels.Elements()
		for k, v := range elements {
			if str, ok := v.(tftypes.String); ok {
				result.MatchLabels[k] = str.ValueString()
			}
		}
	}

	for _, expr := range selector.MatchExpressions {
		matchExpr := &attestationpolicypb.APMatchExpression{
			Key:      expr.Key.ValueString(),
			Operator: expr.Operator.ValueString(),
		}
		for _, v := range expr.Values {
			matchExpr.Values = append(matchExpr.Values, v.ValueString())
		}
		result.MatchExpressions = append(result.MatchExpressions, matchExpr)
	}

	return result
}
