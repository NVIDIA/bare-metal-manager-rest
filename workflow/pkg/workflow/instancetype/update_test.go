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
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	instanceTypeActivity "github.com/nvidia/carbide-rest/workflow/pkg/activity/instancetype"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/testsuite"

	cwssaws "github.com/nvidia/carbide-rest/workflow-schema/schema/site-agent/workflows/v1"
)

type UpdateInstanceTypeTestSuite struct {
	suite.Suite
	testsuite.WorkflowTestSuite

	env *testsuite.TestWorkflowEnvironment
}

func (s *UpdateInstanceTypeTestSuite) SetupTest() {
	s.env = s.NewTestWorkflowEnvironment()
}

func (s *UpdateInstanceTypeTestSuite) AfterTest(suiteName, testName string) {
	s.env.AssertExpectations(s.T())
}

func (s *UpdateInstanceTypeTestSuite) Test_UpdateInstanceTypeInventory_Success() {

	var instanceTypeManager instanceTypeActivity.ManageInstanceType

	siteID := uuid.New()
	instanceTypeInventory := &cwssaws.InstanceTypeInventory{
		InstanceTypes: []*cwssaws.InstanceType{
			{
				Id: uuid.NewString(),
			},
			{
				Id: uuid.NewString(),
			},
		},
	}

	// Mock UpdateInstanceTypeViaSiteAgent activity
	s.env.RegisterActivity(instanceTypeManager.UpdateInstanceTypesInDB)
	s.env.OnActivity(instanceTypeManager.UpdateInstanceTypesInDB, mock.Anything, mock.Anything, mock.Anything).Return(nil)

	// execute UpdateInstanceTypeInventory workflow
	s.env.ExecuteWorkflow(UpdateInstanceTypeInventory, siteID.String(), instanceTypeInventory)
	s.True(s.env.IsWorkflowCompleted())
	s.NoError(s.env.GetWorkflowError())
}

func (s *UpdateInstanceTypeTestSuite) Test_UpdateInstanceTypeInventory_ActivityFails() {

	var instanceTypeManager instanceTypeActivity.ManageInstanceType

	siteID := uuid.New()
	instanceTypeInventory := &cwssaws.InstanceTypeInventory{
		InstanceTypes: []*cwssaws.InstanceType{
			{
				Id: uuid.NewString(),
			},
			{
				Id: uuid.NewString(),
			},
		},
	}

	// Mock UpdateInstanceTypesViaSiteAgent activity failure
	s.env.RegisterActivity(instanceTypeManager.UpdateInstanceTypesInDB)
	s.env.OnActivity(instanceTypeManager.UpdateInstanceTypesInDB, mock.Anything, mock.Anything, mock.Anything).Return(errors.New("UpdateInstanceTypeInventory Failure"))

	// execute UpdateInstanceTypeStatus workflow
	s.env.ExecuteWorkflow(UpdateInstanceTypeInventory, siteID.String(), instanceTypeInventory)
	s.True(s.env.IsWorkflowCompleted())
	err := s.env.GetWorkflowError()
	s.Error(err)

	var applicationErr *temporal.ApplicationError
	s.True(errors.As(err, &applicationErr))
	s.Equal("UpdateInstanceTypeInventory Failure", applicationErr.Error())
}

func TestUpdateInstanceTypeSuite(t *testing.T) {
	suite.Run(t, new(UpdateInstanceTypeTestSuite))
}
