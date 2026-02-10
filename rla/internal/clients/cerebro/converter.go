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
package cerebro

import (
	"net"
	"strings"

	"github.com/google/uuid"

	"github.com/nvidia/carbide-rest/rla/internal/inventory/objects/bmc"
	"github.com/nvidia/carbide-rest/rla/internal/inventory/objects/component"
	"github.com/nvidia/carbide-rest/rla/internal/inventory/objects/rack"
	"github.com/nvidia/carbide-rest/rla/pkg/common/deviceinfo"
	"github.com/nvidia/carbide-rest/rla/pkg/common/devicetypes"
	"github.com/nvidia/carbide-rest/rla/pkg/common/location"
)

type rackConverter struct {
	cerebroRack *Rack
	cache       *Cache
	devicesByID map[string]*Device
}

func newRackConverter(cerebroRack *Rack, cache *Cache) *rackConverter {
	if cerebroRack == nil || cache == nil {
		// Nothing to convert or cache is not initialized.
		return nil
	}

	devicesByID := make(map[string]*Device)
	for _, d := range cerebroRack.Devices {
		devicesByID[d.ID] = &d
	}

	return &rackConverter{
		cerebroRack: cerebroRack,
		cache:       cache,
		devicesByID: devicesByID,
	}
}

func (c *rackConverter) toRack() *rack.Rack {
	r := rack.New(
		deviceinfo.DeviceInfo{
			ID:           uuid.MustParse(c.cerebroRack.ID),
			Name:         c.cerebroRack.Name,
			SerialNumber: c.cerebroRack.Serial,
		},
		c.cerebroLocationToLocation(&c.cerebroRack.Location),
	)

	conflictComps := make([]component.Component, 0)
	for _, d := range c.cerebroRack.Devices {
		c := c.deviceToComponent(&d)
		if c != nil {
			if r.AddComponent(*c) == 0 {
				conflictComps = append(conflictComps, *c)
			}
		}
	}

	// Seal the rack to make sure the components are sorted by in-rack
	// position.
	r.Seal()

	// NOTE: a hacky way to fix the components with same serial numbers.
	if len(conflictComps) > 0 {
		r.Components = append(r.Components, conflictComps...)
	}

	return r
}

func (c *rackConverter) deviceToComponent(d *Device) *component.Component {
	if d == nil {
		return nil
	}

	typ := deviceToComponentType(d)
	if typ == devicetypes.ComponentTypeUnknown {
		// The type is not supported by RLA, skip it.
		return nil
	}

	comp := component.New(
		typ,
		&deviceinfo.DeviceInfo{
			ID:           uuid.MustParse(d.ID),
			Name:         d.Name,
			Manufacturer: d.Type.Manufacturer.Name,
			SerialNumber: d.Serial,
			Model:        d.Type.Model,
		},
		"", // indiate the firmware version is not available
		&component.InRackPosition{
			SlotID:    d.Position,
			TrayIndex: -1, // indiate the information is not available
			HostID:    -1, // indiate the information is not available
		},
	)

	for typ, bmcs := range c.buildAllBMCs(typ, d) {
		for _, bmc := range bmcs {
			comp.AddBMC(typ, bmc)
		}
	}

	return &comp
}

func deviceToComponentType(d *Device) devicetypes.ComponentType {
	switch d.Role.Name {
	case "GPU":
		if strings.Contains(d.Type.Model, "BlueField") {
			return devicetypes.ComponentTypeUnknown
		}

		return devicetypes.ComponentTypeCompute
	case "NVSwitch":
		return devicetypes.ComponentTypeNVLSwitch
	case "SMN-Leaf":
		return devicetypes.ComponentTypeToRSwitch
	case "NSV Device":
		return devicetypes.ComponentTypePowerShelf
	default:
		return devicetypes.ComponentTypeUnknown
	}
}

func (c *rackConverter) buildAllBMCs(
	typ devicetypes.ComponentType,
	d *Device,
) map[devicetypes.BMCType][]bmc.BMC {
	if d == nil || typ == devicetypes.ComponentTypeUnknown {
		// Nothing or unable to build.
		return nil
	}

	bmcsByType := make(map[devicetypes.BMCType][]bmc.BMC)

	buildBMCs := func(interfaces []Interface, subName string) []bmc.BMC {
		// Build the BMCs for the given type
		bmcs := make([]bmc.BMC, 0, 1)
		for _, iface := range interfaces {
			// Only scan the interfaces whose names contain the given sub name
			if !strings.Contains(iface.Name, subName) {
				continue
			}

			mac, err := net.ParseMAC(iface.MACAddress)
			if err != nil {
				continue
			}

			nb := bmc.BMC{MAC: mac}
			if len(iface.IPAddresses) > 0 {
				if ip, _, err := net.ParseCIDR(iface.IPAddresses[0].Address); err == nil {
					nb.IP = ip
				}
			}

			bmcs = append(bmcs, nb)
		}

		return bmcs
	}

	// Build the host BMCs
	if typ == devicetypes.ComponentTypePowerShelf {
		bmcsByType[devicetypes.BMCTypeHost] = buildBMCs(d.Interfaces, "eth")
	} else {
		bmcsByType[devicetypes.BMCTypeHost] = buildBMCs(d.Interfaces, "bmc")
	}

	// Build the DPU BMCs if the component is a compute
	if typ == devicetypes.ComponentTypeCompute {
		dpuBMCs := make([]bmc.BMC, 0, len(d.DeviceBays))
		for _, db := range d.DeviceBays {
			if dpu := c.devicesByID[db.InstalledDevice.ID]; dpu != nil {
				dpuBMCs = append(dpuBMCs, buildBMCs(dpu.Interfaces, "bmc")...)
			}
		}
		bmcsByType[devicetypes.BMCTypeDPU] = dpuBMCs
	}

	return bmcsByType
}

func (c *rackConverter) cerebroLocationToLocation(
	cerebroLocation *Location,
) location.Location {
	var loc location.Location

	cachedLoc := NewCachedLocation(cerebroLocation)
	// Traverse the location chain to fill in the location information.
	// We do not care the the root location type "Provider".
	for cachedLoc != nil && cachedLoc.Type != locationTypeRoot {
		switch cachedLoc.Type {
		case LocationTypeRegion:
			loc.Region = cachedLoc.Name
		case LocationTypeSite:
			loc.DataCenter = cachedLoc.Name
		case LocationTypeModule:
			loc.Room = cachedLoc.Name
		}

		cachedLoc = c.cache.GetLocation(cachedLoc.Parent)
	}

	return loc
}
