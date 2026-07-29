resource "cofide_connect_federation" "example" {
  trust_zone_id        = var.trust_zone_id
  remote_trust_zone_id = var.remote_trust_zone_id
}
