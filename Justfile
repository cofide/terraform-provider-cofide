install: build
    mkdir -p ~/.terraform.d/plugins/local/cofide/cofide/0.1.0/$(go env GOOS)_$(go env GOARCH)
    cp terraform-provider-cofide ~/.terraform.d/plugins/local/cofide/cofide/0.1.0/$(go env GOOS)_$(go env GOARCH)

build: test
    go build -o terraform-provider-cofide ./

test *args:
    go run gotest.tools/gotestsum@latest --format github-actions ./... {{args}}

test-race: (test "--" "-race")

lint *args:
    golangci-lint run --show-stats {{args}}

generate:
    go generate ./...
