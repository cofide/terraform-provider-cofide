package cluster

import "github.com/hashicorp/terraform-plugin-framework/types"

type ClusterModel struct {
	ID                types.String        `tfsdk:"id"`
	Name              types.String        `tfsdk:"name"`
	OrgID             types.String        `tfsdk:"org_id"`
	TrustZoneID       types.String        `tfsdk:"trust_zone_id"`
	KubernetesContext types.String        `tfsdk:"kubernetes_context"`
	TrustProvider     *TrustProviderModel `tfsdk:"trust_provider"`
	ExtraHelmValues   types.String        `tfsdk:"extra_helm_values"`
	Profile           types.String        `tfsdk:"profile"`
	ExternalServer    types.Bool          `tfsdk:"external_server"`
	OidcIssuerURL     types.String        `tfsdk:"oidc_issuer_url"`
	OidcIssuerCaCert  types.String        `tfsdk:"oidc_issuer_ca_cert"`
}

type TrustProviderModel struct {
	Kind          types.String        `tfsdk:"kind"`
	K8sPsatConfig *K8sPsatConfigModel `tfsdk:"k8s_psat_config"`
}

type K8sPsatConfigModel struct {
	Enabled                types.Bool           `tfsdk:"enabled"`
	AllowedServiceAccounts []ServiceAccountModel `tfsdk:"allowed_service_accounts"`
	AllowedNodeLabelKeys   []types.String        `tfsdk:"allowed_node_label_keys"`
	AllowedPodLabelKeys    []types.String        `tfsdk:"allowed_pod_label_keys"`
	ApiServerCaCert        types.String          `tfsdk:"api_server_ca_cert"`
	ApiServerURL           types.String          `tfsdk:"api_server_url"`
	ApiServerTLSServerName types.String          `tfsdk:"api_server_tls_server_name"`
	ApiServerProxyURL      types.String          `tfsdk:"api_server_proxy_url"`
	SpireServerAudience    types.String          `tfsdk:"spire_server_audience"`
}

type ServiceAccountModel struct {
	Namespace          types.String `tfsdk:"namespace"`
	ServiceAccountName types.String `tfsdk:"service_account_name"`
}
