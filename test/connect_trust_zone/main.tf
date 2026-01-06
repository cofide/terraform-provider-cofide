data "cofide_connect_organization" "org" {
  name = "default"
}

resource "cofide_connect_trust_zone" "trust_zone" {
  name         = "test-tz"
  org_id       = data.cofide_connect_organization.org.id
  trust_domain = "test-tz.cofide.dev"
}

data "cofide_connect_trust_zone" "trust_zone" {
  name         = "test-tz"
  org_id       = data.cofide_connect_organization.org.id
  trust_domain = "test-tz.cofide.dev"

  depends_on = [
    cofide_connect_trust_zone.trust_zone
  ]
}

output "trust_zone_id" {
  value = data.cofide_connect_trust_zone.trust_zone.id
}
