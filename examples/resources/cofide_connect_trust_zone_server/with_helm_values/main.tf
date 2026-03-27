resource "cofide_connect_trust_zone_server" "example" {
  trust_zone_id = var.trust_zone_id
  cluster_id    = var.cluster_id

  helm_values = yamlencode({
    replicaCount = 1
    resources = {
      limits = {
        cpu    = "200m"
        memory = "256Mi"
      }
    }
  })

  connect_k8s_psat_config = {
    audiences                   = ["spire-server"]
    spire_server_spiffe_id_path = "/ns/spire/sa/spire-server"
  }
}
