# Overridden by update.tfvars to exercise an in-place update (see run.sh).
variable "trust_zone_name" {
  type    = string
  default = "test-tz"
}

data "cofide_connect_organization_v1alpha1" "org" {
  name = "default"
}

resource "cofide_connect_trust_zone_v1alpha1" "trust_zone" {
  name         = var.trust_zone_name
  org_id       = data.cofide_connect_organization_v1alpha1.org.id
  trust_domain = "test-tz.cofide.dev"
}

data "cofide_connect_trust_zone_v1alpha1" "trust_zone" {
  name         = cofide_connect_trust_zone_v1alpha1.trust_zone.name
  org_id       = data.cofide_connect_organization_v1alpha1.org.id
  trust_domain = "test-tz.cofide.dev"
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

output "trust_zone_bundle_endpoint_url" {
  value = cofide_connect_trust_zone_v1alpha1.trust_zone.bundle_endpoint_url
}

output "trust_zone_bundle_endpoint_profile" {
  value = cofide_connect_trust_zone_v1alpha1.trust_zone.bundle_endpoint_profile
}

output "trust_zone_jwt_issuer" {
  value = cofide_connect_trust_zone_v1alpha1.trust_zone.jwt_issuer
}
