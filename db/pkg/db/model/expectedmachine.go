/*
 * SPDX-FileCopyrightText: Copyright (c) 2021-2023 NVIDIA CORPORATION & AFFILIATES. All rights reserved.
 * SPDX-License-Identifier: LicenseRef-NvidiaProprietary
 *
 * NVIDIA CORPORATION, its affiliates and licensors retain all intellectual
 * property and proprietary rights in and to this material, related
 * documentation and any modifications thereto. Any use, reproduction,
 * disclosure or distribution of this material and related documentation
 * without an express license agreement from NVIDIA CORPORATION or
 * its affiliates is strictly prohibited.
 */


package model

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/nvidia/carbide-rest/db/pkg/db"
	"github.com/nvidia/carbide-rest/db/pkg/db/paginator"

	"github.com/uptrace/bun"

	stracer "github.com/nvidia/carbide-rest/db/pkg/tracer"
)

const (
	// ExpectedMachineOrderByDefault default field to be used for ordering when none specified
	ExpectedMachineOrderByDefault = "created"
)

var (
	// ExpectedMachineOrderByFields is a list of valid order by fields for the ExpectedMachine model
	ExpectedMachineOrderByFields = []string{
		"id",
		"site_id",
		"bmc_mac_address",
		"chassis_serial_number",
		"created",
		"updated",
	}
	// ExpectedMachineRelatedEntities is a list of valid relation by fields for the ExpectedMachine model
	ExpectedMachineRelatedEntities = map[string]bool{
		SiteRelationName:    true,
		SkuRelationName:     true,
		MachineRelationName: true,
	}
)

// ExpectedMachine is a record for each bare-metal host expected to be processed by Forge
type ExpectedMachine struct {
	bun.BaseModel `bun:"table:expected_machine,alias:em"`

	ID                       uuid.UUID         `bun:"id,pk"`
	SiteID                   uuid.UUID         `bun:"site_id,type:uuid,notnull"`
	Site                     *Site             `bun:"rel:belongs-to,join:site_id=id"`
	BmcMacAddress            string            `bun:"bmc_mac_address,notnull"`
	ChassisSerialNumber      string            `bun:"chassis_serial_number,notnull"`
	SkuID                    *string           `bun:"sku_id"`
	Sku                      *SKU              `bun:"rel:belongs-to,join:sku_id=id"`
	MachineID                *string           `bun:"machine_id"`
	Machine                  *Machine          `bun:"rel:belongs-to,join:machine_id=id"`
	FallbackDpuSerialNumbers []string          `bun:"fallback_dpu_serial_numbers,array"`
	Labels                   map[string]string `bun:"labels,type:jsonb"`
	Created                  time.Time         `bun:"created,nullzero,notnull,default:current_timestamp"`
	Updated                  time.Time         `bun:"updated,nullzero,notnull,default:current_timestamp"`
	CreatedBy                uuid.UUID         `bun:"type:uuid,notnull"`
}

// ExpectedMachineCreateInput input parameters for Create method
type ExpectedMachineCreateInput struct {
	ExpectedMachineID        uuid.UUID
	SiteID                   uuid.UUID
	BmcMacAddress            string
	ChassisSerialNumber      string
	SkuID                    *string
	MachineID                *string
	FallbackDpuSerialNumbers []string
	Labels                   map[string]string
	CreatedBy                uuid.UUID
}

// ExpectedMachineUpdateInput input parameters for Update method
type ExpectedMachineUpdateInput struct {
	ExpectedMachineID        uuid.UUID
	BmcMacAddress            *string
	ChassisSerialNumber      *string
	SkuID                    *string
	MachineID                *string
	FallbackDpuSerialNumbers []string
	Labels                   map[string]string
}

// ExpectedMachineClearInput input parameters for Clear method
type ExpectedMachineClearInput struct {
	ExpectedMachineID        uuid.UUID
	SkuID                    bool
	MachineID                bool
	FallbackDpuSerialNumbers bool
	Labels                   bool
}

// ExpectedMachineFilterInput filtering options for GetAll method
type ExpectedMachineFilterInput struct {
	ExpectedMachineIDs   []uuid.UUID
	SiteIDs              []uuid.UUID
	BmcMacAddresses      []string
	ChassisSerialNumbers []string
	SkuIDs               []string
	MachineIDs           []string
	SearchQuery          *string
}

var _ bun.BeforeAppendModelHook = (*ExpectedMachine)(nil)

// BeforeAppendModel is a hook that is called before the model is appended to the query
func (em *ExpectedMachine) BeforeAppendModel(ctx context.Context, query bun.Query) error {
	switch query.(type) {
	case *bun.InsertQuery:
		em.Created = db.GetCurTime()
		em.Updated = db.GetCurTime()
	case *bun.UpdateQuery:
		em.Updated = db.GetCurTime()
	}
	return nil
}

var _ bun.BeforeCreateTableHook = (*ExpectedMachine)(nil)

// BeforeCreateTable is a hook that is called before the table is created
// This is only used in tests
func (em *ExpectedMachine) BeforeCreateTable(ctx context.Context, query *bun.CreateTableQuery) error {
	query.ForeignKey(`("site_id") REFERENCES "site" ("id")`).
		ForeignKey(`("sku_id") REFERENCES "sku" ("id")`).
		ForeignKey(`("machine_id") REFERENCES "machine" ("id")`)
	return nil
}

// ExpectedMachineDAO is an interface for interacting with the ExpectedMachine model
type ExpectedMachineDAO interface {
	// Create used to create new row
	Create(ctx context.Context, tx *db.Tx, input ExpectedMachineCreateInput) (*ExpectedMachine, error)
	// Update used to update row
	Update(ctx context.Context, tx *db.Tx, input ExpectedMachineUpdateInput) (*ExpectedMachine, error)
	// Delete used to delete row
	Delete(ctx context.Context, tx *db.Tx, expectedMachineID uuid.UUID) error
	// Clear used to clear fields in the row
	Clear(ctx context.Context, tx *db.Tx, input ExpectedMachineClearInput) (*ExpectedMachine, error)
	// GetAll returns all the rows based on the filter and page inputs
	GetAll(ctx context.Context, tx *db.Tx, filter ExpectedMachineFilterInput, page paginator.PageInput, includeRelations []string) ([]ExpectedMachine, int, error)
	// Get returns row for specified ID
	Get(ctx context.Context, tx *db.Tx, expectedMachineID uuid.UUID, includeRelations []string, forUpdate bool) (*ExpectedMachine, error)
}

// ExpectedMachineSQLDAO is an implementation of the ExpectedMachineDAO interface
type ExpectedMachineSQLDAO struct {
	dbSession  *db.Session
	tracerSpan *stracer.TracerSpan

	ExpectedMachineDAO
}

// Create creates a new ExpectedMachine from the given parameters
func (emsd ExpectedMachineSQLDAO) Create(ctx context.Context, tx *db.Tx, input ExpectedMachineCreateInput) (*ExpectedMachine, error) {
	// Create a child span and set the attributes for current request
	ctx, expectedMachineDAOSpan := emsd.tracerSpan.CreateChildInCurrentContext(ctx, "ExpectedMachineDAO.Create")
	if expectedMachineDAOSpan != nil {
		defer expectedMachineDAOSpan.End()
	}

	// NOTE: since Expected Machine can be created by Carbide or Cloud API the caller MUST provide the ID.
	em := &ExpectedMachine{
		ID:                       input.ExpectedMachineID,
		SiteID:                   input.SiteID,
		BmcMacAddress:            input.BmcMacAddress,
		ChassisSerialNumber:      input.ChassisSerialNumber,
		SkuID:                    input.SkuID,
		MachineID:                input.MachineID,
		FallbackDpuSerialNumbers: input.FallbackDpuSerialNumbers,
		Labels:                   input.Labels,
		CreatedBy:                input.CreatedBy,
	}

	_, err := db.GetIDB(tx, emsd.dbSession).NewInsert().Model(em).Exec(ctx)
	if err != nil {
		return nil, err
	}

	nv, err := emsd.Get(ctx, tx, em.ID, nil, false)
	if err != nil {
		return nil, err
	}

	return nv, nil
}

// Get returns an ExpectedMachine by ID
// returns db.ErrDoesNotExist error if the record is not found
func (emsd ExpectedMachineSQLDAO) Get(ctx context.Context, tx *db.Tx, expectedMachineID uuid.UUID, includeRelations []string, forUpdate bool) (*ExpectedMachine, error) {
	// Create a child span and set the attributes for current request
	ctx, expectedMachineDAOSpan := emsd.tracerSpan.CreateChildInCurrentContext(ctx, "ExpectedMachineDAO.Get")
	if expectedMachineDAOSpan != nil {
		defer expectedMachineDAOSpan.End()

		emsd.tracerSpan.SetAttribute(expectedMachineDAOSpan, "id", expectedMachineID.String())
	}

	em := &ExpectedMachine{}

	query := db.GetIDB(tx, emsd.dbSession).NewSelect().Model(em).Where("em.id = ?", expectedMachineID)

	if forUpdate {
		query = query.For("UPDATE")
	}

	for _, relation := range includeRelations {
		query = query.Relation(relation)
	}

	err := query.Scan(ctx)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, db.ErrDoesNotExist
		}
		return nil, err
	}

	return em, nil
}

// setQueryWithFilter populates the lookup query based on specified filter
func (emsd ExpectedMachineSQLDAO) setQueryWithFilter(filter ExpectedMachineFilterInput, query *bun.SelectQuery, expectedMachineDAOSpan *stracer.CurrentContextSpan) (*bun.SelectQuery, error) {
	if filter.SiteIDs != nil {
		query = query.Where("em.site_id IN (?)", bun.In(filter.SiteIDs))
		if expectedMachineDAOSpan != nil {
			emsd.tracerSpan.SetAttribute(expectedMachineDAOSpan, "site_ids", filter.SiteIDs)
		}
	}

	if filter.ExpectedMachineIDs != nil {
		query = query.Where("em.id IN (?)", bun.In(filter.ExpectedMachineIDs))
		if expectedMachineDAOSpan != nil {
			emsd.tracerSpan.SetAttribute(expectedMachineDAOSpan, "expected_machine_ids", filter.ExpectedMachineIDs)
		}
	}

	if filter.BmcMacAddresses != nil {
		query = query.Where("em.bmc_mac_address IN (?)", bun.In(filter.BmcMacAddresses))
		if expectedMachineDAOSpan != nil {
			emsd.tracerSpan.SetAttribute(expectedMachineDAOSpan, "bmc_mac_addresses", filter.BmcMacAddresses)
		}
	}

	if filter.ChassisSerialNumbers != nil {
		query = query.Where("em.chassis_serial_number IN (?)", bun.In(filter.ChassisSerialNumbers))
		if expectedMachineDAOSpan != nil {
			emsd.tracerSpan.SetAttribute(expectedMachineDAOSpan, "chassis_serial_numbers", filter.ChassisSerialNumbers)
		}
	}

	if filter.SkuIDs != nil {
		query = query.Where("em.sku_id IN (?)", bun.In(filter.SkuIDs))
		if expectedMachineDAOSpan != nil {
			emsd.tracerSpan.SetAttribute(expectedMachineDAOSpan, "sku_ids", filter.SkuIDs)
		}
	}

	if filter.MachineIDs != nil {
		query = query.Where("em.machine_id IN (?)", bun.In(filter.MachineIDs))
		if expectedMachineDAOSpan != nil {
			emsd.tracerSpan.SetAttribute(expectedMachineDAOSpan, "machine_ids", filter.MachineIDs)
		}
	}

	if filter.SearchQuery != nil {
		normalizedTokens := db.GetStrPtr(db.GetStringToTsQuery(*filter.SearchQuery))
		query = query.WhereGroup(" AND ", func(q *bun.SelectQuery) *bun.SelectQuery {
			return q.
				Where("to_tsvector('english', (coalesce(em.bmc_mac_address, ' ') || ' ' || coalesce(em.chassis_serial_number, ' ') || ' ' || coalesce(em.sku_id, ' ') || ' ' || coalesce(em.machine_id, ' ') || ' ' || coalesce(em.fallback_dpu_serial_numbers::text, ' ') || ' ' || coalesce(em.labels::text, ' '))) @@ to_tsquery('english', ?)", *normalizedTokens).
				WhereOr("em.bmc_mac_address ILIKE ?", "%"+*filter.SearchQuery+"%").
				WhereOr("em.chassis_serial_number ILIKE ?", "%"+*filter.SearchQuery+"%").
				WhereOr("em.sku_id ILIKE ?", "%"+*filter.SearchQuery+"%").
				WhereOr("em.machine_id ILIKE ?", "%"+*filter.SearchQuery+"%").
				WhereOr("em.fallback_dpu_serial_numbers::text ILIKE ?", "%"+*filter.SearchQuery+"%").
				WhereOr("em.labels::text ILIKE ?", "%"+*filter.SearchQuery+"%").
				WhereOr("em.id::text ILIKE ?", "%"+*filter.SearchQuery+"%").
				WhereOr("em.site_id::text ILIKE ?", "%"+*filter.SearchQuery+"%")
		})
		if expectedMachineDAOSpan != nil {
			emsd.tracerSpan.SetAttribute(expectedMachineDAOSpan, "search_query", *filter.SearchQuery)
		}
	}

	return query, nil
}

// GetAll returns all ExpectedMachines based on the filter and paging
// Errors are returned only when there is a db related error
// If records not found, then error is nil, but length of returned slice is 0
// If orderBy is nil, then records are ordered by column specified in ExpectedMachineOrderByDefault in ascending order
func (emsd ExpectedMachineSQLDAO) GetAll(ctx context.Context, tx *db.Tx, filter ExpectedMachineFilterInput, page paginator.PageInput, includeRelations []string) ([]ExpectedMachine, int, error) {
	// Create a child span and set the attributes for current request
	ctx, expectedMachineDAOSpan := emsd.tracerSpan.CreateChildInCurrentContext(ctx, "ExpectedMachineDAO.GetAll")
	if expectedMachineDAOSpan != nil {
		defer expectedMachineDAOSpan.End()
	}

	var expectedMachines []ExpectedMachine

	if filter.ExpectedMachineIDs != nil && len(filter.ExpectedMachineIDs) == 0 {
		return expectedMachines, 0, nil
	}

	query := db.GetIDB(tx, emsd.dbSession).NewSelect().Model(&expectedMachines)

	query, err := emsd.setQueryWithFilter(filter, query, expectedMachineDAOSpan)
	if err != nil {
		return expectedMachines, 0, err
	}

	// Apply relations if requested
	for _, relation := range includeRelations {
		query = query.Relation(relation)
	}

	// If no order is passed, set default order to make sure objects return always in the same order and pagination works properly
	if page.OrderBy == nil {
		page.OrderBy = paginator.NewDefaultOrderBy(ExpectedMachineOrderByDefault)
	}

	expectedMachinePaginator, err := paginator.NewPaginator(ctx, query, page.Offset, page.Limit, page.OrderBy, ExpectedMachineOrderByFields)
	if err != nil {
		return nil, 0, err
	}

	err = expectedMachinePaginator.Query.Limit(expectedMachinePaginator.Limit).Offset(expectedMachinePaginator.Offset).Scan(ctx)
	if err != nil {
		return nil, 0, err
	}

	return expectedMachines, expectedMachinePaginator.Total, nil
}

// Update updates specified fields of an existing ExpectedMachine
// The updated fields are assumed to be set to non-null values
func (emsd ExpectedMachineSQLDAO) Update(ctx context.Context, tx *db.Tx, input ExpectedMachineUpdateInput) (*ExpectedMachine, error) {
	// Create a child span and set the attributes for current request
	ctx, expectedMachineDAOSpan := emsd.tracerSpan.CreateChildInCurrentContext(ctx, "ExpectedMachineDAO.Update")
	if expectedMachineDAOSpan != nil {
		defer expectedMachineDAOSpan.End()
	}

	em := &ExpectedMachine{
		ID: input.ExpectedMachineID,
	}

	updatedFields := []string{}
	if input.BmcMacAddress != nil {
		em.BmcMacAddress = *input.BmcMacAddress
		updatedFields = append(updatedFields, "bmc_mac_address")

		if expectedMachineDAOSpan != nil {
			emsd.tracerSpan.SetAttribute(expectedMachineDAOSpan, "bmc_mac_address", *input.BmcMacAddress)
		}
	}
	if input.ChassisSerialNumber != nil {
		em.ChassisSerialNumber = *input.ChassisSerialNumber
		updatedFields = append(updatedFields, "chassis_serial_number")

		if expectedMachineDAOSpan != nil {
			emsd.tracerSpan.SetAttribute(expectedMachineDAOSpan, "chassis_serial_number", *input.ChassisSerialNumber)
		}
	}
	if input.FallbackDpuSerialNumbers != nil {
		em.FallbackDpuSerialNumbers = input.FallbackDpuSerialNumbers
		updatedFields = append(updatedFields, "fallback_dpu_serial_numbers")

		if expectedMachineDAOSpan != nil {
			emsd.tracerSpan.SetAttribute(expectedMachineDAOSpan, "fallback_dpu_serial_numbers", input.FallbackDpuSerialNumbers)
		}
	}
	if input.Labels != nil {
		em.Labels = input.Labels
		updatedFields = append(updatedFields, "labels")

		if expectedMachineDAOSpan != nil {
			emsd.tracerSpan.SetAttribute(expectedMachineDAOSpan, "labels", input.Labels)
		}
	}
	if input.SkuID != nil {
		em.SkuID = input.SkuID
		updatedFields = append(updatedFields, "sku_id")

		if expectedMachineDAOSpan != nil {
			emsd.tracerSpan.SetAttribute(expectedMachineDAOSpan, "sku_id", *input.SkuID)
		}
	}
	if input.MachineID != nil {
		em.MachineID = input.MachineID
		updatedFields = append(updatedFields, "machine_id")

		if expectedMachineDAOSpan != nil {
			emsd.tracerSpan.SetAttribute(expectedMachineDAOSpan, "machine_id", *input.MachineID)
		}
	}

	if len(updatedFields) > 0 {
		updatedFields = append(updatedFields, "updated")

		_, err := db.GetIDB(tx, emsd.dbSession).NewUpdate().Model(em).Column(updatedFields...).Where("id = ?", input.ExpectedMachineID).Exec(ctx)
		if err != nil {
			return nil, err
		}
	}

	nv, err := emsd.Get(ctx, tx, em.ID, nil, false)

	if err != nil {
		return nil, err
	}

	return nv, nil
}

// Clear sets parameters of an existing ExpectedMachine to null values in db
func (emsd ExpectedMachineSQLDAO) Clear(ctx context.Context, tx *db.Tx, input ExpectedMachineClearInput) (*ExpectedMachine, error) {
	// Create a child span and set the attributes for current request
	ctx, expectedMachineDAOSpan := emsd.tracerSpan.CreateChildInCurrentContext(ctx, "ExpectedMachineDAO.Clear")
	if expectedMachineDAOSpan != nil {
		defer expectedMachineDAOSpan.End()
	}

	em := &ExpectedMachine{
		ID: input.ExpectedMachineID,
	}

	updatedFields := []string{}
	if input.SkuID {
		em.SkuID = nil
		updatedFields = append(updatedFields, "sku_id")
	}
	if input.MachineID {
		em.MachineID = nil
		updatedFields = append(updatedFields, "machine_id")
	}
	if input.FallbackDpuSerialNumbers {
		em.FallbackDpuSerialNumbers = nil
		updatedFields = append(updatedFields, "fallback_dpu_serial_numbers")
	}
	if input.Labels {
		em.Labels = nil
		updatedFields = append(updatedFields, "labels")
	}

	if len(updatedFields) > 0 {
		updatedFields = append(updatedFields, "updated")

		_, err := db.GetIDB(tx, emsd.dbSession).NewUpdate().Model(em).Column(updatedFields...).Where("id = ?", input.ExpectedMachineID).Exec(ctx)
		if err != nil {
			return nil, err
		}
	}

	nv, err := emsd.Get(ctx, tx, input.ExpectedMachineID, nil, false)
	if err != nil {
		return nil, err
	}
	return nv, nil
}

// Delete deletes an ExpectedMachine by ID
// Error is returned only if there is a db error
func (emsd ExpectedMachineSQLDAO) Delete(ctx context.Context, tx *db.Tx, expectedMachineID uuid.UUID) error {
	// Create a child span and set the attributes for current request
	ctx, expectedMachineDAOSpan := emsd.tracerSpan.CreateChildInCurrentContext(ctx, "ExpectedMachineDAO.Delete")
	if expectedMachineDAOSpan != nil {
		defer expectedMachineDAOSpan.End()

		emsd.tracerSpan.SetAttribute(expectedMachineDAOSpan, "id", expectedMachineID.String())
	}

	em := &ExpectedMachine{
		ID: expectedMachineID,
	}

	var err error

	_, err = db.GetIDB(tx, emsd.dbSession).NewDelete().Model(em).Where("id = ?", expectedMachineID).Exec(ctx)
	if err != nil {
		return err
	}

	return nil
}

// NewExpectedMachineDAO returns a new ExpectedMachineDAO
func NewExpectedMachineDAO(dbSession *db.Session) ExpectedMachineDAO {
	return &ExpectedMachineSQLDAO{
		dbSession:  dbSession,
		tracerSpan: stracer.NewTracerSpan(),
	}
}
