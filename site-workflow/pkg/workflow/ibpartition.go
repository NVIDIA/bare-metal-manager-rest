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
	cwssaws "github.com/nvidia/carbide-rest/workflow-schema/schema/site-agent/workflows/v1"
	"github.com/nvidia/carbide-rest/site-workflow/pkg/activity"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

func DiscoverInfiniBandPartitionInventory(ctx workflow.Context) error {
	logger := log.With().Str("Workflow", "DiscoverInfiniBandPartitionInventory").Logger()

	logger.Info().Msg("Starting workflow")

	// RetryPolicy specifies how to automatically handle retries if an Activity fails.
	retrypolicy := &temporal.RetryPolicy{
		InitialInterval:    2 * time.Second,
		BackoffCoefficient: 2.0,
		MaximumInterval:    10 * time.Second,
		// This is executed every 3 minutes, so we don't want too many retry attempts
		MaximumAttempts: 2,
	}
	options := workflow.ActivityOptions{
		// Timeout options specify when to automatically timeout Activity functions.
		StartToCloseTimeout: 2 * time.Minute,
		// Optionally provide a customized RetryPolicy.
		RetryPolicy: retrypolicy,
	}

	ctx = workflow.WithActivityOptions(ctx, options)

	// Invoke activity
	var inventoryManager activity.ManageInfiniBandPartitionInventory

	err := workflow.ExecuteActivity(ctx, inventoryManager.DiscoverInfiniBandPartitionInventory).Get(ctx, nil)
	if err != nil {
		logger.Error().Err(err).Str("Activity", "DiscoverInfiniBandPartitionInventory").Msg("Failed to execute activity from workflow")
		return err
	}

	logger.Info().Msg("Completing workflow")

	return nil
}

// CreateInfiniBandPartitionV2 is a workflow to create new InfiniBand Partitions using the CreateInfiniBandPartitionOnSite activity
// V1 (CreateInfiniBandPartition) is found in cloud-workflow and uses a different activity that does not speak
// to carbide directly.
func CreateInfiniBandPartitionV2(ctx workflow.Context, request *cwssaws.IBPartitionCreationRequest) error {
	logger := log.With().Str("Workflow", "InfiniBandPartition").Str("Action", "Create").Str("IB Partition ID", request.GetId().GetValue()).Str("Name", request.GetConfig().GetName()).Logger()

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

	var ibPartitionManager activity.ManageInfiniBandPartition

	err := workflow.ExecuteActivity(ctx, ibPartitionManager.CreateInfiniBandPartitionOnSite, request).Get(ctx, nil)
	if err != nil {
		logger.Error().Err(err).Str("Activity", "CreateInfiniBandPartitionOnSite").Msg("Failed to execute activity from workflow")
		return err
	}

	logger.Info().Msg("completing workflow")

	return nil
}

// DeleteInfiniBandPartitionV2 is a workflow to Delete InfiniBand Partitions using the DeleteInfiniBandPartitionOnSite activity
// V1 (DeleteInfiniBandPartition) is found in cloud-workflow and uses a different activity that does not speak
// to carbide directly.
func DeleteInfiniBandPartitionV2(ctx workflow.Context, request *cwssaws.IBPartitionDeletionRequest) error {
	logger := log.With().Str("Workflow", "InfiniBandPartition").Str("Action", "Delete").Str("IB Partition ID", request.GetId().GetValue()).Logger()

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

	var ibPartitionManager activity.ManageInfiniBandPartition

	err := workflow.ExecuteActivity(ctx, ibPartitionManager.DeleteInfiniBandPartitionOnSite, request).Get(ctx, nil)
	if err != nil {
		logger.Error().Err(err).Str("Activity", "DeleteInfiniBandPartitionOnSite").Msg("Failed to execute activity from workflow")
		return err
	}

	logger.Info().Msg("completing workflow")

	return nil
}
