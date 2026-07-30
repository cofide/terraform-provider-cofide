variable "org_id" {
  description = "The ID of the Cofide organization."
  type        = string
  default     = "example-org-id"
}

variable "name" {
  description = "The name of the cloud organization."
  type        = string
  default     = "example-cloud-org"
}

variable "aws_org_id" {
  description = "The AWS Organization ID."
  type        = string
  default     = "o-exampleorgid"
}

variable "discovery_role_arn" {
  description = "The ARN of the IAM role to assume for discovery."
  type        = string
  default     = "arn:aws:iam::123456789012:role/cofide-discovery"
}
