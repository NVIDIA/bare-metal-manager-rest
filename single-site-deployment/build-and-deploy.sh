#!/bin/bash

# Master Build and Deploy Script
# Runs all three steps: rewrite → build → deploy

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

# Color output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

echo -e "${BLUE}╔════════════════════════════════════════╗${NC}"
echo -e "${BLUE}║   Carbide Single-Site HA Deployment   ║${NC}"
echo -e "${BLUE}║     Complete Build and Deploy         ║${NC}"
echo -e "${BLUE}╚════════════════════════════════════════╝${NC}"
echo ""

# Check if parent go.mod exists
if [ ! -f "../go.mod" ]; then
    echo -e "${BLUE}═══════════════════════════════════════════${NC}"
    echo -e "${BLUE}Step 1: Creating Monorepo go.mod${NC}"
    echo -e "${BLUE}═══════════════════════════════════════════${NC}"
    echo ""
    ./scripts/create-monorepo.sh
    echo ""
else
    echo -e "${GREEN}✓ Parent go.mod exists, skipping Step 1${NC}"
    echo ""
fi

echo -e "${BLUE}═══════════════════════════════════════════${NC}"
echo -e "${BLUE}Step 2: Building Docker Images${NC}"
echo -e "${BLUE}═══════════════════════════════════════════${NC}"
echo ""
./scripts/build-simple.sh
echo ""

echo -e "${BLUE}═══════════════════════════════════════════${NC}"
echo -e "${BLUE}Step 3: Deploying to Kind Cluster${NC}"
echo -e "${BLUE}═══════════════════════════════════════════${NC}"
echo ""
./scripts/deploy-kind.sh
echo ""

echo -e "${GREEN}╔════════════════════════════════════════╗${NC}"
echo -e "${GREEN}║          DEPLOYMENT COMPLETE!          ║${NC}"
echo -e "${GREEN}╚════════════════════════════════════════╝${NC}"
echo ""

echo -e "${YELLOW}Quick Status Check:${NC}"
echo ""
kubectl get pods -n carbide-system -o wide
echo ""
kubectl get pods -n carbide-site -o wide
echo ""

echo -e "${GREEN}Access Services:${NC}"
echo "  Cloud API:     kubectl port-forward -n carbide-system svc/cloud-api 8388:8388"
echo "  Site Manager:  kubectl port-forward -n carbide-system svc/site-manager 8100:8100"
echo "  Temporal UI:   kubectl port-forward -n carbide-system svc/temporal-web 8088:8088"
echo ""
echo -e "${GREEN}Check Logs:${NC}"
echo "  kubectl logs -f deployment/cloud-api -n carbide-system"
echo "  kubectl logs -f deployment/site-manager -n carbide-system"  
echo "  kubectl logs -f carbide-site-agent-0 -n carbide-site"
echo ""
echo -e "${BLUE}Happy Coding!${NC}"

