data "cofide_connect_organization_v1alpha1" "org" {
  name = "default"
}

resource "cofide_connect_trust_zone_v1alpha1" "trust_zone" {
  name         = "test-role-binding-tz"
  org_id       = data.cofide_connect_organization_v1alpha1.org.id
  trust_domain = "test-rb-tz.cofide.dev"
}

resource "cofide_connect_role_binding_v1alpha1" "role_binding" {
  role_id = "admin"
  user = {
    subject = "test-user-subject"
  }
  resource = {
    type = "TrustZone"
    id   = cofide_connect_trust_zone_v1alpha1.trust_zone.id
  }
}

output "role_binding_id" {
  value = cofide_connect_role_binding_v1alpha1.role_binding.id
}
