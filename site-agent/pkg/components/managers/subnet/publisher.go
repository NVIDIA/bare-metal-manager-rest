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
	"github.com/google/uuid"
	swa "github.com/nvidia/carbide-rest/site-workflow/pkg/activity"
	sww "github.com/nvidia/carbide-rest/site-workflow/pkg/workflow"
	"go.temporal.io/sdk/activity"
)

// RegisterPublisher registers the SubnetWorkflows with the Temporal client
func (api *API) RegisterPublisher() error {

	// Register the publishers here
	ManagerAccess.Data.EB.Log.Info().Msg("Subnet: Registering the publishers")

	// Get Subnet workflow interface
	Subnetinterface := NewSubnetWorkflows(
		ManagerAccess.Data.EB.Managers.Workflow.Temporal.Publisher,
		ManagerAccess.Data.EB.Managers.Workflow.Temporal.Subscriber,
		ManagerAccess.Conf.EB,
	)
	activityRegisterOptions := activity.RegisterOptions{
		Name: "PublishSubnetActivity",
	}
	ManagerAccess.Data.EB.Managers.Workflow.Temporal.Worker.RegisterActivityWithOptions(
		Subnetinterface.PublishSubnetActivity, activityRegisterOptions,
	)
	ManagerAccess.Data.EB.Log.Info().Msg("Subnet: successfully registered the Publish Subnet activity")

	// Subnet Inventory workflow
	ManagerAccess.Data.EB.Managers.Workflow.Temporal.Worker.RegisterWorkflow(sww.DiscoverSubnetInventory)
	ManagerAccess.Data.EB.Log.Info().Msg("Subnet: successfully registered the Discover Subnet Inventory workflow")

	inventoryManager := swa.NewManageSubnetInventory(swa.ManageInventoryConfig{
		SiteID:                uuid.MustParse(ManagerAccess.Conf.EB.Temporal.ClusterID),
		CarbideAtomicClient:   ManagerAccess.Data.EB.Managers.Carbide.Client,
		TemporalPublishClient: ManagerAccess.Data.EB.Managers.Workflow.Temporal.Publisher,
		TemporalPublishQueue:  ManagerAccess.Conf.EB.Temporal.TemporalPublishQueue,
		SitePageSize:          InventoryCarbidePageSize,
		CloudPageSize:         InventoryCloudPageSize,
	})
	ManagerAccess.Data.EB.Managers.Workflow.Temporal.Worker.RegisterActivity(inventoryManager.DiscoverSubnetInventory)
	ManagerAccess.Data.EB.Log.Info().Msg("Subnet: successfully registered the Discover Subnet Inventory activity")

	api.RegisterCron()

	return nil
}
