package federation

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

var _ datasource.DataSource = &FederationDataSource{}

func DataSourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		MarkdownDescription: "Provides information about a federation resource.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The ID of the federation.",
				Computed:    true,
			},
			"org_id": schema.StringAttribute{
				Description: "The ID of the organisation.",
				Required:    true,
			},
			"trust_zone_id": schema.StringAttribute{
				Description: "The ID of the associated trust zone.",
				Required:    true,
			},
			"remote_trust_zone_id": schema.StringAttribute{
				Description: "The ID of the associated remote trust zone.",
				Required:    true,
			},
		},
	}
}

func (f *FederationDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = DataSourceSchema(ctx)
}

func (f *FederationDataSource) ConfigValidators(_ context.Context) []datasource.ConfigValidator {
	return []datasource.ConfigValidator{}
}
