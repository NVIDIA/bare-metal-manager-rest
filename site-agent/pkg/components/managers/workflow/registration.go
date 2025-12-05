// SPDX-FileCopyrightText: Copyright (c) 2021-2023 NVIDIA CORPORATION & AFFILIATES. All rights reserved.
// SPDX-License-Identifier: LicenseRef-NvidiaProprietary
//
// NVIDIA CORPORATION, its affiliates and licensors retain all intellectual
// property and proprietary rights in and to this material, related
// documentation and any modifications thereto. Any use, reproduction,
// disclosure or distribution of this material and related documentation
// without an express license agreement from NVIDIA CORPORATION or
// its affiliates is strictly prohibited.

package workflow

import (
	"errors"
	"reflect"

	"github.com/nvidia/carbide-rest/site-agent/pkg/components/utils"
)

// AddWorkflow - Register all the resource workflow functions here
func (w *API) AddWorkflow(wf interface{}) {
	// Register the workflow here
	ManagerAccess.Data.EB.Log.Info().Str("Function", utils.GetFunctionName(wf)).Msg("Workflow: Registering the workflow")
	ManagerAccess.Data.EB.Managers.Workflow.WorkflowFunctions = append(
		ManagerAccess.Data.EB.Managers.Workflow.WorkflowFunctions,
		wf,
	)
}

// Invoke all the resource workflow functions here
func (w *API) Invoke() error {
	// Invoke the workflow here
	ManagerAccess.Data.EB.Log.Info().Msg("Workflow: Invoking the workflow")
	for _, wf := range ManagerAccess.Data.EB.Managers.Workflow.WorkflowFunctions {
		if err := reflect.ValueOf(wf).Call([]reflect.Value{}); err != nil {
			for _, verr := range err {
				// Add the Iszero utility function here later
				if verr.Interface() != reflect.ValueOf(reflect.Type(nil)) {
					ManagerAccess.Data.EB.Log.Error().Str("Function", utils.GetFunctionName(wf)).Msg("Workflow: Failed to invoke the workflow")
					return errors.New("invoke error")
				}
				ManagerAccess.Data.EB.Log.Info().Str("Function", utils.GetFunctionName(wf)).Msg("Workflow: Invoked the workflow")
			}
		}

	}
	return nil
}
