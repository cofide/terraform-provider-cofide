package attestationpolicy

import tftypes "github.com/hashicorp/terraform-plugin-framework/types"

type AttestationPolicyModel struct {
	ID         tftypes.String     `tfsdk:"id"`
	Name       tftypes.String     `tfsdk:"name"`
	OrgID      tftypes.String     `tfsdk:"org_id"`
	Kubernetes *APKubernetesModel `tfsdk:"kubernetes"`
	Static     *APStaticModel     `tfsdk:"static"`
	TPMNode    *APTPMNodeModel    `tfsdk:"tpm_node"`
}

type APKubernetesModel struct {
	NamespaceSelector    *APLabelSelectorModel `tfsdk:"namespace_selector"`
	PodSelector          *APLabelSelectorModel `tfsdk:"pod_selector"`
	DnsNameTemplates     tftypes.List          `tfsdk:"dns_name_templates"`
	SpiffeIDPathTemplate tftypes.String        `tfsdk:"spiffe_id_path_template"`
}

type APLabelSelectorModel struct {
	MatchLabels      tftypes.Map              `tfsdk:"match_labels"`
	MatchExpressions []APMatchExpressionModel `tfsdk:"match_expressions"`
}

type APMatchExpressionModel struct {
	Key      tftypes.String `tfsdk:"key"`
	Operator tftypes.String `tfsdk:"operator"`
	Values   tftypes.List   `tfsdk:"values"`
}

type APStaticModel struct {
	SpiffeIDPath tftypes.String `tfsdk:"spiffe_id_path"`
	ParentIdPath tftypes.String `tfsdk:"parent_id_path"`
	Selectors    tftypes.List   `tfsdk:"selectors"`
	DNSNames     tftypes.List   `tfsdk:"dns_names"`
}

type APTPMNodeModel struct {
	Attestation    TPMAttestationModel `tfsdk:"attestation"`
	SelectorValues tftypes.List        `tfsdk:"selector_values"`
}

type TPMAttestationModel struct {
	EKHash tftypes.String `tfsdk:"ek_hash"`
}
