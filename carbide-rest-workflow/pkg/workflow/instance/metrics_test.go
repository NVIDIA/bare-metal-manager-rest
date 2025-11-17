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

package instance

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
	instanceActivity "github.com/NVIDIA/carbide-rest-api/carbide-rest-workflow/pkg/activity/instance"
)

type RecordInstanceLifecycleMetricsTestSuite struct {
	suite.Suite
	testsuite.WorkflowTestSuite

	env *testsuite.TestWorkflowEnvironment
}

func (s *RecordInstanceLifecycleMetricsTestSuite) SetupTest() {
	s.env = s.NewTestWorkflowEnvironment()
}

func (s *RecordInstanceLifecycleMetricsTestSuite) AfterTest(suiteName, testName string) {
	s.env.AssertExpectations(s.T())
}

func (s *RecordInstanceLifecycleMetricsTestSuite) Test_RecordInstanceLifecycleMetrics_Success() {
	var lifecycleMetricsManager instanceActivity.ManageInstanceLifecycleMetrics

	siteID := uuid.New()
	objectLifecycleEvents := []cwm.InventoryObjectLifecycleEvent{
		{ObjectID: uuid.New(), Created: cdb.GetTimePtr(cdb.GetCurTime())},
		{ObjectID: uuid.New(), Deleted: cdb.GetTimePtr(cdb.GetCurTime())},
	}

	// Mock RecordInstanceLifecycleMetrics activity success
	s.env.RegisterActivity(lifecycleMetricsManager.RecordInstanceStatusTransitionMetrics)
	s.env.OnActivity(lifecycleMetricsManager.RecordInstanceStatusTransitionMetrics, mock.Anything, siteID, objectLifecycleEvents).Return(nil)

	// Execute CollectAndRecordInstanceOperationMetrics workflow
	s.env.ExecuteWorkflow(RecordInstanceLifecycleMetrics, siteID, objectLifecycleEvents)
	s.True(s.env.IsWorkflowCompleted())
	s.NoError(s.env.GetWorkflowError())
}

func (s *RecordInstanceLifecycleMetricsTestSuite) Test_RecordInstanceLifecycleMetrics_RecordInstanceStatusTransitionMetricsActivityFails() {
	siteID := uuid.New()
	objectLifecycleEvents := []cwm.InventoryObjectLifecycleEvent{
		{ObjectID: uuid.New(), Created: cdb.GetTimePtr(cdb.GetCurTime())},
		{ObjectID: uuid.New(), Deleted: cdb.GetTimePtr(cdb.GetCurTime())},
	}

	var lifecycleMetricsManager instanceActivity.ManageInstanceLifecycleMetrics

	// Mock CollectAndRecordInstanceOperationMetrics activity success
	s.env.RegisterActivity(lifecycleMetricsManager.RecordInstanceStatusTransitionMetrics)
	s.env.OnActivity(lifecycleMetricsManager.RecordInstanceStatusTransitionMetrics, mock.Anything, siteID, objectLifecycleEvents).Return(errors.New("RecordInstanceStatusTransitionMetrics Failure"))

	// Execute CollectAndRecordInstanceOperationMetrics workflow
	s.env.ExecuteWorkflow(RecordInstanceLifecycleMetrics, siteID, objectLifecycleEvents)
	s.True(s.env.IsWorkflowCompleted())
	err := s.env.GetWorkflowError()
	s.Error(err)

	var applicationErr *temporal.ApplicationError
	s.True(errors.As(err, &applicationErr))
	s.Equal("RecordInstanceStatusTransitionMetrics Failure", applicationErr.Error())
}

func TestRecordInstanceLifecycleMetricsSuite(t *testing.T) {
	suite.Run(t, new(RecordInstanceLifecycleMetricsTestSuite))
}
