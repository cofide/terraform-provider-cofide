package workloadsuppressionrule

import (
	"context"
	"time"

	workloadsuppressionrulepb "github.com/cofide/cofide-api-sdk/gen/go/proto/workload_suppression_rule/v1alpha1"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	tftypes "github.com/hashicorp/terraform-plugin-framework/types"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// modelToProto converts a WorkloadSuppressionRuleModel to an equivalent WorkloadSuppressionRule protobuf.
func modelToProto(ctx context.Context, model WorkloadSuppressionRuleModel) (*workloadsuppressionrulepb.WorkloadSuppressionRule, diag.Diagnostics) {
	var diags diag.Diagnostics

	proto := &workloadsuppressionrulepb.WorkloadSuppressionRule{
		Id:          model.ID.ValueString(),
		OrgId:       model.OrgID.ValueString(),
		Name:        model.Name.ValueString(),
		Description: model.Description.ValueString(),
		Enabled:     model.Enabled.ValueBool(),
	}

	if model.KubernetesPod != nil {
		var trustZoneIDs, clusterIDs, namespaces []string
		diags.Append(model.KubernetesPod.TrustZoneIDs.ElementsAs(ctx, &trustZoneIDs, false)...)
		diags.Append(model.KubernetesPod.ClusterIDs.ElementsAs(ctx, &clusterIDs, false)...)
		diags.Append(model.KubernetesPod.Namespaces.ElementsAs(ctx, &namespaces, false)...)

		labels := make(map[string]string)
		if !model.KubernetesPod.Labels.IsNull() {
			for k, v := range model.KubernetesPod.Labels.Elements() {
				if str, ok := v.(tftypes.String); ok {
					labels[k] = str.ValueString()
				}
			}
		}

		proto.Matcher = &workloadsuppressionrulepb.WorkloadSuppressionRule_KubernetesPod{
			KubernetesPod: &workloadsuppressionrulepb.KubernetesPodMatcher{
				TrustZoneIds: trustZoneIDs,
				ClusterIds:   clusterIDs,
				Namespaces:   namespaces,
				Labels:       labels,
			},
		}
	}

	return proto, diags
}

// protoToModel converts a WorkloadSuppressionRule protobuf to an equivalent WorkloadSuppressionRuleModel.
func protoToModel(proto *workloadsuppressionrulepb.WorkloadSuppressionRule) WorkloadSuppressionRuleModel {
	model := WorkloadSuppressionRuleModel{
		ID:            tftypes.StringValue(proto.GetId()),
		OrgID:         tftypes.StringValue(proto.GetOrgId()),
		Name:          tftypes.StringValue(proto.GetName()),
		Description:   stringOrNull(proto.GetDescription()),
		Enabled:       tftypes.BoolValue(proto.GetEnabled()),
		CreatedAt:     timestampToString(proto.GetCreatedAt()),
		LastUpdatedAt: timestampToString(proto.GetLastUpdatedAt()),
	}

	if pod := proto.GetKubernetesPod(); pod != nil {
		model.KubernetesPod = &KubernetesPodMatcherModel{
			TrustZoneIDs: convertProtoStringList(pod.GetTrustZoneIds()),
			ClusterIDs:   convertProtoStringList(pod.GetClusterIds()),
			Namespaces:   convertProtoStringList(pod.GetNamespaces()),
			Labels:       convertProtoLabels(pod.GetLabels()),
		}
	}

	return model
}

// convertProtoStringList converts a slice of strings from protobuf to a Terraform types.List.
// Returns a null list when input is empty.
func convertProtoStringList(input []string) tftypes.List {
	if len(input) == 0 {
		return tftypes.ListNull(tftypes.StringType)
	}
	elems := make([]attr.Value, 0, len(input))
	for _, s := range input {
		elems = append(elems, tftypes.StringValue(s))
	}
	return tftypes.ListValueMust(tftypes.StringType, elems)
}

// convertProtoLabels converts a map of strings from protobuf to a Terraform types.Map.
// Returns a null map when input is empty.
func convertProtoLabels(input map[string]string) tftypes.Map {
	if len(input) == 0 {
		return tftypes.MapNull(tftypes.StringType)
	}
	elements := make(map[string]attr.Value, len(input))
	for k, v := range input {
		elements[k] = tftypes.StringValue(v)
	}
	return tftypes.MapValueMust(tftypes.StringType, elements)
}

// stringOrNull returns a null string when s is empty, since proto3 cannot
// distinguish an unset optional scalar field from its zero value.
func stringOrNull(s string) tftypes.String {
	if s == "" {
		return tftypes.StringNull()
	}
	return tftypes.StringValue(s)
}

// timestampToString formats a protobuf timestamp as RFC3339, returning a null
// string when the timestamp is unset.
func timestampToString(ts *timestamppb.Timestamp) tftypes.String {
	if ts == nil {
		return tftypes.StringNull()
	}
	return tftypes.StringValue(ts.AsTime().Format(time.RFC3339))
}
