variable "cloud_account_id" {
  description = "The ID of the cloud account this configuration applies to."
  type        = string
  default     = "example-cloud-account-id"
}

variable "discovery_role_arn" {
  description = "The ARN of the IAM role to assume for discovery."
  type        = string
  default     = "arn:aws:iam::123456789012:role/cofide-discovery"
}
