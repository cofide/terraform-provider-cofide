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
    store_svid = true
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

resource "cofide_connect_attestation_policy" "attestation_policy_tpm_node" {
  name   = "test-ap-3"
  org_id = data.cofide_connect_organization.org.id

  tpm_node = {
    attestation = {
      ek_hash = "5b3e0a049837688b09028ba84be190720bcc8f6cf74a487dc53b2ce9f376b5fb"
    }
    selector_values = var.selector_values
  }
}

# Use a variable rather than a literal list of strings for selector values to
# cover https://github.com/cofide/terraform-provider-cofide/issues/113.
variable "selector_values" {
  type    = list(string)
  default = ["test-selector"]
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

data "cofide_connect_attestation_policy" "attestation_policy_tpm_node" {
  name   = "test-ap-3"
  org_id = data.cofide_connect_organization.org.id

  depends_on = [
    cofide_connect_attestation_policy.attestation_policy_tpm_node
  ]
}

output "attestation_policy_static_id" {
  value = data.cofide_connect_attestation_policy.attestation_policy_static.id
}

output "attestation_policy_kubernetes_id" {
  value = data.cofide_connect_attestation_policy.attestation_policy_kubernetes.id
}

# store_svid, as read back from the resource state and from the data source.
output "attestation_policy_static_store_svid_resource" {
  value = cofide_connect_attestation_policy.attestation_policy_static.static.store_svid
}

output "attestation_policy_static_store_svid_data_source" {
  value = data.cofide_connect_attestation_policy.attestation_policy_static.static.store_svid
}

# The data source previously ignored tpm_node policies altogether.
output "attestation_policy_tpm_node_ek_hash" {
  value = data.cofide_connect_attestation_policy.attestation_policy_tpm_node.tpm_node.attestation.ek_hash
}

output "attestation_policy_tpm_node_selector_values" {
  value = data.cofide_connect_attestation_policy.attestation_policy_tpm_node.tpm_node.selector_values
}

# The kubernetes policy, read back through the data source.
output "attestation_policy_kubernetes_spiffe_id_path_template" {
  value = data.cofide_connect_attestation_policy.attestation_policy_kubernetes.kubernetes.spiffe_id_path_template
}
