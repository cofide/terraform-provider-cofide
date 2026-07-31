output "cloud_accounts" {
  description = "The list of cloud accounts."
  value       = data.cofide_connect_cloud_accounts.example.cloud_accounts
}
