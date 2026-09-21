resource "cofide_connect_attestation_policy_v1alpha1" "example" {
  name          = var.name
  trust_zone_id = var.trust_zone_id

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
