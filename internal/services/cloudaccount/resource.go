package cloudaccount

import (
	"context"
	"fmt"

	cloudaccountsvcpb "github.com/cofide/cofide-api-sdk/gen/go/proto/connect/cloud_account_service/v1alpha1"
	sdkclient "github.com/cofide/cofide-api-sdk/pkg/connect/client"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ resource.Resource = &CloudAccountResource{}
var _ resource.ResourceWithImportState = &CloudAccountResource{}

type CloudAccountResource struct {
	client sdkclient.ClientSet
}

func NewResource() resource.Resource {
	return &CloudAccountResource{}
}

func (c *CloudAccountResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_connect_cloud_account"
}

func (c *CloudAccountResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	c.client = client
}

func (c *CloudAccountResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan CloudAccountModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	cloudAccount, diags := modelToProto(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	createResp, err := c.client.CloudAccountV1Alpha1().CreateCloudAccount(ctx, cloudAccount)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating cloud account",
			fmt.Sprintf("Could not create cloud account: %s", err),
		)

		return
	}

	state := protoToModel(createResp)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (c *CloudAccountResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state CloudAccountModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	cloudAccountID := state.ID.ValueString()
	cloudAccount, err := c.client.CloudAccountV1Alpha1().GetCloudAccount(ctx, cloudAccountID)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error reading cloud account",
			fmt.Sprintf("Could not read cloud account %q: %s", cloudAccountID, err),
		)

		return
	}

	newState := protoToModel(cloudAccount)
	resp.Diagnostics.Append(resp.State.Set(ctx, &newState)...)
}

func (c *CloudAccountResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan CloudAccountModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var config CloudAccountModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state CloudAccountModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	cloudAccount, diags := modelToProto(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	cloudAccount.Id = state.ID.ValueString()

	// Only fields explicitly set in the practitioner's config are included in
	// the update mask, so a field left unset here (e.g. because it's managed
	// by a standalone aws_lambda_discovery_config/aws_agent_core_discovery_config/
	// cloud_account_suppression_config resource, or by discovery) is never
	// touched by this apply.
	updateMask := updateMaskForConfig(config)

	updateResp, err := c.client.CloudAccountV1Alpha1().UpdateCloudAccount(ctx, cloudAccount, updateMask)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating cloud account",
			fmt.Sprintf("Could not update cloud account: %s", err),
		)

		return
	}

	newState := protoToModel(updateResp)
	resp.Diagnostics.Append(resp.State.Set(ctx, &newState)...)
}

// updateMaskForConfig derives the UpdateCloudAccount update mask from the
// practitioner's config (not the plan, which UseStateForUnknown always
// resolves to a concrete value): a field is only included if the
// practitioner explicitly set it this apply.
func updateMaskForConfig(config CloudAccountModel) *cloudaccountsvcpb.UpdateCloudAccountRequest_UpdateMask {
	mask := &cloudaccountsvcpb.UpdateCloudAccountRequest_UpdateMask{
		Name:               true,
		Suppressed:         !config.Suppressed.IsNull(),
		ManagedByDiscovery: !config.ManagedByDiscovery.IsNull(),
	}

	if config.AWS != nil {
		mask.AwsLambdaDiscoveryConfig = config.AWS.LambdaDiscoveryConfig != nil
		mask.AwsAgentCoreDiscoveryConfig = config.AWS.AgentCoreDiscoveryConfig != nil
	}

	return mask
}

func (c *CloudAccountResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state CloudAccountModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := c.client.CloudAccountV1Alpha1().DestroyCloudAccount(ctx, state.ID.ValueString())
	if err != nil && status.Code(err) != codes.NotFound {
		resp.Diagnostics.AddError(
			"Error deleting cloud account",
			err.Error(),
		)

		return
	}
}

func (c *CloudAccountResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
