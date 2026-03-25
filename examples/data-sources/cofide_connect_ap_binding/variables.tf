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

variable "policy_id" {
  description = "The ID of the attestation policy."
  type        = string
  default     = "example-ap-id"
}
