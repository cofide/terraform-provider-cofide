package exchangepolicy

import (
	"context"
	"testing"

	exchangepolicypb "github.com/cofide/cofide-api-sdk/gen/go/proto/exchange_policy/v1alpha1"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProtoToModel_Minimal(t *testing.T) {
	proto := &exchangepolicypb.ExchangePolicy{
		Id:          "ep-1",
		OrgId:       "org-1",
		Name:        "test-policy",
		TrustZoneId: "tz-1",
	}

	got := protoToModel(proto)

	assert.Equal(t, types.StringValue("ep-1"), got.ID)
	assert.Equal(t, types.StringValue("org-1"), got.OrgID)
	assert.Equal(t, types.StringValue("test-policy"), got.Name)
	assert.Equal(t, types.StringValue("tz-1"), got.TrustZoneID)
	assert.True(t, got.Action.IsNull())
	assert.Equal(t, types.ListNull(stringMatcherObjectType), got.SubjectIdentity)
	assert.Equal(t, types.ListNull(stringMatcherObjectType), got.SubjectIssuer)
	assert.Equal(t, types.ListNull(stringMatcherObjectType), got.ActorIdentity)
	assert.Equal(t, types.ListNull(stringMatcherObjectType), got.ActorIssuer)
	assert.Equal(t, types.ListNull(stringMatcherObjectType), got.ClientID)
	assert.Equal(t, types.ListNull(stringMatcherObjectType), got.TargetAudience)
	assert.Equal(t, types.ListValueMust(types.StringType, []attr.Value{}), got.OutboundScopes)
}

func TestProtoToModel_Full(t *testing.T) {
	action := exchangepolicypb.ExchangePolicyAction_EXCHANGE_POLICY_ACTION_ALLOW
	proto := &exchangepolicypb.ExchangePolicy{
		Id:          "ep-2",
		OrgId:       "org-2",
		Name:        "full-policy",
		TrustZoneId: "tz-2",
		Action:      &action,
		SubjectIdentity: &exchangepolicypb.StringSet{
			Matchers: []*exchangepolicypb.StringMatcher{
				{Match: &exchangepolicypb.StringMatcher_Exact{Exact: "spiffe://example.org/workload"}},
			},
		},
		SubjectIssuer: &exchangepolicypb.StringSet{
			Matchers: []*exchangepolicypb.StringMatcher{
				{Match: &exchangepolicypb.StringMatcher_Glob{Glob: "spiffe://example.org/*"}},
			},
		},
		OutboundScopes: []string{"read", "write"},
	}

	got := protoToModel(proto)

	assert.Equal(t, types.StringValue("ALLOW"), got.Action)

	wantSubjectIdentity := types.ListValueMust(stringMatcherObjectType, []attr.Value{
		types.ObjectValueMust(stringMatcherAttrTypes, map[string]attr.Value{
			"exact": types.StringValue("spiffe://example.org/workload"),
			"glob":  types.StringNull(),
		}),
	})
	assert.Equal(t, wantSubjectIdentity, got.SubjectIdentity)

	wantSubjectIssuer := types.ListValueMust(stringMatcherObjectType, []attr.Value{
		types.ObjectValueMust(stringMatcherAttrTypes, map[string]attr.Value{
			"exact": types.StringNull(),
			"glob":  types.StringValue("spiffe://example.org/*"),
		}),
	})
	assert.Equal(t, wantSubjectIssuer, got.SubjectIssuer)

	wantScopes := types.ListValueMust(types.StringType, []attr.Value{
		types.StringValue("read"),
		types.StringValue("write"),
	})
	assert.Equal(t, wantScopes, got.OutboundScopes)
}

func TestProtoToModel_DenyAction(t *testing.T) {
	action := exchangepolicypb.ExchangePolicyAction_EXCHANGE_POLICY_ACTION_DENY
	proto := &exchangepolicypb.ExchangePolicy{
		Id:     "ep-3",
		Name:   "deny-policy",
		Action: &action,
	}

	got := protoToModel(proto)
	assert.Equal(t, types.StringValue("DENY"), got.Action)
}

func TestModelToProto_Minimal(t *testing.T) {
	model := ExchangePolicyModel{
		ID:          types.StringValue("ep-1"),
		OrgID:       types.StringValue("org-1"),
		Name:        types.StringValue("test-policy"),
		TrustZoneID: types.StringValue("tz-1"),
		Action:      types.StringNull(),
	}

	got, err := modelToProto(context.Background(), model)

	require.NoError(t, err)
	assert.Equal(t, "ep-1", got.GetId())
	assert.Equal(t, "test-policy", got.GetName())
	assert.Equal(t, "tz-1", got.GetTrustZoneId())
	assert.Nil(t, got.Action)
	assert.Nil(t, got.SubjectIdentity)
	assert.Empty(t, got.OutboundScopes)
}

func TestModelToProto_WithAction(t *testing.T) {
	model := ExchangePolicyModel{
		Name:        types.StringValue("allow-policy"),
		TrustZoneID: types.StringValue("tz-1"),
		Action:      types.StringValue("ALLOW"),
	}

	got, err := modelToProto(context.Background(), model)

	require.NoError(t, err)
	require.NotNil(t, got.Action)
	assert.Equal(t, exchangepolicypb.ExchangePolicyAction_EXCHANGE_POLICY_ACTION_ALLOW, *got.Action)
}

func TestModelToProto_WithDenyAction(t *testing.T) {
	model := ExchangePolicyModel{
		Name:        types.StringValue("deny-policy"),
		TrustZoneID: types.StringValue("tz-1"),
		Action:      types.StringValue("DENY"),
	}

	got, err := modelToProto(context.Background(), model)

	require.NoError(t, err)
	require.NotNil(t, got.Action)
	assert.Equal(t, exchangepolicypb.ExchangePolicyAction_EXCHANGE_POLICY_ACTION_DENY, *got.Action)
}

func TestModelToProto_InvalidAction(t *testing.T) {
	model := ExchangePolicyModel{
		Name:        types.StringValue("bad-policy"),
		TrustZoneID: types.StringValue("tz-1"),
		Action:      types.StringValue("INVALID"),
	}

	got, err := modelToProto(context.Background(), model)

	assert.Nil(t, got)
	assert.ErrorContains(t, err, "invalid action")
}

func TestModelToProto_WithStringSetMatchers(t *testing.T) {
	model := ExchangePolicyModel{
		Name:        types.StringValue("policy"),
		TrustZoneID: types.StringValue("tz-1"),
		Action:      types.StringNull(),
		SubjectIdentity: types.ListValueMust(stringMatcherObjectType, []attr.Value{
			types.ObjectValueMust(stringMatcherAttrTypes, map[string]attr.Value{
				"exact": types.StringValue("spiffe://example.org/workload"),
				"glob":  types.StringNull(),
			}),
			types.ObjectValueMust(stringMatcherAttrTypes, map[string]attr.Value{
				"exact": types.StringNull(),
				"glob":  types.StringValue("spiffe://example.org/*"),
			}),
		}),
		OutboundScopes: types.ListValueMust(types.StringType, []attr.Value{
			types.StringValue("openid"),
			types.StringValue("profile"),
		}),
	}

	got, err := modelToProto(context.Background(), model)

	require.NoError(t, err)
	require.NotNil(t, got.SubjectIdentity)
	require.Len(t, got.SubjectIdentity.Matchers, 2)
	assert.Equal(t, "spiffe://example.org/workload", got.SubjectIdentity.Matchers[0].GetExact())
	assert.Equal(t, "spiffe://example.org/*", got.SubjectIdentity.Matchers[1].GetGlob())
	assert.Equal(t, []string{"openid", "profile"}, got.OutboundScopes)
}

// TestRoundTrip verifies that model→proto→model conversion is lossless for all
// fields that are round-tripped. OrgID is a server-computed field that modelToProto
// does not serialize, so round-trip models use OrgID: StringValue("").
func TestRoundTrip(t *testing.T) {
	nullMatchers := types.ListNull(stringMatcherObjectType)
	tests := []struct {
		name  string
		model ExchangePolicyModel
	}{
		{
			name: "minimal policy",
			model: ExchangePolicyModel{
				ID:              types.StringValue("ep-1"),
				OrgID:           types.StringValue(""),
				Name:            types.StringValue("minimal"),
				TrustZoneID:     types.StringValue("tz-1"),
				Action:          types.StringNull(),
				SubjectIdentity: nullMatchers,
				SubjectIssuer:   nullMatchers,
				ActorIdentity:   nullMatchers,
				ActorIssuer:     nullMatchers,
				ClientID:        nullMatchers,
				TargetAudience:  nullMatchers,
				OutboundScopes:  types.ListValueMust(types.StringType, []attr.Value{}),
			},
		},
		{
			name: "allow policy with all string sets",
			model: ExchangePolicyModel{
				ID:          types.StringValue("ep-2"),
				OrgID:       types.StringValue(""),
				Name:        types.StringValue("full-allow"),
				TrustZoneID: types.StringValue("tz-2"),
				Action:      types.StringValue("ALLOW"),
				SubjectIdentity: types.ListValueMust(stringMatcherObjectType, []attr.Value{
					types.ObjectValueMust(stringMatcherAttrTypes, map[string]attr.Value{
						"exact": types.StringValue("spiffe://example.org/subject"),
						"glob":  types.StringNull(),
					}),
				}),
				SubjectIssuer: types.ListValueMust(stringMatcherObjectType, []attr.Value{
					types.ObjectValueMust(stringMatcherAttrTypes, map[string]attr.Value{
						"exact": types.StringNull(),
						"glob":  types.StringValue("spiffe://example.org/*"),
					}),
				}),
				ActorIdentity: types.ListValueMust(stringMatcherObjectType, []attr.Value{
					types.ObjectValueMust(stringMatcherAttrTypes, map[string]attr.Value{
						"exact": types.StringValue("spiffe://example.org/actor"),
						"glob":  types.StringNull(),
					}),
				}),
				ActorIssuer: types.ListValueMust(stringMatcherObjectType, []attr.Value{
					types.ObjectValueMust(stringMatcherAttrTypes, map[string]attr.Value{
						"exact": types.StringNull(),
						"glob":  types.StringValue("spiffe://issuer.example.org/*"),
					}),
				}),
				ClientID: types.ListValueMust(stringMatcherObjectType, []attr.Value{
					types.ObjectValueMust(stringMatcherAttrTypes, map[string]attr.Value{
						"exact": types.StringValue("my-client"),
						"glob":  types.StringNull(),
					}),
				}),
				TargetAudience: types.ListValueMust(stringMatcherObjectType, []attr.Value{
					types.ObjectValueMust(stringMatcherAttrTypes, map[string]attr.Value{
						"exact": types.StringValue("https://api.example.org"),
						"glob":  types.StringNull(),
					}),
				}),
				OutboundScopes: types.ListValueMust(types.StringType, []attr.Value{
					types.StringValue("read"),
					types.StringValue("write"),
				}),
			},
		},
		{
			name: "deny policy",
			model: ExchangePolicyModel{
				ID:              types.StringValue("ep-3"),
				OrgID:           types.StringValue(""),
				Name:            types.StringValue("deny-all"),
				TrustZoneID:     types.StringValue("tz-3"),
				Action:          types.StringValue("DENY"),
				SubjectIdentity: nullMatchers,
				SubjectIssuer:   nullMatchers,
				ActorIdentity:   nullMatchers,
				ActorIssuer:     nullMatchers,
				ClientID:        nullMatchers,
				TargetAudience:  nullMatchers,
				OutboundScopes:  types.ListValueMust(types.StringType, []attr.Value{}),
			},
		},
		{
			name: "policy with multiple matchers per string set",
			model: ExchangePolicyModel{
				ID:          types.StringValue("ep-4"),
				OrgID:       types.StringValue(""),
				Name:        types.StringValue("multi-matcher"),
				TrustZoneID: types.StringValue("tz-4"),
				Action:      types.StringNull(),
				SubjectIdentity: types.ListValueMust(stringMatcherObjectType, []attr.Value{
					types.ObjectValueMust(stringMatcherAttrTypes, map[string]attr.Value{
						"exact": types.StringValue("spiffe://example.org/workload-a"),
						"glob":  types.StringNull(),
					}),
					types.ObjectValueMust(stringMatcherAttrTypes, map[string]attr.Value{
						"exact": types.StringValue("spiffe://example.org/workload-b"),
						"glob":  types.StringNull(),
					}),
					types.ObjectValueMust(stringMatcherAttrTypes, map[string]attr.Value{
						"exact": types.StringNull(),
						"glob":  types.StringValue("spiffe://example.org/ns/*/sa/*"),
					}),
				}),
				SubjectIssuer:  nullMatchers,
				ActorIdentity:  nullMatchers,
				ActorIssuer:    nullMatchers,
				ClientID:       nullMatchers,
				TargetAudience: nullMatchers,
				OutboundScopes: types.ListValueMust(types.StringType, []attr.Value{}),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			proto, err := modelToProto(context.Background(), tt.model)
			require.NoError(t, err)
			got := protoToModel(proto)
			assert.Equal(t, tt.model, got)
		})
	}
}

func TestStringSetToProto_Nil(t *testing.T) {
	got, err := stringSetToProto(context.Background(), types.ListNull(stringMatcherObjectType))
	require.NoError(t, err)
	assert.Nil(t, got)
}

func TestStringSetToProto_Empty(t *testing.T) {
	got, err := stringSetToProto(context.Background(), types.ListValueMust(stringMatcherObjectType, []attr.Value{}))
	require.NoError(t, err)
	assert.NotNil(t, got)
	assert.Nil(t, got.Matchers)
}

func TestStringMatcherToProto_BothNull(t *testing.T) {
	model := StringMatcherModel{
		Exact: types.StringNull(),
		Glob:  types.StringNull(),
	}
	got, err := stringMatcherToProto(model)
	assert.Nil(t, got)
	assert.ErrorContains(t, err, "string matcher must set exactly one of exact or glob")
}

func TestStringMatcherToProto_ExactTakesPrecedence(t *testing.T) {
	model := StringMatcherModel{
		Exact: types.StringValue("exact-value"),
		Glob:  types.StringValue("glob-*"),
	}
	got, err := stringMatcherToProto(model)
	require.NoError(t, err)
	require.NotNil(t, got)
	assert.Equal(t, "exact-value", got.GetExact())
}

func TestStringSetToProto_BothNullMatcher(t *testing.T) {
	list := types.ListValueMust(stringMatcherObjectType, []attr.Value{
		types.ObjectValueMust(stringMatcherAttrTypes, map[string]attr.Value{
			"exact": types.StringNull(),
			"glob":  types.StringNull(),
		}),
	})
	got, err := stringSetToProto(context.Background(), list)
	assert.Nil(t, got)
	assert.ErrorContains(t, err, "string matcher must set exactly one of exact or glob")
}

func TestStringSetFromProto_Nil(t *testing.T) {
	got := stringSetFromProto(nil)
	assert.Equal(t, types.ListNull(stringMatcherObjectType), got)
}

func TestStringSetFromProto_Empty(t *testing.T) {
	// An empty StringSet (no matchers) is treated as absent.
	got := stringSetFromProto(&exchangepolicypb.StringSet{})
	assert.Equal(t, types.ListNull(stringMatcherObjectType), got)
}

func TestStringMatcherFromProto_Unknown(t *testing.T) {
	// A matcher with no match type set should return null fields.
	got := stringMatcherFromProto(&exchangepolicypb.StringMatcher{})
	assert.True(t, got.Exact.IsNull())
	assert.True(t, got.Glob.IsNull())
}
