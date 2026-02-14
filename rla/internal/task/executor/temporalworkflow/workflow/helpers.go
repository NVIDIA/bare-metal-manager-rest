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
package workflow

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
	"go.temporal.io/sdk/workflow"

	taskcommon "github.com/nvidia/bare-metal-manager-rest/rla/internal/task/common"
	"github.com/nvidia/bare-metal-manager-rest/rla/internal/task/executor/temporalworkflow/common"
	"github.com/nvidia/bare-metal-manager-rest/rla/internal/task/task"
	"github.com/nvidia/bare-metal-manager-rest/rla/pkg/common/devicetypes"
)

func updateRunningTaskStatus(
	ctx workflow.Context,
	taskID uuid.UUID,
) error {
	if taskID == uuid.Nil {
		return fmt.Errorf("task ID is not specified")
	}

	arg := &task.TaskStatusUpdate{
		ID:      taskID,
		Status:  taskcommon.TaskStatusRunning,
		Message: "Running",
	}

	return workflow.ExecuteActivity(ctx, "UpdateTaskStatus", arg).Get(ctx, nil)
}

func updateFinishedTaskStatus(
	ctx workflow.Context,
	taskID uuid.UUID,
	err error,
) error {
	if taskID == uuid.Nil {
		return fmt.Errorf("task ID is not specified")
	}

	var arg *task.TaskStatusUpdate

	if err != nil {
		arg = &task.TaskStatusUpdate{
			ID:      taskID,
			Status:  taskcommon.TaskStatusFailed,
			Message: err.Error(),
		}
	} else {
		arg = &task.TaskStatusUpdate{
			ID:      taskID,
			Status:  taskcommon.TaskStatusCompleted,
			Message: "Completed successfully",
		}
	}

	if lerr := workflow.ExecuteActivity(ctx, "UpdateTaskStatus", arg).Get(ctx, nil); lerr != nil { //nolint
		return errors.Join(err, fmt.Errorf("failed to update task status: %v", lerr))
	}

	return err
}

func buildTargets(info *task.ExecutionInfo) map[devicetypes.ComponentType]common.Target {
	if info.Rack == nil {
		return nil
	}

	// Group component IDs by type
	mapOnType := make(map[devicetypes.ComponentType][]string)
	for _, c := range info.Rack.Components {
		if c.ComponentID != "" {
			mapOnType[c.Type] = append(mapOnType[c.Type], c.ComponentID)
		}
	}

	// Build Target for each type with component IDs only
	results := make(map[devicetypes.ComponentType]common.Target)
	for t, componentIDs := range mapOnType {
		results[t] = common.Target{
			Type:         t,
			ComponentIDs: componentIDs,
		}
	}

	return results
}
