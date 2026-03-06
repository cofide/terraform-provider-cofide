package federation

import (
	federationpb "github.com/cofide/cofide-api-sdk/gen/go/proto/federation/v1alpha1"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func modelToProto(plan FederationModel) *federationpb.Federation {
	return &federationpb.Federation{
		OrgId:             plan.OrgID.ValueStringPointer(),
		TrustZoneId:       plan.TrustZoneID.ValueStringPointer(),
		RemoteTrustZoneId: plan.RemoteTrustZoneID.ValueStringPointer(),
	}
}

func protoToModel(federation *federationpb.Federation) FederationModel {
	return FederationModel{
		ID:                types.StringValue(federation.GetId()),
		OrgID:             types.StringValue(federation.GetOrgId()),
		TrustZoneID:       types.StringValue(federation.GetTrustZoneId()),
		RemoteTrustZoneID: types.StringValue(federation.GetRemoteTrustZoneId()),
	}
}
