#!/bin/bash
# demo-site-api.sh
# Quick demo of Forge Cloud API operations
# Shows how to interact with the site you just created

set -e

echo "=========================================="
echo "Forge Cloud API Demo"
echo "=========================================="
echo ""

# Configuration
ORG_NAME="nvidia"
API_BASE="http://localhost:8388/v2/org/$ORG_NAME/carbide"
KEYCLOAK_URL="http://localhost:8080"

# Check prerequisites
if ! command -v jq &> /dev/null; then
    echo "ERROR: jq is required. Install with: brew install jq"
    exit 1
fi

# Setup port forwards (kill any existing ones first)
echo "[1/6] Setting up port forwards..."
pkill -f "kubectl port-forward.*keycloak" 2>/dev/null || true
pkill -f "kubectl port-forward.*cloud-api" 2>/dev/null || true
sleep 1

kubectl port-forward -n carbide-system svc/keycloak 8080:8080 > /dev/null 2>&1 &
KEYCLOAK_PID=$!

kubectl port-forward -n carbide-system svc/cloud-api 8388:8388 > /dev/null 2>&1 &
API_PID=$!

sleep 3
echo "   Port forwards active (PIDs: $KEYCLOAK_PID, $API_PID)"

# Cleanup on exit
cleanup() {
    echo ""
    echo "Cleaning up..."
    kill $KEYCLOAK_PID $API_PID 2>/dev/null || true
}
trap cleanup EXIT

# Authenticate
echo "[2/6] Authenticating..."
TOKEN=$(curl -sf -X POST "$KEYCLOAK_URL/realms/carbide/protocol/openid-connect/token" \
  -H 'Content-Type: application/x-www-form-urlencoded' \
  -d 'username=testuser' \
  -d 'password=testpass' \
  -d 'grant_type=password' \
  -d 'client_id=carbide-api' \
  -d 'client_secret=carbide-secret-dev-only-do-not-use-in-prod' | jq -r '.access_token')

if [ -z "$TOKEN" ] || [ "$TOKEN" == "null" ]; then
    echo "ERROR: Authentication failed"
    exit 1
fi

echo "   Token obtained: ${TOKEN:0:40}..."

# Get infrastructure provider
echo "[3/6] Getting infrastructure provider..."
INFRA_RESPONSE=$(curl -s -H "Authorization: Bearer $TOKEN" \
  "$API_BASE/infrastructure-provider/current")

INFRA_PROVIDER_ID=$(echo "$INFRA_RESPONSE" | jq -r '.id')
INFRA_ORG=$(echo "$INFRA_RESPONSE" | jq -r '.orgDisplayName')

echo "   Provider ID: $INFRA_PROVIDER_ID"
echo "   Org: $INFRA_ORG"

# List all sites
echo "[4/6] Listing all sites..."
SITES_RESPONSE=$(curl -s -H "Authorization: Bearer $TOKEN" \
  "$API_BASE/site?infrastructureProviderId=$INFRA_PROVIDER_ID")

SITE_COUNT=$(echo "$SITES_RESPONSE" | jq '. | length')
echo "   Found $SITE_COUNT site(s)"
echo ""

if [ "$SITE_COUNT" -gt 0 ]; then
    echo "Sites:"
    echo "$SITES_RESPONSE" | jq -r '.[] | "  - \(.name) (\(.id)): \(.status) | Online: \(.isOnline)"'
    echo ""
    
    # Get first site for demo
    DEMO_SITE_ID=$(echo "$SITES_RESPONSE" | jq -r '.[0].id')
    DEMO_SITE_NAME=$(echo "$SITES_RESPONSE" | jq -r '.[0].name')
    
    echo "[5/6] Getting details for site: $DEMO_SITE_NAME"
    SITE_DETAILS=$(curl -s -H "Authorization: Bearer $TOKEN" \
      "$API_BASE/site/$DEMO_SITE_ID?infrastructureProviderId=$INFRA_PROVIDER_ID")
    
    echo "$SITE_DETAILS" | jq '{
      id,
      name,
      description,
      status,
      isOnline,
      siteControllerVersion,
      siteAgentVersion,
      location,
      contact
    }'
    
    echo ""
    echo "[6/6] Getting site status history..."
    STATUS_HISTORY=$(curl -s -H "Authorization: Bearer $TOKEN" \
      "$API_BASE/site/$DEMO_SITE_ID/status-history?infrastructureProviderId=$INFRA_PROVIDER_ID")
    
    echo "$STATUS_HISTORY" | jq '.[] | {status, message, created}'
    
    echo ""
    echo "=========================================="
    echo "Demo Complete!"
    echo "=========================================="
    echo ""
    echo "Your site: $DEMO_SITE_NAME ($DEMO_SITE_ID)"
    echo "Status: $(echo "$SITE_DETAILS" | jq -r '.status')"
    echo "Online: $(echo "$SITE_DETAILS" | jq -r '.isOnline')"
    echo ""
    echo "Useful commands:"
    echo ""
    echo "# Export for more API calls"
    echo "export TOKEN='$TOKEN'"
    echo "export INFRA_PROVIDER_ID='$INFRA_PROVIDER_ID'"
    echo "export SITE_ID='$DEMO_SITE_ID'"
    echo ""
    echo "# Watch site agent logs"
    echo "kubectl logs -f carbide-site-agent-0 -n carbide-site"
    echo ""
    echo "# List machines at site"
    echo "curl -s -H \"Authorization: Bearer \$TOKEN\" \\"
    echo "  \"$API_BASE/machine?siteId=\$SITE_ID&infrastructureProviderId=\$INFRA_PROVIDER_ID\" | jq ."
    echo ""
    echo "# Get site details again"
    echo "curl -s -H \"Authorization: Bearer \$TOKEN\" \\"
    echo "  \"$API_BASE/site/\$SITE_ID?infrastructureProviderId=\$INFRA_PROVIDER_ID\" | jq ."
    echo ""
else
    echo "No sites found. Create one with:"
    echo "  ./scripts/setup-site-api.sh my-site-name"
fi

echo ""

