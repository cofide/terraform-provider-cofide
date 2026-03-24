resource "cofide_connect_role_binding" "example_role_binding" {
  role_id = "example-role-id"
  user = {
    subject = "example-user-subject"
  }
  resource = {
    type = "example-resource-type"
    id   = "example-resource-id"
  }
}
