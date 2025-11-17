#!/bin/bash
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

set -e

VERSIONING_DIR=${SCRIPT_DIR:-$(dirname "$0")}

if [ -f "${VERSIONING_DIR}/settings.sh" ]; then
    source "${VERSIONING_DIR}/settings.sh"
fi

MAIN_BRANCH=${MAIN_BRANCH:-main}
RESTORE_STASH=true

# Grab the current branch
current_branch=$(git rev-parse --abbrev-ref HEAD)

# Set up a trap for EXIT to clean-up
restore_on_exit(){
    set +e

    # Switch back to wherever we were.
    ${DEBUG} git checkout "${current_branch}"
    if ${RESTORE_STASH}; then
        ${DEBUG}  git stash apply
    fi
}
trap restore_on_exit EXIT

# If we're not already on main, stash, switch, and push up the changes.
if [ "${current_branch}" != "${MAIN_BRANCH}" ]; then
    if [ "$(git stash)" == "No local changes to save" ]; then
        RESTORE_STASH=false
    fi 
    ${DEBUG}  git checkout "${MAIN_BRANCH}"
else
    if [ "$(git stash)" == "No local changes to save" ]; then
        RESTORE_STASH=false
    fi 
fi

# Freshen up
${DEBUG} git pull	

echo
echo "adding version tag v$(cat VERSION)"
echo

# Create and push the tag
${DEBUG} git tag "v$(cat VERSION)"
${DEBUG} git push --tags

