// SPDX-FileCopyrightText: Copyright (c) 2021-2023 NVIDIA CORPORATION & AFFILIATES. All rights reserved.
// SPDX-License-Identifier: LicenseRef-NvidiaProprietary
//
// NVIDIA CORPORATION, its affiliates and licensors retain all intellectual
// property and proprietary rights in and to this material, related
// documentation and any modifications thereto. Any use, reproduction,
// disclosure or distribution of this material and related documentation
// without an express license agreement from NVIDIA CORPORATION or
// its affiliates is strictly prohibited.

package elektra

import (
	"context"
	"errors"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"go.opentelemetry.io/otel"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/testsuite"

	log "github.com/rs/zerolog/log"

	"github.com/nvidia/carbide-rest/site-agent/pkg/components/managers/rla"
	swa "github.com/nvidia/carbide-rest/site-workflow/pkg/activity"
	sww "github.com/nvidia/carbide-rest/site-workflow/pkg/workflow"
	rlav1 "github.com/nvidia/carbide-rest/workflow-schema/rla/protobuf/v1"
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

// GetRackFailureTestSuite tests the GetRack workflow failure scenarios
type GetRackFailureTestSuite struct {
	suite.Suite
	testsuite.WorkflowTestSuite

	env *testsuite.TestWorkflowEnvironment
}

func (s *GetRackFailureTestSuite) SetupTest() {
	s.env = s.NewTestWorkflowEnvironment()
}

func (s *GetRackFailureTestSuite) AfterTest(suiteName, testName string) {
	s.env.AssertExpectations(s.T())
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

// GetRacksFailureTestSuite tests the GetRacks workflow failure scenarios
type GetRacksFailureTestSuite struct {
	suite.Suite
	testsuite.WorkflowTestSuite

	env *testsuite.TestWorkflowEnvironment
}

func (s *GetRacksFailureTestSuite) SetupTest() {
	s.env = s.NewTestWorkflowEnvironment()
}

func (s *GetRacksFailureTestSuite) AfterTest(suiteName, testName string) {
	s.env.AssertExpectations(s.T())
}

// TestGetRackWorkflowSuccess tests successful GetRack workflow execution
func (s *GetRackTestSuite) TestGetRackWorkflowSuccess() {
	log.Info().Msg("TestGetRackWorkflowSuccess Start")

	// First, create a rack in the mock server
	rackID := uuid.NewString()
	rlaClient := rla.ManagerAccess.Data.EB.Managers.RLA.Client.GetClient()
	createReq := &rlav1.CreateExpectedRackRequest{
		Rack: &rlav1.Rack{
			Info: &rlav1.DeviceInfo{
				Id:   &rlav1.UUID{Id: rackID},
				Name: "test-rack",
			},
		},
	}
	_, err := rlaClient.Rla().CreateExpectedRack(context.Background(), createReq)
	s.NoError(err)

	// Now test GetRack workflow
	request := &rlav1.GetRackInfoByIDRequest{
		Id: &rlav1.UUID{Id: rackID},
	}

	// Create rack manager with RLA client
	rackManager := swa.NewManageRack(rla.ManagerAccess.Data.EB.Managers.RLA.Client)

	// Register activity - it will call the real RLA mock server
	s.env.RegisterActivity(rackManager.GetRack)

	// Execute GetRack workflow
	s.env.ExecuteWorkflow(sww.GetRack, request)
	s.True(s.env.IsWorkflowCompleted())
	s.NoError(s.env.GetWorkflowError())

	// Verify result
	var response rlav1.GetRackInfoResponse
	err = s.env.GetWorkflowResult(&response)
	s.NoError(err)
	s.NotNil(response.Rack)
	s.NotNil(response.Rack.Info)
	s.Equal(rackID, response.Rack.Info.Id.Id)

	log.Info().Msg("TestGetRackWorkflowSuccess End")
}

// TestGetRackWorkflowFailure tests GetRack workflow with activity failure
func (s *GetRackFailureTestSuite) TestGetRackWorkflowFailure() {
	log.Info().Msg("TestGetRackWorkflowFailure Start")

	rackID := uuid.NewString()
	request := &rlav1.GetRackInfoByIDRequest{
		Id: &rlav1.UUID{Id: rackID},
	}

	// Create rack manager
	rackManager := swa.NewManageRack(rla.ManagerAccess.Data.EB.Managers.RLA.Client)

	// Register activity and mock failure
	s.env.RegisterActivity(rackManager.GetRack)
	s.env.OnActivity(rackManager.GetRack, mock.Anything, mock.Anything).Return(nil, errors.New("RLA connection failed"))

	// Execute GetRack workflow
	s.env.ExecuteWorkflow(sww.GetRack, request)
	s.True(s.env.IsWorkflowCompleted())
	err := s.env.GetWorkflowError()
	s.Error(err)

	var applicationErr *temporal.ApplicationError
	s.True(errors.As(err, &applicationErr))

	log.Info().Msg("TestGetRackWorkflowFailure End")
}

// TestGetRacksWorkflowSuccess tests successful GetRacks workflow execution
func (s *GetRacksTestSuite) TestGetRacksWorkflowSuccess() {
	log.Info().Msg("TestGetRacksWorkflowSuccess Start")

	request := &rlav1.GetListOfRacksRequest{}

	// Create rack manager with RLA client
	rackManager := swa.NewManageRack(rla.ManagerAccess.Data.EB.Managers.RLA.Client)

	// Register activity - it will call the real RLA mock server
	s.env.RegisterActivity(rackManager.GetRacks)

	// Execute GetRacks workflow
	s.env.ExecuteWorkflow(sww.GetRacks, request)
	s.True(s.env.IsWorkflowCompleted())
	s.NoError(s.env.GetWorkflowError())

	// Verify result
	var response rlav1.GetListOfRacksResponse
	err := s.env.GetWorkflowResult(&response)
	s.NoError(err)
	s.NotNil(response.Racks)

	log.Info().Msg("TestGetRacksWorkflowSuccess End")
}

// TestGetRacksWorkflowFailure tests GetRacks workflow with activity failure
func (s *GetRacksFailureTestSuite) TestGetRacksWorkflowFailure() {
	log.Info().Msg("TestGetRacksWorkflowFailure Start")

	request := &rlav1.GetListOfRacksRequest{}

	// Create rack manager
	rackManager := swa.NewManageRack(rla.ManagerAccess.Data.EB.Managers.RLA.Client)

	// Register activity and mock failure
	s.env.RegisterActivity(rackManager.GetRacks)
	s.env.OnActivity(rackManager.GetRacks, mock.Anything, mock.Anything).Return(nil, errors.New("RLA connection failed"))

	// Execute GetRacks workflow
	s.env.ExecuteWorkflow(sww.GetRacks, request)
	s.True(s.env.IsWorkflowCompleted())
	err := s.env.GetWorkflowError()
	s.Error(err)

	var applicationErr *temporal.ApplicationError
	s.True(errors.As(err, &applicationErr))

	log.Info().Msg("TestGetRacksWorkflowFailure End")
}

// TestRackWorkflows tests various Rack workflows
func TestRackWorkflows(t *testing.T) {
	TestInitElektra(t)

	// Reset RLA state counters
	rla.ManagerAccess.Data.EB.Managers.RLA.State.GrpcFail.Store(0)
	rla.ManagerAccess.Data.EB.Managers.RLA.State.GrpcSucc.Store(0)

	_, span := otel.Tracer(os.Getenv("LS_SERVICE_NAME")).Start(context.Background(), "GetRackTestSuite")
	suite.Run(t, new(GetRackTestSuite))
	span.End()

	_, span = otel.Tracer(os.Getenv("LS_SERVICE_NAME")).Start(context.Background(), "GetRackFailureTestSuite")
	suite.Run(t, new(GetRackFailureTestSuite))
	span.End()

	_, span = otel.Tracer(os.Getenv("LS_SERVICE_NAME")).Start(context.Background(), "GetRacksTestSuite")
	suite.Run(t, new(GetRacksTestSuite))
	span.End()

	_, span = otel.Tracer(os.Getenv("LS_SERVICE_NAME")).Start(context.Background(), "GetRacksFailureTestSuite")
	suite.Run(t, new(GetRacksFailureTestSuite))
	span.End()
}
