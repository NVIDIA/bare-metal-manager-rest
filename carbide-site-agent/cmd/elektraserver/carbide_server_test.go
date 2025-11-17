/*
 * SPDX-FileCopyrightText: Copyright (c) 2021-2023 NVIDIA CORPORATION & AFFILIATES. All rights reserved.
 * SPDX-License-Identifier: LicenseRef-NvidiaProprietary
 *
 * NVIDIA CORPORATION, its affiliates and licensors retain all intellectual
 * property and proprietary rights in and to this material, related
 * documentation and any modifications thereto. Any use, reproduction,
 * disclosure or distribution of this material and related documentation
 * without an express license agreement from NVIDIA CORPORATION or
 * its affiliates is strictly prohibited.
 */

package main

import (
	"os/exec"
	"testing"
)

// tstCarbideServer - test the carbide grpc client
func tstCarbideServer(t *testing.T) {
	err := exec.Command("../bin/elektraserver").Run()
	if err != nil {
		t.Log(err.Error())
	}
}
