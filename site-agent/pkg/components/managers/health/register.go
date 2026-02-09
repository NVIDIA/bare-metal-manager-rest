/*
 * SPDX-FileCopyrightText: Copyright (c) 2026 NVIDIA CORPORATION & AFFILIATES. All rights reserved.
 * SPDX-License-Identifier: Apache-2.0
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

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
