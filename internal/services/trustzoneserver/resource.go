package trustzoneserver

import (
	"context"

	"github.com/cofide/terraform-provider-cofide/internal/services/trustzoneserver/v1alpha1"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

var (
	_ resource.Resource                = &TrustZoneServerResource{}
	_ resource.ResourceWithImportState = &TrustZoneServerResource{}
)

type TrustZoneServerResource struct {
	v1alpha1.TrustZoneServerResource
}

func NewResource() resource.Resource {
	return &TrustZoneServerResource{TrustZoneServerResource: v1alpha1.TrustZoneServerResource{}}
}

func (t *TrustZoneServerResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_connect_trust_zone_server"
}

func (t *TrustZoneServerResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	t.TrustZoneServerResource.Schema(ctx, req, resp)
	resp.Schema.DeprecationMessage = "Use cofide_connect_trust_zone_server_v1alpha1 instead. This name is frozen on the v1alpha1 API."
	// The above is shown during terraform plan/apply, the below is shown in the generated docs.
	resp.Schema.MarkdownDescription = "~> **Deprecated:** " + resp.Schema.DeprecationMessage + "\n\n" + resp.Schema.MarkdownDescription
}
