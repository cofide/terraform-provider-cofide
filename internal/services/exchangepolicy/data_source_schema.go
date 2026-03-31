package exchangepolicy

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	tftypes "github.com/hashicorp/terraform-plugin-framework/types"
)

var _ datasource.DataSource = &ExchangePolicyDataSource{}
var _ datasource.DataSource = &ExchangePoliciesDataSource{}

func dataSourceSchema() schema.Schema {
	attrs := exchangePolicyNestedAttributes()
	attrs["id"] = schema.StringAttribute{
		Description: "The ID of the exchange policy.",
		Required:    true,
	}
	return schema.Schema{
		MarkdownDescription: "Provides information about a Cofide Connect exchange policy.",
		Attributes:          attrs,
	}
}

func listDataSourceSchema() schema.Schema {
	return schema.Schema{
		MarkdownDescription: "Provides information about Cofide Connect exchange policies.",
		Attributes: map[string]schema.Attribute{
			"trust_zone_id": schema.StringAttribute{
				Description: "Filter by trust zone ID.",
				Optional:    true,
			},
			"org_id": schema.StringAttribute{
				Description: "Filter by organization ID.",
				Optional:    true,
			},
			"name": schema.StringAttribute{
				Description: "Filter by exchange policy name.",
				Optional:    true,
			},
			"exchange_policies": schema.ListNestedAttribute{
				Description: "The list of exchange policies.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: exchangePolicyNestedAttributes(),
				},
			},
		},
	}
}

// exchangePolicyNestedAttributes returns the schema attributes for a nested ExchangePolicy.
func exchangePolicyNestedAttributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Description: "The ID of the exchange policy.",
			Computed:    true,
		},
		"org_id": schema.StringAttribute{
			Description: "The ID of the organization.",
			Computed:    true,
		},
		"name": schema.StringAttribute{
			Description: "The name of the exchange policy.",
			Computed:    true,
		},
		"trust_zone_id": schema.StringAttribute{
			Description: "The ID of the trust zone to which this policy applies.",
			Computed:    true,
		},
		"action": schema.StringAttribute{
			Description: "Action to take when all conditions match.",
			Computed:    true,
		},
		"subject_identity": stringSetDataSourceAttribute("Match conditions on the subject identity of the inbound token."),
		"subject_issuer":   stringSetDataSourceAttribute("Match conditions on the issuer of the inbound subject token."),
		"actor_identity":   stringSetDataSourceAttribute("Match conditions on the actor identity of the inbound token."),
		"actor_issuer":     stringSetDataSourceAttribute("Match conditions on the issuer of the inbound actor token."),
		"client_id":        stringSetDataSourceAttribute("Match conditions on the OAuth client_id presenting the exchange request."),
		"target_audience":  stringSetDataSourceAttribute("Match conditions on the requested target audience."),
		"outbound_scopes": schema.ListAttribute{
			Description: "Outbound scopes to grant.",
			Computed:    true,
			ElementType: tftypes.StringType,
		},
	}
}

func stringSetDataSourceAttribute(description string) schema.Attribute {
	return schema.ListNestedAttribute{
		Description: description,
		Computed:    true,
		NestedObject: schema.NestedAttributeObject{
			Attributes: map[string]schema.Attribute{
				"exact": schema.StringAttribute{
					Description: "Exact string match.",
					Computed:    true,
				},
				"glob": schema.StringAttribute{
					Description: "Glob pattern match.",
					Computed:    true,
				},
			},
		},
	}
}

func (d *ExchangePolicyDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = dataSourceSchema()
}

func (d *ExchangePoliciesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = listDataSourceSchema()
}
