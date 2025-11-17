#!/bin/bash
# SPDX-FileCopyrightText: Copyright (c) 2021-2023 NVIDIA CORPORATION & AFFILIATES. All rights reserved.
# SPDX-License-Identifier: LicenseRef-NvidiaProprietary
#
# NVIDIA CORPORATION, its affiliates and licensors retain all intellectual
# property and proprietary rights in and to this material, related
# documentation and any modifications thereto. Any use, reproduction,
# disclosure or distribution of this material and related documentation
# without an express license agreement from NVIDIA CORPORATION or
# its affiliates is strictly prohibited.
#

readonly cur_dir="$(cd "$(dirname "$0")";pwd)"
. "${cur_dir}/common.sh"
NS=cloud-site-manager
CERTMGR_NS=cert-manager
set -e

is_csm_ready() {
  kubectl get pods -n ${NS} | grep "1/1" | grep Running >& /dev/null
}

verify_csm_ready() {
  poll_immediate 1 240 is_csm_ready
}

test_pod_ready() {
    kubectl get pods | grep Runn | grep nettools >& /dev/null
}

setup_test_pod() {
  if kubectl apply -f "${cur_dir}"/nettools.yaml; then
    :
  else
    return 1
  fi

  poll_immediate 1 90 test_pod_ready
}

fetch_ca_cert() {
  CMD="curl -k -s https://credsmgr.${CERTMGR_NS}:8000/v1/pki/ca/pem -o /tmp/cacert.pem"
  echo "Command is $CMD"
  kubectl exec -it nettools-pod -- bash -c "$CMD"
}

get_cloud_cert() {
  CMD="curl -s -X POST --cacert /tmp/cacert.pem https://credsmgr.${CERTMGR_NS}:8000/v1/pki/cloud-cert -d '{\"name\": \"cloud\", \"app\": \"server\", \"ttl\": 2160}' -H 'Content-Type: application/json' | jq -r > /tmp/creds.json && jq -r '.certificate' < /tmp/creds.json > /tmp/cert && jq -r '.key' < /tmp/creds.json > /tmp/key"
  echo "Command is $CMD"
  kubectl exec -it nettools-pod -- bash -c "$CMD"

  CMD="openssl x509 -in /tmp/cert -text | grep 'DNS:'"
  echo "Command is $CMD"
  RES=$(kubectl exec -it nettools-pod -- bash -c "$CMD")
  DNS="$(echo $RES | xargs | awk -F"\r" '{ print $1 }')"
  echo "DNS is $DNS"
  if [ "$DNS" = 'DNS:server.cloud.temporal.forge.nvidia.com' ]; then
    echo "Success"
  else
    echo "FAIL Expected DNS:server.cloud.temporal.forge.nvidia.com, got $DNS"
    return 1
  fi

  RES=$(kubectl exec -it nettools-pod -- bash -c "wc -l /tmp/key")
  WC=$(echo "$RES" | awk '{print $1}')
  if [[ "${WC}" -lt 1 ]]; then
    echo "FAIL: key is empty"
    return 1
  fi
}

create_site() {
  CMD="curl -k -s -X POST https://sitemgr.${NS}:8100/v1/site -d '{\"siteuuid\": \"$1\", \"provider\": \"ip1\", \"fcorg\": \"ip1\"}' -H 'Content-Type: application/json' | jq -r"
  echo "Command is $CMD"
  kubectl exec -it nettools-pod -- bash -c "$CMD"
}

get_credentials() {
  CMD="curl -k -s -X POST https://sitemgr.${NS}:8100/v1/sitecreds -d '{\"siteuuid\": \"$1\", \"otp\": \"$2\"}' -H 'Content-Type: application/json' | jq -r > /tmp/site-${1}.json"
  echo "Command is $CMD"
  kubectl exec -it nettools-pod -- bash -c "$CMD"
  kubectl cp nettools-pod:"/tmp/site-${1}.json" "/tmp/site-${1}.json"
  mkdir -p "/tmp/site-${1}"
  jq -r '.["cacertificate"]' < "/tmp/site-${1}.json" > "/tmp/site-${1}/ca.crt"
  jq -r '.["key"]' < "/tmp/site-${1}.json" > "/tmp/site-${1}/tls.key"
  jq -r '.["certificate"]' < "/tmp/site-${1}.json" > "/tmp/site-${1}/tls.crt"
  echo "Certificates are stored in /tmp/site-${1}"
}

read_otp() {
  kubectl get site "site-${1}" -n ${NS} -o jsonpath='{.status.otp.passcode}'
}

main() {
  if [ "$#" -ne 1 ]; then
    echo "Usage $0 <site-uuid>"
    exit 1
  fi
  
  banner "get certs for site-id $1"
  banner "verify csm ready"
  verify_csm_ready

  banner "set up test pod"
  setup_test_pod
  
  banner "fetch and install ca cert in test pod"
  fetch_ca_cert

  banner "create a site"
  create_site "$1"

  banner "read otp"
  otp=$(read_otp ${1})
  echo "otp is $otp"

  banner "get site credentials"
  get_credentials "$1" "$otp"

}

main "$@"

