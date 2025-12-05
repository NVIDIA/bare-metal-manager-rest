// SPDX-FileCopyrightText: Copyright (c) 2021-2023 NVIDIA CORPORATION & AFFILIATES. All rights reserved.
// SPDX-License-Identifier: LicenseRef-NvidiaProprietary
//
// NVIDIA CORPORATION, its affiliates and licensors retain all intellectual
// property and proprietary rights in and to this material, related
// documentation and any modifications thereto. Any use, reproduction,
// disclosure or distribution of this material and related documentation
// without an express license agreement from NVIDIA CORPORATION or
// its affiliates is strictly prohibited.

package activity

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	cwssaws "github.com/nvidia/carbide-rest/workflow-schema/schema/site-agent/workflows/v1"
	cClient "github.com/nvidia/carbide-rest/site-workflow/pkg/grpc/client"
	tmocks "go.temporal.io/sdk/mocks"
)

func TestManageNVLinkLogicalPartitionInventory_DiscoverNVLinkLogicalPartitionInventory(t *testing.T) {
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
			name: "test collecting and publishing nvlink logical partition inventory, empty inventory",
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
			name: "test collecting and publishing nvlink logical partition inventory, normal inventory",
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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tc := &tmocks.Client{}
			tc.Mock.On("ExecuteWorkflow", mock.Anything, mock.AnythingOfType("internal.StartWorkflowOptions"),
				mock.AnythingOfType("string"), mock.AnythingOfType("uuid.UUID"), mock.Anything).Return(wrun, nil)
			tc.AssertNumberOfCalls(t, "ExecuteWorkflow", 0)

			manageInstance := NewManageNVLinkLogicalPartitionInventory(ManageInventoryConfig{
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

			err := manageInstance.DiscoverNVLinkLogicalPartitionInventory(ctx)
			assert.NoError(t, err)

			if tt.args.wantTotalItems == 0 {
				tc.AssertNumberOfCalls(t, "ExecuteWorkflow", 1)
			} else {
				tc.AssertNumberOfCalls(t, "ExecuteWorkflow", totalPages)
			}

			inventory, ok := tc.Calls[0].Arguments[4].(*cwssaws.NVLinkLogicalPartitionInventory)
			assert.True(t, ok)

			if tt.args.wantTotalItems == 0 {
				assert.Equal(t, 0, len(inventory.Partitions))
			} else {
				assert.Equal(t, tt.fields.cloudPageSize, len(inventory.Partitions))
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

func TestManageNVLinkLogicalPartition_CreateNVLinkLogicalPartitionOnSite(t *testing.T) {
	mockCarbide := cClient.NewMockCarbideClient()

	carbideAtomicClient := cClient.NewCarbideAtomicClient(&cClient.CarbideClientConfig{})
	carbideAtomicClient.SwapClient(mockCarbide)

	type fields struct {
		carbideAtomicClient *cClient.CarbideAtomicClient
	}
	type args struct {
		ctx     context.Context
		request *cwssaws.NVLinkLogicalPartitionCreationRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "test create nvlink logical partition success",
			fields: fields{
				carbideAtomicClient: carbideAtomicClient,
			},
			args: args{
				ctx: context.Background(),
				request: &cwssaws.NVLinkLogicalPartitionCreationRequest{
					Id: &cwssaws.NVLinkLogicalPartitionId{Value: "b410867c-655a-11ef-bc4a-0393098e5d09"},
					Config: &cwssaws.NVLinkLogicalPartitionConfig{
						Metadata: &cwssaws.Metadata{
							Name: "test_name",
						},
						TenantOrganizationId: "test_org",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "test create nvlink logical partition fail on missing id",
			fields: fields{
				carbideAtomicClient: carbideAtomicClient,
			},
			args: args{
				ctx: context.Background(),
				request: &cwssaws.NVLinkLogicalPartitionCreationRequest{
					Id: nil,
					Config: &cwssaws.NVLinkLogicalPartitionConfig{
						Metadata: &cwssaws.Metadata{
							Name: "test_name",
						},
						TenantOrganizationId: "test_org",
					},
				},
			},
			wantErr: true,
		},
		{
			name: "test create nvlink logical partition fail on missing name",
			fields: fields{
				carbideAtomicClient: carbideAtomicClient,
			},
			args: args{
				ctx: context.Background(),
				request: &cwssaws.NVLinkLogicalPartitionCreationRequest{
					Id: &cwssaws.NVLinkLogicalPartitionId{Value: "b410867c-655a-11ef-bc4a-0393098e5d09"},
					Config: &cwssaws.NVLinkLogicalPartitionConfig{
						Metadata: &cwssaws.Metadata{
							Name: "",
						},
						TenantOrganizationId: "test_org",
					},
				},
			},
			wantErr: true,
		},
		{
			name: "test create nvlink logical partition fail on missing org",
			fields: fields{
				carbideAtomicClient: carbideAtomicClient,
			},
			args: args{
				ctx: context.Background(),
				request: &cwssaws.NVLinkLogicalPartitionCreationRequest{
					Id: &cwssaws.NVLinkLogicalPartitionId{Value: "b410867c-655a-11ef-bc4a-0393098e5d09"},
					Config: &cwssaws.NVLinkLogicalPartitionConfig{
						Metadata: &cwssaws.Metadata{
							Name: "test_name",
						},
						TenantOrganizationId: "",
					},
				},
			},
			wantErr: true,
		},
		{
			name: "test create nvlink logical partition fail on missing request",
			fields: fields{
				carbideAtomicClient: carbideAtomicClient,
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
			mnvllp := NewManageNVLinkLogicalPartition(tt.fields.carbideAtomicClient)
			err := mnvllp.CreateNVLinkLogicalPartitionOnSite(tt.args.ctx, tt.args.request)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestManageNVLinkLogicalPartition_UpdateNVLinkLogicalPartitionOnSite(t *testing.T) {
	mockCarbide := cClient.NewMockCarbideClient()

	carbideAtomicClient := cClient.NewCarbideAtomicClient(&cClient.CarbideClientConfig{})
	carbideAtomicClient.SwapClient(mockCarbide)

	type fields struct {
		carbideAtomicClient *cClient.CarbideAtomicClient
	}
	type args struct {
		ctx     context.Context
		request *cwssaws.NVLinkLogicalPartitionUpdateRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "test update nvlink logical partition success",
			fields: fields{
				carbideAtomicClient: carbideAtomicClient,
			},
			args: args{
				ctx: context.Background(),
				request: &cwssaws.NVLinkLogicalPartitionUpdateRequest{
					Id: &cwssaws.NVLinkLogicalPartitionId{Value: "b410867c-655a-11ef-bc4a-0393098e5d09"},
					Config: &cwssaws.NVLinkLogicalPartitionConfig{
						Metadata: &cwssaws.Metadata{
							Name: "test_name",
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "test update nvlink logical partition fail on missing id",
			fields: fields{
				carbideAtomicClient: carbideAtomicClient,
			},
			args: args{
				ctx: context.Background(),
				request: &cwssaws.NVLinkLogicalPartitionUpdateRequest{
					Id: nil,
					Config: &cwssaws.NVLinkLogicalPartitionConfig{
						Metadata: &cwssaws.Metadata{
							Name: "test_name",
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "test create nvlink logical partition fail on missing name",
			fields: fields{
				carbideAtomicClient: carbideAtomicClient,
			},
			args: args{
				ctx: context.Background(),
				request: &cwssaws.NVLinkLogicalPartitionUpdateRequest{
					Id: &cwssaws.NVLinkLogicalPartitionId{Value: "b410867c-655a-11ef-bc4a-0393098e5d09"},
					Config: &cwssaws.NVLinkLogicalPartitionConfig{
						Metadata: &cwssaws.Metadata{
							Name: "",
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "test update nvlink logical partition fail on missing request",
			fields: fields{
				carbideAtomicClient: carbideAtomicClient,
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
			mnvllp := NewManageNVLinkLogicalPartition(tt.fields.carbideAtomicClient)
			err := mnvllp.UpdateNVLinkLogicalPartitionOnSite(tt.args.ctx, tt.args.request)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestManageNVLinkLogicalPartition_DeleteNVLinkLogicalPartitionOnSite(t *testing.T) {
	mockCarbide := cClient.NewMockCarbideClient()

	carbideAtomicClient := cClient.NewCarbideAtomicClient(&cClient.CarbideClientConfig{})
	carbideAtomicClient.SwapClient(mockCarbide)

	type fields struct {
		CarbideAtomicClient *cClient.CarbideAtomicClient
	}
	type args struct {
		ctx     context.Context
		request *cwssaws.NVLinkLogicalPartitionDeletionRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "test delete nvlink logical partition success",
			fields: fields{
				CarbideAtomicClient: carbideAtomicClient,
			},
			args: args{
				ctx: context.Background(),
				request: &cwssaws.NVLinkLogicalPartitionDeletionRequest{
					Id: &cwssaws.NVLinkLogicalPartitionId{Value: "b410867c-655a-11ef-bc4a-0393098e5d09"},
				},
			},
			wantErr: false,
		},
		{
			name: "test delete nvlink logical partition fail on blank id",
			fields: fields{
				CarbideAtomicClient: carbideAtomicClient,
			},
			args: args{
				ctx: context.Background(),
				request: &cwssaws.NVLinkLogicalPartitionDeletionRequest{
					Id: &cwssaws.NVLinkLogicalPartitionId{Value: ""},
				},
			},
			wantErr: true,
		},
		{
			name: "test delete nvlink logical partition fail on missing request",
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
			mm := NewManageNVLinkLogicalPartition(tt.fields.CarbideAtomicClient)
			err := mm.DeleteNVLinkLogicalPartitionOnSite(tt.args.ctx, tt.args.request)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
