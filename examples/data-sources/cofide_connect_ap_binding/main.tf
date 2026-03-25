data "cofide_connect_ap_binding" "example" {
  org_id        = var.org_id
  trust_zone_id = var.trust_zone_id
  policy_id     = var.policy_id
}
