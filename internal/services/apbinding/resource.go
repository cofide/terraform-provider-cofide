package apbinding

import (
	"context"
	"fmt"

	apbindingpb "github.com/cofide/cofide-api-sdk/gen/go/proto/ap_binding/v1alpha1"
	apbindinginsvcpb "github.com/cofide/cofide-api-sdk/gen/go/proto/connect/ap_binding_service/v1alpha1"
	sdkclient "github.com/cofide/cofide-api-sdk/pkg/connect/client"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	tftypes "github.com/hashicorp/terraform-plugin-framework/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ resource.Resource = &APBindingResource{}
var _ resource.ResourceWithImportState = &APBindingResource{}
var _ resource.ResourceWithValidateConfig = &APBindingResource{}

type APBindingResource struct {
	client sdkclient.ClientSet
}

func NewResource() resource.Resource {
	return &APBindingResource{}
}

func (r *APBindingResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_connect_ap_binding"
}

func (r *APBindingResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(sdkclient.ClientSet)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected resource configure type",
			fmt.Sprintf("Expected sdkclient.ClientSet, got: %T", req.ProviderData),
		)

		return
	}

	r.client = client
}

func (a *APBindingResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan APBindingModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	federations := make([]*apbindingpb.APBindingFederation, 0)
	for _, federation := range plan.Federations {
		federations = append(federations, &apbindingpb.APBindingFederation{
			TrustZoneId: federation.TrustZoneID.ValueStringPointer(),
		})
	}

	binding := &apbindingpb.APBinding{
		TrustZoneId: plan.TrustZoneID.ValueStringPointer(),
		PolicyId:    plan.PolicyID.ValueStringPointer(),
		Federations: federations,
	}

	if !plan.OrgID.IsNull() {
		binding.OrgId = plan.OrgID.ValueStringPointer()
	}

	createResp, err := a.client.APBindingV1Alpha1().CreateAPBinding(ctx, binding)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating AP binding",
			err.Error(),
		)
		return
	}

	state := APBindingModel{
		ID:          tftypes.StringValue(createResp.GetId()),
		OrgID:       tftypes.StringValue(createResp.GetOrgId()),
		TrustZoneID: tftypes.StringValue(createResp.GetTrustZoneId()),
		PolicyID:    tftypes.StringValue(createResp.GetPolicyId()),
	}

	if createResp.GetFederations() != nil {
		respFederations := make([]APBindingFederationModel, 0, len(createResp.GetFederations()))
		for _, federation := range createResp.GetFederations() {
			respFederations = append(respFederations, APBindingFederationModel{
				TrustZoneID: tftypes.StringValue(federation.GetTrustZoneId()),
			})
		}
		state.Federations = respFederations
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (a *APBindingResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state APBindingModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	stateID := state.ID.ValueString()
	if stateID == "" {
		resp.State.RemoveResource(ctx)
		return
	}

	// The apbinding service does not have a Get method, so we list with a filter.
	filter := &apbindinginsvcpb.ListAPBindingsRequest_Filter{
		OrgId:       state.OrgID.ValueStringPointer(),
		TrustZoneId: state.TrustZoneID.ValueStringPointer(),
		PolicyId:    state.PolicyID.ValueStringPointer(),
	}
	bindings, err := a.client.APBindingV1Alpha1().ListAPBindings(ctx, filter)
	if err != nil {
		resp.Diagnostics.AddError("Error reading AP binding", fmt.Sprintf("Could not list AP bindings: %s", err))
		return
	}

	var foundBinding *apbindingpb.APBinding
	for _, binding := range bindings {
		if binding != nil && binding.GetId() == stateID {
			foundBinding = binding
			break
		}
	}

	if foundBinding == nil {
		// The resource no longer exists or doesn't match the filter.
		resp.State.RemoveResource(ctx)
		return
	}

	newState := APBindingModel{
		ID:          tftypes.StringValue(foundBinding.GetId()),
		OrgID:       tftypes.StringValue(foundBinding.GetOrgId()),
		TrustZoneID: tftypes.StringValue(foundBinding.GetTrustZoneId()),
		PolicyID:    tftypes.StringValue(foundBinding.GetPolicyId()),
	}

	if foundBinding.GetFederations() != nil {
		federations := make([]APBindingFederationModel, 0, len(foundBinding.GetFederations()))
		for _, federation := range foundBinding.GetFederations() {
			federations = append(federations, APBindingFederationModel{
				TrustZoneID: tftypes.StringValue(federation.GetTrustZoneId()),
			})
		}
		newState.Federations = federations
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &newState)...)
}

func (a *APBindingResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan APBindingModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state APBindingModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	bindingID := state.ID.ValueString()

	federations := make([]*apbindingpb.APBindingFederation, 0)
	for _, federation := range plan.Federations {
		federations = append(federations, &apbindingpb.APBindingFederation{
			TrustZoneId: federation.TrustZoneID.ValueStringPointer(),
		})
	}

	binding := &apbindingpb.APBinding{
		Id:          &bindingID,
		TrustZoneId: plan.TrustZoneID.ValueStringPointer(),
		PolicyId:    plan.PolicyID.ValueStringPointer(),
		Federations: federations,
	}

	if !plan.OrgID.IsNull() && plan.OrgID.ValueString() != "" {
		binding.OrgId = plan.OrgID.ValueStringPointer()
	}

	updateResp, err := a.client.APBindingV1Alpha1().UpdateAPBinding(ctx, binding)
	if err != nil {
		resp.Diagnostics.AddError("Error updating AP binding", err.Error())
		return
	}

	var orgIDStr tftypes.String
	if orgID := updateResp.GetOrgId(); orgID != "" {
		orgIDStr = tftypes.StringValue(orgID)
	} else if !plan.OrgID.IsNull() {
		orgIDStr = plan.OrgID
	} else {
		orgIDStr = tftypes.StringNull()
	}

	newState := APBindingModel{
		ID:          tftypes.StringValue(updateResp.GetId()),
		OrgID:       orgIDStr,
		TrustZoneID: tftypes.StringValue(updateResp.GetTrustZoneId()),
		PolicyID:    tftypes.StringValue(updateResp.GetPolicyId()),
	}

	if updateResp.GetFederations() != nil {
		respFederations := make([]APBindingFederationModel, 0, len(updateResp.GetFederations()))
		for _, federation := range updateResp.GetFederations() {
			respFederations = append(respFederations, APBindingFederationModel{
				TrustZoneID: tftypes.StringValue(federation.GetTrustZoneId()),
			})
		}
		newState.Federations = respFederations
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &newState)...)
}

func (a *APBindingResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state APBindingModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := a.client.APBindingV1Alpha1().DestroyAPBinding(ctx, state.ID.ValueString())
	if err != nil {
		if status.Code(err) != codes.NotFound {
			resp.Diagnostics.AddError(
				"Error deleting AP binding",
				err.Error(),
			)

			return
		}
	}
}

func (a *APBindingResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (a *APBindingResource) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	var data APBindingModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
}
