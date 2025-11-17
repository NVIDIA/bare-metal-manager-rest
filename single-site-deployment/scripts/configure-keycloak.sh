#!/bin/bash

# Configure Keycloak for Carbide Development
# Creates realm, client, and test user

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

# Configuration
KEYCLOAK_URL="${KEYCLOAK_URL:-http://localhost:8080}"
ADMIN_USER="${ADMIN_USER:-admin}"
ADMIN_PASSWORD="${ADMIN_PASSWORD:-admin}"
REALM="${REALM:-carbide}"
CLIENT_ID="${CLIENT_ID:-carbide-api}"
CLIENT_SECRET="${CLIENT_SECRET:-carbide-secret-dev-only-do-not-use-in-prod}"

echo "========================================"
echo "Configuring Keycloak"
echo "========================================"
echo ""
echo "Keycloak URL: $KEYCLOAK_URL"
echo "Realm: $REALM"
echo "Client ID: $CLIENT_ID"
echo ""

# Wait for Keycloak to be ready
echo "Waiting for Keycloak to be ready..."
for i in {1..60}; do
    if curl -sf "$KEYCLOAK_URL/health/ready" > /dev/null 2>&1; then
        echo "[SUCCESS] Keycloak is ready"
        break
    fi
    if [ $i -eq 60 ]; then
        echo "[ERROR] Keycloak did not become ready in time"
        exit 1
    fi
    sleep 2
done

# Get admin access token
echo ""
echo "Getting admin access token..."
ADMIN_TOKEN=$(curl -sf -X POST "$KEYCLOAK_URL/realms/master/protocol/openid-connect/token" \
    -H "Content-Type: application/x-www-form-urlencoded" \
    -d "username=$ADMIN_USER" \
    -d "password=$ADMIN_PASSWORD" \
    -d "grant_type=password" \
    -d "client_id=admin-cli" | jq -r '.access_token')

if [ -z "$ADMIN_TOKEN" ] || [ "$ADMIN_TOKEN" == "null" ]; then
    echo "[ERROR] Failed to get admin token"
    exit 1
fi
echo "[SUCCESS] Admin token obtained"

# Create realm
echo ""
echo "Creating realm: $REALM"
REALM_CREATE=$(curl -s -w "%{http_code}" -X POST "$KEYCLOAK_URL/admin/realms" \
    -H "Authorization: Bearer $ADMIN_TOKEN" \
    -H "Content-Type: application/json" \
    -d "{
        \"realm\": \"$REALM\",
        \"enabled\": true,
        \"displayName\": \"Carbide Development\",
        \"accessTokenLifespan\": 3600,
        \"ssoSessionIdleTimeout\": 3600,
        \"ssoSessionMaxLifespan\": 36000
    }")

if echo "$REALM_CREATE" | grep -q "409"; then
    echo "[INFO] Realm already exists, continuing..."
elif echo "$REALM_CREATE" | grep -q "201"; then
    echo "[SUCCESS] Realm created"
else
    echo "[INFO] Realm response: $REALM_CREATE"
fi

# Create client
echo ""
echo "Creating client: $CLIENT_ID"
CLIENT_CREATE=$(curl -s -w "%{http_code}" -X POST "$KEYCLOAK_URL/admin/realms/$REALM/clients" \
    -H "Authorization: Bearer $ADMIN_TOKEN" \
    -H "Content-Type: application/json" \
    -d "{
        \"clientId\": \"$CLIENT_ID\",
        \"enabled\": true,
        \"protocol\": \"openid-connect\",
        \"publicClient\": false,
        \"directAccessGrantsEnabled\": true,
        \"serviceAccountsEnabled\": true,
        \"authorizationServicesEnabled\": true,
        \"secret\": \"$CLIENT_SECRET\",
        \"redirectUris\": [\"*\"],
        \"webOrigins\": [\"*\"],
        \"attributes\": {
            \"access.token.lifespan\": \"3600\"
        }
    }")

if echo "$CLIENT_CREATE" | grep -q "409"; then
    echo "[INFO] Client already exists, continuing..."
elif echo "$CLIENT_CREATE" | grep -q "201"; then
    echo "[SUCCESS] Client created"
else
    echo "[INFO] Client response: continuing..."
fi

# Create test user
echo ""
echo "Creating test user: testuser"
USER_CREATE=$(curl -s -w "%{http_code}" -X POST "$KEYCLOAK_URL/admin/realms/$REALM/users" \
    -H "Authorization: Bearer $ADMIN_TOKEN" \
    -H "Content-Type: application/json" \
    -d "{
        \"username\": \"testuser\",
        \"enabled\": true,
        \"email\": \"testuser@carbide.local\",
        \"firstName\": \"Test\",
        \"lastName\": \"User\",
        \"credentials\": [{
            \"type\": \"password\",
            \"value\": \"testpass\",
            \"temporary\": false
        }]
    }")

if echo "$USER_CREATE" | grep -q "409"; then
    echo "[INFO] User already exists, continuing..."
elif echo "$USER_CREATE" | grep -q "201"; then
    echo "[SUCCESS] Test user created"
else
    echo "[INFO] User response: continuing..."
fi

# Create admin user
echo ""
echo "Creating admin user: admin"
ADMIN_CREATE=$(curl -s -w "%{http_code}" -X POST "$KEYCLOAK_URL/admin/realms/$REALM/users" \
    -H "Authorization: Bearer $ADMIN_TOKEN" \
    -H "Content-Type: application/json" \
    -d "{
        \"username\": \"admin\",
        \"enabled\": true,
        \"email\": \"admin@carbide.local\",
        \"firstName\": \"Admin\",
        \"lastName\": \"User\",
        \"credentials\": [{
            \"type\": \"password\",
            \"value\": \"admin\",
            \"temporary\": false
        }]
    }")

if echo "$ADMIN_CREATE" | grep -q "409"; then
    echo "[INFO] Admin user already exists, continuing..."
elif echo "$ADMIN_CREATE" | grep -q "201"; then
    echo "[SUCCESS] Admin user created"
else
    echo "[INFO] Admin user response: continuing..."
fi

# Create roles for different orgs
# Format: orgname:FORGE_PROVIDER_ADMIN, orgname:FORGE_TENANT_ADMIN
echo ""
echo "Creating realm roles..."
for org in "nvidia" "test-org" "dev-org"; do
    for role_type in "FORGE_PROVIDER_ADMIN" "FORGE_TENANT_ADMIN" "FORGE_PROVIDER_VIEWER"; do
        role_name="${org}:${role_type}"
        ROLE_CREATE=$(curl -s -w "%{http_code}" -X POST "$KEYCLOAK_URL/admin/realms/$REALM/roles" \
            -H "Authorization: Bearer $ADMIN_TOKEN" \
            -H "Content-Type: application/json" \
            -d "{
                \"name\": \"$role_name\",
                \"description\": \"$role_type for $org organization\"
            }")
        
        if echo "$ROLE_CREATE" | grep -q "409"; then
            echo "[INFO] Role '$role_name' already exists"
        elif echo "$ROLE_CREATE" | grep -q "201"; then
            echo "[SUCCESS] Role '$role_name' created"
        else
            echo "[INFO] Role '$role_name': continuing..."
        fi
    done
done

# Get user IDs
echo ""
echo "Getting user IDs..."
TESTUSER_ID=$(curl -sf "$KEYCLOAK_URL/admin/realms/$REALM/users?username=testuser" \
    -H "Authorization: Bearer $ADMIN_TOKEN" | jq -r '.[0].id')
ADMIN_USER_ID=$(curl -sf "$KEYCLOAK_URL/admin/realms/$REALM/users?username=admin" \
    -H "Authorization: Bearer $ADMIN_TOKEN" | jq -r '.[0].id')

if [ -z "$TESTUSER_ID" ] || [ "$TESTUSER_ID" == "null" ]; then
    echo "[ERROR] Failed to get testuser ID"
else
    echo "[SUCCESS] Found testuser ID: $TESTUSER_ID"
fi

if [ -z "$ADMIN_USER_ID" ] || [ "$ADMIN_USER_ID" == "null" ]; then
    echo "[ERROR] Failed to get admin user ID"
else
    echo "[SUCCESS] Found admin user ID: $ADMIN_USER_ID"
fi

# Assign roles to testuser (provider and tenant admin for nvidia and test-org)
echo ""
echo "Assigning roles to testuser..."
for role_name in "nvidia:FORGE_PROVIDER_ADMIN" "nvidia:FORGE_TENANT_ADMIN" "test-org:FORGE_TENANT_ADMIN"; do
    ROLE_ID=$(curl -sf "$KEYCLOAK_URL/admin/realms/$REALM/roles/$role_name" \
        -H "Authorization: Bearer $ADMIN_TOKEN" | jq -r '.id')
    
    if [ -n "$ROLE_ID" ] && [ "$ROLE_ID" != "null" ]; then
        curl -sf -X POST "$KEYCLOAK_URL/admin/realms/$REALM/users/$TESTUSER_ID/role-mappings/realm" \
            -H "Authorization: Bearer $ADMIN_TOKEN" \
            -H "Content-Type: application/json" \
            -d "[{\"id\": \"$ROLE_ID\", \"name\": \"$role_name\"}]" > /dev/null 2>&1
        echo "[SUCCESS] Assigned role '$role_name' to testuser"
    fi
done

# Assign all admin roles to admin user
echo ""
echo "Assigning admin roles to admin user..."
for org in "nvidia" "test-org" "dev-org"; do
    role_name="${org}:FORGE_PROVIDER_ADMIN"
    ROLE_ID=$(curl -sf "$KEYCLOAK_URL/admin/realms/$REALM/roles/$role_name" \
        -H "Authorization: Bearer $ADMIN_TOKEN" | jq -r '.id')
    
    if [ -n "$ROLE_ID" ] && [ "$ROLE_ID" != "null" ]; then
        curl -sf -X POST "$KEYCLOAK_URL/admin/realms/$REALM/users/$ADMIN_USER_ID/role-mappings/realm" \
            -H "Authorization: Bearer $ADMIN_TOKEN" \
            -H "Content-Type: application/json" \
            -d "[{\"id\": \"$ROLE_ID\", \"name\": \"$role_name\"}]" > /dev/null 2>&1
        echo "[SUCCESS] Assigned role '$role_name' to admin"
    fi
done

echo ""
echo "========================================"
echo "Keycloak Configuration Complete!"
echo "========================================"
echo ""
echo "Keycloak Details:"
echo "  URL:            $KEYCLOAK_URL"
echo "  Admin Console:  $KEYCLOAK_URL/admin"
echo "  Admin User:     $ADMIN_USER / $ADMIN_PASSWORD"
echo ""
echo "Realm Details:"
echo "  Realm:          $REALM"
echo "  Client ID:      $CLIENT_ID"
echo "  Client Secret:  $CLIENT_SECRET"
echo ""
echo "Test Users:"
echo "  Username:       testuser"
echo "  Password:       testpass"
echo ""
echo "  Username:       admin"
echo "  Password:       admin"
echo ""
echo "Get Access Token:"
echo "  curl -X POST '$KEYCLOAK_URL/realms/$REALM/protocol/openid-connect/token' \\"
echo "    -H 'Content-Type: application/x-www-form-urlencoded' \\"
echo "    -d 'username=testuser' \\"
echo "    -d 'password=testpass' \\"
echo "    -d 'grant_type=password' \\"
echo "    -d 'client_id=$CLIENT_ID' \\"
echo "    -d 'client_secret=$CLIENT_SECRET' | jq -r '.access_token'"
echo ""

