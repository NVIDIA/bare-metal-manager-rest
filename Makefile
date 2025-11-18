.PHONY: test postgres-up postgres-down postgres-restart test-clean build docker-build clean

# Build configuration
BUILD_DIR := build/binaries
IMAGE_REGISTRY := localhost:5000
IMAGE_TAG := latest
DOCKERFILE_DIR := single-site-deployment/dockerfiles

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
	@echo "Building elektraserver for site agent tests..."
	@cd carbide-site-agent/cmd/elektraserver && go build -race -o ../bin/elektraserver && cd ../../..
	@echo "Starting elektraserver on 127.0.0.1:11079..."
	@./carbide-site-agent/cmd/bin/elektraserver -tout=100 > elektra_server.log 2>&1 & echo $$! > elektra_server.pid
	@echo "Elektraserver started with PID `cat elektra_server.pid`"
	@sleep 1
	@if ps -p `cat elektra_server.pid` > /dev/null 2>&1; then \
		echo "Process is running"; \
	else \
		echo "ERROR: Process failed to start!"; \
		cat elektra_server.log; \
		rm -f elektra_server.pid; \
		exit 1; \
	fi
	@echo "Waiting for elektraserver to be ready..."
	@for i in $$(seq 1 30); do \
		if nc -z 127.0.0.1 11079 2>/dev/null; then \
			echo "Port 11079 is open, waiting for server to fully initialize..."; \
			sleep 3; \
			echo "Elektraserver is ready!"; \
			break; \
		fi; \
		if [ $$i -eq 30 ]; then \
			echo "ERROR: Elektraserver failed to start within 30 seconds"; \
			echo "=== Server logs ==="; \
			cat elektra_server.log 2>/dev/null || echo "No log file"; \
			kill `cat elektra_server.pid` 2>/dev/null || true; \
			rm -f elektra_server.pid elektra_server.log; \
			exit 1; \
		fi; \
		sleep 1; \
	done
	@echo "Running tests..."
	@if DB_NAME=forgetest \
	DB_USER=$(POSTGRES_USER) \
	DB_PASSWORD=$(POSTGRES_PASSWORD) \
	DB_HOST=localhost \
	DB_PORT=$(POSTGRES_PORT) \
	NO_DB_PASSWORD_OK=false \
	CARBIDE_ADDRESS=127.0.0.1:11079 \
	CARBIDE_SEC_OPT=0 \
	TEMPORAL_TLS_ENABLED=false \
	TEMPORAL_SERVER_NAME=test-temporal \
	TEMPORAL_NAMESPACE=test-namespace \
	TEMPORAL_QUEUE=test-queue \
	TEMPORAL_HOST=localhost \
	TEMPORAL_PORT=7233 \
	TEMPORAL_PUBLISH_QUEUE=test-publish-queue \
	TEMPORAL_SUBSCRIBE_QUEUE=test-subscribe-queue \
	TEMPORAL_PUBLISH_NAMESPACE=test-publish-namespace \
	CLUSTER_ID=00000000-0000-0000-0000-000000000000 \
	METRICS_PORT=9090 \
	POD_NAME=test-pod-0 \
	POD_NAMESPACE=default \
	DISABLE_BOOTSTRAP=true \
	CGO_ENABLED=1 go test ./... -race -p 1; then \
		echo "Tests passed"; \
		kill `cat elektra_server.pid` 2>/dev/null || true; \
		rm -f elektra_server.pid elektra_server.log; \
	else \
		echo "Tests failed"; \
		echo "=== Server logs ==="; \
		cat elektra_server.log 2>/dev/null || echo "No log file"; \
		kill `cat elektra_server.pid` 2>/dev/null || true; \
		rm -f elektra_server.pid elektra_server.log; \
		exit 1; \
	fi

# Clean test - stops existing container and starts fresh before running tests
test-clean: postgres-down postgres-up
	@echo "Building elektraserver for site agent tests..."
	@cd carbide-site-agent/cmd/elektraserver && go build -race -o ../bin/elektraserver && cd ../../..
	@echo "Starting elektraserver on 127.0.0.1:11079..."
	@./carbide-site-agent/cmd/bin/elektraserver -tout=100 > elektra_server.log 2>&1 & echo $$! > elektra_server.pid
	@echo "Elektraserver started with PID `cat elektra_server.pid`"
	@sleep 1
	@if ps -p `cat elektra_server.pid` > /dev/null 2>&1; then \
		echo "Process is running"; \
	else \
		echo "ERROR: Process failed to start!"; \
		cat elektra_server.log; \
		rm -f elektra_server.pid; \
		exit 1; \
	fi
	@echo "Waiting for elektraserver to be ready..."
	@for i in $$(seq 1 30); do \
		if nc -z 127.0.0.1 11079 2>/dev/null; then \
			echo "Port 11079 is open, waiting for server to fully initialize..."; \
			sleep 3; \
			echo "Elektraserver is ready!"; \
			break; \
		fi; \
		if [ $$i -eq 30 ]; then \
			echo "ERROR: Elektraserver failed to start within 30 seconds"; \
			echo "=== Server logs ==="; \
			cat elektra_server.log 2>/dev/null || echo "No log file"; \
			kill `cat elektra_server.pid` 2>/dev/null || true; \
			rm -f elektra_server.pid elektra_server.log; \
			exit 1; \
		fi; \
		sleep 1; \
	done
	@echo "Running tests..."
	@if DB_NAME=forgetest \
	DB_USER=$(POSTGRES_USER) \
	DB_PASSWORD=$(POSTGRES_PASSWORD) \
	DB_HOST=localhost \
	DB_PORT=$(POSTGRES_PORT) \
	NO_DB_PASSWORD_OK=false \
	CARBIDE_ADDRESS=127.0.0.1:11079 \
	CARBIDE_SEC_OPT=0 \
	TEMPORAL_TLS_ENABLED=false \
	TEMPORAL_SERVER_NAME=test-temporal \
	TEMPORAL_NAMESPACE=test-namespace \
	TEMPORAL_QUEUE=test-queue \
	TEMPORAL_HOST=localhost \
	TEMPORAL_PORT=7233 \
	TEMPORAL_PUBLISH_QUEUE=test-publish-queue \
	TEMPORAL_SUBSCRIBE_QUEUE=test-subscribe-queue \
	TEMPORAL_PUBLISH_NAMESPACE=test-publish-namespace \
	CLUSTER_ID=00000000-0000-0000-0000-000000000000 \
	METRICS_PORT=9090 \
	POD_NAME=test-pod-0 \
	POD_NAMESPACE=default \
	DISABLE_BOOTSTRAP=true \
	CGO_ENABLED=1 go test ./... -race -p 1 --count=1; then \
		echo ""; \
		echo "Tests completed!"; \
		kill `cat elektra_server.pid` 2>/dev/null || true; \
		rm -f elektra_server.pid elektra_server.log; \
	else \
		echo ""; \
		echo "Tests failed!"; \
		echo "=== Server logs ==="; \
		cat elektra_server.log 2>/dev/null || echo "No log file"; \
		kill `cat elektra_server.pid` 2>/dev/null || true; \
		rm -f elektra_server.pid elektra_server.log; \
		exit 1; \
	fi

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
	@echo "Building: carbide-rest-ipam"
	@cd carbide-rest-ipam && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
		-ldflags "-extldflags '-static'" \
		-o ../$(BUILD_DIR)/ipam-server \
		./cmd/server
	@echo "[SUCCESS]"
	@echo ""
	@echo "Building: carbide-site-agent (elektra)"
	@cd carbide-site-agent && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
		-ldflags "-extldflags '-static'" \
		-o ../$(BUILD_DIR)/elektra \
		./cmd/elektra
	@echo "[SUCCESS]"
	@echo ""
	@echo "Building: carbide-site-agent (elektractl)"
	@cd carbide-site-agent && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
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

# Build all Docker images
docker-build: build
	@echo "========================================"
	@echo "Building Docker Images"
	@echo "========================================"
	@echo ""
	@echo "Building shared base runtime image..."
	@docker build \
		-t $(IMAGE_REGISTRY)/carbide-base-runtime:$(IMAGE_TAG) \
		-f $(DOCKERFILE_DIR)/Dockerfile.base-runtime \
		.
	@echo "[SUCCESS] Base runtime image built"
	@echo ""
	@echo "Building: carbide-rest-api"
	@docker build \
		-t $(IMAGE_REGISTRY)/carbide-rest-api:$(IMAGE_TAG) \
		-f $(DOCKERFILE_DIR)/Dockerfile.carbide-rest-api.fast \
		.
	@echo "[SUCCESS]"
	@echo ""
	@echo "Building: carbide-rest-workflow"
	@docker build \
		-t $(IMAGE_REGISTRY)/carbide-rest-workflow:$(IMAGE_TAG) \
		-f $(DOCKERFILE_DIR)/Dockerfile.carbide-rest-workflow.fast \
		.
	@echo "[SUCCESS]"
	@echo ""
	@echo "Building: carbide-rest-site-manager"
	@docker build \
		-t $(IMAGE_REGISTRY)/carbide-rest-site-manager:$(IMAGE_TAG) \
		-f $(DOCKERFILE_DIR)/Dockerfile.carbide-rest-site-manager.fast \
		.
	@echo "[SUCCESS]"
	@echo ""
	@echo "Building: carbide-rest-ipam"
	@docker build \
		-t $(IMAGE_REGISTRY)/carbide-rest-ipam:$(IMAGE_TAG) \
		-f $(DOCKERFILE_DIR)/Dockerfile.carbide-rest-ipam.fast \
		.
	@echo "[SUCCESS]"
	@echo ""
	@echo "Building: carbide-site-agent"
	@docker build \
		-t $(IMAGE_REGISTRY)/carbide-site-agent:$(IMAGE_TAG) \
		-f $(DOCKERFILE_DIR)/Dockerfile.carbide-site-agent.fast \
		.
	@echo "[SUCCESS]"
	@echo ""
	@echo "Building: carbide-rest-db"
	@docker build \
		-t $(IMAGE_REGISTRY)/carbide-rest-db:$(IMAGE_TAG) \
		-f $(DOCKERFILE_DIR)/Dockerfile.carbide-rest-db.fast \
		.
	@echo "[SUCCESS]"
	@echo ""
	@echo "Building: carbide-rest-cert-manager"
	@docker build \
		-t $(IMAGE_REGISTRY)/carbide-rest-cert-manager:$(IMAGE_TAG) \
		-f $(DOCKERFILE_DIR)/Dockerfile.carbide-rest-cert-manager.fast \
		.
	@echo "[SUCCESS]"
	@echo ""
	@echo "========================================"
	@echo "All Images Built Successfully!"
	@echo "========================================"
	@echo ""
	@docker images --filter "reference=$(IMAGE_REGISTRY)/*"

# Clean up test artifacts and stop any running test servers
clean:
	@echo "Cleaning up test artifacts..."
	@if [ -f elektra_server.pid ]; then \
		kill `cat elektra_server.pid` 2>/dev/null || true; \
		rm -f elektra_server.pid; \
		echo "Stopped elektraserver"; \
	fi
	@rm -f elektra_server.pid
	@echo "Cleanup complete"

