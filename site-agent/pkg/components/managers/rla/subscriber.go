// SPDX-FileCopyrightText: Copyright (c) 2021-2023 NVIDIA CORPORATION & AFFILIATES. All rights reserved.
// SPDX-License-Identifier: LicenseRef-NvidiaProprietary
//
// NVIDIA CORPORATION, its affiliates and licensors retain all intellectual
// property and proprietary rights in and to this material, related
// documentation and any modifications thereto. Any use, reproduction,
// disclosure or distribution of this material and related documentation
// without an express license agreement from NVIDIA CORPORATION or
// its affiliates is strictly prohibited.

package rla

import (
	swa "github.com/nvidia/carbide-rest/site-workflow/pkg/activity"
	sww "github.com/nvidia/carbide-rest/site-workflow/pkg/workflow"
)

// RegisterSubscriber registers the RLA Rack and Tray workflows with the Temporal client
func (api *API) RegisterSubscriber() error {
	// Check if RLA is enabled
	if !ManagerAccess.Conf.EB.RLA.Enabled {
		ManagerAccess.Data.EB.Log.Info().Msg("RLA: RLA is disabled, skipping workflow registration")
		return nil
	}

	rackManager := swa.NewManageRack(ManagerAccess.Data.EB.Managers.RLA.Client)
	trayManager := swa.NewManageTray(ManagerAccess.Data.EB.Managers.RLA.Client)

	// Register the subscribers here
	ManagerAccess.Data.EB.Log.Info().Msg("RLA: Registering the rack workflows")

	/// Register rack workflows

	// GetRack
	ManagerAccess.Data.EB.Managers.Workflow.Temporal.Worker.RegisterWorkflow(sww.GetRack)
	ManagerAccess.Data.EB.Log.Info().Msg("RLA: successfully registered GetRack workflow")

	// GetRacks
	ManagerAccess.Data.EB.Managers.Workflow.Temporal.Worker.RegisterWorkflow(sww.GetRacks)
	ManagerAccess.Data.EB.Log.Info().Msg("RLA: successfully registered GetRacks workflow")

	/// Register rack activities

	// GetRack activity
	ManagerAccess.Data.EB.Managers.Workflow.Temporal.Worker.RegisterActivity(rackManager.GetRack)
	ManagerAccess.Data.EB.Log.Info().Msg("RLA: successfully registered GetRack activity")

	// GetRacks activity
	ManagerAccess.Data.EB.Managers.Workflow.Temporal.Worker.RegisterActivity(rackManager.GetRacks)
	ManagerAccess.Data.EB.Log.Info().Msg("RLA: successfully registered GetRacks activity")

	// Register the tray subscribers here
	ManagerAccess.Data.EB.Log.Info().Msg("RLA: Registering the tray workflows")

	/// Register tray workflows

	// GetTray
	ManagerAccess.Data.EB.Managers.Workflow.Temporal.Worker.RegisterWorkflow(sww.GetTray)
	ManagerAccess.Data.EB.Log.Info().Msg("RLA: successfully registered GetTray workflow")

	// GetTrays (also handles taskId resolution internally via GetTaskComponentIDs activity)
	ManagerAccess.Data.EB.Managers.Workflow.Temporal.Worker.RegisterWorkflow(sww.GetTrays)
	ManagerAccess.Data.EB.Log.Info().Msg("RLA: successfully registered GetTrays workflow")

	/// Register tray activities

	// GetTray activity
	ManagerAccess.Data.EB.Managers.Workflow.Temporal.Worker.RegisterActivity(trayManager.GetTray)
	ManagerAccess.Data.EB.Log.Info().Msg("RLA: successfully registered GetTray activity")

	// GetTrays activity
	ManagerAccess.Data.EB.Managers.Workflow.Temporal.Worker.RegisterActivity(trayManager.GetTrays)
	ManagerAccess.Data.EB.Log.Info().Msg("RLA: successfully registered GetTrays activity")

	// GetTaskComponentIDs activity
	ManagerAccess.Data.EB.Managers.Workflow.Temporal.Worker.RegisterActivity(trayManager.GetTaskComponentIDs)
	ManagerAccess.Data.EB.Log.Info().Msg("RLA: successfully registered GetTaskComponentIDs activity")

	return nil
}
