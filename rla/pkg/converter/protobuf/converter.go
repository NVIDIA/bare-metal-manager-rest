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
package protobuf

import (
	"net"

	"github.com/google/uuid"

	"github.com/nvidia/bare-metal-manager-rest/rla/internal/inventory/objects/bmc"
	"github.com/nvidia/bare-metal-manager-rest/rla/internal/inventory/objects/component"
	"github.com/nvidia/bare-metal-manager-rest/rla/internal/inventory/objects/nvldomain"
	"github.com/nvidia/bare-metal-manager-rest/rla/internal/inventory/objects/rack"
	pb "github.com/nvidia/bare-metal-manager-rest/rla/internal/proto/v1"
	taskcommon "github.com/nvidia/bare-metal-manager-rest/rla/internal/task/common"
	"github.com/nvidia/bare-metal-manager-rest/rla/internal/task/operations"
	taskdef "github.com/nvidia/bare-metal-manager-rest/rla/internal/task/task"
	identifier "github.com/nvidia/bare-metal-manager-rest/rla/pkg/common/Identifier"
	"github.com/nvidia/bare-metal-manager-rest/rla/pkg/common/credential"
	"github.com/nvidia/bare-metal-manager-rest/rla/pkg/common/deviceinfo"
	"github.com/nvidia/bare-metal-manager-rest/rla/pkg/common/devicetypes"
	"github.com/nvidia/bare-metal-manager-rest/rla/pkg/common/location"
	dbquery "github.com/nvidia/bare-metal-manager-rest/rla/pkg/db/query"
)

var (
	componentTypeToMap   map[devicetypes.ComponentType]pb.ComponentType
	bmcTypeToMap         map[devicetypes.BMCType]pb.BMCType
	componentTypeFromMap map[pb.ComponentType]devicetypes.ComponentType
	bmcTypeFromMap       map[pb.BMCType]devicetypes.BMCType
)

func init() {
	// Initialize the mappings between internal component types and protobuf
	// component types
	componentTypeToMap = map[devicetypes.ComponentType]pb.ComponentType{
		devicetypes.ComponentTypeUnknown:    pb.ComponentType_COMPONENT_TYPE_UNKNOWN,
		devicetypes.ComponentTypeCompute:    pb.ComponentType_COMPONENT_TYPE_COMPUTE,
		devicetypes.ComponentTypeNVLSwitch:  pb.ComponentType_COMPONENT_TYPE_NVLSWITCH,
		devicetypes.ComponentTypePowerShelf: pb.ComponentType_COMPONENT_TYPE_POWERSHELF,
		devicetypes.ComponentTypeToRSwitch:  pb.ComponentType_COMPONENT_TYPE_TORSWITCH,
		devicetypes.ComponentTypeUMS:        pb.ComponentType_COMPONENT_TYPE_UMS,
		devicetypes.ComponentTypeCDU:        pb.ComponentType_COMPONENT_TYPE_CDU,
	}

	// Initialize the mappings between internal BMC types and protobuf BMC
	// types
	bmcTypeToMap = map[devicetypes.BMCType]pb.BMCType{
		devicetypes.BMCTypeUnknown: pb.BMCType_BMC_TYPE_UNKNOWN,
		devicetypes.BMCTypeHost:    pb.BMCType_BMC_TYPE_HOST,
		devicetypes.BMCTypeDPU:     pb.BMCType_BMC_TYPE_DPU,
	}

	// Reverse mappings for component types
	componentTypeFromMap = make(map[pb.ComponentType]devicetypes.ComponentType)
	for t, pt := range componentTypeToMap {
		componentTypeFromMap[pt] = t
	}

	// Reverse mappings for BMC types
	bmcTypeFromMap = make(map[pb.BMCType]devicetypes.BMCType)
	for t, pt := range bmcTypeToMap {
		bmcTypeFromMap[pt] = t
	}
}

// ComponentTypeFrom converts a protobuf ComponentType to an internal
// ComponentType
func ComponentTypeFrom(pt pb.ComponentType) devicetypes.ComponentType {
	if t, ok := componentTypeFromMap[pt]; ok {
		return t
	}

	return devicetypes.ComponentTypeUnknown
}

// BMCTypeFrom converts a protobuf BMCType to an internal BMCType
func BMCTypeFrom(pt pb.BMCType) devicetypes.BMCType {
	if t, ok := bmcTypeFromMap[pt]; ok {
		return t
	}

	return devicetypes.BMCTypeUnknown
}

// DeviceInfoFrom converts a protobuf DeviceInfo to an internal DeviceInfo
func DeviceInfoFrom(info *pb.DeviceInfo) deviceinfo.DeviceInfo {
	if info == nil {
		return deviceinfo.DeviceInfo{}
	}

	// No need to check whether the info is nil. They are handled by the
	// methods of pb.DeviceInfo.
	return deviceinfo.DeviceInfo{
		ID:           UUIDFrom(info.GetId()),
		Name:         info.GetName(),
		Manufacturer: info.GetManufacturer(),
		Model:        info.GetModel(),
		SerialNumber: info.GetSerialNumber(),
		Description:  info.GetDescription(),
	}
}

// LocationFrom converts a protobuf Location to an internal Location
func LocationFrom(loc *pb.Location) location.Location {
	if loc == nil {
		return location.Location{}
	}

	return location.Location{
		Region:     loc.GetRegion(),
		DataCenter: loc.GetDatacenter(),
		Room:       loc.GetRoom(),
		Position:   loc.GetPosition(),
	}
}

// UUIDFrom converts a protobuf UUID to an internal uuid.UUID
func UUIDFrom(id *pb.UUID) uuid.UUID {
	if id != nil {
		if parsed, err := uuid.Parse(id.Id); err == nil {
			return parsed
		}
	}

	return uuid.Nil
}

// UUIDsFrom converts a slice of *pb.UUID to a slice of uuid.UUID.
func UUIDsFrom(ids []*pb.UUID) []uuid.UUID {
	result := make([]uuid.UUID, 0, len(ids))
	for _, id := range ids {
		if parsed := UUIDFrom(id); parsed != uuid.Nil {
			result = append(result, parsed)
		}
	}
	return result
}

// RackPositionFrom converts a protobuf RackPosition to an internal
// InRackPosition
func RackPositionFrom(pos *pb.RackPosition) component.InRackPosition {
	if pos == nil {
		// Return invalid position if the postion is not specified.
		return component.InRackPosition{
			SlotID:    -1,
			TrayIndex: -1,
			HostID:    -1,
		}
	}

	return component.InRackPosition{
		SlotID:    int(pos.SlotId),
		TrayIndex: int(pos.TrayIdx),
		HostID:    int(pos.HostId),
	}
}

// BMCFrom converts a protobuf BMCInfo to an internal BMCType and BMC
func BMCFrom(bi *pb.BMCInfo) (devicetypes.BMCType, *bmc.BMC) {
	if bi == nil {
		return devicetypes.BMCTypeUnknown, nil
	}

	var bmc bmc.BMC

	if hwAddr, err := net.ParseMAC(bi.GetMacAddress()); err == nil {
		bmc.MAC = hwAddr
	}

	if ip := bi.GetIpAddress(); ip != "" {
		bmc.IP = net.ParseIP(ip)
	}

	if bi.User != nil && bi.Password != nil {
		// Create a new credential only if the user and password are not nil.
		nc := credential.New(*bi.User, *bi.Password)
		bmc.Credential = &nc
	}

	return BMCTypeFrom(bi.Type), &bmc
}

// ComponentFrom converts a protobuf Component to an internal Component
func ComponentFrom(c *pb.Component) *component.Component {
	if c == nil {
		return nil
	}

	bmcsByType := make(map[devicetypes.BMCType][]bmc.BMC)
	for _, bi := range c.GetBmcs() {
		t, bmc := BMCFrom(bi)
		if bmc != nil {
			bmcsByType[t] = append(bmcsByType[t], *bmc)
		}
	}

	return &component.Component{
		Type:            ComponentTypeFrom(c.GetType()),
		Info:            DeviceInfoFrom(c.GetInfo()),
		FirmwareVersion: c.GetFirmwareVersion(),
		Position:        RackPositionFrom(c.GetPosition()),
		BmcsByType:      bmcsByType,
	}
}

// RackFrom converts a protobuf Rack to an internal Rack
func RackFrom(r *pb.Rack) *rack.Rack {
	if r == nil {
		return nil
	}

	components := make([]component.Component, 0, len(r.GetComponents()))
	for _, c := range r.GetComponents() {
		components = append(components, *ComponentFrom(c))
	}

	return &rack.Rack{
		Info:       DeviceInfoFrom(r.GetInfo()),
		Loc:        LocationFrom(r.GetLocation()),
		Components: components,
	}
}

func PaginationFrom(pg *pb.Pagination) *dbquery.Pagination {
	if pg == nil {
		return nil
	}

	return &dbquery.Pagination{
		Offset: int(pg.GetOffset()),
		Limit:  int(pg.GetLimit()),
	}
}

func StringQueryInfoFrom(info *pb.StringQueryInfo) *dbquery.StringQueryInfo {
	if info == nil {
		return nil
	}

	return &dbquery.StringQueryInfo{
		Patterns:   info.GetPatterns(),
		IsWildcard: info.GetIsWildcard(),
		UseOR:      info.GetUseOr(),
	}
}

func IdentifierFrom(info *pb.Identifier) *identifier.Identifier {
	if info == nil {
		return nil
	}

	return identifier.New(UUIDFrom(info.GetId()), info.GetName())
}

func NVLDomainFrom(info *pb.NVLDomain) *nvldomain.NVLDomain {
	if info == nil || info.GetIdentifier() == nil {
		return nil
	}

	return &nvldomain.NVLDomain{
		Identifier: *IdentifierFrom(info.GetIdentifier()),
	}
}

func PowerControlOpFrom(op pb.PowerControlOp) operations.PowerOperation {
	switch op {
	// Power On
	case pb.PowerControlOp_POWER_CONTROL_OP_ON:
		return operations.PowerOperationPowerOn
	case pb.PowerControlOp_POWER_CONTROL_OP_FORCE_ON:
		return operations.PowerOperationForcePowerOn
	// Power Off
	case pb.PowerControlOp_POWER_CONTROL_OP_OFF:
		return operations.PowerOperationPowerOff
	case pb.PowerControlOp_POWER_CONTROL_OP_FORCE_OFF:
		return operations.PowerOperationForcePowerOff
	// Restart (OS level)
	case pb.PowerControlOp_POWER_CONTROL_OP_RESTART:
		return operations.PowerOperationRestart
	case pb.PowerControlOp_POWER_CONTROL_OP_FORCE_RESTART:
		return operations.PowerOperationForceRestart
	// Reset (hardware level)
	case pb.PowerControlOp_POWER_CONTROL_OP_WARM_RESET:
		return operations.PowerOperationWarmReset
	case pb.PowerControlOp_POWER_CONTROL_OP_COLD_RESET:
		return operations.PowerOperationColdReset
	default:
		return operations.PowerOperationUnknown
	}
}

func TaskExecutorTypeFrom(et pb.TaskExecutorType) taskcommon.ExecutorType {
	switch et {
	case pb.TaskExecutorType_TASK_EXECUTOR_TYPE_TEMPORAL:
		return taskcommon.ExecutorTypeTemporal
	default:
		return taskcommon.ExecutorTypeUnknown
	}
}

func TaskExecutorTypeTo(et taskcommon.ExecutorType) pb.TaskExecutorType {
	switch et {
	case taskcommon.ExecutorTypeTemporal:
		return pb.TaskExecutorType_TASK_EXECUTOR_TYPE_TEMPORAL
	default:
		return pb.TaskExecutorType_TASK_EXECUTOR_TYPE_UNKNOWN
	}
}

func TaskStatusFrom(status pb.TaskStatus) taskcommon.TaskStatus {
	switch status {
	case pb.TaskStatus_TASK_STATUS_PENDING:
		return taskcommon.TaskStatusPending
	case pb.TaskStatus_TASK_STATUS_RUNNING:
		return taskcommon.TaskStatusRunning
	case pb.TaskStatus_TASK_STATUS_COMPLETED:
		return taskcommon.TaskStatusCompleted
	case pb.TaskStatus_TASK_STATUS_FAILED:
		return taskcommon.TaskStatusFailed
	default:
		return taskcommon.TaskStatusUnknown
	}
}

func TaskStatusTo(status taskcommon.TaskStatus) pb.TaskStatus {
	switch status {
	case taskcommon.TaskStatusPending:
		return pb.TaskStatus_TASK_STATUS_PENDING
	case taskcommon.TaskStatusRunning:
		return pb.TaskStatus_TASK_STATUS_RUNNING
	case taskcommon.TaskStatusCompleted:
		return pb.TaskStatus_TASK_STATUS_COMPLETED
	case taskcommon.TaskStatusFailed:
		return pb.TaskStatus_TASK_STATUS_FAILED
	default:
		return pb.TaskStatus_TASK_STATUS_UNKNOWN
	}
}

func TaskTo(task *taskdef.Task) *pb.Task {
	if task == nil {
		return nil
	}

	componentIDs := make([]*pb.UUID, 0, len(task.ComponentUUIDs))
	for _, id := range task.ComponentUUIDs {
		componentIDs = append(componentIDs, UUIDTo(id))
	}

	var opStr string
	operation, err := operations.New(task.Operation.Type, task.Operation.Info)
	if err == nil {
		opStr = operation.Description()
	}

	return &pb.Task{
		Id:             UUIDTo(task.ID),
		Operation:      opStr,
		RackId:         UUIDTo(task.RackID),
		ComponentUuids: componentIDs,
		Description:    task.Description,
		ExecutorType:   TaskExecutorTypeTo(task.ExecutorType),
		ExecutionId:    task.ExecutionID,
		Status:         TaskStatusTo(task.Status),
		Message:        task.Message,
	}
}

// ComponentTypeTo converts an internal ComponentType to a protobuf
// ComponentType
func ComponentTypeTo(t devicetypes.ComponentType) pb.ComponentType {
	if pt, ok := componentTypeToMap[t]; ok {
		return pt
	}

	return pb.ComponentType_COMPONENT_TYPE_UNKNOWN
}

// BMCTypeTo converts an internal BMCType to a protobuf BMCType
func BMCTypeTo(t devicetypes.BMCType) pb.BMCType {
	if pt, ok := bmcTypeToMap[t]; ok {
		return pt
	}

	return pb.BMCType_BMC_TYPE_UNKNOWN
}

// DeviceInfoTo converts an internal DeviceInfo to a protobuf DeviceInfo
func DeviceInfoTo(info *deviceinfo.DeviceInfo) *pb.DeviceInfo {
	if info == nil {
		return nil
	}

	pinfo := pb.DeviceInfo{
		Id:           UUIDTo(info.ID),
		Name:         info.Name,
		Manufacturer: info.Manufacturer,
		SerialNumber: info.SerialNumber,
	}

	if model := info.Model; model != "" {
		pinfo.Model = &model
	}

	if description := info.Description; description != "" {
		pinfo.Description = &description
	}

	return &pinfo
}

// LocationTo converts an internal Location to a protobuf Location
func LocationTo(loc *location.Location) *pb.Location {
	if loc == nil {
		return nil
	}

	return &pb.Location{
		Region:     loc.Region,
		Datacenter: loc.DataCenter,
		Room:       loc.Room,
		Position:   loc.Position,
	}
}

// UUIDTo converts an internal uuid.UUID to a protobuf UUID
func UUIDTo(id uuid.UUID) *pb.UUID {
	if id != uuid.Nil {
		return &pb.UUID{Id: id.String()}
	}

	return nil
}

// UUIDsTo converts a slice of uuid.UUID to a slice of *pb.UUID.
func UUIDsTo(ids []uuid.UUID) []*pb.UUID {
	result := make([]*pb.UUID, 0, len(ids))
	for _, id := range ids {
		if id != uuid.Nil {
			result = append(result, &pb.UUID{Id: id.String()})
		}
	}
	return result
}

// RackPositionTo converts an internal InRackPosition to a protobuf
// RackPosition
func RackPositionTo(pos *component.InRackPosition) *pb.RackPosition {
	if pos == nil {
		return nil
	}

	return &pb.RackPosition{
		SlotId:  int32(pos.SlotID),
		TrayIdx: int32(pos.TrayIndex),
		HostId:  int32(pos.HostID),
	}
}

// BMCTo converts an internal BMCType and BMC to a protobuf BMCInfo
func BMCTo(t devicetypes.BMCType, b *bmc.BMC) *pb.BMCInfo {
	if b == nil {
		return nil
	}

	bmcProto := pb.BMCInfo{
		Type:       BMCTypeTo(t),
		MacAddress: b.MAC.String(),
	}

	if b.IP != nil {
		ip := b.IP.String()
		bmcProto.IpAddress = &ip
	}

	if b.Credential != nil {
		bmcProto.User, bmcProto.Password = b.Credential.Retrieve()
	}

	return &bmcProto
}

// ComponentTo converts an internal Component to a protobuf Component
func ComponentTo(c *component.Component) *pb.Component {
	if c == nil {
		return nil
	}

	bmcInfos := make([]*pb.BMCInfo, 0)
	for t, bmcs := range c.BmcsByType {
		for _, bmc := range bmcs {
			bmcInfos = append(bmcInfos, BMCTo(t, &bmc))
		}
	}

	return &pb.Component{
		Type:            ComponentTypeTo(c.Type),
		Info:            DeviceInfoTo(&c.Info),
		FirmwareVersion: c.FirmwareVersion,
		Position:        RackPositionTo(&c.Position),
		Bmcs:            bmcInfos,
		ComponentId:     c.ComponentID,
		RackId:          UUIDTo(c.RackID),
	}
}

// RackTo converts an internal Rack to a protobuf Rack
func RackTo(r *rack.Rack) *pb.Rack {
	if r == nil {
		return nil
	}

	components := make([]*pb.Component, 0, len(r.Components))
	for _, c := range r.Components {
		components = append(components, ComponentTo(&c))
	}

	return &pb.Rack{
		Info:       DeviceInfoTo(&r.Info),
		Location:   LocationTo(&r.Loc),
		Components: components,
	}
}

func PaginationTo(pg *dbquery.Pagination) *pb.Pagination {
	if pg == nil {
		return nil
	}

	return &pb.Pagination{
		Offset: int32(pg.Offset),
		Limit:  int32(pg.Limit),
	}
}

func StringQueryInfoTo(info *dbquery.StringQueryInfo) *pb.StringQueryInfo {
	if info == nil {
		return nil
	}

	return &pb.StringQueryInfo{
		Patterns:   info.Patterns,
		IsWildcard: info.IsWildcard,
		UseOr:      info.UseOR,
	}
}

func IdentifierTo(info *identifier.Identifier) *pb.Identifier {
	if info == nil {
		return nil
	}

	return &pb.Identifier{
		Id:   UUIDTo(info.ID),
		Name: info.Name,
	}
}

func NVLDomainTo(info *nvldomain.NVLDomain) *pb.NVLDomain {
	if info == nil {
		return nil
	}

	return &pb.NVLDomain{
		Identifier: IdentifierTo(&info.Identifier),
	}
}
