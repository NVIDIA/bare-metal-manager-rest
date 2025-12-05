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
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	wflows "github.com/nvidia/carbide-rest/workflow-schema/schema/site-agent/workflows/v1"
	"github.com/nvidia/carbide-rest/site-agent/pkg/components/managers/carbide"
	"github.com/nvidia/carbide-rest/site-agent/pkg/components/managers/subnet"
	"go.opentelemetry.io/otel"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/testsuite"
	"google.golang.org/protobuf/types/known/timestamppb"

	log "github.com/rs/zerolog/log"
)

type CreateSubnetTestSuite struct {
	suite.Suite
	testsuite.WorkflowTestSuite

	env *testsuite.TestWorkflowEnvironment
}

func (s *CreateSubnetTestSuite) SetupTest() {
	s.env = s.NewTestWorkflowEnvironment()
}

func (s *CreateSubnetTestSuite) AfterTest(suiteName, testName string) {
	s.env.AssertExpectations(s.T())
}

type CreateSubnetFailureTestSuite struct {
	suite.Suite
	testsuite.WorkflowTestSuite

	env *testsuite.TestWorkflowEnvironment
}

func (s *CreateSubnetFailureTestSuite) SetupTest() {
	s.env = s.NewTestWorkflowEnvironment()
}

func (s *CreateSubnetFailureTestSuite) AfterTest(suiteName, testName string) {
	s.env.AssertExpectations(s.T())
}

type UpdateSubnetTestSuite struct {
	suite.Suite
	testsuite.WorkflowTestSuite

	env *testsuite.TestWorkflowEnvironment
}

func (s *UpdateSubnetTestSuite) SetupTest() {
	s.env = s.NewTestWorkflowEnvironment()
}

func (s *UpdateSubnetTestSuite) AfterTest(suiteName, testName string) {
	s.env.AssertExpectations(s.T())
}

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

func (s *CreateSubnetFailureTestSuite) TestCreateSubnetWorkflowFailure() {
	log.Info().Msg("TestCreateSubnetWorkflowFailure Start")
	wflowinstance := subnet.Workflows{}
	subnetID := uuid.NewString()

	gw := "10.0.2.1"
	req := &wflows.CreateSubnetRequest{
		SubnetId:    &wflows.UUID{Value: subnetID},
		Name:        "test-subnet-2",
		VpcId:       &wflows.UUID{Value: uuid.NewString()},
		SubdomainId: &wflows.UUID{Value: uuid.NewString()},
		NetworkPrefixes: []*wflows.NetworkPrefixInfo{
			{
				Prefix:       "10.0.2.0/24",
				Gateway:      &gw,
				ReserveFirst: 1,
			},
		},
	}

	transaction := &wflows.TransactionID{
		ResourceId: subnetID,
		Timestamp: &timestamppb.Timestamp{
			Seconds: time.Now().Unix(),
		},
	}

	// Mock CreateSubnetActivity activity
	s.env.RegisterActivity(wflowinstance.CreateSubnetActivity)
	s.env.RegisterActivity(wflowinstance.PublishSubnetActivity)
	s.env.OnActivity(wflowinstance.PublishSubnetActivity, mock.Anything, transaction, mock.Anything).Return("", errors.New("CreateSubnetViaSiteAgent Failure"))

	// Execute CreateSubnet workflow
	s.env.ExecuteWorkflow(testElektra.manager.API.Subnet.CreateSubnet, transaction, req)
	s.True(s.env.IsWorkflowCompleted())
	err := s.env.GetWorkflowError()
	s.Error(err)

	var applicationErr *temporal.ApplicationError
	s.True(errors.As(err, &applicationErr))
	log.Info().Msg("TestCreateSubnetWorkflowFailure End")
}

func (s *CreateSubnetTestSuite) TestCreateSubnetWorkflowSuccess() {
	log.Info().Msg("TestCreateSubnetWorkflowSuccess Start")
	wflowinstance := subnet.Workflows{}
	subnetID := uuid.NewString()

	gw := "10.0.1.1"
	req := &wflows.CreateSubnetRequest{
		SubnetId:    &wflows.UUID{Value: subnetID},
		Name:        "test-subnet",
		VpcId:       &wflows.UUID{Value: uuid.NewString()},
		SubdomainId: &wflows.UUID{Value: uuid.NewString()},
		NetworkPrefixes: []*wflows.NetworkPrefixInfo{
			{
				Prefix:       "10.0.1.0/24",
				Gateway:      &gw,
				ReserveFirst: 1,
			},
		},
	}

	transaction := &wflows.TransactionID{
		ResourceId: uuid.NewString(),
		Timestamp: &timestamppb.Timestamp{
			Seconds: time.Now().Unix(),
		},
	}

	// Mock CreateSubnetActivity activity
	s.env.RegisterActivity(wflowinstance.CreateSubnetActivity)
	s.env.RegisterActivity(wflowinstance.PublishSubnetActivity)
	s.env.OnActivity(wflowinstance.PublishSubnetActivity, mock.Anything, transaction, mock.Anything).Return("", nil)

	// execute createSubnet workflow
	s.env.ExecuteWorkflow(testElektra.manager.API.Subnet.CreateSubnet, transaction, req)
	s.True(s.env.IsWorkflowCompleted())
	s.NoError(s.env.GetWorkflowError())

	resourceResp := &wflows.SubnetInfo{}
	err := s.env.GetWorkflowResult(resourceResp)
	if err != nil {
		log.Info().Msg(err.Error())
	}
	log.Info().Msg("TestCreateSubnetWorkflowSuccess End")
}

func (s *DeleteSubnetTestSuite) TestDeleteSubnetWorkflowSuccess() {
	log.Info().Msg("TestDeleteSubnetWorkflowSuccess Start")
	wflowinstance := subnet.Workflows{}
	req := &wflows.DeleteSubnetRequest{
		NetworkSegmentId: &wflows.UUID{Value: DefaultTestNetworkSegmentID},
	}

	transaction := &wflows.TransactionID{
		ResourceId: uuid.NewString(),
		Timestamp: &timestamppb.Timestamp{
			Seconds: time.Now().Unix(),
		},
	}

	// Mock DeleteSubnetActivity activity
	s.env.RegisterActivity(wflowinstance.DeleteSubnetActivity)
	s.env.RegisterActivity(wflowinstance.PublishSubnetActivity)
	s.env.OnActivity(wflowinstance.PublishSubnetActivity, mock.Anything, transaction, mock.Anything).Return("", nil)

	// execute DeleteSubnet workflow
	s.env.ExecuteWorkflow(testElektra.manager.API.Subnet.DeleteSubnet, transaction, req)
	s.True(s.env.IsWorkflowCompleted())
	s.NoError(s.env.GetWorkflowError())

	resourceResp := &wflows.SubnetInfo{}
	err := s.env.GetWorkflowResult(resourceResp)
	if err != nil {
		log.Info().Msg(err.Error())
	}
	log.Info().Msg("TestDeleteSubnetWorkflowSuccess End")
}

// TestSubnetWorkflows tests various Subnet workflows
func TestSubnetWorkflows(t *testing.T) {
	TestInitElektra(t)

	stats := subnet.ManagerAccess.Data.EB.Managers.Workflow.SubnetState
	stats.WflowActFail.Store(0)
	stats.WflowActSucc.Store(0)
	stats.WflowPubFail.Store(0)
	stats.WflowPubSucc.Store(0)
	wflowActFail = 0
	wflowActSucc = 1
	wflowPubFail = 0
	wflowPubSucc = 1

	carbide.ManagerAccess.Data.EB.Managers.Carbide.State.GrpcFail.Store(0)
	carbide.ManagerAccess.Data.EB.Managers.Carbide.State.GrpcSucc.Store(0)
	wflowGrpcFail = 0
	wflowGrpcSucc = 1

	_, span := otel.Tracer(os.Getenv("LS_SERVICE_NAME")).Start(context.Background(), "CreateSubnetTestSuite")
	suite.Run(t, new(CreateSubnetTestSuite))
	span.End()

	// failures has multiple tries
	_, span = otel.Tracer(os.Getenv("LS_SERVICE_NAME")).Start(context.Background(), "CreateSubnetFailureTestSuite")
	wflowGrpcFail += 7
	wflowActFail++
	wflowPubFail++
	suite.Run(t, new(CreateSubnetFailureTestSuite))
	span.End()

	_, span = otel.Tracer(os.Getenv("LS_SERVICE_NAME")).Start(context.Background(), "UpdateSubnetTestSuite")
	wflowGrpcSucc++
	wflowActSucc++
	wflowPubSucc++
	suite.Run(t, new(UpdateSubnetTestSuite))
	span.End()

	_, span = otel.Tracer(os.Getenv("LS_SERVICE_NAME")).Start(context.Background(), "DeleteSubnetTestSuite")
	wflowGrpcSucc++
	wflowActSucc++
	wflowPubSucc++
	suite.Run(t, new(DeleteSubnetTestSuite))
	span.End()
}
