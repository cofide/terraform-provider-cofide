terraform {
  required_providers {
    cofide = {
      source  = "cofide/cofide"
      version = "~> 0.8.0"
    }
  }
}

provider "cofide" {}
