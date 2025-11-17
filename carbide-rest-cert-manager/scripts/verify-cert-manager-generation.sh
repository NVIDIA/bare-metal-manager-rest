#!/bin/bash
# SPDX-FileCopyrightText: Copyright (c) 2021-2023 NVIDIA CORPORATION & AFFILIATES. All rights reserved.
# SPDX-License-Identifier: LicenseRef-NvidiaProprietary
#
# Verify that generated cert-manager manifests match existing files
# This script downloads cert-manager.yaml from GitHub and compares with existing files

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "${SCRIPT_DIR}/.." && pwd)"

# cert-manager version
CERT_MANAGER_VERSION="v1.9.1"
GITHUB_RELEASE_URL="https://github.com/cert-manager/cert-manager/releases/download/${CERT_MANAGER_VERSION}/cert-manager.yaml"

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

echo -e "${GREEN}Verifying cert-manager manifest generation...${NC}"
echo -e "${YELLOW}Version: ${CERT_MANAGER_VERSION}${NC}"

# Check curl
if ! command -v curl &> /dev/null; then
    echo -e "${RED}Error: curl is not installed${NC}"
    exit 1
fi

# Download once to temp file
TEMP_FILE=$(mktemp)
trap "rm -f ${TEMP_FILE}" EXIT

echo -e "${YELLOW}Downloading from GitHub...${NC}"
if ! curl -fsSL "${GITHUB_RELEASE_URL}" -o "${TEMP_FILE}"; then
    echo -e "${RED}Error: Failed to download${NC}"
    exit 1
fi

# Check each overlay
OVERLAYS=("local" "cloud-local" "aws-dev" "aws-stg")
ALL_MATCH=true

for overlay in "${OVERLAYS[@]}"; do
    EXISTING_FILE="${REPO_ROOT}/kustomize/cert-manager/overlays/${overlay}/cert-manager.yaml"
    
    echo -e "\n${GREEN}Checking ${overlay}...${NC}"
    
    if [ ! -f "${EXISTING_FILE}" ]; then
        echo -e "${YELLOW}  ⚠ No existing file${NC}"
        continue
    fi
    
    # Skip NVIDIA header (first 12 lines) and compare
    if diff -q <(tail -n +13 "${EXISTING_FILE}") "${TEMP_FILE}" > /dev/null 2>&1; then
        echo -e "${GREEN}  ✓ Files match exactly (excluding NVIDIA header)${NC}"
    else
        # Count differences
        DIFF_LINES=$(diff <(tail -n +13 "${EXISTING_FILE}") "${TEMP_FILE}" | wc -l | tr -d ' ')
        echo -e "${YELLOW}  ⚠ Files differ (${DIFF_LINES} diff lines)${NC}"
        echo -e "${YELLOW}  Showing first 20 differences:${NC}"
        diff <(tail -n +13 "${EXISTING_FILE}") "${TEMP_FILE}" | head -20 || true
        ALL_MATCH=false
    fi
done

echo ""
if [ "$ALL_MATCH" = true ]; then
    echo -e "${GREEN}✓ All files match - safe to remove and regenerate${NC}"
    exit 0
else
    echo -e "${YELLOW}⚠ Some differences found - review above${NC}"
    exit 1
fi
