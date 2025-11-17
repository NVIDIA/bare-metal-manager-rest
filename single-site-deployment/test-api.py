#!/usr/bin/env python3
"""
Comprehensive API Test Suite for Forge Cloud API
Tests all major API endpoints with optional site registration
"""

import subprocess
import sys
import os
import time
import json
import argparse
import urllib.request
import urllib.parse
import urllib.error
import base64
import ssl
from pathlib import Path
from typing import Optional, Dict, Any, List, Tuple
from dataclasses import dataclass

# Configuration
API_BASE_URL = "http://localhost:8388/v2/org/nvidia/carbide"
KEYCLOAK_URL = "http://localhost:8080"
CLUSTER_NAME = "carbide-single-site"
NAMESPACE_SYSTEM = "carbide-system"
NAMESPACE_SITE = "carbide-site"

@dataclass
class TestContext:
    """Holds test context and created resource IDs"""
    token: str = ""
    infra_provider_id: str = ""
    tenant_id: str = ""
    site_id: str = ""
    site_registration_token: str = ""
    tenant_account_id: str = ""
    allocation_id: str = ""
    ip_block_id: str = ""
    vpc_id: str = ""
    vpc_prefix_id: str = ""
    subnet_id: str = ""
    instance_type_id: str = ""
    operating_system_id: str = ""
    ssh_key_id: str = ""
    ssh_key_group_id: str = ""
    infiniband_partition_id: str = ""
    expected_machine_id: str = ""
    network_security_group_id: str = ""
    instance_id: str = ""
    machine_id: str = ""
    allocation_constraint_id: str = ""

class TestStats:
    """Track test statistics"""
    def __init__(self):
        self.total = 0
        self.passed = 0
        self.failed = 0
        self.skipped = 0
        self.failed_tests: List[str] = []
    
    def print_summary(self):
        print("\n" + "="*70)
        print("TEST SUMMARY")
        print("="*70)
        print(f"Total Tests:  {self.total}")
        print(f"Passed:       {self.passed}")
        print(f"Failed:       {self.failed}")
        print(f"Skipped:      {self.skipped}")
        
        if self.failed_tests:
            print("\nFailed Tests:")
            for test_name in self.failed_tests:
                print(f"  - {test_name}")
        
        print("="*70)
        
        if self.failed > 0:
            print(f"\nResult: FAIL ({self.passed}/{self.total} passed)")
            return False
        else:
            print(f"\nResult: PASS (All {self.passed} tests passed)")
            return True

stats = TestStats()
ctx = TestContext()
logs_already_printed = False

def print_header(message: str):
    """Print a formatted header"""
    print(f"\n{'='*70}")
    print(f"{message}")
    print(f"{'='*70}")

def print_test(message: str):
    """Print a test name"""
    print(f"\n[TEST] {message}")

def print_success(message: str):
    """Print a success message"""
    print(f"  [PASS] {message}")

def print_error(message: str):
    """Print an error message"""
    print(f"  [FAIL] {message}")

def print_info(message: str):
    """Print an info message"""
    print(f"  [INFO] {message}")

def print_warning(message: str):
    """Print a warning message"""
    print(f"  [WARN] {message}")

def print_skip(message: str):
    """Print a skip message"""
    print(f"  [SKIP] {message}")

def record_test_failure(test_name: str):
    """Record a test failure and trigger log collection"""
    stats.failed += 1
    stats.failed_tests.append(test_name)
    print_all_service_logs()

# ============================================================================
# Site Registration Functions
# ============================================================================

def run_command(cmd, capture=True, check=True):
    """Run a shell command"""
    try:
        if capture:
            result = subprocess.run(cmd, capture_output=True, text=True, check=check)
            return result.returncode == 0, result.stdout.strip(), result.stderr.strip()
        else:
            result = subprocess.run(cmd, check=check)
            return result.returncode == 0, "", ""
    except subprocess.CalledProcessError as e:
        if capture:
            return False, e.stdout if e.stdout else "", e.stderr if e.stderr else ""
        return False, "", str(e)

def get_pods_by_label(namespace, app_label):
    """Get list of pods by app label"""
    cmd = ["kubectl", "get", "pods", "-n", namespace, 
           "-l", f"app={app_label}", "-o", "jsonpath={{.items[*].metadata.name}}"]
    success, stdout, stderr = run_command(cmd, check=False)
    
    if not success or not stdout:
        return []
    
    return stdout.split()

def print_service_logs(service_name, namespace):
    """Print logs from all pods of a service"""
    print_info(f"Fetching {service_name} pod logs...")
    
    pods = get_pods_by_label(namespace, service_name)
    if not pods:
        print_warning(f"No {service_name} pods found in namespace {namespace}")
        return
    
    for pod_name in pods:
        print_info(f"Logs from pod: {pod_name}")
        print("="*70)
        
        cmd = ["kubectl", "logs", pod_name, "-n", namespace]
        success, stdout, stderr = run_command(cmd, check=False)
        
        if success:
            print(stdout)
        else:
            print_warning(f"Failed to get logs from {pod_name}: {stderr}")
        
        print("="*70)

def print_all_service_logs():
    """Print logs from all relevant services for debugging"""
    global logs_already_printed
    
    if logs_already_printed:
        print_info("Service logs already printed, skipping duplicate collection")
        return
    
    print_header("COLLECTING SERVICE LOGS FOR DIAGNOSIS")
    
    services = [
        ("cloud-api", NAMESPACE_SYSTEM),
        ("site-manager", NAMESPACE_SYSTEM),
        ("site-workflow", NAMESPACE_SYSTEM),
        ("cloud-workflow", NAMESPACE_SYSTEM),
        ("carbide-site-agent", NAMESPACE_SITE),
    ]
    
    for service_name, namespace in services:
        print_service_logs(service_name, namespace)
    
    logs_already_printed = True

def get_ca_cert():
    """Get the CA certificate"""
    cmd = ["kubectl", "get", "secret", "vault-root-ca-certificate", "-n", "carbide-system", 
           "-o", "jsonpath={.data.ca-cert\\.pem}"]
    success, stdout, stderr = run_command(cmd)
    
    if not success:
        print_error(f"Failed to get CA cert: {stderr}")
        return None
    
    try:
        ca_cert = base64.b64decode(stdout).decode('utf-8')
        return ca_cert
    except Exception as e:
        print_error(f"Failed to decode CA cert: {e}")
        return None

def update_site_agent_config(site_uuid, otp, ca_cert):
    """Update the site-agent bootstrap ConfigMap"""
    configmap_data = {
        "apiVersion": "v1",
        "kind": "ConfigMap",
        "metadata": {
            "name": "site-agent-bootstrap",
            "namespace": "carbide-site"
        },
        "data": {
            "site-uuid": site_uuid,
            "otp": otp,
            "creds-url": "https://site-manager.carbide-system.svc.cluster.local:8100/v1/sitecreds",
            "cacert": ca_cert
        }
    }
    
    cmd = ["kubectl", "apply", "-f", "-"]
    proc = subprocess.Popen(cmd, stdin=subprocess.PIPE, stdout=subprocess.PIPE, 
                          stderr=subprocess.PIPE, text=True)
    stdout, stderr = proc.communicate(input=json.dumps(configmap_data))
    
    if proc.returncode != 0:
        print_error(f"Failed to update ConfigMap: {stderr}")
        return False
    
    return True

def restart_site_agent():
    """Restart the site-agent pod to pick up new config"""
    cmd = ["kubectl", "delete", "pod", "carbide-site-agent-0", "-n", "carbide-site"]
    run_command(cmd, check=False)
    
    # Wait for pod to be ready
    cmd = ["kubectl", "wait", "--for=condition=ready", "pod", "carbide-site-agent-0", 
           "-n", "carbide-site", "--timeout=120s"]
    success, stdout, stderr = run_command(cmd, check=False)
    
    if not success:
        print_warning(f"Pod may not be ready yet: {stderr}")
        return False
    
    return True

def call_register_endpoint(site_uuid):
    """Call the site registration endpoint"""
    # Port forward site-manager
    pf_proc = subprocess.Popen(
        ["kubectl", "port-forward", "-n", "carbide-system", "svc/site-manager", "8100:8100"],
        stdout=subprocess.DEVNULL,
        stderr=subprocess.DEVNULL
    )
    
    time.sleep(3)
    
    try:
        url = f"https://localhost:8100/v1/site/register/{site_uuid}"
        
        # Create SSL context that doesn't verify certificates (for testing)
        ssl_ctx = ssl.create_default_context()
        ssl_ctx.check_hostname = False
        ssl_ctx.verify_mode = ssl.CERT_NONE
        
        req = urllib.request.Request(url, method='POST')
        
        try:
            with urllib.request.urlopen(req, context=ssl_ctx, timeout=30) as response:
                if response.status in [200, 202, 204]:
                    return True
                else:
                    print_error(f"Registration returned status {response.status}")
                    return False
        except urllib.error.HTTPError as e:
            print_error(f"Registration failed: {e.code}")
            return False
        except Exception as e:
            print_error(f"Registration failed: {e}")
            return False
    finally:
        pf_proc.terminate()
        pf_proc.wait()

def wait_for_handshake_via_api(site_uuid, timeout=90):
    """Wait for site handshake to complete by checking API"""
    start_time = time.time()
    
    while time.time() - start_time < timeout:
        success, data = api_request("GET", f"/site/{site_uuid}")
        if success:
            status = data.get("status", "")
            # Site moves through: Pending -> Active/Registered
            # Or we can check if isOnline becomes true
            is_online = data.get("isOnline", False)
            
            if is_online or status == "Registered":
                return True
            
            # Also check if status changed from Pending
            if status not in ["Pending", "Error"]:
                return True
        
        time.sleep(3)
    
    return False

def verify_site_registered_via_api(site_uuid):
    """Verify site is registered by checking via API"""
    success, data = api_request("GET", f"/site/{site_uuid}")
    if success:
        status = data.get("status", "")
        is_online = data.get("isOnline", False)
        
        if status == "Registered" or is_online:
            return True
        else:
            print_warning(f"Site status is '{status}', isOnline: {is_online}")
            return False
    
    print_error("Failed to get site status from API")
    return False

def perform_site_registration(site_uuid):
    """Perform complete site registration process"""
    print_header("Site Registration")
    print_info(f"Registering site: {site_uuid}")
    
    # Check we have the registration token from site creation
    if not ctx.site_registration_token:
        print_error("No registration token available (was site creation successful?)")
        return False
    
    # Step 1: Get CA cert
    print_info("Getting CA certificate...")
    ca_cert = get_ca_cert()
    if not ca_cert:
        print_error("Failed to get CA certificate")
        return False
    print_success("CA certificate retrieved")
    
    # Step 2: Update site-agent config with registration token
    print_info("Updating site-agent bootstrap config...")
    if not update_site_agent_config(site_uuid, ctx.site_registration_token, ca_cert):
        print_error("Failed to update site-agent config")
        return False
    print_success("ConfigMap updated with registration token")
    
    # Step 3: Restart site-agent
    print_info("Restarting site-agent to connect to site...")
    if not restart_site_agent():
        print_warning("Site-agent restart had issues, continuing anyway...")
    else:
        print_success("Site-agent restarted")
    
    # Step 4: Wait for handshake (site-agent downloads certs)
    print_info("Waiting for bootstrap handshake to complete...")
    if not wait_for_handshake_via_api(site_uuid):
        print_warning("Handshake monitoring timed out, continuing anyway...")
    else:
        print_success("Handshake completed")
    
    # Step 5: Call register endpoint to mark site as fully registered
    print_info("Calling registration endpoint...")
    if not call_register_endpoint(site_uuid):
        print_error("Failed to call registration endpoint")
        return False
    print_success("Registration endpoint called")
    
    # Step 6: Verify registration via API
    time.sleep(2)  # Give it a moment to update
    print_info("Verifying site registration status via API...")
    if not verify_site_registered_via_api(site_uuid):
        print_warning("Site registration verification failed")
        return False
    
    print_success(f"Site {site_uuid} is now registered and ready!")
    return True

# ============================================================================
# Port Forward and Auth Functions
# ============================================================================

def restart_site_manager_for_rbac():
    """Restart site-manager pods to ensure they have fresh RBAC tokens"""
    print_info("Restarting site-manager pods for RBAC refresh...")
    
    # Delete site-manager pods
    subprocess.run([
        "kubectl", "delete", "pods", "-n", NAMESPACE_SYSTEM,
        "-l", "app=site-manager"
    ], check=False, capture_output=True)
    
    # Wait for them to come back
    subprocess.run([
        "kubectl", "wait", "--for=condition=ready", "pod",
        "-l", "app=site-manager", "-n", NAMESPACE_SYSTEM,
        "--timeout=120s"
    ], check=False, capture_output=True)
    
    print_success("Site-manager pods restarted with fresh RBAC tokens")

def setup_port_forwards() -> Tuple[subprocess.Popen, subprocess.Popen]:
    """Setup port forwards for Keycloak and API"""
    print_info("Setting up port forwards...")
    
    # Wait for pods
    subprocess.run([
        "kubectl", "wait", "--for=condition=ready", "pod",
        "-l", "app=keycloak", "-n", NAMESPACE_SYSTEM,
        "--timeout=120s"
    ], check=False, capture_output=True)
    
    subprocess.run([
        "kubectl", "wait", "--for=condition=ready", "pod",
        "-l", "app=cloud-api", "-n", NAMESPACE_SYSTEM,
        "--timeout=120s"
    ], check=False, capture_output=True)
    
    # Start port forwards
    keycloak_pf = subprocess.Popen([
        "kubectl", "port-forward", "-n", NAMESPACE_SYSTEM,
        "svc/keycloak", "8080:8080"
    ], stdout=subprocess.DEVNULL, stderr=subprocess.DEVNULL)
    
    api_pf = subprocess.Popen([
        "kubectl", "port-forward", "-n", NAMESPACE_SYSTEM,
        "svc/cloud-api", "8388:8388"
    ], stdout=subprocess.DEVNULL, stderr=subprocess.DEVNULL)
    
    time.sleep(3)
    print_success("Port forwards established")
    
    return keycloak_pf, api_pf

def get_auth_token() -> str:
    """Get authentication token from Keycloak"""
    print_info("Getting authentication token...")
    
    data = urllib.parse.urlencode({
        "username": "testuser",
        "password": "testpass",
        "grant_type": "password",
        "client_id": "carbide-api",
        "client_secret": "carbide-secret-dev-only-do-not-use-in-prod"
    }).encode('utf-8')
    
    try:
        req = urllib.request.Request(
            f"{KEYCLOAK_URL}/realms/carbide/protocol/openid-connect/token",
            data=data,
            headers={"Content-Type": "application/x-www-form-urlencoded"}
        )
        
        with urllib.request.urlopen(req) as response:
            response_data = json.loads(response.read().decode('utf-8'))
            token = response_data.get("access_token")
            print_success("Authentication token obtained")
            return token
    except Exception as e:
        print_error(f"Failed to get token: {e}")
        sys.exit(1)

def api_request(method: str, endpoint: str, data: Optional[Dict] = None, 
                params: Optional[Dict] = None, expect_error: bool = False) -> Tuple[bool, Any]:
    """Make an API request and return success status and response data"""
    url = f"{API_BASE_URL}{endpoint}"
    
    # Add query parameters if provided
    if params:
        url = f"{url}?{urllib.parse.urlencode(params)}"
    
    headers = {
        "Authorization": f"Bearer {ctx.token}",
        "Content-Type": "application/json"
    }
    
    try:
        # Prepare request body
        body = None
        if data is not None:
            body = json.dumps(data).encode('utf-8')
        
        # Create request
        req = urllib.request.Request(url, data=body, headers=headers, method=method)
        
        # Make request
        try:
            with urllib.request.urlopen(req) as response:
                response_code = response.status
                response_text = response.read().decode('utf-8')
                response_data = json.loads(response_text) if response_text else {}
                
                if expect_error:
                    return response_code >= 400, response_data
                
                if response_code in [200, 201, 202, 204]:
                    return True, response_data
                else:
                    print_error(f"API request failed: {method} {endpoint}")
                    print_error(f"Status: {response_code}")
                    print_error(f"Response: {response_text[:200]}")
                    return False, {}
        
        except urllib.error.HTTPError as e:
            response_code = e.code
            response_text = e.read().decode('utf-8') if e.fp else ""
            
            if expect_error:
                response_data = json.loads(response_text) if response_text else {}
                return response_code >= 400, response_data
            
            # For DELETE with 202 or 204, these are success
            if method == "DELETE" and response_code in [202, 204]:
                return True, {}
            
            print_error(f"API request failed: {method} {endpoint}")
            print_error(f"Status: {response_code}")
            print_error(f"Response: {response_text[:200]}")
            return False, {}
    
    except Exception as e:
        print_error(f"Request exception: {e}")
        return False, {}

def test_metadata():
    """Test: GET /metadata"""
    print_test("Metadata API")
    stats.total += 1
    
    success, data = api_request("GET", "/metadata")
    if success and "version" in data:
        print_success(f"Metadata retrieved (version: {data['version']})")
        stats.passed += 1
        return True
    else:
        print_error("Metadata test failed")
        record_test_failure("Metadata API")
        return False

def test_user_current():
    """Test: GET /user/current"""
    print_test("Get Current User")
    stats.total += 1
    
    success, data = api_request("GET", "/user/current")
    if success and "email" in data:
        print_success(f"User info retrieved (email: {data['email']})")
        stats.passed += 1
        return True
    else:
        print_error("User current test failed")
        record_test_failure("Get Current User")
        return False

def test_infrastructure_provider():
    """Test: GET /infrastructure-provider/current"""
    print_test("Infrastructure Provider - Get Current")
    stats.total += 1
    
    success, data = api_request("GET", "/infrastructure-provider/current")
    if success and "id" in data:
        ctx.infra_provider_id = data["id"]
        print_success(f"Infrastructure Provider retrieved (ID: {ctx.infra_provider_id})")
        stats.passed += 1
        return True
    else:
        print_error("Infrastructure Provider test failed")
        record_test_failure("Infrastructure Provider - Get Current")
        return False

def test_infrastructure_provider_stats():
    """Test: GET /infrastructure-provider/current/stats"""
    print_test("Infrastructure Provider Stats")
    stats.total += 1
    
    success, data = api_request("GET", "/infrastructure-provider/current/stats")
    if success:
        print_success(f"Infrastructure Provider stats retrieved")
        print_info(f"Stats: {json.dumps(data, indent=2)}")
        stats.passed += 1
        return True
    else:
        print_error("Infrastructure Provider stats test failed")
        record_test_failure("Infrastructure Provider Stats")
        return False

def test_site_create():
    """Test: POST /site"""
    print_test("Site - Create")
    stats.total += 1
    
    payload = {
        "name": f"test-site-{int(time.time())}",
        "description": "Comprehensive test site",
        "location": {
            "city": "Santa Clara",
            "state": "CA",
            "country": "USA"
        },
        "contact": {
            "email": "test@nvidia.com"
        }
    }
    
    success, data = api_request("POST", "/site", data=payload)
    if success and "id" in data:
        ctx.site_id = data["id"]
        ctx.site_registration_token = data.get('registrationToken', '')
        print_success(f"Site created (ID: {ctx.site_id})")
        if ctx.site_registration_token:
            print_info(f"Registration Token: {ctx.site_registration_token[:20]}...")
        else:
            print_info("Registration Token: None")
        stats.passed += 1
        return True
    else:
        print_error("Site create test failed")
        record_test_failure("Site - Create")
        return False

def test_site_get():
    """Test: GET /site/{siteId}"""
    print_test("Site - Get by ID")
    stats.total += 1
    
    if not ctx.site_id:
        print_skip("No site ID available")
        stats.skipped += 1
        return False
    
    params = {"infrastructureProviderId": ctx.infra_provider_id}
    success, data = api_request("GET", f"/site/{ctx.site_id}", params=params)
    if success and data.get("id") == ctx.site_id:
        print_success(f"Site retrieved (name: {data.get('name')})")
        stats.passed += 1
        return True
    else:
        print_error("Site get test failed")
        record_test_failure("Site - Get by ID")
        return False

def test_site_list():
    """Test: GET /site"""
    print_test("Site - List All")
    stats.total += 1
    
    params = {"infrastructureProviderId": ctx.infra_provider_id}
    success, data = api_request("GET", "/site", params=params)
    if success:
        count = len(data) if isinstance(data, list) else 0
        print_success(f"Sites listed ({count} sites)")
        stats.passed += 1
        return True
    else:
        print_error("Site list test failed")
        record_test_failure("Site - List All")
        return False

def test_site_update():
    """Test: PATCH /site/{siteId}"""
    print_test("Site - Update")
    stats.total += 1
    
    if not ctx.site_id:
        print_skip("No site ID available")
        stats.skipped += 1
        return False
    
    payload = {"description": "Updated comprehensive test site description"}
    params = {"infrastructureProviderId": ctx.infra_provider_id}
    success, data = api_request("PATCH", f"/site/{ctx.site_id}", data=payload, params=params)
    if success:
        print_success("Site updated successfully")
        stats.passed += 1
        return True
    else:
        print_error("Site update test failed")
        record_test_failure("Site - Update")
        return False

def test_site_status_history():
    """Test: GET /site/{siteId}/status-history"""
    print_test("Site - Status History")
    stats.total += 1
    
    if not ctx.site_id:
        print_skip("No site ID available")
        stats.skipped += 1
        return False
    
    params = {"infrastructureProviderId": ctx.infra_provider_id}
    success, data = api_request("GET", f"/site/{ctx.site_id}/status-history", params=params)
    if success:
        count = len(data) if isinstance(data, list) else 0
        print_success(f"Status history retrieved ({count} entries)")
        stats.passed += 1
        return True
    else:
        print_error("Site status history test failed")
        record_test_failure("Site - Status History")
        return False

def test_ip_block_create():
    """Test: POST /ipblock"""
    print_test("IP Block - Create")
    stats.total += 1
    
    if not ctx.site_id:
        print_skip("No site ID available")
        stats.skipped += 1
        return False
    
    payload = {
        "name": f"test-ipblock-{int(time.time())}",
        "description": "Test IP block for API testing",
        "siteId": ctx.site_id,
        "infrastructureProviderId": ctx.infra_provider_id,
        "routingType": "Public",
        "prefix": "10.100.0.0",
        "prefixLength": 16,
        "protocolVersion": "IPv4"
    }
    
    success, data = api_request("POST", "/ipblock", data=payload)
    if success and "id" in data:
        ctx.ip_block_id = data["id"]
        print_success(f"IP Block created (ID: {ctx.ip_block_id})")
        stats.passed += 1
        return True
    else:
        print_error("IP Block create test failed")
        record_test_failure("IP Block - Create")
        return False

def test_ip_block_get():
    """Test: GET /ipblock/{ipBlockId}"""
    print_test("IP Block - Get by ID")
    stats.total += 1
    
    if not ctx.ip_block_id:
        print_skip("No IP block ID available")
        stats.skipped += 1
        return False
    
    params = {"infrastructureProviderId": ctx.infra_provider_id}
    success, data = api_request("GET", f"/ipblock/{ctx.ip_block_id}", params=params)
    if success and data.get("id") == ctx.ip_block_id:
        print_success(f"IP Block retrieved (prefix: {data.get('prefix')}/{data.get('prefixLength')})")
        stats.passed += 1
        return True
    else:
        print_error("IP Block get test failed")
        record_test_failure("IP Block - Get by ID")
        return False

def test_ip_block_list():
    """Test: GET /ipblock"""
    print_test("IP Block - List All")
    stats.total += 1
    
    params = {"infrastructureProviderId": ctx.infra_provider_id}
    success, data = api_request("GET", "/ipblock", params=params)
    if success:
        count = len(data) if isinstance(data, list) else 0
        print_success(f"IP Blocks listed ({count} blocks)")
        stats.passed += 1
        return True
    else:
        print_error("IP Block list test failed")
        record_test_failure("IP Block - List All")
        return False

def test_ip_block_update():
    """Test: PATCH /ipblock/{ipBlockId}"""
    print_test("IP Block - Update")
    stats.total += 1
    
    if not ctx.ip_block_id:
        print_skip("No IP block ID available")
        stats.skipped += 1
        return False
    
    payload = {"description": "Updated test IP block description"}
    success, data = api_request("PATCH", f"/ipblock/{ctx.ip_block_id}", data=payload)
    if success:
        print_success("IP Block updated successfully")
        stats.passed += 1
        return True
    else:
        print_error("IP Block update test failed")
        record_test_failure("IP Block - Update")
        return False

def test_instance_type_create():
    """Test: POST /instance/type"""
    print_test("Instance Type - Create")
    stats.total += 1
    
    if not ctx.site_id:
        print_skip("No site ID available")
        stats.skipped += 1
        return False
    
    payload = {
        "name": f"test-type-{int(time.time())}",
        "description": "Test instance type for API testing",
        "siteId": ctx.site_id,
        "machineCapabilities": [
            {
                "type": "CPU",
                "name": "Intel Xeon",
                "count": 2
            }
        ]
    }
    
    success, data = api_request("POST", "/instance/type", data=payload)
    if success and "id" in data:
        ctx.instance_type_id = data["id"]
        print_success(f"Instance Type created (ID: {ctx.instance_type_id})")
        stats.passed += 1
        return True
    else:
        print_error("Instance Type create test failed")
        record_test_failure("Instance Type - Create")
        return False

def test_instance_type_get():
    """Test: GET /instance/type/{instanceTypeId}"""
    print_test("Instance Type - Get by ID")
    stats.total += 1
    
    if not ctx.instance_type_id:
        print_skip("No instance type ID available")
        stats.skipped += 1
        return False
    
    success, data = api_request("GET", f"/instance/type/{ctx.instance_type_id}")
    if success and data.get("id") == ctx.instance_type_id:
        print_success(f"Instance Type retrieved (name: {data.get('name')})")
        stats.passed += 1
        return True
    else:
        print_error("Instance Type get test failed")
        record_test_failure("Instance Type - Get by ID")
        return False

def test_instance_type_list():
    """Test: GET /instance/type"""
    print_test("Instance Type - List All")
    stats.total += 1
    
    if not ctx.site_id:
        print_skip("No site ID available")
        stats.skipped += 1
        return False
    
    params = {
        "siteId": ctx.site_id,
        "infrastructureProviderId": ctx.infra_provider_id
    }
    success, data = api_request("GET", "/instance/type", params=params)
    if success:
        count = len(data) if isinstance(data, list) else 0
        print_success(f"Instance Types listed ({count} types)")
        stats.passed += 1
        return True
    else:
        print_error("Instance Type list test failed")
        record_test_failure("Instance Type - List All")
        return False

def test_instance_type_update():
    """Test: PATCH /instance/type/{instanceTypeId}"""
    print_test("Instance Type - Update")
    stats.total += 1
    
    if not ctx.instance_type_id:
        print_skip("No instance type ID available")
        stats.skipped += 1
        return False
    
    payload = {"description": "Updated test instance type description"}
    success, data = api_request("PATCH", f"/instance/type/{ctx.instance_type_id}", data=payload)
    if success:
        print_success("Instance Type updated successfully")
        stats.passed += 1
        return True
    else:
        print_error("Instance Type update test failed")
        record_test_failure("Instance Type - Update")
        return False

def test_machine_list():
    """Test: GET /machine"""
    print_test("Machine - List All")
    stats.total += 1
    
    if not ctx.site_id:
        print_skip("No site ID available")
        stats.skipped += 1
        return False
    
    params = {"siteId": ctx.site_id}
    success, data = api_request("GET", "/machine", params=params)
    if success:
        count = len(data) if isinstance(data, list) else 0
        print_success(f"Machines listed ({count} machines)")
        
        # Store first machine ID if available
        if isinstance(data, list) and len(data) > 0:
            ctx.machine_id = data[0].get("id", "")
            print_info(f"Found machine: {ctx.machine_id}")
        
        stats.passed += 1
        return True
    else:
        print_error("Machine list test failed")
        record_test_failure("Machine - List All")
        return False

def test_machine_get():
    """Test: GET /machine/{machineId}"""
    print_test("Machine - Get by ID")
    stats.total += 1
    
    if not ctx.machine_id:
        print_skip("No machine ID available")
        stats.skipped += 1
        return False
    
    success, data = api_request("GET", f"/machine/{ctx.machine_id}")
    if success and data.get("id") == ctx.machine_id:
        print_success(f"Machine retrieved (vendor: {data.get('vendor')}, product: {data.get('productName')})")
        stats.passed += 1
        return True
    else:
        print_error("Machine get test failed")
        record_test_failure("Machine - Get by ID")
        return False

def test_machine_capabilities():
    """Test: GET /machine-capability"""
    print_test("Machine Capabilities - List All")
    stats.total += 1
    
    if not ctx.site_id:
        print_skip("No site ID available")
        stats.skipped += 1
        return False
    
    params = {"siteId": ctx.site_id}
    success, data = api_request("GET", "/machine-capability", params=params)
    if success:
        count = len(data) if isinstance(data, list) else 0
        print_success(f"Machine Capabilities listed ({count} capabilities)")
        stats.passed += 1
        return True
    else:
        print_error("Machine Capabilities test failed")
        record_test_failure("Machine Capabilities - List All")
        return False

def test_expected_machine_create():
    """Test: POST /expected-machine"""
    print_test("Expected Machine - Create")
    stats.total += 1
    
    if not ctx.site_id:
        print_skip("No site ID available")
        stats.skipped += 1
        return False
    
    payload = {
        "siteId": ctx.site_id,
        "bmcMacAddress": "AA:BB:CC:DD:EE:FF",
        "bmcUsername": "testadmin",
        "bmcPassword": "testpass123",
        "chassisSerialNumber": f"CHASSIS-{int(time.time())}",
        "labels": {
            "test": "comprehensive"
        }
    }
    
    success, data = api_request("POST", "/expected-machine", data=payload)
    if success and "id" in data:
        ctx.expected_machine_id = data["id"]
        print_success(f"Expected Machine created (ID: {ctx.expected_machine_id})")
        stats.passed += 1
        return True
    else:
        print_error("Expected Machine create test failed")
        record_test_failure("Expected Machine - Create")
        return False

def test_expected_machine_get():
    """Test: GET /expected-machine/{expectedMachineId}"""
    print_test("Expected Machine - Get by ID")
    stats.total += 1
    
    if not ctx.expected_machine_id:
        print_skip("No expected machine ID available")
        stats.skipped += 1
        return False
    
    success, data = api_request("GET", f"/expected-machine/{ctx.expected_machine_id}")
    if success and data.get("id") == ctx.expected_machine_id:
        print_success(f"Expected Machine retrieved")
        stats.passed += 1
        return True
    else:
        print_error("Expected Machine get test failed")
        record_test_failure("Expected Machine - Get by ID")
        return False

def test_expected_machine_list():
    """Test: GET /expected-machine"""
    print_test("Expected Machine - List All")
    stats.total += 1
    
    if not ctx.site_id:
        print_skip("No site ID available")
        stats.skipped += 1
        return False
    
    params = {"siteId": ctx.site_id}
    success, data = api_request("GET", "/expected-machine", params=params)
    if success:
        count = len(data) if isinstance(data, list) else 0
        print_success(f"Expected Machines listed ({count} machines)")
        stats.passed += 1
        return True
    else:
        print_error("Expected Machine list test failed")
        record_test_failure("Expected Machine - List All")
        return False

def test_expected_machine_update():
    """Test: PATCH /expected-machine/{expectedMachineId}"""
    print_test("Expected Machine - Update")
    stats.total += 1
    
    if not ctx.expected_machine_id:
        print_skip("No expected machine ID available")
        stats.skipped += 1
        return False
    
    payload = {
        "labels": {
            "test": "comprehensive",
            "updated": "true"
        }
    }
    success, data = api_request("PATCH", f"/expected-machine/{ctx.expected_machine_id}", data=payload)
    if success:
        print_success("Expected Machine updated successfully")
        stats.passed += 1
        return True
    else:
        print_error("Expected Machine update test failed")
        record_test_failure("Expected Machine - Update")
        return False

def test_tenant_current():
    """Test: GET /tenant/current"""
    print_test("Tenant - Get Current")
    stats.total += 1
    
    success, data = api_request("GET", "/tenant/current")
    if success and "id" in data:
        ctx.tenant_id = data["id"]
        print_success(f"Tenant retrieved (ID: {ctx.tenant_id})")
        stats.passed += 1
        return True
    else:
        print_info("Tenant not found (expected for provider-only org)")
        stats.skipped += 1
        return False

def test_tenant_stats():
    """Test: GET /tenant/current/stats"""
    print_test("Tenant - Get Stats")
    stats.total += 1
    
    if not ctx.tenant_id:
        print_skip("No tenant ID available")
        stats.skipped += 1
        return False
    
    success, data = api_request("GET", "/tenant/current/stats")
    if success:
        print_success(f"Tenant stats retrieved")
        print_info(f"Stats: {json.dumps(data, indent=2)}")
        stats.passed += 1
        return True
    else:
        print_error("Tenant stats test failed")
        record_test_failure("Tenant - Get Stats")
        return False

def test_tenant_account_list():
    """Test: GET /tenant/account"""
    print_test("Tenant Account - List All")
    stats.total += 1
    
    params = {"infrastructureProviderId": ctx.infra_provider_id}
    success, data = api_request("GET", "/tenant/account", params=params)
    if success:
        count = len(data) if isinstance(data, list) else 0
        print_success(f"Tenant Accounts listed ({count} accounts)")
        stats.passed += 1
        return True
    else:
        print_error("Tenant Account list test failed")
        record_test_failure("Tenant Account - List All")
        return False

def test_allocation_create():
    """Test: POST /allocation"""
    print_test("Allocation - Create")
    stats.total += 1
    
    if not ctx.site_id or not ctx.instance_type_id:
        print_skip("Prerequisites not met (site and instance type required)")
        stats.skipped += 1
        return False
    
    # Need a tenant for allocation - skip if provider-only
    if not ctx.tenant_id:
        print_skip("Tenant required for allocation creation")
        stats.skipped += 1
        return False
    
    payload = {
        "name": f"test-allocation-{int(time.time())}",
        "description": "Test allocation for API testing",
        "tenantId": ctx.tenant_id,
        "siteId": ctx.site_id,
        "constraints": [
            {
                "resourceType": "InstanceType",
                "resourceTypeId": ctx.instance_type_id,
                "constraintType": "Reserved",
                "constraintValue": 2
            }
        ]
    }
    
    success, data = api_request("POST", "/allocation", data=payload)
    if success and "id" in data:
        ctx.allocation_id = data["id"]
        print_success(f"Allocation created (ID: {ctx.allocation_id})")
        stats.passed += 1
        return True
    else:
        print_error("Allocation create test failed")
        record_test_failure("Allocation - Create")
        return False

def test_allocation_get():
    """Test: GET /allocation/{allocationId}"""
    print_test("Allocation - Get by ID")
    stats.total += 1
    
    if not ctx.allocation_id:
        print_skip("No allocation ID available")
        stats.skipped += 1
        return False
    
    params = {"infrastructureProviderId": ctx.infra_provider_id}
    success, data = api_request("GET", f"/allocation/{ctx.allocation_id}", params=params)
    if success and data.get("id") == ctx.allocation_id:
        print_success(f"Allocation retrieved (name: {data.get('name')})")
        stats.passed += 1
        return True
    else:
        print_error("Allocation get test failed")
        record_test_failure("Allocation - Get by ID")
        return False

def test_allocation_list():
    """Test: GET /allocation"""
    print_test("Allocation - List All")
    stats.total += 1
    
    params = {"infrastructureProviderId": ctx.infra_provider_id}
    success, data = api_request("GET", "/allocation", params=params)
    if success:
        count = len(data) if isinstance(data, list) else 0
        print_success(f"Allocations listed ({count} allocations)")
        stats.passed += 1
        return True
    else:
        print_error("Allocation list test failed")
        record_test_failure("Allocation - List All")
        return False

def test_allocation_update():
    """Test: PATCH /allocation/{allocationId}"""
    print_test("Allocation - Update")
    stats.total += 1
    
    if not ctx.allocation_id:
        print_skip("No allocation ID available")
        stats.skipped += 1
        return False
    
    payload = {"description": "Updated allocation description"}
    success, data = api_request("PATCH", f"/allocation/{ctx.allocation_id}", data=payload)
    if success:
        print_success("Allocation updated successfully")
        stats.passed += 1
        return True
    else:
        print_error("Allocation update test failed")
        record_test_failure("Allocation - Update")
        return False

def test_allocation_constraint_list():
    """Test: GET /allocation/{allocationId}/constraint"""
    print_test("Allocation Constraint - List All")
    stats.total += 1
    
    if not ctx.allocation_id:
        print_skip("No allocation ID available")
        stats.skipped += 1
        return False
    
    success, data = api_request("GET", f"/allocation/{ctx.allocation_id}/constraint")
    if success:
        count = len(data) if isinstance(data, list) else 0
        print_success(f"Allocation Constraints listed ({count} constraints)")
        
        # Store first constraint ID
        if isinstance(data, list) and len(data) > 0:
            ctx.allocation_constraint_id = data[0].get("id", "")
        
        stats.passed += 1
        return True
    else:
        print_error("Allocation Constraint list test failed")
        record_test_failure("Allocation Constraint - List All")
        return False

def test_vpc_create():
    """Test: POST /vpc"""
    print_test("VPC - Create")
    stats.total += 1
    
    if not ctx.site_id or not ctx.tenant_id:
        print_skip("Prerequisites not met (site and tenant required)")
        stats.skipped += 1
        return False
    
    payload = {
        "name": f"test-vpc-{int(time.time())}",
        "description": "Test VPC for API testing",
        "siteId": ctx.site_id,
        "networkVirtualizationType": "ETHERNET_VIRTUALIZER",
        "labels": {
            "env": "test",
            "purpose": "api-testing"
        }
    }
    
    success, data = api_request("POST", "/vpc", data=payload)
    if success and "id" in data:
        ctx.vpc_id = data["id"]
        print_success(f"VPC created (ID: {ctx.vpc_id})")
        stats.passed += 1
        return True
    else:
        print_error("VPC create test failed")
        record_test_failure("VPC - Create")
        return False

def test_vpc_get():
    """Test: GET /vpc/{vpcId}"""
    print_test("VPC - Get by ID")
    stats.total += 1
    
    if not ctx.vpc_id:
        print_skip("No VPC ID available")
        stats.skipped += 1
        return False
    
    success, data = api_request("GET", f"/vpc/{ctx.vpc_id}")
    if success and data.get("id") == ctx.vpc_id:
        print_success(f"VPC retrieved (name: {data.get('name')})")
        stats.passed += 1
        return True
    else:
        print_error("VPC get test failed")
        record_test_failure("VPC - Get by ID")
        return False

def test_vpc_list():
    """Test: GET /vpc"""
    print_test("VPC - List All")
    stats.total += 1
    
    success, data = api_request("GET", "/vpc")
    if success:
        count = len(data) if isinstance(data, list) else 0
        print_success(f"VPCs listed ({count} VPCs)")
        stats.passed += 1
        return True
    else:
        print_error("VPC list test failed")
        record_test_failure("VPC - List All")
        return False

def test_vpc_update():
    """Test: PATCH /vpc/{vpcId}"""
    print_test("VPC - Update")
    stats.total += 1
    
    if not ctx.vpc_id:
        print_skip("No VPC ID available")
        stats.skipped += 1
        return False
    
    payload = {
        "description": "Updated VPC description",
        "labels": {
            "env": "test",
            "purpose": "api-testing",
            "updated": "true"
        }
    }
    success, data = api_request("PATCH", f"/vpc/{ctx.vpc_id}", data=payload)
    if success:
        print_success("VPC updated successfully")
        stats.passed += 1
        return True
    else:
        print_error("VPC update test failed")
        record_test_failure("VPC - Update")
        return False

def test_vpc_prefix_create():
    """Test: POST /vpc-prefix"""
    print_test("VPC Prefix - Create")
    stats.total += 1
    
    if not ctx.vpc_id or not ctx.ip_block_id:
        print_skip("Prerequisites not met (VPC and IP block required)")
        stats.skipped += 1
        return False
    
    payload = {
        "name": f"test-vpc-prefix-{int(time.time())}",
        "vpcId": ctx.vpc_id,
        "ipBlockId": ctx.ip_block_id,
        "prefixLength": 24
    }
    
    success, data = api_request("POST", "/vpc-prefix", data=payload)
    if success and "id" in data:
        ctx.vpc_prefix_id = data["id"]
        print_success(f"VPC Prefix created (ID: {ctx.vpc_prefix_id})")
        stats.passed += 1
        return True
    else:
        print_error("VPC Prefix create test failed")
        record_test_failure("VPC Prefix - Create")
        return False

def test_vpc_prefix_get():
    """Test: GET /vpc-prefix/{vpcPrefixId}"""
    print_test("VPC Prefix - Get by ID")
    stats.total += 1
    
    if not ctx.vpc_prefix_id:
        print_skip("No VPC prefix ID available")
        stats.skipped += 1
        return False
    
    success, data = api_request("GET", f"/vpc-prefix/{ctx.vpc_prefix_id}")
    if success and data.get("id") == ctx.vpc_prefix_id:
        print_success(f"VPC Prefix retrieved (name: {data.get('name')})")
        stats.passed += 1
        return True
    else:
        print_error("VPC Prefix get test failed")
        record_test_failure("VPC Prefix - Get by ID")
        return False

def test_vpc_prefix_list():
    """Test: GET /vpc-prefix"""
    print_test("VPC Prefix - List All")
    stats.total += 1
    
    if not ctx.site_id:
        print_skip("No site ID available")
        stats.skipped += 1
        return False
    
    params = {"siteId": ctx.site_id}
    success, data = api_request("GET", "/vpc-prefix", params=params)
    if success:
        count = len(data) if isinstance(data, list) else 0
        print_success(f"VPC Prefixes listed ({count} prefixes)")
        stats.passed += 1
        return True
    else:
        print_error("VPC Prefix list test failed")
        record_test_failure("VPC Prefix - List All")
        return False

def test_subnet_create():
    """Test: POST /subnet"""
    print_test("Subnet - Create")
    stats.total += 1
    
    if not ctx.vpc_id or not ctx.ip_block_id:
        print_skip("Prerequisites not met (VPC and IP block required)")
        stats.skipped += 1
        return False
    
    payload = {
        "name": f"test-subnet-{int(time.time())}",
        "description": "Test subnet for API testing",
        "vpcId": ctx.vpc_id,
        "ipv4BlockId": ctx.ip_block_id,
        "prefixLength": 24
    }
    
    success, data = api_request("POST", "/subnet", data=payload)
    if success and "id" in data:
        ctx.subnet_id = data["id"]
        print_success(f"Subnet created (ID: {ctx.subnet_id})")
        stats.passed += 1
        return True
    else:
        print_error("Subnet create test failed")
        record_test_failure("Subnet - Create")
        return False

def test_subnet_get():
    """Test: GET /subnet/{subnetId}"""
    print_test("Subnet - Get by ID")
    stats.total += 1
    
    if not ctx.subnet_id:
        print_skip("No subnet ID available")
        stats.skipped += 1
        return False
    
    success, data = api_request("GET", f"/subnet/{ctx.subnet_id}")
    if success and data.get("id") == ctx.subnet_id:
        print_success(f"Subnet retrieved (name: {data.get('name')})")
        stats.passed += 1
        return True
    else:
        print_error("Subnet get test failed")
        record_test_failure("Subnet - Get by ID")
        return False

def test_subnet_list():
    """Test: GET /subnet"""
    print_test("Subnet - List All")
    stats.total += 1
    
    if not ctx.site_id:
        print_skip("No site ID available")
        stats.skipped += 1
        return False
    
    params = {"siteId": ctx.site_id}
    success, data = api_request("GET", "/subnet", params=params)
    if success:
        count = len(data) if isinstance(data, list) else 0
        print_success(f"Subnets listed ({count} subnets)")
        stats.passed += 1
        return True
    else:
        print_error("Subnet list test failed")
        record_test_failure("Subnet - List All")
        return False

def test_subnet_update():
    """Test: PATCH /subnet/{subnetId}"""
    print_test("Subnet - Update")
    stats.total += 1
    
    if not ctx.subnet_id:
        print_skip("No subnet ID available")
        stats.skipped += 1
        return False
    
    payload = {
        "name": f"updated-subnet-{int(time.time())}",
        "description": "Updated subnet description"
    }
    success, data = api_request("PATCH", f"/subnet/{ctx.subnet_id}", data=payload)
    if success:
        print_success("Subnet updated successfully")
        stats.passed += 1
        return True
    else:
        print_error("Subnet update test failed")
        record_test_failure("Subnet - Update")
        return False

def test_operating_system_create():
    """Test: POST /operating-system"""
    print_test("Operating System - Create")
    stats.total += 1
    
    if not ctx.tenant_id:
        print_skip("Tenant required for operating system creation")
        stats.skipped += 1
        return False
    
    payload = {
        "name": f"test-os-{int(time.time())}",
        "description": "Test operating system for API testing",
        "tenantId": ctx.tenant_id,
        "type": "iPXE",
        "ipxeScript": "#!ipxe\nchain http://example.com/boot.ipxe",
        "isCloudInit": True,
        "phoneHomeEnabled": False,
        "allowOverride": False
    }
    
    success, data = api_request("POST", "/operating-system", data=payload)
    if success and "id" in data:
        ctx.operating_system_id = data["id"]
        print_success(f"Operating System created (ID: {ctx.operating_system_id})")
        stats.passed += 1
        return True
    else:
        print_error("Operating System create test failed")
        record_test_failure("Operating System - Create")
        return False

def test_operating_system_get():
    """Test: GET /operating-system/{operatingSystemId}"""
    print_test("Operating System - Get by ID")
    stats.total += 1
    
    if not ctx.operating_system_id:
        print_skip("No operating system ID available")
        stats.skipped += 1
        return False
    
    success, data = api_request("GET", f"/operating-system/{ctx.operating_system_id}")
    if success and data.get("id") == ctx.operating_system_id:
        print_success(f"Operating System retrieved (name: {data.get('name')})")
        stats.passed += 1
        return True
    else:
        print_error("Operating System get test failed")
        record_test_failure("Operating System - Get by ID")
        return False

def test_operating_system_list():
    """Test: GET /operating-system"""
    print_test("Operating System - List All")
    stats.total += 1
    
    success, data = api_request("GET", "/operating-system")
    if success:
        count = len(data) if isinstance(data, list) else 0
        print_success(f"Operating Systems listed ({count} OSs)")
        stats.passed += 1
        return True
    else:
        print_error("Operating System list test failed")
        record_test_failure("Operating System - List All")
        return False

def test_operating_system_update():
    """Test: PATCH /operating-system/{operatingSystemId}"""
    print_test("Operating System - Update")
    stats.total += 1
    
    if not ctx.operating_system_id:
        print_skip("No operating system ID available")
        stats.skipped += 1
        return False
    
    payload = {
        "description": "Updated operating system description",
        "allowOverride": True
    }
    success, data = api_request("PATCH", f"/operating-system/{ctx.operating_system_id}", data=payload)
    if success:
        print_success("Operating System updated successfully")
        stats.passed += 1
        return True
    else:
        print_error("Operating System update test failed")
        record_test_failure("Operating System - Update")
        return False

def test_ssh_key_create():
    """Test: POST /sshkey"""
    print_test("SSH Key - Create")
    stats.total += 1
    
    if not ctx.tenant_id:
        print_skip("Tenant required for SSH key creation")
        stats.skipped += 1
        return False
    
    payload = {
        "name": f"test-ssh-key-{int(time.time())}",
        "publicKey": "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIOMqqnkVzrm0SdG6UOoqKLsabgH5C9okWi0dh2l9GKJl test@test.com",
        "tenantId": ctx.tenant_id
    }
    
    success, data = api_request("POST", "/sshkey", data=payload)
    if success and "id" in data:
        ctx.ssh_key_id = data["id"]
        print_success(f"SSH Key created (ID: {ctx.ssh_key_id})")
        stats.passed += 1
        return True
    else:
        print_error("SSH Key create test failed")
        record_test_failure("SSH Key - Create")
        return False

def test_ssh_key_get():
    """Test: GET /sshkey/{sshKeyId}"""
    print_test("SSH Key - Get by ID")
    stats.total += 1
    
    if not ctx.ssh_key_id:
        print_skip("No SSH key ID available")
        stats.skipped += 1
        return False
    
    success, data = api_request("GET", f"/sshkey/{ctx.ssh_key_id}")
    if success and data.get("id") == ctx.ssh_key_id:
        print_success(f"SSH Key retrieved (fingerprint: {data.get('fingerprint', 'N/A')[:20]}...)")
        stats.passed += 1
        return True
    else:
        print_error("SSH Key get test failed")
        record_test_failure("SSH Key - Get by ID")
        return False

def test_ssh_key_list():
    """Test: GET /sshkey"""
    print_test("SSH Key - List All")
    stats.total += 1
    
    success, data = api_request("GET", "/sshkey")
    if success:
        count = len(data) if isinstance(data, list) else 0
        print_success(f"SSH Keys listed ({count} keys)")
        stats.passed += 1
        return True
    else:
        print_error("SSH Key list test failed")
        record_test_failure("SSH Key - List All")
        return False

def test_ssh_key_update():
    """Test: PATCH /sshkey/{sshKeyId}"""
    print_test("SSH Key - Update")
    stats.total += 1
    
    if not ctx.ssh_key_id:
        print_skip("No SSH key ID available")
        stats.skipped += 1
        return False
    
    payload = {"name": f"updated-ssh-key-{int(time.time())}"}
    success, data = api_request("PATCH", f"/sshkey/{ctx.ssh_key_id}", data=payload)
    if success:
        print_success("SSH Key updated successfully")
        stats.passed += 1
        return True
    else:
        print_error("SSH Key update test failed")
        record_test_failure("SSH Key - Update")
        return False

def test_ssh_key_group_create():
    """Test: POST /sshkeygroup"""
    print_test("SSH Key Group - Create")
    stats.total += 1
    
    if not ctx.ssh_key_id or not ctx.site_id:
        print_skip("Prerequisites not met (SSH key and site required)")
        stats.skipped += 1
        return False
    
    payload = {
        "name": f"test-ssh-group-{int(time.time())}",
        "description": "Test SSH key group for API testing",
        "siteIds": [ctx.site_id],
        "sshKeyIds": [ctx.ssh_key_id]
    }
    
    success, data = api_request("POST", "/sshkeygroup", data=payload)
    if success and "id" in data:
        ctx.ssh_key_group_id = data["id"]
        print_success(f"SSH Key Group created (ID: {ctx.ssh_key_group_id})")
        stats.passed += 1
        return True
    else:
        print_error("SSH Key Group create test failed")
        record_test_failure("SSH Key Group - Create")
        return False

def test_ssh_key_group_get():
    """Test: GET /sshkeygroup/{sshKeyGroupId}"""
    print_test("SSH Key Group - Get by ID")
    stats.total += 1
    
    if not ctx.ssh_key_group_id:
        print_skip("No SSH key group ID available")
        stats.skipped += 1
        return False
    
    success, data = api_request("GET", f"/sshkeygroup/{ctx.ssh_key_group_id}")
    if success and data.get("id") == ctx.ssh_key_group_id:
        print_success(f"SSH Key Group retrieved (name: {data.get('name')})")
        stats.passed += 1
        return True
    else:
        print_error("SSH Key Group get test failed")
        record_test_failure("SSH Key Group - Get by ID")
        return False

def test_ssh_key_group_list():
    """Test: GET /sshkeygroup"""
    print_test("SSH Key Group - List All")
    stats.total += 1
    
    success, data = api_request("GET", "/sshkeygroup")
    if success:
        count = len(data) if isinstance(data, list) else 0
        print_success(f"SSH Key Groups listed ({count} groups)")
        stats.passed += 1
        return True
    else:
        print_error("SSH Key Group list test failed")
        record_test_failure("SSH Key Group - List All")
        return False

def test_infiniband_partition_create():
    """Test: POST /infiniband-partition"""
    print_test("InfiniBand Partition - Create")
    stats.total += 1
    
    if not ctx.site_id or not ctx.tenant_id:
        print_skip("Prerequisites not met (site and tenant required)")
        stats.skipped += 1
        return False
    
    payload = {
        "name": f"test-ib-partition-{int(time.time())}",
        "description": "Test InfiniBand partition for API testing",
        "siteId": ctx.site_id
    }
    
    success, data = api_request("POST", "/infiniband-partition", data=payload)
    if success and "id" in data:
        ctx.infiniband_partition_id = data["id"]
        print_success(f"InfiniBand Partition created (ID: {ctx.infiniband_partition_id})")
        stats.passed += 1
        return True
    else:
        print_error("InfiniBand Partition create test failed")
        record_test_failure("InfiniBand Partition - Create")
        return False

def test_infiniband_partition_get():
    """Test: GET /infiniband-partition/{infiniBandPartitionId}"""
    print_test("InfiniBand Partition - Get by ID")
    stats.total += 1
    
    if not ctx.infiniband_partition_id:
        print_skip("No InfiniBand partition ID available")
        stats.skipped += 1
        return False
    
    success, data = api_request("GET", f"/infiniband-partition/{ctx.infiniband_partition_id}")
    if success and data.get("id") == ctx.infiniband_partition_id:
        print_success(f"InfiniBand Partition retrieved (name: {data.get('name')})")
        stats.passed += 1
        return True
    else:
        print_error("InfiniBand Partition get test failed")
        record_test_failure("InfiniBand Partition - Get by ID")
        return False

def test_infiniband_partition_list():
    """Test: GET /infiniband-partition"""
    print_test("InfiniBand Partition - List All")
    stats.total += 1
    
    if not ctx.site_id:
        print_skip("No site ID available")
        stats.skipped += 1
        return False
    
    params = {"siteId": ctx.site_id}
    success, data = api_request("GET", "/infiniband-partition", params=params)
    if success:
        count = len(data) if isinstance(data, list) else 0
        print_success(f"InfiniBand Partitions listed ({count} partitions)")
        stats.passed += 1
        return True
    else:
        print_error("InfiniBand Partition list test failed")
        record_test_failure("InfiniBand Partition - List All")
        return False

def test_infiniband_partition_update():
    """Test: PATCH /infiniband-partition/{infiniBandPartitionId}"""
    print_test("InfiniBand Partition - Update")
    stats.total += 1
    
    if not ctx.infiniband_partition_id:
        print_skip("No InfiniBand partition ID available")
        stats.skipped += 1
        return False
    
    payload = {
        "name": f"updated-ib-partition-{int(time.time())}",
        "description": "Updated InfiniBand partition description"
    }
    success, data = api_request("PATCH", f"/infiniband-partition/{ctx.infiniband_partition_id}", data=payload)
    if success:
        print_success("InfiniBand Partition updated successfully")
        stats.passed += 1
        return True
    else:
        print_error("InfiniBand Partition update test failed")
        record_test_failure("InfiniBand Partition - Update")
        return False

def test_network_security_group_create():
    """Test: POST /network-security-group"""
    print_test("Network Security Group - Create")
    stats.total += 1
    
    if not ctx.site_id or not ctx.tenant_id:
        print_skip("Prerequisites not met (site and tenant required)")
        stats.skipped += 1
        return False
    
    payload = {
        "name": f"test-nsg-{int(time.time())}",
        "description": "Test network security group for API testing",
        "siteId": ctx.site_id,
        "rules": [
            {
                "name": "allow-http",
                "direction": "INGRESS",
                "sourcePortRange": "80",
                "destinationPortRange": "80",
                "protocol": "TCP",
                "action": "PERMIT",
                "priority": 100,
                "sourcePrefix": "0.0.0.0/0",
                "destinationPrefix": "0.0.0.0/0"
            }
        ],
        "labels": {
            "test": "comprehensive"
        }
    }
    
    success, data = api_request("POST", "/network-security-group", data=payload)
    if success and "id" in data:
        ctx.network_security_group_id = data["id"]
        print_success(f"Network Security Group created (ID: {ctx.network_security_group_id})")
        stats.passed += 1
        return True
    else:
        print_error("Network Security Group create test failed")
        record_test_failure("Network Security Group - Create")
        return False

def test_network_security_group_get():
    """Test: GET /network-security-group/{networkSecurityGroupId}"""
    print_test("Network Security Group - Get by ID")
    stats.total += 1
    
    if not ctx.network_security_group_id:
        print_skip("No network security group ID available")
        stats.skipped += 1
        return False
    
    success, data = api_request("GET", f"/network-security-group/{ctx.network_security_group_id}")
    if success and data.get("id") == ctx.network_security_group_id:
        print_success(f"Network Security Group retrieved (name: {data.get('name')})")
        stats.passed += 1
        return True
    else:
        print_error("Network Security Group get test failed")
        record_test_failure("Network Security Group - Get by ID")
        return False

def test_network_security_group_list():
    """Test: GET /network-security-group"""
    print_test("Network Security Group - List All")
    stats.total += 1
    
    success, data = api_request("GET", "/network-security-group")
    if success:
        count = len(data) if isinstance(data, list) else 0
        print_success(f"Network Security Groups listed ({count} groups)")
        stats.passed += 1
        return True
    else:
        print_error("Network Security Group list test failed")
        record_test_failure("Network Security Group - List All")
        return False

def test_network_security_group_update():
    """Test: PATCH /network-security-group/{networkSecurityGroupId}"""
    print_test("Network Security Group - Update")
    stats.total += 1
    
    if not ctx.network_security_group_id:
        print_skip("No network security group ID available")
        stats.skipped += 1
        return False
    
    payload = {
        "description": "Updated network security group description",
        "rules": [
            {
                "name": "allow-http",
                "direction": "INGRESS",
                "sourcePortRange": "80",
                "destinationPortRange": "80",
                "protocol": "TCP",
                "action": "PERMIT",
                "priority": 100,
                "sourcePrefix": "0.0.0.0/0",
                "destinationPrefix": "0.0.0.0/0"
            },
            {
                "name": "allow-https",
                "direction": "INGRESS",
                "sourcePortRange": "443",
                "destinationPortRange": "443",
                "protocol": "TCP",
                "action": "PERMIT",
                "priority": 101,
                "sourcePrefix": "0.0.0.0/0",
                "destinationPrefix": "0.0.0.0/0"
            }
        ]
    }
    success, data = api_request("PATCH", f"/network-security-group/{ctx.network_security_group_id}", data=payload)
    if success:
        print_success("Network Security Group updated successfully")
        stats.passed += 1
        return True
    else:
        print_error("Network Security Group update test failed")
        record_test_failure("Network Security Group - Update")
        return False

def test_audit_list():
    """Test: GET /audit"""
    print_test("Audit - List Entries")
    stats.total += 1
    
    params = {"pageSize": 10}
    success, data = api_request("GET", "/audit", params=params)
    if success:
        count = len(data) if isinstance(data, list) else 0
        print_success(f"Audit entries retrieved ({count} entries)")
        stats.passed += 1
        return True
    else:
        print_error("Audit list test failed")
        record_test_failure("Audit - List Entries")
        return False

def test_network_security_group_delete():
    """Test: DELETE /network-security-group/{networkSecurityGroupId}"""
    print_test("Network Security Group - Delete")
    stats.total += 1
    
    if not ctx.network_security_group_id:
        print_skip("No network security group ID available")
        stats.skipped += 1
        return False
    
    success, _ = api_request("DELETE", f"/network-security-group/{ctx.network_security_group_id}")
    if success:
        print_success("Network Security Group deleted successfully")
        stats.passed += 1
        return True
    else:
        print_error("Network Security Group delete test failed")
        record_test_failure("Network Security Group - Delete")
        return False

def test_infiniband_partition_delete():
    """Test: DELETE /infiniband-partition/{infiniBandPartitionId}"""
    print_test("InfiniBand Partition - Delete")
    stats.total += 1
    
    if not ctx.infiniband_partition_id:
        print_skip("No InfiniBand partition ID available")
        stats.skipped += 1
        return False
    
    success, _ = api_request("DELETE", f"/infiniband-partition/{ctx.infiniband_partition_id}")
    if success:
        print_success("InfiniBand Partition deleted successfully")
        stats.passed += 1
        return True
    else:
        print_error("InfiniBand Partition delete test failed")
        record_test_failure("InfiniBand Partition - Delete")
        return False

def test_ssh_key_group_delete():
    """Test: DELETE /sshkeygroup/{sshKeyGroupId}"""
    print_test("SSH Key Group - Delete")
    stats.total += 1
    
    if not ctx.ssh_key_group_id:
        print_skip("No SSH key group ID available")
        stats.skipped += 1
        return False
    
    success, _ = api_request("DELETE", f"/sshkeygroup/{ctx.ssh_key_group_id}")
    if success:
        print_success("SSH Key Group deleted successfully")
        stats.passed += 1
        return True
    else:
        print_error("SSH Key Group delete test failed")
        record_test_failure("SSH Key Group - Delete")
        return False

def test_ssh_key_delete():
    """Test: DELETE /sshkey/{sshKeyId}"""
    print_test("SSH Key - Delete")
    stats.total += 1
    
    if not ctx.ssh_key_id:
        print_skip("No SSH key ID available")
        stats.skipped += 1
        return False
    
    success, _ = api_request("DELETE", f"/sshkey/{ctx.ssh_key_id}")
    if success:
        print_success("SSH Key deleted successfully")
        stats.passed += 1
        return True
    else:
        print_error("SSH Key delete test failed")
        record_test_failure("SSH Key - Delete")
        return False

def test_operating_system_delete():
    """Test: DELETE /operating-system/{operatingSystemId}"""
    print_test("Operating System - Delete")
    stats.total += 1
    
    if not ctx.operating_system_id:
        print_skip("No operating system ID available")
        stats.skipped += 1
        return False
    
    success, _ = api_request("DELETE", f"/operating-system/{ctx.operating_system_id}")
    if success:
        print_success("Operating System deleted successfully")
        stats.passed += 1
        return True
    else:
        print_error("Operating System delete test failed")
        record_test_failure("Operating System - Delete")
        return False

def test_subnet_delete():
    """Test: DELETE /subnet/{subnetId}"""
    print_test("Subnet - Delete")
    stats.total += 1
    
    if not ctx.subnet_id:
        print_skip("No subnet ID available")
        stats.skipped += 1
        return False
    
    success, _ = api_request("DELETE", f"/subnet/{ctx.subnet_id}")
    if success:
        print_success("Subnet deleted successfully")
        stats.passed += 1
        return True
    else:
        print_error("Subnet delete test failed")
        record_test_failure("Subnet - Delete")
        return False

def test_vpc_prefix_delete():
    """Test: DELETE /vpc-prefix/{vpcPrefixId}"""
    print_test("VPC Prefix - Delete")
    stats.total += 1
    
    if not ctx.vpc_prefix_id:
        print_skip("No VPC prefix ID available")
        stats.skipped += 1
        return False
    
    success, _ = api_request("DELETE", f"/vpc-prefix/{ctx.vpc_prefix_id}")
    if success:
        print_success("VPC Prefix deleted successfully")
        stats.passed += 1
        return True
    else:
        print_error("VPC Prefix delete test failed")
        record_test_failure("VPC Prefix - Delete")
        return False

def test_vpc_delete():
    """Test: DELETE /vpc/{vpcId}"""
    print_test("VPC - Delete")
    stats.total += 1
    
    if not ctx.vpc_id:
        print_skip("No VPC ID available")
        stats.skipped += 1
        return False
    
    success, _ = api_request("DELETE", f"/vpc/{ctx.vpc_id}")
    if success:
        print_success("VPC deleted successfully")
        stats.passed += 1
        return True
    else:
        print_error("VPC delete test failed")
        record_test_failure("VPC - Delete")
        return False

def test_allocation_delete():
    """Test: DELETE /allocation/{allocationId}"""
    print_test("Allocation - Delete")
    stats.total += 1
    
    if not ctx.allocation_id:
        print_skip("No allocation ID available")
        stats.skipped += 1
        return False
    
    success, _ = api_request("DELETE", f"/allocation/{ctx.allocation_id}")
    if success:
        print_success("Allocation deleted successfully")
        stats.passed += 1
        return True
    else:
        print_error("Allocation delete test failed")
        record_test_failure("Allocation - Delete")
        return False

def test_expected_machine_delete():
    """Test: DELETE /expected-machine/{expectedMachineId}"""
    print_test("Expected Machine - Delete")
    stats.total += 1
    
    if not ctx.expected_machine_id:
        print_skip("No expected machine ID available")
        stats.skipped += 1
        return False
    
    success, _ = api_request("DELETE", f"/expected-machine/{ctx.expected_machine_id}")
    if success:
        print_success("Expected Machine deleted successfully")
        stats.passed += 1
        return True
    else:
        print_error("Expected Machine delete test failed")
        record_test_failure("Expected Machine - Delete")
        return False

def test_instance_type_delete():
    """Test: DELETE /instance/type/{instanceTypeId}"""
    print_test("Instance Type - Delete")
    stats.total += 1
    
    if not ctx.instance_type_id:
        print_skip("No instance type ID available")
        stats.skipped += 1
        return False
    
    success, _ = api_request("DELETE", f"/instance/type/{ctx.instance_type_id}")
    if success:
        print_success("Instance Type deleted successfully")
        stats.passed += 1
        return True
    else:
        print_error("Instance Type delete test failed")
        record_test_failure("Instance Type - Delete")
        return False

def test_ip_block_delete():
    """Test: DELETE /ipblock/{ipBlockId}"""
    print_test("IP Block - Delete")
    stats.total += 1
    
    if not ctx.ip_block_id:
        print_skip("No IP block ID available")
        stats.skipped += 1
        return False
    
    success, _ = api_request("DELETE", f"/ipblock/{ctx.ip_block_id}")
    if success:
        print_success("IP Block deleted successfully")
        stats.passed += 1
        return True
    else:
        print_error("IP Block delete test failed")
        record_test_failure("IP Block - Delete")
        return False

def test_site_delete():
    """Test: DELETE /site/{siteId}"""
    print_test("Site - Delete")
    stats.total += 1
    
    if not ctx.site_id:
        print_skip("No site ID available")
        stats.skipped += 1
        return False
    
    success, _ = api_request("DELETE", f"/site/{ctx.site_id}")
    if success:
        print_success("Site deleted successfully")
        stats.passed += 1
        return True
    else:
        print_error("Site delete test failed")
        record_test_failure("Site - Delete")
        return False

def run_infrastructure_provider_tests():
    """Run all Infrastructure Provider tests"""
    print_header("Infrastructure Provider Tests")
    test_infrastructure_provider()
    test_infrastructure_provider_stats()

def run_tenant_tests():
    """Run all Tenant tests"""
    print_header("Tenant Tests")
    test_tenant_current()
    test_tenant_stats()
    test_tenant_account_list()

def run_site_tests():
    """Run all Site tests"""
    print_header("Site Tests")
    test_site_create()
    test_site_get()
    test_site_list()
    test_site_update()
    test_site_status_history()

def run_ip_block_tests():
    """Run all IP Block tests"""
    print_header("IP Block Tests")
    test_ip_block_create()
    test_ip_block_get()
    test_ip_block_list()
    test_ip_block_update()

def run_instance_type_tests():
    """Run all Instance Type tests"""
    print_header("Instance Type Tests")
    test_instance_type_create()
    test_instance_type_get()
    test_instance_type_list()
    test_instance_type_update()

def run_machine_tests():
    """Run all Machine tests"""
    print_header("Machine Tests")
    test_machine_list()
    test_machine_get()
    test_machine_capabilities()

def run_expected_machine_tests():
    """Run all Expected Machine tests"""
    print_header("Expected Machine Tests")
    test_expected_machine_create()
    test_expected_machine_get()
    test_expected_machine_list()
    test_expected_machine_update()

def run_allocation_tests():
    """Run all Allocation tests"""
    print_header("Allocation Tests")
    test_allocation_create()
    test_allocation_get()
    test_allocation_list()
    test_allocation_update()
    test_allocation_constraint_list()

def run_vpc_tests():
    """Run all VPC tests"""
    print_header("VPC Tests")
    test_vpc_create()
    test_vpc_get()
    test_vpc_list()
    test_vpc_update()

def run_vpc_prefix_tests():
    """Run all VPC Prefix tests"""
    print_header("VPC Prefix Tests")
    test_vpc_prefix_create()
    test_vpc_prefix_get()
    test_vpc_prefix_list()

def run_subnet_tests():
    """Run all Subnet tests"""
    print_header("Subnet Tests")
    test_subnet_create()
    test_subnet_get()
    test_subnet_list()
    test_subnet_update()

def run_operating_system_tests():
    """Run all Operating System tests"""
    print_header("Operating System Tests")
    test_operating_system_create()
    test_operating_system_get()
    test_operating_system_list()
    test_operating_system_update()

def run_ssh_key_tests():
    """Run all SSH Key tests"""
    print_header("SSH Key Tests")
    test_ssh_key_create()
    test_ssh_key_get()
    test_ssh_key_list()
    test_ssh_key_update()

def run_ssh_key_group_tests():
    """Run all SSH Key Group tests"""
    print_header("SSH Key Group Tests")
    test_ssh_key_group_create()
    test_ssh_key_group_get()
    test_ssh_key_group_list()

def run_infiniband_partition_tests():
    """Run all InfiniBand Partition tests"""
    print_header("InfiniBand Partition Tests")
    test_infiniband_partition_create()
    test_infiniband_partition_get()
    test_infiniband_partition_list()
    test_infiniband_partition_update()

def run_network_security_group_tests():
    """Run all Network Security Group tests"""
    print_header("Network Security Group Tests")
    test_network_security_group_create()
    test_network_security_group_get()
    test_network_security_group_list()
    test_network_security_group_update()

def run_cleanup_tests():
    """Run cleanup/deletion tests in proper order"""
    print_header("Cleanup Tests - Deleting Created Resources")
    
    # Delete in reverse dependency order
    test_network_security_group_delete()
    test_infiniband_partition_delete()
    test_ssh_key_group_delete()
    test_ssh_key_delete()
    test_operating_system_delete()
    test_subnet_delete()
    test_vpc_prefix_delete()
    test_vpc_delete()
    test_allocation_delete()
    test_expected_machine_delete()
    test_instance_type_delete()
    test_ip_block_delete()
    test_site_delete()

def run_misc_tests():
    """Run miscellaneous tests"""
    print_header("Miscellaneous Tests")
    test_metadata()
    test_user_current()
    test_audit_list()

def main():
    """Main test execution flow"""
    parser = argparse.ArgumentParser(
        description="Comprehensive API test suite for Forge Cloud API",
        formatter_class=argparse.RawDescriptionHelpFormatter,
        epilog="""
This test suite covers all major API endpoints from the OpenAPI specification.

Test Categories:
  - Infrastructure Provider (GET current, stats)
  - Tenant (GET current, stats, accounts)
  - Site (CRUD operations, status history)
  - IP Block (CRUD operations)
  - Instance Type (CRUD operations, machine associations)
  - Machine (List, get, capabilities)
  - Expected Machine (CRUD operations)
  - Allocation (CRUD operations, constraints)
  - VPC (CRUD operations)
  - VPC Prefix (CRUD operations)
  - Subnet (CRUD operations)
  - Operating System (CRUD operations)
  - SSH Key (CRUD operations)
  - SSH Key Group (CRUD operations)
  - InfiniBand Partition (CRUD operations)
  - Network Security Group (CRUD operations)
  - User (GET current)
  - Audit (List entries)
  - Metadata (GET server info)

Examples:
  # Run all tests without site registration (basic validation)
  python3 test-api.py

  # Run all tests WITH site registration (full testing - RECOMMENDED)
  python3 test-api.py --register-site

  # Skip port-forward setup (if already running)
  python3 test-api.py --skip-setup --register-site
        """
    )
    parser.add_argument(
        "--skip-setup",
        action="store_true",
        help="Skip port-forward setup (assumes already running)"
    )
    parser.add_argument(
        "--register-site",
        action="store_true",
        help="Register the site after creation to enable full API testing"
    )
    args = parser.parse_args()
    
    print_header("Comprehensive Forge Cloud API Test Suite")
    print_info(f"API Base URL: {API_BASE_URL}")
    print_info(f"Keycloak URL: {KEYCLOAK_URL}")
    print_info(f"Site Registration: {'ENABLED' if args.register_site else 'DISABLED'}")
    print_info(f"Testing {sum(1 for _ in dir() if _.startswith('test_'))} API endpoints")
    
    # Setup
    keycloak_pf = None
    api_pf = None
    
    try:
        if not args.skip_setup:
            # Restart site-manager first to get fresh RBAC tokens
            restart_site_manager_for_rbac()
            keycloak_pf, api_pf = setup_port_forwards()
        else:
            print_info("Skipping port-forward setup")
        
        # Authenticate
        ctx.token = get_auth_token()
        
        # Run test suites in logical order
        # 1. Basic system tests
        run_misc_tests()
        
        # 2. Organization tests
        run_infrastructure_provider_tests()
        run_tenant_tests()
        
        # 3. Site and infrastructure tests
        run_site_tests()
        
        # 3a. Register site if requested (enables full testing)
        if args.register_site and ctx.site_id:
            if perform_site_registration(ctx.site_id):
                print_info("Waiting for registration to propagate...")
                time.sleep(10)
            else:
                print_warning("Site registration failed, continuing with unregistered site")
        
        run_ip_block_tests()
        run_instance_type_tests()
        run_machine_tests()
        run_expected_machine_tests()
        
        # 4. Tenant resource tests (require tenant role)
        run_allocation_tests()
        run_vpc_tests()
        run_vpc_prefix_tests()
        run_subnet_tests()
        run_operating_system_tests()
        run_ssh_key_tests()
        run_ssh_key_group_tests()
        run_infiniband_partition_tests()
        run_network_security_group_tests()
        
        # 5. Cleanup - delete all created resources
        run_cleanup_tests()
        
        # Print summary
        success = stats.print_summary()
        
        # Exit with appropriate code
        sys.exit(0 if success else 1)
        
    except KeyboardInterrupt:
        print("\n\nTest interrupted by user")
        sys.exit(130)
        
    finally:
        # Cleanup port forwards
        if keycloak_pf:
            keycloak_pf.terminate()
            try:
                keycloak_pf.wait(timeout=5)
            except:
                pass
        if api_pf:
            api_pf.terminate()
            try:
                api_pf.wait(timeout=5)
            except:
                pass

if __name__ == "__main__":
    main()

