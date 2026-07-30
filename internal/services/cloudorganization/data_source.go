package cloudorganization

import (
	"context"
	"fmt"

	cloudorganizationpb "github.com/cofide/cofide-api-sdk/gen/go/proto/cloud_organization/v1alpha1"
	cloudorganizationsvcpb "github.com/cofide/cofide-api-sdk/gen/go/proto/connect/cloud_organization_service/v1alpha1"
	sdkclient "github.com/cofide/cofide-api-sdk/pkg/connect/client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

type CloudOrganizationDataSource struct {
	client sdkclient.ClientSet
}

func NewDataSource() datasource.DataSource {
	return &CloudOrganizationDataSource{}
}

func (d *CloudOrganizationDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_connect_cloud_organization"
}

func (d *CloudOrganizationDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *CloudOrganizationDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config CloudOrganizationModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	cloudOrganization, err := getCloudOrganizationByName(ctx, d.client, config.OrgID.ValueString(), config.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", err.Error())
		return
	}

	state := protoToModel(cloudOrganization)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// getCloudOrganizationByName lists cloud organizations, optionally narrowed by
// orgID, and returns the single organization matching name. The
// ListCloudOrganizations filter has no server-side name filter, so matching
// is done client-side.
func getCloudOrganizationByName(ctx context.Context, client sdkclient.ClientSet, orgID, name string) (*cloudorganizationpb.CloudOrganization, error) {
	filter := &cloudorganizationsvcpb.ListCloudOrganizationsRequest_Filter{}
	if orgID != "" {
		filter.OrgIds = []string{orgID}
	}

	cloudOrganizations, err := client.CloudOrganizationV1Alpha1().ListCloudOrganizations(ctx, filter)
	if err != nil {
		return nil, err
	}

	var matches []*cloudorganizationpb.CloudOrganization
	for _, cloudOrganization := range cloudOrganizations {
		if cloudOrganization.GetName() == name {
			matches = append(matches, cloudOrganization)
		}
	}

	if len(matches) == 0 {
		return nil, fmt.Errorf("cloud organization with name '%s' not found", name)
	}
	if len(matches) > 1 {
		return nil, fmt.Errorf("multiple cloud organizations with name '%s' found", name)
	}

	return matches[0], nil
}
