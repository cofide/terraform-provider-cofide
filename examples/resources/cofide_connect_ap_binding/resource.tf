resource "cofide_connect_ap_binding" "example_ap_binding" {
  org_id        = "example-org-id"
  trust_zone_id = "example-tz-id"
  policy_id     = "example-ap-id"
}

resource "cofide_connect_ap_binding" "example_ap_binding_with_federations" {
  org_id        = "example-org-id"
  trust_zone_id = "example-tz-id"
  policy_id     = "example-ap-id"

  federations = [
    {
      trust_zone_id = "example-remote-tz-id"
    }
  ]
}
