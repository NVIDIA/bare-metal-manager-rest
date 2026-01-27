#!/bin/bash
# SPDX-FileCopyrightText: Copyright (c) 2021-2023 NVIDIA CORPORATION & AFFILIATES. All rights reserved.
# SPDX-License-Identifier: LicenseRef-NvidiaProprietary
#
# Test script for Temporal certificate rotation
# This script verifies that cert-manager.io properly rotates certificates
# and that Temporal continues working after rotation.

set -e

NAMESPACE="${NAMESPACE:-carbide}"
TIMEOUT="${TIMEOUT:-120}"

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
    kubectl -n "$NAMESPACE" get pods -l app=temporal,component=frontend -o jsonpath='{.items[0].status.phase}' 2>/dev/null | grep -q "Running"
}

# Step 1: Record current certificate serial numbers
echo "Step 1: Recording current certificate serial numbers..."
echo ""

FRONTEND_SERIAL_BEFORE=$(get_cert_serial "temporal-frontend-tls")
HISTORY_SERIAL_BEFORE=$(get_cert_serial "temporal-history-tls")
MATCHING_SERIAL_BEFORE=$(get_cert_serial "temporal-matching-tls")
WORKER_SERIAL_BEFORE=$(get_cert_serial "temporal-worker-tls")

echo "  temporal-frontend-tls: $FRONTEND_SERIAL_BEFORE"
echo "  temporal-history-tls:  $HISTORY_SERIAL_BEFORE"
echo "  temporal-matching-tls: $MATCHING_SERIAL_BEFORE"
echo "  temporal-worker-tls:   $WORKER_SERIAL_BEFORE"
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
kubectl -n "$NAMESPACE" delete secret temporal-frontend-tls --ignore-not-found
kubectl -n "$NAMESPACE" delete secret temporal-history-tls --ignore-not-found
kubectl -n "$NAMESPACE" delete secret temporal-matching-tls --ignore-not-found
kubectl -n "$NAMESPACE" delete secret temporal-worker-tls --ignore-not-found
echo "  Secrets deleted"
echo ""

# Step 4: Wait for cert-manager to reissue certificates
echo "Step 4: Waiting for cert-manager to reissue certificates..."
echo "  Waiting for temporal-frontend-tls..."
kubectl -n "$NAMESPACE" wait --for=condition=Ready certificate/temporal-frontend-tls --timeout="${TIMEOUT}s"
echo "  Waiting for temporal-history-tls..."
kubectl -n "$NAMESPACE" wait --for=condition=Ready certificate/temporal-history-tls --timeout="${TIMEOUT}s"
echo "  Waiting for temporal-matching-tls..."
kubectl -n "$NAMESPACE" wait --for=condition=Ready certificate/temporal-matching-tls --timeout="${TIMEOUT}s"
echo "  Waiting for temporal-worker-tls..."
kubectl -n "$NAMESPACE" wait --for=condition=Ready certificate/temporal-worker-tls --timeout="${TIMEOUT}s"
echo "  All certificates reissued"
echo ""

# Step 5: Verify new serial numbers are different
echo "Step 5: Verifying new certificate serial numbers..."
echo ""

FRONTEND_SERIAL_AFTER=$(get_cert_serial "temporal-frontend-tls")
HISTORY_SERIAL_AFTER=$(get_cert_serial "temporal-history-tls")
MATCHING_SERIAL_AFTER=$(get_cert_serial "temporal-matching-tls")
WORKER_SERIAL_AFTER=$(get_cert_serial "temporal-worker-tls")

echo "  temporal-frontend-tls: $FRONTEND_SERIAL_AFTER"
echo "  temporal-history-tls:  $HISTORY_SERIAL_AFTER"
echo "  temporal-matching-tls: $MATCHING_SERIAL_AFTER"
echo "  temporal-worker-tls:   $WORKER_SERIAL_AFTER"
echo ""

ROTATION_SUCCESS=true

if [ "$FRONTEND_SERIAL_BEFORE" = "$FRONTEND_SERIAL_AFTER" ]; then
    echo "  ERROR: temporal-frontend-tls was not rotated"
    ROTATION_SUCCESS=false
else
    echo "  temporal-frontend-tls: ROTATED"
fi

if [ "$HISTORY_SERIAL_BEFORE" = "$HISTORY_SERIAL_AFTER" ]; then
    echo "  ERROR: temporal-history-tls was not rotated"
    ROTATION_SUCCESS=false
else
    echo "  temporal-history-tls: ROTATED"
fi

if [ "$MATCHING_SERIAL_BEFORE" = "$MATCHING_SERIAL_AFTER" ]; then
    echo "  ERROR: temporal-matching-tls was not rotated"
    ROTATION_SUCCESS=false
else
    echo "  temporal-matching-tls: ROTATED"
fi

if [ "$WORKER_SERIAL_BEFORE" = "$WORKER_SERIAL_AFTER" ]; then
    echo "  ERROR: temporal-worker-tls was not rotated"
    ROTATION_SUCCESS=false
else
    echo "  temporal-worker-tls: ROTATED"
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
    echo "  frontend: $FRONTEND_SERIAL_BEFORE"
    echo "  history:  $HISTORY_SERIAL_BEFORE"
    echo "  matching: $MATCHING_SERIAL_BEFORE"
    echo "  worker:   $WORKER_SERIAL_BEFORE"
    echo ""
    echo "After rotation:"
    echo "  frontend: $FRONTEND_SERIAL_AFTER"
    echo "  history:  $HISTORY_SERIAL_AFTER"
    echo "  matching: $MATCHING_SERIAL_AFTER"
    echo "  worker:   $WORKER_SERIAL_AFTER"
    exit 0
else
    echo "FAILURE: Some certificates were not rotated"
    exit 1
fi
