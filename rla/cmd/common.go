/*
 * SPDX-FileCopyrightText: Copyright (c) 2026 NVIDIA CORPORATION & AFFILIATES. All rights reserved.
 * SPDX-License-Identifier: Apache-2.0
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
package cmd

import (
	"strings"

	pb "github.com/nvidia/bare-metal-manager-rest/rla/internal/proto/v1"
)

// parseComponentType parses a user-friendly string into a protobuf ComponentType.
// Accepts: "compute", "nvlswitch", "nvl-switch", "powershelf", "power-shelf", etc.
func parseComponentType(s string) pb.ComponentType {
	switch strings.ToLower(s) {
	case "compute":
		return pb.ComponentType_COMPONENT_TYPE_COMPUTE
	case "nvlswitch", "nvl-switch":
		return pb.ComponentType_COMPONENT_TYPE_NVLSWITCH
	case "powershelf", "power-shelf":
		return pb.ComponentType_COMPONENT_TYPE_POWERSHELF
	case "torswitch", "tor-switch":
		return pb.ComponentType_COMPONENT_TYPE_TORSWITCH
	case "ums":
		return pb.ComponentType_COMPONENT_TYPE_UMS
	case "cdu":
		return pb.ComponentType_COMPONENT_TYPE_CDU
	default:
		return pb.ComponentType_COMPONENT_TYPE_UNKNOWN
	}
}
