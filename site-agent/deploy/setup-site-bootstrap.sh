#!/usr/bin/env bash
# -------------------------------------------------------------------
# setup-site-bootstrap.sh — Bootstrap site-agent for a given SITE_UUID
#
# Assumptions:
#   - nettools-pod.yaml is in the same directory as this script
#   - site.sql (with INSERT INTO "public"."site") is in the same dir,
#     OR SITE_UUID is provided via env var.
#   - Uses current kubectl context unless KUBECTL_CONTEXT is set.
#
# Does:
#   - Ensures nettools-pod is running
#   - Calls site-manager to create the Site CR
#   - Fetches CA cert from credsmgr
#   - Reads OTP from Site CR
#   - Creates bootstrap-info and temporal-cert secrets in
#     namespace elektra-site-agent
# -------------------------------------------------------------------

set -eEuo pipefail
[ -n "${DEBUG_JUST:-}" ] && set -xv

die() { echo "❌  $*" >&2; exit 1; }
banner() { printf '\n\033[1;36m%s\033[0m\n' "$*"; }

# kubectl helper with optional context
k() {
  if [[ -n "${KUBECTL_CONTEXT:-}" ]]; then
    kubectl --context "$KUBECTL_CONTEXT" "$@"
  else
    kubectl "$@"
  fi
}

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

NETTOOLS_YAML="${SCRIPT_DIR}/nettools-pod.yaml"
SITE_SQL_PATH="${SCRIPT_DIR}/site.sql"

[[ -f "${NETTOOLS_YAML}" ]] || die "nettools-pod.yaml not found at ${NETTOOLS_YAML}"

# -------------------------------------------------------------------
# 1) Determine SITE_UUID
# -------------------------------------------------------------------

if [[ -z "${SITE_UUID:-}" ]]; then
  if [[ ! -f "${SITE_SQL_PATH}" ]]; then
    die "SITE_UUID not set and ${SITE_SQL_PATH} does not exist. Set SITE_UUID or create site.sql first."
  fi

  SITE_UUID=$(
    awk '
      BEGIN{IGNORECASE=1; inblk=0}
      /INSERT[[:space:]]+INTO[[:space:]]+"public"\."site"/ { inblk=1; next }
      inblk && match($0, /[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}/) {
        print substr($0,RSTART,RLENGTH); exit
      }' "${SITE_SQL_PATH}" || true
  )
fi

SITE_UUID=$(printf '%s' "${SITE_UUID}" | tr -d '\r\n\t ')
if [[ ! "${SITE_UUID}" =~ ^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$ ]]; then
  die "Could not determine a valid SITE_UUID (got: '${SITE_UUID}')."
fi

echo "Using SITE_UUID=${SITE_UUID}"

# -------------------------------------------------------------------
# 2) Ensure nettools-pod exists and is Running
# -------------------------------------------------------------------

banner "Ensuring nettools-pod exists…"
k apply -f "${NETTOOLS_YAML}"

for i in {1..60}; do
  PHASE=$(k get pod nettools-pod -o jsonpath='{.status.phase}' 2>/dev/null || true)
  if [[ "${PHASE}" == "Running" ]]; then
    echo "nettools-pod is Running."
    break
  fi
  if [[ $i -eq 60 ]]; then
    die "nettools-pod not ready after 60s (last phase: ${PHASE})"
  fi
  sleep 1
done

# -------------------------------------------------------------------
# 3) Call site-manager to create the site
# -------------------------------------------------------------------

JSON_PAYLOAD=$(printf '{"siteuuid":"%s","provider":"ip1","fcorg":"ip1"}' "${SITE_UUID}")

banner "Creating site via site-manager…"
k exec nettools-pod -- \
  curl -k -s -X POST \
    'https://sitemgr.cloud-site-manager:8100/v1/site' \
    -H 'Content-Type: application/json' \
    -d "${JSON_PAYLOAD}"

echo # newline

# -------------------------------------------------------------------
# 4) Grab CA cert from credsmgr
# -------------------------------------------------------------------

banner "Fetching CA certificate from credsmgr…"
k exec nettools-pod -- \
  curl -k -s https://credsmgr.cert-manager:8000/v1/pki/ca/pem \
  -o /tmp/cacert.pem

k cp nettools-pod:/tmp/cacert.pem /tmp/cacert.pem

# -------------------------------------------------------------------
# 5) Obtain OTP from the Site CR
# -------------------------------------------------------------------

banner "Reading OTP from Site CR…"
OTP=$(k get site "site-${SITE_UUID}" -n cloud-site-manager \
        -o jsonpath='{.status.otp.passcode}')

if [[ -z "${OTP}" ]]; then
  die "Failed to read OTP from Site CR site-${SITE_UUID} in namespace cloud-site-manager."
fi

echo "OTP=${OTP}"

# -------------------------------------------------------------------
# 6) Ensure namespace + create bootstrap-info & temporal-cert secrets
# -------------------------------------------------------------------

NAMESPACE="elektra-site-agent"

banner "Ensuring namespace ${NAMESPACE} exists…"
if ! k get namespace "${NAMESPACE}" >/dev/null 2>&1; then
  k create namespace "${NAMESPACE}"
fi

banner "Creating / updating bootstrap-info secret in ${NAMESPACE}…"
k -n "${NAMESPACE}" create secret generic bootstrap-info \
  --from-literal=site-uuid="${SITE_UUID}" \
  --from-literal=otp="${OTP}" \
  --from-literal=creds-url="https://sitemgr.cloud-site-manager.svc.cluster.local:8100/v1/sitecreds" \
  --from-file=cacert="/tmp/cacert.pem" \
  --dry-run=client -o yaml | k apply -f -

banner "Creating / updating temporal-cert secret in ${NAMESPACE}…"
k -n "${NAMESPACE}" create secret generic temporal-cert \
  --from-literal=cacertificate="" \
  --from-literal=certificate="" \
  --from-literal=key="" \
  --from-literal=otp="" \
  --dry-run=client -o yaml | k apply -f -

banner "🎉  Site-agent bootstrap secrets created. No restart performed."
