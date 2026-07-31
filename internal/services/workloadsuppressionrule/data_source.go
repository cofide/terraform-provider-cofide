package workloadsuppressionrule

import (
	"context"
	"fmt"

	sdkclient "github.com/cofide/cofide-api-sdk/pkg/connect/client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

type WorkloadSuppressionRuleDataSource struct {
	client sdkclient.ClientSet
}

var _ datasource.DataSourceWithConfigure = (*WorkloadSuppressionRuleDataSource)(nil)

func NewDataSource() datasource.DataSource {
	return &WorkloadSuppressionRuleDataSource{}
}

func (d *WorkloadSuppressionRuleDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_connect_workload_suppression_rule"
}

func (d *WorkloadSuppressionRuleDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *WorkloadSuppressionRuleDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config WorkloadSuppressionRuleModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ruleID := config.ID.ValueString()
	rule, err := d.client.WorkloadSuppressionRuleV1Alpha1().GetWorkloadSuppressionRule(ctx, ruleID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading workload suppression rule",
			fmt.Sprintf("Could not get workload suppression rule %q: %s", ruleID, err),
		)
		return
	}

	state := protoToModel(rule)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
