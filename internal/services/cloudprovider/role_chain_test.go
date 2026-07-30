package cloudprovider

import (
	"testing"

	cloudproviderpb "github.com/cofide/cofide-api-sdk/gen/go/proto/cloud_provider/v1alpha1"
	tftypes "github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
)

func TestRoleChainToProto(t *testing.T) {
	externalID := "external-id"

	tests := []struct {
		name      string
		roleChain []RoleChainModel
		want      []*cloudproviderpb.AWSAssumeRoleConfig
	}{
		{
			name:      "empty",
			roleChain: nil,
			want:      nil,
		},
		{
			name: "single step without external id",
			roleChain: []RoleChainModel{
				{IAMRoleARN: tftypes.StringValue("arn:aws:iam::123456789012:role/first"), ExternalID: tftypes.StringNull()},
			},
			want: []*cloudproviderpb.AWSAssumeRoleConfig{
				{IamRoleArn: "arn:aws:iam::123456789012:role/first"},
			},
		},
		{
			name: "multi-step chain with external id",
			roleChain: []RoleChainModel{
				{IAMRoleARN: tftypes.StringValue("arn:aws:iam::123456789012:role/first"), ExternalID: tftypes.StringValue(externalID)},
				{IAMRoleARN: tftypes.StringValue("arn:aws:iam::123456789012:role/second"), ExternalID: tftypes.StringNull()},
			},
			want: []*cloudproviderpb.AWSAssumeRoleConfig{
				{IamRoleArn: "arn:aws:iam::123456789012:role/first", ExternalId: &externalID},
				{IamRoleArn: "arn:aws:iam::123456789012:role/second"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := RoleChainToProto(tt.roleChain)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestRoleChainFromProto(t *testing.T) {
	externalID := "external-id"

	tests := []struct {
		name      string
		roleChain []*cloudproviderpb.AWSAssumeRoleConfig
		want      []RoleChainModel
	}{
		{
			name:      "empty",
			roleChain: nil,
			want:      nil,
		},
		{
			name: "single step without external id",
			roleChain: []*cloudproviderpb.AWSAssumeRoleConfig{
				{IamRoleArn: "arn:aws:iam::123456789012:role/first"},
			},
			want: []RoleChainModel{
				{IAMRoleARN: tftypes.StringValue("arn:aws:iam::123456789012:role/first"), ExternalID: tftypes.StringNull()},
			},
		},
		{
			name: "multi-step chain with external id",
			roleChain: []*cloudproviderpb.AWSAssumeRoleConfig{
				{IamRoleArn: "arn:aws:iam::123456789012:role/first", ExternalId: &externalID},
				{IamRoleArn: "arn:aws:iam::123456789012:role/second"},
			},
			want: []RoleChainModel{
				{IAMRoleARN: tftypes.StringValue("arn:aws:iam::123456789012:role/first"), ExternalID: tftypes.StringValue(externalID)},
				{IAMRoleARN: tftypes.StringValue("arn:aws:iam::123456789012:role/second"), ExternalID: tftypes.StringNull()},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := RoleChainFromProto(tt.roleChain)
			assert.Equal(t, tt.want, got)
		})
	}
}
