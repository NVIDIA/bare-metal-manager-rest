# Production Docker Images

This directory contains multi-stage Dockerfiles optimized for production deployments.

## Features

### Multi-Stage Builds
- **Build Stage**: `golang:1.25` - Full Go toolchain for compilation
- **Runtime Stage**: `nvcr.io/nvidia/distroless/go:v3.2.1` - Minimal distroless runtime environment

### Optimizations
- Static compilation with CGO disabled
- Strip debug symbols (`-w -s` flags)
- Minimal distroless base image for reduced attack surface
- Non-root user execution for security
- Proper signal handling

### Security Improvements
- Non-root user (`nonroot` from distroless)
- Minimal distroless image with no shell or package manager
- CA certificates included
- Timezone data for proper logging

## Available Images
1. **Dockerfile.carbide-rest-api** - REST API server
2. **Dockerfile.carbide-rest-db** - Database Migrations
3. **Dockerfile.carbide-rest-workflow** - Workflow Service
4. **Dockerfile.carbide-rest-site-manager** - Site Manager service
5. **Dockerfile.carbide-rest-cert-manager** - Certificate Manager
6. **Dockerfile.carbide-rest-site-agent** - Site Agent

## Building Images

### Build from Repository Root

All Dockerfiles must be built from the repository root as they require access to multiple modules:

```bash
cd /path/to/carbide-rest-api

docker build \
  -f docker/production/Dockerfile.carbide-rest-api \
  -t carbide-rest-api:latest \
  .
```

### Build All Images

```bash
make docker-build
```
