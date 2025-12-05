// SPDX-FileCopyrightText: Copyright (c) 2021-2023 NVIDIA CORPORATION & AFFILIATES. All rights reserved.
// SPDX-License-Identifier: LicenseRef-NvidiaProprietary
//
// NVIDIA CORPORATION, its affiliates and licensors retain all intellectual
// property and proprietary rights in and to this material, related
// documentation and any modifications thereto. Any use, reproduction,
// disclosure or distribution of this material and related documentation
// without an express license agreement from NVIDIA CORPORATION or
// its affiliates is strictly prohibited.

package instance

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"

	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"

	"go.temporal.io/sdk/client"

	instanceActivity "github.com/nvidia/carbide-rest/workflow/pkg/activity/instance"
	"github.com/nvidia/carbide-rest/workflow/pkg/queue"
)

// DeleteInstance is a Temporal workflow to delete an existing Instance via Site Agent
func DeleteInstance(ctx workflow.Context, instanceID uuid.UUID) error {
	logger := log.With().Str("Workflow", "Instance").Str("Action", "Delete").Str("Instance ID", instanceID.String()).Logger()

	logger.Info().Msg("starting workflow")

	// RetryPolicy specifies how to automatically handle retries if an Activity fails.
	retrypolicy := &temporal.RetryPolicy{
		InitialInterval:    2 * time.Second,
		BackoffCoefficient: 2.0,
		MaximumInterval:    2 * time.Minute,
		MaximumAttempts:    10,
	}
	options := workflow.ActivityOptions{
		// Timeout options specify when to automatically timeout Activity functions.
		StartToCloseTimeout: 3 * time.Minute,
		// Optionally provide a customized RetryPolicy.
		RetryPolicy: retrypolicy,
	}

	ctx = workflow.WithActivityOptions(ctx, options)

	var instanceManager instanceActivity.ManageInstance

	err := workflow.ExecuteActivity(ctx, instanceManager.DeleteInstanceViaSiteAgent, instanceID).Get(ctx, nil)
	if err != nil {
		logger.Warn().Err(err).Msg("failed to execute activity: DeleteInstanceViaSiteAgent")
		return err
	}

	logger.Info().Msg("completing workflow")

	return nil
}

// ExecuteDeleteInstanceWorkflow is a helper function to trigger execution of delete Instance workflow
func ExecuteDeleteInstanceWorkflow(ctx context.Context, tc client.Client, instanceID uuid.UUID) (*string, error) {
	uid := uuid.New()

	workflowOptions := client.StartWorkflowOptions{
		ID:        "instance-delete-" + uid.String(),
		TaskQueue: queue.CloudTaskQueue,
	}

	we, err := tc.ExecuteWorkflow(ctx, workflowOptions, DeleteInstance, instanceID)

	if err != nil {
		log.Error().Err(err).Msg("failed to execute workflow: DeleteInstance")
		return nil, err
	}

	wid := we.GetID()

	return &wid, nil
}
