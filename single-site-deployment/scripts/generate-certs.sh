#!/bin/bash

# Generate Self-Signed Certificates for Development
# Creates CA and service certificates for local deployment

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ROOT_DIR="$(cd "$SCRIPT_DIR/../.." && pwd)"
CERTS_DIR="$ROOT_DIR/build/certs"

echo "========================================"
echo "Generating Development Certificates"
echo "========================================"
echo ""
echo "Output directory: $CERTS_DIR"
echo ""

# Create certs directory
mkdir -p "$CERTS_DIR"
cd "$CERTS_DIR"

# Generate CA private key
if [ ! -f ca-key.pem ]; then
    echo "Generating CA private key..."
    openssl genrsa -out ca-key.pem 4096
    echo "[SUCCESS] CA private key generated"
fi

# Generate CA certificate
if [ ! -f ca-cert.pem ]; then
    echo "Generating CA certificate..."
    openssl req -new -x509 -days 3650 -key ca-key.pem -out ca-cert.pem \
        -subj "/C=US/ST=California/L=Santa Clara/O=NVIDIA/OU=Carbide/CN=Carbide Root CA"
    echo "[SUCCESS] CA certificate generated"
fi

# Function to generate service certificate
generate_service_cert() {
    local service_name=$1
    local dns_names=$2
    
    echo ""
    echo "Generating certificate for: $service_name"
    
    # Generate private key
    if [ ! -f "${service_name}-key.pem" ]; then
        openssl genrsa -out "${service_name}-key.pem" 2048
        echo "  [SUCCESS] Private key generated"
    fi
    
    # Create config file for SAN
    cat > "${service_name}.cnf" <<EOF
[req]
req_extensions = v3_req
distinguished_name = req_distinguished_name

[req_distinguished_name]

[v3_req]
basicConstraints = CA:FALSE
keyUsage = nonRepudiation, digitalSignature, keyEncipherment
subjectAltName = @alt_names

[alt_names]
EOF
    
    # Add DNS names
    local counter=1
    for dns in $dns_names; do
        echo "DNS.${counter} = ${dns}" >> "${service_name}.cnf"
        counter=$((counter + 1))
    done
    
    # Generate CSR
    openssl req -new -key "${service_name}-key.pem" -out "${service_name}.csr" \
        -subj "/C=US/ST=California/L=Santa Clara/O=NVIDIA/OU=Carbide/CN=${service_name}" \
        -config "${service_name}.cnf"
    
    # Sign certificate with CA
    openssl x509 -req -in "${service_name}.csr" -CA ca-cert.pem -CAkey ca-key.pem \
        -CAcreateserial -out "${service_name}-cert.pem" -days 365 \
        -extensions v3_req -extfile "${service_name}.cnf"
    
    echo "  [SUCCESS] Certificate signed by CA"
    
    # Cleanup
    rm "${service_name}.csr" "${service_name}.cnf"
}

# Generate certificates for services
generate_service_cert "temporal" "temporal temporal.carbide-system temporal.carbide-system.svc temporal.carbide-system.svc.cluster.local localhost"

generate_service_cert "site-manager" "site-manager site-manager.carbide-system site-manager.carbide-system.svc site-manager.carbide-system.svc.cluster.local localhost"

generate_service_cert "cert-manager" "cert-manager cert-manager.carbide-system cert-manager.carbide-system.svc cert-manager.carbide-system.svc.cluster.local localhost"

generate_service_cert "cloud-api" "cloud-api cloud-api.carbide-system cloud-api.carbide-system.svc cloud-api.carbide-system.svc.cluster.local carbide-api carbide-api.carbide-system carbide-api.carbide-system.svc carbide-api.carbide-system.svc.cluster.local localhost"

generate_service_cert "site-agent" "carbide-site-agent carbide-site-agent.carbide-site carbide-site-agent.carbide-site.svc carbide-site-agent.carbide-site.svc.cluster.local carbide-site-agent-0.carbide-site-agent.carbide-site.svc.cluster.local localhost"

# Generate a generic bootstrap OTP token
echo ""
echo "Generating bootstrap OTP token..."
BOOTSTRAP_OTP=$(openssl rand -hex 16)
echo "$BOOTSTRAP_OTP" > bootstrap-otp.txt
echo "[SUCCESS] Bootstrap OTP: $BOOTSTRAP_OTP"

echo ""
echo "========================================"
echo "Certificate Generation Complete!"
echo "========================================"
echo ""
echo "Generated files in $CERTS_DIR:"
ls -lh "$CERTS_DIR"
echo ""
echo "Certificate Summary:"
echo "  CA Certificate:     ca-cert.pem"
echo "  CA Private Key:     ca-key.pem"
echo "  Bootstrap OTP:      bootstrap-otp.txt"
echo "  Service Certs:      *-cert.pem, *-key.pem"
echo ""

