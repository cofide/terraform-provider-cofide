package v1alpha1

import (
	"context"
	"testing"

	rolebindingpb "github.com/cofide/cofide-api-sdk/gen/go/proto/role_binding/v1alpha1"
	sdkclient "github.com/cofide/cofide-api-sdk/pkg/connect/client"
	rolebindingclient "github.com/cofide/cofide-api-sdk/pkg/connect/client/rolebinding/v1alpha1"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// fakeClientSet implements only RoleBindingV1Alpha1; any other method panics
// via the nil embedded interface.
type fakeClientSet struct {
	sdkclient.ClientSet
	roleBindings rolebindingclient.RoleBindingClient
}

func (f *fakeClientSet) RoleBindingV1Alpha1() rolebindingclient.RoleBindingClient {
	return f.roleBindings
}

// fakeRoleBindingClient implements only GetRoleBinding.
type fakeRoleBindingClient struct {
	rolebindingclient.RoleBindingClient
	roleBinding *rolebindingpb.RoleBinding
}

func (f *fakeRoleBindingClient) GetRoleBinding(_ context.Context, _ string) (*rolebindingpb.RoleBinding, error) {
	return f.roleBinding, nil
}

// TestImportThenRead exercises the import flow: ImportState sets only the ID,
// leaving every other attribute (including the `resource` object) null, and
// Read must then populate the rest from the API.
func TestImportThenRead(t *testing.T) {
	ctx := context.Background()

	r := &RoleBindingResource{
		client: &fakeClientSet{
			roleBindings: &fakeRoleBindingClient{
				roleBinding: &rolebindingpb.RoleBinding{
					Id:     "rb-123",
					RoleId: "admin",
					Resource: &rolebindingpb.Resource{
						Type: "Organization",
						Id:   "org-123",
					},
					Principal: &rolebindingpb.RoleBinding_User{
						User: &rolebindingpb.User{Subject: "user-123"},
					},
				},
			},
		},
	}

	schemaResp := &resource.SchemaResponse{}
	r.Schema(ctx, resource.SchemaRequest{}, schemaResp)
	require.False(t, schemaResp.Diagnostics.HasError())

	importResp := &resource.ImportStateResponse{
		State: tfsdk.State{
			Schema: schemaResp.Schema,
			Raw:    tftypes.NewValue(schemaResp.Schema.Type().TerraformType(ctx), nil),
		},
	}
	r.ImportState(ctx, resource.ImportStateRequest{ID: "rb-123"}, importResp)
	require.False(t, importResp.Diagnostics.HasError(), "%v", importResp.Diagnostics)

	readResp := &resource.ReadResponse{State: importResp.State}
	r.Read(ctx, resource.ReadRequest{State: importResp.State}, readResp)
	require.False(t, readResp.Diagnostics.HasError(), "%v", readResp.Diagnostics)

	var got RoleBindingModel
	require.False(t, readResp.State.Get(ctx, &got).HasError())
	assert.Equal(t, "rb-123", got.ID.ValueString())
	assert.Equal(t, "admin", got.RoleID.ValueString())
	require.NotNil(t, got.Resource)
	assert.Equal(t, "Organization", got.Resource.Type.ValueString())
	assert.Equal(t, "org-123", got.Resource.ID.ValueString())
	require.NotNil(t, got.User)
	assert.Equal(t, "user-123", got.User.Subject.ValueString())
	assert.Nil(t, got.Group)
}
