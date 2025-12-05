// SPDX-FileCopyrightText: Copyright (c) 2021-2023 NVIDIA CORPORATION & AFFILIATES. All rights reserved.
// SPDX-License-Identifier: LicenseRef-NvidiaProprietary
//
// NVIDIA CORPORATION, its affiliates and licensors retain all intellectual
// property and proprietary rights in and to this material, related
// documentation and any modifications thereto. Any use, reproduction,
// disclosure or distribution of this material and related documentation
// without an express license agreement from NVIDIA CORPORATION or
// its affiliates is strictly prohibited.

package tenant

import (
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	tenantActivity "github.com/nvidia/carbide-rest/workflow/pkg/activity/tenant"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/testsuite"
	"google.golang.org/protobuf/types/known/timestamppb"

	cwssaws "github.com/nvidia/carbide-rest/workflow-schema/schema/site-agent/workflows/v1"
)

type UpdateTenantTestSuite struct {
	suite.Suite
	testsuite.WorkflowTestSuite

	env *testsuite.TestWorkflowEnvironment
}

func (s *UpdateTenantTestSuite) SetupTest() {
	s.env = s.NewTestWorkflowEnvironment()
}

func (s *UpdateTenantTestSuite) AfterTest(suiteName, testName string) {
	s.env.AssertExpectations(s.T())
}

func (s *UpdateTenantTestSuite) Test_UpdateTenantInventory_Success() {
	var tenantManager tenantActivity.ManageTenant

	siteID := uuid.New()

	tenantInventory := &cwssaws.TenantInventory{
		Tenants:   []*cwssaws.Tenant{},
		Timestamp: timestamppb.Now(),
	}

	// Mock UpdateTenantsInDB activity
	s.env.RegisterActivity(tenantManager.UpdateTenantsInDB)
	s.env.OnActivity(tenantManager.UpdateTenantsInDB, mock.Anything, mock.Anything, mock.Anything).Return(nil)

	// execute UpdateTenantInventory workflow
	s.env.ExecuteWorkflow(UpdateTenantInventory, siteID.String(), tenantInventory)
	s.True(s.env.IsWorkflowCompleted())
	s.NoError(s.env.GetWorkflowError())
}

func (s *UpdateTenantTestSuite) Test_UpdateTenantInventory_ActivityFails() {
	var tenantManager tenantActivity.ManageTenant

	siteID := uuid.New()

	tenantInventory := &cwssaws.TenantInventory{
		Tenants:   []*cwssaws.Tenant{},
		Timestamp: timestamppb.Now(),
	}

	// Mock UpdateTenantsInDB activity
	s.env.RegisterActivity(tenantManager.UpdateTenantsInDB)
	s.env.OnActivity(tenantManager.UpdateTenantsInDB, mock.Anything, mock.Anything, mock.Anything).Return(errors.New("UpdateTenantInventory Failure"))

	// execute UpdateTenantInventory workflow
	s.env.ExecuteWorkflow(UpdateTenantInventory, siteID.String(), tenantInventory)
	s.True(s.env.IsWorkflowCompleted())
	err := s.env.GetWorkflowError()
	s.Error(err)

	var applicationErr *temporal.ApplicationError
	s.True(errors.As(err, &applicationErr))
	s.Equal("UpdateTenantInventory Failure", applicationErr.Error())
}

func TestUpdateTenantInfoSuite(t *testing.T) {
	suite.Run(t, new(UpdateTenantTestSuite))
}
