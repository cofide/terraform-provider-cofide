variable "cloud_organization_ids" {
  description = "Filter by cloud organization IDs."
  type        = list(string)
  default     = ["example-cloud-org-id"]
}
