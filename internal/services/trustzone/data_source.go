package trustzone

import (
	"context"

	"github.com/cofide/terraform-provider-cofide/internal/services/trustzone/v1alpha1"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

type TrustZoneDataSource struct {
	v1alpha1.TrustZoneDataSource
}

var _ datasource.DataSourceWithConfigure = (*TrustZoneDataSource)(nil)

func NewDataSource() datasource.DataSource {
	return &TrustZoneDataSource{TrustZoneDataSource: v1alpha1.TrustZoneDataSource{}}
}

func (d *TrustZoneDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_connect_trust_zone"
}

func (t *TrustZoneDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	t.TrustZoneDataSource.Schema(ctx, req, resp)
	resp.Schema.DeprecationMessage = "Use cofide_connect_trust_zone_v1alpha1 instead. This name is frozen on the v1alpha1 API."
	// The above is shown during terraform plan/apply, the below is shown in the generated docs.
	resp.Schema.MarkdownDescription = "~> **Deprecated:** " + resp.Schema.DeprecationMessage + "\n\n" + resp.Schema.MarkdownDescription
}
