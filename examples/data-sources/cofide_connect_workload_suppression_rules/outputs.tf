output "workload_suppression_rule_ids" {
  value = [for r in data.cofide_connect_workload_suppression_rules.example.workload_suppression_rules : r.id]
}
