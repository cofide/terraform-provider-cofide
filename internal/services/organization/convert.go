package organization

import (
	organizationpb "github.com/cofide/cofide-api-sdk/gen/go/proto/organization/v1alpha1"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// protoToModel converts an Organization protobuf to an equivalent OrganizationModel.
func protoToModel(proto *organizationpb.Organization) OrganizationModel {
	return OrganizationModel{
		ID:   types.StringValue(proto.GetId()),
		Name: types.StringValue(proto.GetName()),
	}
}
