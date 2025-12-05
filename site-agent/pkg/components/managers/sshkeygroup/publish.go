// SPDX-FileCopyrightText: Copyright (c) 2021-2023 NVIDIA CORPORATION & AFFILIATES. All rights reserved.
// SPDX-License-Identifier: LicenseRef-NvidiaProprietary
//
// NVIDIA CORPORATION, its affiliates and licensors retain all intellectual
// property and proprietary rights in and to this material, related
// documentation and any modifications thereto. Any use, reproduction,
// disclosure or distribution of this material and related documentation
// without an express license agreement from NVIDIA CORPORATION or
// its affiliates is strictly prohibited.

package sshkeygroup

import (
	"context"
	wflows "github.com/nvidia/carbide-rest/workflow-schema/schema/site-agent/workflows/v1"
	"go.temporal.io/sdk/activity"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/log"
)

// PublishSSHKeyGroupActivity - Publish SSHKeyGroup Info Activity
func (ac *Workflows) PublishSSHKeyGroupActivity(ctx context.Context, TransactionID *wflows.TransactionID, SSHInfo *wflows.SSHKeyGroupInfo) (workflowID string, err error) {
	ManagerAccess.Data.EB.Log.Info().Interface("Request", TransactionID).Msgf("SSHKeyGroup: Starting Publish Activity %v", SSHInfo)

	// Use temporal logger for temporal logs
	logger := activity.GetLogger(ctx)
	withLogger := log.With(logger, "Activity", "PublishSSHKeyGroupActivity", "ResourceReq", TransactionID)
	withLogger.Info("SSHKeyGroup: Starting the Publish SSHKeyGroup Activity")

	workflowOptions := client.StartWorkflowOptions{
		ID:        TransactionID.ResourceId,
		TaskQueue: ManagerAccess.Conf.EB.Temporal.TemporalPublishQueue,
	}
	var sshkeygroupresponse interface{}
	// Use the response as is
	ManagerAccess.Data.EB.Log.Info().Interface("Request", TransactionID).Msg("SSHKeyGroup: Using the response as is")
	sshkeygroupresponse = SSHInfo

	we, err := ac.tcPublish.ExecuteWorkflow(context.Background(), workflowOptions, "UpdateSSHKeyGroupInfo",
		ManagerAccess.Conf.EB.Temporal.TemporalSubscribeNamespace, TransactionID, sshkeygroupresponse)
	if err != nil {
		return "", err
	}

	wid := we.GetID()
	return wid, nil
}
