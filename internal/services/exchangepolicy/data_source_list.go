package exchangepolicy

import (
	"context"
	"fmt"

	exchangepolicysvcpb "github.com/cofide/cofide-api-sdk/gen/go/proto/connect/exchange_policy_service/v1alpha1"
	sdkclient "github.com/cofide/cofide-api-sdk/pkg/connect/client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

type ExchangePoliciesDataSource struct {
	client sdkclient.ClientSet
}

var _ datasource.DataSourceWithConfigure = (*ExchangePoliciesDataSource)(nil)

func NewListDataSource() datasource.DataSource {
	return &ExchangePoliciesDataSource{}
}

func (d *ExchangePoliciesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_connect_exchange_policies"
}

func (d *ExchangePoliciesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *ExchangePoliciesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config ExchangePoliciesDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	filter := &exchangepolicysvcpb.ListExchangePoliciesRequest_Filter{}
	if !config.TrustZoneID.IsNull() {
		filter.TrustZoneId = config.TrustZoneID.ValueString()
	}
	if !config.OrgID.IsNull() {
		filter.OrgId = config.OrgID.ValueString()
	}
	if !config.Name.IsNull() {
		filter.Name = config.Name.ValueString()
	}

	policies, err := d.client.ExchangePolicyV1Alpha1().ListExchangePolicies(ctx, filter)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading exchange policies",
			fmt.Sprintf("Could not list exchange policies: %s", err),
		)
		return
	}

	state := ExchangePoliciesDataSourceModel{
		TrustZoneID: config.TrustZoneID,
		OrgID:       config.OrgID,
		Name:        config.Name,
	}

	for _, policy := range policies {
		m, err := protoToModel(policy)
		if err != nil {
			resp.Diagnostics.AddError("Invalid exchange policy response", err.Error())
			return
		}
		state.ExchangePolicies = append(state.ExchangePolicies, m)
	}

	if state.ExchangePolicies == nil {
		state.ExchangePolicies = []ExchangePolicyModel{}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
