package rolebinding

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

// oneOfValidator implements resource.ConfigValidator.
type oneOfValidator struct{}

func (v *oneOfValidator) Description(ctx context.Context) string {
	return "Ensures that exactly one of 'user' or 'group' is set."
}

func (v *oneOfValidator) MarkdownDescription(ctx context.Context) string {
	return "Ensures that exactly one of `user` or `group` is set."
}

func (v *oneOfValidator) ValidateResource(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	var data RoleBindingModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	userSet := data.User != nil && !data.User.Subject.IsUnknown() && !data.User.Subject.IsNull()
	groupSet := data.Group != nil && !data.Group.ClaimValue.IsUnknown() && !data.Group.ClaimValue.IsNull()

	if userSet && groupSet {
		resp.Diagnostics.AddAttributeError(
			path.Root("user"),
			"Conflicting Attributes",
			"Cannot set both 'user' and 'group'. Please choose only one.",
		)
		resp.Diagnostics.AddAttributeError(
			path.Root("group"),
			"Conflicting Attributes",
			"Cannot set both 'user' and 'group'. Please choose only one.",
		)
	} else if !userSet && !groupSet {
		resp.Diagnostics.AddAttributeError(
			path.Root("user"),
			"Missing Required Attribute",
			"Either 'user' or 'group' must be set.",
		)
		resp.Diagnostics.AddAttributeError(
			path.Root("group"),
			"Missing Required Attribute",
			"Either 'user' or 'group' must be set.",
		)
	}
}
