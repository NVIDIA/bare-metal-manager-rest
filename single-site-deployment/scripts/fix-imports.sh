#!/bin/bash

# Fix Import Paths Script
# Adds replace directives to go.mod files to use local copies instead of GitLab

set -e

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"

# Color output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

echo -e "${BLUE}========================================${NC}"
echo -e "${BLUE}Fixing Import Paths${NC}"
echo -e "${BLUE}========================================${NC}"
echo ""
echo "Adding replace directives to use local code instead of GitLab"
echo ""

# Map of GitLab import paths to local directories
declare -A REPLACEMENTS=(
    ["gitlab-master.nvidia.com/nvmetal/cloud-api"]="carbide-rest-api"
    ["gitlab-master.nvidia.com/nvmetal/cloud-workflow"]="carbide-rest-workflow"
    ["gitlab-master.nvidia.com/nvmetal/cloud-site-manager"]="carbide-rest-site-manager"
    ["gitlab-master.nvidia.com/nvmetal/cloud-ipam"]="carbide-rest-ipam"
    ["gitlab-master.nvidia.com/nvmetal/elektra-site-agent"]="carbide-site-agent"
    ["gitlab-master.nvidia.com/nvmetal/cloud-db"]="carbide-rest-db"
    ["gitlab-master.nvidia.com/nvmetal/cloud-cert-manager"]="carbide-rest-cert-manager"
    ["gitlab-master.nvidia.com/nvmetal/cloud-auth"]="carbide-rest-auth"
    ["gitlab-master.nvidia.com/nvmetal/cloud-common"]="carbide-rest-common"
    ["gitlab-master.nvidia.com/nvmetal/site-workflow"]="carbide-site-workflow"
    ["gitlab-master.nvidia.com/nvmetal/cloud-api-schema"]="carbide-rest-api-schema"
)

# Services to process
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

for service in "${services[@]}"; do
    service_dir="$ROOT_DIR/$service"
    go_mod="$service_dir/go.mod"
    
    if [ ! -d "$service_dir" ]; then
        echo -e "${YELLOW}Skipping $service (directory not found)${NC}"
        continue
    fi
    
    if [ ! -f "$go_mod" ]; then
        echo -e "${YELLOW}Skipping $service (no go.mod)${NC}"
        continue
    fi
    
    echo -e "${BLUE}Processing: $service${NC}"
    
    # Backup original go.mod
    cp "$go_mod" "$go_mod.backup"
    
    # Check if replace block already exists
    if grep -q "^replace (" "$go_mod"; then
        echo -e "${YELLOW}  Replace block already exists, skipping${NC}"
        continue
    fi
    
    # Add replace directives
    echo "" >> "$go_mod"
    echo "// Local replacements for snapshot" >> "$go_mod"
    echo "replace (" >> "$go_mod"
    
    # Add each replacement
    for import_path in "${!REPLACEMENTS[@]}"; do
        local_dir="${REPLACEMENTS[$import_path]}"
        
        # Check if this go.mod actually uses this import
        if grep -q "$import_path" "$go_mod"; then
            # Calculate relative path from service to dependency
            relative_path="../$local_dir"
            echo "    $import_path => $relative_path" >> "$go_mod"
            echo -e "  ${GREEN}✓${NC} Added: $import_path => $relative_path"
        fi
    done
    
    echo ")" >> "$go_mod"
    
    echo -e "${GREEN}  ✓ Updated go.mod${NC}"
    SUCCESS_COUNT=$((SUCCESS_COUNT + 1))
    
    # Run go mod tidy to clean up
    echo "  Running go mod tidy..."
    cd "$service_dir"
    if go mod tidy 2>&1 | grep -v "warning"; then
        echo -e "${GREEN}  ✓ go mod tidy complete${NC}"
    else
        echo -e "${YELLOW}  ⚠ go mod tidy had warnings${NC}"
    fi
    cd "$ROOT_DIR"
    
    echo ""
done

echo ""
echo -e "${BLUE}========================================${NC}"
echo -e "${BLUE}Import Fix Complete${NC}"
echo -e "${BLUE}========================================${NC}"
echo ""
echo -e "${GREEN}Updated: $SUCCESS_COUNT services${NC}"
echo ""
echo -e "${YELLOW}Next steps:${NC}"
echo "1. Review the changes: git diff */go.mod"
echo "2. Build images: ./scripts/build-all.sh"
echo ""
echo -e "${YELLOW}To revert changes:${NC}"
echo "find . -name 'go.mod.backup' -exec bash -c 'mv \"\$0\" \"\${0%.backup}\"' {} \;"

