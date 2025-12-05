// SPDX-FileCopyrightText: Copyright (c) 2021-2023 NVIDIA CORPORATION & AFFILIATES. All rights reserved.
// SPDX-License-Identifier: LicenseRef-NvidiaProprietary
//
// NVIDIA CORPORATION, its affiliates and licensors retain all intellectual
// property and proprietary rights in and to this material, related
// documentation and any modifications thereto. Any use, reproduction,
// disclosure or distribution of this material and related documentation
// without an express license agreement from NVIDIA CORPORATION or
// its affiliates is strictly prohibited.

package vpc

import (
	swa "github.com/nvidia/carbide-rest/site-workflow/pkg/activity"
	sww "github.com/nvidia/carbide-rest/site-workflow/pkg/workflow"
)

// RegisterSubscriber registers the VPCWorkflows with the Temporal client
func (api *API) RegisterSubscriber() error {

	vpcManager := swa.NewManageVPC(ManagerAccess.Data.EB.Managers.Carbide.Client)

	// Register the subscribers here
	ManagerAccess.Data.EB.Log.Info().Msg("VPC: Registering the subscribers")

	// Get vpc workflow interface
	vpcinterface := NewVPCWorkflows(
		ManagerAccess.Data.EB.Managers.Workflow.Temporal.Publisher,
		ManagerAccess.Data.EB.Managers.Workflow.Temporal.Subscriber,
		ManagerAccess.Conf.EB,
	)

	/// Register worfklows

	// Sync workflows

	// CreateVPCV2
	ManagerAccess.Data.EB.Managers.Workflow.Temporal.Worker.RegisterWorkflow(sww.CreateVPCV2)
	ManagerAccess.Data.EB.Log.Info().Msg("VPC: successfully registered CreateVPC v2 workflow")

	// UpdateVPC
	ManagerAccess.Data.EB.Managers.Workflow.Temporal.Worker.RegisterWorkflow(sww.UpdateVPC)
	ManagerAccess.Data.EB.Log.Info().Msg("VPC: successfully registered UpdateVPC workflow")

	// UpdateVPCVirtualization
	ManagerAccess.Data.EB.Managers.Workflow.Temporal.Worker.RegisterWorkflow(sww.UpdateVPCVirtualization)
	ManagerAccess.Data.EB.Log.Info().Msg("VPC: successfully UpdateVPCVirtualization workflow")

	// DeleteVPCV2
	ManagerAccess.Data.EB.Managers.Workflow.Temporal.Worker.RegisterWorkflow(sww.DeleteVPCV2)
	ManagerAccess.Data.EB.Log.Info().Msg("VPC: successfully registered DeleteVPC v2 workflow")

	/// Legacy workflows

	// CreateVPC
	ManagerAccess.Data.EB.Managers.Workflow.Temporal.Worker.RegisterWorkflow(ManagerAccess.API.VPC.CreateVPC)
	ManagerAccess.Data.EB.Log.Info().Msg("VPC: successfully registered the create VPC workflow")

	// DeleteVPC
	ManagerAccess.Data.EB.Managers.Workflow.Temporal.Worker.RegisterWorkflow(ManagerAccess.API.VPC.DeleteVPC)
	ManagerAccess.Data.EB.Log.Info().Msg("VPC: successfully registered the delete VPC workflow")

	// GetVPCByName
	ManagerAccess.Data.EB.Managers.Workflow.Temporal.Worker.RegisterWorkflow(ManagerAccess.API.VPC.GetVPCByName)
	ManagerAccess.Data.EB.Log.Info().Msg("VPC: successfully registered the GetVPCByName VPC workflow")

	/// Register activities

	// Sync workflow activities

	// CreateVpcOnSite
	ManagerAccess.Data.EB.Managers.Workflow.Temporal.Worker.RegisterActivity(vpcManager.CreateVpcOnSite)
	ManagerAccess.Data.EB.Log.Info().Msg("VPC: successfully registered the CreateVpcOnSite activity")

	// UpdateVpcOnSite
	ManagerAccess.Data.EB.Managers.Workflow.Temporal.Worker.RegisterActivity(vpcManager.UpdateVpcOnSite)
	ManagerAccess.Data.EB.Log.Info().Msg("VPC: successfully registered the UpdateVpcOnSite activity")

	// UpdateVpcVirtualizationOnSite
	ManagerAccess.Data.EB.Managers.Workflow.Temporal.Worker.RegisterActivity(vpcManager.UpdateVpcVirtualizationOnSite)
	ManagerAccess.Data.EB.Log.Info().Msg("VPC: successfully registered the UpdateVpcVirtualizationOnSite activity")

	// UpdateVpcOnSite
	ManagerAccess.Data.EB.Managers.Workflow.Temporal.Worker.RegisterActivity(vpcManager.DeleteVpcOnSite)
	ManagerAccess.Data.EB.Log.Info().Msg("VPC: successfully registered the DeleteVpcOnSite activity")

	// Legacy workflow activities

	// CreateVPCActivity
	ManagerAccess.Data.EB.Managers.Workflow.Temporal.Worker.RegisterActivity(vpcinterface.CreateVPCActivity)
	ManagerAccess.Data.EB.Log.Info().Msg("VPC: successfully registered the Create VPC activity")

	// UpdateVPCActivity
	ManagerAccess.Data.EB.Managers.Workflow.Temporal.Worker.RegisterActivity(vpcinterface.UpdateVPCActivity)
	ManagerAccess.Data.EB.Log.Info().Msg("VPC: successfully registered the Update VPC activity")

	// DeleteVPCActivity
	ManagerAccess.Data.EB.Managers.Workflow.Temporal.Worker.RegisterActivity(vpcinterface.DeleteVPCActivity)
	ManagerAccess.Data.EB.Log.Info().Msg("VPC: successfully registered the Delete VPC activity")

	// GetVPCByNameActivity
	ManagerAccess.Data.EB.Managers.Workflow.Temporal.Worker.RegisterActivity(vpcinterface.GetVPCByNameActivity)
	ManagerAccess.Data.EB.Log.Info().Msg("VPC: successfully registered the GetVPCByName VPC activity")

	return nil
}

// RegisterSubscribers - this is method 2 of registering the subscriber
func RegisterSubscribers() {
	// Register the subscribers here
	ManagerAccess.Data.EB.Log.Info().Msg("VPC: Registering the subscribers")
	ManagerAccess.API.Orchestrator.AddWorkflow(ManagerAccess.API.VPC.CreateVPC)
}
