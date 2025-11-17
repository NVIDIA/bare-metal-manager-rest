#!/bin/bash
# Script to prepare code for partner release
# This will remove all git history and create a fresh repository with a single "Release" commit

set -e

# Get the repository root directory
REPO_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

echo "======================================"
echo "Partner Release Preparation Script"
echo "======================================"
echo ""
echo "This script will:"
echo "  1. Remove all existing git history"
echo "  2. Initialize a fresh git repository"
echo "  3. Create a single commit called 'Release'"
echo ""
echo "WARNING: This will permanently delete all git history!"
echo "Repository root: ${REPO_ROOT}"
echo ""

# Check if we're in a git repository
if [ ! -d "${REPO_ROOT}/.git" ]; then
    echo "Error: Not in a git repository"
    exit 1
fi

# Prompt for confirmation
read -p "Are you sure you want to continue? (yes/no): " -r
echo
if [[ ! $REPLY =~ ^[Yy][Ee][Ss]$ ]]; then
    echo "Aborted."
    exit 0
fi

# Navigate to repository root
cd "${REPO_ROOT}"

echo "Step 1: Removing existing git directory..."
rm -rf .git

echo "Step 2: Initializing new git repository..."
git init

echo "Step 3: Adding all files to staging..."
git add .

echo "Step 4: Creating Release commit..."
git commit -m "Release"

echo ""
echo "======================================"
echo "Success!"
echo "======================================"
echo ""
echo "Repository has been prepared for partner release."
echo "Current status:"
git log --oneline
echo ""
echo "Branch information:"
git branch -v
echo ""
echo "Next steps:"
echo "  - Review the repository contents"
echo "  - Add a remote if needed: git remote add origin <URL>"
echo "  - Push to partner repository: git push -u origin main"
echo ""

