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
		MarkdownDescription: "Provides information about a Cofide Connect attestation policy binding.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The ID of the attestation policy binding.",
				Computed:    true,
			},
			"org_id": schema.StringAttribute{
				Description: "The ID of the organisation.",
				Optional:    true,
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
				Description: "The federated trust zones which will be visible to workloads matching the policy in this binding. Each entry specifies the `trust_zone_id` of a federated trust zone.",
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
