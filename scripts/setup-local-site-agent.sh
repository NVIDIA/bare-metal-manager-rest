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

# Setup script for local site-agent configuration

set -e

API_URL="http://localhost:8388"
KEYCLOAK_URL="http://localhost:8080"
ORG="test-org"
NAMESPACE="carbide"

echo "Setting up site-agent..."

# Wait for API
echo "Waiting for API..."
for i in {1..60}; do
    if curl -sf "$API_URL/healthz" > /dev/null 2>&1; then
        break
    fi
    if [ $i -eq 60 ]; then
        echo "ERROR: API not ready"
        exit 1
    fi
done

# Wait for Keycloak
echo "Waiting for Keycloak..."
for i in {1..30}; do
    if curl -sf "$KEYCLOAK_URL/realms/carbide-dev" > /dev/null 2>&1; then
        break
    fi
    if [ $i -eq 30 ]; then
        echo "ERROR: Keycloak not ready"
        exit 1
    fi
done

# Wait for site-manager to be ready (required for site creation)
echo "Waiting for site-manager..."
if ! kubectl -n $NAMESPACE wait --for=condition=ready pod -l app=carbide-rest-site-manager --timeout=180s; then
    echo "ERROR: Site-manager not ready after 180s. Checking pod status..."
    kubectl -n $NAMESPACE get pods -l app=carbide-rest-site-manager
    kubectl -n $NAMESPACE logs -l app=carbide-rest-site-manager --tail=20 2>/dev/null || true
    exit 1
fi

# Get access token
echo "Acquiring token..."
TOKEN=$(curl -sf -X POST "$KEYCLOAK_URL/realms/carbide-dev/protocol/openid-connect/token" \
    -H "Content-Type: application/x-www-form-urlencoded" \
    -d "client_id=carbide-api" \
    -d "client_secret=carbide-local-secret" \
    -d "grant_type=password" \
    -d "username=admin@example.com" \
    -d "password=adminpassword" | jq -r .access_token)

if [ -z "$TOKEN" ] || [ "$TOKEN" == "null" ]; then
    echo "ERROR: Failed to acquire token"
    exit 1
fi

# Ensure tenant exists (auto-created on first access)
echo "Ensuring tenant exists..."
curl -sf "$API_URL/v2/org/$ORG/carbide/tenant/current" \
    -H "Authorization: Bearer $TOKEN" > /dev/null 2>&1 || echo "WARN: Could not verify tenant"

# Get or create infrastructure provider
echo "Getting infrastructure provider..."
PROVIDER_RESP=$(curl -sf "$API_URL/v2/org/$ORG/carbide/infrastructure-provider/current" \
    -H "Authorization: Bearer $TOKEN" 2>/dev/null || echo "{}")

PROVIDER_ID=$(echo "$PROVIDER_RESP" | jq -r '.id // empty')
if [ -z "$PROVIDER_ID" ]; then
    echo "Creating infrastructure provider..."
    PROVIDER_RESP=$(curl -sf -X POST "$API_URL/v2/org/$ORG/carbide/infrastructure-provider" \
        -H "Authorization: Bearer $TOKEN" \
        -H "Content-Type: application/json" \
        -d '{"name": "Local Dev Provider", "description": "Local development infrastructure provider"}')
    PROVIDER_ID=$(echo "$PROVIDER_RESP" | jq -r '.id')
fi

# Check for existing site
echo "Checking for existing site..."
EXISTING_SITE=$(curl -sf "$API_URL/v2/org/$ORG/carbide/site?infrastructureProviderId=$PROVIDER_ID" \
    -H "Authorization: Bearer $TOKEN" | jq -r '.[] | select(.name == "local-dev-site") | .id' 2>/dev/null || echo "")

if [ -n "$EXISTING_SITE" ] && [ "$EXISTING_SITE" != "null" ]; then
    SITE_ID="$EXISTING_SITE"
    echo "Using existing site: $SITE_ID"
else
    echo "Creating site..."
    # Retry site creation a few times in case site-manager is still starting up
    for attempt in 1 2 3; do
        SITE_RESP=$(curl -s -X POST "$API_URL/v2/org/$ORG/carbide/site?infrastructureProviderId=$PROVIDER_ID" \
            -H "Authorization: Bearer $TOKEN" \
            -H "Content-Type: application/json" \
            -d '{
                "name": "local-dev-site",
                "description": "Local development site",
                "location": {"address": "Local Development", "city": "Santa Clara", "state": "CA", "country": "USA", "postalCode": "95054"},
                "contact": {"name": "Dev Team", "email": "dev@example.com", "phone": "555-0100"}
            }')
        
        SITE_ID=$(echo "$SITE_RESP" | jq -r '.id // empty')
        if [ -n "$SITE_ID" ] && [ "$SITE_ID" != "null" ]; then
            echo "Created site: $SITE_ID"
            break
        fi
        
        if [ $attempt -lt 3 ]; then
            echo "Site creation attempt $attempt failed, retrying in 5 seconds..."
            echo "Response: $SITE_RESP"
            read -t 5 < /dev/null || true  # Wait 5 seconds without sleep
        fi
    done
    
    if [ -z "$SITE_ID" ] || [ "$SITE_ID" == "null" ]; then
        echo "ERROR: Failed to create site after 3 attempts"
        echo "Last response: $SITE_RESP"
        exit 1
    fi
fi

# Create Temporal namespace
echo "Creating Temporal namespace..."
kubectl -n $NAMESPACE exec deployment/temporal -- temporal operator namespace create -n "$SITE_ID" 2>/dev/null || true

# Update site-agent configmap
echo "Updating site-agent configuration..."
kubectl -n $NAMESPACE get configmap carbide-rest-site-agent-config -o yaml | \
    sed "s/CLUSTER_ID: .*/CLUSTER_ID: \"$SITE_ID\"/" | \
    sed "s/TEMPORAL_SUBSCRIBE_NAMESPACE: .*/TEMPORAL_SUBSCRIBE_NAMESPACE: \"$SITE_ID\"/" | \
    sed "s/TEMPORAL_SUBSCRIBE_QUEUE: .*/TEMPORAL_SUBSCRIBE_QUEUE: \"$SITE_ID\"/" | \
    kubectl apply -f -

# Ensure TLS secret exists
if ! kubectl -n $NAMESPACE get secret carbide-tls-certs > /dev/null 2>&1; then
    echo "Creating TLS secret..."
    CERT_DIR=$(mktemp -d)
    openssl req -x509 -newkey rsa:2048 -keyout "$CERT_DIR/tls.key" -out "$CERT_DIR/ca.crt" \
        -days 365 -nodes -subj "/CN=carbide-local-ca" 2>/dev/null
    cp "$CERT_DIR/ca.crt" "$CERT_DIR/tls.crt"
    
    kubectl -n $NAMESPACE create secret generic carbide-tls-certs \
        --from-file=ca.crt="$CERT_DIR/ca.crt" \
        --from-file=tls.crt="$CERT_DIR/tls.crt" \
        --from-file=tls.key="$CERT_DIR/tls.key"
    rm -rf "$CERT_DIR"
fi

# Update site-registration secret
echo "Updating site-registration secret..."
kubectl -n $NAMESPACE get secret site-registration -o yaml 2>/dev/null | \
    sed "s/site-uuid: .*/site-uuid: $(echo -n $SITE_ID | base64)/" | \
    kubectl apply -f - 2>/dev/null || \
    kubectl -n $NAMESPACE create secret generic site-registration \
        --from-literal=site-uuid="$SITE_ID" \
        --from-literal=otp="local-dev-otp" \
        --from-literal=creds-url="http://carbide-rest-site-manager:8100/v1/site/credentials" \
        --from-literal=cacert=""

# Restart site-agent
echo "Restarting site-agent..."
kubectl -n $NAMESPACE rollout restart deployment/carbide-rest-site-agent
kubectl -n $NAMESPACE rollout status deployment/carbide-rest-site-agent --timeout=120s

# Verify
echo "Verifying..."
for i in {1..10}; do
    STATUS=$(kubectl -n $NAMESPACE get pods -l app=carbide-rest-site-agent -o jsonpath='{.items[0].status.phase}' 2>/dev/null)
    READY=$(kubectl -n $NAMESPACE get pods -l app=carbide-rest-site-agent -o jsonpath='{.items[0].status.containerStatuses[0].ready}' 2>/dev/null)
    if [ "$STATUS" == "Running" ] && [ "$READY" == "true" ]; then
        break
    fi
done

kubectl -n $NAMESPACE get pods -l app=carbide-rest-site-agent

echo ""
echo "Site-agent setup complete. Site ID: $SITE_ID"
