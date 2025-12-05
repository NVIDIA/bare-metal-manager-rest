// SPDX-FileCopyrightText: Copyright (c) 2021-2023 NVIDIA CORPORATION & AFFILIATES. All rights reserved.
// SPDX-License-Identifier: LicenseRef-NvidiaProprietary
//
// NVIDIA CORPORATION, its affiliates and licensors retain all intellectual
// property and proprietary rights in and to this material, related
// documentation and any modifications thereto. Any use, reproduction,
// disclosure or distribution of this material and related documentation
// without an express license agreement from NVIDIA CORPORATION or
// its affiliates is strictly prohibited.

package vpcprefix

import (
	swa "github.com/nvidia/carbide-rest/site-workflow/pkg/activity"
	sww "github.com/nvidia/carbide-rest/site-workflow/pkg/workflow"
)

// RegisterSubscriber registers the VpcPrefixWorkflows with the Temporal client
func (api *API) RegisterSubscriber() error {
	// Register the subscribers here
	ManagerAccess.Data.EB.Log.Info().Msg("VpcPrefix: Registering the subscribers")

	vpcPrefixManager := swa.NewManageVpcPrefix(ManagerAccess.Data.EB.Managers.Carbide.Client)

	//  Register Workflows

	// Sync workflows
	// Register CreateVpcPrefix worfklow
	ManagerAccess.Data.EB.Managers.Workflow.Temporal.Worker.RegisterWorkflow(sww.CreateVpcPrefix)
	ManagerAccess.Data.EB.Log.Info().Msg("VpcPrefix: successfully registered Create VpcPrefix workflow")

	// Register UpdateVpcPrefix worfklow
	ManagerAccess.Data.EB.Managers.Workflow.Temporal.Worker.RegisterWorkflow(sww.UpdateVpcPrefix)
	ManagerAccess.Data.EB.Log.Info().Msg("VpcPrefix: successfully registered Update VpcPrefix workflow")

	// Register DeleteVpcPrefix worfklow
	ManagerAccess.Data.EB.Managers.Workflow.Temporal.Worker.RegisterWorkflow(sww.DeleteVpcPrefix)
	ManagerAccess.Data.EB.Log.Info().Msg("VpcPrefix: successfully registered Delete VpcPrefix workflow")

	// Regsiter Activities

	// Sync workflow activities
	// Register CreateVpcPrefixOnSite
	ManagerAccess.Data.EB.Managers.Workflow.Temporal.Worker.RegisterActivity(vpcPrefixManager.CreateVpcPrefixOnSite)
	ManagerAccess.Data.EB.Log.Info().Msg("VpcPrefix: successfully registered Create VpcPrefix activity")

	// Register UpdateVpcPrefixOnSite
	ManagerAccess.Data.EB.Managers.Workflow.Temporal.Worker.RegisterActivity(vpcPrefixManager.UpdateVpcPrefixOnSite)
	ManagerAccess.Data.EB.Log.Info().Msg("VpcPrefix: successfully registered Update VpcPrefix activity")

	// Register DeleteVpcPrefixOnSite
	ManagerAccess.Data.EB.Managers.Workflow.Temporal.Worker.RegisterActivity(vpcPrefixManager.DeleteVpcPrefixOnSite)
	ManagerAccess.Data.EB.Log.Info().Msg("VpcPrefix: successfully registered Delete VpcPrefix activity")

	return nil
}
