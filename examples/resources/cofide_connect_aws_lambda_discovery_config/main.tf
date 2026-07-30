# Useful when the cloud account's existence is managed by automatic
# discovery rather than Terraform: only the Lambda discovery config is
# managed here, referencing the account by ID.
resource "cofide_connect_aws_lambda_discovery_config" "example" {
  cloud_account_id   = var.cloud_account_id
  audience           = "spiffe://example.org/lambda-discovery"
  regions            = ["us-east-1", "eu-west-1"]
  discovery_enabled  = true
  discovery_interval = "5m"

  role_chain = [
    {
      iam_role_arn = var.discovery_role_arn
    }
  ]
}
