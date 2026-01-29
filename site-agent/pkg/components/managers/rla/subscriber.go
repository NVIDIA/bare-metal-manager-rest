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

// RegisterSubscriber registers the RLA Rack workflows with the Temporal client
func (api *API) RegisterSubscriber() error {
	// Check if RLA is enabled
	if !ManagerAccess.Conf.EB.RLA.Enabled {
		ManagerAccess.Data.EB.Log.Info().Msg("RLA: RLA is disabled, skipping workflow registration")
		return nil
	}

	rackManager := swa.NewManageRack(ManagerAccess.Data.EB.Managers.RLA.Client)

	// Register the subscribers here
	ManagerAccess.Data.EB.Log.Info().Msg("RLA: Registering the rack workflows")

	/// Register workflows

	// GetRackByID
	ManagerAccess.Data.EB.Managers.Workflow.Temporal.Worker.RegisterWorkflow(sww.GetRackByID)
	ManagerAccess.Data.EB.Log.Info().Msg("RLA: successfully registered GetRackByID workflow")

	// GetListOfRacks
	ManagerAccess.Data.EB.Managers.Workflow.Temporal.Worker.RegisterWorkflow(sww.GetListOfRacks)
	ManagerAccess.Data.EB.Log.Info().Msg("RLA: successfully registered GetListOfRacks workflow")

	/// Register activities

	// GetRackByID activity
	ManagerAccess.Data.EB.Managers.Workflow.Temporal.Worker.RegisterActivity(rackManager.GetRackByID)
	ManagerAccess.Data.EB.Log.Info().Msg("RLA: successfully registered GetRackByID activity")

	// GetListOfRacks activity
	ManagerAccess.Data.EB.Managers.Workflow.Temporal.Worker.RegisterActivity(rackManager.GetListOfRacks)
	ManagerAccess.Data.EB.Log.Info().Msg("RLA: successfully registered GetListOfRacks activity")

	return nil
}
