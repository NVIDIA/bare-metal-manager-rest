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
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"

	"github.com/nvidia/carbide-rest/rla/pkg/common/deviceinfo"
	"github.com/nvidia/carbide-rest/rla/pkg/common/utils"
	dbquery "github.com/nvidia/carbide-rest/rla/pkg/db/query"
)

var defaultRackPagination = dbquery.Pagination{
	Offset: 0,
	Limit:  100,
	Total:  0,
}

type Rack struct {
	bun.BaseModel `bun:"table:rack,alias:r"`

	ID           uuid.UUID      `bun:"id,pk,type:uuid,default:gen_random_uuid()"`
	Name         string         `bun:"name,notnull,unique:rack_name_idx"`
	Manufacturer string         `bun:"manufacturer,notnull,unique:rack_manufacturer_serial_idx"`
	SerialNumber string         `bun:"serial_number,notnull,unique:rack_manufacturer_serial_idx"`
	Description  map[string]any `bun:"description,type:jsonb,json_use_number"`
	Location     map[string]any `bun:"location,type:jsonb"`
	NVLDomainID  uuid.UUID      `bun:"nvldomain_id,type:uuid"`
	Status       RackStatus     `bun:"status,type:varchar(16),default:'new'"`
	CreatedAt    time.Time      `bun:"created_at,nullzero,notnull,default:current_timestamp"`
	UpdatedAt    time.Time      `bun:"updated_at,nullzero,notnull,default:current_timestamp"`
	IngestedAt   *time.Time     `bun:"ingested_at"`
	DeletedAt    *time.Time     `bun:"deleted_at,soft_delete"`
	Components   []Component    `bun:"rel:has-many,join:id=rack_id"`
	NVLDomain    *NVLDomain     `bun:"rel:belongs-to,join:nvldomain_id=id"`
}

type RackStatus string

const (
	RackStatusNew       RackStatus = "new"
	RackStatusIngesting RackStatus = "ingesting"
	RackStatusIngested  RackStatus = "ingested"
)

/*
var _ bun.BeforeAppendModelHook = (*Rack)(nil)

func (r *Rack) BeforeAppendModel(ctx context.Context, query bun.Query) error {
	switch query.(type) {
	case *bun.InsertQuery:
		r.CreatedAt = db.CurTime()
		r.UpdatedAt = db.CurTime()
	case *bun.UpdateQuery:
		r.UpdatedAt = db.CurTime()
	}
	return nil
}
*/

func (rd *Rack) Create(ctx context.Context, idb bun.IDB) error {
	_, err := idb.NewInsert().Model(rd).Exec(ctx)
	return err
}

func (rd *Rack) Get(
	ctx context.Context,
	idb bun.IDB,
	withComponents bool,
) (*Rack, error) {
	var rack Rack
	var query *bun.SelectQuery

	if rd.ID != uuid.Nil {
		// Get by the ID
		query = idb.NewSelect().Model(&rack).Where("id = ?", rd.ID)
	} else {
		// Get by the serial information
		query = idb.NewSelect().Model(&rack).Where(
			"manufacturer = ? AND serial_number = ?",
			rd.Manufacturer, rd.SerialNumber,
		)
	}

	if withComponents {
		query = query.Relation("Components").Relation("Components.BMCs")
	}

	if err := query.Scan(ctx); err != nil {
		return nil, err
	}

	return &rack, nil
}

func (rd *Rack) Patch(ctx context.Context, idb bun.IDB) error {
	_, err := idb.NewUpdate().Model(rd).Where("id = ?", rd.ID).Exec(ctx)
	return err
}

// BuildPatch builds a patched rack from the current rack and the
// input rack. It goes through the patchable fields and builds the patched
// rack. If there is no change on patchable fields, it returns nil.
func (rd *Rack) BuildPatch(cur *Rack) *Rack {
	if rd == nil || cur == nil {
		return nil
	}

	// The patchable fields include:
	// Name
	// Description
	// Location

	// Make a copy fo the current rack which serves as the base for the
	// patched rack.
	patchedRack := *cur
	patched := false

	if len(rd.Name) > 0 && patchedRack.Name != rd.Name {
		patchedRack.Name = rd.Name
		patched = true
	}

	if desc := utils.CompareAndCopyMaps(rd.Description, cur.Description); desc != nil {
		patchedRack.Description = desc
		patched = true
	}

	if loc := utils.CompareAndCopyMaps(rd.Location, cur.Location); loc != nil {
		patchedRack.Location = loc
		patched = true
	}

	if !patched {
		return nil
	}

	return &patchedRack
}

// SerialInfo returns the serial number information of the rack.
func (rd *Rack) SerialInfo() deviceinfo.SerialInfo {
	return deviceinfo.SerialInfo{
		Manufacturer: rd.Manufacturer,
		SerialNumber: rd.SerialNumber,
	}
}

func GetListOfRacks(
	ctx context.Context,
	idb bun.IDB,
	info dbquery.StringQueryInfo,
	pagination *dbquery.Pagination,
	withComponents bool,
) ([]Rack, int32, error) {
	var racks []Rack
	conf := &dbquery.Config{
		IDB:   idb,
		Model: &racks,
	}

	if pagination != nil {
		conf.Pagination = pagination
	} else {
		conf.Pagination = &defaultRackPagination
	}

	if filterable := info.ToFilterable("name"); filterable != nil {
		conf.Filterables = []dbquery.Filterable{filterable}
	}

	if withComponents {
		conf.Relations = []string{"Components", "Components.BMCs"}
	}

	q, err := dbquery.New(ctx, conf)
	if err != nil {
		return nil, 0, err
	}

	if err := q.Scan(ctx); err != nil {
		return nil, 0, err
	}

	return racks, int32(q.TotalCount()), nil
}

func GetRacksForNVLDomain(
	ctx context.Context,
	idb bun.IDB,
	nvlDomainID uuid.UUID,
) ([]Rack, error) {
	var racks []Rack
	q := idb.NewSelect().Model(&racks).Where("nvldomain_id = ?", nvlDomainID)

	if err := q.Scan(ctx); err != nil {
		return nil, err
	}

	return racks, nil
}

// GetRacksByIDs retrieves multiple racks by their UUIDs
func GetRacksByIDs(
	ctx context.Context,
	idb bun.IDB,
	ids []uuid.UUID,
	withComponents bool,
) ([]Rack, error) {
	var racks []Rack
	q := idb.NewSelect().Model(&racks).Where("id IN (?)", bun.In(ids))

	if withComponents {
		q = q.Relation("Components").Relation("Components.BMCs")
	}

	if err := q.Scan(ctx); err != nil {
		return nil, err
	}

	return racks, nil
}
