#!/usr/bin/env bash
# SPDX-FileCopyrightText: Copyright (c) 2026 NVIDIA CORPORATION & AFFILIATES. All rights reserved.
# SPDX-License-Identifier: Apache-2.0
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
# http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

# -------------------------------------------------------------------
# gen-site-sql.sh  —  Generate site.sql with a fresh site UUID
#
# Usage:
#   ./gen-site-sql.sh
#
# Behavior:
#   - Generates a UUID (or uses SITE_UUID if provided)
#   - Uses directory name as SITE_NAME (or SITE_NAME env if provided)
#   - Writes ./site.sql (same directory as this script)
# -------------------------------------------------------------------

set -eEuo pipefail
die() { echo "❌  $*" >&2; exit 1; }

# Directory where this script lives
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

SITE_SQL_PATH="${SCRIPT_DIR}/site.sql"

# Site name: can override via SITE_NAME env if you want
SITE_NAME="${SITE_NAME:-$(basename "${SCRIPT_DIR}")}"

# Generate or reuse SITE_UUID
if [[ -z "${SITE_UUID:-}" ]]; then
  if ! command -v uuidgen >/dev/null 2>&1; then
    die "uuidgen not found. Install it or export SITE_UUID before running."
  fi
  SITE_UUID=$(uuidgen | tr 'A-Z' 'a-z')
fi

SITE_UUID=$(printf '%s' "${SITE_UUID}" | tr -d '\r\n\t ')
if [[ ! "${SITE_UUID}" =~ ^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$ ]]; then
  die "Generated SITE_UUID is not a valid UUID (got: '${SITE_UUID}')."
fi

echo "🔖  Site name : ${SITE_NAME}"
echo "🔑  SITE_UUID : ${SITE_UUID}"
echo "📌  First-state log: SITE_UUID=${SITE_UUID}"

# Escape site name for SQL
SQL_SITE_NAME=${SITE_NAME//\'/\'\'}

cat > "${SITE_SQL_PATH}" <<EOF
INSERT INTO "public"."infrastructure_provider" ("id", "name", "display_name", "org", "created", "updated", "deleted", "created_by") VALUES
('02306792-249b-49e3-b90b-fd57058d22a2', 'default', 'default', 'nvidia', '2022-09-29 00:40:10.12539+00', '2022-09-29 00:40:10.12539+00', NULL, '${SITE_UUID}');

INSERT INTO "public"."site" ("id", "name", "display_name", "description", "org", "infrastructure_provider_id", "site_controller_version", "site_agent_version", "registration_token", "registration_token_expiration", "is_infinity_enabled", "is_serial_console_enabled", "status", "created", "updated", "deleted", "created_by") VALUES
('${SITE_UUID}', '${SQL_SITE_NAME}', '${SQL_SITE_NAME}', NULL, 'nvidia', '02306792-249b-49e3-b90b-fd57058d22a2', NULL, NULL, '4f3dea95-3108-44eb-9e21-301b9aae8a2c', '2022-10-06 22:15:12.56311+00', false, false, 'Pending', '2022-10-04 22:15:12.562019+00', '2022-10-04 22:15:12.562019+00', NULL, '${SITE_UUID}');
EOF

echo "✅  Generated ${SITE_SQL_PATH}"
