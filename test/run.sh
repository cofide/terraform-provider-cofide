#!/bin/bash

# This script can be used to run the test Terraform configurations in this
# directory against a local development deployment of Connect. It applies then
# destroys each configuration.

set -euo pipefail

TEST_DIR=$(dirname $BASH_SOURCE)

source "$TEST_DIR/test.rc"

function run_test() {
  local dir=${1?Specify test Terraform directory}
  echo "Running Terraform test in \"$dir\""
  pushd "$dir" >/dev/null
  # Ensure popd is called on function return to restore the original directory.
  trap "popd >/dev/null" RETURN

  if ! terraform apply -auto-approve; then
    echo "ERROR: Failed to apply" >&2
    return 1
  fi
  if ! terraform destroy -auto-approve; then
    echo "ERROR: Failed to destroy" >&2
    return 1
  fi
  echo "SUCCESS"
}

if [[ $# -ne 0 ]]; then
  for dir in "$@"; do
    run_test "$TEST_DIR/$dir"
  done
else
  while IFS= read -r -d '' dir; do
    run_test "$dir"
  done < <(find "$TEST_DIR/" -maxdepth 1 -type d -name 'connect_*' -print0)
fi
