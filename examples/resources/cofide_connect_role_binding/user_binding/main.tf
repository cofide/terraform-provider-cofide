resource "cofide_connect_role_binding" "example" {
  role_id = var.role_id
  user = {
    subject = var.user_subject
  }
  resource = {
    type = var.resource_type
    id   = var.resource_id
  }
}
