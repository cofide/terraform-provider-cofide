output "cloud_organization_id" {
  description = "The ID of the cloud organization."
  value       = data.cofide_connect_cloud_organization.example.id
}
