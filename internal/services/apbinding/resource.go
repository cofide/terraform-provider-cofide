package apbinding

import (
	"context"
	"fmt"

	apbindingpb "github.com/cofide/cofide-api-sdk/gen/go/proto/ap_binding/v1alpha1"
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
		OrgId:       plan.OrgID.ValueStringPointer(),
		TrustZoneId: plan.TrustZoneID.ValueStringPointer(),
		PolicyId:    plan.PolicyID.ValueStringPointer(),
		Federations: federations,
	}

	createResp, err := a.client.APBindingV1Alpha1().CreateAPBinding(ctx, binding)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating AP binding",
			err.Error(),
		)
		return
	}

	if createResp.GetId() != "" {
		plan.ID = tftypes.StringValue(createResp.GetId())
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (a *APBindingResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state APBindingModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (a *APBindingResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var state APBindingModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var plan APBindingModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	bindingID := state.ID.ValueString()
	if bindingID == "" {
		resp.Diagnostics.AddError(
			"Error updating AP binding",
			"AP binding ID not found in state. The resource might not have been created properly.",
		)
		return
	}

	federations := make([]*apbindingpb.APBindingFederation, 0)
	for _, federation := range plan.Federations {
		federations = append(federations, &apbindingpb.APBindingFederation{
			TrustZoneId: federation.TrustZoneID.ValueStringPointer(),
		})
	}

	binding := &apbindingpb.APBinding{
		Id:          &bindingID,
		OrgId:       plan.OrgID.ValueStringPointer(),
		TrustZoneId: plan.TrustZoneID.ValueStringPointer(),
		PolicyId:    plan.PolicyID.ValueStringPointer(),
		Federations: federations,
	}

	updateResp, err := a.client.APBindingV1Alpha1().UpdateAPBinding(ctx, binding)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating AP binding",
			err.Error(),
		)
		return
	}

	if updateResp.GetId() != "" {
		plan.ID = tftypes.StringValue(updateResp.GetId())
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
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
