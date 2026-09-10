data "cofide_connect_organization_v1alpha1" "org" {
  name = "default"
}

output "org_id" {
  value = data.cofide_connect_organization_v1alpha1.org.id
}
