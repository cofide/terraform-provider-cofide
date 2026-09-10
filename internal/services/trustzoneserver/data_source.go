package trustzoneserver

import (
	"context"

	"github.com/cofide/terraform-provider-cofide/internal/services/trustzoneserver/v1alpha1"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

type TrustZoneServerDataSource struct {
	v1alpha1.TrustZoneServerDataSource
}

var _ datasource.DataSourceWithConfigure = (*TrustZoneServerDataSource)(nil)

func NewDataSource() datasource.DataSource {
	return &TrustZoneServerDataSource{TrustZoneServerDataSource: v1alpha1.TrustZoneServerDataSource{}}
}

func (d *TrustZoneServerDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_connect_trust_zone_server"
}

func (t *TrustZoneServerDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	t.TrustZoneServerDataSource.Schema(ctx, req, resp)
	resp.Schema.DeprecationMessage = "Use cofide_connect_trust_zone_server_v1alpha1 instead. This name is frozen on the v1alpha1 API."
	// The above is shown during terraform plan/apply, the below is shown in the generated docs.
	resp.Schema.MarkdownDescription = "~> **Deprecated:** " + resp.Schema.DeprecationMessage + "\n\n" + resp.Schema.MarkdownDescription
}
