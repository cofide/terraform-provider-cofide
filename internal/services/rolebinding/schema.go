package rolebinding

import (
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

var _ resource.ResourceWithConfigValidators = (*RoleBindingResource)(nil)

func resourceSchema() schema.Schema {
	return schema.Schema{
		MarkdownDescription: "Manages a Cofide Connect role binding. Grants a user or group a role on a specific resource. Exactly one of `user` or `group` must be provided.",
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
						Description: "The subject identifier of the user (typically an email address or user ID).",
						Required:    true,
					},
				},
			},
			"group": schema.SingleNestedAttribute{
				Description: "The group principal for the role binding. Exactly one of `user` or `group` must be provided.",
				Optional:    true,
				Attributes: map[string]schema.Attribute{
					"claim_value": schema.StringAttribute{
						Description: "The value of the group claim from the identity provider.",
						Required:    true,
					},
				},
			},
			"resource": schema.SingleNestedAttribute{
				Description: "The resource for the role binding.",
				Required:    true,
				Attributes: map[string]schema.Attribute{
					"type": schema.StringAttribute{
						Description: "The type of the resource to bind the role to. e.g. TrustZone, Cluster",
						Required:    true,
					},
					"id": schema.StringAttribute{
						Description: "The ID of the resource to bind the role to.",
						Required:    true,
					},
				},
			},
		},
	}
}
