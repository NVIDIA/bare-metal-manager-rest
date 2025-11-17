#!/usr/bin/env python3
"""
Deployment Diagnostics Script
Systematically checks for deployment issues and gathers debugging information
"""

import subprocess
import sys
import json
from pathlib import Path
from typing import Dict, List

# Configuration
NAMESPACES = ["carbide-system", "carbide-site"]
CLUSTER_NAME = "carbide-single-site"


def print_header(message: str):
    """Print a section header"""
    print(f"\n{'='*80}")
    print(f"{message}")
    print(f"{'='*80}\n")


def print_section(message: str):
    """Print a subsection header"""
    print(f"\n{'-'*80}")
    print(f"{message}")
    print(f"{'-'*80}")


def run_command(cmd: List[str], capture: bool = True) -> str:
    """Run a command and return output"""
    try:
        if capture:
            result = subprocess.run(cmd, capture_output=True, text=True, check=False)
            return result.stdout
        else:
            subprocess.run(cmd, check=False)
            return ""
    except Exception as e:
        return f"Error running command: {e}"


def check_cluster_exists():
    """Check if the cluster exists"""
    print_header("STEP 1: Cluster Status")
    
    output = run_command(["kind", "get", "clusters"])
    if CLUSTER_NAME in output:
        print(f"[OK] Cluster '{CLUSTER_NAME}' exists")
        
        # Get cluster info
        print("\nCluster Info:")
        run_command(["kubectl", "cluster-info", "--context", f"kind-{CLUSTER_NAME}"], capture=False)
        return True
    else:
        print(f"[ERROR] Cluster '{CLUSTER_NAME}' not found")
        print(f"Available clusters: {output.strip()}")
        return False


def get_pod_status():
    """Get status of all pods"""
    print_header("STEP 2: Pod Status Overview")
    
    pod_issues = {}
    
    for namespace in NAMESPACES:
        print_section(f"Namespace: {namespace}")
        
        # Get pods in JSON format for parsing
        output = run_command([
            "kubectl", "get", "pods",
            "-n", namespace,
            "-o", "json"
        ])
        
        try:
            data = json.loads(output)
            pods = data.get('items', [])
            
            if not pods:
                print(f"  No pods found in {namespace}")
                continue
            
            # Show table view
            run_command([
                "kubectl", "get", "pods",
                "-n", namespace,
                "-o", "wide"
            ], capture=False)
            
            # Analyze each pod
            for pod in pods:
                pod_name = pod['metadata']['name']
                pod_status = pod['status']['phase']
                
                # Check container statuses
                container_statuses = pod['status'].get('containerStatuses', [])
                
                for container in container_statuses:
                    container_name = container['name']
                    ready = container.get('ready', False)
                    state = container.get('state', {})
                    
                    # Check for issues
                    if not ready or pod_status != 'Running':
                        issue_key = f"{namespace}/{pod_name}/{container_name}"
                        
                        issue_info = {
                            'namespace': namespace,
                            'pod': pod_name,
                            'container': container_name,
                            'phase': pod_status,
                            'ready': ready,
                            'state': state
                        }
                        
                        # Determine issue type
                        if 'waiting' in state:
                            issue_info['issue_type'] = state['waiting'].get('reason', 'Unknown')
                            issue_info['message'] = state['waiting'].get('message', '')
                        elif 'terminated' in state:
                            issue_info['issue_type'] = 'Terminated'
                            issue_info['exit_code'] = state['terminated'].get('exitCode', 'unknown')
                            issue_info['reason'] = state['terminated'].get('reason', '')
                        else:
                            issue_info['issue_type'] = 'Unknown'
                        
                        pod_issues[issue_key] = issue_info
        
        except json.JSONDecodeError:
            print(f"  Error parsing pod data for {namespace}")
    
    return pod_issues


def diagnose_pod_issues(pod_issues: Dict):
    """Diagnose each problematic pod"""
    print_header("STEP 3: Detailed Pod Diagnostics")
    
    if not pod_issues:
        print("[OK] No pod issues detected!")
        return
    
    print(f"Found {len(pod_issues)} pod issues to diagnose\n")
    
    for issue_key, info in pod_issues.items():
        namespace = info['namespace']
        pod_name = info['pod']
        container_name = info['container']
        issue_type = info['issue_type']
        
        print_section(f"Issue: {issue_key}")
        print(f"  Pod: {pod_name}")
        print(f"  Container: {container_name}")
        print(f"  Namespace: {namespace}")
        print(f"  Phase: {info['phase']}")
        print(f"  Issue Type: {issue_type}")
        
        if 'message' in info and info['message']:
            print(f"  Message: {info['message'][:200]}...")
        
        if 'exit_code' in info:
            print(f"  Exit Code: {info['exit_code']}")
            print(f"  Reason: {info.get('reason', 'N/A')}")
        
        # Get pod description
        print(f"\n  Pod Description (key sections):")
        desc_output = run_command([
            "kubectl", "describe", "pod", pod_name,
            "-n", namespace
        ])
        
        # Extract relevant parts
        lines = desc_output.split('\n')
        in_events = False
        in_conditions = False
        
        for line in lines:
            if 'Conditions:' in line:
                in_conditions = True
                print(f"    {line}")
            elif 'Events:' in line:
                in_events = True
                in_conditions = False
                print(f"    {line}")
            elif in_conditions and line.strip() and not line.startswith('Volumes'):
                print(f"    {line}")
                if not line.startswith(' '):
                    in_conditions = False
            elif in_events:
                print(f"    {line}")
        
        # Get container logs
        print(f"\n  Container Logs (last 30 lines):")
        logs = run_command([
            "kubectl", "logs",
            pod_name,
            "-c", container_name,
            "-n", namespace,
            "--tail=30"
        ])
        
        if logs.strip():
            for line in logs.strip().split('\n')[-30:]:
                print(f"    {line}")
        else:
            print("    [No logs available]")
        
        # Check previous container logs if crashed
        if issue_type in ['CrashLoopBackOff', 'Error']:
            print(f"\n  Previous Container Logs (last 30 lines):")
            prev_logs = run_command([
                "kubectl", "logs",
                pod_name,
                "-c", container_name,
                "-n", namespace,
                "--previous",
                "--tail=30"
            ])
            
            if prev_logs.strip() and "error" not in prev_logs.lower()[:100]:
                for line in prev_logs.strip().split('\n')[-30:]:
                    print(f"    {line}")
            else:
                print("    [No previous logs available or error retrieving them]")
        
        print("\n")


def check_services():
    """Check service endpoints"""
    print_header("STEP 4: Service Status")
    
    for namespace in NAMESPACES:
        print_section(f"Services in {namespace}")
        run_command([
            "kubectl", "get", "svc",
            "-n", namespace
        ], capture=False)


def main():
    """Main diagnostic flow"""
    print_header("Carbide Deployment Diagnostics")
    
    # Check if cluster exists
    if not check_cluster_exists():
        print("\n[ERROR] Cannot proceed without a cluster")
        print("Run: python3 rebuild-and-deploy.py")
        sys.exit(1)
    
    # Get pod status
    pod_issues = get_pod_status()
    
    # Diagnose issues
    diagnose_pod_issues(pod_issues)
    
    # Check services
    check_services()
    
    print_header("Diagnostics Complete")
    
    print("Next Steps:")
    print("  1. Review the logs and errors above")
    print("  2. Fix code/configuration issues")
    print("  3. Rebuild affected services")
    print("  4. Reload images and restart deployments")
    print("\nQuick Commands:")
    print("  Rebuild all:     python3 rebuild-and-deploy.py")
    print("  Deploy only:     python3 rebuild-and-deploy.py --skip-build")
    print("  Watch pods:      kubectl get pods -n carbide-system -w")
    print("  Delete cluster:  kind delete cluster --name carbide-single-site")
    print()


if __name__ == "__main__":
    main()

