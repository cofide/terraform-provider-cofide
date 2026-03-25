resource "cofide_connect_federation" "example" {
  org_id               = var.org_id
  trust_zone_id        = var.trust_zone_id
  remote_trust_zone_id = var.remote_trust_zone_id
}
