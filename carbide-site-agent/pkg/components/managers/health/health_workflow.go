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

package health

import (
	"errors"
	"time"

	wflows "github.com/NVIDIA/carbide-rest-api/carbide-rest-api-schema/schema/site-agent/workflows/v1"
	"github.com/NVIDIA/carbide-rest-api/carbide-site-agent/pkg/conftypes"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/log"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

type HealthWorkflow struct {
	tcPublish   client.Client
	tcSubscribe client.Client
	cfg         *conftypes.Config
}

const (
	// Retry parameters for Temporal workflows
	RetryInterval                 = 2
	RetryCount                    = 10
	MaxTemporalActivityRetryCount = 3
)

// NewHealthWorkflows creates an instance for HealthWorkflows
func NewHealthWorkflows(TMPublish client.Client, TMSubscribe client.Client, CurrentCFG *conftypes.Config) HealthWorkflow {
	return HealthWorkflow{
		tcPublish:   TMPublish,
		tcSubscribe: TMSubscribe,
		cfg:         CurrentCFG,
	}
}

func GetHealth(ctx workflow.Context, TransactionID *wflows.TransactionID) (status wflows.HealthStatus, err error) {
	logger := workflow.GetLogger(ctx)
	withLogger := log.With(logger, "Workflow", "CreateHealthWorkflow", "ResourceRequest", TransactionID)
	withLogger.Info("Health: Starting  the Health Workflow")

	ManagerAccess.Data.EB.Log.Info().Interface("Request", TransactionID).Msg("Health: Starting  the Health Workflow")
	// Validations

	if TransactionID == nil {
		withLogger.Error("Health: TransactionID is nil")
		ManagerAccess.Data.EB.Log.Error().Msg("Health: TransactionID is nil")
		return status, errors.New("Health: TransactionID is nil")
	}
	if TransactionID.ResourceId == "" {
		withLogger.Error("Health: TransactionID.ResourceId is empty")
		ManagerAccess.Data.EB.Log.Error().Msg("Health: TransactionID.ResourceId is empty")
		return status, errors.New("Health: TransactionID.ResourceId is empty")
	}

	// Use default retry interval
	RetryInterval := 1 * time.Second

	retrypolicy := &temporal.RetryPolicy{
		InitialInterval:    RetryInterval,
		BackoffCoefficient: 2.0,
		MaximumInterval:    1 * time.Minute,
		MaximumAttempts:    MaxTemporalActivityRetryCount,
	}
	options := workflow.ActivityOptions{
		// Timeout options specify when to automatically timeout Activity functions.
		StartToCloseTimeout: 20 * time.Second,
		// Optionally provide a customized RetryPolicy.
		RetryPolicy: retrypolicy,
	}
	ctx = workflow.WithActivityOptions(ctx, options)
	Healthwflowinstance := HealthWorkflow{}

	err = workflow.ExecuteActivity(ctx, Healthwflowinstance.GetHealthActivity).Get(ctx, status)
	if err != nil {
		withLogger.Error("Health: Failed to get Health workflow", "Error", err)
		ManagerAccess.Data.EB.Log.Error().Interface("Error", err).Msg("Health: Failed to get Health")
		return status, err
	}

	withLogger.Info("Health: Successfully updated Health")
	ManagerAccess.Data.EB.Log.Info().Interface("Request", TransactionID).Msg("Health: Successfully updated Health")

	return status, err
}
