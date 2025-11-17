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

package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/NVIDIA/carbide-rest-api/carbide-rest-api/internal/config"
	"github.com/NVIDIA/carbide-rest-api/carbide-rest-api/pkg/api/handler/util/common"
	"github.com/NVIDIA/carbide-rest-api/carbide-rest-api/pkg/api/model"
	sc "github.com/NVIDIA/carbide-rest-api/carbide-rest-api/pkg/client/site"
	cdb "github.com/NVIDIA/carbide-rest-api/carbide-rest-db/pkg/db"
	cdbm "github.com/NVIDIA/carbide-rest-api/carbide-rest-db/pkg/db/model"
	cdbu "github.com/NVIDIA/carbide-rest-api/carbide-rest-db/pkg/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/uptrace/bun/extra/bundebug"
	tmocks "go.temporal.io/sdk/mocks"
)

// testExpectedMachineInitDB initializes a test database session
func testExpectedMachineInitDB(t *testing.T) *cdb.Session {
	dbSession := cdbu.GetTestDBSession(t, true)
	dbSession.DB.AddQueryHook(bundebug.NewQueryHook(
		bundebug.WithEnabled(false),
		bundebug.FromEnv("BUNDEBUG"),
	))

	// Reset required tables in dependency order
	ctx := context.Background()

	// First reset parent tables
	err := dbSession.DB.ResetModel(ctx, (*cdbm.InfrastructureProvider)(nil))
	assert.Nil(t, err)
	err = dbSession.DB.ResetModel(ctx, (*cdbm.Site)(nil))
	assert.Nil(t, err)

	// Then reset child tables that depend on parent tables
	err = dbSession.DB.ResetModel(ctx, (*cdbm.SKU)(nil))
	assert.Nil(t, err)
	err = dbSession.DB.ResetModel(ctx, (*cdbm.ExpectedMachine)(nil))
	assert.Nil(t, err)
	err = dbSession.DB.ResetModel(ctx, (*cdbm.StatusDetail)(nil))
	assert.Nil(t, err)

	return dbSession
}

// testExpectedMachineSetupTestData creates test infrastructure provider and site
func testExpectedMachineSetupTestData(t *testing.T, dbSession *cdb.Session, org string) (*cdbm.InfrastructureProvider, *cdbm.Site) {
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

	// Create multiple test SKUs in the database linked to the site
	sku1ID := "test-sku-uuid-1"
	sku1 := &cdbm.SKU{
		ID:     sku1ID,
		SiteID: site.ID,
	}
	_, err = dbSession.DB.NewInsert().Model(sku1).Exec(ctx)
	assert.Nil(t, err)

	sku2ID := "test-sku-uuid-2"
	sku2 := &cdbm.SKU{
		ID:     sku2ID,
		SiteID: site.ID,
	}
	_, err = dbSession.DB.NewInsert().Model(sku2).Exec(ctx)
	assert.Nil(t, err)

	return ip, site
}

func TestCreateExpectedMachineHandler_Handle(t *testing.T) {
	// Setup
	e := echo.New()

	// Initialize test database
	dbSession := testExpectedMachineInitDB(t)
	defer dbSession.Close()

	cfg := common.GetTestConfig()

	// Prepare client pool for workflow calls
	tcfg, _ := cfg.GetTemporalConfig()
	scp := sc.NewClientPool(tcfg)

	// Create test data first to get the site ID
	org := "test-org"
	infraProv, site := testExpectedMachineSetupTestData(t, dbSession, org)

	// Create an unmanaged site (different infrastructure provider)
	ctx := context.Background()
	unmanagedIP := &cdbm.InfrastructureProvider{
		ID:   uuid.New(),
		Name: "unmanaged-provider",
		Org:  "other-org",
	}
	_, err := dbSession.DB.NewInsert().Model(unmanagedIP).Exec(ctx)
	assert.Nil(t, err)

	unmanagedSite := &cdbm.Site{
		ID:                       uuid.New(),
		Name:                     "unmanaged-site",
		Org:                      "other-org",
		InfrastructureProviderID: unmanagedIP.ID,
		Status:                   cdbm.SiteStatusRegistered,
	}
	_, err = dbSession.DB.NewInsert().Model(unmanagedSite).Exec(ctx)
	assert.Nil(t, err)

	// Add mock temporal client for the site
	mockTemporalClient := &tmocks.Client{}
	mockWorkflowRun := &tmocks.WorkflowRun{}
	mockWorkflowRun.On("GetID").Return("test-workflow-id")
	mockWorkflowRun.Mock.On("Get", mock.Anything, mock.Anything).Return(nil)
	mockTemporalClient.Mock.On("ExecuteWorkflow", mock.Anything, mock.Anything, "CreateExpectedMachine", mock.Anything).Return(mockWorkflowRun, nil)
	scp.IDClientMap[site.ID.String()] = mockTemporalClient

	handler := NewCreateExpectedMachineHandler(dbSession, nil, scp, cfg)

	// Helper function to create mock user
	createMockUser := func(org string) *cdbm.User {
		return &cdbm.User{
			StarfleetID: cdb.GetStrPtr("test-user"),
			OrgData: cdbm.OrgData{
				org: cdbm.Org{
					ID:          123,
					Name:        org,
					DisplayName: org,
					OrgType:     "ENTERPRISE",
					Roles:       []string{"FORGE_PROVIDER_ADMIN"},
				},
			},
		}
	}

	// Test cases
	tests := []struct {
		name           string
		requestBody    model.APIExpectedMachineCreateRequest
		setupContext   func(c echo.Context)
		expectedStatus int
	}{
		{
			name: "successful creation",
			requestBody: model.APIExpectedMachineCreateRequest{
				SiteID:                   site.ID.String(),
				BmcMacAddress:            "00:11:22:33:44:55",
				DefaultBmcUsername:       cdb.GetStrPtr("admin"),
				DefaultBmcPassword:       cdb.GetStrPtr("password"),
				ChassisSerialNumber:      "CHASSIS123",
				FallbackDPUSerialNumbers: []string{"DPU001", "DPU002"},
				Labels:                   map[string]string{"env": "test"},
			},
			setupContext: func(c echo.Context) {
				c.Set("user", createMockUser(org))
				c.SetParamNames("orgName")
				c.SetParamValues(org)
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name: "successful creation with SKU",
			requestBody: model.APIExpectedMachineCreateRequest{
				SiteID:                   site.ID.String(),
				BmcMacAddress:            "00:11:22:33:44:66",
				DefaultBmcUsername:       cdb.GetStrPtr("admin"),
				DefaultBmcPassword:       cdb.GetStrPtr("password"),
				ChassisSerialNumber:      "CHASSIS124",
				FallbackDPUSerialNumbers: []string{"DPU001", "DPU002"},
				SkuId:                    cdb.GetStrPtr("test-sku-uuid-1"),
				Labels:                   map[string]string{"env": "test"},
			},
			setupContext: func(c echo.Context) {
				c.Set("user", createMockUser(org))
				c.SetParamNames("orgName")
				c.SetParamValues(org)
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name: "missing user context",
			requestBody: model.APIExpectedMachineCreateRequest{
				SiteID:              site.ID.String(),
				BmcMacAddress:       "00:11:22:33:44:77",
				DefaultBmcUsername:  cdb.GetStrPtr("admin"),
				DefaultBmcPassword:  cdb.GetStrPtr("password"),
				ChassisSerialNumber: "CHASSIS125",
			},
			setupContext: func(c echo.Context) {
				// Don't set user in context - should cause error
				c.SetParamNames("orgName")
				c.SetParamValues(org)
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "invalid mac address length",
			requestBody: model.APIExpectedMachineCreateRequest{
				SiteID:              site.ID.String(),
				BmcMacAddress:       "00:11:22:33:44", // Too short
				DefaultBmcUsername:  cdb.GetStrPtr("admin"),
				DefaultBmcPassword:  cdb.GetStrPtr("password"),
				ChassisSerialNumber: "CHASSIS126",
			},
			setupContext: func(c echo.Context) {
				c.Set("user", createMockUser(org))
				c.SetParamNames("orgName")
				c.SetParamValues(org)
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "site not found",
			requestBody: model.APIExpectedMachineCreateRequest{
				SiteID:              "12345678-1234-1234-1234-123456789099",
				BmcMacAddress:       "00:11:22:33:44:88",
				DefaultBmcUsername:  cdb.GetStrPtr("admin"),
				DefaultBmcPassword:  cdb.GetStrPtr("password"),
				ChassisSerialNumber: "CHASSIS127",
			},
			setupContext: func(c echo.Context) {
				c.Set("user", createMockUser(org))
				c.SetParamNames("orgName")
				c.SetParamValues(org)
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "cannot create on unmanaged site",
			requestBody: model.APIExpectedMachineCreateRequest{
				SiteID:              unmanagedSite.ID.String(),
				BmcMacAddress:       "00:11:22:33:44:99",
				DefaultBmcUsername:  cdb.GetStrPtr("admin"),
				DefaultBmcPassword:  cdb.GetStrPtr("password"),
				ChassisSerialNumber: "CHASSIS128",
			},
			setupContext: func(c echo.Context) {
				c.Set("user", createMockUser(org))
				c.SetParamNames("orgName")
				c.SetParamValues(org)
			},
			expectedStatus: http.StatusForbidden,
		},
	}

	_ = infraProv // Ensure infraProv is used to avoid compiler warning

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create request
			reqBody, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest(http.MethodPost, "/v2/org/carbide/expected-machine", bytes.NewReader(reqBody))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			req = req.WithContext(context.Background())

			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			// Setup context
			tt.setupContext(c)

			// Execute
			err := handler.Handle(c)

			// Assert
			assert.Nil(t, err)
			assert.Equal(t, tt.expectedStatus, rec.Code)
			if tt.expectedStatus != rec.Code {
				t.Errorf("Response: %v", rec.Body.String())
			}

			// For successful creations, verify labels are returned in response
			if tt.expectedStatus == http.StatusCreated && tt.requestBody.Labels != nil {
				var response model.APIExpectedMachine
				err := json.Unmarshal(rec.Body.Bytes(), &response)
				assert.Nil(t, err)
				assert.NotNil(t, response.Labels, "Labels should not be nil in response")
				assert.Equal(t, tt.requestBody.Labels, response.Labels, "Labels in response should match request")
			}
		})
	}
}

func TestGetAllExpectedMachineHandler_Handle(t *testing.T) {
	// Setup
	e := echo.New()
	dbSession := testExpectedMachineInitDB(t)
	defer dbSession.Close()

	ctx := context.Background()
	cfg := &config.Config{}
	handler := NewGetAllExpectedMachineHandler(dbSession, nil, cfg)

	org := "test-org"
	infraProv, site := testExpectedMachineSetupTestData(t, dbSession, org)

	// Create an unmanaged site
	unmanagedIP := &cdbm.InfrastructureProvider{
		ID:   uuid.New(),
		Name: "unmanaged-provider",
		Org:  "other-org",
	}
	_, err := dbSession.DB.NewInsert().Model(unmanagedIP).Exec(ctx)
	assert.Nil(t, err)

	unmanagedSite := &cdbm.Site{
		ID:                       uuid.New(),
		Name:                     "unmanaged-site",
		Org:                      "other-org",
		InfrastructureProviderID: unmanagedIP.ID,
		Status:                   cdbm.SiteStatusRegistered,
	}
	_, err = dbSession.DB.NewInsert().Model(unmanagedSite).Exec(ctx)
	assert.Nil(t, err)

	// Create expected machines - one on managed site, one on unmanaged site
	emDAO := cdbm.NewExpectedMachineDAO(dbSession)
	managedEM, err := emDAO.Create(ctx, nil, cdbm.ExpectedMachineCreateInput{
		ExpectedMachineID:        uuid.New(),
		SiteID:                   site.ID,
		BmcMacAddress:            "00:11:22:33:44:AA",
		ChassisSerialNumber:      "MANAGED-CHASSIS",
		FallbackDpuSerialNumbers: []string{"DPU001"},
	})
	assert.Nil(t, err)
	assert.NotNil(t, managedEM)

	unmanagedEM, err := emDAO.Create(ctx, nil, cdbm.ExpectedMachineCreateInput{
		ExpectedMachineID:        uuid.New(),
		SiteID:                   unmanagedSite.ID,
		BmcMacAddress:            "00:11:22:33:44:BB",
		ChassisSerialNumber:      "UNMANAGED-CHASSIS",
		FallbackDpuSerialNumbers: []string{"DPU002"},
	})
	assert.Nil(t, err)
	assert.NotNil(t, unmanagedEM)

	// Helper function to create mock user
	createMockUser := func(org string) *cdbm.User {
		return &cdbm.User{
			StarfleetID: cdb.GetStrPtr("test-user"),
			OrgData: cdbm.OrgData{
				org: cdbm.Org{
					ID:          123,
					Name:        org,
					DisplayName: org,
					OrgType:     "ENTERPRISE",
					Roles:       []string{"FORGE_PROVIDER_VIEWER"},
				},
			},
		}
	}

	tests := []struct {
		name                 string
		siteId               string
		includeRelations     []string
		setupContext         func(c echo.Context)
		expectedStatus       int
		checkResponseContent func(t *testing.T, body []byte)
	}{
		{
			name:   "successful GetAll without siteId (lists only managed sites)",
			siteId: "",
			setupContext: func(c echo.Context) {
				c.Set("user", createMockUser(org))
				c.SetParamNames("orgName")
				c.SetParamValues(org)
			},
			expectedStatus: http.StatusOK,
			checkResponseContent: func(t *testing.T, body []byte) {
				var response []model.APIExpectedMachine
				err := json.Unmarshal(body, &response)
				assert.Nil(t, err)
				// Should only return the managed machine
				for _, em := range response {
					assert.NotEqual(t, unmanagedEM.ID, em.ID, "Unmanaged machine should not be in response")
				}
			},
		},
		{
			name:   "successful GetAll with valid siteId",
			siteId: site.ID.String(),
			setupContext: func(c echo.Context) {
				c.Set("user", createMockUser(org))
				c.SetParamNames("orgName")
				c.SetParamValues(org)
			},
			expectedStatus: http.StatusOK,
			checkResponseContent: func(t *testing.T, body []byte) {
				var response []model.APIExpectedMachine
				err := json.Unmarshal(body, &response)
				assert.Nil(t, err)
				// Verify we get results from the specified site only
				for _, em := range response {
					assert.Equal(t, site.ID, em.SiteID, "All results should be from the specified site")
				}
			},
		},
		{
			name:             "successful GetAll with includeRelation=Site",
			siteId:           "",
			includeRelations: []string{"Site"},
			setupContext: func(c echo.Context) {
				c.Set("user", createMockUser(org))
				c.SetParamNames("orgName")
				c.SetParamValues(org)
			},
			expectedStatus: http.StatusOK,
			checkResponseContent: func(t *testing.T, body []byte) {
				var response []model.APIExpectedMachine
				err := json.Unmarshal(body, &response)
				assert.Nil(t, err)
				assert.Greater(t, len(response), 0, "Should return at least one expected machine")
				// Verify Site relation is loaded
				for _, em := range response {
					assert.NotNil(t, em.Site, "Site relation should be loaded")
					assert.Equal(t, em.SiteID.String(), em.Site.ID, "Site ID should match")
				}
			},
		},
		{
			name:             "successful GetAll with includeRelation=Sku",
			siteId:           "",
			includeRelations: []string{"Sku"},
			setupContext: func(c echo.Context) {
				c.Set("user", createMockUser(org))
				c.SetParamNames("orgName")
				c.SetParamValues(org)
			},
			expectedStatus: http.StatusOK,
			checkResponseContent: func(t *testing.T, body []byte) {
				var response []model.APIExpectedMachine
				err := json.Unmarshal(body, &response)
				assert.Nil(t, err)
				// Verify we can include Sku relation without errors
				assert.Greater(t, len(response), 0, "Should return at least one expected machine")
			},
		},
		{
			name:             "successful GetAll with includeRelation=Site,Sku (both relations)",
			siteId:           "",
			includeRelations: []string{"Site", "Sku"},
			setupContext: func(c echo.Context) {
				c.Set("user", createMockUser(org))
				c.SetParamNames("orgName")
				c.SetParamValues(org)
			},
			expectedStatus: http.StatusOK,
			checkResponseContent: func(t *testing.T, body []byte) {
				var response []model.APIExpectedMachine
				err := json.Unmarshal(body, &response)
				assert.Nil(t, err)
				assert.Greater(t, len(response), 0, "Should return at least one expected machine")
				// Verify both Site and Sku relations can be loaded together
				for _, em := range response {
					assert.NotNil(t, em.Site, "Site relation should be loaded")
					assert.Equal(t, em.SiteID.String(), em.Site.ID, "Site ID should match")
					// Sku is optional, so we just verify no error occurred
				}
			},
		},
		{
			name:   "cannot retrieve from unmanaged site",
			siteId: unmanagedSite.ID.String(),
			setupContext: func(c echo.Context) {
				c.Set("user", createMockUser(org))
				c.SetParamNames("orgName")
				c.SetParamValues(org)
			},
			expectedStatus: http.StatusForbidden,
			checkResponseContent: func(t *testing.T, body []byte) {
				// Should return forbidden error
			},
		},
	}

	_ = infraProv // Ensure infraProv is used to avoid compiler warning

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			url := "/v2/org/" + org + "/carbide/expected-machine"
			params := []string{}
			if tt.siteId != "" {
				params = append(params, "siteId="+tt.siteId)
			}
			for _, relation := range tt.includeRelations {
				params = append(params, "includeRelation="+relation)
			}
			if len(params) > 0 {
				url += "?" + params[0]
				for i := 1; i < len(params); i++ {
					url += "&" + params[i]
				}
			}
			req := httptest.NewRequest(http.MethodGet, url, nil)
			req = req.WithContext(context.Background())

			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			// Setup context
			tt.setupContext(c)

			// Execute
			err := handler.Handle(c)

			// Assert
			assert.Nil(t, err)
			assert.Equal(t, tt.expectedStatus, rec.Code)
			if tt.expectedStatus != rec.Code {
				t.Errorf("Response: %v", rec.Body.String())
			}

			// Check response content if provided
			if tt.checkResponseContent != nil && rec.Code == http.StatusOK {
				tt.checkResponseContent(t, rec.Body.Bytes())
			}
		})
	}
}

func TestGetExpectedMachineHandler_Handle(t *testing.T) {
	// Setup
	e := echo.New()
	dbSession := testExpectedMachineInitDB(t)
	defer dbSession.Close()

	ctx := context.Background()

	cfg := &config.Config{}
	handler := NewGetExpectedMachineHandler(dbSession, nil, cfg)

	org := "test-org"
	infraProv, site := testExpectedMachineSetupTestData(t, dbSession, org)

	// Create an unmanaged site
	unmanagedIP := &cdbm.InfrastructureProvider{
		ID:   uuid.New(),
		Name: "unmanaged-provider",
		Org:  "other-org",
	}
	_, err := dbSession.DB.NewInsert().Model(unmanagedIP).Exec(ctx)
	assert.Nil(t, err)

	unmanagedSite := &cdbm.Site{
		ID:                       uuid.New(),
		Name:                     "unmanaged-site",
		Org:                      "other-org",
		InfrastructureProviderID: unmanagedIP.ID,
		Status:                   cdbm.SiteStatusRegistered,
	}
	_, err = dbSession.DB.NewInsert().Model(unmanagedSite).Exec(ctx)
	assert.Nil(t, err)

	// Create a test ExpectedMachine on managed site
	emDAO := cdbm.NewExpectedMachineDAO(dbSession)
	testEM, err := emDAO.Create(ctx, nil, cdbm.ExpectedMachineCreateInput{
		ExpectedMachineID:        uuid.New(),
		SiteID:                   site.ID,
		BmcMacAddress:            "00:11:22:33:44:55",
		ChassisSerialNumber:      "TEST-CHASSIS-123",
		FallbackDpuSerialNumbers: []string{"DPU001"},
		Labels:                   map[string]string{"env": "test"},
	})
	assert.Nil(t, err)
	assert.NotNil(t, testEM)

	// Create a test ExpectedMachine on unmanaged site
	unmanagedEM, err := emDAO.Create(ctx, nil, cdbm.ExpectedMachineCreateInput{
		ExpectedMachineID:        uuid.New(),
		SiteID:                   unmanagedSite.ID,
		BmcMacAddress:            "00:11:22:33:44:CC",
		ChassisSerialNumber:      "UNMANAGED-CHASSIS-456",
		FallbackDpuSerialNumbers: []string{"DPU002"},
		Labels:                   map[string]string{"env": "unmanaged"},
	})
	assert.Nil(t, err)
	assert.NotNil(t, unmanagedEM)

	// Helper function to create mock user
	createMockUser := func(org string) *cdbm.User {
		return &cdbm.User{
			StarfleetID: cdb.GetStrPtr("test-user"),
			OrgData: cdbm.OrgData{
				org: cdbm.Org{
					ID:          123,
					Name:        org,
					DisplayName: org,
					OrgType:     "ENTERPRISE",
					Roles:       []string{"FORGE_PROVIDER_ADMIN"},
				},
			},
		}
	}

	tests := []struct {
		name           string
		id             string
		setupContext   func(c echo.Context)
		expectedStatus int
	}{
		{
			name: "invalid ID",
			id:   "invalid-id",
			setupContext: func(c echo.Context) {
				c.Set("user", createMockUser(org))
				c.SetParamNames("orgName", "id")
				c.SetParamValues(org, "invalid-id")
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "successful retrieval",
			id:   testEM.ID.String(),
			setupContext: func(c echo.Context) {
				c.Set("user", createMockUser(org))
				c.SetParamNames("orgName", "id")
				c.SetParamValues(org, testEM.ID.String())
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "machine not found",
			id:   "12345678-1234-1234-1234-123456789099",
			setupContext: func(c echo.Context) {
				c.Set("user", createMockUser(org))
				c.SetParamNames("orgName", "id")
				c.SetParamValues(org, "12345678-1234-1234-1234-123456789099")
			},
			expectedStatus: http.StatusNotFound,
		},
		{
			name: "cannot retrieve from unmanaged site",
			id:   unmanagedEM.ID.String(),
			setupContext: func(c echo.Context) {
				c.Set("user", createMockUser(org))
				c.SetParamNames("orgName", "id")
				c.SetParamValues(org, unmanagedEM.ID.String())
			},
			expectedStatus: http.StatusForbidden,
		},
	}

	_ = infraProv // Ensure infraProv is used to avoid compiler warning

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			url := "/v2/org/" + org + "/carbide/expected-machine/" + tt.id
			req := httptest.NewRequest(http.MethodGet, url, nil)
			req = req.WithContext(context.Background())

			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			// Setup context
			tt.setupContext(c)

			// Execute
			err := handler.Handle(c)

			// Assert
			assert.Nil(t, err)
			assert.Equal(t, tt.expectedStatus, rec.Code)
			if tt.expectedStatus != rec.Code {
				t.Errorf("Response: %v", rec.Body.String())
			}

			// For successful retrieval, verify labels are returned in response
			if tt.expectedStatus == http.StatusOK && tt.name == "successful retrieval" {
				var response model.APIExpectedMachine
				err := json.Unmarshal(rec.Body.Bytes(), &response)
				assert.Nil(t, err)
				assert.NotNil(t, response.Labels, "Labels should not be nil in response")
				assert.Equal(t, "test", response.Labels["env"], "Labels in response should contain the 'env' label with value 'test'")
			}
		})
	}
}

func TestUpdateExpectedMachineHandler_Handle(t *testing.T) {
	// Setup
	e := echo.New()
	dbSession := testExpectedMachineInitDB(t)
	defer dbSession.Close()

	ctx := context.Background()
	cfg := common.GetTestConfig()

	// Prepare client pool for workflow calls
	tcfg, _ := cfg.GetTemporalConfig()
	scp := sc.NewClientPool(tcfg)

	org := "test-org"
	infraProv, site := testExpectedMachineSetupTestData(t, dbSession, org)

	// Create an unmanaged site
	unmanagedIP := &cdbm.InfrastructureProvider{
		ID:   uuid.New(),
		Name: "unmanaged-provider",
		Org:  "other-org",
	}
	_, err := dbSession.DB.NewInsert().Model(unmanagedIP).Exec(ctx)
	assert.Nil(t, err)

	unmanagedSite := &cdbm.Site{
		ID:                       uuid.New(),
		Name:                     "unmanaged-site",
		Org:                      "other-org",
		InfrastructureProviderID: unmanagedIP.ID,
		Status:                   cdbm.SiteStatusRegistered,
	}
	_, err = dbSession.DB.NewInsert().Model(unmanagedSite).Exec(ctx)
	assert.Nil(t, err)

	// Create a test ExpectedMachine on managed site
	emDAO := cdbm.NewExpectedMachineDAO(dbSession)
	testEM, err := emDAO.Create(ctx, nil, cdbm.ExpectedMachineCreateInput{
		ExpectedMachineID:        uuid.New(),
		SiteID:                   site.ID,
		BmcMacAddress:            "00:11:22:33:44:DD",
		ChassisSerialNumber:      "UPDATE-CHASSIS-123",
		FallbackDpuSerialNumbers: []string{"DPU001"},
		Labels:                   map[string]string{"env": "test"},
	})
	assert.Nil(t, err)
	assert.NotNil(t, testEM)

	// Create a test ExpectedMachine on unmanaged site
	unmanagedEM, err := emDAO.Create(ctx, nil, cdbm.ExpectedMachineCreateInput{
		ExpectedMachineID:        uuid.New(),
		SiteID:                   unmanagedSite.ID,
		BmcMacAddress:            "00:11:22:33:44:EE",
		ChassisSerialNumber:      "UNMANAGED-UPDATE-456",
		FallbackDpuSerialNumbers: []string{"DPU002"},
		Labels:                   map[string]string{"env": "unmanaged"},
	})
	assert.Nil(t, err)
	assert.NotNil(t, unmanagedEM)

	// Add mock temporal client for the site
	mockTemporalClient := &tmocks.Client{}
	mockWorkflowRun := &tmocks.WorkflowRun{}
	mockWorkflowRun.On("GetID").Return("test-workflow-id")
	mockWorkflowRun.Mock.On("Get", mock.Anything, mock.Anything).Return(nil)
	mockTemporalClient.Mock.On("ExecuteWorkflow", mock.Anything, mock.Anything, "UpdateExpectedMachine", mock.Anything).Return(mockWorkflowRun, nil)
	scp.IDClientMap[site.ID.String()] = mockTemporalClient

	handler := NewUpdateExpectedMachineHandler(dbSession, nil, scp, cfg)

	// Helper function to create mock user
	createMockUser := func(org string) *cdbm.User {
		return &cdbm.User{
			StarfleetID: cdb.GetStrPtr("test-user"),
			OrgData: cdbm.OrgData{
				org: cdbm.Org{
					ID:          123,
					Name:        org,
					DisplayName: org,
					OrgType:     "ENTERPRISE",
					Roles:       []string{"FORGE_PROVIDER_ADMIN"},
				},
			},
		}
	}

	tests := []struct {
		name                 string
		id                   string
		requestBody          model.APIExpectedMachineUpdateRequest
		setupContext         func(c echo.Context)
		expectedStatus       int
		checkResponseContent func(t *testing.T, body []byte)
	}{
		{
			name: "successful update",
			id:   testEM.ID.String(),
			requestBody: model.APIExpectedMachineUpdateRequest{
				ChassisSerialNumber: cdb.GetStrPtr("UPDATED-CHASSIS-123"),
				Labels:              map[string]string{"env": "updated"},
			},
			setupContext: func(c echo.Context) {
				c.Set("user", createMockUser(org))
				c.SetParamNames("orgName", "id")
				c.SetParamValues(org, testEM.ID.String())
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "successful MAC address update",
			id:   testEM.ID.String(),
			requestBody: model.APIExpectedMachineUpdateRequest{
				BmcMacAddress: cdb.GetStrPtr("AA:BB:CC:DD:EE:FF"),
			},
			setupContext: func(c echo.Context) {
				c.Set("user", createMockUser(org))
				c.SetParamNames("orgName", "id")
				c.SetParamValues(org, testEM.ID.String())
			},
			expectedStatus: http.StatusOK,
			checkResponseContent: func(t *testing.T, body []byte) {
				var response model.APIExpectedMachine
				err := json.Unmarshal(body, &response)
				assert.Nil(t, err)
				assert.Equal(t, "AA:BB:CC:DD:EE:FF", response.BmcMacAddress, "MAC address in response should match the updated value")
			},
		},
		{
			name: "cannot update on unmanaged site",
			id:   unmanagedEM.ID.String(),
			requestBody: model.APIExpectedMachineUpdateRequest{
				ChassisSerialNumber: cdb.GetStrPtr("SHOULD-NOT-UPDATE"),
				Labels:              map[string]string{"env": "fail"},
			},
			setupContext: func(c echo.Context) {
				c.Set("user", createMockUser(org))
				c.SetParamNames("orgName", "id")
				c.SetParamValues(org, unmanagedEM.ID.String())
			},
			expectedStatus: http.StatusForbidden,
		},
	}

	_ = infraProv // Ensure infraProv is used to avoid compiler warning

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reqBody, _ := json.Marshal(tt.requestBody)
			url := "/v2/org/" + org + "/carbide/expected-machine/" + tt.id
			req := httptest.NewRequest(http.MethodPatch, url, bytes.NewReader(reqBody))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			req = req.WithContext(context.Background())

			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			// Setup context
			tt.setupContext(c)

			// Execute
			err := handler.Handle(c)

			// Assert
			assert.Nil(t, err)
			assert.Equal(t, tt.expectedStatus, rec.Code)
			if tt.expectedStatus != rec.Code {
				t.Errorf("Response: %v", rec.Body.String())
			}

			// Check response content if provided
			if tt.checkResponseContent != nil && rec.Code == http.StatusOK {
				tt.checkResponseContent(t, rec.Body.Bytes())
			}
		})
	}
}

func TestDeleteExpectedMachineHandler_Handle(t *testing.T) {
	// Setup
	e := echo.New()
	dbSession := testExpectedMachineInitDB(t)
	defer dbSession.Close()

	ctx := context.Background()
	cfg := common.GetTestConfig()

	// Prepare client pool for workflow calls
	tcfg, _ := cfg.GetTemporalConfig()
	scp := sc.NewClientPool(tcfg)

	org := "test-org"
	infraProv, site := testExpectedMachineSetupTestData(t, dbSession, org)

	// Create an unmanaged site
	unmanagedIP := &cdbm.InfrastructureProvider{
		ID:   uuid.New(),
		Name: "unmanaged-provider",
		Org:  "other-org",
	}
	_, err := dbSession.DB.NewInsert().Model(unmanagedIP).Exec(ctx)
	assert.Nil(t, err)

	unmanagedSite := &cdbm.Site{
		ID:                       uuid.New(),
		Name:                     "unmanaged-site",
		Org:                      "other-org",
		InfrastructureProviderID: unmanagedIP.ID,
		Status:                   cdbm.SiteStatusRegistered,
	}
	_, err = dbSession.DB.NewInsert().Model(unmanagedSite).Exec(ctx)
	assert.Nil(t, err)

	// Create a test ExpectedMachine on managed site (to be deleted)
	emDAO := cdbm.NewExpectedMachineDAO(dbSession)
	testEM, err := emDAO.Create(ctx, nil, cdbm.ExpectedMachineCreateInput{
		ExpectedMachineID:        uuid.New(),
		SiteID:                   site.ID,
		BmcMacAddress:            "00:11:22:33:44:FF",
		ChassisSerialNumber:      "DELETE-CHASSIS-123",
		FallbackDpuSerialNumbers: []string{"DPU001"},
	})
	assert.Nil(t, err)
	assert.NotNil(t, testEM)

	// Create a test ExpectedMachine on unmanaged site (should not be deletable)
	unmanagedEM, err := emDAO.Create(ctx, nil, cdbm.ExpectedMachineCreateInput{
		ExpectedMachineID:        uuid.New(),
		SiteID:                   unmanagedSite.ID,
		BmcMacAddress:            "00:11:22:33:55:00",
		ChassisSerialNumber:      "UNMANAGED-DELETE-456",
		FallbackDpuSerialNumbers: []string{"DPU002"},
	})
	assert.Nil(t, err)
	assert.NotNil(t, unmanagedEM)

	// Add mock temporal client for the site
	mockTemporalClient := &tmocks.Client{}
	mockWorkflowRun := &tmocks.WorkflowRun{}
	mockWorkflowRun.On("GetID").Return("test-workflow-id")
	mockWorkflowRun.Mock.On("Get", mock.Anything, mock.Anything).Return(nil)
	mockTemporalClient.Mock.On("ExecuteWorkflow", mock.Anything, mock.Anything, "DeleteExpectedMachine", mock.Anything).Return(mockWorkflowRun, nil)
	scp.IDClientMap[site.ID.String()] = mockTemporalClient

	handler := NewDeleteExpectedMachineHandler(dbSession, nil, scp, cfg)

	// Helper function to create mock user
	createMockUser := func(org string) *cdbm.User {
		return &cdbm.User{
			StarfleetID: cdb.GetStrPtr("test-user"),
			OrgData: cdbm.OrgData{
				org: cdbm.Org{
					ID:          123,
					Name:        org,
					DisplayName: org,
					OrgType:     "ENTERPRISE",
					Roles:       []string{"FORGE_PROVIDER_ADMIN"},
				},
			},
		}
	}

	tests := []struct {
		name           string
		id             string
		setupContext   func(c echo.Context)
		expectedStatus int
	}{
		{
			name: "successful delete",
			id:   testEM.ID.String(),
			setupContext: func(c echo.Context) {
				c.Set("user", createMockUser(org))
				c.SetParamNames("orgName", "id")
				c.SetParamValues(org, testEM.ID.String())
			},
			expectedStatus: http.StatusNoContent,
		},
		{
			name: "cannot delete on unmanaged site",
			id:   unmanagedEM.ID.String(),
			setupContext: func(c echo.Context) {
				c.Set("user", createMockUser(org))
				c.SetParamNames("orgName", "id")
				c.SetParamValues(org, unmanagedEM.ID.String())
			},
			expectedStatus: http.StatusForbidden,
		},
	}

	_ = infraProv // Ensure infraProv is used to avoid compiler warning

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			url := "/v2/org/" + org + "/carbide/expected-machine/" + tt.id
			req := httptest.NewRequest(http.MethodDelete, url, nil)
			req = req.WithContext(context.Background())

			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			// Setup context
			tt.setupContext(c)

			// Execute
			err := handler.Handle(c)

			// Assert
			assert.Nil(t, err)
			assert.Equal(t, tt.expectedStatus, rec.Code)
			if tt.expectedStatus != rec.Code {
				t.Errorf("Response: %v", rec.Body.String())
			}
		})
	}
}
