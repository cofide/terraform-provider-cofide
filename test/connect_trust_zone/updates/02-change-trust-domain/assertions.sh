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

# The recreated trust zone has the new trust domain and keeps the name, as seen
# by both the resource and a fresh read via the data source.
assert_eq "trust_zone_trust_domain" "test-tz-changed.cofide.dev"
assert_eq "trust_zone_trust_domain_data_source" "test-tz-changed.cofide.dev"
assert_eq "trust_zone_name" "test-tz-renamed"
assert_eq "trust_zone_name_data_source" "test-tz-renamed"

# Connect-computed fields are populated on the recreated trust zone.
assert_not_empty "trust_zone_bundle_endpoint_url"
assert_not_empty "trust_zone_bundle_endpoint_profile"
assert_not_empty "trust_zone_jwt_issuer"
