resource "cofide_connect_exchange_policy" "example" {
  name         = var.name
  trust_zone_id = var.trust_zone_id
  action        = "ALLOW"

  subject_identity = [
    { glob = "spiffe://example.org/ns/foo/sa/*" },
    { glob = "spiffe://example.org/ns/bar/sa/*" }
  ]

  target_audience = [
    { exact = "https://api.example.org" }
  ]

  outbound_scopes = ["read", "write"]
}
