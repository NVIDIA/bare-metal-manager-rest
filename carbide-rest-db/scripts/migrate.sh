#!/bin/sh
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

# Exit on error
set -e

# Get the current directory
SCRIPT_DIR=$(dirname "$0")

# Setup migration tables
echo "Creating migration tables..."
$SCRIPT_DIR/migrations db init

# Run migrations
echo "Running migrations..."
$SCRIPT_DIR/migrations db migrate
