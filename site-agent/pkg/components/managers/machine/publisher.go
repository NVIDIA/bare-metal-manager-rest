// SPDX-FileCopyrightText: Copyright (c) 2021-2023 NVIDIA CORPORATION & AFFILIATES. All rights reserved.
// SPDX-License-Identifier: LicenseRef-NvidiaProprietary
//
// NVIDIA CORPORATION, its affiliates and licensors retain all intellectual
// property and proprietary rights in and to this material, related
// documentation and any modifications thereto. Any use, reproduction,
// disclosure or distribution of this material and related documentation
// without an express license agreement from NVIDIA CORPORATION or
// its affiliates is strictly prohibited.

package machine

import (
	"github.com/google/uuid"

	swa "github.com/nvidia/carbide-rest/site-workflow/pkg/activity"
	sww "github.com/nvidia/carbide-rest/site-workflow/pkg/workflow"
)

// RegisterPublisher registers the MachineWorkflows with the Temporal client
func (api *API) RegisterPublisher() error {
	// Register publisher workflows

	// Collect and Publish Machine Inventory workflow
	ManagerAccess.Data.EB.Managers.Workflow.Temporal.Worker.RegisterWorkflow(sww.CollectAndPublishMachineInventory)
	ManagerAccess.Data.EB.Log.Info().Msg("Machine: successfully registered the Collect and Publish Machine Inventory workflow")

	// Register Machine activity for Collect and Publish Machine Inventory
	machineInventoryManager := swa.NewManageMachineInventory(
		uuid.MustParse(ManagerAccess.Conf.EB.Temporal.ClusterID),
		ManagerAccess.Data.EB.Managers.Carbide.Client,
		ManagerAccess.Data.EB.Managers.Workflow.Temporal.Publisher,
		ManagerAccess.Conf.EB.Temporal.TemporalPublishQueue,
		InventoryCarbidePageSize,
		InventoryCloudPageSize,
	)
	ManagerAccess.Data.EB.Managers.Workflow.Temporal.Worker.RegisterActivity(machineInventoryManager.CollectAndPublishMachineInventory)
	ManagerAccess.Data.EB.Log.Info().Msg("Machine: successfully registered the Collect and Publish Machine Inventory activity")

	api.RegisterCron()
	return nil
}
