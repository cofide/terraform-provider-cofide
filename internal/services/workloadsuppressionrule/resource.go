package workloadsuppressionrule

import (
	"context"
	"fmt"

	workloadsuppressionrulesvcpb "github.com/cofide/cofide-api-sdk/gen/go/proto/connect/workload_suppression_rule_service/v1alpha1"
	sdkclient "github.com/cofide/cofide-api-sdk/pkg/connect/client"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	_ resource.Resource                     = &WorkloadSuppressionRuleResource{}
	_ resource.ResourceWithImportState      = &WorkloadSuppressionRuleResource{}
	_ resource.ResourceWithConfigValidators = &WorkloadSuppressionRuleResource{}
)

type WorkloadSuppressionRuleResource struct {
	client sdkclient.ClientSet
}

func NewResource() resource.Resource {
	return &WorkloadSuppressionRuleResource{}
}

func (r *WorkloadSuppressionRuleResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_connect_workload_suppression_rule"
}

func (r *WorkloadSuppressionRuleResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *WorkloadSuppressionRuleResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan WorkloadSuppressionRuleModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	rule, diags := modelToProto(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	createResp, err := r.client.WorkloadSuppressionRuleV1Alpha1().CreateWorkloadSuppressionRule(ctx, rule)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating workload suppression rule",
			fmt.Sprintf("Could not create workload suppression rule: %s", err.Error()),
		)
		return
	}

	state := protoToModel(createResp)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *WorkloadSuppressionRuleResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state WorkloadSuppressionRuleModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ruleID := state.ID.ValueString()
	rule, err := r.client.WorkloadSuppressionRuleV1Alpha1().GetWorkloadSuppressionRule(ctx, ruleID)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error reading workload suppression rule",
			fmt.Sprintf("Could not read workload suppression rule %q: %s", ruleID, err),
		)
		return
	}

	newState := protoToModel(rule)
	resp.Diagnostics.Append(resp.State.Set(ctx, &newState)...)
}

func (r *WorkloadSuppressionRuleResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var state WorkloadSuppressionRuleModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var plan WorkloadSuppressionRuleModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ruleID := state.ID.ValueString()
	if ruleID == "" {
		resp.Diagnostics.AddError(
			"Error updating workload suppression rule",
			"Rule ID not found in state. The resource might not have been created properly.",
		)
		return
	}

	rule, diags := modelToProto(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	rule.Id = ruleID

	updateMask := &workloadsuppressionrulesvcpb.UpdateWorkloadSuppressionRuleRequest_UpdateMask{
		Name:        true,
		Description: true,
		Enabled:     true,
		Matcher:     true,
	}

	updateResp, err := r.client.WorkloadSuppressionRuleV1Alpha1().UpdateWorkloadSuppressionRule(ctx, rule, updateMask)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating workload suppression rule",
			fmt.Sprintf("Could not update workload suppression rule: %s", err.Error()),
		)
		return
	}

	newState := protoToModel(updateResp)
	resp.Diagnostics.Append(resp.State.Set(ctx, &newState)...)
}

func (r *WorkloadSuppressionRuleResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state WorkloadSuppressionRuleModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.WorkloadSuppressionRuleV1Alpha1().DestroyWorkloadSuppressionRule(ctx, state.ID.ValueString())
	if err != nil {
		if status.Code(err) != codes.NotFound {
			resp.Diagnostics.AddError(
				"Error deleting workload suppression rule",
				fmt.Sprintf("Could not delete workload suppression rule: %s", err),
			)
			return
		}
	}
}

func (r *WorkloadSuppressionRuleResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
