output "cluster_id" {
  value = data.cofide_connect_cluster.example.id
}

output "cluster_trust_zone_id" {
  value = data.cofide_connect_cluster.example.trust_zone_id
}

output "cluster_trust_provider_kind" {
  value = data.cofide_connect_cluster.example.trust_provider.kind
}

output "cluster_oidc_issuer_url" {
  value = data.cofide_connect_cluster.example.oidc_issuer_url
}
