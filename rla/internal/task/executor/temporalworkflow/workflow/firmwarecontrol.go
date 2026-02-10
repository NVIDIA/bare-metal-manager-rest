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
	"fmt"
	"time"

	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"

	"github.com/nvidia/carbide-rest/rla/internal/task/operations"
	"github.com/nvidia/carbide-rest/rla/internal/task/task"
)

var firmwareControlActivityOptions = workflow.ActivityOptions{
	StartToCloseTimeout: 30 * time.Minute,
	RetryPolicy: &temporal.RetryPolicy{
		MaximumAttempts:    3,
		InitialInterval:    5 * time.Second,
		MaximumInterval:    2 * time.Minute,
		BackoffCoefficient: 2,
	},
}

func FirmwareControl(
	ctx workflow.Context,
	reqInfo task.ExecutionInfo,
	info *operations.FirmwareControlTaskInfo,
) error {
	if reqInfo.Rack == nil || len(reqInfo.Rack.Components) == 0 {
		return fmt.Errorf("no components provided")
	}

	if err := info.Validate(); err != nil {
		return fmt.Errorf("invalid firmware control info: %w", err)
	}

	ctx = workflow.WithActivityOptions(ctx, firmwareControlActivityOptions)

	if err := updateRunningTaskStatus(ctx, reqInfo.TaskID); err != nil {
		return err
	}

	// Collect all component IDs (external IDs like machine_id in Carbide)
	var componentIDs []string
	for _, comp := range reqInfo.Rack.Components {
		if comp.ComponentID != "" {
			componentIDs = append(componentIDs, comp.ComponentID)
		}
	}

	if len(componentIDs) == 0 {
		return fmt.Errorf("no component IDs found in components")
	}

	// Single batch call to set firmware update time window
	req := operations.SetFirmwareUpdateTimeWindowRequest{
		ComponentIDs: componentIDs,
		StartTime:    time.Unix(info.StartTime, 0),
		EndTime:      time.Unix(info.EndTime, 0),
	}

	err := workflow.ExecuteActivity(ctx, "SetFirmwareUpdateTimeWindow", req).Get(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to set firmware update time window: %w", err)
	}

	return nil
}
