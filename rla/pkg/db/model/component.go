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
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"

	"github.com/nvidia/carbide-rest/rla/internal/carbideapi"
	"github.com/nvidia/carbide-rest/rla/pkg/common/deviceinfo"
	"github.com/nvidia/carbide-rest/rla/pkg/common/devicetypes"
	"github.com/nvidia/carbide-rest/rla/pkg/common/utils"
)

type Component struct {
	bun.BaseModel `bun:"table:component,alias:c"`

	ID              uuid.UUID              `bun:"id,pk,type:uuid,default:gen_random_uuid()"`
	Name            string                 `bun:"name"`
	Type            string                 `bun:"type,type:varchar(16),default:'Compute'"`
	Manufacturer    string                 `bun:"manufacturer,notnull,unique:component_manufacturer_serial_idx"`
	Model           string                 `bun:"model"`
	SerialNumber    string                 `bun:"serial_number,notnull,notnull,unique:component_manufacturer_serial_idx"`
	Description     map[string]any         `bun:"description,type:jsonb,json_use_number"`
	FirmwareVersion string                 `bun:"firmware_version,nullzero"`
	RackID          uuid.UUID              `bun:"rack_id,type:uuid,notnull"`
	SlotID          int                    `bun:"slot_id"`
	TrayIndex       int                    `bun:"tray_index"`
	HostID          int                    `bun:"host_id"`
	IngestedAt      *time.Time             `bun:"ingested_at"`
	DeletedAt       *time.Time             `bun:"deleted_at,soft_delete"`
	Rack            *Rack                  `bun:"rel:belongs-to,join:rack_id=id"`
	BMCs            []BMC                  `bun:"rel:has-many,join:id=component_id"`
	ComponentID     *string                `bun:"external_id"`
	PowerState      *carbideapi.PowerState `bun:"power_state"`
}

func (cd *Component) Create(ctx context.Context, idb bun.IDB) error {
	_, err := idb.NewInsert().Model(cd).Exec(ctx)
	return err
}

func (cd *Component) Get(
	ctx context.Context,
	idb bun.IDB,
) (*Component, error) {
	var component Component
	var query *bun.SelectQuery

	if cd.ID != uuid.Nil {
		query = idb.NewSelect().Model(&component).Where("id = ?", cd.ID)
	} else {
		query = idb.NewSelect().Model(&component).Where(
			"manufacturer = ? AND serial_number = ?",
			cd.Manufacturer,
			cd.SerialNumber,
		)
	}

	query = query.Relation("BMCs")

	if err := query.Scan(ctx); err != nil {
		return nil, err
	}

	return &component, nil
}

func GetAllComponents(ctx context.Context, idb bun.IDB) (ret []Component, err error) {
	err = idb.NewSelect().Model(&Component{}).Scan(ctx, &ret)
	return ret, err
}

// GetComponentsByType returns all components of a specific type with their associated BMCs
func GetComponentsByType(ctx context.Context, idb bun.IDB, componentType devicetypes.ComponentType) (ret []Component, err error) {
	err = idb.NewSelect().Model(&ret).Where("type = ?", componentType).Relation("BMCs").Scan(ctx)
	return ret, err
}

func (cd *Component) Patch(ctx context.Context, idb bun.IDB) error {
	_, err := idb.NewUpdate().Model(cd).Where("id = ?", cd.ID).Exec(ctx)
	return err
}

// BuildPatch builds a patched component from the current component
// and the input component. It goes through the patchable fields and builds
// the patched component. If there is no change on patchable fields, it returns
// nil.
func (cd *Component) BuildPatch(cur *Component) *Component {
	if cd == nil || cur == nil {
		return nil
	}

	// Make a copy fo the current component which serves as the base for the
	// patched component.
	patchedComp := *cur
	patched := false

	// Go through the patchable fields which include:
	// Description
	// FirmwareVersion
	// RackID
	// SlotID
	// TrayIndex
	// HostID

	if len(cd.FirmwareVersion) > 0 &&
		patchedComp.FirmwareVersion != cd.FirmwareVersion {
		patchedComp.FirmwareVersion = cd.FirmwareVersion
		patched = true
	}

	if desc := utils.CompareAndCopyMaps(cd.Description, cur.Description); desc != nil {
		patchedComp.Description = desc
		patched = true
	}

	if cd.RackID != uuid.Nil && cd.RackID != cur.RackID {
		patchedComp.RackID = cd.RackID
		patched = true
	}

	if cd.SlotID >= 0 && cd.SlotID != cur.SlotID {
		patchedComp.SlotID = cd.SlotID
		patched = true
	}

	if cd.TrayIndex >= 0 && cd.TrayIndex != cur.TrayIndex {
		patchedComp.TrayIndex = cd.TrayIndex
		patched = true
	}

	if cd.HostID >= 0 && cd.HostID != cur.HostID {
		patchedComp.HostID = cd.HostID
		patched = true
	}

	if !patched {
		return nil
	}

	return &patchedComp
}

// SerialInfo returns the serial number information of the component.
func (cd *Component) SerialInfo() deviceinfo.SerialInfo {
	return deviceinfo.SerialInfo{
		Manufacturer: cd.Manufacturer,
		SerialNumber: cd.SerialNumber,
	}
}

// InvalidType returns true if the component type is unknown.
func (cd *Component) InvalidType() bool {
	return !devicetypes.IsValidComponentTypeString(cd.Type)
}

func (cd *Component) SetComponentIDBySerial(ctx context.Context, idb bun.IDB) error {
	if cd.ComponentID == nil {
		return errors.New("component ID not set")
	}
	_, err := idb.NewUpdate().Model(cd).Set("external_id = ?", *cd.ComponentID).Where("serial_number = ?", cd.SerialNumber).Exec(ctx)
	return err
}

func (cd *Component) SetPowerStateByComponentID(ctx context.Context, idb bun.IDB) error {
	if cd.ComponentID == nil {
		return errors.New("component ID not set")
	}
	if cd.PowerState == nil {
		return errors.New("power state not set")
	}
	_, err := idb.NewUpdate().Model(cd).Set("power_state = ?", *cd.PowerState).Where("external_id = ?", *cd.ComponentID).Exec(ctx)
	return err
}
