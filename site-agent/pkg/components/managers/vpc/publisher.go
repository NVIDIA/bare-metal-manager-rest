// SPDX-FileCopyrightText: Copyright (c) 2021-2023 NVIDIA CORPORATION & AFFILIATES. All rights reserved.
// SPDX-License-Identifier: LicenseRef-NvidiaProprietary
//
// NVIDIA CORPORATION, its affiliates and licensors retain all intellectual
// property and proprietary rights in and to this material, related
// documentation and any modifications thereto. Any use, reproduction,
// disclosure or distribution of this material and related documentation
// without an express license agreement from NVIDIA CORPORATION or
// its affiliates is strictly prohibited.

package vpc

import (
	"github.com/google/uuid"
	swa "github.com/nvidia/carbide-rest/site-workflow/pkg/activity"
	sww "github.com/nvidia/carbide-rest/site-workflow/pkg/workflow"
	"go.temporal.io/sdk/activity"
)

// RegisterPublisher registers the VPCWorkflows with the Temporal client
func (api *API) RegisterPublisher() error {
	// Register the publishers here
	ManagerAccess.Data.EB.Log.Info().Msg("VPC: Registering the publishers")

	// Get VPC workflow interface
	VPCinterface := NewVPCWorkflows(
		ManagerAccess.Data.EB.Managers.Workflow.Temporal.Publisher,
		ManagerAccess.Data.EB.Managers.Workflow.Temporal.Subscriber,
		ManagerAccess.Conf.EB,
	)

	activityRegisterOptions := activity.RegisterOptions{
		Name: "PublishVPCActivity",
	}

	ManagerAccess.Data.EB.Managers.Workflow.Temporal.Worker.RegisterActivityWithOptions(
		VPCinterface.PublishVPCActivity, activityRegisterOptions,
	)
	ManagerAccess.Data.EB.Log.Info().Msg("VPC: successfully registered the Publish VPC activity")

	activityRegisterOptions = activity.RegisterOptions{
		Name: "PublishVPCListActivity",
	}
	ManagerAccess.Data.EB.Managers.Workflow.Temporal.Worker.RegisterActivityWithOptions(
		VPCinterface.PublishVPCListActivity, activityRegisterOptions,
	)
	ManagerAccess.Data.EB.Log.Info().Msg("VPC: successfully registered the Publish VPC List activity")

	// Instance Inventory workflow
	ManagerAccess.Data.EB.Managers.Workflow.Temporal.Worker.RegisterWorkflow(sww.DiscoverVPCInventory)
	ManagerAccess.Data.EB.Log.Info().Msg("VPC: successfully registered the Discover VPC Inventory workflow")

	inventoryManager := swa.NewManageVPCInventory(swa.ManageInventoryConfig{
		SiteID:                uuid.MustParse(ManagerAccess.Conf.EB.Temporal.ClusterID),
		CarbideAtomicClient:   ManagerAccess.Data.EB.Managers.Carbide.Client,
		TemporalPublishClient: ManagerAccess.Data.EB.Managers.Workflow.Temporal.Publisher,
		TemporalPublishQueue:  ManagerAccess.Conf.EB.Temporal.TemporalPublishQueue,
		SitePageSize:          InventoryCarbidePageSize,
		CloudPageSize:         InventoryCloudPageSize,
	})
	ManagerAccess.Data.EB.Managers.Workflow.Temporal.Worker.RegisterActivity(inventoryManager.DiscoverVPCInventory)
	ManagerAccess.Data.EB.Log.Info().Msg("VPC: successfully registered the Discover VPC Inventory activity")

	api.RegisterCron()
	return nil
}
