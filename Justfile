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

# Updates the terraform provider version in all test, example, and documentation files
# Usage: just update-tf-version 0.9.0
update-tf-version version:
    @echo "Updating Terraform provider version to ~> {{version}}..."
    @find ./examples ./test -type f -name "*.tf" -exec perl -i -pe 's/(version\s*=\s*")~>\s*[0-9\.]+"/$1~> {{version}}"/' {} +
    @perl -i -pe 's/(version\s*=\s*")~>\s*[0-9\.]+"/$1~> {{version}}"/' README.md
    @echo "Update complete."
