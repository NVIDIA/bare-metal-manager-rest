// SPDX-FileCopyrightText: Copyright (c) 2021-2023 NVIDIA CORPORATION & AFFILIATES. All rights reserved.
// SPDX-License-Identifier: LicenseRef-NvidiaProprietary
//
// NVIDIA CORPORATION, its affiliates and licensors retain all intellectual
// property and proprietary rights in and to this material, related
// documentation and any modifications thereto. Any use, reproduction,
// disclosure or distribution of this material and related documentation
// without an express license agreement from NVIDIA CORPORATION or
// its affiliates is strictly prohibited.

package expectedmachine

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

// RegisterCron - Register Cron
func (api *API) RegisterCron() error {
	ManagerAccess.Data.EB.Log.Info().Msg("ExpectedMachine: Registering Inventory Discovery Cron")

	workflowID := "inventory-expected-machine-" + ManagerAccess.Conf.EB.Temporal.TemporalSubscribeNamespace

	cronSchedule := InventoryDefaultSchedule
	if ManagerAccess.Conf.EB.Temporal.TemporalInventorySchedule != "" {
		cronSchedule = ManagerAccess.Conf.EB.Temporal.TemporalInventorySchedule
	}

	ManagerAccess.Data.EB.Log.Info().Str("Schedule", cronSchedule).Msg("ExpectedMachine: Inventory Discovery Cron Schedule")

	workflowOptions := client.StartWorkflowOptions{
		ID:           workflowID,
		TaskQueue:    ManagerAccess.Conf.EB.Temporal.TemporalSubscribeQueue,
		CronSchedule: cronSchedule,
	}

	we, err := ManagerAccess.Data.EB.Managers.Workflow.Temporal.Subscriber.ExecuteWorkflow(
		context.Background(),
		workflowOptions,
		sww.DiscoverExpectedMachineInventory,
	)
	if err != nil {
		ManagerAccess.Data.EB.Log.Error().Err(err).Msg("ExpectedMachine: Error registering Inventory Discovery Cron")
		return err
	}

	ManagerAccess.Data.EB.Log.Info().Interface("Workflow ID", we.GetID()).Msg("ExpectedMachine: successfully registered Inventory Discovery Cron")

	return nil
}
