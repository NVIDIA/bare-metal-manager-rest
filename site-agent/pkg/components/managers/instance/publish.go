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
	wflows "github.com/nvidia/carbide-rest/workflow-schema/schema/site-agent/workflows/v1"
	"go.temporal.io/sdk/activity"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/log"
)

// PublishInstanceActivity - Publish Instance Activity
func (ac *Workflows) PublishInstanceActivity(ctx context.Context, TransactionID *wflows.TransactionID, InstanceInfo *wflows.InstanceInfo) (workflowID string, err error) {
	ManagerAccess.Data.EB.Log.Info().Interface("Request", TransactionID).Msg("Instance: Starting  the Publish Instance Activity")

	// Use temporal logger for temporal logs
	logger := activity.GetLogger(ctx)
	withLogger := log.With(logger, "Activity", "PublishInstanceActivity", "ResourceReq", TransactionID)
	withLogger.Info("Instance: Starting the Publish Instance Activity")

	workflowOptions := client.StartWorkflowOptions{
		ID:        TransactionID.ResourceId,
		TaskQueue: ManagerAccess.Conf.EB.Temporal.TemporalPublishQueue,
	}
	we, err := ac.tcPublish.ExecuteWorkflow(ctx, workflowOptions, "UpdateInstanceInfo", ManagerAccess.Conf.EB.Temporal.TemporalSubscribeNamespace, TransactionID, InstanceInfo)
	if err != nil {
		return "", err
	}
	wid := we.GetID()
	return wid, nil
}

// PublishInstancePowerStatus - Publish Instance Power Status
func (ac *Workflows) PublishInstancePowerStatus(ctx context.Context, TransactionID *wflows.TransactionID, InstanceInfo *wflows.InstanceRebootInfo) (workflowID string, err error) {
	ManagerAccess.Data.EB.Log.Info().Interface("Request", TransactionID).Msg("Instance: Starting  the Publish Instance power status Activity")

	// Use temporal logger for temporal logs
	logger := activity.GetLogger(ctx)
	withLogger := log.With(logger, "Activity", "PublishInstancePowerStatus", "ResourceReq", TransactionID)
	withLogger.Info("Instance: Starting the Publish Instance power status Activity")

	workflowOptions := client.StartWorkflowOptions{
		ID:        TransactionID.ResourceId,
		TaskQueue: ManagerAccess.Conf.EB.Temporal.TemporalPublishQueue,
	}

	we, err := ac.tcPublish.ExecuteWorkflow(ctx, workflowOptions, "UpdateInstanceRebootInfo", ManagerAccess.Conf.EB.Temporal.TemporalSubscribeNamespace, TransactionID, InstanceInfo)
	if err != nil {
		return "", err
	}

	wid := we.GetID()
	return wid, nil
}
