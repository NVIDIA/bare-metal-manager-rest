#!/bin/bash
# SPDX-FileCopyrightText: Copyright (c) 2021-2023 NVIDIA CORPORATION & AFFILIATES. All rights reserved.
# SPDX-License-Identifier: LicenseRef-NvidiaProprietary
#
# Test script for Temporal certificate rotation
# This script verifies that cert-manager.io properly rotates certificates
# and that Temporal continues working after rotation.

set -e

NAMESPACE="${NAMESPACE:-carbide}"
TIMEOUT="${TIMEOUT:-300}"

echo "=========================================="
echo "Temporal Certificate Rotation Test"
echo "=========================================="
echo ""

# Function to get certificate serial number from a secret
get_cert_serial() {
    local secret_name=$1
    kubectl -n "$NAMESPACE" get secret "$secret_name" -o jsonpath='{.data.tls\.crt}' 2>/dev/null | \
        base64 -d | openssl x509 -noout -serial 2>/dev/null | cut -d= -f2
}

# Function to check if Temporal frontend is healthy
check_temporal_health() {
    kubectl -n "$NAMESPACE" get pods -l app.kubernetes.io/name=temporal,app.kubernetes.io/component=frontend -o jsonpath='{.items[0].status.phase}' 2>/dev/null | grep -q "Running"
}

# Step 1: Record current certificate serial numbers
echo "Step 1: Recording current certificate serial numbers..."
echo ""

INTERSERVICE_SERIAL_BEFORE=$(get_cert_serial "server-interservice-certs")
CLOUD_SERIAL_BEFORE=$(get_cert_serial "server-cloud-certs")
SITE_SERIAL_BEFORE=$(get_cert_serial "server-site-certs")
CLIENT_SERIAL_BEFORE=$(get_cert_serial "temporal-client-certs")

echo "  server-interservice-certs: $INTERSERVICE_SERIAL_BEFORE"
echo "  server-cloud-certs:        $CLOUD_SERIAL_BEFORE"
echo "  server-site-certs:         $SITE_SERIAL_BEFORE"
echo "  temporal-client-certs:     $CLIENT_SERIAL_BEFORE"
echo ""

# Step 2: Verify Temporal is healthy before rotation
echo "Step 2: Verifying Temporal is healthy before rotation..."
if check_temporal_health; then
    echo "  Temporal frontend is healthy"
else
    echo "  ERROR: Temporal frontend is not healthy"
    exit 1
fi
echo ""

# Step 3: Delete certificate secrets to trigger rotation
echo "Step 3: Deleting certificate secrets to trigger rotation..."
kubectl -n "$NAMESPACE" delete secret server-interservice-certs --ignore-not-found
kubectl -n "$NAMESPACE" delete secret server-cloud-certs --ignore-not-found
kubectl -n "$NAMESPACE" delete secret server-site-certs --ignore-not-found
kubectl -n "$NAMESPACE" delete secret temporal-client-certs --ignore-not-found
echo "  Secrets deleted"
echo ""

# Step 4: Wait for cert-manager to reissue certificates
echo "Step 4: Waiting for cert-manager to reissue certificates..."
echo "  Waiting for server-interservice-cert..."
kubectl -n "$NAMESPACE" wait --for=condition=Ready certificate/server-interservice-cert --timeout="${TIMEOUT}s"
echo "  Waiting for server-cloud-cert..."
kubectl -n "$NAMESPACE" wait --for=condition=Ready certificate/server-cloud-cert --timeout="${TIMEOUT}s"
echo "  Waiting for server-site-cert..."
kubectl -n "$NAMESPACE" wait --for=condition=Ready certificate/server-site-cert --timeout="${TIMEOUT}s"
echo "  Waiting for temporal-client-cert..."
kubectl -n "$NAMESPACE" wait --for=condition=Ready certificate/temporal-client-cert --timeout="${TIMEOUT}s"
echo "  All certificates reissued"
echo ""

# Step 5: Verify new serial numbers are different
echo "Step 5: Verifying new certificate serial numbers..."
echo ""

INTERSERVICE_SERIAL_AFTER=$(get_cert_serial "server-interservice-certs")
CLOUD_SERIAL_AFTER=$(get_cert_serial "server-cloud-certs")
SITE_SERIAL_AFTER=$(get_cert_serial "server-site-certs")
CLIENT_SERIAL_AFTER=$(get_cert_serial "temporal-client-certs")

echo "  server-interservice-certs: $INTERSERVICE_SERIAL_AFTER"
echo "  server-cloud-certs:        $CLOUD_SERIAL_AFTER"
echo "  server-site-certs:         $SITE_SERIAL_AFTER"
echo "  temporal-client-certs:     $CLIENT_SERIAL_AFTER"
echo ""

ROTATION_SUCCESS=true

if [ "$INTERSERVICE_SERIAL_BEFORE" = "$INTERSERVICE_SERIAL_AFTER" ]; then
    echo "  ERROR: server-interservice-certs was not rotated"
    ROTATION_SUCCESS=false
else
    echo "  server-interservice-certs: ROTATED"
fi

if [ "$CLOUD_SERIAL_BEFORE" = "$CLOUD_SERIAL_AFTER" ]; then
    echo "  ERROR: server-cloud-certs was not rotated"
    ROTATION_SUCCESS=false
else
    echo "  server-cloud-certs: ROTATED"
fi

if [ "$SITE_SERIAL_BEFORE" = "$SITE_SERIAL_AFTER" ]; then
    echo "  ERROR: server-site-certs was not rotated"
    ROTATION_SUCCESS=false
else
    echo "  server-site-certs: ROTATED"
fi

if [ "$CLIENT_SERIAL_BEFORE" = "$CLIENT_SERIAL_AFTER" ]; then
    echo "  ERROR: temporal-client-certs was not rotated"
    ROTATION_SUCCESS=false
else
    echo "  temporal-client-certs: ROTATED"
fi
echo ""

# Step 6: Restart Temporal pods to pick up new certificates
echo "Step 6: Restarting Temporal pods to pick up new certificates..."
kubectl -n "$NAMESPACE" rollout restart deployment/temporal-frontend
kubectl -n "$NAMESPACE" rollout restart deployment/temporal-history
kubectl -n "$NAMESPACE" rollout restart deployment/temporal-matching
kubectl -n "$NAMESPACE" rollout restart deployment/temporal-worker
echo "  Rollout restart initiated"
echo ""

# Step 7: Wait for Temporal to be ready again
echo "Step 7: Waiting for Temporal pods to be ready..."
kubectl -n "$NAMESPACE" rollout status deployment/temporal-frontend --timeout="${TIMEOUT}s"
kubectl -n "$NAMESPACE" rollout status deployment/temporal-history --timeout="${TIMEOUT}s"
kubectl -n "$NAMESPACE" rollout status deployment/temporal-matching --timeout="${TIMEOUT}s"
kubectl -n "$NAMESPACE" rollout status deployment/temporal-worker --timeout="${TIMEOUT}s"
echo "  All Temporal pods ready"
echo ""

# Step 8: Verify Temporal is healthy after rotation
echo "Step 8: Verifying Temporal is healthy after rotation..."
if check_temporal_health; then
    echo "  Temporal frontend is healthy"
else
    echo "  ERROR: Temporal frontend is not healthy after rotation"
    exit 1
fi
echo ""

# Summary
echo "=========================================="
echo "Certificate Rotation Test Results"
echo "=========================================="
if [ "$ROTATION_SUCCESS" = true ]; then
    echo "SUCCESS: All certificates were rotated and Temporal is healthy"
    echo ""
    echo "Before rotation:"
    echo "  interservice: $INTERSERVICE_SERIAL_BEFORE"
    echo "  cloud:        $CLOUD_SERIAL_BEFORE"
    echo "  site:         $SITE_SERIAL_BEFORE"
    echo "  client:       $CLIENT_SERIAL_BEFORE"
    echo ""
    echo "After rotation:"
    echo "  interservice: $INTERSERVICE_SERIAL_AFTER"
    echo "  cloud:        $CLOUD_SERIAL_AFTER"
    echo "  site:         $SITE_SERIAL_AFTER"
    echo "  client:       $CLIENT_SERIAL_AFTER"
    exit 0
else
    echo "FAILURE: Some certificates were not rotated"
    exit 1
fi
