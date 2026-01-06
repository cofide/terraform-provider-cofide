package trustzone

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

var _ datasource.DataSource = &TrustZoneDataSource{}

func DataSourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		MarkdownDescription: "Provides information about a trust zone resource.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The ID of the trust zone.",
				Computed:    true,
			},
			"name": schema.StringAttribute{
				Description: "The name of the trust zone.",
				Optional:    true,
			},
			"org_id": schema.StringAttribute{
				Description: "The ID of the organisation.",
				Optional:    true,
			},
			"trust_domain": schema.StringAttribute{
				Description: "The trust domain of the trust zone.",
				Optional:    true,
			},
			"is_management_zone": schema.BoolAttribute{
				Description: "Whether or not this is a management trust zone.",
				Computed:    true,
			},
			"bundle_endpoint_url": schema.StringAttribute{
				Description: "The bundle endpoint URL of the trust zone.",
				Computed:    true,
			},
			"bundle_endpoint_profile": schema.StringAttribute{
				Description: "The bundle endpoint profile of the trust zone.",
				Computed:    true,
			},
			"jwt_issuer": schema.StringAttribute{
				Description: "The JWT issuer of the trust zone.",
				Computed:    true,
			},
		},
	}
}

func (t *TrustZoneDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = DataSourceSchema(ctx)
}

func (t *TrustZoneDataSource) ConfigValidators(_ context.Context) []datasource.ConfigValidator {
	return []datasource.ConfigValidator{}
}
