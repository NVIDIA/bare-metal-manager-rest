// SPDX-FileCopyrightText: Copyright (c) 2021-2023 NVIDIA CORPORATION & AFFILIATES. All rights reserved.
// SPDX-License-Identifier: LicenseRef-NvidiaProprietary
//
// NVIDIA CORPORATION, its affiliates and licensors retain all intellectual
// property and proprietary rights in and to this material, related
// documentation and any modifications thereto. Any use, reproduction,
// disclosure or distribution of this material and related documentation
// without an express license agreement from NVIDIA CORPORATION or
// its affiliates is strictly prohibited.

package sku

import (
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	skuActivity "github.com/nvidia/carbide-rest/workflow/pkg/activity/sku"

	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/testsuite"

	cwssaws "github.com/nvidia/carbide-rest/workflow-schema/schema/site-agent/workflows/v1"
)

type UpdateSkuTestSuite struct {
	suite.Suite
	testsuite.WorkflowTestSuite

	env *testsuite.TestWorkflowEnvironment
}

func (s *UpdateSkuTestSuite) SetupTest() {
	s.env = s.NewTestWorkflowEnvironment()
}

func (s *UpdateSkuTestSuite) AfterTest(suiteName, testName string) {
	s.env.AssertExpectations(s.T())
}

func (s *UpdateSkuTestSuite) Test_UpdateSkuInventory_Success() {
	var skuManager skuActivity.ManageSku

	siteID := uuid.New()

	inv := &cwssaws.SkuInventory{Skus: []*cwssaws.Sku{}}

	s.env.RegisterActivity(skuManager.UpdateSkusInDB)
	s.env.OnActivity(skuManager.UpdateSkusInDB, mock.Anything, mock.Anything, mock.Anything).Return(nil)

	s.env.ExecuteWorkflow(UpdateSkuInventory, siteID.String(), inv)
	s.True(s.env.IsWorkflowCompleted())
	s.NoError(s.env.GetWorkflowError())
}

func (s *UpdateSkuTestSuite) Test_UpdateSkuInventory_ActivityFails() {
	var skuManager skuActivity.ManageSku

	siteID := uuid.New()

	inv := &cwssaws.SkuInventory{Skus: []*cwssaws.Sku{}}

	s.env.RegisterActivity(skuManager.UpdateSkusInDB)
	s.env.OnActivity(skuManager.UpdateSkusInDB, mock.Anything, mock.Anything, mock.Anything).Return(errors.New("UpdateSkuInventory Failure"))

	s.env.ExecuteWorkflow(UpdateSkuInventory, siteID.String(), inv)
	s.True(s.env.IsWorkflowCompleted())
	err := s.env.GetWorkflowError()
	s.NotNil(err)

	var applicationErr *temporal.ApplicationError
	s.True(errors.As(err, &applicationErr))
	s.Equal("UpdateSkuInventory Failure", applicationErr.Error())
}

func TestUpdateSkuTestSuite(t *testing.T) {
	suite.Run(t, new(UpdateSkuTestSuite))
}

