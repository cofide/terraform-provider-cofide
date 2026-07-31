// Package cloudprovider contains schema, model, and conversion helpers shared
// across the cloud discovery resources (cloud_organization, cloud_account,
// and the AWS discovery config resources), since they all embed the same
// AWS IAM role chain and cloud-provider-agnostic status/timestamp fields.
package cloudprovider

import (
	cloudproviderpb "github.com/cofide/cofide-api-sdk/gen/go/proto/cloud_provider/v1alpha1"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	tftypes "github.com/hashicorp/terraform-plugin-framework/types"
)

// RoleChainModel describes a single step in an AWS IAM role assumption chain.
type RoleChainModel struct {
	IAMRoleARN tftypes.String `tfsdk:"iam_role_arn"`
	ExternalID tftypes.String `tfsdk:"external_id"`
}

// RoleChainAttribute returns the reusable schema for an AWS IAM role chain.
// Whether the first step is assumed via AssumeRoleWithWebIdentity using a
// SPIFFE JWT or via plain AssumeRole using ambient credentials is controlled
// by the assume_through_oidc field on the owning config; each subsequent
// step is always assumed using the credentials from the prior step.
func RoleChainAttribute(description string) schema.ListNestedAttribute {
	return schema.ListNestedAttribute{
		Description: description,
		Required:    true,
		NestedObject: schema.NestedAttributeObject{
			Attributes: map[string]schema.Attribute{
				"iam_role_arn": schema.StringAttribute{
					Description: "ARN of the IAM role to assume.",
					Required:    true,
				},
				"external_id": schema.StringAttribute{
					Description: "Optional external ID for confused-deputy protection.",
					Optional:    true,
				},
			},
		},
	}
}

// RoleChainToProto converts a Terraform role chain model list to its proto representation.
func RoleChainToProto(roleChain []RoleChainModel) []*cloudproviderpb.AWSAssumeRoleConfig {
	if len(roleChain) == 0 {
		return nil
	}

	protoRoleChain := make([]*cloudproviderpb.AWSAssumeRoleConfig, 0, len(roleChain))
	for _, step := range roleChain {
		protoRoleChain = append(protoRoleChain, &cloudproviderpb.AWSAssumeRoleConfig{
			IamRoleArn: step.IAMRoleARN.ValueString(),
			ExternalId: step.ExternalID.ValueStringPointer(),
		})
	}
	return protoRoleChain
}

// RoleChainFromProto converts a proto role chain to the Terraform model list.
func RoleChainFromProto(protoRoleChain []*cloudproviderpb.AWSAssumeRoleConfig) []RoleChainModel {
	if len(protoRoleChain) == 0 {
		return nil
	}

	roleChain := make([]RoleChainModel, 0, len(protoRoleChain))
	for _, step := range protoRoleChain {
		model := RoleChainModel{
			IAMRoleARN: tftypes.StringValue(step.GetIamRoleArn()),
			ExternalID: tftypes.StringNull(),
		}
		if step.ExternalId != nil {
			model.ExternalID = tftypes.StringValue(step.GetExternalId())
		}
		roleChain = append(roleChain, model)
	}
	return roleChain
}
