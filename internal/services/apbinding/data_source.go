package apbinding

import (
	"context"
	"fmt"

	apbindinginsvcpb "github.com/cofide/cofide-api-sdk/gen/go/proto/connect/ap_binding_service/v1alpha1"
	sdkclient "github.com/cofide/cofide-api-sdk/pkg/connect/client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type APBindingDataSource struct {
	client sdkclient.ClientSet
}

var _ datasource.DataSourceWithConfigure = (*APBindingDataSource)(nil)

func NewDataSource() datasource.DataSource {
	return &APBindingDataSource{}
}

func (a *APBindingDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_connect_ap_binding"
}

func (a *APBindingDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	a.client = client
}

func (a *APBindingDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
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
	bindings, err := a.client.APBindingV1Alpha1().ListAPBindings(ctx, filter)
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

	state := APBindingModel{
		ID:          types.StringValue(binding.GetId()),
		OrgID:       types.StringValue(binding.GetOrgId()),
		TrustZoneID: types.StringValue(binding.GetTrustZoneId()),
		PolicyID:    types.StringValue(binding.GetPolicyId()),
	}

	federations := make([]APBindingFederationModel, 0)
	for _, federation := range binding.GetFederations() {
		federations = append(federations, APBindingFederationModel{
			TrustZoneID: types.StringValue(federation.GetTrustZoneId()),
		})
	}
	state.Federations = federations

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
