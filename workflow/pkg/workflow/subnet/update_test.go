// SPDX-FileCopyrightText: Copyright (c) 2021-2023 NVIDIA CORPORATION & AFFILIATES. All rights reserved.
// SPDX-License-Identifier: LicenseRef-NvidiaProprietary
//
// NVIDIA CORPORATION, its affiliates and licensors retain all intellectual
// property and proprietary rights in and to this material, related
// documentation and any modifications thereto. Any use, reproduction,
// disclosure or distribution of this material and related documentation
// without an express license agreement from NVIDIA CORPORATION or
// its affiliates is strictly prohibited.

package subnet

import (
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	cwm "github.com/nvidia/carbide-rest/workflow/internal/metrics"
	subnetActivity "github.com/nvidia/carbide-rest/workflow/pkg/activity/subnet"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/testsuite"
	"google.golang.org/protobuf/types/known/timestamppb"

	cwssaws "github.com/nvidia/carbide-rest/workflow-schema/schema/site-agent/workflows/v1"
)

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

func (s *UpdateSubnetTestSuite) Test_UpdateSubnetInfo_Success() {
	var subnetManager subnetActivity.ManageSubnet

	siteID := uuid.New()

	transactionID := &cwssaws.TransactionID{
		ResourceId: uuid.New().String(),
		Timestamp:  timestamppb.Now(),
	}

	subnetInfo := &cwssaws.SubnetInfo{
		Status:    cwssaws.WorkflowStatus_WORKFLOW_STATUS_IN_PROGRESS,
		StatusMsg: "Subnet creation in progress",
		NetworkSegment: &cwssaws.NetworkSegment{
			Id:   &cwssaws.NetworkSegmentId{Value: uuid.New().String()},
			Name: uuid.New().String(),
		},
	}

	// Mock UpdateSubnetInDB activity
	s.env.RegisterActivity(subnetManager.UpdateSubnetInDB)
	s.env.OnActivity(subnetManager.UpdateSubnetInDB, mock.Anything, mock.Anything, mock.Anything).Return(nil)

	// Execute UpdateSubnetInfo workflow
	s.env.ExecuteWorkflow(UpdateSubnetInfo, siteID.String(), transactionID, subnetInfo)
	s.True(s.env.IsWorkflowCompleted())
	s.NoError(s.env.GetWorkflowError())
}

func (s *UpdateSubnetTestSuite) Test_UpdateSubnetInfo_ActivityFails() {
	var subnetManager subnetActivity.ManageSubnet

	siteID := uuid.New()

	transactionID := &cwssaws.TransactionID{
		ResourceId: uuid.New().String(),
		Timestamp:  timestamppb.Now(),
	}

	subnetInfo := &cwssaws.SubnetInfo{
		Status:    cwssaws.WorkflowStatus_WORKFLOW_STATUS_IN_PROGRESS,
		StatusMsg: "Subnet creation in progress",
		NetworkSegment: &cwssaws.NetworkSegment{
			Id:   &cwssaws.NetworkSegmentId{Value: uuid.New().String()},
			Name: uuid.New().String(),
		},
	}

	// Mock UpdateSubnetInDB activity failure
	s.env.RegisterActivity(subnetManager.UpdateSubnetInDB)
	s.env.OnActivity(subnetManager.UpdateSubnetInDB, mock.Anything, mock.Anything, mock.Anything).Return(errors.New("UpdateSubnetInfo Failure"))

	// Execute UpdateSubnetInfo workflow
	s.env.ExecuteWorkflow(UpdateSubnetInfo, siteID.String(), transactionID, subnetInfo)
	s.True(s.env.IsWorkflowCompleted())
	err := s.env.GetWorkflowError()
	s.Error(err)

	var applicationErr *temporal.ApplicationError
	s.True(errors.As(err, &applicationErr))
	s.Equal("UpdateSubnetInfo Failure", applicationErr.Error())
}

func (s *UpdateSubnetTestSuite) Test_UpdateSubnetInventory_Success() {
	var subnetManager subnetActivity.ManageSubnet
	var lifecycleMetricsManager subnetActivity.ManageSubnetLifecycleMetrics
	var inventoryMetricsManager cwm.ManageInventoryMetrics

	siteID := uuid.New()

	subnetInventory := &cwssaws.SubnetInventory{
		Segments:  []*cwssaws.NetworkSegment{},
		Timestamp: timestamppb.Now(),
	}

	// Mock UpdateSubnetsInDB activity
	s.env.RegisterActivity(subnetManager.UpdateSubnetsInDB)
	s.env.OnActivity(subnetManager.UpdateSubnetsInDB, mock.Anything, siteID, mock.Anything).Return([]cwm.InventoryObjectLifecycleEvent{}, nil)

	// Mock RecordSubnetStatusTransitionMetrics activity
	s.env.RegisterActivity(lifecycleMetricsManager.RecordSubnetStatusTransitionMetrics)
	s.env.OnActivity(lifecycleMetricsManager.RecordSubnetStatusTransitionMetrics, mock.Anything, siteID, mock.Anything).Return(nil)

	// Mock RecordLatency activity
	s.env.RegisterActivity(inventoryMetricsManager.RecordLatency)
	s.env.OnActivity(inventoryMetricsManager.RecordLatency, mock.Anything, siteID, "UpdateSubnetInventory", false, mock.Anything).Return(nil)

	// Execute UpdateSubnetInventory workflow
	s.env.ExecuteWorkflow(UpdateSubnetInventory, siteID.String(), subnetInventory)
	s.True(s.env.IsWorkflowCompleted())
	s.NoError(s.env.GetWorkflowError())
}

func (s *UpdateSubnetTestSuite) Test_UpdateSubnetInventory_ActivityFails() {
	var subnetManager subnetActivity.ManageSubnet

	siteID := uuid.New()

	subnetInventory := &cwssaws.SubnetInventory{
		Segments:  []*cwssaws.NetworkSegment{},
		Timestamp: timestamppb.Now(),
	}

	// Mock UpdateSubnetsInDB activity
	s.env.RegisterActivity(subnetManager.UpdateSubnetsInDB)
	s.env.OnActivity(subnetManager.UpdateSubnetsInDB, mock.Anything, siteID, mock.Anything).Return([]cwm.InventoryObjectLifecycleEvent{}, errors.New("UpdateSubnetInventory Failure"))

	// Execute UpdateSubnetInventory workflow
	s.env.ExecuteWorkflow(UpdateSubnetInventory, siteID.String(), subnetInventory)
	s.True(s.env.IsWorkflowCompleted())
	err := s.env.GetWorkflowError()
	s.Error(err)

	var applicationErr *temporal.ApplicationError
	s.True(errors.As(err, &applicationErr))
	s.Equal("UpdateSubnetInventory Failure", applicationErr.Error())
}

func TestUpdateSubnetInfoSuite(t *testing.T) {
	suite.Run(t, new(UpdateSubnetTestSuite))
}
