data "cofide_connect_workload_suppression_rules" "example" {
  org_ids = [var.org_id]
  enabled = true
}
