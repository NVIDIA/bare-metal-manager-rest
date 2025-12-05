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
	"fmt"
	"time"

	cwm "github.com/nvidia/carbide-rest/workflow/internal/metrics"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"

	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"

	ibpActivity "github.com/nvidia/carbide-rest/workflow/pkg/activity/infinibandpartition"

	cwssaws "github.com/nvidia/carbide-rest/workflow-schema/schema/site-agent/workflows/v1"
)

// UpdateInfiniBandPartitionInfo is a Temporal workflow that Site Agent calls to update InfiniBandPartition information
func UpdateInfiniBandPartitionInfo(ctx workflow.Context, siteID string, transactionID *cwssaws.TransactionID, ibpInfo *cwssaws.InfiniBandPartitionInfo) error {
	logger := log.With().Str("Workflow", "UpdateInfiniBandPartitionInfo").Str("Site ID", siteID).Logger()

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

	err := workflow.ExecuteActivity(ctx, ibpManager.UpdateInfiniBandPartitionInDB, transactionID, ibpInfo).Get(ctx, nil)
	if err != nil {
		logger.Warn().Err(err).Msg("failed to execute activity: UpdateInfiniBandPartitionInDB")
		return err
	}

	logger.Info().Msg("completing workflow")

	return nil
}

// UpdateInfiniBandPartitionInventory is a workflow called by Site Agent to update InfiniBandPartition inventory for a Site
func UpdateInfiniBandPartitionInventory(ctx workflow.Context, siteID string, ibpInventory *cwssaws.InfiniBandPartitionInventory) (err error) {
	logger := log.With().Str("Workflow", "UpdateInfiniBandPartitionInventory").Str("Site ID", siteID).Logger()

	startTime := time.Now()

	logger.Info().Msg("starting workflow")

	parsedSiteID, err := uuid.Parse(siteID)
	if err != nil {
		logger.Warn().Err(err).Msg(fmt.Sprintf("workflow triggered with invalid site ID: %s", siteID))
		return err
	}

	// RetryPolicy specifies how to automatically handle retries if an Activity fails.
	retrypolicy := &temporal.RetryPolicy{
		InitialInterval:    5 * time.Second,
		BackoffCoefficient: 2.0,
		MaximumInterval:    30 * time.Second,
		MaximumAttempts:    2,
	}
	options := workflow.ActivityOptions{
		// Timeout options specify when to automatically timeout Activity functions.
		StartToCloseTimeout: 30 * time.Second,
		// Optionally provide a customized RetryPolicy.
		RetryPolicy: retrypolicy,
	}

	ctx = workflow.WithActivityOptions(ctx, options)

	var ibpManager ibpActivity.ManageInfiniBandPartition

	err = workflow.ExecuteActivity(ctx, ibpManager.UpdateInfiniBandPartitionsInDB, parsedSiteID, ibpInventory).Get(ctx, nil)
	if err != nil {
		logger.Warn().Err(err).Msg("failed to execute activity: UpdateInfiniBandPartitionsInDB")
	}

	// Record latency for this inventory call
	var inventoryMetricsManager cwm.ManageInventoryMetrics

	serr := workflow.ExecuteActivity(ctx, inventoryMetricsManager.RecordLatency, parsedSiteID, "UpdateInfiniBandPartitionInventory", err != nil, time.Since(startTime)).Get(ctx, nil)
	if serr != nil {
		logger.Warn().Err(serr).Msg("failed to execute activity: RecordLatency")
	}

	logger.Info().Msg("completing workflow")

	// Return original error from inventory activity, if any
	return err
}
