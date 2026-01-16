package rolebinding

import (
	"context"
	"fmt"

	rolebindingsvcpb "github.com/cofide/cofide-api-sdk/gen/go/proto/connect/role_binding_service/v1alpha1"
	sdkclient "github.com/cofide/cofide-api-sdk/pkg/connect/client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

type RoleBindingDataSource struct {
	client sdkclient.ClientSet
}

var _ datasource.DataSourceWithConfigure = (*RoleBindingDataSource)(nil)

func NewDataSource() datasource.DataSource {
	return &RoleBindingDataSource{}
}

func (d *RoleBindingDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_connect_role_binding"
}

func (d *RoleBindingDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(sdkclient.ClientSet)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected data source configure type",
			fmt.Sprintf("Expected sdkclient.ClientSet, got: %T", req.ProviderData),
		)
		return
	}

	d.client = client
}

func (d *RoleBindingDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config RoleBindingModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	filter := &rolebindingsvcpb.ListRoleBindingsRequest_Filter{
		RoleId:       config.RoleID.ValueStringPointer(),
		ResourceType: config.Resource.Type.ValueStringPointer(),
		ResourceId:   config.Resource.ID.ValueStringPointer(),
	}

	if config.User != nil {
		filter.UserSubject = config.User.Subject.ValueStringPointer()
	}
	if config.Group != nil {
		filter.GroupClaimValue = config.Group.ClaimValue.ValueStringPointer()
	}

	roleBindings, err := d.client.RoleBindingV1Alpha1().ListRoleBindings(ctx, filter)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading role binding",
			fmt.Sprintf("Could not list role bindings: %s", err),
		)
		return
	}

	if len(roleBindings) == 0 {
		resp.Diagnostics.AddError(
			"Error reading role binding",
			"No matching role binding found",
		)
		return
	}

	if len(roleBindings) > 1 {
		resp.Diagnostics.AddError(
			"Error reading role binding",
			"Multiple role bindings found",
		)
		return
	}

	roleBinding := roleBindings[0]

	if roleBinding == nil {
		resp.Diagnostics.AddError(
			"Error reading role binding",
			"API returned a nil role binding. This is unexpected, please report this issue.",
		)
		return
	}

	state := protoToModel(roleBinding)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (d *RoleBindingDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = dataSourceSchema()
}

func (d *RoleBindingDataSource) ConfigValidators(_ context.Context) []datasource.ConfigValidator {
	return []datasource.ConfigValidator{}
}
