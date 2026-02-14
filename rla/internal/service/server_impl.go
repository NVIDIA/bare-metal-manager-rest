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
// Package service implements the gRPC server for the RLA (Rack Level Asset) management system.
// It provides APIs for managing rack-level assets including creating, retrieving, and updating
// rack and component information.
//
// TODO: This file is getting large. Consider splitting into multiple module files
package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/nvidia/bare-metal-manager-rest/rla/internal/carbideapi"
	inventorymanager "github.com/nvidia/bare-metal-manager-rest/rla/internal/inventory/manager"
	"github.com/nvidia/bare-metal-manager-rest/rla/internal/inventory/objects/bmc"
	"github.com/nvidia/bare-metal-manager-rest/rla/internal/inventory/objects/component"
	"github.com/nvidia/bare-metal-manager-rest/rla/internal/inventory/objects/rack"
	"github.com/nvidia/bare-metal-manager-rest/rla/internal/operation"
	pb "github.com/nvidia/bare-metal-manager-rest/rla/internal/proto/v1"
	"github.com/nvidia/bare-metal-manager-rest/rla/internal/psmapi"
	taskcommon "github.com/nvidia/bare-metal-manager-rest/rla/internal/task/common"
	taskmanager "github.com/nvidia/bare-metal-manager-rest/rla/internal/task/manager"
	operations "github.com/nvidia/bare-metal-manager-rest/rla/internal/task/operations"
	taskstore "github.com/nvidia/bare-metal-manager-rest/rla/internal/task/store"
	identifier "github.com/nvidia/bare-metal-manager-rest/rla/pkg/common/Identifier"
	"github.com/nvidia/bare-metal-manager-rest/rla/pkg/common/devicetypes"
	"github.com/nvidia/bare-metal-manager-rest/rla/pkg/converter/protobuf"
	"github.com/nvidia/bare-metal-manager-rest/rla/pkg/metadata"
)

// RLAServerImpl implements the gRPC RLA server interface.
// It acts as an adapter between gRPC protobuf messages and the internal managers,
// handling protobuf conversion and delegating business logic to the InventoryManager.
type RLAServerImpl struct {
	inventoryManager          inventorymanager.Manager // Business logic manager for inventory operations
	taskManager               *taskmanager.Manager     // Task manager for orchestrating task lifecycle
	taskStore                 taskstore.Store          // Task store for task queries
	carbideClient             carbideapi.Client        // Carbide API client for actual component data
	psmClient                 psmapi.Client            // PSM API client for powershelf operations
	pb.UnimplementedRLAServer                          // Embedded protobuf server interface for forward compatibility
}

// newServerImplementation creates a new RLA gRPC server implementation.
// It initializes the server with the provided managers for handling business logic.
//
// Parameters:
//   - inventoryManager: The inventory manager instance for handling rack and component topology
//   - taskManager: The Task manager for orchestrating task lifecycle
//   - taskStore: The task store for task queries
//   - carbideClient: The Carbide API client for actual component data
//   - psmClient: The PSM API client for powershelf operations
//
// Returns:
//   - *RLAServerImpl: A new server implementation instance
//   - error: Always nil in current implementation, reserved for future error handling
func newServerImplementation(
	inventoryManager inventorymanager.Manager,
	taskManager *taskmanager.Manager,
	taskStore taskstore.Store,
	carbideClient carbideapi.Client,
	psmClient psmapi.Client,
) (*RLAServerImpl, error) {
	return &RLAServerImpl{
		inventoryManager: inventoryManager,
		taskManager:      taskManager,
		taskStore:        taskStore,
		carbideClient:    carbideClient,
		psmClient:        psmClient,
	}, nil
}

// Version returns the build information for this RLA service.
// This includes the version, build time, and git commit hash.
func (rs *RLAServerImpl) Version(
	ctx context.Context,
	req *pb.VersionRequest,
) (*pb.BuildInfo, error) {
	return &pb.BuildInfo{
		Version:   metadata.Version,
		BuildTime: metadata.BuildTime,
		GitCommit: metadata.GitCommit,
	}, nil
}

// CreateExpectedRack creates a new expected rack configuration in the system.
// It converts the protobuf rack definition to internal format and stores it
// for later matching against physical rack discoveries.
//
// Parameters:
//   - ctx: Request context for cancellation and deadline management
//   - req: CreateExpectedRackRequest containing the rack configuration
//
// Returns:
//   - *pb.CreateExpectedRackResponse: Response containing the generated or
//     existing or given rack ID
//   - error: Any error that occurred during rack creation
func (rs *RLAServerImpl) CreateExpectedRack(
	ctx context.Context,
	req *pb.CreateExpectedRackRequest,
) (*pb.CreateExpectedRackResponse, error) {
	id, err := rs.inventoryManager.CreateExpectedRack(ctx, protobuf.RackFrom(req.GetRack()))

	return &pb.CreateExpectedRackResponse{Id: protobuf.UUIDTo(id)}, err
}

// GetRackInfoByID retrieves rack information by its unique identifier.
// Optionally includes component information if requested.
//
// Parameters:
//   - ctx: Request context for cancellation and deadline management
//   - req: GetRackInfoByIDRequest containing the rack ID and options
//
// Returns:
//   - *pb.GetRackInfoResponse: Response containing the rack information
//   - error: Any error that occurred during rack retrieval
func (rs *RLAServerImpl) GetRackInfoByID(
	ctx context.Context,
	req *pb.GetRackInfoByIDRequest,
) (*pb.GetRackInfoResponse, error) {
	r, err := rs.inventoryManager.GetRackByID(
		ctx,
		protobuf.UUIDFrom(req.GetId()),
		req.GetWithComponents(),
	)

	return &pb.GetRackInfoResponse{Rack: protobuf.RackTo(r)}, err
}

// GetRackInfoBySerial retrieves rack information by its manufacturer and serial number.
// This allows lookup of racks using their physical identification rather than system-generated ID.
// Optionally includes component information if requested.
//
// Parameters:
//   - ctx: Request context for cancellation and deadline management
//   - req: GetRackInfoBySerialRequest containing manufacturer, serial number, and options
//
// Returns:
//   - *pb.GetRackInfoResponse: Response containing the rack information
//   - error: Any error that occurred during rack retrieval
func (rs *RLAServerImpl) GetRackInfoBySerial(
	ctx context.Context,
	req *pb.GetRackInfoBySerialRequest,
) (*pb.GetRackInfoResponse, error) {
	r, err := rs.inventoryManager.GetRackBySerial(
		ctx,
		req.GetSerialInfo().GetManufacturer(),
		req.GetSerialInfo().GetSerialNumber(),
		req.GetWithComponents(),
	)

	return &pb.GetRackInfoResponse{Rack: protobuf.RackTo(r)}, err
}

// PatchRack updates an existing rack configuration with new information.
// This method performs intelligent merging of rack and component data, creating
// new components as needed and updating existing ones. It returns a detailed
// report of all operations performed during the patching process.
//
// Parameters:
//   - ctx: Request context for cancellation and deadline management
//   - req: PatchRackRequest containing the updated rack configuration
//
// Returns:
//   - *pb.PatchRackResponse: Response containing a JSON report of patch operations
//   - error: Any error that occurred during rack patching
func (rs *RLAServerImpl) PatchRack(
	ctx context.Context,
	req *pb.PatchRackRequest,
) (*pb.PatchRackResponse, error) {
	r := protobuf.RackFrom(req.GetRack())

	report, err := rs.inventoryManager.PatchRack(ctx, r)

	return &pb.PatchRackResponse{
		Report: report,
	}, err
}

// GetComponentInfoByID retrieves component information by its unique identifier.
// Optionally includes the parent rack information if requested. This method
// performs a two-step lookup: first retrieving the component and its rack ID,
// then fetching rack details if requested.
//
// Parameters:
//   - ctx: Request context for cancellation and deadline management
//   - req: GetComponentInfoByIDRequest containing the component ID and options
//
// Returns:
//   - *pb.GetComponentInfoResponse: Response containing component and optionally rack information
//   - error: Any error that occurred during component or rack retrieval
func (rs *RLAServerImpl) GetComponentInfoByID(
	ctx context.Context,
	req *pb.GetComponentInfoByIDRequest,
) (*pb.GetComponentInfoResponse, error) {
	c, err := rs.inventoryManager.GetComponentByID(
		ctx,
		protobuf.UUIDFrom(req.GetId()),
	)

	if err != nil {
		return nil, err
	}

	var r *rack.Rack

	if req.GetWithRack() {
		// Get the rack information
		r, err = rs.inventoryManager.GetRackByID(ctx, c.RackID, false)
		if err != nil {
			return nil, err
		}
	}

	return &pb.GetComponentInfoResponse{
		Component: protobuf.ComponentTo(c),
		Rack:      protobuf.RackTo(r),
	}, nil
}

// GetComponentInfoBySerial retrieves component information by its manufacturer and serial number.
// This allows lookup of components using their physical identification rather than system-generated ID.
// Optionally includes the parent rack information if requested. Like GetComponentInfoByID,
// this method performs a two-step lookup when rack information is requested.
//
// Parameters:
//   - ctx: Request context for cancellation and deadline management
//   - req: GetComponentInfoBySerialRequest containing manufacturer, serial number, and options
//
// Returns:
//   - *pb.GetComponentInfoResponse: Response containing component and optionally rack information
//   - error: Any error that occurred during component or rack retrieval
func (rs *RLAServerImpl) GetComponentInfoBySerial(
	ctx context.Context,
	req *pb.GetComponentInfoBySerialRequest,
) (*pb.GetComponentInfoResponse, error) {
	c, err := rs.inventoryManager.GetComponentBySerial(
		ctx,
		req.GetSerialInfo().GetManufacturer(),
		req.GetSerialInfo().GetSerialNumber(),
		req.GetWithRack(),
	)

	if err != nil {
		return nil, err
	}

	var r *rack.Rack

	if req.GetWithRack() {
		// Get the rack information
		r, err = rs.inventoryManager.GetRackByID(ctx, c.RackID, false)
		if err != nil {
			return nil, err
		}
	}

	return &pb.GetComponentInfoResponse{
		Component: protobuf.ComponentTo(c),
		Rack:      protobuf.RackTo(r),
	}, nil
}

func (rs *RLAServerImpl) GetListOfRacks(
	ctx context.Context,
	req *pb.GetListOfRacksRequest,
) (*pb.GetListOfRacksResponse, error) {
	pg := protobuf.PaginationFrom(req.GetPagination())
	if err := pg.Validate(); err != nil {
		return nil, fmt.Errorf("invalid pagination information: %v", err)
	}

	if req.GetInfo() == nil {
		return nil, errors.New("info is required")
	}

	racks, total, err := rs.inventoryManager.GetListOfRacks(
		ctx,
		*protobuf.StringQueryInfoFrom(req.GetInfo()),
		pg,
		req.GetWithComponents(),
	)

	results := make([]*pb.Rack, 0, len(racks))
	for _, r := range racks {
		results = append(results, protobuf.RackTo(r))
	}

	return &pb.GetListOfRacksResponse{
		Racks: results,
		Total: total,
	}, err
}

func (rs *RLAServerImpl) CreateNVLDomain(
	ctx context.Context,
	req *pb.CreateNVLDomainRequest,
) (*pb.CreateNVLDomainResponse, error) {
	id, err := rs.inventoryManager.CreateNVLDomain(
		ctx,
		protobuf.NVLDomainFrom(req.GetNvlDomain()),
	)

	return &pb.CreateNVLDomainResponse{Id: protobuf.UUIDTo(id)}, err
}

func (rs *RLAServerImpl) AttachRacksToNVLDomain(
	ctx context.Context,
	req *pb.AttachRacksToNVLDomainRequest,
) (*emptypb.Empty, error) {
	if req.GetNvlDomainIdentifier() == nil {
		return &emptypb.Empty{}, errors.New(
			"nvl domain identifier is required",
		)
	}

	if req.GetRackIdentifiers() == nil {
		// Nothing to do, return as no error.
		return &emptypb.Empty{}, nil
	}

	rackIDs := make([]identifier.Identifier, 0, len(req.GetRackIdentifiers()))
	for _, pbRackID := range req.GetRackIdentifiers() {
		rackIDs = append(rackIDs, *protobuf.IdentifierFrom(pbRackID))
	}

	return &emptypb.Empty{}, rs.inventoryManager.AttachRacksToNVLDomain(
		ctx,
		*protobuf.IdentifierFrom(req.GetNvlDomainIdentifier()),
		rackIDs,
	)
}

func (rs *RLAServerImpl) DetachRacksFromNVLDomain(
	ctx context.Context,
	req *pb.DetachRacksFromNVLDomainRequest,
) (*emptypb.Empty, error) {
	if req.GetRackIdentifiers() == nil {
		// Nothing to do, return as no error.
		return &emptypb.Empty{}, nil
	}

	rackIDs := make([]identifier.Identifier, 0, len(req.GetRackIdentifiers()))
	for _, pbRackID := range req.GetRackIdentifiers() {
		rackIDs = append(rackIDs, *protobuf.IdentifierFrom(pbRackID))
	}

	return &emptypb.Empty{}, rs.inventoryManager.DetachRacksFromNVLDomain(
		ctx,
		rackIDs,
	)
}

func (rs *RLAServerImpl) GetListOfNVLDomains(
	ctx context.Context,
	req *pb.GetListOfNVLDomainsRequest,
) (*pb.GetListOfNVLDomainsResponse, error) {
	pg := protobuf.PaginationFrom(req.GetPagination())
	if err := pg.Validate(); err != nil {
		return nil, fmt.Errorf("invalid pagination information: %v", err)
	}

	if req.GetInfo() == nil {
		return nil, errors.New("info is required")
	}

	nvlDomains, total, err := rs.inventoryManager.GetListOfNVLDomains(
		ctx,
		*protobuf.StringQueryInfoFrom(req.GetInfo()),
		pg,
	)

	results := make([]*pb.NVLDomain, 0, len(nvlDomains))
	for _, n := range nvlDomains {
		results = append(results, protobuf.NVLDomainTo(n))
	}

	return &pb.GetListOfNVLDomainsResponse{
		NvlDomains: results,
		Total:      total,
	}, err
}

func (rs *RLAServerImpl) GetRacksForNVLDomain(
	ctx context.Context,
	req *pb.GetRacksForNVLDomainRequest,
) (*pb.GetRacksForNVLDomainResponse, error) {
	if req.GetNvlDomainIdentifier() == nil {
		return nil, errors.New(
			"nvl domain identifier is required",
		)
	}

	racks, err := rs.inventoryManager.GetRacksForNVLDomain(
		ctx,
		*protobuf.IdentifierFrom(req.GetNvlDomainIdentifier()),
	)

	results := make([]*pb.Rack, 0, len(racks))
	for _, r := range racks {
		results = append(results, protobuf.RackTo(r))
	}

	return &pb.GetRacksForNVLDomainResponse{Racks: results}, err
}

func (rs *RLAServerImpl) PowerOnRack(
	ctx context.Context,
	req *pb.PowerOnRackRequest,
) (*pb.SubmitTaskResponse, error) {
	return rs.handlePowerControlTask(
		ctx,
		req.GetTargetSpec(),
		req.GetDescription(),
		&operations.PowerControlTaskInfo{
			Operation: operations.PowerOperationPowerOn,
		},
	)
}

func (rs *RLAServerImpl) PowerOffRack(
	ctx context.Context,
	req *pb.PowerOffRackRequest,
) (*pb.SubmitTaskResponse, error) {
	op := operations.PowerOperationPowerOff
	if req.GetForced() {
		op = operations.PowerOperationForcePowerOff
	}
	return rs.handlePowerControlTask(
		ctx,
		req.GetTargetSpec(),
		req.GetDescription(),
		&operations.PowerControlTaskInfo{
			Operation: op,
			Forced:    req.GetForced(),
		},
	)
}

func (rs *RLAServerImpl) PowerResetRack(
	ctx context.Context,
	req *pb.PowerResetRackRequest,
) (*pb.SubmitTaskResponse, error) {
	op := operations.PowerOperationRestart
	if req.GetForced() {
		op = operations.PowerOperationForceRestart
	}
	return rs.handlePowerControlTask(
		ctx,
		req.GetTargetSpec(),
		req.GetDescription(),
		&operations.PowerControlTaskInfo{
			Operation: op,
			Forced:    req.GetForced(),
		},
	)
}

func (rs *RLAServerImpl) handlePowerControlTask(
	ctx context.Context,
	targetSpec *pb.OperationTargetSpec,
	description string,
	info *operations.PowerControlTaskInfo,
) (*pb.SubmitTaskResponse, error) {
	if rs.taskManager == nil {
		return nil, errors.New("task manager is not available")
	}

	if targetSpec == nil {
		return nil, errors.New("target_spec is required")
	}

	// Convert pb.OperationTargetSpec to internal operation.Request
	req, err := rs.convertTargetSpecToOperationRequest(targetSpec, description, info)
	if err != nil {
		return nil, err
	}

	// Task Manager handles resolve + split by rack + create tasks
	taskIDs, err := rs.taskManager.SubmitTask(ctx, req)
	if err != nil {
		return nil, err
	}

	if len(taskIDs) == 0 {
		return nil, errors.New("failed to create any tasks")
	}

	return &pb.SubmitTaskResponse{TaskIds: protobuf.UUIDsTo(taskIDs)}, nil
}

// convertTargetSpecToOperationRequest converts pb.OperationTargetSpec to internal operation.Request.
// This is a simple translation without DB queries. Task Manager handles resolve + split.
func (rs *RLAServerImpl) convertTargetSpecToOperationRequest(
	targetSpec *pb.OperationTargetSpec,
	description string,
	info operations.Operation,
) (*operation.Request, error) {
	raw, err := info.Marshal()
	if err != nil {
		return nil, err
	}

	req := &operation.Request{
		Operation: operation.Wrapper{
			Type: info.Type(),
			Info: raw,
		},
		Description: description,
	}

	// Convert pb targets to internal targets based on the oneof type
	switch targets := targetSpec.GetTargets().(type) {
	case *pb.OperationTargetSpec_Racks:
		for _, pbRack := range targets.Racks.GetTargets() {
			rackTarget, err := rs.convertPbRackTargetToRackTarget(pbRack)
			if err != nil {
				return nil, fmt.Errorf("failed to convert rack target: %w", err)
			}
			req.TargetSpec.Racks = append(req.TargetSpec.Racks, *rackTarget)
		}

	case *pb.OperationTargetSpec_Components:
		for _, pbComp := range targets.Components.GetTargets() {
			compTarget, err := rs.convertPbComponentTargetToComponentTarget(pbComp)
			if err != nil {
				return nil, fmt.Errorf("failed to convert component target: %w", err)
			}
			req.TargetSpec.Components = append(req.TargetSpec.Components, *compTarget)
		}

	default:
		return nil, errors.New("target_spec must have either racks or components set")
	}

	return req, nil
}

// convertPbRackTargetToRackTarget converts a protobuf RackTarget to an internal RackTarget
func (rs *RLAServerImpl) convertPbRackTargetToRackTarget(rt *pb.RackTarget) (*operation.RackTarget, error) {
	if rt == nil {
		return nil, fmt.Errorf("rack target is nil")
	}

	target := &operation.RackTarget{}

	switch id := rt.GetIdentifier().(type) {
	case *pb.RackTarget_Id:
		parsed, err := uuid.Parse(id.Id)
		if err != nil {
			return nil, fmt.Errorf("invalid rack id %q: %w", id.Id, err)
		}
		target.Identifier.ID = parsed

	case *pb.RackTarget_Name:
		target.Identifier.Name = id.Name

	default:
		return nil, fmt.Errorf("rack target must have either id or name set")
	}

	// Convert component types
	for _, pbType := range rt.GetComponentTypes() {
		target.ComponentTypes = append(target.ComponentTypes, protobuf.ComponentTypeFrom(pbType))
	}

	return target, nil
}

// convertPbComponentTargetToComponentTarget converts a protobuf ComponentTarget to an internal ComponentTarget
func (rs *RLAServerImpl) convertPbComponentTargetToComponentTarget(ct *pb.ComponentTarget) (*operation.ComponentTarget, error) {
	if ct == nil {
		return nil, fmt.Errorf("component target is nil")
	}

	target := &operation.ComponentTarget{}

	switch id := ct.GetIdentifier().(type) {
	case *pb.ComponentTarget_Uuid:
		parsed, err := uuid.Parse(id.Uuid)
		if err != nil {
			return nil, fmt.Errorf("invalid component uuid %q: %w", id.Uuid, err)
		}
		target.UUID = parsed

	case *pb.ComponentTarget_External:
		target.External = &operation.ExternalRef{
			Type: protobuf.ComponentTypeFrom(id.External.GetType()),
			ID:   id.External.GetId(),
		}

	default:
		return nil, fmt.Errorf("component target must have either uuid or external set")
	}

	return target, nil
}

func (rs *RLAServerImpl) ListTasks(
	ctx context.Context,
	req *pb.ListTasksRequest,
) (*pb.ListTasksResponse, error) {
	options := &taskcommon.TaskListOptions{
		TaskType:   taskcommon.TaskTypeUnknown,
		RackID:     protobuf.UUIDFrom(req.GetRackId()),
		ActiveOnly: req.GetActiveOnly(),
	}

	pagination := protobuf.PaginationFrom(req.GetPagination())
	if err := pagination.Validate(); err != nil {
		return nil, fmt.Errorf("invalid pagination information: %v", err)
	}

	tasks, total, err := rs.taskStore.ListTasks(ctx, options, pagination)
	if err != nil {
		return nil, err
	}

	results := make([]*pb.Task, 0, len(tasks))
	for _, t := range tasks {
		results = append(results, protobuf.TaskTo(t))
	}

	return &pb.ListTasksResponse{Tasks: results, Total: total}, nil
}

func (rs *RLAServerImpl) GetTasksByIDs(
	ctx context.Context,
	req *pb.GetTasksByIDsRequest,
) (*pb.GetTasksByIDsResponse, error) {
	taskIDs := make([]uuid.UUID, 0, len(req.GetTaskIds()))
	for _, tid := range req.GetTaskIds() {
		taskIDs = append(taskIDs, protobuf.UUIDFrom(tid))
	}

	tasks, err := rs.taskStore.GetTasks(ctx, taskIDs)
	if err != nil {
		return nil, err
	}

	results := make([]*pb.Task, 0, len(tasks))
	for _, t := range tasks {
		results = append(results, protobuf.TaskTo(t))
	}

	return &pb.GetTasksByIDsResponse{Tasks: results}, nil
}

// UpgradeFirmware upgrades firmware for components.
// It uses OperationTargetSpec to specify targets and creates a task via the Task framework.
func (rs *RLAServerImpl) UpgradeFirmware(
	ctx context.Context,
	req *pb.UpgradeFirmwareRequest,
) (*pb.SubmitTaskResponse, error) {
	if rs.taskManager == nil {
		return nil, errors.New("task manager is not available")
	}

	targetSpec := req.GetTargetSpec()
	if targetSpec == nil {
		return nil, errors.New("target_spec is required")
	}

	// Build FirmwareControlTaskInfo
	info := &operations.FirmwareControlTaskInfo{
		Operation:     operations.FirmwareOperationUpgrade,
		TargetVersion: req.GetTargetVersion(),
	}

	// Parse optional time parameters for scheduled upgrade
	if req.GetStartTime() != nil {
		info.StartTime = req.GetStartTime().AsTime().Unix()
	}
	if req.GetEndTime() != nil {
		info.EndTime = req.GetEndTime().AsTime().Unix()
	}

	// Convert pb.OperationTargetSpec to internal operation.Request
	opReq, err := rs.convertTargetSpecToOperationRequest(targetSpec, req.GetDescription(), info)
	if err != nil {
		return nil, err
	}

	// Task Manager handles resolve + split by rack + create tasks
	taskIDs, err := rs.taskManager.SubmitTask(ctx, opReq)
	if err != nil {
		return nil, err
	}

	if len(taskIDs) == 0 {
		return nil, errors.New("failed to create any tasks")
	}

	return &pb.SubmitTaskResponse{TaskIds: protobuf.UUIDsTo(taskIDs)}, nil
}

// GetExpectedComponents retrieves expected components from local database.
// It accepts OperationTargetSpec with rack targets (with optional type filter) or component targets.
func (rs *RLAServerImpl) GetExpectedComponents(
	ctx context.Context,
	req *pb.GetExpectedComponentsRequest,
) (*pb.GetExpectedComponentsResponse, error) {
	// Extract components using the shared helper
	components, err := rs.extractComponentsFromTargetSpec(ctx, req.GetTargetSpec())
	if err != nil {
		return nil, err
	}

	// Convert to protobuf components
	pbComponents := make([]*pb.Component, 0, len(components))
	for _, comp := range components {
		pbComponents = append(pbComponents, protobuf.ComponentTo(comp))
	}

	// TODO: Apply pagination if provided
	// pagination := req.GetPagination()

	return &pb.GetExpectedComponentsResponse{
		Components: pbComponents,
		Total:      int32(len(pbComponents)),
	}, nil
}

// GetActualComponents returns actual component data from external systems (e.g., Carbide, PSM).
// Currently supports Compute and Powershelf component types.
//
// TODO: Currently finding actual data by component_id. Consider using MAC or serial number instead.
func (rs *RLAServerImpl) GetActualComponents(
	ctx context.Context,
	req *pb.GetActualComponentsRequest,
) (*pb.GetActualComponentsResponse, error) {
	targetSpec := req.GetTargetSpec()
	if targetSpec == nil {
		return nil, errors.New("target_spec is required")
	}

	// First, get expected components
	expectedComponents, err := rs.extractComponentsFromTargetSpec(ctx, targetSpec)
	if err != nil {
		return nil, fmt.Errorf("failed to extract components: %w", err)
	}

	if len(expectedComponents) == 0 {
		return &pb.GetActualComponentsResponse{
			Components: []*pb.ActualComponent{},
			Total:      0,
		}, nil
	}

	// Group components by type
	componentsByType := make(map[devicetypes.ComponentType][]*component.Component)
	for _, comp := range expectedComponents {
		componentsByType[comp.Type] = append(componentsByType[comp.Type], comp)
	}

	// Process each type and collect results
	var allActualComponents []*pb.ActualComponent
	for compType, comps := range componentsByType {
		switch compType {
		case devicetypes.ComponentTypeCompute:
			resp, err := rs.getActualComputeComponents(ctx, comps)
			if err != nil {
				return nil, err
			}
			allActualComponents = append(allActualComponents, resp.GetComponents()...)
		case devicetypes.ComponentTypePowerShelf:
			resp, err := rs.getActualPowershelfComponents(ctx, comps)
			if err != nil {
				return nil, err
			}
			allActualComponents = append(allActualComponents, resp.GetComponents()...)
		default:
			return nil, fmt.Errorf("component type %s is not supported for GetActualComponents, only Compute and Powershelf types are supported",
				devicetypes.ComponentTypeToString(compType))
		}
	}

	return &pb.GetActualComponentsResponse{
		Components: allActualComponents,
		Total:      int32(len(allActualComponents)),
	}, nil
}

// getActualComputeComponents returns actual compute component data from Carbide.
func (rs *RLAServerImpl) getActualComputeComponents(
	ctx context.Context,
	expectedComponents []*component.Component,
) (*pb.GetActualComponentsResponse, error) {
	if rs.carbideClient == nil {
		return nil, errors.New("carbide client is not available")
	}

	// Extract component IDs from expected components
	var componentIDs []string
	componentIDToRackID := make(map[string]uuid.UUID) // Map component_id -> rack_id for lookup
	for _, comp := range expectedComponents {
		if comp.ComponentID != "" {
			componentIDs = append(componentIDs, comp.ComponentID)
			componentIDToRackID[comp.ComponentID] = comp.RackID
		}
	}

	if len(componentIDs) == 0 {
		return &pb.GetActualComponentsResponse{
			Components: []*pb.ActualComponent{},
			Total:      0,
		}, nil
	}

	// Query Carbide for machine details (Carbide uses "machine" terminology)
	machineDetails, err := rs.carbideClient.FindMachinesByIds(ctx, componentIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to get machine details from Carbide: %w", err)
	}

	// Query Carbide for machine positions
	machinePositions, err := rs.carbideClient.GetMachinePositionInfo(ctx, componentIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to get machine positions from Carbide: %w", err)
	}

	// Query Carbide for power states
	powerStates, err := rs.carbideClient.GetPowerStates(ctx, componentIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to get power states from Carbide: %w", err)
	}

	// Build position map for quick lookup
	positionMap := make(map[string]carbideapi.MachinePosition)
	for _, pos := range machinePositions {
		positionMap[pos.MachineID] = pos
	}

	// Build power state map for quick lookup
	powerStateMap := make(map[string]carbideapi.PowerState)
	for _, ps := range powerStates {
		powerStateMap[ps.MachineID] = ps.PowerState
	}

	// Convert to ActualComponent
	actualComponents := make([]*pb.ActualComponent, 0, len(machineDetails))
	for _, detail := range machineDetails {
		ac := &pb.ActualComponent{
			Type:            pb.ComponentType_COMPONENT_TYPE_COMPUTE,
			ComponentId:     detail.MachineID,
			FirmwareVersion: detail.FirmwareVersion,
			HealthStatus:    detail.HealthStatus,
			Source:          "carbide",
		}

		// Set info
		ac.Info = &pb.DeviceInfo{}
		if detail.ChassisSerial != nil {
			ac.Info.SerialNumber = *detail.ChassisSerial
		}

		// Set rack_id from our local lookup
		if rackID, ok := componentIDToRackID[detail.MachineID]; ok {
			ac.RackId = protobuf.UUIDTo(rackID)
		}

		// Set position from Carbide
		if pos, ok := positionMap[detail.MachineID]; ok {
			ac.Position = &pb.RackPosition{}
			if pos.PhysicalSlotNum != nil {
				ac.Position.SlotId = *pos.PhysicalSlotNum
			}
			if pos.ComputeTrayIndex != nil {
				ac.Position.TrayIdx = *pos.ComputeTrayIndex
			}
			if pos.TopologyID != nil {
				ac.Position.HostId = *pos.TopologyID
			}
		}

		// Set BMC info
		if detail.BmcIP != "" || detail.BmcMac != "" {
			ac.Bmcs = []*pb.BMCInfo{
				{
					IpAddress:  &detail.BmcIP,
					MacAddress: detail.BmcMac,
				},
			}
		}

		// Set last seen time
		if detail.LastObservationTime != nil {
			ac.LastSeen = timestamppb.New(*detail.LastObservationTime)
		}

		// Set power state from Carbide
		if ps, ok := powerStateMap[detail.MachineID]; ok {
			ac.PowerState = powerStateToString(ps)
		}

		actualComponents = append(actualComponents, ac)
	}

	return &pb.GetActualComponentsResponse{
		Components: actualComponents,
		Total:      int32(len(actualComponents)),
	}, nil
}

// getActualPowershelfComponents returns actual powershelf component data from PSM.
func (rs *RLAServerImpl) getActualPowershelfComponents(
	ctx context.Context,
	expectedComponents []*component.Component,
) (*pb.GetActualComponentsResponse, error) {
	// Extract PMC MAC addresses from expected components
	var pmcMacs []string
	pmcMacToRackID := make(map[string]uuid.UUID) // Map pmc_mac -> rack_id for lookup
	pmcMacToPosition := make(map[string]component.InRackPosition)

	for _, comp := range expectedComponents {
		bmcs := getBmcsFromComponent(comp)
		if len(bmcs) != 1 {
			return nil, fmt.Errorf("component %s: expected exactly 1 BMC, found %d", comp.Info.Name, len(bmcs))
		}
		if len(bmcs[0].MAC) == 0 {
			return nil, fmt.Errorf("component %s: BMC has invalid MAC address", comp.Info.Name)
		}

		// For powershelves, the first BMC is the PMC
		pmcMac := bmcs[0].MAC.String()
		if pmcMac != "" {
			pmcMacs = append(pmcMacs, pmcMac)
			pmcMacToRackID[pmcMac] = comp.RackID
			pmcMacToPosition[pmcMac] = comp.Position
		}
	}

	if len(pmcMacs) == 0 {
		return &pb.GetActualComponentsResponse{
			Components: []*pb.ActualComponent{},
			Total:      0,
		}, nil
	}

	if rs.psmClient == nil {
		return nil, errors.New("PSM client is not available")
	}

	// Query PSM for powershelf details
	powershelves, err := rs.psmClient.GetPowershelves(ctx, pmcMacs)
	if err != nil {
		return nil, fmt.Errorf("failed to get powershelf details from PSM: %w", err)
	}

	// Convert to ActualComponent
	actualComponents := make([]*pb.ActualComponent, 0, len(powershelves))
	for _, ps := range powershelves {
		pmcMac := ps.PMC.MACAddress

		ac := &pb.ActualComponent{
			Type:            pb.ComponentType_COMPONENT_TYPE_POWERSHELF,
			ComponentId:     pmcMac, // Use PMC MAC as component ID
			FirmwareVersion: ps.PMC.FirmwareVersion,
			Source:          "psm",
		}

		// Set info from chassis
		ac.Info = &pb.DeviceInfo{
			SerialNumber: ps.Chassis.SerialNumber,
			Model:        &ps.Chassis.Model,
			Manufacturer: ps.Chassis.Manufacturer,
		}

		// Set rack_id from our local lookup
		if rackID, ok := pmcMacToRackID[pmcMac]; ok {
			ac.RackId = protobuf.UUIDTo(rackID)
		}

		// Set position from expected component
		if pos, ok := pmcMacToPosition[pmcMac]; ok {
			ac.Position = &pb.RackPosition{
				SlotId:  int32(pos.SlotID),
				TrayIdx: int32(pos.TrayIndex),
				HostId:  int32(pos.HostID),
			}
		}

		// Set BMC info (PMC)
		ac.Bmcs = []*pb.BMCInfo{
			{
				IpAddress:  &ps.PMC.IPAddress,
				MacAddress: ps.PMC.MACAddress,
			},
		}

		// Determine power state from PSUs
		// If any PSU is powered on, consider the powershelf as "on"
		powerState := "off"
		for _, psu := range ps.PSUs {
			if psu.PowerState {
				powerState = "on"
				break
			}
		}
		ac.PowerState = powerState

		actualComponents = append(actualComponents, ac)
	}

	return &pb.GetActualComponentsResponse{
		Components: actualComponents,
		Total:      int32(len(actualComponents)),
	}, nil
}

// powerStateToString converts a Carbide PowerState to a string representation
func powerStateToString(ps carbideapi.PowerState) string {
	switch ps {
	case carbideapi.PowerStateOn:
		return "on"
	case carbideapi.PowerStateOff:
		return "off"
	case carbideapi.PowerStateDisabled:
		return "disabled"
	default:
		return "unknown"
	}
}

// ValidateComponents compares expected (local DB) vs actual (from external source systems) components.
// Returns validation result indicating if components match and details of any differences.
// Each component type queries its own source system (e.g., Carbide for Compute, PSM for PowerShelf).
// Currently only supports Compute component type.
//
// TODO: Currently matching by component_id. Consider using MAC or serial number instead.
func (rs *RLAServerImpl) ValidateComponents(
	ctx context.Context,
	req *pb.ValidateComponentsRequest,
) (*pb.ValidateComponentsResponse, error) {
	targetSpec := req.GetTargetSpec()
	if targetSpec == nil {
		return nil, errors.New("target_spec is required")
	}

	if rs.carbideClient == nil {
		return nil, errors.New("carbide client is not available")
	}

	// Get expected components from local DB
	expectedComponents, err := rs.extractComponentsFromTargetSpec(ctx, targetSpec)
	if err != nil {
		return nil, fmt.Errorf("failed to extract expected components: %w", err)
	}

	// Currently only Compute type is supported - check all components are Compute
	for _, comp := range expectedComponents {
		if comp.Type != devicetypes.ComponentTypeCompute {
			return nil, fmt.Errorf("component type %s is not supported for ValidateComponents, only Compute is supported",
				devicetypes.ComponentTypeToString(comp.Type))
		}
	}

	// Build expected component map by component_id
	expectedByComponentID := make(map[string]*component.Component)
	var componentIDs []string
	for _, comp := range expectedComponents {
		if comp.ComponentID != "" {
			expectedByComponentID[comp.ComponentID] = comp
			componentIDs = append(componentIDs, comp.ComponentID)
		}
	}

	// TODO: Refactor - extract this source system querying logic to a shared helper function
	// to avoid duplication with GetActualComponents. Each component type should query its own
	// source system (Carbide for Compute, PSM for PowerShelf, etc.)

	// Query source system for actual data (currently Carbide for Compute type)
	var machineDetails []carbideapi.MachineDetail
	var machinePositions []carbideapi.MachinePosition
	var powerStates []carbideapi.MachinePowerState

	if len(componentIDs) > 0 {
		machineDetails, err = rs.carbideClient.FindMachinesByIds(ctx, componentIDs)
		if err != nil {
			return nil, fmt.Errorf("failed to get machine details from Carbide: %w", err)
		}

		machinePositions, err = rs.carbideClient.GetMachinePositionInfo(ctx, componentIDs)
		if err != nil {
			return nil, fmt.Errorf("failed to get machine positions from Carbide: %w", err)
		}

		powerStates, err = rs.carbideClient.GetPowerStates(ctx, componentIDs)
		if err != nil {
			return nil, fmt.Errorf("failed to get power states from Carbide: %w", err)
		}
	}

	// Build lookup maps for actual data (keyed by Carbide's MachineID which is our ComponentID)
	positionMap := make(map[string]carbideapi.MachinePosition)
	for _, pos := range machinePositions {
		positionMap[pos.MachineID] = pos
	}

	powerStateMap := make(map[string]carbideapi.PowerState)
	for _, ps := range powerStates {
		powerStateMap[ps.MachineID] = ps.PowerState
	}

	actualByComponentID := make(map[string]carbideapi.MachineDetail)
	for _, detail := range machineDetails {
		actualByComponentID[detail.MachineID] = detail
	}

	// Compare and generate diffs
	var diffs []*pb.ComponentDiff
	var onlyInExpectedCount, onlyInActualCount, driftCount, matchCount int32

	// Check expected components
	for compID, expected := range expectedByComponentID {
		actual, existsInActual := actualByComponentID[compID]

		if !existsInActual {
			// Only in expected (missing from source system)
			diffs = append(diffs, &pb.ComponentDiff{
				Type:        pb.DiffType_DIFF_TYPE_ONLY_IN_EXPECTED,
				ComponentId: compID,
				Expected:    protobuf.ComponentTo(expected),
			})
			onlyInExpectedCount++
		} else {
			// Exists in both - check for drift
			fieldDiffs := compareComponentFields(expected, actual, positionMap[compID])

			if len(fieldDiffs) > 0 {
				diffs = append(diffs, &pb.ComponentDiff{
					Type:        pb.DiffType_DIFF_TYPE_DRIFT,
					ComponentId: compID,
					FieldDiffs:  fieldDiffs,
				})
				driftCount++
			} else {
				matchCount++
			}
		}
	}

	// Check for components only in actual (not in expected)
	for compID, actual := range actualByComponentID {
		if _, existsInExpected := expectedByComponentID[compID]; !existsInExpected {
			// Only in actual (exists in source system but not in local DB)
			ac := buildActualComponentFromDetail(actual, positionMap[compID], powerStateMap[compID], uuid.Nil)
			diffs = append(diffs, &pb.ComponentDiff{
				Type:        pb.DiffType_DIFF_TYPE_ONLY_IN_ACTUAL,
				ComponentId: compID,
				Actual:      ac,
			})
			onlyInActualCount++
		}
	}

	return &pb.ValidateComponentsResponse{
		Diffs:               diffs,
		TotalDiffs:          int32(len(diffs)),
		OnlyInExpectedCount: onlyInExpectedCount,
		OnlyInActualCount:   onlyInActualCount,
		DriftCount:          driftCount,
		MatchCount:          matchCount,
	}, nil
}

// compareComponentFields compares expected and actual component fields and returns differences.
// Does not compare frequently changing fields like power_state and health_status.
func compareComponentFields(
	expected *component.Component,
	actual carbideapi.MachineDetail,
	position carbideapi.MachinePosition,
) []*pb.FieldDiff {
	var diffs []*pb.FieldDiff

	// Compare position.slot_id
	if position.PhysicalSlotNum != nil && expected.Position.SlotID != int(*position.PhysicalSlotNum) {
		diffs = append(diffs, &pb.FieldDiff{
			FieldName:     "position.slot_id",
			ExpectedValue: fmt.Sprintf("%d", expected.Position.SlotID),
			ActualValue:   fmt.Sprintf("%d", *position.PhysicalSlotNum),
		})
	}

	// Compare position.tray_idx
	if position.ComputeTrayIndex != nil && expected.Position.TrayIndex != int(*position.ComputeTrayIndex) {
		diffs = append(diffs, &pb.FieldDiff{
			FieldName:     "position.tray_idx",
			ExpectedValue: fmt.Sprintf("%d", expected.Position.TrayIndex),
			ActualValue:   fmt.Sprintf("%d", *position.ComputeTrayIndex),
		})
	}

	// Compare position.host_id
	if position.TopologyID != nil && expected.Position.HostID != int(*position.TopologyID) {
		diffs = append(diffs, &pb.FieldDiff{
			FieldName:     "position.host_id",
			ExpectedValue: fmt.Sprintf("%d", expected.Position.HostID),
			ActualValue:   fmt.Sprintf("%d", *position.TopologyID),
		})
	}

	// Compare firmware_version
	if expected.FirmwareVersion != "" && actual.FirmwareVersion != "" &&
		expected.FirmwareVersion != actual.FirmwareVersion {
		diffs = append(diffs, &pb.FieldDiff{
			FieldName:     "firmware_version",
			ExpectedValue: expected.FirmwareVersion,
			ActualValue:   actual.FirmwareVersion,
		})
	}

	// Compare serial_number (chassis_serial)
	if actual.ChassisSerial != nil && expected.Info.SerialNumber != "" &&
		expected.Info.SerialNumber != *actual.ChassisSerial {
		diffs = append(diffs, &pb.FieldDiff{
			FieldName:     "serial_number",
			ExpectedValue: expected.Info.SerialNumber,
			ActualValue:   *actual.ChassisSerial,
		})
	}

	return diffs
}

// buildActualComponentFromDetail builds an ActualComponent proto from source system data.
// TODO: Refactor - consolidate with similar logic in GetActualComponents.
// This should be generalized to support different source systems per component type.
func buildActualComponentFromDetail(
	detail carbideapi.MachineDetail,
	position carbideapi.MachinePosition,
	powerState carbideapi.PowerState,
	rackID uuid.UUID,
) *pb.ActualComponent {
	ac := &pb.ActualComponent{
		Type:            pb.ComponentType_COMPONENT_TYPE_COMPUTE,
		ComponentId:     detail.MachineID,
		FirmwareVersion: detail.FirmwareVersion,
		HealthStatus:    detail.HealthStatus,
		Source:          "carbide",
		PowerState:      powerStateToString(powerState),
	}

	ac.Info = &pb.DeviceInfo{}
	if detail.ChassisSerial != nil {
		ac.Info.SerialNumber = *detail.ChassisSerial
	}

	if rackID != uuid.Nil {
		ac.RackId = protobuf.UUIDTo(rackID)
	}

	ac.Position = &pb.RackPosition{}
	if position.PhysicalSlotNum != nil {
		ac.Position.SlotId = *position.PhysicalSlotNum
	}
	if position.ComputeTrayIndex != nil {
		ac.Position.TrayIdx = *position.ComputeTrayIndex
	}
	if position.TopologyID != nil {
		ac.Position.HostId = *position.TopologyID
	}

	if detail.BmcIP != "" || detail.BmcMac != "" {
		ac.Bmcs = []*pb.BMCInfo{
			{
				IpAddress:  &detail.BmcIP,
				MacAddress: detail.BmcMac,
			},
		}
	}

	if detail.LastObservationTime != nil {
		ac.LastSeen = timestamppb.New(*detail.LastObservationTime)
	}

	return ac
}

// getBmcsFromComponent returns all BMCs found in the component.
func getBmcsFromComponent(comp *component.Component) []bmc.BMC {
	var result []bmc.BMC
	for _, bmcs := range comp.BmcsByType {
		result = append(result, bmcs...)
	}
	return result
}

// now returns the current time (extracted for testing)
var now = time.Now

// extractComponentsByType extracts components of the specified type from a rack
func extractComponentsByType(r *rack.Rack, compType devicetypes.ComponentType) []*component.Component {
	var result []*component.Component
	for i := range r.Components {
		if r.Components[i].Type == compType {
			result = append(result, &r.Components[i])
		}
	}
	return result
}

// extractComponentsByTypes extracts components matching any of the specified types from a rack.
// If compTypes is empty, returns all components.
func extractComponentsByTypes(r *rack.Rack, compTypes []devicetypes.ComponentType) []*component.Component {
	var result []*component.Component

	// If no types specified, return all components
	if len(compTypes) == 0 {
		for i := range r.Components {
			result = append(result, &r.Components[i])
		}
		return result
	}

	// Build a set of allowed types for O(1) lookup
	allowedTypes := make(map[devicetypes.ComponentType]bool)
	for _, ct := range compTypes {
		allowedTypes[ct] = true
	}

	for i := range r.Components {
		if allowedTypes[r.Components[i].Type] {
			result = append(result, &r.Components[i])
		}
	}
	return result
}

// extractComponentsFromTargetSpec extracts components from an OperationTargetSpec message.
// It handles rack targets (with optional type filtering) or component targets (by UUID or external ref).
func (rs *RLAServerImpl) extractComponentsFromTargetSpec(
	ctx context.Context,
	targetSpec *pb.OperationTargetSpec,
) ([]*component.Component, error) {
	if targetSpec == nil {
		return nil, errors.New("target_spec is required")
	}

	var components []*component.Component

	switch targets := targetSpec.GetTargets().(type) {
	case *pb.OperationTargetSpec_Racks:
		if len(targets.Racks.GetTargets()) == 0 {
			return nil, errors.New("racks targets is empty")
		}
		for _, rt := range targets.Racks.GetTargets() {
			resolved, err := rs.extractComponentsFromRackTarget(ctx, rt)
			if err != nil {
				return nil, err
			}
			components = append(components, resolved...)
		}

	case *pb.OperationTargetSpec_Components:
		if len(targets.Components.GetTargets()) == 0 {
			return nil, errors.New("components targets is empty")
		}
		for _, ct := range targets.Components.GetTargets() {
			resolved, err := rs.extractComponentsFromComponentTarget(ctx, ct)
			if err != nil {
				return nil, err
			}
			components = append(components, resolved...)
		}

	default:
		return nil, errors.New("target_spec must have either racks or components set")
	}

	return components, nil
}

// extractComponentsFromRackTarget extracts components from a rack target.
func (rs *RLAServerImpl) extractComponentsFromRackTarget(
	ctx context.Context,
	rt *pb.RackTarget,
) ([]*component.Component, error) {
	var r *rack.Rack
	var err error

	switch id := rt.GetIdentifier().(type) {
	case *pb.RackTarget_Id:
		rackUUID, parseErr := uuid.Parse(id.Id)
		if parseErr != nil {
			return nil, fmt.Errorf("invalid rack id %q: %w", id.Id, parseErr)
		}
		r, err = rs.inventoryManager.GetRackByID(ctx, rackUUID, true)
	case *pb.RackTarget_Name:
		r, err = rs.inventoryManager.GetRackByIdentifier(ctx, identifier.Identifier{Name: id.Name}, true)
	default:
		return nil, errors.New("rack target must have either id or name set")
	}

	if err != nil {
		return nil, err
	}

	// Convert pb component types to internal types
	compTypes := make([]devicetypes.ComponentType, 0, len(rt.GetComponentTypes()))
	for _, pbType := range rt.GetComponentTypes() {
		compTypes = append(compTypes, protobuf.ComponentTypeFrom(pbType))
	}

	return extractComponentsByTypes(r, compTypes), nil
}

// extractComponentsFromComponentTarget extracts components from a component target.
func (rs *RLAServerImpl) extractComponentsFromComponentTarget(
	ctx context.Context,
	ct *pb.ComponentTarget,
) ([]*component.Component, error) {
	switch id := ct.GetIdentifier().(type) {
	case *pb.ComponentTarget_Uuid:
		compUUID, parseErr := uuid.Parse(id.Uuid)
		if parseErr != nil {
			return nil, fmt.Errorf("invalid component uuid %q: %w", id.Uuid, parseErr)
		}
		comp, err := rs.inventoryManager.GetComponentByID(ctx, compUUID)
		if err != nil {
			return nil, fmt.Errorf("failed to get component by uuid %s: %w", id.Uuid, err)
		}
		return []*component.Component{comp}, nil

	case *pb.ComponentTarget_External:
		extRef := id.External
		compType := protobuf.ComponentTypeFrom(extRef.GetType())
		externalID := extRef.GetId()

		// Use GetComponentsByExternalIDs which looks up by external_id field
		comps, err := rs.inventoryManager.GetComponentsByExternalIDs(ctx, []string{externalID})
		if err != nil {
			return nil, fmt.Errorf("failed to get component by external id %s: %w", externalID, err)
		}
		if len(comps) == 0 {
			return nil, fmt.Errorf("component with external id %s not found", externalID)
		}
		// Filter by type if specified
		if compType != devicetypes.ComponentTypeUnknown {
			for _, comp := range comps {
				if comp.Type == compType {
					return []*component.Component{comp}, nil
				}
			}
			return nil, fmt.Errorf("component with external id %s and type %s not found",
				externalID, devicetypes.ComponentTypeToString(compType))
		}
		return comps[:1], nil // Return first match if no type filter

	default:
		return nil, errors.New("component target must have either uuid or external set")
	}
}
