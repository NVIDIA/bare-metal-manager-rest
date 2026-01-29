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
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	rlav1 "github.com/nvidia/carbide-rest/workflow-schema/rla/protobuf/v1"
	cClient "github.com/nvidia/carbide-rest/site-workflow/pkg/grpc/client"
)

func TestManageRack_GetRackByID(t *testing.T) {
	tests := []struct {
		name        string
		request     *rlav1.GetRackInfoByIDRequest
		mockResp    *rlav1.GetRackInfoResponse
		mockErr     error
		wantErr     bool
		errContains string
	}{
		{
			name:        "nil request returns error",
			request:     nil,
			mockResp:    nil,
			mockErr:     nil,
			wantErr:     true,
			errContains: "empty get rack request",
		},
		{
			name: "request with nil ID returns error",
			request: &rlav1.GetRackInfoByIDRequest{
				Id: nil,
			},
			mockResp:    nil,
			mockErr:     nil,
			wantErr:     true,
			errContains: "missing rack ID",
		},
		{
			name: "request with empty ID returns error",
			request: &rlav1.GetRackInfoByIDRequest{
				Id: &rlav1.UUID{Id: ""},
			},
			mockResp:    nil,
			mockErr:     nil,
			wantErr:     true,
			errContains: "missing rack ID",
		},
		{
			name: "successful request",
			request: &rlav1.GetRackInfoByIDRequest{
				Id:             &rlav1.UUID{Id: "test-rack-id"},
				WithComponents: true,
			},
			mockResp: &rlav1.GetRackInfoResponse{
				Rack: &rlav1.Rack{
					Info: &rlav1.DeviceInfo{
						Id:   &rlav1.UUID{Id: "test-rack-id"},
						Name: "Test Rack",
					},
				},
			},
			mockErr: nil,
			wantErr: false,
		},
		{
			name: "RLA client error",
			request: &rlav1.GetRackInfoByIDRequest{
				Id: &rlav1.UUID{Id: "test-rack-id"},
			},
			mockResp:    nil,
			mockErr:     errors.New("connection refused"),
			wantErr:     true,
			errContains: "connection refused",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mock RLA client
			mockRlaClient := cClient.NewMockRlaClient()
			mockRla := cClient.GetMockRla(mockRlaClient)

			// Setup mock expectation if request is valid
			if tt.request != nil && tt.request.GetId() != nil && tt.request.GetId().GetId() != "" {
				mockRla.On("GetRackInfoByID", mock.Anything, tt.request).Return(tt.mockResp, tt.mockErr)
			}

			// Create atomic client and swap with mock
			rlaAtomicClient := cClient.NewRlaAtomicClient(&cClient.RlaClientConfig{})
			rlaAtomicClient.SwapClient(mockRlaClient)

			// Create ManageRack instance
			manageRack := NewManageRack(rlaAtomicClient)

			// Execute activity
			ctx := context.Background()
			result, err := manageRack.GetRackByID(ctx, tt.request)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errContains != "" {
					assert.Contains(t, err.Error(), tt.errContains)
				}
				return
			}

			assert.NoError(t, err)
			assert.NotNil(t, result)
			assert.Equal(t, tt.mockResp.GetRack().GetInfo().GetId().GetId(), result.GetRack().GetInfo().GetId().GetId())
		})
	}
}

func TestManageRack_GetListOfRacks(t *testing.T) {
	tests := []struct {
		name        string
		request     *rlav1.GetListOfRacksRequest
		mockResp    *rlav1.GetListOfRacksResponse
		mockErr     error
		wantErr     bool
		errContains string
	}{
		{
			name:    "successful request - empty list",
			request: &rlav1.GetListOfRacksRequest{},
			mockResp: &rlav1.GetListOfRacksResponse{
				Racks: []*rlav1.Rack{},
				Total: 0,
			},
			mockErr: nil,
			wantErr: false,
		},
		{
			name: "successful request - multiple racks",
			request: &rlav1.GetListOfRacksRequest{
				WithComponents: true,
			},
			mockResp: &rlav1.GetListOfRacksResponse{
				Racks: []*rlav1.Rack{
					{
						Info: &rlav1.DeviceInfo{
							Id:   &rlav1.UUID{Id: "rack-1"},
							Name: "Rack 1",
						},
					},
					{
						Info: &rlav1.DeviceInfo{
							Id:   &rlav1.UUID{Id: "rack-2"},
							Name: "Rack 2",
						},
					},
				},
				Total: 2,
			},
			mockErr: nil,
			wantErr: false,
		},
		{
			name:        "RLA client error",
			request:     &rlav1.GetListOfRacksRequest{},
			mockResp:    nil,
			mockErr:     errors.New("internal server error"),
			wantErr:     true,
			errContains: "internal server error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mock RLA client
			mockRlaClient := cClient.NewMockRlaClient()
			mockRla := cClient.GetMockRla(mockRlaClient)

			// Setup mock expectation if request is valid
			if tt.request != nil {
				mockRla.On("GetListOfRacks", mock.Anything, tt.request).Return(tt.mockResp, tt.mockErr)
			}

			// Create atomic client and swap with mock
			rlaAtomicClient := cClient.NewRlaAtomicClient(&cClient.RlaClientConfig{})
			rlaAtomicClient.SwapClient(mockRlaClient)

			// Create ManageRack instance
			manageRack := NewManageRack(rlaAtomicClient)

			// Execute activity
			ctx := context.Background()
			result, err := manageRack.GetListOfRacks(ctx, tt.request)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errContains != "" {
					assert.Contains(t, err.Error(), tt.errContains)
				}
				return
			}

			assert.NoError(t, err)
			assert.NotNil(t, result)
			assert.Equal(t, tt.mockResp.GetTotal(), result.GetTotal())
			assert.Equal(t, len(tt.mockResp.GetRacks()), len(result.GetRacks()))
		})
	}
}
