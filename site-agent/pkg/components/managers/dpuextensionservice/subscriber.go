// SPDX-FileCopyrightText: Copyright (c) 2021-2023 NVIDIA CORPORATION & AFFILIATES. All rights reserved.
// SPDX-License-Identifier: LicenseRef-NvidiaProprietary
//
// NVIDIA CORPORATION, its affiliates and licensors retain all intellectual
// property and proprietary rights in and to this material, related
// documentation and any modifications thereto. Any use, reproduction,
// disclosure or distribution of this material and related documentation
// without an express license agreement from NVIDIA CORPORATION or
// its affiliates is strictly prohibited.

package dpuextensionservice

import (
	swa "github.com/nvidia/carbide-rest/site-workflow/pkg/activity"
	sww "github.com/nvidia/carbide-rest/site-workflow/pkg/workflow"
)

// RegisterSubscriber registers DPU Extension Service workflows Site Agent subscribes to execute
func (api *API) RegisterSubscriber() error {
	//  Register Workflows

	// Sync workflows
	// Register CreateDpuExtensionService workflow
	ManagerAccess.Data.EB.Managers.Workflow.Temporal.Worker.RegisterWorkflow(sww.CreateDpuExtensionService)
	ManagerAccess.Data.EB.Log.Info().Msg("DpuExtensionService: successfully registered CreateDpuExtensionService workflow")

	// Register UpdateDpuExtensionService workflow
	ManagerAccess.Data.EB.Managers.Workflow.Temporal.Worker.RegisterWorkflow(sww.UpdateDpuExtensionService)
	ManagerAccess.Data.EB.Log.Info().Msg("DpuExtensionService: successfully registered UpdateDpuExtensionService workflow")

	// Register DeleteDpuExtensionService workflow
	ManagerAccess.Data.EB.Managers.Workflow.Temporal.Worker.RegisterWorkflow(sww.DeleteDpuExtensionService)
	ManagerAccess.Data.EB.Log.Info().Msg("DpuExtensionService: successfully registered DeleteDpuExtensionService workflow")

	// Register GetDpuExtensionServiceVersionsInfo workflow
	ManagerAccess.Data.EB.Managers.Workflow.Temporal.Worker.RegisterWorkflow(sww.GetDpuExtensionServiceVersionsInfo)
	ManagerAccess.Data.EB.Log.Info().Msg("DpuExtensionService: successfully registered GetDpuExtensionServiceVersionsInfo workflow")

	// Register Activities
	dpuExtServiceManager := swa.NewManageDpuExtensionService(ManagerAccess.Data.EB.Managers.Carbide.Client)

	// Sync workflow activities
	// Register CreateDpuExtensionServiceOnSite
	ManagerAccess.Data.EB.Managers.Workflow.Temporal.Worker.RegisterActivity(dpuExtServiceManager.CreateDpuExtensionServiceOnSite)
	ManagerAccess.Data.EB.Log.Info().Msg("DpuExtensionService: successfully registered CreateDpuExtensionServiceOnSite activity")

	// Register UpdateDpuExtensionServiceOnSite
	ManagerAccess.Data.EB.Managers.Workflow.Temporal.Worker.RegisterActivity(dpuExtServiceManager.UpdateDpuExtensionServiceOnSite)
	ManagerAccess.Data.EB.Log.Info().Msg("DpuExtensionService: successfully registered UpdateDpuExtensionServiceOnSite activity")

	// Register DeleteDpuExtensionServiceOnSite
	ManagerAccess.Data.EB.Managers.Workflow.Temporal.Worker.RegisterActivity(dpuExtServiceManager.DeleteDpuExtensionServiceOnSite)
	ManagerAccess.Data.EB.Log.Info().Msg("DpuExtensionService: successfully registered DeleteDpuExtensionServiceOnSite activity")

	// Register GetDpuExtensionServiceVersionsInfoOnSite activity
	ManagerAccess.Data.EB.Managers.Workflow.Temporal.Worker.RegisterActivity(dpuExtServiceManager.GetDpuExtensionServiceVersionsInfoOnSite)
	ManagerAccess.Data.EB.Log.Info().Msg("DpuExtensionService: successfully registered GetDpuExtensionServiceVersionsInfoOnSite activity")

	return nil
}
