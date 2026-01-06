data "cofide_connect_organization" "org" {
  name = "default"
}

resource "cofide_connect_attestation_policy" "attestation_policy_static" {
  name   = "test-ap-1"
  org_id = data.cofide_connect_organization.org.id

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

resource "cofide_connect_attestation_policy" "attestation_policy_kubernetes" {
  name   = "test-ap-2"
  org_id = data.cofide_connect_organization.org.id

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
    # TODO: dns_name_templates is not yet supported in Connect.
    #dns_name_templates = [
    #  "test.workload"
    #]
    spiffe_id_path_template = "test/workload"
  }
}

data "cofide_connect_attestation_policy" "attestation_policy_static" {
  name   = "test-ap-1"
  org_id = data.cofide_connect_organization.org.id

  depends_on = [
    cofide_connect_attestation_policy.attestation_policy_static
  ]
}

data "cofide_connect_attestation_policy" "attestation_policy_kubernetes" {
  name   = "test-ap-2"
  org_id = data.cofide_connect_organization.org.id

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
