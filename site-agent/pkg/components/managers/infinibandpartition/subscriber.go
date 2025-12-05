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
	"go.temporal.io/sdk/activity"
	workflow "go.temporal.io/sdk/workflow"

	swa "github.com/nvidia/carbide-rest/site-workflow/pkg/activity"
	sww "github.com/nvidia/carbide-rest/site-workflow/pkg/workflow"
)

// RegisterSubscriber registers the InfiniBandPartitionWorkflows with the Temporal client
func (api *API) RegisterSubscriber() error {

	// Register the subscribers here
	ManagerAccess.Data.EB.Log.Info().Msg("InfiniBandPartition: Registering the subscribers")

	ibpManager := swa.NewManageInfiniBandPartition(ManagerAccess.Data.EB.Managers.Carbide.Client)

	// Get InfiniBandPartition workflow interface
	infiniBandPartitionInterface := NewInfiniBandPartitionWorkflows(
		ManagerAccess.Data.EB.Managers.Workflow.Temporal.Publisher,
		ManagerAccess.Data.EB.Managers.Workflow.Temporal.Subscriber,
		ManagerAccess.Conf.EB,
	)

	// Register worfklow

	// Sync workflows

	// CreateInfiniBandPartitionV2
	ManagerAccess.Data.EB.Managers.Workflow.Temporal.Worker.RegisterWorkflow(sww.CreateInfiniBandPartitionV2)
	ManagerAccess.Data.EB.Log.Info().Msg("InfiniBandPartition: successfully registered CreateInfiniBandPartitionV2 workflow")

	// DeleteInfiniBandPartitionV2
	ManagerAccess.Data.EB.Managers.Workflow.Temporal.Worker.RegisterWorkflow(sww.DeleteInfiniBandPartitionV2)
	ManagerAccess.Data.EB.Log.Info().Msg("InfiniBandPartition: successfully registered DeleteInfiniBandPartitionV2 workflow")

	wflowRegisterOptions := workflow.RegisterOptions{
		Name: "CreateInfiniBandPartition",
	}
	ManagerAccess.Data.EB.Managers.Workflow.Temporal.Worker.RegisterWorkflowWithOptions(
		ManagerAccess.API.InfiniBandPartition.CreateInfiniBandPartition, wflowRegisterOptions,
	)

	ManagerAccess.Data.EB.Log.Info().Msg("InfiniBandPartition: successfully registered the create InfiniBandPartition workflow")

	wflowRegisterOptions = workflow.RegisterOptions{
		Name: "DeleteInfiniBandPartition",
	}
	ManagerAccess.Data.EB.Managers.Workflow.Temporal.Worker.RegisterWorkflowWithOptions(
		ManagerAccess.API.InfiniBandPartition.DeleteInfiniBandPartition, wflowRegisterOptions,
	)

	ManagerAccess.Data.EB.Log.Info().Msg("InfiniBandPartition: successfully registered the delete InfiniBandPartition workflow")

	// Register activity

	// Sync workflow activities

	// CreateInfiniBandPartitionOnSite
	ManagerAccess.Data.EB.Managers.Workflow.Temporal.Worker.RegisterActivity(ibpManager.CreateInfiniBandPartitionOnSite)
	ManagerAccess.Data.EB.Log.Info().Msg("VPC: successfully registered the CreateInfiniBandPartitionOnSite activity")

	// DeleteInfiniBandPartitionOnSite
	ManagerAccess.Data.EB.Managers.Workflow.Temporal.Worker.RegisterActivity(ibpManager.DeleteInfiniBandPartitionOnSite)
	ManagerAccess.Data.EB.Log.Info().Msg("VPC: successfully registered the DeleteInfiniBandPartitionOnSite activity")

	activityRegisterOptions := activity.RegisterOptions{
		Name: "CreateInfiniBandPartitionActivity",
	}

	ManagerAccess.Data.EB.Managers.Workflow.Temporal.Worker.RegisterActivityWithOptions(
		infiniBandPartitionInterface.CreateInfiniBandPartitionActivity, activityRegisterOptions,
	)
	ManagerAccess.Data.EB.Log.Info().Msg("InfiniBandPartition: successfully registered the Create InfiniBandPartition activity")

	activityRegisterOptions = activity.RegisterOptions{
		Name: "DeleteInfiniBandPartitionActivity",
	}

	ManagerAccess.Data.EB.Managers.Workflow.Temporal.Worker.RegisterActivityWithOptions(
		infiniBandPartitionInterface.DeleteInfiniBandPartitionActivity, activityRegisterOptions,
	)
	ManagerAccess.Data.EB.Log.Info().Msg("InfiniBandPartition: successfully registered the Delete InfiniBandPartition activity")

	return nil
}

// RegisterSubscribers - this is method 2 of registering the subscriber
func RegisterSubscribers() {
	// Register the subscribers here
	ManagerAccess.Data.EB.Log.Info().Msg("InfiniBandPartition: Registering the subscribers")
	ManagerAccess.API.Orchestrator.AddWorkflow(ManagerAccess.API.InfiniBandPartition.CreateInfiniBandPartition)
}
