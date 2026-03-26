resource "cofide_connect_trust_zone_server" "example" {
  trust_zone_id = var.trust_zone_id
  cluster_id    = var.cluster_id
}
