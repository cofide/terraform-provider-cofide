package rolebinding

import (
	rolebindingpb "github.com/cofide/cofide-api-sdk/gen/go/proto/role_binding/v1alpha1"
	tftypes "github.com/hashicorp/terraform-plugin-framework/types"
)

// modelToProto converts an RoleBindingModel to an equivalent RoleBinding protobuf.
func modelToProto(model RoleBindingModel) *rolebindingpb.RoleBinding {
	proto := &rolebindingpb.RoleBinding{
		Id:     model.ID.ValueString(),
		RoleId: model.RoleID.ValueString(),
		Resource: &rolebindingpb.Resource{
			Type: model.Resource.Type.ValueString(),
			Id:   model.Resource.ID.ValueString(),
		},
	}

	if model.User != nil {
		proto.Principal = &rolebindingpb.RoleBinding_User{
			User: &rolebindingpb.User{
				Subject: model.User.Subject.ValueString(),
			},
		}
	}

	if model.Group != nil {
		proto.Principal = &rolebindingpb.RoleBinding_Group{
			Group: &rolebindingpb.Group{
				ClaimValue: model.Group.ClaimValue.ValueString(),
			},
		}
	}

	return proto
}

// protoToModel converts an RoleBinding protobuf to an equivalent RoleBindingModel.
func protoToModel(proto *rolebindingpb.RoleBinding) RoleBindingModel {
	model := RoleBindingModel{
		ID:     tftypes.StringValue(proto.GetId()),
		RoleID: tftypes.StringValue(proto.GetRoleId()),
		Resource: ResourceModel{
			Type: tftypes.StringValue(proto.GetResource().GetType()),
			ID:   tftypes.StringValue(proto.GetResource().GetId()),
		},
	}

	if userProto := proto.GetUser(); userProto != nil {
		model.User = &UserModel{
			Subject: tftypes.StringValue(userProto.GetSubject()),
		}
	}

	if groupProto := proto.GetGroup(); groupProto != nil {
		model.Group = &GroupModel{
			ClaimValue: tftypes.StringValue(groupProto.GetClaimValue()),
		}
	}

	return model
}
