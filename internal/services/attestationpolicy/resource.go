package attestationpolicy

import (
	"context"

	"github.com/cofide/terraform-provider-cofide/internal/services/attestationpolicy/v1alpha1"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

var (
	_ resource.Resource                     = &AttestationPolicyResource{}
	_ resource.ResourceWithImportState      = &AttestationPolicyResource{}
	_ resource.ResourceWithConfigValidators = &AttestationPolicyResource{}
)

type AttestationPolicyResource struct {
	v1alpha1.AttestationPolicyResource
}

func NewResource() resource.Resource {
	return &AttestationPolicyResource{AttestationPolicyResource: v1alpha1.AttestationPolicyResource{}}
}

func (t *AttestationPolicyResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_connect_attestation_policy"
}

func (t *AttestationPolicyResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	t.AttestationPolicyResource.Schema(ctx, req, resp)
	resp.Schema.DeprecationMessage = "Use cofide_connect_attestation_policy_v1alpha1 instead. This name is frozen on the v1alpha1 API."
	// The above is shown during terraform plan/apply, the below is shown in the generated docs.
	resp.Schema.MarkdownDescription = "~> **Deprecated:** " + resp.Schema.DeprecationMessage + "\n\n" + resp.Schema.MarkdownDescription
}
