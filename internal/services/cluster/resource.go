package cluster

import (
	"context"
	"fmt"

	sdkclient "github.com/cofide/cofide-api-sdk/pkg/connect/client"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ resource.Resource = &ClusterResource{}
var _ resource.ResourceWithImportState = &ClusterResource{}
var _ resource.ResourceWithValidateConfig = &ClusterResource{}

type ClusterResource struct {
	client sdkclient.ClientSet
}

func NewResource() resource.Resource {
	return &ClusterResource{}
}

func (c *ClusterResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_connect_cluster"
}

func (c *ClusterResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (c *ClusterResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan ClusterModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	cluster, err := modelToProto(plan)
	if err != nil {
		resp.Diagnostics.AddError("Error converting model to proto", err.Error())
		return
	}

	createResp, err := c.client.ClusterV1Alpha1().CreateCluster(ctx, cluster)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating cluster",
			fmt.Sprintf("Could not create cluster: %s", err.Error()),
		)

		return
	}

	if err := validateHelmValues(plan.ExtraHelmValues, createResp.GetExtraHelmValues()); err != nil {
		resp.Diagnostics.AddError("Inconsistent Helm values", err.Error())
		return
	}

	newState, err := protoToModel(createResp)
	if err != nil {
		resp.Diagnostics.AddError("Error converting proto to model", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, newState)...)
}

func (c *ClusterResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state ClusterModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	clusterID := state.ID.ValueString()
	if clusterID == "" {
		resp.Diagnostics.AddError(
			"Error reading cluster",
			"Cluster ID not found in state.",
		)
		return
	}

	getResp, err := c.client.ClusterV1Alpha1().GetCluster(ctx, clusterID)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error reading cluster",
			err.Error(),
		)
		return
	}

	if err := validateHelmValues(state.ExtraHelmValues, getResp.GetExtraHelmValues()); err != nil {
		resp.Diagnostics.AddError("Inconsistent Helm values", err.Error())
		return
	}

	newState, err := protoToModel(getResp)
	if err != nil {
		resp.Diagnostics.AddError("Error converting proto to model", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &newState)...)
}

func (c *ClusterResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var state ClusterModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var plan ClusterModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	clusterID := state.ID.ValueString()
	if clusterID == "" {
		resp.Diagnostics.AddError(
			"Error updating cluster",
			"Cluster ID not found in state. The resource might not have been created properly.",
		)
		return
	}

	cluster, err := modelToProto(plan)
	if err != nil {
		resp.Diagnostics.AddError("Error converting model to proto", err.Error())
		return
	}
	cluster.Id = &clusterID

	updateResp, err := c.client.ClusterV1Alpha1().UpdateCluster(ctx, cluster)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating cluster",
			err.Error(),
		)

		return
	}

	if err := validateHelmValues(plan.ExtraHelmValues, updateResp.GetExtraHelmValues()); err != nil {
		resp.Diagnostics.AddError("Inconsistent Helm values", err.Error())
		return
	}

	newState, err := protoToModel(updateResp)
	if err != nil {
		resp.Diagnostics.AddError("Error converting proto to model", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &newState)...)
}

func (c *ClusterResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state ClusterModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := c.client.ClusterV1Alpha1().DestroyCluster(ctx, state.ID.ValueString())
	if err != nil {
		if status.Code(err) != codes.NotFound {
			resp.Diagnostics.AddError(
				"Error deleting cluster",
				err.Error(),
			)
			return
		}
	}
}

func (c *ClusterResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (c *ClusterResource) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	var data ClusterModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
}
