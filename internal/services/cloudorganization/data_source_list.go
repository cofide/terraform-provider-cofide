package cloudorganization

import (
	"context"
	"fmt"

	cloudorganizationsvcpb "github.com/cofide/cofide-api-sdk/gen/go/proto/connect/cloud_organization_service/v1alpha1"
	sdkclient "github.com/cofide/cofide-api-sdk/pkg/connect/client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

type CloudOrganizationsDataSource struct {
	client sdkclient.ClientSet
}

func NewListDataSource() datasource.DataSource {
	return &CloudOrganizationsDataSource{}
}

func (d *CloudOrganizationsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_connect_cloud_organizations"
}

func (d *CloudOrganizationsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *CloudOrganizationsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config CloudOrganizationsDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	filter := &cloudorganizationsvcpb.ListCloudOrganizationsRequest_Filter{}
	if !config.OrgIDs.IsNull() {
		var orgIDs []string
		resp.Diagnostics.Append(config.OrgIDs.ElementsAs(ctx, &orgIDs, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
		filter.OrgIds = orgIDs
	}

	cloudOrganizations, err := d.client.CloudOrganizationV1Alpha1().ListCloudOrganizations(ctx, filter)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading cloud organizations",
			fmt.Sprintf("Could not list cloud organizations: %s", err),
		)
		return
	}

	state := CloudOrganizationsDataSourceModel{
		OrgIDs: config.OrgIDs,
	}

	for _, cloudOrganization := range cloudOrganizations {
		state.CloudOrganizations = append(state.CloudOrganizations, protoToModel(cloudOrganization))
	}

	if state.CloudOrganizations == nil {
		state.CloudOrganizations = []CloudOrganizationModel{}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
