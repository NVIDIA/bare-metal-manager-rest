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
	swa "github.com/nvidia/carbide-rest/site-workflow/pkg/activity"
	sww "github.com/nvidia/carbide-rest/site-workflow/pkg/workflow"
)

// RegisterSubscriber registers the Machine workflows/activities with the Temporal client
func (api *API) RegisterSubscriber() error {
	// Register subscriber workflows
	ManagerAccess.Data.EB.Log.Info().Msg("Machine: Registering the subscribers")

	osImageManager := swa.NewManageOperatingSystem(ManagerAccess.Data.EB.Managers.Carbide.Client)

	// Register workflows

	// Sync workflows
	// Register CreateOsImage workflow
	ManagerAccess.Data.EB.Managers.Workflow.Temporal.Worker.RegisterWorkflow(sww.CreateOsImage)
	ManagerAccess.Data.EB.Log.Info().Msg("OsImage: successfully registered the CreateOsImage workflow")

	// Register UpdateOsImage workflow
	ManagerAccess.Data.EB.Managers.Workflow.Temporal.Worker.RegisterWorkflow(sww.UpdateOsImage)
	ManagerAccess.Data.EB.Log.Info().Msg("OsImage: successfully registered the UpdateOsImage workflow")

	// Register DeleteOsImage workflow
	ManagerAccess.Data.EB.Managers.Workflow.Temporal.Worker.RegisterWorkflow(sww.DeleteOsImage)
	ManagerAccess.Data.EB.Log.Info().Msg("OsImage: successfully registered the DeleteOsImage workflow")

	// Register activities

	// Sync workflow activities
	// Register CreateOsImageOnSite
	ManagerAccess.Data.EB.Managers.Workflow.Temporal.Worker.RegisterActivity(osImageManager.CreateOsImageOnSite)
	ManagerAccess.Data.EB.Log.Info().Msg("OsImage: successfully registered CreateOsImageOnSite activity")

	// Register UpdateOsImageOnSite
	ManagerAccess.Data.EB.Managers.Workflow.Temporal.Worker.RegisterActivity(osImageManager.UpdateOsImageOnSite)
	ManagerAccess.Data.EB.Log.Info().Msg("OsImage: successfully registered UpdateOsImageOnSite activity")

	// Register DeleteOsImageOnSite
	ManagerAccess.Data.EB.Managers.Workflow.Temporal.Worker.RegisterActivity(osImageManager.DeleteOsImageOnSite)
	ManagerAccess.Data.EB.Log.Info().Msg("OsImage: successfully registered DeleteOsImageOnSite activity")

	return nil
}
