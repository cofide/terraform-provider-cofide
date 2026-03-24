resource "cofide_connect_attestation_policy" "example_attestation_policy_static" {
  name   = "example-ap-static"
  org_id = "example-org-id"

  static = {
    spiffe_id_path = "test/workload"
    parent_id_path = "test/agent"
    selectors = [
      {
        type  = "k8s"
        value = "ns:test"
      },
      {
        type  = "k8s"
        value = "sa:test-sa"
      }
    ]
    dns_names = [
      "test.workload"
    ]
  }
}

resource "cofide_connect_attestation_policy" "example_attestation_policy_kubernetes" {
  name   = "example-ap-kubernetes"
  org_id = "example-org-id"

  kubernetes = {
    namespace_selector = {
      match_labels = {
        "kubernetes.io/metadata.name" = "test"
      }
    }
    pod_selector = {
      match_labels = {
        "test-label" = "test"
      }
    }
    spiffe_id_path_template = "test/workload"
  }
}
