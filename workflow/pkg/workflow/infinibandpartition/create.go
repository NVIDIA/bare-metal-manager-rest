// SPDX-FileCopyrightText: Copyright (c) 2021-2023 NVIDIA CORPORATION & AFFILIATES. All rights reserved.
// SPDX-License-Identifier: LicenseRef-NvidiaProprietary
//
// NVIDIA CORPORATION, its affiliates and licensors retain all intellectual
// property and proprietary rights in and to this material, related
// documentation and any modifications thereto. Any use, reproduction,
// disclosure or distribution of this material and related documentation
// without an express license agreement from NVIDIA CORPORATION or
// its affiliates is strictly prohibited.

package infinibandpartition

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"

	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"

	"go.temporal.io/sdk/client"

	ibpActivity "github.com/nvidia/carbide-rest/workflow/pkg/activity/infinibandpartition"
	"github.com/nvidia/carbide-rest/workflow/pkg/queue"
)

// CreateInfiniBandPartition is a Temporal workflow to create a new InfiniBandPartition via Site Agent
func CreateInfiniBandPartition(ctx workflow.Context, siteID uuid.UUID, ibpID uuid.UUID) error {
	logger := log.With().Str("Workflow", "CreateInfiniBandPartition").Str("Site ID", siteID.String()).
		Str("Partition ID", ibpID.String()).Logger()

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

	var ibpManager ibpActivity.ManageInfiniBandPartition

	err := workflow.ExecuteActivity(ctx, ibpManager.CreateInfiniBandPartitionViaSiteAgent, siteID, ibpID).Get(ctx, nil)
	if err != nil {
		logger.Warn().Err(err).Msg("failed to execute activity: CreateInfiniBandPartitionViaSiteAgent")
		return err
	}

	logger.Info().Msg("completing workflow")

	return nil
}

// ExecuteCreateInfiniBandPartitionWorkflow is a helper function to trigger execution of create InfiniBandPartition workflow
func ExecuteCreateInfiniBandPartitionWorkflow(ctx context.Context, tc client.Client, siteID uuid.UUID, ibpID uuid.UUID) (*string, error) {
	uid := uuid.New()

	workflowOptions := client.StartWorkflowOptions{
		ID:        "infiniband-partition-create-" + uid.String(),
		TaskQueue: queue.CloudTaskQueue,
	}

	we, err := tc.ExecuteWorkflow(ctx, workflowOptions, CreateInfiniBandPartition, siteID, ibpID)

	if err != nil {
		log.Error().Err(err).Msg("failed to execute workflow: CreateInfiniBandPartition")
		return nil, err
	}

	wid := we.GetID()

	return &wid, nil
}
