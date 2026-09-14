resource "cofide_connect_federation_v1alpha1" "example" {
  trust_zone_id        = var.trust_zone_id
  remote_trust_zone_id = var.remote_trust_zone_id
}
