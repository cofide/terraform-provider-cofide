---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "cofide_connect_cluster Data Source - terraform-provider-cofide"
subcategory: ""
description: |-
  Provides information about a cluster resource.
---

# cofide_connect_cluster (Data Source)

Provides information about a cluster resource.



<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String) The name of the cluster.
- `org_id` (String) The ID of the organisation.

### Optional

- `trust_zone_id` (String) The ID of the associated trust zone.

### Read-Only

- `external_server` (Boolean) Whether or not the SPIRE server runs externally.
- `extra_helm_values` (String) The extra Helm values to provide to the cluster.
- `id` (String) The ID of the cluster.
- `kubernetes_context` (String) The Kubernetes context of the cluster.
- `oidc_issuer_ca_cert` (String) The CA certificate (base64-encoded) to validate the cluster's OIDC issuer URL.
- `oidc_issuer_url` (String) The OIDC issuer URL of the cluster.
- `profile` (String) The Cofide profile used by the cluster.
- `trust_provider` (Attributes) The trust provider of the cluster. (see [below for nested schema](#nestedatt--trust_provider))

<a id="nestedatt--trust_provider"></a>
### Nested Schema for `trust_provider`

Read-Only:

- `kind` (String) The kind of trust provider.
