data "cofide_connect_organization" "org" {
  name = "default"
}

resource "cofide_connect_trust_zone" "trust_zone_a" {
  name         = "test-tz-a"
  org_id       = data.cofide_connect_organization.org.id
  trust_domain = "test-tz-a.cofide.dev"
}

resource "cofide_connect_trust_zone" "trust_zone_b" {
  name         = "test-tz-b"
  org_id       = data.cofide_connect_organization.org.id
  trust_domain = "test-tz-b.cofide.dev"
}

resource "cofide_connect_federation" "federation" {
  org_id               = cofide_connect_trust_zone.trust_zone_a.org_id
  trust_zone_id        = cofide_connect_trust_zone.trust_zone_a.id
  remote_trust_zone_id = cofide_connect_trust_zone.trust_zone_b.id
}
