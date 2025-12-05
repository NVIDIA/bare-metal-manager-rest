// SPDX-FileCopyrightText: Copyright (c) 2021-2023 NVIDIA CORPORATION & AFFILIATES. All rights reserved.
// SPDX-License-Identifier: LicenseRef-NvidiaProprietary
//
// NVIDIA CORPORATION, its affiliates and licensors retain all intellectual
// property and proprietary rights in and to this material, related
// documentation and any modifications thereto. Any use, reproduction,
// disclosure or distribution of this material and related documentation
// without an express license agreement from NVIDIA CORPORATION or
// its affiliates is strictly prohibited.

package vpc

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

	vpcActivity "github.com/nvidia/carbide-rest/workflow/pkg/activity/vpc"
)

type DeleteVpcTestSuite struct {
	suite.Suite
	testsuite.WorkflowTestSuite

	env *testsuite.TestWorkflowEnvironment
}

func (s *DeleteVpcTestSuite) SetupTest() {
	s.env = s.NewTestWorkflowEnvironment()
}

func (s *DeleteVpcTestSuite) AfterTest(suiteName, testName string) {
	s.env.AssertExpectations(s.T())
}

func (s *DeleteVpcTestSuite) Test_DeleteVPCWorkflow_Success() {
	var vpcManager vpcActivity.ManageVpc

	siteID := uuid.New()
	vpcID := uuid.New()

	// Mock DeleteVpcViaSiteAgent activity
	s.env.RegisterActivity(vpcManager.DeleteVpcViaSiteAgent)
	s.env.OnActivity(vpcManager.DeleteVpcViaSiteAgent, mock.Anything, siteID, vpcID).Return(nil)

	// execute deleteVPC workflow
	s.env.ExecuteWorkflow(DeleteVpc, siteID, vpcID)
	s.True(s.env.IsWorkflowCompleted())
	s.NoError(s.env.GetWorkflowError())
}

func (s *DeleteVpcTestSuite) Test_DeleteVPCWorkflow_ActivityFails() {
	var vpcManager vpcActivity.ManageVpc

	siteID := uuid.New()
	vpcID := uuid.New()

	// Mock DeleteVpcViaSiteAgent activity failure
	s.env.RegisterActivity(vpcManager.DeleteVpcViaSiteAgent)
	s.env.OnActivity(vpcManager.DeleteVpcViaSiteAgent, mock.Anything, siteID, vpcID).Return(errors.New("DeleteVpcViaSiteAgent Failure"))

	// execute createVPC workflow
	s.env.ExecuteWorkflow(DeleteVpc, siteID, vpcID)
	s.True(s.env.IsWorkflowCompleted())
	err := s.env.GetWorkflowError()
	s.Error(err)

	var applicationErr *temporal.ApplicationError
	s.True(errors.As(err, &applicationErr))
	s.Equal("DeleteVpcViaSiteAgent Failure", applicationErr.Error())
}

func (s *DeleteVpcTestSuite) Test_ExecuteDeleteVpcWorkflow_Success() {
	ctx := context.Background()
	siteID := uuid.New()
	vpcID := uuid.New()

	wid := "test-workflow-id"

	wrun := &tmocks.WorkflowRun{}
	wrun.On("GetID").Return(wid)

	tc := &tmocks.Client{}

	tc.Mock.On("ExecuteWorkflow", context.Background(), mock.AnythingOfType("internal.StartWorkflowOptions"),
		mock.Anything, siteID, vpcID).Return(wrun, nil)

	rwid, err := ExecuteDeleteVpcWorkflow(ctx, tc, siteID, vpcID)
	s.NoError(err)
	s.Equal(wid, *rwid)
}

func TestDeleteVpcSuite(t *testing.T) {
	suite.Run(t, new(DeleteVpcTestSuite))
}
