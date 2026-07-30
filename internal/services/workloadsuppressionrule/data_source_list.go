package workloadsuppressionrule

import (
	"context"
	"fmt"

	workloadsuppressionrulesvcpb "github.com/cofide/cofide-api-sdk/gen/go/proto/connect/workload_suppression_rule_service/v1alpha1"
	sdkclient "github.com/cofide/cofide-api-sdk/pkg/connect/client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

type WorkloadSuppressionRulesDataSource struct {
	client sdkclient.ClientSet
}

var _ datasource.DataSourceWithConfigure = (*WorkloadSuppressionRulesDataSource)(nil)

func NewListDataSource() datasource.DataSource {
	return &WorkloadSuppressionRulesDataSource{}
}

func (d *WorkloadSuppressionRulesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_connect_workload_suppression_rules"
}

func (d *WorkloadSuppressionRulesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(sdkclient.ClientSet)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected data source configure type",
			fmt.Sprintf("Expected sdkclient.ClientSet, got: %T", req.ProviderData),
		)
		return
	}

	d.client = client
}

func (d *WorkloadSuppressionRulesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config WorkloadSuppressionRulesDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	filter := &workloadsuppressionrulesvcpb.ListWorkloadSuppressionRulesRequest_Filter{}
	if !config.OrgIDs.IsNull() {
		var orgIDs []string
		resp.Diagnostics.Append(config.OrgIDs.ElementsAs(ctx, &orgIDs, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
		filter.OrgIds = orgIDs
	}
	if !config.Enabled.IsNull() {
		enabled := config.Enabled.ValueBool()
		filter.Enabled = &enabled
	}

	rules, err := d.client.WorkloadSuppressionRuleV1Alpha1().ListWorkloadSuppressionRules(ctx, filter)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading workload suppression rules",
			fmt.Sprintf("Could not list workload suppression rules: %s", err),
		)
		return
	}

	state := WorkloadSuppressionRulesDataSourceModel{
		OrgIDs:  config.OrgIDs,
		Enabled: config.Enabled,
	}

	for _, rule := range rules {
		state.WorkloadSuppressionRules = append(state.WorkloadSuppressionRules, protoToModel(rule))
	}

	if state.WorkloadSuppressionRules == nil {
		state.WorkloadSuppressionRules = []WorkloadSuppressionRuleModel{}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
