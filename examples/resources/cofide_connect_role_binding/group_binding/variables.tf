variable "role_id" {
  description = "The ID of the role."
  type        = string
  default     = "example-role-id"
}

variable "group_claim_value" {
  description = "The claim value of the group."
  type        = string
  default     = "platform-engineers"
}

variable "resource_type" {
  description = "The type of the resource."
  type        = string
  default     = "Organization"
}

variable "resource_id" {
  description = "The ID of the resource."
  type        = string
  default     = "example-org-id"
}
