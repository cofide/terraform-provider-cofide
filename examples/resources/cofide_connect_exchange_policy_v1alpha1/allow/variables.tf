variable "name" {
  description = "The name of the exchange policy."
  type        = string
  default     = "example-allow-policy"
}

variable "trust_zone_id" {
  description = "The ID of the trust zone."
  type        = string
  default     = "example-tz-id"
}
