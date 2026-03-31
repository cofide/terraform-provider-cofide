package exchangepolicy

import (
	"context"
	"fmt"

	exchangepolicysvcpb "github.com/cofide/cofide-api-sdk/gen/go/proto/connect/exchange_policy_service/v1alpha1"
	sdkclient "github.com/cofide/cofide-api-sdk/pkg/connect/client"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	_ resource.Resource                = &ExchangePolicyResource{}
	_ resource.ResourceWithImportState = &ExchangePolicyResource{}
)

type ExchangePolicyResource struct {
	client sdkclient.ClientSet
}

func NewResource() resource.Resource {
	return &ExchangePolicyResource{}
}

func (r *ExchangePolicyResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_connect_exchange_policy"
}

func (r *ExchangePolicyResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *ExchangePolicyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan ExchangePolicyModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	policy, err := modelToProto(ctx, plan)
	if err != nil {
		resp.Diagnostics.AddError("Invalid exchange policy", err.Error())
		return
	}
	createResp, err := r.client.ExchangePolicyV1Alpha1().CreateExchangePolicy(ctx, policy)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating exchange policy",
			err.Error(),
		)
		return
	}

	state := protoToModel(createResp)
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *ExchangePolicyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state ExchangePolicyModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := state.ID.ValueString()
	if id == "" {
		resp.Diagnostics.AddError(
			"Error reading exchange policy",
			"Exchange policy ID not found in state.",
		)
		return
	}

	getResp, err := r.client.ExchangePolicyV1Alpha1().GetExchangePolicy(ctx, id)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error reading exchange policy",
			err.Error(),
		)
		return
	}

	newState := protoToModel(getResp)
	resp.Diagnostics.Append(resp.State.Set(ctx, newState)...)
}

func (r *ExchangePolicyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan ExchangePolicyModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state ExchangePolicyModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	policy, err := modelToProto(ctx, plan)
	if err != nil {
		resp.Diagnostics.AddError("Invalid exchange policy", err.Error())
		return
	}
	policy.Id = state.ID.ValueString()

	updateMask := &exchangepolicysvcpb.UpdateExchangePolicyRequest_UpdateMask{
		Name:            true,
		Action:          true,
		SubjectIdentity: true,
		SubjectIssuer:   true,
		ActorIdentity:   true,
		ActorIssuer:     true,
		ClientId:        true,
		TargetAudience:  true,
		OutboundScopes:  true,
	}

	updateResp, err := r.client.ExchangePolicyV1Alpha1().UpdateExchangePolicy(ctx, policy, updateMask)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating exchange policy",
			err.Error(),
		)
		return
	}

	newState := protoToModel(updateResp)
	resp.Diagnostics.Append(resp.State.Set(ctx, newState)...)
}

func (r *ExchangePolicyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state ExchangePolicyModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.ExchangePolicyV1Alpha1().DestroyExchangePolicy(ctx, state.ID.ValueString())
	if err != nil {
		if status.Code(err) != codes.NotFound {
			resp.Diagnostics.AddError(
				"Error deleting exchange policy",
				err.Error(),
			)
			return
		}
	}
}

func (r *ExchangePolicyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
