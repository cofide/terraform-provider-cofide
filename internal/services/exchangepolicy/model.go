package exchangepolicy

import tftypes "github.com/hashicorp/terraform-plugin-framework/types"

type ExchangePolicyModel struct {
	ID              tftypes.String `tfsdk:"id"`
	OrgID           tftypes.String `tfsdk:"org_id"`
	Name            tftypes.String `tfsdk:"name"`
	TrustZoneID     tftypes.String `tfsdk:"trust_zone_id"`
	Action          tftypes.String `tfsdk:"action"`
	SubjectIdentity tftypes.List   `tfsdk:"subject_identity"`
	SubjectIssuer   tftypes.List   `tfsdk:"subject_issuer"`
	ActorIdentity   tftypes.List   `tfsdk:"actor_identity"`
	ActorIssuer     tftypes.List   `tfsdk:"actor_issuer"`
	ClientID        tftypes.List   `tfsdk:"client_id"`
	TargetAudience  tftypes.List   `tfsdk:"target_audience"`
	OutboundScopes  tftypes.List   `tfsdk:"outbound_scopes"`
}

type StringMatcherModel struct {
	Exact tftypes.String `tfsdk:"exact"`
	Glob  tftypes.String `tfsdk:"glob"`
}

type ExchangePoliciesDataSourceModel struct {
	TrustZoneID      tftypes.String        `tfsdk:"trust_zone_id"`
	OrgID            tftypes.String        `tfsdk:"org_id"`
	Name             tftypes.String        `tfsdk:"name"`
	ExchangePolicies []ExchangePolicyModel `tfsdk:"exchange_policies"`
}
