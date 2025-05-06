package federation

import (
	"context"
	"fmt"

	federationsvcpb "github.com/cofide/cofide-api-sdk/gen/go/proto/connect/federation_service/v1alpha1"

	sdkclient "github.com/cofide/cofide-api-sdk/pkg/connect/client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type FederationDataSource struct {
	client sdkclient.ClientSet
}

var _ datasource.DataSourceWithConfigure = (*FederationDataSource)(nil)

func NewDataSource() datasource.DataSource {
	return &FederationDataSource{}
}

func (f *FederationDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_connect_federation"
}

func (f *FederationDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(sdkclient.ClientSet)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected sdkclient.ClientSet, got: %T", req.ProviderData),
		)
		return
	}

	f.client = client
}

func (f *FederationDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config FederationModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	filter := &federationsvcpb.ListFederationsRequest_Filter{
		OrgId:             config.OrgID.ValueStringPointer(),
		TrustZoneId:       config.TrustZoneID.ValueStringPointer(),
		RemoteTrustZoneId: config.RemoteTrustZoneID.ValueStringPointer(),
	}

	federations, err := f.client.FederationV1Alpha1().ListFederations(ctx, filter)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Federation",
			fmt.Sprintf("Could not list federations: %s", err),
		)

		return
	}

	if len(federations) == 0 {
		resp.Diagnostics.AddError(
			"Error Reading Federation",
			"No matching federation found",
		)

		return
	}

	if len(federations) > 1 {
		resp.Diagnostics.AddError(
			"Error Reading Federation",
			"Multiple federations found",
		)

		return
	}

	federation := federations[0]

	state := FederationModel{
		ID:                types.StringValue(federation.GetId()),
		OrgID:             types.StringValue(federation.GetOrgId()),
		TrustZoneID:       types.StringValue(federation.GetTrustZoneId()),
		RemoteTrustZoneID: types.StringValue(federation.GetRemoteTrustZoneId()),
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
