package cloudorganization

import (
	"testing"

	cloudorganizationpb "github.com/cofide/cofide-api-sdk/gen/go/proto/cloud_organization/v1alpha1"
	cloudproviderpb "github.com/cofide/cofide-api-sdk/gen/go/proto/cloud_provider/v1alpha1"
	"github.com/cofide/terraform-provider-cofide/internal/services/cloudprovider"
	tftypes "github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestProtoToModel(t *testing.T) {
	now := timestamppb.Now()
	proto := &cloudorganizationpb.CloudOrganization{
		Id:     "co-1",
		OrgId:  "org-1",
		Name:   "test-cloud-org",
		Status: cloudproviderpb.DiscoveryStatus_DISCOVERY_STATUS_DISCOVERING,
		Provider: &cloudorganizationpb.CloudOrganization_Aws{
			Aws: &cloudorganizationpb.AWSOrganization{
				AwsOrgId: "o-fakeorgid12",
				Audience: "aud",
				RoleChain: []*cloudproviderpb.AWSAssumeRoleConfig{
					{IamRoleArn: "arn:aws:iam::123456789012:role/first"},
				},
			},
		},
		DiscoveryEnabled: true,
		LastDiscoveredAt: now,
		CreatedAt:        now,
		LastUpdatedAt:    now,
	}

	got := protoToModel(proto)

	assert.Equal(t, tftypes.StringValue("co-1"), got.ID)
	assert.Equal(t, tftypes.StringValue("org-1"), got.OrgID)
	assert.Equal(t, tftypes.StringValue("test-cloud-org"), got.Name)
	assert.Equal(t, tftypes.BoolValue(true), got.DiscoveryEnabled)
	assert.Equal(t, tftypes.StringValue("DISCOVERING"), got.Status)
	assert.Equal(t, cloudprovider.TimestampToString(now), got.LastDiscoveredAt)
	assert.Equal(t, cloudprovider.TimestampToString(now), got.CreatedAt)
	assert.Equal(t, cloudprovider.TimestampToString(now), got.LastUpdatedAt)

	if assert.NotNil(t, got.AWS) {
		assert.Equal(t, tftypes.StringValue("o-fakeorgid12"), got.AWS.AWSOrgID)
		assert.Equal(t, tftypes.StringValue("aud"), got.AWS.Audience)
		assert.Equal(t, []cloudprovider.RoleChainModel{
			{IAMRoleARN: tftypes.StringValue("arn:aws:iam::123456789012:role/first"), ExternalID: tftypes.StringNull()},
		}, got.AWS.RoleChain)
	}
}

func TestModelToProto(t *testing.T) {
	model := CloudOrganizationModel{
		OrgID: tftypes.StringValue("org-1"),
		Name:  tftypes.StringValue("test-cloud-org"),
		AWS: &AWSOrganizationModel{
			AWSOrgID: tftypes.StringValue("o-fakeorgid12"),
			Audience: tftypes.StringValue("aud"),
			RoleChain: []cloudprovider.RoleChainModel{
				{IAMRoleARN: tftypes.StringValue("arn:aws:iam::123456789012:role/first"), ExternalID: tftypes.StringNull()},
			},
		},
		DiscoveryEnabled: tftypes.BoolValue(true),
	}

	got := modelToProto(model)

	assert.Equal(t, "org-1", got.GetOrgId())
	assert.Equal(t, "test-cloud-org", got.GetName())
	assert.True(t, got.GetDiscoveryEnabled())
	assert.Equal(t, "o-fakeorgid12", got.GetAws().GetAwsOrgId())
	assert.Equal(t, "aud", got.GetAws().GetAudience())
	assert.Equal(t, []*cloudproviderpb.AWSAssumeRoleConfig{
		{IamRoleArn: "arn:aws:iam::123456789012:role/first"},
	}, got.GetAws().GetRoleChain())
}
