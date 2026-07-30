resource "cofide_connect_workload_suppression_rule" "example" {
  name        = var.name
  org_id      = var.org_id
  description = "Suppress noisy debug pods in the dev namespace."
  enabled     = true

  kubernetes_pod = {
    trust_zone_ids = [var.trust_zone_id]
    namespaces     = ["dev"]
    labels = {
      "app.kubernetes.io/component" = "debug"
    }
  }
}
