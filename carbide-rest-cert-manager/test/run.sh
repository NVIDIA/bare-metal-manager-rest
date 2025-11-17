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
CERTMGR_NS=cert-manager
set -e

is_cluster_ready() {
  kubectl get nodes | grep Ready | grep -v NotReady >& /dev/null
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
  STATUS="$(kubectl get certificate test-cert -o jsonpath='{.status.conditions[0].status}')"
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

verify_cert_manager_ready() {
  poll_immediate 1 480 is_cert_manager_ready
}

setup_ca() {
  openssl req -x509 -sha256 -nodes -newkey rsa:2048 -keyout /tmp/ca.key -out /tmp/cacert.pem -days 3650 -subj "/C=US/ST=New Sweden/L=Stockholm /O=ut/OU=ut/CN=integ-test-ca/emailAddress=..."

  kubectl create secret -n "${CERTMGR_NS}" generic vault-root-ca-certificate --from-file=certificate=/tmp/cacert.pem
  kubectl create secret -n "${CERTMGR_NS}" generic vault-root-ca-private-key --from-file=privatekey=/tmp/ca.key
}

cleanup() {
  echo ""
  echo "|| clean up kind cluster ||"
  kind delete cluster --name "${CLUSTER_NAME}"
  rm -f /tmp/cacert.pem /tmp/ca.key
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
  LSIP=$(ifconfig en0 | grep netmask | awk '{print $2}')
  OTEL_EXPORTER_OTLP_SPAN_ENDPOINT="${LSIP}:8360" kubectl apply -k "${cur_dir}/../kustomize/cert-manager/overlays/local"
  tc_banner "verify cert-manager ready"
  verify_cert_manager_ready

  tc_banner "deploy cluster-issuer"
  kubectl apply -f "${cur_dir}/manifests/cert-manager-resources/"
  tc_banner "verify cluster-issuer"
  verify_cluster_issuer

  tc_banner "verify certificate"
  verify_certificate
}

main "$@"

