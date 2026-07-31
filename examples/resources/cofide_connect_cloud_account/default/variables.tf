variable "org_id" {
  description = "The ID of the Cofide organization."
  type        = string
  default     = "example-org-id"
}

variable "cloud_organization_id" {
  description = "The ID of the cloud organization this account belongs to."
  type        = string
  default     = "example-cloud-org-id"
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
