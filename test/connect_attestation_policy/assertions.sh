#!/bin/bash

set -euo pipefail

dir=${1?Specify test Terraform directory}

outputs=$(terraform -chdir="$dir" output -json)

assert_eq() {
  local key=$1
  local expected=$2
  local actual
  actual=$(echo "$outputs" | jq -r ".$key.value")
  if [[ "$actual" != "$expected" ]]; then
    echo "ERROR: output '$key': expected '$expected', got '$actual'" >&2
    return 1
  fi
}

assert_jq_eq() {
  local filter=$1
  local expected=$2
  local actual
  actual=$(echo "$outputs" | jq -r "$filter")
  if [[ "$actual" != "$expected" ]]; then
    echo "ERROR: '$filter': expected '$expected', got '$actual'" >&2
    return 1
  fi
}

assert_not_empty() {
  local key=$1
  local actual
  actual=$(echo "$outputs" | jq -r ".$key.value")
  if [[ -z "$actual" || "$actual" == "null" ]]; then
    echo "ERROR: output '$key': expected non-empty value" >&2
    return 1
  fi
}

assert_not_empty "attestation_policy_static_id"
assert_not_empty "attestation_policy_kubernetes_id"

# store_svid round-trips through both the resource and the data source. The
# data source read is the interesting half: it is served by a separate schema,
# which omitted store_svid entirely until this was covered.
assert_eq "attestation_policy_static_store_svid_resource" "true"
assert_eq "attestation_policy_static_store_svid_data_source" "true"

# The data source ignored tpm_node policies entirely before it was switched to
# the shared conversion.
assert_eq "attestation_policy_tpm_node_ek_hash" \
  "5b3e0a049837688b09028ba84be190720bcc8f6cf74a487dc53b2ce9f376b5fb"

selector_count=$(echo "$outputs" | jq '.attestation_policy_tpm_node_selector_values.value | length')
if [[ "$selector_count" != "1" ]]; then
  echo "ERROR: expected 1 tpm_node selector value, got $selector_count" >&2
  exit 1
fi
assert_jq_eq ".attestation_policy_tpm_node_selector_values.value[0]" "test-selector"

assert_eq "attestation_policy_kubernetes_spiffe_id_path_template" "test/workload"
