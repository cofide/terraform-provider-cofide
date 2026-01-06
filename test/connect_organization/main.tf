data "cofide_connect_organization" "org" {
  name = "default"
}

output "org_id" {
  value = data.cofide_connect_organization.org.id
}
