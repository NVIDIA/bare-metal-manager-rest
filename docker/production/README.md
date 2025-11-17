# Production Docker Images

This directory contains multi-stage Dockerfiles optimized for production deployments.

## Features

### Multi-Stage Builds
- **Build Stage**: `golang:1.25` - Full Go toolchain for compilation
- **Runtime Stage**: `alpine:latest` - Minimal runtime environment

### Optimizations
- Static compilation with CGO disabled
- Strip debug symbols (`-w -s` flags)
- Minimal base image for reduced attack surface
- Non-root user execution for security
- Health checks included
- Proper signal handling

### Security Improvements
- Non-root user (`appuser:1000`)
- No unnecessary packages
- CA certificates included
- Timezone data for proper logging

## Available Images

1. **Dockerfile.carbide-rest-api** - REST API server
2. **Dockerfile.carbide-rest-db** - Database migrations
3. **Dockerfile.carbide-rest-ipam** - IPAM server
4. **Dockerfile.carbide-rest-site-manager** - Site manager service
5. **Dockerfile.carbide-rest-workflow** - Workflow service
6. **Dockerfile.carbide-rest-cert-manager** - Certificate manager
7. **Dockerfile.carbide-site-agent** - Site agent (elektra)

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

### Build with Private Repository Access

If your dependencies are in private repositories, you need to provide Git credentials:

```bash
docker build \
  --secret id=gitcreds,src=$HOME/.netrc \
  -f docker/production/Dockerfile.carbide-rest-api \
  -t carbide-rest-api:latest \
  .
```

Create a `.netrc` file with:

```
machine gitlab-master.nvidia.com
login <your-username>
password <your-token>
```

### Build All Images

```bash
make build-production-images
```

## Image Sizes

These production images are significantly smaller than development images:

| Image | Approximate Size |
|-------|-----------------|
| carbide-rest-api | ~20-30 MB |
| carbide-rest-db | ~25-35 MB |
| carbide-rest-ipam | ~20-30 MB |
| carbide-rest-site-manager | ~20-30 MB |
| carbide-rest-workflow | ~20-30 MB |
| carbide-site-agent | ~35-45 MB |
| carbide-rest-cert-manager | ~20-30 MB |

## Running Images

### Basic Run

```bash
docker run -p 8388:8388 carbide-rest-api:latest
```

### With Environment Variables

```bash
docker run \
  -e DB_HOST=postgres \
  -e DB_PORT=5432 \
  -e DB_NAME=carbide \
  -p 8388:8388 \
  carbide-rest-api:latest
```

### With Volumes

```bash
docker run \
  -v /path/to/certs:/var/secrets/temporal/certs:ro \
  -p 8388:8388 \
  carbide-rest-api:latest
```

## Health Checks

All images include health checks that run every 30 seconds:

```bash
docker ps
```

The STATUS column will show "healthy" or "unhealthy".

## Troubleshooting

### Build Failures

#### Missing Dependencies
If you see errors about missing modules, ensure you're building from the repository root:

```bash
pwd  # Should be /path/to/carbide-rest-api
docker build -f docker/production/Dockerfile.carbide-rest-api .
```

#### Private Repository Access
If `go mod download` fails:

```bash
docker build --secret id=gitcreds,src=$HOME/.netrc ...
```

#### Version Not Found
If VERSION file is missing, the build will use "dev" as the version.

### Runtime Issues

#### Permission Denied
Images run as non-root user (UID 1000). Ensure mounted volumes have correct permissions:

```bash
chown -R 1000:1000 /path/to/volume
```

#### Health Check Failures
Check if the health endpoint is accessible:

```bash
docker exec <container-id> /app/<binary> health
```

### Image Inspection

View image details:

```bash
docker inspect carbide-rest-api:latest

docker history carbide-rest-api:latest

docker run --rm carbide-rest-api:latest --version
```

## Differences from Development Images

| Feature | Development | Production |
|---------|-------------|------------|
| Base Image | golang:1.24-alpine | golang:1.25 + alpine:latest |
| Build Tools | Included | Build stage only |
| User | root | appuser (UID 1000) |
| Debug Symbols | Included | Stripped (-w -s) |
| CGO | Enabled (some) | Disabled |
| Size | Larger | Minimal |
| Health Checks | Optional | Included |

## GitHub Actions Integration

These Dockerfiles are automatically built and pushed by the GitHub Actions workflow in `.github/workflows/build-push-docker.yml`.

Each commit to main/master triggers:
1. Version generation from git
2. Multi-platform builds
3. Tagging with git SHA, version, and latest
4. Push to NVCR

## Best Practices

1. **Always build from repository root**
2. **Use build secrets for private repos**
3. **Tag with specific versions** (not just latest)
4. **Scan images for vulnerabilities**
5. **Use multi-stage builds** (these Dockerfiles already do)
6. **Run as non-root** (these images already do)
7. **Keep base images updated**

## Version Updates

To update Go version:

1. Change `FROM golang:1.25` to `FROM golang:1.26` (or desired version)
2. Test builds locally
3. Update this README
4. Commit changes

## Security Notes

These images:
- Run as non-root user (UID 1000)
- Have minimal attack surface (alpine:latest)
- Include only necessary CA certificates
- Use static compilation (no dynamic linking vulnerabilities)
- Strip debug symbols to reduce information disclosure

## Support

For issues:
1. Verify you're building from repository root
2. Check Docker version (requires 20.10+)
3. Verify BuildKit is enabled
4. Check .netrc for private repo access
5. Review build logs for specific errors

