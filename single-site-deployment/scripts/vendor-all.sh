#!/bin/bash

# Vendor All Dependencies Script
# Downloads and vendors all Go dependencies using your local credentials

set -e

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"

# Color output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

echo -e "${BLUE}========================================${NC}"
echo -e "${BLUE}Vendoring All Dependencies${NC}"
echo -e "${BLUE}========================================${NC}"
echo ""
echo "This will download and vendor Go modules for all services"
echo "using your local Git credentials."
echo ""

services=(
    "carbide-rest-api"
    "carbide-rest-workflow"
    "carbide-rest-site-manager"
    "carbide-rest-ipam"
    "carbide-site-agent"
    "carbide-rest-db"
    "carbide-rest-cert-manager"
    "carbide-rest-auth"
    "carbide-rest-common"
    "carbide-site-workflow"
    "carbide-rest-api-schema"
)

SUCCESS_COUNT=0
FAIL_COUNT=0

for service in "${services[@]}"; do
    service_dir="$ROOT_DIR/$service"
    
    if [ ! -d "$service_dir" ]; then
        echo -e "${YELLOW}Skipping $service (directory not found)${NC}"
        continue
    fi
    
    if [ ! -f "$service_dir/go.mod" ]; then
        echo -e "${YELLOW}Skipping $service (no go.mod)${NC}"
        continue
    fi
    
    echo -e "${BLUE}Processing: $service${NC}"
    
    cd "$service_dir"
    
    # Download dependencies
    echo "  Downloading dependencies..."
    if go mod download 2>&1 | tee "/tmp/vendor-$service.log"; then
        echo -e "${GREEN}  ✓ Downloaded${NC}"
    else
        echo -e "${RED}  ✗ Download failed${NC}"
        echo "  Check: /tmp/vendor-$service.log"
        FAIL_COUNT=$((FAIL_COUNT + 1))
        cd "$ROOT_DIR"
        continue
    fi
    
    # Vendor dependencies
    echo "  Creating vendor directory..."
    if go mod vendor 2>&1 | tee -a "/tmp/vendor-$service.log"; then
        echo -e "${GREEN}  ✓ Vendored${NC}"
        SUCCESS_COUNT=$((SUCCESS_COUNT + 1))
    else
        echo -e "${RED}  ✗ Vendor failed${NC}"
        echo "  Check: /tmp/vendor-$service.log"
        FAIL_COUNT=$((FAIL_COUNT + 1))
    fi
    
    cd "$ROOT_DIR"
    echo ""
done

echo ""
echo -e "${BLUE}========================================${NC}"
echo -e "${BLUE}Vendoring Complete${NC}"
echo -e "${BLUE}========================================${NC}"
echo ""
echo -e "${GREEN}Successful: $SUCCESS_COUNT${NC}"
echo -e "${RED}Failed: $FAIL_COUNT${NC}"
echo ""

if [ $FAIL_COUNT -eq 0 ]; then
    echo -e "${GREEN}All services vendored successfully!${NC}"
    echo ""
    echo -e "${YELLOW}Next steps:${NC}"
    echo "1. Dockerfiles need to be updated to use vendor directory"
    echo "2. Add 'COPY ./vendor ./vendor' to Dockerfiles"
    echo "3. Add '-mod=vendor' to go build commands"
    echo ""
    echo "Or run: ./scripts/update-dockerfiles-for-vendor.sh"
else
    echo -e "${RED}Some services failed to vendor.${NC}"
    echo "Check the log files in /tmp/vendor-*.log for details"
    exit 1
fi

