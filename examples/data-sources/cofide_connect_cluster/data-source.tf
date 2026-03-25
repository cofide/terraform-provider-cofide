data "cofide_connect_cluster" "example_cluster" {
  name   = "example-cluster"
  org_id = "example-org-id"
}

output "cluster_id" {
  value = data.cofide_connect_cluster.example_cluster.id
}

output "cluster_trust_zone_id" {
  value = data.cofide_connect_cluster.example_cluster.trust_zone_id
}

output "cluster_trust_provider_kind" {
  value = data.cofide_connect_cluster.example_cluster.trust_provider.kind
}

output "cluster_oidc_issuer_url" {
  value = data.cofide_connect_cluster.example_cluster.oidc_issuer_url
}
