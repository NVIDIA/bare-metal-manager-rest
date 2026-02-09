#!/usr/bin/env bash
# -------------------------------------------------------------------
# create-site-in-db.sh — Apply site.sql and ensure Temporal namespace
#
# Usage:
#   ./create-site-in-db.sh            # uses ./site.sql (same dir as script)
#   ./create-site-in-db.sh /path/to/site.sql
#
# Notes:
#   - Uses current kubectl context by default.
#   - You can override with: export KUBECTL_CONTEXT=your-context
# -------------------------------------------------------------------

set -eEuo pipefail
[ -n "${DEBUG_JUST:-}" ] && set -xv

die() { echo "❌  $*" >&2; exit 1; }
banner() { printf '\n\033[1;36m%s\033[0m\n' "$*"; }

# Helper: kubectl with optional context
k() {
  if [[ -n "${KUBECTL_CONTEXT:-}" ]]; then
    kubectl --context "$KUBECTL_CONTEXT" "$@"
  else
    kubectl "$@"
  fi
}

# Determine site.sql path
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

if [[ $# -gt 0 ]]; then
  SITE_SQL_PATH="$1"
else
  SITE_SQL_PATH="${SCRIPT_DIR}/site.sql"
fi

[[ -f "${SITE_SQL_PATH}" ]] || die "site.sql not found at ${SITE_SQL_PATH}"

# --------- extract SITE_UUID from site.sql --------------------------

SITE_UUID=$(
  awk '
    BEGIN{IGNORECASE=1; inblk=0}
    /INSERT[[:space:]]+INTO[[:space:]]+"public"\."site"/ { inblk=1; next }
    inblk && match($0, /[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}/) {
      print substr($0,RSTART,RLENGTH); exit
    }' "${SITE_SQL_PATH}" || true
)

SITE_UUID=$(printf '%s' "${SITE_UUID}" | tr -d '\r\n\t ')
if [[ ! "${SITE_UUID}" =~ ^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$ ]]; then
  die "Could not reliably determine SITE_UUID from ${SITE_SQL_PATH} (got: '${SITE_UUID}')."
fi

echo "Using SITE_UUID=${SITE_UUID}"

# --------- DB name + master pod ------------------------------------

# NOTE: key is dbName, not name
DB_NAME=$(k get configmap cloud-db-config -n cloud-db -o jsonpath='{.data.name}')
if [[ -z "${DB_NAME}" ]]; then
  die "Could not read DB name from configmap cloud-db-config in namespace cloud-db."
fi

MASTER_POD=$(k -n postgres get pods -l spilo-role=master -o jsonpath='{.items[0].metadata.name}')
[[ -n "${MASTER_POD}" ]] || die "Could not locate master Postgres pod"

# --------- Insert site row if needed --------------------------------

banner "⏳  Checking if site rows already exist…"
EXISTS=$(k -n postgres exec "${MASTER_POD}" -- \
          psql -tAc "SELECT 1 FROM site WHERE id = '${SITE_UUID}'" -U postgres -d "${DB_NAME}" 2>/dev/null || true)
EXISTS=$(printf '%s' "${EXISTS}" | tr -d '[:space:]')

if [[ "${EXISTS}" != "1" ]]; then
  echo "Inserting site data into Postgres…"
  k -n postgres exec -i "${MASTER_POD}" -- \
    psql -U postgres -d "${DB_NAME}" < "${SITE_SQL_PATH}"
else
  echo "Site rows for SITE_UUID=${SITE_UUID} already present — skipping SQL import."
fi

# --------- Ensure Temporal namespace exists -------------------------

banner "⏳  Ensuring Temporal namespace ${SITE_UUID} exists…"

ADM_POD=$(k -n temporal get pods \
          -l app.kubernetes.io/component=admintools -o jsonpath='{.items[0].metadata.name}')

[[ -n "${ADM_POD}" ]] || die "Could not locate Temporal admintools pod"

if ! k exec -n temporal "$ADM_POD" \
      -- tctl --ns "${SITE_UUID}" namespace describe &>/dev/null; then
  k exec -n temporal "$ADM_POD" \
    -- tctl --ns "${SITE_UUID}" namespace register
  echo "✓  Temporal namespace ${SITE_UUID} registered."
else
  echo "Temporal namespace ${SITE_UUID} already exists — skipping."
fi

banner "🎉  Site DB rows & Temporal namespace ensured for SITE_UUID=${SITE_UUID}."
