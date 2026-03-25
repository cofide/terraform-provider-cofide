variable "name" {
  description = "The name of the cluster."
  type        = string
  default     = "example-cluster"
}

variable "trust_zone_id" {
  description = "The ID of the trust zone."
  type        = string
  default     = "example-tz-id"
}

variable "org_id" {
  description = "The ID of the organization."
  type        = string
  default     = "example-org-id"
}

variable "kubernetes_context" {
  description = "The Kubernetes context to use."
  type        = string
  default     = "example-cluster-context"
}
