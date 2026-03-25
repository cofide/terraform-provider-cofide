resource "cofide_connect_attestation_policy" "example_attestation_policy_static" {
  name   = "example-ap-static"
  org_id = "example-org-id"

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

resource "cofide_connect_attestation_policy" "example_attestation_policy_kubernetes" {
  name   = "example-ap-kubernetes"
  org_id = "example-org-id"

  kubernetes = {
    namespace_selector = {
      match_labels = {
        "kubernetes.io/metadata.name" = "default"
      }
    }
    pod_selector = {
      match_labels = {
        "app" = "my-app"
      }
      match_expressions = [
        {
          key      = "environment"
          operator = "In"
          values   = ["production", "staging"]
        }
      ]
    }
    spiffe_id_path_template = "ns/default/sa/my-service-account"
  }
}

resource "cofide_connect_attestation_policy" "example_attestation_policy_tpm_node" {
  name   = "example-ap-tpm"
  org_id = "example-org-id"

  tpm_node = {
    attestation = {
      ek_hash = "a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2"
    }
    selector_values = ["plugin_name:tpm"]
  }
}
