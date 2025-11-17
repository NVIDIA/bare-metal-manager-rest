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

package expectedmachine

import (
	"github.com/google/uuid"

	swa "github.com/NVIDIA/carbide-rest-api/carbide-site-workflow/pkg/activity"
	sww "github.com/NVIDIA/carbide-rest-api/carbide-site-workflow/pkg/workflow"
)

// RegisterPublisher registers the ExpectedMachineWorkflows with the Temporal client
func (api *API) RegisterPublisher() error {
	// Register the publishers here

	// Collect and Publish ExpectedMachine Inventory workflow
	ManagerAccess.Data.EB.Managers.Workflow.Temporal.Worker.RegisterWorkflow(sww.DiscoverExpectedMachineInventory)
	ManagerAccess.Data.EB.Log.Info().Msg("ExpectedMachine: successfully registered the DiscoverExpectedMachineInventory workflow")

	inventoryManager := swa.NewManageExpectedMachineInventory(swa.ManageInventoryConfig{
		SiteID:                uuid.MustParse(ManagerAccess.Conf.EB.Temporal.ClusterID),
		CarbideAtomicClient:   ManagerAccess.Data.EB.Managers.Carbide.Client,
		TemporalPublishClient: ManagerAccess.Data.EB.Managers.Workflow.Temporal.Publisher,
		TemporalPublishQueue:  ManagerAccess.Conf.EB.Temporal.TemporalPublishQueue,
		SitePageSize:          InventoryCarbidePageSize,
		CloudPageSize:         InventoryCloudPageSize,
	})
	ManagerAccess.Data.EB.Managers.Workflow.Temporal.Worker.RegisterActivity(inventoryManager.DiscoverExpectedMachineInventory)
	ManagerAccess.Data.EB.Log.Info().Msg("ExpectedMachine: successfully registered the DiscoverExpectedMachineInventory activity")

	_ = api.RegisterCron()

	return nil
}
