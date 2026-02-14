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
package builder

import (
	"fmt"

	"github.com/google/uuid"

	"github.com/nvidia/bare-metal-manager-rest/rla/internal/inventory/objects/bmc"
	"github.com/nvidia/bare-metal-manager-rest/rla/internal/inventory/objects/component"
	"github.com/nvidia/bare-metal-manager-rest/rla/internal/inventory/objects/rack"
	"github.com/nvidia/bare-metal-manager-rest/rla/pkg/common/deviceinfo"
	"github.com/nvidia/bare-metal-manager-rest/rla/pkg/common/devicetypes"
)

type rackBuilder struct {
	rackExcel   *rack.Rack
	rackCerebro *rack.Rack

	compCerebroBySerial map[string]*component.Component
	compCerebroBySlotID map[int]*component.Component
	compCerebroByMAC    map[string]*component.Component

	reports []buildReport
}

type buildReport struct {
	description string
	compExcel   *component.Component
	compCerebro *component.Component
}

func (r *buildReport) String() string {
	return fmt.Sprintf(
		"description: %s\ncompExcel: %+v\ncompCerebro: %+v",
		r.description,
		r.compExcel,
		r.compCerebro,
	)
}

func newRackBuilder(
	rackExcel *rack.Rack,
	rackCerebro *rack.Rack,
) *rackBuilder {
	compCerebroBySerial := make(map[string]*component.Component)
	compCerebroBySlotID := make(map[int]*component.Component)
	compCerebroByMAC := make(map[string]*component.Component)
	if rackCerebro != nil {
		for _, c := range rackCerebro.Components {
			compCerebroBySerial[c.Info.SerialNumber] = &c
			compCerebroBySlotID[c.Position.SlotID] = &c
			if hostBmcs := c.BmcsByType[devicetypes.BMCTypeHost]; len(hostBmcs) > 0 {
				compCerebroByMAC[hostBmcs[0].MAC.String()] = &c
			}
		}
	}

	return &rackBuilder{
		rackExcel:           rackExcel,
		rackCerebro:         rackCerebro,
		compCerebroBySerial: compCerebroBySerial,
		compCerebroBySlotID: compCerebroBySlotID,
		compCerebroByMAC:    compCerebroByMAC,
	}
}

// build builds the rack from the Excel and Cerebro sources.
// There are a couple of points to note:
//  1. The Cerebro source does not provide information for all the components
//     in a rack, so we use the Excel source as the source of truth for the
//     rack components if the information is missing from the Cerebro source;
//     othewise, we use the Cerebro source as the source of truth.
//  2. The serial numbers of the components may be mismatched between the Excel
//     and Cerebro sources, so we use the slot ID to match the component and
//     use the information from the Cerebro as the "correct" one.
func (b *rackBuilder) build() *rack.Rack {
	getCompCerebro := func(ec *component.Component) *component.Component {
		// Try to find the matching componet from Cereburo by the following
		// order:
		// 1. By the Slot ID
		// 2. By the host BMC MAC address
		// 3. By the serial number
		// Normally, serial number is the most reliable way for matching.
		// Unfortunately, there are mistakes to input the serial number
		// by human.
		if cc := b.compCerebroBySlotID[ec.Position.SlotID]; cc != nil {
			return cc
		}

		if hostBmcs := ec.BmcsByType[devicetypes.BMCTypeHost]; len(hostBmcs) > 0 {
			if cc := b.compCerebroByMAC[hostBmcs[0].MAC.String()]; cc != nil {
				return cc
			}
		}

		if cc := b.compCerebroBySerial[ec.Info.SerialNumber]; cc != nil {
			return cc
		}

		return nil
	}

	components := make([]component.Component, 0, len(b.rackExcel.Components))
	for _, compExcel := range b.rackExcel.Components {
		compCerebro := getCompCerebro(&compExcel)
		components = append(components, b.buildComponent(&compExcel, compCerebro))
	}

	return &rack.Rack{
		Info: deviceinfo.DeviceInfo{
			ID:           b.rackCerebro.Info.ID,
			Name:         b.rackCerebro.Info.Name,
			Model:        b.rackExcel.Info.Model,
			Manufacturer: b.rackExcel.Info.Manufacturer,
			SerialNumber: b.rackExcel.Info.SerialNumber,
			Description:  b.rackExcel.Info.Description,
		},
		Loc:        b.rackExcel.Loc,
		Components: components,
	}
}

func (b *rackBuilder) buildComponent(
	ce *component.Component,
	cc *component.Component,
) component.Component {
	if ce == nil {
		panic("component from Excel should not be nil")
	}

	template := buildReport{
		compExcel:   ce,
		compCerebro: cc,
	}

	if cc == nil {
		if cerebroSupportedComponentType(ce.Type) {
			b.newBuildReport(template, "component from Cerebro is missing")
		}

		// For the components which are not found in Cerebro, use the
		// component from Excel.
		if ce.Info.ID == uuid.Nil {
			// Create a new ID for the component.
			ce.Info.ID = uuid.New()
		}

		return *ce
	}

	// Create the base component by cherry-picking information from the
	// Cerebro and Excel
	postion := cc.Position
	if cc.Position.SlotID == 0 {
		// The Cerebro does not provide the slot ID, so we use the slot ID
		// from Excel.
		postion = ce.Position
	}

	comp := component.New(
		cc.Type,
		&deviceinfo.DeviceInfo{
			ID:           cc.Info.ID,
			Name:         cc.Info.Name,
			SerialNumber: cc.Info.SerialNumber,
			Manufacturer: cc.Info.Manufacturer,
			Model:        cc.Info.Model,
			Description:  cc.Info.Description,
		},
		ce.FirmwareVersion,
		&postion,
	)

	if comp.Info.ID == uuid.Nil {
		b.newBuildReport(template, "component from Cerebro has no ID")

		// Create a new ID for the component.
		comp.Info.ID = uuid.New()
	}

	if cc.Type != ce.Type {
		b.newBuildReport(
			template,
			"component from Cerebro does not match the type from Excel",
		)

		// Use the component type from Cerebro, no further action.
	}

	if cc.Info.SerialNumber != ce.Info.SerialNumber {
		b.newBuildReport(
			template,
			"component from Cerebro does not match the serial number from Excel",
		)

		if matchBmcs(ce.BmcsByType[devicetypes.BMCTypeHost], cc.BmcsByType[devicetypes.BMCTypeHost]) {
			// The host BMCs match, so it is not a replacement. Most likely it
			// is caused by homan mistakes, so use the excel serial number.
			comp.Info.SerialNumber = ce.Info.SerialNumber
		}

		// Use the component serial from Cerebro, no further action.
	}

	if !matchBmcs(ce.BmcsByType[devicetypes.BMCTypeHost], cc.BmcsByType[devicetypes.BMCTypeHost]) {
		b.newBuildReport(
			template,
			"component from Cerebro does not match the Host BMCs from Excel",
		)
	} else {
		if !matchBmcs(ce.BmcsByType[devicetypes.BMCTypeDPU], cc.BmcsByType[devicetypes.BMCTypeDPU]) {
			b.newBuildReport(
				template,
				"component from Cerebro does not match the DPU BMCs from Excel",
			)
		}
	}

	// Use the BMCs from Cerebro if available, otherwise, use the BMCs
	// from Excel.
	bmcsByType := cc.BmcsByType
	if len(cc.BmcsByType[devicetypes.BMCTypeHost]) == 0 {
		bmcsByType = ce.BmcsByType
	}

	for typ, bmcs := range bmcsByType {
		for _, bmc := range bmcs {
			comp.AddBMC(typ, bmc)
		}
	}

	return comp
}

func (b *rackBuilder) newBuildReport(t buildReport, d string) {
	t.description = d
	b.reports = append(b.reports, t)
}

func matchBmcs(bmcs1 []bmc.BMC, bmcs2 []bmc.BMC) bool {
	if len(bmcs1) != len(bmcs2) {
		return false
	}

	bmcs2Map := make(map[string]bmc.BMC)
	for _, bmc2 := range bmcs2 {
		bmcs2Map[bmc2.MAC.String()] = bmc2
	}

	for _, bmc1 := range bmcs1 {
		if _, ok := bmcs2Map[bmc1.MAC.String()]; !ok {
			return false
		}
	}

	return true
}

func cerebroSupportedComponentType(typ devicetypes.ComponentType) bool {
	return typ == devicetypes.ComponentTypeCompute ||
		typ == devicetypes.ComponentTypeNVLSwitch ||
		typ == devicetypes.ComponentTypeToRSwitch ||
		typ == devicetypes.ComponentTypePowerShelf
}
