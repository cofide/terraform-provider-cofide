package organization

import (
	organizationpb "github.com/cofide/cofide-api-sdk/gen/go/proto/organization/v1alpha1"
	tftypes "github.com/hashicorp/terraform-plugin-framework/types"
)

func protoToModel(proto *organizationpb.Organization) OrganizationModel {
	model := OrganizationModel{
		ID:   tftypes.StringValue(proto.GetId()),
		Name: tftypes.StringValue(proto.GetName()),
	}
	return model
}
