// SPDX-FileCopyrightText: Copyright (c) 2021-2023 NVIDIA CORPORATION & AFFILIATES. All rights reserved.
// SPDX-License-Identifier: LicenseRef-NvidiaProprietary
//
// NVIDIA CORPORATION, its affiliates and licensors retain all intellectual
// property and proprietary rights in and to this material, related
// documentation and any modifications thereto. Any use, reproduction,
// disclosure or distribution of this material and related documentation
// without an express license agreement from NVIDIA CORPORATION or
// its affiliates is strictly prohibited.

package workflow

import (
	"time"

	"github.com/rs/zerolog/log"

	"github.com/nvidia/carbide-rest/site-workflow/pkg/activity"
	rlav1 "github.com/nvidia/carbide-rest/workflow-schema/rla/protobuf/v1"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

// GetTray is a workflow to get a tray by its UUID from RLA
func GetTray(ctx workflow.Context, request *rlav1.GetComponentInfoByIDRequest) (*rlav1.GetComponentInfoResponse, error) {
	logger := log.With().Str("Workflow", "Tray").Str("Action", "Get").Logger()
	if request != nil && request.Id != nil {
		logger = log.With().Str("Workflow", "Tray").Str("Action", "Get").Str("TrayID", request.Id.Id).Logger()
	}

	logger.Info().Msg("starting workflow")

	// RetryPolicy specifies how to automatically handle retries if an Activity fails.
	retrypolicy := &temporal.RetryPolicy{
		InitialInterval:    1 * time.Second,
		BackoffCoefficient: 2.0,
		MaximumInterval:    10 * time.Second,
		MaximumAttempts:    2,
	}
	options := workflow.ActivityOptions{
		// Timeout options specify when to automatically timeout Activity functions.
		StartToCloseTimeout: 2 * time.Minute,
		// Optionally provide a customized RetryPolicy.
		RetryPolicy: retrypolicy,
	}

	ctx = workflow.WithActivityOptions(ctx, options)

	var trayManager activity.ManageTray
	var response rlav1.GetComponentInfoResponse

	err := workflow.ExecuteActivity(ctx, trayManager.GetTray, request).Get(ctx, &response)
	if err != nil {
		logger.Error().Err(err).Str("Activity", "GetTray").Msg("Failed to execute activity from workflow")
		return nil, err
	}

	logger.Info().Msg("completing workflow")

	return &response, nil
}

// GetTrays is a workflow to get a list of trays from RLA with optional filters.
// If taskID is non-empty, it first resolves the task to its component UUIDs via the
// GetTaskComponentIDs activity and uses those to filter the returned trays.
func GetTrays(ctx workflow.Context, request *rlav1.GetComponentsRequest, taskID string) (*rlav1.GetComponentsResponse, error) {
	logger := log.With().Str("Workflow", "Tray").Str("Action", "GetAll").Logger()

	logger.Info().Msg("starting workflow")

	// RetryPolicy specifies how to automatically handle retries if an Activity fails.
	retrypolicy := &temporal.RetryPolicy{
		InitialInterval:    1 * time.Second,
		BackoffCoefficient: 2.0,
		MaximumInterval:    10 * time.Second,
		MaximumAttempts:    2,
	}
	options := workflow.ActivityOptions{
		// Timeout options specify when to automatically timeout Activity functions.
		StartToCloseTimeout: 2 * time.Minute,
		// Optionally provide a customized RetryPolicy.
		RetryPolicy: retrypolicy,
	}

	ctx = workflow.WithActivityOptions(ctx, options)

	var trayManager activity.ManageTray

	// If a taskID was provided, resolve it to component UUIDs first
	if taskID != "" {
		logger.Info().Str("TaskID", taskID).Msg("resolving task to component UUIDs")

		var componentIDs []string
		err := workflow.ExecuteActivity(ctx, trayManager.GetTaskComponentIDs, taskID).Get(ctx, &componentIDs)
		if err != nil {
			logger.Error().Err(err).Str("Activity", "GetTaskComponentIDs").Msg("Failed to resolve task component IDs")
			return nil, err
		}

		// Task has no components — return empty response immediately
		if len(componentIDs) == 0 {
			logger.Info().Str("TaskID", taskID).Msg("task has no component UUIDs, returning empty response")
			return &rlav1.GetComponentsResponse{}, nil
		}

		// Build component targets from the resolved IDs and set them on the request
		componentTargets := make([]*rlav1.ComponentTarget, 0, len(componentIDs))
		for _, id := range componentIDs {
			componentTargets = append(componentTargets, &rlav1.ComponentTarget{
				Identifier: &rlav1.ComponentTarget_Id{
					Id: &rlav1.UUID{Id: id},
				},
			})
		}

		if request == nil {
			request = &rlav1.GetComponentsRequest{}
		}
		request.TargetSpec = &rlav1.OperationTargetSpec{
			Targets: &rlav1.OperationTargetSpec_Components{
				Components: &rlav1.ComponentTargets{
					Targets: componentTargets,
				},
			},
		}
	}

	var response rlav1.GetComponentsResponse

	err := workflow.ExecuteActivity(ctx, trayManager.GetTrays, request).Get(ctx, &response)
	if err != nil {
		logger.Error().Err(err).Str("Activity", "GetTrays").Msg("Failed to execute activity from workflow")
		return nil, err
	}

	logger.Info().Int32("total", response.GetTotal()).Msg("completing workflow")

	return &response, nil
}
