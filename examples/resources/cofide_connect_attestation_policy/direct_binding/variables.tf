variable "name" {
  description = "The name of the attestation policy."
  type        = string
  default     = "example-ap-direct-binding"
}

variable "org_id" {
  description = "The ID of the organization."
  type        = string
  default     = "example-org-id"
}

variable "trust_zone_id" {
  description = "The ID of the trust zone."
  type        = string
  default     = "example-tz-id"
}

variable "remote_trust_zone_id" {
  description = "The ID of the remote (federated) trust zone."
  type        = string
  default     = "example-remote-tz-id"
}
