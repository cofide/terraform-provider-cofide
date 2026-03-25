variable "org_id" {
  description = "The ID of the organization."
  type        = string
  default     = "example-org-id"
}

variable "trust_zone_name" {
  description = "The name for the trust zone."
  type        = string
  default     = "example-tz"
}

variable "trust_domain" {
  description = "The trust domain for the trust zone."
  type        = string
  default     = "example.cofide.dev"
}
