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

echo "Downloading cert-manager manifests from GitHub release..."
echo "Version: ${CERT_MANAGER_VERSION}"

# Check if curl is installed
if ! command -v curl &> /dev/null; then
    echo "Error: curl is not installed"
    echo "Please install curl to download cert-manager manifests"
    exit 1
fi

# Download the cert-manager manifest once to a temporary file
TEMP_FILE=$(mktemp)
trap "rm -f ${TEMP_FILE}" EXIT

echo "Downloading from ${GITHUB_RELEASE_URL}..."
if ! curl -fsSL "${GITHUB_RELEASE_URL}" -o "${TEMP_FILE}"; then
    echo "Error: Failed to download cert-manager manifest"
    echo "URL: ${GITHUB_RELEASE_URL}"
    exit 1
fi

echo "Downloaded successfully"

# Target directories that need cert-manager.yaml
TARGETS=(
    "charts/cert-manager/templates"
    "test/manifests/cert-manager/templates"
)

for target in "${TARGETS[@]}"; do
    TARGET_DIR="${REPO_ROOT}/${target}"
    OUTPUT_FILE="${TARGET_DIR}/cert-manager.yaml"
    
    if [ ! -d "${TARGET_DIR}" ]; then
        echo "Warning: Target directory ${target} does not exist, skipping..."
        continue
    fi
    
    echo "Copying to ${OUTPUT_FILE}..."
    cp "${TEMP_FILE}" "${OUTPUT_FILE}"
    echo "Generated ${OUTPUT_FILE}"
done

echo ""
echo "All cert-manager manifests generated successfully"
echo "Generated from: ${GITHUB_RELEASE_URL}"

