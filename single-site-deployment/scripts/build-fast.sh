#!/bin/bash

# Fast Build Script - Builds binaries once, then creates Docker images
# This is much faster than the traditional multi-stage Docker builds

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ROOT_DIR="$(cd "$SCRIPT_DIR/../.." && pwd)"
DOCKERFILE_DIR="$SCRIPT_DIR/../dockerfiles"

# Configuration
IMAGE_REGISTRY="${IMAGE_REGISTRY:-localhost:5000}"
IMAGE_TAG="${IMAGE_TAG:-latest}"
SKIP_BINARY_BUILD="${SKIP_BINARY_BUILD:-false}"

echo "========================================"
echo "Carbide Stack - Fast Build"
echo "========================================"
echo ""
echo "Registry: $IMAGE_REGISTRY"
echo "Tag: $IMAGE_TAG"
echo "Root: $ROOT_DIR"
echo ""

# Step 1: Build Go binaries (unless skipped)
if [ "$SKIP_BINARY_BUILD" = "false" ]; then
    echo "Step 1: Building Go Binaries"
    echo "----------------------------------------"
    "$SCRIPT_DIR/build-binaries.sh"
    echo ""
else
    echo "Skipping binary build (SKIP_BINARY_BUILD=true)"
    echo ""
fi

# Verify binaries exist
BINARY_DIR="$ROOT_DIR/build/binaries"
if [ ! -d "$BINARY_DIR" ]; then
    echo "ERROR: Binary directory not found: $BINARY_DIR"
    echo "Run without SKIP_BINARY_BUILD=true"
    exit 1
fi

echo "Step 2: Building Docker Images"
echo "----------------------------------------"
echo ""

build_image() {
    local service_name=$1
    local dockerfile=$2
    local image_name="$IMAGE_REGISTRY/$service_name:$IMAGE_TAG"
    
    echo "Building: ${service_name}"
    echo "  Dockerfile: $dockerfile"
    echo "  Image: $image_name"
    
    if [ ! -f "$dockerfile" ]; then
        echo "  ERROR: Dockerfile not found: $dockerfile"
        exit 1
    fi
    
    # Build from ROOT with the service-specific Dockerfile
    # These builds are FAST because they just copy pre-built binaries
    docker build \
        -t "$image_name" \
        -f "$dockerfile" \
        "$ROOT_DIR" \
        2>&1 | tee "/tmp/build-$service_name.log"
    
    build_exit_code=${PIPESTATUS[0]}
    
    if [ $build_exit_code -eq 0 ]; then
        echo "  [SUCCESS]"
        echo ""
        return 0
    else
        echo ""
        echo "  [FAILED] (exit code: $build_exit_code)"
        echo "Last 50 lines of build output:"
        echo "----------------------------------------"
        tail -50 "/tmp/build-$service_name.log"
        echo "----------------------------------------"
        echo ""
        echo "Full log: /tmp/build-$service_name.log"
        exit 1
    fi
}

# First, build the shared base runtime image that all services will use
echo "Building shared base runtime image..."
echo ""
BASE_RUNTIME_IMAGE="$IMAGE_REGISTRY/carbide-base-runtime:$IMAGE_TAG"
BASE_DOCKERFILE="$DOCKERFILE_DIR/Dockerfile.base-runtime"

echo "  Image: $BASE_RUNTIME_IMAGE"
echo "  Dockerfile: $BASE_DOCKERFILE"

if [ ! -f "$BASE_DOCKERFILE" ]; then
    echo "  ERROR: Base runtime Dockerfile not found: $BASE_DOCKERFILE"
    exit 1
fi

docker build \
    -t "$BASE_RUNTIME_IMAGE" \
    -f "$BASE_DOCKERFILE" \
    "$ROOT_DIR"

if [ $? -eq 0 ]; then
    echo "  [SUCCESS]"
    echo ""
else
    echo "  [FAILED]"
    exit 1
fi

# Now build all services using their fast Dockerfiles (which use the base runtime image)
build_image "carbide-rest-api" "$DOCKERFILE_DIR/Dockerfile.carbide-rest-api.fast"
build_image "carbide-rest-workflow" "$DOCKERFILE_DIR/Dockerfile.carbide-rest-workflow.fast"
build_image "carbide-rest-site-manager" "$DOCKERFILE_DIR/Dockerfile.carbide-rest-site-manager.fast"
build_image "carbide-site-agent" "$DOCKERFILE_DIR/Dockerfile.carbide-site-agent.fast"
build_image "carbide-rest-db" "$DOCKERFILE_DIR/Dockerfile.carbide-rest-db.fast"
build_image "carbide-rest-cert-manager" "$DOCKERFILE_DIR/Dockerfile.carbide-rest-cert-manager.fast"

echo ""
echo "========================================"
echo "All Builds Successful!"
echo "========================================"
echo ""
echo "Images built:"
docker images | grep "$IMAGE_REGISTRY" | grep carbide
echo ""
echo "Next: Deploy to kind"
echo "  ./scripts/deploy-kind.sh"
echo ""
echo "Tip: For subsequent builds, if source hasn't changed:"
echo "  SKIP_BINARY_BUILD=true ./scripts/build-fast.sh"

