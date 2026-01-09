# Terraform Provider for Cofide

This project is the official Terraform provider for the Cofide Connect workload identity platform. It is written in Go and uses the `terraform-plugin-framework`.

## Project Overview

*   **Purpose**: To manage Cofide Connect resources using Terraform.
*   **Technologies**: Go, Terraform
*   **Architecture**: The provider is a Go binary that communicates with the Cofide Connect API. It defines a set of resources and data sources that map to the API's objects.

## Building and Running

This project uses `just` as a command runner.

*   **Build the provider**:

    ```bash
    just build
    ```

*   **Run tests**:

    ```bash
    just test
    ```

*   **Run integration tests**:

    Requires a local development deployment of Connect and active login.

    ```bash
    just integration
    ```

*   **Lint the code**:

    ```bash
    just lint
    ```

*   **Generate documentation**:

    ```bash
    just generate
    ```

## Development Conventions

*   **Code Style**: Standard Go formatting is enforced using `go fmt`.
*   **Static Analysis**: `go vet` and `golangci-lint` are used to find potential issues.
*   **Releasing**: The project uses `goreleaser` for building and releasing binaries.
*   **Local Development**: For local development, a `dev.tfrc` file is used to override the provider installation. See the `README.md` for more details.
