package trustzoneserver

import (
	"context"
	"fmt"

	trustzoneserversvcpb "github.com/cofide/cofide-api-sdk/gen/go/proto/connect/trust_zone_server_service/v1alpha1"
	sdkclient "github.com/cofide/cofide-api-sdk/pkg/connect/client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

type TrustZoneServersDataSource struct {
	client sdkclient.ClientSet
}

var _ datasource.DataSourceWithConfigure = (*TrustZoneServersDataSource)(nil)

func NewListDataSource() datasource.DataSource {
	return &TrustZoneServersDataSource{}
}

func (d *TrustZoneServersDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_connect_trust_zone_servers"
}

func (d *TrustZoneServersDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *TrustZoneServersDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config TrustZoneServersDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	filter := &trustzoneserversvcpb.ListTrustZoneServersRequest_Filter{}
	if !config.TrustZoneID.IsNull() {
		filter.TrustZoneId = config.TrustZoneID.ValueString()
	}
	if !config.ClusterID.IsNull() {
		filter.ClusterId = config.ClusterID.ValueString()
	}
	if !config.OrgID.IsNull() {
		filter.OrgId = config.OrgID.ValueString()
	}

	servers, err := d.client.TrustZoneServerV1Alpha1().ListTrustZoneServers(ctx, filter)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading trust zone servers",
			fmt.Sprintf("Could not list trust zone servers: %s", err),
		)
		return
	}

	state := TrustZoneServersDataSourceModel{
		TrustZoneID: config.TrustZoneID,
		ClusterID:   config.ClusterID,
		OrgID:       config.OrgID,
	}

	for _, server := range servers {
		helmValues := helmValuesFromProto(server.GetHelmValues())
		state.TrustZoneServers = append(state.TrustZoneServers, trustZoneServerFromProto(server, helmValues))
	}

	if state.TrustZoneServers == nil {
		state.TrustZoneServers = []TrustZoneServerModel{}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
