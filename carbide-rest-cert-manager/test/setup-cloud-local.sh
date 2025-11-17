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
CERTMGR_NS=cert-manager
set -e

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

verify_cert_manager_ready() {
  poll_immediate 1 240 is_cert_manager_ready
}

setup_ca() {
  openssl req -x509 -sha256 -nodes -newkey rsa:4096 -keyout /tmp/ca.key -out /tmp/cacert.pem -days 3650 -subj "/C=US/ST=CA /L=Local /O=Nvidia Dev/OU=Dev/CN=local.carbide-test.nvidia.com/emailAddress=..."

  kubectl delete secret -n "${CERTMGR_NS}" generic vault-root-ca-certificate --ignore-not-found
  kubectl create secret -n "${CERTMGR_NS}" generic vault-root-ca-certificate --from-file=certificate=/tmp/cacert.pem

  kubectl delete secret -n "${CERTMGR_NS}" generic vault-root-ca-private-key --ignore-not-found
  kubectl create secret -n "${CERTMGR_NS}" generic vault-root-ca-private-key --from-file=privatekey=/tmp/ca.key
}

cleanup() {
  rm -f /tmp/cacert.pem /tmp/ca.key
}

main() {
  trap cleanup EXIT

  kubectl get ns "${CERTMGR_NS}" || kubectl create ns "${CERTMGR_NS}"

  banner "🔄 Setting up CA"
  setup_ca

  banner "🔄 Deploying cert-manager"
  kubectl apply -k "${cur_dir}/../kustomize/cert-manager/overlays/cloud-local"
  banner "🔄 Verifying cert-manager ready"
  verify_cert_manager_ready

  banner "🔄 Deploying cluster-issuer"
  kubectl apply -f "${cur_dir}/manifests/cert-manager-resources/"
  banner "🔄 Verifying cluster-issuer"
  verify_cluster_issuer

  echo "✅ Cloud Cert Manager installation completed successfully"
}

main "$@"
