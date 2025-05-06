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
}

type TrustProviderModel struct {
	Kind types.String `tfsdk:"kind"`
}
