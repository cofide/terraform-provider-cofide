terraform {
  required_providers {
    cofide = {
      source  = "cofide/cofide"
      version = "~> 0.16.0"
    }
  }
}

provider "cofide" {}
