package attestationpolicy

import (
	"context"

	"github.com/cofide/terraform-provider-cofide/internal/services/attestationpolicy/v1alpha1"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

type AttestationPolicyDataSource struct {
	v1alpha1.AttestationPolicyDataSource
}

var _ datasource.DataSourceWithConfigure = (*AttestationPolicyDataSource)(nil)

func NewDataSource() datasource.DataSource {
	return &AttestationPolicyDataSource{AttestationPolicyDataSource: v1alpha1.AttestationPolicyDataSource{}}
}

func (d *AttestationPolicyDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_connect_attestation_policy"
}

func (t *AttestationPolicyDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	t.AttestationPolicyDataSource.Schema(ctx, req, resp)
	resp.Schema.DeprecationMessage = "Use cofide_connect_attestation_policy_v1alpha1 instead. This name is frozen on the v1alpha1 API."
	// The above is shown during terraform plan/apply, the below is shown in the generated docs.
	resp.Schema.MarkdownDescription = "~> **Deprecated:** " + resp.Schema.DeprecationMessage + "\n\n" + resp.Schema.MarkdownDescription
}
