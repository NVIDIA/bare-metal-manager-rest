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

package instance

import (
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"

	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"

	instanceActivity "github.com/NVIDIA/carbide-rest-api/carbide-rest-workflow/pkg/activity/instance"

	cwm "github.com/NVIDIA/carbide-rest-api/carbide-rest-workflow/internal/metrics"
)

// RecordInstanceLifecycleMetrics is a Temporal workflow that collects and records Instance lifecycle metrics
func RecordInstanceLifecycleMetrics(ctx workflow.Context, siteID uuid.UUID, instanceLifecycleEvents []cwm.InventoryObjectLifecycleEvent) error {
	logger := log.With().Str("Workflow", "RecordInstanceLifecycleMetrics").Str("Site ID", siteID.String()).Logger()

	logger.Info().Msg("starting workflow")

	retrypolicy := &temporal.RetryPolicy{
		InitialInterval:    5 * time.Second,
		BackoffCoefficient: 2.0,
		MaximumInterval:    30 * time.Second,
		MaximumAttempts:    2,
	}
	options := workflow.ActivityOptions{
		StartToCloseTimeout: 30 * time.Second,
		RetryPolicy:         retrypolicy,
	}

	ctx = workflow.WithActivityOptions(ctx, options)

	var lifecycleMetricsManager instanceActivity.ManageInstanceLifecycleMetrics

	err := workflow.ExecuteActivity(ctx, lifecycleMetricsManager.RecordInstanceStatusTransitionMetrics, siteID, instanceLifecycleEvents).Get(ctx, nil)
	if err != nil {
		logger.Warn().Err(err).Msg("failed to execute activity: RecordInstanceStatusTransitionMetrics")
		return err
	}

	logger.Info().Msg("completing workflow")

	return nil
}
