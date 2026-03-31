terraform {
  required_providers {
    cofide = {
      source  = "cofide/cofide"
      version = "~> 0.9.0"
    }
  }
}

provider "cofide" {
  connect_url = "cofide.security:8443"
}
