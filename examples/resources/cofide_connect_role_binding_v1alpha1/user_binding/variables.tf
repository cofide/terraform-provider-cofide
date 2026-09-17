variable "role_id" {
  description = "The ID of the role."
  type        = string
  default     = "example-role-id"
}

variable "user_subject" {
  description = "The subject of the user."
  type        = string
  default     = "user@example.com"
}

variable "resource_type" {
  description = "The type of the resource."
  type        = string
  default     = "TrustZone"
}

variable "resource_id" {
  description = "The ID of the resource."
  type        = string
  default     = "example-tz-id"
}
