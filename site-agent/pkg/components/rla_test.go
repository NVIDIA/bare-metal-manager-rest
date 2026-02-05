// SPDX-FileCopyrightText: Copyright (c) 2021-2023 NVIDIA CORPORATION & AFFILIATES. All rights reserved.
// SPDX-License-Identifier: LicenseRef-NvidiaProprietary
//
// NVIDIA CORPORATION, its affiliates and licensors retain all intellectual
// property and proprietary rights in and to this material, related
// documentation and any modifications thereto. Any use, reproduction,
// disclosure or distribution of this material and related documentation
// without an express license agreement from NVIDIA CORPORATION or
// its affiliates is strictly prohibited.

package elektra

import (
	"context"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/nvidia/carbide-rest/site-agent/pkg/components/managers/rla"
	rlav1 "github.com/nvidia/carbide-rest/workflow-schema/rla/protobuf/v1"
	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel"
)

// TestRlaRack - test the RLA grpc client
func TestRlaRack(t *testing.T) {
	TestInitElektra(t)
	grpcClient := testElektra.manager.API.RLA.GetGrpcClient()

	var rack *rlav1.Rack

	tcs := []struct {
		descr     string
		expectErr bool
		index     int
	}{{
		descr:     "get",
		expectErr: false,
		index:     0,
	}, {
		descr:     "list",
		expectErr: false,
		index:     0,
	},
	}
	rpcSucc := 0
	for _, tc := range tcs {
		t.Run(tc.descr, func(t *testing.T) {
			switch tc.descr {
			case "get":
				rackID := uuid.NewString()
				ctx := context.Background()

				// First create the rack in mock server (setup, not counted in metrics)
				createReq := &rlav1.CreateExpectedRackRequest{
					Rack: &rlav1.Rack{
						Info: &rlav1.DeviceInfo{
							Id:   &rlav1.UUID{Id: rackID},
							Name: "test-rack",
						},
					},
				}
				_, createErr := grpcClient.Rla().CreateExpectedRack(ctx, createReq)
				assert.Nil(t, createErr)

				// Now test GetRackInfoByID
				ctx, span := otel.Tracer(os.Getenv("LS_SERVICE_NAME")).Start(ctx, "RlaTest-GetRack")

				getRequest := &rlav1.GetRackInfoByIDRequest{
					Id: &rlav1.UUID{Id: rackID},
				}

				response, err := grpcClient.Rla().GetRackInfoByID(ctx, getRequest)
				span.End()
				rla.ManagerAccess.API.RLA.UpdateGrpcClientState(err)
				if err != nil {
					t.Log(err.Error())
				}
				assert.Nil(t, err)
				assert.NotNil(t, response)
				assert.NotNil(t, response.Rack)
				assert.NotNil(t, response.Rack.Info)
				assert.Equal(t, rackID, response.Rack.Info.Id.Id)
				rpcSucc++
				assert.Equal(t, 0,
					int(rla.ManagerAccess.Data.EB.Managers.RLA.State.GrpcFail.Load()))
				assert.Equal(t, rpcSucc,
					int(rla.ManagerAccess.Data.EB.Managers.RLA.State.GrpcSucc.Load()))
				rack = response.Rack
				t.Log("GRPCResponse", response)
			case "list":
				ctx := context.Background()
				ctx, span := otel.Tracer(os.Getenv("LS_SERVICE_NAME")).Start(ctx, "RlaTest-GetListOfRacks")
				listRequest := &rlav1.GetListOfRacksRequest{}
				resq, err := grpcClient.Rla().GetListOfRacks(ctx, listRequest)
				span.End()
				rla.ManagerAccess.API.RLA.UpdateGrpcClientState(err)
				if err != nil {
					t.Log(err.Error())
				}
				assert.Nil(t, err)
				assert.NotNil(t, resq)
				assert.NotNil(t, resq.Racks)
				// Verify that the rack we got earlier is in the list
				if rack != nil && len(resq.Racks) > 0 {
					found := false
					for _, r := range resq.Racks {
						if r.Info != nil && r.Info.Id != nil && r.Info.Id.Id == rack.Info.Id.Id {
							found = true
							break
						}
					}
					assert.True(t, found, "Previously retrieved rack should be in the list")
				}
				rpcSucc++
				assert.Equal(t, 0,
					int(rla.ManagerAccess.Data.EB.Managers.RLA.State.GrpcFail.Load()))
				assert.Equal(t, rpcSucc,
					int(rla.ManagerAccess.Data.EB.Managers.RLA.State.GrpcSucc.Load()))
				t.Log("GRPCResponse", resq)
			default:
				panic("invalid operation name")
			}
		})
	}
}
