package cloudorganization

import (
	"context"
	"fmt"

	sdkclient "github.com/cofide/cofide-api-sdk/pkg/connect/client"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ resource.Resource = &CloudOrganizationResource{}
var _ resource.ResourceWithImportState = &CloudOrganizationResource{}

type CloudOrganizationResource struct {
	client sdkclient.ClientSet
}

func NewResource() resource.Resource {
	return &CloudOrganizationResource{}
}

func (c *CloudOrganizationResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_connect_cloud_organization"
}

func (c *CloudOrganizationResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (c *CloudOrganizationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan CloudOrganizationModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	createResp, err := c.client.CloudOrganizationV1Alpha1().CreateCloudOrganization(ctx, modelToProto(plan))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating cloud organization",
			fmt.Sprintf("Could not create cloud organization: %s", err),
		)

		return
	}

	state := protoToModel(createResp)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (c *CloudOrganizationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state CloudOrganizationModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	cloudOrganizationID := state.ID.ValueString()
	cloudOrganization, err := c.client.CloudOrganizationV1Alpha1().GetCloudOrganization(ctx, cloudOrganizationID)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error reading cloud organization",
			fmt.Sprintf("Could not read cloud organization %q: %s", cloudOrganizationID, err),
		)

		return
	}

	newState := protoToModel(cloudOrganization)
	resp.Diagnostics.Append(resp.State.Set(ctx, &newState)...)
}

func (c *CloudOrganizationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan CloudOrganizationModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state CloudOrganizationModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	cloudOrganization := modelToProto(plan)
	cloudOrganization.Id = state.ID.ValueString()

	// No update_mask is passed: the API performs a full replacement of the
	// mutable fields, which is what we want here since cloud_organization is
	// always fully owned by this resource (unlike cloud_account, which has
	// fields that can also be managed by standalone resources).
	updateResp, err := c.client.CloudOrganizationV1Alpha1().UpdateCloudOrganization(ctx, cloudOrganization, nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating cloud organization",
			fmt.Sprintf("Could not update cloud organization: %s", err),
		)

		return
	}

	newState := protoToModel(updateResp)
	resp.Diagnostics.Append(resp.State.Set(ctx, &newState)...)
}

func (c *CloudOrganizationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state CloudOrganizationModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := c.client.CloudOrganizationV1Alpha1().DestroyCloudOrganization(ctx, state.ID.ValueString())
	if err != nil && status.Code(err) != codes.NotFound {
		resp.Diagnostics.AddError(
			"Error deleting cloud organization",
			err.Error(),
		)

		return
	}
}

func (c *CloudOrganizationResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
