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

import (
	"encoding/json"
	"fmt"
	"time"

	taskcommon "github.com/nvidia/bare-metal-manager-rest/rla/internal/task/common"
)

type Operation interface {
	Validate() error
	Marshal() (json.RawMessage, error)
	Unmarshal(data json.RawMessage) error
	Type() taskcommon.TaskType
	Description() string
}

func New(typ taskcommon.TaskType, info json.RawMessage) (Operation, error) {
	switch typ {
	case taskcommon.TaskTypePowerControl:
		var taskInfo PowerControlTaskInfo
		if err := json.Unmarshal(info, &taskInfo); err != nil {
			return nil, fmt.Errorf("failed to unmarshal power control task info: %w", err) //nolint
		}
		return &taskInfo, nil
	case taskcommon.TaskTypeFirmwareControl:
		var taskInfo FirmwareControlTaskInfo
		if err := json.Unmarshal(info, &taskInfo); err != nil {
			return nil, fmt.Errorf("failed to unmarshal firmware control task info: %w", err) //nolint
		}
		return &taskInfo, nil
	case taskcommon.TaskTypeInjectExpectation:
		var taskInfo InjectExpectationTaskInfo
		if err := json.Unmarshal(info, &taskInfo); err != nil {
			return nil, fmt.Errorf("failed to unmarshal inject expectation task info: %w", err) //nolint
		}
		return &taskInfo, nil
	default:
		return nil, fmt.Errorf("unsupported task type: %s", typ)
	}
}

type PowerControlTaskInfo struct {
	Operation PowerOperation `json:"operation"`
	Forced    bool           `json:"forced"`
}

func (t *PowerControlTaskInfo) Validate() error {
	if t.Operation == PowerOperationUnknown {
		return fmt.Errorf("invalid power control operation")
	}

	return nil
}

func (t *PowerControlTaskInfo) Marshal() (json.RawMessage, error) {
	raw, err := json.Marshal(t)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal power control task info: %w", err)
	}
	return raw, nil
}

func (t *PowerControlTaskInfo) Unmarshal(data json.RawMessage) error {
	if err := json.Unmarshal(data, t); err != nil {
		return fmt.Errorf("failed to unmarshal power control task info: %w", err)
	}
	return nil
}

func (t *PowerControlTaskInfo) Type() taskcommon.TaskType {
	return taskcommon.TaskTypePowerControl
}

func (t *PowerControlTaskInfo) Description() string {
	return fmt.Sprintf("%s, forced %t", t.Operation.String(), t.Forced)
}

type InjectExpectationTaskInfo struct {
	Info json.RawMessage `json:"info"`
}

func (t *InjectExpectationTaskInfo) Validate() error {
	if t.Info == nil {
		return fmt.Errorf("invalid info")
	}

	return nil
}

func (t *InjectExpectationTaskInfo) Marshal() (json.RawMessage, error) {
	raw, err := json.Marshal(t)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal inject expectation task info: %w", err)
	}
	return raw, nil
}

func (t *InjectExpectationTaskInfo) Unmarshal(data json.RawMessage) error {
	if err := json.Unmarshal(data, t); err != nil {
		return fmt.Errorf("failed to unmarshal inject expectation task info: %w", err)
	}
	return nil
}

func (t *InjectExpectationTaskInfo) Type() taskcommon.TaskType {
	return taskcommon.TaskTypeInjectExpectation
}

func (t *InjectExpectationTaskInfo) Description() string {
	return fmt.Sprintf("inject expectation: %s", t.Info)
}

type FirmwareControlTaskInfo struct {
	Operation     FirmwareOperation `json:"operation"`
	TargetVersion string            `json:"target_version,omitempty"`
	StartTime     int64             `json:"start_time,omitempty"` // Unix timestamp
	EndTime       int64             `json:"end_time,omitempty"`   // Unix timestamp
}

func (t *FirmwareControlTaskInfo) Validate() error {
	if t.Operation == FirmwareOperationUnknown {
		return fmt.Errorf("invalid firmware control operation")
	}

	return nil
}

func (t *FirmwareControlTaskInfo) Marshal() (json.RawMessage, error) {
	raw, err := json.Marshal(t)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal firmware control task info: %w", err)
	}
	return raw, nil
}

func (t *FirmwareControlTaskInfo) Unmarshal(data json.RawMessage) error {
	if err := json.Unmarshal(data, t); err != nil {
		return fmt.Errorf("failed to unmarshal firmware control task info: %w", err)
	}
	return nil
}

func (t *FirmwareControlTaskInfo) Type() taskcommon.TaskType {
	return taskcommon.TaskTypeFirmwareControl
}

func (t *FirmwareControlTaskInfo) Description() string {
	return fmt.Sprintf("%s, target version %s", t.Operation.String(), t.TargetVersion)
}

// SetFirmwareUpdateTimeWindowRequest is the request for setting firmware update time window.
// ComponentIDs are external IDs (e.g., machine_id in Carbide) that identify the components.
type SetFirmwareUpdateTimeWindowRequest struct {
	ComponentIDs []string
	StartTime    time.Time
	EndTime      time.Time
}
