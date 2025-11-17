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

package subnet

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
	subnetActivity "github.com/NVIDIA/carbide-rest-api/carbide-rest-workflow/pkg/activity/subnet"
)

type RecordSubnetLifecycleMetricsTestSuite struct {
	suite.Suite
	testsuite.WorkflowTestSuite

	env *testsuite.TestWorkflowEnvironment
}

func (s *RecordSubnetLifecycleMetricsTestSuite) SetupTest() {
	s.env = s.NewTestWorkflowEnvironment()
}

func (s *RecordSubnetLifecycleMetricsTestSuite) AfterTest(suiteName, testName string) {
	s.env.AssertExpectations(s.T())
}

func (s *RecordSubnetLifecycleMetricsTestSuite) Test_RecordSubnetLifecycleMetrics_Success() {
	var lifecycleMetricsManager subnetActivity.ManageSubnetLifecycleMetrics

	siteID := uuid.New()
	objectLifecycleEvents := []cwm.InventoryObjectLifecycleEvent{
		{ObjectID: uuid.New(), Created: cdb.GetTimePtr(cdb.GetCurTime())},
		{ObjectID: uuid.New(), Deleted: cdb.GetTimePtr(cdb.GetCurTime())},
	}

	// Mock CollectAndRecordSubnetOperationMetrics activity success
	s.env.RegisterActivity(lifecycleMetricsManager.RecordSubnetStatusTransitionMetrics)
	s.env.OnActivity(lifecycleMetricsManager.RecordSubnetStatusTransitionMetrics, mock.Anything, siteID, objectLifecycleEvents).Return(nil)

	// Execute CollectAndRecordSubnetOperationMetrics workflow
	s.env.ExecuteWorkflow(RecordSubnetLifecycleMetrics, siteID, objectLifecycleEvents)
	s.True(s.env.IsWorkflowCompleted())
	s.NoError(s.env.GetWorkflowError())
}

func (s *RecordSubnetLifecycleMetricsTestSuite) Test_RecordSubnetLifecycleMetrics_RecordSubnetStatusTransitionMetricsActivityFails() {
	siteID := uuid.New()
	objectLifecycleEvents := []cwm.InventoryObjectLifecycleEvent{
		{ObjectID: uuid.New(), Created: cdb.GetTimePtr(cdb.GetCurTime())},
		{ObjectID: uuid.New(), Deleted: cdb.GetTimePtr(cdb.GetCurTime())},
	}

	var lifecycleMetricsManager subnetActivity.ManageSubnetLifecycleMetrics

	// Mock CollectAndRecordSubnetOperationMetrics activity failure
	s.env.RegisterActivity(lifecycleMetricsManager.RecordSubnetStatusTransitionMetrics)
	s.env.OnActivity(lifecycleMetricsManager.RecordSubnetStatusTransitionMetrics, mock.Anything, siteID, objectLifecycleEvents).Return(errors.New("RecordSubnetStatusTransitionMetrics Failure"))

	// Execute CollectAndRecordSubnetOperationMetrics workflow
	s.env.ExecuteWorkflow(RecordSubnetLifecycleMetrics, siteID, objectLifecycleEvents)
	s.True(s.env.IsWorkflowCompleted())
	err := s.env.GetWorkflowError()
	s.Error(err)

	var applicationErr *temporal.ApplicationError
	s.True(errors.As(err, &applicationErr))
	s.Equal("RecordSubnetStatusTransitionMetrics Failure", applicationErr.Error())
}

func TestRecordSubnetLifecycleMetricsSuite(t *testing.T) {
	suite.Run(t, new(RecordSubnetLifecycleMetricsTestSuite))
}
