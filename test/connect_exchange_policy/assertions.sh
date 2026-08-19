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

assert_not_empty "allow_policy_id"

policy_count=$(echo "$outputs" | jq '.exchange_policy_ids.value | length')
if [[ "$policy_count" != "6" ]]; then
  echo "ERROR: expected 6 exchange policies, got $policy_count" >&2
  exit 1
fi

OUTBOUND_IDENTITY="spiffe://ep-tz.cofide.dev/ns/outbound/sa/exchanged"

# outbound_identity round-trips through the resource, the single-policy data
# source and the list data source.
assert_eq "outbound_identity_resource" "$OUTBOUND_IDENTITY"
assert_eq "outbound_identity_data_source" "$OUTBOUND_IDENTITY"
assert_eq "outbound_identity_from_list" "$OUTBOUND_IDENTITY"

# outbound_issuer.oauth_as round-trips through the resource, the single-policy
# data source and the list data source.
for out in outbound_issuer_resource outbound_issuer_data_source outbound_issuer_from_list; do
  assert_jq_eq ".$out.value.oauth_as.grant_type" "client_credentials"
  assert_jq_eq ".$out.value.oauth_as.issuer_url" "https://as.example.com"
  assert_jq_eq ".$out.value.oauth_as.token_url" "https://as.example.com/oauth2/token"
  assert_jq_eq ".$out.value.oauth_as.timeout" "30"
  assert_jq_eq ".$out.value.oauth_as.audiences | length" "1"
  assert_jq_eq ".$out.value.oauth_as.audiences[0]" "https://external-api.example.com"
done

# outbound_issuer.spiffe round-trips through the resource, the single-policy
# data source and the list data source. It is a marker with no fields, so the
# spiffe variant is a non-null empty object and oauth_as is null.
for out in spiffe_issuer_resource spiffe_issuer_data_source spiffe_issuer_from_list; do
  assert_jq_eq ".$out.value.spiffe != null" "true"
  assert_jq_eq ".$out.value.oauth_as" "null"
done

# A policy with neither field set reports both as null.
assert_jq_eq ".allow_policy_outbound_issuer.value" "null"
assert_jq_eq ".allow_policy_outbound_identity.value" "null"
