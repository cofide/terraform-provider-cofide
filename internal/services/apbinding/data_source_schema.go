package apbinding

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ datasource.DataSource = &APBindingDataSource{}

func DataSourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		MarkdownDescription: "Provides information about an attestation policy binding resource.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The ID of the attestation policy binding.",
				Computed:    true,
			},
			"org_id": schema.StringAttribute{
				Description: "The ID of the organisation.",
				Required:    true,
			},
			"trust_zone_id": schema.StringAttribute{
				Description: "The ID of the trust zone.",
				Required:    true,
			},
			"policy_id": schema.StringAttribute{
				Description: "The ID of the attestation policy.",
				Required:    true,
			},
			"federations": schema.ListAttribute{
				Description: "The list of associated federations.",
				Computed:    true,
				ElementType: types.ObjectType{
					AttrTypes: map[string]attr.Type{
						"trust_zone_id": types.StringType,
					},
				},
			},
		},
	}
}

func (a *APBindingDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = DataSourceSchema(ctx)
}

func (a *APBindingDataSource) ConfigValidators(_ context.Context) []datasource.ConfigValidator {
	return []datasource.ConfigValidator{}
}
