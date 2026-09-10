package rolebinding

import (
	"context"

	"github.com/cofide/terraform-provider-cofide/internal/services/rolebinding/v1alpha1"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

var (
	_ resource.Resource                = &RleBindingResource{}
	_ resource.ResourceWithImportState = &RleBindingResource{}
)

type RleBindingResource struct {
	v1alpha1.RoleBindingResource
}

func NewResource() resource.Resource {
	return &RleBindingResource{RoleBindingResource: v1alpha1.RoleBindingResource{}}
}

func (t *RleBindingResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_connect_role_binding"
}

func (t *RleBindingResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	t.RoleBindingResource.Schema(ctx, req, resp)
	resp.Schema.DeprecationMessage = "Use cofide_connect_role_binding_v1alpha1 instead. This name is frozen on the v1alpha1 API."
	// The above is shown during terraform plan/apply, the below is shown in the generated docs.
	resp.Schema.MarkdownDescription = "~> **Deprecated:** " + resp.Schema.DeprecationMessage + "\n\n" + resp.Schema.MarkdownDescription
}
