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
package client

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/nvidia/carbide-rest/rla/internal/inventory/objects/component"
	"github.com/nvidia/carbide-rest/rla/internal/inventory/objects/nvldomain"
	"github.com/nvidia/carbide-rest/rla/internal/inventory/objects/rack"
	pb "github.com/nvidia/carbide-rest/rla/internal/proto/v1"
	identifier "github.com/nvidia/carbide-rest/rla/pkg/common/Identifier"
	"github.com/nvidia/carbide-rest/rla/pkg/converter/protobuf"
	dbquery "github.com/nvidia/carbide-rest/rla/pkg/db/query"
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
	rack *rack.Rack,
) (uuid.UUID, error) {
	rsp, err := c.client.CreateExpectedRack(
		ctx,
		&pb.CreateExpectedRackRequest{
			Rack: protobuf.RackTo(rack),
		},
	)

	if err != nil {
		return uuid.Nil, err
	}

	return protobuf.UUIDFrom(rsp.GetId()), nil
}

// GetRackInfoByID retrieves rack information by its UUID.
func (c *Client) GetRackInfoByID(
	ctx context.Context,
	id uuid.UUID,
	withComponents bool,
) (*rack.Rack, error) {
	rsp, err := c.client.GetRackInfoByID(
		ctx,
		&pb.GetRackInfoByIDRequest{
			Id:             protobuf.UUIDTo(id),
			WithComponents: withComponents,
		},
	)

	if err != nil {
		return nil, err
	}

	return protobuf.RackFrom(rsp.Rack), nil
}

// GetRackInfoBySerial retrieves rack information by its manufacturer and
// serial number.
func (c *Client) GetRackInfoBySerial(
	ctx context.Context,
	manufacturer string,
	serial string,
	withComponents bool,
) (*rack.Rack, error) {
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

	return protobuf.RackFrom(rsp.Rack), nil
}

// GetComponentInfoByID retrieves component information by its UUID.
func (c *Client) GetComponentInfoByID(
	ctx context.Context,
	id uuid.UUID,
	withRack bool,
) (*component.Component, *rack.Rack, error) {
	rsp, err := c.client.GetComponentInfoByID(
		ctx,
		&pb.GetComponentInfoByIDRequest{
			Id:       protobuf.UUIDTo(id),
			WithRack: withRack,
		},
	)

	if err != nil {
		return nil, nil, err
	}

	return protobuf.ComponentFrom(rsp.Component), protobuf.RackFrom(rsp.Rack), nil //nolint
}

// GetComponentInfoBySerial retrieves component information by its
// manufacturer and serail number.
func (c *Client) GetComponentInfoBySerial(
	ctx context.Context,
	manufacturer string,
	serial string,
	withRack bool,
) (*component.Component, *rack.Rack, error) {
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

	return protobuf.ComponentFrom(rsp.Component), protobuf.RackFrom(rsp.Rack), nil //nolint
}

func (c *Client) GetListOfRacks(
	ctx context.Context,
	info *dbquery.StringQueryInfo,
	pagination *dbquery.Pagination,
	withComponents bool,
) ([]*rack.Rack, int32, error) {
	rsp, err := c.client.GetListOfRacks(
		ctx,
		&pb.GetListOfRacksRequest{
			Info:           protobuf.StringQueryInfoTo(info),
			Pagination:     protobuf.PaginationTo(pagination),
			WithComponents: withComponents,
		},
	)

	if err != nil {
		return nil, 0, err
	}

	results := make([]*rack.Rack, 0, len(rsp.Racks))
	for _, rack := range rsp.Racks {
		results = append(results, protobuf.RackFrom(rack))
	}

	return results, rsp.Total, nil
}

func (c *Client) CreateNVLDomain(
	ctx context.Context,
	nvlDomain *nvldomain.NVLDomain,
) (uuid.UUID, error) {
	rsp, err := c.client.CreateNVLDomain(
		ctx,
		&pb.CreateNVLDomainRequest{NvlDomain: protobuf.NVLDomainTo(nvlDomain)},
	)

	if err != nil {
		return uuid.Nil, err
	}

	return protobuf.UUIDFrom(rsp.Id), nil
}

func (c *Client) AttachRacksToNVLDomain(
	ctx context.Context,
	nvlDomainID identifier.Identifier,
	rackIDs []identifier.Identifier,
) error {
	pbRackIDs := make([]*pb.Identifier, 0, len(rackIDs))
	for _, rackID := range rackIDs {
		pbRackIDs = append(pbRackIDs, protobuf.IdentifierTo(&rackID))
	}

	_, err := c.client.AttachRacksToNVLDomain(
		ctx,
		&pb.AttachRacksToNVLDomainRequest{
			NvlDomainIdentifier: protobuf.IdentifierTo(&nvlDomainID),
			RackIdentifiers:     pbRackIDs,
		},
	)

	return err
}

func (c *Client) DetachRacksFromNVLDomain(
	ctx context.Context,
	rackIDs []identifier.Identifier,
) error {
	pbRackIDs := make([]*pb.Identifier, 0, len(rackIDs))
	for _, rackID := range rackIDs {
		pbRackIDs = append(pbRackIDs, protobuf.IdentifierTo(&rackID))
	}

	_, err := c.client.DetachRacksFromNVLDomain(
		ctx,
		&pb.DetachRacksFromNVLDomainRequest{RackIdentifiers: pbRackIDs},
	)

	return err
}

func (c *Client) GetListOfNVLDomains(
	ctx context.Context,
	info *dbquery.StringQueryInfo,
	pagination *dbquery.Pagination,
) ([]*nvldomain.NVLDomain, int32, error) {
	rsp, err := c.client.GetListOfNVLDomains(
		ctx,
		&pb.GetListOfNVLDomainsRequest{
			Info:       protobuf.StringQueryInfoTo(info),
			Pagination: protobuf.PaginationTo(pagination),
		},
	)

	if err != nil {
		return nil, 0, err
	}

	results := make([]*nvldomain.NVLDomain, 0, len(rsp.NvlDomains))
	for _, nvlDomain := range rsp.NvlDomains {
		results = append(results, protobuf.NVLDomainFrom(nvlDomain))
	}

	return results, rsp.Total, nil
}

func (c *Client) GetRacksForNVLDomain(
	ctx context.Context,
	nvlDomainID identifier.Identifier,
) ([]*rack.Rack, error) {
	rsp, err := c.client.GetRacksForNVLDomain(
		ctx,
		&pb.GetRacksForNVLDomainRequest{
			NvlDomainIdentifier: protobuf.IdentifierTo(&nvlDomainID),
		},
	)

	if err != nil {
		return nil, err
	}

	results := make([]*rack.Rack, 0, len(rsp.Racks))
	for _, rack := range rsp.Racks {
		results = append(results, protobuf.RackFrom(rack))
	}

	return results, nil
}

// UpgradeFirmwareResult represents the result of a firmware upgrade operation.
type UpgradeFirmwareResult struct {
	TaskIDs []uuid.UUID // Multiple task IDs (1 task per rack)
}

// UpgradeFirmwareByRackIDs upgrades firmware for components in the given rack IDs.
// Note: target_version is not yet supported in client API.
func (c *Client) UpgradeFirmwareByRackIDs(
	ctx context.Context,
	rackIDs []uuid.UUID,
	componentType pb.ComponentType,
	startTime, endTime *time.Time,
) (*UpgradeFirmwareResult, error) {
	rackTargets := make([]*pb.RackTarget, 0, len(rackIDs))
	for _, id := range rackIDs {
		rt := &pb.RackTarget{
			Identifier: &pb.RackTarget_Id{Id: id.String()},
		}
		if componentType != pb.ComponentType_COMPONENT_TYPE_UNKNOWN {
			rt.ComponentTypes = []pb.ComponentType{componentType}
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
		TaskIDs: protobuf.UUIDsFrom(rsp.GetTaskIds()),
	}, nil
}

// UpgradeFirmwareByRackNames upgrades firmware for components in the given rack names.
// Note: target_version is not yet supported in client API.
func (c *Client) UpgradeFirmwareByRackNames(
	ctx context.Context,
	rackNames []string,
	componentType pb.ComponentType,
	startTime, endTime *time.Time,
) (*UpgradeFirmwareResult, error) {
	rackTargets := make([]*pb.RackTarget, 0, len(rackNames))
	for _, name := range rackNames {
		rt := &pb.RackTarget{
			Identifier: &pb.RackTarget_Name{Name: name},
		}
		if componentType != pb.ComponentType_COMPONENT_TYPE_UNKNOWN {
			rt.ComponentTypes = []pb.ComponentType{componentType}
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
		TaskIDs: protobuf.UUIDsFrom(rsp.GetTaskIds()),
	}, nil
}

// UpgradeFirmwareByMachineIDs upgrades firmware for the given machine IDs (external component IDs).
// Machine IDs are Carbide machine_id values, which are external system references for Compute components.
// Note: target_version is not yet supported in client API.
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
		TaskIDs: protobuf.UUIDsFrom(rsp.GetTaskIds()),
	}, nil
}

// PowerControlResult represents the result of a power control operation.
type PowerControlResult struct {
	TaskIDs []uuid.UUID // Multiple task IDs (1 task per rack)
}

// PowerControlByRackIDs performs power control on components in the given rack IDs.
func (c *Client) PowerControlByRackIDs(
	ctx context.Context,
	rackIDs []uuid.UUID,
	componentType pb.ComponentType,
	op pb.PowerControlOp,
) (*PowerControlResult, error) {
	rackTargets := make([]*pb.RackTarget, 0, len(rackIDs))
	for _, id := range rackIDs {
		rt := &pb.RackTarget{
			Identifier: &pb.RackTarget_Id{Id: id.String()},
		}
		if componentType != pb.ComponentType_COMPONENT_TYPE_UNKNOWN {
			rt.ComponentTypes = []pb.ComponentType{componentType}
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
	componentType pb.ComponentType,
	op pb.PowerControlOp,
) (*PowerControlResult, error) {
	rackTargets := make([]*pb.RackTarget, 0, len(rackNames))
	for _, name := range rackNames {
		rt := &pb.RackTarget{
			Identifier: &pb.RackTarget_Name{Name: name},
		}
		if componentType != pb.ComponentType_COMPONENT_TYPE_UNKNOWN {
			rt.ComponentTypes = []pb.ComponentType{componentType}
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

// PowerControlByMachineIDs performs power control on the given machine IDs (component IDs).
func (c *Client) PowerControlByMachineIDs(
	ctx context.Context,
	machineIDs []string,
	op pb.PowerControlOp,
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
	op pb.PowerControlOp,
) (*PowerControlResult, error) {
	var rsp *pb.SubmitTaskResponse
	var err error

	switch op {
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
		TaskIDs: protobuf.UUIDsFrom(rsp.GetTaskIds()),
	}, nil
}

// GetExpectedComponentsResult contains the result of GetExpectedComponents operation
type GetExpectedComponentsResult struct {
	Components []*component.Component
	Total      int
}

// GetExpectedComponentsByRackIDs retrieves expected components from local database by rack IDs.
func (c *Client) GetExpectedComponentsByRackIDs(
	ctx context.Context,
	rackIDs []uuid.UUID,
	componentType pb.ComponentType,
) (*GetExpectedComponentsResult, error) {
	rackTargets := make([]*pb.RackTarget, 0, len(rackIDs))
	for _, id := range rackIDs {
		rt := &pb.RackTarget{
			Identifier: &pb.RackTarget_Id{Id: id.String()},
		}
		if componentType != pb.ComponentType_COMPONENT_TYPE_UNKNOWN {
			rt.ComponentTypes = []pb.ComponentType{componentType}
		}
		rackTargets = append(rackTargets, rt)
	}

	rsp, err := c.client.GetExpectedComponents(
		ctx,
		&pb.GetExpectedComponentsRequest{
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

	return convertGetExpectedComponentsResponse(rsp), nil
}

// GetExpectedComponentsByRackNames retrieves expected components from local database by rack names.
func (c *Client) GetExpectedComponentsByRackNames(
	ctx context.Context,
	rackNames []string,
	componentType pb.ComponentType,
) (*GetExpectedComponentsResult, error) {
	rackTargets := make([]*pb.RackTarget, 0, len(rackNames))
	for _, name := range rackNames {
		rt := &pb.RackTarget{
			Identifier: &pb.RackTarget_Name{Name: name},
		}
		if componentType != pb.ComponentType_COMPONENT_TYPE_UNKNOWN {
			rt.ComponentTypes = []pb.ComponentType{componentType}
		}
		rackTargets = append(rackTargets, rt)
	}

	rsp, err := c.client.GetExpectedComponents(
		ctx,
		&pb.GetExpectedComponentsRequest{
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

	return convertGetExpectedComponentsResponse(rsp), nil
}

// GetExpectedComponentsByComponentIDs retrieves expected components from local database by external component IDs.
// Component IDs are external system IDs (e.g., Carbide machine_id for Compute components).
func (c *Client) GetExpectedComponentsByComponentIDs(
	ctx context.Context,
	componentIDs []string,
	componentType pb.ComponentType,
) (*GetExpectedComponentsResult, error) {
	compTargets := make([]*pb.ComponentTarget, 0, len(componentIDs))
	for _, compID := range componentIDs {
		compTargets = append(compTargets, &pb.ComponentTarget{
			Identifier: &pb.ComponentTarget_External{
				External: &pb.ExternalRef{
					Type: componentType,
					Id:   compID,
				},
			},
		})
	}

	rsp, err := c.client.GetExpectedComponents(
		ctx,
		&pb.GetExpectedComponentsRequest{
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

	return convertGetExpectedComponentsResponse(rsp), nil
}

func convertGetExpectedComponentsResponse(rsp *pb.GetExpectedComponentsResponse) *GetExpectedComponentsResult {
	components := make([]*component.Component, 0, len(rsp.Components))
	for _, c := range rsp.Components {
		components = append(components, protobuf.ComponentFrom(c))
	}

	return &GetExpectedComponentsResult{
		Components: components,
		Total:      int(rsp.Total),
	}
}

// GetActualComponentsResult represents the result of GetActualComponents call
type GetActualComponentsResult struct {
	Components []*pb.ActualComponent
	Total      int
}

// GetActualComponentsByRackIDs retrieves actual components from external systems by rack IDs.
// Currently only supports Compute component type.
func (c *Client) GetActualComponentsByRackIDs(
	ctx context.Context,
	rackIDs []uuid.UUID,
	componentType pb.ComponentType,
) (*GetActualComponentsResult, error) {
	rackTargets := make([]*pb.RackTarget, 0, len(rackIDs))
	for _, id := range rackIDs {
		rt := &pb.RackTarget{
			Identifier: &pb.RackTarget_Id{Id: id.String()},
		}
		if componentType != pb.ComponentType_COMPONENT_TYPE_UNKNOWN {
			rt.ComponentTypes = []pb.ComponentType{componentType}
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

	return &GetActualComponentsResult{
		Components: rsp.Components,
		Total:      int(rsp.Total),
	}, nil
}

// GetActualComponentsByRackNames retrieves actual components from external systems by rack names.
// Currently only supports Compute component type.
func (c *Client) GetActualComponentsByRackNames(
	ctx context.Context,
	rackNames []string,
	componentType pb.ComponentType,
) (*GetActualComponentsResult, error) {
	rackTargets := make([]*pb.RackTarget, 0, len(rackNames))
	for _, name := range rackNames {
		rt := &pb.RackTarget{
			Identifier: &pb.RackTarget_Name{Name: name},
		}
		if componentType != pb.ComponentType_COMPONENT_TYPE_UNKNOWN {
			rt.ComponentTypes = []pb.ComponentType{componentType}
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

	return &GetActualComponentsResult{
		Components: rsp.Components,
		Total:      int(rsp.Total),
	}, nil
}

// GetActualComponentsByComponentIDs retrieves actual components from external systems by external component IDs.
// Component IDs are external system IDs (e.g., Carbide machine_id for Compute components).
// Currently only supports Compute component type.
func (c *Client) GetActualComponentsByComponentIDs(
	ctx context.Context,
	componentIDs []string,
	componentType pb.ComponentType,
) (*GetActualComponentsResult, error) {
	compTargets := make([]*pb.ComponentTarget, 0, len(componentIDs))
	for _, compID := range componentIDs {
		compTargets = append(compTargets, &pb.ComponentTarget{
			Identifier: &pb.ComponentTarget_External{
				External: &pb.ExternalRef{
					Type: componentType,
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

	return &GetActualComponentsResult{
		Components: rsp.Components,
		Total:      int(rsp.Total),
	}, nil
}

// ValidateComponentsResult represents the result of ValidateComponents call
type ValidateComponentsResult struct {
	Diffs               []*pb.ComponentDiff
	TotalDiffs          int
	OnlyInExpectedCount int
	OnlyInActualCount   int
	DriftCount          int
	MatchCount          int
}

// ValidateComponentsByRackIDs validates expected vs actual components by rack IDs.
// Currently only supports Compute component type.
func (c *Client) ValidateComponentsByRackIDs(
	ctx context.Context,
	rackIDs []uuid.UUID,
	componentType pb.ComponentType,
) (*ValidateComponentsResult, error) {
	rackTargets := make([]*pb.RackTarget, 0, len(rackIDs))
	for _, id := range rackIDs {
		rt := &pb.RackTarget{
			Identifier: &pb.RackTarget_Id{Id: id.String()},
		}
		if componentType != pb.ComponentType_COMPONENT_TYPE_UNKNOWN {
			rt.ComponentTypes = []pb.ComponentType{componentType}
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
// Currently only supports Compute component type.
func (c *Client) ValidateComponentsByRackNames(
	ctx context.Context,
	rackNames []string,
	componentType pb.ComponentType,
) (*ValidateComponentsResult, error) {
	rackTargets := make([]*pb.RackTarget, 0, len(rackNames))
	for _, name := range rackNames {
		rt := &pb.RackTarget{
			Identifier: &pb.RackTarget_Name{Name: name},
		}
		if componentType != pb.ComponentType_COMPONENT_TYPE_UNKNOWN {
			rt.ComponentTypes = []pb.ComponentType{componentType}
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
// Component IDs are external system IDs (e.g., Carbide machine_id for Compute components).
// Currently only supports Compute component type.
func (c *Client) ValidateComponentsByComponentIDs(
	ctx context.Context,
	componentIDs []string,
	componentType pb.ComponentType,
) (*ValidateComponentsResult, error) {
	compTargets := make([]*pb.ComponentTarget, 0, len(componentIDs))
	for _, compID := range componentIDs {
		compTargets = append(compTargets, &pb.ComponentTarget{
			Identifier: &pb.ComponentTarget_External{
				External: &pb.ExternalRef{
					Type: componentType,
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
	return &ValidateComponentsResult{
		Diffs:               rsp.Diffs,
		TotalDiffs:          int(rsp.TotalDiffs),
		OnlyInExpectedCount: int(rsp.OnlyInExpectedCount),
		OnlyInActualCount:   int(rsp.OnlyInActualCount),
		DriftCount:          int(rsp.DriftCount),
		MatchCount:          int(rsp.MatchCount),
	}
}
