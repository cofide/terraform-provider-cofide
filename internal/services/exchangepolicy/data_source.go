package exchangepolicy

import (
	"context"
	"fmt"

	sdkclient "github.com/cofide/cofide-api-sdk/pkg/connect/client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

type ExchangePolicyDataSource struct {
	client sdkclient.ClientSet
}

var _ datasource.DataSourceWithConfigure = (*ExchangePolicyDataSource)(nil)

func NewDataSource() datasource.DataSource {
	return &ExchangePolicyDataSource{}
}

func (d *ExchangePolicyDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_connect_exchange_policy"
}

func (d *ExchangePolicyDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *ExchangePolicyDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config ExchangePolicyModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	policy, err := d.client.ExchangePolicyV1Alpha1().GetExchangePolicy(ctx, config.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading exchange policy",
			fmt.Sprintf("Could not get exchange policy %q: %s", config.ID.ValueString(), err),
		)
		return
	}

	state, err := protoToModel(policy)
	if err != nil {
		resp.Diagnostics.AddError("Invalid exchange policy response", err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
