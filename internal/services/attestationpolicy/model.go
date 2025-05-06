package attestationpolicy

import tftypes "github.com/hashicorp/terraform-plugin-framework/types"

type AttestationPolicyModel struct {
	ID         tftypes.String     `tfsdk:"id"`
	Name       tftypes.String     `tfsdk:"name"`
	OrgID      tftypes.String     `tfsdk:"org_id"`
	Kubernetes *APKubernetesModel `tfsdk:"kubernetes"`
	Static     *APStaticModel     `tfsdk:"static"`
}

type APKubernetesModel struct {
	NamespaceSelector *APLabelSelectorModel `tfsdk:"namespace_selector"`
	PodSelector       *APLabelSelectorModel `tfsdk:"pod_selector"`
}

type APLabelSelectorModel struct {
	MatchLabels      tftypes.Map              `tfsdk:"match_labels"`
	MatchExpressions []APMatchExpressionModel `tfsdk:"match_expressions"`
}

type APMatchExpressionModel struct {
	Key      tftypes.String   `tfsdk:"key"`
	Operator tftypes.String   `tfsdk:"operator"`
	Values   []tftypes.String `tfsdk:"values"`
}

type APStaticModel struct {
	SpiffeID  tftypes.String          `tfsdk:"spiffe_id"`
	Selectors []APStaticSelectorModel `tfsdk:"selectors"`
}

type APStaticSelectorModel struct {
	Type  tftypes.String `tfsdk:"type"`
	Value tftypes.String `tfsdk:"value"`
}
