package attestationpolicy

import (
	"context"
	"fmt"

	attestationpolicypb "github.com/cofide/cofide-api-sdk/gen/go/proto/attestation_policy/v1alpha1"
	attestationpolicysvcpb "github.com/cofide/cofide-api-sdk/gen/go/proto/connect/attestation_policy_service/v1alpha1"
	sdkclient "github.com/cofide/cofide-api-sdk/pkg/connect/client"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type AttestationPolicyDataSource struct {
	client sdkclient.ClientSet
}

var _ datasource.DataSourceWithConfigure = (*AttestationPolicyDataSource)(nil)

func NewDataSource() datasource.DataSource {
	return &AttestationPolicyDataSource{}
}

func (d *AttestationPolicyDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_connect_attestation_policy"
}

func (d *AttestationPolicyDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(sdkclient.ClientSet)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected data source configure type",
			fmt.Sprintf("Expected sdkclient.ClientSet, got: %T", req.ProviderData),
		)
		return
	}

	d.client = client
}

func (d *AttestationPolicyDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config AttestationPolicyModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	filter := &attestationpolicysvcpb.ListAttestationPoliciesRequest_Filter{
		Name:  config.Name.ValueStringPointer(),
		OrgId: config.OrgID.ValueStringPointer(),
	}
	policies, err := d.client.AttestationPolicyV1Alpha1().ListAttestationPolicies(ctx, filter)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading attestation policy",
			fmt.Sprintf("Could not list attestation policies: %s", err),
		)
		return
	}

	if len(policies) == 0 {
		resp.Diagnostics.AddError(
			"Error reading attestation policy",
			"No matching attestation policy found",
		)
		return
	}

	if len(policies) > 1 {
		resp.Diagnostics.AddError(
			"Error reading attestation policy",
			"Multiple attestation policies found",
		)
		return
	}

	policy := policies[0]

	if policy == nil {
		resp.Diagnostics.AddError(
			"Error reading attestation policy",
			"No matching attestation policy found",
		)
		return
	}

	state := AttestationPolicyModel{
		ID:    types.StringValue(policy.GetId()),
		Name:  types.StringValue(policy.GetName()),
		OrgID: types.StringValue(policy.GetOrgId()),
	}

	if k8s := policy.GetKubernetes(); k8s != nil {
		state.Kubernetes = &APKubernetesModel{}
		if ns := k8s.GetNamespaceSelector(); ns != nil {
			state.Kubernetes.NamespaceSelector = convertProtoLabelSelector(ns)
		}
		if ps := k8s.GetPodSelector(); ps != nil {
			state.Kubernetes.PodSelector = convertProtoLabelSelector(ps)
		}
	}

	if static := policy.GetStatic(); static != nil {
		state.Static = &APStaticModel{
			SpiffeID: types.StringValue(static.GetSpiffeId()),
		}
		for _, selector := range static.GetSelectors() {
			state.Static.Selectors = append(state.Static.Selectors, APStaticSelectorModel{
				Type:  types.StringValue(selector.GetType()),
				Value: types.StringValue(selector.GetValue()),
			})
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func convertProtoLabelSelector(selector *attestationpolicypb.APLabelSelector) *APLabelSelectorModel {
	if selector == nil {
		return nil
	}

	result := &APLabelSelectorModel{
		MatchLabels: types.MapNull(types.StringType),
	}

	// Convert match labels
	if len(selector.MatchLabels) > 0 {
		elements := make(map[string]attr.Value)
		for k, v := range selector.MatchLabels {
			elements[k] = types.StringValue(v)
		}
		result.MatchLabels = types.MapValueMust(types.StringType, elements)
	}

	// Convert match expressions
	for _, expr := range selector.MatchExpressions {
		matchExpr := APMatchExpressionModel{
			Key:      types.StringValue(expr.GetKey()),
			Operator: types.StringValue(expr.GetOperator()),
		}
		for _, v := range expr.GetValues() {
			matchExpr.Values = append(matchExpr.Values, types.StringValue(v))
		}
		result.MatchExpressions = append(result.MatchExpressions, matchExpr)
	}

	return result
}
