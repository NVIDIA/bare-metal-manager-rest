#!/bin/bash

# Create and Register Test Site
# This script creates a site via cloud-api and updates site-agent with the registration details

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ROOT_DIR="$(cd "$SCRIPT_DIR/../.." && pwd)"

# Configuration
SITE_NAME="${SITE_NAME:-test-site-1}"
ORG_NAME="${ORG_NAME:-nvidia}"

echo "========================================"
echo "Creating Test Site"
echo "========================================"
echo ""
echo "Site Name: $SITE_NAME"
echo "Organization: $ORG_NAME"
echo ""

# Start port forwards
echo "Starting port forwards..."
kubectl port-forward -n carbide-system svc/keycloak 8080:8080 > /dev/null 2>&1 &
PF_KC_PID=$!
kubectl port-forward -n carbide-system svc/cloud-api 8388:8388 > /dev/null 2>&1 &
PF_API_PID=$!
kubectl port-forward -n carbide-system svc/site-manager 8100:8100 > /dev/null 2>&1 &
PF_SM_PID=$!

# Cleanup function
cleanup() {
    echo ""
    echo "Cleaning up port forwards..."
    kill $PF_KC_PID $PF_API_PID $PF_SM_PID 2>/dev/null || true
}
trap cleanup EXIT

sleep 3

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
echo "[SUCCESS] Authentication token obtained"

# Create site
echo ""
echo "Creating site via cloud-api..."
SITE_RESPONSE=$(curl -s -X POST -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  http://localhost:8388/v2/org/$ORG_NAME/carbide/site \
  -d "{
    \"name\": \"$SITE_NAME\",
    \"description\": \"Test site created by automation script\",
    \"serialConsoleHostname\": \"console.${SITE_NAME}.local\",
    \"location\": {
      \"city\": \"Santa Clara\",
      \"state\": \"CA\",
      \"country\": \"US\"
    },
    \"contact\": {
      \"email\": \"admin@${SITE_NAME}.local\"
    }
  }")

SITE_ID=$(echo "$SITE_RESPONSE" | jq -r '.id // empty')

if [ -z "$SITE_ID" ]; then
    echo "[ERROR] Failed to create site"
    echo "Response: $SITE_RESPONSE"
    exit 1
fi

echo "[SUCCESS] Site created with ID: $SITE_ID"

# Get site details from site-manager (includes OTP)
echo ""
echo "Retrieving site OTP from site-manager..."
sleep 5  # Give site-manager time to process

SITE_DETAILS=$(curl -sk https://localhost:8100/v1/site/$SITE_ID)
SITE_OTP=$(echo "$SITE_DETAILS" | jq -r '.otp // empty')

if [ -z "$SITE_OTP" ]; then
    echo "[ERROR] Failed to retrieve site OTP"
    echo "Response: $SITE_DETAILS"
    exit 1
fi

echo "[SUCCESS] Site OTP retrieved: ${SITE_OTP:0:16}..."

# Update site-agent bootstrap ConfigMap
echo ""
echo "Updating site-agent bootstrap ConfigMap..."

kubectl create configmap site-agent-bootstrap -n carbide-site \
  --from-literal=site-uuid="$SITE_ID" \
  --from-literal=otp="$SITE_OTP" \
  --from-literal=creds-url="https://site-manager.carbide-system.svc.cluster.local:8100/v1/sitecreds" \
  --from-file=cacert="$ROOT_DIR/build/certs/ca-cert.pem" \
  --dry-run=client -o yaml | kubectl apply -f -

echo "[SUCCESS] Site-agent bootstrap ConfigMap updated"

# Restart site-agent
echo ""
echo "Restarting site-agent..."
kubectl rollout restart statefulset/carbide-site-agent -n carbide-site

echo ""
echo "========================================"
echo "Site Creation Complete!"
echo "========================================"
echo ""
echo "Site Details:"
echo "  Site ID:    $SITE_ID"
echo "  Site Name:  $SITE_NAME"
echo "  Site OTP:   ${SITE_OTP:0:16}... (hidden)"
echo ""
echo "Site-agent will now bootstrap with this site."
echo ""
echo "Watch site-agent logs:"
echo "  kubectl logs -f carbide-site-agent-0 -n carbide-site"
echo ""
echo "Check site status:"
echo "  curl -s -H \"Authorization: Bearer \$TOKEN\" \\"
echo "    http://localhost:8388/v2/org/$ORG_NAME/carbide/site/$SITE_ID | jq ."
echo ""

