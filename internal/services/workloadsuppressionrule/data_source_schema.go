package workloadsuppressionrule

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	tftypes "github.com/hashicorp/terraform-plugin-framework/types"
)

var _ datasource.DataSource = &WorkloadSuppressionRuleDataSource{}
var _ datasource.DataSource = &WorkloadSuppressionRulesDataSource{}

// kubernetesPodMatcherDataSourceAttributes returns the Computed attributes shared
// by the kubernetes_pod matcher block in both the singular and list data sources.
func kubernetesPodMatcherDataSourceAttributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"trust_zone_ids": schema.ListAttribute{
			Description: "Matches only pods within these trust zones.",
			Computed:    true,
			ElementType: tftypes.StringType,
		},
		"cluster_ids": schema.ListAttribute{
			Description: "Matches only pods within these clusters.",
			Computed:    true,
			ElementType: tftypes.StringType,
		},
		"namespaces": schema.ListAttribute{
			Description: "Matches pods in these Kubernetes namespaces.",
			Computed:    true,
			ElementType: tftypes.StringType,
		},
		"labels": schema.MapAttribute{
			Description: "Matches pods with these Kubernetes labels.",
			Computed:    true,
			ElementType: tftypes.StringType,
		},
	}
}

func DataSourceSchema(_ context.Context) schema.Schema {
	return schema.Schema{
		MarkdownDescription: "Provides information about a Cofide Connect workload suppression rule.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The ID of the workload suppression rule.",
				Required:    true,
			},
			"org_id": schema.StringAttribute{
				Description: "The ID of the organization.",
				Computed:    true,
			},
			"name": schema.StringAttribute{
				Description: "The name of the workload suppression rule.",
				Computed:    true,
			},
			"description": schema.StringAttribute{
				Description: "A human-readable explanation of why the rule exists.",
				Computed:    true,
			},
			"enabled": schema.BoolAttribute{
				Description: "Whether the rule is currently active.",
				Computed:    true,
			},
			"kubernetes_pod": schema.SingleNestedAttribute{
				Description: "Matches Kubernetes pod workloads.",
				Computed:    true,
				Attributes:  kubernetesPodMatcherDataSourceAttributes(),
			},
			"created_at": schema.StringAttribute{
				Description: "The time the rule was created, in RFC3339 format.",
				Computed:    true,
			},
			"last_updated_at": schema.StringAttribute{
				Description: "The time the rule was last updated, in RFC3339 format.",
				Computed:    true,
			},
		},
	}
}

func (d *WorkloadSuppressionRuleDataSource) Schema(ctx context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = DataSourceSchema(ctx)
}

func ListDataSourceSchema(_ context.Context) schema.Schema {
	return schema.Schema{
		MarkdownDescription: "Provides information about Cofide Connect workload suppression rules.",
		Attributes: map[string]schema.Attribute{
			"org_ids": schema.ListAttribute{
				Description: "Filter by organization IDs.",
				Optional:    true,
				ElementType: tftypes.StringType,
			},
			"enabled": schema.BoolAttribute{
				Description: "Filter by enabled state. When unset, rules are returned regardless of enabled state.",
				Optional:    true,
			},
			"workload_suppression_rules": schema.ListNestedAttribute{
				Description: "The list of workload suppression rules.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Description: "The ID of the workload suppression rule.",
							Computed:    true,
						},
						"org_id": schema.StringAttribute{
							Description: "The ID of the organization.",
							Computed:    true,
						},
						"name": schema.StringAttribute{
							Description: "The name of the workload suppression rule.",
							Computed:    true,
						},
						"description": schema.StringAttribute{
							Description: "A human-readable explanation of why the rule exists.",
							Computed:    true,
						},
						"enabled": schema.BoolAttribute{
							Description: "Whether the rule is currently active.",
							Computed:    true,
						},
						"kubernetes_pod": schema.SingleNestedAttribute{
							Description: "Matches Kubernetes pod workloads.",
							Computed:    true,
							Attributes:  kubernetesPodMatcherDataSourceAttributes(),
						},
						"created_at": schema.StringAttribute{
							Description: "The time the rule was created, in RFC3339 format.",
							Computed:    true,
						},
						"last_updated_at": schema.StringAttribute{
							Description: "The time the rule was last updated, in RFC3339 format.",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func (d *WorkloadSuppressionRulesDataSource) Schema(ctx context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = ListDataSourceSchema(ctx)
}

func (d *WorkloadSuppressionRulesDataSource) ConfigValidators(_ context.Context) []datasource.ConfigValidator {
	return []datasource.ConfigValidator{}
}
