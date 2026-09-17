package exchangepolicy

import (
	"context"

	"github.com/cofide/terraform-provider-cofide/internal/services/exchangepolicy/v1alpha1"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

type ExchangePoliciesDataSource struct {
	v1alpha1.ExchangePoliciesDataSource
}

var _ datasource.DataSourceWithConfigure = (*ExchangePoliciesDataSource)(nil)

func NewListDataSource() datasource.DataSource {
	return &ExchangePoliciesDataSource{ExchangePoliciesDataSource: v1alpha1.ExchangePoliciesDataSource{}}
}

func (d *ExchangePoliciesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_connect_exchange_policies"
}

func (t *ExchangePoliciesDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	t.ExchangePoliciesDataSource.Schema(ctx, req, resp)
	resp.Schema.DeprecationMessage = "Use cofide_connect_exchange_policies_v1alpha1 instead. This name is frozen on the v1alpha1 API."
	// The above is shown during terraform plan/apply, the below is shown in the generated docs.
	resp.Schema.MarkdownDescription = "~> **Deprecated:** " + resp.Schema.DeprecationMessage + "\n\n" + resp.Schema.MarkdownDescription
}
