resource "cofide_connect_attestation_policy" "example" {
  name   = var.name
  org_id = var.org_id

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
