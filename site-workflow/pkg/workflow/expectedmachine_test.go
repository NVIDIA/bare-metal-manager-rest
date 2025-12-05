// SPDX-FileCopyrightText: Copyright (c) 2021-2023 NVIDIA CORPORATION & AFFILIATES. All rights reserved.
// SPDX-License-Identifier: LicenseRef-NvidiaProprietary
//
// NVIDIA CORPORATION, its affiliates and licensors retain all intellectual
// property and proprietary rights in and to this material, related
// documentation and any modifications thereto. Any use, reproduction,
// disclosure or distribution of this material and related documentation
// without an express license agreement from NVIDIA CORPORATION or
// its affiliates is strictly prohibited.

package workflow

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	cwssaws "github.com/nvidia/carbide-rest/workflow-schema/schema/site-agent/workflows/v1"
	iActivity "github.com/nvidia/carbide-rest/site-workflow/pkg/activity"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/testsuite"
)

type InventoryExpectedMachineTestSuite struct {
	suite.Suite
	testsuite.WorkflowTestSuite

	env *testsuite.TestWorkflowEnvironment
}

func (s *InventoryExpectedMachineTestSuite) SetupTest() {
	s.env = s.NewTestWorkflowEnvironment()
}

func (s *InventoryExpectedMachineTestSuite) AfterTest(suiteName, testName string) {
	s.env.AssertExpectations(s.T())
}

func (s *InventoryExpectedMachineTestSuite) Test_DiscoverExpectedMachineInventory_Success() {
	var inventoryManager iActivity.ManageExpectedMachineInventory

	s.env.RegisterActivity(inventoryManager.DiscoverExpectedMachineInventory)
	s.env.OnActivity(inventoryManager.DiscoverExpectedMachineInventory, mock.Anything).Return(nil)

	// execute workflow
	s.env.ExecuteWorkflow(DiscoverExpectedMachineInventory)
	s.True(s.env.IsWorkflowCompleted())
	s.NoError(s.env.GetWorkflowError())
}

func (s *InventoryExpectedMachineTestSuite) Test_DiscoverExpectedMachineInventory_ActivityFails() {
	var inventoryManager iActivity.ManageExpectedMachineInventory

	errMsg := "Site Controller communication error"

	s.env.RegisterActivity(inventoryManager.DiscoverExpectedMachineInventory)
	s.env.OnActivity(inventoryManager.DiscoverExpectedMachineInventory, mock.Anything).Return(errors.New(errMsg))

	// Execute workflow
	s.env.ExecuteWorkflow(DiscoverExpectedMachineInventory)
	s.True(s.env.IsWorkflowCompleted())
	err := s.env.GetWorkflowError()
	s.Error(err)

	var applicationErr *temporal.ApplicationError
	s.True(errors.As(err, &applicationErr))
	s.Equal(errMsg, applicationErr.Error())
}

func TestInventoryExpectedMachineTestSuite(t *testing.T) {
	suite.Run(t, new(InventoryExpectedMachineTestSuite))
}

type CreateExpectedMachineTestSuite struct {
	suite.Suite
	testsuite.WorkflowTestSuite

	env *testsuite.TestWorkflowEnvironment
}

func (cemts *CreateExpectedMachineTestSuite) SetupTest() {
	cemts.env = cemts.NewTestWorkflowEnvironment()
}

func (cemts *CreateExpectedMachineTestSuite) AfterTest(suiteName, testName string) {
	cemts.env.AssertExpectations(cemts.T())
}

func (cemts *CreateExpectedMachineTestSuite) Test_CreateExpectedMachine_Success() {
	var expectedMachineManager iActivity.ManageExpectedMachine

	request := &cwssaws.ExpectedMachine{
		Id:            &cwssaws.UUID{Value: "test-create-workflow-001"},
		BmcMacAddress: "00:11:22:33:44:55",
	}

	// Mock CreateExpectedMachineOnSite activity
	cemts.env.RegisterActivity(expectedMachineManager.CreateExpectedMachineOnSite)
	cemts.env.OnActivity(expectedMachineManager.CreateExpectedMachineOnSite, mock.Anything, mock.Anything).Return(nil)

	// Execute CreateExpectedMachine workflow
	cemts.env.ExecuteWorkflow(CreateExpectedMachine, request)
	cemts.True(cemts.env.IsWorkflowCompleted())
	cemts.NoError(cemts.env.GetWorkflowError())
}

func (cemts *CreateExpectedMachineTestSuite) Test_CreateExpectedMachine_Failure() {
	var expectedMachineManager iActivity.ManageExpectedMachine

	request := &cwssaws.ExpectedMachine{
		Id:            &cwssaws.UUID{Value: "test-create-workflow-001"},
		BmcMacAddress: "00:11:22:33:44:55",
	}

	errMsg := "Site Controller communication error"

	// Mock CreateExpectedMachineOnSite activity
	cemts.env.RegisterActivity(expectedMachineManager.CreateExpectedMachineOnSite)
	cemts.env.OnActivity(expectedMachineManager.CreateExpectedMachineOnSite, mock.Anything, mock.Anything).Return(errors.New(errMsg))

	// execute CreateExpectedMachine workflow
	cemts.env.ExecuteWorkflow(CreateExpectedMachine, request)
	cemts.True(cemts.env.IsWorkflowCompleted())
	cemts.Error(cemts.env.GetWorkflowError())
}

func TestCreateExpectedMachineTestSuite(t *testing.T) {
	suite.Run(t, new(CreateExpectedMachineTestSuite))
}

type UpdateExpectedMachineTestSuite struct {
	suite.Suite
	testsuite.WorkflowTestSuite

	env *testsuite.TestWorkflowEnvironment
}

func (uemts *UpdateExpectedMachineTestSuite) SetupTest() {
	uemts.env = uemts.NewTestWorkflowEnvironment()
}

func (uemts *UpdateExpectedMachineTestSuite) AfterTest(suiteName, testName string) {
	uemts.env.AssertExpectations(uemts.T())
}

func (uemts *UpdateExpectedMachineTestSuite) Test_UpdateExpectedMachine_Success() {
	var expectedMachineManager iActivity.ManageExpectedMachine

	request := &cwssaws.ExpectedMachine{
		Id:            &cwssaws.UUID{Value: "test-create-workflow-001"},
		BmcMacAddress: "00:11:22:33:44:55",
	}

	// Mock UpdateExpectedMachineOnSite activity
	uemts.env.RegisterActivity(expectedMachineManager.UpdateExpectedMachineOnSite)
	uemts.env.OnActivity(expectedMachineManager.UpdateExpectedMachineOnSite, mock.Anything, mock.Anything).Return(nil)

	// Execute workflow
	uemts.env.ExecuteWorkflow(UpdateExpectedMachine, request)
	uemts.True(uemts.env.IsWorkflowCompleted())
	uemts.NoError(uemts.env.GetWorkflowError())
}

func (uemts *UpdateExpectedMachineTestSuite) Test_UpdateExpectedMachine_Failure() {
	var expectedMachineManager iActivity.ManageExpectedMachine

	request := &cwssaws.ExpectedMachine{
		Id:            &cwssaws.UUID{Value: "test-create-workflow-001"},
		BmcMacAddress: "00:11:22:33:44:55",
	}

	errMsg := "Site Controller communication error"

	// Mock UpdateExpectedMachineOnSite activity
	uemts.env.RegisterActivity(expectedMachineManager.UpdateExpectedMachineOnSite)
	uemts.env.OnActivity(expectedMachineManager.UpdateExpectedMachineOnSite, mock.Anything, mock.Anything).Return(errors.New(errMsg))

	// execute UpdateExpectedMachine workflow
	uemts.env.ExecuteWorkflow(UpdateExpectedMachine, request)
	uemts.True(uemts.env.IsWorkflowCompleted())
	uemts.Error(uemts.env.GetWorkflowError())
}

func TestUpdateExpectedMachineTestSuite(t *testing.T) {
	suite.Run(t, new(UpdateExpectedMachineTestSuite))
}

type DeleteExpectedMachineTestSuite struct {
	suite.Suite
	testsuite.WorkflowTestSuite

	env *testsuite.TestWorkflowEnvironment
}

func (demts *DeleteExpectedMachineTestSuite) SetupTest() {
	demts.env = demts.NewTestWorkflowEnvironment()
}

func (demts *DeleteExpectedMachineTestSuite) AfterTest(suiteName, testName string) {
	demts.env.AssertExpectations(demts.T())
}

func (demts *DeleteExpectedMachineTestSuite) Test_DeleteExpectedMachine_Success() {
	var expectedMachineManager iActivity.ManageExpectedMachine

	request := &cwssaws.ExpectedMachineRequest{
		Id:            &cwssaws.UUID{Value: "test-delete-workflow-001"},
		BmcMacAddress: "00:11:22:33:44:55",
	}

	// Mock DeleteExpectedMachineOnSite activity
	demts.env.RegisterActivity(expectedMachineManager.DeleteExpectedMachineOnSite)
	demts.env.OnActivity(expectedMachineManager.DeleteExpectedMachineOnSite, mock.Anything, mock.Anything).Return(nil)

	// execute workflow
	demts.env.ExecuteWorkflow(DeleteExpectedMachine, request)
	demts.True(demts.env.IsWorkflowCompleted())
	demts.NoError(demts.env.GetWorkflowError())
}

func (demts *DeleteExpectedMachineTestSuite) Test_DeleteExpectedMachine_Failure() {
	var expectedMachineManager iActivity.ManageExpectedMachine

	request := &cwssaws.ExpectedMachineRequest{
		Id:            &cwssaws.UUID{Value: "test-delete-workflow-001"},
		BmcMacAddress: "00:11:22:33:44:55",
	}

	errMsg := "Site Controller communication error"

	// Mock DeleteExpectedMachineOnSite activity
	demts.env.RegisterActivity(expectedMachineManager.DeleteExpectedMachineOnSite)
	demts.env.OnActivity(expectedMachineManager.DeleteExpectedMachineOnSite, mock.Anything, mock.Anything).Return(errors.New(errMsg))

	// execute DeleteExpectedMachine workflow
	demts.env.ExecuteWorkflow(DeleteExpectedMachine, request)
	demts.True(demts.env.IsWorkflowCompleted())
	demts.Error(demts.env.GetWorkflowError())
}

func TestDeleteExpectedMachineTestSuite(t *testing.T) {
	suite.Run(t, new(DeleteExpectedMachineTestSuite))
}
