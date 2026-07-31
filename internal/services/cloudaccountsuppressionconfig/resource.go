package cloudaccountsuppressionconfig

import (
	"context"
	"fmt"

	cloudaccountpb "github.com/cofide/cofide-api-sdk/gen/go/proto/cloud_account/v1alpha1"
	cloudaccountsvcpb "github.com/cofide/cofide-api-sdk/gen/go/proto/connect/cloud_account_service/v1alpha1"
	sdkclient "github.com/cofide/cofide-api-sdk/pkg/connect/client"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	tftypes "github.com/hashicorp/terraform-plugin-framework/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ resource.Resource = &Resource{}
var _ resource.ResourceWithImportState = &Resource{}

type Resource struct {
	client sdkclient.ClientSet
}

func NewResource() resource.Resource {
	return &Resource{}
}

func (r *Resource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_connect_cloud_account_suppression_config"
}

func (r *Resource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *Resource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan Model

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	newState, err := r.setSuppressed(ctx, plan.CloudAccountID.ValueString(), plan.Suppressed.ValueBool())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error setting cloud account suppression config",
			fmt.Sprintf("Could not update cloud account %q: %s", plan.CloudAccountID.ValueString(), err),
		)

		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &newState)...)
}

func (r *Resource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan Model

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	newState, err := r.setSuppressed(ctx, plan.CloudAccountID.ValueString(), plan.Suppressed.ValueBool())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error setting cloud account suppression config",
			fmt.Sprintf("Could not update cloud account %q: %s", plan.CloudAccountID.ValueString(), err),
		)

		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &newState)...)
}

// setSuppressed sets the cloud account's suppressed flag, using an update
// mask so no other field on the cloud account is touched.
func (r *Resource) setSuppressed(ctx context.Context, cloudAccountID string, suppressed bool) (Model, error) {
	cloudAccount := &cloudaccountpb.CloudAccount{
		Id:         cloudAccountID,
		Suppressed: suppressed,
	}
	updateMask := &cloudaccountsvcpb.UpdateCloudAccountRequest_UpdateMask{Suppressed: true}

	updateResp, err := r.client.CloudAccountV1Alpha1().UpdateCloudAccount(ctx, cloudAccount, updateMask)
	if err != nil {
		return Model{}, err
	}

	return Model{
		ID:             tftypes.StringValue(cloudAccountID),
		CloudAccountID: tftypes.StringValue(cloudAccountID),
		Suppressed:     tftypes.BoolValue(updateResp.GetSuppressed()),
	}, nil
}

func (r *Resource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state Model

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	cloudAccountID := state.ID.ValueString()
	cloudAccount, err := r.client.CloudAccountV1Alpha1().GetCloudAccount(ctx, cloudAccountID)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error reading cloud account suppression config",
			fmt.Sprintf("Could not read cloud account %q: %s", cloudAccountID, err),
		)

		return
	}

	newState := Model{
		ID:             tftypes.StringValue(cloudAccountID),
		CloudAccountID: tftypes.StringValue(cloudAccountID),
		Suppressed:     tftypes.BoolValue(cloudAccount.GetSuppressed()),
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &newState)...)
}

func (r *Resource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state Model

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if _, err := r.setSuppressed(ctx, state.CloudAccountID.ValueString(), false); err != nil && status.Code(err) != codes.NotFound {
		resp.Diagnostics.AddError(
			"Error deleting cloud account suppression config",
			err.Error(),
		)

		return
	}
}

func (r *Resource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
