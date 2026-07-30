package cloudaccount

import (
	"testing"

	cloudaccountsvcpb "github.com/cofide/cofide-api-sdk/gen/go/proto/connect/cloud_account_service/v1alpha1"
	tftypes "github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
)

func TestUpdateMaskForConfig(t *testing.T) {
	tests := []struct {
		name   string
		config CloudAccountModel
		want   *cloudaccountsvcpb.UpdateCloudAccountRequest_UpdateMask
	}{
		{
			name: "nothing set beyond required fields",
			config: CloudAccountModel{
				Name:               tftypes.StringValue("test"),
				AWS:                &AWSAccountModel{AccountID: tftypes.StringValue("123456789012")},
				Suppressed:         tftypes.BoolNull(),
				ManagedByDiscovery: tftypes.BoolNull(),
			},
			want: &cloudaccountsvcpb.UpdateCloudAccountRequest_UpdateMask{
				Name: true,
			},
		},
		{
			name: "suppressed and managed_by_discovery explicitly set",
			config: CloudAccountModel{
				Name:               tftypes.StringValue("test"),
				AWS:                &AWSAccountModel{AccountID: tftypes.StringValue("123456789012")},
				Suppressed:         tftypes.BoolValue(true),
				ManagedByDiscovery: tftypes.BoolValue(false),
			},
			want: &cloudaccountsvcpb.UpdateCloudAccountRequest_UpdateMask{
				Name:               true,
				Suppressed:         true,
				ManagedByDiscovery: true,
			},
		},
		{
			name: "lambda discovery config set inline, agent core config omitted",
			config: CloudAccountModel{
				Name: tftypes.StringValue("test"),
				AWS: &AWSAccountModel{
					AccountID:             tftypes.StringValue("123456789012"),
					LambdaDiscoveryConfig: &AWSDiscoveryConfigModel{Audience: tftypes.StringValue("aud")},
				},
			},
			want: &cloudaccountsvcpb.UpdateCloudAccountRequest_UpdateMask{
				Name:                     true,
				AwsLambdaDiscoveryConfig: true,
			},
		},
		{
			name: "both discovery configs set inline",
			config: CloudAccountModel{
				Name: tftypes.StringValue("test"),
				AWS: &AWSAccountModel{
					AccountID:                tftypes.StringValue("123456789012"),
					LambdaDiscoveryConfig:    &AWSDiscoveryConfigModel{Audience: tftypes.StringValue("aud")},
					AgentCoreDiscoveryConfig: &AWSDiscoveryConfigModel{Audience: tftypes.StringValue("aud")},
				},
			},
			want: &cloudaccountsvcpb.UpdateCloudAccountRequest_UpdateMask{
				Name:                        true,
				AwsLambdaDiscoveryConfig:    true,
				AwsAgentCoreDiscoveryConfig: true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := updateMaskForConfig(tt.config)
			assert.Equal(t, tt.want, got)
		})
	}
}
