resource "cofide_connect_cluster" "example_cluster" {
  name               = "example-cluster"
  trust_zone_id      = "example-tz-id"
  org_id             = "example-org-id"
  profile            = "kubernetes"
  kubernetes_context = "example-cluster-context"
  external_server    = false

  trust_provider = {
    kind = "kubernetes"
  }
}

resource "cofide_connect_cluster" "example_cluster_with_helm_values" {
  name                = "example-cluster-helm"
  trust_zone_id       = "example-tz-id"
  org_id              = "example-org-id"
  profile             = "kubernetes"
  kubernetes_context  = "example-cluster-context"
  external_server     = false
  oidc_issuer_url     = "https://oidc.example.com"
  oidc_issuer_ca_cert = base64encode(file("ca.pem"))

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
