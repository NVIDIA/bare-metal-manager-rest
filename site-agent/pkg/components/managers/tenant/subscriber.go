// SPDX-FileCopyrightText: Copyright (c) 2021-2023 NVIDIA CORPORATION & AFFILIATES. All rights reserved.
// SPDX-License-Identifier: LicenseRef-NvidiaProprietary
//
// NVIDIA CORPORATION, its affiliates and licensors retain all intellectual
// property and proprietary rights in and to this material, related
// documentation and any modifications thereto. Any use, reproduction,
// disclosure or distribution of this material and related documentation
// without an express license agreement from NVIDIA CORPORATION or
// its affiliates is strictly prohibited.

package tenant

import (
	swa "github.com/nvidia/carbide-rest/site-workflow/pkg/activity"
	sww "github.com/nvidia/carbide-rest/site-workflow/pkg/workflow"
)

// RegisterSubscriber registers Tenant workflows Site Agent subscribes to execute
func (api *API) RegisterSubscriber() error {
	//  Register Workflows

	// Sync workflows
	// Register CreateTenant worfklow
	ManagerAccess.Data.EB.Managers.Workflow.Temporal.Worker.RegisterWorkflow(sww.CreateTenant)
	ManagerAccess.Data.EB.Log.Info().Msg("Tenant: successfully registered CreateTenant workflow")

	// Register UpdateTenant worfklow
	ManagerAccess.Data.EB.Managers.Workflow.Temporal.Worker.RegisterWorkflow(sww.UpdateTenant)
	ManagerAccess.Data.EB.Log.Info().Msg("Tenant: successfully registered UpdateTenant workflow")

	// Regsiter Activities
	tenantManager := swa.NewManageTenant(ManagerAccess.Data.EB.Managers.Carbide.Client)

	// Sync workflow activities
	// Register CreateTenantOnSite
	ManagerAccess.Data.EB.Managers.Workflow.Temporal.Worker.RegisterActivity(tenantManager.CreateTenantOnSite)
	ManagerAccess.Data.EB.Log.Info().Msg("Tenant: successfully registered CreateTenantOnSite activity")

	// Register UpdateTenantOnSite
	ManagerAccess.Data.EB.Managers.Workflow.Temporal.Worker.RegisterActivity(tenantManager.UpdateTenantOnSite)
	ManagerAccess.Data.EB.Log.Info().Msg("Tenant: successfully registered UpdateTenantOnSite activity")

	return nil
}

// RegisterSubscribers - this is method 2 of registering the subscriber
func RegisterSubscribers() {
	// Register the subscribers here
	ManagerAccess.Data.EB.Log.Info().Msg("Subnet: Registering the subscribers")
	ManagerAccess.API.Orchestrator.AddWorkflow(ManagerAccess.API.Subnet.CreateSubnet)
}
