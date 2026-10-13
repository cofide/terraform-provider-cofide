package v1alpha1

import (
	"context"

	attestationpolicypb "github.com/cofide/cofide-api-sdk/gen/go/proto/attestation_policy/v1alpha1"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	tftypes "github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	spiretypes "github.com/spiffe/spire-api-sdk/proto/spire/api/types"
)

// The types of the model's objects and list elements, taken from the resource schema so they can't drift from it.
// The data source schema has the same shape (see internal/schema_shape_test.go).
var (
	policyAttrTypes         = ResourceSchema(context.Background()).Type().(tftypes.ObjectType).AttrTypes
	federationElemType      = policyAttrTypes["federations"].(tftypes.ListType).ElemType.(tftypes.ObjectType)
	kubernetesType          = policyAttrTypes["kubernetes"].(tftypes.ObjectType)
	labelSelectorType       = kubernetesType.AttrTypes["namespace_selector"].(tftypes.ObjectType)
	matchExpressionElemType = labelSelectorType.AttrTypes["match_expressions"].(tftypes.ListType).ElemType.(tftypes.ObjectType)
	staticType              = policyAttrTypes["static"].(tftypes.ObjectType)
	tpmNodeType             = policyAttrTypes["tpm_node"].(tftypes.ObjectType)
	tpmAttestationType      = tpmNodeType.AttrTypes["attestation"].(tftypes.ObjectType)
)

// isSet returns whether v has a known, non-null value.
func isSet(v attr.Value) bool {
	return !v.IsNull() && !v.IsUnknown()
}

// modelToProto converts an AttestationPolicyModel to an equivalent AttestationPolicy protobuf. It's called with a
// plan at apply, when every configured value is known.
func modelToProto(ctx context.Context, model AttestationPolicyModel) (*attestationpolicypb.AttestationPolicy, diag.Diagnostics) {
	var diags diag.Diagnostics

	proto := &attestationpolicypb.AttestationPolicy{
		Id:          model.ID.ValueStringPointer(),
		Name:        model.Name.ValueString(),
		TrustZoneId: model.TrustZoneID.ValueString(),
	}

	if isSet(model.Federations) {
		var federations []APFederationModel
		diags.Append(model.Federations.ElementsAs(ctx, &federations, false)...)
		proto.Federations = make([]*attestationpolicypb.Federation, 0, len(federations))
		for _, federation := range federations {
			proto.Federations = append(proto.Federations, &attestationpolicypb.Federation{
				TrustZoneId: federation.TrustZoneID.ValueStringPointer(),
			})
		}
	}

	if isSet(model.Kubernetes) {
		var k8s APKubernetesModel
		diags.Append(model.Kubernetes.As(ctx, &k8s, basetypes.ObjectAsOptions{})...)
		var dnsNameTemplates []string
		diags.Append(k8s.DnsNameTemplates.ElementsAs(ctx, &dnsNameTemplates, false)...)
		k8sPolicy := &attestationpolicypb.APKubernetes{
			NamespaceSelector:    convertLabelSelector(ctx, k8s.NamespaceSelector, &diags),
			PodSelector:          convertLabelSelector(ctx, k8s.PodSelector, &diags),
			DnsNameTemplates:     dnsNameTemplates,
			SpiffeIdPathTemplate: k8s.SpiffeIDPathTemplate.ValueStringPointer(),
		}
		proto.Policy = &attestationpolicypb.AttestationPolicy_Kubernetes{
			Kubernetes: k8sPolicy,
		}
	}

	if isSet(model.Static) {
		var static APStaticModel
		diags.Append(model.Static.As(ctx, &static, basetypes.ObjectAsOptions{})...)
		var dnsNames []string
		diags.Append(static.DNSNames.ElementsAs(ctx, &dnsNames, false)...)
		staticPolicy := &attestationpolicypb.APStatic{
			SpiffeIdPath: static.SpiffeIDPath.ValueStringPointer(),
			ParentIdPath: static.ParentIdPath.ValueStringPointer(),
			Selectors:    convertSelectors(static.Selectors),
			DnsNames:     dnsNames,
			StoreSvid:    isSet(static.StoreSvid) && static.StoreSvid.ValueBool(),
		}
		proto.Policy = &attestationpolicypb.AttestationPolicy_Static{
			Static: staticPolicy,
		}
	}

	if isSet(model.TPMNode) {
		var tpmNode APTPMNodeModel
		diags.Append(model.TPMNode.As(ctx, &tpmNode, basetypes.ObjectAsOptions{})...)
		var attestation TPMAttestationModel
		if isSet(tpmNode.Attestation) {
			diags.Append(tpmNode.Attestation.As(ctx, &attestation, basetypes.ObjectAsOptions{})...)
		}
		var selectorValues []string
		diags.Append(tpmNode.SelectorValues.ElementsAs(ctx, &selectorValues, false)...)
		tpmNodePolicy := &attestationpolicypb.APTPMNode{
			Attestation: &attestationpolicypb.TPMAttestation{
				EkHash: attestation.EKHash.ValueStringPointer(),
			},
			SelectorValues: selectorValues,
		}
		proto.Policy = &attestationpolicypb.AttestationPolicy_TpmNode{
			TpmNode: tpmNodePolicy,
		}
	}

	return proto, diags
}

// protoToModel converts an AttestationPolicy protobuf to an equivalent AttestationPolicyModel.
func protoToModel(ctx context.Context, proto *attestationpolicypb.AttestationPolicy) (AttestationPolicyModel, diag.Diagnostics) {
	var diags diag.Diagnostics

	model := AttestationPolicyModel{
		ID:          optionalStringValue(proto.Id),
		Name:        tftypes.StringValue(proto.GetName()),
		OrgID:       optionalStringValue(proto.OrgId),
		TrustZoneID: nonEmptyStringValue(proto.GetTrustZoneId()),
		Federations: tftypes.ListNull(federationElemType),
		Kubernetes:  tftypes.ObjectNull(kubernetesType.AttrTypes),
		Static:      tftypes.ObjectNull(staticType.AttrTypes),
		TPMNode:     tftypes.ObjectNull(tpmNodeType.AttrTypes),
	}

	if federations := proto.GetFederations(); federations != nil {
		respFederations := make([]APFederationModel, 0, len(federations))
		for _, federation := range federations {
			respFederations = append(respFederations, APFederationModel{
				TrustZoneID: optionalStringValue(federation.TrustZoneId),
			})
		}
		model.Federations = listValueFrom(ctx, federationElemType, respFederations, &diags)
	}

	if k8s := proto.GetKubernetes(); k8s != nil {
		model.Kubernetes = objectValueFrom(ctx, kubernetesType, APKubernetesModel{
			NamespaceSelector:    convertProtoLabelSelector(ctx, k8s.NamespaceSelector, &diags),
			PodSelector:          convertProtoLabelSelector(ctx, k8s.PodSelector, &diags),
			DnsNameTemplates:     convertProtoSelectorValues(k8s.GetDnsNameTemplates()),
			SpiffeIDPathTemplate: optionalStringValue(k8s.SpiffeIdPathTemplate),
		}, &diags)
	}

	if static := proto.GetStatic(); static != nil {
		model.Static = objectValueFrom(ctx, staticType, APStaticModel{
			SpiffeIDPath: optionalStringValue(static.SpiffeIdPath),
			ParentIdPath: optionalStringValue(static.ParentIdPath),
			Selectors:    convertProtoSelectors(static.GetSelectors()),
			DNSNames:     convertProtoSelectorValues(static.GetDnsNames()),
			StoreSvid:    tftypes.BoolValue(static.GetStoreSvid()),
		}, &diags)
	}

	if tpmNode := proto.GetTpmNode(); tpmNode != nil {
		attestation := TPMAttestationModel{EKHash: tftypes.StringNull()}
		if protoAttestation := tpmNode.GetAttestation(); protoAttestation != nil {
			attestation.EKHash = optionalStringValue(protoAttestation.EkHash)
		}
		model.TPMNode = objectValueFrom(ctx, tpmNodeType, APTPMNodeModel{
			Attestation:    objectValueFrom(ctx, tpmAttestationType, attestation, &diags),
			SelectorValues: convertProtoSelectorValues(tpmNode.GetSelectorValues()),
		}, &diags)
	}

	return model, diags
}

// objectValueFrom converts value, a model struct, to an object of type t, appending any errors to diags.
func objectValueFrom(ctx context.Context, t tftypes.ObjectType, value any, diags *diag.Diagnostics) tftypes.Object {
	obj, d := tftypes.ObjectValueFrom(ctx, t.AttrTypes, value)
	diags.Append(d...)
	return obj
}

// listValueFrom converts elements, a slice of model structs, to a list of elemType, appending any errors to diags.
func listValueFrom(ctx context.Context, elemType tftypes.ObjectType, elements any, diags *diag.Diagnostics) tftypes.List {
	list, d := tftypes.ListValueFrom(ctx, elemType, elements)
	diags.Append(d...)
	return list
}

func convertLabelSelector(ctx context.Context, obj tftypes.Object, diags *diag.Diagnostics) *attestationpolicypb.APLabelSelector {
	if !isSet(obj) {
		return nil
	}

	var selector APLabelSelectorModel
	diags.Append(obj.As(ctx, &selector, basetypes.ObjectAsOptions{})...)

	result := &attestationpolicypb.APLabelSelector{
		MatchLabels: make(map[string]string),
	}

	if isSet(selector.MatchLabels) {
		elements := selector.MatchLabels.Elements()
		for k, v := range elements {
			if str, ok := v.(tftypes.String); ok {
				result.MatchLabels[k] = str.ValueString()
			}
		}
	}

	if isSet(selector.MatchExpressions) {
		var exprs []APMatchExpressionModel
		diags.Append(selector.MatchExpressions.ElementsAs(ctx, &exprs, false)...)
		for _, expr := range exprs {
			matchExpr := &attestationpolicypb.APMatchExpression{
				Key:      expr.Key.ValueString(),
				Operator: expr.Operator.ValueString(),
			}
			var values []string
			diags.Append(expr.Values.ElementsAs(ctx, &values, false)...)
			matchExpr.Values = values
			result.MatchExpressions = append(result.MatchExpressions, matchExpr)
		}
	}

	return result
}

func convertProtoLabelSelector(ctx context.Context, selector *attestationpolicypb.APLabelSelector, diags *diag.Diagnostics) tftypes.Object {
	if selector == nil {
		return tftypes.ObjectNull(labelSelectorType.AttrTypes)
	}

	result := APLabelSelectorModel{
		MatchLabels:      tftypes.MapNull(tftypes.StringType),
		MatchExpressions: tftypes.ListNull(matchExpressionElemType),
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
	if len(selector.MatchExpressions) > 0 {
		exprs := make([]APMatchExpressionModel, 0, len(selector.MatchExpressions))
		for _, expr := range selector.MatchExpressions {
			exprs = append(exprs, APMatchExpressionModel{
				Key:      tftypes.StringValue(expr.GetKey()),
				Operator: tftypes.StringValue(expr.GetOperator()),
				Values:   convertProtoSelectorValues(expr.GetValues()),
			})
		}
		result.MatchExpressions = listValueFrom(ctx, matchExpressionElemType, exprs, diags)
	}

	return objectValueFrom(ctx, labelSelectorType, result, diags)
}

// convertProtoSelectorValues converts a slice of strings from protobuf to a Terraform types.List.
// Returns a null list when input is empty.
func convertProtoSelectorValues(input []string) tftypes.List {
	if len(input) == 0 {
		return tftypes.ListNull(tftypes.StringType)
	}
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

// nonEmptyStringValue converts a plain (non-pointer) proto string field that has no
// presence tracking into a Terraform string, treating the empty zero value as null.
// This keeps state consistent with an unconfigured Optional attribute.
func nonEmptyStringValue(s string) basetypes.StringValue {
	if s == "" {
		return tftypes.StringNull()
	}
	return tftypes.StringValue(s)
}
