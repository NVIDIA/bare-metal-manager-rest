/*
 * SPDX-FileCopyrightText: Copyright (c) 2021-2023 NVIDIA CORPORATION & AFFILIATES. All rights reserved.
 * SPDX-License-Identifier: LicenseRef-NvidiaProprietary
 *
 * NVIDIA CORPORATION, its affiliates and licensors retain all intellectual
 * property and proprietary rights in and to this material, related
 * documentation and any modifications thereto. Any use, reproduction,
 * disclosure or distribution of this material and related documentation
 * without an express license agreement from NVIDIA CORPORATION or
 * its affiliates is strictly prohibited.
 */

package vpc

import (
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/testsuite"

	cdb "github.com/NVIDIA/carbide-rest-api/carbide-rest-db/pkg/db"
	cwm "github.com/NVIDIA/carbide-rest-api/carbide-rest-workflow/internal/metrics"
	vpcActivity "github.com/NVIDIA/carbide-rest-api/carbide-rest-workflow/pkg/activity/vpc"
)

type RecordVpcLifecycleMetricsTestSuite struct {
	suite.Suite
	testsuite.WorkflowTestSuite

	env *testsuite.TestWorkflowEnvironment
}

func (s *RecordVpcLifecycleMetricsTestSuite) SetupTest() {
	s.env = s.NewTestWorkflowEnvironment()
}

func (s *RecordVpcLifecycleMetricsTestSuite) AfterTest(suiteName, testName string) {
	s.env.AssertExpectations(s.T())
}

func (s *RecordVpcLifecycleMetricsTestSuite) Test_RecordVpcLifecycleMetrics_Success() {
	var lifecycleMetricsManager vpcActivity.ManageVpcLifecycleMetrics

	siteID := uuid.New()
	objectLifecycleEvents := []cwm.InventoryObjectLifecycleEvent{
		{ObjectID: uuid.New(), Deleted: cdb.GetTimePtr(cdb.GetCurTime())},
	}

	// Mock CollectAndRecordVpcOperationMetrics activity success
	s.env.RegisterActivity(lifecycleMetricsManager.RecordVpcStatusTransitionMetrics)
	s.env.OnActivity(lifecycleMetricsManager.RecordVpcStatusTransitionMetrics, mock.Anything, siteID, objectLifecycleEvents).Return(nil)

	// Execute CollectAndRecordVpcOperationMetrics workflow
	s.env.ExecuteWorkflow(RecordVpcLifecycleMetrics, siteID, objectLifecycleEvents)
	s.True(s.env.IsWorkflowCompleted())
	s.NoError(s.env.GetWorkflowError())
}

func (s *RecordVpcLifecycleMetricsTestSuite) Test_RecordVpcLifecycleMetrics_RecordVpcStatusTransitionMetricsActivityFails() {
	siteID := uuid.New()
	objectLifecycleEvents := []cwm.InventoryObjectLifecycleEvent{
		{ObjectID: uuid.New(), Deleted: cdb.GetTimePtr(cdb.GetCurTime())},
	}

	var lifecycleMetricsManager vpcActivity.ManageVpcLifecycleMetrics

	// Mock CollectAndRecordVpcOperationMetrics activity failure
	s.env.RegisterActivity(lifecycleMetricsManager.RecordVpcStatusTransitionMetrics)
	s.env.OnActivity(lifecycleMetricsManager.RecordVpcStatusTransitionMetrics, mock.Anything, siteID, objectLifecycleEvents).Return(errors.New("RecordVpcStatusTransitionMetrics Failure"))

	// Execute CollectAndRecordVpcOperationMetrics workflow
	s.env.ExecuteWorkflow(RecordVpcLifecycleMetrics, siteID, objectLifecycleEvents)
	s.True(s.env.IsWorkflowCompleted())
	err := s.env.GetWorkflowError()
	s.Error(err)

	var applicationErr *temporal.ApplicationError
	s.True(errors.As(err, &applicationErr))
	s.Equal("RecordVpcStatusTransitionMetrics Failure", applicationErr.Error())
}

func TestCollectAndRecordVpcOperationMetricsSuite(t *testing.T) {
	suite.Run(t, new(RecordVpcLifecycleMetricsTestSuite))
}
