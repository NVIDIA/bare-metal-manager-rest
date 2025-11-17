# GitHub Actions Workflows for Docker Build and Push

This directory contains GitHub Actions workflows for building and pushing Docker images to NVIDIA Container Registry (NVCR).

## Workflows Overview

### 1. `main-build.yml` - Main Build Workflow
The primary workflow that orchestrates the entire build and push process.

**Triggers:**
- Push to `main`, `master`, or `develop` branches (builds and pushes)
- Push to `feature/docker-build-push-nvcr` branch (builds only, does not push)
- Push of version tags (e.g., `v1.2.3`) (builds and pushes)
- Pull requests to `main` or `master` (builds only, does not push)
- Manual trigger via workflow_dispatch

**Build-Only Mode:**
The workflow supports building without pushing to test Dockerfiles before secrets are configured. This happens automatically on:
- Feature branches
- Pull requests
- Any branch that isn't main/master/develop or a tag

### 2. `prepare-build-info.yml` - Build Information Preparation
Reusable workflow that generates version information and build metadata.

**Outputs:**
- `version`: Generated from git describe or tag name
- `short_sha`: 7-character git commit SHA
- `full_sha`: Full git commit SHA
- `target_registry`: Target NVCR registry path

### 3. `build-push-docker.yml` - Docker Build and Push
Reusable workflow that builds and pushes all Docker images.

**Images Built:**
1. `carbide-rest-api`
2. `carbide-rest-db`
3. `carbide-rest-ipam`
4. `carbide-rest-site-manager`
5. `carbide-rest-workflow`
6. `carbide-site-agent`
7. `carbide-rest-cert-manager`

**Note:** This workflow uses production-optimized multi-stage Dockerfiles located in `docker/production/`. These Dockerfiles:
- Build with `golang:1.25`
- Use `alpine:latest` as the runtime base
- Include security hardening (non-root user, minimal packages)
- Produce smaller, production-ready images

## Required Secrets

You must configure the following secrets in your GitHub repository:

### NVCR_USERNAME
- **Description**: NVIDIA Container Registry username
- **Value**: `$oauthtoken` (for NGC API tokens) or your NVCR username
- **Required**: Yes

### NVCR_TOKEN
- **Description**: NVIDIA Container Registry API token
- **Value**: Your NGC API token
- **Required**: Yes

### How to Configure Secrets

1. **Get your NGC API Token:**
   - Go to https://ngc.nvidia.com/
   - Sign in with your NVIDIA account
   - Navigate to "Setup" -> "Generate API Key"
   - Copy your API key

2. **Add secrets to GitHub:**
   - Go to your GitHub repository
   - Click "Settings" -> "Secrets and variables" -> "Actions"
   - Click "New repository secret"
   - Add both secrets:
     - Name: `NVCR_USERNAME`, Value: `$oauthtoken`
     - Name: `NVCR_TOKEN`, Value: `<your-ngc-api-key>`

## Configuration

### Customize Registry Path

Edit `.github/workflows/prepare-build-info.yml` and update the `target_registry` value:

```yaml
# Line ~57
target_registry="nvcr.io/your-org/carbide"
```

Replace `your-org` with your NVCR organization or team name.

### Customize Runner

By default, workflows use `ubuntu-latest`. To use custom runners, edit `main-build.yml`:

```yaml
prepare:
  uses: ./.github/workflows/prepare-build-info.yml
  with:
    runner: your-custom-runner  # Change this

build-and-push:
  uses: ./.github/workflows/build-push-docker.yml
  with:
    runner: your-custom-runner  # Change this
```

### Enable/Disable Push

Images are pushed only when:
- Event is NOT a pull request
- Branch is `main`, `master`, `develop`, or a version tag
- Secrets are properly configured (login step is skipped if secrets are missing)

**Build-Only Mode (No Push):**
The workflow automatically runs in build-only mode on:
- Feature branches (like `feature/docker-build-push-nvcr`)
- Pull requests
- Any branch not in the allowed list

This allows you to:
- Test Dockerfiles before secrets are configured
- Verify builds succeed in CI
- See build logs and results in GitHub Actions

To always build without pushing (dry run), set `push_enabled: false` in `main-build.yml`.

## Docker Image Tagging

Each image is tagged with three tags:

1. **Git SHA Tag**: `<image>:<short-sha>` (e.g., `carbide-rest-api:abc1234`)
   - Primary tag for reproducibility and traceability
   - Always unique per commit

2. **Version Tag**: `<image>:<version>` (e.g., `carbide-rest-api:v1.2.3-5-gabc1234`)
   - Generated from git describe
   - Shows version, commits since tag, and SHA

3. **Latest Tag**: `<image>:latest`
   - Points to the most recent successful build
   - Convenient for development/testing

## Usage Examples

### Trigger on Push
```bash
git add .
git commit -m "Update API"
git push origin main
```

### Trigger on Tag
```bash
git tag -a v1.2.3 -m "Release v1.2.3"
git push origin v1.2.3
```

### Manual Trigger
1. Go to "Actions" tab in GitHub
2. Select "Build and Push All Docker Images"
3. Click "Run workflow"
4. Select branch and click "Run workflow"

## Pull Images

After successful build and push:

```bash
# Log in to NVCR
docker login nvcr.io
# Username: $oauthtoken
# Password: <your-ngc-api-key>

# Pull by SHA (recommended for production)
docker pull nvcr.io/your-org/carbide/carbide-rest-api:abc1234

# Pull by version
docker pull nvcr.io/your-org/carbide/carbide-rest-api:v1.2.3

# Pull latest
docker pull nvcr.io/your-org/carbide/carbide-rest-api:latest
```

## Build Cache

The workflows use GitHub Actions cache to speed up builds:
- Cache is scoped per image
- Uses layer caching via Docker Buildx
- Reduces build times on subsequent runs

## Troubleshooting

### Authentication Failures
- Verify `NVCR_USERNAME` is set to `$oauthtoken`
- Verify `NVCR_TOKEN` contains your valid NGC API key
- Ensure your NGC API key has push permissions

### Build Failures
- Check that all Dockerfiles exist and are valid
- Verify Docker build contexts are correct
- Review build logs in GitHub Actions

### Push Not Happening
- Ensure you're not on a pull request (PRs don't push by default)
- Verify `push_enabled: true` in the workflow
- Check that secrets are configured correctly

### Registry Path Issues
- Update `target_registry` in `prepare-build-info.yml`
- Format should be: `nvcr.io/<org>/<project>`
- No trailing slashes

## Multi-Architecture Support

To add ARM64 support, modify the `platforms` field in `build-push-docker.yml`:

```yaml
platforms: linux/amd64,linux/arm64
```

Note: ARM64 builds may require ARM64 runners or QEMU emulation.

## Advanced: Conditional Builds

To build only specific images based on file changes, you can add path filters:

```yaml
- name: Filter paths
  uses: dorny/paths-filter@v3
  id: filter
  with:
    filters: |
      api:
        - 'carbide-rest-api/**'
      db:
        - 'carbide-rest-db/**'

- name: Build API
  if: steps.filter.outputs.api == 'true'
  # ... build steps
```

## Security Considerations

1. **Secrets Management**: Never commit secrets to the repository
2. **API Token Rotation**: Regularly rotate your NGC API tokens
3. **Access Control**: Limit who can trigger workflows manually
4. **Image Scanning**: Consider adding vulnerability scanning steps
5. **Tag Immutability**: Use SHA tags for production deployments

## Support

For issues or questions:
- Check GitHub Actions logs for detailed error messages
- Review NVCR documentation: https://docs.nvidia.com/ngc/
- Verify Docker build locally before pushing

