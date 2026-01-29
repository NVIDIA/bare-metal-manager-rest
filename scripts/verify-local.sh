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

# Verification script for Carbide REST local development environment
set -e

API_URL="http://localhost:8388"
KEYCLOAK_URL="http://localhost:8080"
TEMPORAL_URL="http://localhost:8233"
ADMINER_URL="http://localhost:8081"
NAMESPACE="${NAMESPACE:-carbide}"

pass() { echo "[OK] $1"; }
fail() { echo "[FAIL] $1"; exit 1; }
warn() { echo "[WARN] $1"; }

echo "Verifying local deployment..."
echo ""

# Health check
echo -n "API health... "
if curl -sf "$API_URL/healthz" 2>/dev/null | jq -e '.is_healthy == true' > /dev/null 2>&1; then
  pass "healthy"
else
  fail "not healthy"
fi

# Readiness check
echo -n "API readiness... "
if curl -sf "$API_URL/readyz" 2>/dev/null | jq -e '.is_healthy == true' > /dev/null 2>&1; then
  pass "ready"
else
  warn "not ready"
fi

# Keycloak realm
echo -n "Keycloak realm... "
if curl -sf "$KEYCLOAK_URL/realms/carbide-dev" 2>/dev/null | jq -e '.realm == "carbide-dev"' > /dev/null 2>&1; then
  pass "available"
else
  fail "not available"
fi

# Get access token
echo -n "Token acquisition... "
TOKEN=$(curl -sf -X POST "$KEYCLOAK_URL/realms/carbide-dev/protocol/openid-connect/token" \
  -H "Content-Type: application/x-www-form-urlencoded" \
  -d "client_id=carbide-api" \
  -d "client_secret=carbide-local-secret" \
  -d "grant_type=password" \
  -d "username=testuser@example.com" \
  -d "password=testpassword" 2>/dev/null | jq -r .access_token)

if [ -n "$TOKEN" ] && [ "$TOKEN" != "null" ]; then
  pass "success"
else
  fail "failed"
fi

# Authenticated request
echo -n "Authenticated request... "
if curl -sf "$API_URL/healthz" -H "Authorization: Bearer $TOKEN" 2>/dev/null | jq -e '.is_healthy == true' > /dev/null 2>&1; then
  pass "success"
else
  fail "failed"
fi

# Temporal UI
echo -n "Temporal UI... "
if curl -sf -o /dev/null "$TEMPORAL_URL" 2>/dev/null; then
  pass "accessible"
else
  warn "not accessible"
fi

# Native cert-manager service
echo -n "Native cert manager... "
CERT_MANAGER_READY=$(kubectl -n "$NAMESPACE" get deployment carbide-rest-cert-manager -o jsonpath='{.status.readyReplicas}' 2>/dev/null || echo "")
if [ -n "$CERT_MANAGER_READY" ] && [ "$CERT_MANAGER_READY" -ge 1 ] 2>/dev/null; then
  pass "ready"
else
  warn "not ready"
fi

# Adminer
echo -n "Adminer... "
if curl -sf -o /dev/null "$ADMINER_URL" 2>/dev/null; then
  pass "accessible"
else
  warn "not accessible"
fi

echo ""
kubectl -n carbide get pods 2>/dev/null || warn "Could not get pod status"
echo ""
echo "Verification complete."
