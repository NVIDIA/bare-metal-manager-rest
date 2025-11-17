# Carbide REST API Snapshot

This is a snapshot of the NVIDIA Carbide cloud infrastructure codebase, set up as a monorepo for local development and single-site deployment.

## Repository Structure

This repository contains all Carbide services in a single monorepo:

```
carbide-rest-api-snapshot/
├── go.mod                          # Parent module with replace directives
├── go.sum                          # All dependencies
├── carbide-rest-api/               # Cloud REST API
├── carbide-rest-api-schema/        # API schema definitions
├── carbide-rest-auth/              # Authentication library
├── carbide-rest-cert-manager/      # Certificate management (Vault)
├── carbide-rest-common/            # Shared utilities
├── carbide-rest-db/                # Database layer and migrations
├── carbide-rest-ipam/              # IP Address Management
├── carbide-rest-site-manager/      # Site lifecycle management
├── carbide-rest-workflow/          # Workflow orchestration
├── carbide-site-agent/             # Edge site agent (Elektra)
├── carbide-site-workflow/          # Site workflow library
└── single-site-deployment/         # Complete deployment setup
    ├── scripts/
    │   ├── create-monorepo.sh     # Initialize monorepo
    │   ├── build-simple.sh        # Build all images
    │   └── deploy-kind.sh         # Deploy to kind
    └── kubernetes/                 # K8s manifests
```

## Monorepo Setup

This repository uses a **parent go.mod** at the root that includes all child modules. The key features:

### Parent go.mod
- Declares `module carbide.local/snapshot`
- Contains **replace directives** for all child modules
- Points to local `./carbide-*` directories
- Enables offline builds (no GitLab access needed)

### Child Modules
- Keep their original module names (e.g., `gitlab-master.nvidia.com/nvmetal/cloud-api`)
- Unchanged from upstream
- Resolved via parent replace directives

### Benefits
- ✅ Works completely offline
- ✅ No code changes needed
- ✅ Clean dependency management
- ✅ Easy Docker builds
- ✅ Supports local development

## Quick Start

### Development Setup

```bash
# The root go.mod is already set up
# Just verify it exists
cat go.mod

# For IDE support, ensure Go workspace is recognized
# Most IDEs auto-detect go.mod at root
```

### Building Services

```bash
cd single-site-deployment
./build-and-deploy.sh
```

This runs:
1. Creates/verifies parent go.mod
2. Builds all Docker images
3. Deploys to local kind cluster with HA

Or run steps individually:
```bash
./scripts/create-monorepo.sh    # Ensure parent go.mod exists
./scripts/build-simple.sh        # Build all images
./scripts/deploy-kind.sh         # Deploy to kind
```

### Local Development

```bash
# Make changes in any carbide-* directory
cd carbide-rest-api
# ... edit code ...

# Build just that service
cd ../single-site-deployment
docker build -t localhost:5000/carbide-rest-api:latest \
    -f /tmp/Dockerfile.carbide-rest-api \
    ..

# Or rebuild all
./scripts/build-simple.sh
```

## Documentation

- **single-site-deployment/START_HERE.md** - Quick start guide
- **single-site-deployment/README_COMPLETE.md** - Full reference
- **single-site-deployment/SETUP_STEPS.md** - Step-by-step guide
- **SITE_MANAGER_DETAILED.md** - Site Manager architecture
- **ELEKTRA_SITE_AGENT_DETAILED.md** - Site Agent architecture

## Architecture

### Services
- **Cloud API**: Main REST API (port 8388)
- **Workflow Service**: Temporal workers (port 9360)
- **Site Manager**: Site lifecycle management (port 8100)
- **IPAM**: IP address management (port 9090)
- **Cert Manager**: PKI/Vault service (port 8000)
- **Site Agent (Elektra)**: Edge orchestration (port 9360)
- **Database**: PostgreSQL with 3 databases

### Databases
- **forge**: Cloud control plane data
- **elektra**: Site agent data
- **goipam**: IP management data

### HA Configuration
All services run with 3 replicas except databases and Temporal.

## Prerequisites

- Docker Desktop or Colima
- kind: `brew install kind`
- kubectl: `brew install kubectl`
- Go 1.21+
- 16GB+ RAM

## Contributing

### Making Changes

1. Edit code in any `carbide-*` directory
2. Test locally
3. Rebuild and test: `./single-site-deployment/scripts/build-simple.sh`
4. Commit changes

### Adding Dependencies

```bash
# Add to specific service
cd carbide-rest-api
go get github.com/new/dependency

# Parent go.mod will pick it up automatically
cd ..
go mod tidy
```

### Updating Module Versions

```bash
# Update specific dependency across all modules
cd carbide-rest-api
go get -u github.com/some/package

cd ..
go mod tidy
```

## Important Files

- **go.mod / go.sum**: Parent module (commit these!)
- **single-site-deployment/**: Complete deployment setup
- **.gitignore**: Excludes build artifacts and temp files

## Support

For issues:
1. Check build logs: `/tmp/build-*.log`
2. Check deployment: `kubectl get pods --all-namespaces`
3. Review documentation in `single-site-deployment/`

## License

See individual service LICENSE files and THIRD-PARTY-LICENSES.

