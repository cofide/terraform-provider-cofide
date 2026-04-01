package cluster

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"

	clusterpb "github.com/cofide/cofide-api-sdk/gen/go/proto/cluster/v1alpha1"
	trustproviderpb "github.com/cofide/cofide-api-sdk/gen/go/proto/trust_provider/v1alpha1"
	sdkclient "github.com/cofide/cofide-api-sdk/pkg/connect/client"
	"github.com/cofide/terraform-provider-cofide/internal/util"
	"github.com/hashicorp/terraform-plugin-framework/attr"
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
		TrustZoneId:       plan.TrustZoneID.ValueStringPointer(),
		KubernetesContext: plan.KubernetesContext.ValueStringPointer(),
		Profile:           plan.Profile.ValueStringPointer(),
		ExternalServer:    plan.ExternalServer.ValueBoolPointer(),
	}

	if util.IsStringAttributeNonEmpty(plan.OrgID) {
		cluster.OrgId = plan.OrgID.ValueStringPointer()
	}

	if util.IsStringAttributeNonEmpty(plan.OidcIssuerURL) {
		cluster.OidcIssuerUrl = plan.OidcIssuerURL.ValueStringPointer()
	}

	if util.IsStringAttributeNonEmpty(plan.OidcIssuerCaCert) {
		decodedCert, err := base64.StdEncoding.DecodeString(plan.OidcIssuerCaCert.ValueString())
		if err != nil {
			resp.Diagnostics.AddError(
				"Error decoding oidc_issuer_ca_cert",
				fmt.Sprintf("Failed to decode oidc_issuer_ca_cert from base64: %s", err),
			)

			return
		}
		cluster.OidcIssuerCaCert = decodedCert
	}

	trustProvider, err := trustProviderToProto(ctx, plan.TrustProvider)
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

	// Store the plan value for extra_helm_values to preserve the original format
	// (YAML or JSON) and avoid a plan/state inconsistency on apply.
	extraHelmValues := plan.ExtraHelmValues

	var oidcIssuerURL tftypes.String
	if url := createResp.GetOidcIssuerUrl(); url != "" {
		oidcIssuerURL = tftypes.StringValue(url)
	} else {
		oidcIssuerURL = tftypes.StringNull()
	}

	var oidcIssuerCaCert tftypes.String
	if certBytes := createResp.GetOidcIssuerCaCert(); len(certBytes) > 0 {
		oidcIssuerCaCert = tftypes.StringValue(base64.StdEncoding.EncodeToString(certBytes))
	} else if !plan.OidcIssuerCaCert.IsNull() && !plan.OidcIssuerCaCert.IsUnknown() {
		oidcIssuerCaCert = plan.OidcIssuerCaCert
	} else {
		oidcIssuerCaCert = tftypes.StringNull()
	}

	state := ClusterModel{
		ID:                tftypes.StringValue(createResp.GetId()),
		Name:              tftypes.StringValue(createResp.GetName()),
		OrgID:             tftypes.StringValue(createResp.GetOrgId()),
		TrustZoneID:       tftypes.StringValue(createResp.GetTrustZoneId()),
		KubernetesContext: tftypes.StringValue(createResp.GetKubernetesContext()),
		TrustProvider:     trustProviderForState(createResp.GetTrustProvider(), plan.TrustProvider),
		ExtraHelmValues:   extraHelmValues,
		Profile:           tftypes.StringValue(createResp.GetProfile()),
		ExternalServer:    tftypes.BoolValue(createResp.GetExternalServer()),
		OidcIssuerURL:     oidcIssuerURL,
		OidcIssuerCaCert:  oidcIssuerCaCert,
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (c *ClusterResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state ClusterModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	clusterID := state.ID.ValueString()
	cluster, err := c.client.ClusterV1Alpha1().GetCluster(ctx, clusterID)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error reading cluster",
			fmt.Sprintf("Could not read cluster %q: %s", clusterID, err),
		)
		return
	}

	extraHelmValues, err := helmValuesForState(cluster.GetExtraHelmValues(), state.ExtraHelmValues)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error processing cluster data",
			fmt.Sprintf("Could not process extra_helm_values: %s", err),
		)
		return
	}

	var oidcIssuerURL tftypes.String
	if url := cluster.GetOidcIssuerUrl(); url != "" {
		oidcIssuerURL = tftypes.StringValue(url)
	} else {
		oidcIssuerURL = tftypes.StringNull()
	}

	var oidcIssuerCaCert tftypes.String
	if certBytes := cluster.GetOidcIssuerCaCert(); len(certBytes) > 0 {
		oidcIssuerCaCert = tftypes.StringValue(base64.StdEncoding.EncodeToString(certBytes))
	} else {
		oidcIssuerCaCert = tftypes.StringNull()
	}

	newState := ClusterModel{
		ID:                tftypes.StringValue(cluster.GetId()),
		Name:              tftypes.StringValue(cluster.GetName()),
		OrgID:             tftypes.StringValue(cluster.GetOrgId()),
		TrustZoneID:       tftypes.StringValue(cluster.GetTrustZoneId()),
		KubernetesContext: tftypes.StringValue(cluster.GetKubernetesContext()),
		TrustProvider:     trustProviderForState(cluster.GetTrustProvider(), state.TrustProvider),
		ExtraHelmValues:   extraHelmValues,
		Profile:           tftypes.StringValue(cluster.GetProfile()),
		ExternalServer:    tftypes.BoolValue(cluster.GetExternalServer()),
		OidcIssuerURL:     oidcIssuerURL,
		OidcIssuerCaCert:  oidcIssuerCaCert,
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &newState)...)
}

func (c *ClusterResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan ClusterModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state ClusterModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	clusterID := state.ID.ValueString()

	cluster := &clusterpb.Cluster{
		Id:                &clusterID,
		Name:              plan.Name.ValueStringPointer(),
		TrustZoneId:       plan.TrustZoneID.ValueStringPointer(),
		KubernetesContext: plan.KubernetesContext.ValueStringPointer(),
		Profile:           plan.Profile.ValueStringPointer(),
		ExternalServer:    plan.ExternalServer.ValueBoolPointer(),
	}

	if util.IsStringAttributeNonEmpty(plan.OrgID) {
		cluster.OrgId = plan.OrgID.ValueStringPointer()
	}

	if util.IsStringAttributeNonEmpty(plan.OidcIssuerURL) {
		cluster.OidcIssuerUrl = plan.OidcIssuerURL.ValueStringPointer()
	}

	if util.IsStringAttributeNonEmpty(plan.OidcIssuerCaCert) {
		decodedCert, err := base64.StdEncoding.DecodeString(plan.OidcIssuerCaCert.ValueString())
		if err != nil {
			resp.Diagnostics.AddError("Error decoding oidc_issuer_ca_cert", fmt.Sprintf("Failed to decode oidc_issuer_ca_cert from base64: %s", err))
			return
		}
		cluster.OidcIssuerCaCert = decodedCert
	}

	if plan.TrustProvider != nil && !plan.TrustProvider.Kind.IsNull() {
		trustProvider, err := trustProviderToProto(ctx, plan.TrustProvider)
		if err != nil {
			resp.Diagnostics.AddError("Error updating trust provider", fmt.Sprintf("Failed to create trust provider: %s", err))
			return
		}
		cluster.TrustProvider = trustProvider
	} else {
		trustProvider, _ := trustProviderToProto(ctx, state.TrustProvider)
		cluster.TrustProvider = trustProvider
	}

	if !plan.ExtraHelmValues.IsNull() {
		parsedHelmValues, err := parseExtraHelmValues(plan.ExtraHelmValues)
		if err != nil {
			resp.Diagnostics.AddError("Error parsing extra_helm_values", fmt.Sprintf("Failed to parse extra_helm_values: %s", err))
			return
		}
		cluster.ExtraHelmValues = parsedHelmValues
	}

	updateResp, err := c.client.ClusterV1Alpha1().UpdateCluster(ctx, cluster)
	if err != nil {
		resp.Diagnostics.AddError("Error updating cluster", err.Error())
		return
	}

	// Store the plan value for extra_helm_values to preserve the original format
	// (YAML or JSON) and avoid a plan/state inconsistency on apply.
	extraHelmValuesStr := plan.ExtraHelmValues

	var oidcIssuerURLStr tftypes.String
	if url := updateResp.GetOidcIssuerUrl(); url != "" {
		oidcIssuerURLStr = tftypes.StringValue(url)
	} else if !plan.OidcIssuerURL.IsNull() && !plan.OidcIssuerURL.IsUnknown() {
		oidcIssuerURLStr = plan.OidcIssuerURL
	} else {
		oidcIssuerURLStr = tftypes.StringNull()
	}

	var oidcIssuerCaCertStr tftypes.String
	if certBytes := updateResp.GetOidcIssuerCaCert(); len(certBytes) > 0 {
		oidcIssuerCaCertStr = tftypes.StringValue(base64.StdEncoding.EncodeToString(certBytes))
	} else if !plan.OidcIssuerCaCert.IsNull() && !plan.OidcIssuerCaCert.IsUnknown() {
		oidcIssuerCaCertStr = plan.OidcIssuerCaCert
	} else {
		oidcIssuerCaCertStr = tftypes.StringNull()
	}

	newState := ClusterModel{
		ID:                tftypes.StringValue(updateResp.GetId()),
		Name:              tftypes.StringValue(updateResp.GetName()),
		OrgID:             tftypes.StringValue(updateResp.GetOrgId()),
		TrustZoneID:       tftypes.StringValue(updateResp.GetTrustZoneId()),
		KubernetesContext: tftypes.StringValue(updateResp.GetKubernetesContext()),
		TrustProvider:     trustProviderForState(updateResp.GetTrustProvider(), plan.TrustProvider),
		ExtraHelmValues:   extraHelmValuesStr,
		Profile:           tftypes.StringValue(updateResp.GetProfile()),
		ExternalServer:    tftypes.BoolValue(updateResp.GetExternalServer()),
		OidcIssuerURL:     oidcIssuerURLStr,
		OidcIssuerCaCert:  oidcIssuerCaCertStr,
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &newState)...)
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

// trustProviderToProto converts a TrustProviderModel to a TrustProvider proto message.
func trustProviderToProto(ctx context.Context, model *TrustProviderModel) (*trustproviderpb.TrustProvider, error) {
	tp, err := newTrustProvider(model.Kind.ValueString())
	if err != nil {
		return nil, err
	}
	if model.K8sPsatConfig != nil {
		cfg, err := k8sPsatConfigToProto(ctx, model.K8sPsatConfig)
		if err != nil {
			return nil, err
		}
		tp.K8SPsatConfig = cfg
	}
	return tp, nil
}

// k8sPsatConfigToProto converts a K8sPsatConfigModel to a K8SPsatConfig proto message.
func k8sPsatConfigToProto(ctx context.Context, model *K8sPsatConfigModel) (*trustproviderpb.K8SPsatConfig, error) {
	cfg := &trustproviderpb.K8SPsatConfig{
		Enabled: model.Enabled.ValueBool(),
	}

	for _, sa := range model.AllowedServiceAccounts {
		cfg.AllowedServiceAccounts = append(cfg.AllowedServiceAccounts, &trustproviderpb.K8SPsatConfig_ServiceAccount{
			Namespace:          sa.Namespace.ValueString(),
			ServiceAccountName: sa.ServiceAccountName.ValueString(),
		})
	}

	var allowedNodeLabelKeys []string
	if diags := model.AllowedNodeLabelKeys.ElementsAs(ctx, &allowedNodeLabelKeys, false); diags.HasError() {
		return nil, fmt.Errorf("error reading allowed_node_label_keys")
	}
	cfg.AllowedNodeLabelKeys = allowedNodeLabelKeys

	var allowedPodLabelKeys []string
	if diags := model.AllowedPodLabelKeys.ElementsAs(ctx, &allowedPodLabelKeys, false); diags.HasError() {
		return nil, fmt.Errorf("error reading allowed_pod_label_keys")
	}
	cfg.AllowedPodLabelKeys = allowedPodLabelKeys

	if !model.ApiServerCaCert.IsNull() && model.ApiServerCaCert.ValueString() != "" {
		decoded, err := base64.StdEncoding.DecodeString(model.ApiServerCaCert.ValueString())
		if err != nil {
			return nil, fmt.Errorf("failed to decode api_server_ca_cert from base64: %w", err)
		}
		cfg.ApiServerCaCert = decoded
	}

	if !model.ApiServerURL.IsNull() {
		cfg.ApiServerUrl = model.ApiServerURL.ValueString()
	}
	if !model.ApiServerTLSServerName.IsNull() {
		cfg.ApiServerTlsServerName = model.ApiServerTLSServerName.ValueString()
	}
	if !model.ApiServerProxyURL.IsNull() {
		cfg.ApiServerProxyUrl = model.ApiServerProxyURL.ValueString()
	}
	if !model.SpireServerAudience.IsNull() {
		cfg.SpireServerAudience = model.SpireServerAudience.ValueString()
	}

	return cfg, nil
}

// trustProviderForState returns a TrustProviderModel for storage in state, using prev to
// preserve user intent for fields where nil and empty are semantically equivalent in proto3.
func trustProviderForState(apiTp *trustproviderpb.TrustProvider, prev *TrustProviderModel) *TrustProviderModel {
	model := trustProviderFromProto(apiTp)
	if prev == nil || prev.K8sPsatConfig == nil || model.K8sPsatConfig == nil {
		return model
	}
	model.K8sPsatConfig = k8sPsatConfigForState(model.K8sPsatConfig, prev.K8sPsatConfig)
	return model
}

// k8sPsatConfigForState returns a K8sPsatConfigModel for storage in state.
// Proto3 does not distinguish nil from empty slice, but Terraform does. If the API
// returns nil/empty for a list field and prev has an explicit empty list, the empty
// list is preserved to avoid spurious plan diffs.
func k8sPsatConfigForState(model, prev *K8sPsatConfigModel) *K8sPsatConfigModel {
	if len(model.AllowedServiceAccounts) == 0 && len(prev.AllowedServiceAccounts) == 0 && prev.AllowedServiceAccounts != nil {
		model.AllowedServiceAccounts = prev.AllowedServiceAccounts
	}
	if model.AllowedNodeLabelKeys.IsNull() && !prev.AllowedNodeLabelKeys.IsNull() && len(prev.AllowedNodeLabelKeys.Elements()) == 0 {
		model.AllowedNodeLabelKeys = prev.AllowedNodeLabelKeys
	}
	if model.AllowedPodLabelKeys.IsNull() && !prev.AllowedPodLabelKeys.IsNull() && len(prev.AllowedPodLabelKeys.Elements()) == 0 {
		model.AllowedPodLabelKeys = prev.AllowedPodLabelKeys
	}
	return model
}

// trustProviderFromProto converts a TrustProvider proto message to a TrustProviderModel.
func trustProviderFromProto(tp *trustproviderpb.TrustProvider) *TrustProviderModel {
	model := &TrustProviderModel{
		Kind: tftypes.StringValue(tp.GetKind()),
	}
	if cfg := tp.GetK8SPsatConfig(); cfg != nil {
		model.K8sPsatConfig = k8sPsatConfigFromProto(cfg)
	}
	return model
}

// k8sPsatConfigFromProto converts a K8SPsatConfig proto message to a K8sPsatConfigModel.
func k8sPsatConfigFromProto(cfg *trustproviderpb.K8SPsatConfig) *K8sPsatConfigModel {
	model := &K8sPsatConfigModel{
		Enabled: tftypes.BoolValue(cfg.Enabled),
	}

	for _, sa := range cfg.AllowedServiceAccounts {
		model.AllowedServiceAccounts = append(model.AllowedServiceAccounts, ServiceAccountModel{
			Namespace:          tftypes.StringValue(sa.Namespace),
			ServiceAccountName: tftypes.StringValue(sa.ServiceAccountName),
		})
	}

	if len(cfg.AllowedNodeLabelKeys) > 0 {
		nodeLabelKeyElems := make([]attr.Value, 0, len(cfg.AllowedNodeLabelKeys))
		for _, key := range cfg.AllowedNodeLabelKeys {
			nodeLabelKeyElems = append(nodeLabelKeyElems, tftypes.StringValue(key))
		}
		model.AllowedNodeLabelKeys = tftypes.ListValueMust(tftypes.StringType, nodeLabelKeyElems)
	} else {
		model.AllowedNodeLabelKeys = tftypes.ListNull(tftypes.StringType)
	}

	if len(cfg.AllowedPodLabelKeys) > 0 {
		podLabelKeyElems := make([]attr.Value, 0, len(cfg.AllowedPodLabelKeys))
		for _, key := range cfg.AllowedPodLabelKeys {
			podLabelKeyElems = append(podLabelKeyElems, tftypes.StringValue(key))
		}
		model.AllowedPodLabelKeys = tftypes.ListValueMust(tftypes.StringType, podLabelKeyElems)
	} else {
		model.AllowedPodLabelKeys = tftypes.ListNull(tftypes.StringType)
	}

	if len(cfg.ApiServerCaCert) > 0 {
		model.ApiServerCaCert = tftypes.StringValue(base64.StdEncoding.EncodeToString(cfg.ApiServerCaCert))
	} else {
		model.ApiServerCaCert = tftypes.StringNull()
	}

	if cfg.ApiServerUrl != "" {
		model.ApiServerURL = tftypes.StringValue(cfg.ApiServerUrl)
	} else {
		model.ApiServerURL = tftypes.StringNull()
	}
	if cfg.ApiServerTlsServerName != "" {
		model.ApiServerTLSServerName = tftypes.StringValue(cfg.ApiServerTlsServerName)
	} else {
		model.ApiServerTLSServerName = tftypes.StringNull()
	}
	if cfg.ApiServerProxyUrl != "" {
		model.ApiServerProxyURL = tftypes.StringValue(cfg.ApiServerProxyUrl)
	} else {
		model.ApiServerProxyURL = tftypes.StringNull()
	}
	if cfg.SpireServerAudience != "" {
		model.SpireServerAudience = tftypes.StringValue(cfg.SpireServerAudience)
	} else {
		model.SpireServerAudience = tftypes.StringNull()
	}

	return model
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
