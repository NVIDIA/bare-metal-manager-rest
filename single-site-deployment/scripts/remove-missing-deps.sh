#!/bin/bash

# Remove Missing Dependencies Script
# Removes references to dependencies not in the snapshot

set -e

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"

# Color output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

echo -e "${BLUE}========================================${NC}"
echo -e "${BLUE}Removing Missing Dependencies${NC}"
echo -e "${BLUE}========================================${NC}"
echo ""

# Dependencies that are not in the snapshot
MISSING_DEPS=(
    "gitlab-master.nvidia.com/ngcsecurity/nvtpm"
)

services=(
    "carbide-rest-api"
    "carbide-rest-workflow"
    "carbide-site-agent"
)

for service in "${services[@]}"; do
    service_dir="$ROOT_DIR/$service"
    go_mod="$service_dir/go.mod"
    go_sum="$service_dir/go.sum"
    
    if [ ! -f "$go_mod" ]; then
        continue
    fi
    
    echo -e "${BLUE}Processing: $service${NC}"
    
    # Backup if not already backed up
    if [ ! -f "$go_mod.backup2" ]; then
        cp "$go_mod" "$go_mod.backup2"
    fi
    
    if [ -f "$go_sum" ] && [ ! -f "$go_sum.backup" ]; then
        cp "$go_sum" "$go_sum.backup"
    fi
    
    for dep in "${MISSING_DEPS[@]}"; do
        if grep -q "$dep" "$go_mod" 2>/dev/null; then
            echo "  Removing $dep from go.mod"
            # Remove lines containing this dependency
            sed -i.tmp "/$dep/d" "$go_mod"
            rm -f "$go_mod.tmp"
        fi
        
        if [ -f "$go_sum" ] && grep -q "$dep" "$go_sum" 2>/dev/null; then
            echo "  Removing $dep from go.sum"
            sed -i.tmp "/$dep/d" "$go_sum"
            rm -f "$go_sum.tmp"
        fi
    done
    
    # Run go mod tidy to fix everything
    echo "  Running go mod tidy..."
    cd "$service_dir"
    go mod tidy 2>&1 | head -10
    cd "$ROOT_DIR"
    
    echo -e "${GREEN}  ✓ Cleaned $service${NC}"
    echo ""
done

echo -e "${GREEN}Missing dependencies removed!${NC}"
echo ""
echo -e "${YELLOW}Next: Try building again${NC}"
echo "  cd single-site-deployment"
echo "  ./scripts/build-all.sh"

