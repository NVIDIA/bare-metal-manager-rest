/*
 * SPDX-FileCopyrightText: Copyright (c) 2026 NVIDIA CORPORATION & AFFILIATES. All rights reserved.
 * SPDX-License-Identifier: Apache-2.0
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
// Package client provides a gRPC client for interacting with the RLA service.
// This package can be imported by external modules to communicate with RLA.
//
// The client uses types from pkg/types, which can be imported independently
// for interface definitions and mocking without gRPC dependencies.
package client

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/timestamppb"

	pb "github.com/nvidia/carbide-rest/rla/pkg/proto/v1"
	"github.com/nvidia/carbide-rest/rla/pkg/types"
)

// Client is the gRPC client for interacting with the RLA service.
type Client struct {
	client pb.RLAClient
	conn   *grpc.ClientConn
}

// New creates a new client with the given configuration.
func New(c Config) (*Client, error) {
	if err := c.Validate(); err != nil {
		return nil, err
	}

	conn, err := grpc.NewClient(
		c.Target(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, err
	}

	return &Client{
		conn:   conn,
		client: pb.NewRLAClient(conn),
	}, nil
}

// Close closes the gRPC connection
func (c *Client) Close() error {
	return c.conn.Close()
}

// CreateExpectedRack creates a new expected rack and returns its UUID.
func (c *Client) CreateExpectedRack(
	ctx context.Context,
	rack *types.Rack,
) (uuid.UUID, error) {
	rsp, err := c.client.CreateExpectedRack(
		ctx,
		&pb.CreateExpectedRackRequest{
			Rack: rackToProto(rack),
		},
	)

	if err != nil {
		return uuid.Nil, err
	}

	return uuidFromProto(rsp.GetId()), nil
}

// GetRackInfoByID retrieves rack information by its UUID.
func (c *Client) GetRackInfoByID(
	ctx context.Context,
	id uuid.UUID,
	withComponents bool,
) (*types.Rack, error) {
	rsp, err := c.client.GetRackInfoByID(
		ctx,
		&pb.GetRackInfoByIDRequest{
			Id:             uuidToProto(id),
			WithComponents: withComponents,
		},
	)

	if err != nil {
		return nil, err
	}

	return rackFromProto(rsp.Rack), nil
}

// GetRackInfoBySerial retrieves rack information by its manufacturer and
// serial number.
func (c *Client) GetRackInfoBySerial(
	ctx context.Context,
	manufacturer string,
	serial string,
	withComponents bool,
) (*types.Rack, error) {
	rsp, err := c.client.GetRackInfoBySerial(
		ctx,
		&pb.GetRackInfoBySerialRequest{
			SerialInfo: &pb.DeviceSerialInfo{
				Manufacturer: manufacturer,
				SerialNumber: serial,
			},
			WithComponents: withComponents,
		},
	)

	if err != nil {
		return nil, err
	}

	return rackFromProto(rsp.Rack), nil
}

// GetComponentInfoByID retrieves component information by its UUID.
func (c *Client) GetComponentInfoByID(
	ctx context.Context,
	id uuid.UUID,
	withRack bool,
) (*types.Component, *types.Rack, error) {
	rsp, err := c.client.GetComponentInfoByID(
		ctx,
		&pb.GetComponentInfoByIDRequest{
			Id:       uuidToProto(id),
			WithRack: withRack,
		},
	)

	if err != nil {
		return nil, nil, err
	}

	return componentFromProto(rsp.Component), rackFromProto(rsp.Rack), nil
}

// GetComponentInfoBySerial retrieves component information by its
// manufacturer and serial number.
func (c *Client) GetComponentInfoBySerial(
	ctx context.Context,
	manufacturer string,
	serial string,
	withRack bool,
) (*types.Component, *types.Rack, error) {
	rsp, err := c.client.GetComponentInfoBySerial(
		ctx,
		&pb.GetComponentInfoBySerialRequest{
			SerialInfo: &pb.DeviceSerialInfo{
				Manufacturer: manufacturer,
				SerialNumber: serial,
			},
			WithRack: withRack,
		},
	)

	if err != nil {
		return nil, nil, err
	}

	return componentFromProto(rsp.Component), rackFromProto(rsp.Rack), nil
}

// GetListOfRacks retrieves a list of racks matching the query.
func (c *Client) GetListOfRacks(
	ctx context.Context,
	info *types.StringQueryInfo,
	pagination *types.Pagination,
	withComponents bool,
) ([]*types.Rack, int32, error) {
	filters := make([]*pb.Filter, 0)
	if info != nil {
		filters = append(filters, &pb.Filter{
			Field: &pb.Filter_RackField{
				RackField: pb.RackFilterField_RACK_FILTER_FIELD_NAME,
			},
			QueryInfo: stringQueryInfoToProto(info),
		})
	}

	rsp, err := c.client.GetListOfRacks(
		ctx,
		&pb.GetListOfRacksRequest{
			Filters:        filters,
			Pagination:     paginationToProto(pagination),
			WithComponents: withComponents,
		},
	)

	if err != nil {
		return nil, 0, err
	}

	results := make([]*types.Rack, 0, len(rsp.Racks))
	for _, rack := range rsp.Racks {
		results = append(results, rackFromProto(rack))
	}

	return results, rsp.Total, nil
}

// CreateNVLDomain creates a new NVL domain and returns its UUID.
func (c *Client) CreateNVLDomain(
	ctx context.Context,
	nvlDomain *types.NVLDomain,
) (uuid.UUID, error) {
	rsp, err := c.client.CreateNVLDomain(
		ctx,
		&pb.CreateNVLDomainRequest{NvlDomain: nvlDomainToProto(nvlDomain)},
	)

	if err != nil {
		return uuid.Nil, err
	}

	return uuidFromProto(rsp.Id), nil
}

// AttachRacksToNVLDomain attaches racks to an NVL domain.
func (c *Client) AttachRacksToNVLDomain(
	ctx context.Context,
	nvlDomainID types.Identifier,
	rackIDs []types.Identifier,
) error {
	pbRackIDs := make([]*pb.Identifier, 0, len(rackIDs))
	for _, rackID := range rackIDs {
		pbRackIDs = append(pbRackIDs, identifierToProto(&rackID))
	}

	_, err := c.client.AttachRacksToNVLDomain(
		ctx,
		&pb.AttachRacksToNVLDomainRequest{
			NvlDomainIdentifier: identifierToProto(&nvlDomainID),
			RackIdentifiers:     pbRackIDs,
		},
	)

	return err
}

// DetachRacksFromNVLDomain detaches racks from their NVL domain.
func (c *Client) DetachRacksFromNVLDomain(
	ctx context.Context,
	rackIDs []types.Identifier,
) error {
	pbRackIDs := make([]*pb.Identifier, 0, len(rackIDs))
	for _, rackID := range rackIDs {
		pbRackIDs = append(pbRackIDs, identifierToProto(&rackID))
	}

	_, err := c.client.DetachRacksFromNVLDomain(
		ctx,
		&pb.DetachRacksFromNVLDomainRequest{RackIdentifiers: pbRackIDs},
	)

	return err
}

// GetListOfNVLDomains retrieves a list of NVL domains matching the query.
func (c *Client) GetListOfNVLDomains(
	ctx context.Context,
	info *types.StringQueryInfo,
	pagination *types.Pagination,
) ([]*types.NVLDomain, int32, error) {
	rsp, err := c.client.GetListOfNVLDomains(
		ctx,
		&pb.GetListOfNVLDomainsRequest{
			Info:       stringQueryInfoToProto(info),
			Pagination: paginationToProto(pagination),
		},
	)

	if err != nil {
		return nil, 0, err
	}

	results := make([]*types.NVLDomain, 0, len(rsp.NvlDomains))
	for _, nvlDomain := range rsp.NvlDomains {
		results = append(results, nvlDomainFromProto(nvlDomain))
	}

	return results, rsp.Total, nil
}

// GetRacksForNVLDomain retrieves racks belonging to an NVL domain.
func (c *Client) GetRacksForNVLDomain(
	ctx context.Context,
	nvlDomainID types.Identifier,
) ([]*types.Rack, error) {
	rsp, err := c.client.GetRacksForNVLDomain(
		ctx,
		&pb.GetRacksForNVLDomainRequest{
			NvlDomainIdentifier: identifierToProto(&nvlDomainID),
		},
	)

	if err != nil {
		return nil, err
	}

	results := make([]*types.Rack, 0, len(rsp.Racks))
	for _, rack := range rsp.Racks {
		results = append(results, rackFromProto(rack))
	}

	return results, nil
}

// UpgradeFirmwareByRackIDs upgrades firmware for components in the given rack IDs.
func (c *Client) UpgradeFirmwareByRackIDs(
	ctx context.Context,
	rackIDs []uuid.UUID,
	componentType types.ComponentType,
	startTime, endTime *time.Time,
) (*UpgradeFirmwareResult, error) {
	rackTargets := make([]*pb.RackTarget, 0, len(rackIDs))
	for _, id := range rackIDs {
		rt := &pb.RackTarget{
			Identifier: &pb.RackTarget_Id{Id: uuidToProto(id)},
		}
		if componentType != types.ComponentTypeUnknown {
			rt.ComponentTypes = []pb.ComponentType{componentTypeToProto(componentType)}
		}
		rackTargets = append(rackTargets, rt)
	}

	req := &pb.UpgradeFirmwareRequest{
		TargetSpec: &pb.OperationTargetSpec{
			Targets: &pb.OperationTargetSpec_Racks{
				Racks: &pb.RackTargets{Targets: rackTargets},
			},
		},
	}
	if startTime != nil {
		req.StartTime = timestamppb.New(*startTime)
	}
	if endTime != nil {
		req.EndTime = timestamppb.New(*endTime)
	}

	rsp, err := c.client.UpgradeFirmware(ctx, req)
	if err != nil {
		return nil, err
	}

	return &UpgradeFirmwareResult{
		TaskIDs: uuidsFromProto(rsp.GetTaskIds()),
	}, nil
}

// UpgradeFirmwareByRackNames upgrades firmware for components in the given rack names.
func (c *Client) UpgradeFirmwareByRackNames(
	ctx context.Context,
	rackNames []string,
	componentType types.ComponentType,
	startTime, endTime *time.Time,
) (*UpgradeFirmwareResult, error) {
	rackTargets := make([]*pb.RackTarget, 0, len(rackNames))
	for _, name := range rackNames {
		rt := &pb.RackTarget{
			Identifier: &pb.RackTarget_Name{Name: name},
		}
		if componentType != types.ComponentTypeUnknown {
			rt.ComponentTypes = []pb.ComponentType{componentTypeToProto(componentType)}
		}
		rackTargets = append(rackTargets, rt)
	}

	req := &pb.UpgradeFirmwareRequest{
		TargetSpec: &pb.OperationTargetSpec{
			Targets: &pb.OperationTargetSpec_Racks{
				Racks: &pb.RackTargets{Targets: rackTargets},
			},
		},
	}
	if startTime != nil {
		req.StartTime = timestamppb.New(*startTime)
	}
	if endTime != nil {
		req.EndTime = timestamppb.New(*endTime)
	}

	rsp, err := c.client.UpgradeFirmware(ctx, req)
	if err != nil {
		return nil, err
	}

	return &UpgradeFirmwareResult{
		TaskIDs: uuidsFromProto(rsp.GetTaskIds()),
	}, nil
}

// UpgradeFirmwareByMachineIDs upgrades firmware for the given machine IDs (external component IDs).
func (c *Client) UpgradeFirmwareByMachineIDs(
	ctx context.Context,
	machineIDs []string,
	startTime, endTime *time.Time,
) (*UpgradeFirmwareResult, error) {
	compTargets := make([]*pb.ComponentTarget, 0, len(machineIDs))
	for _, machineID := range machineIDs {
		compTargets = append(compTargets, &pb.ComponentTarget{
			Identifier: &pb.ComponentTarget_External{
				External: &pb.ExternalRef{
					Type: pb.ComponentType_COMPONENT_TYPE_COMPUTE,
					Id:   machineID,
				},
			},
		})
	}

	req := &pb.UpgradeFirmwareRequest{
		TargetSpec: &pb.OperationTargetSpec{
			Targets: &pb.OperationTargetSpec_Components{
				Components: &pb.ComponentTargets{Targets: compTargets},
			},
		},
	}
	if startTime != nil {
		req.StartTime = timestamppb.New(*startTime)
	}
	if endTime != nil {
		req.EndTime = timestamppb.New(*endTime)
	}

	rsp, err := c.client.UpgradeFirmware(ctx, req)
	if err != nil {
		return nil, err
	}

	return &UpgradeFirmwareResult{
		TaskIDs: uuidsFromProto(rsp.GetTaskIds()),
	}, nil
}

// PowerControlByRackIDs performs power control on components in the given rack IDs.
func (c *Client) PowerControlByRackIDs(
	ctx context.Context,
	rackIDs []uuid.UUID,
	componentType types.ComponentType,
	op types.PowerControlOp,
) (*PowerControlResult, error) {
	rackTargets := make([]*pb.RackTarget, 0, len(rackIDs))
	for _, id := range rackIDs {
		rt := &pb.RackTarget{
			Identifier: &pb.RackTarget_Id{Id: uuidToProto(id)},
		}
		if componentType != types.ComponentTypeUnknown {
			rt.ComponentTypes = []pb.ComponentType{componentTypeToProto(componentType)}
		}
		rackTargets = append(rackTargets, rt)
	}

	targetSpec := &pb.OperationTargetSpec{
		Targets: &pb.OperationTargetSpec_Racks{
			Racks: &pb.RackTargets{Targets: rackTargets},
		},
	}

	return c.executePowerControl(ctx, targetSpec, op)
}

// PowerControlByRackNames performs power control on components in the given rack names.
func (c *Client) PowerControlByRackNames(
	ctx context.Context,
	rackNames []string,
	componentType types.ComponentType,
	op types.PowerControlOp,
) (*PowerControlResult, error) {
	rackTargets := make([]*pb.RackTarget, 0, len(rackNames))
	for _, name := range rackNames {
		rt := &pb.RackTarget{
			Identifier: &pb.RackTarget_Name{Name: name},
		}
		if componentType != types.ComponentTypeUnknown {
			rt.ComponentTypes = []pb.ComponentType{componentTypeToProto(componentType)}
		}
		rackTargets = append(rackTargets, rt)
	}

	targetSpec := &pb.OperationTargetSpec{
		Targets: &pb.OperationTargetSpec_Racks{
			Racks: &pb.RackTargets{Targets: rackTargets},
		},
	}

	return c.executePowerControl(ctx, targetSpec, op)
}

// PowerControlByMachineIDs performs power control on the given machine IDs.
func (c *Client) PowerControlByMachineIDs(
	ctx context.Context,
	machineIDs []string,
	op types.PowerControlOp,
) (*PowerControlResult, error) {
	compTargets := make([]*pb.ComponentTarget, 0, len(machineIDs))
	for _, machineID := range machineIDs {
		compTargets = append(compTargets, &pb.ComponentTarget{
			Identifier: &pb.ComponentTarget_External{
				External: &pb.ExternalRef{
					Type: pb.ComponentType_COMPONENT_TYPE_COMPUTE,
					Id:   machineID,
				},
			},
		})
	}

	targetSpec := &pb.OperationTargetSpec{
		Targets: &pb.OperationTargetSpec_Components{
			Components: &pb.ComponentTargets{Targets: compTargets},
		},
	}

	return c.executePowerControl(ctx, targetSpec, op)
}

// executePowerControl executes a power control operation with the given target spec.
func (c *Client) executePowerControl(
	ctx context.Context,
	targetSpec *pb.OperationTargetSpec,
	op types.PowerControlOp,
) (*PowerControlResult, error) {
	var rsp *pb.SubmitTaskResponse
	var err error

	pbOp := powerControlOpToProto(op)

	switch pbOp {
	case pb.PowerControlOp_POWER_CONTROL_OP_ON, pb.PowerControlOp_POWER_CONTROL_OP_FORCE_ON:
		rsp, err = c.client.PowerOnRack(ctx, &pb.PowerOnRackRequest{
			TargetSpec: targetSpec,
		})

	case pb.PowerControlOp_POWER_CONTROL_OP_OFF:
		rsp, err = c.client.PowerOffRack(ctx, &pb.PowerOffRackRequest{
			TargetSpec: targetSpec,
			Forced:     false,
		})

	case pb.PowerControlOp_POWER_CONTROL_OP_FORCE_OFF:
		rsp, err = c.client.PowerOffRack(ctx, &pb.PowerOffRackRequest{
			TargetSpec: targetSpec,
			Forced:     true,
		})

	case pb.PowerControlOp_POWER_CONTROL_OP_RESTART, pb.PowerControlOp_POWER_CONTROL_OP_WARM_RESET:
		rsp, err = c.client.PowerResetRack(ctx, &pb.PowerResetRackRequest{
			TargetSpec: targetSpec,
			Forced:     false,
		})

	case pb.PowerControlOp_POWER_CONTROL_OP_FORCE_RESTART, pb.PowerControlOp_POWER_CONTROL_OP_COLD_RESET:
		rsp, err = c.client.PowerResetRack(ctx, &pb.PowerResetRackRequest{
			TargetSpec: targetSpec,
			Forced:     true,
		})

	default:
		return nil, fmt.Errorf("unsupported power control operation: %v", op)
	}

	if err != nil {
		return nil, err
	}

	return &PowerControlResult{
		TaskIDs: uuidsFromProto(rsp.GetTaskIds()),
	}, nil
}

// GetExpectedComponentsByRackIDs retrieves expected components from local database by rack IDs.
func (c *Client) GetExpectedComponentsByRackIDs(
	ctx context.Context,
	rackIDs []uuid.UUID,
	componentType types.ComponentType,
) (*GetExpectedComponentsResult, error) {
	rackTargets := make([]*pb.RackTarget, 0, len(rackIDs))
	for _, id := range rackIDs {
		rt := &pb.RackTarget{
			Identifier: &pb.RackTarget_Id{Id: uuidToProto(id)},
		}
		if componentType != types.ComponentTypeUnknown {
			rt.ComponentTypes = []pb.ComponentType{componentTypeToProto(componentType)}
		}
		rackTargets = append(rackTargets, rt)
	}

	// Build filters array
	filters := make([]*pb.Filter, 0)
	if componentType != types.ComponentTypeUnknown {
		// Convert ComponentType enum to string
		typeStr := componentTypeToString(componentType)
		filters = append(filters, &pb.Filter{
			Field: &pb.Filter_ComponentField{
				ComponentField: pb.ComponentFilterField_COMPONENT_FILTER_FIELD_TYPE,
			},
			QueryInfo: &pb.StringQueryInfo{
				Patterns:   []string{typeStr},
				IsWildcard: false,
				UseOr:      false,
			},
		})
	}

	rsp, err := c.client.GetComponents(
		ctx,
		&pb.GetComponentsRequest{
			TargetSpec: &pb.OperationTargetSpec{
				Targets: &pb.OperationTargetSpec_Racks{
					Racks: &pb.RackTargets{Targets: rackTargets},
				},
			},
			Filters: filters,
		},
	)
	if err != nil {
		return nil, err
	}

	return convertGetComponentsResponse(rsp), nil
}

// GetExpectedComponentsByRackNames retrieves expected components from local database by rack names.
func (c *Client) GetExpectedComponentsByRackNames(
	ctx context.Context,
	rackNames []string,
	componentType types.ComponentType,
) (*GetExpectedComponentsResult, error) {
	rackTargets := make([]*pb.RackTarget, 0, len(rackNames))
	for _, name := range rackNames {
		rt := &pb.RackTarget{
			Identifier: &pb.RackTarget_Name{Name: name},
		}
		if componentType != types.ComponentTypeUnknown {
			rt.ComponentTypes = []pb.ComponentType{componentTypeToProto(componentType)}
		}
		rackTargets = append(rackTargets, rt)
	}

	// Build filters array
	filters := make([]*pb.Filter, 0)
	if componentType != types.ComponentTypeUnknown {
		// Convert ComponentType enum to string
		typeStr := componentTypeToString(componentType)
		filters = append(filters, &pb.Filter{
			Field: &pb.Filter_ComponentField{
				ComponentField: pb.ComponentFilterField_COMPONENT_FILTER_FIELD_TYPE,
			},
			QueryInfo: &pb.StringQueryInfo{
				Patterns:   []string{typeStr},
				IsWildcard: false,
				UseOr:      false,
			},
		})
	}

	rsp, err := c.client.GetComponents(
		ctx,
		&pb.GetComponentsRequest{
			TargetSpec: &pb.OperationTargetSpec{
				Targets: &pb.OperationTargetSpec_Racks{
					Racks: &pb.RackTargets{Targets: rackTargets},
				},
			},
			Filters: filters,
		},
	)
	if err != nil {
		return nil, err
	}

	return convertGetComponentsResponse(rsp), nil
}

// GetExpectedComponentsByComponentIDs retrieves expected components by external component IDs.
func (c *Client) GetExpectedComponentsByComponentIDs(
	ctx context.Context,
	componentIDs []string,
	componentType types.ComponentType,
) (*GetExpectedComponentsResult, error) {
	compTargets := make([]*pb.ComponentTarget, 0, len(componentIDs))
	for _, compID := range componentIDs {
		compTargets = append(compTargets, &pb.ComponentTarget{
			Identifier: &pb.ComponentTarget_External{
				External: &pb.ExternalRef{
					Type: componentTypeToProto(componentType),
					Id:   compID,
				},
			},
		})
	}

	// Build filters array
	filters := make([]*pb.Filter, 0)
	if componentType != types.ComponentTypeUnknown {
		// Convert ComponentType enum to string
		typeStr := componentTypeToString(componentType)
		filters = append(filters, &pb.Filter{
			Field: &pb.Filter_ComponentField{
				ComponentField: pb.ComponentFilterField_COMPONENT_FILTER_FIELD_TYPE,
			},
			QueryInfo: &pb.StringQueryInfo{
				Patterns:   []string{typeStr},
				IsWildcard: false,
				UseOr:      false,
			},
		})
	}

	rsp, err := c.client.GetComponents(
		ctx,
		&pb.GetComponentsRequest{
			TargetSpec: &pb.OperationTargetSpec{
				Targets: &pb.OperationTargetSpec_Components{
					Components: &pb.ComponentTargets{Targets: compTargets},
				},
			},
			Filters: filters,
		},
	)
	if err != nil {
		return nil, err
	}

	return convertGetComponentsResponse(rsp), nil
}

func convertGetComponentsResponse(rsp *pb.GetComponentsResponse) *GetExpectedComponentsResult {
	components := make([]*types.Component, 0, len(rsp.Components))
	for _, c := range rsp.Components {
		components = append(components, componentFromProto(c))
	}

	return &GetExpectedComponentsResult{
		Components: components,
		Total:      int(rsp.Total),
	}
}

// GetActualComponentsByRackIDs retrieves actual components from external systems by rack IDs.
// NOTE: This method is temporarily unavailable because GetActualComponents rpc is commented out in proto.
// The implementation remains for future use when the rpc is re-enabled.
// TODO: Uncomment when GetActualComponents rpc is re-enabled in proto.
func (c *Client) GetActualComponentsByRackIDs(
	ctx context.Context,
	rackIDs []uuid.UUID,
	componentType types.ComponentType,
) (*GetActualComponentsResult, error) {
	// TODO: Uncomment when GetActualComponents rpc is re-enabled in proto
	return nil, fmt.Errorf("GetActualComponents rpc is temporarily disabled in proto")
	/*
		rackTargets := make([]*pb.RackTarget, 0, len(rackIDs))
		for _, id := range rackIDs {
			rt := &pb.RackTarget{
				Identifier: &pb.RackTarget_Id{Id: uuidToProto(id)},
			}
			if componentType != types.ComponentTypeUnknown {
				rt.ComponentTypes = []pb.ComponentType{componentTypeToProto(componentType)}
			}
			rackTargets = append(rackTargets, rt)
		}

		rsp, err := c.client.GetActualComponents(
			ctx,
			&pb.GetActualComponentsRequest{
				TargetSpec: &pb.OperationTargetSpec{
					Targets: &pb.OperationTargetSpec_Racks{
						Racks: &pb.RackTargets{Targets: rackTargets},
					},
				},
			},
		)
		if err != nil {
			return nil, err
		}

		return convertGetActualComponentsResponse(rsp), nil
	*/
}

// GetActualComponentsByRackNames retrieves actual components from external systems by rack names.
// NOTE: This method is temporarily unavailable because GetActualComponents rpc is commented out in proto.
// TODO: Uncomment when GetActualComponents rpc is re-enabled in proto.
func (c *Client) GetActualComponentsByRackNames(
	ctx context.Context,
	rackNames []string,
	componentType types.ComponentType,
) (*GetActualComponentsResult, error) {
	// TODO: Uncomment when GetActualComponents rpc is re-enabled in proto
	return nil, fmt.Errorf("GetActualComponents rpc is temporarily disabled in proto")
	/*
		rackTargets := make([]*pb.RackTarget, 0, len(rackNames))
		for _, name := range rackNames {
			rt := &pb.RackTarget{
				Identifier: &pb.RackTarget_Name{Name: name},
			}
			if componentType != types.ComponentTypeUnknown {
				rt.ComponentTypes = []pb.ComponentType{componentTypeToProto(componentType)}
			}
			rackTargets = append(rackTargets, rt)
		}

		rsp, err := c.client.GetActualComponents(
			ctx,
			&pb.GetActualComponentsRequest{
				TargetSpec: &pb.OperationTargetSpec{
					Targets: &pb.OperationTargetSpec_Racks{
						Racks: &pb.RackTargets{Targets: rackTargets},
					},
				},
			},
		)
		if err != nil {
			return nil, err
		}

		return convertGetActualComponentsResponse(rsp), nil
	*/
}

// GetActualComponentsByComponentIDs retrieves actual components by external component IDs.
// NOTE: This method is temporarily unavailable because GetActualComponents rpc is commented out in proto.
// TODO: Uncomment when GetActualComponents rpc is re-enabled in proto.
func (c *Client) GetActualComponentsByComponentIDs(
	ctx context.Context,
	componentIDs []string,
	componentType types.ComponentType,
) (*GetActualComponentsResult, error) {
	// TODO: Uncomment when GetActualComponents rpc is re-enabled in proto
	return nil, fmt.Errorf("GetActualComponents rpc is temporarily disabled in proto")
	/*
		compTargets := make([]*pb.ComponentTarget, 0, len(componentIDs))
		for _, compID := range componentIDs {
			compTargets = append(compTargets, &pb.ComponentTarget{
				Identifier: &pb.ComponentTarget_External{
					External: &pb.ExternalRef{
						Type: componentTypeToProto(componentType),
						Id:   compID,
					},
				},
			})
		}

		rsp, err := c.client.GetActualComponents(
			ctx,
			&pb.GetActualComponentsRequest{
				TargetSpec: &pb.OperationTargetSpec{
					Targets: &pb.OperationTargetSpec_Components{
						Components: &pb.ComponentTargets{Targets: compTargets},
					},
				},
			},
		)
		if err != nil {
			return nil, err
		}

		return convertGetActualComponentsResponse(rsp), nil
	*/
}

func convertGetActualComponentsResponse(rsp *pb.GetActualComponentsResponse) *GetActualComponentsResult {
	components := make([]*types.ActualComponent, 0, len(rsp.Components))
	for _, c := range rsp.Components {
		components = append(components, actualComponentFromProto(c))
	}

	return &GetActualComponentsResult{
		Components: components,
		Total:      int(rsp.Total),
	}
}

// ValidateComponentsByRackIDs validates expected vs actual components by rack IDs.
func (c *Client) ValidateComponentsByRackIDs(
	ctx context.Context,
	rackIDs []uuid.UUID,
	componentType types.ComponentType,
) (*ValidateComponentsResult, error) {
	rackTargets := make([]*pb.RackTarget, 0, len(rackIDs))
	for _, id := range rackIDs {
		rt := &pb.RackTarget{
			Identifier: &pb.RackTarget_Id{Id: uuidToProto(id)},
		}
		if componentType != types.ComponentTypeUnknown {
			rt.ComponentTypes = []pb.ComponentType{componentTypeToProto(componentType)}
		}
		rackTargets = append(rackTargets, rt)
	}

	rsp, err := c.client.ValidateComponents(
		ctx,
		&pb.ValidateComponentsRequest{
			TargetSpec: &pb.OperationTargetSpec{
				Targets: &pb.OperationTargetSpec_Racks{
					Racks: &pb.RackTargets{Targets: rackTargets},
				},
			},
		},
	)
	if err != nil {
		return nil, err
	}

	return convertValidateComponentsResponse(rsp), nil
}

// ValidateComponentsByRackNames validates expected vs actual components by rack names.
func (c *Client) ValidateComponentsByRackNames(
	ctx context.Context,
	rackNames []string,
	componentType types.ComponentType,
) (*ValidateComponentsResult, error) {
	rackTargets := make([]*pb.RackTarget, 0, len(rackNames))
	for _, name := range rackNames {
		rt := &pb.RackTarget{
			Identifier: &pb.RackTarget_Name{Name: name},
		}
		if componentType != types.ComponentTypeUnknown {
			rt.ComponentTypes = []pb.ComponentType{componentTypeToProto(componentType)}
		}
		rackTargets = append(rackTargets, rt)
	}

	rsp, err := c.client.ValidateComponents(
		ctx,
		&pb.ValidateComponentsRequest{
			TargetSpec: &pb.OperationTargetSpec{
				Targets: &pb.OperationTargetSpec_Racks{
					Racks: &pb.RackTargets{Targets: rackTargets},
				},
			},
		},
	)
	if err != nil {
		return nil, err
	}

	return convertValidateComponentsResponse(rsp), nil
}

// ValidateComponentsByComponentIDs validates expected vs actual components by external component IDs.
func (c *Client) ValidateComponentsByComponentIDs(
	ctx context.Context,
	componentIDs []string,
	componentType types.ComponentType,
) (*ValidateComponentsResult, error) {
	compTargets := make([]*pb.ComponentTarget, 0, len(componentIDs))
	for _, compID := range componentIDs {
		compTargets = append(compTargets, &pb.ComponentTarget{
			Identifier: &pb.ComponentTarget_External{
				External: &pb.ExternalRef{
					Type: componentTypeToProto(componentType),
					Id:   compID,
				},
			},
		})
	}

	rsp, err := c.client.ValidateComponents(
		ctx,
		&pb.ValidateComponentsRequest{
			TargetSpec: &pb.OperationTargetSpec{
				Targets: &pb.OperationTargetSpec_Components{
					Components: &pb.ComponentTargets{Targets: compTargets},
				},
			},
		},
	)
	if err != nil {
		return nil, err
	}

	return convertValidateComponentsResponse(rsp), nil
}

func convertValidateComponentsResponse(rsp *pb.ValidateComponentsResponse) *ValidateComponentsResult {
	diffs := make([]*types.ComponentDiff, 0, len(rsp.Diffs))
	for _, d := range rsp.Diffs {
		diffs = append(diffs, componentDiffFromProto(d))
	}

	return &ValidateComponentsResult{
		Diffs:               diffs,
		TotalDiffs:          int(rsp.TotalDiffs),
		OnlyInExpectedCount: int(rsp.OnlyInExpectedCount),
		OnlyInActualCount:   int(rsp.OnlyInActualCount),
		DriftCount:          int(rsp.DriftCount),
		MatchCount:          int(rsp.MatchCount),
	}
}

// ListTasks lists tasks matching the query.
func (c *Client) ListTasks(
	ctx context.Context,
	rackID *uuid.UUID,
	activeOnly bool,
	pagination *types.Pagination,
) (*ListTasksResult, error) {
	req := &pb.ListTasksRequest{
		ActiveOnly: activeOnly,
		Pagination: paginationToProto(pagination),
	}

	if rackID != nil {
		req.RackId = uuidToProto(*rackID)
	}

	rsp, err := c.client.ListTasks(ctx, req)
	if err != nil {
		return nil, err
	}

	tasks := make([]*types.Task, 0, len(rsp.Tasks))
	for _, t := range rsp.Tasks {
		tasks = append(tasks, taskFromProto(t))
	}

	return &ListTasksResult{
		Tasks: tasks,
		Total: int(rsp.Total),
	}, nil
}

// GetTasksByIDs retrieves tasks by their IDs.
func (c *Client) GetTasksByIDs(
	ctx context.Context,
	taskIDs []uuid.UUID,
) ([]*types.Task, error) {
	pbIDs := make([]*pb.UUID, 0, len(taskIDs))
	for _, id := range taskIDs {
		pbIDs = append(pbIDs, uuidToProto(id))
	}

	rsp, err := c.client.GetTasksByIDs(ctx, &pb.GetTasksByIDsRequest{
		TaskIds: pbIDs,
	})
	if err != nil {
		return nil, err
	}

	tasks := make([]*types.Task, 0, len(rsp.Tasks))
	for _, t := range rsp.Tasks {
		tasks = append(tasks, taskFromProto(t))
	}

	return tasks, nil
}
