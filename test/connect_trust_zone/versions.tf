terraform {
  required_providers {
    cofide = {
      source  = "cofide/cofide"
      version = "~> 0.16.0"
    }
  }
}

provider "cofide" {
  connect_tls_grpc_target = "connect.cofide.security:8443"
}
