/*
 * SPDX-FileCopyrightText: Copyright (c) 2021-2023 NVIDIA CORPORATION & AFFILIATES. All rights reserved.
 * SPDX-License-Identifier: LicenseRef-NvidiaProprietary
 *
 * NVIDIA CORPORATION, its affiliates and licensors retain all intellectual
 * property and proprietary rights in and to this material, related
 * documentation and any modifications thereto. Any use, reproduction,
 * disclosure or distribution of this material and related documentation
 * without an express license agreement from NVIDIA CORPORATION or
 * its affiliates is strictly prohibited.
 */

package subnet

import (
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"

	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"

	subnetActivity "github.com/NVIDIA/carbide-rest-api/carbide-rest-workflow/pkg/activity/subnet"

	cwm "github.com/NVIDIA/carbide-rest-api/carbide-rest-workflow/internal/metrics"
)

// RecordSubnetLifecycleMetrics is a Temporal workflow that collects and records Subnet lifecycle metrics
func RecordSubnetLifecycleMetrics(ctx workflow.Context, siteID uuid.UUID, subnetLifecycleEvents []cwm.InventoryObjectLifecycleEvent) error {
	logger := log.With().Str("Workflow", "RecordSubnetLifecycleMetrics").Str("Site ID", siteID.String()).Logger()

	logger.Info().Msg("starting workflow")

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

	var lifecycleMetricsManager subnetActivity.ManageSubnetLifecycleMetrics

	err := workflow.ExecuteActivity(ctx, lifecycleMetricsManager.RecordSubnetStatusTransitionMetrics, siteID, subnetLifecycleEvents).Get(ctx, nil)
	if err != nil {
		logger.Warn().Err(err).Msg("failed to execute activity: RecordSubnetStatusTransitionMetrics")
		return err
	}

	logger.Info().Msg("completing workflow")

	return nil
}
