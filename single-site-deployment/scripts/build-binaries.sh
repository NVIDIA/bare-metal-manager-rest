#!/bin/bash

# Build All Go Binaries Script
# Compiles all service binaries outside Docker for faster builds
# Uses shared Go build cache and downloads dependencies only once

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ROOT_DIR="$(cd "$SCRIPT_DIR/../.." && pwd)"
BUILD_DIR="$ROOT_DIR/build/binaries"

echo "========================================"
echo "Building All Go Binaries"
echo "========================================"
echo ""
echo "Root: $ROOT_DIR"
echo "Output: $BUILD_DIR"
echo ""

# Create build directory
mkdir -p "$BUILD_DIR"

# Verify go.mod exists
if [ ! -f "$ROOT_DIR/go.mod" ]; then
    echo "ERROR: go.mod not found in $ROOT_DIR"
    echo "Run ../scripts/create-monorepo.sh first"
    exit 1
fi

cd "$ROOT_DIR"

# Download dependencies once for all services
echo "Downloading Go dependencies..."
go mod download
echo "[SUCCESS] Dependencies downloaded"
echo ""

# Function to build a binary
build_binary() {
    local service_name=$1
    local service_dir=$2
    local cmd_dir=$3
    local binary_name=$4
    local version_package=$5
    
    echo "Building: ${service_name}"
    echo "  Directory: $service_dir"
    echo "  Command: $cmd_dir"
    echo "  Binary: $binary_name"
    
    local output_path="$BUILD_DIR/$binary_name"
    local build_time=$(date +"%Y-%m-%d %H:%M:%S")
    local version="dev"
    
    # Try to read VERSION file if it exists
    if [ -f "$ROOT_DIR/$service_dir/VERSION" ]; then
        version=$(cat "$ROOT_DIR/$service_dir/VERSION")
    fi
    
    # Build with CGO_ENABLED=0 for fully static binaries
    local ldflags="-extldflags '-static'"
    if [ -n "$version_package" ]; then
        ldflags="$ldflags -X '${version_package}.Version=${version}' -X '${version_package}.BuildTime=${build_time}'"
    fi
    
    cd "$ROOT_DIR/$service_dir"
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
        -ldflags "$ldflags" \
        -o "$output_path" \
        "./$cmd_dir"
    
    if [ $? -eq 0 ]; then
        echo "  [SUCCESS]"
        ls -lh "$output_path"
        echo ""
        return 0
    else
        echo "  [FAILED]"
        exit 1
    fi
}

# Build all services
echo "Building services..."
echo ""

build_binary "carbide-rest-api" \
    "carbide-rest-api" \
    "cmd/api" \
    "api" \
    "github.com/nvidia/carbide-rest-api-snapshot/carbide-rest-api/pkg/metadata"

build_binary "carbide-rest-workflow" \
    "carbide-rest-workflow" \
    "cmd/workflow" \
    "workflow" \
    "github.com/nvidia/carbide-rest-api-snapshot/carbide-rest-workflow/pkg/metadata"

build_binary "carbide-rest-site-manager" \
    "carbide-rest-site-manager" \
    "cmd/sitemgr" \
    "sitemgr" \
    "github.com/nvidia/carbide-rest-api-snapshot/carbide-rest-site-manager/pkg/metadata"

build_binary "carbide-site-agent (elektra)" \
    "carbide-site-agent" \
    "cmd/elektra" \
    "elektra" \
    "github.com/nvidia/carbide-rest-api-snapshot/carbide-site-agent/pkg/metadata"

build_binary "carbide-site-agent (elektractl)" \
    "carbide-site-agent" \
    "cmd/elektractl" \
    "elektractl" \
    "github.com/nvidia/carbide-rest-api-snapshot/carbide-site-agent/pkg/metadata"

build_binary "carbide-rest-db" \
    "carbide-rest-db" \
    "cmd/migrations" \
    "migrations" \
    ""

build_binary "carbide-rest-cert-manager" \
    "carbide-rest-cert-manager" \
    "cmd/credsmgr" \
    "credsmgr" \
    "github.com/nvidia/carbide-rest-api-snapshot/carbide-rest-cert-manager/pkg/metadata"

echo ""
echo "========================================"
echo "All Binaries Built Successfully!"
echo "========================================"
echo ""
echo "Binaries available at:"
ls -lh "$BUILD_DIR"
echo ""
echo "Next: Build Docker images"
echo "  ./scripts/build-fast.sh"

