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
	"go.opentelemetry.io/otel"

	wflows "github.com/nvidia/carbide-rest/workflow-schema/schema/site-agent/workflows/v1"
	"github.com/nvidia/carbide-rest/site-agent/pkg/components/managers/carbide"
	"github.com/nvidia/carbide-rest/site-agent/pkg/components/managers/instance"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/testsuite"
	"google.golang.org/protobuf/types/known/timestamppb"

	log "github.com/rs/zerolog/log"
)

var (
	testInstance *wflows.Instance
)

type CreateDeprecatedInstanceTestSuite struct {
	suite.Suite
	testsuite.WorkflowTestSuite

	env *testsuite.TestWorkflowEnvironment
}

func (s *CreateDeprecatedInstanceTestSuite) SetupTest() {
	s.env = s.NewTestWorkflowEnvironment()
}

func (s *CreateDeprecatedInstanceTestSuite) AfterTest(suiteName, testName string) {
	s.env.AssertExpectations(s.T())
}

type CreateDeprecatedInstanceFailureTestSuite struct {
	suite.Suite
	testsuite.WorkflowTestSuite

	env *testsuite.TestWorkflowEnvironment
}

func (s *CreateDeprecatedInstanceFailureTestSuite) SetupTest() {
	s.env = s.NewTestWorkflowEnvironment()
}

func (s *CreateDeprecatedInstanceFailureTestSuite) AfterTest(suiteName, testName string) {
	s.env.AssertExpectations(s.T())
}

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

type DeprecatedRebootTestSuite struct {
	suite.Suite
	testsuite.WorkflowTestSuite

	env *testsuite.TestWorkflowEnvironment
}

func (s *DeprecatedRebootTestSuite) SetupTest() {
	s.env = s.NewTestWorkflowEnvironment()
}

func (s *DeprecatedRebootTestSuite) AfterTest(suiteName, testName string) {
	s.env.AssertExpectations(s.T())
}

type DeprecatedRebootFailureTestSuite struct {
	suite.Suite
	testsuite.WorkflowTestSuite

	env *testsuite.TestWorkflowEnvironment
}

func (s *DeprecatedRebootFailureTestSuite) SetupTest() {
	s.env = s.NewTestWorkflowEnvironment()
}

func (s *DeprecatedRebootFailureTestSuite) AfterTest(suiteName, testName string) {
	s.env.AssertExpectations(s.T())
}

func (s *CreateDeprecatedInstanceTestSuite) TestDeprecatedCreateInstanceWorkflowSuccess() {
	log.Info().Msg("TestCreateInstanceWorkflowSuccess Start")
	wflowinstance := instance.Workflows{}
	req := &wflows.CreateInstanceRequest{
		InstanceId:      &wflows.UUID{Value: testInstance.Id.Value},
		MachineId:       &wflows.MachineId{Id: testInstance.MachineId.Id},
		TenantOrg:       "test-tenant-org",
		Interfaces:      testInstance.Config.Network.Interfaces,
		CustomIpxe:      &testInstance.Config.Tenant.CustomIpxe,
		UserData:        testInstance.Config.Tenant.UserData,
		TenantKeysetIds: []string{uuid.NewString()},
	}

	transaction := &wflows.TransactionID{
		ResourceId: req.InstanceId.Value,
		Timestamp: &timestamppb.Timestamp{
			Seconds: time.Now().Unix(),
		},
	}

	// Mock CreateInstanceActivity activity
	s.env.RegisterActivity(wflowinstance.CreateInstanceActivity)
	s.env.RegisterActivity(wflowinstance.PublishInstanceActivity)
	s.env.OnActivity(wflowinstance.PublishInstanceActivity, mock.Anything, transaction, mock.Anything).Return("", nil)

	// execute createInstance workflow
	s.env.ExecuteWorkflow(testElektra.manager.API.Instance.CreateInstance, transaction, req)
	s.True(s.env.IsWorkflowCompleted())
	s.NoError(s.env.GetWorkflowError())

	resourceResp := &wflows.InstanceInfo{}
	err := s.env.GetWorkflowResult(resourceResp)
	if err != nil {
		log.Info().Msg(err.Error())
	}
	log.Info().Msg(resourceResp.String())
	log.Info().Msg("TestCreateInstanceWorkflowSuccess End")
}

func (s *CreateDeprecatedInstanceFailureTestSuite) TestDeprecatedCreateInstanceWorkflowFailure() {
	log.Info().Msg("TestCreateInstanceWorkflowFailure Start")
	wflowinstance := instance.Workflows{}
	req := &wflows.CreateInstanceRequest{
		InstanceId:      &wflows.UUID{Value: uuid.NewString()},
		MachineId:       &wflows.MachineId{Id: uuid.NewString()},
		TenantOrg:       "test-tenant-org",
		Interfaces:      testInstance.Config.Network.Interfaces,
		CustomIpxe:      &testInstance.Config.Tenant.CustomIpxe,
		UserData:        testInstance.Config.Tenant.UserData,
		TenantKeysetIds: []string{uuid.NewString()},
	}

	transaction := &wflows.TransactionID{
		ResourceId: req.InstanceId.Value,
		Timestamp: &timestamppb.Timestamp{
			Seconds: time.Now().Unix(),
		},
	}

	// Mock CreateInstanceActivity activity
	s.env.RegisterActivity(wflowinstance.CreateInstanceActivity)
	s.env.RegisterActivity(wflowinstance.PublishInstanceActivity)
	s.env.OnActivity(wflowinstance.PublishInstanceActivity, mock.Anything, transaction, mock.Anything).Return("", errors.New("PublishInstanceActivity Failure"))

	// Execute CreateInstance workflow
	s.env.ExecuteWorkflow(testElektra.manager.API.Instance.CreateInstance, transaction, req)
	s.True(s.env.IsWorkflowCompleted())
	err := s.env.GetWorkflowError()
	s.Error(err)

	var applicationErr *temporal.ApplicationError
	s.True(errors.As(err, &applicationErr))

	log.Info().Msg("TestCreateInstanceWorkflowFailure End")
}

func (s *DeleteInstanceTestSuite) TestDeleteInstanceWorkflowSuccess() {
	log.Info().Msg("TestDeleteInstanceWorkflowSuccess Start")
	wflowinstance := instance.Workflows{}
	req := &wflows.DeleteInstanceRequest{
		InstanceId: &wflows.UUID{Value: testInstance.Id.Value},
	}

	transaction := &wflows.TransactionID{
		ResourceId: testInstance.Id.Value,
		Timestamp: &timestamppb.Timestamp{
			Seconds: time.Now().Unix(),
		},
	}

	// Mock DeleteInstanceActivity activity
	s.env.RegisterActivity(wflowinstance.DeleteInstanceActivity)
	s.env.RegisterActivity(wflowinstance.PublishInstanceActivity)
	s.env.OnActivity(wflowinstance.PublishInstanceActivity, mock.Anything, transaction, mock.Anything).Return("", nil)

	// execute DeleteInstance workflow
	s.env.ExecuteWorkflow(testElektra.manager.API.Instance.DeleteInstance, transaction, req)
	s.True(s.env.IsWorkflowCompleted())
	s.NoError(s.env.GetWorkflowError())

	resourceResp := &wflows.InstanceInfo{}
	err := s.env.GetWorkflowResult(resourceResp)
	if err != nil {
		log.Info().Msg(err.Error())
	}
	log.Info().Msg(resourceResp.String())
	log.Info().Msg("TestDeleteInstanceWorkflowSuccess End")
}

func (s *DeprecatedRebootTestSuite) TestDeprecatedRebootWorkflowSuccess() {
	log.Info().Msg("TestDeprecatedRebootWorkflowSuccess Start")
	wflowinstance := instance.Workflows{}
	machineID := &wflows.MachineId{}
	machineID.Id = testInstance.MachineId.Id
	req := &wflows.RebootInstanceRequest{
		MachineId:          machineID,
		BootWithCustomIpxe: true,
	}

	transaction := &wflows.TransactionID{
		ResourceId: testInstance.Id.Value,
		Timestamp: &timestamppb.Timestamp{
			Seconds: time.Now().Unix(),
		},
	}

	// Mock DeprecatedRebootActivity activity
	s.env.RegisterActivity(wflowinstance.RebootInstanceActivity)
	s.env.RegisterActivity(wflowinstance.PublishInstancePowerStatus)
	s.env.OnActivity(wflowinstance.PublishInstancePowerStatus, mock.Anything, transaction, mock.Anything).Return("", nil)

	// execute DeprecatedReboot workflow
	s.env.ExecuteWorkflow(testElektra.manager.API.Instance.RebootInstance, transaction, req)
	s.True(s.env.IsWorkflowCompleted())
	s.NoError(s.env.GetWorkflowError())

	resourceResp := &wflows.InstanceInfo{}
	err := s.env.GetWorkflowResult(resourceResp)
	if err != nil {
		log.Info().Msg(err.Error())
	}
	log.Info().Msg(resourceResp.String())
	log.Info().Msg("TestDeprecatedRebootWorkflowSuccess End")
}

func (s *DeprecatedRebootFailureTestSuite) TestDeprecatedRebootWorkflowFailure() {
	log.Info().Msg("TestDeprecatedRebootWorkflowFailure Start")
	wflowinstance := instance.Workflows{}
	req := &wflows.RebootInstanceRequest{
		MachineId:          nil,
		BootWithCustomIpxe: true,
	}

	transaction := &wflows.TransactionID{
		ResourceId: testInstance.Id.Value,
		Timestamp: &timestamppb.Timestamp{
			Seconds: time.Now().Unix(),
		},
	}

	// Mock DeprecatedRebootActivity activity
	s.env.RegisterActivity(wflowinstance.RebootInstanceActivity)
	s.env.RegisterActivity(wflowinstance.PublishInstancePowerStatus)
	s.env.OnActivity(wflowinstance.PublishInstancePowerStatus, mock.Anything, transaction, mock.Anything).Return("", nil)

	// execute DeprecatedReboot workflow
	s.env.ExecuteWorkflow(testElektra.manager.API.Instance.RebootInstance, transaction, req)
	s.True(s.env.IsWorkflowCompleted())
	err := s.env.GetWorkflowError()
	s.Error(err)

	var applicationErr *temporal.ApplicationError
	s.True(errors.As(err, &applicationErr))
	log.Info().Msg("TestDeprecatedRebootWorkflowFailure End")
}

// TestInstanceWorkflows tests various Instance workflows
func TestInstanceWorkflows(t *testing.T) {
	TestInitElektra(t)
	carbide.ManagerAccess.Data.EB.Managers.Carbide.State.GrpcFail.Store(0)
	carbide.ManagerAccess.Data.EB.Managers.Carbide.State.GrpcSucc.Store(0)
	wflowGrpcFail = 0
	wflowGrpcSucc = 1

	stats := instance.ManagerAccess.Data.EB.Managers.Workflow.InstanceState
	stats.WflowActFail.Store(0)
	stats.WflowActSucc.Store(0)
	stats.WflowPubFail.Store(0)
	stats.WflowPubSucc.Store(0)
	wflowActFail = 0
	wflowActSucc = 1
	wflowPubFail = 0
	wflowPubSucc = 1

	// Send a request to the GRPC server
	testUserData := "test-user-data"
	testInstance = &wflows.Instance{
		Id:        &wflows.InstanceId{Value: uuid.NewString()},
		MachineId: &wflows.MachineId{Id: uuid.NewString()},
		Config: &wflows.InstanceConfig{
			Tenant: &wflows.TenantConfig{
				TenantOrganizationId: "test-tenant-org",
				CustomIpxe:           "sample-ipxe-script",
				UserData:             &testUserData,
			},
			Network: &wflows.InstanceNetworkConfig{
				Interfaces: []*wflows.InstanceInterfaceConfig{
					&wflows.InstanceInterfaceConfig{
						FunctionType:     1,
						NetworkSegmentId: &wflows.NetworkSegmentId{Value: uuid.NewString()},
					},
					&wflows.InstanceInterfaceConfig{
						FunctionType:     0,
						NetworkSegmentId: &wflows.NetworkSegmentId{Value: uuid.NewString()},
					},
				},
			},
			Os: &wflows.OperatingSystem{
				RunProvisioningInstructionsOnEveryBoot: true,
				Variant: &wflows.OperatingSystem_Ipxe{
					Ipxe: &wflows.IpxeOperatingSystem{
						IpxeScript: "#!ipxe",
						UserData:   &testUserData,
					},
				},
			},
		},
		Status: &wflows.InstanceStatus{},
	}
	_, span := otel.Tracer(os.Getenv("LS_SERVICE_NAME")).Start(context.Background(), "CreateDeprecatedInstanceTestSuite")
	suite.Run(t, new(CreateDeprecatedInstanceTestSuite))
	span.End()

	// Failure triggers multiple retries
	_, span = otel.Tracer(os.Getenv("LS_SERVICE_NAME")).Start(context.Background(), "CreateDeprecatedInstanceFailureTestSuite")
	wflowGrpcSucc++
	wflowActSucc++
	wflowPubFail++
	suite.Run(t, new(CreateDeprecatedInstanceFailureTestSuite))
	span.End()

	_, span = otel.Tracer(os.Getenv("LS_SERVICE_NAME")).Start(context.Background(), "DeprecatedRebootTestSuite")
	wflowGrpcSucc++
	wflowActSucc++
	wflowPubSucc++
	suite.Run(t, new(DeprecatedRebootTestSuite))
	span.End()

	_, span = otel.Tracer(os.Getenv("LS_SERVICE_NAME")).Start(context.Background(), "DeprecatedRebootFailureTestSuite")
	wflowGrpcFail += 7
	wflowActFail++
	wflowPubSucc++
	suite.Run(t, new(DeprecatedRebootFailureTestSuite))
	span.End()

	_, span = otel.Tracer(os.Getenv("LS_SERVICE_NAME")).Start(context.Background(), "DeleteInstanceTestSuite")
	wflowGrpcSucc++
	wflowActSucc++
	wflowPubSucc++
	suite.Run(t, new(DeleteInstanceTestSuite))
	span.End()
}
