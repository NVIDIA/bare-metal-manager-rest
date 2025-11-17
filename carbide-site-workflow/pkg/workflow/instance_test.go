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

	iActivity "github.com/NVIDIA/carbide-rest-api/carbide-site-workflow/pkg/activity"
)

type UpdateInstanceTestSuite struct {
	suite.Suite
	testsuite.WorkflowTestSuite

	env *testsuite.TestWorkflowEnvironment
}

func (s *UpdateInstanceTestSuite) SetupTest() {
	s.env = s.NewTestWorkflowEnvironment()
}

func (s *UpdateInstanceTestSuite) AfterTest(suiteName, testName string) {
	s.env.AssertExpectations(s.T())
}

func (s *UpdateInstanceTestSuite) Test_UpdateInstance_Success() {
	var machineManager iActivity.ManageInstance

	ipxeScript := "#!ipxe"
	userData := "echo"
	labelKey := "key1"
	labelValue := "value1"

	request := &cwssaws.InstanceConfigUpdateRequest{
		InstanceId: &cwssaws.UUID{Value: uuid.NewString()},
		Metadata: &cwssaws.Metadata{
			Name:        "updated_name",
			Description: "updated_description",
			Labels: []*cwssaws.Label{
				{
					Key:   labelKey,
					Value: &labelValue,
				},
			},
		},
		Config: &cwssaws.InstanceConfig{
			Os: &cwssaws.OperatingSystem{
				RunProvisioningInstructionsOnEveryBoot: true,
				Variant: &cwssaws.OperatingSystem_Ipxe{
					Ipxe: &cwssaws.IpxeOperatingSystem{
						IpxeScript: ipxeScript,
						UserData:   &userData,
					},
				},
			},
		},
	}

	// Mock UpdateInstanceOnSiteActivity activity
	s.env.RegisterActivity(machineManager.UpdateInstanceOnSite)
	s.env.OnActivity(machineManager.UpdateInstanceOnSite, mock.Anything, mock.Anything).Return(nil)

	// execute UpdateMachineInventory workflow
	s.env.ExecuteWorkflow(UpdateInstance, request)
	s.True(s.env.IsWorkflowCompleted())
	s.NoError(s.env.GetWorkflowError())
}

func (s *UpdateInstanceTestSuite) Test_UpdateInstance_Failure() {
	var machineManager iActivity.ManageInstance

	ipxeScript := "#!ipxe"
	userData := "echo"
	labelKey := "key1"
	labelValue := "value1"

	request := &cwssaws.InstanceConfigUpdateRequest{
		InstanceId: &cwssaws.UUID{Value: uuid.NewString()},
		Metadata: &cwssaws.Metadata{
			Name:        "updated_name",
			Description: "updated_description",
			Labels: []*cwssaws.Label{
				{
					Key:   labelKey,
					Value: &labelValue,
				},
			},
		},
		Config: &cwssaws.InstanceConfig{
			Os: &cwssaws.OperatingSystem{
				RunProvisioningInstructionsOnEveryBoot: true,
				Variant: &cwssaws.OperatingSystem_Ipxe{
					Ipxe: &cwssaws.IpxeOperatingSystem{
						IpxeScript: ipxeScript,
						UserData:   &userData,
					},
				},
			},
		},
	}

	errMsg := "Site Controller communication error"

	// Mock UpdateInstanceOnSiteActivity activity
	s.env.RegisterActivity(machineManager.UpdateInstanceOnSite)
	s.env.OnActivity(machineManager.UpdateInstanceOnSite, mock.Anything, mock.Anything).Return(errors.New(errMsg))

	// execute UpdateMachineInventory workflow
	s.env.ExecuteWorkflow(UpdateInstance, request)
	s.True(s.env.IsWorkflowCompleted())
	s.Error(s.env.GetWorkflowError())
}

func TestUpdateInstanceTestSuite(t *testing.T) {
	suite.Run(t, new(UpdateInstanceTestSuite))
}

type CreateInstanceV2TestSuite struct {
	suite.Suite
	testsuite.WorkflowTestSuite

	env *testsuite.TestWorkflowEnvironment
}

func (s *CreateInstanceV2TestSuite) SetupTest() {
	s.env = s.NewTestWorkflowEnvironment()
}

func (s *CreateInstanceV2TestSuite) AfterTest(suiteName, testName string) {
	s.env.AssertExpectations(s.T())
}

func (s *CreateInstanceV2TestSuite) Test_CreateInstanceV2_Success() {
	var machineManager iActivity.ManageInstance

	ipxeScript := "#!ipxe"
	userData := "echo"
	labelKey := "key1"
	labelValue := "value1"

	request := &cwssaws.InstanceAllocationRequest{
		MachineId: &cwssaws.MachineId{Id: uuid.NewString()},
		Metadata: &cwssaws.Metadata{
			Name:        "updated_name",
			Description: "updated_description",
			Labels: []*cwssaws.Label{
				{
					Key:   labelKey,
					Value: &labelValue,
				},
			},
		},
		Config: &cwssaws.InstanceConfig{
			Os: &cwssaws.OperatingSystem{
				RunProvisioningInstructionsOnEveryBoot: true,
				Variant: &cwssaws.OperatingSystem_Ipxe{
					Ipxe: &cwssaws.IpxeOperatingSystem{
						IpxeScript: ipxeScript,
						UserData:   &userData,
					},
				},
			},
		},
	}

	// Mock CreateInstanceOnSiteActivity activity
	s.env.RegisterActivity(machineManager.CreateInstanceOnSite)
	s.env.OnActivity(machineManager.CreateInstanceOnSite, mock.Anything, mock.Anything).Return(nil)

	// execute workflow
	s.env.ExecuteWorkflow(CreateInstanceV2, request)
	s.True(s.env.IsWorkflowCompleted())
	s.NoError(s.env.GetWorkflowError())
}

func (s *CreateInstanceV2TestSuite) Test_CreateInstanceV2_Failure() {
	var machineManager iActivity.ManageInstance

	ipxeScript := "#!ipxe"
	userData := "echo"
	labelKey := "key1"
	labelValue := "value1"

	request := &cwssaws.InstanceAllocationRequest{
		MachineId: &cwssaws.MachineId{Id: uuid.NewString()},
		Metadata: &cwssaws.Metadata{
			Name:        "updated_name",
			Description: "updated_description",
			Labels: []*cwssaws.Label{
				{
					Key:   labelKey,
					Value: &labelValue,
				},
			},
		},
		Config: &cwssaws.InstanceConfig{
			Os: &cwssaws.OperatingSystem{
				RunProvisioningInstructionsOnEveryBoot: true,
				Variant: &cwssaws.OperatingSystem_Ipxe{
					Ipxe: &cwssaws.IpxeOperatingSystem{
						IpxeScript: ipxeScript,
						UserData:   &userData,
					},
				},
			},
		},
	}

	errMsg := "Site Controller communication error"

	// Mock CreateInstanceOnSiteActivity activity
	s.env.RegisterActivity(machineManager.CreateInstanceOnSite)
	s.env.OnActivity(machineManager.CreateInstanceOnSite, mock.Anything, mock.Anything).Return(errors.New(errMsg))

	// execute CreateMachineInventory workflow
	s.env.ExecuteWorkflow(CreateInstanceV2, request)
	s.True(s.env.IsWorkflowCompleted())
	s.Error(s.env.GetWorkflowError())
}

func TestCreateInstanceV2TestSuite(t *testing.T) {
	suite.Run(t, new(CreateInstanceV2TestSuite))
}

type DeleteInstanceV2TestSuite struct {
	suite.Suite
	testsuite.WorkflowTestSuite

	env *testsuite.TestWorkflowEnvironment
}

func (s *DeleteInstanceV2TestSuite) SetupTest() {
	s.env = s.NewTestWorkflowEnvironment()
}

func (s *DeleteInstanceV2TestSuite) AfterTest(suiteName, testName string) {
	s.env.AssertExpectations(s.T())
}

func (s *DeleteInstanceV2TestSuite) Test_DeleteInstanceV2_Success() {
	var instanceManager iActivity.ManageInstance

	request := &cwssaws.InstanceReleaseRequest{
		Id: &cwssaws.UUID{Value: uuid.NewString()},
	}

	// Mock DeleteInstanceOnSiteActivity activity
	s.env.RegisterActivity(instanceManager.DeleteInstanceOnSite)
	s.env.OnActivity(instanceManager.DeleteInstanceOnSite, mock.Anything, mock.Anything).Return(nil)

	// execute workflow
	s.env.ExecuteWorkflow(DeleteInstanceV2, request)
	s.True(s.env.IsWorkflowCompleted())
	s.NoError(s.env.GetWorkflowError())
}

func (s *DeleteInstanceV2TestSuite) Test_DeleteInstanceV2_Failure() {
	var machineManager iActivity.ManageInstance

	request := &cwssaws.InstanceReleaseRequest{
		Id: &cwssaws.UUID{Value: uuid.NewString()},
	}

	errMsg := "Site Controller communication error"

	// Mock DeleteInstanceOnSiteActivity activity
	s.env.RegisterActivity(machineManager.DeleteInstanceOnSite)
	s.env.OnActivity(machineManager.DeleteInstanceOnSite, mock.Anything, mock.Anything).Return(errors.New(errMsg))

	// execute DeleteMachineInventory workflow
	s.env.ExecuteWorkflow(DeleteInstanceV2, request)
	s.True(s.env.IsWorkflowCompleted())
	s.Error(s.env.GetWorkflowError())
}

func TestDeleteInstanceV2TestSuite(t *testing.T) {
	suite.Run(t, new(DeleteInstanceV2TestSuite))
}

type RebootInstanceTestSuite struct {
	suite.Suite
	testsuite.WorkflowTestSuite

	env *testsuite.TestWorkflowEnvironment
}

func (s *RebootInstanceTestSuite) SetupTest() {
	s.env = s.NewTestWorkflowEnvironment()
}

func (s *RebootInstanceTestSuite) AfterTest(suiteName, testName string) {
	s.env.AssertExpectations(s.T())
}

func (s *RebootInstanceTestSuite) Test_RebootInstance_Success() {
	var machineManager iActivity.ManageInstance

	request := &cwssaws.InstancePowerRequest{
		MachineId: &cwssaws.MachineId{Id: uuid.NewString()},
		Operation: cwssaws.InstancePowerRequest_POWER_RESET,
	}

	// Mock RebootInstanceOnSiteActivity activity
	s.env.RegisterActivity(machineManager.RebootInstanceOnSite)
	s.env.OnActivity(machineManager.RebootInstanceOnSite, mock.Anything, mock.Anything).Return(nil)

	// execute workflow
	s.env.ExecuteWorkflow(RebootInstance, request)
	s.True(s.env.IsWorkflowCompleted())
	s.NoError(s.env.GetWorkflowError())
}

func (s *RebootInstanceTestSuite) Test_RebootInstance_Failure() {
	var machineManager iActivity.ManageInstance

	request := &cwssaws.InstancePowerRequest{
		MachineId: &cwssaws.MachineId{Id: uuid.NewString()},
		Operation: cwssaws.InstancePowerRequest_POWER_RESET,
	}

	errMsg := "Site Controller communication error"

	// Mock RebootInstanceOnSiteActivity activity
	s.env.RegisterActivity(machineManager.RebootInstanceOnSite)
	s.env.OnActivity(machineManager.RebootInstanceOnSite, mock.Anything, mock.Anything).Return(errors.New(errMsg))

	// execute RebootMachineInventory workflow
	s.env.ExecuteWorkflow(RebootInstance, request)
	s.True(s.env.IsWorkflowCompleted())
	s.Error(s.env.GetWorkflowError())
}

func TestRebootInstanceTestSuite(t *testing.T) {
	suite.Run(t, new(RebootInstanceTestSuite))
}

type InventoryInstanceTestSuite struct {
	suite.Suite
	testsuite.WorkflowTestSuite

	env *testsuite.TestWorkflowEnvironment
}

func (s *InventoryInstanceTestSuite) SetupTest() {
	s.env = s.NewTestWorkflowEnvironment()
}

func (s *InventoryInstanceTestSuite) AfterTest(suiteName, testName string) {
	s.env.AssertExpectations(s.T())
}

func (s *InventoryInstanceTestSuite) Test_DiscoverInstanceInventory_Success() {
	var inventoryManager iActivity.ManageInstanceInventory

	s.env.RegisterActivity(inventoryManager.DiscoverInstanceInventory)
	s.env.OnActivity(inventoryManager.DiscoverInstanceInventory, mock.Anything).Return(nil)

	// execute workflow
	s.env.ExecuteWorkflow(DiscoverInstanceInventory)
	s.True(s.env.IsWorkflowCompleted())
	s.NoError(s.env.GetWorkflowError())
}

func (s *InventoryInstanceTestSuite) Test_DiscoverInstanceInventory_ActivityFails() {
	var inventoryManager iActivity.ManageInstanceInventory

	errMsg := "Site Controller communication error"

	s.env.RegisterActivity(inventoryManager.DiscoverInstanceInventory)
	s.env.OnActivity(inventoryManager.DiscoverInstanceInventory, mock.Anything).Return(errors.New(errMsg))

	// Execute workflow
	s.env.ExecuteWorkflow(DiscoverInstanceInventory)
	s.True(s.env.IsWorkflowCompleted())
	err := s.env.GetWorkflowError()
	s.Error(err)

	var applicationErr *temporal.ApplicationError
	s.True(errors.As(err, &applicationErr))
	s.Equal(errMsg, applicationErr.Error())
}

func TestInventoryInstanceTestSuite(t *testing.T) {
	suite.Run(t, new(InventoryInstanceTestSuite))
}
