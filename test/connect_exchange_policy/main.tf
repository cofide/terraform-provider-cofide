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

resource "cofide_connect_exchange_policy" "minimal_policy" {
  name          = "test-ep-minimal"
  trust_zone_id = cofide_connect_trust_zone.trust_zone.id
}

data "cofide_connect_exchange_policy" "allow_policy" {
  id = cofide_connect_exchange_policy.allow_policy.id
}

data "cofide_connect_exchange_policies" "by_trust_zone" {
  trust_zone_id = cofide_connect_trust_zone.trust_zone.id

  depends_on = [
    cofide_connect_exchange_policy.allow_policy,
    cofide_connect_exchange_policy.deny_policy,
    cofide_connect_exchange_policy.minimal_policy,
  ]
}

output "allow_policy_id" {
  value = data.cofide_connect_exchange_policy.allow_policy.id
}

output "exchange_policy_ids" {
  value = [for p in data.cofide_connect_exchange_policies.by_trust_zone.exchange_policies : p.id]
}
