package trustzone

import (
	"context"

	"github.com/cofide/terraform-provider-cofide/internal/services/trustzone/v1alpha1"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

var (
	_ resource.Resource                   = &TrustZoneResource{}
	_ resource.ResourceWithImportState    = &TrustZoneResource{}
	_ resource.ResourceWithValidateConfig = &TrustZoneResource{}
)

type TrustZoneResource struct {
	v1alpha1.TrustZoneResource
}

func NewResource() resource.Resource {
	return &TrustZoneResource{TrustZoneResource: v1alpha1.TrustZoneResource{}}
}

func (t *TrustZoneResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_connect_trust_zone"
}

func (t *TrustZoneResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	t.TrustZoneResource.Schema(ctx, req, resp)
	resp.Schema.DeprecationMessage = "Use cofide_connect_trust_zone_v1alpha1 instead. This name is frozen on the v1alpha1 API."
	// The above is shown during terraform plan/apply, the below is shown in the generated docs.
	resp.Schema.MarkdownDescription = "~> **Deprecated:** " + resp.Schema.DeprecationMessage + "\n\n" + resp.Schema.MarkdownDescription
}
