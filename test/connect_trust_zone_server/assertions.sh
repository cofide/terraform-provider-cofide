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

assert_not_empty() {
  local key=$1
  local actual
  actual=$(echo "$outputs" | jq -r ".$key.value")
  if [[ -z "$actual" || "$actual" == "null" ]]; then
    echo "ERROR: output '$key': expected non-empty value" >&2
    return 1
  fi
}

assert_not_empty "server_id"
assert_not_empty "server_trust_zone_id"
assert_not_empty "server_cluster_id"
assert_not_empty "server_org_id"

assert_eq "server_k8s_psat_spiffe_id_path" "/ns/spire/sa/spire-server"

audience_count=$(echo "$outputs" | jq '.server_k8s_psat_audiences.value | length')
if [[ "$audience_count" != "1" ]]; then
  echo "ERROR: expected 1 audience, got $audience_count" >&2
  exit 1
fi
audience=$(echo "$outputs" | jq -r '.server_k8s_psat_audiences.value[0]')
if [[ "$audience" != "spire-server" ]]; then
  echo "ERROR: expected audience 'spire-server', got '$audience'" >&2
  exit 1
fi

# Verify the list data source returns exactly the one server in the trust zone.
assert_eq "servers_by_trust_zone_count" "1"

server_id=$(echo "$outputs" | jq -r '.server_id.value')
assert_eq "servers_by_trust_zone_first_id" "$server_id"
