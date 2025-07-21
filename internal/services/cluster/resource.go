package cluster

import (
	"context"
	"encoding/json"
	"fmt"

	clusterpb "github.com/cofide/cofide-api-sdk/gen/go/proto/cluster/v1alpha1"
	trustproviderpb "github.com/cofide/cofide-api-sdk/gen/go/proto/trust_provider/v1alpha1"
	sdkclient "github.com/cofide/cofide-api-sdk/pkg/connect/client"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	tftypes "github.com/hashicorp/terraform-plugin-framework/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/structpb"
	"gopkg.in/yaml.v3"
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

	cluster := &clusterpb.Cluster{
		Name:              plan.Name.ValueStringPointer(),
		OrgId:             plan.OrgID.ValueStringPointer(),
		TrustZoneId:       plan.TrustZoneID.ValueStringPointer(),
		KubernetesContext: plan.KubernetesContext.ValueStringPointer(),
		Profile:           plan.Profile.ValueStringPointer(),
		ExternalServer:    plan.ExternalServer.ValueBoolPointer(),
	}

	if !plan.OidcIssuerURL.IsNull() && plan.OidcIssuerURL.ValueString() != "" {
		cluster.OidcIssuerUrl = plan.OidcIssuerURL.ValueStringPointer()
	}

	if !plan.OidcIssuerCaCert.IsNull() && plan.OidcIssuerCaCert.ValueString() != "" {
		cluster.OidcIssuerCaCert = []byte(plan.OidcIssuerCaCert.ValueString())
	}

	trustProvider, err := newTrustProvider(plan.TrustProvider.Kind.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating trust provider",
			fmt.Sprintf("Failed to create trust provider: %s", err),
		)

		return
	}

	cluster.TrustProvider = trustProvider

	if !plan.ExtraHelmValues.IsNull() {
		parsedHelmValues, err := parseExtraHelmValues(plan.ExtraHelmValues)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error parsing extra_helm_values",
				fmt.Sprintf("Failed to parse extra_helm_values: %s", err),
			)

			return
		}

		cluster.ExtraHelmValues = parsedHelmValues
	}

	createResp, err := c.client.ClusterV1Alpha1().CreateCluster(ctx, cluster)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating cluster",
			fmt.Sprintf("Could not create cluster: %s", err.Error()),
		)

		return
	}

	plan.ID = tftypes.StringValue(createResp.GetId())
	plan.Name = tftypes.StringValue(createResp.GetName())
	plan.OrgID = tftypes.StringValue(createResp.GetOrgId())
	plan.TrustZoneID = tftypes.StringValue(createResp.GetTrustZoneId())
	plan.KubernetesContext = tftypes.StringValue(createResp.GetKubernetesContext())
	plan.Profile = tftypes.StringValue(createResp.GetProfile())
	plan.ExternalServer = tftypes.BoolValue(createResp.GetExternalServer())

	if createResp.GetOidcIssuerUrl() != "" {
		plan.OidcIssuerURL = tftypes.StringValue(createResp.GetOidcIssuerUrl())
	} else {
		plan.OidcIssuerURL = tftypes.StringNull()
	}

	if len(createResp.GetOidcIssuerCaCert()) > 0 {
		plan.OidcIssuerCaCert = tftypes.StringValue(string(createResp.GetOidcIssuerCaCert()))
	} else {
		plan.OidcIssuerCaCert = tftypes.StringNull()
	}

	plan.TrustProvider = &TrustProviderModel{
		Kind: tftypes.StringValue(createResp.GetTrustProvider().GetKind()),
	}

	extraHelmValues := createResp.GetExtraHelmValues()
	if err := validateHelmValues(plan.ExtraHelmValues, extraHelmValues); err != nil {
		resp.Diagnostics.AddError("Inconsistent Helm values", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (c *ClusterResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state ClusterModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
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

	cluster := &clusterpb.Cluster{
		Id:          &clusterID,
		Name:        plan.Name.ValueStringPointer(),
		OrgId:       plan.OrgID.ValueStringPointer(),
		TrustZoneId: plan.TrustZoneID.ValueStringPointer(),
	}

	if !plan.KubernetesContext.IsNull() && plan.KubernetesContext.ValueString() != "" {
		cluster.KubernetesContext = plan.KubernetesContext.ValueStringPointer()
	} else {
		cluster.KubernetesContext = state.KubernetesContext.ValueStringPointer()
	}

	if !plan.Profile.IsNull() && plan.Profile.ValueString() != "" {
		cluster.Profile = plan.Profile.ValueStringPointer()
	} else {
		cluster.Profile = state.Profile.ValueStringPointer()
	}

	if !plan.ExternalServer.IsNull() {
		cluster.ExternalServer = plan.ExternalServer.ValueBoolPointer()
	} else {
		cluster.ExternalServer = state.ExternalServer.ValueBoolPointer()
	}

	if !plan.OidcIssuerURL.IsNull() && plan.OidcIssuerURL.ValueString() != "" {
		url := plan.OidcIssuerURL.ValueString()
		cluster.OidcIssuerUrl = &url
	}
	if !plan.OidcIssuerCaCert.IsNull() {
		cluster.OidcIssuerCaCert = []byte(plan.OidcIssuerCaCert.ValueString())
	} else {
		cluster.OidcIssuerCaCert = nil
	}

	if !plan.TrustProvider.Kind.IsNull() && plan.TrustProvider.Kind.ValueString() != "" {
		trustProvider, err := newTrustProvider(plan.TrustProvider.Kind.ValueString())
		if err != nil {
			resp.Diagnostics.AddError(
				"Error updating trust provider",
				fmt.Sprintf("Failed to create trust provider: %s", err),
			)
			return
		}

		cluster.TrustProvider = trustProvider
	} else {
		trustProviderKind := state.TrustProvider.Kind.ValueString()
		cluster.TrustProvider = &trustproviderpb.TrustProvider{
			Kind: &trustProviderKind,
		}
	}

	if !plan.ExtraHelmValues.IsNull() {
		parsedHelmValues, err := parseExtraHelmValues(plan.ExtraHelmValues)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error parsing extra_helm_values",
				fmt.Sprintf("Failed to parse extra_helm_values: %s", err),
			)

			return
		}

		cluster.ExtraHelmValues = parsedHelmValues
	}

	updateResp, err := c.client.ClusterV1Alpha1().UpdateCluster(ctx, cluster)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating cluster",
			err.Error(),
		)

		return
	}

	plan.ID = tftypes.StringValue(updateResp.GetId())
	plan.Name = tftypes.StringValue(updateResp.GetName())
	plan.OrgID = tftypes.StringValue(updateResp.GetOrgId())
	plan.TrustZoneID = tftypes.StringValue(updateResp.GetTrustZoneId())
	plan.KubernetesContext = tftypes.StringValue(updateResp.GetKubernetesContext())
	plan.Profile = tftypes.StringValue(updateResp.GetProfile())
	plan.ExternalServer = tftypes.BoolValue(updateResp.GetExternalServer())

	if updateResp.GetOidcIssuerUrl() != "" {
		plan.OidcIssuerURL = tftypes.StringValue(updateResp.GetOidcIssuerUrl())
	} else {
		plan.OidcIssuerURL = tftypes.StringNull()
	}

	if len(updateResp.GetOidcIssuerCaCert()) > 0 {
		plan.OidcIssuerCaCert = tftypes.StringValue(string(updateResp.GetOidcIssuerCaCert()))
	} else {
		plan.OidcIssuerCaCert = tftypes.StringNull()
	}

	plan.TrustProvider = &TrustProviderModel{
		Kind: tftypes.StringValue(updateResp.GetTrustProvider().GetKind()),
	}

	extraHelmValues := updateResp.GetExtraHelmValues()
	if err := validateHelmValues(plan.ExtraHelmValues, extraHelmValues); err != nil {
		resp.Diagnostics.AddError("Inconsistent Helm values", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
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

// newTrustProvider validates the trust provider kind and returns the corresponding trust provider.
func newTrustProvider(kind string) (*trustproviderpb.TrustProvider, error) {
	switch kind {
	case "kubernetes":
		return &trustproviderpb.TrustProvider{
			Kind: &kind,
		}, nil
	default:
		return nil, fmt.Errorf("invalid trust provider kind: %s", kind)
	}
}

// parseExtraHelmValues parses the extra_helm_values field from a string to a structpb.Struct object.
func parseExtraHelmValues(valueStr tftypes.String) (*structpb.Struct, error) {
	if valueStr.IsNull() || valueStr.ValueString() == "" {
		return nil, nil
	}

	var helmValues map[string]interface{}

	if err := yaml.Unmarshal([]byte(valueStr.ValueString()), &helmValues); err != nil {
		return nil, fmt.Errorf("invalid YAML in extra_helm_values: %w", err)
	}

	extraHelmValuesStruct, err := structpb.NewStruct(helmValues)
	if err != nil {
		return nil, fmt.Errorf("failed to convert extra_helm_values to Struct: %w", err)
	}

	return extraHelmValuesStruct, nil
}

// validateHelmValues compares the Helm values from the plan and the cluster API response.
func validateHelmValues(planValues tftypes.String, responseValues *structpb.Struct) error {
	if planValues.ValueString() == "" && responseValues == nil {
		return nil
	}

	if planValues.ValueString() == "" && responseValues != nil {
		return fmt.Errorf("invalid Helm values: plan is empty while the response is not")
	}

	if planValues.ValueString() != "" && responseValues == nil {
		return fmt.Errorf("invalid Helm values: plan is not empty while the response is nil")
	}

	var planValuesMap map[string]interface{}

	if err := yaml.Unmarshal([]byte(planValues.ValueString()), &planValuesMap); err != nil {
		return fmt.Errorf("error parsing plan extra_helm_values: %w", err)
	}

	responseValuesMap := responseValues.AsMap()

	planValuesJSON, err := json.Marshal(planValuesMap)
	if err != nil {
		return fmt.Errorf("error marshaling plan values: %w", err)
	}

	responseValuesJSON, err := json.Marshal(responseValuesMap)
	if err != nil {
		return fmt.Errorf("error marshaling response values: %w", err)
	}

	if string(planValuesJSON) != string(responseValuesJSON) {
		return fmt.Errorf("a Helm values mismatch: plan values don't match response values")
	}

	return nil
}
