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

package workflow

import (
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/testsuite"

	cwssaws "github.com/NVIDIA/carbide-rest-api/carbide-rest-api-schema/schema/site-agent/workflows/v1"

	mActivity "github.com/NVIDIA/carbide-rest-api/carbide-site-workflow/pkg/activity"
	"github.com/NVIDIA/carbide-rest-api/carbide-site-workflow/pkg/util"
)

type MachineWorkflowTestSuite struct {
	suite.Suite
	testsuite.WorkflowTestSuite

	env *testsuite.TestWorkflowEnvironment
}

func (s *MachineWorkflowTestSuite) SetupTest() {
	s.env = s.NewTestWorkflowEnvironment()
}

func (s *MachineWorkflowTestSuite) AfterTest(suiteName, testName string) {
	s.env.AssertExpectations(s.T())
}

func (s *MachineWorkflowTestSuite) Test_UpdateMachineInventory_Success() {
	var machineManager mActivity.ManageMachine

	request := &cwssaws.MaintenanceRequest{
		Operation: cwssaws.MaintenanceOperation_Enable,
		HostId:    &cwssaws.MachineId{Id: uuid.New().String()},
		Reference: util.GetStrPtr("Machine needs to taken offline to re-cable the network"),
	}

	// Mock UpdateVpcViaSiteAgent activity
	s.env.RegisterActivity(machineManager.SetMachineMaintenanceOnSite)
	s.env.OnActivity(machineManager.SetMachineMaintenanceOnSite, mock.Anything, mock.Anything).Return(nil)

	// execute UpdateMachineInventory workflow
	s.env.ExecuteWorkflow(SetMachineMaintenance, request)
	s.True(s.env.IsWorkflowCompleted())
	s.NoError(s.env.GetWorkflowError())
}

func (s *MachineWorkflowTestSuite) Test_UpdateMachineInventory_ActivityFails() {
	var machineManager mActivity.ManageMachine

	request := &cwssaws.MaintenanceRequest{
		Operation: cwssaws.MaintenanceOperation_Enable,
		HostId:    &cwssaws.MachineId{Id: uuid.New().String()},
		Reference: util.GetStrPtr("Machine needs to taken offline to re-cable the network"),
	}

	errMsg := "Site Controller communication error"

	// Mock SetMachineMaintenanceOnSite activity failure
	s.env.RegisterActivity(machineManager.SetMachineMaintenanceOnSite)
	s.env.OnActivity(machineManager.SetMachineMaintenanceOnSite, mock.Anything, mock.Anything).Return(errors.New(errMsg))

	// Execute SetMachineMaintenanceOnSite workflow
	s.env.ExecuteWorkflow(SetMachineMaintenance, request)
	s.True(s.env.IsWorkflowCompleted())
	err := s.env.GetWorkflowError()
	s.Error(err)

	var applicationErr *temporal.ApplicationError
	s.True(errors.As(err, &applicationErr))
	s.Equal(errMsg, applicationErr.Error())
}

func (s *MachineWorkflowTestSuite) Test_CollectAndPublishMachineInventory_Success() {
	var machineInventoryManager mActivity.ManageMachineInventory

	// Mock SetMachineMaintenanceOnSite activity failure
	s.env.RegisterActivity(machineInventoryManager.CollectAndPublishMachineInventory)
	s.env.OnActivity(machineInventoryManager.CollectAndPublishMachineInventory, mock.Anything).Return(nil)

	// execute UpdateMachineInventory workflow
	s.env.ExecuteWorkflow(CollectAndPublishMachineInventory)
	s.True(s.env.IsWorkflowCompleted())
	s.NoError(s.env.GetWorkflowError())
}

func (s *MachineWorkflowTestSuite) Test_CollectAndPublishMachineInventory_ActivityFails() {
	var machineInventoryManager mActivity.ManageMachineInventory

	errMsg := "Site Controller communication error"

	// Mock SetMachineMaintenanceOnSite activity failure
	s.env.RegisterActivity(machineInventoryManager.CollectAndPublishMachineInventory)
	s.env.OnActivity(machineInventoryManager.CollectAndPublishMachineInventory, mock.Anything).Return(errors.New(errMsg))

	// Execute SetMachineMaintenanceOnSite workflow
	s.env.ExecuteWorkflow(CollectAndPublishMachineInventory)
	s.True(s.env.IsWorkflowCompleted())
	err := s.env.GetWorkflowError()
	s.Error(err)

	var applicationErr *temporal.ApplicationError
	s.True(errors.As(err, &applicationErr))
	s.Equal(errMsg, applicationErr.Error())
}

func (s *MachineWorkflowTestSuite) Test_UpdateMachineMetadata_Success() {
	var machineManager mActivity.ManageMachine

	request := &cwssaws.MachineMetadataUpdateRequest{
		MachineId: &cwssaws.MachineId{Id: uuid.New().String()},
		Metadata: &cwssaws.Metadata{
			Labels: []*cwssaws.Label{
				{
					Key:   "test-key",
					Value: util.GetStrPtr("test-value"),
				},
			},
		},
	}

	// Mock UpdateMachineMetadataOnSite activity
	s.env.RegisterActivity(machineManager.UpdateMachineMetadataOnSite)
	s.env.OnActivity(machineManager.UpdateMachineMetadataOnSite, mock.Anything, mock.Anything).Return(nil)

	// execute UpdateMachineMetadata workflow
	s.env.ExecuteWorkflow(UpdateMachineMetadata, request)
	s.True(s.env.IsWorkflowCompleted())
	s.NoError(s.env.GetWorkflowError())
}

func (s *MachineWorkflowTestSuite) Test_UpdateMachineMetadata_ActivityFails() {
	var machineManager mActivity.ManageMachine

	errMsg := "Site Controller communication error"

	request := &cwssaws.MachineMetadataUpdateRequest{
		MachineId: &cwssaws.MachineId{Id: uuid.New().String()},
		Metadata: &cwssaws.Metadata{
			Labels: []*cwssaws.Label{
				{
					Key:   "test-key",
					Value: util.GetStrPtr("test-value"),
				},
			},
		},
	}

	// Mock UpdateMachineMetadataOnSite activity failure
	s.env.RegisterActivity(machineManager.UpdateMachineMetadataOnSite)
	s.env.OnActivity(machineManager.UpdateMachineMetadataOnSite, mock.Anything, mock.Anything).Return(errors.New(errMsg))

	// Execute UpdateMachineMetadata workflow
	s.env.ExecuteWorkflow(UpdateMachineMetadata, request)
	s.True(s.env.IsWorkflowCompleted())
	err := s.env.GetWorkflowError()
	s.Error(err)

	var applicationErr *temporal.ApplicationError
	s.True(errors.As(err, &applicationErr))
	s.Equal(errMsg, applicationErr.Error())
}

func TestMachineWorkflowSuite(t *testing.T) {
	suite.Run(t, new(MachineWorkflowTestSuite))
}
