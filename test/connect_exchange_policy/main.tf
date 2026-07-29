data "cofide_connect_organization" "org" {
  name = "default"
}

resource "cofide_connect_trust_zone" "trust_zone" {
  name         = "ep-tz"
  org_id       = data.cofide_connect_organization.org.id
  trust_domain = "ep-tz.cofide.dev"
}

resource "cofide_connect_exchange_policy" "allow_policy" {
  name          = "test-ep-allow"
  trust_zone_id = cofide_connect_trust_zone.trust_zone.id
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

resource "cofide_connect_exchange_policy" "deny_policy" {
  name          = "test-ep-deny"
  trust_zone_id = cofide_connect_trust_zone.trust_zone.id
  action        = "DENY"

  subject_identity = [
    { exact = "spiffe://ep-tz.cofide.dev/untrusted-workload" }
  ]
}

resource "cofide_connect_exchange_policy" "outbound_identity_policy" {
  name          = "test-ep-outbound-identity"
  trust_zone_id = cofide_connect_trust_zone.trust_zone.id
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

resource "cofide_connect_exchange_policy" "outbound_issuer_policy" {
  name          = "test-ep-outbound-issuer"
  trust_zone_id = cofide_connect_trust_zone.trust_zone.id
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

resource "cofide_connect_exchange_policy" "minimal_policy" {
  name          = "test-ep-minimal"
  trust_zone_id = cofide_connect_trust_zone.trust_zone.id
}

data "cofide_connect_exchange_policy" "allow_policy" {
  id = cofide_connect_exchange_policy.allow_policy.id
}

data "cofide_connect_exchange_policy" "outbound_identity_policy" {
  id = cofide_connect_exchange_policy.outbound_identity_policy.id
}

data "cofide_connect_exchange_policy" "outbound_issuer_policy" {
  id = cofide_connect_exchange_policy.outbound_issuer_policy.id
}

data "cofide_connect_exchange_policies" "by_trust_zone" {
  trust_zone_id = cofide_connect_trust_zone.trust_zone.id

  depends_on = [
    cofide_connect_exchange_policy.allow_policy,
    cofide_connect_exchange_policy.deny_policy,
    cofide_connect_exchange_policy.outbound_identity_policy,
    cofide_connect_exchange_policy.outbound_issuer_policy,
    cofide_connect_exchange_policy.minimal_policy,
  ]
}

output "allow_policy_id" {
  value = data.cofide_connect_exchange_policy.allow_policy.id
}

output "exchange_policy_ids" {
  value = [for p in data.cofide_connect_exchange_policies.by_trust_zone.exchange_policies : p.id]
}

# outbound_identity, as read back from the resource state.
output "outbound_identity_resource" {
  value = cofide_connect_exchange_policy.outbound_identity_policy.outbound_identity
}

# outbound_identity, as returned by the single-policy data source.
output "outbound_identity_data_source" {
  value = data.cofide_connect_exchange_policy.outbound_identity_policy.outbound_identity
}

# outbound_identity, as returned by the list data source.
output "outbound_identity_from_list" {
  value = one([
    for p in data.cofide_connect_exchange_policies.by_trust_zone.exchange_policies :
    p.outbound_identity
    if p.id == cofide_connect_exchange_policy.outbound_identity_policy.id
  ])
}

# outbound_issuer, as read back from the resource state.
output "outbound_issuer_resource" {
  value = cofide_connect_exchange_policy.outbound_issuer_policy.outbound_issuer
}

# outbound_issuer, as returned by the single-policy data source.
output "outbound_issuer_data_source" {
  value = data.cofide_connect_exchange_policy.outbound_issuer_policy.outbound_issuer
}

# outbound_issuer, as returned by the list data source.
output "outbound_issuer_from_list" {
  value = one([
    for p in data.cofide_connect_exchange_policies.by_trust_zone.exchange_policies :
    p.outbound_issuer
    if p.id == cofide_connect_exchange_policy.outbound_issuer_policy.id
  ])
}

# Policies without an outbound issuer should report a null issuer.
output "allow_policy_outbound_issuer" {
  value = data.cofide_connect_exchange_policy.allow_policy.outbound_issuer
}

output "allow_policy_outbound_identity" {
  value = data.cofide_connect_exchange_policy.allow_policy.outbound_identity
}
