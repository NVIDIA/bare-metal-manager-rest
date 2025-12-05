.PHONY: test postgres-up postgres-down postgres-restart test-clean build docker-build

# Build configuration
BUILD_DIR := build/binaries
IMAGE_REGISTRY := localhost:5000
IMAGE_TAG := latest
DOCKERFILE_DIR := docker/production

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

# Build all Go binaries
build:
	@echo "========================================"
	@echo "Building All Go Binaries"
	@echo "========================================"
	@echo ""
	@echo "Output: $(BUILD_DIR)"
	@echo ""
	@mkdir -p $(BUILD_DIR)
	@echo "Downloading Go dependencies..."
	@go mod download
	@echo "[SUCCESS] Dependencies downloaded"
	@echo ""
	@echo "Building services..."
	@echo ""
	@echo "Building: carbide-rest-api"
	@cd carbide-rest-api && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
		-ldflags "-extldflags '-static'" \
		-o ../$(BUILD_DIR)/api \
		./cmd/api
	@echo "[SUCCESS]"
	@echo ""
	@echo "Building: carbide-rest-workflow"
	@cd carbide-rest-workflow && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
		-ldflags "-extldflags '-static'" \
		-o ../$(BUILD_DIR)/workflow \
		./cmd/workflow
	@echo "[SUCCESS]"
	@echo ""
	@echo "Building: carbide-rest-site-manager"
	@cd carbide-rest-site-manager && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
		-ldflags "-extldflags '-static'" \
		-o ../$(BUILD_DIR)/sitemgr \
		./cmd/sitemgr
	@echo "[SUCCESS]"
	@echo ""
	@echo "Building: carbide-rest-site-agent (elektra)"
	@cd carbide-rest-site-agent && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
		-ldflags "-extldflags '-static'" \
		-o ../$(BUILD_DIR)/elektra \
		./cmd/elektra
	@echo "[SUCCESS]"
	@echo ""
	@echo "Building: carbide-rest-site-agent (elektractl)"
	@cd carbide-rest-site-agent && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
		-ldflags "-extldflags '-static'" \
		-o ../$(BUILD_DIR)/elektractl \
		./cmd/elektractl
	@echo "[SUCCESS]"
	@echo ""
	@echo "Building: carbide-rest-db"
	@cd carbide-rest-db && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
		-ldflags "-extldflags '-static'" \
		-o ../$(BUILD_DIR)/migrations \
		./cmd/migrations
	@echo "[SUCCESS]"
	@echo ""
	@echo "Building: carbide-rest-cert-manager"
	@cd carbide-rest-cert-manager && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
		-ldflags "-extldflags '-static'" \
		-o ../$(BUILD_DIR)/credsmgr \
		./cmd/credsmgr
	@echo "[SUCCESS]"
	@echo ""
	@echo "========================================"
	@echo "All Binaries Built Successfully!"
	@echo "========================================"
	@echo ""
	@ls -lh $(BUILD_DIR)

# Build all Docker images (production distroless images)
docker-build:
	@echo "========================================"
	@echo "Building Production Docker Images"
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

