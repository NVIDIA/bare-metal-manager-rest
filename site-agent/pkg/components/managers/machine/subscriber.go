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
	swa "github.com/nvidia/carbide-rest/site-workflow/pkg/activity"
	sww "github.com/nvidia/carbide-rest/site-workflow/pkg/workflow"
)

// RegisterSubscriber registers the Machine workflows/activities with the Temporal client
func (api *API) RegisterSubscriber() error {
	// Register subscriber workflows
	ManagerAccess.Data.EB.Log.Info().Msg("Machine: Registering the subscribers")

	machineManager := swa.NewManageMachine(ManagerAccess.Data.EB.Managers.Carbide.Client)

	// Register workflows

	// Sync workflows
	// Set Maintenance Mode workflow
	ManagerAccess.Data.EB.Managers.Workflow.Temporal.Worker.RegisterWorkflow(sww.SetMachineMaintenance)
	ManagerAccess.Data.EB.Log.Info().Msg("Machine: successfully registered the Set Machine Maintenance workflow")

	ManagerAccess.Data.EB.Managers.Workflow.Temporal.Worker.RegisterWorkflow(sww.UpdateMachineMetadata)
	ManagerAccess.Data.EB.Log.Info().Msg("Machine: successfully registered the Update Machine Metadata workflow")

	// Register activities

	// Sync workflow activities
	// Register Machine activities
	ManagerAccess.Data.EB.Managers.Workflow.Temporal.Worker.RegisterActivity(machineManager.SetMachineMaintenanceOnSite)
	ManagerAccess.Data.EB.Log.Info().Msg("Machine: successfully registered the Set Machine Maintenance activity")

	ManagerAccess.Data.EB.Managers.Workflow.Temporal.Worker.RegisterActivity(machineManager.UpdateMachineMetadataOnSite)
	ManagerAccess.Data.EB.Log.Info().Msg("Machine: successfully registered the Update Machine Metadata activity")

	return nil
}
