package trustzoneserver

import (
	"context"
	"fmt"

	sdkclient "github.com/cofide/cofide-api-sdk/pkg/connect/client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

type TrustZoneServerDataSource struct {
	client sdkclient.ClientSet
}

var _ datasource.DataSourceWithConfigure = (*TrustZoneServerDataSource)(nil)

func NewDataSource() datasource.DataSource {
	return &TrustZoneServerDataSource{}
}

func (d *TrustZoneServerDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_connect_trust_zone_server"
}

func (d *TrustZoneServerDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *TrustZoneServerDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config TrustZoneServerModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	server, err := d.client.TrustZoneServerV1Alpha1().GetTrustZoneServer(ctx, config.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading trust zone server",
			fmt.Sprintf("Could not get trust zone server %q: %s", config.ID.ValueString(), err),
		)
		return
	}

	helmValues := helmValuesFromProto(server.GetHelmValues())
	state, diags := trustZoneServerFromProto(server, helmValues)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
