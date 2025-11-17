#!/bin/bash

# Tidy All Go Modules Script
# Runs go mod tidy on all services to clean up indirect dependencies

set -e

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"

# Color output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

echo -e "${BLUE}========================================${NC}"
echo -e "${BLUE}Running go mod tidy on all services${NC}"
echo -e "${BLUE}========================================${NC}"
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
    
    echo -e "${BLUE}Tidying: $service${NC}"
    
    cd "$service_dir"
    
    if go mod tidy 2>&1 | tee "/tmp/tidy-$service.log"; then
        # Check if nvtpm was removed
        if grep -q "nvtpm" go.mod 2>/dev/null; then
            echo -e "${YELLOW}  ⚠ Still has nvtpm reference${NC}"
        else
            echo -e "${GREEN}  ✓ nvtpm removed (if it was there)${NC}"
        fi
        echo -e "${GREEN}  ✓ Tidied${NC}"
        SUCCESS_COUNT=$((SUCCESS_COUNT + 1))
    else
        echo -e "${RED}  ✗ Tidy failed${NC}"
        echo "  Check: /tmp/tidy-$service.log"
        FAIL_COUNT=$((FAIL_COUNT + 1))
    fi
    
    cd "$ROOT_DIR"
    echo ""
done

echo ""
echo -e "${BLUE}========================================${NC}"
echo -e "${BLUE}Tidy Complete${NC}"
echo -e "${BLUE}========================================${NC}"
echo ""
echo -e "${GREEN}Successful: $SUCCESS_COUNT${NC}"
echo -e "${RED}Failed: $FAIL_COUNT${NC}"
echo ""

if [ $FAIL_COUNT -eq 0 ]; then
    echo -e "${GREEN}All services tidied successfully!${NC}"
    echo ""
    echo -e "${YELLOW}Next step: Try building${NC}"
    echo "  ./scripts/build-all.sh"
else
    echo -e "${RED}Some services failed to tidy.${NC}"
    echo "Check the log files in /tmp/tidy-*.log for details"
    exit 1
fi

