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
	SubjectAudience tftypes.List   `tfsdk:"subject_audience"`
	ClientID        tftypes.List   `tfsdk:"client_id"`
	TargetAudience  tftypes.List   `tfsdk:"target_audience"`
	OutboundScopes  tftypes.List   `tfsdk:"outbound_scopes"`
	ExternalHooks   tftypes.List   `tfsdk:"external_hooks"`
}

type StringMatcherModel struct {
	Exact tftypes.String `tfsdk:"exact"`
	Glob  tftypes.String `tfsdk:"glob"`
}

// ExternalHookModel represents a post-matching hook. The Auth field is a nested
// object whose attributes map one-to-one to proto oneof variants. Adding a new
// auth variant requires updating authAttrTypes and the convert functions but not
// this struct.
type ExternalHookModel struct {
	Name        tftypes.String `tfsdk:"name"`
	Description tftypes.String `tfsdk:"description"`
	URL         tftypes.String `tfsdk:"url"`
	Auth        tftypes.Object `tfsdk:"auth"`
	Timeout     tftypes.Int64  `tfsdk:"timeout"`
}

type ExchangePoliciesDataSourceModel struct {
	TrustZoneID      tftypes.String        `tfsdk:"trust_zone_id"`
	OrgID            tftypes.String        `tfsdk:"org_id"`
	Name             tftypes.String        `tfsdk:"name"`
	ExchangePolicies []ExchangePolicyModel `tfsdk:"exchange_policies"`
}
