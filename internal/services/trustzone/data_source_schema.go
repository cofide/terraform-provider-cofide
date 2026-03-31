package trustzone

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

var _ datasource.DataSource = &TrustZoneDataSource{}

func DataSourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		MarkdownDescription: "Provides information about a Cofide Connect trust zone.",
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
				Description: "The ID of the organization.",
				Optional:    true,
			},
			"trust_domain": schema.StringAttribute{
				Description: "The SPIFFE trust domain for this trust zone (e.g. `example.cofide.dev`).",
				Optional:    true,
			},
			"is_management_zone": schema.BoolAttribute{
				Description: "Whether this is a management trust zone. Cannot be changed after creation.",
				Computed:    true,
			},
			"bundle_endpoint_url": schema.StringAttribute{
				Description: "The URL of the SPIFFE bundle endpoint for this trust zone. Set by Cofide Connect.",
				Computed:    true,
			},
			"bundle_endpoint_profile": schema.StringAttribute{
				Description: "The SPIFFE bundle endpoint profile for this trust zone (`BUNDLE_ENDPOINT_PROFILE_HTTPS_SPIFFE` or `BUNDLE_ENDPOINT_PROFILE_HTTPS_WEB`). Set by Cofide Connect.",
				Computed:    true,
			},
			"jwt_issuer": schema.StringAttribute{
				Description: "The JWT issuer URL for this trust zone. Set by Cofide Connect.",
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
