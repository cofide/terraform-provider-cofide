package attestationpolicy

import (
	"testing"

	attestationpolicypb "github.com/cofide/cofide-api-sdk/gen/go/proto/attestation_policy/v1alpha1"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	spiretypes "github.com/spiffe/spire-api-sdk/proto/spire/api/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"
)

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
			name: "minimal: no trust zone, federations, or policy variant",
			proto: &attestationpolicypb.AttestationPolicy{
				Id:    new("ap-1"),
				Name:  "test-policy",
				OrgId: new("org-1"),
			},
			want: AttestationPolicyModel{
				ID:          types.StringValue("ap-1"),
				Name:        types.StringValue("test-policy"),
				OrgID:       types.StringValue("org-1"),
				TrustZoneID: types.StringNull(),
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
				Federations: []APFederationModel{
					{TrustZoneID: types.StringValue("tz-2")},
					{TrustZoneID: types.StringValue("tz-3")},
				},
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
				Kubernetes: &APKubernetesModel{
					NamespaceSelector: &APLabelSelectorModel{
						MatchLabels: types.MapValueMust(types.StringType, map[string]attr.Value{"env": types.StringValue("prod")}),
						MatchExpressions: []APMatchExpressionModel{
							{
								Key:      types.StringValue("tier"),
								Operator: types.StringValue("In"),
								Values:   types.ListValueMust(types.StringType, []attr.Value{types.StringValue("web"), types.StringValue("api")}),
							},
						},
					},
					PodSelector: &APLabelSelectorModel{
						MatchLabels:      types.MapNull(types.StringType),
						MatchExpressions: nil,
					},
					DnsNameTemplates:     types.ListValueMust(types.StringType, []attr.Value{types.StringValue("{{ .PodMeta.Name }}.example.org")}),
					SpiffeIDPathTemplate: types.StringValue("ns/{{ .PodMeta.Namespace }}/sa/{{ .PodSpec.ServiceAccountName }}"),
				},
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
				Static: &APStaticModel{
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
				},
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
				TPMNode: &APTPMNodeModel{
					Attestation:    TPMAttestationModel{EKHash: types.StringValue("deadbeef")},
					SelectorValues: types.ListValueMust(types.StringType, []attr.Value{types.StringValue("selector-a")}),
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := protoToModel(tt.proto)
			assert.Equal(t, tt.want, got)
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
			name: "minimal: no trust zone, federations, or policy variant",
			model: AttestationPolicyModel{
				ID:    types.StringValue("ap-1"),
				Name:  types.StringValue("test-policy"),
				OrgID: types.StringValue("org-1"),
			},
			want: &attestationpolicypb.AttestationPolicy{
				Id:    new("ap-1"),
				Name:  "test-policy",
				OrgId: new("org-1"),
			},
		},
		{
			name: "trust zone and federations set",
			model: AttestationPolicyModel{
				Name:        types.StringValue("test-policy"),
				TrustZoneID: types.StringValue("tz-1"),
				Federations: []APFederationModel{
					{TrustZoneID: types.StringValue("tz-2")},
				},
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
				Kubernetes: &APKubernetesModel{
					NamespaceSelector: &APLabelSelectorModel{
						MatchLabels: types.MapValueMust(types.StringType, map[string]attr.Value{"env": types.StringValue("prod")}),
						MatchExpressions: []APMatchExpressionModel{
							{
								Key:      types.StringValue("tier"),
								Operator: types.StringValue("In"),
								Values:   types.ListValueMust(types.StringType, []attr.Value{types.StringValue("web")}),
							},
						},
					},
					DnsNameTemplates:     types.ListValueMust(types.StringType, []attr.Value{types.StringValue("{{ .PodMeta.Name }}.example.org")}),
					SpiffeIDPathTemplate: types.StringValue("ns/{{ .PodMeta.Namespace }}/sa/{{ .PodSpec.ServiceAccountName }}"),
				},
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
				Static: &APStaticModel{
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
				},
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
			name: "tpm_node policy",
			model: AttestationPolicyModel{
				Name: types.StringValue("tpm-policy"),
				TPMNode: &APTPMNodeModel{
					Attestation:    TPMAttestationModel{EKHash: types.StringValue("deadbeef")},
					SelectorValues: types.ListValueMust(types.StringType, []attr.Value{types.StringValue("selector-a")}),
				},
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
			got, diags := modelToProto(t.Context(), tt.model)
			require.False(t, diags.HasError())
			assert.Truef(t, proto.Equal(tt.want, got), "got: %v\nwant: %v", got, tt.want)
		})
	}
}

// TestRoundTrip verifies that model -> proto -> model conversion is lossless
// for every attestation policy variant, including the trust_zone_id and
// federations fields.
func TestRoundTrip(t *testing.T) {
	t.Parallel()
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
				Federations: []APFederationModel{
					{TrustZoneID: types.StringValue("tz-2")},
				},
				Kubernetes: &APKubernetesModel{
					NamespaceSelector:    nil,
					PodSelector:          nil,
					DnsNameTemplates:     types.ListNull(types.StringType),
					SpiffeIDPathTemplate: types.StringNull(),
				},
			},
		},
		{
			name: "trust zone owned policy with explicitly empty federations",
			model: AttestationPolicyModel{
				ID:          types.StringValue("ap-4"),
				Name:        types.StringValue("k8s-policy"),
				TrustZoneID: types.StringValue("tz-1"),
				Federations: []APFederationModel{},
				Kubernetes: &APKubernetesModel{
					NamespaceSelector:    nil,
					PodSelector:          nil,
					DnsNameTemplates:     types.ListNull(types.StringType),
					SpiffeIDPathTemplate: types.StringNull(),
				},
			},
		},
		{
			name: "static policy with org_id, no trust zone",
			model: AttestationPolicyModel{
				ID:    types.StringValue("ap-2"),
				Name:  types.StringValue("static-policy"),
				OrgID: types.StringValue("org-1"),
				Static: &APStaticModel{
					SpiffeIDPath: types.StringValue("ns/default/sa/my-service-account"),
					ParentIdPath: types.StringValue("spire/agent/join_token/abc123"),
					Selectors:    types.ListValueMust(selectorElemType, []attr.Value{}),
					DNSNames:     types.ListNull(types.StringType),
					StoreSvid:    types.BoolValue(false),
				},
			},
		},
		{
			name: "tpm_node policy",
			model: AttestationPolicyModel{
				ID:   types.StringValue("ap-3"),
				Name: types.StringValue("tpm-policy"),
				TPMNode: &APTPMNodeModel{
					Attestation:    TPMAttestationModel{EKHash: types.StringValue("deadbeef")},
					SelectorValues: types.ListNull(types.StringType),
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			gotProto, diags := modelToProto(t.Context(), tt.model)
			require.False(t, diags.HasError())

			got := protoToModel(gotProto)

			assert.Equal(t, tt.model, got)
		})
	}
}
