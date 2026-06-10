# AGENTS.md

This file provides guidance to AI coding agents when working with code in this repository.

## Commands

```bash
just build          # Build the provider binary to bin/terraform-provider-cofide
just test           # Run unit tests (uses gotestsum)
just test-race      # Run unit tests with race detector
just lint           # Run golangci-lint
just generate       # Regenerate docs from templates
just integration    # Run all integration tests against a local Connect deployment
just integration connect_trust_zone  # Run a specific integration test
just update-tf-version 0.9.0  # Update provider version in all tf/docs files
```

Run a single unit test:
```bash
go test ./internal/services/cluster/... -run TestFunctionName
```

For local development, build the provider and set `TF_CLI_CONFIG_FILE=./dev.tfrc` to use the local binary instead of the registry version. The `dev.tfrc` file points Terraform to `bin/` for the `cofide/cofide` provider.

## Architecture

This is a Terraform provider built with `terraform-plugin-framework` that manages Cofide Connect workload identity resources via gRPC.

**Provider entry point**: `main.go` → `internal/provider.go` (`CofideProvider`)

**Transport**: `internal/client/client.go` — creates a gRPC TLS connection to the Cofide Connect API. The SDK client (`github.com/cofide/cofide-api-sdk/pkg/connect/client.ClientSet`) is initialized once in `provider.Configure()` and passed to all resources/data sources via `ProviderData`.

**Resource structure**: Each resource type lives in `internal/services/<name>/` with a consistent file layout:
- `resource.go` — implements `resource.Resource` CRUD methods
- `data_source.go` — implements `datasource.DataSource` Read method
- `schema.go` — defines the Terraform schema for the resource
- `data_source_schema.go` — defines the schema for the data source
- `model.go` — defines the Go struct (`*Model`) that maps to Terraform state via `tfsdk:` tags
- `convert.go` — bidirectional conversion between the model and the protobuf API types (not all resources have this; simpler ones do conversion inline)

**Protobuf types**: Come from `github.com/cofide/cofide-api-sdk/gen/go/proto/...`. Each resource calls the corresponding versioned SDK client (e.g. `t.client.TrustZoneV1Alpha1()`).

**Shared utilities**:
- `internal/planmodifiers/optional_computed_modifier.go` — `OptionalComputedModifier` for fields that are optional in config but computed by the API; marks value unknown during plan when config omits it
- `internal/util/util.go` — `HelmValuesForState` preserves user-supplied YAML/JSON helm values format across Read/Update to avoid spurious diffs; `IsStringAttributeNonEmpty` for nil-safe optional string handling

**Resources registered** (in `provider.go`):
- `cofide_connect_trust_zone` — top-level SPIFFE trust domain container
- `cofide_connect_cluster` — Kubernetes cluster registered in a trust zone
- `cofide_connect_attestation_policy` — workload attestation policy (Kubernetes, static, or TPM)
- `cofide_connect_ap_binding` — binds an attestation policy to a trust zone
- `cofide_connect_federation` — federates two trust zones
- `cofide_connect_exchange_policy` — controls workload-to-workload communication
- `cofide_connect_role_binding` — assigns roles within the platform
- `cofide_connect_trust_zone_server` — external SPIRE server attached to a trust zone

**Data sources** mirror the resources above, with some supporting list variants (e.g. `exchangepolicy.NewListDataSource`, `trustzoneserver.NewListDataSource`).

## Key conventions

- Resources that support `ImportState` use `resource.ImportStatePassthroughID` to import by ID.
- gRPC `codes.NotFound` on Read removes the resource from state (drift detection), rather than erroring.
- Provider config can be supplied via HCL attributes or environment variables (`COFIDE_API_TOKEN`, `COFIDE_CONNECT_URL`, `COFIDE_INSECURE_SKIP_VERIFY`).
- The gRPC client uses a retry policy on `UNAUTHENTICATED` errors (up to 10 attempts) to handle transient JWKS fetch failures.

## Integration tests

Tests live in `test/<resource>/main.tf` and are applied/destroyed by `test/run.sh`. They require a running local Connect deployment and an active `cofidectl connect login` session. The `test/test.rc` file sources environment variables for the test run.

## Releasing

Uses `goreleaser`. Before tagging a release, run `just update-tf-version <version>` to update version constraints in all `.tf` example/test files and `README.md`. After changes to resource schemas or descriptions, run `just generate` to regenerate the docs in `docs/`.
