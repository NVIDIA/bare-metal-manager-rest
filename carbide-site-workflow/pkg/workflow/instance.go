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

package workflow

import (
	"time"

	"github.com/rs/zerolog/log"

	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"

	cwssaws "github.com/NVIDIA/carbide-rest-api/carbide-rest-api-schema/schema/site-agent/workflows/v1"

	"github.com/NVIDIA/carbide-rest-api/carbide-site-workflow/pkg/activity"
)

// UpdateInstance is a workflow to update Instance data using then UpdateInstanceOnSite activity
func UpdateInstance(ctx workflow.Context, updateRequest *cwssaws.InstanceConfigUpdateRequest) error {
	logger := log.With().Str("Workflow", "Instance").Str("Action", "Update").Str("Instance ID", updateRequest.InstanceId.String()).Logger()

	logger.Info().Msg("Starting workflow")

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

	var instanceManager activity.ManageInstance

	err := workflow.ExecuteActivity(ctx, instanceManager.UpdateInstanceOnSite, updateRequest).Get(ctx, nil)
	if err != nil {
		logger.Error().Err(err).Str("Activity", "UpdateInstanceOnSite").Msg("Failed to execute activity from workflow")
		return err
	}

	logger.Info().Msg("Completing workflow")

	return nil
}

// CreateInstanceV2 is a workflow to create (allocate) new Instances using the CreateInstanceOnSite activity
// V1 (CreateInstance) is found in cloud-workflow and uses a different activity that does not speak
// to carbide directly.
func CreateInstanceV2(ctx workflow.Context, request *cwssaws.InstanceAllocationRequest) error {
	logger := log.With().Str("Workflow", "Instance").Str("Action", "Create").Str("Machine ID", request.MachineId.Id).Logger()

	logger.Info().Msg("Starting workflow")

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

	var instanceManager activity.ManageInstance

	err := workflow.ExecuteActivity(ctx, instanceManager.CreateInstanceOnSite, request).Get(ctx, nil)
	if err != nil {
		logger.Error().Err(err).Str("Activity", "CreateInstanceOnSite").Msg("Failed to execute activity from workflow")
		return err
	}

	logger.Info().Msg("Completing workflow")

	return nil
}

// DeleteInstanceV2 is a workflow to delete new Instances using the DeleteInstanceOnSite activity
// V1 (DeleteInstance) is found in cloud-workflow and uses a different activity that does not speak
// to carbide directly.
func DeleteInstanceV2(ctx workflow.Context, request *cwssaws.InstanceReleaseRequest) error {

	logger := log.With().Str("Workflow", "Instance").Str("Action", "Delete").Str("Request", request.String()).Logger()

	logger.Info().Msg("Starting workflow")

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

	var instanceManager activity.ManageInstance

	err := workflow.ExecuteActivity(ctx, instanceManager.DeleteInstanceOnSite, request).Get(ctx, nil)
	if err != nil {
		logger.Error().Err(err).Str("Activity", "DeleteInstanceOnSite").Msg("Failed to execute activity from workflow")
		return err
	}

	logger.Info().Msg("Completing workflow")

	return nil
}

// RebootInstance is a workflow to reboot Instances using the RebootInstanceOnSite activity
func RebootInstance(ctx workflow.Context, request *cwssaws.InstancePowerRequest) error {
	logger := log.With().Str("Workflow", "Instance").Str("Action", "Reboot").Str("Machine ID", request.MachineId.Id).Logger()

	logger.Info().Msg("Starting workflow")

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

	var instanceManager activity.ManageInstance

	err := workflow.ExecuteActivity(ctx, instanceManager.RebootInstanceOnSite, request).Get(ctx, nil)
	if err != nil {
		logger.Error().Err(err).Str("Activity", "RebootInstanceOnSite").Msg("Failed to execute activity from workflow")
		return err
	}

	logger.Info().Msg("Completing workflow")

	return nil
}

func DiscoverInstanceInventory(ctx workflow.Context) error {
	logger := log.With().Str("Workflow", "DiscoverInstanceInventory").Logger()

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

	// Invoke DiscoverInstanceInventory activity
	var instanceInventoryManager activity.ManageInstanceInventory

	err := workflow.ExecuteActivity(ctx, instanceInventoryManager.DiscoverInstanceInventory).Get(ctx, nil)
	if err != nil {
		logger.Error().Err(err).Str("Activity", "DiscoverInstanceInventory").Msg("Failed to execute activity from workflow")
		return err
	}

	logger.Info().Msg("Completing workflow")

	return nil
}
