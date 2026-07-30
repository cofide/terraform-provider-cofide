package cloudaccount

import (
	"context"
	"fmt"

	cloudaccountsvcpb "github.com/cofide/cofide-api-sdk/gen/go/proto/connect/cloud_account_service/v1alpha1"
	sdkclient "github.com/cofide/cofide-api-sdk/pkg/connect/client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

type CloudAccountsDataSource struct {
	client sdkclient.ClientSet
}

func NewListDataSource() datasource.DataSource {
	return &CloudAccountsDataSource{}
}

func (d *CloudAccountsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_connect_cloud_accounts"
}

func (d *CloudAccountsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *CloudAccountsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config CloudAccountsDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	filter := &cloudaccountsvcpb.ListCloudAccountsRequest_Filter{
		IncludeSuppressed: config.IncludeSuppressed.ValueBool(),
	}

	if !config.OrgIDs.IsNull() {
		var orgIDs []string
		resp.Diagnostics.Append(config.OrgIDs.ElementsAs(ctx, &orgIDs, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
		filter.OrgIds = orgIDs
	}

	if !config.CloudOrganizationIDs.IsNull() {
		var cloudOrganizationIDs []string
		resp.Diagnostics.Append(config.CloudOrganizationIDs.ElementsAs(ctx, &cloudOrganizationIDs, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
		filter.CloudOrganizationIds = cloudOrganizationIDs
	}

	cloudAccounts, err := d.client.CloudAccountV1Alpha1().ListCloudAccounts(ctx, filter)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading cloud accounts",
			fmt.Sprintf("Could not list cloud accounts: %s", err),
		)
		return
	}

	state := CloudAccountsDataSourceModel{
		OrgIDs:               config.OrgIDs,
		CloudOrganizationIDs: config.CloudOrganizationIDs,
		IncludeSuppressed:    config.IncludeSuppressed,
	}

	for _, cloudAccount := range cloudAccounts {
		state.CloudAccounts = append(state.CloudAccounts, protoToModel(cloudAccount))
	}

	if state.CloudAccounts == nil {
		state.CloudAccounts = []CloudAccountModel{}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
