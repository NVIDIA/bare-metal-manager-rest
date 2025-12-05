// SPDX-FileCopyrightText: Copyright (c) 2021-2023 NVIDIA CORPORATION & AFFILIATES. All rights reserved.
// SPDX-License-Identifier: LicenseRef-NvidiaProprietary
//
// NVIDIA CORPORATION, its affiliates and licensors retain all intellectual
// property and proprietary rights in and to this material, related
// documentation and any modifications thereto. Any use, reproduction,
// disclosure or distribution of this material and related documentation
// without an express license agreement from NVIDIA CORPORATION or
// its affiliates is strictly prohibited.

package site

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/testsuite"

	siteActivity "github.com/nvidia/carbide-rest/workflow/pkg/activity/site"
	tmocks "go.temporal.io/sdk/mocks"
)

type DeleteSiteComponentsTestSuite struct {
	suite.Suite
	testsuite.WorkflowTestSuite

	env *testsuite.TestWorkflowEnvironment
}

func (s *DeleteSiteComponentsTestSuite) SetupTest() {
	s.env = s.NewTestWorkflowEnvironment()
}

func (s *DeleteSiteComponentsTestSuite) AfterTest(suiteName, testName string) {
	s.env.AssertExpectations(s.T())
}

func (s *DeleteSiteComponentsTestSuite) Test_DeleteSiteComponentsWorkflow_Success() {
	var siteManager siteActivity.ManageSite

	siteID := uuid.New()
	infrastructureProviderID := uuid.New()
	purgeMachines := true

	// Mock DeleteSiteComponentsFromDB activity success
	s.env.RegisterActivity(siteManager.DeleteSiteComponentsFromDB)
	s.env.OnActivity(siteManager.DeleteSiteComponentsFromDB, mock.Anything, siteID, infrastructureProviderID, purgeMachines).Return(nil)

	// Execute DeleteSiteComponents workflow
	s.env.ExecuteWorkflow(DeleteSiteComponents, siteID, infrastructureProviderID, purgeMachines)
	s.env.AssertCalled(s.T(), "DeleteSiteComponentsFromDB", mock.Anything, siteID, infrastructureProviderID, purgeMachines)
	s.True(s.env.IsWorkflowCompleted())
	s.NoError(s.env.GetWorkflowError())
}

func (s *DeleteSiteComponentsTestSuite) Test_DeleteSiteComponentsWorkflow_ActivityFails() {
	var siteManager siteActivity.ManageSite

	siteID := uuid.New()
	infrastructureProviderID := uuid.New()
	purgeMachines := false

	// Mock DeleteSiteComponentFromDB activity failure
	s.env.RegisterActivity(siteManager.DeleteSiteComponentsFromDB)
	s.env.OnActivity(siteManager.DeleteSiteComponentsFromDB, mock.Anything, siteID, infrastructureProviderID, purgeMachines).Return(errors.New("DeleteSiteComponentsFromDB Failure"))

	// Execute DeleteSiteComponents workflow
	s.env.ExecuteWorkflow(DeleteSiteComponents, siteID, infrastructureProviderID, purgeMachines)
	s.True(s.env.IsWorkflowCompleted())
	err := s.env.GetWorkflowError()
	s.Error(err)

	var applicationErr *temporal.ApplicationError
	s.True(errors.As(err, &applicationErr))
	s.Equal("DeleteSiteComponentsFromDB Failure", applicationErr.Error())
}

func (s *DeleteSiteComponentsTestSuite) Test_ExecuteDeleteSiteComponentsWorkflow_Success() {
	ctx := context.Background()
	siteID := uuid.New()
	infrastructureProviderID := uuid.New()
	purgeMachines := true

	wid := "test-workflow-id"

	wrun := &tmocks.WorkflowRun{}
	wrun.On("GetID").Return(wid)

	tc := &tmocks.Client{}

	tc.Mock.On("ExecuteWorkflow", context.Background(), mock.AnythingOfType("internal.StartWorkflowOptions"),
		mock.Anything, siteID, infrastructureProviderID, purgeMachines).Return(wrun, nil)

	rwid, err := ExecuteDeleteSiteComponentsWorkflow(ctx, tc, siteID, infrastructureProviderID, purgeMachines)
	s.NoError(err)
	s.Equal(wid, *rwid)
}

func TestDeleteSiteComponentsSuite(t *testing.T) {
	suite.Run(t, new(DeleteSiteComponentsTestSuite))
}
