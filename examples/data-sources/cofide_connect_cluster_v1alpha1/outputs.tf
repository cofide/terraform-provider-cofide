output "cluster_id" {
  value = data.cofide_connect_cluster_v1alpha1.example.id
}

output "cluster_trust_zone_id" {
  value = data.cofide_connect_cluster_v1alpha1.example.trust_zone_id
}

output "cluster_trust_provider_kind" {
  value = data.cofide_connect_cluster_v1alpha1.example.trust_provider.kind
}

output "cluster_oidc_issuer_url" {
  value = data.cofide_connect_cluster_v1alpha1.example.oidc_issuer_url
}
