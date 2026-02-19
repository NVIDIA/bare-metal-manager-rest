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
package operations

type PowerOperation int

const (
	PowerOperationUnknown PowerOperation = iota
	// Power On
	PowerOperationPowerOn
	PowerOperationForcePowerOn
	// Power Off
	PowerOperationPowerOff
	PowerOperationForcePowerOff
	// Restart (OS level)
	PowerOperationRestart
	PowerOperationForceRestart
	// Reset (hardware level)
	PowerOperationWarmReset
	PowerOperationColdReset
)

func (o PowerOperation) String() string {
	switch o {
	case PowerOperationPowerOn:
		return "PowerOn"
	case PowerOperationForcePowerOn:
		return "ForcePowerOn"
	case PowerOperationPowerOff:
		return "PowerOff"
	case PowerOperationForcePowerOff:
		return "ForcePowerOff"
	case PowerOperationRestart:
		return "Restart"
	case PowerOperationForceRestart:
		return "ForceRestart"
	case PowerOperationWarmReset:
		return "WarmReset"
	case PowerOperationColdReset:
		return "ColdReset"
	}

	return "Unknown"
}

type PowerStatus string

const (
	PowerStatusUnknown   PowerStatus = "Unknown"
	PowerStatusOn        PowerStatus = "On"
	PowerStatusOff       PowerStatus = "Off"
	PowerStatusRebooting PowerStatus = "Rebooting"
)

type FirmwareOperation int

const (
	FirmwareOperationUnknown FirmwareOperation = iota
	FirmwareOperationUpgrade
	FirmwareOperationDowngrade
	FirmwareOperationRollback
	FirmwareOperationVersion
)

func (o FirmwareOperation) String() string {
	switch o {
	case FirmwareOperationUpgrade:
		return "FirmwareUpgrade"
	case FirmwareOperationDowngrade:
		return "FirmwareDowngrade"
	case FirmwareOperationRollback:
		return "FirmwareRollback"
	case FirmwareOperationVersion:
		return "FirmwareVersion"
	}

	return "Unknown"
}
