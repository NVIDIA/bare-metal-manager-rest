#!/bin/bash
# SPDX-FileCopyrightText: Copyright (c) 2021-2023 NVIDIA CORPORATION & AFFILIATES. All rights reserved.
# SPDX-License-Identifier: LicenseRef-NvidiaProprietary
#
# Generate cert-manager manifests from official GitHub release
# This script downloads the cert-manager.yaml files from the cert-manager GitHub releases

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "${SCRIPT_DIR}/.." && pwd)"

# cert-manager version - update this when upgrading cert-manager
CERT_MANAGER_VERSION="v1.9.1"
GITHUB_RELEASE_URL="https://github.com/cert-manager/cert-manager/releases/download/${CERT_MANAGER_VERSION}/cert-manager.yaml"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${GREEN}Downloading cert-manager manifests from GitHub release...${NC}"
echo -e "${YELLOW}Version: ${CERT_MANAGER_VERSION}${NC}"

# Check if curl is installed
if ! command -v curl &> /dev/null; then
    echo -e "${RED}Error: curl is not installed${NC}"
    echo "Please install curl to download cert-manager manifests"
    exit 1
fi

# Download the cert-manager manifest once to a temporary file
TEMP_FILE=$(mktemp)
trap "rm -f ${TEMP_FILE}" EXIT

echo -e "${YELLOW}Downloading from ${GITHUB_RELEASE_URL}...${NC}"
if ! curl -fsSL "${GITHUB_RELEASE_URL}" -o "${TEMP_FILE}"; then
    echo -e "${RED}Error: Failed to download cert-manager manifest${NC}"
    echo -e "${RED}URL: ${GITHUB_RELEASE_URL}${NC}"
    exit 1
fi

echo -e "${GREEN}✓ Downloaded successfully${NC}"

# Copy to each overlay that needs cert-manager.yaml
OVERLAYS=(
    "local"
    "cloud-local"
)

for overlay in "${OVERLAYS[@]}"; do
    OVERLAY_DIR="${REPO_ROOT}/kustomize/cert-manager/overlays/${overlay}"
    OUTPUT_FILE="${OVERLAY_DIR}/cert-manager.yaml"
    
    if [ ! -d "${OVERLAY_DIR}" ]; then
        echo -e "${YELLOW}Warning: Overlay directory ${overlay} does not exist, skipping...${NC}"
        continue
    fi
    
    echo -e "${GREEN}Copying to ${OUTPUT_FILE}...${NC}"
    cp "${TEMP_FILE}" "${OUTPUT_FILE}"
    echo -e "${GREEN}✓ Generated ${OUTPUT_FILE}${NC}"
done

echo ""
echo -e "${GREEN}✓ All cert-manager manifests generated successfully${NC}"
echo -e "${YELLOW}They are generated from: ${GITHUB_RELEASE_URL}${NC}"
