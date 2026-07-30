package internal_test

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/require"

	"github.com/cofide/terraform-provider-cofide/internal/services/apbinding"
	"github.com/cofide/terraform-provider-cofide/internal/services/attestationpolicy"
	"github.com/cofide/terraform-provider-cofide/internal/services/cluster"
	"github.com/cofide/terraform-provider-cofide/internal/services/exchangepolicy"
	"github.com/cofide/terraform-provider-cofide/internal/services/federation"
	"github.com/cofide/terraform-provider-cofide/internal/services/trustzone"
	"github.com/cofide/terraform-provider-cofide/internal/services/trustzoneserver"
)

// TestSchemaShapesMatch asserts that each resource schema and the data source
// schema(s) for the same resource describe an identical object shape.
//
// A resource and its data sources share a single *Model struct, so a field
// present in one schema but not the other makes reading fail at runtime with a
// Value Conversion Error. The two schemas differ legitimately only in
// Required/Optional/Computed and plan modifiers, neither of which affects the
// shape compared here.
func TestSchemaShapesMatch(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name       string
		resource   attr.Type
		dataSource attr.Type
	}{
		{
			name:       "apbinding",
			resource:   apbinding.ResourceSchema(ctx).Type(),
			dataSource: apbinding.DataSourceSchema(ctx).Type(),
		},
		{
			name:       "attestationpolicy",
			resource:   attestationpolicy.ResourceSchema(ctx).Type(),
			dataSource: attestationpolicy.DataSourceSchema(ctx).Type(),
		},
		{
			name:       "cluster",
			resource:   cluster.ResourceSchema(ctx).Type(),
			dataSource: cluster.DataSourceSchema(ctx).Type(),
		},
		{
			name:       "exchangepolicy",
			resource:   exchangepolicy.ResourceSchema().Type(),
			dataSource: exchangepolicy.DataSourceSchema().Type(),
		},
		{
			name:       "federation",
			resource:   federation.ResourceSchema(ctx).Type(),
			dataSource: federation.DataSourceSchema(ctx).Type(),
		},
		{
			name:       "trustzone",
			resource:   trustzone.ResourceSchema(ctx).Type(),
			dataSource: trustzone.DataSourceSchema(ctx).Type(),
		},
		{
			name:       "trustzoneserver",
			resource:   trustzoneserver.ResourceSchema(ctx).Type(),
			dataSource: trustzoneserver.DataSourceSchema(ctx).Type(),
		},
		// List data sources nest the same model under a list attribute, so
		// compare the element type rather than the top-level schema.
		{
			name:       "exchangepolicy list",
			resource:   exchangepolicy.ResourceSchema().Type(),
			dataSource: listElementType(t, exchangepolicy.ListDataSourceSchema().Type(), "exchange_policies"),
		},
		{
			name:       "trustzoneserver list",
			resource:   trustzoneserver.ResourceSchema(ctx).Type(),
			dataSource: listElementType(t, trustzoneserver.ListDataSourceSchema(ctx).Type(), "trust_zone_servers"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if !test.resource.Equal(test.dataSource) {
				t.Errorf("resource and data source schemas describe different shapes\nresource:    %s\ndata source: %s",
					test.resource, test.dataSource)
			}
		})
	}
}

// listElementType returns the element type of the named list attribute of an
// object type.
func listElementType(t *testing.T, schemaType attr.Type, attrName string) attr.Type {
	t.Helper()

	object, ok := schemaType.(types.ObjectType)
	require.True(t, ok, "expected an object type, got %s", schemaType)

	attribute, ok := object.AttributeTypes()[attrName]
	require.True(t, ok, "no %q attribute in %s", attrName, schemaType)

	list, ok := attribute.(types.ListType)
	require.True(t, ok, "expected %q to be a list type, got %s", attrName, attribute)

	return list.ElementType()
}
