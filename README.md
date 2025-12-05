# Carbide REST API

A collection of microservices that comprise the management backend for Carbide, exposed as a REST API.

## Architecture

This monorepo contains the following services:

| Service | Description | Binary |
|---------|-------------|--------|
| **carbide-rest-api** | Main REST API server | `api` |
| **carbide-rest-workflow** | Temporal workflow service | `workflow` |
| **carbide-rest-site-manager** | Site management service | `sitemgr` |
| **carbide-site-agent** | On-site agent (elektra) | `elektra` |
| **carbide-rest-db** | Database migrations | `migrations` |
| **carbide-rest-cert-manager** | Certificate/credentials manager | `credsmgr` |

### Supporting Modules

- **carbide-rest-common** - Shared utilities and configurations
- **carbide-rest-auth** - Authentication and authorization
- **carbide-rest-ipam** - IP Address Management
- **carbide-rest-api-schema** - Protocol buffer schemas
- **carbide-site-workflow** - Site-side workflow definitions

## Prerequisites

- Go 1.25.1 or later
- Docker 20.10+ with BuildKit enabled
- Make
- PostgreSQL 14+ (for running tests)

## Building

### Build All Binaries

```bash
make build
```

Binaries are output to `build/binaries/`:

```
build/binaries/
  api          # carbide-rest-api
  workflow     # carbide-rest-workflow
  sitemgr      # carbide-rest-site-manager
  elektra      # carbide-site-agent
  migrations   # carbide-rest-db
  credsmgr     # carbide-rest-cert-manager
```

### Build Docker Images

Build all production Docker images:

```bash
make docker-build
```

By default, images are tagged with `localhost:5000` registry and `latest` tag.

## Using a Private Container Registry

### Configure Registry and Tag

Override the default registry and tag via environment variables or Make arguments:

```bash
# Using Make arguments
make docker-build IMAGE_REGISTRY=my-registry.example.com/carbide IMAGE_TAG=v1.0.0

# Or export environment variables
export IMAGE_REGISTRY=my-registry.example.com/carbide
export IMAGE_TAG=v1.0.0
make docker-build
```

### Build and Push to Private Registry

1. **Authenticate with your registry:**

```bash
# Docker Hub
docker login

# AWS ECR
aws ecr get-login-password --region us-east-1 | docker login --username AWS --password-stdin 123456789.dkr.ecr.us-east-1.amazonaws.com

# Google Container Registry
gcloud auth configure-docker

# Azure Container Registry
az acr login --name myregistry

# Harbor or other private registries
docker login my-registry.example.com
```

2. **Build images with your registry prefix:**

```bash
make docker-build IMAGE_REGISTRY=my-registry.example.com/carbide IMAGE_TAG=v1.0.0
```

3. **Push all images:**

```bash
# Push each image
docker push my-registry.example.com/carbide/carbide-rest-api:v1.0.0
docker push my-registry.example.com/carbide/carbide-rest-workflow:v1.0.0
docker push my-registry.example.com/carbide/carbide-rest-site-manager:v1.0.0
docker push my-registry.example.com/carbide/carbide-site-agent:v1.0.0
docker push my-registry.example.com/carbide/carbide-rest-db:v1.0.0
docker push my-registry.example.com/carbide/carbide-rest-cert-manager:v1.0.0
```

### Quick Build and Push Script

```bash
#!/bin/bash
REGISTRY="${1:-my-registry.example.com/carbide}"
TAG="${2:-latest}"

# Build all images
make docker-build IMAGE_REGISTRY="$REGISTRY" IMAGE_TAG="$TAG"

# Push all images
for image in carbide-rest-api carbide-rest-workflow carbide-rest-site-manager carbide-site-agent carbide-rest-db carbide-rest-cert-manager; do
    docker push "$REGISTRY/$image:$TAG"
done
```

### Available Images

| Image | Port | Description |
|-------|------|-------------|
| `carbide-rest-api` | 8388 | Main REST API |
| `carbide-rest-workflow` | - | Temporal workflow worker |
| `carbide-rest-site-manager` | - | Site management worker |
| `carbide-site-agent` | - | On-site agent |
| `carbide-rest-db` | - | Database migrations (run to completion) |
| `carbide-rest-cert-manager` | - | Certificate manager |

## Running Tests

Tests are currently being ported over from closed source and should not be considered reliable. Some may pass and some may fail while we work through the transition to open source.

Tests require a PostgreSQL database. The Makefile manages a test container automatically.

```bash
# Run tests (auto-starts PostgreSQL if needed)
make test

# Run tests with a fresh database
make test-clean

# Manually manage the test database
make postgres-up      # Start PostgreSQL
make postgres-down    # Stop PostgreSQL
make postgres-restart # Restart PostgreSQL
```

Test database configuration:
- Host: `localhost`
- Port: `30432`
- User: `postgres`
- Password: `postgres`
- Database: `postgres`

## Configuration

Services are configured via environment variables. Common variables:

| Variable | Description | Default |
|----------|-------------|---------|
| `DB_HOST` | PostgreSQL host | - |
| `DB_PORT` | PostgreSQL port | `5432` |
| `DB_NAME` | Database name | - |
| `DB_USER` | Database user | - |
| `DB_PASSWORD` | Database password | - |
| `TEMPORAL_TLS_ENABLED` | Enable Temporal TLS | `true` |
| `TEMPORAL_SERVER_NAME` | Temporal server address | - |
| `TEMPORAL_NAMESPACE` | Temporal namespace | - |
| `TEMPORAL_QUEUE` | Temporal task queue | - |

## Running Containers

### Basic Example

```bash
docker run -p 8388:8388 \
  -e DB_HOST=postgres \
  -e DB_PORT=5432 \
  -e DB_NAME=carbide \
  -e DB_USER=user \
  -e DB_PASSWORD=password \
  my-registry.example.com/carbide/carbide-rest-api:v1.0.0
```

### With Certificate Volumes

```bash
docker run -p 8388:8388 \
  -v /path/to/certs:/var/secrets/temporal/certs:ro \
  -e TEMPORAL_TLS_ENABLED=true \
  my-registry.example.com/carbide/carbide-rest-api:v1.0.0
```

## Docker Image Details

Production images use multi-stage builds with:

- **Build stage:** `golang:1.25.1` for compilation
- **Runtime stage:** `nvcr.io/nvidia/distroless/go:v3.2.1` minimal runtime

Features:
- Static compilation (CGO disabled)
- Debug symbols stripped for smaller size
- Non-root user execution
- No shell or package manager (distroless)

Approximate image sizes: 20-45 MB depending on service.

## Project Structure

```
carbide-rest-api/
  carbide-rest-api/         # Main REST API service
  carbide-rest-workflow/    # Temporal workflow service
  carbide-rest-site-manager/# Site manager service
  carbide-site-agent/       # On-site agent (elektra)
  carbide-rest-db/          # Database migrations
  carbide-rest-cert-manager/# Certificate manager
  carbide-rest-common/      # Shared utilities
  carbide-rest-auth/        # Authentication
  carbide-rest-ipam/        # IP Address Management
  carbide-rest-api-schema/  # Protocol buffers
  carbide-site-workflow/    # Site workflow definitions
  docker/production/        # Production Dockerfiles
  build/binaries/           # Compiled binaries (generated)
```

## License

See [LICENSE](LICENSE) for details.
