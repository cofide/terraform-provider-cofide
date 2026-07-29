variable "name" {
  description = "The name of the workload suppression rule."
  type        = string
  default     = "example-suppression-rule"
}

variable "org_id" {
  description = "The ID of the organization."
  type        = string
  default     = "example-org-id"
}

variable "trust_zone_id" {
  description = "The ID of the trust zone to match."
  type        = string
  default     = "example-tz-id"
}
