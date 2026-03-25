resource "cofide_connect_role_binding" "example_role_binding_user" {
  role_id = "example-role-id"
  user = {
    subject = "user@example.com"
  }
  resource = {
    type = "TrustZone"
    id   = "example-tz-id"
  }
}

resource "cofide_connect_role_binding" "example_role_binding_group" {
  role_id = "example-role-id"
  group = {
    claim_value = "platform-engineers"
  }
  resource = {
    type = "Organization"
    id   = "example-org-id"
  }
}
