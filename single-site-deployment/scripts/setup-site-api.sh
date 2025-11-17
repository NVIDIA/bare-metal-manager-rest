#!/bin/bash
# setup-site-api.sh
# Complete workflow to create a Forge site via API and configure site agent

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/../.." && pwd)"

# Configuration
ORG_NAME="${ORG_NAME:-nvidia}"
SITE_NAME="${1:-test-site-$(date +%s)}"

# API endpoints
KEYCLOAK_URL="http://localhost:8080"
API_BASE_URL="http://localhost:8388/v2/org/$ORG_NAME/carbide"

log() {
    echo "[$(date +%H:%M:%S)] $1"
}

error_exit() {
    echo "[ERROR] $1"
    exit 1
}

# Cleanup function
cleanup() {
    log "Cleaning up port forwards..."
    [ ! -z "$KEYCLOAK_PF_PID" ] && kill $KEYCLOAK_PF_PID 2>/dev/null || true
    [ ! -z "$API_PF_PID" ] && kill $API_PF_PID 2>/dev/null || true
}
trap cleanup EXIT

echo "=========================================="
echo "Forge Site Setup via Cloud API"
echo "=========================================="
echo "Organization: $ORG_NAME"
echo "Site Name:    $SITE_NAME"
echo ""

# Check prerequisites
command -v jq >/dev/null 2>&1 || error_exit "jq is required but not installed"
command -v kubectl >/dev/null 2>&1 || error_exit "kubectl is required but not installed"

# Step 1: Setup port forwards
log "[1/9] Setting up port forwards..."
pkill -f "kubectl port-forward.*keycloak" 2>/dev/null || true
pkill -f "kubectl port-forward.*cloud-api" 2>/dev/null || true

kubectl port-forward -n carbide-system svc/keycloak 8080:8080 > /dev/null 2>&1 &
KEYCLOAK_PF_PID=$!

kubectl port-forward -n carbide-system svc/cloud-api 8388:8388 > /dev/null 2>&1 &
API_PF_PID=$!

sleep 3

# Step 2: Authenticate
log "[2/9] Authenticating with Keycloak..."
TOKEN=$(curl -sf -X POST "$KEYCLOAK_URL/realms/carbide/protocol/openid-connect/token" \
  -H 'Content-Type: application/x-www-form-urlencoded' \
  -d 'username=testuser' \
  -d 'password=testpass' \
  -d 'grant_type=password' \
  -d 'client_id=carbide-api' \
  -d 'client_secret=carbide-secret-dev-only-do-not-use-in-prod' | jq -r '.access_token')

[ -z "$TOKEN" ] || [ "$TOKEN" == "null" ] && error_exit "Failed to get authentication token"
echo "       Token: ${TOKEN:0:40}..."

# Step 3: Check API connectivity
log "[3/9] Verifying API connectivity..."
API_VERSION=$(curl -s -H "Authorization: Bearer $TOKEN" \
  "$API_BASE_URL/metadata" | jq -r '.version')
echo "       API Version: $API_VERSION"

USER_EMAIL=$(curl -s -H "Authorization: Bearer $TOKEN" \
  "$API_BASE_URL/user/current" | jq -r '.email')
echo "       User: $USER_EMAIL"

# Step 4: Get infrastructure provider
log "[4/9] Getting infrastructure provider..."
INFRA_RESPONSE=$(curl -s -H "Authorization: Bearer $TOKEN" \
  "$API_BASE_URL/infrastructure-provider/current")

INFRA_PROVIDER_ID=$(echo "$INFRA_RESPONSE" | jq -r '.id')
INFRA_ORG_DISPLAY=$(echo "$INFRA_RESPONSE" | jq -r '.orgDisplayName')

[ -z "$INFRA_PROVIDER_ID" ] || [ "$INFRA_PROVIDER_ID" == "null" ] && \
  error_exit "Failed to get infrastructure provider"

echo "       Provider ID: $INFRA_PROVIDER_ID"
echo "       Provider Org: $INFRA_ORG_DISPLAY"

# Step 5: Create site
log "[5/9] Creating site: $SITE_NAME"
SITE_RESPONSE=$(curl -s -X POST \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d "{
    \"name\": \"$SITE_NAME\",
    \"description\": \"Test site created via API on $(date)\",
    \"location\": {
      \"city\": \"Santa Clara\",
      \"state\": \"CA\",
      \"country\": \"USA\"
    },
    \"contact\": {
      \"email\": \"admin@${SITE_NAME}.local\"
    }
  }" \
  "$API_BASE_URL/site")

SITE_ID=$(echo "$SITE_RESPONSE" | jq -r '.id')
REGISTRATION_TOKEN=$(echo "$SITE_RESPONSE" | jq -r '.registrationToken')
TOKEN_EXPIRATION=$(echo "$SITE_RESPONSE" | jq -r '.registrationTokenExpiration')

if [ -z "$SITE_ID" ] || [ "$SITE_ID" == "null" ]; then
    echo "       ERROR: Site creation failed"
    echo "$SITE_RESPONSE" | jq .
    exit 1
fi

echo "       Site ID: $SITE_ID"
echo "       Status: $(echo "$SITE_RESPONSE" | jq -r '.status')"
echo "       Registration Token: ${REGISTRATION_TOKEN:0:30}..."
echo "       Token Expires: $TOKEN_EXPIRATION"

# Step 6: Configure site settings
log "[6/9] Configuring site settings..."
UPDATE_RESPONSE=$(curl -s -X PATCH \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d "{
    \"description\": \"Test site created via API - Updated\",
    \"serialConsoleHostname\": \"console.${SITE_NAME}.local\",
    \"isSerialConsoleEnabled\": true,
    \"serialConsoleIdleTimeout\": 300,
    \"serialConsoleMaxSessionLength\": 3600,
    \"location\": {
      \"city\": \"Santa Clara\",
      \"state\": \"CA\",
      \"country\": \"USA\"
    }
  }" \
  "$API_BASE_URL/site/$SITE_ID?infrastructureProviderId=$INFRA_PROVIDER_ID")

UPDATED_DESC=$(echo "$UPDATE_RESPONSE" | jq -r '.description')
echo "       Description: $UPDATED_DESC"
echo "       Serial Console: Enabled"

# Step 7: Get site details
log "[7/9] Retrieving site details..."
SITE_DETAILS=$(curl -s -H "Authorization: Bearer $TOKEN" \
  "$API_BASE_URL/site/$SITE_ID?infrastructureProviderId=$INFRA_PROVIDER_ID")

echo "       Site Details:"
echo "$SITE_DETAILS" | jq '{
  id,
  name,
  description,
  status,
  isOnline,
  infrastructureProviderId,
  serialConsoleHostname,
  isSerialConsoleEnabled,
  location,
  contact
}'

# Step 8: Update site agent bootstrap
log "[8/9] Updating site agent bootstrap configuration..."

# Check if CA cert exists
CA_CERT_PATH="$PROJECT_ROOT/build/certs/ca-cert.pem"
if [ -f "$CA_CERT_PATH" ]; then
    kubectl create configmap site-agent-bootstrap \
      -n carbide-site \
      --from-literal=site-uuid="$SITE_ID" \
      --from-literal=otp="$REGISTRATION_TOKEN" \
      --from-literal=creds-url="https://site-manager.carbide-system.svc.cluster.local:8100/v1/sitecreds" \
      --from-file=cacert="$CA_CERT_PATH" \
      --dry-run=client -o yaml | kubectl apply -f -
    echo "       ConfigMap updated with CA cert"
else
    kubectl create configmap site-agent-bootstrap \
      -n carbide-site \
      --from-literal=site-uuid="$SITE_ID" \
      --from-literal=otp="$REGISTRATION_TOKEN" \
      --from-literal=creds-url="https://site-manager.carbide-system.svc.cluster.local:8100/v1/sitecreds" \
      --from-literal=cacert="" \
      --dry-run=client -o yaml | kubectl apply -f -
    echo "       ConfigMap updated (no CA cert found)"
fi

# Step 9: Restart site agent
log "[9/9] Restarting site agent..."
kubectl rollout restart statefulset/carbide-site-agent -n carbide-site > /dev/null 2>&1
echo "       Site agent restarting (3 replicas)"

# Wait a moment for restart to begin
sleep 2

echo ""
echo "=========================================="
echo "Site Setup Complete!"
echo "=========================================="
echo ""
echo "Site Information:"
echo "  Name:       $SITE_NAME"
echo "  ID:         $SITE_ID"
echo "  Status:     Pending -> Will become 'Registered' when agent connects"
echo "  Provider:   $INFRA_PROVIDER_ID"
echo ""
echo "Export these variables for further API calls:"
echo ""
echo "  export TOKEN='$TOKEN'"
echo "  export INFRA_PROVIDER_ID='$INFRA_PROVIDER_ID'"
echo "  export SITE_ID='$SITE_ID'"
echo "  export SITE_NAME='$SITE_NAME'"
echo "  export API_BASE_URL='$API_BASE_URL'"
echo ""
echo "Verification Commands:"
echo ""
echo "  # Check site status (isOnline should become true)"
echo "  curl -s -H \"Authorization: Bearer \$TOKEN\" \\"
echo "    \"\$API_BASE_URL/site/\$SITE_ID?infrastructureProviderId=\$INFRA_PROVIDER_ID\" | jq '.status, .isOnline'"
echo ""
echo "  # Watch site agent logs"
echo "  kubectl logs -f carbide-site-agent-0 -n carbide-site"
echo ""
echo "  # Get site status history"
echo "  curl -s -H \"Authorization: Bearer \$TOKEN\" \\"
echo "    \"\$API_BASE_URL/site/\$SITE_ID/status-history?infrastructureProviderId=\$INFRA_PROVIDER_ID\" | jq ."
echo ""
echo "  # List all sites"
echo "  curl -s -H \"Authorization: Bearer \$TOKEN\" \\"
echo "    \"\$API_BASE_URL/site?infrastructureProviderId=\$INFRA_PROVIDER_ID\" | jq ."
echo ""
echo "  # List machines (after agent connects and discovers hardware)"
echo "  curl -s -H \"Authorization: Bearer \$TOKEN\" \\"
echo "    \"\$API_BASE_URL/machine?siteId=\$SITE_ID&infrastructureProviderId=\$INFRA_PROVIDER_ID\" | jq ."
echo ""
echo "Next Steps:"
echo "  1. Wait for site agent to connect (watch logs)"
echo "  2. Verify site status changes to 'Registered' and isOnline=true"
echo "  3. Check machine discovery via /machine endpoint"
echo "  4. Create instance types for resource allocation"
echo "  5. Set up networking (IP blocks, VPCs)"
echo ""

