// SPDX-FileCopyrightText: Copyright (c) 2021-2023 NVIDIA CORPORATION & AFFILIATES. All rights reserved.
// SPDX-License-Identifier: LicenseRef-NvidiaProprietary
//
// NVIDIA CORPORATION, its affiliates and licensors retain all intellectual
// property and proprietary rights in and to this material, related
// documentation and any modifications thereto. Any use, reproduction,
// disclosure or distribution of this material and related documentation
// without an express license agreement from NVIDIA CORPORATION or
// its affiliates is strictly prohibited.

package instance

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/testsuite"

	tmocks "go.temporal.io/sdk/mocks"

	instanceActivity "github.com/nvidia/carbide-rest/workflow/pkg/activity/instance"
)

type DeleteInstanceTestSuite struct {
	suite.Suite
	testsuite.WorkflowTestSuite

	env *testsuite.TestWorkflowEnvironment
}

func (s *DeleteInstanceTestSuite) SetupTest() {
	s.env = s.NewTestWorkflowEnvironment()
}

func (s *DeleteInstanceTestSuite) AfterTest(suiteName, testName string) {
	s.env.AssertExpectations(s.T())
}

func (s *DeleteInstanceTestSuite) Test_DeleteInstanceWorkflow_Success() {
	var instanceManager instanceActivity.ManageInstance

	instanceID := uuid.New()

	// Mock DeleteInstanceViaSiteAgent activity
	s.env.RegisterActivity(instanceManager.DeleteInstanceViaSiteAgent)
	s.env.OnActivity(instanceManager.DeleteInstanceViaSiteAgent, mock.Anything, instanceID).Return(nil)

	// execute deleteVPC workflow
	s.env.ExecuteWorkflow(DeleteInstance, instanceID)
	s.True(s.env.IsWorkflowCompleted())
	s.NoError(s.env.GetWorkflowError())
}

func (s *DeleteInstanceTestSuite) Test_DeleteInstanceWorkflow_DeleteInstanceViaSiteAgentActivityFails() {
	var instanceManager instanceActivity.ManageInstance

	instanceID := uuid.New()

	// Mock DeleteInstanceViaSiteAgent activity failure
	s.env.RegisterActivity(instanceManager.DeleteInstanceViaSiteAgent)
	s.env.OnActivity(instanceManager.DeleteInstanceViaSiteAgent, mock.Anything, instanceID).Return(errors.New("DeleteInstanceViaSiteAgent Failure"))

	// execute DeleteInstance workflow
	s.env.ExecuteWorkflow(DeleteInstance, instanceID)
	s.True(s.env.IsWorkflowCompleted())
	err := s.env.GetWorkflowError()
	s.Error(err)

	var applicationErr *temporal.ApplicationError
	s.True(errors.As(err, &applicationErr))
	s.Equal("DeleteInstanceViaSiteAgent Failure", applicationErr.Error())
}

func (s *DeleteInstanceTestSuite) Test_ExecuteDeleteInstanceWorkflow_Success() {
	ctx := context.Background()

	instanceID := uuid.New()

	wid := "test-workflow-id"

	wrun := &tmocks.WorkflowRun{}
	wrun.On("GetID").Return(wid)

	tc := &tmocks.Client{}

	tc.Mock.On("ExecuteWorkflow", context.Background(), mock.AnythingOfType("internal.StartWorkflowOptions"),
		mock.Anything, instanceID).Return(wrun, nil)

	rwid, err := ExecuteDeleteInstanceWorkflow(ctx, tc, instanceID)
	s.NoError(err)
	s.Equal(wid, *rwid)
}

func (s *DeleteInstanceTestSuite) Test_ExecuteDeleteInstanceWorkflow_Failure() {
	ctx := context.Background()

	instanceID := uuid.New()

	wid := "test-workflow-id"

	wrun := &tmocks.WorkflowRun{}
	wrun.On("GetID").Return(wid)

	tc := &tmocks.Client{}

	tc.Mock.On("ExecuteWorkflow", context.Background(), mock.AnythingOfType("internal.StartWorkflowOptions"),
		mock.Anything, instanceID).Return(wrun, fmt.Errorf("failed to execute workflow"))

	_, err := ExecuteDeleteInstanceWorkflow(ctx, tc, instanceID)
	s.Error(err)
}

func TestDeleteInstanceSuite(t *testing.T) {
	suite.Run(t, new(DeleteInstanceTestSuite))
}
