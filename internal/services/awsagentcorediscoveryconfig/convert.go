package awsagentcorediscoveryconfig

import (
	"context"

	cloudaccountpb "github.com/cofide/cofide-api-sdk/gen/go/proto/cloud_account/v1alpha1"
	"github.com/cofide/terraform-provider-cofide/internal/services/cloudprovider"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	tftypes "github.com/hashicorp/terraform-plugin-framework/types"
)

// modelToProto converts a Model to an AWSAgentCoreDiscoveryConfig proto message.
func modelToProto(ctx context.Context, model Model) (*cloudaccountpb.AWSAgentCoreDiscoveryConfig, diag.Diagnostics) {
	regions, diags := cloudprovider.StringListToProto(ctx, model.Regions)
	if diags.HasError() {
		return nil, diags
	}

	discoveryInterval, intervalDiags := cloudprovider.DurationToProto(model.DiscoveryInterval)
	diags.Append(intervalDiags...)
	if diags.HasError() {
		return nil, diags
	}

	return &cloudaccountpb.AWSAgentCoreDiscoveryConfig{
		Audience:          model.Audience.ValueString(),
		Regions:           regions,
		DiscoveryEnabled:  model.DiscoveryEnabled.ValueBool(),
		RoleChain:         cloudprovider.RoleChainToProto(model.RoleChain),
		DiscoveryInterval: discoveryInterval,
	}, diags
}

// protoToModel converts an AWSAgentCoreDiscoveryConfig proto message to a Model.
func protoToModel(cloudAccountID string, proto *cloudaccountpb.AWSAgentCoreDiscoveryConfig) Model {
	return Model{
		ID:                      tftypes.StringValue(cloudAccountID),
		CloudAccountID:          tftypes.StringValue(cloudAccountID),
		Audience:                tftypes.StringValue(proto.GetAudience()),
		Regions:                 cloudprovider.StringListFromProto(proto.GetRegions()),
		DiscoveryEnabled:        tftypes.BoolValue(proto.GetDiscoveryEnabled()),
		RoleChain:               cloudprovider.RoleChainFromProto(proto.GetRoleChain()),
		DiscoveryInterval:       cloudprovider.DurationToString(proto.GetDiscoveryInterval()),
		Status:                  cloudprovider.DiscoveryStatusToString(proto.GetStatus()),
		LastSuccessfulDiscovery: cloudprovider.TimestampToString(proto.GetLastSuccessfulDiscovery()),
		StatusLastUpdatedAt:     cloudprovider.TimestampToString(proto.GetStatusLastUpdatedAt()),
	}
}
