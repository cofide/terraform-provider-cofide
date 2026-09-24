package internal_test

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/stretchr/testify/require"

	"github.com/cofide/terraform-provider-cofide/internal"
	attestationpolicy "github.com/cofide/terraform-provider-cofide/internal/services/attestationpolicy/v1alpha1"
	cluster "github.com/cofide/terraform-provider-cofide/internal/services/cluster/v1alpha1"
	exchangepolicy "github.com/cofide/terraform-provider-cofide/internal/services/exchangepolicy/v1alpha1"
	federation "github.com/cofide/terraform-provider-cofide/internal/services/federation/v1alpha1"
	rolebinding "github.com/cofide/terraform-provider-cofide/internal/services/rolebinding/v1alpha1"
	trustzone "github.com/cofide/terraform-provider-cofide/internal/services/trustzone/v1alpha1"
	trustzoneserver "github.com/cofide/terraform-provider-cofide/internal/services/trustzoneserver/v1alpha1"
)

// TestImportedStateDecodesIntoModel asserts that the state left behind by
// ImportState can be decoded into each resource's *Model.
//
// ImportStatePassthroughID sets only `id`, leaving every other attribute null,
// and Read then runs against that state. Resources whose Read decodes the
// whole prior state (e.g. cluster and trustzoneserver, which need it to
// preserve user formatting) fail with a Value Conversion Error if the model
// holds a nested object as a non-pointer struct, because a struct cannot
// represent null. Nested objects must be pointers or types.Object.
func TestImportedStateDecodesIntoModel(t *testing.T) {
	t.Parallel()
	// decoders maps each resource type name to its model type. The model type
	// is not reachable through the resource.Resource interface, so it has to
	// be listed here.
	decoders := map[string]func(context.Context, tfsdk.State) diag.Diagnostics{
		"cofide_connect_attestation_policy_v1alpha1": decodeAs[attestationpolicy.AttestationPolicyModel](),
		"cofide_connect_cluster_v1alpha1":            decodeAs[cluster.ClusterModel](),
		"cofide_connect_exchange_policy_v1alpha1":    decodeAs[exchangepolicy.ExchangePolicyModel](),
		"cofide_connect_federation_v1alpha1":         decodeAs[federation.FederationModel](),
		"cofide_connect_role_binding_v1alpha1":       decodeAs[rolebinding.RoleBindingModel](),
		"cofide_connect_trust_zone_v1alpha1":         decodeAs[trustzone.TrustZoneModel](),
		"cofide_connect_trust_zone_server_v1alpha1":  decodeAs[trustzoneserver.TrustZoneServerModel](),
	}

	// Iterate over the provider's registered resources rather than the table,
	// so that a new resource without a decoder fails the test.
	for _, newResource := range internal.NewProvider("test")().Resources(t.Context()) {
		r := newResource()

		metadataResp := &resource.MetadataResponse{}
		r.Metadata(t.Context(), resource.MetadataRequest{ProviderTypeName: "cofide"}, metadataResp)
		name := metadataResp.TypeName

		t.Run(name, func(t *testing.T) {
			t.Parallel()
			decode, ok := decoders[name]
			require.True(t, ok, "no decoder for %s; add one to this test", name)

			importer, ok := r.(resource.ResourceWithImportState)
			require.True(t, ok, "%s does not support import", name)

			schemaResp := &resource.SchemaResponse{}
			r.Schema(t.Context(), resource.SchemaRequest{}, schemaResp)
			require.False(t, schemaResp.Diagnostics.HasError(), "%v", schemaResp.Diagnostics)

			// Terraform starts an import from a null state.
			importResp := &resource.ImportStateResponse{
				State: tfsdk.State{
					Schema: schemaResp.Schema,
					Raw:    tftypes.NewValue(schemaResp.Schema.Type().TerraformType(t.Context()), nil),
				},
			}
			importer.ImportState(t.Context(), resource.ImportStateRequest{ID: "test-id"}, importResp)
			require.False(t, importResp.Diagnostics.HasError(), "%v", importResp.Diagnostics)

			diags := decode(t.Context(), importResp.State)
			require.False(t, diags.HasError(), "imported state does not decode into the model: %v", diags)
		})
	}
}

// decodeAs returns a function that decodes a state into a zero value of T.
func decodeAs[T any]() func(context.Context, tfsdk.State) diag.Diagnostics {
	return func(ctx context.Context, s tfsdk.State) diag.Diagnostics {
		var m T
		return s.Get(ctx, &m)
	}
}
