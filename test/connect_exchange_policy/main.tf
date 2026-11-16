data "cofide_connect_organization_v1alpha1" "org" {
  name = "default"
}

resource "cofide_connect_trust_zone_v1alpha1" "trust_zone" {
  name         = "ep-tz"
  org_id       = data.cofide_connect_organization_v1alpha1.org.id
  trust_domain = "ep-tz.cofide.dev"
}

resource "cofide_connect_exchange_policy_v1alpha1" "allow_policy" {
  name          = "test-ep-allow"
  trust_zone_id = cofide_connect_trust_zone_v1alpha1.trust_zone.id
  action        = "ALLOW"

  subject_identity = [
    { glob = "spiffe://ep-tz.cofide.dev/ns/foo/sa/*" },
    { glob = "spiffe://ep-tz.cofide.dev/ns/bar/sa/*" }
  ]

  subject_audience = [
    { exact = "https://audience.ep-tz.cofide.dev" }
  ]

  target_audience = [
    { exact = "https://api.ep-tz.cofide.dev" }
  ]

  outbound_scopes = ["read", "write"]
}

resource "cofide_connect_exchange_policy_v1alpha1" "deny_policy" {
  name          = "test-ep-deny"
  trust_zone_id = cofide_connect_trust_zone_v1alpha1.trust_zone.id
  action        = "DENY"

  subject_identity = [
    { exact = "spiffe://ep-tz.cofide.dev/untrusted-workload" }
  ]
}

resource "cofide_connect_exchange_policy_v1alpha1" "outbound_identity_policy" {
  name          = "test-ep-outbound-identity"
  trust_zone_id = cofide_connect_trust_zone_v1alpha1.trust_zone.id
  action        = "ALLOW"

  # Connect requires subject_issuer to use exact matching when
  # outbound_identity is set.
  subject_issuer = [
    { exact = "https://issuer.ep-tz.cofide.dev" }
  ]

  subject_identity = [
    { exact = "spiffe://ep-tz.cofide.dev/ns/outbound/sa/caller" }
  ]

  outbound_identity = "spiffe://ep-tz.cofide.dev/ns/outbound/sa/exchanged"
}

resource "cofide_connect_exchange_policy_v1alpha1" "outbound_issuer_policy" {
  name          = "test-ep-outbound-issuer"
  trust_zone_id = cofide_connect_trust_zone_v1alpha1.trust_zone.id
  action        = "ALLOW"

  subject_identity = [
    { glob = "spiffe://ep-tz.cofide.dev/ns/outbound/sa/*" }
  ]

  target_audience = [
    { exact = "https://external-api.example.com" }
  ]

  outbound_issuer = {
    oauth_as = {
      grant_type = "client_credentials"
      issuer_url = "https://as.example.com"
      token_url  = "https://as.example.com/oauth2/token"
      audiences  = ["https://external-api.example.com"]
      timeout    = 30
    }
  }
}

resource "cofide_connect_exchange_policy_v1alpha1" "spiffe_issuer_policy" {
  name          = "test-ep-spiffe-issuer"
  trust_zone_id = cofide_connect_trust_zone_v1alpha1.trust_zone.id
  action        = "ALLOW"

  # Connect requires subject_issuer to use exact matching when
  # outbound_identity is set. The SPIFFE outbound issuer derives the issued
  # SVID's SPIFFE ID from outbound_identity.
  subject_issuer = [
    { exact = "https://issuer.ep-tz.cofide.dev" }
  ]

  subject_identity = [
    { exact = "user@example.com" }
  ]

  outbound_identity = "spiffe://ep-tz.cofide.dev/ns/spiffe/sa/exchanged"

  outbound_issuer = {
    spiffe = {}
  }
}

# Credex rejects client assertion policies that set outbound_identity,
# outbound_scopes or target_audience, or that use an inexact subject_issuer.
resource "cofide_connect_exchange_policy_v1alpha1" "client_assertion_policy" {
  name          = "test-ep-client-assertion"
  trust_zone_id = cofide_connect_trust_zone_v1alpha1.trust_zone.id
  action        = "ALLOW"

  subject_issuer = [
    { exact = "https://issuer.ep-tz.cofide.dev" }
  ]

  subject_identity = [
    { exact = "spiffe://ep-tz.cofide.dev/ns/assertion/sa/caller" }
  ]

  outbound_issuer = {
    client_assertion = {
      client_id  = "oauth-client-123"
      issuer_url = "https://as.example.com"
      audiences  = ["https://as.example.com/oauth2/token"]
      timeout    = 30
    }
  }
}

# client_id may be omitted when an external hook supplies it.
resource "cofide_connect_exchange_policy_v1alpha1" "client_assertion_hook_policy" {
  name          = "test-ep-client-assertion-hook"
  trust_zone_id = cofide_connect_trust_zone_v1alpha1.trust_zone.id
  action        = "ALLOW"

  subject_issuer = [
    { exact = "https://issuer.ep-tz.cofide.dev" }
  ]

  subject_identity = [
    { exact = "spiffe://ep-tz.cofide.dev/ns/assertion/sa/hooked" }
  ]

  outbound_issuer = {
    client_assertion = {
      audiences = ["https://as.example.com/oauth2/token"]
    }
  }

  external_hooks = [
    {
      name = "client-id-mapper"
      url  = "https://hooks.ep-tz.cofide.dev/client-id"
      auth = {
        spiffe_mtls = {
          spiffe_id = "spiffe://ep-tz.cofide.dev/ns/hooks/sa/client-id-mapper"
        }
      }
    }
  ]
}

resource "cofide_connect_exchange_policy_v1alpha1" "minimal_policy" {
  name          = "test-ep-minimal"
  trust_zone_id = cofide_connect_trust_zone_v1alpha1.trust_zone.id
}

data "cofide_connect_exchange_policy_v1alpha1" "allow_policy" {
  id = cofide_connect_exchange_policy_v1alpha1.allow_policy.id
}

data "cofide_connect_exchange_policy_v1alpha1" "outbound_identity_policy" {
  id = cofide_connect_exchange_policy_v1alpha1.outbound_identity_policy.id
}

data "cofide_connect_exchange_policy_v1alpha1" "outbound_issuer_policy" {
  id = cofide_connect_exchange_policy_v1alpha1.outbound_issuer_policy.id
}

data "cofide_connect_exchange_policy_v1alpha1" "spiffe_issuer_policy" {
  id = cofide_connect_exchange_policy_v1alpha1.spiffe_issuer_policy.id
}

data "cofide_connect_exchange_policy_v1alpha1" "client_assertion_policy" {
  id = cofide_connect_exchange_policy_v1alpha1.client_assertion_policy.id
}

data "cofide_connect_exchange_policy_v1alpha1" "client_assertion_hook_policy" {
  id = cofide_connect_exchange_policy_v1alpha1.client_assertion_hook_policy.id
}

data "cofide_connect_exchange_policies_v1alpha1" "by_trust_zone" {
  trust_zone_id = cofide_connect_trust_zone_v1alpha1.trust_zone.id

  depends_on = [
    cofide_connect_exchange_policy_v1alpha1.allow_policy,
    cofide_connect_exchange_policy_v1alpha1.deny_policy,
    cofide_connect_exchange_policy_v1alpha1.outbound_identity_policy,
    cofide_connect_exchange_policy_v1alpha1.outbound_issuer_policy,
    cofide_connect_exchange_policy_v1alpha1.spiffe_issuer_policy,
    cofide_connect_exchange_policy_v1alpha1.client_assertion_policy,
    cofide_connect_exchange_policy_v1alpha1.client_assertion_hook_policy,
    cofide_connect_exchange_policy_v1alpha1.minimal_policy,
  ]
}

output "allow_policy_id" {
  value = data.cofide_connect_exchange_policy_v1alpha1.allow_policy.id
}

output "exchange_policy_ids" {
  value = [for p in data.cofide_connect_exchange_policies_v1alpha1.by_trust_zone.exchange_policies : p.id]
}

# outbound_identity, as read back from the resource state.
output "outbound_identity_resource" {
  value = cofide_connect_exchange_policy_v1alpha1.outbound_identity_policy.outbound_identity
}

# outbound_identity, as returned by the single-policy data source.
output "outbound_identity_data_source" {
  value = data.cofide_connect_exchange_policy_v1alpha1.outbound_identity_policy.outbound_identity
}

# outbound_identity, as returned by the list data source.
output "outbound_identity_from_list" {
  value = one([
    for p in data.cofide_connect_exchange_policies_v1alpha1.by_trust_zone.exchange_policies :
    p.outbound_identity
    if p.id == cofide_connect_exchange_policy_v1alpha1.outbound_identity_policy.id
  ])
}

# outbound_issuer, as read back from the resource state.
output "outbound_issuer_resource" {
  value = cofide_connect_exchange_policy_v1alpha1.outbound_issuer_policy.outbound_issuer
}

# outbound_issuer, as returned by the single-policy data source.
output "outbound_issuer_data_source" {
  value = data.cofide_connect_exchange_policy_v1alpha1.outbound_issuer_policy.outbound_issuer
}

# outbound_issuer, as returned by the list data source.
output "outbound_issuer_from_list" {
  value = one([
    for p in data.cofide_connect_exchange_policies_v1alpha1.by_trust_zone.exchange_policies :
    p.outbound_issuer
    if p.id == cofide_connect_exchange_policy_v1alpha1.outbound_issuer_policy.id
  ])
}

# outbound_issuer.spiffe round-trips through the resource, the single-policy
# data source and the list data source.
output "spiffe_issuer_resource" {
  value = cofide_connect_exchange_policy_v1alpha1.spiffe_issuer_policy.outbound_issuer
}

output "spiffe_issuer_data_source" {
  value = data.cofide_connect_exchange_policy_v1alpha1.spiffe_issuer_policy.outbound_issuer
}

output "spiffe_issuer_from_list" {
  value = one([
    for p in data.cofide_connect_exchange_policies_v1alpha1.by_trust_zone.exchange_policies :
    p.outbound_issuer
    if p.id == cofide_connect_exchange_policy_v1alpha1.spiffe_issuer_policy.id
  ])
}

# outbound_issuer.client_assertion round-trips through the resource, the
# single-policy data source and the list data source.
output "client_assertion_issuer_resource" {
  value = cofide_connect_exchange_policy_v1alpha1.client_assertion_policy.outbound_issuer
}

output "client_assertion_issuer_data_source" {
  value = data.cofide_connect_exchange_policy_v1alpha1.client_assertion_policy.outbound_issuer
}

output "client_assertion_issuer_from_list" {
  value = one([
    for p in data.cofide_connect_exchange_policies_v1alpha1.by_trust_zone.exchange_policies :
    p.outbound_issuer
    if p.id == cofide_connect_exchange_policy_v1alpha1.client_assertion_policy.id
  ])
}

# A client assertion policy with no client_id, relying on an external hook.
output "client_assertion_hook_issuer_resource" {
  value = cofide_connect_exchange_policy_v1alpha1.client_assertion_hook_policy.outbound_issuer
}

output "client_assertion_hook_issuer_data_source" {
  value = data.cofide_connect_exchange_policy_v1alpha1.client_assertion_hook_policy.outbound_issuer
}

output "client_assertion_hook_issuer_from_list" {
  value = one([
    for p in data.cofide_connect_exchange_policies_v1alpha1.by_trust_zone.exchange_policies :
    p.outbound_issuer
    if p.id == cofide_connect_exchange_policy_v1alpha1.client_assertion_hook_policy.id
  ])
}

output "client_assertion_hook_names" {
  value = [for h in cofide_connect_exchange_policy_v1alpha1.client_assertion_hook_policy.external_hooks : h.name]
}

# Policies without an outbound issuer should report a null issuer.
output "allow_policy_outbound_issuer" {
  value = data.cofide_connect_exchange_policy_v1alpha1.allow_policy.outbound_issuer
}

output "allow_policy_outbound_identity" {
  value = data.cofide_connect_exchange_policy_v1alpha1.allow_policy.outbound_identity
}
