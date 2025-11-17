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

echo "Elektra container starting up..."
echo "executing elektra process"

# Execute elektra process
if [ "$1" == 'elektra' ]; then
    echo "Listing the executable directory"
    ls /usr/local/bin
    echo "Listing the executable binary"
    ls /usr/local/bin/elektra
    # Add future prequisites here
    nohup sh /usr/local/bin/elektraserver &
    exec /usr/local/bin/elektra
fi

# Run the command as is
exec "$@"
