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
CLUSTERID=${CI_PIPELINE_ID:-"local"}
CLUSTER_NAME="it-${CLUSTERID}"
. "${cur_dir}/common.sh"
NS=cloud-site-manager
CERTMGR_NS=cert-manager
set -e

is_cluster_ready() {
  kubectl get nodes | grep Ready | grep -v NotReady >& /dev/null
}

is_csm_ready() {
  kubectl get pods -n ${NS} | grep "1/1" | grep Running >& /dev/null
}

is_cert_manager_ready() {
  kubectl get pods -n ${CERTMGR_NS} | grep "2/2" | grep Running >& /dev/null
}

is_cluster_issuer_ready() {
  STATUS="$(kubectl get clusterissuer vault-issuer -o jsonpath='{.status.conditions[0].status}')"
  if [ "$STATUS" = 'True' ]; then
    return 0
  fi
  return 1
}

verify_cluster_issuer() {
  poll_immediate 1 60 is_cluster_issuer_ready
}

is_certificate_ready() {
  kubectl get certificate test-cert -o jsonpath='{.status.conditions[0].status}'
  if [ "$STATUS" = 'True' ]; then
    return 0
  fi
  return 1
}

verify_certificate() {
  poll_immediate 1 60 is_certificate_ready
}

wait_for_cluster() {
  poll_immediate 1 120 is_cluster_ready
}

verify_csm_ready() {
  poll_immediate 1 240 is_csm_ready
}

verify_cert_manager_ready() {
  poll_immediate 1 240 is_cert_manager_ready
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
  CMD="curl -k -s -X POST https://sitemgr.${NS}:8100/v1/sitecreds -d '{\"siteuuid\": \"$1\", \"otp\": \"$2\"}' -H 'Content-Type: application/json'"
  echo "Command is $CMD"
  kubectl exec -it nettools-pod -- bash -c "$CMD"
}

get_credentials_resp() {
  CMD="curl -k -s -X POST -o /dev/null -w \"%{http_code}\" https://sitemgr.${NS}:8100/v1/sitecreds -d '{\"siteuuid\": \"$1\", \"otp\": \"$2\"}' -H 'Content-Type: application/json'"
  echo "Command is $CMD"
  HTTP_RESP=$(kubectl exec -it nettools-pod -- bash -c "$CMD")
  if [ "$HTTP_RESP" = '500' ]; then
    echo "PASS"
  else
    echo "FAIL: HTTP_RESP is ${HTTP_RESP}, exp: 500."
    return 1
  fi
}

cleanup() {
  echo ""
  echo "|| clean up kind cluster ||"
  kind delete cluster --name "${CLUSTER_NAME}"
  rm -f /tmp/cacert.pem /tmp/ca.key
}

read_otp() {
  kubectl get site site-test-site-1234567890 -n ${NS} -o jsonpath='{.status.otp.passcode}'
}

setup_ca() {
  openssl req -x509 -sha256 -nodes -newkey rsa:2048 -keyout /tmp/ca.key -out /tmp/cacert.pem -days 3650 -subj "/C=US/ST=New Sweden/L=Stockholm /O=ut/OU=ut/CN=integ-test-ca/emailAddress=..."

  kubectl create secret -n "${CERTMGR_NS}" generic vault-root-ca-certificate --from-file=certificate=/tmp/cacert.pem
  kubectl create secret -n "${CERTMGR_NS}" generic vault-root-ca-private-key --from-file=privatekey=/tmp/ca.key
}

update_cert_manager() {
  LSIP=$(ifconfig en0 | grep netmask | awk '{print $2}')
  LSEP=$(echo -n "${LSIP}:8360" | base64)
  kubectl patch secret otel-lightstep -n cert-manager -p="{\"data\":{\"OTEL_EXPORTER_OTLP_SPAN_ENDPOINT\": \"$LSEP\"}}" -v=1
  kubectl scale deployment credsmgr -n cert-manager --replicas=0
  sleep 1
  kubectl scale deployment credsmgr -n cert-manager --replicas=1
  sleep 5
}

main() {
  tc_banner "create a kind cluster"
  trap cleanup EXIT
  kind create cluster --name "${CLUSTER_NAME}"

  tc_banner "wait for cluster ready"
  wait_for_cluster

  tc_banner "deploy cert-manager"
  kubectl create ns "${CERTMGR_NS}"
  setup_ca
  kubectl create secret -n "${CERTMGR_NS}" docker-registry credsmgr-image-pull-secret --from-file=.dockerconfigjson=$HOME/.docker/config.json
  kubectl apply -f "${cur_dir}/manifests/cert-manager/templates/"
  update_cert_manager
  tc_banner "verify cert-manager ready"
  verify_cert_manager_ready

  tc_banner "deploy csm"
  kubectl create ns "${NS}"
  kubectl create secret -n "${NS}" docker-registry csm-image-pull-secret --from-file=.dockerconfigjson=$HOME/.docker/config.json
  LSIP=$(ifconfig en0 | grep netmask | awk '{print $2}')
  OTEL_EXPORTER_OTLP_SPAN_ENDPOINT="${LSIP}:8360" kubectl apply -k "${cur_dir}/../kustomize/site-manager/overlays/local"

  tc_banner "verify csm ready"
  verify_csm_ready

  tc_banner "set up test pod"
  setup_test_pod
  
  tc_banner "fetch and install ca cert in test pod"
  fetch_ca_cert

  tc_banner "generate a cloud cert from test pod"
  get_cloud_cert

  tc_banner "create a site"
  create_site "test-site-1234567890"

  tc_banner "read otp"
  otp=$(read_otp)
  echo "otp is $otp"

  tc_banner "get site credentials"
  get_credentials "test-site-1234567890" "$otp"

  tc_banner "repeat get site credentials"

  get_credentials_resp "test-site-1234567890" "$otp"

}

main "$@"

