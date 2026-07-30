package cloudaccount

import (
	"context"
	"testing"
	"time"

	cloudaccountpb "github.com/cofide/cofide-api-sdk/gen/go/proto/cloud_account/v1alpha1"
	cloudproviderpb "github.com/cofide/cofide-api-sdk/gen/go/proto/cloud_provider/v1alpha1"
	"github.com/cofide/terraform-provider-cofide/internal/services/cloudprovider"
	tftypes "github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestProtoToModel_Minimal(t *testing.T) {
	proto := &cloudaccountpb.CloudAccount{
		Id:    "ca-1",
		OrgId: "org-1",
		Name:  "test-cloud-account",
	}

	got := protoToModel(proto)

	assert.Equal(t, tftypes.StringValue("ca-1"), got.ID)
	assert.Equal(t, tftypes.StringValue("org-1"), got.OrgID)
	assert.True(t, got.CloudOrganizationID.IsNull())
	assert.Equal(t, tftypes.StringValue("test-cloud-account"), got.Name)
	assert.Nil(t, got.AWS)
	assert.Equal(t, tftypes.BoolValue(false), got.Suppressed)
	assert.Equal(t, tftypes.BoolValue(false), got.ManagedByDiscovery)
}

func TestProtoToModel_Full(t *testing.T) {
	now := timestamppb.Now()
	cloudOrgID := "co-1"

	proto := &cloudaccountpb.CloudAccount{
		Id:                  "ca-2",
		OrgId:               "org-2",
		CloudOrganizationId: &cloudOrgID,
		Name:                "full-cloud-account",
		Provider: &cloudaccountpb.CloudAccount_Aws{
			Aws: &cloudaccountpb.AWSAccount{
				AccountId: "123456789012",
				LambdaDiscoveryConfig: &cloudaccountpb.AWSLambdaDiscoveryConfig{
					Audience:          "lambda-aud",
					Regions:           []string{"eu-west-1"},
					DiscoveryEnabled:  true,
					Status:            cloudproviderpb.DiscoveryStatus_DISCOVERY_STATUS_DISCOVERING,
					DiscoveryInterval: durationpb.New(time.Minute),
					RoleChain: []*cloudproviderpb.AWSAssumeRoleConfig{
						{IamRoleArn: "arn:aws:iam::123456789012:role/lambda"},
					},
				},
				AgentCoreDiscoveryConfig: &cloudaccountpb.AWSAgentCoreDiscoveryConfig{
					Audience: "agent-core-aud",
					RoleChain: []*cloudproviderpb.AWSAssumeRoleConfig{
						{IamRoleArn: "arn:aws:iam::123456789012:role/agent-core"},
					},
				},
			},
		},
		Suppressed:         true,
		ManagedByDiscovery: true,
		CreatedAt:          now,
		LastUpdatedAt:      now,
	}

	got := protoToModel(proto)

	assert.Equal(t, tftypes.StringValue("co-1"), got.CloudOrganizationID)
	assert.True(t, got.Suppressed.ValueBool())
	assert.True(t, got.ManagedByDiscovery.ValueBool())
	assert.Equal(t, cloudprovider.TimestampToString(now), got.CreatedAt)

	require.NotNil(t, got.AWS)
	assert.Equal(t, tftypes.StringValue("123456789012"), got.AWS.AccountID)

	require.NotNil(t, got.AWS.LambdaDiscoveryConfig)
	assert.Equal(t, tftypes.StringValue("lambda-aud"), got.AWS.LambdaDiscoveryConfig.Audience)
	assert.Equal(t, tftypes.StringValue("DISCOVERING"), got.AWS.LambdaDiscoveryConfig.Status)
	assert.Equal(t, tftypes.StringValue("1m0s"), got.AWS.LambdaDiscoveryConfig.DiscoveryInterval)
	assert.Equal(t, []cloudprovider.RoleChainModel{
		{IAMRoleARN: tftypes.StringValue("arn:aws:iam::123456789012:role/lambda"), ExternalID: tftypes.StringNull()},
	}, got.AWS.LambdaDiscoveryConfig.RoleChain)

	require.NotNil(t, got.AWS.AgentCoreDiscoveryConfig)
	assert.Equal(t, tftypes.StringValue("agent-core-aud"), got.AWS.AgentCoreDiscoveryConfig.Audience)
	assert.True(t, got.AWS.AgentCoreDiscoveryConfig.Regions.IsNull())
}

func TestModelToProto(t *testing.T) {
	ctx := context.Background()

	model := CloudAccountModel{
		OrgID:               tftypes.StringValue("org-1"),
		CloudOrganizationID: tftypes.StringValue("co-1"),
		Name:                tftypes.StringValue("test-cloud-account"),
		AWS: &AWSAccountModel{
			AccountID: tftypes.StringValue("123456789012"),
			LambdaDiscoveryConfig: &AWSDiscoveryConfigModel{
				Audience:          tftypes.StringValue("lambda-aud"),
				Regions:           tftypes.ListNull(tftypes.StringType),
				DiscoveryEnabled:  tftypes.BoolValue(true),
				DiscoveryInterval: tftypes.StringValue("1m"),
				RoleChain: []cloudprovider.RoleChainModel{
					{IAMRoleARN: tftypes.StringValue("arn:aws:iam::123456789012:role/lambda"), ExternalID: tftypes.StringNull()},
				},
			},
		},
		Suppressed:         tftypes.BoolValue(true),
		ManagedByDiscovery: tftypes.BoolValue(false),
	}

	got, diags := modelToProto(ctx, model)
	require.False(t, diags.HasError())

	assert.Equal(t, "org-1", got.GetOrgId())
	assert.Equal(t, "co-1", got.GetCloudOrganizationId())
	assert.Equal(t, "test-cloud-account", got.GetName())
	assert.True(t, got.GetSuppressed())
	assert.False(t, got.GetManagedByDiscovery())
	assert.Equal(t, "123456789012", got.GetAws().GetAccountId())
	assert.Equal(t, "lambda-aud", got.GetAws().GetLambdaDiscoveryConfig().GetAudience())
	assert.Equal(t, durationpb.New(time.Minute), got.GetAws().GetLambdaDiscoveryConfig().GetDiscoveryInterval())
	assert.Nil(t, got.GetAws().GetAgentCoreDiscoveryConfig())
}
