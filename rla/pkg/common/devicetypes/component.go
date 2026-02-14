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
package devicetypes

import "strings"

// Define component types
type ComponentType int

const (
	ComponentTypeUnknown ComponentType = iota
	ComponentTypeCompute
	ComponentTypeNVLSwitch
	ComponentTypePowerShelf
	ComponentTypeToRSwitch
	ComponentTypeUMS
	ComponentTypeCDU
)

var (
	componentTypeStrings = map[ComponentType]string{
		ComponentTypeUnknown:    "Unknown",
		ComponentTypeCompute:    "Compute",
		ComponentTypeNVLSwitch:  "NVLSwitch",
		ComponentTypePowerShelf: "PowerShelf",
		ComponentTypeToRSwitch:  "ToRSwitch",
		ComponentTypeUMS:        "UMS",
		ComponentTypeCDU:        "CDU",
	}

	componentTypeStringMaxLen int
)

func init() {
	for _, str := range componentTypeStrings {
		if len(str) > componentTypeStringMaxLen {
			componentTypeStringMaxLen = len(str)
		}
	}
}

// ComponentTypes returns all the supported Component types
func ComponentTypes() []ComponentType {
	return []ComponentType{
		ComponentTypeUnknown,
		ComponentTypeCompute,
		ComponentTypeNVLSwitch,
		ComponentTypePowerShelf,
		ComponentTypeToRSwitch,
		ComponentTypeUMS,
		ComponentTypeCDU,
	}
}

// ComponentTypeFromString returns the Component type from the given string.
func ComponentTypeFromString(str string) ComponentType {
	for ct, componentTypeStr := range componentTypeStrings {
		if strings.EqualFold(str, componentTypeStr) {
			return ct
		}
	}
	return ComponentTypeUnknown
}

// ComponentTypeToString returns the string representation for the given
// component type.
func ComponentTypeToString(ct ComponentType) string {
	return componentTypeStrings[ct]
}

func IsValidComponentTypeString(str string) bool {
	return ComponentTypeFromString(str) != ComponentTypeUnknown
}

// String return the aligned string representation for the given component
// type
func (ct ComponentType) String() string {
	spaces := componentTypeStringMaxLen - len(ComponentTypeToString(ct))
	return strings.Repeat(" ", spaces) + ComponentTypeToString(ct)
}
