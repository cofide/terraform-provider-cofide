output "trust_zone_server_id" {
  value = data.cofide_connect_trust_zone_server_v1alpha1.example.id
}

output "trust_zone_server_trust_zone_id" {
  value = data.cofide_connect_trust_zone_server_v1alpha1.example.trust_zone_id
}

output "trust_zone_server_cluster_id" {
  value = data.cofide_connect_trust_zone_server_v1alpha1.example.cluster_id
}
