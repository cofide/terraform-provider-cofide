package rolebinding

import (
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

var _ resource.ResourceWithConfigValidators = (*RoleBindingResource)(nil)

func resourceSchema() schema.Schema {
	return schema.Schema{
		MarkdownDescription: "Provides a role binding resource.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The ID of the role binding.",
				Computed:    true,
			},
			"role_id": schema.StringAttribute{
				Description: "The ID of the role.",
				Required:    true,
			},
			"user": schema.SingleNestedAttribute{
				Description: "The user principal for the role binding. Exactly one of `user` or `group` must be provided.",
				Optional:    true,
				Attributes: map[string]schema.Attribute{
					"subject": schema.StringAttribute{
						Description: "The subject of the user.",
						Required:    true,
					},
				},
			},
			"group": schema.SingleNestedAttribute{
				Description: "The group principal for the role binding. Exactly one of `user` or `group` must be provided.",
				Optional:    true,
				Attributes: map[string]schema.Attribute{
					"claim_value": schema.StringAttribute{
						Description: "The claim value of the group.",
						Required:    true,
					},
				},
			},
			"resource": schema.SingleNestedAttribute{
				Description: "The resource for the role binding.",
				Required:    true,
				Attributes: map[string]schema.Attribute{
					"type": schema.StringAttribute{
						Description: "The type of the resource.",
						Required:    true,
					},
					"id": schema.StringAttribute{
						Description: "The ID of the resource.",
						Required:    true,
					},
				},
			},
		},
	}
}
