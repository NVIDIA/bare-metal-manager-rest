.PHONY: test postgres-up postgres-down ensure-postgres postgres-wait
.PHONY: build docker-build docker-build-local
.PHONY: test-ipam test-site-agent test-site-manager test-workflow test-db test-api test-auth test-common test-cert-manager test-site-workflow migrate carbide-mock-server-build carbide-mock-server-start carbide-mock-server-stop

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
POSTGRES_DB := forgetest
POSTGRES_IMAGE := postgres:14.4-alpine

postgres-up:
	docker run -d --rm \
		--name $(POSTGRES_CONTAINER_NAME) \
		-p $(POSTGRES_PORT):5432 \
		-e POSTGRES_USER=$(POSTGRES_USER) \
		-e POSTGRES_PASSWORD=$(POSTGRES_PASSWORD) \
		-e POSTGRES_DB=$(POSTGRES_DB) \
		$(POSTGRES_IMAGE)

postgres-down:
	-docker rm -f $(POSTGRES_CONTAINER_NAME)

ensure-postgres:
	@docker inspect $(POSTGRES_CONTAINER_NAME) > /dev/null 2>&1 || $(MAKE) postgres-up
	@$(MAKE) postgres-wait

postgres-wait:
	@until docker exec $(POSTGRES_CONTAINER_NAME) psql -U $(POSTGRES_USER) -d $(POSTGRES_DB) -c "SELECT 1" > /dev/null 2>&1; do :; done

migrate:
	docker exec $(POSTGRES_CONTAINER_NAME) psql -U $(POSTGRES_USER) -d $(POSTGRES_DB) -c "CREATE EXTENSION IF NOT EXISTS pg_trgm;"
	cd db/cmd/migrations && go build -o migrations .
	PGHOST=localhost \
	PGPORT=$(POSTGRES_PORT) \
	PGDATABASE=$(POSTGRES_DB) \
	PGUSER=$(POSTGRES_USER) \
	PGPASSWORD=$(POSTGRES_PASSWORD) \
	./db/cmd/migrations/migrations db init_migrate

test-ipam:
	$(MAKE) ensure-postgres
	cd ipam && go test ./... -count=1

test-site-manager:
	cd site-manager && CGO_ENABLED=1 go test -race -p 1 ./... -count=1

test-workflow:
	$(MAKE) ensure-postgres
	cd workflow && go test -p 1 ./... -count=1

test-db:
	$(MAKE) ensure-postgres
	cd db && go test -p 1 ./... -count=1

carbide-mock-server-build:
	mkdir -p build
	cd site-agent/cmd/elektraserver && go build -o ../../../build/elektraserver .

carbide-mock-server-start: carbide-mock-server-build
	-lsof -ti:11079 | xargs kill -9 2>/dev/null
	./build/elektraserver -tout 0 > build/elektraserver.log 2>&1 & echo $$! > build/elektraserver.pid
	@echo "Waiting for gRPC server to start..."
	@for i in 1 2 3 4 5 6 7 8 9 10; do \
		if grep -q "Started API server" build/elektraserver.log 2>/dev/null; then \
			sleep 0.1; \
			echo "gRPC server is ready"; \
			exit 0; \
		fi; \
		sleep 0.2; \
	done; \
	echo "Timeout waiting for gRPC server to start"; \
	exit 1

carbide-mock-server-stop:
	-kill $$(cat build/elektraserver.pid) 2>/dev/null
	-rm -f build/elektraserver.pid

test-site-agent: carbide-mock-server-start
	cd site-agent/pkg/components && CGO_ENABLED=1 go test -race -p 1 ./... -count=1 ; \
	ret=$$? ; cd ../../.. && $(MAKE) carbide-mock-server-stop ; exit $$ret

test-api:
	$(MAKE) ensure-postgres
	cd api && go test -p 1 ./... -count=1

test-auth:
	$(MAKE) ensure-postgres
	cd auth && go test -p 1 ./... -count=1

test-common:
	cd common && go test -p 1 ./... -count=1

test-cert-manager:
	cd cert-manager && go test -p 1 ./... -count=1

test-site-workflow:
	cd site-workflow && go test -p 1 ./... -count=1

test: test-ipam test-db test-api test-auth test-common test-cert-manager test-site-workflow test-site-manager test-workflow test-site-agent

build:
	mkdir -p $(BUILD_DIR)
	go mod download
	cd api && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-extldflags '-static'" -o ../$(BUILD_DIR)/api ./cmd/api
	cd workflow && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-extldflags '-static'" -o ../$(BUILD_DIR)/workflow ./cmd/workflow
	cd site-manager && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-extldflags '-static'" -o ../$(BUILD_DIR)/sitemgr ./cmd/sitemgr
	cd site-agent && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-extldflags '-static'" -o ../$(BUILD_DIR)/elektra ./cmd/elektra
	cd db && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-extldflags '-static'" -o ../$(BUILD_DIR)/migrations ./cmd/migrations
	cd cert-manager && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-extldflags '-static'" -o ../$(BUILD_DIR)/credsmgr ./cmd/credsmgr

docker-build:
	docker build -t $(IMAGE_REGISTRY)/carbide-rest-api:$(IMAGE_TAG) -f $(DOCKERFILE_DIR)/Dockerfile.carbide-rest-api .
	docker build -t $(IMAGE_REGISTRY)/carbide-rest-workflow:$(IMAGE_TAG) -f $(DOCKERFILE_DIR)/Dockerfile.carbide-rest-workflow .
	docker build -t $(IMAGE_REGISTRY)/carbide-rest-site-manager:$(IMAGE_TAG) -f $(DOCKERFILE_DIR)/Dockerfile.carbide-rest-site-manager .
	docker build -t $(IMAGE_REGISTRY)/carbide-rest-site-agent:$(IMAGE_TAG) -f $(DOCKERFILE_DIR)/Dockerfile.carbide-rest-site-agent .
	docker build -t $(IMAGE_REGISTRY)/carbide-rest-db:$(IMAGE_TAG) -f $(DOCKERFILE_DIR)/Dockerfile.carbide-rest-db .
	docker build -t $(IMAGE_REGISTRY)/carbide-rest-cert-manager:$(IMAGE_TAG) -f $(DOCKERFILE_DIR)/Dockerfile.carbide-rest-cert-manager .

proto:
	if [ -d "carbide-core" ]; then rm -rf carbide-core; fi
	git clone ssh://git@github.com/nvidia/carbide-core.git
	ls carbide-core/rpc/proto
	@for file in carbide-core/rpc/proto/*.proto; do \
		cp "$$file" "workflow-schema/site-agent/workflows/v1/$$(basename "$$file" .proto)_carbide.proto"; \
		echo "Copied: $$file"; \
		./workflow-schema/scripts/add-go-package-option.sh "workflow-schema/site-agent/workflows/v1/$$(basename "$$file" .proto)_carbide.proto"; \
	done
	rm -rf carbide-core

protogen:
	cd workflow-schema
	# Lint is disabled for now
	# echo "Checking validity of proto files"
	# buf lint
	echo "Generating go proto files now"
	buf generate
	cd ..

# =============================================================================
# Kind Local Deployment Targets
# =============================================================================

.PHONY: kind-up kind-down kind-deploy kind-load kind-apply kind-redeploy kind-status kind-logs kind-reset kind-verify setup-site-agent

# Kind cluster configuration
KIND_CLUSTER_NAME := carbide-local
KUSTOMIZE_OVERLAY := deploy/kustomize/overlays/local
LOCAL_DOCKERFILE_DIR := docker/local

# Build images using local Dockerfiles (public base images for local dev)
docker-build-local:
	docker build -t $(IMAGE_REGISTRY)/carbide-rest-api:$(IMAGE_TAG) -f $(LOCAL_DOCKERFILE_DIR)/Dockerfile.carbide-rest-api .
	docker build -t $(IMAGE_REGISTRY)/carbide-rest-workflow:$(IMAGE_TAG) -f $(LOCAL_DOCKERFILE_DIR)/Dockerfile.carbide-rest-workflow .
	docker build -t $(IMAGE_REGISTRY)/carbide-rest-site-manager:$(IMAGE_TAG) -f $(LOCAL_DOCKERFILE_DIR)/Dockerfile.carbide-rest-site-manager .
	docker build -t $(IMAGE_REGISTRY)/carbide-rest-site-agent:$(IMAGE_TAG) -f $(LOCAL_DOCKERFILE_DIR)/Dockerfile.carbide-rest-site-agent .
	docker build -t $(IMAGE_REGISTRY)/carbide-rest-elektraserver:$(IMAGE_TAG) -f $(LOCAL_DOCKERFILE_DIR)/Dockerfile.carbide-rest-elektraserver .
	docker build -t $(IMAGE_REGISTRY)/carbide-rest-db:$(IMAGE_TAG) -f $(LOCAL_DOCKERFILE_DIR)/Dockerfile.carbide-rest-db .
	docker build -t $(IMAGE_REGISTRY)/carbide-rest-cert-manager:$(IMAGE_TAG) -f $(LOCAL_DOCKERFILE_DIR)/Dockerfile.carbide-rest-cert-manager .

# Create kind cluster with port mappings
kind-up:
	kind create cluster --name $(KIND_CLUSTER_NAME) --config deploy/kind/cluster-config.yaml
	kubectl apply -f deploy/kustomize/base/crds/

# Delete kind cluster
kind-down:
	kind delete cluster --name $(KIND_CLUSTER_NAME)

# Full deployment: build images, load into kind, apply manifests
kind-deploy: docker-build-local kind-load kind-apply
	@echo "Deployment complete. Run 'make kind-verify' to verify."

# Load all images into kind cluster
kind-load:
	kind load docker-image $(IMAGE_REGISTRY)/carbide-rest-api:$(IMAGE_TAG) --name $(KIND_CLUSTER_NAME)
	kind load docker-image $(IMAGE_REGISTRY)/carbide-rest-workflow:$(IMAGE_TAG) --name $(KIND_CLUSTER_NAME)
	kind load docker-image $(IMAGE_REGISTRY)/carbide-rest-site-manager:$(IMAGE_TAG) --name $(KIND_CLUSTER_NAME)
	kind load docker-image $(IMAGE_REGISTRY)/carbide-rest-site-agent:$(IMAGE_TAG) --name $(KIND_CLUSTER_NAME)
	kind load docker-image $(IMAGE_REGISTRY)/carbide-rest-elektraserver:$(IMAGE_TAG) --name $(KIND_CLUSTER_NAME)
	kind load docker-image $(IMAGE_REGISTRY)/carbide-rest-db:$(IMAGE_TAG) --name $(KIND_CLUSTER_NAME)
	kind load docker-image $(IMAGE_REGISTRY)/carbide-rest-cert-manager:$(IMAGE_TAG) --name $(KIND_CLUSTER_NAME)

# Apply kustomize manifests and wait for readiness
kind-apply:
	kubectl apply -k $(KUSTOMIZE_OVERLAY)
	@echo "Waiting for PostgreSQL..."
	kubectl -n carbide wait --for=condition=ready pod -l app=postgres --timeout=120s
	@echo "Waiting for Cert Manager (with embedded Vault - PKI auto-initialized)..."
	kubectl -n carbide wait --for=condition=ready pod -l app=carbide-rest-cert-manager --timeout=180s
	@echo "Waiting for Temporal..."
	kubectl -n carbide wait --for=condition=ready pod -l app=temporal --timeout=120s
	@echo "Waiting for Keycloak..."
	kubectl -n carbide wait --for=condition=ready pod -l app=keycloak --timeout=180s
	@echo "Running database migrations..."
	kubectl -n carbide wait --for=condition=complete job/db-migrations --timeout=120s
	@echo "Waiting for API service..."
	kubectl -n carbide wait --for=condition=ready pod -l app=carbide-rest-api --timeout=120s || true
	@echo "Waiting for Site Manager..."
	kubectl -n carbide rollout restart deployment/carbide-rest-site-manager || true
	kubectl -n carbide wait --for=condition=ready pod -l app=carbide-rest-site-manager --timeout=120s || true
	@echo "Waiting for Site Agent..."
	kubectl -n carbide rollout restart deployment/carbide-rest-site-agent || true
	kubectl -n carbide wait --for=condition=ready pod -l app=carbide-rest-site-agent --timeout=120s || true

# Rebuild and redeploy apps only (faster iteration)
kind-redeploy: docker-build-local kind-load
	kubectl -n carbide rollout restart deployment/carbide-rest-api
	kubectl -n carbide rollout restart deployment/carbide-rest-workflow
	kubectl -n carbide rollout restart deployment/carbide-rest-site-agent
	kubectl -n carbide rollout restart deployment/carbide-rest-elektraserver
	kubectl -n carbide rollout restart deployment/carbide-rest-cert-manager
	kubectl -n carbide rollout restart deployment/carbide-rest-site-manager

# Show status of all pods and services
kind-status:
	kubectl -n carbide get pods,svc,jobs

# View logs from API service
kind-logs:
	kubectl -n carbide logs -l app=carbide-rest-api -f --tail=100

# Full reset: tear down cluster, rebuild images, and redeploy everything
kind-reset:
	-kind delete cluster --name $(KIND_CLUSTER_NAME)
	kind create cluster --name $(KIND_CLUSTER_NAME) --config deploy/kind/cluster-config.yaml
	kubectl apply -f deploy/kustomize/base/crds/
	@echo "Installing cert-manager.io..."
	kubectl apply -f https://github.com/cert-manager/cert-manager/releases/download/v1.14.4/cert-manager.yaml
	@echo "Waiting for cert-manager deployments..."
	kubectl -n cert-manager rollout status deployment/cert-manager --timeout=120s
	kubectl -n cert-manager rollout status deployment/cert-manager-webhook --timeout=120s
	kubectl -n cert-manager rollout status deployment/cert-manager-cainjector --timeout=120s
	docker build -t $(IMAGE_REGISTRY)/carbide-rest-api:$(IMAGE_TAG) -f $(LOCAL_DOCKERFILE_DIR)/Dockerfile.carbide-rest-api .
	docker build -t $(IMAGE_REGISTRY)/carbide-rest-workflow:$(IMAGE_TAG) -f $(LOCAL_DOCKERFILE_DIR)/Dockerfile.carbide-rest-workflow .
	docker build -t $(IMAGE_REGISTRY)/carbide-rest-site-manager:$(IMAGE_TAG) -f $(LOCAL_DOCKERFILE_DIR)/Dockerfile.carbide-rest-site-manager .
	docker build -t $(IMAGE_REGISTRY)/carbide-rest-site-agent:$(IMAGE_TAG) -f $(LOCAL_DOCKERFILE_DIR)/Dockerfile.carbide-rest-site-agent .
	docker build -t $(IMAGE_REGISTRY)/carbide-rest-elektraserver:$(IMAGE_TAG) -f $(LOCAL_DOCKERFILE_DIR)/Dockerfile.carbide-rest-elektraserver .
	docker build -t $(IMAGE_REGISTRY)/carbide-rest-db:$(IMAGE_TAG) -f $(LOCAL_DOCKERFILE_DIR)/Dockerfile.carbide-rest-db .
	docker build -t $(IMAGE_REGISTRY)/carbide-rest-cert-manager:$(IMAGE_TAG) -f $(LOCAL_DOCKERFILE_DIR)/Dockerfile.carbide-rest-cert-manager .
	kind load docker-image $(IMAGE_REGISTRY)/carbide-rest-api:$(IMAGE_TAG) --name $(KIND_CLUSTER_NAME)
	kind load docker-image $(IMAGE_REGISTRY)/carbide-rest-workflow:$(IMAGE_TAG) --name $(KIND_CLUSTER_NAME)
	kind load docker-image $(IMAGE_REGISTRY)/carbide-rest-site-manager:$(IMAGE_TAG) --name $(KIND_CLUSTER_NAME)
	kind load docker-image $(IMAGE_REGISTRY)/carbide-rest-site-agent:$(IMAGE_TAG) --name $(KIND_CLUSTER_NAME)
	kind load docker-image $(IMAGE_REGISTRY)/carbide-rest-elektraserver:$(IMAGE_TAG) --name $(KIND_CLUSTER_NAME)
	kind load docker-image $(IMAGE_REGISTRY)/carbide-rest-db:$(IMAGE_TAG) --name $(KIND_CLUSTER_NAME)
	kind load docker-image $(IMAGE_REGISTRY)/carbide-rest-cert-manager:$(IMAGE_TAG) --name $(KIND_CLUSTER_NAME)
	kubectl apply -k $(KUSTOMIZE_OVERLAY)
	kubectl -n carbide wait --for=condition=ready pod -l app=postgres --timeout=120s
	kubectl -n carbide wait --for=condition=ready pod -l app=carbide-rest-cert-manager --timeout=180s
	@echo "Configuring cert-manager.io ClusterIssuer for Vault..."
	kubectl apply -k deploy/kustomize/base/cert-manager-io/
	kubectl -n carbide wait --for=condition=ready pod -l app=temporal --timeout=120s
	kubectl -n carbide wait --for=condition=ready pod -l app=keycloak --timeout=180s
	kubectl -n carbide wait --for=condition=complete job/db-migrations --timeout=120s
	-kubectl -n carbide wait --for=condition=ready pod -l app=carbide-rest-api --timeout=120s
	-kubectl -n carbide wait --for=condition=ready pod -l app=carbide-rest-workflow --timeout=60s
	@echo "Waiting for Site Manager (fetches TLS cert from Vault)..."
	kubectl -n carbide wait --for=condition=ready pod -l app=carbide-rest-site-manager --timeout=180s
	./scripts/setup-local-site-agent.sh
	@echo ""
	@echo "================================================================================"
	@echo "Deployment complete!"
	@echo ""
	@echo "Temporal UI: http://localhost:8233"
	@echo "  - View running workflows and their execution history"
	@echo ""
	@echo "API: http://localhost:8388"
	@echo "Keycloak: http://localhost:8080"
	@echo "Vault: http://localhost:8200"
	@echo "================================================================================"
	@echo ""
	@echo "Example: Get a token and list all sites:"
	@echo ""
	@echo '  TOKEN=$$(curl -s -X POST "http://localhost:8080/realms/carbide-dev/protocol/openid-connect/token" \'
	@echo '    -H "Content-Type: application/x-www-form-urlencoded" \'
	@echo '    -d "client_id=carbide-api" \'
	@echo '    -d "client_secret=carbide-local-secret" \'
	@echo '    -d "grant_type=password" \'
	@echo '    -d "username=admin@example.com" \'
	@echo '    -d "password=adminpassword" | jq -r .access_token)'
	@echo ""
	@echo '  curl -s "http://localhost:8388/v2/org/test-org/carbide/site" \'
	@echo '    -H "Authorization: Bearer $$TOKEN" | jq ".[].name"'
	@echo ""

# Setup site-agent with a real site created via the API
setup-site-agent:
	./scripts/setup-local-site-agent.sh

# Verify the local deployment (health checks)
kind-verify:
	./scripts/verify-local.sh
