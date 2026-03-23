#!/bin/bash

# This script can be used to run the test Terraform configurations in this
# directory against a local development deployment of Connect. It applies then
# destroys each configuration.

set -euo pipefail

TEST_DIR=$(dirname $BASH_SOURCE)
REPO_DIR=$(dirname "$TEST_DIR")

source "$TEST_DIR/test.rc"

DEV_TFRC="$REPO_DIR/dev.tfrc"
if [[ -f "$DEV_TFRC" ]]; then
  export TF_CLI_CONFIG_FILE="$DEV_TFRC"
fi

function run_test() {
  local dir=${1?Specify test Terraform directory}
  echo "Running Terraform test in \"$dir\""

  if ! terraform -chdir="$dir" init -upgrade; then
    echo "ERROR: Failed to init" >&2
    return 1
  fi
  terraform -chdir="$dir" destroy -auto-approve || true
  if ! terraform -chdir="$dir" apply -auto-approve; then
    echo "ERROR: Failed to apply" >&2
    return 1
  fi
  if [[ -f "$dir/assertions.sh" ]]; then
    if ! bash "$dir/assertions.sh" "$dir"; then
      echo "ERROR: Assertions failed" >&2
      return 1
    fi
  fi
  if ! terraform -chdir="$dir" destroy -auto-approve; then
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
