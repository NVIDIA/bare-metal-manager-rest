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

	"github.com/rs/zerolog/log"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"

	"github.com/nvidia/carbide-rest/rla/internal/task/executor/temporalworkflow/common"
	"github.com/nvidia/carbide-rest/rla/internal/task/operations"
	"github.com/nvidia/carbide-rest/rla/internal/task/task"
	"github.com/nvidia/carbide-rest/rla/pkg/common/devicetypes"
)

var (
	powerOnSequence = []devicetypes.ComponentType{
		devicetypes.ComponentTypePowerShelf,
		devicetypes.ComponentTypeNVLSwitch,
		devicetypes.ComponentTypeCompute,
	}

	powerOffSequence = []devicetypes.ComponentType{
		devicetypes.ComponentTypeCompute,
		devicetypes.ComponentTypeNVLSwitch,
		devicetypes.ComponentTypePowerShelf,
	}

	powerControlActivityOptions = workflow.ActivityOptions{
		StartToCloseTimeout: 20 * time.Minute,
		RetryPolicy: &temporal.RetryPolicy{
			MaximumAttempts:    3,
			InitialInterval:    1 * time.Second,
			MaximumInterval:    1 * time.Minute,
			BackoffCoefficient: 2,
		},
	}
)

func PowerControl(
	ctx workflow.Context,
	reqInfo task.ExecutionInfo,
	info operations.PowerControlTaskInfo,
) (err error) {
	if reqInfo.Rack == nil || len(reqInfo.Rack.Components) == 0 {
		return nil
	}

	ctx = workflow.WithActivityOptions(ctx, powerControlActivityOptions)

	if err := updateRunningTaskStatus(ctx, reqInfo.TaskID); err != nil {
		// XXX -- The workflow will be terminated, but the task status won't be
		// updated. We need to add a background process to try to detect and
		// fix this situation in the future.
		return err
	}

	typeToTargets := buildTargets(&reqInfo)

	switch info.Operation {
	case operations.PowerOperationPowerOn, operations.PowerOperationForcePowerOn:
		for _, t := range powerOnSequence {
			if err = powerControlTarget(ctx, typeToTargets, t, info); err != nil {
				break
			}
		}
	case operations.PowerOperationPowerOff, operations.PowerOperationForcePowerOff:
		for _, t := range powerOffSequence {
			if err = powerControlTarget(ctx, typeToTargets, t, info); err != nil {
				break
			}
		}
	case operations.PowerOperationRestart, operations.PowerOperationForceRestart:
		// Restart: power off then power on
		for _, t := range powerOffSequence {
			if err = powerControlTarget(ctx, typeToTargets, t, info); err != nil {
				break
			}
		}
		for _, t := range powerOnSequence {
			if err = powerControlTarget(ctx, typeToTargets, t, info); err != nil {
				break
			}
		}
	case operations.PowerOperationWarmReset, operations.PowerOperationColdReset:
		err = fmt.Errorf("hardware reset is not supported yet")
	default:
		err = fmt.Errorf("unknown power operation: %v", info.Operation)
	}

	return updateFinishedTaskStatus(ctx, reqInfo.TaskID, err)
}

func powerControlTarget(
	ctx workflow.Context,
	typeToTarget map[devicetypes.ComponentType]common.Target,
	targetType devicetypes.ComponentType,
	info operations.PowerControlTaskInfo,
) error {
	target, ok := typeToTarget[targetType]
	if !ok {
		return nil
	}

	log.Debug().Msgf(
		"power control target %s op %s workflow started",
		target.String(), info.Operation.String(),
	)

	err := workflow.ExecuteActivity(ctx, "PowerControl", target, info).Get(ctx, nil) //nolint
	if err != nil {
		log.Debug().Msgf(
			"power control target %s op %s workflow completed with error: %v",
			target.String(), info.Operation.String(), err,
		)

		return err
	}

	log.Debug().Msgf(
		"power control target %s op %s workflow completed successfully",
		target.String(), info.Operation.String(),
	)

	return nil
}
