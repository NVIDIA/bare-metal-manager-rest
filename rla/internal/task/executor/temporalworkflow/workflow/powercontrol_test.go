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

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.temporal.io/sdk/activity"
	"go.temporal.io/sdk/testsuite"

	"github.com/nvidia/carbide-rest/rla/internal/inventory/objects/component"
	"github.com/nvidia/carbide-rest/rla/internal/inventory/objects/rack"
	taskcommon "github.com/nvidia/carbide-rest/rla/internal/task/common"
	"github.com/nvidia/carbide-rest/rla/internal/task/operations"
	taskdef "github.com/nvidia/carbide-rest/rla/internal/task/task"
	"github.com/nvidia/carbide-rest/rla/pkg/common/deviceinfo"
	"github.com/nvidia/carbide-rest/rla/pkg/common/devicetypes"
	"github.com/nvidia/carbide-rest/rla/pkg/common/location"
)

// mockPowerControl is a mock activity function for testing
func mockPowerControl(
	ctx context.Context,
	info taskcommon.ComponentInfo,
	pcInfo *operations.PowerControlTaskInfo,
) error {
	return nil
}

func mockUpdateTaskStatus(ctx context.Context, arg *taskdef.TaskStatusUpdate) error {
	return nil
}

// Helper function for creating test components
// - id: internal UUID (RLA database primary key)
// - name: human-readable name for DeviceInfo.Name
// - externalID: external ID for ComponentID (e.g., Carbide machine_id)
// - compType: component type
func newTestComponent(id uuid.UUID, name string, externalID string, compType devicetypes.ComponentType) *component.Component {
	return &component.Component{
		Type:        compType,
		ComponentID: externalID,
		Info: deviceinfo.DeviceInfo{
			ID:   id,
			Name: name,
		},
	}
}

// Helper function to build a rack from components for testing
func buildTestRack(components []*component.Component) *rack.Rack {
	r := rack.New(deviceinfo.DeviceInfo{ID: uuid.New(), Name: "test-rack"}, location.Location{})
	for _, c := range components {
		r.AddComponent(*c)
	}
	return r
}

func TestPowerControlWorkflow(t *testing.T) {
	computeID1 := uuid.New()
	computeID2 := uuid.New()
	powershelfID := uuid.New()
	nvlswitchID := uuid.New()

	// Full set of components (PowerShelf, NVLSwitch, Compute)
	fullComponents := []*component.Component{
		newTestComponent(powershelfID, "powershelf-1", "ext-powershelf-1", devicetypes.ComponentTypePowerShelf),
		newTestComponent(nvlswitchID, "nvlswitch-1", "ext-nvlswitch-1", devicetypes.ComponentTypeNVLSwitch),
		newTestComponent(computeID1, "compute-1", "ext-compute-1", devicetypes.ComponentTypeCompute),
		newTestComponent(computeID2, "compute-2", "ext-compute-2", devicetypes.ComponentTypeCompute),
	}

	// Compute-only components
	computeOnlyComponents := []*component.Component{
		newTestComponent(computeID1, "compute-1", "ext-compute-1", devicetypes.ComponentTypeCompute),
		newTestComponent(computeID2, "compute-2", "ext-compute-2", devicetypes.ComponentTypeCompute),
	}

	testCases := map[string]struct {
		components    []*component.Component
		op            operations.PowerOperation
		activityError error
		expectError   bool
		errorContains string
	}{
		"power on full components success": {
			components:    fullComponents,
			op:            operations.PowerOperationPowerOn,
			activityError: nil,
			expectError:   false,
		},
		"power off full components success": {
			components:    fullComponents,
			op:            operations.PowerOperationPowerOff,
			activityError: nil,
			expectError:   false,
		},
		"force power on success": {
			components:    fullComponents,
			op:            operations.PowerOperationForcePowerOn,
			activityError: nil,
			expectError:   false,
		},
		"force power off success": {
			components:    fullComponents,
			op:            operations.PowerOperationForcePowerOff,
			activityError: nil,
			expectError:   false,
		},
		"power on compute only": {
			components:    computeOnlyComponents,
			op:            operations.PowerOperationPowerOn,
			activityError: nil,
			expectError:   false,
		},
		"empty components returns nil": {
			components:    nil,
			op:            operations.PowerOperationPowerOn,
			activityError: nil,
			expectError:   false,
		},
		"restart success": {
			components:    computeOnlyComponents,
			op:            operations.PowerOperationRestart,
			activityError: nil,
			expectError:   false,
		},
		"force restart success": {
			components:    computeOnlyComponents,
			op:            operations.PowerOperationForceRestart,
			activityError: nil,
			expectError:   false,
		},
		"warm reset not supported": {
			components:    fullComponents,
			op:            operations.PowerOperationWarmReset,
			activityError: nil,
			expectError:   true,
			errorContains: "hardware reset is not supported yet",
		},
		"cold reset not supported": {
			components:    fullComponents,
			op:            operations.PowerOperationColdReset,
			activityError: nil,
			expectError:   true,
			errorContains: "hardware reset is not supported yet",
		},
		"unknown operation returns error": {
			components:    fullComponents,
			op:            operations.PowerOperationUnknown,
			activityError: nil,
			expectError:   true,
			errorContains: "unknown power operation",
		},
		"activity failure returns error": {
			components:    computeOnlyComponents,
			op:            operations.PowerOperationPowerOn,
			activityError: errors.New("BMC connection failed"),
			expectError:   true,
			errorContains: "BMC connection failed",
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			testSuite := &testsuite.WorkflowTestSuite{}
			env := testSuite.NewTestWorkflowEnvironment()

			env.RegisterActivityWithOptions(mockPowerControl, activity.RegisterOptions{
				Name: "PowerControl",
			})
			env.RegisterActivityWithOptions(mockUpdateTaskStatus, activity.RegisterOptions{
				Name: "UpdateTaskStatus",
			})

			env.OnActivity(mockPowerControl, mock.Anything, mock.Anything, mock.Anything).Return(tc.activityError)
			env.OnActivity(mockUpdateTaskStatus, mock.Anything, mock.Anything).Return(nil)

			info := operations.PowerControlTaskInfo{Operation: tc.op}
			reqInfo := taskdef.ExecutionInfo{
				TaskID: uuid.New(),
				Rack:   buildTestRack(tc.components),
			}

			env.ExecuteWorkflow(PowerControl, reqInfo, info)

			assert.True(t, env.IsWorkflowCompleted())

			wfErr := env.GetWorkflowError()
			if tc.expectError {
				assert.Error(t, wfErr)
				if tc.errorContains != "" && wfErr != nil {
					assert.Contains(t, wfErr.Error(), tc.errorContains)
				}
			} else {
				assert.NoError(t, wfErr)
			}
		})
	}
}
