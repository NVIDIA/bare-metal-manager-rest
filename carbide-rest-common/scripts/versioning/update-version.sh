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

# Grab any additional settings
VERSIONING_DIR=${SCRIPT_DIR:-$(dirname "$0")}
if [ -f "${VERSIONING_DIR}/settings.sh" ]; then
    source "${VERSIONING_DIR}/settings.sh"
fi

# Set global defaults
if [ "${#TARGET_DIRS[@]}" == "0" ]; then
    TARGET_DIRS=(deploy)
fi

MAIN_BRANCH=${MAIN_BRANCH:-main}
RESTORE_STASH=true

USAGE="
usage:
    ${0} [--bump-level <major|minor|patch>] [--new-version <x.x.x>]
"

# Parse arguments
options_used=0;
while [ "${1}" != "" ];do
    case "${1}" in
        --bump-level|-l)    shift; BUMP_LEVEL="${1}"; options_used=$((options_used+1)); shift;;
        --new-version|-v)   shift; NEW_VERSION="${1}"; options_used=$((options_used+1)); shift;;
        --help|-h)
            echo -e "${USAGE}"
            exit 0
        ;;
        *)
            echo "unknown option: ${1}"
            echo -e "${USAGE}"
            exit 1
        ;;
    esac
done;

# Make sure at least one option was used
if [ "${options_used}" -lt "1" ]; then
    echo "too few arguments"
    echo -e "${USAGE}"
    exit 1
fi

# Grab the current branch
current_branch=$(git rev-parse --abbrev-ref HEAD)

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

# Grab the current version
# Slice it up into the version tag parts
IFS=. read -r v_major v_minor v_patch < VERSION

export NEW_V_MAJOR="${v_major}"
export NEW_V_MINOR="${v_minor}"
export NEW_V_PATCH="${v_patch}"

# If a desired new version was explicity passed in, use that.
if [ "${NEW_VERSION}" != "" ]; then
    echo
    echo "using explicitly requested new version of ${NEW_VERSION}"
    echo

    IFS=. read -r NEW_V_MAJOR NEW_V_MINOR NEW_V_PATCH <<< "${NEW_VERSION}"
else
    # Increment the relevant version section.
    if [ "${BUMP_LEVEL}" == "patch" ] ; then
        NEW_V_PATCH="$((v_patch+1))"
    elif [ "${BUMP_LEVEL}" == "minor" ]; then
        NEW_V_MINOR="$((v_minor+1))"
        NEW_V_PATCH="0"
    elif [ "${BUMP_LEVEL}" == "major" ]; then
        NEW_V_MAJOR="$((v_major+1))"
        NEW_V_MINOR="0"
        NEW_V_PATCH="0"
    else
        echo "error: bump level must be one of major, minor, patch"
        exit 1
    fi
fi

# If the updated patch version is over 99, then
# it's time to rollover to the next MINOR instead.
if [ "${NEW_V_PATCH}" -gt "99" ]; then
    NEW_V_MINOR="$((NEW_V_MINOR+1))"
    NEW_V_PATCH="0"
fi

# If the updated minor version is over 99, either from a
# direct bump or because patch went over 99 and minor was already
# 99, then it's time to rollover to the next MAJOR instead.
if [ "${NEW_V_MINOR}" -gt "99" ]; then
    NEW_V_MAJOR="$((NEW_V_MAJOR+1))"
    NEW_V_MINOR="0"
    NEW_V_PATCH="0"
fi

export NEW_VERSION="${NEW_V_MAJOR}.${NEW_V_MINOR}.${NEW_V_PATCH}"

echo
echo "updating version to ${NEW_VERSION}"
echo

# Generate a new branch name for an MR
new_branch="version-${NEW_VERSION}-$(date +%s)"

# Set up a trap for EXIT to clean-up
restore_on_exit(){
    set +e

    # Switch back to wherever we were.
    ${DEBUG} git checkout "${current_branch}"

    if ${RESTORE_STASH}; then
        ${DEBUG}  git stash apply
    fi
    ${DEBUG}  git branch -D "${new_branch}"
}
trap restore_on_exit EXIT

# Deal with expected MacOS vs Linux...
SED_LINE='s/'"${v_major}"'\([-.]\)'"${v_minor}"'\([-.]\)'"${v_patch}"'/'"${NEW_V_MAJOR}"'\1'"${NEW_V_MINOR}"'\2'"${NEW_V_PATCH}"'/g'
if [ "$(uname)" == "Darwin" ]; then
    # Require gnu-sed.
    if ! [ -x "$(command -v gsed)" ]; then
        echo "Error: 'gsed' is not installed." >&2
        echo "If you are using Homebrew, install with 'brew install gnu-sed'." >&2
        exit 1
    fi
    SED_COMMAND="gsed -i -e"
else
    SED_COMMAND="sed -i -e"
fi

COMMIT_TARGETS=()

# Find everything that needs an update.
for target_dir in "${TARGET_DIRS[@]}";do
    if [ -d "${target_dir}" ]; then
        for f in $(find "${target_dir}" -type f);do
            ${DEBUG} ${SED_COMMAND} "${SED_LINE}" $f
        done
        COMMIT_TARGETS+=( "${target_dir}" )
    else
        echo
        echo "'${target_dir}' directory not found; skipping..."
        echo
    fi
done

# Update the version file itself
if [ "${DEBUG}" == "" ]; then
    echo "${NEW_VERSION}" > VERSION
fi

${DEBUG} git checkout -b "${new_branch}"
${DEBUG} git commit VERSION ${COMMIT_TARGETS[@]} -m"[Version] Bump version to ${NEW_VERSION}"
${DEBUG} git push origin "${new_branch}"
