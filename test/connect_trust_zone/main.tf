# Overridden by the steps in updates/ to exercise updates (see run.sh).
variable "trust_zone_name" {
  type    = string
  default = "test-tz"
}

variable "trust_domain" {
  type    = string
  default = "test-tz.cofide.dev"
}

data "cofide_connect_organization_v1alpha1" "org" {
  name = "default"
}

resource "cofide_connect_trust_zone_v1alpha1" "trust_zone" {
  name         = var.trust_zone_name
  org_id       = data.cofide_connect_organization_v1alpha1.org.id
  trust_domain = var.trust_domain
}

data "cofide_connect_trust_zone_v1alpha1" "trust_zone" {
  name         = cofide_connect_trust_zone_v1alpha1.trust_zone.name
  org_id       = data.cofide_connect_organization_v1alpha1.org.id
  trust_domain = cofide_connect_trust_zone_v1alpha1.trust_zone.trust_domain
}

output "trust_zone_id" {
  value = data.cofide_connect_trust_zone_v1alpha1.trust_zone.id
}

output "trust_zone_name" {
  value = cofide_connect_trust_zone_v1alpha1.trust_zone.name
}

output "trust_zone_name_data_source" {
  value = data.cofide_connect_trust_zone_v1alpha1.trust_zone.name
}

output "trust_zone_trust_domain" {
  value = cofide_connect_trust_zone_v1alpha1.trust_zone.trust_domain
}

output "trust_zone_trust_domain_data_source" {
  value = data.cofide_connect_trust_zone_v1alpha1.trust_zone.trust_domain
}

output "trust_zone_bundle_endpoint_url" {
  value = cofide_connect_trust_zone_v1alpha1.trust_zone.bundle_endpoint_url
}

output "trust_zone_bundle_endpoint_profile" {
  value = cofide_connect_trust_zone_v1alpha1.trust_zone.bundle_endpoint_profile
}

output "trust_zone_jwt_issuer" {
  value = cofide_connect_trust_zone_v1alpha1.trust_zone.jwt_issuer
}
