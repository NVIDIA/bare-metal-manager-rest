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

set -eu -o pipefail
readonly script_dir="$(cd "$(dirname "$0")";pwd)"
# This script contains common shell library functions.

# Polling for command stdout until command succeeded with given
# interval and timeout.
#
# Usage:
#   poll_immediate <interval_seconds> <timeout-seconds> <cmd> [<cmd_args>]
#
# Example usage:
#   poll_immediate 5 120 curl -ksS https://master.egx.nvidia.com/healthz
#
# Contract for the command used in polling:
#   - Exit with non-zero code if failed, the command will be retried
#     if timeout is not reached, therefore command is responsible make
#     sure itself is retriable.
#   - Only produce output to stdout when command succeed, otherwise
#     output to stderr.
#
# Inspired by:
# https://godoc.org/gopkg.in/kubernetes/client-go.v1/1.5/pkg/util/wait#PollImmediate
poll_immediate() {
  local -r interval="$1" timeout="$2"
  shift 2

  local start_time now deadline
  start_time="$(date +%s)"
  now="${start_time}"
  deadline="$((now+timeout))"
  while [[ "${now}" -lt "${deadline}" ]]; do
    if "$@"; then
      echo "'$*' SUCCEEDED at $((now-start_time)) seconds since start !" >&2
      return 0
    fi
    echo "'$*' FAILED at $((now-start_time)) seconds since start, will retry in ${interval} seconds ..." >&2

    sleep "${interval}"
    now="$(date +%s)"
  done

  echo "Polling '$*' FAILED: timeout ${timeout} seconds exceeded !" >&2
  return 1
}

# displays a banner for a test case
tc_banner() {
  echo ""
  echo "||"
  echo "|| Test case: \"$1\""
  echo "||"
  echo ""
}

banner() {
  echo ""
  echo "||"
  echo "|| \"$1\""
  echo "||"
  echo ""
}
