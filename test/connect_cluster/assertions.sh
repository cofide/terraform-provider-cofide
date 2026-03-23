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

assert_not_empty "cluster_id"
assert_eq "cluster_trust_provider_kind" "kubernetes"
assert_eq "cluster_trust_provider_k8s_psat_enabled" "true"
assert_eq "cluster_trust_provider_k8s_psat_api_server_url" "https://kubernetes.default.svc"
assert_eq "cluster_trust_provider_k8s_psat_spire_server_audience" "spire-server"
assert_not_empty "cluster_trust_provider_k8s_psat_api_server_ca_cert"
assert_not_empty "cluster_oidc_issuer_url"
assert_not_empty "cluster_oidc_issuer_ca_cert"

sa_namespace=$(echo "$outputs" | jq -r '.cluster_trust_provider_k8s_psat_allowed_service_accounts.value[0].namespace')
sa_name=$(echo "$outputs" | jq -r '.cluster_trust_provider_k8s_psat_allowed_service_accounts.value[0].service_account_name')
sa_count=$(echo "$outputs" | jq '.cluster_trust_provider_k8s_psat_allowed_service_accounts.value | length')

if [[ "$sa_count" != "1" ]]; then
  echo "ERROR: expected 1 allowed_service_account, got $sa_count" >&2
  exit 1
fi
if [[ "$sa_namespace" != "spire" ]]; then
  echo "ERROR: expected service account namespace 'spire', got '$sa_namespace'" >&2
  exit 1
fi
if [[ "$sa_name" != "spire-agent" ]]; then
  echo "ERROR: expected service account name 'spire-agent', got '$sa_name'" >&2
  exit 1
fi
