package cloudaccount

import (
	"context"
	"fmt"

	cloudaccountpb "github.com/cofide/cofide-api-sdk/gen/go/proto/cloud_account/v1alpha1"
	cloudaccountsvcpb "github.com/cofide/cofide-api-sdk/gen/go/proto/connect/cloud_account_service/v1alpha1"
	sdkclient "github.com/cofide/cofide-api-sdk/pkg/connect/client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

type CloudAccountDataSource struct {
	client sdkclient.ClientSet
}

func NewDataSource() datasource.DataSource {
	return &CloudAccountDataSource{}
}

func (d *CloudAccountDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_connect_cloud_account"
}

func (d *CloudAccountDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *CloudAccountDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config CloudAccountModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	cloudAccount, err := getCloudAccountByName(ctx, d.client, config.OrgID.ValueString(), config.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", err.Error())
		return
	}

	state := protoToModel(cloudAccount)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// getCloudAccountByName lists cloud accounts, optionally narrowed by orgID,
// and returns the single account matching name. The ListCloudAccounts filter
// has no server-side name filter, so matching is done client-side.
func getCloudAccountByName(ctx context.Context, client sdkclient.ClientSet, orgID, name string) (*cloudaccountpb.CloudAccount, error) {
	filter := &cloudaccountsvcpb.ListCloudAccountsRequest_Filter{
		IncludeSuppressed: true,
	}
	if orgID != "" {
		filter.OrgIds = []string{orgID}
	}

	cloudAccounts, err := client.CloudAccountV1Alpha1().ListCloudAccounts(ctx, filter)
	if err != nil {
		return nil, err
	}

	var matches []*cloudaccountpb.CloudAccount
	for _, cloudAccount := range cloudAccounts {
		if cloudAccount.GetName() == name {
			matches = append(matches, cloudAccount)
		}
	}

	if len(matches) == 0 {
		return nil, fmt.Errorf("cloud account with name '%s' not found", name)
	}
	if len(matches) > 1 {
		return nil, fmt.Errorf("multiple cloud accounts with name '%s' found", name)
	}

	return matches[0], nil
}
