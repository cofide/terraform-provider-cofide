output "trust_zone_server_ids" {
  value = [for s in data.cofide_connect_trust_zone_servers.example.trust_zone_servers : s.id]
}
