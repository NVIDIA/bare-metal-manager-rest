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

	tmocks "go.temporal.io/sdk/mocks"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/testsuite"

	instanceActivity "github.com/nvidia/carbide-rest/workflow/pkg/activity/instance"
)

type RebootInstanceTestSuite struct {
	suite.Suite
	testsuite.WorkflowTestSuite

	env *testsuite.TestWorkflowEnvironment
}

func (s *RebootInstanceTestSuite) SetupTest() {
	s.env = s.NewTestWorkflowEnvironment()
}

func (s *RebootInstanceTestSuite) AfterTest(suiteName, testName string) {
	s.env.AssertExpectations(s.T())
}

func (s *RebootInstanceTestSuite) Test_RebootInstanceWorkflow_Success() {
	var instanceManager instanceActivity.ManageInstance

	instanceID := uuid.New()

	// Mock RebootInstanceViaSiteAgent activity
	s.env.RegisterActivity(instanceManager.RebootInstanceViaSiteAgent)
	s.env.OnActivity(instanceManager.RebootInstanceViaSiteAgent, mock.Anything, instanceID, true, true).Return(nil)

	// execute RebootInstance workflow
	s.env.ExecuteWorkflow(RebootInstance, instanceID, true, true)
	s.True(s.env.IsWorkflowCompleted())
	s.NoError(s.env.GetWorkflowError())
}

func (s *RebootInstanceTestSuite) Test_RebootInstanceWorkflow_ActivityFailsErrorActivityFails() {
	var instanceManager instanceActivity.ManageInstance

	instanceID := uuid.New()

	// Mock RebootInstanceViaSiteAgent activity failure
	s.env.RegisterActivity(instanceManager.RebootInstanceViaSiteAgent)
	s.env.OnActivity(instanceManager.RebootInstanceViaSiteAgent, mock.Anything, instanceID, true, true).Return(errors.New("RebootInstanceViaSiteAgent Failure"))

	// execute RebootInstance workflow
	s.env.ExecuteWorkflow(RebootInstance, instanceID, true, true)
	s.True(s.env.IsWorkflowCompleted())
	err := s.env.GetWorkflowError()
	s.Error(err)

	var applicationErr *temporal.ApplicationError
	s.True(errors.As(err, &applicationErr))
	s.Equal("RebootInstanceViaSiteAgent Failure", applicationErr.Error())
}

func (s *RebootInstanceTestSuite) Test_ExecuteRebootInstanceWorkflow_Success() {
	ctx := context.Background()

	instanceID := uuid.New()

	wid := "test-workflow-id"

	wrun := &tmocks.WorkflowRun{}
	wrun.On("GetID").Return(wid)

	tc := &tmocks.Client{}

	tc.Mock.On("ExecuteWorkflow", context.Background(), mock.AnythingOfType("internal.StartWorkflowOptions"), mock.Anything,
		instanceID, true, true).Return(wrun, nil)

	rwid, err := ExecuteRebootInstanceWorkflow(ctx, tc, instanceID, true, true)
	s.NoError(err)
	s.Equal(wid, *rwid)
}

func (s *RebootInstanceTestSuite) Test_ExecuteRebootInstanceWorkflow_Failure() {
	ctx := context.Background()

	instanceID := uuid.New()

	wid := "test-workflow-id"

	wrun := &tmocks.WorkflowRun{}
	wrun.On("GetID").Return(wid)

	tc := &tmocks.Client{}

	tc.Mock.On("ExecuteWorkflow", context.Background(), mock.AnythingOfType("internal.StartWorkflowOptions"), mock.Anything,
		instanceID, true, true).Return(wrun, fmt.Errorf("failed to execute workflow"))

	_, err := ExecuteRebootInstanceWorkflow(ctx, tc, instanceID, true, true)
	s.Error(err)
}

func TestRebootInstanceSuite(t *testing.T) {
	suite.Run(t, new(RebootInstanceTestSuite))
}
