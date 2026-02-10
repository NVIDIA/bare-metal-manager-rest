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
package common

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/nvidia/carbide-rest/rla/pkg/common/deviceinfo"
	"github.com/nvidia/carbide-rest/rla/pkg/common/devicetypes"
)

type TaskType string

const (
	TaskTypeUnknown           TaskType = "unknown"
	TaskTypeInjectExpectation TaskType = "inject_expectation"
	TaskTypePowerControl      TaskType = "power_control"
	TaskTypeFirmwareControl   TaskType = "firmware_control"
)

func (tt TaskType) IsValid() bool {
	return tt != TaskTypeUnknown
}

type ExecutorType string

const (
	ExecutorTypeUnknown  ExecutorType = "unknown"
	ExecutorTypeTemporal ExecutorType = "temporal"
)

func (et ExecutorType) IsValid() bool {
	return et != ExecutorTypeUnknown
}

type TaskStatus string

const (
	TaskStatusUnknown    TaskStatus = "unknown"
	TaskStatusPending    TaskStatus = "pending"
	TaskStatusRunning    TaskStatus = "running"
	TaskStatusCompleted  TaskStatus = "completed"
	TaskStatusFailed     TaskStatus = "failed"
	TaskStatusTerminated TaskStatus = "terminated"
)

func (s TaskStatus) IsFinished() bool {
	return s == TaskStatusCompleted ||
		s == TaskStatusFailed ||
		s == TaskStatusTerminated
}

type TaskListOptions struct {
	TaskType   TaskType
	RackID     uuid.UUID
	ActiveOnly bool
}

type ComponentInfo struct {
	Type        devicetypes.ComponentType
	DeviceInfo  deviceinfo.DeviceInfo
	ComponentID string // Component service's own ID (e.g., machine_id from Carbide)
}

func (ci *ComponentInfo) Validate() error {
	if ci.Type == devicetypes.ComponentTypeUnknown {
		return fmt.Errorf("component type is unknown")
	}

	if !ci.DeviceInfo.VerifyIDOrSerial() {
		return fmt.Errorf("component device info is invalid")
	}

	return nil
}
