resource "cofide_connect_trust_zone" "example" {
  name         = var.trust_zone_name
  org_id       = var.org_id
  trust_domain = var.trust_domain
}
