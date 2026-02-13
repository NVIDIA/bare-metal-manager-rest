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
package excel

import (
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/nvidia/bare-metal-manager-rest/rla/internal/inventory/objects/bmc"
	"github.com/nvidia/bare-metal-manager-rest/rla/internal/inventory/objects/component"
	"github.com/nvidia/bare-metal-manager-rest/rla/pkg/common/deviceinfo"
	"github.com/nvidia/bare-metal-manager-rest/rla/pkg/common/devicetypes"
	"github.com/nvidia/bare-metal-manager-rest/rla/pkg/common/location"
)

type RackInvetory struct {
	Manufacturer  string `title:"Manufacturer"`
	Model         string `title:"Model Name"`
	ManufactoryPN string `title:"MANUFACTURE_PN"`
	// Description      string `title:"Description"`
	SerialNum    string `title:"System Serial Number"`
	Region       string `title:"Region Name"`
	DataCenter   string `title:"Name of DC"`
	Room         string `title:"DC Room"`
	RackPosition string `title:"Rack Position in DC"`
	RackSlot     string `title:"Rack Sub Position"`
	// RackSubSlot      string `title:"Rack Sub Position Slot"`
	RackSerialNumber string `title:"RACK_SN"`
	BMCMAC           string `title:"MAC address-1"`
	BMCMAC1          string `title:"MAC address-2"`
}

func (inv *RackInvetory) rackInfo(name string) deviceinfo.DeviceInfo {
	return deviceinfo.DeviceInfo{
		ID:           uuid.Nil,
		Name:         name,
		Manufacturer: inv.Manufacturer,
		SerialNumber: inv.RackSerialNumber,
	}
}

func (inv *RackInvetory) locationInfo() location.Location {
	return location.Location{
		Region:     inv.Region,
		DataCenter: inv.DataCenter,
		Room:       inv.Room,
		Position:   inv.RackPosition,
	}
}

func (inv *RackInvetory) componentType() devicetypes.ComponentType {
	switch inv.ManufactoryPN {
	case "B81.11801.0008":
		return devicetypes.ComponentTypeCompute
	case "BCS.118K2.0001":
		return devicetypes.ComponentTypeNVLSwitch
	case "BK1.00111.0003", "BK1.00118.0001":
		return devicetypes.ComponentTypePowerShelf
	case "BK2.02201.MS01":
		return devicetypes.ComponentTypeToRSwitch
	case "B81.10X01.0008":
		return devicetypes.ComponentTypeUMS
	}

	return devicetypes.ComponentTypeUnknown
}

func (inv *RackInvetory) componentInfo() component.Component {
	c := component.New(
		inv.componentType(),
		&deviceinfo.DeviceInfo{
			Manufacturer: inv.Manufacturer,
			Model:        inv.Model,
			SerialNumber: inv.SerialNum,
		},
		"",
		&component.InRackPosition{
			SlotID: strToID(inv.RackSlot, false),
		},
	)

	if mac, err := bmc.New(inv.BMCMAC, nil, ""); err == nil {
		c.AddBMC(devicetypes.BMCTypeHost, *mac)
	}

	if mac, err := bmc.New(inv.BMCMAC1, nil, ""); err == nil {
		c.AddBMC(devicetypes.BMCTypeHost, *mac)
	}

	return c
}

type ComponentInvetory struct {
	// Site            string `title:"Site Name"`
	ParentTag    string `title:"Parent Asset Tag # (TBD)"`
	Manufacturer string `title:"Manufacturer"`
	SerialNum    string `title:"System Serial Number"`
	// Type             string `title:"Model Name"`
	FirmwareVersion string `title:"Firmware Version"`
	BMCMAC          string `title:"BMC MAC -1"`
	DPU1BMCMAC      string `title:"DPU-1 BMC MAC"`
	DPU2BMCMAC      string `title:"DPU-2 BMC MAC"`
}

func (ci *ComponentInvetory) componentInfo() component.Component {
	c := component.New(
		devicetypes.ComponentTypeUnknown,
		&deviceinfo.DeviceInfo{
			Manufacturer: ci.Manufacturer,
			SerialNumber: ci.SerialNum,
		},
		ci.FirmwareVersion,
		nil,
	)

	// NOTE: the information from the component inventory excel file may not have
	// the correct host BMC information, so we skip it. It is available in the
	// rack inventory excel file anyway.
	/*
		if info, err := bmc.New(ici.BMCMAC, nil, ""); err == nil {
			c.AddBMC(bmc.TypeHost, *info)
		}
	*/
	if ci.DPU1BMCMAC != "" {
		if info, err := bmc.New(ci.DPU1BMCMAC, nil, ""); err == nil {
			c.AddBMC(devicetypes.BMCTypeDPU, *info)
		}
	}

	if ci.DPU2BMCMAC != "" {
		if info, err := bmc.New(ci.DPU2BMCMAC, nil, ""); err == nil {
			c.AddBMC(devicetypes.BMCTypeDPU, *info)
		}
	}

	return c
}

func strToID(s string, allowEmpty bool) int {
	if len(s) == 0 {
		if allowEmpty {
			return 0
		}
	} else {
		if i, err := strconv.Atoi(s); err == nil {
			return i
		}
	}

	return -1
}

func positionToRackName(Position string, expectedSegs int) string {
	segs := strings.Split(Position, "-")
	if len(segs) != expectedSegs {
		return ""
	}

	return strings.ToLower(segs[expectedSegs-1])
}
