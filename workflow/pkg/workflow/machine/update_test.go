// SPDX-FileCopyrightText: Copyright (c) 2021-2023 NVIDIA CORPORATION & AFFILIATES. All rights reserved.
// SPDX-License-Identifier: LicenseRef-NvidiaProprietary
//
// NVIDIA CORPORATION, its affiliates and licensors retain all intellectual
// property and proprietary rights in and to this material, related
// documentation and any modifications thereto. Any use, reproduction,
// disclosure or distribution of this material and related documentation
// without an express license agreement from NVIDIA CORPORATION or
// its affiliates is strictly prohibited.

package machine

import (
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/testsuite"
	"google.golang.org/protobuf/types/known/timestamppb"

	cwssaws "github.com/nvidia/carbide-rest/workflow-schema/schema/site-agent/workflows/v1"

	machineActivity "github.com/nvidia/carbide-rest/workflow/pkg/activity/machine"
)

type UpdateMachineInventoryTestSuite struct {
	suite.Suite
	testsuite.WorkflowTestSuite

	env *testsuite.TestWorkflowEnvironment
}

func (s *UpdateMachineInventoryTestSuite) SetupTest() {
	s.env = s.NewTestWorkflowEnvironment()
}

func (s *UpdateMachineInventoryTestSuite) AfterTest(suiteName, testName string) {
	s.env.AssertExpectations(s.T())
}

func (s *UpdateMachineInventoryTestSuite) Test_UpdateMachineInventory_Success() {
	var machineManager machineActivity.ManageMachine

	siteID := uuid.New()

	machineInfo := &cwssaws.MachineInfo{
		Machine: &cwssaws.Machine{
			Id:    &cwssaws.MachineId{Id: uuid.New().String()},
			State: "Running",
		},
	}

	machineInventory := &cwssaws.MachineInventory{
		Machines:  []*cwssaws.MachineInfo{machineInfo},
		Timestamp: timestamppb.Now(),
	}

	// Mock UpdateVpcViaSiteAgent activity
	s.env.RegisterActivity(machineManager.UpdateMachinesInDB)
	s.env.OnActivity(machineManager.UpdateMachinesInDB, mock.Anything, mock.Anything, mock.Anything).Return(nil)

	// execute UpdateMachineInventory workflow
	s.env.ExecuteWorkflow(UpdateMachineInventory, siteID.String(), machineInventory)
	s.True(s.env.IsWorkflowCompleted())
	s.NoError(s.env.GetWorkflowError())
}

func (s *UpdateMachineInventoryTestSuite) Test_UpdateMachineInventory_ActivityFails() {
	var machineManager machineActivity.ManageMachine

	siteID := uuid.New()

	machineInfo := &cwssaws.MachineInfo{
		Machine: &cwssaws.Machine{
			Id:    &cwssaws.MachineId{Id: uuid.New().String()},
			State: "Running",
		},
	}

	machineInventory := &cwssaws.MachineInventory{
		Machines:  []*cwssaws.MachineInfo{machineInfo},
		Timestamp: timestamppb.Now(),
	}

	// Mock UpdateMachinesInDB activity failure
	s.env.RegisterActivity(machineManager.UpdateMachinesInDB)
	s.env.OnActivity(machineManager.UpdateMachinesInDB, mock.Anything, mock.Anything, mock.Anything).Return(errors.New("UpdateMachineInventory Failure"))

	// execute UpdateMachineInventory workflow
	s.env.ExecuteWorkflow(UpdateMachineInventory, siteID.String(), machineInventory)
	s.True(s.env.IsWorkflowCompleted())
	err := s.env.GetWorkflowError()
	s.Error(err)

	var applicationErr *temporal.ApplicationError
	s.True(errors.As(err, &applicationErr))
	s.Equal("UpdateMachineInventory Failure", applicationErr.Error())
}

func TestUpdateMachineInventorySuite(t *testing.T) {
	suite.Run(t, new(UpdateMachineInventoryTestSuite))
}
