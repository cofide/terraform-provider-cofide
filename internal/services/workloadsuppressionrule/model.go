package workloadsuppressionrule

import tftypes "github.com/hashicorp/terraform-plugin-framework/types"

type WorkloadSuppressionRuleModel struct {
	ID            tftypes.String             `tfsdk:"id"`
	OrgID         tftypes.String             `tfsdk:"org_id"`
	Name          tftypes.String             `tfsdk:"name"`
	Description   tftypes.String             `tfsdk:"description"`
	Enabled       tftypes.Bool               `tfsdk:"enabled"`
	KubernetesPod *KubernetesPodMatcherModel `tfsdk:"kubernetes_pod"`
	CreatedAt     tftypes.String             `tfsdk:"created_at"`
	LastUpdatedAt tftypes.String             `tfsdk:"last_updated_at"`
}

type KubernetesPodMatcherModel struct {
	TrustZoneIDs tftypes.List `tfsdk:"trust_zone_ids"`
	ClusterIDs   tftypes.List `tfsdk:"cluster_ids"`
	Namespaces   tftypes.List `tfsdk:"namespaces"`
	Labels       tftypes.Map  `tfsdk:"labels"`
}

type WorkloadSuppressionRulesDataSourceModel struct {
	OrgIDs                   tftypes.List                   `tfsdk:"org_ids"`
	Enabled                  tftypes.Bool                   `tfsdk:"enabled"`
	WorkloadSuppressionRules []WorkloadSuppressionRuleModel `tfsdk:"workload_suppression_rules"`
}
