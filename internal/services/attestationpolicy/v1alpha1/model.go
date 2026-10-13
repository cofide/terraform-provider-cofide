package v1alpha1

import tftypes "github.com/hashicorp/terraform-plugin-framework/types"

// AttestationPolicyModel is the Terraform model of an attestation policy.
//
// Every attribute that can be set in configuration is a framework type (tftypes.List, tftypes.Object, ...), not a Go
// slice, struct or pointer, so that it can hold an unknown value. Terraform sends unknown values while validating
// and planning for anything not known until apply, e.g. a value built from a module variable. The structs below are
// the shapes of those objects and list elements, for converting them with As and ElementsAs.
type AttestationPolicyModel struct {
	ID          tftypes.String `tfsdk:"id"`
	Name        tftypes.String `tfsdk:"name"`
	OrgID       tftypes.String `tfsdk:"org_id"`
	TrustZoneID tftypes.String `tfsdk:"trust_zone_id"`
	// Federations is a list of APFederationModel.
	Federations tftypes.List `tfsdk:"federations"`
	// Kubernetes is an APKubernetesModel.
	Kubernetes tftypes.Object `tfsdk:"kubernetes"`
	// Static is an APStaticModel.
	Static tftypes.Object `tfsdk:"static"`
	// TPMNode is an APTPMNodeModel.
	TPMNode tftypes.Object `tfsdk:"tpm_node"`
}

type APFederationModel struct {
	TrustZoneID tftypes.String `tfsdk:"trust_zone_id"`
}

type APKubernetesModel struct {
	// NamespaceSelector is an APLabelSelectorModel.
	NamespaceSelector tftypes.Object `tfsdk:"namespace_selector"`
	// PodSelector is an APLabelSelectorModel.
	PodSelector          tftypes.Object `tfsdk:"pod_selector"`
	DnsNameTemplates     tftypes.List   `tfsdk:"dns_name_templates"`
	SpiffeIDPathTemplate tftypes.String `tfsdk:"spiffe_id_path_template"`
}

type APLabelSelectorModel struct {
	MatchLabels tftypes.Map `tfsdk:"match_labels"`
	// MatchExpressions is a list of APMatchExpressionModel.
	MatchExpressions tftypes.List `tfsdk:"match_expressions"`
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
	StoreSvid    tftypes.Bool   `tfsdk:"store_svid"`
}

type APTPMNodeModel struct {
	// Attestation is a TPMAttestationModel.
	Attestation    tftypes.Object `tfsdk:"attestation"`
	SelectorValues tftypes.List   `tfsdk:"selector_values"`
}

type TPMAttestationModel struct {
	EKHash tftypes.String `tfsdk:"ek_hash"`
}
