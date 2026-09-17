package exchangepolicy

import (
	"context"

	"github.com/cofide/terraform-provider-cofide/internal/services/exchangepolicy/v1alpha1"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

type ExchangePolicyDataSource struct {
	v1alpha1.ExchangePolicyDataSource
}

var _ datasource.DataSourceWithConfigure = (*ExchangePolicyDataSource)(nil)

func NewDataSource() datasource.DataSource {
	return &ExchangePolicyDataSource{ExchangePolicyDataSource: v1alpha1.ExchangePolicyDataSource{}}
}

func (d *ExchangePolicyDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_connect_exchange_policy"
}

func (t *ExchangePolicyDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	t.ExchangePolicyDataSource.Schema(ctx, req, resp)
	resp.Schema.DeprecationMessage = "Use cofide_connect_exchange_policy_v1alpha1 instead. This name is frozen on the v1alpha1 API."
	// The above is shown during terraform plan/apply, the below is shown in the generated docs.
	resp.Schema.MarkdownDescription = "~> **Deprecated:** " + resp.Schema.DeprecationMessage + "\n\n" + resp.Schema.MarkdownDescription
}
