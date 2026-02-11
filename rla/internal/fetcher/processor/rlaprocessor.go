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
package processor

import (
	"context"

	"github.com/rs/zerolog/log"

	"github.com/nvidia/carbide-rest/rla/internal/inventory/objects/bmc"
	"github.com/nvidia/carbide-rest/rla/internal/inventory/objects/component"
	"github.com/nvidia/carbide-rest/rla/internal/inventory/objects/nvldomain"
	"github.com/nvidia/carbide-rest/rla/internal/inventory/objects/rack"
	"github.com/nvidia/carbide-rest/rla/pkg/client"
	identifier "github.com/nvidia/carbide-rest/rla/pkg/common/Identifier"
	"github.com/nvidia/carbide-rest/rla/pkg/common/devicetypes"
	"github.com/nvidia/carbide-rest/rla/pkg/types"
)

type RLAProcessor struct {
	rlaClient *client.Client
}

func newRLAProcessor(c client.Config) (*RLAProcessor, error) {
	rlaClient, err := client.New(c)
	if err != nil {
		return nil, err
	}

	return &RLAProcessor{
		rlaClient: rlaClient,
	}, nil
}

func (rp *RLAProcessor) Process(ctx context.Context, item *ProcessableItem) error {
	for _, r := range item.Racks {
		if r == nil {
			continue
		}

		_, err := rp.rlaClient.CreateExpectedRack(ctx, rackToTypes(r))
		if err != nil {
			log.Error().Err(err).Msgf("failed to create expected rack %s", r.Info.Name)
			return err
		}
	}

	if item.NVLDomain != nil {
		_, err := rp.rlaClient.CreateNVLDomain(ctx, nvlDomainToTypes(item.NVLDomain))
		if err != nil {
			log.Error().Err(err).Msg("failed to create NVL domain")
			return err
		}

		err = rp.rlaClient.AttachRacksToNVLDomain(
			ctx,
			identifierToTypes(item.NVLDomain.Identifier),
			identifiersToTypes(item.NVLDomain.RackIdentifiers),
		)

		if err != nil {
			log.Error().Err(err).Msg("failed to attach racks to NVL domain")
			return err
		}
	}

	return nil
}

func (rp *RLAProcessor) Type() Type {
	return TypeRLA
}

func (rp *RLAProcessor) Done() error {
	return rp.rlaClient.Close()
}

// Conversion functions from internal types to pkg/types

func rackToTypes(r *rack.Rack) *types.Rack {
	if r == nil {
		return nil
	}

	tr := &types.Rack{
		Info: types.DeviceInfo{
			ID:           r.Info.ID,
			Name:         r.Info.Name,
			Manufacturer: r.Info.Manufacturer,
			Model:        r.Info.Model,
			SerialNumber: r.Info.SerialNumber,
			Description:  r.Info.Description,
		},
		Location: types.Location{
			Region:     r.Loc.Region,
			Datacenter: r.Loc.DataCenter,
			Room:       r.Loc.Room,
			Position:   r.Loc.Position,
		},
	}

	if len(r.Components) > 0 {
		tr.Components = make([]types.Component, 0, len(r.Components))
		for i := range r.Components {
			tr.Components = append(tr.Components, componentToTypes(&r.Components[i]))
		}
	}

	return tr
}

func componentToTypes(c *component.Component) types.Component {
	if c == nil {
		return types.Component{}
	}

	comp := types.Component{
		Type: componentTypeToTypes(c.Type),
		Info: types.DeviceInfo{
			ID:           c.Info.ID,
			Name:         c.Info.Name,
			Manufacturer: c.Info.Manufacturer,
			Model:        c.Info.Model,
			SerialNumber: c.Info.SerialNumber,
			Description:  c.Info.Description,
		},
		FirmwareVersion: c.FirmwareVersion,
		Position: types.InRackPosition{
			SlotID:    c.Position.SlotID,
			TrayIndex: c.Position.TrayIndex,
			HostID:    c.Position.HostID,
		},
		ComponentID: c.ComponentID,
		RackID:      c.RackID,
	}

	// Convert BMCs
	for _, bmcType := range devicetypes.BMCTypes() {
		for _, b := range c.BmcsByType[bmcType] {
			comp.BMCs = append(comp.BMCs, bmcToTypes(&b, bmcType))
		}
	}

	return comp
}

func bmcToTypes(b *bmc.BMC, bmcType devicetypes.BMCType) types.BMC {
	if b == nil {
		return types.BMC{}
	}

	tb := types.BMC{
		Type: bmcTypeToTypes(bmcType),
		MAC:  b.MAC,
		IP:   b.IP,
	}

	if b.Credential != nil {
		tb.User = b.Credential.User
		tb.Password = b.Credential.Password.String()
	}

	return tb
}

func componentTypeToTypes(ct devicetypes.ComponentType) types.ComponentType {
	switch ct {
	case devicetypes.ComponentTypeCompute:
		return types.ComponentTypeCompute
	case devicetypes.ComponentTypeNVLSwitch:
		return types.ComponentTypeNVSwitch
	case devicetypes.ComponentTypePowerShelf:
		return types.ComponentTypePowerShelf
	case devicetypes.ComponentTypeToRSwitch:
		return types.ComponentTypeTORSwitch
	case devicetypes.ComponentTypeUMS:
		return types.ComponentTypeUMS
	case devicetypes.ComponentTypeCDU:
		return types.ComponentTypeCDU
	default:
		return types.ComponentTypeUnknown
	}
}

func bmcTypeToTypes(bt devicetypes.BMCType) types.BMCType {
	switch bt {
	case devicetypes.BMCTypeHost:
		return types.BMCTypeHost
	case devicetypes.BMCTypeDPU:
		return types.BMCTypeDPU
	default:
		return types.BMCTypeUnknown
	}
}

func nvlDomainToTypes(d *nvldomain.NVLDomain) *types.NVLDomain {
	if d == nil {
		return nil
	}

	return &types.NVLDomain{
		ID:   d.Identifier.ID,
		Name: d.Identifier.Name,
	}
}

func identifierToTypes(id identifier.Identifier) types.Identifier {
	return types.Identifier{
		ID:   id.ID,
		Name: id.Name,
	}
}

func identifiersToTypes(ids []identifier.Identifier) []types.Identifier {
	result := make([]types.Identifier, 0, len(ids))
	for _, id := range ids {
		result = append(result, identifierToTypes(id))
	}
	return result
}
