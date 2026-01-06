package trustzone

import (
	"context"
	"fmt"

	trustzonesvcpb "github.com/cofide/cofide-api-sdk/gen/go/proto/connect/trust_zone_service/v1alpha1"
	sdkclient "github.com/cofide/cofide-api-sdk/pkg/connect/client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type TrustZoneDataSource struct {
	client sdkclient.ClientSet
}

var _ datasource.DataSourceWithConfigure = (*TrustZoneDataSource)(nil)

func NewDataSource() datasource.DataSource {
	return &TrustZoneDataSource{}
}

func (d *TrustZoneDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_connect_trust_zone"
}

func (t *TrustZoneDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	t.client = client
}

func (t *TrustZoneDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config TrustZoneModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	filter := &trustzonesvcpb.ListTrustZonesRequest_Filter{
		Name:        config.Name.ValueStringPointer(),
		OrgId:       config.OrgID.ValueStringPointer(),
		TrustDomain: config.TrustDomain.ValueStringPointer(),
	}
	trustZones, err := t.client.TrustZoneV1Alpha1().ListTrustZones(ctx, filter)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading trust zone",
			fmt.Sprintf("Could not list trust zones: %s", err),
		)
		return
	}

	if len(trustZones) == 0 {
		resp.Diagnostics.AddError(
			"Error reading trust zone",
			"No matching trust zone found",
		)
		return
	}

	if len(trustZones) > 1 {
		resp.Diagnostics.AddError(
			"Error reading trust zone",
			"Multiple trust zones found",
		)
		return
	}

	trustZone := trustZones[0]

	state := TrustZoneModel{
		ID:                    types.StringValue(trustZone.GetId()),
		Name:                  types.StringValue(trustZone.GetName()),
		TrustDomain:           types.StringValue(trustZone.GetTrustDomain()),
		OrgID:                 types.StringValue(trustZone.GetOrgId()),
		IsManagementZone:      types.BoolValue(trustZone.GetIsManagementZone()),
		BundleEndpointURL:     types.StringValue(trustZone.GetBundleEndpointUrl()),
		BundleEndpointProfile: types.StringValue(trustZone.GetBundleEndpointProfile().String()),
		JWTIssuer:             types.StringValue(trustZone.GetJwtIssuer()),
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
