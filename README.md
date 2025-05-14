# Cofide Terraform Provider

This is the repository for the Cofide Terraform Provider. Learn more about the Cofide Connect workload identity platform at https://www.cofide.io/.

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
  api_token   = "your_api_token"
  connect_url = "foo.cofide.dev:8443"
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

1. **Build the provider locally:**

   ```bash
   just build
   ```
   This will:
   - Build the provider and install it in `bin/terraform-provider-cofide`

2. **Create a `dev.tfrc` file to enable local provider development**  
   (see [Terraform docs](https://developer.hashicorp.com/terraform/cli/config/config-file#development-overrides-for-provider-developers) for more info):

   ```
   provider_installation {
     dev_overrides {
       "cofide/cofide" = "/<PATH-TO-REPO>/terraform-provider-cofide/bin"
     }

     direct {}
   }
   ```

3. **Configure your Terraform project:**

   ```hcl
   terraform {
     required_providers {
       cofide = {
         source  = "cofide/cofide"
         version = "0.1.0"
       }
     }
   }

   provider "cofide" {
     api_token            = "your_api_token"
     connect_url          = "foo.cofide.dev:8443"
     insecure_skip_verify = true                     # Only use this for local development
   }
   ```

   - The `connect_url` will be provided to you by Cofide, or use your local instance URL for development.
   - Instead of using the `connect_url` and `api_token` attributes, you can set the `COFIDE_CONNECT_URL` and `COFIDE_API_TOKEN` environment variables.
   - To retrieve an API token, authenticate with Connect using `cofidectl connect login`. The token can be found in `~/.cofide/credentials`.

   - If you are running a local instance of Connect, update your `/etc/hosts` file to include:
     ```
     <connect-api-load-balancer-service-ip> connect.cofide.security
     ```

4. **Initialize your Terraform project:**

   ```bash
   TF_CLI_CONFIG_FILE=./dev.tfrc terraform init
   ```

5. **Use the provider as normal:**

   ```bash
   TF_CLI_CONFIG_FILE=./dev.tfrc terraform plan
   TF_CLI_CONFIG_FILE=./dev.tfrc terraform apply
   ```

To generate or update documentation for the provider, run the following command from the project root:

```bash
just generate
```
