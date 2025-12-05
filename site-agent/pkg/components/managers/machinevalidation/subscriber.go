// SPDX-FileCopyrightText: Copyright (c) 2021-2023 NVIDIA CORPORATION & AFFILIATES. All rights reserved.
// SPDX-License-Identifier: LicenseRef-NvidiaProprietary
//
// NVIDIA CORPORATION, its affiliates and licensors retain all intellectual
// property and proprietary rights in and to this material, related
// documentation and any modifications thereto. Any use, reproduction,
// disclosure or distribution of this material and related documentation
// without an express license agreement from NVIDIA CORPORATION or
// its affiliates is strictly prohibited.

package machinevalidation

import (
	swa "github.com/nvidia/carbide-rest/site-workflow/pkg/activity"
	sww "github.com/nvidia/carbide-rest/site-workflow/pkg/workflow"
)

// RegisterSubscriber registers the MachineValidation workflows/activities with the Temporal client
func (MachineValidation *API) RegisterSubscriber() error {
	// Register subscriber workflows
	ManagerAccess.Data.EB.Log.Info().Msg("MachineValidation: Registering the subscribers")

	manager := swa.NewManageMachineValidation(ManagerAccess.Data.EB.Managers.Carbide.Client)

	// Register workflows

	// Sync workflows
	// Register EnableDisableMachineValidationTest workflow
	ManagerAccess.Data.EB.Managers.Workflow.Temporal.Worker.RegisterWorkflow(sww.EnableDisableMachineValidationTest)
	ManagerAccess.Data.EB.Log.Info().Msg("MachineValidation: successfully registered the EnableDisableMachineValidationTest workflow")

	// Register PersistValidationResult workflow
	ManagerAccess.Data.EB.Managers.Workflow.Temporal.Worker.RegisterWorkflow(sww.PersistValidationResult)
	ManagerAccess.Data.EB.Log.Info().Msg("MachineValidation: successfully registered the PersistValidationResult workflow")

	// Register GetMachineValidationResults workflow
	ManagerAccess.Data.EB.Managers.Workflow.Temporal.Worker.RegisterWorkflow(sww.GetMachineValidationResults)
	ManagerAccess.Data.EB.Log.Info().Msg("MachineValidation: successfully registered the GetMachineValidationResults workflow")

	// Register GetMachineValidationRuns workflow
	ManagerAccess.Data.EB.Managers.Workflow.Temporal.Worker.RegisterWorkflow(sww.GetMachineValidationRuns)
	ManagerAccess.Data.EB.Log.Info().Msg("MachineValidation: successfully registered the GetMachineValidationRuns workflow")

	// Register GetMachineValidationTests workflow
	ManagerAccess.Data.EB.Managers.Workflow.Temporal.Worker.RegisterWorkflow(sww.GetMachineValidationTests)
	ManagerAccess.Data.EB.Log.Info().Msg("MachineValidation: successfully registered the GetMachineValidationTests workflow")

	// Register AddMachineValidationTest workflow
	ManagerAccess.Data.EB.Managers.Workflow.Temporal.Worker.RegisterWorkflow(sww.AddMachineValidationTest)
	ManagerAccess.Data.EB.Log.Info().Msg("MachineValidation: successfully registered the AddMachineValidationTest workflow")

	// Register UpdateMachineValidationTest workflow
	ManagerAccess.Data.EB.Managers.Workflow.Temporal.Worker.RegisterWorkflow(sww.UpdateMachineValidationTest)
	ManagerAccess.Data.EB.Log.Info().Msg("MachineValidation: successfully registered the UpdateMachineValidationTest workflow")

	// Register GetMachineValidationExternalConfigs workflow
	ManagerAccess.Data.EB.Managers.Workflow.Temporal.Worker.RegisterWorkflow(sww.GetMachineValidationExternalConfigs)
	ManagerAccess.Data.EB.Log.Info().Msg("MachineValidation: successfully registered the GetMachineValidationExternalConfigs workflow")

	// Register AddUpdateMachineValidationExternalConfig workflow
	ManagerAccess.Data.EB.Managers.Workflow.Temporal.Worker.RegisterWorkflow(sww.AddUpdateMachineValidationExternalConfig)
	ManagerAccess.Data.EB.Log.Info().Msg("MachineValidation: successfully registered the AddUpdateMachineValidationExternalConfig workflow")

	// Register RemoveMachineValidationExternalConfig workflow
	ManagerAccess.Data.EB.Managers.Workflow.Temporal.Worker.RegisterWorkflow(sww.RemoveMachineValidationExternalConfig)
	ManagerAccess.Data.EB.Log.Info().Msg("MachineValidation: successfully registered the RemoveMachineValidationExternalConfig workflow")

	// Register activities

	// Sync workflow activities
	// Register EnableDisableMachineValidationTestOnSite
	ManagerAccess.Data.EB.Managers.Workflow.Temporal.Worker.RegisterActivity(manager.EnableDisableMachineValidationTestOnSite)
	ManagerAccess.Data.EB.Log.Info().Msg("MachineValidation: successfully registered EnableDisableMachineValidationTestOnSite activity")

	// Register PersistValidationResultOnSite
	ManagerAccess.Data.EB.Managers.Workflow.Temporal.Worker.RegisterActivity(manager.PersistValidationResultOnSite)
	ManagerAccess.Data.EB.Log.Info().Msg("MachineValidation: successfully registered PersistValidationResultOnSite activity")

	// Register GetMachineValidationResultsFromSite
	ManagerAccess.Data.EB.Managers.Workflow.Temporal.Worker.RegisterActivity(manager.GetMachineValidationResultsFromSite)
	ManagerAccess.Data.EB.Log.Info().Msg("MachineValidation: successfully registered GetMachineValidationResultsFromSite activity")

	// Register GetMachineValidationRunsFromSite
	ManagerAccess.Data.EB.Managers.Workflow.Temporal.Worker.RegisterActivity(manager.GetMachineValidationRunsFromSite)
	ManagerAccess.Data.EB.Log.Info().Msg("MachineValidation: successfully registered GetMachineValidationRunsFromSite activity")

	// Register GetMachineValidationTestsFromSite
	ManagerAccess.Data.EB.Managers.Workflow.Temporal.Worker.RegisterActivity(manager.GetMachineValidationTestsFromSite)
	ManagerAccess.Data.EB.Log.Info().Msg("MachineValidation: successfully registered GetMachineValidationTestsFromSite activity")

	// Register AddMachineValidationTestOnSite
	ManagerAccess.Data.EB.Managers.Workflow.Temporal.Worker.RegisterActivity(manager.AddMachineValidationTestOnSite)
	ManagerAccess.Data.EB.Log.Info().Msg("MachineValidation: successfully registered AddMachineValidationTestOnSite activity")

	// Register UpdateMachineValidationTestOnSite
	ManagerAccess.Data.EB.Managers.Workflow.Temporal.Worker.RegisterActivity(manager.UpdateMachineValidationTestOnSite)
	ManagerAccess.Data.EB.Log.Info().Msg("MachineValidation: successfully registered UpdateMachineValidationTestOnSite activity")

	// Register GetMachineValidationExternalConfigsFromSite
	ManagerAccess.Data.EB.Managers.Workflow.Temporal.Worker.RegisterActivity(manager.GetMachineValidationExternalConfigsFromSite)
	ManagerAccess.Data.EB.Log.Info().Msg("MachineValidation: successfully registered GetMachineValidationExternalConfigsFromSite activity")

	// Register AddUpdateMachineValidationExternalConfigOnSite
	ManagerAccess.Data.EB.Managers.Workflow.Temporal.Worker.RegisterActivity(manager.AddUpdateMachineValidationExternalConfigOnSite)
	ManagerAccess.Data.EB.Log.Info().Msg("MachineValidation: successfully registered AddUpdateMachineValidationExternalConfigOnSite activity")

	// Register RemoveMachineValidationExternalConfigOnSite
	ManagerAccess.Data.EB.Managers.Workflow.Temporal.Worker.RegisterActivity(manager.RemoveMachineValidationExternalConfigOnSite)
	ManagerAccess.Data.EB.Log.Info().Msg("MachineValidation: successfully registered RemoveMachineValidationExternalConfigOnSite activity")

	return nil
}
