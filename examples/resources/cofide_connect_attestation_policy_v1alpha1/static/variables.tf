variable "name" {
  description = "The name of the attestation policy."
  type        = string
  default     = "example-ap-static"
}

variable "trust_zone_id" {
  description = "The ID of the trust zone."
  type        = string
  default     = "example-tz-id"
}
