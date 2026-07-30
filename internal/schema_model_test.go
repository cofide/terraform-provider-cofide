package internal_test

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/stretchr/testify/require"

	"github.com/cofide/terraform-provider-cofide/internal/services/apbinding"
	"github.com/cofide/terraform-provider-cofide/internal/services/attestationpolicy"
	"github.com/cofide/terraform-provider-cofide/internal/services/cluster"
	"github.com/cofide/terraform-provider-cofide/internal/services/exchangepolicy"
	"github.com/cofide/terraform-provider-cofide/internal/services/federation"
	"github.com/cofide/terraform-provider-cofide/internal/services/organization"
	"github.com/cofide/terraform-provider-cofide/internal/services/rolebinding"
	"github.com/cofide/terraform-provider-cofide/internal/services/trustzone"
	"github.com/cofide/terraform-provider-cofide/internal/services/trustzoneserver"
)

// schemaModel pairs a schema with the model struct that Read populates from it.
type schemaModel struct {
	name   string
	schema attr.Type
	// model must be a pointer to a fresh zero value of the model struct.
	model any
}

// TestSchemaMatchesModel asserts that every schema can be read into its model
// struct without conversion errors.
//
// Every Read ends in State.Set, which reflects the model struct against the
// schema — the same reflection exercised here — and fails if either side has a
// field the other lacks. That failure
// otherwise only appears against a live Connect deployment, as a Value
// Conversion Error at plan or apply time — which is how the attestation policy
// data source shipped without static.store_svid.
//
// This complements TestSchemaShapesMatch: that test compares a resource schema
// against its data source schema, so it cannot catch a field missing from both.
// This one anchors each schema to the model independently, and so also covers
// rolebinding and organization, which have no resource/data source pair.
//
// One limitation: fields typed as an untyped collection value (tftypes.List
// rather than a []SomeModel slice) are opaque to this reflection, so their
// element attributes are not checked. attestationpolicy's static.selectors is
// one such field.
func TestSchemaMatchesModel(t *testing.T) {
	ctx := context.Background()

	tests := []schemaModel{
		// Resources.
		{"apbinding resource", apbinding.ResourceSchema(ctx).Type(), &apbinding.APBindingModel{}},
		{"attestationpolicy resource", attestationpolicy.ResourceSchema(ctx).Type(), &attestationpolicy.AttestationPolicyModel{}},
		{"cluster resource", cluster.ResourceSchema(ctx).Type(), &cluster.ClusterModel{}},
		{"exchangepolicy resource", exchangepolicy.ResourceSchema().Type(), &exchangepolicy.ExchangePolicyModel{}},
		{"federation resource", federation.ResourceSchema(ctx).Type(), &federation.FederationModel{}},
		{"rolebinding resource", rolebinding.ResourceSchema().Type(), &rolebinding.RoleBindingModel{}},
		{"trustzone resource", trustzone.ResourceSchema(ctx).Type(), &trustzone.TrustZoneModel{}},
		{"trustzoneserver resource", trustzoneserver.ResourceSchema(ctx).Type(), &trustzoneserver.TrustZoneServerModel{}},

		// Data sources.
		{"apbinding data source", apbinding.DataSourceSchema(ctx).Type(), &apbinding.APBindingModel{}},
		{"attestationpolicy data source", attestationpolicy.DataSourceSchema(ctx).Type(), &attestationpolicy.AttestationPolicyModel{}},
		{"cluster data source", cluster.DataSourceSchema(ctx).Type(), &cluster.ClusterModel{}},
		{"exchangepolicy data source", exchangepolicy.DataSourceSchema().Type(), &exchangepolicy.ExchangePolicyModel{}},
		{"exchangepolicy list data source", exchangepolicy.ListDataSourceSchema().Type(), &exchangepolicy.ExchangePoliciesDataSourceModel{}},
		{"federation data source", federation.DataSourceSchema(ctx).Type(), &federation.FederationModel{}},
		{"organization data source", organization.DataSourceSchema().Type(), &organization.OrganizationModel{}},
		{"trustzone data source", trustzone.DataSourceSchema(ctx).Type(), &trustzone.TrustZoneModel{}},
		{"trustzoneserver data source", trustzoneserver.DataSourceSchema(ctx).Type(), &trustzoneserver.TrustZoneServerModel{}},
		{"trustzoneserver list data source", trustzoneserver.ListDataSourceSchema(ctx).Type(), &trustzoneserver.TrustZoneServersDataSourceModel{}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			value, err := test.schema.ValueFromTerraform(ctx, populatedValue(test.schema.TerraformType(ctx)))
			require.NoError(t, err)

			object, ok := value.(types.Object)
			require.True(t, ok, "expected schema to produce an object value, got %T", value)

			diags := object.As(ctx, test.model, basetypes.ObjectAsOptions{})
			for _, diagnostic := range diags.Errors() {
				t.Errorf("%s: %s", diagnostic.Summary(), diagnostic.Detail())
			}
		})
	}
}

// populatedValue builds a value of the given type in which every object is
// present with null attributes, and every collection holds a single element.
//
// Nested structs are only reflected against when the value containing them is
// non-null, so a wholly null value would skip most of the model. Collections
// need an element for the same reason.
func populatedValue(tfType tftypes.Type) tftypes.Value {
	switch t := tfType.(type) {
	case tftypes.Object:
		attrs := make(map[string]tftypes.Value, len(t.AttributeTypes))
		for name, attrType := range t.AttributeTypes {
			attrs[name] = populatedValue(attrType)
		}
		return tftypes.NewValue(t, attrs)
	case tftypes.List:
		return tftypes.NewValue(t, []tftypes.Value{populatedValue(t.ElementType)})
	case tftypes.Set:
		return tftypes.NewValue(t, []tftypes.Value{populatedValue(t.ElementType)})
	case tftypes.Map:
		return tftypes.NewValue(t, map[string]tftypes.Value{"key": populatedValue(t.ElementType)})
	default:
		return tftypes.NewValue(tfType, nil)
	}
}
