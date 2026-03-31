package attestationpolicy

import (
	"context"

	attestationpolicypb "github.com/cofide/cofide-api-sdk/gen/go/proto/attestation_policy/v1alpha1"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	tftypes "github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	spiretypes "github.com/spiffe/spire-api-sdk/proto/spire/api/types"
)

// modelToProto converts an AttestationPolicyModel to an equivalent AttestationPolicy protobuf.
func modelToProto(ctx context.Context, model AttestationPolicyModel) *attestationpolicypb.AttestationPolicy {
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

	if model.TPMNode != nil {
		var selectorValues []tftypes.String
		model.TPMNode.SelectorValues.ElementsAs(ctx, &selectorValues, false)
		tpmNodePolicy := &attestationpolicypb.APTPMNode{
			Attestation: &attestationpolicypb.TPMAttestation{
				EkHash: model.TPMNode.Attestation.EKHash.ValueStringPointer(),
			},
			SelectorValues: convertStringSlice(selectorValues),
		}
		proto.Policy = &attestationpolicypb.AttestationPolicy_TpmNode{
			TpmNode: tpmNodePolicy,
		}
	}

	return proto
}

// protoToModel converts an AttestationPolicy protobuf to an equivalent AttestationPolicyModel.
func protoToModel(proto *attestationpolicypb.AttestationPolicy) AttestationPolicyModel {
	model := AttestationPolicyModel{
		ID:    optionalStringValue(proto.Id),
		Name:  tftypes.StringValue(proto.GetName()),
		OrgID: optionalStringValue(proto.OrgId),
	}

	if k8s := proto.GetKubernetes(); k8s != nil {
		model.Kubernetes = &APKubernetesModel{
			NamespaceSelector:    convertProtoLabelSelector(k8s.NamespaceSelector),
			PodSelector:          convertProtoLabelSelector(k8s.PodSelector),
			DnsNameTemplates:     convertProtoStringSlice(k8s.GetDnsNameTemplates()),
			SpiffeIDPathTemplate: optionalStringValue(k8s.SpiffeIdPathTemplate),
		}
	}

	if static := proto.GetStatic(); static != nil {
		model.Static = &APStaticModel{
			SpiffeIDPath: optionalStringValue(static.SpiffeIdPath),
			ParentIdPath: optionalStringValue(static.ParentIdPath),
			Selectors:    convertProtoSelectors(static.GetSelectors()),
			DNSNames:     convertProtoStringSlice(static.GetDnsNames()),
		}
	}

	if tpmNode := proto.GetTpmNode(); tpmNode != nil {
		model.TPMNode = &APTPMNodeModel{
			SelectorValues: convertProtoSelectorValues(tpmNode.GetSelectorValues()),
		}
		if attestation := tpmNode.GetAttestation(); attestation != nil {
			model.TPMNode.Attestation.EKHash = optionalStringValue(attestation.EkHash)
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

// convertProtoSelectorValues converts a slice of strings from protobuf to a Terraform types.List.
func convertProtoSelectorValues(input []string) tftypes.List {
	elems := make([]attr.Value, 0, len(input))
	for _, s := range input {
		elems = append(elems, tftypes.StringValue(s))
	}
	return tftypes.ListValueMust(tftypes.StringType, elems)
}

// selectorAttrTypes defines the attribute types for a single selector object.
var selectorAttrTypes = map[string]attr.Type{
	"type":  tftypes.StringType,
	"value": tftypes.StringType,
}

// selectorElemType is the Terraform object type for a single selector.
var selectorElemType = tftypes.ObjectType{AttrTypes: selectorAttrTypes}

func convertSelectors(selectors tftypes.List) []*spiretypes.Selector {
	var protoSelectors []*spiretypes.Selector
	for _, elem := range selectors.Elements() {
		obj, ok := elem.(tftypes.Object)
		if !ok {
			continue
		}
		attrs := obj.Attributes()
		protoSelectors = append(protoSelectors, &spiretypes.Selector{
			Type:  attrs["type"].(tftypes.String).ValueString(),
			Value: attrs["value"].(tftypes.String).ValueString(),
		})
	}
	return protoSelectors
}

func convertProtoSelectors(selectors []*spiretypes.Selector) tftypes.List {
	elems := make([]attr.Value, 0, len(selectors))
	for _, s := range selectors {
		obj, diags := tftypes.ObjectValue(selectorAttrTypes, map[string]attr.Value{
			"type":  tftypes.StringValue(s.Type),
			"value": tftypes.StringValue(s.Value),
		})
		if diags.HasError() {
			continue
		}
		elems = append(elems, obj)
	}
	return tftypes.ListValueMust(selectorElemType, elems)
}

func optionalStringValue(s *string) basetypes.StringValue {
	if s == nil {
		return tftypes.StringNull()
	}
	return tftypes.StringValue(*s)
}
