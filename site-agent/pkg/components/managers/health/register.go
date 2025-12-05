// SPDX-FileCopyrightText: Copyright (c) 2021-2023 NVIDIA CORPORATION & AFFILIATES. All rights reserved.
// SPDX-License-Identifier: LicenseRef-NvidiaProprietary
//
// NVIDIA CORPORATION, its affiliates and licensors retain all intellectual
// property and proprietary rights in and to this material, related
// documentation and any modifications thereto. Any use, reproduction,
// disclosure or distribution of this material and related documentation
// without an express license agreement from NVIDIA CORPORATION or
// its affiliates is strictly prohibited.

package health

import (
	"go.temporal.io/sdk/activity"
	workflow "go.temporal.io/sdk/workflow"
)

// RegisterSubscriber registers the HealthWorkflows with the Temporal client
func (api *API) RegisterSubscriber() error {

	// Register the subscribers here
	ManagerAccess.Data.EB.Log.Info().Msg("Health: Registering the subscribers")

	// Get Health workflow interface
	Healthinterface := NewHealthWorkflows(
		ManagerAccess.Data.EB.Managers.Workflow.Temporal.Publisher,
		ManagerAccess.Data.EB.Managers.Workflow.Temporal.Subscriber,
		ManagerAccess.Conf.EB,
	)

	// Register worfklow
	wflowRegisterOptions := workflow.RegisterOptions{
		Name: "GetHealth",
	}
	ManagerAccess.Data.EB.Managers.Workflow.Temporal.Worker.RegisterWorkflowWithOptions(
		GetHealth, wflowRegisterOptions,
	)

	ManagerAccess.Data.EB.Log.Info().Msg("Health: successfully registered the get Health workflow")

	// Register activity
	activityRegisterOptions := activity.RegisterOptions{
		Name: "CreateHealthActivity",
	}

	ManagerAccess.Data.EB.Managers.Workflow.Temporal.Worker.RegisterActivityWithOptions(
		Healthinterface.GetHealthActivity, activityRegisterOptions,
	)
	ManagerAccess.Data.EB.Log.Info().Msg("Health: successfully registered the get Health activity")

	return nil
}
