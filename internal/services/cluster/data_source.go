package cluster

import (
	"context"
	"fmt"

	clustersvcpb "github.com/cofide/cofide-api-sdk/gen/go/proto/connect/cluster_service/v1alpha1"
	sdkclient "github.com/cofide/cofide-api-sdk/pkg/connect/client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

type ClusterDataSource struct {
	client sdkclient.ClientSet
}

var _ datasource.DataSourceWithConfigure = (*ClusterDataSource)(nil)

func NewDataSource() datasource.DataSource {
	return &ClusterDataSource{}
}

func (c *ClusterDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_connect_cluster"
}

func (c *ClusterDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	c.client = client
}

func (c *ClusterDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config ClusterModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	filter := &clustersvcpb.ListClustersRequest_Filter{
		Name:        config.Name.ValueStringPointer(),
		OrgId:       config.OrgID.ValueStringPointer(),
		TrustZoneId: config.TrustZoneID.ValueStringPointer(),
	}

	clusters, err := c.client.ClusterV1Alpha1().ListClusters(ctx, filter)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading cluster",
			fmt.Sprintf("Could not list clusters: %s", err),
		)

		return
	}

	if len(clusters) == 0 {
		resp.Diagnostics.AddError(
			"Error reading cluster",
			"No matching cluster found",
		)

		return
	}

	if len(clusters) > 1 {
		resp.Diagnostics.AddError(
			"Error reading cluster",
			"Multiple clusters found",
		)

		return
	}

	cluster := clusters[0]

	if err := validateHelmValues(config.ExtraHelmValues, cluster.GetExtraHelmValues()); err != nil {
		resp.Diagnostics.AddError("Inconsistent Helm values", err.Error())
		return
	}

	newState, err := protoToModel(cluster)
	if err != nil {
		resp.Diagnostics.AddError("Error converting proto to model", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &newState)...)
}
