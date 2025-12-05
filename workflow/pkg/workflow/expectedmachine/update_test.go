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
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	expectedMachineActivity "github.com/nvidia/carbide-rest/workflow/pkg/activity/expectedmachine"

	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/testsuite"

	cwssaws "github.com/nvidia/carbide-rest/workflow-schema/schema/site-agent/workflows/v1"
)

type UpdateExpectedMachineTestSuite struct {
	suite.Suite
	testsuite.WorkflowTestSuite

	env *testsuite.TestWorkflowEnvironment
}

func (s *UpdateExpectedMachineTestSuite) SetupTest() {
	s.env = s.NewTestWorkflowEnvironment()
}

func (s *UpdateExpectedMachineTestSuite) AfterTest(suiteName, testName string) {
	s.env.AssertExpectations(s.T())
}

func (s *UpdateExpectedMachineTestSuite) Test_UpdateExpectedMachineInventory_Success() {
	var expectedMachineManager expectedMachineActivity.ManageExpectedMachine

	siteID := uuid.New()

	inv := &cwssaws.ExpectedMachineInventory{
		ExpectedMachines: []*cwssaws.ExpectedMachine{},
	}

	s.env.RegisterActivity(expectedMachineManager.UpdateExpectedMachinesInDB)
	s.env.OnActivity(expectedMachineManager.UpdateExpectedMachinesInDB, mock.Anything, mock.Anything, mock.Anything).Return(nil)

	s.env.ExecuteWorkflow(UpdateExpectedMachineInventory, siteID.String(), inv)
	s.True(s.env.IsWorkflowCompleted())
	s.NoError(s.env.GetWorkflowError())
}

func (s *UpdateExpectedMachineTestSuite) Test_UpdateExpectedMachineInventory_ActivityFails() {
	var expectedMachineManager expectedMachineActivity.ManageExpectedMachine

	siteID := uuid.New()

	inv := &cwssaws.ExpectedMachineInventory{
		ExpectedMachines: []*cwssaws.ExpectedMachine{},
	}

	s.env.RegisterActivity(expectedMachineManager.UpdateExpectedMachinesInDB)
	s.env.OnActivity(expectedMachineManager.UpdateExpectedMachinesInDB, mock.Anything, mock.Anything, mock.Anything).Return(errors.New("UpdateExpectedMachineInventory Failure"))

	s.env.ExecuteWorkflow(UpdateExpectedMachineInventory, siteID.String(), inv)
	s.True(s.env.IsWorkflowCompleted())
	err := s.env.GetWorkflowError()
	s.NotNil(err)

	var applicationErr *temporal.ApplicationError
	s.True(errors.As(err, &applicationErr))
	s.Equal("UpdateExpectedMachineInventory Failure", applicationErr.Error())
}

func TestUpdateExpectedMachineTestSuite(t *testing.T) {
	suite.Run(t, new(UpdateExpectedMachineTestSuite))
}
