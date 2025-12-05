// SPDX-FileCopyrightText: Copyright (c) 2021-2023 NVIDIA CORPORATION & AFFILIATES. All rights reserved.
// SPDX-License-Identifier: LicenseRef-NvidiaProprietary
//
// NVIDIA CORPORATION, its affiliates and licensors retain all intellectual
// property and proprietary rights in and to this material, related
// documentation and any modifications thereto. Any use, reproduction,
// disclosure or distribution of this material and related documentation
// without an express license agreement from NVIDIA CORPORATION or
// its affiliates is strictly prohibited.

package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/nvidia/carbide-rest/api/internal/config"
	"github.com/nvidia/carbide-rest/api/pkg/api/handler/util/common"
	"github.com/nvidia/carbide-rest/api/pkg/api/model"
	"github.com/nvidia/carbide-rest/api/pkg/api/pagination"
	"github.com/nvidia/carbide-rest/common/pkg/otelecho"
	sutil "github.com/nvidia/carbide-rest/common/pkg/util"
	cdb "github.com/nvidia/carbide-rest/db/pkg/db"
	cdbm "github.com/nvidia/carbide-rest/db/pkg/db/model"
	oteltrace "go.opentelemetry.io/otel/trace"
	temporalClient "go.temporal.io/sdk/client"
	tmocks "go.temporal.io/sdk/mocks"
)

func TestGetAllNVLinkInterfaceHandler_Handle(t *testing.T) {
	ctx := context.Background()
	type fields struct {
		dbSession *cdb.Session
		tc        temporalClient.Client
		cfg       *config.Config
	}
	type args struct {
		reqInstance *cdbm.Instance

		reqInstanceID string
		reqOrg        string
		reqUser       *cdbm.User
		reqMachine    *cdbm.Machine
		respCode      int
	}

	dbSession := testInstanceInitDB(t)
	defer dbSession.Close()

	testInstanceSetupSchema(t, dbSession)

	ipOrg := "test-provider-org"
	ipOrgRoles := []string{"FORGE_PROVIDER_ADMIN"}

	tnOrg1 := "test-tenant-org-1"
	tnOrgRoles1 := []string{"FORGE_TENANT_ADMIN"}

	tnOrg2 := "test-tenant-org-2"
	tnOrgRoles2 := []string{"FORGE_TENANT_ADMIN"}

	ipu := testInstanceBuildUser(t, dbSession, "test-starfleet-id-1", ipOrg, ipOrgRoles)
	ip := testInstanceSiteBuildInfrastructureProvider(t, dbSession, "test-infrastructure-provider", ipOrg, ipu)

	st1 := testInstanceBuildSite(t, dbSession, ip, "test-site-1", cdbm.SiteStatusRegistered, true, ipu)
	assert.NotNil(t, st1)

	tnu1 := testInstanceBuildUser(t, dbSession, "test-starfleet-id-2", tnOrg1, tnOrgRoles1)
	tn1 := testInstanceBuildTenant(t, dbSession, "test-tenant", tnOrg1, tnu1)

	tnu2 := testInstanceBuildUser(t, dbSession, "test-starfleet-id-3", tnOrg2, tnOrgRoles2)

	al1 := testInstanceSiteBuildAllocation(t, dbSession, st1, tn1, "test-allocation-1", ipu)
	assert.NotNil(t, al1)

	ist1 := testInstanceBuildInstanceType(t, dbSession, ip, "test-instance-type-1", st1, cdbm.InstanceStatusReady)
	assert.NotNil(t, ist1)

	alc1 := testInstanceSiteBuildAllocationContraints(t, dbSession, al1, cdbm.AllocationResourceTypeInstanceType, ist1.ID, cdbm.AllocationConstraintTypeReserved, 5, ipu)
	assert.NotNil(t, alc1)

	mc1 := testInstanceBuildMachine(t, dbSession, ip.ID, st1.ID, cdb.GetBoolPtr(false), nil)
	assert.NotNil(t, mc1)

	mcinst1 := testInstanceBuildMachineInstanceType(t, dbSession, mc1, ist1)
	assert.NotNil(t, mcinst1)

	os1 := testInstanceBuildOperatingSystem(t, dbSession, "test-operating-system-1", tn1, cdbm.OperatingSystemTypeImage, false, nil, false, cdbm.OperatingSystemStatusReady, tnu1)
	assert.NotNil(t, os1)

	vpc1 := testInstanceBuildVPC(t, dbSession, "test-vpc-1", ip, tn1, st1, cdb.GetUUIDPtr(uuid.New()), nil, cdb.GetStrPtr(cdbm.VpcEthernetVirtualizer), nil, cdbm.VpcStatusReady, tnu1)
	assert.NotNil(t, vpc1)

	vpc2 := testInstanceBuildVPC(t, dbSession, "test-vpc-2", ip, tn1, st1, nil, nil, cdb.GetStrPtr(cdbm.VpcEthernetVirtualizer), nil, cdbm.VpcStatusPending, tnu1)
	assert.NotNil(t, vpc2)

	inst1 := testInstanceBuildInstance(t, dbSession, "test-instance-2", al1.ID, alc1.ID, tn1.ID, ip.ID, st1.ID, &ist1.ID, vpc1.ID, cdb.GetStrPtr(mc1.ID), &os1.ID, nil, cdbm.InstanceStatusReady)
	assert.NotNil(t, inst1)

	nvlinklogicalpartitions := []*cdbm.NVLinkLogicalPartition{}
	for i := 0; i < 3; i++ {
		nvlinklogicalpartition1 := testBuildNVLinkLogicalPartition(t, dbSession, fmt.Sprintf("test-nvlinklogicalpartition-%d", i), tn1.Org, st1, tn1, cdb.GetStrPtr(cdbm.NVLinkLogicalPartitionStatusReady), false)
		assert.NotNil(t, nvlinklogicalpartition1)
		nvlinklogicalpartitions = append(nvlinklogicalpartitions, nvlinklogicalpartition1)
	}

	nvlifcs := []*cdbm.NVLinkInterface{}
	for i := 0; i < 25; i++ {
		nvlinklogicalpartition := nvlinklogicalpartitions[i%3]
		nvlifc := testInstanceBuildInstanceNVLinkInterface(t, dbSession, st1.ID, inst1.ID, nvlinklogicalpartition.ID, cdb.GetStrPtr("NVIDIA GB200"), i%4, cdbm.NVLinkInterfaceStatusProvisioning)
		assert.NotNil(t, nvlifc)
		nvlifcs = append(nvlifcs, nvlifc)
	}

	e := echo.New()
	cfg := common.GetTestConfig()
	tc := &tmocks.Client{}

	// OTEL Spanner configuration
	tracer, _, ctx := common.TestCommonTraceProviderSetup(t, ctx)

	tests := []struct {
		name                             string
		fields                           fields
		args                             args
		wantErr                          bool
		queryStatus                      *string
		queryIncludeRelations1           *string
		queryIncludeRelations2           *string
		pageNumber                       *int
		pageSize                         *int
		orderBy                          *string
		expectedNVLinkLogicalPartitionID *uuid.UUID
		expectedDeviceInstance           *int
		expectedCount                    int
		expectedTotal                    int
		verifyChildSpanner               bool
	}{
		{
			name: "test NVLinkInterface getall nvlink interfaces API endpoint success",
			fields: fields{
				dbSession: dbSession,
				tc:        tc,
				cfg:       cfg,
			},
			args: args{
				reqInstance:   inst1,
				reqInstanceID: inst1.ID.String(),
				reqOrg:        tnOrg1,
				reqUser:       tnu1,
				respCode:      http.StatusOK,
			},
			wantErr:                          false,
			orderBy:                          cdb.GetStrPtr("CREATED_ASC"),
			expectedCount:                    20,
			expectedTotal:                    25,
			expectedNVLinkLogicalPartitionID: cdb.GetUUIDPtr(nvlinklogicalpartitions[0].ID),
			expectedDeviceInstance:           cdb.GetIntPtr(nvlifcs[0].DeviceInstance),
			verifyChildSpanner:               true,
		},
		{
			name: "test NVLinkInterface getall success with paging",
			fields: fields{
				dbSession: dbSession,
				tc:        tc,
				cfg:       cfg,
			},
			args: args{
				reqInstance:   inst1,
				reqInstanceID: inst1.ID.String(),
				reqOrg:        tnOrg1,
				reqUser:       tnu1,
				respCode:      http.StatusOK,
			},
			wantErr:                          false,
			pageNumber:                       cdb.GetIntPtr(1),
			pageSize:                         cdb.GetIntPtr(10),
			orderBy:                          cdb.GetStrPtr("CREATED_ASC"),
			expectedCount:                    10,
			expectedTotal:                    25,
			expectedNVLinkLogicalPartitionID: cdb.GetUUIDPtr(nvlinklogicalpartitions[0].ID),
			expectedDeviceInstance:           cdb.GetIntPtr(nvlifcs[0].DeviceInstance),
		},
		{
			name: "test NVLinkInterface getall success with paging on page 2",
			fields: fields{
				dbSession: dbSession,
				tc:        tc,
				cfg:       cfg,
			},
			args: args{
				reqInstance:   inst1,
				reqInstanceID: inst1.ID.String(),
				reqOrg:        tnOrg1,
				reqUser:       tnu1,
				respCode:      http.StatusOK,
			},
			wantErr:                          false,
			pageNumber:                       cdb.GetIntPtr(2),
			pageSize:                         cdb.GetIntPtr(10),
			orderBy:                          cdb.GetStrPtr("CREATED_ASC"),
			expectedCount:                    10,
			expectedTotal:                    25,
			expectedNVLinkLogicalPartitionID: cdb.GetUUIDPtr(nvlinklogicalpartitions[1].ID),
			expectedDeviceInstance:           cdb.GetIntPtr(nvlifcs[10].DeviceInstance),
		},
		{
			name: "test NVLinkInterface getall error with paging bad orderby",
			fields: fields{
				dbSession: dbSession,
				tc:        tc,
				cfg:       cfg,
			},
			args: args{
				reqInstance:   inst1,
				reqInstanceID: inst1.ID.String(),
				reqOrg:        tnOrg1,
				reqUser:       tnu1,
				respCode:      http.StatusBadRequest,
			},
			wantErr:    false,
			pageNumber: cdb.GetIntPtr(2),
			pageSize:   cdb.GetIntPtr(10),
			orderBy:    cdb.GetStrPtr("TEST_ASC"),
		},
		{
			name: "test Instance getall nvlink interfaces API failure, org does not have a Tenant associated",
			fields: fields{
				dbSession: dbSession,
				tc:        tc,
				cfg:       cfg,
			},
			args: args{
				reqInstance: inst1,
				reqOrg:      ipOrg,
				reqUser:     ipu,
				respCode:    http.StatusForbidden,
			},
			wantErr: false,
		},
		{
			name: "test Instance getall nvlink interfaces API failure, invalid Instance ID in request",
			fields: fields{
				dbSession: dbSession,
				tc:        tc,
				cfg:       cfg,
			},
			args: args{
				reqInstance: inst1,
				reqOrg:      tnOrg1,
				reqUser:     tnu1,
				respCode:    http.StatusBadRequest,
			},
			wantErr: false,
		},
		{
			name: "test Instance getall nvlink interfaces API failure, Instance ID in request not found",
			fields: fields{
				dbSession: dbSession,
				tc:        tc,
				cfg:       cfg,
			},
			args: args{
				reqInstance:   inst1,
				reqInstanceID: uuid.New().String(),
				reqOrg:        tnOrg1,
				reqUser:       tnu1,
				respCode:      http.StatusNotFound,
			},
			wantErr: false,
		},
		{
			name: "test Instance getall nvlink interfaces API failure, Instance not belong to current tenant",
			fields: fields{
				dbSession: dbSession,
				tc:        tc,
				cfg:       cfg,
			},
			args: args{
				reqInstance:   inst1,
				reqInstanceID: inst1.ID.String(),
				reqOrg:        tnOrg2,
				reqUser:       tnu2,
				respCode:      http.StatusForbidden,
			},
			wantErr: false,
		},
		{
			name: "test Instance getall nvlink interfaces API endpoint success include relation",
			fields: fields{
				dbSession: dbSession,
				tc:        tc,
				cfg:       cfg,
			},
			args: args{
				reqInstance:   inst1,
				reqInstanceID: inst1.ID.String(),
				reqOrg:        tnOrg1,
				reqUser:       tnu1,
				respCode:      http.StatusOK,
			},
			queryIncludeRelations1:           cdb.GetStrPtr(cdbm.NVLinkLogicalPartitionRelationName),
			queryIncludeRelations2:           cdb.GetStrPtr(cdbm.InstanceRelationName),
			expectedCount:                    20,
			expectedTotal:                    25,
			orderBy:                          cdb.GetStrPtr("CREATED_ASC"),
			expectedNVLinkLogicalPartitionID: cdb.GetUUIDPtr(nvlinklogicalpartitions[0].ID),
			expectedDeviceInstance:           cdb.GetIntPtr(nvlifcs[0].DeviceInstance),
			wantErr:                          false,
		},
		{
			name: "test Instance getall nvlink interfaces NVLinkInterfaceStatusProvisioning status success",
			fields: fields{
				dbSession: dbSession,
				tc:        tc,
				cfg:       cfg,
			},
			args: args{
				reqInstance:   inst1,
				reqInstanceID: inst1.ID.String(),
				reqOrg:        tnOrg1,
				reqUser:       tnu1,
				respCode:      http.StatusOK,
			},
			queryStatus:                      cdb.GetStrPtr(cdbm.NVLinkInterfaceStatusProvisioning),
			expectedCount:                    20,
			expectedTotal:                    25,
			orderBy:                          cdb.GetStrPtr("CREATED_ASC"),
			expectedNVLinkLogicalPartitionID: cdb.GetUUIDPtr(nvlinklogicalpartitions[0].ID),
			expectedDeviceInstance:           cdb.GetIntPtr(nvlifcs[0].DeviceInstance),
			wantErr:                          false,
		},
		{
			name: "test Instance getall nvlink interfaces BadStatus status success",
			fields: fields{
				dbSession: dbSession,
				tc:        tc,
				cfg:       cfg,
			},
			args: args{
				reqInstance:   inst1,
				reqInstanceID: inst1.ID.String(),
				reqOrg:        tnOrg1,
				reqUser:       tnu1,
				respCode:      http.StatusBadRequest,
			},
			queryStatus:   cdb.GetStrPtr("BadStatus"),
			expectedCount: 0,
			expectedTotal: 0,
			wantErr:       false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			csh := GetAllNVLinkInterfaceHandler{
				dbSession: tt.fields.dbSession,
				tc:        tt.fields.tc,
				cfg:       tt.fields.cfg,
			}

			// Setup echo server/context
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()

			q := req.URL.Query()
			if tt.queryIncludeRelations1 != nil {
				q.Add("includeRelation", *tt.queryIncludeRelations1)
			}
			if tt.queryStatus != nil {
				q.Add("status", *tt.queryStatus)
			}
			if tt.queryIncludeRelations2 != nil {
				q.Add("includeRelation", *tt.queryIncludeRelations2)
			}
			if tt.pageNumber != nil {
				q.Set("pageNumber", fmt.Sprintf("%v", *tt.pageNumber))
			}
			if tt.pageSize != nil {
				q.Set("pageSize", fmt.Sprintf("%v", *tt.pageSize))
			}
			if tt.orderBy != nil {
				q.Set("orderBy", *tt.orderBy)
			}
			req.URL.RawQuery = q.Encode()

			ec := e.NewContext(req, rec)
			ec.SetPath(fmt.Sprintf("/v2/org/%v/carbide/instance/%v/nvlink-interface", tt.args.reqOrg, tt.args.reqInstanceID))
			ec.SetParamNames("orgName", "instanceId")
			ec.SetParamValues(tt.args.reqOrg, tt.args.reqInstanceID)
			ec.Set("user", tt.args.reqUser)

			ctx = context.WithValue(ctx, otelecho.TracerKey, tracer)
			ec.SetRequest(ec.Request().WithContext(ctx))

			if err := csh.Handle(ec); (err != nil) != tt.wantErr {
				t.Errorf("GetAllNVLinkInterfaceHandler.Handle() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.args.respCode != rec.Code {
				t.Errorf("GetAllNVLinkInterfaceHandler.Handle() resp = %v", rec.Body.String())
			}

			require.Equal(t, tt.args.respCode, rec.Code)
			if tt.args.respCode != http.StatusOK {
				return
			}

			rst := []model.APINVLinkInterface{}
			serr := json.Unmarshal(rec.Body.Bytes(), &rst)
			if serr != nil {
				t.Fatal(serr)
			}

			assert.Equal(t, tt.expectedCount, len(rst))
			if tt.expectedNVLinkLogicalPartitionID != nil && tt.expectedNVLinkLogicalPartitionID.String() != "" && len(rst) > 0 {
				assert.Equal(t, tt.expectedNVLinkLogicalPartitionID.String(), rst[0].NVLinkLogicalPartitionID)
			}

			if tt.expectedDeviceInstance != nil && len(rst) > 0 {
				assert.Equal(t, *tt.expectedDeviceInstance, rst[0].DeviceInstance)
			}

			ph := rec.Header().Get(pagination.ResponseHeaderName)
			assert.NotEmpty(t, ph)

			pr := &pagination.PageResponse{}
			err := json.Unmarshal([]byte(ph), pr)
			assert.NoError(t, err)

			assert.Equal(t, tt.expectedTotal, pr.Total)

			if len(rst) > 0 {
				assert.Equal(t, rst[0].InstanceID, tt.args.reqInstance.ID.String())
			}

			if tt.queryIncludeRelations1 != nil || tt.queryIncludeRelations2 != nil {
				if tt.expectedNVLinkLogicalPartitionID != nil && tt.expectedNVLinkLogicalPartitionID.String() != "" {
					assert.Equal(t, tt.expectedNVLinkLogicalPartitionID.String(), rst[0].NVLinkLogicalPartitionID)
				}
			} else {
				if len(rst) > 0 {
					assert.Nil(t, rst[0].Instance)
					if tt.expectedNVLinkLogicalPartitionID != nil && tt.expectedNVLinkLogicalPartitionID.String() != "" {
						assert.Nil(t, rst[0].NVLinkLogicalPartition)
					}
				}
			}

			if tt.verifyChildSpanner {
				span := oteltrace.SpanFromContext(ec.Request().Context())
				assert.True(t, span.SpanContext().IsValid())
			}
		})
	}
}

func TestNewGetAllNVLinkInterfaceHandler(t *testing.T) {
	type args struct {
		dbSession *cdb.Session
		tc        temporalClient.Client
		cfg       *config.Config
	}

	dbSession := testInstanceInitDB(t)
	defer dbSession.Close()
	tc := &tmocks.Client{}
	cfg := common.GetTestConfig()

	tests := []struct {
		name string
		args args
		want GetAllNVLinkInterfaceHandler
	}{
		{
			name: "test GetAllNVLinkInterfaceHandler initialization",
			args: args{
				dbSession: dbSession,
				tc:        tc,
				cfg:       cfg,
			},
			want: GetAllNVLinkInterfaceHandler{
				dbSession:  dbSession,
				tc:         tc,
				cfg:        cfg,
				tracerSpan: sutil.NewTracerSpan(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewGetAllNVLinkInterfaceHandler(tt.args.dbSession, tt.args.tc, tt.args.cfg); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewGetAllNVLinkInterfaceHandler() = %v, want %v", got, tt.want)
			}
		})
	}
}
