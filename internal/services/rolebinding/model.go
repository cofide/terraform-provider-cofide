package rolebinding

import tftypes "github.com/hashicorp/terraform-plugin-framework/types"

type RoleBindingModel struct {
	ID       tftypes.String `tfsdk:"id"`
	RoleID   tftypes.String `tfsdk:"role_id"`
	User     *UserModel     `tfsdk:"user"`
	Group    *GroupModel    `tfsdk:"group"`
	Resource ResourceModel  `tfsdk:"resource"`
}

type UserModel struct {
	Subject tftypes.String `tfsdk:"subject"`
}

type GroupModel struct {
	ClaimValue tftypes.String `tfsdk:"claim_value"`
}

type ResourceModel struct {
	Type tftypes.String `tfsdk:"type"`
	ID   tftypes.String `tfsdk:"id"`
}
