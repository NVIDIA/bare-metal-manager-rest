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

package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/nvidia/bare-metal-manager-rest/api/pkg/api/handler/util/common"
	"github.com/nvidia/bare-metal-manager-rest/api/pkg/api/model"
	"github.com/nvidia/bare-metal-manager-rest/api/pkg/api/pagination"
	sc "github.com/nvidia/bare-metal-manager-rest/api/pkg/client/site"
	"github.com/nvidia/bare-metal-manager-rest/common/pkg/otelecho"
	cdb "github.com/nvidia/bare-metal-manager-rest/db/pkg/db"
	cdbm "github.com/nvidia/bare-metal-manager-rest/db/pkg/db/model"
	cdbu "github.com/nvidia/bare-metal-manager-rest/db/pkg/util"
	rlav1 "github.com/nvidia/bare-metal-manager-rest/workflow-schema/rla/protobuf/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/uptrace/bun/extra/bundebug"
	oteltrace "go.opentelemetry.io/otel/trace"
	tmocks "go.temporal.io/sdk/mocks"
)

func testTrayInitDB(t *testing.T) *cdb.Session {
	dbSession := cdbu.GetTestDBSession(t, false)
	dbSession.DB.AddQueryHook(bundebug.NewQueryHook(
		bundebug.WithEnabled(false),
		bundebug.FromEnv("BUNDEBUG"),
	))

	ctx := context.Background()

	// Reset required tables in dependency order
	err := dbSession.DB.ResetModel(ctx, (*cdbm.User)(nil))
	assert.Nil(t, err)
	err = dbSession.DB.ResetModel(ctx, (*cdbm.InfrastructureProvider)(nil))
	assert.Nil(t, err)
	err = dbSession.DB.ResetModel(ctx, (*cdbm.Tenant)(nil))
	assert.Nil(t, err)
	err = dbSession.DB.ResetModel(ctx, (*cdbm.Site)(nil))
	assert.Nil(t, err)
	err = dbSession.DB.ResetModel(ctx, (*cdbm.TenantSite)(nil))
	assert.Nil(t, err)
	err = dbSession.DB.ResetModel(ctx, (*cdbm.TenantAccount)(nil))
	assert.Nil(t, err)

	return dbSession
}

func testTraySetupTestData(t *testing.T, dbSession *cdb.Session, org string) (*cdbm.InfrastructureProvider, *cdbm.Site, *cdbm.Tenant) {
	ctx := context.Background()

	// Create infrastructure provider
	ip := &cdbm.InfrastructureProvider{
		ID:   uuid.New(),
		Name: "test-provider",
		Org:  org,
	}
	_, err := dbSession.DB.NewInsert().Model(ip).Exec(ctx)
	assert.Nil(t, err)

	// Create site
	site := &cdbm.Site{
		ID:                       uuid.New(),
		Name:                     "test-site",
		Org:                      org,
		InfrastructureProviderID: ip.ID,
		Status:                   cdbm.SiteStatusRegistered,
	}
	_, err = dbSession.DB.NewInsert().Model(site).Exec(ctx)
	assert.Nil(t, err)

	// Create tenant with TargetedInstanceCreation enabled (privileged tenant)
	tenant := &cdbm.Tenant{
		ID:  uuid.New(),
		Org: org,
		Config: &cdbm.TenantConfig{
			TargetedInstanceCreation: true,
		},
		CreatedBy: uuid.New(),
	}
	_, err = dbSession.DB.NewInsert().Model(tenant).Exec(ctx)
	assert.Nil(t, err)

	// Create tenant account for privileged tenant
	ta := &cdbm.TenantAccount{
		ID:                       uuid.New(),
		TenantID:                 &tenant.ID,
		InfrastructureProviderID: ip.ID,
	}
	_, err = dbSession.DB.NewInsert().Model(ta).Exec(ctx)
	assert.Nil(t, err)

	return ip, site, tenant
}

func testTrayBuildUser(t *testing.T, dbSession *cdb.Session, starfleetID string, org string, roles []string) *cdbm.User {
	uDAO := cdbm.NewUserDAO(dbSession)

	OrgData := cdbm.OrgData{}
	OrgData[org] = cdbm.Org{
		ID:          123,
		Name:        org,
		DisplayName: org,
		OrgType:     "ENTERPRISE",
		Roles:       roles,
	}
	u, err := uDAO.Create(
		context.Background(),
		nil,
		cdbm.UserCreateInput{
			AuxiliaryID: nil,
			StarfleetID: &starfleetID,
			Email:       cdb.GetStrPtr("test@test.com"),
			FirstName:   cdb.GetStrPtr("Test"),
			LastName:    cdb.GetStrPtr("User"),
			OrgData:     OrgData,
		},
	)
	assert.Nil(t, err)

	return u
}

// createMockComponent creates a mock RLA Component for testing
func createMockComponent(id, name, manufacturer, modelStr, componentID string, compType rlav1.ComponentType, rackID string) *rlav1.Component {
	comp := &rlav1.Component{
		Type:        compType,
		ComponentId: componentID,
		Info: &rlav1.DeviceInfo{
			Id:           &rlav1.UUID{Id: id},
			Name:         name,
			Manufacturer: manufacturer,
			Model:        &modelStr,
		},
	}
	if rackID != "" {
		comp.RackId = &rlav1.UUID{Id: rackID}
	}
	return comp
}

func TestGetTrayHandler_Handle(t *testing.T) {
	// Setup
	e := echo.New()
	dbSession := testTrayInitDB(t)
	defer dbSession.Close()

	cfg := common.GetTestConfig()
	tcfg, _ := cfg.GetTemporalConfig()
	scp := sc.NewClientPool(tcfg)

	org := "test-org"
	_, site, _ := testTraySetupTestData(t, dbSession, org)

	// Create provider user
	providerUser := testTrayBuildUser(t, dbSession, "provider-user-tray-get", org, []string{"FORGE_PROVIDER_ADMIN"})

	// Create tenant user (no provider role, no site access)
	tenantUser := testTrayBuildUser(t, dbSession, "tenant-user-tray-get", org, []string{"FORGE_TENANT_ADMIN"})

	handler := NewGetTrayHandler(dbSession, nil, scp, cfg)

	trayID := uuid.New().String()

	// Create mock component for success cases
	mockComponent := createMockComponent(
		trayID, "compute-tray-1", "NVIDIA", "GB200", "carbide-machine-001",
		rlav1.ComponentType_COMPONENT_TYPE_COMPUTE, "rack-id-1",
	)

	tracer := oteltrace.NewNoopTracerProvider().Tracer("test")
	ctx := context.Background()

	tests := []struct {
		name           string
		reqOrg         string
		user           *cdbm.User
		trayID         string
		queryParams    map[string]string
		mockComponent  *rlav1.Component
		expectedStatus int
		wantErr        bool
	}{
		{
			name:   "success - get tray by ID",
			reqOrg: org,
			user:   providerUser,
			trayID: trayID,
			queryParams: map[string]string{
				"siteId": site.ID.String(),
			},
			mockComponent:  mockComponent,
			expectedStatus: http.StatusOK,
			wantErr:        false,
		},
		{
			name:   "failure - missing siteId",
			reqOrg: org,
			user:   providerUser,
			trayID: trayID,
			queryParams: map[string]string{
				// no siteId
			},
			expectedStatus: http.StatusBadRequest,
			wantErr:        true,
		},
		{
			name:   "failure - invalid siteId",
			reqOrg: org,
			user:   providerUser,
			trayID: trayID,
			queryParams: map[string]string{
				"siteId": uuid.New().String(), // non-existent site
			},
			expectedStatus: http.StatusBadRequest,
			wantErr:        true,
		},
		{
			name:   "failure - tenant access denied (no site access)",
			reqOrg: org,
			user:   tenantUser,
			trayID: trayID,
			queryParams: map[string]string{
				"siteId": site.ID.String(),
			},
			expectedStatus: http.StatusForbidden,
			wantErr:        true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock Temporal client
			mockTemporalClient := &tmocks.Client{}
			mockWorkflowRun := &tmocks.WorkflowRun{}
			mockWorkflowRun.On("GetID").Return("test-workflow-id")
			if tt.mockComponent != nil {
				mockWorkflowRun.Mock.On("Get", mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
					resp := args.Get(1).(*rlav1.GetComponentInfoResponse)
					resp.Component = tt.mockComponent
				}).Return(nil)
			} else {
				mockWorkflowRun.Mock.On("Get", mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
					resp := args.Get(1).(*rlav1.GetComponentInfoResponse)
					resp.Component = nil
				}).Return(nil)
			}
			mockTemporalClient.Mock.On("ExecuteWorkflow", mock.Anything, mock.Anything, "GetTray", mock.Anything).Return(mockWorkflowRun, nil)
			scp.IDClientMap[site.ID.String()] = mockTemporalClient

			// Build query string
			q := url.Values{}
			for k, v := range tt.queryParams {
				q.Set(k, v)
			}
			path := fmt.Sprintf("/v2/org/%s/carbide/tray/%s?%s", tt.reqOrg, tt.trayID, q.Encode())

			req := httptest.NewRequest(http.MethodGet, path, nil)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()

			ec := e.NewContext(req, rec)
			ec.SetParamNames("orgName", "id")
			ec.SetParamValues(tt.reqOrg, tt.trayID)
			ec.Set("user", tt.user)

			ctx = context.WithValue(ctx, otelecho.TracerKey, tracer)
			ec.SetRequest(ec.Request().WithContext(ctx))

			err := handler.Handle(ec)
			if tt.expectedStatus != rec.Code {
				t.Errorf("GetTrayHandler.Handle() status = %v, want %v, response: %v, err: %v", rec.Code, tt.expectedStatus, rec.Body.String(), err)
			}

			require.Equal(t, tt.expectedStatus, rec.Code)
			if tt.expectedStatus != http.StatusOK {
				return
			}

			// Verify response
			var apiTray model.APITray
			err = json.Unmarshal(rec.Body.Bytes(), &apiTray)
			assert.NoError(t, err)
			assert.Equal(t, trayID, apiTray.ID)
			assert.Equal(t, "compute", apiTray.Type)
			assert.Equal(t, "NVIDIA", apiTray.Manufacturer)
		})
	}
}

func TestGetAllTrayHandler_Handle(t *testing.T) {
	// Setup
	e := echo.New()
	dbSession := testTrayInitDB(t)
	defer dbSession.Close()

	cfg := common.GetTestConfig()
	tcfg, _ := cfg.GetTemporalConfig()
	scp := sc.NewClientPool(tcfg)

	org := "test-org"
	_, site, _ := testTraySetupTestData(t, dbSession, org)

	// Create provider user
	providerUser := testTrayBuildUser(t, dbSession, "provider-user-tray", org, []string{"FORGE_PROVIDER_ADMIN"})

	// Create tenant user (no provider role, no site access)
	tenantUser := testTrayBuildUser(t, dbSession, "tenant-user-tray", org, []string{"FORGE_TENANT_ADMIN"})

	handler := NewGetAllTrayHandler(dbSession, nil, scp, cfg)

	rackID := uuid.New().String()

	// Helper to create mock RLA response
	createMockRLAResponse := func(components []*rlav1.Component, total int32) *rlav1.GetComponentsResponse {
		return &rlav1.GetComponentsResponse{
			Components: components,
			Total:      total,
		}
	}

	// Create test components (trays)
	testComponents := []*rlav1.Component{
		createMockComponent("tray-1", "Compute-001", "NVIDIA", "GB200", "comp-1", rlav1.ComponentType_COMPONENT_TYPE_COMPUTE, rackID),
		createMockComponent("tray-2", "Compute-002", "NVIDIA", "GB200", "comp-2", rlav1.ComponentType_COMPONENT_TYPE_COMPUTE, rackID),
		createMockComponent("tray-3", "Switch-001", "NVIDIA", "NVL-Switch", "comp-3", rlav1.ComponentType_COMPONENT_TYPE_NVLSWITCH, rackID),
		createMockComponent("tray-4", "Power-001", "NVIDIA", "PowerShelf", "comp-4", rlav1.ComponentType_COMPONENT_TYPE_POWERSHELF, rackID),
		createMockComponent("tray-5", "ToRSwitch-001", "Dell", "S5248", "comp-5", rlav1.ComponentType_COMPONENT_TYPE_TORSWITCH, rackID),
	}

	tracer := oteltrace.NewNoopTracerProvider().Tracer("test")
	ctx := context.Background()

	tests := []struct {
		name           string
		reqOrg         string
		user           *cdbm.User
		queryParams    map[string]string
		mockResponse   *rlav1.GetComponentsResponse
		expectedStatus int
		expectedCount  int
		expectedTotal  *int
		wantErr        bool
	}{
		{
			name:   "success - get all trays",
			reqOrg: org,
			user:   providerUser,
			queryParams: map[string]string{
				"siteId": site.ID.String(),
			},
			mockResponse:   createMockRLAResponse(testComponents, int32(len(testComponents))),
			expectedStatus: http.StatusOK,
			expectedCount:  len(testComponents),
			expectedTotal:  cdb.GetIntPtr(len(testComponents)),
			wantErr:        false,
		},
		{
			name:   "success - filter by rackId",
			reqOrg: org,
			user:   providerUser,
			queryParams: map[string]string{
				"siteId": site.ID.String(),
				"rackId": rackID,
			},
			mockResponse:   createMockRLAResponse(testComponents, int32(len(testComponents))),
			expectedStatus: http.StatusOK,
			expectedCount:  len(testComponents),
			expectedTotal:  cdb.GetIntPtr(len(testComponents)),
			wantErr:        false,
		},
		{
			name:   "success - filter by type compute",
			reqOrg: org,
			user:   providerUser,
			queryParams: map[string]string{
				"siteId": site.ID.String(),
				"type":   "compute",
			},
			mockResponse:   createMockRLAResponse(testComponents[:2], 2),
			expectedStatus: http.StatusOK,
			expectedCount:  2,
			expectedTotal:  cdb.GetIntPtr(2),
			wantErr:        false,
		},
		{
			name:   "success - filter by rackName",
			reqOrg: org,
			user:   providerUser,
			queryParams: map[string]string{
				"siteId":   site.ID.String(),
				"rackName": "Rack-001",
			},
			mockResponse:   createMockRLAResponse(testComponents, int32(len(testComponents))),
			expectedStatus: http.StatusOK,
			expectedCount:  len(testComponents),
			expectedTotal:  cdb.GetIntPtr(len(testComponents)),
			wantErr:        false,
		},
		{
			name:   "success - pagination",
			reqOrg: org,
			user:   providerUser,
			queryParams: map[string]string{
				"siteId":     site.ID.String(),
				"pageNumber": "1",
				"pageSize":   "2",
			},
			mockResponse:   createMockRLAResponse(testComponents[:2], int32(len(testComponents))),
			expectedStatus: http.StatusOK,
			expectedCount:  2,
			expectedTotal:  cdb.GetIntPtr(len(testComponents)),
			wantErr:        false,
		},
		{
			name:   "success - orderBy name ASC",
			reqOrg: org,
			user:   providerUser,
			queryParams: map[string]string{
				"siteId":  site.ID.String(),
				"orderBy": "NAME_ASC",
			},
			mockResponse:   createMockRLAResponse(testComponents, int32(len(testComponents))),
			expectedStatus: http.StatusOK,
			expectedCount:  len(testComponents),
			expectedTotal:  cdb.GetIntPtr(len(testComponents)),
			wantErr:        false,
		},
		{
			name:   "success - orderBy manufacturer DESC",
			reqOrg: org,
			user:   providerUser,
			queryParams: map[string]string{
				"siteId":  site.ID.String(),
				"orderBy": "MANUFACTURER_DESC",
			},
			mockResponse:   createMockRLAResponse(testComponents, int32(len(testComponents))),
			expectedStatus: http.StatusOK,
			expectedCount:  len(testComponents),
			expectedTotal:  cdb.GetIntPtr(len(testComponents)),
			wantErr:        false,
		},
		{
			name:   "failure - missing siteId",
			reqOrg: org,
			user:   providerUser,
			queryParams: map[string]string{
				// no siteId
			},
			mockResponse:   nil,
			expectedStatus: http.StatusBadRequest,
			wantErr:        true,
		},
		{
			name:   "failure - tenant access denied",
			reqOrg: org,
			user:   tenantUser,
			queryParams: map[string]string{
				"siteId": site.ID.String(),
			},
			mockResponse:   nil,
			expectedStatus: http.StatusForbidden,
			wantErr:        true,
		},
		{
			name:   "failure - invalid orderBy",
			reqOrg: org,
			user:   providerUser,
			queryParams: map[string]string{
				"siteId":  site.ID.String(),
				"orderBy": "INVALID_FIELD_ASC",
			},
			expectedStatus: http.StatusBadRequest,
			wantErr:        true,
		},
		{
			name:   "failure - invalid pagination (negative pageNumber)",
			reqOrg: org,
			user:   providerUser,
			queryParams: map[string]string{
				"siteId":     site.ID.String(),
				"pageNumber": "-1",
			},
			mockResponse:   nil,
			expectedStatus: http.StatusBadRequest,
			wantErr:        true,
		},
		{
			name:   "failure - invalid rackId (not UUID)",
			reqOrg: org,
			user:   providerUser,
			queryParams: map[string]string{
				"siteId": site.ID.String(),
				"rackId": "not-a-uuid",
			},
			mockResponse:   nil,
			expectedStatus: http.StatusBadRequest,
			wantErr:        true,
		},
		{
			name:   "failure - invalid type",
			reqOrg: org,
			user:   providerUser,
			queryParams: map[string]string{
				"siteId": site.ID.String(),
				"type":   "invalid-type",
			},
			mockResponse:   nil,
			expectedStatus: http.StatusBadRequest,
			wantErr:        true,
		},
		{
			name:   "failure - componentId without type",
			reqOrg: org,
			user:   providerUser,
			queryParams: map[string]string{
				"siteId":      site.ID.String(),
				"componentId": "comp-1",
			},
			mockResponse:   nil,
			expectedStatus: http.StatusBadRequest,
			wantErr:        true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock Temporal client
			mockTemporalClient := &tmocks.Client{}
			mockWorkflowRun := &tmocks.WorkflowRun{}
			mockWorkflowRun.On("GetID").Return("test-workflow-id")
			// Always set up Get mock, even for error cases, as handler may still call it
			if tt.mockResponse != nil {
				mockWorkflowRun.Mock.On("Get", mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
					resp := args.Get(1).(*rlav1.GetComponentsResponse)
					resp.Components = tt.mockResponse.Components
					resp.Total = tt.mockResponse.Total
				}).Return(nil)
			} else {
				// For error cases, set up a mock that returns empty response
				mockWorkflowRun.Mock.On("Get", mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
					resp := args.Get(1).(*rlav1.GetComponentsResponse)
					resp.Components = []*rlav1.Component{}
					resp.Total = 0
				}).Return(nil)
			}
			mockTemporalClient.Mock.On("ExecuteWorkflow", mock.Anything, mock.Anything, "GetTrays", mock.Anything, mock.Anything).Return(mockWorkflowRun, nil)
			scp.IDClientMap[site.ID.String()] = mockTemporalClient

			// Build query string
			q := url.Values{}
			for k, v := range tt.queryParams {
				q.Set(k, v)
			}
			path := fmt.Sprintf("/v2/org/%s/carbide/tray?%s", tt.reqOrg, q.Encode())

			req := httptest.NewRequest(http.MethodGet, path, nil)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()

			ec := e.NewContext(req, rec)
			ec.SetParamNames("orgName")
			ec.SetParamValues(tt.reqOrg)
			ec.Set("user", tt.user)

			ctx = context.WithValue(ctx, otelecho.TracerKey, tracer)
			ec.SetRequest(ec.Request().WithContext(ctx))

			err := handler.Handle(ec)
			// In Echo, c.JSON() returns nil on success, so err can be nil even when returning error response
			// Check status code instead of err for error cases
			if tt.expectedStatus != rec.Code {
				t.Errorf("GetAllTrayHandler.Handle() status = %v, want %v, response: %v, err: %v", rec.Code, tt.expectedStatus, rec.Body.String(), err)
			}

			require.Equal(t, tt.expectedStatus, rec.Code)
			if tt.expectedStatus != http.StatusOK {
				return
			}

			// Verify response
			var apiTrays []*model.APITray
			err = json.Unmarshal(rec.Body.Bytes(), &apiTrays)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedCount, len(apiTrays))

			// Verify pagination header
			ph := rec.Header().Get(pagination.ResponseHeaderName)
			assert.NotEmpty(t, ph)

			pr := &pagination.PageResponse{}
			err = json.Unmarshal([]byte(ph), pr)
			assert.NoError(t, err)

			if tt.expectedTotal != nil {
				assert.Equal(t, *tt.expectedTotal, pr.Total)
			}
		})
	}
}

func TestBuildTrayFilterInput(t *testing.T) {
	e := echo.New()

	validRackID := uuid.New().String()
	validID1 := uuid.New().String()
	validID2 := uuid.New().String()

	tests := []struct {
		name        string
		queryParams map[string]string
		wantErr     bool
		validate    func(t *testing.T, filter *model.TrayFilterInput)
	}{
		{
			name:        "empty params - no error",
			queryParams: map[string]string{},
			wantErr:     false,
			validate: func(t *testing.T, filter *model.TrayFilterInput) {
				assert.Nil(t, filter.RackID)
				assert.Nil(t, filter.RackName)
				assert.Nil(t, filter.Type)
				assert.Empty(t, filter.ComponentIDs)
				assert.Empty(t, filter.IDs)
			},
		},
		{
			name: "valid rackId",
			queryParams: map[string]string{
				"rackId": validRackID,
			},
			wantErr: false,
			validate: func(t *testing.T, filter *model.TrayFilterInput) {
				require.NotNil(t, filter.RackID)
				assert.Equal(t, validRackID, *filter.RackID)
			},
		},
		{
			name: "invalid rackId",
			queryParams: map[string]string{
				"rackId": "not-a-uuid",
			},
			wantErr: true,
		},
		{
			name: "valid rackName",
			queryParams: map[string]string{
				"rackName": "Rack-001",
			},
			wantErr: false,
			validate: func(t *testing.T, filter *model.TrayFilterInput) {
				require.NotNil(t, filter.RackName)
				assert.Equal(t, "Rack-001", *filter.RackName)
			},
		},
		{
			name: "rackname lowercase alias",
			queryParams: map[string]string{
				"rackname": "Rack-002",
			},
			wantErr: false,
			validate: func(t *testing.T, filter *model.TrayFilterInput) {
				require.NotNil(t, filter.RackName)
				assert.Equal(t, "Rack-002", *filter.RackName)
			},
		},
		{
			name: "valid type - compute",
			queryParams: map[string]string{
				"type": "compute",
			},
			wantErr: false,
			validate: func(t *testing.T, filter *model.TrayFilterInput) {
				require.NotNil(t, filter.Type)
				assert.Equal(t, "compute", *filter.Type)
			},
		},
		{
			name: "valid type - switch",
			queryParams: map[string]string{
				"type": "switch",
			},
			wantErr: false,
			validate: func(t *testing.T, filter *model.TrayFilterInput) {
				require.NotNil(t, filter.Type)
				assert.Equal(t, "switch", *filter.Type)
			},
		},
		{
			name: "invalid type",
			queryParams: map[string]string{
				"type": "invalid-type",
			},
			wantErr: true,
		},
		{
			name: "valid componentId with type",
			queryParams: map[string]string{
				"componentId": "comp-1,comp-2",
				"type":        "compute",
			},
			wantErr: false,
			validate: func(t *testing.T, filter *model.TrayFilterInput) {
				assert.Len(t, filter.ComponentIDs, 2)
				assert.Contains(t, filter.ComponentIDs, "comp-1")
				assert.Contains(t, filter.ComponentIDs, "comp-2")
			},
		},
		{
			name: "componentId without type - error",
			queryParams: map[string]string{
				"componentId": "comp-1",
			},
			wantErr: true,
		},
		{
			name: "valid UUID ids",
			queryParams: map[string]string{
				"id": fmt.Sprintf("%s,%s", validID1, validID2),
			},
			wantErr: false,
			validate: func(t *testing.T, filter *model.TrayFilterInput) {
				assert.Len(t, filter.IDs, 2)
			},
		},
		{
			name: "invalid UUID id",
			queryParams: map[string]string{
				"id": "not-a-uuid",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q := url.Values{}
			for k, v := range tt.queryParams {
				q.Set(k, v)
			}
			path := fmt.Sprintf("/test?%s", q.Encode())

			req := httptest.NewRequest(http.MethodGet, path, nil)
			rec := httptest.NewRecorder()
			ec := e.NewContext(req, rec)

			filter, err := buildTrayFilterInput(ec)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			require.NotNil(t, filter)
			if tt.validate != nil {
				tt.validate(t, filter)
			}
		})
	}
}

func TestBuildRLARequestFromFilter(t *testing.T) {
	rackID := uuid.New().String()
	rackName := "Rack-001"
	trayType := "compute"
	id1 := uuid.New().String()
	id2 := uuid.New().String()

	tests := []struct {
		name     string
		filter   *model.TrayFilterInput
		validate func(t *testing.T, req *rlav1.GetComponentsRequest)
	}{
		{
			name:   "empty filter - no target spec",
			filter: &model.TrayFilterInput{},
			validate: func(t *testing.T, req *rlav1.GetComponentsRequest) {
				assert.Nil(t, req.TargetSpec)
			},
		},
		{
			name: "rackId only - rack-level targeting",
			filter: &model.TrayFilterInput{
				RackID: &rackID,
			},
			validate: func(t *testing.T, req *rlav1.GetComponentsRequest) {
				require.NotNil(t, req.TargetSpec)
				rackTargets := req.TargetSpec.GetRacks()
				require.NotNil(t, rackTargets)
				require.Len(t, rackTargets.Targets, 1)
				assert.Equal(t, rackID, rackTargets.Targets[0].GetId().GetId())
			},
		},
		{
			name: "rackName only - rack-level targeting",
			filter: &model.TrayFilterInput{
				RackName: &rackName,
			},
			validate: func(t *testing.T, req *rlav1.GetComponentsRequest) {
				require.NotNil(t, req.TargetSpec)
				rackTargets := req.TargetSpec.GetRacks()
				require.NotNil(t, rackTargets)
				require.Len(t, rackTargets.Targets, 1)
				assert.Equal(t, rackName, rackTargets.Targets[0].GetName())
			},
		},
		{
			name: "type only - rack-level targeting with component type",
			filter: &model.TrayFilterInput{
				Type: &trayType,
			},
			validate: func(t *testing.T, req *rlav1.GetComponentsRequest) {
				require.NotNil(t, req.TargetSpec)
				rackTargets := req.TargetSpec.GetRacks()
				require.NotNil(t, rackTargets)
				require.Len(t, rackTargets.Targets, 1)
				assert.Contains(t, rackTargets.Targets[0].ComponentTypes, rlav1.ComponentType_COMPONENT_TYPE_COMPUTE)
			},
		},
		{
			name: "rackId with type - rack-level targeting with filter",
			filter: &model.TrayFilterInput{
				RackID: &rackID,
				Type:   &trayType,
			},
			validate: func(t *testing.T, req *rlav1.GetComponentsRequest) {
				require.NotNil(t, req.TargetSpec)
				rackTargets := req.TargetSpec.GetRacks()
				require.NotNil(t, rackTargets)
				require.Len(t, rackTargets.Targets, 1)
				assert.Equal(t, rackID, rackTargets.Targets[0].GetId().GetId())
				assert.Contains(t, rackTargets.Targets[0].ComponentTypes, rlav1.ComponentType_COMPONENT_TYPE_COMPUTE)
			},
		},
		{
			name: "IDs - component-level targeting",
			filter: &model.TrayFilterInput{
				IDs: []string{id1, id2},
			},
			validate: func(t *testing.T, req *rlav1.GetComponentsRequest) {
				require.NotNil(t, req.TargetSpec)
				compTargets := req.TargetSpec.GetComponents()
				require.NotNil(t, compTargets)
				assert.Len(t, compTargets.Targets, 2)
				assert.Equal(t, id1, compTargets.Targets[0].GetId().GetId())
				assert.Equal(t, id2, compTargets.Targets[1].GetId().GetId())
			},
		},
		{
			name: "componentIDs with type - component-level targeting via ExternalRef",
			filter: &model.TrayFilterInput{
				ComponentIDs: []string{"comp-1", "comp-2"},
				Type:         &trayType,
			},
			validate: func(t *testing.T, req *rlav1.GetComponentsRequest) {
				require.NotNil(t, req.TargetSpec)
				compTargets := req.TargetSpec.GetComponents()
				require.NotNil(t, compTargets)
				assert.Len(t, compTargets.Targets, 2)
				// Verify they are ExternalRef targets
				for _, target := range compTargets.Targets {
					ext := target.GetExternal()
					require.NotNil(t, ext)
					assert.Equal(t, rlav1.ComponentType_COMPONENT_TYPE_COMPUTE, ext.Type)
				}
			},
		},
		{
			name: "IDs take priority over rackId (component-level targeting)",
			filter: &model.TrayFilterInput{
				IDs:    []string{id1},
				RackID: &rackID,
			},
			validate: func(t *testing.T, req *rlav1.GetComponentsRequest) {
				require.NotNil(t, req.TargetSpec)
				compTargets := req.TargetSpec.GetComponents()
				require.NotNil(t, compTargets, "IDs should produce component-level targeting, not rack-level")
				assert.Len(t, compTargets.Targets, 1)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := buildRLARequestFromFilter(tt.filter)
			require.NotNil(t, req)
			if tt.validate != nil {
				tt.validate(t, req)
			}
		})
	}
}
