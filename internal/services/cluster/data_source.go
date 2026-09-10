package cluster

import (
	"context"

	"github.com/cofide/terraform-provider-cofide/internal/services/cluster/v1alpha1"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

type ClusterDataSource struct {
	v1alpha1.ClusterDataSource
}

var _ datasource.DataSourceWithConfigure = (*ClusterDataSource)(nil)

func NewDataSource() datasource.DataSource {
	return &ClusterDataSource{ClusterDataSource: v1alpha1.ClusterDataSource{}}
}

func (d *ClusterDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_connect_cluster"
}

func (t *ClusterDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	t.ClusterDataSource.Schema(ctx, req, resp)
	resp.Schema.DeprecationMessage = "Use cofide_connect_cluster_v1alpha1 instead. This name is frozen on the v1alpha1 API."
	// The above is shown during terraform plan/apply, the below is shown in the generated docs.
	resp.Schema.MarkdownDescription = "~> **Deprecated:** " + resp.Schema.DeprecationMessage + "\n\n" + resp.Schema.MarkdownDescription
}
