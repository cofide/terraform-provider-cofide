package cluster

import (
	"context"

	"github.com/cofide/terraform-provider-cofide/internal/services/cluster/v1alpha1"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

var (
	_ resource.Resource                   = &ClusterResource{}
	_ resource.ResourceWithImportState    = &ClusterResource{}
	_ resource.ResourceWithValidateConfig = &ClusterResource{}
)

type ClusterResource struct {
	v1alpha1.ClusterResource
}

func NewResource() resource.Resource {
	return &ClusterResource{ClusterResource: v1alpha1.ClusterResource{}}
}

func (t *ClusterResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_connect_cluster"
}

func (t *ClusterResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	t.ClusterResource.Schema(ctx, req, resp)
	resp.Schema.DeprecationMessage = "Use cofide_connect_cluster_v1alpha1 instead. This name is frozen on the v1alpha1 API."
	// The above is shown during terraform plan/apply, the below is shown in the generated docs.
	resp.Schema.MarkdownDescription = "~> **Deprecated:** " + resp.Schema.DeprecationMessage + "\n\n" + resp.Schema.MarkdownDescription
}
