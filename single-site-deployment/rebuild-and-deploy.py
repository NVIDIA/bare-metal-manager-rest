#!/usr/bin/env python3
"""
Rebuild and Deploy Script
Rebuilds all Carbide Docker images and deploys them to a fresh Kind cluster

Uses fast build mode:
- Compiles all Go binaries outside Docker first
- Creates minimal Docker images using pre-built binaries
- Much faster builds with better caching
- Recommended for development
"""

import subprocess
import sys
import os
import argparse
from pathlib import Path
from typing import List, Tuple

# Configuration
CLUSTER_NAME = "carbide-single-site"
IMAGE_REGISTRY = "localhost:5000"
IMAGE_TAG = "latest"
ROOT_DIR = Path(__file__).parent.parent.absolute()
DOCKERFILE_DIR = Path(__file__).parent / "dockerfiles"
SCRIPTS_DIR = Path(__file__).parent / "scripts"

# Services to build (service_name, dockerfile)
SERVICES = [
    ("carbide-rest-api", "Dockerfile.carbide-rest-api.fast"),
    ("carbide-rest-workflow", "Dockerfile.carbide-rest-workflow.fast"),
    ("carbide-rest-site-manager", "Dockerfile.carbide-rest-site-manager.fast"),
    ("carbide-site-agent", "Dockerfile.carbide-site-agent.fast"),
    ("carbide-rest-db", "Dockerfile.carbide-rest-db.fast"),
    ("carbide-rest-cert-manager", "Dockerfile.carbide-rest-cert-manager.fast"),
]


def print_header(message: str):
    """Print a formatted header"""
    print(f"\n{'='*60}")
    print(f"{message}")
    print(f"{'='*60}\n")


def print_step(step_num: int, message: str):
    """Print a step header"""
    print(f"\nStep {step_num}: {message}")


def print_success(message: str):
    """Print a success message"""
    print(f"[SUCCESS] {message}")


def print_error(message: str):
    """Print an error message"""
    print(f"[ERROR] {message}")


def print_info(message: str):
    """Print an info message"""
    print(f"[INFO] {message}")


def run_command(cmd: List[str], description: str = "", check: bool = True, show_output: bool = False) -> Tuple[bool, str]:
    """Run a command and return success status and output"""
    if description:
        print(f"  {description}...")
    
    try:
        if show_output:
            # Show output in real-time
            result = subprocess.run(cmd, cwd=ROOT_DIR, check=check)
            return True, ""
        else:
            result = subprocess.run(
                cmd,
                cwd=ROOT_DIR,
                capture_output=True,
                text=True,
                check=check
            )
            return True, result.stdout
    except subprocess.CalledProcessError as e:
        print_error(f"Command failed: {' '.join(cmd)}")
        if hasattr(e, 'stderr') and e.stderr:
            print(f"  Error: {e.stderr}")
        return False, e.stderr if hasattr(e, 'stderr') else ""


def check_prerequisites():
    """Check if required tools are installed"""
    print_step(0, "Checking prerequisites")
    
    tools = ["docker", "kind", "kubectl"]
    missing = []
    
    for tool in tools:
        success, _ = run_command(["which", tool], check=False)
        if success:
            print_success(f"{tool} found")
        else:
            print_error(f"{tool} not found")
            missing.append(tool)
    
    if missing:
        print_error(f"Missing tools: {', '.join(missing)}")
        print(f"\nInstall with:")
        for tool in missing:
            print(f"  brew install {tool}")
        return False
    
    return True


def clean_everything(remove_base_image=False):
    """Clean up everything from previous runs"""
    print_step(1, "Cleaning up previous installations")
    
    # Delete ALL kind clusters to be safe
    print_info("Getting list of all Kind clusters...")
    result = subprocess.run(["kind", "get", "clusters"], capture_output=True, text=True, check=False)
    
    if result.returncode == 0 and result.stdout.strip():
        clusters = result.stdout.strip().split('\n')
        print(f"  Found clusters: {', '.join(clusters)}")
        
        # Delete each cluster
        for cluster in clusters:
            if cluster.strip():
                print(f"  Deleting cluster '{cluster.strip()}'...")
                subprocess.run(["kind", "delete", "cluster", "--name", cluster.strip()], check=False)
        
        print_success("All Kind clusters deleted")
    else:
        print_info("No existing Kind clusters found")
    
    # Clean up old Docker images
    print_info("Cleaning up old Docker images...")
    
    # Remove base runtime image (only if rebuilding)
    if remove_base_image:
        base_image_name = f"{IMAGE_REGISTRY}/carbide-base-runtime:{IMAGE_TAG}"
        subprocess.run(
            ["docker", "rmi", "-f", base_image_name],
            capture_output=True,
            check=False
        )
        print_info("  Base runtime image removed (will be rebuilt)")
    else:
        print_info("  Keeping base runtime image (use --rebuild-base to remove)")
    
    # Remove service images
    for service_info in SERVICES:
        service_name = service_info[0]
        image_name = f"{IMAGE_REGISTRY}/{service_name}:{IMAGE_TAG}"
        subprocess.run(
            ["docker", "rmi", "-f", image_name],
            capture_output=True,
            check=False
        )
    
    # Extra cleanup - remove any dangling build cache
    print_info("Pruning Docker build cache...")
    subprocess.run(["docker", "builder", "prune", "-f"], capture_output=True, check=False)
    
    print_success("Cleanup complete")
    return True


def build_go_binaries():
    """Build all Go binaries outside Docker (for fast build mode)"""
    print("\n  Building Go binaries...")
    
    build_script = SCRIPTS_DIR / "build-binaries.sh"
    
    if not build_script.exists():
        print_error(f"Build script not found: {build_script}")
        return False
    
    try:
        result = subprocess.run(
            ["bash", str(build_script)],
            cwd=ROOT_DIR,
            check=True
        )
        print_success("Go binaries built successfully")
        return True
    except subprocess.CalledProcessError as e:
        print_error("Failed to build Go binaries")
        return False


def check_base_image_exists():
    """Check if the base runtime image already exists"""
    base_image_name = f"{IMAGE_REGISTRY}/carbide-base-runtime:{IMAGE_TAG}"
    
    try:
        result = subprocess.run(
            ["docker", "images", "-q", base_image_name],
            capture_output=True,
            text=True,
            check=True
        )
        return bool(result.stdout.strip())
    except subprocess.CalledProcessError:
        return False


def build_base_runtime_image(force_rebuild=False):
    """Build the shared base runtime image (for fast build mode)"""
    base_image_name = f"{IMAGE_REGISTRY}/carbide-base-runtime:{IMAGE_TAG}"
    
    # Check if image exists and we're not forcing a rebuild
    if not force_rebuild and check_base_image_exists():
        print("\n  Base runtime image already exists, skipping rebuild...")
        print(f"    Image: {base_image_name}")
        print_info("    Use --rebuild-base to force rebuild")
        return True
    
    print("\n  Building shared base runtime image...")
    
    dockerfile_path = DOCKERFILE_DIR / "Dockerfile.base-runtime"
    
    if not dockerfile_path.exists():
        print_error(f"Base runtime Dockerfile not found: {dockerfile_path}")
        return False
    
    print(f"    Image: {base_image_name}")
    print(f"    Dockerfile: {dockerfile_path.name}")
    
    cmd = [
        "docker", "build",
        "-t", base_image_name,
        "-f", str(dockerfile_path),
        str(ROOT_DIR)
    ]
    
    try:
        result = subprocess.run(cmd, cwd=ROOT_DIR, check=True)
        print_success("Base runtime image built successfully")
        return True
    except subprocess.CalledProcessError as e:
        print_error("Failed to build base runtime image")
        return False


def build_images(rebuild_base=False):
    """Build all Docker images using fast build mode"""
    print_step(2, "Building Docker images (FAST MODE)")
    
    print_info("Fast build mode: Building Go binaries first, then minimal Docker images")
    print_info("This is much faster as we only download dependencies once!\n")
    
    # Build all Go binaries first
    if not build_go_binaries():
        print_error("Failed to build Go binaries")
        return False
    
    # Build the shared base runtime image
    if not build_base_runtime_image(force_rebuild=rebuild_base):
        print_error("Failed to build base runtime image")
        return False
    
    print("")
    
    failed_services = []
    
    for service_info in SERVICES:
        service_name = service_info[0]
        dockerfile = service_info[1]
        
        image_name = f"{IMAGE_REGISTRY}/{service_name}:{IMAGE_TAG}"
        dockerfile_path = DOCKERFILE_DIR / dockerfile
        
        print(f"\n{'='*60}")
        print(f"Building: {service_name}")
        print(f"  Dockerfile: {dockerfile_path.name}")
        print(f"  Image: {image_name}")
        print(f"{'='*60}")
        
        if not dockerfile_path.exists():
            print_error(f"Dockerfile not found: {dockerfile_path}")
            failed_services.append(service_name)
            continue
        
        cmd = [
            "docker", "build",
            "-t", image_name,
            "-f", str(dockerfile_path),
            str(ROOT_DIR)
        ]
        
        # Run with output shown in real-time
        try:
            result = subprocess.run(cmd, cwd=ROOT_DIR, check=True)
            print_success(f"{service_name} built successfully\n")
        except subprocess.CalledProcessError as e:
            print_error(f"{service_name} build failed")
            failed_services.append(service_name)
    
    if failed_services:
        print_error(f"\nFailed to build: {', '.join(failed_services)}")
        return False
    
    print_success("\nAll images built successfully")
    
    # Show built images
    print("\nBuilt images:")
    subprocess.run(["docker", "images", "--filter", f"reference={IMAGE_REGISTRY}/*"])
    
    return True


def create_cluster():
    """Create a new kind cluster"""
    print_step(3, "Creating Kind cluster")
    
    kind_config = Path(__file__).parent / "kubernetes/base/kind-config.yaml"
    
    cmd = [
        "kind", "create", "cluster",
        "--name", CLUSTER_NAME,
        "--config", str(kind_config)
    ]
    
    success, _ = run_command(cmd, f"Creating cluster '{CLUSTER_NAME}'", show_output=True)
    
    if success:
        print_success("Cluster created successfully")
        return True
    else:
        print_error("Failed to create cluster")
        return False


def load_images():
    """Load Docker images into kind cluster"""
    print_step(4, "Loading images into Kind cluster")
    
    # Note: We skip pre-loading Temporal images due to multi-arch manifest issues
    # on Docker Desktop for Mac. Kind nodes will pull these directly from Docker Hub.
    # Temporal images: temporalio/auto-setup:1.25.1, temporalio/ui:2.31.2
    print("\n  Note: Temporal images will be pulled directly by Kubernetes from Docker Hub")
    print("  (Skipping pre-load due to multi-arch manifest issues with kind on Mac)\n")
    
    # Load the base runtime image first (needed by all services)
    print("\n  Loading base runtime image...")
    base_image_name = f"{IMAGE_REGISTRY}/carbide-base-runtime:{IMAGE_TAG}"
    cmd = [
        "kind", "load", "docker-image",
        base_image_name,
        "--name", CLUSTER_NAME
    ]
    success, _ = run_command(cmd, check=False)
    if not success:
        print_error(f"Failed to load base runtime image: {base_image_name}")
        return False
    
    # Load our built images
    print("\n  Loading Carbide service images...")
    failed_loads = []
    
    for service_info in SERVICES:
        service_name = service_info[0]
        image_name = f"{IMAGE_REGISTRY}/{service_name}:{IMAGE_TAG}"
        
        print(f"    Loading: {image_name}")
        
        cmd = [
            "kind", "load", "docker-image",
            image_name,
            "--name", CLUSTER_NAME
        ]
        
        success, _ = run_command(cmd, check=False)
        
        if not success:
            failed_loads.append(service_name)
    
    if failed_loads:
        print_error(f"Failed to load: {', '.join(failed_loads)}")
        return False
    
    print_success("All images loaded successfully")
    return True


def create_namespaces():
    """Create Kubernetes namespaces"""
    print_step(5, "Creating namespaces")
    
    k8s_dir = Path(__file__).parent / "kubernetes"
    namespaces_file = k8s_dir / "base/namespaces.yaml"
    
    if not namespaces_file.exists():
        print_error(f"Namespaces file not found: {namespaces_file}")
        return False
    
    print("\nCreating namespaces...")
    result = subprocess.run(
        ["kubectl", "apply", "-f", str(namespaces_file)],
        capture_output=True,
        text=True,
        check=False
    )
    
    if result.returncode == 0:
        print_success("Namespaces created successfully")
        return True
    else:
        print_error(f"Failed to create namespaces: {result.stderr}")
        return False


def generate_certificates():
    """Generate development certificates"""
    print_step(6, "Generating development certificates")
    
    cert_script = Path(__file__).parent / "scripts/generate-certs.sh"
    
    if not cert_script.exists():
        print_error(f"Certificate generation script not found: {cert_script}")
        return False
    
    success, _ = run_command(
        ["bash", str(cert_script)],
        "Generating certificates",
        show_output=True
    )
    
    if success:
        print_success("Certificates generated successfully")
        return True
    else:
        print_error("Failed to generate certificates")
        return False


def create_cert_secrets():
    """Create Kubernetes secrets from generated certificates"""
    print("\n  Creating certificate secrets...")
    
    certs_dir = ROOT_DIR / "build/certs"
    
    if not certs_dir.exists():
        print_error(f"Certificates directory not found: {certs_dir}")
        return False
    
    # Read generated OTP
    otp_file = certs_dir / "bootstrap-otp.txt"
    if otp_file.exists():
        with open(otp_file, 'r') as f:
            bootstrap_otp = f.read().strip()
    else:
        bootstrap_otp = "dev-default-otp"
        print_error(f"OTP file not found, using default: {bootstrap_otp}")
    
    # Read CA certificate
    ca_cert_file = certs_dir / "ca-cert.pem"
    if ca_cert_file.exists():
        with open(ca_cert_file, 'r') as f:
            ca_cert = f.read()
    else:
        print_error("CA certificate not found")
        return False
    
    # Read CA private key
    ca_key_file = certs_dir / "ca-key.pem"
    if ca_key_file.exists():
        with open(ca_key_file, 'r') as f:
            ca_key = f.read()
    else:
        print_error("CA private key not found")
        return False
    
    # Update site-agent bootstrap ConfigMap with generated values
    print("  Updating site-agent bootstrap ConfigMap...")
    bootstrap_data = {
        "site-uuid": "00000000-0000-0000-0000-000000000001",
        "otp": bootstrap_otp,
        "creds-url": "https://site-manager.carbide-system.svc.cluster.local:8100/v1/sitecreds",
        "cacert": ca_cert
    }
    
    # Create the ConfigMap using kubectl
    import json
    import tempfile
    
    configmap = {
        "apiVersion": "v1",
        "kind": "ConfigMap",
        "metadata": {
            "name": "site-agent-bootstrap",
            "namespace": "carbide-site"
        },
        "data": bootstrap_data
    }
    
    with tempfile.NamedTemporaryFile(mode='w', suffix='.json', delete=False) as f:
        json.dump(configmap, f)
        temp_file = f.name
    
    try:
        subprocess.run(["kubectl", "apply", "-f", temp_file], check=False)
        print_success("Site-agent bootstrap ConfigMap created")
    finally:
        Path(temp_file).unlink()
    
    # Create vault CA certificate secret for cert-manager
    print("  Creating vault CA certificate secret...")
    subprocess.run([
        "kubectl", "create", "secret", "generic",
        "vault-root-ca-certificate",
        "-n", "carbide-system",
        f"--from-literal=certificate={ca_cert}",
        "--dry-run=client",
        "-o", "yaml"
    ], stdout=subprocess.PIPE, check=False)
    
    result = subprocess.run([
        "kubectl", "create", "secret", "generic",
        "vault-root-ca-certificate",
        "-n", "carbide-system",
        f"--from-literal=certificate={ca_cert}",
        "--dry-run=client",
        "-o", "yaml"
    ], capture_output=True, text=True)
    
    subprocess.run([
        "kubectl", "apply", "-f", "-"
    ], input=result.stdout, text=True, check=False)
    
    # Create vault CA private key secret for cert-manager
    print("  Creating vault CA private key secret...")
    result = subprocess.run([
        "kubectl", "create", "secret", "generic",
        "vault-root-ca-private-key",
        "-n", "carbide-system",
        f"--from-literal=privatekey={ca_key}",
        "--dry-run=client",
        "-o", "yaml"
    ], capture_output=True, text=True)
    
    subprocess.run([
        "kubectl", "apply", "-f", "-"
    ], input=result.stdout, text=True, check=False)
    
    print_success("Vault CA secrets created for cert-manager")
    
    # Create temporal certs secret (empty for now, will be populated by cert-manager)
    print("  Creating temporal certs secret...")
    subprocess.run([
        "kubectl", "create", "secret", "generic",
        "site-agent-temporal-certs",
        "-n", "carbide-site",
        "--from-literal=otp=",
        "--from-literal=ca.crt=",
        "--from-literal=tls.crt=",
        "--from-literal=tls.key=",
        "--dry-run=client",
        "-o", "yaml"
    ], stdout=subprocess.PIPE, check=False)
    
    subprocess.run([
        "kubectl", "apply", "-f", "-"
    ], input=subprocess.run([
        "kubectl", "create", "secret", "generic",
        "site-agent-temporal-certs",
        "-n", "carbide-site",
        "--from-literal=otp=",
        "--from-literal=ca.crt=",
        "--from-literal=tls.crt=",
        "--from-literal=tls.key=",
        "--dry-run=client",
        "-o", "yaml"
    ], capture_output=True, text=True).stdout, text=True, check=False)
    
    print_success("Certificate secrets created")
    return True


def deploy_services():
    """Deploy services to the cluster"""
    print_step(7, "Deploying services to cluster")
    
    k8s_dir = Path(__file__).parent / "kubernetes"
    
    # Deploy CRDs first
    print("\nDeploying Custom Resource Definitions...")
    crd_dir = k8s_dir / "crds"
    if crd_dir.exists():
        subprocess.run(["kubectl", "apply", "-f", str(crd_dir)])
    
    # Deploy databases
    print("\nDeploying databases...")
    subprocess.run(["kubectl", "apply", "-f", str(k8s_dir / "databases/")])
    
    print("\nWaiting for databases to be ready (15s)...")
    subprocess.run(["sleep", "15"])
    subprocess.run([
        "kubectl", "wait",
        "--for=condition=ready",
        "pod", "-l", "app=postgres",
        "-n", "carbide-system",
        "--timeout=60s"
    ], check=False)
    
    # Show postgres status
    print("\nPostgres pods:")
    subprocess.run(["kubectl", "get", "pods", "-n", "carbide-system", "-l", "app=postgres"])
    
    # Create Temporal databases manually (required before Temporal can start)
    print("\nCreating Temporal databases...")
    subprocess.run([
        "kubectl", "exec", "-n", "carbide-system", "postgres-0", "--",
        "psql", "-U", "postgres", "-c",
        "CREATE DATABASE temporal WITH ENCODING 'UTF8';"
    ], check=False, capture_output=True)
    
    subprocess.run([
        "kubectl", "exec", "-n", "carbide-system", "postgres-0", "--",
        "psql", "-U", "postgres", "-c",
        "CREATE DATABASE temporal_visibility WITH ENCODING 'UTF8';"
    ], check=False, capture_output=True)
    
    subprocess.run([
        "kubectl", "exec", "-n", "carbide-system", "postgres-0", "--",
        "psql", "-U", "postgres", "-c",
        "GRANT ALL PRIVILEGES ON DATABASE temporal TO forge;"
    ], check=False, capture_output=True)
    
    subprocess.run([
        "kubectl", "exec", "-n", "carbide-system", "postgres-0", "--",
        "psql", "-U", "postgres", "-c",
        "GRANT ALL PRIVILEGES ON DATABASE temporal_visibility TO forge;"
    ], check=False, capture_output=True)
    
    # Create Keycloak database
    print("\nCreating Keycloak database and user...")
    subprocess.run([
        "kubectl", "exec", "-n", "carbide-system", "postgres-0", "--",
        "psql", "-U", "postgres", "-c",
        "CREATE DATABASE keycloak WITH ENCODING 'UTF8';"
    ], check=False)
    
    subprocess.run([
        "kubectl", "exec", "-n", "carbide-system", "postgres-0", "--",
        "psql", "-U", "postgres", "-c",
        "CREATE USER keycloak WITH PASSWORD 'keycloak';"
    ], check=False)
    
    subprocess.run([
        "kubectl", "exec", "-n", "carbide-system", "postgres-0", "--",
        "psql", "-U", "postgres", "-c",
        "GRANT ALL PRIVILEGES ON DATABASE keycloak TO keycloak;"
    ], check=False)
    
    # Grant schema permissions for PostgreSQL 14+ compatibility
    subprocess.run([
        "kubectl", "exec", "-n", "carbide-system", "postgres-0", "--",
        "psql", "-U", "postgres", "-d", "keycloak", "-c",
        "GRANT ALL ON SCHEMA public TO keycloak;"
    ], check=False)
    
    # Make keycloak the database owner
    subprocess.run([
        "kubectl", "exec", "-n", "carbide-system", "postgres-0", "--",
        "psql", "-U", "postgres", "-c",
        "ALTER DATABASE keycloak OWNER TO keycloak;"
    ], check=False)
    
    print_success("Temporal and Keycloak databases created")
    
    # Run migrations
    print("\nRunning database migrations...")
    subprocess.run(["kubectl", "apply", "-f", str(k8s_dir / "base/migrations-job.yaml")])
    
    print("\nWaiting for migrations to complete...")
    subprocess.run([
        "kubectl", "wait",
        "--for=condition=complete",
        "job/db-migrations",
        "-n", "carbide-system",
        "--timeout=120s"
    ], check=False)
    
    # Show migration logs
    print("\nMigration logs:")
    subprocess.run([
        "kubectl", "logs",
        "job/db-migrations",
        "-n", "carbide-system",
        "--tail=50"
    ], check=False)
    
    # Deploy Temporal
    print("\nDeploying Temporal...")
    subprocess.run(["kubectl", "apply", "-f", str(k8s_dir / "services/temporal.yaml")])
    
    print("\nWaiting for Temporal to be ready...")
    subprocess.run([
        "kubectl", "wait",
        "--for=condition=ready",
        "pod", "-l", "app=temporal",
        "-n", "carbide-system",
        "--timeout=180s"
    ], check=False)
    
    print("\nTemporal pods:")
    subprocess.run(["kubectl", "get", "pods", "-n", "carbide-system", "-l", "app=temporal"])
    
    # Check temporal logs for any errors
    print("\nTemporal startup status:")
    subprocess.run([
        "kubectl", "logs", "-l", "app=temporal",
        "-c", "temporal-server",
        "-n", "carbide-system",
        "--tail=5"
    ], check=False)
    
    # Deploy Keycloak first (auth dependency)
    print("\nDeploying Keycloak...")
    keycloak_file = k8s_dir / "services/keycloak.yaml"
    if keycloak_file.exists():
        subprocess.run(["kubectl", "apply", "-f", str(keycloak_file)], check=False)
        
        print("\nWaiting for Keycloak to be ready...")
        subprocess.run([
            "kubectl", "wait",
            "--for=condition=ready",
            "pod", "-l", "app=keycloak",
            "-n", "carbide-system",
            "--timeout=180s"
        ], check=False)
        
        print("\nKeycloak pods:")
        subprocess.run(["kubectl", "get", "pods", "-n", "carbide-system", "-l", "app=keycloak"])
        
        # Configure Keycloak
        print("\nConfiguring Keycloak...")
        config_script = Path(__file__).parent / "scripts/configure-keycloak.sh"
        if config_script.exists():
            # Port forward Keycloak temporarily for configuration
            import signal
            pf_process = subprocess.Popen([
                "kubectl", "port-forward", "-n", "carbide-system",
                "svc/keycloak", "8080:8080"
            ], stdout=subprocess.DEVNULL, stderr=subprocess.DEVNULL)
            
            try:
                import time
                time.sleep(3)  # Give port-forward time to establish
                
                result = subprocess.run(
                    ["bash", str(config_script)],
                    capture_output=True,
                    text=True
                )
                if result.returncode == 0:
                    print_success("Keycloak configured successfully")
                    print(result.stdout)
                else:
                    print_error(f"Keycloak configuration failed: {result.stderr}")
            finally:
                pf_process.send_signal(signal.SIGTERM)
                pf_process.wait(timeout=5)
        else:
            print_error(f"Keycloak config script not found: {config_script}")
    else:
        print_error(f"Keycloak manifest not found: {keycloak_file}")
    
    # Deploy core services
    print("\nDeploying core services...")
    
    # First deploy services without site-manager
    services_first = [
        "cert-manager.yaml",
        "cloud-api.yaml",
        "workflow.yaml",
    ]
    
    for service in services_first:
        service_file = k8s_dir / "services" / service
        if service_file.exists():
            print(f"\n  Deploying {service}...")
            subprocess.run(["kubectl", "apply", "-f", str(service_file)], check=False)
        else:
            print_error(f"Service file not found: {service}")
    
    # Deploy site-manager RBAC first
    print(f"\n  Deploying site-manager RBAC...")
    site_manager_rbac_file = k8s_dir / "services/site-manager-rbac.yaml"
    if site_manager_rbac_file.exists():
        subprocess.run(["kubectl", "apply", "-f", str(site_manager_rbac_file)], check=False)
        print("  Waiting for RBAC to fully propagate (30s)...")
        print("  (This ensures service account tokens have correct permissions)")
        subprocess.run(["sleep", "30"])
        
        # Verify RBAC is working
        result = subprocess.run([
            "kubectl", "auth", "can-i", "create", "sites.forge.nvidia.com",
            "--as=system:serviceaccount:carbide-system:site-manager",
            "-n", "carbide-system"
        ], capture_output=True, text=True)
        if result.stdout.strip() == "yes":
            print("  RBAC permissions verified - ready to deploy workload")
        else:
            print("  WARNING: RBAC permissions not yet active")
    else:
        print_error(f"Site manager RBAC file not found: {site_manager_rbac_file}")
    
    # Now deploy site-manager workload
    print(f"\n  Deploying site-manager workload...")
    site_manager_file = k8s_dir / "services/site-manager.yaml"
    if site_manager_file.exists():
        subprocess.run(["kubectl", "apply", "-f", str(site_manager_file)], check=False)
    else:
        print_error(f"Site manager file not found: {site_manager_file}")
    
    # Wait for cloud-api to be ready before creating site
    print("\nWaiting for cloud-api to be ready...")
    subprocess.run([
        "kubectl", "wait",
        "--for=condition=ready",
        "pod", "-l", "app=cloud-api",
        "-n", "carbide-system",
        "--timeout=120s"
    ], check=False)
    
    print_info("\nNote: Site creation should be done manually after deployment")
    print_info("  See TESTING_GUIDE.md for instructions on creating and registering a site")
    
    # Deploy site agent
    print("\nDeploying Site Agent...")
    site_agent_file = k8s_dir / "services/site-agent.yaml"
    if site_agent_file.exists():
        subprocess.run(["kubectl", "apply", "-f", str(site_agent_file)], check=False)
    else:
        print_error(f"Site agent file not found: {site_agent_file}")
    
    print_success("\nDeployment commands completed")


def show_status():
    """Show deployment status"""
    print_step(8, "Deployment Status")
    
    print("\nPods in carbide-system:")
    subprocess.run(["kubectl", "get", "pods", "-n", "carbide-system", "-o", "wide"])
    
    print("\nPods in carbide-site:")
    subprocess.run(["kubectl", "get", "pods", "-n", "carbide-site", "-o", "wide"])
    
    print("\nUseful commands:")
    print("  Watch pods:          kubectl get pods -n carbide-system -w")
    print("  Check logs:          kubectl logs -f deployment/cloud-api -n carbide-system")
    print("  Port forward API:    kubectl port-forward -n carbide-system svc/cloud-api 8388:8388")
    print("  Port forward UI:     kubectl port-forward -n carbide-system svc/temporal-web 8088:8088")
    print("  Delete cluster:      kind delete cluster --name carbide-single-site")


def main():
    """Main execution flow"""
    # Parse command-line arguments
    parser = argparse.ArgumentParser(
        description="Rebuild and deploy Carbide services to Kind cluster",
        formatter_class=argparse.RawDescriptionHelpFormatter,
        epilog="""
Build Mode:
  Always uses fast build - builds Go binaries outside Docker first, then creates minimal images.
  This is MUCH faster with better caching.

Base Image Caching:
  By default, the base runtime image is reused if it exists (much faster).
  Use --rebuild-base to force rebuilding the base image.

Examples:
  # Rebuild and deploy (reuses base image if found)
  ./rebuild-and-deploy.py

  # Force rebuild of base image
  ./rebuild-and-deploy.py --rebuild-base

  # Only redeploy (skip build)
  ./rebuild-and-deploy.py --skip-build

  # Update deployment without recreating cluster
  ./rebuild-and-deploy.py --skip-cleanup
        """
    )
    parser.add_argument(
        "--skip-build",
        action="store_true",
        help="Skip Docker image builds and only deploy existing images"
    )
    parser.add_argument(
        "--skip-cleanup",
        action="store_true",
        help="Skip cleanup of existing cluster (useful for updating deployments)"
    )
    parser.add_argument(
        "--rebuild-base",
        action="store_true",
        help="Force rebuild of the base runtime image (by default, reuses existing image if found)"
    )
    args = parser.parse_args()
    
    print_header("Carbide Rebuild and Deploy")
    
    print("Configuration:")
    print(f"  Cluster Name: {CLUSTER_NAME}")
    print(f"  Image Registry: {IMAGE_REGISTRY}")
    print(f"  Image Tag: {IMAGE_TAG}")
    print(f"  Root Directory: {ROOT_DIR}")
    print(f"  Build Mode: FAST")
    print(f"  Skip Build: {args.skip_build}")
    print(f"  Skip Cleanup: {args.skip_cleanup}")
    print(f"  Rebuild Base: {args.rebuild_base}")
    
    # Check prerequisites
    if not check_prerequisites():
        sys.exit(1)
    
    # Clean everything (unless skipped)
    if not args.skip_cleanup:
        if not clean_everything(remove_base_image=args.rebuild_base):
            print_error("Failed to clean up")
            sys.exit(1)
    else:
        print_info("\nSkipping cleanup as requested")
    
    # Build images (unless skipped)
    if not args.skip_build:
        if not build_images(rebuild_base=args.rebuild_base):
            print_error("\nBuild failed! Fix the errors and try again.")
            sys.exit(1)
    else:
        print_info("\nSkipping build as requested")
        print("\nExisting images:")
        subprocess.run(["docker", "images", "--filter", f"reference={IMAGE_REGISTRY}/*"])
    
    # Create cluster (unless it exists and we're skipping cleanup)
    if not args.skip_cleanup:
        if not create_cluster():
            print_error("Failed to create cluster")
            sys.exit(1)
    else:
        print_info("\nUsing existing cluster")
        subprocess.run(["kubectl", "cluster-info", "--context", f"kind-{CLUSTER_NAME}"], check=False)
    
    # Load images (unless skipping build and cleanup)
    if not args.skip_build or not args.skip_cleanup:
        if not load_images():
            print_error("Failed to load images")
            sys.exit(1)
    else:
        print_info("\nSkipping image loading (assuming images already in cluster)")
    
    # Create namespaces (required before certificate secrets)
    if not create_namespaces():
        print_error("Failed to create namespaces")
        sys.exit(1)
    
    # Generate certificates
    if not generate_certificates():
        print_error("Failed to generate certificates")
        sys.exit(1)
    
    # Create certificate secrets
    if not create_cert_secrets():
        print_error("Failed to create certificate secrets")
        sys.exit(1)
    
    # Deploy services
    deploy_services()
    
    # Show status
    show_status()
    
    print(f"\n{'='*60}")
    print(f"Deployment Complete!")
    print(f"{'='*60}\n")
    
    print_info("To run API tests: ./comprehensive-api-test.py")


if __name__ == "__main__":
    main()
