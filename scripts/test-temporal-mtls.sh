#!/bin/bash
# SPDX-FileCopyrightText: Copyright (c) 2021-2023 NVIDIA CORPORATION & AFFILIATES. All rights reserved.
# SPDX-License-Identifier: LicenseRef-NvidiaProprietary
#
# Test script for Temporal mTLS verification
# This script verifies that the Temporal multi-pod deployment with mTLS is working correctly.

set -e

NAMESPACE="${NAMESPACE:-carbide}"
TIMEOUT="${TIMEOUT:-60}"

echo "=========================================="
echo "Temporal mTLS Verification Test"
echo "=========================================="
echo ""

# Colors for output (disabled to follow no-emoji/no-color rule)
PASS="PASS"
FAIL="FAIL"

test_passed=0
test_failed=0

# Function to run a test
run_test() {
    local name=$1
    local cmd=$2
    
    echo -n "  $name... "
    if eval "$cmd" > /dev/null 2>&1; then
        echo "$PASS"
        ((test_passed++))
        return 0
    else
        echo "$FAIL"
        ((test_failed++))
        return 1
    fi
}

# Step 1: Check Temporal pods are running
echo "Step 1: Checking Temporal pods..."
echo ""

run_test "Frontend pod running" \
    "kubectl -n $NAMESPACE get pods -l app.kubernetes.io/name=temporal,app.kubernetes.io/component=frontend -o jsonpath='{.items[0].status.phase}' | grep -q Running"

run_test "History pod running" \
    "kubectl -n $NAMESPACE get pods -l app.kubernetes.io/name=temporal,app.kubernetes.io/component=history -o jsonpath='{.items[0].status.phase}' | grep -q Running"

run_test "Matching pod running" \
    "kubectl -n $NAMESPACE get pods -l app.kubernetes.io/name=temporal,app.kubernetes.io/component=matching -o jsonpath='{.items[0].status.phase}' | grep -q Running"

run_test "Worker pod running" \
    "kubectl -n $NAMESPACE get pods -l app.kubernetes.io/name=temporal,app.kubernetes.io/component=worker -o jsonpath='{.items[0].status.phase}' | grep -q Running"

echo ""

# Step 2: Check TLS certificates exist
echo "Step 2: Checking TLS certificates..."
echo ""

run_test "server-interservice-certs secret exists" \
    "kubectl -n $NAMESPACE get secret server-interservice-certs"

run_test "server-cloud-certs secret exists" \
    "kubectl -n $NAMESPACE get secret server-cloud-certs"

run_test "server-site-certs secret exists" \
    "kubectl -n $NAMESPACE get secret server-site-certs"

run_test "temporal-client-certs secret exists" \
    "kubectl -n $NAMESPACE get secret temporal-client-certs"

run_test "temporal-client-certs contains client cert" \
    "kubectl -n $NAMESPACE get secret temporal-client-certs -o jsonpath='{.data.tls\\.crt}' | base64 -d | openssl x509 -noout > /dev/null"

echo ""

# Step 3: Verify certificates are issued by carbide-ca
echo "Step 3: Verifying certificate chain..."
echo ""

verify_issuer() {
    local secret_name=$1
    # Check for either the CA name used in local dev or the ClusterIssuer name
    local expected_issuer="Carbide"
    
    issuer=$(kubectl -n $NAMESPACE get secret "$secret_name" -o jsonpath='{.data.tls\.crt}' | \
        base64 -d | openssl x509 -noout -issuer 2>/dev/null | grep -o "CN=[^,]*" | head -1)
    
    echo "$issuer" | grep -qi "$expected_issuer"
}

run_test "server-interservice-certs issued by CA" \
    "verify_issuer server-interservice-certs"

run_test "temporal-client-certs issued by CA" \
    "verify_issuer temporal-client-certs"

echo ""

# Step 4: Check certificate status via cert-manager
echo "Step 4: Checking cert-manager Certificate status..."
echo ""

run_test "server-interservice-cert Certificate Ready" \
    "kubectl -n $NAMESPACE get certificate server-interservice-cert -o jsonpath='{.status.conditions[?(@.type==\"Ready\")].status}' | grep -q True"

run_test "server-cloud-cert Certificate Ready" \
    "kubectl -n $NAMESPACE get certificate server-cloud-cert -o jsonpath='{.status.conditions[?(@.type==\"Ready\")].status}' | grep -q True"

run_test "server-site-cert Certificate Ready" \
    "kubectl -n $NAMESPACE get certificate server-site-cert -o jsonpath='{.status.conditions[?(@.type==\"Ready\")].status}' | grep -q True"

run_test "temporal-client-cert Certificate Ready" \
    "kubectl -n $NAMESPACE get certificate temporal-client-cert -o jsonpath='{.status.conditions[?(@.type==\"Ready\")].status}' | grep -q True"

echo ""

# Step 5: Check Temporal frontend logs for TLS
echo "Step 5: Checking Temporal logs for TLS configuration..."
echo ""

frontend_pod=$(kubectl -n $NAMESPACE get pods -l app.kubernetes.io/name=temporal,app.kubernetes.io/component=frontend -o jsonpath='{.items[0].metadata.name}' 2>/dev/null)

if [ -n "$frontend_pod" ]; then
    run_test "Frontend loading TLS certificates" \
        "kubectl -n $NAMESPACE logs $frontend_pod --tail=100 2>/dev/null | grep -q 'loading certificate from file'"
    
    run_test "Frontend loading CA certs" \
        "kubectl -n $NAMESPACE logs $frontend_pod --tail=100 2>/dev/null | grep -q 'loading CA certs from'"
    
    run_test "Frontend is healthy" \
        "kubectl -n $NAMESPACE logs $frontend_pod --tail=100 2>/dev/null | grep -q 'Frontend is now healthy'"
else
    echo "  WARNING: Could not find frontend pod"
fi

echo ""

# Step 6: Check worker pod for SDK client connection
echo "Step 6: Checking worker SDK client connection..."
echo ""

worker_pod=$(kubectl -n $NAMESPACE get pods -l app.kubernetes.io/name=temporal,app.kubernetes.io/component=worker -o jsonpath='{.items[0].metadata.name}' 2>/dev/null)

if [ -n "$worker_pod" ]; then
    run_test "Worker started successfully" \
        "kubectl -n $NAMESPACE logs $worker_pod --tail=100 2>/dev/null | grep -q 'worker service started'"
    
    run_test "Worker workflows running" \
        "kubectl -n $NAMESPACE logs $worker_pod --tail=100 2>/dev/null | grep -q 'workflow successfully started'"
else
    echo "  WARNING: Could not find worker pod"
fi

echo ""

# Step 7: Check service discovery
echo "Step 7: Checking service discovery..."
echo ""

run_test "temporal-frontend service exists" \
    "kubectl -n $NAMESPACE get service temporal-frontend"

run_test "temporal-history service exists" \
    "kubectl -n $NAMESPACE get service temporal-history-headless"

run_test "temporal-matching service exists" \
    "kubectl -n $NAMESPACE get service temporal-matching-headless"

run_test "temporal-worker service exists" \
    "kubectl -n $NAMESPACE get service temporal-worker-headless"

echo ""

# Summary
echo "=========================================="
echo "Test Results"
echo "=========================================="
echo ""
echo "  Passed: $test_passed"
echo "  Failed: $test_failed"
echo ""

if [ $test_failed -eq 0 ]; then
    echo "SUCCESS: All tests passed"
    exit 0
else
    echo "FAILURE: $test_failed test(s) failed"
    exit 1
fi
