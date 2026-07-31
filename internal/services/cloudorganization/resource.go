package cloudorganization

import (
	"context"
	"fmt"

	cloudorganizationsvcpb "github.com/cofide/cofide-api-sdk/gen/go/proto/connect/cloud_organization_service/v1alpha1"
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

	cloudOrganization, diags := modelToProto(plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	createResp, err := c.client.CloudOrganizationV1Alpha1().CreateCloudOrganization(ctx, cloudOrganization)
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

	cloudOrganization, diags := modelToProto(plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	cloudOrganization.Id = state.ID.ValueString()

	updateResp, err := c.client.CloudOrganizationV1Alpha1().UpdateCloudOrganization(ctx, cloudOrganization, fullUpdateMask())
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

// fullUpdateMask returns an update mask covering every mutable field of a
// cloud organization. The API rejects a nil update_mask (it does not treat
// that as "update all fields"), so every field must be masked explicitly to
// achieve a full replacement. This is safe because cloud_organization is
// always fully owned by this resource (unlike cloud_account, which has
// fields that can also be managed by standalone resources).
func fullUpdateMask() *cloudorganizationsvcpb.UpdateCloudOrganizationRequest_UpdateMask {
	return &cloudorganizationsvcpb.UpdateCloudOrganizationRequest_UpdateMask{
		Name:                 true,
		AwsAudience:          true,
		DiscoveryEnabled:     true,
		AwsRoleChain:         true,
		DiscoveryInterval:    true,
		AwsAssumeThroughOidc: true,
	}
}

func (c *CloudOrganizationResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
