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
package workflow

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.temporal.io/sdk/activity"
	"go.temporal.io/sdk/testsuite"

	"github.com/nvidia/bare-metal-manager-rest/rla/internal/inventory/objects/component"
	"github.com/nvidia/bare-metal-manager-rest/rla/internal/inventory/objects/rack"
	"github.com/nvidia/bare-metal-manager-rest/rla/internal/task/operations"
	"github.com/nvidia/bare-metal-manager-rest/rla/internal/task/task"
	"github.com/nvidia/bare-metal-manager-rest/rla/pkg/common/deviceinfo"
	"github.com/nvidia/bare-metal-manager-rest/rla/pkg/common/location"
)

// mockSetFirmwareUpdateTimeWindowForFirmwareControl is a mock activity function for testing
func mockSetFirmwareUpdateTimeWindowForFirmwareControl(
	ctx context.Context,
	req operations.SetFirmwareUpdateTimeWindowRequest,
) error {
	return nil
}

// mockUpdateTaskStatusForFirmwareControl is a mock activity for updating task status
func mockUpdateTaskStatusForFirmwareControl(ctx context.Context, arg *task.TaskStatusUpdate) error {
	return nil
}

// createTestRackForFirmwareControl creates a test rack with components having the given external IDs.
// externalIDs are the external component IDs (e.g., Carbide machine_id) used for activity calls.
func createTestRackForFirmwareControl(externalIDs ...string) *rack.Rack {
	r := rack.New(deviceinfo.DeviceInfo{ID: uuid.New(), Name: "test-rack"}, location.Location{})
	for _, extID := range externalIDs {
		r.AddComponent(component.Component{
			ComponentID: extID, // External ID for Carbide API calls
		})
	}
	return r
}

func TestFirmwareControlWorkflow(t *testing.T) {
	now := time.Now()
	baseInfo := &operations.FirmwareControlTaskInfo{
		Operation: operations.FirmwareOperationUpgrade,
		StartTime: now.Unix(),
		EndTime:   now.Add(time.Hour * 2).Unix(),
	}
	baseReqInfo := task.ExecutionInfo{
		TaskID: uuid.New(),
		Rack:   createTestRackForFirmwareControl("comp1", "comp2"),
	}

	testCases := map[string]struct {
		reqInfo       task.ExecutionInfo
		info          *operations.FirmwareControlTaskInfo
		activityError error
		expectError   bool
	}{
		"success": {
			reqInfo:       baseReqInfo,
			info:          baseInfo,
			activityError: nil,
			expectError:   false,
		},
		"activity fails": {
			reqInfo:       baseReqInfo,
			info:          baseInfo,
			activityError: errors.New("connection timeout"),
			expectError:   true,
		},
		"single machine success": {
			reqInfo: task.ExecutionInfo{
				TaskID: uuid.New(),
				Rack:   createTestRackForFirmwareControl("single-component"),
			},
			info:          baseInfo,
			activityError: nil,
			expectError:   false,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			testSuite := &testsuite.WorkflowTestSuite{}
			env := testSuite.NewTestWorkflowEnvironment()

			env.RegisterActivityWithOptions(mockSetFirmwareUpdateTimeWindowForFirmwareControl, activity.RegisterOptions{
				Name: "SetFirmwareUpdateTimeWindow",
			})
			env.RegisterActivityWithOptions(mockUpdateTaskStatusForFirmwareControl, activity.RegisterOptions{
				Name: "UpdateTaskStatus",
			})

			env.OnActivity(mockSetFirmwareUpdateTimeWindowForFirmwareControl, mock.Anything, mock.Anything).Return(tc.activityError)
			env.OnActivity(mockUpdateTaskStatusForFirmwareControl, mock.Anything, mock.Anything).Return(nil)

			env.ExecuteWorkflow(FirmwareControl, tc.reqInfo, tc.info)

			assert.True(t, env.IsWorkflowCompleted())

			if tc.expectError {
				assert.Error(t, env.GetWorkflowError())
			} else {
				assert.NoError(t, env.GetWorkflowError())
			}
		})
	}
}

func TestFirmwareControlWorkflowEmptyComponents(t *testing.T) {
	testSuite := &testsuite.WorkflowTestSuite{}
	env := testSuite.NewTestWorkflowEnvironment()

	now := time.Now()
	// Create rack with no components
	emptyRack := rack.New(deviceinfo.DeviceInfo{ID: uuid.New(), Name: "empty-rack"}, location.Location{})
	reqInfo := task.ExecutionInfo{
		TaskID: uuid.New(),
		Rack:   emptyRack,
	}
	info := &operations.FirmwareControlTaskInfo{
		Operation: operations.FirmwareOperationUpgrade,
		StartTime: now.Unix(),
		EndTime:   now.Add(time.Hour * 2).Unix(),
	}

	env.ExecuteWorkflow(FirmwareControl, reqInfo, info)

	assert.True(t, env.IsWorkflowCompleted())
	assert.Error(t, env.GetWorkflowError()) // Should error because no components
}

func TestFirmwareControlWorkflowNoComponentIDs(t *testing.T) {
	testSuite := &testsuite.WorkflowTestSuite{}
	env := testSuite.NewTestWorkflowEnvironment()

	now := time.Now()
	// Components without ComponentID
	r := rack.New(deviceinfo.DeviceInfo{ID: uuid.New(), Name: "test-rack"}, location.Location{})
	r.AddComponent(component.Component{}) // Component without ComponentID
	r.AddComponent(component.Component{}) // Component without ComponentID
	reqInfo := task.ExecutionInfo{
		TaskID: uuid.New(),
		Rack:   r,
	}
	info := &operations.FirmwareControlTaskInfo{
		Operation: operations.FirmwareOperationUpgrade,
		StartTime: now.Unix(),
		EndTime:   now.Add(time.Hour * 2).Unix(),
	}

	env.ExecuteWorkflow(FirmwareControl, reqInfo, info)

	assert.True(t, env.IsWorkflowCompleted())
	assert.Error(t, env.GetWorkflowError()) // Should error because no component IDs
}
