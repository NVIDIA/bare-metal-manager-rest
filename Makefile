.PHONY: test postgres-up postgres-down postgres-restart test-clean build build-arm64 build-all docker-build docker-build-multiarch

# Build configuration
BUILD_DIR := build/binaries
IMAGE_REGISTRY := localhost:5000
IMAGE_TAG := latest
DOCKERFILE_DIR := docker/production

# Architecture configuration (default to current platform)
GOARCH ?= $(shell go env GOARCH)
GOOS ?= linux

# PostgreSQL container configuration
POSTGRES_CONTAINER_NAME := project-test
POSTGRES_PORT := 30432
POSTGRES_USER := postgres
POSTGRES_PASSWORD := postgres
POSTGRES_DB := postgres
POSTGRES_IMAGE := postgres:14.4-alpine

# Start PostgreSQL container for testing (idempotent - only starts if not running)
postgres-up:
	@if docker ps -q -f name=$(POSTGRES_CONTAINER_NAME) | grep -q .; then \
		echo "PostgreSQL container already running"; \
	else \
		echo "Starting PostgreSQL container..."; \
		docker run -d --rm \
			--name $(POSTGRES_CONTAINER_NAME) \
			-p $(POSTGRES_PORT):5432 \
			-e POSTGRES_USER=$(POSTGRES_USER) \
			-e POSTGRES_PASSWORD=$(POSTGRES_PASSWORD) \
			-e POSTGRES_DB=$(POSTGRES_DB) \
			$(POSTGRES_IMAGE); \
		echo "PostgreSQL container started on port $(POSTGRES_PORT)"; \
		echo "Waiting for PostgreSQL to be ready..."; \
		until docker exec $(POSTGRES_CONTAINER_NAME) pg_isready -U $(POSTGRES_USER) > /dev/null 2>&1; do \
			printf '.'; \
		done; \
		echo ""; \
		echo "PostgreSQL is ready!"; \
	fi

# Stop PostgreSQL container
postgres-down:
	@echo "Stopping PostgreSQL container..."
	@docker stop $(POSTGRES_CONTAINER_NAME) 2>/dev/null || true
	@echo "PostgreSQL container stopped"

# Restart PostgreSQL container
postgres-restart: postgres-down postgres-up

# Run tests with race detector, sequential execution, and no test caching
# Automatically starts PostgreSQL if needed
test: postgres-up
	DB_NAME=forgetest \
	DB_USER=$(POSTGRES_USER) \
	DB_PASSWORD=$(POSTGRES_PASSWORD) \
	DB_HOST=localhost \
	DB_PORT=$(POSTGRES_PORT) \
	NO_DB_PASSWORD_OK=false \
	TEMPORAL_TLS_ENABLED=false \
	TEMPORAL_SERVER_NAME=test-temporal \
	TEMPORAL_NAMESPACE=test-namespace \
	TEMPORAL_QUEUE=test-queue \
	CGO_ENABLED=1 go test ./... -race -p 1

# Clean test - stops existing container and starts fresh before running tests
test-clean: postgres-down postgres-up
	DB_NAME=forgetest \
	DB_USER=$(POSTGRES_USER) \
	DB_PASSWORD=$(POSTGRES_PASSWORD) \
	DB_HOST=localhost \
	DB_PORT=$(POSTGRES_PORT) \
	NO_DB_PASSWORD_OK=false \
	TEMPORAL_TLS_ENABLED=false \
	TEMPORAL_SERVER_NAME=test-temporal \
	TEMPORAL_NAMESPACE=test-namespace \
	TEMPORAL_QUEUE=test-queue \
	CGO_ENABLED=1 go test ./... -race -p 1 --count=1
	@echo ""
	@echo "Tests completed!"

# Build all Go binaries for the specified architecture (default: current platform)
# Usage: make build [GOARCH=amd64|arm64] [GOOS=linux|darwin]
build:
	@echo "========================================"
	@echo "Building All Go Binaries ($(GOOS)/$(GOARCH))"
	@echo "========================================"
	@echo ""
	@echo "Output: $(BUILD_DIR)/$(GOARCH)"
	@echo ""
	@mkdir -p $(BUILD_DIR)/$(GOARCH)
	@echo "Downloading Go dependencies..."
	@go mod download
	@echo "[SUCCESS] Dependencies downloaded"
	@echo ""
	@echo "Building services..."
	@echo ""
	@echo "Building: api ($(GOARCH))"
	@cd api && CGO_ENABLED=0 GOOS=$(GOOS) GOARCH=$(GOARCH) go build \
		-ldflags "-extldflags '-static'" \
		-o ../$(BUILD_DIR)/$(GOARCH)/api \
		./cmd/api
	@echo "[SUCCESS]"
	@echo ""
	@echo "Building: workflow ($(GOARCH))"
	@cd workflow && CGO_ENABLED=0 GOOS=$(GOOS) GOARCH=$(GOARCH) go build \
		-ldflags "-extldflags '-static'" \
		-o ../$(BUILD_DIR)/$(GOARCH)/workflow \
		./cmd/workflow
	@echo "[SUCCESS]"
	@echo ""
	@echo "Building: site-manager ($(GOARCH))"
	@cd site-manager && CGO_ENABLED=0 GOOS=$(GOOS) GOARCH=$(GOARCH) go build \
		-ldflags "-extldflags '-static'" \
		-o ../$(BUILD_DIR)/$(GOARCH)/sitemgr \
		./cmd/sitemgr
	@echo "[SUCCESS]"
	@echo ""
	@echo "Building: site-agent (elektra) ($(GOARCH))"
	@cd site-agent && CGO_ENABLED=0 GOOS=$(GOOS) GOARCH=$(GOARCH) go build \
		-ldflags "-extldflags '-static'" \
		-o ../$(BUILD_DIR)/$(GOARCH)/elektra \
		./cmd/elektra
	@echo "[SUCCESS]"
	@echo ""
	@echo "Building: db ($(GOARCH))"
	@cd db && CGO_ENABLED=0 GOOS=$(GOOS) GOARCH=$(GOARCH) go build \
		-ldflags "-extldflags '-static'" \
		-o ../$(BUILD_DIR)/$(GOARCH)/migrations \
		./cmd/migrations
	@echo "[SUCCESS]"
	@echo ""
	@echo "Building: cert-manager ($(GOARCH))"
	@cd cert-manager && CGO_ENABLED=0 GOOS=$(GOOS) GOARCH=$(GOARCH) go build \
		-ldflags "-extldflags '-static'" \
		-o ../$(BUILD_DIR)/$(GOARCH)/credsmgr \
		./cmd/credsmgr
	@echo "[SUCCESS]"
	@echo ""
	@echo "========================================"
	@echo "All Binaries Built Successfully! ($(GOARCH))"
	@echo "========================================"
	@echo ""
	@ls -lh $(BUILD_DIR)/$(GOARCH)

# Build binaries for AMD64
build-amd64:
	@$(MAKE) build GOARCH=amd64

# Build binaries for ARM64
build-arm64:
	@$(MAKE) build GOARCH=arm64

# Build binaries for both architectures
build-all:
	@echo "========================================"
	@echo "Building for All Architectures"
	@echo "========================================"
	@$(MAKE) build-amd64
	@echo ""
	@$(MAKE) build-arm64
	@echo ""
	@echo "========================================"
	@echo "Multi-Architecture Build Complete!"
	@echo "========================================"
	@echo ""
	@echo "AMD64 binaries:"
	@ls -lh $(BUILD_DIR)/amd64 2>/dev/null || echo "  (none)"
	@echo ""
	@echo "ARM64 binaries:"
	@ls -lh $(BUILD_DIR)/arm64 2>/dev/null || echo "  (none)"

# Build all Docker images for current platform (production distroless images)
docker-build:
	@echo "========================================"
	@echo "Building Production Docker Images (current platform)"
	@echo "========================================"
	@echo ""
	@echo "Building: carbide-rest-api"
	@docker build \
		-t $(IMAGE_REGISTRY)/carbide-rest-api:$(IMAGE_TAG) \
		-f $(DOCKERFILE_DIR)/Dockerfile.carbide-rest-api \
		.
	@echo "[SUCCESS]"
	@echo ""
	@echo "Building: carbide-rest-workflow"
	@docker build \
		-t $(IMAGE_REGISTRY)/carbide-rest-workflow:$(IMAGE_TAG) \
		-f $(DOCKERFILE_DIR)/Dockerfile.carbide-rest-workflow \
		.
	@echo "[SUCCESS]"
	@echo ""
	@echo "Building: carbide-rest-site-manager"
	@docker build \
		-t $(IMAGE_REGISTRY)/carbide-rest-site-manager:$(IMAGE_TAG) \
		-f $(DOCKERFILE_DIR)/Dockerfile.carbide-rest-site-manager \
		.
	@echo "[SUCCESS]"
	@echo ""
	@echo "Building: carbide-rest-site-agent"
	@docker build \
		-t $(IMAGE_REGISTRY)/carbide-rest-site-agent:$(IMAGE_TAG) \
		-f $(DOCKERFILE_DIR)/Dockerfile.carbide-rest-site-agent \
		.
	@echo "[SUCCESS]"
	@echo ""
	@echo "Building: carbide-rest-db"
	@docker build \
		-t $(IMAGE_REGISTRY)/carbide-rest-db:$(IMAGE_TAG) \
		-f $(DOCKERFILE_DIR)/Dockerfile.carbide-rest-db \
		.
	@echo "[SUCCESS]"
	@echo ""
	@echo "Building: carbide-rest-cert-manager"
	@docker build \
		-t $(IMAGE_REGISTRY)/carbide-rest-cert-manager:$(IMAGE_TAG) \
		-f $(DOCKERFILE_DIR)/Dockerfile.carbide-rest-cert-manager \
		.
	@echo "[SUCCESS]"
	@echo ""
	@echo "========================================"
	@echo "All Production Images Built Successfully!"
	@echo "========================================"
	@echo ""
	@docker images --filter "reference=$(IMAGE_REGISTRY)/*"

# Build multi-architecture Docker images (amd64 + arm64) using Buildx
# Note: Requires Docker Buildx and optionally QEMU for cross-platform emulation
docker-build-multiarch:
	@echo "========================================"
	@echo "Building Multi-Architecture Docker Images"
	@echo "========================================"
	@echo "Platforms: linux/amd64, linux/arm64"
	@echo ""
	@echo "Setting up Docker Buildx..."
	@docker buildx create --name carbide-multiarch --use 2>/dev/null || docker buildx use carbide-multiarch
	@docker buildx inspect --bootstrap
	@echo ""
	@echo "Building: carbide-rest-api (multi-arch)"
	@docker buildx build \
		--platform linux/amd64,linux/arm64 \
		-t $(IMAGE_REGISTRY)/carbide-rest-api:$(IMAGE_TAG) \
		-f $(DOCKERFILE_DIR)/Dockerfile.carbide-rest-api \
		--push \
		.
	@echo "[SUCCESS]"
	@echo ""
	@echo "Building: carbide-rest-workflow (multi-arch)"
	@docker buildx build \
		--platform linux/amd64,linux/arm64 \
		-t $(IMAGE_REGISTRY)/carbide-rest-workflow:$(IMAGE_TAG) \
		-f $(DOCKERFILE_DIR)/Dockerfile.carbide-rest-workflow \
		--push \
		.
	@echo "[SUCCESS]"
	@echo ""
	@echo "Building: carbide-rest-site-manager (multi-arch)"
	@docker buildx build \
		--platform linux/amd64,linux/arm64 \
		-t $(IMAGE_REGISTRY)/carbide-rest-site-manager:$(IMAGE_TAG) \
		-f $(DOCKERFILE_DIR)/Dockerfile.carbide-rest-site-manager \
		--push \
		.
	@echo "[SUCCESS]"
	@echo ""
	@echo "Building: carbide-rest-site-agent (multi-arch)"
	@docker buildx build \
		--platform linux/amd64,linux/arm64 \
		-t $(IMAGE_REGISTRY)/carbide-rest-site-agent:$(IMAGE_TAG) \
		-f $(DOCKERFILE_DIR)/Dockerfile.carbide-rest-site-agent \
		--push \
		.
	@echo "[SUCCESS]"
	@echo ""
	@echo "Building: carbide-rest-db (multi-arch)"
	@docker buildx build \
		--platform linux/amd64,linux/arm64 \
		-t $(IMAGE_REGISTRY)/carbide-rest-db:$(IMAGE_TAG) \
		-f $(DOCKERFILE_DIR)/Dockerfile.carbide-rest-db \
		--push \
		.
	@echo "[SUCCESS]"
	@echo ""
	@echo "Building: carbide-rest-cert-manager (multi-arch)"
	@docker buildx build \
		--platform linux/amd64,linux/arm64 \
		-t $(IMAGE_REGISTRY)/carbide-rest-cert-manager:$(IMAGE_TAG) \
		-f $(DOCKERFILE_DIR)/Dockerfile.carbide-rest-cert-manager \
		--push \
		.
	@echo "[SUCCESS]"
	@echo ""
	@echo "========================================"
	@echo "All Multi-Arch Images Built and Pushed!"
	@echo "========================================"
	@echo ""
	@echo "Images pushed to: $(IMAGE_REGISTRY)"
	@echo "Platforms: linux/amd64, linux/arm64"
