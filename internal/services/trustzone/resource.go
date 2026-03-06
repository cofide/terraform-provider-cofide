package trustzone

import (
	"context"
	"fmt"

	sdkclient "github.com/cofide/cofide-api-sdk/pkg/connect/client"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	_ resource.Resource                   = &TrustZoneResource{}
	_ resource.ResourceWithImportState    = &TrustZoneResource{}
	_ resource.ResourceWithValidateConfig = &TrustZoneResource{}
)

type TrustZoneResource struct {
	client sdkclient.ClientSet
}

func NewResource() resource.Resource {
	return &TrustZoneResource{}
}

func (t *TrustZoneResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_connect_trust_zone"
}

func (t *TrustZoneResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	t.client = client
}

func (t *TrustZoneResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan TrustZoneModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	trustZone, err := modelToProto(plan)
	if err != nil {
		resp.Diagnostics.AddError("Error converting model to proto", err.Error())
		return
	}

	createResp, err := t.client.TrustZoneV1Alpha1().CreateTrustZone(ctx, trustZone)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating trust zone",
			fmt.Sprintf("Could not create trust zone: %s", err.Error()),
		)
		return
	}

	newState := protoToModel(createResp)
	resp.Diagnostics.Append(resp.State.Set(ctx, &newState)...)
}

func (t *TrustZoneResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state TrustZoneModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	trustZoneID := state.ID.ValueString()
	if trustZoneID == "" {
		resp.Diagnostics.AddError(
			"Error reading trust zone",
			"Trust zone ID not found in state.",
		)
		return
	}

	getResp, err := t.client.TrustZoneV1Alpha1().GetTrustZone(ctx, trustZoneID)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error reading trust zone",
			err.Error(),
		)
		return
	}

	newState := protoToModel(getResp)
	resp.Diagnostics.Append(resp.State.Set(ctx, &newState)...)
}

func (t *TrustZoneResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var state TrustZoneModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var plan TrustZoneModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	trustZoneID := state.ID.ValueString()
	if trustZoneID == "" {
		resp.Diagnostics.AddError(
			"Error updating trust zone",
			"Trust zone ID not found in state. The resource might not have been created properly.",
		)
		return
	}

	trustZone, err := modelToProto(plan)
	if err != nil {
		resp.Diagnostics.AddError("Error converting model to proto", err.Error())
		return
	}
	trustZone.Id = &trustZoneID

	updateResp, err := t.client.TrustZoneV1Alpha1().UpdateTrustZone(ctx, trustZone)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating trust zone",
			fmt.Sprintf("Could not update trust zone: %s", err.Error()),
		)
		return
	}

	newState := protoToModel(updateResp)
	resp.Diagnostics.Append(resp.State.Set(ctx, &newState)...)
}

func (t *TrustZoneResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state TrustZoneModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := t.client.TrustZoneV1Alpha1().DestroyTrustZone(ctx, state.ID.ValueString())
	if err != nil {
		if status.Code(err) != codes.NotFound {
			resp.Diagnostics.AddError(
				"Error deleting trust zone",
				err.Error(),
			)

			return
		}
	}
}

func (t *TrustZoneResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (t *TrustZoneResource) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	var data TrustZoneModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if !data.IsManagementZone.IsNull() {
		resp.Diagnostics.AddWarning(
			"is_management_zone is immutable",
			"The is_management_zone field cannot be modified after creation. Create a new trust zone instead.",
		)
	}
}
