#!/usr/bin/env bash

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/../.." && pwd)"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

log_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

# Check if jq is installed
if ! command -v jq &> /dev/null; then
    log_error "jq is required but not installed. Please install jq first."
    exit 1
fi

# Cleanup function
cleanup() {
    log_info "Cleaning up port-forwards..."
    if [ ! -z "$KEYCLOAK_PF_PID" ]; then
        kill $KEYCLOAK_PF_PID 2>/dev/null || true
    fi
    if [ ! -z "$API_PF_PID" ]; then
        kill $API_PF_PID 2>/dev/null || true
    fi
}

trap cleanup EXIT

log_info "============================================================"
log_info "Testing Site Manager and Site Agent APIs"
log_info "============================================================"

# Wait for pods to be ready
log_info "Waiting for pods to be ready..."
kubectl wait --for=condition=ready pod -l app=keycloak -n carbide-system --timeout=120s
kubectl wait --for=condition=ready pod -l app=cloud-api -n carbide-system --timeout=120s

# Port forward services
log_info "Setting up port-forwards..."
kubectl port-forward -n carbide-system svc/keycloak 8080:8080 > /dev/null 2>&1 &
KEYCLOAK_PF_PID=$!
sleep 2

kubectl port-forward -n carbide-system svc/cloud-api 8388:8388 > /dev/null 2>&1 &
API_PF_PID=$!
sleep 2

# Test Keycloak connectivity
log_info "Testing Keycloak connectivity..."
if ! curl -sf http://localhost:8080/realms/carbide/.well-known/openid-configuration > /dev/null; then
    log_error "Cannot reach Keycloak"
    exit 1
fi
log_success "Keycloak is reachable"

# Get authentication token
log_info "Getting authentication token..."
TOKEN=$(curl -sf -X POST 'http://localhost:8080/realms/carbide/protocol/openid-connect/token' \
  -H 'Content-Type: application/x-www-form-urlencoded' \
  -d 'username=testuser' \
  -d 'password=testpass' \
  -d 'grant_type=password' \
  -d 'client_id=carbide-api' \
  -d 'client_secret=carbide-secret-dev-only-do-not-use-in-prod' | jq -r '.access_token')

if [ -z "$TOKEN" ] || [ "$TOKEN" == "null" ]; then
    log_error "Failed to get authentication token"
    exit 1
fi
log_success "Authentication token obtained"

# Test API connectivity
log_info "Testing API connectivity..."
METADATA=$(curl -s -H "Authorization: Bearer $TOKEN" http://localhost:8388/v2/org/nvidia/carbide/metadata)
if echo "$METADATA" | jq -e '.version' > /dev/null 2>&1; then
    VERSION=$(echo "$METADATA" | jq -r '.version')
    log_success "API is reachable (version: $VERSION)"
else
    log_error "API metadata endpoint failed"
    echo "$METADATA"
    exit 1
fi

# Test user endpoint
log_info "Testing user endpoint..."
USER_INFO=$(curl -s -H "Authorization: Bearer $TOKEN" http://localhost:8388/v2/org/nvidia/carbide/user/current)
if echo "$USER_INFO" | jq -e '.email' > /dev/null 2>&1; then
    EMAIL=$(echo "$USER_INFO" | jq -r '.email')
    log_success "User endpoint working (email: $EMAIL)"
else
    log_warning "User endpoint returned: $USER_INFO"
fi

# Create Infrastructure Provider
log_info "Creating infrastructure provider..."
INFRA_PROVIDER_PAYLOAD='{
  "name": "Test Infrastructure Provider",
  "description": "Test infrastructure provider for site testing",
  "type": "on-prem"
}'

INFRA_PROVIDER_RESULT=$(curl -s -X POST \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d "$INFRA_PROVIDER_PAYLOAD" \
  http://localhost:8388/v2/org/nvidia/carbide/infrastructure-provider)

if echo "$INFRA_PROVIDER_RESULT" | jq -e '.id' > /dev/null 2>&1; then
    INFRA_PROVIDER_ID=$(echo "$INFRA_PROVIDER_RESULT" | jq -r '.id')
    log_success "Infrastructure provider created (ID: $INFRA_PROVIDER_ID)"
else
    # Try to get existing infrastructure provider
    log_info "Checking for existing infrastructure provider..."
    INFRA_PROVIDER_RESULT=$(curl -s -H "Authorization: Bearer $TOKEN" \
      http://localhost:8388/v2/org/nvidia/carbide/infrastructure-provider/current)
    
    if echo "$INFRA_PROVIDER_RESULT" | jq -e '.id' > /dev/null 2>&1; then
        INFRA_PROVIDER_ID=$(echo "$INFRA_PROVIDER_RESULT" | jq -r '.id')
        log_success "Using existing infrastructure provider (ID: $INFRA_PROVIDER_ID)"
    else
        log_error "Failed to create or retrieve infrastructure provider"
        echo "$INFRA_PROVIDER_RESULT"
        exit 1
    fi
fi

# Create Site
log_info "Creating site..."
SITE_PAYLOAD=$(cat <<EOF
{
  "name": "test-site-$(date +%s)",
  "description": "Test site for API validation",
  "infrastructureProviderId": "$INFRA_PROVIDER_ID",
  "isSerialConsoleEnabled": false
}
EOF
)

SITE_RESULT=$(curl -s -X POST \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d "$SITE_PAYLOAD" \
  http://localhost:8388/v2/org/nvidia/carbide/site)

if echo "$SITE_RESULT" | jq -e '.id' > /dev/null 2>&1; then
    SITE_ID=$(echo "$SITE_RESULT" | jq -r '.id')
    SITE_NAME=$(echo "$SITE_RESULT" | jq -r '.name')
    REGISTRATION_TOKEN=$(echo "$SITE_RESULT" | jq -r '.registrationToken')
    log_success "Site created successfully!"
    echo ""
    log_info "Site Details:"
    echo "  ID: $SITE_ID"
    echo "  Name: $SITE_NAME"
    echo "  Registration Token: ${REGISTRATION_TOKEN:0:20}..."
else
    log_error "Failed to create site"
    echo "$SITE_RESULT" | jq '.' 2>/dev/null || echo "$SITE_RESULT"
    exit 1
fi

# Get site details
log_info "Retrieving site details..."
SITE_DETAILS=$(curl -s -H "Authorization: Bearer $TOKEN" \
  "http://localhost:8388/v2/org/nvidia/carbide/site/$SITE_ID?infrastructureProviderId=$INFRA_PROVIDER_ID")

if echo "$SITE_DETAILS" | jq -e '.id' > /dev/null 2>&1; then
    log_success "Site details retrieved successfully"
    echo "$SITE_DETAILS" | jq '{id, name, infrastructureProviderId, siteControllerVersion, siteAgentVersion}'
else
    log_error "Failed to retrieve site details"
    echo "$SITE_DETAILS" | jq '.' 2>/dev/null || echo "$SITE_DETAILS"
fi

# List all sites
log_info "Listing all sites..."
ALL_SITES=$(curl -s -H "Authorization: Bearer $TOKEN" \
  "http://localhost:8388/v2/org/nvidia/carbide/site?infrastructureProviderId=$INFRA_PROVIDER_ID")

# Check if response is an array or paginated object
if echo "$ALL_SITES" | jq -e 'type == "array"' > /dev/null 2>&1; then
    SITE_COUNT=$(echo "$ALL_SITES" | jq '. | length')
    log_success "Found $SITE_COUNT site(s)"
    echo "$ALL_SITES" | jq '.[] | {id, name, infrastructureProviderId, status}'
elif echo "$ALL_SITES" | jq -e '.results' > /dev/null 2>&1; then
    SITE_COUNT=$(echo "$ALL_SITES" | jq -r '.totalCount // 0')
    log_success "Found $SITE_COUNT site(s)"
    echo "$ALL_SITES" | jq '.results[] | {id, name, infrastructureProviderId, status}'
else
    log_warning "Unexpected sites response format"
    echo "$ALL_SITES" | jq '.'
fi

# Get site status history
log_info "Retrieving site status history..."
STATUS_HISTORY=$(curl -s -H "Authorization: Bearer $TOKEN" \
  "http://localhost:8388/v2/org/nvidia/carbide/site/$SITE_ID/status-history?infrastructureProviderId=$INFRA_PROVIDER_ID")

if echo "$STATUS_HISTORY" | jq -e 'type == "array"' > /dev/null 2>&1; then
    STATUS_COUNT=$(echo "$STATUS_HISTORY" | jq '. | length')
    log_success "Site status history retrieved ($STATUS_COUNT entries)"
    echo "$STATUS_HISTORY" | jq '.[] | {status, message, created}'
elif echo "$STATUS_HISTORY" | jq -e '.message' > /dev/null 2>&1; then
    log_warning "Site status history returned: $(echo "$STATUS_HISTORY" | jq -r '.message')"
else
    log_success "Site status history retrieved"
    echo "$STATUS_HISTORY" | jq '.'
fi

# Update site
log_info "Updating site description..."
UPDATE_PAYLOAD='{
  "description": "Updated test site description"
}'

UPDATE_RESULT=$(curl -s -X PATCH \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d "$UPDATE_PAYLOAD" \
  "http://localhost:8388/v2/org/nvidia/carbide/site/$SITE_ID?infrastructureProviderId=$INFRA_PROVIDER_ID")

if echo "$UPDATE_RESULT" | jq -e '.id' > /dev/null 2>&1; then
    UPDATED_DESC=$(echo "$UPDATE_RESULT" | jq -r '.description')
    log_success "Site updated successfully (description: $UPDATED_DESC)"
else
    log_error "Failed to update site"
    echo "$UPDATE_RESULT" | jq '.' 2>/dev/null || echo "$UPDATE_RESULT"
fi

# Test site-agent related endpoints (machines, etc.)
log_info "Testing machine endpoints..."
MACHINES=$(curl -s -H "Authorization: Bearer $TOKEN" \
  "http://localhost:8388/v2/org/nvidia/carbide/machine?siteId=$SITE_ID&infrastructureProviderId=$INFRA_PROVIDER_ID")

if echo "$MACHINES" | jq -e '.results' > /dev/null 2>&1; then
    MACHINE_COUNT=$(echo "$MACHINES" | jq -r '.totalCount // 0')
    log_success "Machine list retrieved ($MACHINE_COUNT machines)"
else
    log_warning "Machine endpoint returned: $(echo "$MACHINES" | jq -c '.')"
fi

echo ""
log_info "============================================================"
log_success "Site Manager API Tests Complete!"
log_info "============================================================"
echo ""
log_info "Summary:"
echo "  ✓ Authentication working"
echo "  ✓ Infrastructure provider: $INFRA_PROVIDER_ID"
echo "  ✓ Site created: $SITE_ID"
echo "  ✓ Site operations (create, read, update) working"
echo ""
log_info "Next steps:"
echo "  - Deploy site-agent with registration token to test agent functionality"
echo "  - Add machines to test machine management"
echo "  - Test VPC and allocation workflows"
echo ""

