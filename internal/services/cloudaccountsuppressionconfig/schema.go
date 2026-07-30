package cloudaccountsuppressionconfig

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
)

func ResourceSchema(_ context.Context) schema.Schema {
	return schema.Schema{
		MarkdownDescription: "Manages the suppression state of a Cofide Connect cloud account, independently of the `cofide_connect_cloud_account` resource. Useful when the cloud account's existence is managed by automatic discovery rather than Terraform: only the `suppressed` flag is managed here, referencing the account by ID. Do not manage this flag both here and inline on `cofide_connect_cloud_account.suppressed`.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The ID of this resource. Matches `cloud_account_id`.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"cloud_account_id": schema.StringAttribute{
				Description: "The ID of the cloud account this suppression setting applies to. Cannot be changed after creation.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"suppressed": schema.BoolAttribute{
				Description: "When true, discovery is suspended for this account and existing discovered resources are hidden from findings.",
				Required:    true,
			},
		},
	}
}

func (r *Resource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = ResourceSchema(ctx)
}
