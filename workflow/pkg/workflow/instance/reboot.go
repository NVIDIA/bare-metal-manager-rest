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

// RebootInstance is a Temporal workflow to reboot a machine associated with Instance via Site Agent
func RebootInstance(ctx workflow.Context, instanceID uuid.UUID, rebootWithCustomIpxe bool, applyUpdatesOnReboot bool) error {
	logger := log.With().Str("Workflow", "Instance").Str("Action", "Reboot").Str("Instance ID", instanceID.String()).Logger()

	logger.Info().Msg("starting workflow")

	// RetryPolicy specifies how to automatically handle retries if an Activity fails.
	retrypolicy := &temporal.RetryPolicy{
		InitialInterval:    2 * time.Second,
		BackoffCoefficient: 2.0,
		MaximumInterval:    2 * time.Minute,
		MaximumAttempts:    15,
	}
	options := workflow.ActivityOptions{
		// Timeout options specify when to automatically timeout Activity functions.
		StartToCloseTimeout: 2 * time.Minute,
		// Optionally provide a customized RetryPolicy.
		RetryPolicy: retrypolicy,
	}

	ctx = workflow.WithActivityOptions(ctx, options)

	var instanceManager instanceActivity.ManageInstance

	err := workflow.ExecuteActivity(ctx, instanceManager.RebootInstanceViaSiteAgent, instanceID, rebootWithCustomIpxe, applyUpdatesOnReboot).Get(ctx, nil)
	if err != nil {
		logger.Warn().Err(err).Msg("failed to execute activity: RebootInstanceViaSiteAgent")
		return err
	}

	logger.Info().Msg("completing workflow")

	return nil
}

// ExecuteRebootInstanceWorkflow is a helper function to trigger execution of reboot Instance workflow
func ExecuteRebootInstanceWorkflow(ctx context.Context, tc client.Client, instanceID uuid.UUID, rebootWithCustomIpxe bool, applyUpdatesOnReboot bool) (*string, error) {
	uid := uuid.New()

	workflowOptions := client.StartWorkflowOptions{
		ID:        "instance-reboot-" + uid.String(),
		TaskQueue: queue.CloudTaskQueue,
	}

	we, err := tc.ExecuteWorkflow(ctx, workflowOptions, RebootInstance, instanceID, rebootWithCustomIpxe, applyUpdatesOnReboot)

	if err != nil {
		log.Error().Err(err).Msg("failed to execute workflow: RebootInstance")
		return nil, err
	}

	wid := we.GetID()

	return &wid, nil
}
