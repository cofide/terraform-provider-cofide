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

# Exercises updates to an applied configuration. A test opts in by providing an
# updates/ directory containing one subdirectory per update step. Steps run in
# lexical order (so prefix them, e.g. 01-rename), each starting from the state
# the previous step left behind. A step directory contains:
#
#   update.tfvars   Variable values to apply. Variables it doesn't set revert
#                   to their defaults in main.tf, so a step must also carry any
#                   overrides from earlier steps that it wants to keep.
#   expect_replace  Optional. Addresses of the managed resources the update is
#                   expected to replace, one per line. Without it, every change
#                   must be an in-place update.
#   assertions.sh   Optional. Run against the test directory after the update.
#
# Each step must plan exactly the expected replacements, apply cleanly, and
# converge (no diff on a subsequent plan).
function run_update() {
  local dir=${1?Specify test Terraform directory}
  local step_dir=${2?Specify update step directory}
  local step
  step=$(basename "$step_dir")
  local var_file="updates/$step/update.tfvars"
  local plan_file="update.tfplan"
  echo "Running Terraform update step \"$step\" in \"$dir\""

  if ! terraform -chdir="$dir" plan -var-file="$var_file" -out="$plan_file"; then
    echo "ERROR: Failed to plan update" >&2
    return 1
  fi
  local replaced
  replaced=$(terraform -chdir="$dir" show -json "$plan_file" |
    jq -r '.resource_changes[]? | select(.mode == "managed" and (.change.actions | index("delete"))) | .address' |
    sort) || {
    echo "ERROR: Failed to inspect update plan" >&2
    return 1
  }
  local expected=""
  if [[ -f "$step_dir/expect_replace" ]]; then
    expected=$(grep -v '^[[:space:]]*$' "$step_dir/expect_replace" | sort)
  fi
  if [[ "$replaced" != "$expected" ]]; then
    echo "ERROR: Update step \"$step\" replaces the wrong resources" >&2
    echo "Expected to replace: ${expected:-<none>}" >&2
    echo "Planned to replace:  ${replaced:-<none>}" >&2
    return 1
  fi
  if ! terraform -chdir="$dir" apply -auto-approve "$plan_file"; then
    echo "ERROR: Failed to apply update" >&2
    return 1
  fi
  rm -f "$dir/$plan_file"
  if ! terraform -chdir="$dir" plan -var-file="$var_file" -detailed-exitcode; then
    echo "ERROR: Update did not converge (plan after update is not empty)" >&2
    return 1
  fi
  if [[ -f "$step_dir/assertions.sh" ]]; then
    if ! bash "$step_dir/assertions.sh" "$dir"; then
      echo "ERROR: Update assertions failed" >&2
      return 1
    fi
  fi
}

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
  if [[ -d "$dir/updates" ]]; then
    local step_dir
    for step_dir in "$dir"/updates/*/; do
      if ! run_update "$dir" "${step_dir%/}"; then
        return 1
      fi
    done
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
