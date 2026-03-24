resource "cofide_connect_cluster" "example_cluster" {
  name               = "example-cluster"
  trust_zone_id      = "example-tz-id"
  org_id             = "example-org-id"
  profile            = "kubernetes"
  kubernetes_context = "example-cluster-context"

  trust_provider = {
    kind = "kubernetes"
  }
}
