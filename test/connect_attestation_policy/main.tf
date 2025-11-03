resource "cofide_connect_attestation_policy" "attestation_policy_static" {
  name   = "test-ap-1"
  org_id = "test-org-id"

  static = {
    spiffe_id = "spiffe://example.org/workload"
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
  }
}

resource "cofide_connect_attestation_policy" "attestation_policy_kubernetes" {
  name   = "test-ap-2"
  org_id = "test-org-id"

  kubernetes = {
    namespace_selector = {
      match_labels = {
        "kubernetes.io/metadata.name" = "test"
      }
    }
    dns_name_templates = ["example.namespace.svc.cluster.local"]
  }
}

data "cofide_connect_attestation_policy" "attestation_policy_static" {
  name   = "test-ap-1"
  org_id = "test-org-id"

  depends_on = [
    cofide_connect_attestation_policy.attestation_policy_static
  ]
}

data "cofide_connect_attestation_policy" "attestation_policy_kubernetes" {
  name   = "test-ap-2"
  org_id = "test-org-id"

  depends_on = [
    cofide_connect_attestation_policy.attestation_policy_kubernetes
  ]
}

output "attestation_policy_static_id" {
  value = data.cofide_connect_attestation_policy.attestation_policy_static.id
}

output "attestation_policy_kubernetes_id" {
  value = data.cofide_connect_attestation_policy.attestation_policy_kubernetes.id
}
