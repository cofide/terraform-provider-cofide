terraform {
  required_providers {
    cofide = {
      source  = "local/cofide/cofide"
      version = "0.1.0"
    }
  }
}

provider "cofide" {
  connect_url = "cofide.security:8443"
}
