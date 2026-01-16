data "cofide_connect_organization" "org" {
  name = "default"
}

resource "cofide_connect_trust_zone" "trust_zone" {
  name         = "test-role-binding-tz"
  org_id       = data.cofide_connect_organization.org.id
  trust_domain = "test-rb-tz.cofide.dev"
}

resource "cofide_connect_role_binding" "role_binding" {
  role_id = "admin"
  user = {
    subject = "test-user-subject"
  }
  resource = {
    type = "TrustZone"
    id   = cofide_connect_trust_zone.trust_zone.id
  }
}

data "cofide_connect_role_binding" "role_binding" {
  role_id = "admin"
  user = {
    subject = "test-user-subject"
  }
  resource = {
    type = "TrustZone"
    id   = cofide_connect_trust_zone.trust_zone.id
  }

  depends_on = [
    cofide_connect_role_binding.role_binding
  ]
}

output "role_binding_id" {
  value = data.cofide_connect_role_binding.role_binding.id
}
