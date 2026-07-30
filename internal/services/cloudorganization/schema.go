package cloudorganization

import (
	"context"

	"github.com/cofide/terraform-provider-cofide/internal/services/cloudprovider"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
)

func ResourceSchema(_ context.Context) schema.Schema {
	return schema.Schema{
		MarkdownDescription: "Manages a Cofide Connect cloud organization. A cloud organization represents a cloud provider organization (e.g. an AWS Organization) linked to a Cofide organization and acts as the root for cloud account and resource discovery.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The ID of the cloud organization.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"org_id": schema.StringAttribute{
				Description: "The ID of the Cofide organization this cloud organization belongs to. Cannot be changed after creation.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"name": schema.StringAttribute{
				Description: "The name of the cloud organization.",
				Required:    true,
			},
			"aws": schema.SingleNestedAttribute{
				Description: "AWS-specific configuration for the cloud organization.",
				Required:    true,
				Attributes: map[string]schema.Attribute{
					"aws_org_id": schema.StringAttribute{
						Description: "The AWS Organization ID (e.g. `o-xxxxxxxxxx`).",
						Required:    true,
					},
					"audience": schema.StringAttribute{
						Description: "Audience value for the initial SPIFFE JWT-based assume role call.",
						Required:    true,
					},
					"role_chain": cloudprovider.RoleChainAttribute("Ordered chain of IAM roles to assume when discovering resources in this cloud organization."),
				},
			},
			"discovery_enabled": schema.BoolAttribute{
				Description: "Whether discovery is enabled for this cloud organization.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
			"status": schema.StringAttribute{
				Description: "The current status of cloud resource discovery (e.g. `DISCOVERING`, `ERROR`).",
				Computed:    true,
			},
			"last_discovered_at": schema.StringAttribute{
				Description: "The timestamp (RFC3339) of the last successful discovery.",
				Computed:    true,
			},
			"created_at": schema.StringAttribute{
				Description: "The timestamp (RFC3339) of resource creation.",
				Computed:    true,
			},
			"last_updated_at": schema.StringAttribute{
				Description: "The timestamp (RFC3339) of the last resource update.",
				Computed:    true,
			},
		},
	}
}

func (c *CloudOrganizationResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = ResourceSchema(ctx)
}
