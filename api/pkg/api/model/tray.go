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

package model

import (
	"fmt"
	"maps"
	"net/url"
	"slices"

	"github.com/google/uuid"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	validationis "github.com/go-ozzo/ozzo-validation/v4/is"

	rlav1 "github.com/nvidia/bare-metal-manager-rest/workflow-schema/rla/protobuf/v1"
)

// APIToProtoComponentTypeName maps API tray type strings to protobuf ComponentType enum names.
var APIToProtoComponentTypeName = map[string]string{
	"compute":    "COMPONENT_TYPE_COMPUTE",
	"switch":     "COMPONENT_TYPE_NVLSWITCH",
	"powershelf": "COMPONENT_TYPE_POWERSHELF",
}

// ProtoToAPIComponentTypeName maps protobuf ComponentType enum names to API tray type strings.
var ProtoToAPIComponentTypeName = map[string]string{
	"COMPONENT_TYPE_COMPUTE":    "compute",
	"COMPONENT_TYPE_NVLSWITCH":  "switch",
	"COMPONENT_TYPE_POWERSHELF": "powershelf",
}

var validTrayTypesAny, ValidProtoComponentTypes = func() ([]interface{}, []rlav1.ComponentType) {
	anyTypes := make([]interface{}, 0, len(APIToProtoComponentTypeName))
	protoTypes := make([]rlav1.ComponentType, 0, len(APIToProtoComponentTypeName))
	for apiName, protoName := range APIToProtoComponentTypeName {
		anyTypes = append(anyTypes, apiName)
		protoTypes = append(protoTypes, rlav1.ComponentType(rlav1.ComponentType_value[protoName]))
	}
	return anyTypes, protoTypes
}()

// TrayOrderByFieldMap maps API field names to RLA protobuf ComponentOrderByField enum
var TrayOrderByFieldMap = map[string]rlav1.ComponentOrderByField{
	"name":         rlav1.ComponentOrderByField_COMPONENT_ORDER_BY_FIELD_NAME,
	"manufacturer": rlav1.ComponentOrderByField_COMPONENT_ORDER_BY_FIELD_MANUFACTURER,
	"model":        rlav1.ComponentOrderByField_COMPONENT_ORDER_BY_FIELD_MODEL,
	"type":         rlav1.ComponentOrderByField_COMPONENT_ORDER_BY_FIELD_TYPE,
}

// GetProtoTrayOrderByFromQueryParam creates an RLA protobuf OrderBy from API query parameters for tray (component) queries
func GetProtoTrayOrderByFromQueryParam(fieldName, direction string) *rlav1.OrderBy {
	field, ok := TrayOrderByFieldMap[fieldName]
	if !ok {
		return nil
	}
	return &rlav1.OrderBy{
		Field: &rlav1.OrderBy_ComponentField{
			ComponentField: field,
		},
		Direction: direction,
	}
}

// APITrayGetAllRequest captures query parameters for listing trays from RLA.
type APITrayGetAllRequest struct {
	RackID       *string
	RackName     *string
	Type         *string
	ComponentIDs []string
	IDs          []string
}

// FromQueryParams populates the request from URL query parameters.
func (r *APITrayGetAllRequest) FromQueryParams(params url.Values) {
	if v := params.Get("rackId"); v != "" {
		r.RackID = &v
	}
	if v := params.Get("rackName"); v != "" {
		r.RackName = &v
	}
	if v := params.Get("type"); v != "" {
		r.Type = &v
	}
	if vals := params["componentId"]; len(vals) > 0 {
		r.ComponentIDs = vals
	}
	if vals := params["id"]; len(vals) > 0 {
		r.IDs = vals
	}
}

// Validate checks field formats and enforces the RLA protobuf oneof constraints:
//   - rackId must be a valid UUID
//   - rackId and rackName are mutually exclusive (RackTarget.oneof identifier)
//   - rackId/rackName cannot be combined with id/componentId (OperationTargetSpec.oneof targets)
//   - componentId requires type (ExternalRef needs type)
//   - type must be one of the supported tray types
//   - each entry in IDs must be a valid UUID
func (r *APITrayGetAllRequest) Validate() error {
	err := validation.ValidateStruct(r,
		validation.Field(&r.RackID,
			validation.When(r.RackID != nil, validationis.UUID.Error(validationErrorInvalidUUID))),
		validation.Field(&r.Type,
			validation.When(r.Type != nil, validation.In(validTrayTypesAny...).Error(
				fmt.Sprintf("must be one of %v", slices.Collect(maps.Keys(APIToProtoComponentTypeName)))))),
	)
	if err != nil {
		return err
	}

	for _, id := range r.IDs {
		if _, parseErr := uuid.Parse(id); parseErr != nil {
			return validation.Errors{"id": fmt.Errorf("%s: %s", validationErrorInvalidUUID, id)}
		}
	}

	if r.RackID != nil && r.RackName != nil {
		return validation.Errors{"rackId": fmt.Errorf("rackId and rackName are mutually exclusive")}
	}

	hasRackParams := r.RackID != nil || r.RackName != nil
	hasComponentParams := len(r.IDs) > 0 || len(r.ComponentIDs) > 0
	if hasRackParams && hasComponentParams {
		return validation.Errors{"rackId": fmt.Errorf("rackId/rackName cannot be combined with id/componentId")}
	}

	if len(r.ComponentIDs) > 0 && r.Type == nil {
		return validation.Errors{"componentId": fmt.Errorf("type is required when componentId is provided")}
	}

	return nil
}

// ToProto converts a validated APITrayGetAllRequest to an RLA GetComponentsRequest.
func (r *APITrayGetAllRequest) ToProto() *rlav1.GetComponentsRequest {
	rlaRequest := &rlav1.GetComponentsRequest{}

	hasIDs := len(r.IDs) > 0
	hasComponentIDsWithType := len(r.ComponentIDs) > 0 && r.Type != nil

	if hasIDs || hasComponentIDsWithType {
		componentTargets := make([]*rlav1.ComponentTarget, 0, len(r.IDs)+len(r.ComponentIDs))

		for _, id := range r.IDs {
			componentTargets = append(componentTargets, &rlav1.ComponentTarget{
				Identifier: &rlav1.ComponentTarget_Id{
					Id: &rlav1.UUID{Id: id},
				},
			})
		}

		if hasComponentIDsWithType {
			if protoName, ok := APIToProtoComponentTypeName[*r.Type]; ok {
				protoType := rlav1.ComponentType(rlav1.ComponentType_value[protoName])
				for _, cid := range r.ComponentIDs {
					componentTargets = append(componentTargets, &rlav1.ComponentTarget{
						Identifier: &rlav1.ComponentTarget_External{
							External: &rlav1.ExternalRef{
								Type: protoType,
								Id:   cid,
							},
						},
					})
				}
			}
		}

		rlaRequest.TargetSpec = &rlav1.OperationTargetSpec{
			Targets: &rlav1.OperationTargetSpec_Components{
				Components: &rlav1.ComponentTargets{
					Targets: componentTargets,
				},
			},
		}
		return rlaRequest
	}

	rackTarget := &rlav1.RackTarget{}

	if r.RackID != nil {
		rackTarget.Identifier = &rlav1.RackTarget_Id{
			Id: &rlav1.UUID{Id: *r.RackID},
		}
	} else if r.RackName != nil {
		rackTarget.Identifier = &rlav1.RackTarget_Name{
			Name: *r.RackName,
		}
	}

	if r.Type != nil {
		if protoName, ok := APIToProtoComponentTypeName[*r.Type]; ok {
			rackTarget.ComponentTypes = []rlav1.ComponentType{
				rlav1.ComponentType(rlav1.ComponentType_value[protoName]),
			}
		}
	} else {
		rackTarget.ComponentTypes = ValidProtoComponentTypes
	}

	rlaRequest.TargetSpec = &rlav1.OperationTargetSpec{
		Targets: &rlav1.OperationTargetSpec_Racks{
			Racks: &rlav1.RackTargets{
				Targets: []*rlav1.RackTarget{rackTarget},
			},
		},
	}

	return rlaRequest
}

// APITrayPosition represents the position of a tray within a rack
type APITrayPosition struct {
	SlotID  int32 `json:"slotId"`
	TrayIdx int32 `json:"trayIdx"`
	HostID  int32 `json:"hostId"`
}

// FromProto converts a proto RackPosition to an APITrayPosition
func (atp *APITrayPosition) FromProto(protoPosition *rlav1.RackPosition) {
	if protoPosition == nil {
		return
	}
	atp.SlotID = protoPosition.GetSlotId()
	atp.TrayIdx = protoPosition.GetTrayIdx()
	atp.HostID = protoPosition.GetHostId()
}

// APITray is the API representation of a Tray (Component) from RLA
type APITray struct {
	ID              string           `json:"id"`
	ComponentID     string           `json:"componentId"`
	Type            string           `json:"type"`
	Name            string           `json:"name"`
	Manufacturer    string           `json:"manufacturer"`
	Model           string           `json:"model"`
	SerialNumber    string           `json:"serialNumber"`
	Description     string           `json:"description"`
	FirmwareVersion string           `json:"firmwareVersion"`
	PowerState      string           `json:"powerState"`
	Position        *APITrayPosition `json:"position"`
	RackID          string           `json:"rackId"`
}

// FromProto converts an RLA protobuf Component to an APITray
func (at *APITray) FromProto(comp *rlav1.Component) {
	if comp == nil {
		return
	}

	at.Type = ProtoToAPIComponentTypeName[rlav1.ComponentType_name[int32(comp.GetType())]]
	at.FirmwareVersion = comp.GetFirmwareVersion()
	at.PowerState = comp.GetPowerState()
	at.ComponentID = comp.GetComponentId()

	// Get info from DeviceInfo
	if comp.GetInfo() != nil {
		info := comp.GetInfo()
		if info.GetId() != nil {
			at.ID = info.GetId().GetId()
		}
		at.Name = info.GetName()
		at.Manufacturer = info.GetManufacturer()
		if info.Model != nil {
			at.Model = *info.Model
		}
		at.SerialNumber = info.GetSerialNumber()
		if info.Description != nil {
			at.Description = *info.Description
		}
	}

	// Get position
	if comp.GetPosition() != nil {
		at.Position = &APITrayPosition{}
		at.Position.FromProto(comp.GetPosition())
	}

	// Get rack ID
	if comp.GetRackId() != nil {
		at.RackID = comp.GetRackId().GetId()
	}
}

// NewAPITray creates an APITray from the RLA protobuf Component
func NewAPITray(comp *rlav1.Component) *APITray {
	if comp == nil {
		return nil
	}
	apiTray := &APITray{}
	apiTray.FromProto(comp)
	return apiTray
}
