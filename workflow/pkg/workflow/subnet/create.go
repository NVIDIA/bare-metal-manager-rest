// SPDX-FileCopyrightText: Copyright (c) 2021-2023 NVIDIA CORPORATION & AFFILIATES. All rights reserved.
// SPDX-License-Identifier: LicenseRef-NvidiaProprietary
//
// NVIDIA CORPORATION, its affiliates and licensors retain all intellectual
// property and proprietary rights in and to this material, related
// documentation and any modifications thereto. Any use, reproduction,
// disclosure or distribution of this material and related documentation
// without an express license agreement from NVIDIA CORPORATION or
// its affiliates is strictly prohibited.

package subnet

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"

	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"

	"go.temporal.io/sdk/client"

	subnetActivity "github.com/nvidia/carbide-rest/workflow/pkg/activity/subnet"
	"github.com/nvidia/carbide-rest/workflow/pkg/queue"
)

// CreateSubnet is a Temporal workflow to create a new Subnet via Site Agent
func CreateSubnet(ctx workflow.Context, subnetID uuid.UUID, vpcID uuid.UUID) error {
	logger := log.With().Str("Workflow", "Subnet").Str("Action", "Create").Str("Subnet ID", subnetID.String()).
		Str("VPC ID", vpcID.String()).Logger()

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

	var subnetManager subnetActivity.ManageSubnet

	err := workflow.ExecuteActivity(ctx, subnetManager.CreateSubnetViaSiteAgent, subnetID, vpcID).Get(ctx, nil)
	if err != nil {
		logger.Warn().Err(err).Msg("failed to execute activity: CreateSubnetViaSiteAgent")
		return err
	}

	logger.Info().Msg("completing workflow")

	return nil
}

// ExecuteCreateSubnetWorkflow is a helper function to trigger execution of create Subnet workflow
func ExecuteCreateSubnetWorkflow(ctx context.Context, tc client.Client, subnetID uuid.UUID, vpcID uuid.UUID) (*string, error) {
	uid := uuid.New()

	workflowOptions := client.StartWorkflowOptions{
		ID:        "subnet-create-" + uid.String(),
		TaskQueue: queue.CloudTaskQueue,
	}

	we, err := tc.ExecuteWorkflow(ctx, workflowOptions, CreateSubnet, subnetID, vpcID)

	if err != nil {
		log.Error().Err(err).Msg("failed to execute workflow: CreateSubnet")
		return nil, err
	}

	wid := we.GetID()

	return &wid, nil
}
