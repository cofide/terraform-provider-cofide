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

# Colour the banners and results only when writing to a terminal.
if [[ -t 1 ]]; then
  BOLD=$'\e[1m' CYAN=$'\e[36m' GREEN=$'\e[32m' RED=$'\e[31m' RESET=$'\e[0m'
else
  BOLD="" CYAN="" GREEN="" RED="" RESET=""
fi

# Name of the test currently running, used to label its output.
CURRENT_TEST=""

# Records, within a test directory, the var file (relative to that directory)
# of the update step most recently applied. Destroy still evaluates the
# configuration (e.g. data sources are re-read), so it must use the same
# variable values as the state it's destroying, not the defaults in main.tf.
# The marker persists across runs so that the cleanup at the start of a run
# also works after a run that failed part way through its update steps.
LAST_VAR_FILE_MARKER=".last_var_file"

# Destroys everything in a test directory's state, using the variable values it
# was last applied with.
function destroy_test() {
  local dir=${1?Specify test Terraform directory}
  local marker="$dir/$LAST_VAR_FILE_MARKER"
  local args=()
  if [[ -f "$marker" ]]; then
    args+=("-var-file=$(cat "$marker")")
  fi
  terraform -chdir="$dir" destroy -compact-warnings -auto-approve "${args[@]}" || return 1
  rm -f "$marker"
}

# Prints a banner marking the start of a phase of the current test, e.g.
# "==> connect_trust_zone | update 01-rename | plan".
function section() {
  local label="$CURRENT_TEST"
  local part
  for part in "$@"; do
    label+=" | $part"
  done
  echo
  echo "${BOLD}${CYAN}==> ${label}${RESET}"
}

function fail() {
  echo "${BOLD}${RED}FAIL: ${CURRENT_TEST}: $*${RESET}" >&2
}

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

  section "update $step" "plan"
  if ! terraform -chdir="$dir" plan -compact-warnings -var-file="$var_file" -out="$plan_file"; then
    fail "update $step: failed to plan"
    return 1
  fi
  local replaced
  replaced=$(terraform -chdir="$dir" show -json "$plan_file" |
    jq -r '.resource_changes[]? | select(.mode == "managed" and (.change.actions | index("delete"))) | .address' |
    sort) || {
    fail "update $step: failed to inspect plan"
    return 1
  }
  local expected=""
  if [[ -f "$step_dir/expect_replace" ]]; then
    expected=$(grep -v '^[[:space:]]*$' "$step_dir/expect_replace" | sort)
  fi
  if [[ "$replaced" != "$expected" ]]; then
    fail "update $step: plan replaces the wrong resources"
    echo "Expected to replace: ${expected:-<none>}" >&2
    echo "Planned to replace:  ${replaced:-<none>}" >&2
    return 1
  fi
  section "update $step" "apply"
  echo "$var_file" >"$dir/$LAST_VAR_FILE_MARKER"
  if ! terraform -chdir="$dir" apply -compact-warnings -auto-approve "$plan_file"; then
    fail "update $step: failed to apply"
    return 1
  fi
  rm -f "$dir/$plan_file"
  section "update $step" "check convergence"
  if ! terraform -chdir="$dir" plan -compact-warnings -var-file="$var_file" -detailed-exitcode; then
    fail "update $step: did not converge (plan after update is not empty)"
    return 1
  fi
  if [[ -f "$step_dir/assertions.sh" ]]; then
    section "update $step" "assertions"
    if ! bash "$step_dir/assertions.sh" "$dir"; then
      fail "update $step: assertions failed"
      return 1
    fi
  fi
}

function run_test() {
  local dir=${1?Specify test Terraform directory}
  CURRENT_TEST=$(basename "$dir")
  echo
  echo "${BOLD}${CYAN}######## ${CURRENT_TEST} ########${RESET}"

  section "init"
  if ! terraform -chdir="$dir" init -upgrade; then
    fail "failed to init"
    return 1
  fi
  section "clean up previous run"
  destroy_test "$dir" || true
  section "apply"
  rm -f "$dir/$LAST_VAR_FILE_MARKER"
  if ! terraform -chdir="$dir" apply -compact-warnings -auto-approve; then
    fail "failed to apply"
    return 1
  fi
  if [[ -f "$dir/assertions.sh" ]]; then
    section "assertions"
    if ! bash "$dir/assertions.sh" "$dir"; then
      fail "assertions failed"
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
  section "destroy"
  if ! destroy_test "$dir"; then
    fail "failed to destroy"
    return 1
  fi
  echo
  echo "${BOLD}${GREEN}PASS: ${CURRENT_TEST}${RESET}"
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
