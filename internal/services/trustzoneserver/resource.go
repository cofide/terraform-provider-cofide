package trustzoneserver

import (
	"context"
	"fmt"

	trustzoneserversvcpb "github.com/cofide/cofide-api-sdk/gen/go/proto/connect/trust_zone_server_service/v1alpha1"
	trustzoneserverpb "github.com/cofide/cofide-api-sdk/gen/go/proto/trust_zone_server/v1alpha1"
	sdkclient "github.com/cofide/cofide-api-sdk/pkg/connect/client"
	"github.com/cofide/terraform-provider-cofide/internal/util"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	tftypes "github.com/hashicorp/terraform-plugin-framework/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/structpb"
	"gopkg.in/yaml.v3"
)

var _ resource.Resource = &TrustZoneServerResource{}
var _ resource.ResourceWithImportState = &TrustZoneServerResource{}
var _ resource.ResourceWithValidateConfig = &TrustZoneServerResource{}

type TrustZoneServerResource struct {
	client sdkclient.ClientSet
}

func NewResource() resource.Resource {
	return &TrustZoneServerResource{}
}

func (r *TrustZoneServerResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_connect_trust_zone_server"
}

func (r *TrustZoneServerResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *TrustZoneServerResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan TrustZoneServerModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	server := &trustzoneserverpb.TrustZoneServer{
		TrustZoneId: plan.TrustZoneID.ValueString(),
		ClusterId:   plan.ClusterID.ValueString(),
	}

	if util.IsStringAttributeNonEmpty(plan.KubernetesNamespace) {
		server.KubernetesNamespace = plan.KubernetesNamespace.ValueString()
	}

	if util.IsStringAttributeNonEmpty(plan.KubernetesServiceAccount) {
		server.KubernetesServiceAccount = plan.KubernetesServiceAccount.ValueString()
	}

	if !plan.HelmValues.IsNull() {
		helmValues, err := parseHelmValues(plan.HelmValues)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error parsing helm_values",
				fmt.Sprintf("Failed to parse helm_values: %s", err),
			)
			return
		}
		server.HelmValues = helmValues
	}

	if plan.ConnectK8sPsatConfig != nil {
		server.ConnectK8SPsatConfig = connectK8sPsatConfigToProto(plan.ConnectK8sPsatConfig)
	}

	createResp, err := r.client.TrustZoneServerV1Alpha1().CreateTrustZoneServer(ctx, server)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating trust zone server",
			fmt.Sprintf("Could not create trust zone server: %s", err.Error()),
		)
		return
	}

	state := trustZoneServerFromProto(createResp, plan.HelmValues)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *TrustZoneServerResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state TrustZoneServerModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	serverID := state.ID.ValueString()
	server, err := r.client.TrustZoneServerV1Alpha1().GetTrustZoneServer(ctx, serverID)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error reading trust zone server",
			fmt.Sprintf("Could not read trust zone server %q: %s", serverID, err),
		)
		return
	}

	helmValues := helmValuesForState(server.GetHelmValues(), state.HelmValues)
	newState := trustZoneServerFromProto(server, helmValues)
	resp.Diagnostics.Append(resp.State.Set(ctx, &newState)...)
}

func (r *TrustZoneServerResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan TrustZoneServerModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state TrustZoneServerModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	serverID := state.ID.ValueString()
	server := &trustzoneserverpb.TrustZoneServer{
		Id:          serverID,
		TrustZoneId: plan.TrustZoneID.ValueString(),
		ClusterId:   plan.ClusterID.ValueString(),
	}

	if !plan.HelmValues.IsNull() {
		helmValues, err := parseHelmValues(plan.HelmValues)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error parsing helm_values",
				fmt.Sprintf("Failed to parse helm_values: %s", err),
			)
			return
		}
		server.HelmValues = helmValues
	}

	if plan.ConnectK8sPsatConfig != nil {
		server.ConnectK8SPsatConfig = connectK8sPsatConfigToProto(plan.ConnectK8sPsatConfig)
	}

	updateMask := &trustzoneserversvcpb.UpdateTrustZoneServerRequest_UpdateMask{
		HelmValues: true,
	}

	updateResp, err := r.client.TrustZoneServerV1Alpha1().UpdateTrustZoneServer(ctx, server, updateMask)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating trust zone server",
			err.Error(),
		)
		return
	}

	newState := trustZoneServerFromProto(updateResp, plan.HelmValues)
	resp.Diagnostics.Append(resp.State.Set(ctx, &newState)...)
}

func (r *TrustZoneServerResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state TrustZoneServerModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.TrustZoneServerV1Alpha1().DestroyTrustZoneServer(ctx, state.ID.ValueString())
	if err != nil {
		if status.Code(err) != codes.NotFound {
			resp.Diagnostics.AddError(
				"Error deleting trust zone server",
				err.Error(),
			)
			return
		}
	}
}

func (r *TrustZoneServerResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *TrustZoneServerResource) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	var data TrustZoneServerModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// trustZoneServerFromProto converts a TrustZoneServer proto to a TrustZoneServerModel.
// helmValues is passed separately to preserve the original YAML/JSON format from the plan.
func trustZoneServerFromProto(server *trustzoneserverpb.TrustZoneServer, helmValues tftypes.String) TrustZoneServerModel {
	model := TrustZoneServerModel{
		ID:          tftypes.StringValue(server.GetId()),
		TrustZoneID: tftypes.StringValue(server.GetTrustZoneId()),
		ClusterID:   tftypes.StringValue(server.GetClusterId()),
		OrgID:       tftypes.StringValue(server.GetOrgId()),
		HelmValues:  helmValues,
	}

	if ns := server.GetKubernetesNamespace(); ns != "" {
		model.KubernetesNamespace = tftypes.StringValue(ns)
	} else {
		model.KubernetesNamespace = tftypes.StringNull()
	}

	if sa := server.GetKubernetesServiceAccount(); sa != "" {
		model.KubernetesServiceAccount = tftypes.StringValue(sa)
	} else {
		model.KubernetesServiceAccount = tftypes.StringNull()
	}

	if s := server.GetStatus(); s != nil {
		model.Status = statusFromProto(s)
	} else {
		model.Status = &TrustZoneServerStatusModel{
			Status:             tftypes.StringNull(),
			LastTransitionTime: tftypes.StringNull(),
		}
	}

	if cfg := server.GetConnectK8SPsatConfig(); cfg != nil {
		model.ConnectK8sPsatConfig = connectK8sPsatConfigFromProto(cfg)
	}

	return model
}

// statusFromProto converts a TrustZoneServer_Status proto to a TrustZoneServerStatusModel.
func statusFromProto(s *trustzoneserverpb.TrustZoneServer_Status) *TrustZoneServerStatusModel {
	model := &TrustZoneServerStatusModel{
		Status: tftypes.StringValue(s.GetStatus().String()),
	}

	if t := s.GetLastTransitionTime(); t != nil {
		model.LastTransitionTime = tftypes.StringValue(t.AsTime().Format("2006-01-02T15:04:05Z07:00"))
	} else {
		model.LastTransitionTime = tftypes.StringNull()
	}

	return model
}

// connectK8sPsatConfigToProto converts a ConnectK8sPsatConfigModel to a ConnectK8SPsatConfig proto.
func connectK8sPsatConfigToProto(model *ConnectK8sPsatConfigModel) *trustzoneserverpb.ConnectK8SPsatConfig {
	cfg := &trustzoneserverpb.ConnectK8SPsatConfig{
		SpireServerSpiffeIdPath: model.SpireServerSpiffeIDPath.ValueString(),
	}
	for _, a := range model.Audiences {
		cfg.Audiences = append(cfg.Audiences, a.ValueString())
	}
	return cfg
}

// connectK8sPsatConfigFromProto converts a ConnectK8SPsatConfig proto to a ConnectK8sPsatConfigModel.
func connectK8sPsatConfigFromProto(cfg *trustzoneserverpb.ConnectK8SPsatConfig) *ConnectK8sPsatConfigModel {
	model := &ConnectK8sPsatConfigModel{
		SpireServerSpiffeIDPath: tftypes.StringValue(cfg.GetSpireServerSpiffeIdPath()),
	}
	for _, a := range cfg.GetAudiences() {
		model.Audiences = append(model.Audiences, tftypes.StringValue(a))
	}
	return model
}

// parseHelmValues parses the helm_values field from a YAML/JSON string to a structpb.Struct.
func parseHelmValues(valueStr tftypes.String) (*structpb.Struct, error) {
	if valueStr.IsNull() || valueStr.ValueString() == "" {
		return nil, nil
	}

	var helmValues map[string]any
	if err := yaml.Unmarshal([]byte(valueStr.ValueString()), &helmValues); err != nil {
		return nil, fmt.Errorf("invalid YAML in helm_values: %w", err)
	}

	helmValuesStruct, err := structpb.NewStruct(helmValues)
	if err != nil {
		return nil, fmt.Errorf("failed to convert helm_values to Struct: %w", err)
	}

	return helmValuesStruct, nil
}

// helmValuesForState returns a helm_values string for storage in state, preserving the original
// YAML/JSON format from the plan when the API response matches semantically.
func helmValuesForState(apiValues *structpb.Struct, prev tftypes.String) tftypes.String {
	if apiValues == nil || len(apiValues.Fields) == 0 {
		return tftypes.StringNull()
	}
	// Preserve previous format (YAML or JSON) to avoid spurious diffs.
	if !prev.IsNull() && prev.ValueString() != "" {
		return prev
	}
	jsonBytes, err := apiValues.MarshalJSON()
	if err != nil {
		return tftypes.StringNull()
	}
	return tftypes.StringValue(string(jsonBytes))
}

// helmValuesFromProto converts a structpb.Struct to a JSON string for use in data sources.
func helmValuesFromProto(apiValues *structpb.Struct) tftypes.String {
	if apiValues == nil || len(apiValues.Fields) == 0 {
		return tftypes.StringNull()
	}
	jsonBytes, err := apiValues.MarshalJSON()
	if err != nil {
		return tftypes.StringNull()
	}
	return tftypes.StringValue(string(jsonBytes))
}
