// SPDX-FileCopyrightText: Copyright (c) 2021-2023 NVIDIA CORPORATION & AFFILIATES. All rights reserved.
// SPDX-License-Identifier: LicenseRef-NvidiaProprietary
//
// NVIDIA CORPORATION, its affiliates and licensors retain all intellectual
// property and proprietary rights in and to this material, related
// documentation and any modifications thereto. Any use, reproduction,
// disclosure or distribution of this material and related documentation
// without an express license agreement from NVIDIA CORPORATION or
// its affiliates is strictly prohibited.

package infinibandpartition

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

	ibpActivity "github.com/nvidia/carbide-rest/workflow/pkg/activity/infinibandpartition"
)

type DeleteInfiniBandPartitionTestSuite struct {
	suite.Suite
	testsuite.WorkflowTestSuite

	env *testsuite.TestWorkflowEnvironment
}

func (s *DeleteInfiniBandPartitionTestSuite) SetupTest() {
	s.env = s.NewTestWorkflowEnvironment()
}

func (s *DeleteInfiniBandPartitionTestSuite) AfterTest(suiteName, testName string) {
	s.env.AssertExpectations(s.T())
}

func (s *DeleteInfiniBandPartitionTestSuite) Test_DeleteInfiniBandPartitionWorkflow_Success() {
	var ibpManager ibpActivity.ManageInfiniBandPartition

	siteID := uuid.New()
	ibpID := uuid.New()

	// Mock DeleteInfiniBandPartitionViaSiteAgent activity
	s.env.RegisterActivity(ibpManager.DeleteInfiniBandPartitionViaSiteAgent)
	s.env.OnActivity(ibpManager.DeleteInfiniBandPartitionViaSiteAgent, mock.Anything, siteID, ibpID).Return(nil)

	// execute DeleteInfiniBandPartition workflow
	s.env.ExecuteWorkflow(DeleteInfiniBandPartition, siteID, ibpID)
	s.True(s.env.IsWorkflowCompleted())
	s.NoError(s.env.GetWorkflowError())
}

func (s *DeleteInfiniBandPartitionTestSuite) Test_DeleteInfiniBandPartitionWorkflow_ActivityFails() {
	var ibpManager ibpActivity.ManageInfiniBandPartition

	siteID := uuid.New()
	ibpID := uuid.New()

	// Mock DeleteInfiniBandPartitionViaSiteAgent activity failure
	s.env.RegisterActivity(ibpManager.DeleteInfiniBandPartitionViaSiteAgent)
	s.env.OnActivity(ibpManager.DeleteInfiniBandPartitionViaSiteAgent, mock.Anything, siteID, ibpID).Return(errors.New("DeleteInfiniBandPartitionViaSiteAgent Failure"))

	// execute createVPC workflow
	s.env.ExecuteWorkflow(DeleteInfiniBandPartition, siteID, ibpID)
	s.True(s.env.IsWorkflowCompleted())
	err := s.env.GetWorkflowError()
	s.Error(err)

	var applicationErr *temporal.ApplicationError
	s.True(errors.As(err, &applicationErr))
	s.Equal("DeleteInfiniBandPartitionViaSiteAgent Failure", applicationErr.Error())
}

func (s *DeleteInfiniBandPartitionTestSuite) Test_ExecuteDeleteInfiniBandPartitionWorkflow_Success() {
	ctx := context.Background()
	siteID := uuid.New()
	ibpID := uuid.New()

	wid := "test-workflow-id"

	wrun := &tmocks.WorkflowRun{}
	wrun.On("GetID").Return(wid)

	tc := &tmocks.Client{}

	tc.Mock.On("ExecuteWorkflow", context.Background(), mock.AnythingOfType("internal.StartWorkflowOptions"),
		mock.Anything, siteID, ibpID).Return(wrun, nil)

	rwid, err := ExecuteDeleteInfiniBandPartitionWorkflow(ctx, tc, siteID, ibpID)
	s.NoError(err)
	s.Equal(wid, *rwid)
}

func TestDeleteInfiniBandPartitionSuite(t *testing.T) {
	suite.Run(t, new(DeleteInfiniBandPartitionTestSuite))
}
