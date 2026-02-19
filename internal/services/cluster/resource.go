package cluster

import (
	"context"
	"encoding/base64"
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
		TrustZoneId:       plan.TrustZoneID.ValueStringPointer(),
		KubernetesContext: plan.KubernetesContext.ValueStringPointer(),
		Profile:           plan.Profile.ValueStringPointer(),
		ExternalServer:    plan.ExternalServer.ValueBoolPointer(),
	}

	if !plan.OrgID.IsNull() {
		cluster.OrgId = plan.OrgID.ValueStringPointer()
	}

	if !plan.OidcIssuerURL.IsNull() {
		cluster.OidcIssuerUrl = plan.OidcIssuerURL.ValueStringPointer()
	}

	if !plan.OidcIssuerCaCert.IsNull() && plan.OidcIssuerCaCert.ValueString() != "" {
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

	var extraHelmValues tftypes.String
	if helmValues := createResp.GetExtraHelmValues(); helmValues != nil && len(helmValues.Fields) > 0 {
		jsonBytes, err := helmValues.MarshalJSON()
		if err != nil {
			resp.Diagnostics.AddError(
				"Error processing cluster data",
				fmt.Sprintf("Could not marshal extra_helm_values to JSON: %s", err),
			)
			return
		}
		extraHelmValues = tftypes.StringValue(string(jsonBytes))
	} else {
		extraHelmValues = tftypes.StringNull()
	}

	var oidcIssuerURL tftypes.String
	if url := createResp.GetOidcIssuerUrl(); url != "" {
		oidcIssuerURL = tftypes.StringValue(url)
	} else {
		oidcIssuerURL = tftypes.StringNull()
	}

	var oidcIssuerCaCert tftypes.String
	if certBytes := createResp.GetOidcIssuerCaCert(); len(certBytes) > 0 {
		oidcIssuerCaCert = tftypes.StringValue(base64.StdEncoding.EncodeToString(certBytes))
	} else {
		oidcIssuerCaCert = tftypes.StringNull()
	}

	state := ClusterModel{
		ID:                tftypes.StringValue(createResp.GetId()),
		Name:              tftypes.StringValue(createResp.GetName()),
		OrgID:             tftypes.StringValue(createResp.GetOrgId()),
		TrustZoneID:       tftypes.StringValue(createResp.GetTrustZoneId()),
		KubernetesContext: tftypes.StringValue(createResp.GetKubernetesContext()),
		TrustProvider: &TrustProviderModel{
			Kind: tftypes.StringValue(createResp.GetTrustProvider().GetKind()),
		},
		ExtraHelmValues:  extraHelmValues,
		Profile:          tftypes.StringValue(createResp.GetProfile()),
		ExternalServer:   tftypes.BoolValue(createResp.GetExternalServer()),
		OidcIssuerURL:    oidcIssuerURL,
		OidcIssuerCaCert: oidcIssuerCaCert,
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

	var extraHelmValues tftypes.String
	if helmValues := cluster.GetExtraHelmValues(); helmValues != nil && len(helmValues.Fields) > 0 {
		jsonBytes, err := helmValues.MarshalJSON()
		if err != nil {
			resp.Diagnostics.AddError(
				"Error processing cluster data",
				fmt.Sprintf("Could not marshal extra_helm_values to JSON: %s", err),
			)
			return
		}
		extraHelmValues = tftypes.StringValue(string(jsonBytes))
	} else {
		extraHelmValues = tftypes.StringNull()
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
		TrustProvider: &TrustProviderModel{
			Kind: tftypes.StringValue(cluster.GetTrustProvider().GetKind()),
		},
		ExtraHelmValues:  extraHelmValues,
		Profile:          tftypes.StringValue(cluster.GetProfile()),
		ExternalServer:   tftypes.BoolValue(cluster.GetExternalServer()),
		OidcIssuerURL:    oidcIssuerURL,
		OidcIssuerCaCert: oidcIssuerCaCert,
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

	if !plan.OrgID.IsNull() && plan.OrgID.ValueString() != "" {
		cluster.OrgId = plan.OrgID.ValueStringPointer()
	}

	if !plan.OidcIssuerURL.IsNull() && plan.OidcIssuerURL.ValueString() != "" {
		cluster.OidcIssuerUrl = plan.OidcIssuerURL.ValueStringPointer()
	}

	if !plan.OidcIssuerCaCert.IsNull() && plan.OidcIssuerCaCert.ValueString() != "" {
		decodedCert, err := base64.StdEncoding.DecodeString(plan.OidcIssuerCaCert.ValueString())
		if err != nil {
			resp.Diagnostics.AddError("Error decoding oidc_issuer_ca_cert", fmt.Sprintf("Failed to decode oidc_issuer_ca_cert from base64: %s", err))
			return
		}
		cluster.OidcIssuerCaCert = decodedCert
	}

	if plan.TrustProvider != nil && !plan.TrustProvider.Kind.IsNull() {
		trustProvider, err := newTrustProvider(plan.TrustProvider.Kind.ValueString())
		if err != nil {
			resp.Diagnostics.AddError("Error updating trust provider", fmt.Sprintf("Failed to create trust provider: %s", err))
			return
		}
		cluster.TrustProvider = trustProvider
	} else {
		trustProvider, _ := newTrustProvider(state.TrustProvider.Kind.ValueString())
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

	// Create new state object from response
	var extraHelmValuesStr tftypes.String
	if helmValues := updateResp.GetExtraHelmValues(); helmValues != nil && len(helmValues.Fields) > 0 {
		jsonBytes, err := helmValues.MarshalJSON()
		if err != nil {
			resp.Diagnostics.AddError("Error processing cluster data", fmt.Sprintf("Could not marshal extra_helm_values to JSON: %s", err))
			return
		}
		extraHelmValuesStr = tftypes.StringValue(string(jsonBytes))
	} else {
		extraHelmValuesStr = tftypes.StringNull()
	}

	var oidcIssuerURLStr tftypes.String
	if url := updateResp.GetOidcIssuerUrl(); url != "" {
		oidcIssuerURLStr = tftypes.StringValue(url)
	} else if !plan.OidcIssuerURL.IsNull() {
		oidcIssuerURLStr = plan.OidcIssuerURL
	} else {
		oidcIssuerURLStr = tftypes.StringNull()
	}

	var oidcIssuerCaCertStr tftypes.String
	if certBytes := updateResp.GetOidcIssuerCaCert(); len(certBytes) > 0 {
		oidcIssuerCaCertStr = tftypes.StringValue(base64.StdEncoding.EncodeToString(certBytes))
	} else if !plan.OidcIssuerCaCert.IsNull() {
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
		TrustProvider: &TrustProviderModel{
			Kind: tftypes.StringValue(updateResp.GetTrustProvider().GetKind()),
		},
		ExtraHelmValues:  extraHelmValuesStr,
		Profile:          tftypes.StringValue(updateResp.GetProfile()),
		ExternalServer:   tftypes.BoolValue(updateResp.GetExternalServer()),
		OidcIssuerURL:    oidcIssuerURLStr,
		OidcIssuerCaCert: oidcIssuerCaCertStr,
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
