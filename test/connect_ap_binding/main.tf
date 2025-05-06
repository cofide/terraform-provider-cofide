resource "cofide_connect_trust_zone" "trust_zone" {
  name         = "test-tz"
  trust_domain = "test-tz.cofide.dev"
}

resource "cofide_connect_attestation_policy" "attestation_policy_static" {
  name   = "test-ap"
  org_id = cofide_connect_trust_zone.trust_zone.org_id

  static = {
    spiffe_id = "spiffe://example.org/workload"
    selectors = [
      {
        type  = "k8s"
        value = "ns:demo"
      },
      {
        type  = "k8s"
        value = "sa:demo-sa"
      }
    ]
  }

  depends_on = [
    cofide_connect_trust_zone.trust_zone
  ]
}

resource "cofide_connect_ap_binding" "ap_binding" {
  org_id        = cofide_connect_trust_zone.trust_zone.org_id
  trust_zone_id = cofide_connect_trust_zone.trust_zone.id
  policy_id     = cofide_connect_attestation_policy.attestation_policy_static.id
  federations = [
    {
      trust_zone_id = cofide_connect_trust_zone.trust_zone.id
    }
  ]

  depends_on = [
    cofide_connect_trust_zone.trust_zone,
    cofide_connect_attestation_policy.attestation_policy_static
  ]
}

data "cofide_connect_ap_binding" "ap_binding" {
  org_id        = cofide_connect_trust_zone.trust_zone.org_id
  trust_zone_id = cofide_connect_trust_zone.trust_zone.id
  policy_id     = cofide_connect_attestation_policy.attestation_policy_static.id

  depends_on = [
    cofide_connect_ap_binding.ap_binding
  ]
}

output "trust_zone_id" {
  value = cofide_connect_trust_zone.trust_zone.id
}

output "attestation_policy_id" {
  value = cofide_connect_attestation_policy.attestation_policy_static.id
}

output "ap_binding_id" {
  value = data.cofide_connect_ap_binding.ap_binding.id
}
