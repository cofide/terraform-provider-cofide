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

var _ resource.Resource = &TrustZoneResource{}
var _ resource.ResourceWithImportState = &TrustZoneResource{}
var _ resource.ResourceWithValidateConfig = &TrustZoneResource{}

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
		Name:             plan.Name.ValueString(),
		TrustDomain:      plan.TrustDomain.ValueString(),
		IsManagementZone: plan.IsManagementZone.ValueBool(),
	}

	createResp, err := t.client.TrustZoneV1Alpha1().CreateTrustZone(ctx, trustZone)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating trust zone",
			fmt.Sprintf("Could not create trust zone: %s", err.Error()),
		)
		return
	}

	plan.ID = tftypes.StringValue(createResp.GetId())
	plan.Name = tftypes.StringValue(createResp.GetName())
	plan.TrustDomain = tftypes.StringValue(createResp.GetTrustDomain())
	plan.OrgID = tftypes.StringValue(createResp.GetOrgId())
	plan.IsManagementZone = tftypes.BoolValue(createResp.GetIsManagementZone())
	plan.BundleEndpointURL = tftypes.StringValue(createResp.GetBundleEndpointUrl())
	plan.BundleEndpointProfile = tftypes.StringValue(createResp.GetBundleEndpointProfile().String())
	plan.JWTIssuer = tftypes.StringValue(createResp.GetJwtIssuer())

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
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
}

func (t *TrustZoneResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var state TrustZoneModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var plan TrustZoneModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	trustZoneID := state.ID.ValueString()
	if trustZoneID == "" {
		resp.Diagnostics.AddError(
			"Error updating trust zone",
			"Trust zone ID not found in state. The resource might not have been created properly.",
		)
		return
	}

	orgID := state.OrgID.ValueString()
	if orgID == "" {
		resp.Diagnostics.AddError(
			"Error updating trust zone",
			"Org ID not found in state. The resource might not have been created properly.",
		)
		return
	}

	trustZone := &trustzonepb.TrustZone{
		Id:    &trustZoneID,
		Name:  plan.Name.ValueString(),
		OrgId: &orgID,
	}

	if !plan.TrustDomain.IsNull() && plan.TrustDomain.ValueString() != "" {
		trustZone.TrustDomain = plan.TrustDomain.ValueString()
	} else {
		trustZone.TrustDomain = state.TrustDomain.ValueString()
	}

	if !plan.IsManagementZone.IsNull() {
		trustZone.IsManagementZone = plan.IsManagementZone.ValueBool()
	} else {
		trustZone.IsManagementZone = state.IsManagementZone.ValueBool()
	}

	if !plan.BundleEndpointURL.IsNull() && plan.BundleEndpointURL.ValueString() != "" {
		trustZone.BundleEndpointUrl = plan.BundleEndpointURL.ValueStringPointer()
	} else {
		trustZone.BundleEndpointUrl = state.BundleEndpointURL.ValueStringPointer()
	}

	var profileStr string

	if !plan.BundleEndpointProfile.IsNull() && plan.BundleEndpointProfile.ValueString() != "" {
		profileStr = plan.BundleEndpointProfile.ValueString()
	} else {
		profileStr = state.BundleEndpointProfile.ValueString()
	}

	if profile, ok := getBundleEndpointProfile(profileStr); ok {
		trustZone.BundleEndpointProfile = profile
	} else {
		resp.Diagnostics.AddWarning(
			"Unknown BundleEndpointProfile",
			fmt.Sprintf("Value '%s' is not recognized. This might be due to API version mismatch.", profileStr),
		)
	}

	if !plan.JWTIssuer.IsNull() && plan.JWTIssuer.ValueString() != "" {
		trustZone.JwtIssuer = plan.JWTIssuer.ValueStringPointer()
	} else {
		trustZone.JwtIssuer = state.JWTIssuer.ValueStringPointer()
	}

	updateResp, err := t.client.TrustZoneV1Alpha1().UpdateTrustZone(ctx, trustZone)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating trust zone",
			fmt.Sprintf("Could not update trust zone: %s", err.Error()),
		)
		return
	}

	plan.ID = tftypes.StringValue(updateResp.GetId())
	plan.Name = tftypes.StringValue(updateResp.GetName())
	plan.TrustDomain = tftypes.StringValue(updateResp.GetTrustDomain())
	plan.OrgID = tftypes.StringValue(updateResp.GetOrgId())
	plan.IsManagementZone = tftypes.BoolValue(updateResp.GetIsManagementZone())
	plan.BundleEndpointURL = tftypes.StringValue(updateResp.GetBundleEndpointUrl())
	plan.BundleEndpointProfile = tftypes.StringValue(updateResp.GetBundleEndpointProfile().String())
	plan.JWTIssuer = tftypes.StringValue(updateResp.GetJwtIssuer())

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
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
