#!/bin/bash

# Rewrite Imports Script V2
# Uses perl for more reliable string replacement

set -e

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"

# Color output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

echo -e "${BLUE}========================================${NC}"
echo -e "${BLUE}Rewriting All Import Paths (V2)${NC}"
echo -e "${BLUE}========================================${NC}"
echo ""

# Replacement mappings
declare -A REPLACEMENTS=(
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

echo -e "${BLUE}Step 1: Processing go.mod files${NC}"
echo ""

for gomod in "$ROOT_DIR"/carbide-*/go.mod; do
    if [ ! -f "$gomod" ]; then
        continue
    fi
    
    service=$(basename $(dirname "$gomod"))
    echo "Processing: $service"
    
    # Create backup
    cp "$gomod" "$gomod.pre-rewrite"
    
    # Apply all replacements using perl
    for old in "${!REPLACEMENTS[@]}"; do
        new="${REPLACEMENTS[$old]}"
        perl -pi -e "s|\\Q$old\\E|$new|g" "$gomod"
    done
    
    echo -e "  ${GREEN}✓ Rewritten${NC}"
done

echo ""
echo -e "${BLUE}Step 2: Processing .go files${NC}"
echo ""

for service_dir in "$ROOT_DIR"/carbide-*; do
    if [ ! -d "$service_dir" ]; then
        continue
    fi
    
    service=$(basename "$service_dir")
    echo "Processing: $service"
    
    # Find and process all .go files
    find "$service_dir" -name "*.go" -type f -print0 | while IFS= read -r -d '' gofile; do
        for old in "${!REPLACEMENTS[@]}"; do
            new="${REPLACEMENTS[$old]}"
            perl -pi -e "s|\\Q$old\\E|$new|g" "$gofile"
        done
    done
    
    echo -e "  ${GREEN}✓ Updated${NC}"
done

echo ""
echo -e "${BLUE}Step 3: Cleaning up go.mod files${NC}"
echo ""

# Remove replace blocks that reference gitlab
for gomod in "$ROOT_DIR"/carbide-*/go.mod; do
    if [ ! -f "$gomod" ]; then
        continue
    fi
    
    service=$(basename $(dirname "$gomod"))
    
    # Remove replace blocks using perl
    perl -i -pe 'BEGIN{undef $/;} s|// Local replacements for snapshot\nreplace \([^)]*\)\n||sg' "$gomod"
    
    # Remove any remaining standalone replace lines with gitlab
    perl -pi -e '/^replace .*gitlab-master\.nvidia\.com/d' "$gomod"
done

echo -e "${GREEN}✓ Removed replace directives${NC}"
echo ""

echo -e "${BLUE}Step 4: Running go mod tidy${NC}"
echo ""

for service_dir in "$ROOT_DIR"/carbide-*; do
    if [ ! -f "$service_dir/go.mod" ]; then
        continue
    fi
    
    service=$(basename "$service_dir")
    echo "Tidying: $service"
    
    cd "$service_dir"
    go mod tidy 2>&1 | head -3 || true
    cd "$ROOT_DIR"
done

echo ""
echo -e "${BLUE}Step 5: Creating go.work${NC}"
echo ""

cd "$ROOT_DIR"
if [ -f "go.work" ]; then
    echo -e "${YELLOW}go.work exists, recreating...${NC}"
    rm go.work
fi

go work init
for service_dir in "$ROOT_DIR"/carbide-*; do
    if [ -f "$service_dir/go.mod" ]; then
        go work use "./$( basename "$service_dir")"
    fi
done

echo -e "${GREEN}✓ Created go.work${NC}"

echo ""
echo -e "${BLUE}========================================${NC}"
echo -e "${GREEN}✓ Import Rewrite Complete!${NC}"
echo -e "${BLUE}========================================${NC}"
echo ""

# Final verification
echo "Verification:"
gitlab_count=$(grep -r "gitlab-master.nvidia.com/nvmetal" "$ROOT_DIR"/carbide-*/go.mod 2>/dev/null | wc -l | tr -d ' ')

if [ "$gitlab_count" -eq "0" ]; then
    echo -e "${GREEN}✓ No gitlab references in go.mod files!${NC}"
else
    echo -e "${RED}✗ Still found $gitlab_count gitlab references${NC}"
    echo "Remaining references:"
    grep -n "gitlab-master.nvidia.com/nvmetal" "$ROOT_DIR"/carbide-*/go.mod | head -10
fi

echo ""
echo -e "${YELLOW}Next: Build images${NC}"
echo "  ./scripts/build-simple.sh"

