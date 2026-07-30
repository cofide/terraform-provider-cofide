# This cloud account manages its Lambda discovery config and suppression
# state inline. Do not also manage these fields with the standalone
# cofide_connect_aws_lambda_discovery_config or
# cofide_connect_cloud_account_suppression_config resources for the same account.
resource "cofide_connect_cloud_account" "example" {
  org_id = var.org_id
  name   = var.name

  aws = {
    account_id = var.aws_account_id

    lambda_discovery_config = {
      audience          = "spiffe://example.org/lambda-discovery"
      regions           = ["us-east-1", "eu-west-1"]
      discovery_enabled = true

      role_chain = [
        {
          iam_role_arn = var.discovery_role_arn
        }
      ]
    }
  }

  suppressed = false
}
