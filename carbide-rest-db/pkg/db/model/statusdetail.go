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
	"github.com/NVIDIA/carbide-rest-api/carbide-rest-db/pkg/db"
	"github.com/NVIDIA/carbide-rest-api/carbide-rest-db/pkg/db/paginator"
	"github.com/uptrace/bun"

	stracer "github.com/NVIDIA/carbide-rest-api/carbide-rest-db/pkg/tracer"
)

// StatusDetail represents entries in the status_detail table
type StatusDetail struct {
	bun.BaseModel `bun:"table:status_detail,alias:sd"`

	ID       uuid.UUID `bun:"type:uuid,pk"`
	EntityID string    `bun:"entity_id"`
	Status   string    `bun:"status,notnull"`
	Message  *string   `bun:"message"`
	Count    int       `bun:"count,notnull"`
	Created  time.Time `bun:"created,nullzero,notnull,default:current_timestamp"`
	Updated  time.Time `bun:"updated,nullzero,notnull,default:current_timestamp"`
}

const (
	// StatusDetailRelationName is the relation name for the StatusDetail model
	StatusDetailRelationName = "StatusDetail"
)

var (
	// StatusDetailOrderByFields is a list of valid order by fields for the StatusDetail model
	StatusDetailOrderByFields = []string{"status", "created", "updated"}
	// StatusDetailOrderByDefault default field to be used for ordering when none specified
	StatusDetailOrderByDefault = "created"
)

var _ bun.BeforeAppendModelHook = (*InfrastructureProvider)(nil)

// BeforeAppendModel is a hook that is called before the model is appended to the query
func (sd *StatusDetail) BeforeAppendModel(ctx context.Context, query bun.Query) error {
	switch query.(type) {
	case *bun.InsertQuery:
		sd.Created = db.GetCurTime()
		sd.Updated = db.GetCurTime()
	case *bun.UpdateQuery:
		sd.Updated = db.GetCurTime()
	}
	return nil
}

// StatusDetailDAO is the data access interface for StatusDetail
type StatusDetailDAO interface {
	//
	GetAllByEntityID(ctx context.Context, tx *db.Tx, entityID string, offset *int, limit *int, orderBy *paginator.OrderBy) ([]StatusDetail, int, error)
	//
	GetAllByEntityIDs(ctx context.Context, tx *db.Tx, entityIDs []string, offset *int, limit *int, orderBy *paginator.OrderBy) ([]StatusDetail, int, error)
	//
	GetByID(ctx context.Context, tx *db.Tx, id uuid.UUID) (*StatusDetail, error)
	//
	CreateFromParams(ctx context.Context, tx *db.Tx, entityID string, status string, message *string) (*StatusDetail, error)
	//
	UpdateFromParams(ctx context.Context, tx *db.Tx, id uuid.UUID, status string, message *string) (*StatusDetail, error)
	// GetRecentByEntityIDs returns most recent status records for specified entity IDs
	GetRecentByEntityIDs(ctx context.Context, tx *db.Tx, entityIDs []string, recentCount int) ([]StatusDetail, error)
}

// StatusDetailSQLDAO is the data access object for StatusDetail
type StatusDetailSQLDAO struct {
	dbSession  *db.Session
	tracerSpan *stracer.TracerSpan
}

// GetByID returns a StatusDetail by ID
func (sdd StatusDetailSQLDAO) GetByID(ctx context.Context, tx *db.Tx, id uuid.UUID) (*StatusDetail, error) {
	// Create a child span and set the attributes for current request
	ctx, sdDAOSpan := sdd.tracerSpan.CreateChildInCurrentContext(ctx, "StatusDetailDAO.GetByID")
	if sdDAOSpan != nil {
		defer sdDAOSpan.End()

		sdd.tracerSpan.SetAttribute(sdDAOSpan, "id", id.String())
	}

	sd := &StatusDetail{}

	err := db.GetIDB(tx, sdd.dbSession).NewSelect().Model(sd).Where("id = ?", id).Scan(ctx)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, db.ErrDoesNotExist
		}
		return nil, err
	}

	return sd, nil
}

// GetAllByEntityID returns status details for the given entity ID
func (sdd StatusDetailSQLDAO) GetAllByEntityID(ctx context.Context, tx *db.Tx, entityID string, offset *int, limit *int, orderBy *paginator.OrderBy) ([]StatusDetail, int, error) {
	// Create a child span and set the attributes for current request
	ctx, sdDAOSpan := sdd.tracerSpan.CreateChildInCurrentContext(ctx, "StatusDetailDAO.GetAllByEntityID")
	if sdDAOSpan != nil {
		defer sdDAOSpan.End()

		sdd.tracerSpan.SetAttribute(sdDAOSpan, "entityID", entityID)
	}

	sds := []StatusDetail{}

	query := db.GetIDB(tx, sdd.dbSession).NewSelect().Model(&sds).Where("entity_id = ?", entityID)

	// StatusDetail has a default order by of created desc
	normalizedOrderBy := &paginator.OrderBy{
		Field: "created",
		Order: paginator.OrderDescending,
	}
	if orderBy != nil {
		normalizedOrderBy = orderBy
	}

	paginator, err := paginator.NewPaginator(ctx, query, offset, limit, normalizedOrderBy, StatusDetailOrderByFields)
	if err != nil {
		return nil, 0, err
	}

	err = paginator.Query.Limit(paginator.Limit).Offset(paginator.Offset).Scan(ctx)
	if err != nil {
		return nil, 0, err
	}

	return sds, paginator.Total, nil
}

// GetAllByEntityIDs returns status details for the given set of entity IDs
func (sdd StatusDetailSQLDAO) GetAllByEntityIDs(ctx context.Context, tx *db.Tx, entityIDs []string, offset *int, limit *int, orderBy *paginator.OrderBy) ([]StatusDetail, int, error) {
	// Create a child span and set the attributes for current request
	ctx, sdDAOSpan := sdd.tracerSpan.CreateChildInCurrentContext(ctx, "StatusDetailDAO.GetAllByEntityIDs")
	if sdDAOSpan != nil {
		defer sdDAOSpan.End()
	}

	sds := []StatusDetail{}

	if len(entityIDs) == 0 {
		return sds, 0, nil
	}

	query := db.GetIDB(tx, sdd.dbSession).NewSelect().Model(&sds).Where("entity_id IN (?)", bun.In(entityIDs))

	// StatusDetail has a default order by of created desc
	normalizedOrderBy := &paginator.OrderBy{
		Field: "created",
		Order: paginator.OrderDescending,
	}
	if orderBy != nil {
		normalizedOrderBy = orderBy
	}

	paginator, err := paginator.NewPaginator(ctx, query, offset, limit, normalizedOrderBy, StatusDetailOrderByFields)
	if err != nil {
		return nil, 0, err
	}

	err = paginator.Query.Limit(paginator.Limit).Offset(paginator.Offset).Scan(ctx)
	if err != nil {
		return nil, 0, err
	}

	return sds, paginator.Total, nil
}

// CreateFromParams creates a new StatusDetail from the given parameters
func (sdd StatusDetailSQLDAO) CreateFromParams(ctx context.Context, tx *db.Tx, entityID string, status string, message *string) (*StatusDetail, error) {
	// Create a child span and set the attributes for current request
	ctx, sdDAOSpan := sdd.tracerSpan.CreateChildInCurrentContext(ctx, "StatusDetailDAO.CreateFromParams")
	if sdDAOSpan != nil {
		defer sdDAOSpan.End()
		sdd.tracerSpan.SetAttribute(sdDAOSpan, "entityID", entityID)

	}

	sd := &StatusDetail{
		ID:       uuid.New(),
		EntityID: entityID,
		Status:   status,
		Message:  message,
		Count:    1,
	}

	_, err := db.GetIDB(tx, sdd.dbSession).NewInsert().Model(sd).Exec(ctx)
	if err != nil {
		return nil, err
	}

	return sdd.GetByID(ctx, tx, sd.ID)
}

// UpdateFromParams updates the given StatusDetail with the given parameters
func (sdd StatusDetailSQLDAO) UpdateFromParams(ctx context.Context, tx *db.Tx, id uuid.UUID, status string, message *string) (*StatusDetail, error) {
	// Create a child span and set the attributes for current request
	ctx, sdDAOSpan := sdd.tracerSpan.CreateChildInCurrentContext(ctx, "StatusDetailDAO.UpdateFromParams")
	if sdDAOSpan != nil {
		defer sdDAOSpan.End()

		sdd.tracerSpan.SetAttribute(sdDAOSpan, "id", id.String())
	}

	sd := &StatusDetail{}

	err := db.GetIDB(tx, sdd.dbSession).NewSelect().Model(sd).Where("id = ?", id).Scan(ctx)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, db.ErrDoesNotExist
		}
		return nil, err
	}

	if status == "" {
		return nil, db.ErrInvalidValue
	}

	upsd := &StatusDetail{
		ID:       sd.ID,
		EntityID: sd.EntityID,
	}

	updatedFields := []string{}
	if sd.Status != status {
		upsd.Status = status
		updatedFields = append(updatedFields, "status")
		sdd.tracerSpan.SetAttribute(sdDAOSpan, "status", status)
	}

	if sd.Message != message {
		upsd.Message = message
		updatedFields = append(updatedFields, "message")
		sdd.tracerSpan.SetAttribute(sdDAOSpan, "message", message)
	}

	if len(updatedFields) == 0 {
		return sd, nil
	}

	upsd.Count = sd.Count + 1
	updatedFields = append(updatedFields, "count")

	updatedFields = append(updatedFields, "updated")

	_, err = db.GetIDB(tx, sdd.dbSession).NewUpdate().Model(upsd).Column(updatedFields...).Where("entity_id = ?", sd.EntityID).Exec(ctx)
	if err != nil {
		return nil, err
	}

	nsd, err := sdd.GetByID(ctx, tx, sd.ID)
	if err != nil {
		return nil, err
	}

	return nsd, nil
}

// GetRecentByEntityIDs returns most recent status records for specified entity IDs
func (sdd StatusDetailSQLDAO) GetRecentByEntityIDs(ctx context.Context, tx *db.Tx, entityIDs []string, recentCount int) ([]StatusDetail, error) {
	// Create a child span and set the attributes for current request
	ctx, sdDAOSpan := sdd.tracerSpan.CreateChildInCurrentContext(ctx, "StatusDetailDAO.GetRecentByEntityIDs")
	if sdDAOSpan != nil {
		defer sdDAOSpan.End()
	}

	sds := []StatusDetail{}

	if len(entityIDs) == 0 {
		return sds, nil
	}

	sqlQuery := `SELECT id, entity_id, status, message, count, created, updated FROM (
					SELECT id, entity_id, status, message, count, created, updated, row_number() OVER (
						PARTITION BY entity_id ORDER BY created DESC
					) rn FROM status_detail WHERE entity_id IN (?)
				) t WHERE rn <= ?;`

	query := db.GetIDB(tx, sdd.dbSession).NewRaw(sqlQuery, bun.In(entityIDs), recentCount)

	err := query.Scan(ctx, &sds)
	return sds, err
}

// NewStatusDetailDAO creates and returns a new data access object for StatusDetail
func NewStatusDetailDAO(dbSession *db.Session) StatusDetailDAO {
	return StatusDetailSQLDAO{
		dbSession:  dbSession,
		tracerSpan: stracer.NewTracerSpan(),
	}
}
