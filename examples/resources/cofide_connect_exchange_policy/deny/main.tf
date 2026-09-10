resource "cofide_connect_exchange_policy_v1alpha1" "example" {
  name          = var.name
  trust_zone_id = var.trust_zone_id
  action        = "DENY"

  subject_identity = [
    { exact = "spiffe://example.org/untrusted-workload" }
  ]
}
