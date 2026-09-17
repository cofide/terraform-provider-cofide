package apbinding

import (
	"context"

	"github.com/cofide/terraform-provider-cofide/internal/services/apbinding/v1alpha1"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

type APBindingDataSource struct {
	v1alpha1.APBindingDataSource
}

var _ datasource.DataSourceWithConfigure = (*APBindingDataSource)(nil)

func NewDataSource() datasource.DataSource {
	return &APBindingDataSource{APBindingDataSource: v1alpha1.APBindingDataSource{}}
}

func (d *APBindingDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_connect_ap_binding"
}

func (t *APBindingDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	t.APBindingDataSource.Schema(ctx, req, resp)
	resp.Schema.DeprecationMessage = "Use cofide_connect_ap_binding_v1alpha1 instead. This name is frozen on the v1alpha1 API."
	// The above is shown during terraform plan/apply, the below is shown in the generated docs.
	resp.Schema.MarkdownDescription = "~> **Deprecated:** " + resp.Schema.DeprecationMessage + "\n\n" + resp.Schema.MarkdownDescription
}
