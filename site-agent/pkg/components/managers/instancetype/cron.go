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
	"context"

	"go.temporal.io/sdk/client"

	sww "github.com/nvidia/carbide-rest/site-workflow/pkg/workflow"
)

const (
	// InventoryCarbidePageSize is the number of items to be fetched from Carbide API at a time
	InventoryCarbidePageSize = 100
	// InventoryCloudPageSize is the number of items to be sent to Cloud at a time
	InventoryCloudPageSize = 25
	// InventoryDefaultSchedule is the default schedule for inventory discovery
	InventoryDefaultSchedule = "@every 3m"
)

// RegisterCron - Register cron workflow
func (api *API) RegisterCron() error {
	// Validate the OS Image config later
	ManagerAccess.Data.EB.Log.Info().Msg("InstanceType: Registering Inventory Discovery Cron")

	workflowID := "inventory-instance-type-" + ManagerAccess.Conf.EB.Temporal.TemporalSubscribeNamespace

	cronSchedule := InventoryDefaultSchedule
	if ManagerAccess.Conf.EB.Temporal.TemporalInventorySchedule != "" {
		cronSchedule = ManagerAccess.Conf.EB.Temporal.TemporalInventorySchedule
	}

	ManagerAccess.Data.EB.Log.Info().Str("Schedule", cronSchedule).Msg("InstanceType: Inventory Discovery Cron Schedule")

	workflowOptions := client.StartWorkflowOptions{
		ID:           workflowID,
		TaskQueue:    ManagerAccess.Conf.EB.Temporal.TemporalSubscribeQueue,
		CronSchedule: cronSchedule,
	}

	we, err := ManagerAccess.Data.EB.Managers.Workflow.Temporal.Subscriber.ExecuteWorkflow(
		context.Background(),
		workflowOptions,
		sww.DiscoverInstanceTypeInventory,
	)

	if err != nil {
		ManagerAccess.Data.EB.Log.Error().Err(err).Msg("InstanceType: Error registering Inventory Collect/Publish Cron")
		return err
	}

	wid := ""
	if !ManagerAccess.Data.EB.Conf.UtMode {
		wid = we.GetID()
	}

	ManagerAccess.Data.EB.Log.Info().Interface("Workflow ID", wid).Msg("InstanceType: successfully registered Inventory Collect/Publish Cron")

	return nil
}
