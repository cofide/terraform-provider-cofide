package trustzone

import (
	"fmt"

	trustzonepb "github.com/cofide/cofide-api-sdk/gen/go/proto/trust_zone/v1alpha1"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func modelToProto(plan TrustZoneModel) (*trustzonepb.TrustZone, error) {
	trustZone := &trustzonepb.TrustZone{
		Name:             plan.Name.ValueString(),
		OrgId:            plan.OrgID.ValueStringPointer(),
		TrustDomain:      plan.TrustDomain.ValueString(),
		IsManagementZone: plan.IsManagementZone.ValueBool(),
	}

	if !plan.BundleEndpointURL.IsNull() && plan.BundleEndpointURL.ValueString() != "" {
		trustZone.BundleEndpointUrl = plan.BundleEndpointURL.ValueStringPointer()
	}

	if !plan.BundleEndpointProfile.IsNull() && plan.BundleEndpointProfile.ValueString() != "" {
		if profile, ok := getBundleEndpointProfile(plan.BundleEndpointProfile.ValueString()); ok {
			trustZone.BundleEndpointProfile = profile
		} else {
			return nil, fmt.Errorf("unknown BundleEndpointProfile: %s", plan.BundleEndpointProfile.ValueString())
		}
	}

	if !plan.JWTIssuer.IsNull() && plan.JWTIssuer.ValueString() != "" {
		trustZone.JwtIssuer = plan.JWTIssuer.ValueStringPointer()
	}

	return trustZone, nil
}

func protoToModel(trustZone *trustzonepb.TrustZone) TrustZoneModel {
	return TrustZoneModel{
		ID:                    types.StringValue(trustZone.GetId()),
		Name:                  types.StringValue(trustZone.GetName()),
		TrustDomain:           types.StringValue(trustZone.GetTrustDomain()),
		OrgID:                 types.StringValue(trustZone.GetOrgId()),
		IsManagementZone:      types.BoolValue(trustZone.GetIsManagementZone()),
		BundleEndpointURL:     types.StringValue(trustZone.GetBundleEndpointUrl()),
		BundleEndpointProfile: types.StringValue(trustZone.GetBundleEndpointProfile().String()),
		JWTIssuer:             types.StringValue(trustZone.GetJwtIssuer()),
	}
}

// getBundleEndpointProfile converts a string to a BundleEndpointProfile enum pointer
func getBundleEndpointProfile(value string) (*trustzonepb.BundleEndpointProfile, bool) {
	if profileVal, ok := trustzonepb.BundleEndpointProfile_value[value]; ok {
		profile := trustzonepb.BundleEndpointProfile(profileVal)
		return &profile, true
	}
	return nil, false
}
