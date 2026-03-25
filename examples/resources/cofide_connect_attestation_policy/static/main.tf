resource "cofide_connect_attestation_policy" "example" {
  name   = var.name
  org_id = var.org_id

  static = {
    spiffe_id_path = "ns/default/sa/my-service-account"
    parent_id_path = "test/agent"
    selectors = [
      {
        type  = "k8s"
        value = "ns:default"
      },
      {
        type  = "k8s"
        value = "sa:my-service-account"
      }
    ]
    dns_names = [
      "test.workload"
    ]
  }
}
