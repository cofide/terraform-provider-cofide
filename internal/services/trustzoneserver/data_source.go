package trustzoneserver

import (
	"context"
	"fmt"

	trustzoneserversvcpb "github.com/cofide/cofide-api-sdk/gen/go/proto/connect/trust_zone_server_service/v1alpha1"
	trustzoneserverpb "github.com/cofide/cofide-api-sdk/gen/go/proto/trust_zone_server/v1alpha1"
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

	var server *trustzoneserverpb.TrustZoneServer

	if !config.ID.IsNull() && config.ID.ValueString() != "" {
		s, err := d.client.TrustZoneServerV1Alpha1().GetTrustZoneServer(ctx, config.ID.ValueString())
		if err != nil {
			resp.Diagnostics.AddError(
				"Error reading trust zone server",
				fmt.Sprintf("Could not get trust zone server %q: %s", config.ID.ValueString(), err),
			)
			return
		}
		server = s
	} else {
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
				"Error reading trust zone server",
				fmt.Sprintf("Could not list trust zone servers: %s", err),
			)
			return
		}

		if len(servers) == 0 {
			resp.Diagnostics.AddError(
				"Error reading trust zone server",
				"No matching trust zone server found",
			)
			return
		}

		if len(servers) > 1 {
			resp.Diagnostics.AddError(
				"Error reading trust zone server",
				"Multiple trust zone servers found",
			)
			return
		}

		server = servers[0]
	}

	helmValues := helmValuesFromProto(server.GetHelmValues())
	state := trustZoneServerFromProto(server, helmValues)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
