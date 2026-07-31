output "cloud_organizations" {
  description = "The list of cloud organizations."
  value       = data.cofide_connect_cloud_organizations.example.cloud_organizations
}
