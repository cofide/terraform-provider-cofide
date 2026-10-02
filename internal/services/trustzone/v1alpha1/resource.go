package v1alpha1

import (
	"context"
	"fmt"

	trustzonepb "github.com/cofide/cofide-api-sdk/gen/go/proto/trust_zone/v1alpha1"
	sdkclient "github.com/cofide/cofide-api-sdk/pkg/connect/client"
	"github.com/cofide/terraform-provider-cofide/internal/util"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	tftypes "github.com/hashicorp/terraform-plugin-framework/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	_ resource.Resource                = &TrustZoneResource{}
	_ resource.ResourceWithImportState = &TrustZoneResource{}
)

type TrustZoneResource struct {
	client sdkclient.ClientSet
}

func NewResource() resource.Resource {
	return &TrustZoneResource{}
}

func (t *TrustZoneResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_connect_trust_zone_v1alpha1"
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

	if util.IsStringAttributeNonEmpty(plan.OrgID) {
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

	trustZone, err := newUpdateRequest(plan, state)
	if err != nil {
		resp.Diagnostics.AddError("Error updating trust zone", fmt.Sprintf("Could not build trust zone update request: %s", err.Error()))
		return
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

	newState := TrustZoneModel{
		ID:                    tftypes.StringValue(updateResp.GetId()),
		Name:                  tftypes.StringValue(updateResp.GetName()),
		TrustDomain:           tftypes.StringValue(updateResp.GetTrustDomain()),
		OrgID:                 orgIDStr,
		IsManagementZone:      tftypes.BoolValue(updateResp.GetIsManagementZone()),
		BundleEndpointURL:     tftypes.StringValue(updateResp.GetBundleEndpointUrl()),
		BundleEndpointProfile: tftypes.StringValue(updateResp.GetBundleEndpointProfile().String()),
		JWTIssuer:             tftypes.StringValue(updateResp.GetJwtIssuer()),
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &newState)...)
}

// newUpdateRequest builds the trust zone sent to UpdateTrustZone. Connect
// replaces the whole trust zone on update, so the computed fields that aren't
// configurable (bundle endpoint URL and profile, JWT issuer) are carried over
// from state; omitting them would clear them, which Connect rejects.
func newUpdateRequest(plan, state TrustZoneModel) (*trustzonepb.TrustZone, error) {
	trustZoneID := state.ID.ValueString()

	trustZone := &trustzonepb.TrustZone{
		Id:               &trustZoneID,
		Name:             plan.Name.ValueString(),
		TrustDomain:      plan.TrustDomain.ValueString(),
		IsManagementZone: plan.IsManagementZone.ValueBool(),
	}

	if util.IsStringAttributeNonEmpty(plan.OrgID) {
		trustZone.OrgId = plan.OrgID.ValueStringPointer()
	}

	if util.IsStringAttributeNonEmpty(state.BundleEndpointURL) {
		trustZone.BundleEndpointUrl = state.BundleEndpointURL.ValueStringPointer()
	}

	if util.IsStringAttributeNonEmpty(state.JWTIssuer) {
		trustZone.JwtIssuer = state.JWTIssuer.ValueStringPointer()
	}

	if util.IsStringAttributeNonEmpty(state.BundleEndpointProfile) {
		profileName := state.BundleEndpointProfile.ValueString()
		profile, ok := trustzonepb.BundleEndpointProfile_value[profileName]
		if !ok {
			return nil, fmt.Errorf("unknown bundle endpoint profile %q in state", profileName)
		}
		trustZone.BundleEndpointProfile = trustzonepb.BundleEndpointProfile(profile).Enum()
	}

	return trustZone, nil
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
