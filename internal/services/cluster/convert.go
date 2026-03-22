package cluster

import (
	"encoding/base64"
	"encoding/json"
	"fmt"

	clusterpb "github.com/cofide/cofide-api-sdk/gen/go/proto/cluster/v1alpha1"
	trustproviderpb "github.com/cofide/cofide-api-sdk/gen/go/proto/trust_provider/v1alpha1"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"google.golang.org/protobuf/types/known/structpb"
	"gopkg.in/yaml.v3"
)

func modelToProto(plan ClusterModel) (*clusterpb.Cluster, error) {
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
		decodedCert, err := base64.StdEncoding.DecodeString(plan.OidcIssuerCaCert.ValueString())
		if err != nil {
			return nil, fmt.Errorf("failed to decode oidc_issuer_ca_cert from base64: %w", err)
		}
		cluster.OidcIssuerCaCert = decodedCert
	}

	trustProvider, err := newTrustProvider(plan.TrustProvider.Kind.ValueString())
	if err != nil {
		return nil, fmt.Errorf("failed to create trust provider: %w", err)
	}

	cluster.TrustProvider = trustProvider

	if !plan.ExtraHelmValues.IsNull() {
		parsedHelmValues, err := parseExtraHelmValues(plan.ExtraHelmValues)
		if err != nil {
			return nil, fmt.Errorf("failed to parse extra_helm_values: %w", err)
		}

		cluster.ExtraHelmValues = parsedHelmValues
	}

	return cluster, nil
}

func protoToModel(cluster *clusterpb.Cluster) (ClusterModel, error) {
	var extraHelmValues types.String
	if helmValues := cluster.GetExtraHelmValues(); helmValues != nil && len(helmValues.Fields) > 0 {
		jsonBytes, err := helmValues.MarshalJSON()
		if err != nil {
			return ClusterModel{}, fmt.Errorf("could not marshal extra_helm_values to JSON: %w", err)
		}
		extraHelmValues = types.StringValue(string(jsonBytes))
	} else {
		extraHelmValues = types.StringNull()
	}

	return ClusterModel{
		ID:                types.StringValue(cluster.GetId()),
		Name:              types.StringValue(cluster.GetName()),
		OrgID:             types.StringValue(cluster.GetOrgId()),
		TrustZoneID:       types.StringValue(cluster.GetTrustZoneId()),
		KubernetesContext: types.StringValue(cluster.GetKubernetesContext()),
		TrustProvider: &TrustProviderModel{
			Kind: types.StringValue(cluster.GetTrustProvider().GetKind()),
		},
		ExtraHelmValues:  extraHelmValues,
		Profile:          types.StringValue(cluster.GetProfile()),
		ExternalServer:   types.BoolValue(cluster.GetExternalServer()),
		OidcIssuerURL:    types.StringValue(cluster.GetOidcIssuerUrl()),
		OidcIssuerCaCert: types.StringValue(base64.StdEncoding.EncodeToString(cluster.GetOidcIssuerCaCert())),
	}, nil
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
func parseExtraHelmValues(valueStr types.String) (*structpb.Struct, error) {
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
func validateHelmValues(planValues types.String, responseValues *structpb.Struct) error {
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
