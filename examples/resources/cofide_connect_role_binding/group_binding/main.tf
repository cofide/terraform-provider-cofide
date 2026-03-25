resource "cofide_connect_role_binding" "example" {
  role_id = var.role_id
  group = {
    claim_value = var.group_claim_value
  }
  resource = {
    type = var.resource_type
    id   = var.resource_id
  }
}
