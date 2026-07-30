output "cloud_account_id" {
  description = "The ID of the cloud account."
  value       = data.cofide_connect_cloud_account.example.id
}
