output "exchange_policy_ids" {
  description = "The IDs of the exchange policies."
  value       = [for p in data.cofide_connect_exchange_policies.example.exchange_policies : p.id]
}
