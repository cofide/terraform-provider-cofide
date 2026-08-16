resource "cofide_connect_attestation_policy" "example" {
  name   = var.name
  org_id = var.org_id

  # Binds this policy directly to a trust zone (and optionally federates it
  # with other trust zones) without needing a separate
  # cofide_connect_ap_binding resource.
  trust_zone_id = var.trust_zone_id
  federations = [
    {
      trust_zone_id = var.remote_trust_zone_id
    }
  ]

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
    }
    spiffe_id_path_template = "ns/default/sa/my-service-account"
  }
}
