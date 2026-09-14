package federation

import (
	"context"

	"github.com/cofide/terraform-provider-cofide/internal/services/federation/v1alpha1"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

var (
	_ resource.Resource                   = &FederationResource{}
	_ resource.ResourceWithImportState    = &FederationResource{}
	_ resource.ResourceWithValidateConfig = &FederationResource{}
)

type FederationResource struct {
	v1alpha1.FederationResource
}

func NewResource() resource.Resource {
	return &FederationResource{FederationResource: v1alpha1.FederationResource{}}
}

func (t *FederationResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_connect_federation"
}

func (t *FederationResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	t.FederationResource.Schema(ctx, req, resp)
	resp.Schema.DeprecationMessage = "Use cofide_connect_federation_v1alpha1 instead. This name is frozen on the v1alpha1 API."
	// The above is shown during terraform plan/apply, the below is shown in the generated docs.
	resp.Schema.MarkdownDescription = "~> **Deprecated:** " + resp.Schema.DeprecationMessage + "\n\n" + resp.Schema.MarkdownDescription
}
