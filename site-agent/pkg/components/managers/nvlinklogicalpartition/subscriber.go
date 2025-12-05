// SPDX-FileCopyrightText: Copyright (c) 2021-2023 NVIDIA CORPORATION & AFFILIATES. All rights reserved.
// SPDX-License-Identifier: LicenseRef-NvidiaProprietary
//
// NVIDIA CORPORATION, its affiliates and licensors retain all intellectual
// property and proprietary rights in and to this material, related
// documentation and any modifications thereto. Any use, reproduction,
// disclosure or distribution of this material and related documentation
// without an express license agreement from NVIDIA CORPORATION or
// its affiliates is strictly prohibited.

package nvlinklogicalpartition

import (
	swa "github.com/nvidia/carbide-rest/site-workflow/pkg/activity"
	sww "github.com/nvidia/carbide-rest/site-workflow/pkg/workflow"
)

// RegisterSubscriber registers the NVLinkLogicalPartitionWorkflows with the Temporal client
func (api *API) RegisterSubscriber() error {
	// Register the subscribers here
	ManagerAccess.Data.EB.Log.Info().Msg("NVLinkLogicalPartition: Registering the subscribers")

	// Register workflows
	// Register CreateNVLinkLogicalPartition workflow
	ManagerAccess.Data.EB.Managers.Workflow.Temporal.Worker.RegisterWorkflow(sww.CreateNVLinkLogicalPartition)
	ManagerAccess.Data.EB.Log.Info().Msg("NVLinkLogicalPartition: successfully registered the CreateNVLinkLogicalPartition workflow")

	// Register UpdateNVLinkLogicalPartition workflow
	ManagerAccess.Data.EB.Managers.Workflow.Temporal.Worker.RegisterWorkflow(sww.UpdateNVLinkLogicalPartition)
	ManagerAccess.Data.EB.Log.Info().Msg("NVLinkLogicalPartition: successfully registered the UpdateNVLinkLogicalPartition workflow")

	// Register DeleteNVLinkLogicalPartition workflow
	ManagerAccess.Data.EB.Managers.Workflow.Temporal.Worker.RegisterWorkflow(sww.DeleteNVLinkLogicalPartition)
	ManagerAccess.Data.EB.Log.Info().Msg("NVLinkLogicalPartition: successfully registered the DeleteNVLinkLogicalPartition workflow")

	// Register activities
	nvlinkLogicalPartitionManager := swa.NewManageNVLinkLogicalPartition(ManagerAccess.Data.EB.Managers.Carbide.Client)

	// Register CreateNVLinkLogicalPartitionOnSite activity
	ManagerAccess.Data.EB.Managers.Workflow.Temporal.Worker.RegisterActivity(nvlinkLogicalPartitionManager.CreateNVLinkLogicalPartitionOnSite)
	ManagerAccess.Data.EB.Log.Info().Msg("NVLinkLogicalPartition: successfully registered the CreateNVLinkLogicalPartitionOnSite activity")

	// Register UpdateNVLinkLogicalPartitionOnSite activity
	ManagerAccess.Data.EB.Managers.Workflow.Temporal.Worker.RegisterActivity(nvlinkLogicalPartitionManager.UpdateNVLinkLogicalPartitionOnSite)
	ManagerAccess.Data.EB.Log.Info().Msg("NVLinkLogicalPartition: successfully registered the UpdateNVLinkLogicalPartitionOnSite activity")

	// Register DeleteNVLinkLogicalPartitionOnSite activity
	ManagerAccess.Data.EB.Managers.Workflow.Temporal.Worker.RegisterActivity(nvlinkLogicalPartitionManager.DeleteNVLinkLogicalPartitionOnSite)
	ManagerAccess.Data.EB.Log.Info().Msg("NVLinkLogicalPartition: successfully registered the DeleteNVLinkLogicalPartitionOnSite activity")

	return nil
}
