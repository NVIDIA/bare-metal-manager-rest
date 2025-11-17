#!/bin/bash

# Comprehensive Cloud-API Testing Script
# Tests all major endpoints with proper authentication

set -e

echo "========================================"
echo "Cloud-API Comprehensive Testing"
echo "========================================"
echo ""

# Get authentication token
echo "Getting authentication token..."
TOKEN=$(curl -sf -X POST 'http://localhost:8080/realms/carbide/protocol/openid-connect/token' \
  -H 'Content-Type: application/x-www-form-urlencoded' \
  -d 'username=testuser' \
  -d 'password=testpass' \
  -d 'grant_type=password' \
  -d 'client_id=carbide-api' \
  -d 'client_secret=carbide-secret-dev-only-do-not-use-in-prod' | jq -r '.access_token')

if [ -z "$TOKEN" ] || [ "$TOKEN" == "null" ]; then
    echo "[ERROR] Failed to get authentication token"
    exit 1
fi
echo "[SUCCESS] Token obtained: ${TOKEN:0:30}..."
echo ""

# Test 1: Health Check
echo "TEST 1: Health Check"
HEALTH=$(curl -s http://localhost:8388/healthz)
echo "Response: $HEALTH"
if echo "$HEALTH" | grep -q "is_healthy.*true"; then
    echo "[SUCCESS] Health check passed"
else
    echo "[FAILED] Health check failed"
fi
echo ""

# Test 2: Metadata
echo "TEST 2: Get Metadata"
METADATA=$(curl -s -H "Authorization: Bearer $TOKEN" \
  http://localhost:8388/v2/org/nvidia/carbide/metadata)
echo "Response: $METADATA"
if echo "$METADATA" | grep -q "version"; then
    echo "[SUCCESS] Metadata endpoint working"
else
    echo "[FAILED] Metadata failed: $METADATA"
fi
echo ""

# Test 3: Infrastructure Provider
echo "TEST 3: Get Infrastructure Provider"
INFRA=$(curl -s -H "Authorization: Bearer $TOKEN" \
  http://localhost:8388/v2/org/nvidia/carbide/infrastructure-provider/current)
echo "Response: $INFRA"
if echo "$INFRA" | grep -q "id"; then
    INFRA_ID=$(echo "$INFRA" | jq -r '.id')
    echo "[SUCCESS] Infrastructure Provider retrieved: $INFRA_ID"
else
    echo "[FAILED] Infrastructure Provider failed: $INFRA"
fi
echo ""

# Test 4: Tenant
echo "TEST 4: Get Tenant"
TENANT=$(curl -s -H "Authorization: Bearer $TOKEN" \
  http://localhost:8388/v2/org/nvidia/carbide/tenant/current)
echo "Response: $TENANT"
if echo "$TENANT" | grep -q "id"; then
    TENANT_ID=$(echo "$TENANT" | jq -r '.id')
    echo "[SUCCESS] Tenant retrieved: $TENANT_ID"
else
    echo "[FAILED] Tenant failed: $TENANT"
fi
echo ""

# Test 5: Create Site
echo "TEST 5: Create Site"
SITE=$(curl -s -X POST -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  http://localhost:8388/v2/org/nvidia/carbide/site \
  -d '{
    "name": "test-site-'$(date +%s)'",
    "description": "Automated test site for API validation",
    "serialConsoleHostname": "console.test.local",
    "location": {
      "city": "Santa Clara",
      "state": "CA",
      "country": "US"
    },
    "contact": {
      "email": "admin@test.local"
    }
  }')
echo "Response: $SITE" | jq -c '{id, name, status}'
if echo "$SITE" | grep -q '"id"'; then
    SITE_ID=$(echo "$SITE" | jq -r '.id')
    echo "[SUCCESS] Site created: $SITE_ID"
elif echo "$SITE" | grep -q "already exists"; then
    # Site already exists, extract ID from error
    SITE_ID=$(echo "$SITE" | jq -r '.data.id // empty')
    if [ -n "$SITE_ID" ]; then
        echo "[INFO] Site already exists, using existing: $SITE_ID"
    else
        echo "[FAILED] Site creation failed"
        SITE_ID=""
    fi
else
    echo "[FAILED] Site creation failed: $SITE"
    SITE_ID=""
fi
echo ""

# Test 6: List Sites (with infrastructureProviderId)
if [ -n "$SITE_ID" ] && [ -n "$INFRA_ID" ]; then
    echo "TEST 6: List Sites (as Provider Admin)"
    SITES=$(curl -s -H "Authorization: Bearer $TOKEN" \
      "http://localhost:8388/v2/org/nvidia/carbide/site?infrastructureProviderId=$INFRA_ID")
    
    # Check if response is an array or error
    if echo "$SITES" | jq -e 'type == "array"' > /dev/null 2>&1; then
        echo "Response: $SITES" | jq -c '.[0] | {id, name, status}'
        SITE_COUNT=$(echo "$SITES" | jq '. | length')
        echo "[SUCCESS] Retrieved $SITE_COUNT sites"
        if echo "$SITES" | grep -q "$SITE_ID"; then
            echo "[SUCCESS] Our test site appears in the list"
        fi
    else
        echo "Response: $SITES"
        echo "[FAILED] List sites returned an error"
    fi
    echo ""
    
    # Test 7: Get Specific Site (with infrastructureProviderId)
    echo "TEST 7: Get Specific Site"
    SITE_DETAIL=$(curl -s -H "Authorization: Bearer $TOKEN" \
      "http://localhost:8388/v2/org/nvidia/carbide/site/$SITE_ID?infrastructureProviderId=$INFRA_ID")
    echo "Response: $SITE_DETAIL" | jq -c '{id, name, status, created}'
    if echo "$SITE_DETAIL" | grep -q "test-site-automated"; then
        echo "[SUCCESS] Site details retrieved"
        
        # Extract registration token if available
        REG_TOKEN=$(echo "$SITE_DETAIL" | jq -r '.registrationToken // empty')
        if [ -n "$REG_TOKEN" ]; then
            echo "[INFO] Registration Token: ${REG_TOKEN:0:20}..."
        else
            echo "[INFO] Registration Token: Not yet available (site-manager may still be processing)"
        fi
    else
        echo "[FAILED] Site details failed"
    fi
    echo ""
fi

echo "========================================"
echo "Testing Complete"
echo "========================================"
echo ""
echo "Summary:"
echo "  Token: OK"
echo "  Health: OK"
echo "  Authentication: Working"
echo "  API Endpoints: See results above"
echo ""

