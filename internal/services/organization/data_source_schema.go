package organization

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

func getOrganizationDataSourceSchema() schema.Schema {
	return schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Provides information about an organization resource.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Organization identifier",
				Computed:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Organization name",
				Required:            true,
			},
		},
	}
}
