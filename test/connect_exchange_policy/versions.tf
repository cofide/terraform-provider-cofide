terraform {
  required_providers {
    cofide = {
      source  = "cofide/cofide"
      version = "~> v0.13.0"
    }
  }
}

provider "cofide" {
  connect_tls_grpc_target = "connect.cofide.security:8443"
}
