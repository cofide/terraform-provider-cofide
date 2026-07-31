# Useful when the cloud account's existence is managed by automatic
# discovery rather than Terraform: only the suppressed flag is managed here,
# referencing the account by ID.
resource "cofide_connect_cloud_account_suppression_config" "example" {
  cloud_account_id = var.cloud_account_id
  suppressed       = true
}
