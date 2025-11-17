#! /usr/bin/env bash
# SPDX-FileCopyrightText: Copyright (c) 2021-2023 NVIDIA CORPORATION & AFFILIATES. All rights reserved.
# SPDX-License-Identifier: LicenseRef-NvidiaProprietary
#
# NVIDIA CORPORATION, its affiliates and licensors retain all intellectual
# property and proprietary rights in and to this material, related
# documentation and any modifications thereto. Any use, reproduction,
# disclosure or distribution of this material and related documentation
# without an express license agreement from NVIDIA CORPORATION or
# its affiliates is strictly prohibited.
#

gen_code() {
  local -r header_file="${repo_root}/scripts/boilerplate.go.txt"
  local -r pkg_name="gitlab-master.nvidia.com/nvmetal/cloud-site-manager"
  cd "${repo_root}"

  # Remove old generated file to avoid conflicts
  rm -f "${repo_root}/pkg/crds/v1/zz_generated.deepcopy.go"
  
  # Download only Kubernetes dependencies needed for code generation
  # Ignore private dependencies that aren't needed for CRD type generation
  export GOPROXY=${GOPROXY:-https://proxy.golang.org,direct}
  go mod download k8s.io/apimachinery || true
  go mod download k8s.io/client-go || true
  go mod download k8s.io/code-generator || true

  echo "Generating pkg/crds/v1/zz_generated.deepcopy.go"
  deepcopy-gen \
    --input-dirs "${pkg_name}/pkg/crds/v1" \
    -O zz_generated.deepcopy \
    --bounding-dirs "${pkg_name}/pkg/crds" \
    --go-header-file "${header_file}"

  local output_base
  output_base="$(mktemp -d)"
  rm -fr "${repo_root}/pkg/client"
  mkdir -p "${repo_root}/pkg/client"

  echo "Generating pkg/client/clientset"
  client-gen \
    --clientset-name versioned \
    --input-base "" \
    --input "${pkg_name}/pkg/crds/v1" \
    --output-package "${pkg_name}/pkg/client/clientset" \
    --go-header-file "${header_file}" \
    --output-base "${output_base}"
  mv "${output_base}/${pkg_name}/pkg/client/clientset" "${repo_root}/pkg/client/"

  echo "Generating pkg/client/listers"
  lister-gen \
    --input-dirs "${pkg_name}/pkg/crds/v1" \
    --output-package "${pkg_name}/pkg/client/listers" \
    --go-header-file "${header_file}" \
    --output-base "${output_base}"
  mv "${output_base}/${pkg_name}/pkg/client/listers" "${repo_root}/pkg/client/"

  echo "Generating pkg/client/informers"
  informer-gen \
    --input-dirs "${pkg_name}/pkg/crds/v1" \
    --versioned-clientset-package "${pkg_name}/pkg/client/clientset/versioned" \
    --listers-package "${pkg_name}/pkg/client/listers" \
    --output-package "${pkg_name}/pkg/client/informers" \
    --go-header-file "${header_file}" \
    --output-base "${output_base}"
  mv "${output_base}/${pkg_name}/pkg/client/informers" "${repo_root}/pkg/client/"

  rm -fr "${output_base}"
}
