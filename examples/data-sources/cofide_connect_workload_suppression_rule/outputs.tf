output "workload_suppression_rule_name" {
  value = data.cofide_connect_workload_suppression_rule.example.name
}

output "workload_suppression_rule_enabled" {
  value = data.cofide_connect_workload_suppression_rule.example.enabled
}
