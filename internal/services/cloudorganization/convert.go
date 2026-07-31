package cloudorganization

import (
	cloudorganizationpb "github.com/cofide/cofide-api-sdk/gen/go/proto/cloud_organization/v1alpha1"
	"github.com/cofide/terraform-provider-cofide/internal/services/cloudprovider"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	tftypes "github.com/hashicorp/terraform-plugin-framework/types"
)

// modelToProto converts a CloudOrganizationModel to a CloudOrganization proto message.
func modelToProto(model CloudOrganizationModel) (*cloudorganizationpb.CloudOrganization, diag.Diagnostics) {
	discoveryInterval, diags := cloudprovider.DurationToProto(model.DiscoveryInterval)
	if diags.HasError() {
		return nil, diags
	}

	return &cloudorganizationpb.CloudOrganization{
		OrgId: model.OrgID.ValueString(),
		Name:  model.Name.ValueString(),
		Provider: &cloudorganizationpb.CloudOrganization_Aws{
			Aws: awsOrganizationToProto(model.AWS),
		},
		DiscoveryEnabled:  model.DiscoveryEnabled.ValueBool(),
		DiscoveryInterval: discoveryInterval,
	}, diags
}

func awsOrganizationToProto(model *AWSOrganizationModel) *cloudorganizationpb.AWSOrganization {
	if model == nil {
		return nil
	}

	return &cloudorganizationpb.AWSOrganization{
		AwsOrgId:          model.AWSOrgID.ValueString(),
		Audience:          model.Audience.ValueString(),
		AssumeThroughOidc: model.AssumeThroughOidc.ValueBool(),
		RoleChain:         cloudprovider.RoleChainToProto(model.RoleChain),
	}
}

// protoToModel converts a CloudOrganization proto message to a CloudOrganizationModel.
func protoToModel(proto *cloudorganizationpb.CloudOrganization) CloudOrganizationModel {
	return CloudOrganizationModel{
		ID:                tftypes.StringValue(proto.GetId()),
		OrgID:             tftypes.StringValue(proto.GetOrgId()),
		Name:              tftypes.StringValue(proto.GetName()),
		AWS:               awsOrganizationFromProto(proto.GetAws()),
		DiscoveryEnabled:  tftypes.BoolValue(proto.GetDiscoveryEnabled()),
		DiscoveryInterval: cloudprovider.DurationToString(proto.GetDiscoveryInterval()),
		Status:            cloudprovider.DiscoveryStatusToString(proto.GetStatus()),
		LastDiscoveredAt:  cloudprovider.TimestampToString(proto.GetLastDiscoveredAt()),
		CreatedAt:         cloudprovider.TimestampToString(proto.GetCreatedAt()),
		LastUpdatedAt:     cloudprovider.TimestampToString(proto.GetLastUpdatedAt()),
	}
}

func awsOrganizationFromProto(proto *cloudorganizationpb.AWSOrganization) *AWSOrganizationModel {
	if proto == nil {
		return nil
	}

	return &AWSOrganizationModel{
		AWSOrgID:          tftypes.StringValue(proto.GetAwsOrgId()),
		Audience:          tftypes.StringValue(proto.GetAudience()),
		AssumeThroughOidc: tftypes.BoolValue(proto.GetAssumeThroughOidc()),
		RoleChain:         cloudprovider.RoleChainFromProto(proto.GetRoleChain()),
	}
}
