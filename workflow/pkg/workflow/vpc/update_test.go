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
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	cwm "github.com/nvidia/carbide-rest/workflow/internal/metrics"
	vpcActivity "github.com/nvidia/carbide-rest/workflow/pkg/activity/vpc"

	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/testsuite"
	"google.golang.org/protobuf/types/known/timestamppb"

	cwssaws "github.com/nvidia/carbide-rest/workflow-schema/schema/site-agent/workflows/v1"
)

type UpdateVpcTestSuite struct {
	suite.Suite
	testsuite.WorkflowTestSuite

	env *testsuite.TestWorkflowEnvironment
}

func (s *UpdateVpcTestSuite) SetupTest() {
	s.env = s.NewTestWorkflowEnvironment()
}

func (s *UpdateVpcTestSuite) AfterTest(suiteName, testName string) {
	s.env.AssertExpectations(s.T())
}

func (s *UpdateVpcTestSuite) Test_UpdateVpcInfo_Success() {
	var vpcManager vpcActivity.ManageVpc

	siteID := uuid.New()

	transactionID := &cwssaws.TransactionID{
		ResourceId: uuid.New().String(),
		Timestamp:  timestamppb.Now(),
	}

	vpcInfo := &cwssaws.VPCInfo{
		Status:    cwssaws.WorkflowStatus_WORKFLOW_STATUS_IN_PROGRESS,
		StatusMsg: "VPC creation in progress",
		Vpc: &cwssaws.Vpc{
			Id:                   &cwssaws.VpcId{Value: uuid.New().String()},
			Name:                 uuid.New().String(),
			TenantOrganizationId: uuid.NewString(),
		},
	}

	// Mock UpdateVpcViaSiteAgent activity
	s.env.RegisterActivity(vpcManager.UpdateVpcInDB)
	s.env.OnActivity(vpcManager.UpdateVpcInDB, mock.Anything, mock.Anything, mock.Anything).Return(nil)

	// execute UpdateVpcInfo workflow
	s.env.ExecuteWorkflow(UpdateVpcInfo, siteID.String(), transactionID, vpcInfo)
	s.True(s.env.IsWorkflowCompleted())
	s.NoError(s.env.GetWorkflowError())
}

func (s *UpdateVpcTestSuite) Test_UpdateVpcInfo_ActivityFails() {
	var vpcManager vpcActivity.ManageVpc

	siteID := uuid.New()

	transactionID := &cwssaws.TransactionID{
		ResourceId: uuid.New().String(),
		Timestamp:  timestamppb.Now(),
	}

	vpcInfo := &cwssaws.VPCInfo{
		Status:    cwssaws.WorkflowStatus_WORKFLOW_STATUS_IN_PROGRESS,
		StatusMsg: "VPC creation in progress",
		Vpc: &cwssaws.Vpc{
			Id:                   &cwssaws.VpcId{Value: uuid.New().String()},
			Name:                 uuid.New().String(),
			TenantOrganizationId: uuid.NewString(),
		},
	}

	// Mock UpdateVpcViaSiteAgent activity failure
	s.env.RegisterActivity(vpcManager.UpdateVpcInDB)
	s.env.OnActivity(vpcManager.UpdateVpcInDB, mock.Anything, mock.Anything, mock.Anything).Return(errors.New("UpdateVpcInfo Failure"))

	// execute UpdateVPCStatus workflow
	s.env.ExecuteWorkflow(UpdateVpcInfo, siteID.String(), transactionID, vpcInfo)
	s.True(s.env.IsWorkflowCompleted())
	err := s.env.GetWorkflowError()
	s.Error(err)

	var applicationErr *temporal.ApplicationError
	s.True(errors.As(err, &applicationErr))
	s.Equal("UpdateVpcInfo Failure", applicationErr.Error())
}

func (s *UpdateVpcTestSuite) Test_UpdateVpcInventory_Success() {
	var vpcManager vpcActivity.ManageVpc
	var lifecycleMetricsManager vpcActivity.ManageVpcLifecycleMetrics
	var inventoryMetricsManager cwm.ManageInventoryMetrics

	siteID := uuid.New()
	vpcInventory := &cwssaws.VPCInventory{
		Vpcs: []*cwssaws.Vpc{
			{
				Id: &cwssaws.VpcId{Value: uuid.NewString()},
			},
			{
				Id: &cwssaws.VpcId{Value: uuid.NewString()},
			},
		},
	}

	// Mock UpdateVpcsInDB activity
	s.env.RegisterActivity(vpcManager.UpdateVpcsInDB)
	s.env.OnActivity(vpcManager.UpdateVpcsInDB, mock.Anything, siteID, mock.Anything).Return([]cwm.InventoryObjectLifecycleEvent{}, nil)

	// Mock RecordVpcStatusTransitionMetrics activity
	s.env.RegisterActivity(lifecycleMetricsManager.RecordVpcStatusTransitionMetrics)
	s.env.OnActivity(lifecycleMetricsManager.RecordVpcStatusTransitionMetrics, mock.Anything, siteID, mock.Anything).Return(nil)

	// Mock RecordLatency activity
	s.env.RegisterActivity(inventoryMetricsManager.RecordLatency)
	s.env.OnActivity(inventoryMetricsManager.RecordLatency, mock.Anything, siteID, "UpdateVpcInventory", false, mock.Anything).Return(nil)

	// execute UpdateVpcInventory workflow
	s.env.ExecuteWorkflow(UpdateVpcInventory, siteID.String(), vpcInventory)
	s.True(s.env.IsWorkflowCompleted())
	s.NoError(s.env.GetWorkflowError())
}

func (s *UpdateVpcTestSuite) Test_UpdateVpcInventory_ActivityFails() {
	var vpcManager vpcActivity.ManageVpc

	siteID := uuid.New()
	vpcInventory := &cwssaws.VPCInventory{
		Vpcs: []*cwssaws.Vpc{
			{
				Id: &cwssaws.VpcId{Value: uuid.NewString()},
			},
			{
				Id: &cwssaws.VpcId{Value: uuid.NewString()},
			},
		},
	}

	// Mock UpdateVpcsInDB activity failure
	s.env.RegisterActivity(vpcManager.UpdateVpcsInDB)
	s.env.OnActivity(vpcManager.UpdateVpcsInDB, mock.Anything, siteID, mock.Anything).Return([]cwm.InventoryObjectLifecycleEvent{}, errors.New("UpdateVpcInventory Failure"))

	// execute UpdateVPCStatus workflow
	s.env.ExecuteWorkflow(UpdateVpcInventory, siteID.String(), vpcInventory)
	s.True(s.env.IsWorkflowCompleted())
	err := s.env.GetWorkflowError()
	s.Error(err)

	var applicationErr *temporal.ApplicationError
	s.True(errors.As(err, &applicationErr))
	s.Equal("UpdateVpcInventory Failure", applicationErr.Error())
}

func TestUpdateVpcSuite(t *testing.T) {
	suite.Run(t, new(UpdateVpcTestSuite))
}
