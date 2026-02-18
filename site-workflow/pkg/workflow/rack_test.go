// SPDX-FileCopyrightText: Copyright (c) 2021-2023 NVIDIA CORPORATION & AFFILIATES. All rights reserved.
// SPDX-License-Identifier: LicenseRef-NvidiaProprietary
//
// NVIDIA CORPORATION, its affiliates and licensors retain all intellectual
// property and proprietary rights in and to this material, related
// documentation and any modifications thereto. Any use, reproduction,
// disclosure or distribution of this material and related documentation
// without an express license agreement from NVIDIA CORPORATION or
// its affiliates is strictly prohibited.

package workflow

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/testsuite"

	rActivity "github.com/nvidia/bare-metal-manager-rest/site-workflow/pkg/activity"
	rlav1 "github.com/nvidia/bare-metal-manager-rest/workflow-schema/rla/protobuf/v1"
)

// GetRackTestSuite tests the GetRack workflow
type GetRackTestSuite struct {
	suite.Suite
	testsuite.WorkflowTestSuite

	env *testsuite.TestWorkflowEnvironment
}

func (s *GetRackTestSuite) SetupTest() {
	s.env = s.NewTestWorkflowEnvironment()
}

func (s *GetRackTestSuite) AfterTest(suiteName, testName string) {
	s.env.AssertExpectations(s.T())
}

func (s *GetRackTestSuite) Test_GetRack_Success() {
	var rackManager rActivity.ManageRack

	rackID := "test-rack-id"
	request := &rlav1.GetRackInfoByIDRequest{
		Id: &rlav1.UUID{Id: rackID},
	}

	expectedResponse := &rlav1.GetRackInfoResponse{
		Rack: &rlav1.Rack{
			Info: &rlav1.DeviceInfo{
				Id:   &rlav1.UUID{Id: rackID},
				Name: "test-rack",
			},
		},
	}

	// Mock GetRack activity
	s.env.RegisterActivity(rackManager.GetRack)
	s.env.OnActivity(rackManager.GetRack, mock.Anything, mock.Anything).Return(expectedResponse, nil)

	// Execute GetRack workflow
	s.env.ExecuteWorkflow(GetRack, request)
	s.True(s.env.IsWorkflowCompleted())
	s.NoError(s.env.GetWorkflowError())

	// Verify result
	var response rlav1.GetRackInfoResponse
	s.NoError(s.env.GetWorkflowResult(&response))
	s.NotNil(response.Rack)
	s.NotNil(response.Rack.Info)
	s.Equal(rackID, response.Rack.Info.Id.Id)
}

func (s *GetRackTestSuite) Test_GetRack_ActivityFails() {
	var rackManager rActivity.ManageRack

	rackID := "test-rack-id"
	request := &rlav1.GetRackInfoByIDRequest{
		Id: &rlav1.UUID{Id: rackID},
	}

	errMsg := "RLA connection failed"

	// Mock GetRack activity failure
	s.env.RegisterActivity(rackManager.GetRack)
	s.env.OnActivity(rackManager.GetRack, mock.Anything, mock.Anything).Return(nil, errors.New(errMsg))

	// Execute GetRack workflow
	s.env.ExecuteWorkflow(GetRack, request)
	s.True(s.env.IsWorkflowCompleted())
	err := s.env.GetWorkflowError()
	s.Error(err)

	var applicationErr *temporal.ApplicationError
	s.True(errors.As(err, &applicationErr))
	s.Equal(errMsg, applicationErr.Error())
}

func TestGetRackTestSuite(t *testing.T) {
	suite.Run(t, new(GetRackTestSuite))
}

// GetRacksTestSuite tests the GetRacks workflow
type GetRacksTestSuite struct {
	suite.Suite
	testsuite.WorkflowTestSuite

	env *testsuite.TestWorkflowEnvironment
}

func (s *GetRacksTestSuite) SetupTest() {
	s.env = s.NewTestWorkflowEnvironment()
}

func (s *GetRacksTestSuite) AfterTest(suiteName, testName string) {
	s.env.AssertExpectations(s.T())
}

func (s *GetRacksTestSuite) Test_GetRacks_Success() {
	var rackManager rActivity.ManageRack

	request := &rlav1.GetListOfRacksRequest{}

	expectedResponse := &rlav1.GetListOfRacksResponse{
		Racks: []*rlav1.Rack{
			{
				Info: &rlav1.DeviceInfo{
					Id:   &rlav1.UUID{Id: "rack-1"},
					Name: "Rack 1",
				},
			},
			{
				Info: &rlav1.DeviceInfo{
					Id:   &rlav1.UUID{Id: "rack-2"},
					Name: "Rack 2",
				},
			},
		},
		Total: 2,
	}

	// Mock GetRacks activity
	s.env.RegisterActivity(rackManager.GetRacks)
	s.env.OnActivity(rackManager.GetRacks, mock.Anything, mock.Anything).Return(expectedResponse, nil)

	// Execute GetRacks workflow
	s.env.ExecuteWorkflow(GetRacks, request)
	s.True(s.env.IsWorkflowCompleted())
	s.NoError(s.env.GetWorkflowError())

	// Verify result
	var response rlav1.GetListOfRacksResponse
	s.NoError(s.env.GetWorkflowResult(&response))
	s.NotNil(response.Racks)
	s.Equal(int32(2), response.Total)
	s.Equal(2, len(response.Racks))
}

func (s *GetRacksTestSuite) Test_GetRacks_ActivityFails() {
	var rackManager rActivity.ManageRack

	request := &rlav1.GetListOfRacksRequest{}

	errMsg := "RLA connection failed"

	// Mock GetRacks activity failure
	s.env.RegisterActivity(rackManager.GetRacks)
	s.env.OnActivity(rackManager.GetRacks, mock.Anything, mock.Anything).Return(nil, errors.New(errMsg))

	// Execute GetRacks workflow
	s.env.ExecuteWorkflow(GetRacks, request)
	s.True(s.env.IsWorkflowCompleted())
	err := s.env.GetWorkflowError()
	s.Error(err)

	var applicationErr *temporal.ApplicationError
	s.True(errors.As(err, &applicationErr))
	s.Equal(errMsg, applicationErr.Error())
}

func TestGetRacksTestSuite(t *testing.T) {
	suite.Run(t, new(GetRacksTestSuite))
}

// ValidateComponentsTestSuite tests the ValidateComponents workflow
type ValidateComponentsTestSuite struct {
	suite.Suite
	testsuite.WorkflowTestSuite

	env *testsuite.TestWorkflowEnvironment
}

func (s *ValidateComponentsTestSuite) SetupTest() {
	s.env = s.NewTestWorkflowEnvironment()
}

func (s *ValidateComponentsTestSuite) AfterTest(suiteName, testName string) {
	s.env.AssertExpectations(s.T())
}

func (s *ValidateComponentsTestSuite) Test_ValidateComponents_Success_NoDiffs() {
	var rackManager rActivity.ManageRack

	request := &rlav1.ValidateComponentsRequest{
		TargetSpec: &rlav1.OperationTargetSpec{
			Targets: &rlav1.OperationTargetSpec_Racks{
				Racks: &rlav1.RackTargets{
					Targets: []*rlav1.RackTarget{
						{
							Identifier: &rlav1.RackTarget_Id{
								Id: &rlav1.UUID{Id: "test-rack-id"},
							},
						},
					},
				},
			},
		},
	}

	expectedResponse := &rlav1.ValidateComponentsResponse{
		Diffs:               []*rlav1.ComponentDiff{},
		TotalDiffs:          0,
		OnlyInExpectedCount: 0,
		OnlyInActualCount:   0,
		DriftCount:          0,
		MatchCount:          5,
	}

	// Mock ValidateComponents activity
	s.env.RegisterActivity(rackManager.ValidateComponents)
	s.env.OnActivity(rackManager.ValidateComponents, mock.Anything, mock.Anything).Return(expectedResponse, nil)

	// Execute ValidateComponents workflow
	s.env.ExecuteWorkflow(ValidateComponents, request)
	s.True(s.env.IsWorkflowCompleted())
	s.NoError(s.env.GetWorkflowError())

	// Verify result
	var response rlav1.ValidateComponentsResponse
	s.NoError(s.env.GetWorkflowResult(&response))
	s.Equal(int32(0), response.TotalDiffs)
	s.Equal(int32(5), response.MatchCount)
	s.Equal(0, len(response.Diffs))
}

func (s *ValidateComponentsTestSuite) Test_ValidateComponents_Success_WithDiffs() {
	var rackManager rActivity.ManageRack

	request := &rlav1.ValidateComponentsRequest{
		TargetSpec: &rlav1.OperationTargetSpec{
			Targets: &rlav1.OperationTargetSpec_Racks{
				Racks: &rlav1.RackTargets{
					Targets: []*rlav1.RackTarget{
						{
							Identifier: &rlav1.RackTarget_Id{
								Id: &rlav1.UUID{Id: "test-rack-id"},
							},
						},
					},
				},
			},
		},
	}

	expectedResponse := &rlav1.ValidateComponentsResponse{
		Diffs: []*rlav1.ComponentDiff{
			{
				Type:        rlav1.DiffType_DIFF_TYPE_ONLY_IN_EXPECTED,
				ComponentId: "comp-1",
			},
			{
				Type:        rlav1.DiffType_DIFF_TYPE_DRIFT,
				ComponentId: "comp-2",
				FieldDiffs: []*rlav1.FieldDiff{
					{
						FieldName:     "firmware_version",
						ExpectedValue: "1.0.0",
						ActualValue:   "2.0.0",
					},
				},
			},
		},
		TotalDiffs:          2,
		OnlyInExpectedCount: 1,
		OnlyInActualCount:   0,
		DriftCount:          1,
		MatchCount:          3,
	}

	// Mock ValidateComponents activity
	s.env.RegisterActivity(rackManager.ValidateComponents)
	s.env.OnActivity(rackManager.ValidateComponents, mock.Anything, mock.Anything).Return(expectedResponse, nil)

	// Execute ValidateComponents workflow
	s.env.ExecuteWorkflow(ValidateComponents, request)
	s.True(s.env.IsWorkflowCompleted())
	s.NoError(s.env.GetWorkflowError())

	// Verify result
	var response rlav1.ValidateComponentsResponse
	s.NoError(s.env.GetWorkflowResult(&response))
	s.Equal(int32(2), response.TotalDiffs)
	s.Equal(int32(1), response.OnlyInExpectedCount)
	s.Equal(int32(1), response.DriftCount)
	s.Equal(int32(3), response.MatchCount)
	s.Equal(2, len(response.Diffs))
}

func (s *ValidateComponentsTestSuite) Test_ValidateComponents_ActivityFails() {
	var rackManager rActivity.ManageRack

	request := &rlav1.ValidateComponentsRequest{
		TargetSpec: &rlav1.OperationTargetSpec{
			Targets: &rlav1.OperationTargetSpec_Racks{
				Racks: &rlav1.RackTargets{
					Targets: []*rlav1.RackTarget{
						{
							Identifier: &rlav1.RackTarget_Id{
								Id: &rlav1.UUID{Id: "test-rack-id"},
							},
						},
					},
				},
			},
		},
	}

	errMsg := "RLA connection failed"

	// Mock ValidateComponents activity failure
	s.env.RegisterActivity(rackManager.ValidateComponents)
	s.env.OnActivity(rackManager.ValidateComponents, mock.Anything, mock.Anything).Return(nil, errors.New(errMsg))

	// Execute ValidateComponents workflow
	s.env.ExecuteWorkflow(ValidateComponents, request)
	s.True(s.env.IsWorkflowCompleted())
	err := s.env.GetWorkflowError()
	s.Error(err)

	var applicationErr *temporal.ApplicationError
	s.True(errors.As(err, &applicationErr))
	s.Equal(errMsg, applicationErr.Error())
}

func TestValidateComponentsTestSuite(t *testing.T) {
	suite.Run(t, new(ValidateComponentsTestSuite))
}
