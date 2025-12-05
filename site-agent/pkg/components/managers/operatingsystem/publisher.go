// SPDX-FileCopyrightText: Copyright (c) 2021-2023 NVIDIA CORPORATION & AFFILIATES. All rights reserved.
// SPDX-License-Identifier: LicenseRef-NvidiaProprietary
//
// NVIDIA CORPORATION, its affiliates and licensors retain all intellectual
// property and proprietary rights in and to this material, related
// documentation and any modifications thereto. Any use, reproduction,
// disclosure or distribution of this material and related documentation
// without an express license agreement from NVIDIA CORPORATION or
// its affiliates is strictly prohibited.

package operatingsystem

import (
	"github.com/google/uuid"

	swa "github.com/nvidia/carbide-rest/site-workflow/pkg/activity"
	sww "github.com/nvidia/carbide-rest/site-workflow/pkg/workflow"
)

// RegisterPublisher registers the OsImage Workflows with the Temporal client
func (api *API) RegisterPublisher() error {
	// Register publisher workflows

	// Collect and Publish OsImage Inventory workflow
	ManagerAccess.Data.EB.Managers.Workflow.Temporal.Worker.RegisterWorkflow(sww.DiscoverOsImageInventory)
	ManagerAccess.Data.EB.Log.Info().Msg("OsImage: successfully registered DiscoverOsImageInventory workflow")

	// Register OsImage activity for discovering and publishing OsImage Inventory
	OsImageInventoryManager := swa.NewManageOsImageInventory(swa.ManageInventoryConfig{
		SiteID:                uuid.MustParse(ManagerAccess.Conf.EB.Temporal.ClusterID),
		CarbideAtomicClient:   ManagerAccess.Data.EB.Managers.Carbide.Client,
		TemporalPublishClient: ManagerAccess.Data.EB.Managers.Workflow.Temporal.Publisher,
		TemporalPublishQueue:  ManagerAccess.Conf.EB.Temporal.TemporalPublishQueue,
		SitePageSize:          InventoryCarbidePageSize,
		CloudPageSize:         InventoryCloudPageSize,
	})

	ManagerAccess.Data.EB.Managers.Workflow.Temporal.Worker.RegisterActivity(OsImageInventoryManager.DiscoverOsImageInventory)
	ManagerAccess.Data.EB.Log.Info().Msg("OsImage: successfully registered DiscoverOsImageInventory activity")

	api.RegisterCron()
	return nil
}
