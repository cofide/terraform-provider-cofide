package cluster

import (
	"context"
	"fmt"

	clustersvcpb "github.com/cofide/cofide-api-sdk/gen/go/proto/connect/cluster_service/v1alpha1"
	sdkclient "github.com/cofide/cofide-api-sdk/pkg/connect/client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"
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

	state := ClusterModel{
		ID:                types.StringValue(cluster.GetId()),
		Name:              types.StringValue(cluster.GetName()),
		OrgID:             types.StringValue(cluster.GetOrgId()),
		TrustZoneID:       types.StringValue(cluster.GetTrustZoneId()),
		KubernetesContext: types.StringValue(cluster.GetKubernetesContext()),
		TrustProvider: &TrustProviderModel{
			Kind: types.StringValue(cluster.GetTrustProvider().GetKind()),
		},
		ExtraHelmValues:  types.StringValue(cluster.GetExtraHelmValues().String()),
		Profile:          types.StringValue(cluster.GetProfile()),
		ExternalServer:   types.BoolValue(cluster.GetExternalServer()),
		OidcIssuerURL:    types.StringValue(cluster.GetOidcIssuerUrl()),
		OidcIssuerCaCert: types.StringValue(string(cluster.GetOidcIssuerCaCert())),
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
