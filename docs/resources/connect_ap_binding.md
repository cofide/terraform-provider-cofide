---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "cofide_connect_ap_binding Resource - terraform-provider-cofide"
subcategory: ""
description: |-
  Provides an attestation policy binding resource.
---

# cofide_connect_ap_binding (Resource)

Provides an attestation policy binding resource.



<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `org_id` (String) The ID of the organisation.
- `policy_id` (String) The ID of the attestation policy.
- `trust_zone_id` (String) The ID of the trust zone.

### Optional

- `federations` (List of Object) The list of associated federations. (see [below for nested schema](#nestedatt--federations))

### Read-Only

- `id` (String) The ID of the attestation policy binding.

<a id="nestedatt--federations"></a>
### Nested Schema for `federations`

Optional:

- `trust_zone_id` (String)
