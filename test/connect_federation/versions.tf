terraform {
  required_providers {
    cofide = {
      source  = "cofide/cofide"
      version = "~> 0.9.0"
    }
  }
}

provider "cofide" {
  connect_api_address = "connect.cofide.security:8443"
}
