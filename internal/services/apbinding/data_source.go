package apbinding

import (
	"context"
	"fmt"

	apbindinginsvcpb "github.com/cofide/cofide-api-sdk/gen/go/proto/connect/ap_binding_service/v1alpha1"
	sdkclient "github.com/cofide/cofide-api-sdk/pkg/connect/client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

type APBindingDataSource struct {
	client sdkclient.ClientSet
}

var _ datasource.DataSourceWithConfigure = (*APBindingDataSource)(nil)

func NewDataSource() datasource.DataSource {
	return &APBindingDataSource{}
}

func (d *APBindingDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_connect_ap_binding"
}

func (d *APBindingDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *APBindingDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config APBindingModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	filter := &apbindinginsvcpb.ListAPBindingsRequest_Filter{
		OrgId:       config.OrgID.ValueStringPointer(),
		TrustZoneId: config.TrustZoneID.ValueStringPointer(),
		PolicyId:    config.PolicyID.ValueStringPointer(),
	}
	bindings, err := d.client.APBindingV1Alpha1().ListAPBindings(ctx, filter)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading attestation policy binding",
			fmt.Sprintf("Could not list attestation policy bindings: %s", err),
		)
		return
	}

	if len(bindings) == 0 {
		resp.Diagnostics.AddError(
			"Error reading attestation policy binding",
			"No matching attestation policy binding found",
		)
		return
	}

	if len(bindings) > 1 {
		resp.Diagnostics.AddError(
			"Error reading attestation policy binding",
			"Multiple attestation policy bindings found",
		)
		return
	}

	binding := bindings[0]

	if binding == nil {
		resp.Diagnostics.AddError(
			"Error reading attestation policy binding",
			"No matching attestation policy binding found",
		)
		return
	}

	state := protoToModel(binding)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
