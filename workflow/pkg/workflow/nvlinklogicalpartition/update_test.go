// SPDX-FileCopyrightText: Copyright (c) 2021-2023 NVIDIA CORPORATION & AFFILIATES. All rights reserved.
// SPDX-License-Identifier: LicenseRef-NvidiaProprietary
//
// NVIDIA CORPORATION, its affiliates and licensors retain all intellectual
// property and proprietary rights in and to this material, related
// documentation and any modifications thereto. Any use, reproduction,
// disclosure or distribution of this material and related documentation
// without an express license agreement from NVIDIA CORPORATION or
// its affiliates is strictly prohibited.

package nvlinklogicalpartition

import (
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	nvlinklogicalpartitionActivity "github.com/nvidia/carbide-rest/workflow/pkg/activity/nvlinklogicalpartition"

	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/testsuite"

	cwssaws "github.com/nvidia/carbide-rest/workflow-schema/schema/site-agent/workflows/v1"
)

type UpdateNVLinkLogicalPartitionTestSuite struct {
	suite.Suite
	testsuite.WorkflowTestSuite

	env *testsuite.TestWorkflowEnvironment
}

func (s *UpdateNVLinkLogicalPartitionTestSuite) SetupTest() {
	s.env = s.NewTestWorkflowEnvironment()
}

func (s *UpdateNVLinkLogicalPartitionTestSuite) AfterTest(suiteName, testName string) {
	s.env.AssertExpectations(s.T())
}

func (s *UpdateNVLinkLogicalPartitionTestSuite) Test_UpdateNVLinkLogicalPartitionInventory_Success() {
	var nvlinklogicalpartitionManager nvlinklogicalpartitionActivity.ManageNVLinkLogicalPartition

	siteID := uuid.New()

	inv := &cwssaws.NVLinkLogicalPartitionInventory{
		Partitions: []*cwssaws.NVLinkLogicalPartition{},
	}

	s.env.RegisterActivity(nvlinklogicalpartitionManager.UpdateNVLinkLogicalPartitionsInDB)
	s.env.OnActivity(nvlinklogicalpartitionManager.UpdateNVLinkLogicalPartitionsInDB, mock.Anything, mock.Anything, mock.Anything).Return(nil)

	s.env.ExecuteWorkflow(UpdateNVLinkLogicalPartitionInventory, siteID.String(), inv)
	s.True(s.env.IsWorkflowCompleted())
	s.NoError(s.env.GetWorkflowError())
}

func (s *UpdateNVLinkLogicalPartitionTestSuite) Test_UpdateNVLinkLogicalPartitionInventory_ActivityFails() {
	var nvlinklogicalpartitionManager nvlinklogicalpartitionActivity.ManageNVLinkLogicalPartition

	siteID := uuid.New()

	inv := &cwssaws.NVLinkLogicalPartitionInventory{
		Partitions: []*cwssaws.NVLinkLogicalPartition{},
	}

	s.env.RegisterActivity(nvlinklogicalpartitionManager.UpdateNVLinkLogicalPartitionsInDB)
	s.env.OnActivity(nvlinklogicalpartitionManager.UpdateNVLinkLogicalPartitionsInDB, mock.Anything, mock.Anything, mock.Anything).Return(errors.New("UpdateNVLinkLogicalPartitionInventory Failure"))

	s.env.ExecuteWorkflow(UpdateNVLinkLogicalPartitionInventory, siteID.String(), inv)
	s.True(s.env.IsWorkflowCompleted())
	err := s.env.GetWorkflowError()
	s.NotNil(err)

	var applicationErr *temporal.ApplicationError
	s.True(errors.As(err, &applicationErr))
	s.Equal("UpdateNVLinkLogicalPartitionInventory Failure", applicationErr.Error())
}

func TestUpdateNVLinkLogicalPartitionTestSuite(t *testing.T) {
	suite.Run(t, new(UpdateNVLinkLogicalPartitionTestSuite))
}
