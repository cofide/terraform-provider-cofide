resource "cofide_connect_cluster" "example" {
  name               = var.name
  trust_zone_id      = var.trust_zone_id
  org_id             = var.org_id
  profile            = "kubernetes"
  kubernetes_context = var.kubernetes_context
  external_server    = false

  trust_provider = {
    kind = "kubernetes"
  }
}
