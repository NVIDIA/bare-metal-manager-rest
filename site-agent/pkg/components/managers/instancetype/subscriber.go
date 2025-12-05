// SPDX-FileCopyrightText: Copyright (c) 2021-2023 NVIDIA CORPORATION & AFFILIATES. All rights reserved.
// SPDX-License-Identifier: LicenseRef-NvidiaProprietary
//
// NVIDIA CORPORATION, its affiliates and licensors retain all intellectual
// property and proprietary rights in and to this material, related
// documentation and any modifications thereto. Any use, reproduction,
// disclosure or distribution of this material and related documentation
// without an express license agreement from NVIDIA CORPORATION or
// its affiliates is strictly prohibited.

package instancetype

import (
	swa "github.com/nvidia/carbide-rest/site-workflow/pkg/activity"
	sww "github.com/nvidia/carbide-rest/site-workflow/pkg/workflow"
)

// RegisterSubscriber registers the InstanceType workflows and activities with the Temporal client
func (api *API) RegisterSubscriber() error {
	// Register the subscribers here
	ManagerAccess.Data.EB.Log.Info().Msg("InstanceType: Registering the subscribers")

	instanceTypeManager := swa.NewManageInstanceType(ManagerAccess.Data.EB.Managers.Carbide.Client)

	//  Register Workflows

	// Sync workflows
	// Register CreateInstanceType worfklow
	ManagerAccess.Data.EB.Managers.Workflow.Temporal.Worker.RegisterWorkflow(sww.CreateInstanceType)
	ManagerAccess.Data.EB.Log.Info().Msg("InstanceType: successfully registered Create InstanceType workflow")

	// Register UpdateInstanceType worfklow
	ManagerAccess.Data.EB.Managers.Workflow.Temporal.Worker.RegisterWorkflow(sww.UpdateInstanceType)
	ManagerAccess.Data.EB.Log.Info().Msg("InstanceType: successfully registered Update InstanceType workflow")

	// Register DeleteInstanceType worfklow
	ManagerAccess.Data.EB.Managers.Workflow.Temporal.Worker.RegisterWorkflow(sww.DeleteInstanceType)
	ManagerAccess.Data.EB.Log.Info().Msg("InstanceType: successfully registered Delete InstanceType workflow")

	// Register AssociateMachinesWithInstanceType worfklow
	ManagerAccess.Data.EB.Managers.Workflow.Temporal.Worker.RegisterWorkflow(sww.AssociateMachinesWithInstanceType)

	// Register RemoveMachineInstanceTypeAssociation worfklow
	ManagerAccess.Data.EB.Managers.Workflow.Temporal.Worker.RegisterWorkflow(sww.RemoveMachineInstanceTypeAssociation)

	// Regsiter Activities

	// Sync workflow activities
	// Register CreateInstanceTypeOnSite
	ManagerAccess.Data.EB.Managers.Workflow.Temporal.Worker.RegisterActivity(instanceTypeManager.CreateInstanceTypeOnSite)
	ManagerAccess.Data.EB.Log.Info().Msg("InstanceType: successfully registered Create InstanceType activity")

	// Register UpdateInstanceTypeOnSite
	ManagerAccess.Data.EB.Managers.Workflow.Temporal.Worker.RegisterActivity(instanceTypeManager.UpdateInstanceTypeOnSite)
	ManagerAccess.Data.EB.Log.Info().Msg("InstanceType: successfully registered Update InstanceType activity")

	// Register DeleteInstanceTypeOnSite
	ManagerAccess.Data.EB.Managers.Workflow.Temporal.Worker.RegisterActivity(instanceTypeManager.DeleteInstanceTypeOnSite)
	ManagerAccess.Data.EB.Log.Info().Msg("InstanceType: successfully registered Delete InstanceType activity")

	// Register AssociateMachinesWithInstanceTypeOnSite
	ManagerAccess.Data.EB.Managers.Workflow.Temporal.Worker.RegisterActivity(instanceTypeManager.AssociateMachinesWithInstanceTypeOnSite)
	ManagerAccess.Data.EB.Log.Info().Msg("InstanceType: successfully registered AssociateMachinesWithInstanceType activity")

	// Register RemoveMachineInstanceTypeAssociationOnSite
	ManagerAccess.Data.EB.Managers.Workflow.Temporal.Worker.RegisterActivity(instanceTypeManager.RemoveMachineInstanceTypeAssociationOnSite)
	ManagerAccess.Data.EB.Log.Info().Msg("InstanceType: successfully registered RemoveMachineInstanceTypeAssociation activity")

	return nil
}
