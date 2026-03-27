output "trust_zone_server_id" {
  value = data.cofide_connect_trust_zone_server.example.id
}

output "trust_zone_server_trust_zone_id" {
  value = data.cofide_connect_trust_zone_server.example.trust_zone_id
}

output "trust_zone_server_cluster_id" {
  value = data.cofide_connect_trust_zone_server.example.cluster_id
}
