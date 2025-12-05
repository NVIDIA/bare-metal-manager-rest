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

	tmocks "go.temporal.io/sdk/mocks"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/testsuite"

	vpcActivity "github.com/nvidia/carbide-rest/workflow/pkg/activity/vpc"
)

type CreateVpcTestSuite struct {
	suite.Suite
	testsuite.WorkflowTestSuite

	env *testsuite.TestWorkflowEnvironment
}

func (s *CreateVpcTestSuite) SetupTest() {
	s.env = s.NewTestWorkflowEnvironment()
}

func (s *CreateVpcTestSuite) AfterTest(suiteName, testName string) {
	s.env.AssertExpectations(s.T())
}

func (s *CreateVpcTestSuite) Test_CreateVPCWorkflow_Success() {
	var vpcManager vpcActivity.ManageVpc

	siteID := uuid.New()
	vpcID := uuid.New()

	// Mock CreateVpcViaSiteAgent activity
	s.env.RegisterActivity(vpcManager.CreateVpcViaSiteAgent)
	s.env.OnActivity(vpcManager.CreateVpcViaSiteAgent, mock.Anything, siteID, vpcID).Return(nil)

	// execute createVPC workflow
	s.env.ExecuteWorkflow(CreateVpc, siteID, vpcID)
	s.True(s.env.IsWorkflowCompleted())
	s.NoError(s.env.GetWorkflowError())
}

func (s *CreateVpcTestSuite) Test_CreateVPCWorkflow_ActivityFails() {

	var vpcManager vpcActivity.ManageVpc

	siteID := uuid.New()
	vpcID := uuid.New()

	// Mock CreateVpcViaSiteAgent activity failure
	s.env.RegisterActivity(vpcManager.CreateVpcViaSiteAgent)
	s.env.OnActivity(vpcManager.CreateVpcViaSiteAgent, mock.Anything, siteID, vpcID).Return(errors.New("CreateVpcViaSiteAgent Failure"))

	// execute createVPC workflow
	s.env.ExecuteWorkflow(CreateVpc, siteID, vpcID)
	s.True(s.env.IsWorkflowCompleted())
	err := s.env.GetWorkflowError()
	s.Error(err)

	var applicationErr *temporal.ApplicationError
	s.True(errors.As(err, &applicationErr))
	s.Equal("CreateVpcViaSiteAgent Failure", applicationErr.Error())
}

func (s *CreateVpcTestSuite) Test_ExecuteCreateVpcWorkflow_Success() {
	ctx := context.Background()
	siteID := uuid.New()
	vpcID := uuid.New()

	wid := "test-workflow-id"

	wrun := &tmocks.WorkflowRun{}
	wrun.On("GetID").Return(wid)

	tc := &tmocks.Client{}

	tc.Mock.On("ExecuteWorkflow", context.Background(), mock.AnythingOfType("internal.StartWorkflowOptions"),
		mock.Anything, siteID, vpcID).Return(wrun, nil)

	rwid, err := ExecuteCreateVpcWorkflow(ctx, tc, siteID, vpcID)
	s.NoError(err)
	s.Equal(wid, *rwid)
}

func TestCreateVpcSuite(t *testing.T) {
	suite.Run(t, new(CreateVpcTestSuite))
}
