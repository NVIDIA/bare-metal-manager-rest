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
	wflows "github.com/nvidia/carbide-rest/workflow-schema/schema/site-agent/workflows/v1"
	"go.temporal.io/sdk/activity"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/log"
)

// PublishSubnetActivity - Publish Subnet Activity
func (ac *Workflows) PublishSubnetActivity(ctx context.Context, TransactionID *wflows.TransactionID, SubnetInfo *wflows.SubnetInfo) (workflowID string, err error) {
	ManagerAccess.Data.EB.Log.Info().Interface("Request", TransactionID).Msg("Subnet: Starting  the Publish Subnet Activity")

	// Use temporal logger for temporal logs
	logger := activity.GetLogger(ctx)
	withLogger := log.With(logger, "Activity", "PublishSubnetActivity", "ResourceReq", TransactionID)
	withLogger.Info("Subnet: Starting the Publish Subnet Activity")

	workflowOptions := client.StartWorkflowOptions{
		ID:        TransactionID.ResourceId,
		TaskQueue: ManagerAccess.Conf.EB.Temporal.TemporalPublishQueue,
	}

	we, err := ac.tcPublish.ExecuteWorkflow(ctx, workflowOptions, "UpdateSubnetInfo", ManagerAccess.Conf.EB.Temporal.TemporalSubscribeNamespace, TransactionID, SubnetInfo)
	if err != nil {
		return "", err
	}

	wid := we.GetID()
	return wid, nil
}
