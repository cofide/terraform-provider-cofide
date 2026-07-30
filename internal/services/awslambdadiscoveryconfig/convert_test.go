package awslambdadiscoveryconfig

import (
	"context"
	"testing"

	cloudaccountpb "github.com/cofide/cofide-api-sdk/gen/go/proto/cloud_account/v1alpha1"
	cloudproviderpb "github.com/cofide/cofide-api-sdk/gen/go/proto/cloud_provider/v1alpha1"
	"github.com/cofide/terraform-provider-cofide/internal/services/cloudprovider"
	tftypes "github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestModelToProto(t *testing.T) {
	ctx := context.Background()

	model := Model{
		Audience:         tftypes.StringValue("aud"),
		Regions:          tftypes.ListNull(tftypes.StringType),
		DiscoveryEnabled: tftypes.BoolValue(true),
		RoleChain: []cloudprovider.RoleChainModel{
			{IAMRoleARN: tftypes.StringValue("arn:aws:iam::123456789012:role/lambda"), ExternalID: tftypes.StringNull()},
		},
	}

	got, diags := modelToProto(ctx, model)
	require.False(t, diags.HasError())

	assert.Equal(t, "aud", got.GetAudience())
	assert.True(t, got.GetDiscoveryEnabled())
	assert.Equal(t, []*cloudproviderpb.AWSAssumeRoleConfig{
		{IamRoleArn: "arn:aws:iam::123456789012:role/lambda"},
	}, got.GetRoleChain())
}

func TestProtoToModel(t *testing.T) {
	proto := &cloudaccountpb.AWSLambdaDiscoveryConfig{
		Audience:         "aud",
		Regions:          []string{"eu-west-1"},
		DiscoveryEnabled: true,
		Status:           cloudproviderpb.DiscoveryStatus_DISCOVERY_STATUS_ERROR,
		RoleChain: []*cloudproviderpb.AWSAssumeRoleConfig{
			{IamRoleArn: "arn:aws:iam::123456789012:role/lambda"},
		},
	}

	got := protoToModel("ca-1", proto)

	assert.Equal(t, tftypes.StringValue("ca-1"), got.ID)
	assert.Equal(t, tftypes.StringValue("ca-1"), got.CloudAccountID)
	assert.Equal(t, tftypes.StringValue("aud"), got.Audience)
	assert.Equal(t, tftypes.StringValue("ERROR"), got.Status)
	assert.True(t, got.DiscoveryEnabled.ValueBool())
}
