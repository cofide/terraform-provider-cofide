package v1alpha1

import (
	"context"
	"testing"

	attestationpolicypb "github.com/cofide/cofide-api-sdk/gen/go/proto/attestation_policy/v1alpha1"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	spiretypes "github.com/spiffe/spire-api-sdk/proto/spire/api/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"
)

// mustObject converts value, a model struct, to an object of type t, as the
// model holds it.
func mustObject(t *testing.T, objType types.ObjectType, value any) types.Object {
	t.Helper()

	obj, diags := types.ObjectValueFrom(context.Background(), objType.AttrTypes, value)
	require.False(t, diags.HasError(), "failed to build object: %v", diags)

	return obj
}

// mustList converts elements, a slice of model structs, to a list of
// elemType, as the model holds it.
func mustList(t *testing.T, elemType types.ObjectType, elements any) types.List {
	t.Helper()

	list, diags := types.ListValueFrom(context.Background(), elemType, elements)
	require.False(t, diags.HasError(), "failed to build list: %v", diags)

	return list
}

// withTypedNulls returns model with each of its object and list attributes
// that a test case leaves unset (Go's zero value, which has no type) set to a
// typed null, as protoToModel and the framework produce. This keeps test
// cases to the attributes they're about.
func withTypedNulls(model AttestationPolicyModel) AttestationPolicyModel {
	ctx := context.Background()
	if model.Federations.ElementType(ctx) == nil {
		model.Federations = types.ListNull(federationElemType)
	}
	if len(model.Kubernetes.AttributeTypes(ctx)) == 0 {
		model.Kubernetes = types.ObjectNull(kubernetesType.AttrTypes)
	}
	if len(model.Static.AttributeTypes(ctx)) == 0 {
		model.Static = types.ObjectNull(staticType.AttrTypes)
	}
	if len(model.TPMNode.AttributeTypes(ctx)) == 0 {
		model.TPMNode = types.ObjectNull(tpmNodeType.AttrTypes)
	}
	return model
}

// TestProtoToModel exercises protoToModel across every attestation policy
// variant. Each case asserts on the full resulting AttestationPolicyModel,
// not just the fields the case is nominally about, so a regression in any
// field — not only the one under test — fails the case.
func TestProtoToModel(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		proto *attestationpolicypb.AttestationPolicy
		want  AttestationPolicyModel
	}{
		{
			name: "minimal: no federations or policy variant",
			proto: &attestationpolicypb.AttestationPolicy{
				Id:          new("ap-1"),
				Name:        "test-policy",
				OrgId:       new("org-1"),
				TrustZoneId: "tz-1",
			},
			want: AttestationPolicyModel{
				ID:          types.StringValue("ap-1"),
				Name:        types.StringValue("test-policy"),
				OrgID:       types.StringValue("org-1"),
				TrustZoneID: types.StringValue("tz-1"),
			},
		},
		{
			name: "trust zone and federations set",
			proto: &attestationpolicypb.AttestationPolicy{
				Id:          new("ap-1"),
				Name:        "test-policy",
				TrustZoneId: "tz-1",
				Federations: []*attestationpolicypb.Federation{
					{TrustZoneId: new("tz-2")},
					{TrustZoneId: new("tz-3")},
				},
			},
			want: AttestationPolicyModel{
				ID:          types.StringValue("ap-1"),
				Name:        types.StringValue("test-policy"),
				OrgID:       types.StringNull(),
				TrustZoneID: types.StringValue("tz-1"),
				Federations: mustList(t, federationElemType, []APFederationModel{
					{TrustZoneID: types.StringValue("tz-2")},
					{TrustZoneID: types.StringValue("tz-3")},
				}),
			},
		},
		{
			name: "kubernetes policy",
			proto: &attestationpolicypb.AttestationPolicy{
				Id:   new("ap-1"),
				Name: "k8s-policy",
				Policy: &attestationpolicypb.AttestationPolicy_Kubernetes{
					Kubernetes: &attestationpolicypb.APKubernetes{
						NamespaceSelector: &attestationpolicypb.APLabelSelector{
							MatchLabels: map[string]string{"env": "prod"},
							MatchExpressions: []*attestationpolicypb.APMatchExpression{
								{Key: "tier", Operator: "In", Values: []string{"web", "api"}},
							},
						},
						PodSelector:          &attestationpolicypb.APLabelSelector{},
						DnsNameTemplates:     []string{"{{ .PodMeta.Name }}.example.org"},
						SpiffeIdPathTemplate: new("ns/{{ .PodMeta.Namespace }}/sa/{{ .PodSpec.ServiceAccountName }}"),
					},
				},
			},
			want: AttestationPolicyModel{
				ID:          types.StringValue("ap-1"),
				Name:        types.StringValue("k8s-policy"),
				OrgID:       types.StringNull(),
				TrustZoneID: types.StringNull(),
				Kubernetes: mustObject(t, kubernetesType, APKubernetesModel{
					NamespaceSelector: mustObject(t, labelSelectorType, APLabelSelectorModel{
						MatchLabels: types.MapValueMust(types.StringType, map[string]attr.Value{"env": types.StringValue("prod")}),
						MatchExpressions: mustList(t, matchExpressionElemType, []APMatchExpressionModel{
							{
								Key:      types.StringValue("tier"),
								Operator: types.StringValue("In"),
								Values:   types.ListValueMust(types.StringType, []attr.Value{types.StringValue("web"), types.StringValue("api")}),
							},
						}),
					}),
					PodSelector: mustObject(t, labelSelectorType, APLabelSelectorModel{
						MatchLabels:      types.MapNull(types.StringType),
						MatchExpressions: types.ListNull(matchExpressionElemType),
					}),
					DnsNameTemplates:     types.ListValueMust(types.StringType, []attr.Value{types.StringValue("{{ .PodMeta.Name }}.example.org")}),
					SpiffeIDPathTemplate: types.StringValue("ns/{{ .PodMeta.Namespace }}/sa/{{ .PodSpec.ServiceAccountName }}"),
				}),
			},
		},
		{
			name: "static policy",
			proto: &attestationpolicypb.AttestationPolicy{
				Id:   new("ap-1"),
				Name: "static-policy",
				Policy: &attestationpolicypb.AttestationPolicy_Static{
					Static: &attestationpolicypb.APStatic{
						SpiffeIdPath: new("ns/default/sa/my-service-account"),
						ParentIdPath: new("spire/agent/join_token/abc123"),
						Selectors: []*spiretypes.Selector{
							{Type: "k8s", Value: "ns:default"},
						},
						DnsNames:  []string{"my-service.example.org"},
						StoreSvid: true,
					},
				},
			},
			want: AttestationPolicyModel{
				ID:          types.StringValue("ap-1"),
				Name:        types.StringValue("static-policy"),
				OrgID:       types.StringNull(),
				TrustZoneID: types.StringNull(),
				Static: mustObject(t, staticType, APStaticModel{
					SpiffeIDPath: types.StringValue("ns/default/sa/my-service-account"),
					ParentIdPath: types.StringValue("spire/agent/join_token/abc123"),
					Selectors: types.ListValueMust(selectorElemType, []attr.Value{
						types.ObjectValueMust(selectorAttrTypes, map[string]attr.Value{
							"type":  types.StringValue("k8s"),
							"value": types.StringValue("ns:default"),
						}),
					}),
					DNSNames:  types.ListValueMust(types.StringType, []attr.Value{types.StringValue("my-service.example.org")}),
					StoreSvid: types.BoolValue(true),
				}),
			},
		},
		{
			name: "tpm_node policy",
			proto: &attestationpolicypb.AttestationPolicy{
				Id:   new("ap-1"),
				Name: "tpm-policy",
				Policy: &attestationpolicypb.AttestationPolicy_TpmNode{
					TpmNode: &attestationpolicypb.APTPMNode{
						Attestation:    &attestationpolicypb.TPMAttestation{EkHash: new("deadbeef")},
						SelectorValues: []string{"selector-a"},
					},
				},
			},
			want: AttestationPolicyModel{
				ID:          types.StringValue("ap-1"),
				Name:        types.StringValue("tpm-policy"),
				OrgID:       types.StringNull(),
				TrustZoneID: types.StringNull(),
				TPMNode: mustObject(t, tpmNodeType, APTPMNodeModel{
					Attestation:    mustObject(t, tpmAttestationType, TPMAttestationModel{EKHash: types.StringValue("deadbeef")}),
					SelectorValues: types.ListValueMust(types.StringType, []attr.Value{types.StringValue("selector-a")}),
				}),
			},
		},
		{
			name: "tpm_node policy without attestation",
			proto: &attestationpolicypb.AttestationPolicy{
				Id:   new("ap-1"),
				Name: "tpm-policy",
				Policy: &attestationpolicypb.AttestationPolicy_TpmNode{
					TpmNode: &attestationpolicypb.APTPMNode{},
				},
			},
			want: AttestationPolicyModel{
				ID:          types.StringValue("ap-1"),
				Name:        types.StringValue("tpm-policy"),
				OrgID:       types.StringNull(),
				TrustZoneID: types.StringNull(),
				TPMNode: mustObject(t, tpmNodeType, APTPMNodeModel{
					Attestation:    mustObject(t, tpmAttestationType, TPMAttestationModel{EKHash: types.StringNull()}),
					SelectorValues: types.ListNull(types.StringType),
				}),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, diags := protoToModel(t.Context(), tt.proto)
			require.False(t, diags.HasError(), "diagnostics: %v", diags)
			assert.Equal(t, withTypedNulls(tt.want), got)
		})
	}
}

// TestModelToProto exercises modelToProto across every attestation policy
// variant. Each case asserts on the full resulting proto message via
// proto.Equal — the correct way to compare protobuf messages, since it
// treats absent/empty repeated fields as equivalent and ignores the
// internal bookkeeping fields that make reflect-based struct equality
// unreliable for generated proto types.
func TestModelToProto(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		model AttestationPolicyModel
		want  *attestationpolicypb.AttestationPolicy
	}{
		{
			name: "minimal: no federations or policy variant",
			model: AttestationPolicyModel{
				ID:          types.StringValue("ap-1"),
				Name:        types.StringValue("test-policy"),
				TrustZoneID: types.StringValue("tz-1"),
			},
			want: &attestationpolicypb.AttestationPolicy{
				Id:          new("ap-1"),
				Name:        "test-policy",
				TrustZoneId: "tz-1",
			},
		},
		{
			name: "trust zone and federations set",
			model: AttestationPolicyModel{
				Name:        types.StringValue("test-policy"),
				TrustZoneID: types.StringValue("tz-1"),
				Federations: mustList(t, federationElemType, []APFederationModel{
					{TrustZoneID: types.StringValue("tz-2")},
				}),
			},
			want: &attestationpolicypb.AttestationPolicy{
				Name:        "test-policy",
				TrustZoneId: "tz-1",
				Federations: []*attestationpolicypb.Federation{
					{TrustZoneId: new("tz-2")},
				},
			},
		},
		{
			name: "kubernetes policy",
			model: AttestationPolicyModel{
				Name: types.StringValue("k8s-policy"),
				Kubernetes: mustObject(t, kubernetesType, APKubernetesModel{
					NamespaceSelector: mustObject(t, labelSelectorType, APLabelSelectorModel{
						MatchLabels: types.MapValueMust(types.StringType, map[string]attr.Value{"env": types.StringValue("prod")}),
						MatchExpressions: mustList(t, matchExpressionElemType, []APMatchExpressionModel{
							{
								Key:      types.StringValue("tier"),
								Operator: types.StringValue("In"),
								Values:   types.ListValueMust(types.StringType, []attr.Value{types.StringValue("web")}),
							},
						}),
					}),
					PodSelector:          types.ObjectNull(labelSelectorType.AttrTypes),
					DnsNameTemplates:     types.ListValueMust(types.StringType, []attr.Value{types.StringValue("{{ .PodMeta.Name }}.example.org")}),
					SpiffeIDPathTemplate: types.StringValue("ns/{{ .PodMeta.Namespace }}/sa/{{ .PodSpec.ServiceAccountName }}"),
				}),
			},
			want: &attestationpolicypb.AttestationPolicy{
				Name: "k8s-policy",
				Policy: &attestationpolicypb.AttestationPolicy_Kubernetes{
					Kubernetes: &attestationpolicypb.APKubernetes{
						NamespaceSelector: &attestationpolicypb.APLabelSelector{
							MatchLabels: map[string]string{"env": "prod"},
							MatchExpressions: []*attestationpolicypb.APMatchExpression{
								{Key: "tier", Operator: "In", Values: []string{"web"}},
							},
						},
						DnsNameTemplates:     []string{"{{ .PodMeta.Name }}.example.org"},
						SpiffeIdPathTemplate: new("ns/{{ .PodMeta.Namespace }}/sa/{{ .PodSpec.ServiceAccountName }}"),
					},
				},
			},
		},
		{
			name: "static policy",
			model: AttestationPolicyModel{
				Name: types.StringValue("static-policy"),
				Static: mustObject(t, staticType, APStaticModel{
					SpiffeIDPath: types.StringValue("ns/default/sa/my-service-account"),
					ParentIdPath: types.StringValue("spire/agent/join_token/abc123"),
					Selectors: types.ListValueMust(selectorElemType, []attr.Value{
						types.ObjectValueMust(selectorAttrTypes, map[string]attr.Value{
							"type":  types.StringValue("k8s"),
							"value": types.StringValue("ns:default"),
						}),
					}),
					DNSNames:  types.ListValueMust(types.StringType, []attr.Value{types.StringValue("my-service.example.org")}),
					StoreSvid: types.BoolValue(true),
				}),
			},
			want: &attestationpolicypb.AttestationPolicy{
				Name: "static-policy",
				Policy: &attestationpolicypb.AttestationPolicy_Static{
					Static: &attestationpolicypb.APStatic{
						SpiffeIdPath: new("ns/default/sa/my-service-account"),
						ParentIdPath: new("spire/agent/join_token/abc123"),
						Selectors: []*spiretypes.Selector{
							{Type: "k8s", Value: "ns:default"},
						},
						DnsNames:  []string{"my-service.example.org"},
						StoreSvid: true,
					},
				},
			},
		},
		{
			name: "static policy with store_svid unknown, as planned when it isn't configured",
			model: AttestationPolicyModel{
				Name: types.StringValue("static-policy"),
				Static: mustObject(t, staticType, APStaticModel{
					SpiffeIDPath: types.StringValue("ns/default/sa/my-service-account"),
					ParentIdPath: types.StringValue("spire/agent/join_token/abc123"),
					Selectors:    types.ListValueMust(selectorElemType, []attr.Value{}),
					DNSNames:     types.ListNull(types.StringType),
					StoreSvid:    types.BoolUnknown(),
				}),
			},
			want: &attestationpolicypb.AttestationPolicy{
				Name: "static-policy",
				Policy: &attestationpolicypb.AttestationPolicy_Static{
					Static: &attestationpolicypb.APStatic{
						SpiffeIdPath: new("ns/default/sa/my-service-account"),
						ParentIdPath: new("spire/agent/join_token/abc123"),
					},
				},
			},
		},
		{
			name: "tpm_node policy",
			model: AttestationPolicyModel{
				Name: types.StringValue("tpm-policy"),
				TPMNode: mustObject(t, tpmNodeType, APTPMNodeModel{
					Attestation:    mustObject(t, tpmAttestationType, TPMAttestationModel{EKHash: types.StringValue("deadbeef")}),
					SelectorValues: types.ListValueMust(types.StringType, []attr.Value{types.StringValue("selector-a")}),
				}),
			},
			want: &attestationpolicypb.AttestationPolicy{
				Name: "tpm-policy",
				Policy: &attestationpolicypb.AttestationPolicy_TpmNode{
					TpmNode: &attestationpolicypb.APTPMNode{
						Attestation:    &attestationpolicypb.TPMAttestation{EkHash: new("deadbeef")},
						SelectorValues: []string{"selector-a"},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, diags := modelToProto(t.Context(), withTypedNulls(tt.model))
			require.False(t, diags.HasError(), "diagnostics: %v", diags)
			assert.Truef(t, proto.Equal(tt.want, got), "got: %v\nwant: %v", got, tt.want)
		})
	}
}

// TestRoundTrip verifies that model -> proto -> model conversion is lossless
// for every attestation policy variant, including the trust_zone_id and
// federations fields.
func TestRoundTrip(t *testing.T) {
	t.Parallel()

	emptyKubernetes := APKubernetesModel{
		NamespaceSelector:    types.ObjectNull(labelSelectorType.AttrTypes),
		PodSelector:          types.ObjectNull(labelSelectorType.AttrTypes),
		DnsNameTemplates:     types.ListNull(types.StringType),
		SpiffeIDPathTemplate: types.StringNull(),
	}

	tests := []struct {
		name  string
		model AttestationPolicyModel
	}{
		{
			name: "kubernetes policy owned by a trust zone with federations",
			model: AttestationPolicyModel{
				ID:          types.StringValue("ap-1"),
				Name:        types.StringValue("k8s-policy"),
				TrustZoneID: types.StringValue("tz-1"),
				Federations: mustList(t, federationElemType, []APFederationModel{
					{TrustZoneID: types.StringValue("tz-2")},
				}),
				Kubernetes: mustObject(t, kubernetesType, emptyKubernetes),
			},
		},
		{
			name: "trust zone owned policy with explicitly empty federations",
			model: AttestationPolicyModel{
				ID:          types.StringValue("ap-4"),
				Name:        types.StringValue("k8s-policy"),
				TrustZoneID: types.StringValue("tz-1"),
				Federations: mustList(t, federationElemType, []APFederationModel{}),
				Kubernetes:  mustObject(t, kubernetesType, emptyKubernetes),
			},
		},
		{
			name: "static policy owned by a trust zone",
			model: AttestationPolicyModel{
				ID:          types.StringValue("ap-2"),
				Name:        types.StringValue("static-policy"),
				TrustZoneID: types.StringValue("tz-1"),
				Static: mustObject(t, staticType, APStaticModel{
					SpiffeIDPath: types.StringValue("ns/default/sa/my-service-account"),
					ParentIdPath: types.StringValue("spire/agent/join_token/abc123"),
					Selectors:    types.ListValueMust(selectorElemType, []attr.Value{}),
					DNSNames:     types.ListNull(types.StringType),
					StoreSvid:    types.BoolValue(false),
				}),
			},
		},
		{
			name: "tpm_node policy",
			model: AttestationPolicyModel{
				ID:   types.StringValue("ap-3"),
				Name: types.StringValue("tpm-policy"),
				TPMNode: mustObject(t, tpmNodeType, APTPMNodeModel{
					Attestation:    mustObject(t, tpmAttestationType, TPMAttestationModel{EKHash: types.StringValue("deadbeef")}),
					SelectorValues: types.ListNull(types.StringType),
				}),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			model := withTypedNulls(tt.model)

			gotProto, diags := modelToProto(t.Context(), model)
			require.False(t, diags.HasError(), "diagnostics: %v", diags)

			got, diags := protoToModel(t.Context(), gotProto)
			require.False(t, diags.HasError(), "diagnostics: %v", diags)

			assert.Equal(t, model, got)
		})
	}
}
