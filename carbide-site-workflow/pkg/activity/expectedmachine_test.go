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

package activity

import (
	"context"
	"testing"

	"github.com/google/uuid"
	cwssaws "github.com/NVIDIA/carbide-rest-api/carbide-rest-api-schema/schema/site-agent/workflows/v1"
	cClient "github.com/NVIDIA/carbide-rest-api/carbide-site-workflow/pkg/grpc/client"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	tmocks "go.temporal.io/sdk/mocks"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestManageExpectedMachineInventory_DiscoverExpectedMachineInventory(t *testing.T) {
	mockCarbide := cClient.NewMockCarbideClient()

	carbideAtomicClient := cClient.NewCarbideAtomicClient(&cClient.CarbideClientConfig{})
	carbideAtomicClient.SwapClient(mockCarbide)

	wid := "test-workflow-id"
	wrun := &tmocks.WorkflowRun{}
	wrun.On("GetID").Return(wid)

	type fields struct {
		siteID               uuid.UUID
		carbideAtomicClient  *cClient.CarbideAtomicClient
		temporalPublishQueue string
		sitePageSize         int
		cloudPageSize        int
	}
	type args struct {
		wantTotalItems int
		findIDsError   error
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "test collecting and publishing expected machine inventory, empty inventory",
			fields: fields{
				siteID:               uuid.New(),
				carbideAtomicClient:  carbideAtomicClient,
				temporalPublishQueue: "test-queue",
				sitePageSize:         100,
				cloudPageSize:        25,
			},
			args: args{
				wantTotalItems: 0,
			},
		},
		{
			name: "test collecting and publishing expected machine inventory, normal inventory",
			fields: fields{
				siteID:               uuid.New(),
				carbideAtomicClient:  carbideAtomicClient,
				temporalPublishQueue: "test-queue",
				sitePageSize:         100,
				cloudPageSize:        25,
			},
			args: args{
				wantTotalItems: 195,
			},
		},
		{
			name: "test collecting and publishing expected machine inventory fallback, empty inventory",
			fields: fields{
				siteID:               uuid.New(),
				carbideAtomicClient:  carbideAtomicClient,
				temporalPublishQueue: "test-queue",
				sitePageSize:         100,
				cloudPageSize:        25,
			},
			args: args{
				wantTotalItems: 0,
				findIDsError:   status.Error(codes.Unimplemented, "not implemented"),
			},
		},
		{
			name: "test collecting and publishing expected machine inventory fallback, normal inventory",
			fields: fields{
				siteID:               uuid.New(),
				carbideAtomicClient:  carbideAtomicClient,
				temporalPublishQueue: "test-queue",
				sitePageSize:         100,
				cloudPageSize:        25,
			},
			args: args{
				wantTotalItems: 195,
				findIDsError:   status.Error(codes.Unimplemented, "not implemented"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tc := &tmocks.Client{}
			tc.Mock.On("ExecuteWorkflow", mock.Anything, mock.AnythingOfType("internal.StartWorkflowOptions"),
				mock.AnythingOfType("string"), mock.AnythingOfType("uuid.UUID"), mock.Anything).Return(wrun, nil)
			tc.AssertNumberOfCalls(t, "ExecuteWorkflow", 0)

			manageInstance := NewManageExpectedMachineInventory(ManageInventoryConfig{
				SiteID:                tt.fields.siteID,
				CarbideAtomicClient:   tt.fields.carbideAtomicClient,
				TemporalPublishClient: tc,
				TemporalPublishQueue:  tt.fields.temporalPublishQueue,
				SitePageSize:          tt.fields.sitePageSize,
				CloudPageSize:         tt.fields.cloudPageSize,
			})

			ctx := context.Background()
			ctx = context.WithValue(ctx, "wantCount", tt.args.wantTotalItems)
			if tt.args.findIDsError != nil {
				ctx = context.WithValue(ctx, "wantError", tt.args.findIDsError)
			}

			totalPages := tt.args.wantTotalItems / tt.fields.cloudPageSize
			if tt.args.wantTotalItems%tt.fields.cloudPageSize > 0 {
				totalPages++
			}

			err := manageInstance.DiscoverExpectedMachineInventory(ctx)
			assert.NoError(t, err)

			if tt.args.wantTotalItems == 0 {
				tc.AssertNumberOfCalls(t, "ExecuteWorkflow", 1)
			} else {
				tc.AssertNumberOfCalls(t, "ExecuteWorkflow", totalPages)
			}

			inventory, ok := tc.Calls[0].Arguments[4].(*cwssaws.ExpectedMachineInventory)
			assert.True(t, ok)

			if tt.args.wantTotalItems == 0 {
				assert.Equal(t, 0, len(inventory.ExpectedMachines))
			} else {
				assert.Equal(t, tt.fields.cloudPageSize, len(inventory.ExpectedMachines))
			}

			assert.Equal(t, cwssaws.InventoryStatus_INVENTORY_STATUS_SUCCESS, inventory.InventoryStatus)
			assert.Equal(t, totalPages, int(inventory.InventoryPage.TotalPages))
			assert.Equal(t, 1, int(inventory.InventoryPage.CurrentPage))
			assert.Equal(t, tt.fields.cloudPageSize, int(inventory.InventoryPage.PageSize))
			assert.Equal(t, tt.args.wantTotalItems, int(inventory.InventoryPage.TotalItems))
			assert.Equal(t, tt.args.wantTotalItems, len(inventory.InventoryPage.ItemIds))
		})
	}
}

func TestManageExpectedMachine_CreateExpectedMachineOnSite(t *testing.T) {
	mockCarbide := cClient.NewMockCarbideClient()

	carbideAtomicClient := cClient.NewCarbideAtomicClient(&cClient.CarbideClientConfig{})
	carbideAtomicClient.SwapClient(mockCarbide)

	type fields struct {
		CarbideAtomicClient *cClient.CarbideAtomicClient
	}
	type args struct {
		ctx     context.Context
		request *cwssaws.ExpectedMachine
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "test create expected machine success",
			fields: fields{
				CarbideAtomicClient: carbideAtomicClient,
			},
			args: args{
				ctx: context.Background(),
				request: &cwssaws.ExpectedMachine{
					Id:                  &cwssaws.UUID{Value: "test-machine-001"},
					BmcMacAddress:       "00:11:22:33:44:55",
					ChassisSerialNumber: "SN123456789",
				},
			},
			wantErr: false,
		},
		{
			name: "test create expected machine fail on missing MAC address",
			fields: fields{
				CarbideAtomicClient: carbideAtomicClient,
			},
			args: args{
				ctx: context.Background(),
				request: &cwssaws.ExpectedMachine{
					Id:                  &cwssaws.UUID{Value: "test-machine-002"},
					BmcMacAddress:       "",
					ChassisSerialNumber: "SN123456789",
				},
			},
			wantErr: true, // This should fail since MAC address is missing (now required)
		},
		{
			name: "test create expected machine fail on missing serial number",
			fields: fields{
				CarbideAtomicClient: carbideAtomicClient,
			},
			args: args{
				ctx: context.Background(),
				request: &cwssaws.ExpectedMachine{
					Id:                  &cwssaws.UUID{Value: "test-machine-003"},
					BmcMacAddress:       "00:11:22:33:44:55",
					ChassisSerialNumber: "",
				},
			},
			wantErr: true, // This should fail since serial number is missing (now required)
		},
		{
			name: "test create expected machine fail on missing id",
			fields: fields{
				CarbideAtomicClient: carbideAtomicClient,
			},
			args: args{
				ctx: context.Background(),
				request: &cwssaws.ExpectedMachine{
					Id:                  nil,
					BmcMacAddress:       "00:11:22:33:44:55",
					ChassisSerialNumber: "SN123456789",
				},
			},
			wantErr: true,
		},
		{
			name: "test create expected machine fail on missing identifying information",
			fields: fields{
				CarbideAtomicClient: carbideAtomicClient,
			},
			args: args{
				ctx: context.Background(),
				request: &cwssaws.ExpectedMachine{
					Id:                  &cwssaws.UUID{Value: "test-machine-004"},
					BmcMacAddress:       "",
					ChassisSerialNumber: "",
				},
			},
			wantErr: true,
		},
		{
			name: "test create expected machine fail on missing request",
			fields: fields{
				CarbideAtomicClient: carbideAtomicClient,
			},
			args: args{
				ctx:     context.Background(),
				request: nil,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mm := NewManageExpectedMachine(tt.fields.CarbideAtomicClient)
			err := mm.CreateExpectedMachineOnSite(tt.args.ctx, tt.args.request)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestManageExpectedMachine_UpdateExpectedMachineOnSite(t *testing.T) {
	mockCarbide := cClient.NewMockCarbideClient()

	carbideAtomicClient := cClient.NewCarbideAtomicClient(&cClient.CarbideClientConfig{})
	carbideAtomicClient.SwapClient(mockCarbide)

	type fields struct {
		CarbideAtomicClient *cClient.CarbideAtomicClient
	}
	type args struct {
		ctx     context.Context
		request *cwssaws.ExpectedMachine
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "test update expected machine success",
			fields: fields{
				CarbideAtomicClient: carbideAtomicClient,
			},
			args: args{
				ctx: context.Background(),
				request: &cwssaws.ExpectedMachine{
					Id:                  &cwssaws.UUID{Value: "test-update-001"},
					BmcMacAddress:       "00:11:22:33:44:55",
					ChassisSerialNumber: "SN123456789",
				},
			},
			wantErr: false,
		},
		{
			name: "test update expected machine fail on missing id",
			fields: fields{
				CarbideAtomicClient: carbideAtomicClient,
			},
			args: args{
				ctx: context.Background(),
				request: &cwssaws.ExpectedMachine{
					Id:                  nil,
					BmcMacAddress:       "00:11:22:33:44:55",
					ChassisSerialNumber: "SN123456789",
				},
			},
			wantErr: true,
		},
		{
			name: "test update expected machine fail on missing MAC address",
			fields: fields{
				CarbideAtomicClient: carbideAtomicClient,
			},
			args: args{
				ctx: context.Background(),
				request: &cwssaws.ExpectedMachine{
					Id:                  &cwssaws.UUID{Value: "test-update-002"},
					BmcMacAddress:       "",
					ChassisSerialNumber: "SN123456789",
				},
			},
			wantErr: true, // This should fail since MAC address is missing (now required)
		},
		{
			name: "test update expected machine fail on missing serial number",
			fields: fields{
				CarbideAtomicClient: carbideAtomicClient,
			},
			args: args{
				ctx: context.Background(),
				request: &cwssaws.ExpectedMachine{
					Id:                  &cwssaws.UUID{Value: "test-update-003"},
					BmcMacAddress:       "00:11:22:33:44:55",
					ChassisSerialNumber: "",
				},
			},
			wantErr: true, // This should fail since serial number is missing (now required)
		},
		{
			name: "test update expected machine fail on missing both MAC and serial",
			fields: fields{
				CarbideAtomicClient: carbideAtomicClient,
			},
			args: args{
				ctx: context.Background(),
				request: &cwssaws.ExpectedMachine{
					Id:                  &cwssaws.UUID{Value: "test-update-004"},
					BmcMacAddress:       "",
					ChassisSerialNumber: "",
				},
			},
			wantErr: true, // This should fail since both MAC address and serial number are missing
		},
		{
			name: "test update expected machine fail on missing request",
			fields: fields{
				CarbideAtomicClient: carbideAtomicClient,
			},
			args: args{
				ctx:     context.Background(),
				request: nil,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mm := NewManageExpectedMachine(tt.fields.CarbideAtomicClient)
			err := mm.UpdateExpectedMachineOnSite(tt.args.ctx, tt.args.request)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestManageExpectedMachine_DeleteExpectedMachineOnSite(t *testing.T) {
	mockCarbide := cClient.NewMockCarbideClient()

	carbideAtomicClient := cClient.NewCarbideAtomicClient(&cClient.CarbideClientConfig{})
	carbideAtomicClient.SwapClient(mockCarbide)

	type fields struct {
		CarbideAtomicClient *cClient.CarbideAtomicClient
	}
	type args struct {
		ctx     context.Context
		request *cwssaws.ExpectedMachineRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "test delete expected machine success",
			fields: fields{
				CarbideAtomicClient: carbideAtomicClient,
			},
			args: args{
				ctx: context.Background(),
				request: &cwssaws.ExpectedMachineRequest{
					Id:            &cwssaws.UUID{Value: "test-delete-001"},
					BmcMacAddress: "00:11:22:33:44:55",
				},
			},
			wantErr: false,
		},
		{
			name: "test delete expected machine fail on missing id",
			fields: fields{
				CarbideAtomicClient: carbideAtomicClient,
			},
			args: args{
				ctx: context.Background(),
				request: &cwssaws.ExpectedMachineRequest{
					Id:            nil,
					BmcMacAddress: "00:11:22:33:44:55",
				},
			},
			wantErr: true,
		},
		{
			name: "test delete expected machine success with missing BMC MAC address",
			fields: fields{
				CarbideAtomicClient: carbideAtomicClient,
			},
			args: args{
				ctx: context.Background(),
				request: &cwssaws.ExpectedMachineRequest{
					Id:            &cwssaws.UUID{Value: "test-delete-002"},
					BmcMacAddress: "",
				},
			},
			wantErr: false, // MAC address is no longer required, only ID is needed
		},
		{
			name: "test delete expected machine fail on missing request",
			fields: fields{
				CarbideAtomicClient: carbideAtomicClient,
			},
			args: args{
				ctx:     context.Background(),
				request: nil,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mm := NewManageExpectedMachine(tt.fields.CarbideAtomicClient)
			err := mm.DeleteExpectedMachineOnSite(tt.args.ctx, tt.args.request)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
