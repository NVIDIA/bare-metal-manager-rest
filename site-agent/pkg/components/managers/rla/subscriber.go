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
	swa "github.com/nvidia/bare-metal-manager-rest/site-workflow/pkg/activity"
	sww "github.com/nvidia/bare-metal-manager-rest/site-workflow/pkg/workflow"
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

	// GetRack
	ManagerAccess.Data.EB.Managers.Workflow.Temporal.Worker.RegisterWorkflow(sww.GetRack)
	ManagerAccess.Data.EB.Log.Info().Msg("RLA: successfully registered GetRack workflow")

	// GetRacks
	ManagerAccess.Data.EB.Managers.Workflow.Temporal.Worker.RegisterWorkflow(sww.GetRacks)
	ManagerAccess.Data.EB.Log.Info().Msg("RLA: successfully registered GetRacks workflow")

	// ValidateComponents
	ManagerAccess.Data.EB.Managers.Workflow.Temporal.Worker.RegisterWorkflow(sww.ValidateComponents)
	ManagerAccess.Data.EB.Log.Info().Msg("RLA: successfully registered ValidateComponents workflow")

	/// Register activities

	// GetRack activity
	ManagerAccess.Data.EB.Managers.Workflow.Temporal.Worker.RegisterActivity(rackManager.GetRack)
	ManagerAccess.Data.EB.Log.Info().Msg("RLA: successfully registered GetRack activity")

	// GetRacks activity
	ManagerAccess.Data.EB.Managers.Workflow.Temporal.Worker.RegisterActivity(rackManager.GetRacks)
	ManagerAccess.Data.EB.Log.Info().Msg("RLA: successfully registered GetRacks activity")

	// ValidateComponents activity
	ManagerAccess.Data.EB.Managers.Workflow.Temporal.Worker.RegisterActivity(rackManager.ValidateComponents)
	ManagerAccess.Data.EB.Log.Info().Msg("RLA: successfully registered ValidateComponents activity")

	return nil
}
