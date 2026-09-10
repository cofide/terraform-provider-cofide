package federation

import (
	"context"

	"github.com/cofide/terraform-provider-cofide/internal/services/federation/v1alpha1"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

type FederationDataSource struct {
	v1alpha1.FederationDataSource
}

var _ datasource.DataSourceWithConfigure = (*FederationDataSource)(nil)

func NewDataSource() datasource.DataSource {
	return &FederationDataSource{FederationDataSource: v1alpha1.FederationDataSource{}}
}

func (d *FederationDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_connect_federation"
}

func (t *FederationDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	t.FederationDataSource.Schema(ctx, req, resp)
	resp.Schema.DeprecationMessage = "Use cofide_connect_federation_v1alpha1 instead. This name is frozen on the v1alpha1 API."
	// The above is shown during terraform plan/apply, the below is shown in the generated docs.
	resp.Schema.MarkdownDescription = "~> **Deprecated:** " + resp.Schema.DeprecationMessage + "\n\n" + resp.Schema.MarkdownDescription
}
