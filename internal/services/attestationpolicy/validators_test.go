package attestationpolicy

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	rschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// testProvider is a minimal provider.Provider that registers only the
// attestation policy resource, so its real ValidateResourceConfig RPC —
// the same one Terraform itself calls — can be driven in-process without a
// live Connect client.
type testProvider struct{}

func (testProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "cofide"
}

func (testProvider) Schema(context.Context, provider.SchemaRequest, *provider.SchemaResponse) {}

func (testProvider) Configure(context.Context, provider.ConfigureRequest, *provider.ConfigureResponse) {
}

func (testProvider) DataSources(context.Context) []func() datasource.DataSource { return nil }

func (testProvider) Resources(context.Context) []func() resource.Resource {
	return []func() resource.Resource{NewResource}
}

// rawFromModel converts model into the tftypes.Value the framework would
// send on the wire for it, via the same reflection-based Set used in
// production (e.g. resource.go's resp.State.Set(ctx, &state)). This lets
// test cases describe a config as an AttestationPolicyModel literal — the
// same type used throughout convert_test.go — instead of a bespoke set of
// booleans/pointers that need extending whenever a new attribute is added.
func rawFromModel(t *testing.T, rSchema rschema.Schema, model AttestationPolicyModel) tftypes.Value {
	t.Helper()

	plan := tfsdk.Plan{Schema: rSchema}
	diags := plan.Set(t.Context(), &model)
	require.False(t, diags.HasError(), "failed to build raw value from model: %v", diags)

	return plan.Raw
}

// validateResourceConfig drives the attestation policy resource's real
// ValidateResourceConfig RPC against raw, via an in-process protocol server.
// Unlike hand-invoking specific validator instances, this exercises the
// framework's own schema walk — every attribute type and nesting depth,
// generically — so a validator added on a future attribute (of any type, at
// any depth) is covered automatically with no changes needed here.
func validateResourceConfig(t *testing.T, rSchema rschema.Schema, raw tftypes.Value) []*tfprotov6.Diagnostic {
	t.Helper()

	srv, err := providerserver.NewProtocol6WithError(testProvider{})()
	require.NoError(t, err)

	schemaType, ok := rSchema.Type().TerraformType(t.Context()).(tftypes.Object)
	require.True(t, ok, "expected resource schema type to be an Object")

	dv, err := tfprotov6.NewDynamicValue(schemaType, raw)
	require.NoError(t, err)

	resp, err := srv.ValidateResourceConfig(t.Context(), &tfprotov6.ValidateResourceConfigRequest{
		TypeName: "cofide_connect_attestation_policy",
		Config:   &dv,
	})
	require.NoError(t, err)

	return resp.Diagnostics
}

// wantDiagnostic is the comparable subset of a tfprotov6.Diagnostic used to
// assert on validation errors: which attribute it's attached to (empty for
// a resource-level diagnostic), its summary, and its detail. Asserting on
// this rather than just a count means a test case can't keep passing after
// a behavior change that swaps in a different error with the same count
// (e.g. the wrong variant-count message, or a conflict reported against the
// wrong attribute).
type wantDiagnostic struct {
	attribute string
	summary   string
	detail    string
}

// errorDiagnostics filters diags down to error-severity entries (ignoring
// any warnings) and reduces each to its comparable wantDiagnostic shape.
func errorDiagnostics(diags []*tfprotov6.Diagnostic) []wantDiagnostic {
	var errs []wantDiagnostic
	for _, d := range diags {
		if d.Severity != tfprotov6.DiagnosticSeverityError {
			continue
		}
		var attr string
		if d.Attribute != nil {
			attr = d.Attribute.String()
		}
		errs = append(errs, wantDiagnostic{attribute: attr, summary: d.Summary, detail: d.Detail})
	}
	return errs
}

// attrPath renders name the same way tfprotov6.Diagnostic.Attribute does,
// so expected paths are built from the same library call that produces the
// real ones rather than a hand-guessed string format.
func attrPath(name string) string {
	return tftypes.NewAttributePath().WithAttributeName(name).String()
}

// errNoVariant is the resource-level diagnostic exactlyOneOfValidator
// produces when none of kubernetes/static/tpm_node is set.
func errNoVariant() wantDiagnostic {
	return wantDiagnostic{
		summary: "Invalid configuration",
		detail:  "Exactly one of kubernetes, static or tpm_node blocks must be configured, but none were provided.",
	}
}

// errMultipleVariants is the resource-level diagnostic exactlyOneOfValidator
// produces when more than one of kubernetes/static/tpm_node is set.
func errMultipleVariants() wantDiagnostic {
	return wantDiagnostic{
		summary: "Invalid configuration",
		detail:  "Exactly one of kubernetes, static or tpm_node blocks must be configured, but multiple were provided.",
	}
}

// errConflict is the diagnostic ConflictsWith attaches to attribute when
// other is also set. Terraform runs each attribute's validators
// independently, so a conflicting pair produces one of these per side —
// e.g. errConflict("org_id", "trust_zone_id") is org_id's own complaint,
// distinct from errConflict("trust_zone_id", "org_id").
func errConflict(attribute, other string) wantDiagnostic {
	return wantDiagnostic{
		attribute: attrPath(attribute),
		summary:   "Invalid Attribute Combination",
		detail:    fmt.Sprintf("Attribute %q cannot be specified when %q is specified", other, attribute),
	}
}

// minimalKubernetes is a `kubernetes` policy variant with none of its own
// fields set. Valid on its own: kubernetes has no Required nested
// attributes, so this is enough to test the exactly-one-variant rule
// without incidental "attribute is required" noise from other rules.
func minimalKubernetes() *APKubernetesModel {
	return &APKubernetesModel{
		DnsNameTemplates:     types.ListNull(types.StringType),
		SpiffeIDPathTemplate: types.StringNull(),
	}
}

// minimalStatic is a `static` policy variant with the minimum content its
// Required attributes (spiffe_id_path, parent_id_path, selectors) need.
func minimalStatic() *APStaticModel {
	return &APStaticModel{
		SpiffeIDPath: types.StringValue("ns/default/sa/my-service-account"),
		ParentIdPath: types.StringValue("spire/agent/join_token/abc123"),
		Selectors:    types.ListValueMust(selectorElemType, []attr.Value{}),
		DNSNames:     types.ListNull(types.StringType),
		StoreSvid:    types.BoolNull(),
	}
}

// minimalTPMNode is a `tpm_node` policy variant with the minimum content its
// Required attribute (attestation.ek_hash) needs.
func minimalTPMNode() *APTPMNodeModel {
	return &APTPMNodeModel{
		Attestation:    TPMAttestationModel{EKHash: types.StringValue("deadbeef")},
		SelectorValues: types.ListNull(types.StringType),
	}
}

// TestValidateResourceConfig exercises the attestation policy schema's
// validation rules as a whole, through the real ValidateResourceConfig RPC:
// for a given config, is it valid, and if not, what errors does it
// produce. This covers both business rules the schema enforces.
func TestValidateResourceConfig(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		model    AttestationPolicyModel
		wantErrs []wantDiagnostic
	}{
		{
			name: "valid: kubernetes policy, no org or trust zone set",
			model: AttestationPolicyModel{
				Name:       types.StringValue("test-policy"),
				Kubernetes: minimalKubernetes(),
			},
		},
		{
			name: "valid: static policy with org_id",
			model: AttestationPolicyModel{
				Name:   types.StringValue("test-policy"),
				OrgID:  types.StringValue("org-1"),
				Static: minimalStatic(),
			},
		},
		{
			name: "valid: tpm_node policy owned by a trust zone",
			model: AttestationPolicyModel{
				Name:        types.StringValue("test-policy"),
				TrustZoneID: types.StringValue("tz-1"),
				TPMNode:     minimalTPMNode(),
			},
		},
		{
			name: "invalid: no policy variant set",
			model: AttestationPolicyModel{
				Name: types.StringValue("test-policy"),
			},
			wantErrs: []wantDiagnostic{errNoVariant()},
		},
		{
			name: "invalid: two policy variants set",
			model: AttestationPolicyModel{
				Name:       types.StringValue("test-policy"),
				Kubernetes: minimalKubernetes(),
				Static:     minimalStatic(),
			},
			wantErrs: []wantDiagnostic{errMultipleVariants()},
		},
		{
			name: "invalid: all three policy variants set",
			model: AttestationPolicyModel{
				Name:       types.StringValue("test-policy"),
				Kubernetes: minimalKubernetes(),
				Static:     minimalStatic(),
				TPMNode:    minimalTPMNode(),
			},
			wantErrs: []wantDiagnostic{errMultipleVariants()},
		},
		{
			name: "invalid: org_id and trust_zone_id both set",
			model: AttestationPolicyModel{
				Name:        types.StringValue("test-policy"),
				OrgID:       types.StringValue("org-1"),
				TrustZoneID: types.StringValue("tz-1"),
				Kubernetes:  minimalKubernetes(),
			},
			wantErrs: []wantDiagnostic{
				errConflict("org_id", "trust_zone_id"),
				errConflict("trust_zone_id", "org_id"),
			},
		},
		{
			name: "invalid: no policy variant, and org_id/trust_zone_id both set",
			model: AttestationPolicyModel{
				Name:        types.StringValue("test-policy"),
				OrgID:       types.StringValue("org-1"),
				TrustZoneID: types.StringValue("tz-1"),
			},
			wantErrs: []wantDiagnostic{
				errNoVariant(),
				errConflict("org_id", "trust_zone_id"),
				errConflict("trust_zone_id", "org_id"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			rSchema := ResourceSchema(t.Context())
			raw := rawFromModel(t, rSchema, tt.model)

			diags := validateResourceConfig(t, rSchema, raw)
			errs := errorDiagnostics(diags)

			assert.ElementsMatchf(t, tt.wantErrs, errs, "diagnostics: %v", diags)
		})
	}
}
