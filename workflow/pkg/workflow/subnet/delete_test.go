// SPDX-FileCopyrightText: Copyright (c) 2021-2023 NVIDIA CORPORATION & AFFILIATES. All rights reserved.
// SPDX-License-Identifier: LicenseRef-NvidiaProprietary
//
// NVIDIA CORPORATION, its affiliates and licensors retain all intellectual
// property and proprietary rights in and to this material, related
// documentation and any modifications thereto. Any use, reproduction,
// disclosure or distribution of this material and related documentation
// without an express license agreement from NVIDIA CORPORATION or
// its affiliates is strictly prohibited.

package subnet

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/testsuite"

	tmocks "go.temporal.io/sdk/mocks"

	subnetActivity "github.com/nvidia/carbide-rest/workflow/pkg/activity/subnet"
)

type DeleteSubnetTestSuite struct {
	suite.Suite
	testsuite.WorkflowTestSuite

	env *testsuite.TestWorkflowEnvironment
}

func (s *DeleteSubnetTestSuite) SetupTest() {
	s.env = s.NewTestWorkflowEnvironment()
}

func (s *DeleteSubnetTestSuite) AfterTest(suiteName, testName string) {
	s.env.AssertExpectations(s.T())
}

func (s *DeleteSubnetTestSuite) Test_DeleteSubnetWorkflow_Success() {
	var subnetManager subnetActivity.ManageSubnet

	vpcID := uuid.New()
	subnetID := uuid.New()

	// Mock DeleteSubnetViaSiteAgent activity
	s.env.RegisterActivity(subnetManager.DeleteSubnetViaSiteAgent)
	s.env.OnActivity(subnetManager.DeleteSubnetViaSiteAgent, mock.Anything, subnetID, vpcID).Return(nil)

	// execute deleteVPC workflow
	s.env.ExecuteWorkflow(DeleteSubnet, subnetID, vpcID)
	s.True(s.env.IsWorkflowCompleted())
	s.NoError(s.env.GetWorkflowError())
}

func (s *DeleteSubnetTestSuite) Test_DeleteSubnetWorkflow_ActivityFails() {
	var subnetManager subnetActivity.ManageSubnet

	vpcID := uuid.New()
	subnetID := uuid.New()

	// Mock DeleteSubnetViaSiteAgent activity failure
	s.env.RegisterActivity(subnetManager.DeleteSubnetViaSiteAgent)
	s.env.OnActivity(subnetManager.DeleteSubnetViaSiteAgent, mock.Anything, subnetID, vpcID).Return(errors.New("DeleteSubnetViaSiteAgent Failure"))

	// execute DeleteSubnet workflow
	s.env.ExecuteWorkflow(DeleteSubnet, subnetID, vpcID)
	s.True(s.env.IsWorkflowCompleted())
	err := s.env.GetWorkflowError()
	s.Error(err)

	var applicationErr *temporal.ApplicationError
	s.True(errors.As(err, &applicationErr))
	s.Equal("DeleteSubnetViaSiteAgent Failure", applicationErr.Error())
}

func (s *DeleteSubnetTestSuite) Test_ExecuteDeleteSubnetWorkflow_Success() {
	ctx := context.Background()

	vpcID := uuid.New()
	subnetID := uuid.New()

	wid := "test-workflow-id"

	wrun := &tmocks.WorkflowRun{}
	wrun.On("GetID").Return(wid)

	tc := &tmocks.Client{}

	tc.Mock.On("ExecuteWorkflow", context.Background(), mock.AnythingOfType("internal.StartWorkflowOptions"),
		mock.Anything, subnetID, vpcID).Return(wrun, nil)

	rwid, err := ExecuteDeleteSubnetWorkflow(ctx, tc, subnetID, vpcID)
	s.NoError(err)
	s.Equal(wid, *rwid)
}

func TestDeleteSubnetSuite(t *testing.T) {
	suite.Run(t, new(DeleteSubnetTestSuite))
}
