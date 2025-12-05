// SPDX-FileCopyrightText: Copyright (c) 2021-2023 NVIDIA CORPORATION & AFFILIATES. All rights reserved.
// SPDX-License-Identifier: LicenseRef-NvidiaProprietary
//
// NVIDIA CORPORATION, its affiliates and licensors retain all intellectual
// property and proprietary rights in and to this material, related
// documentation and any modifications thereto. Any use, reproduction,
// disclosure or distribution of this material and related documentation
// without an express license agreement from NVIDIA CORPORATION or
// its affiliates is strictly prohibited.

package sshkeygroup

import (
	"github.com/google/uuid"
	swa "github.com/nvidia/carbide-rest/site-workflow/pkg/activity"
	sww "github.com/nvidia/carbide-rest/site-workflow/pkg/workflow"
	"go.temporal.io/sdk/activity"
)

// RegisterPublisher registers the SSHKeyGroupWorkflows with the Temporal client
func (api *API) RegisterPublisher() error {
	// Register the publishers here
	ManagerAccess.Data.EB.Log.Info().Msg("SSHKeyGroup: Registering the publishers")

	// Get SSHKeyGroup workflow interface
	SSHKeyGroupinterface := NewSSHKeyGroupWorkflows(
		ManagerAccess.Data.EB.Managers.Workflow.Temporal.Publisher,
		ManagerAccess.Data.EB.Managers.Workflow.Temporal.Subscriber,
		ManagerAccess.Conf.EB,
	)

	activityRegisterOptions := activity.RegisterOptions{
		Name: "PublishSSHKeyGroupActivity",
	}

	ManagerAccess.Data.EB.Managers.Workflow.Temporal.Worker.RegisterActivityWithOptions(
		SSHKeyGroupinterface.PublishSSHKeyGroupActivity, activityRegisterOptions,
	)
	ManagerAccess.Data.EB.Log.Info().Msg("SSHKeyGroup: successfully registered the Publish SSHKeyGroup activity")

	// SSHKeyGroup Inventory workflow
	ManagerAccess.Data.EB.Managers.Workflow.Temporal.Worker.RegisterWorkflow(sww.DiscoverSSHKeyGroupInventory)
	ManagerAccess.Data.EB.Log.Info().Msg("SSHKeyGroup: successfully registered the Discover SSHKeyGroup Inventory workflow")

	inventoryManager := swa.NewManageSSHKeyGroupInventory(swa.ManageInventoryConfig{
		SiteID:                uuid.MustParse(ManagerAccess.Conf.EB.Temporal.ClusterID),
		CarbideAtomicClient:   ManagerAccess.Data.EB.Managers.Carbide.Client,
		TemporalPublishClient: ManagerAccess.Data.EB.Managers.Workflow.Temporal.Publisher,
		TemporalPublishQueue:  ManagerAccess.Conf.EB.Temporal.TemporalPublishQueue,
		SitePageSize:          InventoryCarbidePageSize,
		CloudPageSize:         InventoryCloudPageSize,
	})
	ManagerAccess.Data.EB.Managers.Workflow.Temporal.Worker.RegisterActivity(inventoryManager.DiscoverSSHKeyGroupInventory)
	ManagerAccess.Data.EB.Log.Info().Msg("SSHKeyGroup: successfully registered the Discover SSHKeyGroup Inventory activity")

	api.RegisterCron()

	return nil
}
