package cloudaccount

import (
	"context"

	cloudaccountpb "github.com/cofide/cofide-api-sdk/gen/go/proto/cloud_account/v1alpha1"
	"github.com/cofide/terraform-provider-cofide/internal/services/cloudprovider"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	tftypes "github.com/hashicorp/terraform-plugin-framework/types"
)

// modelToProto converts a CloudAccountModel to a CloudAccount proto message.
func modelToProto(ctx context.Context, model CloudAccountModel) (*cloudaccountpb.CloudAccount, diag.Diagnostics) {
	aws, diags := awsAccountToProto(ctx, model.AWS)
	if diags.HasError() {
		return nil, diags
	}

	cloudAccount := &cloudaccountpb.CloudAccount{
		OrgId:               model.OrgID.ValueString(),
		CloudOrganizationId: model.CloudOrganizationID.ValueStringPointer(),
		Name:                model.Name.ValueString(),
		Provider: &cloudaccountpb.CloudAccount_Aws{
			Aws: aws,
		},
		Suppressed:         model.Suppressed.ValueBool(),
		ManagedByDiscovery: model.ManagedByDiscovery.ValueBool(),
	}

	return cloudAccount, diags
}

func awsAccountToProto(ctx context.Context, model *AWSAccountModel) (*cloudaccountpb.AWSAccount, diag.Diagnostics) {
	if model == nil {
		return nil, nil
	}

	var diags diag.Diagnostics

	lambdaConfig, lambdaDiags := lambdaDiscoveryConfigToProto(ctx, model.LambdaDiscoveryConfig)
	diags.Append(lambdaDiags...)

	agentCoreConfig, agentCoreDiags := agentCoreDiscoveryConfigToProto(ctx, model.AgentCoreDiscoveryConfig)
	diags.Append(agentCoreDiags...)

	if diags.HasError() {
		return nil, diags
	}

	return &cloudaccountpb.AWSAccount{
		AccountId:                model.AccountID.ValueString(),
		LambdaDiscoveryConfig:    lambdaConfig,
		AgentCoreDiscoveryConfig: agentCoreConfig,
	}, diags
}

func lambdaDiscoveryConfigToProto(ctx context.Context, model *AWSDiscoveryConfigModel) (*cloudaccountpb.AWSLambdaDiscoveryConfig, diag.Diagnostics) {
	if model == nil {
		return nil, nil
	}

	regions, diags := cloudprovider.StringListToProto(ctx, model.Regions)
	if diags.HasError() {
		return nil, diags
	}

	return &cloudaccountpb.AWSLambdaDiscoveryConfig{
		Audience:         model.Audience.ValueString(),
		Regions:          regions,
		DiscoveryEnabled: model.DiscoveryEnabled.ValueBool(),
		RoleChain:        cloudprovider.RoleChainToProto(model.RoleChain),
	}, diags
}

func agentCoreDiscoveryConfigToProto(ctx context.Context, model *AWSDiscoveryConfigModel) (*cloudaccountpb.AWSAgentCoreDiscoveryConfig, diag.Diagnostics) {
	if model == nil {
		return nil, nil
	}

	regions, diags := cloudprovider.StringListToProto(ctx, model.Regions)
	if diags.HasError() {
		return nil, diags
	}

	return &cloudaccountpb.AWSAgentCoreDiscoveryConfig{
		Audience:         model.Audience.ValueString(),
		Regions:          regions,
		DiscoveryEnabled: model.DiscoveryEnabled.ValueBool(),
		RoleChain:        cloudprovider.RoleChainToProto(model.RoleChain),
	}, diags
}

// protoToModel converts a CloudAccount proto message to a CloudAccountModel.
func protoToModel(proto *cloudaccountpb.CloudAccount) CloudAccountModel {
	return CloudAccountModel{
		ID:                  tftypes.StringValue(proto.GetId()),
		OrgID:               tftypes.StringValue(proto.GetOrgId()),
		CloudOrganizationID: optionalStringValue(proto.CloudOrganizationId),
		Name:                tftypes.StringValue(proto.GetName()),
		AWS:                 awsAccountFromProto(proto.GetAws()),
		Suppressed:          tftypes.BoolValue(proto.GetSuppressed()),
		ManagedByDiscovery:  tftypes.BoolValue(proto.GetManagedByDiscovery()),
		CreatedAt:           cloudprovider.TimestampToString(proto.GetCreatedAt()),
		LastUpdatedAt:       cloudprovider.TimestampToString(proto.GetLastUpdatedAt()),
	}
}

func awsAccountFromProto(proto *cloudaccountpb.AWSAccount) *AWSAccountModel {
	if proto == nil {
		return nil
	}

	return &AWSAccountModel{
		AccountID:                tftypes.StringValue(proto.GetAccountId()),
		LambdaDiscoveryConfig:    lambdaDiscoveryConfigFromProto(proto.GetLambdaDiscoveryConfig()),
		AgentCoreDiscoveryConfig: agentCoreDiscoveryConfigFromProto(proto.GetAgentCoreDiscoveryConfig()),
	}
}

func lambdaDiscoveryConfigFromProto(proto *cloudaccountpb.AWSLambdaDiscoveryConfig) *AWSDiscoveryConfigModel {
	if proto == nil {
		return nil
	}

	return &AWSDiscoveryConfigModel{
		Audience:                tftypes.StringValue(proto.GetAudience()),
		Regions:                 cloudprovider.StringListFromProto(proto.GetRegions()),
		DiscoveryEnabled:        tftypes.BoolValue(proto.GetDiscoveryEnabled()),
		RoleChain:               cloudprovider.RoleChainFromProto(proto.GetRoleChain()),
		Status:                  cloudprovider.DiscoveryStatusToString(proto.GetStatus()),
		LastSuccessfulDiscovery: cloudprovider.TimestampToString(proto.GetLastSuccessfulDiscovery()),
		StatusLastUpdatedAt:     cloudprovider.TimestampToString(proto.GetStatusLastUpdatedAt()),
	}
}

func agentCoreDiscoveryConfigFromProto(proto *cloudaccountpb.AWSAgentCoreDiscoveryConfig) *AWSDiscoveryConfigModel {
	if proto == nil {
		return nil
	}

	return &AWSDiscoveryConfigModel{
		Audience:                tftypes.StringValue(proto.GetAudience()),
		Regions:                 cloudprovider.StringListFromProto(proto.GetRegions()),
		DiscoveryEnabled:        tftypes.BoolValue(proto.GetDiscoveryEnabled()),
		RoleChain:               cloudprovider.RoleChainFromProto(proto.GetRoleChain()),
		Status:                  cloudprovider.DiscoveryStatusToString(proto.GetStatus()),
		LastSuccessfulDiscovery: cloudprovider.TimestampToString(proto.GetLastSuccessfulDiscovery()),
		StatusLastUpdatedAt:     cloudprovider.TimestampToString(proto.GetStatusLastUpdatedAt()),
	}
}

// optionalStringValue converts a *string to a null-safe Terraform string.
func optionalStringValue(s *string) tftypes.String {
	if s == nil {
		return tftypes.StringNull()
	}
	return tftypes.StringValue(*s)
}
