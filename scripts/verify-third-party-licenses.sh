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

# Verify that every Go module dependency (go list -m all) is listed in
# THIRD-PARTY-LICENSES. Exits 0 if all are named, 1 if any are missing.
# Run from repository root.

set -euo pipefail

REPO_ROOT="$(cd "$(dirname "$0")/.." && pwd)"
cd "$REPO_ROOT"
LICENSE_FILE="${1:-THIRD-PARTY-LICENSES}"

if [ ! -f "$LICENSE_FILE" ]; then
  echo "Error: $LICENSE_FILE not found" >&2
  exit 1
fi

MAIN_MODULE="$(go list -m)"

# Build set of module paths from THIRD-PARTY-LICENSES (repository = "https://path" -> path)
# Use temp files for portability (no associative arrays in baseline sh).
LISTED="$(mktemp)"
DEPS="$(mktemp)"
trap 'rm -f "$LISTED" "$DEPS"' EXIT
grep -E '^repository = "https://' "$LICENSE_FILE" | sed 's/^repository = "https:\/\///; s/".*$//' | sort -u > "$LISTED"

go list -m all > "$DEPS"

MISSING=""
while read -r path version; do
  [ -n "$path" ] || continue
  [ "$path" = "$MAIN_MODULE" ] && continue
  if ! grep -Fxq "$path" "$LISTED"; then
    MISSING="${MISSING:+$MISSING$newline}$path"
    newline='
'
  fi
done < "$DEPS"

if [ -n "$MISSING" ]; then
  echo "The following Go module dependencies are not listed in $LICENSE_FILE:" >&2
  echo "$MISSING" | while read -r line; do echo "  - $line"; done >&2
  echo "Run: make generate-third-party-licenses" >&2
  exit 1
fi

echo "All Go module dependencies are named in $LICENSE_FILE."
