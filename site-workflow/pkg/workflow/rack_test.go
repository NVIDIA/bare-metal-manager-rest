/*
 * SPDX-FileCopyrightText: Copyright (c) 2026 NVIDIA CORPORATION & AFFILIATES. All rights reserved.
 * SPDX-License-Identifier: Apache-2.0
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

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
