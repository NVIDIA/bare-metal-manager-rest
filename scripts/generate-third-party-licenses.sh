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

# Generate THIRD-PARTY-LICENSES from Go module dependencies.
# Run from repository root. After running, the temporal-helm-charts entry is
# prepended (manually maintained; not a Go module). All other entries come
# from "go list -m all" after "go mod download".

set -euo pipefail

REPO_ROOT="$(cd "$(dirname "$0")/.." && pwd)"
cd "$REPO_ROOT"
OUTPUT="${1:-THIRD-PARTY-LICENSES}"

echo "Downloading Go modules..."
go mod download

MAIN_MODULE="$(go list -m)"
TEMPORAL_HELM_BLOCK='# When regenerating this file from go.mod (or your license/compliance tool), re-add
# the temporal-helm-charts entry below manually (it is not a Go module; see
# https://github.com/temporalio/helm-charts, MIT, used in temporal-helm/).
root_name = "bare-metal-manager-rest"

[[third_party_libraries]]
package_name = "temporal-helm-charts"
package_version = "0.35.0"
repository = "https://github.com/temporalio/helm-charts"
license = "MIT"
[[third_party_libraries.licenses]]
license = "MIT"
text = """
The MIT License

Copyright (c) 2020 Temporal Technologies Inc.  All rights reserved.

Copyright (c) 2020 Uber Technologies, Inc.

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
"""

'

# Detect license type from file contents (first 4K)
detect_license() {
  local f="$1"
  if ! [ -f "$f" ] || ! [ -r "$f" ]; then
    echo "Unknown"
    return
  fi
  local head
  head="$(head -c 4096 "$f")"
  if echo "$head" | grep -qi "Apache License"; then
    echo "Apache-2.0"
  elif echo "$head" | grep -qi "MIT License"; then
    echo "MIT"
  elif echo "$head" | grep -qi "BSD 3-Clause\|BSD-3-Clause"; then
    echo "BSD-3-Clause"
  elif echo "$head" | grep -qi "BSD 2-Clause\|BSD-2-Clause"; then
    echo "BSD-2-Clause"
  elif echo "$head" | grep -qi "Mozilla Public License\|MPL"; then
    echo "MPL-2.0"
  else
    echo "Unknown"
  fi
}

# Find first LICENSE* or LICENCE* in dir
find_license_file() {
  local dir="$1"
  [ -d "$dir" ] || return 1
  for name in LICENSE LICENSE.md LICENCE LICENCE.md; do
    if [ -f "$dir/$name" ]; then
      echo "$dir/$name"
      return 0
    fi
  done
  return 1
}

# Escape content for TOML multiline string: \ -> \\, " -> \"
escape_toml_text() {
  sed 's/\\/\\\\/g; s/"/\\"/g'
}

# Last path component for package_name (e.g. connectrpc.com/connect -> connect)
package_name_from_path() {
  local path="$1"
  echo "${path##*/}"
}

# Repository URL from module path
repo_url_from_path() {
  local path="$1"
  if [[ "$path" == github.com/* ]]; then
    echo "https://$path"
  elif [[ "$path" == gitlab.com/* ]]; then
    echo "https://$path"
  elif [[ "$path" == gitea.com/* ]]; then
    echo "https://$path"
  elif [[ "$path" == buf.build/* ]]; then
    echo "https://$path"
  else
    echo "https://$path"
  fi
}

{
  printf '%s' "$TEMPORAL_HELM_BLOCK"

  # Sort by module path then version so order is stable and diffs are minimal on regeneration
  go list -m all | sort -k1,1 -k2,2 | while read -r path version; do
    [ -n "$path" ] || continue
    [ "$path" = "$MAIN_MODULE" ] && continue
    [ -n "$version" ] || version="(devel)"

    pkgname="$(package_name_from_path "$path")"
    repo="$(repo_url_from_path "$path")"

    dir=""
    dir="$(go list -m -f '{{.Dir}}' "$path@$version" 2>/dev/null)" || true
    license_type="Unknown"
    license_file=""
    if [ -n "$dir" ] && license_file="$(find_license_file "$dir")"; then
      license_type="$(detect_license "$license_file")"
    fi

    printf '\n[[third_party_libraries]]\n'
    printf 'package_name = "%s"\n' "$pkgname"
    printf 'package_version = "%s"\n' "$version"
    printf 'repository = "https://%s"\n' "$path"
    printf 'license = "%s"\n' "$license_type"

    if [ -n "$license_file" ] && [ -r "$license_file" ]; then
      # Include full license text for every dependency that has a license file
      printf '[[third_party_libraries.licenses]]\n'
      printf 'license = "%s"\n' "$license_type"
      printf 'text = """\n'
      cat "$license_file" | escape_toml_text
      printf '\n"""\n'
    fi
  done
} > "$OUTPUT.tmp"

mv "$OUTPUT.tmp" "$OUTPUT"
echo "Wrote $OUTPUT"
