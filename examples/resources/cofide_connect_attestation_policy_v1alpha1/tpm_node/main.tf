resource "cofide_connect_attestation_policy_v1alpha1" "example" {
  name          = var.name
  trust_zone_id = var.trust_zone_id

  tpm_node = {
    attestation = {
      ek_hash = "a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2"
    }
  }
}
