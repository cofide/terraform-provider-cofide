# This cloud account only manages identity/name and does not set discovery
# config or suppression inline, leaving them under external management (e.g.
# automatic discovery, or the standalone cofide_connect_aws_lambda_discovery_config,
# cofide_connect_aws_agent_core_discovery_config, and
# cofide_connect_cloud_account_suppression_config resources).
resource "cofide_connect_cloud_account" "example" {
  org_id                = var.org_id
  cloud_organization_id = var.cloud_organization_id
  name                  = var.name

  aws = {
    account_id = var.aws_account_id
  }
}
