package apbinding

import (
	apbindingpb "github.com/cofide/cofide-api-sdk/gen/go/proto/ap_binding/v1alpha1"
	tftypes "github.com/hashicorp/terraform-plugin-framework/types"
)

// modelToProto converts an APBindingModel to an equivalent APBinding protobuf.
func modelToProto(model APBindingModel) *apbindingpb.APBinding {
	federations := make([]*apbindingpb.APBindingFederation, 0, len(model.Federations))
	for _, federation := range model.Federations {
		federations = append(federations, &apbindingpb.APBindingFederation{
			TrustZoneId: federation.TrustZoneID.ValueStringPointer(),
		})
	}

	return &apbindingpb.APBinding{
		OrgId:       model.OrgID.ValueStringPointer(),
		TrustZoneId: model.TrustZoneID.ValueStringPointer(),
		PolicyId:    model.PolicyID.ValueStringPointer(),
		Federations: federations,
	}
}

// protoToModel converts an APBinding protobuf to an equivalent APBindingModel.
func protoToModel(proto *apbindingpb.APBinding) APBindingModel {
	federations := make([]APBindingFederationModel, 0, len(proto.GetFederations()))
	for _, federation := range proto.GetFederations() {
		federations = append(federations, APBindingFederationModel{
			TrustZoneID: tftypes.StringValue(federation.GetTrustZoneId()),
		})
	}

	return APBindingModel{
		ID:          tftypes.StringValue(proto.GetId()),
		OrgID:       tftypes.StringValue(proto.GetOrgId()),
		TrustZoneID: tftypes.StringValue(proto.GetTrustZoneId()),
		PolicyID:    tftypes.StringValue(proto.GetPolicyId()),
		Federations: federations,
	}
}
