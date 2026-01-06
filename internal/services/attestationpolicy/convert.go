package attestationpolicy

import (
	attestationpolicypb "github.com/cofide/cofide-api-sdk/gen/go/proto/attestation_policy/v1alpha1"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	tftypes "github.com/hashicorp/terraform-plugin-framework/types"
	spiretypes "github.com/spiffe/spire-api-sdk/proto/spire/api/types"
)

// modelToProto converts an AttestationPolicyModel to an equivalent AttestationPolicy protobuf.
func modelToProto(model AttestationPolicyModel) *attestationpolicypb.AttestationPolicy {
	proto := &attestationpolicypb.AttestationPolicy{
		Id:    model.ID.ValueStringPointer(),
		Name:  model.Name.ValueString(),
		OrgId: model.OrgID.ValueStringPointer(),
	}

	if model.Kubernetes != nil {
		k8sPolicy := &attestationpolicypb.APKubernetes{
			NamespaceSelector:    convertLabelSelector(model.Kubernetes.NamespaceSelector),
			PodSelector:          convertLabelSelector(model.Kubernetes.PodSelector),
			DnsNameTemplates:     convertStringSlice(model.Kubernetes.DnsNameTemplates),
			SpiffeIdPathTemplate: model.Kubernetes.SpiffeIDPathTemplate.ValueStringPointer(),
		}
		proto.Policy = &attestationpolicypb.AttestationPolicy_Kubernetes{
			Kubernetes: k8sPolicy,
		}
	}

	if model.Static != nil {
		staticPolicy := &attestationpolicypb.APStatic{
			SpiffeIdPath: model.Static.SpiffeIDPath.ValueStringPointer(),
			ParentIdPath: model.Static.ParentIdPath.ValueStringPointer(),
			Selectors:    convertSelectors(model.Static.Selectors),
			DnsNames:     convertStringSlice(model.Static.DNSNames),
		}
		proto.Policy = &attestationpolicypb.AttestationPolicy_Static{
			Static: staticPolicy,
		}
	}
	return proto
}

// protoToModel converts an AttestationPolicy protobuf to an equivalent AttestationPolicyModel.
func protoToModel(proto *attestationpolicypb.AttestationPolicy) AttestationPolicyModel {
	model := AttestationPolicyModel{
		ID:    tftypes.StringValue(proto.GetId()),
		Name:  tftypes.StringValue(proto.GetName()),
		OrgID: tftypes.StringValue(proto.GetOrgId()),
	}

	if k8s := proto.GetKubernetes(); k8s != nil {
		model.Kubernetes = &APKubernetesModel{
			NamespaceSelector:    convertProtoLabelSelector(k8s.NamespaceSelector),
			PodSelector:          convertProtoLabelSelector(k8s.PodSelector),
			DnsNameTemplates:     convertProtoStringSlice(k8s.GetDnsNameTemplates()),
			SpiffeIDPathTemplate: tftypes.StringValue(k8s.GetSpiffeIdPathTemplate()),
		}
	}

	if static := proto.GetStatic(); static != nil {
		model.Static = &APStaticModel{
			SpiffeIDPath: tftypes.StringValue(static.GetSpiffeIdPath()),
			ParentIdPath: tftypes.StringValue(static.GetParentIdPath()),
			Selectors:    convertProtoSelectors(static.GetSelectors()),
			DNSNames:     convertProtoStringSlice(static.GetDnsNames()),
		}
	}
	return model
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

func convertProtoLabelSelector(selector *attestationpolicypb.APLabelSelector) *APLabelSelectorModel {
	if selector == nil {
		return nil
	}

	result := &APLabelSelectorModel{
		MatchLabels: tftypes.MapNull(tftypes.StringType),
	}

	// Convert match labels
	if len(selector.MatchLabels) > 0 {
		elements := make(map[string]attr.Value)
		for k, v := range selector.MatchLabels {
			elements[k] = tftypes.StringValue(v)
		}
		result.MatchLabels = tftypes.MapValueMust(tftypes.StringType, elements)
	}

	// Convert match expressions
	for _, expr := range selector.MatchExpressions {
		matchExpr := APMatchExpressionModel{
			Key:      tftypes.StringValue(expr.GetKey()),
			Operator: tftypes.StringValue(expr.GetOperator()),
		}
		for _, v := range expr.GetValues() {
			matchExpr.Values = append(matchExpr.Values, tftypes.StringValue(v))
		}
		result.MatchExpressions = append(result.MatchExpressions, matchExpr)
	}

	return result
}

// convertStringSlice converts a slice of strings from Terraform to protobuf types.
func convertStringSlice(input []tftypes.String) []string {
	var result []string
	for _, s := range input {
		result = append(result, s.ValueString())
	}
	return result
}

// convertProtoStringSlice converts a slice of strings from protobuf to Terraform types.
func convertProtoStringSlice(input []string) []tftypes.String {
	var result []tftypes.String
	for _, t := range input {
		result = append(result, tftypes.StringValue(t))
	}
	return result
}

func convertSelectors(selectors []APStaticSelectorModel) []*spiretypes.Selector {
	protoSelectors := []*spiretypes.Selector{}
	for _, s := range selectors {
		protoSelectors = append(protoSelectors, &spiretypes.Selector{
			Type:  s.Type.ValueString(),
			Value: s.Value.ValueString(),
		})
	}
	return protoSelectors
}

func convertProtoSelectors(selectors []*spiretypes.Selector) []APStaticSelectorModel {
	modelSelectors := []APStaticSelectorModel{}
	for _, s := range selectors {
		modelSelectors = append(modelSelectors, APStaticSelectorModel{
			Type:  tftypes.StringValue(s.Type),
			Value: tftypes.StringValue(s.Value),
		})
	}
	return modelSelectors
}
