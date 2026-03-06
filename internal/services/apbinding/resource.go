package apbinding

import (
	"context"
	"fmt"

	sdkclient "github.com/cofide/cofide-api-sdk/pkg/connect/client"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
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

func (r *APBindingResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan APBindingModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	binding := modelToProto(plan)

	createResp, err := r.client.APBindingV1Alpha1().CreateAPBinding(ctx, binding)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating AP binding",
			err.Error(),
		)
		return
	}

	newState := protoToModel(createResp)
	resp.Diagnostics.Append(resp.State.Set(ctx, &newState)...)
}

func (r *APBindingResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state APBindingModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	bindingID := state.ID.ValueString()
	if bindingID == "" {
		resp.Diagnostics.AddError(
			"Error reading AP binding",
			"AP binding ID not found in state.",
		)
		return
	}

	getResp, err := r.client.APBindingV1Alpha1().GetAPBinding(ctx, bindingID)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error reading AP binding",
			err.Error(),
		)
		return
	}

	newState := protoToModel(getResp)
	resp.Diagnostics.Append(resp.State.Set(ctx, &newState)...)
}

func (r *APBindingResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
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

	binding := modelToProto(plan)
	binding.Id = &bindingID

	updateResp, err := r.client.APBindingV1Alpha1().UpdateAPBinding(ctx, binding)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating AP binding",
			err.Error(),
		)
		return
	}

	newState := protoToModel(updateResp)
	resp.Diagnostics.Append(resp.State.Set(ctx, &newState)...)
}

func (r *APBindingResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state APBindingModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.APBindingV1Alpha1().DestroyAPBinding(ctx, state.ID.ValueString())
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

func (r *APBindingResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *APBindingResource) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	var data APBindingModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
}
