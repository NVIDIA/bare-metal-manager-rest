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

# Setup PKI for local development

set -e

NAMESPACE="${NAMESPACE:-carbide}"
VAULT_ADDR="${VAULT_ADDR:-http://localhost:8200}"
VAULT_TOKEN="${VAULT_TOKEN:-root}"

echo "Setting up local PKI..."
echo ""

generate_ca() {
  echo "Generating CA certificate..."
  
  CA_DIR=$(mktemp -d)
  trap "rm -rf $CA_DIR" EXIT
  
  # Create CA config with proper key usage for CRL signing
  cat > "$CA_DIR/ca.cnf" << 'EOFCNF'
[req]
distinguished_name = req_distinguished_name
x509_extensions = v3_ca
prompt = no

[req_distinguished_name]
C = US
ST = CA
L = Local
O = Carbide Dev
OU = Dev
CN = carbide-local-ca

[v3_ca]
basicConstraints = critical,CA:TRUE
keyUsage = critical,keyCertSign,cRLSign,digitalSignature
subjectKeyIdentifier = hash
authorityKeyIdentifier = keyid:always,issuer:always
EOFCNF

  openssl req -x509 -sha256 -nodes -newkey rsa:4096 \
    -keyout "$CA_DIR/ca.key" \
    -out "$CA_DIR/ca.crt" \
    -days 3650 \
    -config "$CA_DIR/ca.cnf" \
    -extensions v3_ca
  
  kubectl create secret generic vault-root-ca-certificate \
    --from-file=certificate="$CA_DIR/ca.crt" \
    -n "$NAMESPACE" \
    --dry-run=client -o yaml | kubectl apply -f -
  
  kubectl create secret generic vault-root-ca-private-key \
    --from-file=privatekey="$CA_DIR/ca.key" \
    -n "$NAMESPACE" \
    --dry-run=client -o yaml | kubectl apply -f -
  
  VAULT_INIT_RESPONSE='{"root_token":"root","keys":["aabbccdd"],"keys_base64":["aabbccdd"]}'
  kubectl create secret generic vault-token \
    --from-literal=vault-token="$VAULT_INIT_RESPONSE" \
    --from-literal=certmgr-token="root" \
    -n "$NAMESPACE" \
    --dry-run=client -o yaml | kubectl apply -f -
  
  echo "CA secrets created"
}

configure_vault_pki() {
  echo "Configuring Vault PKI..."
  
  if ! curl -sf "$VAULT_ADDR/v1/sys/health" > /dev/null 2>&1; then
    echo "Vault not accessible, skipping"
    return 0
  fi
  
  curl -sf -X POST \
    -H "X-Vault-Token: $VAULT_TOKEN" \
    "$VAULT_ADDR/v1/sys/mounts/pki" \
    -d '{"type":"pki"}' 2>/dev/null || true
  
  curl -sf -X POST \
    -H "X-Vault-Token: $VAULT_TOKEN" \
    "$VAULT_ADDR/v1/sys/mounts/pki/tune" \
    -d '{"max_lease_ttl":"87600h"}'
  
  curl -sf -X POST \
    -H "X-Vault-Token: $VAULT_TOKEN" \
    "$VAULT_ADDR/v1/pki/root/generate/internal" \
    -d '{"common_name":"Carbide Local Dev CA","issuer_name":"root-2024","ttl":"87600h"}' > /dev/null 2>&1 || true
  
  curl -sf -X POST \
    -H "X-Vault-Token: $VAULT_TOKEN" \
    "$VAULT_ADDR/v1/pki/config/urls" \
    -d '{"issuing_certificates":["http://vault:8200/v1/pki/ca"],"crl_distribution_points":["http://vault:8200/v1/pki/crl"]}'
  
  curl -sf -X POST \
    -H "X-Vault-Token: $VAULT_TOKEN" \
    "$VAULT_ADDR/v1/pki/roles/cloud-cert" \
    -d '{"allowed_domains":"carbide.local,localhost,carbide,svc.cluster.local,carbide.svc.cluster.local,carbide-rest-cert-manager,carbide-rest-site-manager","allow_subdomains":true,"allow_any_name":true,"max_ttl":"720h"}'
  
  echo "Vault PKI configured"
}

update_site_agent_secret() {
  echo "Updating site-agent secrets..."
  
  CERT_DIR=$(mktemp -d)
  
  CA_CERT=$(kubectl get secret vault-root-ca-certificate -n "$NAMESPACE" -o jsonpath='{.data.certificate}' | base64 -d 2>/dev/null || echo "")
  CA_KEY=$(kubectl get secret vault-root-ca-private-key -n "$NAMESPACE" -o jsonpath='{.data.privatekey}' | base64 -d 2>/dev/null || echo "")
  
  if [ -n "$CA_CERT" ] && [ -n "$CA_KEY" ]; then
    echo "$CA_CERT" > "$CERT_DIR/ca.crt"
    echo "$CA_KEY" > "$CERT_DIR/ca.key"
    
    openssl genrsa -out "$CERT_DIR/tls.key" 2048 2>/dev/null
    openssl req -new -key "$CERT_DIR/tls.key" -out "$CERT_DIR/tls.csr" \
      -subj "/C=US/ST=CA/L=Local/O=Carbide Dev/CN=carbide-rest-site-agent" 2>/dev/null
    openssl x509 -req -in "$CERT_DIR/tls.csr" -CA "$CERT_DIR/ca.crt" -CAkey "$CERT_DIR/ca.key" \
      -CAcreateserial -out "$CERT_DIR/tls.crt" -days 365 2>/dev/null
    
    kubectl create secret generic carbide-tls-certs \
      --from-file=ca.crt="$CERT_DIR/ca.crt" \
      --from-file=tls.crt="$CERT_DIR/tls.crt" \
      --from-file=tls.key="$CERT_DIR/tls.key" \
      -n "$NAMESPACE" \
      --dry-run=client -o yaml | kubectl apply -f -
    
    kubectl create secret generic site-registration \
      --from-literal=site-uuid="00000000-0000-4000-8000-000000000001" \
      --from-literal=otp="local-dev-otp-token" \
      --from-literal=creds-url="http://carbide-rest-site-manager:8100/v1/site/credentials" \
      --from-literal=cacert="$CA_CERT" \
      -n "$NAMESPACE" \
      --dry-run=client -o yaml | kubectl apply -f -
    
    echo "Site-agent secrets created"
  else
    echo "Warning: Could not retrieve CA certificate"
  fi
  
  rm -rf "$CERT_DIR"
}

create_site_manager_certs() {
  echo "Creating site-manager TLS certificates..."
  
  CERT_DIR=$(mktemp -d)
  
  CA_CERT=$(kubectl get secret vault-root-ca-certificate -n "$NAMESPACE" -o jsonpath='{.data.certificate}' | base64 -d 2>/dev/null || echo "")
  CA_KEY=$(kubectl get secret vault-root-ca-private-key -n "$NAMESPACE" -o jsonpath='{.data.privatekey}' | base64 -d 2>/dev/null || echo "")
  
  if [ -n "$CA_CERT" ] && [ -n "$CA_KEY" ]; then
    echo "$CA_CERT" > "$CERT_DIR/ca.crt"
    echo "$CA_KEY" > "$CERT_DIR/ca.key"
    
    openssl genrsa -out "$CERT_DIR/tls.key" 2048 2>/dev/null
    
    cat > "$CERT_DIR/san.cnf" << EOF
[req]
distinguished_name = req_distinguished_name
req_extensions = v3_req
prompt = no

[req_distinguished_name]
C = US
ST = CA
L = Local
O = Carbide Dev
CN = carbide-rest-site-manager

[v3_req]
basicConstraints = CA:FALSE
keyUsage = nonRepudiation, digitalSignature, keyEncipherment
subjectAltName = @alt_names

[alt_names]
DNS.1 = carbide-rest-site-manager
DNS.2 = carbide-rest-site-manager.carbide
DNS.3 = carbide-rest-site-manager.carbide.svc.cluster.local
DNS.4 = localhost
EOF

    openssl req -new -key "$CERT_DIR/tls.key" -out "$CERT_DIR/tls.csr" \
      -config "$CERT_DIR/san.cnf" 2>/dev/null
    openssl x509 -req -in "$CERT_DIR/tls.csr" -CA "$CERT_DIR/ca.crt" -CAkey "$CERT_DIR/ca.key" \
      -CAcreateserial -out "$CERT_DIR/tls.crt" -days 365 \
      -extensions v3_req -extfile "$CERT_DIR/san.cnf" 2>/dev/null
    
    kubectl create secret tls site-manager-tls \
      --cert="$CERT_DIR/tls.crt" \
      --key="$CERT_DIR/tls.key" \
      -n "$NAMESPACE" \
      --dry-run=client -o yaml | kubectl apply -f -
    
    echo "Site-manager TLS secret created"
  else
    echo "Warning: Could not create site-manager TLS certs"
  fi
  
  rm -rf "$CERT_DIR"
}

main() {
  kubectl get ns "$NAMESPACE" > /dev/null 2>&1 || kubectl create ns "$NAMESPACE"
  
  generate_ca
  configure_vault_pki
  update_site_agent_secret
  create_site_manager_certs
  
  echo ""
  echo "PKI setup complete."
}

main "$@"
