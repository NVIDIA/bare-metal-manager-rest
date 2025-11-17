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
	swa "github.com/NVIDIA/carbide-rest-api/carbide-site-workflow/pkg/activity"
	sww "github.com/NVIDIA/carbide-rest-api/carbide-site-workflow/pkg/workflow"
)

// RegisterSubscriber registers the ExpectedMachineWorkflows with the Temporal client
func (api *API) RegisterSubscriber() error {
	// Register the subscribers here
	ManagerAccess.Data.EB.Log.Info().Msg("ExpectedMachine: Registering the subscribers")

	// Register workflows
	// Register CreateExpectedMachine workflow
	ManagerAccess.Data.EB.Managers.Workflow.Temporal.Worker.RegisterWorkflow(sww.CreateExpectedMachine)
	ManagerAccess.Data.EB.Log.Info().Msg("ExpectedMachine: successfully registered the CreateExpectedMachine workflow")

	// Register UpdateExpectedMachine workflow
	ManagerAccess.Data.EB.Managers.Workflow.Temporal.Worker.RegisterWorkflow(sww.UpdateExpectedMachine)
	ManagerAccess.Data.EB.Log.Info().Msg("ExpectedMachine: successfully registered the UpdateExpectedMachine workflow")

	// Register DeleteExpectedMachine workflow
	ManagerAccess.Data.EB.Managers.Workflow.Temporal.Worker.RegisterWorkflow(sww.DeleteExpectedMachine)
	ManagerAccess.Data.EB.Log.Info().Msg("ExpectedMachine: successfully registered the DeleteExpectedMachine workflow")

	// Register activities
	expectedMachineManager := swa.NewManageExpectedMachine(ManagerAccess.Data.EB.Managers.Carbide.Client)

	// Register CreateExpectedMachineOnSite activity
	ManagerAccess.Data.EB.Managers.Workflow.Temporal.Worker.RegisterActivity(expectedMachineManager.CreateExpectedMachineOnSite)
	ManagerAccess.Data.EB.Log.Info().Msg("ExpectedMachine: successfully registered the CreateExpectedMachineOnSite activity")

	// Register UpdateExpectedMachineOnSite activity
	ManagerAccess.Data.EB.Managers.Workflow.Temporal.Worker.RegisterActivity(expectedMachineManager.UpdateExpectedMachineOnSite)
	ManagerAccess.Data.EB.Log.Info().Msg("ExpectedMachine: successfully registered the UpdateExpectedMachineOnSite activity")

	// Register DeleteExpectedMachineOnSite activity
	ManagerAccess.Data.EB.Managers.Workflow.Temporal.Worker.RegisterActivity(expectedMachineManager.DeleteExpectedMachineOnSite)
	ManagerAccess.Data.EB.Log.Info().Msg("ExpectedMachine: successfully registered the DeleteExpectedMachineOnSite activity")

	return nil
}
