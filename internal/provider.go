package internal

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	sdkclient "github.com/cofide/cofide-api-sdk/pkg/connect/client"
	"github.com/cofide/terraform-provider-cofide/internal/client"
	"github.com/cofide/terraform-provider-cofide/internal/consts"
	"github.com/cofide/terraform-provider-cofide/internal/services/apbinding"
	"github.com/cofide/terraform-provider-cofide/internal/services/attestationpolicy"
	"github.com/cofide/terraform-provider-cofide/internal/services/cluster"
	"github.com/cofide/terraform-provider-cofide/internal/services/federation"
	"github.com/cofide/terraform-provider-cofide/internal/services/organization"
	"github.com/cofide/terraform-provider-cofide/internal/services/rolebinding"
	"github.com/cofide/terraform-provider-cofide/internal/services/trustzone"
)

var _ provider.Provider = &CofideProvider{}

func NewProvider(version string) func() provider.Provider {
	return func() provider.Provider {
		return &CofideProvider{
			version: version,
		}
	}
}

// CofideProvider defines the provider implementation.
type CofideProvider struct {
	Client  sdkclient.ClientSet
	version string
}

// CofideProviderModel describes the provider data model.
type CofideProviderModel struct {
	APIToken           types.String `tfsdk:"api_token"`
	ConnectURL         types.String `tfsdk:"connect_url"`
	InsecureSkipVerify types.Bool   `tfsdk:"insecure_skip_verify"`
}

func (p *CofideProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "cofide"
	resp.Version = p.version
}

func (p *CofideProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "This project is the official Terraform provider for Cofide.",
		Attributes: map[string]schema.Attribute{
			"api_token": schema.StringAttribute{
				Description: fmt.Sprintf("API token used to communicate with the Cofide Connect API. Alternatively, can be configured using the `%s` environment variable.", consts.APITokenEnvVarKey),
				Optional:    true,
				Sensitive:   true,
			},
			"connect_url": schema.StringAttribute{
				Description: fmt.Sprintf("Cofide Connect service URL. Alternatively, can be configured using the `%s` environment variable.", consts.ConnectURLEnvVarKey),
				Optional:    true,
			},
			"insecure_skip_verify": schema.BoolAttribute{
				Description: fmt.Sprintf("Skip TLS certificate verification (should only be used for local testing). Alternatively, can be configured using the `%s` environment variable.", consts.InsecureSkipVerifyEnvVar),
				Optional:    true,
			},
		},
	}
}

func (p *CofideProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	if req.ClientCapabilities.DeferralAllowed && !req.Config.Raw.IsFullyKnown() {
		resp.Deferred = &provider.Deferred{
			Reason: provider.DeferredReasonProviderConfigUnknown,
		}
	}

	var config CofideProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Attempts to get configuration from environment variables if not provided in the provider block.
	apiToken := config.APIToken.ValueString()
	if apiToken == "" {
		apiToken = os.Getenv(consts.APITokenEnvVarKey)
	}

	connectURL := config.ConnectURL.ValueString()
	if connectURL == "" {
		connectURL = os.Getenv(consts.ConnectURLEnvVarKey)
	}

	insecureSkipVerify := config.InsecureSkipVerify.ValueBool()
	if config.InsecureSkipVerify.IsNull() || config.InsecureSkipVerify.IsUnknown() {
		if envVal, ok := os.LookupEnv(consts.InsecureSkipVerifyEnvVar); ok {
			if parsed, err := strconv.ParseBool(envVal); err == nil {
				insecureSkipVerify = parsed
			}
		}
	}

	if apiToken == "" {
		resp.Diagnostics.AddError(
			"Missing API Token Configuration",
			"API token must be specified either in provider configuration or via COFIDE_API_TOKEN environment variable",
		)
		return
	}

	if connectURL == "" {
		resp.Diagnostics.AddError(
			"Missing Connect URL Configuration",
			"Connect URL must be specified in provider configuration or via the COFIDE_CONNECT_URL environment variable",
		)
		return
	}

	log := hclog.New(&hclog.LoggerOptions{
		Name: "cofide",
	})

	client, err := client.NewTLSClient(connectURL, apiToken, insecureSkipVerify, log, p.version)
	if err != nil {
		resp.Diagnostics.AddError("Failed to create TLS client", err.Error())
		return
	}

	p.Client = client

	resp.DataSourceData = p.Client
	resp.ResourceData = p.Client

	tflog.Debug(ctx, "Configure method completed successfully")
}

func (p *CofideProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		attestationpolicy.NewResource,
		apbinding.NewResource,
		cluster.NewResource,
		federation.NewResource,
		rolebinding.NewResource,
		trustzone.NewResource,
	}
}

func (p *CofideProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		attestationpolicy.NewDataSource,
		apbinding.NewDataSource,
		cluster.NewDataSource,
		federation.NewDataSource,
		trustzone.NewDataSource,
		organization.NewDataSource,
	}
}
