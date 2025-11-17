#!/bin/bash

# Convert to True Monorepo
# Changes all imports to github.com/nvidia/carbide-rest-api-snapshot
# Removes child go.mod files
# Creates single go.mod at root

set -e

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"

# Color output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

echo -e "${BLUE}╔════════════════════════════════════════╗${NC}"
echo -e "${BLUE}║  Converting to True Monorepo          ║${NC}"
echo -e "${BLUE}╚════════════════════════════════════════╝${NC}"
echo ""

GITHUB_BASE="github.com/nvidia/carbide-rest-api-snapshot"

# Mapping of old gitlab paths to new github paths
declare -A REPLACEMENTS=(
    ["gitlab-master.nvidia.com/nvmetal/cloud-api"]="$GITHUB_BASE/carbide-rest-api"
    ["gitlab-master.nvidia.com/nvmetal/cloud-api-schema"]="$GITHUB_BASE/carbide-rest-api-schema"
    ["gitlab-master.nvidia.com/nvmetal/cloud-auth"]="$GITHUB_BASE/carbide-rest-auth"
    ["gitlab-master.nvidia.com/nvmetal/cloud-cert-manager"]="$GITHUB_BASE/carbide-rest-cert-manager"
    ["gitlab-master.nvidia.com/nvmetal/cloud-common"]="$GITHUB_BASE/carbide-rest-common"
    ["gitlab-master.nvidia.com/nvmetal/cloud-db"]="$GITHUB_BASE/carbide-rest-db"
    ["gitlab-master.nvidia.com/nvmetal/cloud-ipam"]="$GITHUB_BASE/carbide-rest-ipam"
    ["gitlab-master.nvidia.com/nvmetal/cloud-site-manager"]="$GITHUB_BASE/carbide-rest-site-manager"
    ["gitlab-master.nvidia.com/nvmetal/cloud-workflow"]="$GITHUB_BASE/carbide-rest-workflow"
    ["gitlab-master.nvidia.com/nvmetal/cloud-workflow-schema"]="$GITHUB_BASE/carbide-rest-api-schema"
    ["gitlab-master.nvidia.com/nvmetal/elektra-site-agent"]="$GITHUB_BASE/carbide-site-agent"
    ["gitlab-master.nvidia.com/nvmetal/site-workflow"]="$GITHUB_BASE/carbide-site-workflow"
)

echo -e "${BLUE}Step 1: Updating all import paths in .go files${NC}"
echo ""

cd "$ROOT_DIR"

for old_path in "${!REPLACEMENTS[@]}"; do
    new_path="${REPLACEMENTS[$old_path]}"
    echo "  $old_path"
    echo "    → $new_path"
    
    # Use find and sed to replace in all .go files
    find ./carbide-* -name "*.go" -type f -exec sed -i '' "s|$old_path|$new_path|g" {} +
done

echo ""
echo -e "${GREEN}✓ Updated all .go files${NC}"
echo ""

echo -e "${BLUE}Step 2: Backing up and removing child go.mod files${NC}"
echo ""

mkdir -p /tmp/carbide-gomod-backup

for gomod in ./carbide-*/go.mod; do
    if [ -f "$gomod" ]; then
        service=$(basename $(dirname "$gomod"))
        echo "  Removing: $gomod"
        cp "$gomod" "/tmp/carbide-gomod-backup/$service-go.mod.backup"
        rm "$gomod"
    fi
done

for gosum in ./carbide-*/go.sum; do
    if [ -f "$gosum" ]; then
        rm "$gosum"
    fi
done

echo -e "${GREEN}✓ Removed all child go.mod and go.sum files${NC}"
echo -e "${YELLOW}  Backups in: /tmp/carbide-gomod-backup/${NC}"
echo ""

echo -e "${BLUE}Step 3: Creating new root go.mod${NC}"
echo ""

cat > go.mod << EOF
module github.com/nvidia/carbide-rest-api-snapshot

go 1.21

// This is a monorepo containing all Carbide services
// All internal dependencies are resolved within this module
EOF

echo -e "${GREEN}✓ Created root go.mod${NC}"
echo ""

echo -e "${BLUE}Step 4: Running go mod tidy (this may take a minute...)${NC}"
echo ""

go mod tidy

echo ""
echo -e "${GREEN}✓ go mod tidy complete${NC}"
echo ""

echo -e "${BLUE}Step 5: Verification${NC}"
echo ""

# Check for any remaining gitlab references
gitlab_count=$(find ./carbide-* -name "*.go" -type f -exec grep -l "gitlab-master.nvidia.com/nvmetal" {} + 2>/dev/null | wc -l | tr -d ' ')

if [ "$gitlab_count" -eq "0" ]; then
    echo -e "${GREEN}✓ No gitlab references in .go files!${NC}"
else
    echo -e "${RED}✗ Found $gitlab_count files with gitlab references${NC}"
    find ./carbide-* -name "*.go" -type f -exec grep -l "gitlab-master.nvidia.com/nvmetal" {} + | head -5
fi

echo ""
echo -e "${BLUE}╔════════════════════════════════════════╗${NC}"
echo -e "${GREEN}║     MONOREPO CONVERSION COMPLETE!     ║${NC}"
echo -e "${BLUE}╚════════════════════════════════════════╝${NC}"
echo ""

echo "Summary:"
echo "  • All imports changed to: $GITHUB_BASE/carbide-*"
echo "  • Removed all child go.mod files"
echo "  • Created single root go.mod"
echo "  • Ran go mod tidy"
echo ""

echo "Files:"
ls -lh go.mod go.sum
echo ""

echo -e "${YELLOW}Next: Build and test${NC}"
echo "  cd single-site-deployment"
echo "  ./scripts/build-simple.sh"

