output "trust_zone_id" {
  description = "The ID of the trust zone."
  value       = data.cofide_connect_trust_zone.example.id
}
