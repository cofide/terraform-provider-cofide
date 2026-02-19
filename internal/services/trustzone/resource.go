package trustzone

import (
	"context"
	"fmt"

	trustzonepb "github.com/cofide/cofide-api-sdk/gen/go/proto/trust_zone/v1alpha1"
	sdkclient "github.com/cofide/cofide-api-sdk/pkg/connect/client"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	tftypes "github.com/hashicorp/terraform-plugin-framework/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	_ resource.Resource                   = &TrustZoneResource{}
	_ resource.ResourceWithImportState    = &TrustZoneResource{}
	_ resource.ResourceWithValidateConfig = &TrustZoneResource{}
)

type TrustZoneResource struct {
	client sdkclient.ClientSet
}

func NewResource() resource.Resource {
	return &TrustZoneResource{}
}

func (t *TrustZoneResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_connect_trust_zone"
}

func (t *TrustZoneResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	t.client = client
}

func (t *TrustZoneResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan TrustZoneModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	trustZone := &trustzonepb.TrustZone{
		Name:        plan.Name.ValueString(),
		TrustDomain: plan.TrustDomain.ValueString(),
	}

	if !plan.OrgID.IsNull() {
		trustZone.OrgId = plan.OrgID.ValueStringPointer()
	}

	if !plan.IsManagementZone.IsNull() {
		trustZone.IsManagementZone = plan.IsManagementZone.ValueBool()
	}

	createResp, err := t.client.TrustZoneV1Alpha1().CreateTrustZone(ctx, trustZone)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating trust zone",
			fmt.Sprintf("Could not create trust zone: %s", err.Error()),
		)
		return
	}

	state := TrustZoneModel{
		ID:                    tftypes.StringValue(createResp.GetId()),
		Name:                  tftypes.StringValue(createResp.GetName()),
		TrustDomain:           tftypes.StringValue(createResp.GetTrustDomain()),
		OrgID:                 tftypes.StringValue(createResp.GetOrgId()),
		IsManagementZone:      tftypes.BoolValue(createResp.GetIsManagementZone()),
		BundleEndpointURL:     tftypes.StringValue(createResp.GetBundleEndpointUrl()),
		BundleEndpointProfile: tftypes.StringValue(createResp.GetBundleEndpointProfile().String()),
		JWTIssuer:             tftypes.StringValue(createResp.GetJwtIssuer()),
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (t *TrustZoneResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state TrustZoneModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	trustZoneID := state.ID.ValueString()
	trustZone, err := t.client.TrustZoneV1Alpha1().GetTrustZone(ctx, trustZoneID)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error reading trust zone",
			fmt.Sprintf("Could not read trust zone %q: %s", trustZoneID, err),
		)
		return
	}

	newState := TrustZoneModel{
		ID:                    tftypes.StringValue(trustZone.GetId()),
		Name:                  tftypes.StringValue(trustZone.GetName()),
		TrustDomain:           tftypes.StringValue(trustZone.GetTrustDomain()),
		OrgID:                 tftypes.StringValue(trustZone.GetOrgId()),
		IsManagementZone:      tftypes.BoolValue(trustZone.GetIsManagementZone()),
		BundleEndpointURL:     tftypes.StringValue(trustZone.GetBundleEndpointUrl()),
		BundleEndpointProfile: tftypes.StringValue(trustZone.GetBundleEndpointProfile().String()),
		JWTIssuer:             tftypes.StringValue(trustZone.GetJwtIssuer()),
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &newState)...)
}

func (t *TrustZoneResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan TrustZoneModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state TrustZoneModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	trustZoneID := state.ID.ValueString()

	trustZone := &trustzonepb.TrustZone{
		Id:               &trustZoneID,
		Name:             plan.Name.ValueString(),
		TrustDomain:      plan.TrustDomain.ValueString(),
		IsManagementZone: plan.IsManagementZone.ValueBool(),
	}

	if !plan.OrgID.IsNull() && plan.OrgID.ValueString() != "" {
		trustZone.OrgId = plan.OrgID.ValueStringPointer()
	}

	updateResp, err := t.client.TrustZoneV1Alpha1().UpdateTrustZone(ctx, trustZone)
	if err != nil {
		resp.Diagnostics.AddError("Error updating trust zone", fmt.Sprintf("Could not update trust zone: %s", err.Error()))
		return
	}

	var orgIDStr tftypes.String
	if orgID := updateResp.GetOrgId(); orgID != "" {
		orgIDStr = tftypes.StringValue(orgID)
	} else if !plan.OrgID.IsNull() {
		orgIDStr = plan.OrgID
	} else {
		orgIDStr = tftypes.StringNull()
	}

	// is_management_zone is a bool, so we can't check for empty string, but we can preserve the plan value
	var isMgmtZoneBool tftypes.Bool
	if updateResp.GetIsManagementZone() || !plan.IsManagementZone.IsNull() {
		// If API returns true, or if there was a value in the plan, use the API response.
		// If API is false AND plan was null, this will result in false, which is the only option.
		isMgmtZoneBool = tftypes.BoolValue(updateResp.GetIsManagementZone())
	} else {
		isMgmtZoneBool = tftypes.BoolNull()
	}

	newState := TrustZoneModel{
		ID:                    tftypes.StringValue(updateResp.GetId()),
		Name:                  tftypes.StringValue(updateResp.GetName()),
		TrustDomain:           tftypes.StringValue(updateResp.GetTrustDomain()),
		OrgID:                 orgIDStr,
		IsManagementZone:      isMgmtZoneBool,
		BundleEndpointURL:     tftypes.StringValue(updateResp.GetBundleEndpointUrl()),
		BundleEndpointProfile: tftypes.StringValue(updateResp.GetBundleEndpointProfile().String()),
		JWTIssuer:             tftypes.StringValue(updateResp.GetJwtIssuer()),
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &newState)...)
}

func (t *TrustZoneResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state TrustZoneModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := t.client.TrustZoneV1Alpha1().DestroyTrustZone(ctx, state.ID.ValueString())
	if err != nil {
		if status.Code(err) != codes.NotFound {
			resp.Diagnostics.AddError(
				"Error deleting trust zone",
				err.Error(),
			)

			return
		}
	}
}

func (t *TrustZoneResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (t *TrustZoneResource) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	var data TrustZoneModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if !data.IsManagementZone.IsNull() {
		resp.Diagnostics.AddWarning(
			"is_management_zone is immutable",
			"The is_management_zone field cannot be modified after creation. Create a new trust zone instead.",
		)
	}
}

// getBundleEndpointProfile converts a string to a BundleEndpointProfile enum pointer
func getBundleEndpointProfile(value string) (*trustzonepb.BundleEndpointProfile, bool) {
	if profileVal, ok := trustzonepb.BundleEndpointProfile_value[value]; ok {
		profile := trustzonepb.BundleEndpointProfile(profileVal)
		return &profile, true
	}
	return nil, false
}
