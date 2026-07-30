resource "cofide_connect_cloud_organization" "example" {
  org_id = var.org_id
  name   = var.name

  aws = {
    aws_org_id = var.aws_org_id
    audience   = "spiffe://example.org/cloud-discovery"

    role_chain = [
      {
        iam_role_arn = var.discovery_role_arn
      }
    ]
  }

  discovery_enabled  = true
  discovery_interval = "5m"
}
