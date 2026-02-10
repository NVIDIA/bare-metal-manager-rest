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
// Package store provides the storage layer for inventory management.
// It defines the InventoryStore interface for persisting and retrieving
// inventory data (racks, components, NVL domains).
package store

import (
	"context"

	"github.com/google/uuid"

	"github.com/nvidia/carbide-rest/rla/internal/inventory/objects/component"
	"github.com/nvidia/carbide-rest/rla/internal/inventory/objects/nvldomain"
	"github.com/nvidia/carbide-rest/rla/internal/inventory/objects/rack"
	identifier "github.com/nvidia/carbide-rest/rla/pkg/common/Identifier"
	dbquery "github.com/nvidia/carbide-rest/rla/pkg/db/query"
)

// Store defines the interface for inventory data persistence.
// It provides operations for managing racks (with components) and NVL domains.
type Store interface {
	// Lifecycle
	Start(ctx context.Context) error
	Stop(ctx context.Context) error

	// Rack operations
	CreateExpectedRack(ctx context.Context, rack *rack.Rack) (uuid.UUID, error)
	GetRackByID(ctx context.Context, id uuid.UUID, withComponents bool) (*rack.Rack, error)
	GetRacksByIDs(ctx context.Context, ids []uuid.UUID, withComponents bool) ([]*rack.Rack, error)
	GetRackBySerial(ctx context.Context, manufacturer string, serial string, withComponents bool) (*rack.Rack, error)
	GetRackByIdentifier(ctx context.Context, identifier identifier.Identifier, withComponents bool) (*rack.Rack, error)
	PatchRack(ctx context.Context, rack *rack.Rack) (string, error)
	GetListOfRacks(ctx context.Context, info dbquery.StringQueryInfo, pagination *dbquery.Pagination, withComponents bool) ([]*rack.Rack, int32, error)

	// Component operations
	GetComponentByID(ctx context.Context, id uuid.UUID) (*component.Component, error)
	GetComponentBySerial(ctx context.Context, manufacturer string, serial string, withRack bool) (*component.Component, error)
	GetComponentByBMCMAC(ctx context.Context, macAddress string) (*component.Component, error)
	GetComponentsByExternalIDs(ctx context.Context, externalIDs []string) ([]*component.Component, error)

	// NVL Domain operations
	CreateNVLDomain(ctx context.Context, nvlDomain *nvldomain.NVLDomain) (uuid.UUID, error)
	AttachRacksToNVLDomain(ctx context.Context, nvlDomainID identifier.Identifier, rackIDs []identifier.Identifier) error
	DetachRacksFromNVLDomain(ctx context.Context, rackIDs []identifier.Identifier) error
	GetListOfNVLDomains(ctx context.Context, info dbquery.StringQueryInfo, pagination *dbquery.Pagination) ([]*nvldomain.NVLDomain, int32, error)
	GetRacksForNVLDomain(ctx context.Context, nvlDomainID identifier.Identifier) ([]*rack.Rack, error)
}
