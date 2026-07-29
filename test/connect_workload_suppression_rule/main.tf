data "cofide_connect_organization" "org" {
  name = "default"
}

resource "cofide_connect_workload_suppression_rule" "test" {
  name        = "test-workload-suppression-rule"
  org_id      = data.cofide_connect_organization.org.id
  description = "test rule created by integration tests"
  enabled     = true

  kubernetes_pod = {
    namespaces = ["test"]
    labels = {
      "test-label" = "test"
    }
  }
}

data "cofide_connect_workload_suppression_rule" "test" {
  id = cofide_connect_workload_suppression_rule.test.id

  depends_on = [
    cofide_connect_workload_suppression_rule.test
  ]
}

data "cofide_connect_workload_suppression_rules" "test" {
  org_ids = [data.cofide_connect_organization.org.id]

  depends_on = [
    cofide_connect_workload_suppression_rule.test
  ]
}

output "workload_suppression_rule_id" {
  value = data.cofide_connect_workload_suppression_rule.test.id
}

output "workload_suppression_rule_ids" {
  value = [for r in data.cofide_connect_workload_suppression_rules.test.workload_suppression_rules : r.id]
}
