// SPDX-FileCopyrightText: Copyright (c) 2021-2023 NVIDIA CORPORATION & AFFILIATES. All rights reserved.
// SPDX-License-Identifier: LicenseRef-NvidiaProprietary
//
// NVIDIA CORPORATION, its affiliates and licensors retain all intellectual
// property and proprietary rights in and to this material, related
// documentation and any modifications thereto. Any use, reproduction,
// disclosure or distribution of this material and related documentation
// without an express license agreement from NVIDIA CORPORATION or
// its affiliates is strictly prohibited.

package networksecuritygroup

import (
	swa "github.com/nvidia/carbide-rest/site-workflow/pkg/activity"
	sww "github.com/nvidia/carbide-rest/site-workflow/pkg/workflow"
)

// RegisterSubscriber registers the NetworkSecurityGroup workflows and activities with the Temporal client
func (api *API) RegisterSubscriber() error {
	// Register the subscribers here
	ManagerAccess.Data.EB.Log.Info().Msg("NetworkSecurityGroup: Registering the subscribers")

	networkSecurityGroupManager := swa.NewManageNetworkSecurityGroup(ManagerAccess.Data.EB.Managers.Carbide.Client)

	//  Register Workflows

	// Sync workflows
	// Register CreateNetworkSecurityGroup worfklow
	ManagerAccess.Data.EB.Managers.Workflow.Temporal.Worker.RegisterWorkflow(sww.CreateNetworkSecurityGroup)
	ManagerAccess.Data.EB.Log.Info().Msg("NetworkSecurityGroup: successfully registered Create NetworkSecurityGroup workflow")

	// Register UpdateNetworkSecurityGroup worfklow
	ManagerAccess.Data.EB.Managers.Workflow.Temporal.Worker.RegisterWorkflow(sww.UpdateNetworkSecurityGroup)
	ManagerAccess.Data.EB.Log.Info().Msg("NetworkSecurityGroup: successfully registered Update NetworkSecurityGroup workflow")

	// Register DeleteNetworkSecurityGroup worfklow
	ManagerAccess.Data.EB.Managers.Workflow.Temporal.Worker.RegisterWorkflow(sww.DeleteNetworkSecurityGroup)
	ManagerAccess.Data.EB.Log.Info().Msg("NetworkSecurityGroup: successfully registered Delete NetworkSecurityGroup workflow")

	// Register Activities

	// Sync workflow activities
	// Register CreateNetworkSecurityGroupOnSite
	ManagerAccess.Data.EB.Managers.Workflow.Temporal.Worker.RegisterActivity(networkSecurityGroupManager.CreateNetworkSecurityGroupOnSite)
	ManagerAccess.Data.EB.Log.Info().Msg("NetworkSecurityGroup: successfully registered Create NetworkSecurityGroup activity")

	// Register UpdateNetworkSecurityGroupOnSite
	ManagerAccess.Data.EB.Managers.Workflow.Temporal.Worker.RegisterActivity(networkSecurityGroupManager.UpdateNetworkSecurityGroupOnSite)
	ManagerAccess.Data.EB.Log.Info().Msg("NetworkSecurityGroup: successfully registered Update NetworkSecurityGroup activity")

	// Register DeleteNetworkSecurityGroupOnSite
	ManagerAccess.Data.EB.Managers.Workflow.Temporal.Worker.RegisterActivity(networkSecurityGroupManager.DeleteNetworkSecurityGroupOnSite)
	ManagerAccess.Data.EB.Log.Info().Msg("NetworkSecurityGroup: successfully registered Delete NetworkSecurityGroup activity")

	return nil
}
