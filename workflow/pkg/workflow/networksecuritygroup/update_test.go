// SPDX-FileCopyrightText: Copyright (c) 2021-2023 NVIDIA CORPORATION & AFFILIATES. All rights reserved.
// SPDX-License-Identifier: LicenseRef-NvidiaProprietary
//
// NVIDIA CORPORATION, its affiliates and licensors retain all intellectual
// property and proprietary rights in and to this material, related
// documentation and any modifications thereto. Any use, reproduction,
// disclosure or distribution of this material and related documentation
// without an express license agreement from NVIDIA CORPORATION or
// its affiliates is strictly prohibited.

package networksecuritygroup

import (
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	networkSecurityGroupActivity "github.com/nvidia/carbide-rest/workflow/pkg/activity/networksecuritygroup"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/testsuite"

	cwssaws "github.com/nvidia/carbide-rest/workflow-schema/schema/site-agent/workflows/v1"
)

type UpdateNetworkSecurityGroupTestSuite struct {
	suite.Suite
	testsuite.WorkflowTestSuite

	env *testsuite.TestWorkflowEnvironment
}

func (s *UpdateNetworkSecurityGroupTestSuite) SetupTest() {
	s.env = s.NewTestWorkflowEnvironment()
}

func (s *UpdateNetworkSecurityGroupTestSuite) AfterTest(suiteName, testName string) {
	s.env.AssertExpectations(s.T())
}

func (s *UpdateNetworkSecurityGroupTestSuite) Test_UpdateNetworkSecurityGroupInventory_Success() {

	var networkSecurityGroupManager networkSecurityGroupActivity.ManageNetworkSecurityGroup

	siteID := uuid.New()
	networkSecurityGroupInventory := &cwssaws.NetworkSecurityGroupInventory{
		NetworkSecurityGroups: []*cwssaws.NetworkSecurityGroup{
			{
				Id: uuid.NewString(),
			},
			{
				Id: uuid.NewString(),
			},
		},
	}

	// Mock UpdateNetworkSecurityGroupViaSiteAgent activity
	s.env.RegisterActivity(networkSecurityGroupManager.UpdateNetworkSecurityGroupsInDB)
	s.env.OnActivity(networkSecurityGroupManager.UpdateNetworkSecurityGroupsInDB, mock.Anything, mock.Anything, mock.Anything).Return(nil)

	// execute UpdateNetworkSecurityGroupInventory workflow
	s.env.ExecuteWorkflow(UpdateNetworkSecurityGroupInventory, siteID.String(), networkSecurityGroupInventory)
	s.True(s.env.IsWorkflowCompleted())
	s.NoError(s.env.GetWorkflowError())
}

func (s *UpdateNetworkSecurityGroupTestSuite) Test_UpdateNetworkSecurityGroupInventory_ActivityFails() {

	var networkSecurityGroupManager networkSecurityGroupActivity.ManageNetworkSecurityGroup

	siteID := uuid.New()
	networkSecurityGroupInventory := &cwssaws.NetworkSecurityGroupInventory{
		NetworkSecurityGroups: []*cwssaws.NetworkSecurityGroup{
			{
				Id: uuid.NewString(),
			},
			{
				Id: uuid.NewString(),
			},
		},
	}

	// Mock UpdateNetworkSecurityGroupsViaSiteAgent activity failure
	s.env.RegisterActivity(networkSecurityGroupManager.UpdateNetworkSecurityGroupsInDB)
	s.env.OnActivity(networkSecurityGroupManager.UpdateNetworkSecurityGroupsInDB, mock.Anything, mock.Anything, mock.Anything).Return(errors.New("UpdateNetworkSecurityGroupInventory Failure"))

	// execute UpdateNetworkSecurityGroupStatus workflow
	s.env.ExecuteWorkflow(UpdateNetworkSecurityGroupInventory, siteID.String(), networkSecurityGroupInventory)
	s.True(s.env.IsWorkflowCompleted())
	err := s.env.GetWorkflowError()
	s.Error(err)

	var applicationErr *temporal.ApplicationError
	s.True(errors.As(err, &applicationErr))
	s.Equal("UpdateNetworkSecurityGroupInventory Failure", applicationErr.Error())
}

func TestUpdateNetworkSecurityGroupSuite(t *testing.T) {
	suite.Run(t, new(UpdateNetworkSecurityGroupTestSuite))
}
