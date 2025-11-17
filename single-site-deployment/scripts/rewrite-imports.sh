#!/bin/bash

# Rewrite Imports Script
# Changes all module names and imports from gitlab-master.nvidia.com/nvmetal to carbide
# This script is idempotent and can be re-run safely

set -e

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"

# Color output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

echo -e "${BLUE}========================================${NC}"
echo -e "${BLUE}Rewriting All Import Paths${NC}"
echo -e "${BLUE}========================================${NC}"
echo ""
echo "Changing gitlab-master.nvidia.com/nvmetal/* to carbide/*"
echo ""

# Mapping of old paths to new paths
declare -A PATH_MAP=(
    ["gitlab-master.nvidia.com/nvmetal/cloud-api"]="carbide/rest-api"
    ["gitlab-master.nvidia.com/nvmetal/cloud-api-schema"]="carbide/rest-api-schema"
    ["gitlab-master.nvidia.com/nvmetal/cloud-auth"]="carbide/rest-auth"
    ["gitlab-master.nvidia.com/nvmetal/cloud-cert-manager"]="carbide/rest-cert-manager"
    ["gitlab-master.nvidia.com/nvmetal/cloud-common"]="carbide/rest-common"
    ["gitlab-master.nvidia.com/nvmetal/cloud-db"]="carbide/rest-db"
    ["gitlab-master.nvidia.com/nvmetal/cloud-ipam"]="carbide/rest-ipam"
    ["gitlab-master.nvidia.com/nvmetal/cloud-site-manager"]="carbide/rest-site-manager"
    ["gitlab-master.nvidia.com/nvmetal/cloud-workflow"]="carbide/rest-workflow"
    ["gitlab-master.nvidia.com/nvmetal/cloud-workflow-schema"]="carbide/rest-api-schema"
    ["gitlab-master.nvidia.com/nvmetal/elektra-site-agent"]="carbide/site-agent"
    ["gitlab-master.nvidia.com/nvmetal/site-workflow"]="carbide/site-workflow"
)

# Directories to process
services=(
    "carbide-rest-api"
    "carbide-rest-api-schema"
    "carbide-rest-auth"
    "carbide-rest-cert-manager"
    "carbide-rest-common"
    "carbide-rest-db"
    "carbide-rest-ipam"
    "carbide-rest-site-manager"
    "carbide-rest-workflow"
    "carbide-site-agent"
    "carbide-site-workflow"
)

echo -e "${BLUE}Step 1: Rewriting module declarations in go.mod files${NC}"
echo ""

for service in "${services[@]}"; do
    service_dir="$ROOT_DIR/$service"
    go_mod="$service_dir/go.mod"
    
    if [ ! -f "$go_mod" ]; then
        echo -e "${YELLOW}Skipping $service (no go.mod)${NC}"
        continue
    fi
    
    echo "Processing: $service"
    
    # Rewrite module declaration and all imports in go.mod
    for old_path in "${!PATH_MAP[@]}"; do
        new_path="${PATH_MAP[$old_path]}"
        
        # Update module declaration
        sed -i.bak "s|^module $old_path|module $new_path|g" "$go_mod"
        
        # Update require statements
        sed -i.bak "s|$old_path |$new_path |g" "$go_mod"
        sed -i.bak "s|$old_path\$|$new_path|g" "$go_mod"
    done
    
    # Remove any replace directives that reference gitlab
    sed -i.bak '/^\/\/ Local replacements for snapshot/,/^)$/d' "$go_mod"
    sed -i.bak '/^replace ($/,/^)$/{/gitlab-master.nvidia.com/d;}' "$go_mod"
    
    # Clean up empty replace blocks
    sed -i.bak '/^replace ($/N;/^replace (\n)$/d' "$go_mod"
    
    rm -f "$go_mod.bak"
    echo -e "  ${GREEN}✓ Updated go.mod${NC}"
done

echo ""
echo -e "${BLUE}Step 2: Rewriting import statements in .go files${NC}"
echo ""

for service in "${services[@]}"; do
    service_dir="$ROOT_DIR/$service"
    
    if [ ! -d "$service_dir" ]; then
        continue
    fi
    
    echo "Processing: $service"
    
    # Find all .go files and rewrite imports
    find "$service_dir" -name "*.go" -type f | while read -r go_file; do
        for old_path in "${!PATH_MAP[@]}"; do
            new_path="${PATH_MAP[$old_path]}"
            sed -i.bak "s|\"$old_path|\"$new_path|g" "$go_file"
            sed -i.bak "s|'$old_path|'$new_path|g" "$go_file"
        done
        rm -f "$go_file.bak"
    done
    
    echo -e "  ${GREEN}✓ Updated .go files${NC}"
done

echo ""
echo -e "${BLUE}Step 3: Running go mod tidy on all services${NC}"
echo ""

for service in "${services[@]}"; do
    service_dir="$ROOT_DIR/$service"
    
    if [ ! -f "$service_dir/go.mod" ]; then
        continue
    fi
    
    echo "Tidying: $service"
    cd "$service_dir"
    go mod tidy 2>&1 | head -5 || true
    echo -e "  ${GREEN}✓ Tidied${NC}"
    cd "$ROOT_DIR"
done

echo ""
echo -e "${BLUE}========================================${NC}"
echo -e "${BLUE}Import Rewrite Complete!${NC}"
echo -e "${BLUE}========================================${NC}"
echo ""

echo -e "${GREEN}Summary:${NC}"
echo "  - Changed all module names from gitlab-master.nvidia.com/nvmetal/* to carbide/*"
echo "  - Updated all import statements in .go files"
echo "  - Removed replace directives"
echo "  - Ran go mod tidy"
echo ""

echo -e "${YELLOW}Next steps:${NC}"
echo "  1. Review changes: git diff"
echo "  2. Build images: ./scripts/build-all.sh"
echo ""

echo -e "${BLUE}Verification:${NC}"
echo "Checking for remaining gitlab references..."
remaining=$(grep -r "gitlab-master.nvidia.com/nvmetal" "$ROOT_DIR"/carbide-*/go.mod | grep -v "^Binary" | wc -l)
if [ "$remaining" -eq 0 ]; then
    echo -e "${GREEN}✓ No gitlab references found in go.mod files!${NC}"
else
    echo -e "${YELLOW}⚠ Found $remaining remaining gitlab references${NC}"
    grep -r "gitlab-master.nvidia.com/nvmetal" "$ROOT_DIR"/carbide-*/go.mod | head -10
fi

