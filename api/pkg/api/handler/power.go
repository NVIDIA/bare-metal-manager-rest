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

package handler

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
	temporalEnums "go.temporal.io/api/enums/v1"
	tClient "go.temporal.io/sdk/client"
	tp "go.temporal.io/sdk/temporal"

	"github.com/nvidia/bare-metal-manager-rest/api/pkg/api/handler/util/common"
	"github.com/nvidia/bare-metal-manager-rest/api/pkg/api/model"
	cerr "github.com/nvidia/bare-metal-manager-rest/common/pkg/util"
	rlav1 "github.com/nvidia/bare-metal-manager-rest/workflow-schema/rla/protobuf/v1"
	"github.com/nvidia/bare-metal-manager-rest/workflow/pkg/queue"
)

// executePowerControlWorkflow determines the appropriate power control workflow based on state,
// executes it via Temporal, and returns the API response with task IDs.
func executePowerControlWorkflow(
	ctx context.Context,
	c echo.Context,
	logger zerolog.Logger,
	stc tClient.Client,
	targetSpec *rlav1.OperationTargetSpec,
	state string,
	workflowID string,
	entityName string,
) error {
	var workflowName string
	var rlaRequest interface{}

	switch state {
	case "on":
		workflowName = "PowerOnRack"
		rlaRequest = &rlav1.PowerOnRackRequest{
			TargetSpec:  targetSpec,
			Description: fmt.Sprintf("API power on %s", entityName),
		}
	case "off":
		workflowName = "PowerOffRack"
		rlaRequest = &rlav1.PowerOffRackRequest{
			TargetSpec:  targetSpec,
			Description: fmt.Sprintf("API power off %s", entityName),
		}
	case "cycle":
		workflowName = "PowerResetRack"
		rlaRequest = &rlav1.PowerResetRackRequest{
			TargetSpec:  targetSpec,
			Description: fmt.Sprintf("API power cycle %s", entityName),
		}
	case "forceoff":
		workflowName = "PowerOffRack"
		rlaRequest = &rlav1.PowerOffRackRequest{
			TargetSpec:  targetSpec,
			Forced:      true,
			Description: fmt.Sprintf("API force power off %s", entityName),
		}
	case "forcecycle":
		workflowName = "PowerResetRack"
		rlaRequest = &rlav1.PowerResetRackRequest{
			TargetSpec:  targetSpec,
			Forced:      true,
			Description: fmt.Sprintf("API force power cycle %s", entityName),
		}
	default:
		return cerr.NewAPIErrorResponse(c, http.StatusBadRequest, fmt.Sprintf("Invalid power control state: %s", state), nil)
	}

	workflowOptions := tClient.StartWorkflowOptions{
		ID:                       workflowID,
		WorkflowIDReusePolicy:    temporalEnums.WORKFLOW_ID_REUSE_POLICY_ALLOW_DUPLICATE,
		WorkflowExecutionTimeout: common.WorkflowExecutionTimeout,
		TaskQueue:                queue.SiteTaskQueue,
	}

	ctx, cancel := context.WithTimeout(ctx, common.WorkflowContextTimeout)
	defer cancel()

	we, err := stc.ExecuteWorkflow(ctx, workflowOptions, workflowName, rlaRequest)
	if err != nil {
		logger.Error().Err(err).Msg(fmt.Sprintf("failed to execute %s workflow", workflowName))
		return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, fmt.Sprintf("Failed to power control %s", entityName), nil)
	}

	var rlaResponse rlav1.SubmitTaskResponse
	err = we.Get(ctx, &rlaResponse)
	if err != nil {
		var timeoutErr *tp.TimeoutError
		if errors.As(err, &timeoutErr) || err == context.DeadlineExceeded || ctx.Err() != nil {
			return common.TerminateWorkflowOnTimeOut(c, logger, stc, workflowID, err, entityName, workflowName)
		}
		logger.Error().Err(err).Msg(fmt.Sprintf("failed to get result from %s workflow", workflowName))
		return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, fmt.Sprintf("Failed to power control %s", entityName), nil)
	}

	apiResponse := model.NewAPIPowerControlResponse(&rlaResponse)

	logger.Info().Str("state", state).Msg("finishing API handler")

	return c.JSON(http.StatusOK, apiResponse)
}
