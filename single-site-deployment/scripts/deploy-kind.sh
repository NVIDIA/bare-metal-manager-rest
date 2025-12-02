#!/bin/bash

# Carbide Single-Site Deployment - Deploy to Kind Cluster
# This script deploys the entire Carbide stack to a local kind cluster

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
DEPLOYMENT_DIR="$(cd "$SCRIPT_DIR/.." && pwd)"
K8S_DIR="$DEPLOYMENT_DIR/kubernetes"

# Color output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
CLUSTER_NAME="${CLUSTER_NAME:-carbide-single-site}"
IMAGE_REGISTRY="${IMAGE_REGISTRY:-localhost:5000}"
IMAGE_TAG="${IMAGE_TAG:-latest}"
NAMESPACE_BASE="carbide-system"
NAMESPACE_SITE="carbide-site"

echo -e "${BLUE}========================================${NC}"
echo -e "${BLUE}Carbide Stack - Deploy to Kind${NC}"
echo -e "${BLUE}========================================${NC}"
echo ""
echo "Cluster: $CLUSTER_NAME"
echo "Registry: $IMAGE_REGISTRY"
echo "Tag: $IMAGE_TAG"
echo ""

# Check if kind is installed
if ! command -v kind &> /dev/null; then
    echo -e "${RED}ERROR: kind is not installed${NC}"
    echo "Install with: brew install kind"
    exit 1
fi

# Check if kubectl is installed
if ! command -v kubectl &> /dev/null; then
    echo -e "${RED}ERROR: kubectl is not installed${NC}"
    echo "Install with: brew install kubectl"
    exit 1
fi

# Function to wait for pods
wait_for_pods() {
    local namespace=$1
    local label=$2
    local timeout=${3:-300}
    
    echo -e "${YELLOW}Waiting for pods with label $label in namespace $namespace...${NC}"
    kubectl wait --for=condition=ready pod -l "$label" -n "$namespace" --timeout="${timeout}s" || true
}

# Function to check deployment status
check_deployment() {
    local namespace=$1
    echo ""
    echo -e "${BLUE}Status of namespace: $namespace${NC}"
    kubectl get pods -n "$namespace" -o wide || true
    echo ""
}

# Step 1: Check if cluster exists
echo -e "${YELLOW}Step 1: Checking for existing cluster${NC}"
if kind get clusters | grep -q "^${CLUSTER_NAME}$"; then
    echo -e "${GREEN}Cluster $CLUSTER_NAME exists${NC}"
else
    echo -e "${YELLOW}Cluster not found. Creating new cluster...${NC}"
    kind create cluster --name "$CLUSTER_NAME" --config "$K8S_DIR/base/kind-config.yaml"
fi
echo ""

# Step 2: Load images into kind
echo -e "${YELLOW}Step 2: Loading images into kind cluster${NC}"
echo "This may take a few minutes..."

declare -a IMAGES=(
    "carbide-rest-api"
    "carbide-rest-workflow"
    "carbide-rest-site-manager"
    "carbide-site-agent"
    "carbide-rest-db"
    "carbide-rest-cert-manager"
)

for image in "${IMAGES[@]}"; do
    echo "  Loading: $IMAGE_REGISTRY/$image:$IMAGE_TAG"
    kind load docker-image "$IMAGE_REGISTRY/$image:$IMAGE_TAG" --name "$CLUSTER_NAME" || echo "    (skip if not found)"
done
echo ""

# Step 3: Create namespaces
echo -e "${YELLOW}Step 3: Creating namespaces${NC}"
kubectl apply -f "$K8S_DIR/base/namespaces.yaml"
echo ""

# Step 4: Deploy databases
echo -e "${YELLOW}Step 4: Deploying databases${NC}"
kubectl apply -f "$K8S_DIR/databases/"
echo "Waiting for databases to be ready..."
sleep 10
wait_for_pods "$NAMESPACE_BASE" "app=postgres" 60
check_deployment "$NAMESPACE_BASE"

# Step 5: Run database migrations
echo -e "${YELLOW}Step 5: Running database migrations${NC}"
kubectl apply -f "$K8S_DIR/base/migrations-job.yaml"
echo "Waiting for migrations to complete..."
sleep 5
kubectl wait --for=condition=complete job/db-migrations -n "$NAMESPACE_BASE" --timeout=300s || true
kubectl logs job/db-migrations -n "$NAMESPACE_BASE" || true
echo ""

# Step 6: Deploy Temporal
echo -e "${YELLOW}Step 6: Deploying Temporal${NC}"
kubectl apply -f "$K8S_DIR/services/temporal.yaml"
echo "Waiting for Temporal to be ready..."
sleep 10
wait_for_pods "$NAMESPACE_BASE" "app=temporal" 120
check_deployment "$NAMESPACE_BASE"

# Step 7: Deploy core services
echo -e "${YELLOW}Step 7: Deploying core services${NC}"

echo "  Deploying Cert Manager (Vault)..."
kubectl apply -f "$K8S_DIR/services/cert-manager.yaml"
sleep 5

echo "  Deploying Cloud API..."
kubectl apply -f "$K8S_DIR/services/cloud-api.yaml"
sleep 5

echo "  Deploying Workflow Service..."
kubectl apply -f "$K8S_DIR/services/workflow.yaml"
sleep 5

echo "  Deploying Site Manager..."
kubectl apply -f "$K8S_DIR/services/site-manager.yaml"
sleep 5

echo "Waiting for core services..."
wait_for_pods "$NAMESPACE_BASE" "tier=api" 120
check_deployment "$NAMESPACE_BASE"

# Step 8: Deploy Site Agent
echo -e "${YELLOW}Step 8: Deploying Site Agent${NC}"
kubectl apply -f "$K8S_DIR/services/site-agent.yaml"
echo "Waiting for Site Agent..."
sleep 10
wait_for_pods "$NAMESPACE_SITE" "app=elektra-site-agent" 120
check_deployment "$NAMESPACE_SITE"

# Step 9: Deploy Ingress (optional)
echo -e "${YELLOW}Step 9: Deploying Ingress${NC}"
kubectl apply -f "$K8S_DIR/ingress/" 2>/dev/null || echo "  (ingress optional - skip if not needed)"
echo ""

# Final status
echo ""
echo -e "${BLUE}========================================${NC}"
echo -e "${BLUE}Deployment Complete!${NC}"
echo -e "${BLUE}========================================${NC}"
echo ""

echo -e "${GREEN}Check deployment status:${NC}"
echo ""
echo "All pods:"
echo "  kubectl get pods --all-namespaces"
echo ""
echo "Base services:"
echo "  kubectl get pods -n $NAMESPACE_BASE"
echo ""
echo "Site services:"
echo "  kubectl get pods -n $NAMESPACE_SITE"
echo ""

echo -e "${GREEN}View logs:${NC}"
echo "  kubectl logs -f deployment/cloud-api -n $NAMESPACE_BASE"
echo "  kubectl logs -f deployment/site-manager -n $NAMESPACE_BASE"
echo "  kubectl logs -f statefulset/carbide-site-agent -n $NAMESPACE_SITE"
echo ""

echo -e "${GREEN}Access services:${NC}"
echo "  Cloud API:     kubectl port-forward -n $NAMESPACE_BASE svc/cloud-api 8388:8388"
echo "  Site Manager:  kubectl port-forward -n $NAMESPACE_BASE svc/site-manager 8100:8100"
echo "  Temporal UI:   kubectl port-forward -n $NAMESPACE_BASE svc/temporal-web 8088:8088"
echo ""

echo -e "${GREEN}Test connectivity:${NC}"
echo "  curl -k https://localhost:8100/health  # Site Manager"
echo "  curl http://localhost:8388/health      # Cloud API"
echo ""

echo -e "${YELLOW}Showing current pod status:${NC}"
kubectl get pods -n "$NAMESPACE_BASE"
echo ""
kubectl get pods -n "$NAMESPACE_SITE"
echo ""

echo -e "${BLUE}Deployment script complete!${NC}"

