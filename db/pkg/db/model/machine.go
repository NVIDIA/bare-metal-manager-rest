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
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/mitchellh/mapstructure"
	"github.com/nvidia/carbide-rest/db/pkg/db"
	"github.com/nvidia/carbide-rest/db/pkg/db/paginator"
	cwssaws "github.com/nvidia/carbide-rest/workflow-schema/schema/site-agent/workflows/v1"
	"google.golang.org/protobuf/encoding/protojson"

	"github.com/uptrace/bun"

	stracer "github.com/nvidia/carbide-rest/db/pkg/tracer"
)

// Represents status of the machine
const (
	// MachineStatusInitializing indicates that the Machine is Initializing
	MachineStatusInitializing = "Initializing"
	// MachineStatusReady indicates that the Machine is ready
	MachineStatusReady = "Ready"
	// MachineStatusReset indicates that the Machine is being reset
	MachineStatusReset = "Reset"
	// MachineStatusMaintenance indicates that the Machine is in maintenance mode
	MachineStatusMaintenance = "Maintenance"
	// MachineStatusInUse indicates that the Machine is being used by an Instance
	MachineStatusInUse = "InUse"
	// MachineStatusError indicates that the Machine is in error state
	MachineStatusError = "Error"
	// MachineStatusDecommissioned indicates that the Machine was decommissioned
	MachineStatusDecommissioned = "Decommissioned"
	// MachineStatusUnknown indicates that the Machine status cannot be determined
	MachineStatusUnknown = "Unknown"
	// MachineRelationName is the relation name for the Machine model
	MachineRelationName = "Machine"

	// MachineOrderByDefault default field to be used for ordering when none specified
	MachineOrderByDefault = "created"
)

var (
	// MachineOrderByFields is a list of valid order by fields for the Machine model
	MachineOrderByFields = []string{"id", "status", "created", "updated"}
	// MachineRelatedEntities is a list of valid relation by fields for the Machine model
	MachineRelatedEntities = map[string]bool{
		InfrastructureProviderRelationName: true,
		SiteRelationName:                   true,
		InstanceTypeRelationName:           true,
	}
	// MachineStatusMap is a list of valid status for the Machine model
	MachineStatusMap = map[string]bool{
		MachineStatusInitializing:   true,
		MachineStatusReady:          true,
		MachineStatusReset:          true,
		MachineStatusMaintenance:    true,
		MachineStatusInUse:          true,
		MachineStatusError:          true,
		MachineStatusDecommissioned: true,
		MachineStatusUnknown:        true,
	}
)

// A light wrapper around the protobuf so
// that we can implement our own marshal/unmarshal
// that understands how to work with protobuf messages
type SiteControllerMachine struct {
	*cwssaws.Machine
}

func (s *SiteControllerMachine) UnmarshalJSON(b []byte) error {
	if s.Machine == nil {
		s.Machine = &cwssaws.Machine{}
	}

	// We intentionally ignore the error here.
	// If the "metadata" column can't be successfully parsed,
	// we don't want to prevent the entire machine record from being read.
	// It _would_ be nice to see if we can log this somehow, though.
	_ = protoJsonUnmarshalOptions.Unmarshal(b, s)

	return nil
}

func (s *SiteControllerMachine) MarshalJSON() ([]byte, error) {
	return protojson.Marshal(s)
}

// Machine is the baremetal server that sits in the datacenter
type Machine struct {
	bun.BaseModel `bun:"table:machine,alias:m"`

	ID                       string                  `bun:"id,pk"`
	InfrastructureProviderID uuid.UUID               `bun:"infrastructure_provider_id,type:uuid,notnull"`
	InfrastructureProvider   *InfrastructureProvider `bun:"rel:belongs-to,join:infrastructure_provider_id=id"`
	SiteID                   uuid.UUID               `bun:"site_id,type:uuid,notnull"`
	Site                     *Site                   `bun:"rel:belongs-to,join:site_id=id"`
	InstanceTypeID           *uuid.UUID              `bun:"instance_type_id,type:uuid"`
	InstanceType             *InstanceType           `bun:"rel:belongs-to,join:instance_type_id=id"`
	ControllerMachineID      string                  `bun:"controller_machine_id,notnull"`
	ControllerMachineType    *string                 `bun:"controller_machine_type"`
	HwSkuDeviceType          *string                 `bun:"hw_sku_device_type"`
	Vendor                   *string                 `bun:"vendor"`
	ProductName              *string                 `bun:"product_name"`
	SerialNumber             *string                 `bun:"serial_number"`
	Metadata                 *SiteControllerMachine  `bun:"metadata,type:jsonb"`
	IsInMaintenance          bool                    `bun:"is_in_maintenance,notnull"`
	IsAllocatable            bool                    `bun:"is_allocatable"`
	MaintenanceMessage       *string                 `bun:"maintenance_message"`
	IsNetworkDegraded        bool                    `bun:"is_network_degraded,notnull"`
	NetworkHealthMessage     *string                 `bun:"network_health_message"`
	Health                   map[string]interface{}  `bun:"health,type:jsonb,json_use_number"`
	DefaultMacAddress        *string                 `bun:"default_mac_address"`
	Hostname                 *string                 `bun:"hostname"`
	IsAssigned               bool                    `bun:"is_assigned,notnull"` // To be deprecate in favor of checking InstanceTypeID not null
	Labels                   map[string]string       `bun:"labels,type:jsonb"`   // Labels are used to store additional metadata about the machine
	Status                   string                  `bun:"status,notnull"`
	IsMissingOnSite          bool                    `bun:"is_missing_on_site,notnull"`
	Created                  time.Time               `bun:"created,nullzero,notnull,default:current_timestamp"`
	Updated                  time.Time               `bun:"updated,nullzero,notnull,default:current_timestamp"`
	Deleted                  *time.Time              `bun:"deleted,soft_delete"`
}

// MachineCreateInput input parameters for Create method
type MachineCreateInput struct {
	MachineID                string
	InfrastructureProviderID uuid.UUID
	SiteID                   uuid.UUID
	InstanceTypeID           *uuid.UUID
	ControllerMachineID      string
	ControllerMachineType    *string
	HwSkuDeviceType          *string
	Vendor                   *string
	ProductName              *string
	SerialNumber             *string
	Metadata                 *SiteControllerMachine
	IsInMaintenance          bool
	IsAllocatable            bool
	MaintenanceMessage       *string
	IsNetworkDegraded        bool
	NetworkHealthMessage     *string
	Health                   map[string]interface{}
	DefaultMacAddress        *string
	Hostname                 *string
	Status                   string
	Labels                   map[string]string
}

// MachineUpdateInput input parameters for Update method
type MachineUpdateInput struct {
	MachineID                string
	InfrastructureProviderID *uuid.UUID
	SiteID                   *uuid.UUID
	InstanceTypeID           *uuid.UUID
	ControllerMachineID      *string
	ControllerMachineType    *string
	HwSkuDeviceType          *string
	Vendor                   *string
	ProductName              *string
	SerialNumber             *string
	Metadata                 *SiteControllerMachine
	IsInMaintenance          *bool
	IsAllocatable            *bool
	MaintenanceMessage       *string
	IsNetworkDegraded        *bool
	NetworkHealthMessage     *string
	Health                   map[string]interface{}
	DefaultMacAddress        *string
	Hostname                 *string
	IsAssigned               *bool
	Status                   *string
	Labels                   map[string]string
	IsMissingOnSite          *bool
}

// MachineClearInput input parameters for Clear method
type MachineClearInput struct {
	MachineID             string
	InstanceTypeID        bool
	ControllerMachineType bool
	HwSkuDeviceType       bool
	Vendor                bool
	ProductName           bool
	SerialNumber          bool
	Metadata              bool
	MaintenanceMessage    bool
	Health                bool
	NetworkHealthMessage  bool
	DefaultMacAddress     bool
	Hostname              bool
}

// MachineFilterInput filtering options for GetAll method
type MachineFilterInput struct {
	InfrastructureProviderID *uuid.UUID
	SiteID                   *uuid.UUID
	HasInstanceType          *bool
	InstanceTypeIDs          []uuid.UUID
	ControllerMachineID      *string
	HwSkuDeviceTypes         []string
	IsAssigned               *bool
	Hostname                 *string
	CapabilityType           *string
	CapabilityNames          []string
	Statuses                 []string
	SearchQuery              *string
	MachineIDs               []string
}

type MachineHealth struct {
	Source     string               `json:"source"`
	ObservedAt *string              `json:"observed_at"`
	Successes  []HealthProbeSuccess `json:"successes"`
	Alerts     []HealthProbeAlert   `json:"alerts"`
}

type HealthProbeSuccess struct {
	Id     string  `json:"id"`
	Target *string `json:"target"`
}

type HealthProbeAlert struct {
	Id              string   `json:"id"`
	Target          *string  `json:"target"`
	InAlertSince    *string  `json:"in_alert_since"`
	Message         string   `json:"message"`
	TenantMessage   *string  `json:"tenant_message"`
	Classifications []string `json:"classifications"`
}

// GetIndentedJSON returns formatted json of Machine
func (m *Machine) GetIndentedJSON() ([]byte, error) {
	return json.MarshalIndent(m, "", "  ")
}

var _ bun.BeforeAppendModelHook = (*Machine)(nil)

// BeforeAppendModel is a hook that is called before the model is appended to the query
func (m *Machine) BeforeAppendModel(ctx context.Context, query bun.Query) error {
	switch query.(type) {
	case *bun.InsertQuery:
		m.Created = db.GetCurTime()
		m.Updated = db.GetCurTime()
	case *bun.UpdateQuery:
		m.Updated = db.GetCurTime()
	}
	return nil
}

var _ bun.BeforeCreateTableHook = (*Machine)(nil)

// BeforeCreateTable is a hook that is called before the table is created
// This is only used in tests
func (m *Machine) BeforeCreateTable(ctx context.Context, query *bun.CreateTableQuery) error {
	query.ForeignKey(`("infrastructure_provider_id") REFERENCES "infrastructure_provider" ("id")`).
		ForeignKey(`("instance_type_id") REFERENCES "instance_type" ("id")`).
		ForeignKey(`("site_id") REFERENCES "site" ("id")`)
	return nil
}

// returns db.ErrDoesNotExist error if the record is not found
func (m *Machine) GetHealth() (*MachineHealth, error) {
	var health *MachineHealth
	serr := mapstructure.Decode(m.Health, &health)

	if serr != nil {
		return nil, db.ErrInvalidValue
	}

	return health, nil
}

// MachineDAO is an interface for interacting with the Machine model
type MachineDAO interface {
	// Create used to create new row
	Create(ctx context.Context, tx *db.Tx, input MachineCreateInput) (*Machine, error)
	// Update used to update row
	Update(ctx context.Context, tx *db.Tx, input MachineUpdateInput) (*Machine, error)
	// Delete used to delete row
	Delete(ctx context.Context, tx *db.Tx, machineID string, purge bool) error
	// Clear used to clear fields in the row
	Clear(ctx context.Context, tx *db.Tx, input MachineClearInput) (*Machine, error)
	// GetAll returns all the rows based on the filter and page inputs
	GetAll(ctx context.Context, tx *db.Tx, filter MachineFilterInput, page paginator.PageInput, includeRelations []string) ([]Machine, int, error)
	// GetByID returns row for specified ID
	GetByID(ctx context.Context, tx *db.Tx, machineID string, includeRelations []string, forUpdate bool) (*Machine, error)
	// GetCountByStatus returns row counts per status
	GetCountByStatus(ctx context.Context, tx *db.Tx, infrastructureProviderID *uuid.UUID, siteID *uuid.UUID, instanceTypeID *uuid.UUID) (map[string]int, error)
	// GetCount returns total count of rows for specified filter
	GetCount(ctx context.Context, tx *db.Tx, filter MachineFilterInput) (count int, err error)
	// GetHealth returns the machine's health deserialized from json
	GetHealth(ctx context.Context, tx *db.Tx, machineID string, includeRelations []string) (*MachineHealth, error)
}

// MachineSQLDAO is an implementation of the MachineDAO interface
type MachineSQLDAO struct {
	dbSession *db.Session
	MachineDAO
	tracerSpan *stracer.TracerSpan
}

// Create creates a new Machine from the given parameters
// The returned Machine will not have any related structs filled in
// since there are 2 operations (INSERT, SELECT), in this, it is required that
// this library call happens within a transaction
func (msd MachineSQLDAO) Create(ctx context.Context, tx *db.Tx, input MachineCreateInput) (*Machine, error) {
	// Create a child span and set the attributes for current request
	ctx, machineDAOSpan := msd.tracerSpan.CreateChildInCurrentContext(ctx, "MachineDAO.Create")
	if machineDAOSpan != nil {
		defer machineDAOSpan.End()
	}

	m := &Machine{
		ID:                       input.MachineID,
		InfrastructureProviderID: input.InfrastructureProviderID,
		SiteID:                   input.SiteID,
		InstanceTypeID:           input.InstanceTypeID,
		ControllerMachineID:      input.ControllerMachineID,
		ControllerMachineType:    input.ControllerMachineType,
		HwSkuDeviceType:          input.HwSkuDeviceType,
		Vendor:                   input.Vendor,
		ProductName:              input.ProductName,
		SerialNumber:             input.SerialNumber,
		Metadata:                 input.Metadata,
		IsInMaintenance:          input.IsInMaintenance,
		IsAllocatable:            input.IsAllocatable,
		MaintenanceMessage:       input.MaintenanceMessage,
		IsNetworkDegraded:        input.IsNetworkDegraded,
		NetworkHealthMessage:     input.NetworkHealthMessage,
		Health:                   input.Health,
		DefaultMacAddress:        input.DefaultMacAddress,
		Hostname:                 input.Hostname,
		IsAssigned:               false,
		Status:                   input.Status,
		Labels:                   input.Labels,
		IsMissingOnSite:          false,
	}

	_, err := db.GetIDB(tx, msd.dbSession).NewInsert().Model(m).Exec(ctx)
	if err != nil {
		return nil, err
	}

	nv, err := msd.GetByID(ctx, tx, m.ID, nil, false)
	if err != nil {
		return nil, err
	}

	return nv, nil
}

// GetByID returns a Machine by ID
// returns db.ErrDoesNotExist error if the record is not found
func (msd MachineSQLDAO) GetByID(ctx context.Context, tx *db.Tx, id string, includeRelations []string, forUpdate bool) (*Machine, error) {
	// Create a child span and set the attributes for current request
	ctx, machineDAOSpan := msd.tracerSpan.CreateChildInCurrentContext(ctx, "MachineDAO.GetByID")
	if machineDAOSpan != nil {
		defer machineDAOSpan.End()

		msd.tracerSpan.SetAttribute(machineDAOSpan, "id", id)
	}

	m := &Machine{}

	query := db.GetIDB(tx, msd.dbSession).NewSelect().Model(m).Where("m.id = ?", id)

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

	return m, nil
}

// GetCountByStatus returns count of Machines for given status
// Errors are returned only when there is a db related error
// if records not found, then error is nil, but length of returned map is 0
func (msd MachineSQLDAO) GetCountByStatus(ctx context.Context, tx *db.Tx, infrastructureProviderID *uuid.UUID, siteID *uuid.UUID, instanceTypeID *uuid.UUID) (map[string]int, error) {
	// Create a child span and set the attributes for current request
	ctx, machineDAOSpan := msd.tracerSpan.CreateChildInCurrentContext(ctx, "MachineDAO.GetCountByStatus")
	if machineDAOSpan != nil {
		defer machineDAOSpan.End()
	}

	m := &Machine{}
	var statusQueryResults []map[string]interface{}

	query := db.GetIDB(tx, msd.dbSession).NewSelect().Model(m)
	if infrastructureProviderID != nil {
		query = query.Where("m.infrastructure_provider_id = ?", *infrastructureProviderID)

		if machineDAOSpan != nil {
			msd.tracerSpan.SetAttribute(machineDAOSpan, "infrastructure_provider_id", infrastructureProviderID.String())
		}
	}
	if siteID != nil {
		query = query.Where("m.site_id = ?", *siteID)

		if machineDAOSpan != nil {
			msd.tracerSpan.SetAttribute(machineDAOSpan, "site_id", siteID.String())
		}
	}
	if instanceTypeID != nil {
		query = query.Where("m.instance_type_id = ?", *instanceTypeID)

		if machineDAOSpan != nil {
			msd.tracerSpan.SetAttribute(machineDAOSpan, "instance_type_id", instanceTypeID.String())
		}
	}

	err := query.Column("m.status").ColumnExpr("COUNT(*) AS total_count").GroupExpr("m.status").Scan(ctx, &statusQueryResults)
	if err != nil {
		return nil, err
	}

	// creare results map by holding key as status value with total count
	results := map[string]int{
		"total":                     0,
		MachineStatusUnknown:        0,
		MachineStatusInitializing:   0,
		MachineStatusReady:          0,
		MachineStatusInUse:          0,
		MachineStatusDecommissioned: 0,
		MachineStatusError:          0,
		MachineStatusReset:          0,
		MachineStatusMaintenance:    0,
	}
	if len(statusQueryResults) > 0 {
		for _, statusMap := range statusQueryResults {
			results[statusMap["status"].(string)] = int(statusMap["total_count"].(int64))
			results["total"] += int(statusMap["total_count"].(int64))
		}
	}
	return results, nil
}

func (msd MachineSQLDAO) setQueryWithFilter(filter MachineFilterInput, query *bun.SelectQuery, machineDAOSpan *stracer.CurrentContextSpan) (*bun.SelectQuery, error) {
	if filter.InfrastructureProviderID != nil {
		query = query.Where("m.infrastructure_provider_id = ?", *filter.InfrastructureProviderID)

		if machineDAOSpan != nil {
			msd.tracerSpan.SetAttribute(machineDAOSpan, "infrastructure_provider_id", filter.InfrastructureProviderID.String())
		}
	}

	if filter.SiteID != nil {
		query = query.Where("m.site_id = ?", *filter.SiteID)

		if machineDAOSpan != nil {
			msd.tracerSpan.SetAttribute(machineDAOSpan, "site_id", filter.SiteID.String())
		}
	}

	if filter.HasInstanceType != nil {
		if *filter.HasInstanceType {
			query = query.Where("m.instance_type_id IS NOT NULL")
		} else if filter.InstanceTypeIDs != nil {
			return nil, errors.New("InstanceTypeID cannot be specified when hasInstanceType is false")
		} else {
			query = query.Where("m.instance_type_id IS NULL")
		}

		if machineDAOSpan != nil {
			msd.tracerSpan.SetAttribute(machineDAOSpan, "has_instancetype", *filter.HasInstanceType)
		}
	}

	if filter.InstanceTypeIDs != nil {
		if len(filter.InstanceTypeIDs) == 1 {
			query = query.Where("m.instance_type_id = ?", filter.InstanceTypeIDs[0])
		} else {
			query = query.Where("m.instance_type_id IN (?)", bun.In(filter.InstanceTypeIDs))
		}

		if machineDAOSpan != nil {
			msd.tracerSpan.SetAttribute(machineDAOSpan, "instance_type_ids", filter.InstanceTypeIDs)
		}
	}

	if filter.ControllerMachineID != nil {
		query = query.Where("m.controller_machine_id = ?", *filter.ControllerMachineID)

		if machineDAOSpan != nil {
			msd.tracerSpan.SetAttribute(machineDAOSpan, "controller_machine_id", *filter.ControllerMachineID)
		}
	}

	if filter.HwSkuDeviceTypes != nil {
		query = query.Where("m.hw_sku_device_type IN (?)", bun.In(filter.HwSkuDeviceTypes))

		if machineDAOSpan != nil {
			msd.tracerSpan.SetAttribute(machineDAOSpan, "hw_sku_device_types", filter.HwSkuDeviceTypes)
		}
	}

	if filter.Hostname != nil {
		query = query.Where("m.hostname = ?", *filter.Hostname)

		if machineDAOSpan != nil {
			msd.tracerSpan.SetAttribute(machineDAOSpan, "hostname", *filter.Hostname)
		}
	}

	if filter.IsAssigned != nil {
		query = query.Where("m.is_assigned = ?", *filter.IsAssigned)

		if machineDAOSpan != nil {
			msd.tracerSpan.SetAttribute(machineDAOSpan, "is_assigned", *filter.IsAssigned)
		}
	}

	if filter.CapabilityType != nil || filter.CapabilityNames != nil {

		query = query.Join("JOIN machine_capability AS mc ON mc.machine_id = m.id AND mc.deleted IS NULL").
			Distinct()

		if filter.CapabilityType != nil {
			query = query.Where("mc.type = ?", *filter.CapabilityType)
			if machineDAOSpan != nil {
				msd.tracerSpan.SetAttribute(machineDAOSpan, "capability_type", *filter.CapabilityType)
			}
		}
		if filter.CapabilityNames != nil {
			if len(filter.CapabilityNames) == 1 {
				query = query.Where("mc.name = ?", filter.CapabilityNames[0])
			} else {
				query = query.Where("mc.name IN (?)", bun.In(filter.CapabilityNames))
			}
			if machineDAOSpan != nil {
				msd.tracerSpan.SetAttribute(machineDAOSpan, "capability_names", filter.CapabilityNames)
			}
		}
	}

	if filter.Statuses != nil {
		if len(filter.Statuses) == 1 {
			query = query.Where("m.status = ?", filter.Statuses[0])
		} else {
			query = query.Where("m.status IN (?)", bun.In(filter.Statuses))
		}

		if machineDAOSpan != nil {
			msd.tracerSpan.SetAttribute(machineDAOSpan, "statuses", filter.Statuses)
		}
	}

	if filter.SearchQuery != nil {
		normalizedTokens := db.GetStrPtr(db.GetStringToTsQuery(*filter.SearchQuery))
		query = query.WhereGroup(" AND ", func(q *bun.SelectQuery) *bun.SelectQuery {
			return q.
				Where("to_tsvector('english', (coalesce(m.id, ' ') || ' ' || coalesce(m.vendor, ' ') || ' ' || coalesce(m.product_name, ' ') || ' ' || coalesce(m.hostname, ' ') || ' ' || coalesce(m.status, ' ') || ' ' || coalesce(m.labels::text, ' '))) @@ to_tsquery('english', ?)", *normalizedTokens).
				WhereOr("m.id ILIKE ?", "%"+*filter.SearchQuery+"%").
				WhereOr("m.vendor ILIKE ?", "%"+*filter.SearchQuery+"%").
				WhereOr("m.product_name ILIKE ?", "%"+*filter.SearchQuery+"%").
				WhereOr("m.hostname ILIKE ?", "%"+*filter.SearchQuery+"%").
				WhereOr("m.status ILIKE ?", "%"+*filter.SearchQuery+"%").
				WhereOr("m.labels::text ILIKE ?", "%"+*filter.SearchQuery+"%")
		})
		if machineDAOSpan != nil {
			msd.tracerSpan.SetAttribute(machineDAOSpan, "search_query", *filter.SearchQuery)
		}
	}

	if filter.MachineIDs != nil {
		query = query.Where("m.id IN (?)", bun.In(filter.MachineIDs))
	}
	return query, nil
}

// GetAll returns all Machines based on the filter and paging
// Errors are returned only when there is a db related error
// if records not found, then error is nil, but length of returned slice is 0
// if orderBy is nil, then records are ordered by column specified in MachineOrderByDefault in ascending order
func (msd MachineSQLDAO) GetAll(ctx context.Context, tx *db.Tx, filter MachineFilterInput, page paginator.PageInput, includeRelations []string) ([]Machine, int, error) {
	// Create a child span and set the attributes for current request
	ctx, machineDAOSpan := msd.tracerSpan.CreateChildInCurrentContext(ctx, "MachineDAO.GetAll")
	if machineDAOSpan != nil {
		defer machineDAOSpan.End()
	}

	var machines []Machine

	if filter.MachineIDs != nil && len(filter.MachineIDs) == 0 {
		return machines, 0, nil
	}

	if filter.HwSkuDeviceTypes != nil && len(filter.HwSkuDeviceTypes) == 0 {
		return machines, 0, nil
	}

	query := db.GetIDB(tx, msd.dbSession).NewSelect().Model(&machines)

	query, err := msd.setQueryWithFilter(filter, query, machineDAOSpan)
	if err != nil {
		return machines, 0, err
	}

	for _, relation := range includeRelations {
		query = query.Relation(relation)
	}

	// if no order is passed, set default to make sure objects return always in the same order and pagination works properly
	if page.OrderBy == nil {
		page.OrderBy = paginator.NewDefaultOrderBy(MachineOrderByDefault)
	}

	machinePaginator, err := paginator.NewPaginator(ctx, query, page.Offset, page.Limit, page.OrderBy, MachineOrderByFields)
	if err != nil {
		return nil, 0, err
	}

	err = machinePaginator.Query.Limit(machinePaginator.Limit).Offset(machinePaginator.Offset).Scan(ctx)
	if err != nil {
		return nil, 0, err
	}

	return machines, machinePaginator.Total, nil
}

// Update updates specified fields of an existing Machine
// The updated fields are assumed to be set to non-null values
// since there are 2 operations (UPDATE, SELECT), in this, it is required that
// this library call happens within a transaction
func (msd MachineSQLDAO) Update(ctx context.Context, tx *db.Tx, input MachineUpdateInput) (*Machine, error) {
	// Create a child span and set the attributes for current request
	ctx, machineDAOSpan := msd.tracerSpan.CreateChildInCurrentContext(ctx, "MachineDAO.Update")
	if machineDAOSpan != nil {
		defer machineDAOSpan.End()
	}

	m := &Machine{
		ID: input.MachineID,
	}

	updatedFields := []string{}
	if input.InfrastructureProviderID != nil {
		m.InfrastructureProviderID = *input.InfrastructureProviderID
		updatedFields = append(updatedFields, "infrastructure_provider_id")

		if machineDAOSpan != nil {
			msd.tracerSpan.SetAttribute(machineDAOSpan, "infrastructure_provider_id", input.InfrastructureProviderID.String())
		}
	}
	if input.SiteID != nil {
		m.SiteID = *input.SiteID
		updatedFields = append(updatedFields, "site_id")

		if machineDAOSpan != nil {
			msd.tracerSpan.SetAttribute(machineDAOSpan, "site_id", input.SiteID.String())
		}
	}
	if input.InstanceTypeID != nil {
		m.InstanceTypeID = input.InstanceTypeID
		updatedFields = append(updatedFields, "instance_type_id")

		if machineDAOSpan != nil {
			msd.tracerSpan.SetAttribute(machineDAOSpan, "instance_type_id", input.InstanceTypeID.String())
		}
	}
	if input.ControllerMachineID != nil {
		m.ControllerMachineID = *input.ControllerMachineID
		updatedFields = append(updatedFields, "controller_machine_id")

		if machineDAOSpan != nil {
			msd.tracerSpan.SetAttribute(machineDAOSpan, "controller_machine_id", *input.ControllerMachineID)
		}
	}
	if input.ControllerMachineType != nil {
		m.ControllerMachineType = input.ControllerMachineType
		updatedFields = append(updatedFields, "controller_machine_type")

		if machineDAOSpan != nil {
			msd.tracerSpan.SetAttribute(machineDAOSpan, "controller_machine_type", *input.ControllerMachineType)
		}
	}
	if input.HwSkuDeviceType != nil {
		m.HwSkuDeviceType = input.HwSkuDeviceType
		updatedFields = append(updatedFields, "hw_sku_device_type")

		if machineDAOSpan != nil {
			msd.tracerSpan.SetAttribute(machineDAOSpan, "hw_sku_device_type", *input.HwSkuDeviceType)
		}
	}
	if input.Vendor != nil {
		m.Vendor = input.Vendor
		updatedFields = append(updatedFields, "vendor")

		if machineDAOSpan != nil {
			msd.tracerSpan.SetAttribute(machineDAOSpan, "vendor", *input.Vendor)
		}
	}
	if input.ProductName != nil {
		m.ProductName = input.ProductName
		updatedFields = append(updatedFields, "product_name")

		if machineDAOSpan != nil {
			msd.tracerSpan.SetAttribute(machineDAOSpan, "product_name", *input.ProductName)
		}
	}
	if input.SerialNumber != nil {
		m.SerialNumber = input.SerialNumber
		updatedFields = append(updatedFields, "serial_number")

		if machineDAOSpan != nil {
			msd.tracerSpan.SetAttribute(machineDAOSpan, "serial_number", *input.SerialNumber)
		}
	}
	if input.Metadata != nil {
		m.Metadata = input.Metadata
		updatedFields = append(updatedFields, "metadata")
	}
	if input.IsAllocatable != nil {
		m.IsAllocatable = *input.IsAllocatable
		updatedFields = append(updatedFields, "is_allocatable")

		if machineDAOSpan != nil {
			msd.tracerSpan.SetAttribute(machineDAOSpan, "is_allocatable", *input.IsAllocatable)
		}
	}
	if input.IsInMaintenance != nil {
		m.IsInMaintenance = *input.IsInMaintenance
		updatedFields = append(updatedFields, "is_in_maintenance")

		if machineDAOSpan != nil {
			msd.tracerSpan.SetAttribute(machineDAOSpan, "is_in_maintenance", *input.IsInMaintenance)
		}
	}
	if input.MaintenanceMessage != nil {
		m.MaintenanceMessage = input.MaintenanceMessage
		updatedFields = append(updatedFields, "maintenance_message")

		if machineDAOSpan != nil {
			msd.tracerSpan.SetAttribute(machineDAOSpan, "maintenance_message", *input.MaintenanceMessage)
		}
	}
	if input.IsNetworkDegraded != nil {
		m.IsNetworkDegraded = *input.IsNetworkDegraded
		updatedFields = append(updatedFields, "is_network_degraded")

		if machineDAOSpan != nil {
			msd.tracerSpan.SetAttribute(machineDAOSpan, "is_network_degraded", *input.IsNetworkDegraded)
		}
	}
	if input.NetworkHealthMessage != nil {
		m.NetworkHealthMessage = input.NetworkHealthMessage
		updatedFields = append(updatedFields, "network_health_message")

		if machineDAOSpan != nil {
			msd.tracerSpan.SetAttribute(machineDAOSpan, "network_health_message", *input.NetworkHealthMessage)
		}
	}
	if input.Health != nil {
		m.Health = input.Health
		updatedFields = append(updatedFields, "health")
	}
	if input.DefaultMacAddress != nil {
		m.DefaultMacAddress = input.DefaultMacAddress
		updatedFields = append(updatedFields, "default_mac_address")

		if machineDAOSpan != nil {
			msd.tracerSpan.SetAttribute(machineDAOSpan, "default_mac_address", *input.DefaultMacAddress)
		}
	}
	if input.Hostname != nil {
		m.Hostname = input.Hostname
		updatedFields = append(updatedFields, "hostname")

		if machineDAOSpan != nil {
			msd.tracerSpan.SetAttribute(machineDAOSpan, "hostname", *input.Hostname)
		}
	}
	if input.IsAssigned != nil {
		m.IsAssigned = *input.IsAssigned
		updatedFields = append(updatedFields, "is_assigned")

		if machineDAOSpan != nil {
			msd.tracerSpan.SetAttribute(machineDAOSpan, "is_assigned", *input.IsAssigned)
		}
	}
	if input.Status != nil {
		m.Status = *input.Status
		updatedFields = append(updatedFields, "status")

		if machineDAOSpan != nil {
			msd.tracerSpan.SetAttribute(machineDAOSpan, "status", *input.Status)
		}
	}
	if input.Labels != nil {
		m.Labels = input.Labels
		updatedFields = append(updatedFields, "labels")

		if machineDAOSpan != nil {
			msd.tracerSpan.SetAttribute(machineDAOSpan, "labels", input.Labels)
		}
	}
	if input.IsMissingOnSite != nil {
		m.IsMissingOnSite = *input.IsMissingOnSite
		updatedFields = append(updatedFields, "is_missing_on_site")

		if machineDAOSpan != nil {
			msd.tracerSpan.SetAttribute(machineDAOSpan, "is_missing_on_site", *input.IsMissingOnSite)
		}
	}

	if len(updatedFields) > 0 {
		updatedFields = append(updatedFields, "updated")

		_, err := db.GetIDB(tx, msd.dbSession).NewUpdate().Model(m).Column(updatedFields...).Where("id = ?", input.MachineID).Exec(ctx)
		if err != nil {
			return nil, err
		}
	}

	nv, err := msd.GetByID(ctx, tx, m.ID, nil, false)

	if err != nil {
		return nil, err
	}
	return nv, nil
}

// Clear sets parameters of an existing Machine to null values in db
// since there are 2 operations (UPDATE, SELECT), it is required that
// this must be within a transaction
func (msd MachineSQLDAO) Clear(ctx context.Context, tx *db.Tx, input MachineClearInput) (*Machine, error) {
	// Create a child span and set the attributes for current request
	ctx, machineDAOSpan := msd.tracerSpan.CreateChildInCurrentContext(ctx, "MachineDAO.Clear")
	if machineDAOSpan != nil {
		defer machineDAOSpan.End()
	}

	m := &Machine{
		ID: input.MachineID,
	}

	updatedFields := []string{}
	if input.InstanceTypeID {
		m.InstanceTypeID = nil
		updatedFields = append(updatedFields, "instance_type_id")
	}
	if input.ControllerMachineType {
		m.ControllerMachineType = nil
		updatedFields = append(updatedFields, "controller_machine_type")
	}
	if input.HwSkuDeviceType {
		m.HwSkuDeviceType = nil
		updatedFields = append(updatedFields, "hw_sku_device_type")
	}
	if input.Vendor {
		m.Vendor = nil
		updatedFields = append(updatedFields, "vendor")
	}
	if input.ProductName {
		m.ProductName = nil
		updatedFields = append(updatedFields, "product_name")
	}
	if input.SerialNumber {
		m.SerialNumber = nil
		updatedFields = append(updatedFields, "serial_number")
	}
	if input.Metadata {
		m.Metadata = nil
		updatedFields = append(updatedFields, "metadata")
	}
	if input.MaintenanceMessage {
		m.MaintenanceMessage = nil
		updatedFields = append(updatedFields, "maintenance_message")
	}
	if input.NetworkHealthMessage {
		m.NetworkHealthMessage = nil
		updatedFields = append(updatedFields, "network_health_message")
	}
	if input.Health {
		m.Health = nil
		updatedFields = append(updatedFields, "health")
	}
	if input.DefaultMacAddress {
		m.DefaultMacAddress = nil
		updatedFields = append(updatedFields, "default_mac_address")
	}
	if input.Hostname {
		m.Hostname = nil
		updatedFields = append(updatedFields, "hostname")
	}

	if len(updatedFields) > 0 {
		updatedFields = append(updatedFields, "updated")

		_, err := db.GetIDB(tx, msd.dbSession).NewUpdate().Model(m).Column(updatedFields...).Where("id = ?", input.MachineID).Exec(ctx)
		if err != nil {
			return nil, err
		}
	}

	nv, err := msd.GetByID(ctx, tx, input.MachineID, nil, false)
	if err != nil {
		return nil, err
	}
	return nv, nil
}

// Delete deletes an Machine by ID
// error is returned only if there is a db error
// if the object being deleted doesnt exist, error is not returned (idempotent delete)
func (msd MachineSQLDAO) Delete(ctx context.Context, tx *db.Tx, machineID string, purge bool) error {
	// Create a child span and set the attributes for current request
	ctx, machineDAOSpan := msd.tracerSpan.CreateChildInCurrentContext(ctx, "MachineDAO.Delete")
	if machineDAOSpan != nil {
		defer machineDAOSpan.End()

		msd.tracerSpan.SetAttribute(machineDAOSpan, "id", machineID)
	}

	m := &Machine{
		ID: machineID,
	}

	var err error

	if purge {
		_, err = db.GetIDB(tx, msd.dbSession).NewDelete().Model(m).Where("id = ?", machineID).ForceDelete().Exec(ctx)
	} else {
		_, err = db.GetIDB(tx, msd.dbSession).NewDelete().Model(m).Where("id = ?", machineID).Exec(ctx)
	}
	if err != nil {
		return err
	}

	return nil
}

func (msd MachineSQLDAO) GetCount(ctx context.Context, tx *db.Tx, filter MachineFilterInput) (count int, err error) {
	// Create a child span and set the attributes for current request
	ctx, machineDAOSpan := msd.tracerSpan.CreateChildInCurrentContext(ctx, "MachineDAO.GetCount")
	if machineDAOSpan != nil {
		defer machineDAOSpan.End()
	}

	query := db.GetIDB(tx, msd.dbSession).NewSelect().Model((*Machine)(nil))
	query, err = msd.setQueryWithFilter(filter, query, machineDAOSpan)
	if err != nil {
		return 0, err
	}

	return query.Count(ctx)
}

// NewMachineDAO returns a new MachineDAO
func NewMachineDAO(dbSession *db.Session) MachineDAO {
	return &MachineSQLDAO{
		dbSession:  dbSession,
		tracerSpan: stracer.NewTracerSpan(),
	}
}
