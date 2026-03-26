package trustzoneserver

import "github.com/hashicorp/terraform-plugin-framework/types"

type TrustZoneServerModel struct {
	ID                       types.String               `tfsdk:"id"`
	TrustZoneID              types.String               `tfsdk:"trust_zone_id"`
	ClusterID                types.String               `tfsdk:"cluster_id"`
	KubernetesNamespace      types.String               `tfsdk:"kubernetes_namespace"`
	KubernetesServiceAccount types.String               `tfsdk:"kubernetes_service_account"`
	OrgID                    types.String               `tfsdk:"org_id"`
	HelmValues               types.String               `tfsdk:"helm_values"`
	Status                   types.Object               `tfsdk:"status"`
	ConnectK8sPsatConfig     *ConnectK8sPsatConfigModel `tfsdk:"connect_k8s_psat_config"`
}

type TrustZoneServerStatusModel struct {
	Status             types.String `tfsdk:"status"`
	LastTransitionTime types.String `tfsdk:"last_transition_time"`
}

type ConnectK8sPsatConfigModel struct {
	Audiences               []types.String `tfsdk:"audiences"`
	SpireServerSpiffeIDPath types.String   `tfsdk:"spire_server_spiffe_id_path"`
}

type TrustZoneServersDataSourceModel struct {
	TrustZoneID      types.String           `tfsdk:"trust_zone_id"`
	ClusterID        types.String           `tfsdk:"cluster_id"`
	OrgID            types.String           `tfsdk:"org_id"`
	TrustZoneServers []TrustZoneServerModel `tfsdk:"trust_zone_servers"`
}
