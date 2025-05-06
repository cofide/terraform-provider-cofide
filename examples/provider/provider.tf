terraform {
  required_providers {
    cofide = {
      source  = "local/cofide/cofide"
      version = "0.1.0"
    }
  }
}

provider "cofide" {}
