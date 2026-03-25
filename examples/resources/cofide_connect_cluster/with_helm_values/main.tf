resource "cofide_connect_cluster" "example" {
  name                = var.name
  trust_zone_id       = var.trust_zone_id
  org_id              = var.org_id
  profile             = "kubernetes"
  kubernetes_context  = var.kubernetes_context
  external_server     = false
  oidc_issuer_url     = var.oidc_issuer_url
  oidc_issuer_ca_cert = base64encode(file("${path.module}/ca.pem"))

  trust_provider = {
    kind = "kubernetes"
  }

  extra_helm_values = yamlencode({
    replicaCount = 2
    resources = {
      limits = {
        cpu    = "200m"
        memory = "256Mi"
      }
    }
  })
}
