#!/bin/bash

# Build All Images Using Go Workspace
# This approach uses go.work to handle all local dependencies

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ROOT_DIR="$(cd "$SCRIPT_DIR/../.." && pwd)"

# Color output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

# Configuration
IMAGE_REGISTRY="${IMAGE_REGISTRY:-localhost:5000}"
IMAGE_TAG="${IMAGE_TAG:-latest}"

echo -e "${BLUE}========================================${NC}"
echo -e "${BLUE}Carbide Stack - Build with Workspace${NC}"
echo -e "${BLUE}========================================${NC}"
echo ""
echo "Using Go Workspace (go.work) for local dependencies"
echo "Registry: $IMAGE_REGISTRY"
echo "Tag: $IMAGE_TAG"
echo ""

# Check if go.work exists
if [ ! -f "$ROOT_DIR/go.work" ]; then
    echo -e "${RED}ERROR: go.work not found in $ROOT_DIR${NC}"
    echo "Run this first:"
    echo "  go work init"
    echo "  go work use ./carbide-*"
    exit 1
fi

echo -e "${GREEN}Found go.work - workspace mode enabled${NC}"
echo ""

# Build function using multi-stage dockerfile
build_service() {
    local service_name=$1
    local service_dir=$2
    local binary_name=$3
    local cmd_path=${4:-./cmd/api}
    local port=${5:-8080}
    
    local image_name="$IMAGE_REGISTRY/$service_name:$IMAGE_TAG"
    
    echo -e "${YELLOW}Building: $service_name${NC}"
    echo "  Service: $service_dir"
    echo "  Binary: $binary_name"
    echo "  Image: $image_name"
    
    # Create a Dockerfile on-the-fly
    cat > "/tmp/Dockerfile.$service_name" << EOF
FROM golang:1.24-alpine AS builder

RUN apk add --no-cache git gcc musl-dev

WORKDIR /workspace

# Copy the entire workspace (all carbide-* directories and go.work)
COPY go.work go.work.sum* ./
COPY carbide-rest-api ./carbide-rest-api
COPY carbide-rest-api-schema ./carbide-rest-api-schema
COPY carbide-rest-auth ./carbide-rest-auth
COPY carbide-rest-cert-manager ./carbide-rest-cert-manager
COPY carbide-rest-common ./carbide-rest-common
COPY carbide-rest-db ./carbide-rest-db
COPY carbide-rest-ipam ./carbide-rest-ipam
COPY carbide-rest-site-manager ./carbide-rest-site-manager
COPY carbide-rest-workflow ./carbide-rest-workflow
COPY carbide-site-agent ./carbide-site-agent
COPY carbide-site-workflow ./carbide-site-workflow

# Build from service directory
WORKDIR /workspace/$service_dir
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /app/$binary_name $cmd_path

# Final stage
FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /app
COPY --from=builder /app/$binary_name .
EXPOSE $port
CMD ["./$binary_name"]
EOF
    
    # Build
    docker build -t "$image_name" -f "/tmp/Dockerfile.$service_name" "$ROOT_DIR" 2>&1 | tee "/tmp/build-$service_name.log"
    build_exit_code=${PIPESTATUS[0]}
    
    if [ $build_exit_code -eq 0 ]; then
        echo -e "${GREEN}  ✓ SUCCESS${NC}"
        echo ""
        return 0
    else
        echo ""
        echo -e "${RED}  ✗ FAILED (exit code: $build_exit_code)${NC}"
        echo -e "${RED}Last 30 lines:${NC}"
        tail -30 "/tmp/build-$service_name.log"
        echo ""
        exit 1
    fi
}

# Build all services
build_service "carbide-rest-api" "carbide-rest-api" "api" "./cmd/api" "8388"
build_service "carbide-rest-workflow" "carbide-rest-workflow" "workflow" "./cmd/workflow" "9360"
build_service "carbide-rest-site-manager" "carbide-rest-site-manager" "sitemgr" "./cmd/sitemgr" "8100"
build_service "carbide-rest-ipam" "carbide-rest-ipam" "ipam-server" "./cmd/server" "9090"
build_service "carbide-site-agent" "carbide-site-agent" "elektra" "./cmd/elektra" "9360"
build_service "carbide-rest-db" "carbide-rest-db" "migrations" "./cmd/migrations" "8080"

# Cert manager uses a different structure
echo -e "${YELLOW}Building: carbide-rest-cert-manager${NC}"
docker build -t "$IMAGE_REGISTRY/carbide-rest-cert-manager:$IMAGE_TAG" \
    -f "$ROOT_DIR/carbide-rest-cert-manager/docker/Dockerfile.credsmgr" \
    "$ROOT_DIR/carbide-rest-cert-manager" 2>&1 | tee "/tmp/build-carbide-rest-cert-manager.log"
if [ ${PIPESTATUS[0]} -eq 0 ]; then
    echo -e "${GREEN}  ✓ SUCCESS${NC}"
else
    echo -e "${RED}  ✗ FAILED${NC}"
    tail -30 "/tmp/build-carbide-rest-cert-manager.log"
    exit 1
fi

echo ""
echo -e "${GREEN}========================================${NC}"
echo -e "${GREEN}All builds successful!${NC}"
echo -e "${GREEN}========================================${NC}"
echo ""
echo "Next: ./scripts/deploy-kind.sh"

