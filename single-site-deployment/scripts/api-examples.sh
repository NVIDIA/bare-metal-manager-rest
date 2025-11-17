#!/bin/bash
# api-examples.sh
# Collection of example API calls for Forge Cloud API
# Source this file to load functions, then call them as needed
#
# Usage:
#   source scripts/api-examples.sh
#   setup_auth
#   create_site "my-site"
#   list_sites
#   get_site_details $SITE_ID

ORG_NAME="${ORG_NAME:-nvidia}"
KEYCLOAK_URL="http://localhost:8080"
API_BASE="http://localhost:8388/v2/org/$ORG_NAME/carbide"

# Authentication
setup_auth() {
    echo "Getting authentication token..."
    TOKEN=$(curl -sf -X POST "$KEYCLOAK_URL/realms/carbide/protocol/openid-connect/token" \
      -H 'Content-Type: application/x-www-form-urlencoded' \
      -d 'username=testuser' \
      -d 'password=testpass' \
      -d 'grant_type=password' \
      -d 'client_id=carbide-api' \
      -d 'client_secret=carbide-secret-dev-only-do-not-use-in-prod' | jq -r '.access_token')
    
    if [ -z "$TOKEN" ] || [ "$TOKEN" == "null" ]; then
        echo "ERROR: Authentication failed"
        return 1
    fi
    export TOKEN
    echo "Authentication successful"
    
    # Get infrastructure provider
    INFRA_RESPONSE=$(curl -s -H "Authorization: Bearer $TOKEN" \
      "$API_BASE/infrastructure-provider/current")
    INFRA_PROVIDER_ID=$(echo "$INFRA_RESPONSE" | jq -r '.id')
    export INFRA_PROVIDER_ID
    echo "Infrastructure Provider ID: $INFRA_PROVIDER_ID"
}

# API Info
get_api_version() {
    curl -s -H "Authorization: Bearer $TOKEN" \
      "$API_BASE/metadata" | jq .
}

get_current_user() {
    curl -s -H "Authorization: Bearer $TOKEN" \
      "$API_BASE/user/current" | jq .
}

# Infrastructure Provider
get_infrastructure_provider() {
    curl -s -H "Authorization: Bearer $TOKEN" \
      "$API_BASE/infrastructure-provider/current" | jq .
}

get_infrastructure_provider_stats() {
    curl -s -H "Authorization: Bearer $TOKEN" \
      "$API_BASE/infrastructure-provider/current/stats" | jq .
}

# Site Operations
create_site() {
    local site_name="${1:-test-site-$(date +%s)}"
    echo "Creating site: $site_name"
    
    SITE_RESPONSE=$(curl -s -X POST \
      -H "Authorization: Bearer $TOKEN" \
      -H "Content-Type: application/json" \
      -d "{
        \"name\": \"$site_name\",
        \"description\": \"Test site created via API\",
        \"location\": {
          \"city\": \"Santa Clara\",
          \"state\": \"CA\",
          \"country\": \"USA\"
        },
        \"contact\": {
          \"email\": \"admin@${site_name}.local\"
        }
      }" \
      "$API_BASE/site")
    
    SITE_ID=$(echo "$SITE_RESPONSE" | jq -r '.id')
    REGISTRATION_TOKEN=$(echo "$SITE_RESPONSE" | jq -r '.registrationToken')
    export SITE_ID
    export REGISTRATION_TOKEN
    
    echo "Site ID: $SITE_ID"
    echo "Registration Token: ${REGISTRATION_TOKEN:0:30}..."
    echo "$SITE_RESPONSE" | jq .
}

list_sites() {
    echo "Listing all sites..."
    curl -s -H "Authorization: Bearer $TOKEN" \
      "$API_BASE/site?infrastructureProviderId=$INFRA_PROVIDER_ID" | jq .
}

get_site_details() {
    local site_id="${1:-$SITE_ID}"
    echo "Getting site details: $site_id"
    curl -s -H "Authorization: Bearer $TOKEN" \
      "$API_BASE/site/$site_id?infrastructureProviderId=$INFRA_PROVIDER_ID" | jq .
}

update_site() {
    local site_id="${1:-$SITE_ID}"
    echo "Updating site: $site_id"
    curl -s -X PATCH \
      -H "Authorization: Bearer $TOKEN" \
      -H "Content-Type: application/json" \
      -d '{
        "description": "Updated via API",
        "serialConsoleHostname": "console.example.local",
        "isSerialConsoleEnabled": true,
        "serialConsoleIdleTimeout": 300,
        "serialConsoleMaxSessionLength": 3600
      }' \
      "$API_BASE/site/$site_id?infrastructureProviderId=$INFRA_PROVIDER_ID" | jq .
}

delete_site() {
    local site_id="${1:-$SITE_ID}"
    echo "Deleting site: $site_id"
    curl -s -X DELETE \
      -H "Authorization: Bearer $TOKEN" \
      "$API_BASE/site/$site_id?infrastructureProviderId=$INFRA_PROVIDER_ID"
}

get_site_status_history() {
    local site_id="${1:-$SITE_ID}"
    echo "Getting site status history: $site_id"
    curl -s -H "Authorization: Bearer $TOKEN" \
      "$API_BASE/site/$site_id/status-history?infrastructureProviderId=$INFRA_PROVIDER_ID" | jq .
}

# Machine Operations
list_machines() {
    local site_id="${1:-$SITE_ID}"
    echo "Listing machines for site: $site_id"
    curl -s -H "Authorization: Bearer $TOKEN" \
      "$API_BASE/machine?siteId=$site_id&infrastructureProviderId=$INFRA_PROVIDER_ID" | jq .
}

get_machine() {
    local machine_id="$1"
    echo "Getting machine details: $machine_id"
    curl -s -H "Authorization: Bearer $TOKEN" \
      "$API_BASE/machine/$machine_id?infrastructureProviderId=$INFRA_PROVIDER_ID&includeMetadata=true" | jq .
}

list_machine_capabilities() {
    local site_id="${1:-$SITE_ID}"
    echo "Listing machine capabilities for site: $site_id"
    curl -s -H "Authorization: Bearer $TOKEN" \
      "$API_BASE/machine-capability?siteId=$site_id" | jq .
}

# Expected Machine Operations
create_expected_machine() {
    local site_id="${1:-$SITE_ID}"
    echo "Creating expected machine for site: $site_id"
    curl -s -X POST \
      -H "Authorization: Bearer $TOKEN" \
      -H "Content-Type: application/json" \
      -d "{
        \"siteId\": \"$site_id\",
        \"bmcMacAddress\": \"00:1A:2B:3C:4D:5E\",
        \"bmcUsername\": \"admin\",
        \"bmcPassword\": \"admin123\",
        \"chassisSerialNumber\": \"CHASSIS-001\",
        \"labels\": {
          \"rack\": \"A1\",
          \"row\": \"1\"
        }
      }" \
      "$API_BASE/expected-machine" | jq .
}

list_expected_machines() {
    local site_id="${1:-$SITE_ID}"
    echo "Listing expected machines for site: $site_id"
    curl -s -H "Authorization: Bearer $TOKEN" \
      "$API_BASE/expected-machine?siteId=$site_id" | jq .
}

# IP Block Operations
create_ip_block() {
    local site_id="${1:-$SITE_ID}"
    echo "Creating IP block for site: $site_id"
    curl -s -X POST \
      -H "Authorization: Bearer $TOKEN" \
      -H "Content-Type: application/json" \
      -d "{
        \"name\": \"Primary Network Block\",
        \"description\": \"Main IP block for site\",
        \"siteId\": \"$site_id\",
        \"routingType\": \"DatacenterOnly\",
        \"prefix\": \"10.0.0.0\",
        \"prefixLength\": 16,
        \"protocolVersion\": \"IPv4\"
      }" \
      "$API_BASE/ipblock" | jq .
}

list_ip_blocks() {
    echo "Listing IP blocks..."
    curl -s -H "Authorization: Bearer $TOKEN" \
      "$API_BASE/ipblock?infrastructureProviderId=$INFRA_PROVIDER_ID" | jq .
}

# Instance Type Operations
create_instance_type() {
    local site_id="${1:-$SITE_ID}"
    echo "Creating instance type for site: $site_id"
    curl -s -X POST \
      -H "Authorization: Bearer $TOKEN" \
      -H "Content-Type: application/json" \
      -d "{
        \"name\": \"standard.large\",
        \"description\": \"Standard large instance\",
        \"siteId\": \"$site_id\",
        \"machineCapabilities\": [
          {
            \"type\": \"CPU\",
            \"name\": \"Intel Xeon\",
            \"count\": 2
          },
          {
            \"type\": \"Memory\",
            \"name\": \"DDR4\",
            \"capacity\": \"32GB\",
            \"count\": 4
          }
        ]
      }" \
      "$API_BASE/instance/type" | jq .
}

list_instance_types() {
    local site_id="${1:-$SITE_ID}"
    echo "Listing instance types for site: $site_id"
    curl -s -H "Authorization: Bearer $TOKEN" \
      "$API_BASE/instance/type?siteId=$site_id&infrastructureProviderId=$INFRA_PROVIDER_ID" | jq .
}

# Tenant Operations (when acting as tenant)
get_tenant() {
    curl -s -H "Authorization: Bearer $TOKEN" \
      "$API_BASE/tenant/current" | jq .
}

get_tenant_stats() {
    curl -s -H "Authorization: Bearer $TOKEN" \
      "$API_BASE/tenant/current/stats" | jq .
}

# Allocation Operations
create_allocation() {
    local site_id="${1:-$SITE_ID}"
    local tenant_id="$2"
    local instance_type_id="$3"
    
    echo "Creating allocation for site: $site_id"
    curl -s -X POST \
      -H "Authorization: Bearer $TOKEN" \
      -H "Content-Type: application/json" \
      -d "{
        \"name\": \"Test Allocation\",
        \"description\": \"Test resource allocation\",
        \"tenantId\": \"$tenant_id\",
        \"siteId\": \"$site_id\",
        \"constraints\": [
          {
            \"resourceType\": \"InstanceType\",
            \"resourceTypeId\": \"$instance_type_id\",
            \"constraintType\": \"Reserved\",
            \"constraintValue\": 10
          }
        ]
      }" \
      "$API_BASE/allocation" | jq .
}

list_allocations() {
    echo "Listing allocations..."
    curl -s -H "Authorization: Bearer $TOKEN" \
      "$API_BASE/allocation?infrastructureProviderId=$INFRA_PROVIDER_ID" | jq .
}

# Audit Operations
list_audit_logs() {
    echo "Listing audit logs..."
    curl -s -H "Authorization: Bearer $TOKEN" \
      "$API_BASE/audit?pageNumber=1&pageSize=20&orderBy=TIMESTAMP_DESC" | jq .
}

list_failed_audit_logs() {
    echo "Listing failed audit logs..."
    curl -s -H "Authorization: Bearer $TOKEN" \
      "$API_BASE/audit?failedOnly=true&pageNumber=1&pageSize=20" | jq .
}

# Helper function to configure site agent
configure_site_agent() {
    local site_id="${1:-$SITE_ID}"
    local reg_token="${2:-$REGISTRATION_TOKEN}"
    
    echo "Configuring site agent for site: $site_id"
    
    kubectl create configmap site-agent-bootstrap \
      -n carbide-site \
      --from-literal=site-uuid="$site_id" \
      --from-literal=otp="$reg_token" \
      --from-literal=creds-url="https://site-manager.carbide-system.svc.cluster.local:8100/v1/sitecreds" \
      --dry-run=client -o yaml | kubectl apply -f -
    
    kubectl rollout restart statefulset/carbide-site-agent -n carbide-site
    
    echo "Site agent restarted with new configuration"
    echo ""
    echo "Watch logs with:"
    echo "  kubectl logs -f carbide-site-agent-0 -n carbide-site"
}

# Print available functions
show_functions() {
    echo ""
    echo "Available API Functions:"
    echo ""
    echo "Setup & Auth:"
    echo "  setup_auth                              - Authenticate and get provider ID"
    echo ""
    echo "API Info:"
    echo "  get_api_version                         - Get API version"
    echo "  get_current_user                        - Get current user details"
    echo ""
    echo "Infrastructure Provider:"
    echo "  get_infrastructure_provider             - Get provider details"
    echo "  get_infrastructure_provider_stats       - Get provider statistics"
    echo ""
    echo "Site Operations:"
    echo "  create_site <name>                      - Create new site"
    echo "  list_sites                              - List all sites"
    echo "  get_site_details <site_id>              - Get site details"
    echo "  update_site <site_id>                   - Update site configuration"
    echo "  delete_site <site_id>                   - Delete site"
    echo "  get_site_status_history <site_id>       - Get site status history"
    echo ""
    echo "Machine Operations:"
    echo "  list_machines <site_id>                 - List machines at site"
    echo "  get_machine <machine_id>                - Get machine details"
    echo "  list_machine_capabilities <site_id>     - List machine capabilities"
    echo ""
    echo "Expected Machine:"
    echo "  create_expected_machine <site_id>       - Pre-register expected machine"
    echo "  list_expected_machines <site_id>        - List expected machines"
    echo ""
    echo "IP Block & Networking:"
    echo "  create_ip_block <site_id>               - Create IP block"
    echo "  list_ip_blocks                          - List IP blocks"
    echo ""
    echo "Instance Types:"
    echo "  create_instance_type <site_id>          - Create instance type"
    echo "  list_instance_types <site_id>           - List instance types"
    echo ""
    echo "Tenant Operations:"
    echo "  get_tenant                              - Get tenant details"
    echo "  get_tenant_stats                        - Get tenant statistics"
    echo ""
    echo "Allocation:"
    echo "  create_allocation <site_id> <tenant_id> <instance_type_id>"
    echo "  list_allocations                        - List allocations"
    echo ""
    echo "Audit:"
    echo "  list_audit_logs                         - List audit logs"
    echo "  list_failed_audit_logs                  - List failed audit logs"
    echo ""
    echo "Site Agent:"
    echo "  configure_site_agent <site_id> <token>  - Configure site agent"
    echo ""
    echo "Environment Variables:"
    echo "  TOKEN                                   - Authentication token"
    echo "  INFRA_PROVIDER_ID                       - Infrastructure provider ID"
    echo "  SITE_ID                                 - Current site ID"
    echo "  REGISTRATION_TOKEN                      - Site registration token"
    echo ""
}

# Auto-show functions when sourced
if [ "${BASH_SOURCE[0]}" != "${0}" ]; then
    show_functions
    echo "Functions loaded. Start with: setup_auth"
else
    echo "This script should be sourced, not executed directly:"
    echo "  source $0"
    echo ""
    show_functions
fi

