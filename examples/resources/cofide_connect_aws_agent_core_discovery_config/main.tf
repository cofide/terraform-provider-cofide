# Useful when the cloud account's existence is managed by automatic
# discovery rather than Terraform: only the AgentCore discovery config is
# managed here, referencing the account by ID.
resource "cofide_connect_aws_agent_core_discovery_config" "example" {
  cloud_account_id  = var.cloud_account_id
  audience          = "spiffe://example.org/agent-core-discovery"
  regions           = ["us-east-1"]
  discovery_enabled = true

  role_chain = [
    {
      iam_role_arn = var.discovery_role_arn
    }
  ]
}
