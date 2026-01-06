package organization

import (
	"context"
	"fmt"

	organizationsvcpb "github.com/cofide/cofide-api-sdk/gen/go/proto/connect/organization_service/v1alpha1"
	organizationpb "github.com/cofide/cofide-api-sdk/gen/go/proto/organization/v1alpha1"
	sdkclient "github.com/cofide/cofide-api-sdk/pkg/connect/client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

func NewDataSource() datasource.DataSource {
	return &OrganizationDataSource{}
}

// OrganizationDataSource defines the data source implementation.
type OrganizationDataSource struct {
	client sdkclient.ClientSet
}

func (d *OrganizationDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_connect_organization"
}

func (d *OrganizationDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = getOrganizationDataSourceSchema()
}

func (d *OrganizationDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Prevent panic if the provider has no client configured
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(sdkclient.ClientSet)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			"Expected sdkclient.ClientSet, got: %T. Please report this issue to the provider developers.",
		)

		return
	}

	d.client = client
}

func (d *OrganizationDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config OrganizationModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Organization datasource: getting organization", map[string]any{
		"org_name": config.Name.ValueString(),
	})

	org, err := getOrganizationByName(ctx, d.client, config.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", err.Error())
		return
	}

	// Save data into Terraform state
	state := protoToModel(org)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func getOrganizationByName(ctx context.Context, client sdkclient.ClientSet, name string) (*organizationpb.Organization, error) {
	filter := &organizationsvcpb.ListOrganizationsRequest_Filter{
		Name: &name,
	}
	resp, err := client.OrganizationV1Alpha1().ListOrganizations(ctx, filter)
	if err != nil {
		return nil, err
	}
	if len(resp) == 0 {
		return nil, fmt.Errorf("organization with name '%s' not found", name)
	}
	if len(resp) > 1 {
		return nil, fmt.Errorf("multiple organizations with name '%s' found", name)
	}
	return resp[0], nil
}
