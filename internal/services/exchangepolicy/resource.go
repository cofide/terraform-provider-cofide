package exchangepolicy

import (
	"context"

	"github.com/cofide/terraform-provider-cofide/internal/services/exchangepolicy/v1alpha1"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

var (
	_ resource.Resource                = &ExchangePolicyResource{}
	_ resource.ResourceWithImportState = &ExchangePolicyResource{}
)

type ExchangePolicyResource struct {
	v1alpha1.ExchangePolicyResource
}

func NewResource() resource.Resource {
	return &ExchangePolicyResource{ExchangePolicyResource: v1alpha1.ExchangePolicyResource{}}
}

func (t *ExchangePolicyResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_connect_exchange_policy"
}

func (t *ExchangePolicyResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	t.ExchangePolicyResource.Schema(ctx, req, resp)
	resp.Schema.DeprecationMessage = "Use cofide_connect_exchange_policy_v1alpha1 instead. This name is frozen on the v1alpha1 API."
	// The above is shown during terraform plan/apply, the below is shown in the generated docs.
	resp.Schema.MarkdownDescription = "~> **Deprecated:** " + resp.Schema.DeprecationMessage + "\n\n" + resp.Schema.MarkdownDescription
}
