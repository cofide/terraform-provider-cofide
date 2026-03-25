resource "cofide_connect_attestation_policy" "example" {
  name   = var.name
  org_id = var.org_id

  tpm_node = {
    attestation = {
      ek_hash = "a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2"
    }
  }
}
