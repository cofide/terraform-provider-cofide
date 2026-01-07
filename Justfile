build: test
    go build -o bin/terraform-provider-cofide .

test *args:
    go run gotest.tools/gotestsum@latest --format github-actions ./... {{args}}

test-race: (test "--" "-race")

integration *args:
    {{justfile_directory()}}/test/run.sh {{args}}

lint *args:
    golangci-lint run --show-stats {{args}}

generate:
    go generate ./...
