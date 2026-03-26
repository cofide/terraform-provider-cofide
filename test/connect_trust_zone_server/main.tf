data "cofide_connect_organization" "org" {
  name = "default"
}

resource "cofide_connect_trust_zone" "trust_zone" {
  name         = "tzserver-tz"
  org_id       = data.cofide_connect_organization.org.id
  trust_domain = "tzserver-tz.cofide.dev"
}

resource "cofide_connect_cluster" "cluster" {
  name               = "tzserver-cluster"
  org_id             = data.cofide_connect_organization.org.id
  trust_zone_id      = cofide_connect_trust_zone.trust_zone.id
  profile            = "kubernetes"
  kubernetes_context = "tzserver-cluster-context"
  external_server    = true

  trust_provider = {
    kind = "kubernetes"
    k8s_psat_config = {
      enabled = true
    }
  }
}

resource "cofide_connect_trust_zone_server" "server" {
  trust_zone_id = cofide_connect_trust_zone.trust_zone.id
  cluster_id    = cofide_connect_cluster.cluster.id

  helm_values = yamlencode({
    spire-server = {
      controllerManager = {
        enabled = false
      }
    }
  })

  connect_k8s_psat_config = {
    audiences                  = ["spire-server"]
    spire_server_spiffe_id_path = "/ns/spire/sa/spire-server"
  }
}

data "cofide_connect_trust_zone_server" "server" {
  id = cofide_connect_trust_zone_server.server.id
}

data "cofide_connect_trust_zone_servers" "by_trust_zone" {
  trust_zone_id = cofide_connect_trust_zone.trust_zone.id

  depends_on = [
    cofide_connect_trust_zone_server.server
  ]
}

output "server_id" {
  value = cofide_connect_trust_zone_server.server.id
}

output "server_trust_zone_id" {
  value = data.cofide_connect_trust_zone_server.server.trust_zone_id
}

output "server_cluster_id" {
  value = data.cofide_connect_trust_zone_server.server.cluster_id
}

output "server_org_id" {
  value = data.cofide_connect_trust_zone_server.server.org_id
}

output "server_k8s_psat_audiences" {
  value = cofide_connect_trust_zone_server.server.connect_k8s_psat_config.audiences
}

output "server_k8s_psat_spiffe_id_path" {
  value = cofide_connect_trust_zone_server.server.connect_k8s_psat_config.spire_server_spiffe_id_path
}

output "servers_by_trust_zone_count" {
  value = length(data.cofide_connect_trust_zone_servers.by_trust_zone.trust_zone_servers)
}

output "servers_by_trust_zone_first_id" {
  value = data.cofide_connect_trust_zone_servers.by_trust_zone.trust_zone_servers[0].id
}
