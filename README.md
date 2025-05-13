# Cofide Terraform Provider

## Requirements

Terraform 1.10 or newer. We recommend running the [latest version](https://developer.hashicorp.com/terraform/downloads?product_intent=terraform) for optimal compatibility with the Cofide provider.

## Usage

<!-- x-release-please-start-version -->

```hcl
# Declare the provider and version
terraform {
  required_providers {
    cofide = {
      source  = "cofide/cofide"
      version = "0.1.0"
    }
  }
}

# Initialize the provider
provider "cofide" {
  api_token = "your_api_token"
}

# Configure a resource
resource "cofide_connect_trust_zone" "example_trust_zone" {
  name         = "example-tz"
  trust_domain = "example.cofide.dev"
}
```

<!-- x-release-please-end -->

Initialize your project by running `terraform init` in the directory.

## Local Development

To use this provider locally:

1. Build and install the provider locally:
   ```bash
   just install
   ```
   This will:
   - Build the provider
   - Install it to `~/.terraform.d/plugins/local/cofide/cofide/0.1.0/<os>_<arch>`

3. In your Terraform configuration, use the local provider:
   ```hcl
   terraform {
     required_providers {
       cofide = {
         source  = "local/cofide/cofide"
         version = "0.1.0"
       }
     }
   }

   provider "cofide" {
     api_token   = "your_api_token"
     connect_url = "connect.cofide.security:8443"
   }
   ```

   Instead of using the `connect_url` attribute, you can also set the `COFIDE_CONNECT_URL` environment variable. As an alternative to setting the `api_token` attribute, you can use the `COFIDE_API_TOKEN` environment variable instead. To retrieve an API token, you must authenticate with Connect in the usual way via the `cofidectl connect login` command. The API token can be extracted from the generated `~/.cofide/credentials` file.

   If you are running a local instance of Connect, you'll need to update your `etc/hosts` file to include the following:
   ```
   <connect-api-load-balancer-service-ip> connect.cofide.security
   ```

4. Initialize your Terraform project as normal:
   ```bash
   terraform init
   ```
   
5. Now you can use the provider with `terraform plan` and `terraform apply` as normal.

To generate or update documentation for the provider, run `just generate` from the project root.
