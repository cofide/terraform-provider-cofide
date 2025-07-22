resource "cofide_connect_trust_zone" "trust_zone" {
  name         = "test-tz"
  trust_domain = "test-tz.cofide.dev"
}

resource "cofide_connect_cluster" "cluster" {
  name               = "test-cluster"
  trust_zone_id      = cofide_connect_trust_zone.trust_zone.id
  org_id             = cofide_connect_trust_zone.trust_zone.org_id
  profile            = "kubernetes"
  kubernetes_context = "test-cluster-context"

  trust_provider = {
    kind = "kubernetes"
  }

  extra_helm_values = yamlencode({
    spire-server = {
      controllerManager = {
        enabled = false
      },
      extraEnv = [
        {
          name  = "CLUSTER_NAME",
          value = "test-cluster-a"
        },
      ]
    }
  })

  external_server = true

  oidc_issuer_url     = "https://oidc.example.com"
  oidc_issuer_ca_cert = base64encode("DUMMY_CERTIFICATE_DATA")

  depends_on = [
    cofide_connect_trust_zone.trust_zone
  ]
}

data "cofide_connect_cluster" "cluster" {
  name          = cofide_connect_cluster.cluster.name
  org_id        = cofide_connect_trust_zone.trust_zone.org_id

  depends_on = [
    cofide_connect_cluster.cluster
  ]
}

output "cluster_id" {
  value = data.cofide_connect_cluster.cluster.id
}

output "cluster_trust_provider_kind" {
  value = data.cofide_connect_cluster.cluster.trust_provider.kind
}

output "cluster_oidc_issuer_url" {
  value = data.cofide_connect_cluster.cluster.oidc_issuer_url
}

output "cluster_oidc_issuer_ca_cert" {
  value     = data.cofide_connect_cluster.cluster.oidc_issuer_ca_cert
  sensitive = true
}
