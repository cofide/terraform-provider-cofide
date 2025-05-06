resource "cofide_connect_trust_zone" "trust_zone" {
  name         = "test-tz"
  trust_domain = "test-tz.cofide.dev"
}

data "cofide_connect_trust_zone" "trust_zone" {
  name         = "test-tz"
  trust_domain = "test-tz.cofide.dev"

  depends_on = [
    cofide_connect_trust_zone.trust_zone
  ]
}

output "trust_zone_id" {
  value = data.cofide_connect_trust_zone.trust_zone.id
}
