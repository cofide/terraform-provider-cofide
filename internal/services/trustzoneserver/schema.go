package trustzoneserver

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ resource.ResourceWithConfigValidators = (*TrustZoneServerResource)(nil)

func ResourceSchema(_ context.Context) schema.Schema {
	return schema.Schema{
		MarkdownDescription: "Manages a Cofide Connect trust zone server. A trust zone server defines how the SPIRE server managing a trust zone should be deployed on a cluster.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The ID of the trust zone server.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"trust_zone_id": schema.StringAttribute{
				Description: "The ID of the trust zone managed by this server. Cannot be changed after creation.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"cluster_id": schema.StringAttribute{
				Description: "The ID of the cluster on which the server should be deployed. Cannot be changed after creation.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"kubernetes_namespace": schema.StringAttribute{
				Description: "The Kubernetes namespace in which the server should be deployed. Set by Cofide Connect if not provided. Cannot be changed after creation.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"kubernetes_service_account": schema.StringAttribute{
				Description: "The name of the Kubernetes service account to deploy with the server. Set by Cofide Connect if not provided. Cannot be changed after creation.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"org_id": schema.StringAttribute{
				Description: "The ID of the organisation. Derived from the trust zone by Cofide Connect.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"helm_values": schema.StringAttribute{
				Description: "Additional Helm values for the SPIRE server Helm chart installation, in YAML format. Use `yamlencode()` to generate from a Terraform map.",
				Optional:    true,
			},
			"status": schema.SingleNestedAttribute{
				Description: "The current lifecycle status of the trust zone server. Set by Cofide Connect.",
				Computed:    true,
				Attributes: map[string]schema.Attribute{
					"status": schema.StringAttribute{
						Description: "The status of the trust zone server (e.g. `TRUST_ZONE_SERVER_STATUS_PROVISIONED`).",
						Computed:    true,
					},
					"last_transition_time": schema.StringAttribute{
						Description: "The time of the last status transition (RFC3339).",
						Computed:    true,
					},
				},
			},
			"connect_k8s_psat_config": schema.SingleNestedAttribute{
				Description: "Configuration for the k8s PSAT node attestor plugin when using a Connect datasource with remote clusters.",
				Optional:    true,
				Attributes: map[string]schema.Attribute{
					"audiences": schema.ListAttribute{
						Description: "Audiences that SPIRE agents in remote clusters can present for node attestation. At least one must be provided if there are remote clusters in the trust zone.",
						Required:    true,
						ElementType: types.StringType,
					},
					"spire_server_spiffe_id_path": schema.StringAttribute{
						Description: "SPIFFE ID path used in the JWT presented by the SPIRE server to the cluster's API server (e.g. `/ns/spire/sa/spire-server`).",
						Required:    true,
					},
				},
			},
		},
	}
}

func (r *TrustZoneServerResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = ResourceSchema(ctx)
}

func (r *TrustZoneServerResource) ConfigValidators(_ context.Context) []resource.ConfigValidator {
	return []resource.ConfigValidator{}
}
