package attestationpolicy

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
	_ resource.Resource                   = &AttestationPolicyResource{}
	_ resource.ResourceWithImportState    = &AttestationPolicyResource{}
	_ resource.ResourceWithValidateConfig = &AttestationPolicyResource{}
)

type AttestationPolicyResource struct {
	client sdkclient.ClientSet
}

func NewResource() resource.Resource {
	return &AttestationPolicyResource{}
}

func (r *AttestationPolicyResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_connect_attestation_policy"
}

func (r *AttestationPolicyResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *AttestationPolicyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan AttestationPolicyModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	policy := modelToProto(plan)
	createResp, err := r.client.AttestationPolicyV1Alpha1().CreateAttestationPolicy(ctx, policy)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating attestation policy",
			fmt.Sprintf("Could not create attestation policy: %s", err.Error()),
		)
		return
	}

	state := protoToModel(createResp)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *AttestationPolicyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state AttestationPolicyModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *AttestationPolicyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var state AttestationPolicyModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var plan AttestationPolicyModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	policyID := state.ID.ValueString()
	if policyID == "" {
		resp.Diagnostics.AddError(
			"Error updating attestation policy",
			"Policy ID not found in state. The resource might not have been created properly.",
		)
		return
	}

	policy := modelToProto(plan)
	policy.Id = &policyID

	updateResp, err := r.client.AttestationPolicyV1Alpha1().UpdateAttestationPolicy(ctx, policy)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating attestation policy",
			fmt.Sprintf("Could not update attestation policy: %s", err.Error()),
		)
		return
	}

	newState := protoToModel(updateResp)
	resp.Diagnostics.Append(resp.State.Set(ctx, &newState)...)
}

func (r *AttestationPolicyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state AttestationPolicyModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.AttestationPolicyV1Alpha1().DestroyAttestationPolicy(ctx, state.ID.ValueString())
	if err != nil {
		if status.Code(err) != codes.NotFound {
			resp.Diagnostics.AddError(
				"Error deleting attestation policy",
				fmt.Sprintf("Could not delete attestation policy: %s", err),
			)
			return
		}
	}
}

func (r *AttestationPolicyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *AttestationPolicyResource) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	var data AttestationPolicyModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	hasKubernetes := data.Kubernetes != nil
	hasStatic := data.Static != nil

	if !hasKubernetes && !hasStatic {
		resp.Diagnostics.AddError(
			"Invalid configuration",
			"Either kubernetes or static block must be configured, but neither was provided.",
		)
		return
	}

	if hasKubernetes && hasStatic {
		resp.Diagnostics.AddError(
			"Invalid configuration",
			"Only one of kubernetes or static block can be configured, but both were provided.",
		)
		return
	}
}
