data "cofide_connect_organization_v1alpha1" "org" {
  name = "default"
}

resource "cofide_connect_trust_zone_v1alpha1" "trust_zone_a" {
  name         = "test-tz-a"
  org_id       = data.cofide_connect_organization_v1alpha1.org.id
  trust_domain = "test-tz-a.cofide.dev"
}

resource "cofide_connect_trust_zone_v1alpha1" "trust_zone_b" {
  name         = "test-tz-b"
  org_id       = data.cofide_connect_organization_v1alpha1.org.id
  trust_domain = "test-tz-b.cofide.dev"
}

resource "cofide_connect_federation_v1alpha1" "federation" {
  trust_zone_id        = cofide_connect_trust_zone_v1alpha1.trust_zone_a.id
  remote_trust_zone_id = cofide_connect_trust_zone_v1alpha1.trust_zone_b.id
}
