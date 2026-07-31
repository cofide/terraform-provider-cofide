package awslambdadiscoveryconfig

import (
	"context"
	"fmt"

	cloudaccountpb "github.com/cofide/cofide-api-sdk/gen/go/proto/cloud_account/v1alpha1"
	cloudaccountsvcpb "github.com/cofide/cofide-api-sdk/gen/go/proto/connect/cloud_account_service/v1alpha1"
	sdkclient "github.com/cofide/cofide-api-sdk/pkg/connect/client"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
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
	resp.TypeName = req.ProviderTypeName + "_connect_aws_lambda_discovery_config"
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

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	newState, diags := r.upsert(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &newState)...)
}

func (r *Resource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan Model

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	newState, diags := r.upsert(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &newState)...)
}

// upsert sets the cloud account's lambda_discovery_config to exactly what's
// in the plan, using an update mask so no other field on the cloud account
// is touched.
func (r *Resource) upsert(ctx context.Context, plan Model) (Model, diag.Diagnostics) {
	lambdaConfig, diags := modelToProto(ctx, plan)
	if diags.HasError() {
		return Model{}, diags
	}

	cloudAccountID := plan.CloudAccountID.ValueString()
	cloudAccount := &cloudaccountpb.CloudAccount{
		Id: cloudAccountID,
		Provider: &cloudaccountpb.CloudAccount_Aws{
			Aws: &cloudaccountpb.AWSAccount{LambdaDiscoveryConfig: lambdaConfig},
		},
	}
	updateMask := &cloudaccountsvcpb.UpdateCloudAccountRequest_UpdateMask{AwsLambdaDiscoveryConfig: true}

	updateResp, err := r.client.CloudAccountV1Alpha1().UpdateCloudAccount(ctx, cloudAccount, updateMask)
	if err != nil {
		diags.AddError(
			"Error setting AWS Lambda discovery config",
			fmt.Sprintf("Could not update cloud account %q: %s", cloudAccountID, err),
		)
		return Model{}, diags
	}

	return protoToModel(cloudAccountID, updateResp.GetAws().GetLambdaDiscoveryConfig()), diags
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
			"Error reading AWS Lambda discovery config",
			fmt.Sprintf("Could not read cloud account %q: %s", cloudAccountID, err),
		)

		return
	}

	lambdaConfig := cloudAccount.GetAws().GetLambdaDiscoveryConfig()
	if lambdaConfig == nil {
		// The config no longer exists on the cloud account (e.g. cleared
		// inline on the cloud_account resource, or by discovery).
		resp.State.RemoveResource(ctx)
		return
	}

	newState := protoToModel(cloudAccountID, lambdaConfig)
	resp.Diagnostics.Append(resp.State.Set(ctx, &newState)...)
}

func (r *Resource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state Model

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	cloudAccount := &cloudaccountpb.CloudAccount{
		Id: state.CloudAccountID.ValueString(),
		Provider: &cloudaccountpb.CloudAccount_Aws{
			Aws: &cloudaccountpb.AWSAccount{LambdaDiscoveryConfig: nil},
		},
	}
	updateMask := &cloudaccountsvcpb.UpdateCloudAccountRequest_UpdateMask{AwsLambdaDiscoveryConfig: true}

	_, err := r.client.CloudAccountV1Alpha1().UpdateCloudAccount(ctx, cloudAccount, updateMask)
	if err != nil && status.Code(err) != codes.NotFound {
		resp.Diagnostics.AddError(
			"Error deleting AWS Lambda discovery config",
			err.Error(),
		)

		return
	}
}

func (r *Resource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
