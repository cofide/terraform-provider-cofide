variable "org_id" {
  description = "The ID of the Cofide organization."
  type        = string
  default     = "example-org-id"
}

variable "name" {
  description = "The name of the cloud account."
  type        = string
  default     = "example-cloud-account"
}

variable "aws_account_id" {
  description = "The 12-digit AWS account ID."
  type        = string
  default     = "123456789012"
}

variable "discovery_role_arn" {
  description = "The ARN of the IAM role to assume for discovery."
  type        = string
  default     = "arn:aws:iam::123456789012:role/cofide-discovery"
}
