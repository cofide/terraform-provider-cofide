package attestationpolicy

import (
	"context"
	"fmt"

	sdkclient "github.com/cofide/cofide-api-sdk/pkg/connect/client"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	_ resource.Resource                     = &AttestationPolicyResource{}
	_ resource.ResourceWithImportState      = &AttestationPolicyResource{}
	_ resource.ResourceWithConfigValidators = &AttestationPolicyResource{}
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

	policyID := state.ID.ValueString()
	policy, err := r.client.AttestationPolicyV1Alpha1().GetAttestationPolicy(ctx, policyID)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error reading attestation policy",
			fmt.Sprintf("Could not read attestation policy %q: %s", policyID, err),
		)
		return
	}

	newState := protoToModel(policy)
	resp.Diagnostics.Append(resp.State.Set(ctx, &newState)...)
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
	name := req.ID

	policies, err := r.client.AttestationPolicyV1Alpha1().ListAttestationPolicies(ctx, nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error importing attestation policy",
			fmt.Sprintf("Could not list attestation policies: %s", err),
		)
		return
	}

	for _, p := range policies {
		if p.GetName() == name {
			state := protoToModel(p)
			resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
			return
		}
	}

	resp.Diagnostics.AddError(
		"Error importing attestation policy",
		fmt.Sprintf("Could not find attestation policy with name %q", name),
	)
}
