package federation

import (
	"context"
	"fmt"

	federationpb "github.com/cofide/cofide-api-sdk/gen/go/proto/federation/v1alpha1"
	sdkclient "github.com/cofide/cofide-api-sdk/pkg/connect/client"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	tftypes "github.com/hashicorp/terraform-plugin-framework/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ resource.Resource = &FederationResource{}
var _ resource.ResourceWithImportState = &FederationResource{}
var _ resource.ResourceWithValidateConfig = &FederationResource{}

type FederationResource struct {
	client sdkclient.ClientSet
}

func NewResource() resource.Resource {
	return &FederationResource{}
}

func (f *FederationResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_connect_federation"
}

func (f *FederationResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(sdkclient.ClientSet)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected sdkclient.ClientSet, got: %T", req.ProviderData),
		)

		return
	}

	f.client = client
}

func (f *FederationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan FederationModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	federation := &federationpb.Federation{
		TrustZoneId:       plan.TrustZoneID.ValueStringPointer(),
		RemoteTrustZoneId: plan.RemoteTrustZoneID.ValueStringPointer(),
	}

	if !plan.OrgID.IsNull() {
		federation.OrgId = plan.OrgID.ValueStringPointer()
	}

	createResp, err := f.client.FederationV1Alpha1().CreateFederation(ctx, federation)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating Federation",
			fmt.Sprintf("Could not create federation: %s", err.Error()),
		)

		return
	}

	state := FederationModel{
		ID:                tftypes.StringValue(createResp.GetId()),
		OrgID:             tftypes.StringValue(createResp.GetOrgId()),
		TrustZoneID:       tftypes.StringValue(createResp.GetTrustZoneId()),
		RemoteTrustZoneID: tftypes.StringValue(createResp.GetRemoteTrustZoneId()),
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (f *FederationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state FederationModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	federationID := state.ID.ValueString()
	federation, err := f.client.FederationV1Alpha1().GetFederation(ctx, federationID)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error reading federation",
			fmt.Sprintf("Could not read federation %q: %s", federationID, err),
		)
		return
	}

	newState := FederationModel{
		ID:                tftypes.StringValue(federation.GetId()),
		OrgID:             tftypes.StringValue(federation.GetOrgId()),
		TrustZoneID:       tftypes.StringValue(federation.GetTrustZoneId()),
		RemoteTrustZoneID: tftypes.StringValue(federation.GetRemoteTrustZoneId()),
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &newState)...)
}

func (f *FederationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.Diagnostics.AddError(
		"Federation Update Not Supported",
		"The Connect API does not support updating federations.",
	)
}

func (f *FederationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state FederationModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := f.client.FederationV1Alpha1().DestroyFederation(ctx, state.ID.ValueString())
	if err != nil {
		if status.Code(err) != codes.NotFound {
			resp.Diagnostics.AddError(
				"Error deleting federation",
				err.Error(),
			)

			return
		}
	}
}

func (f *FederationResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (f *FederationResource) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	var data FederationModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
}
