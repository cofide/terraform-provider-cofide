package organization

import (
	"context"

	"github.com/cofide/terraform-provider-cofide/internal/services/organization/v1alpha1"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

type OrganizationDataSource struct {
	v1alpha1.OrganizationDataSource
}

var _ datasource.DataSourceWithConfigure = (*OrganizationDataSource)(nil)

func NewDataSource() datasource.DataSource {
	return &OrganizationDataSource{OrganizationDataSource: v1alpha1.OrganizationDataSource{}}
}

func (d *OrganizationDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_connect_organization"
}

func (t *OrganizationDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	t.OrganizationDataSource.Schema(ctx, req, resp)
	resp.Schema.DeprecationMessage = "Use cofide_connect_organization_v1alpha1 instead. This name is frozen on the v1alpha1 API."
	// The above is shown during terraform plan/apply, the below is shown in the generated docs.
	resp.Schema.MarkdownDescription = "~> **Deprecated:** " + resp.Schema.DeprecationMessage + "\n\n" + resp.Schema.MarkdownDescription
}
