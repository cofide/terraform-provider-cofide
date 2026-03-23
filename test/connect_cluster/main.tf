data "cofide_connect_organization" "org" {
  name = "default"
}

resource "cofide_connect_trust_zone" "trust_zone" {
  name         = "cluster-tz"
  org_id       = data.cofide_connect_organization.org.id
  trust_domain = "cluster-tz.cofide.dev"
}

resource "cofide_connect_cluster" "cluster" {
  name               = "test-cluster"
  org_id             = data.cofide_connect_organization.org.id
  trust_zone_id      = cofide_connect_trust_zone.trust_zone.id
  profile            = "kubernetes"
  kubernetes_context = "test-cluster-context"

  trust_provider = {
    kind = "kubernetes"
    k8s_psat_config = {
      enabled = true
      allowed_service_accounts = [
        {
          namespace            = "spire"
          service_account_name = "spire-agent"
        }
      ]
      api_server_url        = "https://kubernetes.default.svc"
      api_server_ca_cert    = base64encode(file("oidc-issuer-ca.crt"))
      spire_server_audience = "spire-server"
    }
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
  oidc_issuer_ca_cert = base64encode(file("oidc-issuer-ca.crt"))

  depends_on = [
    cofide_connect_trust_zone.trust_zone
  ]
}

data "cofide_connect_cluster" "cluster" {
  name          = cofide_connect_cluster.cluster.name
  org_id        = data.cofide_connect_organization.org.id

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

output "cluster_trust_provider_k8s_psat_enabled" {
  value = data.cofide_connect_cluster.cluster.trust_provider.k8s_psat_config.enabled
}

output "cluster_trust_provider_k8s_psat_allowed_service_accounts" {
  value = data.cofide_connect_cluster.cluster.trust_provider.k8s_psat_config.allowed_service_accounts
}

output "cluster_trust_provider_k8s_psat_api_server_url" {
  value = data.cofide_connect_cluster.cluster.trust_provider.k8s_psat_config.api_server_url
}

output "cluster_trust_provider_k8s_psat_api_server_ca_cert" {
  value     = data.cofide_connect_cluster.cluster.trust_provider.k8s_psat_config.api_server_ca_cert
  sensitive = true
}

output "cluster_trust_provider_k8s_psat_spire_server_audience" {
  value = data.cofide_connect_cluster.cluster.trust_provider.k8s_psat_config.spire_server_audience
}
