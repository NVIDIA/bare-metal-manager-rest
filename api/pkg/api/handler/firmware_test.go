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
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/nvidia/bare-metal-manager-rest/api/pkg/api/handler/util/common"
	"github.com/nvidia/bare-metal-manager-rest/api/pkg/api/model"
	sc "github.com/nvidia/bare-metal-manager-rest/api/pkg/client/site"
	"github.com/nvidia/bare-metal-manager-rest/common/pkg/otelecho"
	cdbm "github.com/nvidia/bare-metal-manager-rest/db/pkg/db/model"
	rlav1 "github.com/nvidia/bare-metal-manager-rest/workflow-schema/rla/protobuf/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	oteltrace "go.opentelemetry.io/otel/trace"
	tmocks "go.temporal.io/sdk/mocks"
)

func TestFirmwareUpgradeRackHandler_Handle(t *testing.T) {
	e := echo.New()
	dbSession := testRackInitDB(t)
	defer dbSession.Close()

	cfg := common.GetTestConfig()
	tcfg, _ := cfg.GetTemporalConfig()
	scp := sc.NewClientPool(tcfg)

	org := "test-org"
	_, site, _ := testRackSetupTestData(t, dbSession, org)

	providerUser := testRackBuildUser(t, dbSession, "provider-user-fw-rack", org, []string{"FORGE_PROVIDER_ADMIN"})
	tenantUser := testRackBuildUser(t, dbSession, "tenant-user-fw-rack", org, []string{"FORGE_TENANT_ADMIN"})

	handler := NewFirmwareUpgradeRackHandler(dbSession, nil, scp, cfg)

	rackID := uuid.New().String()

	tracer := oteltrace.NewNoopTracerProvider().Tracer("test")
	ctx := context.Background()

	tests := []struct {
		name           string
		reqOrg         string
		user           *cdbm.User
		rackID         string
		queryParams    map[string]string
		body           string
		mockTaskIDs    []*rlav1.UUID
		expectedStatus int
	}{
		{
			name:   "success - firmware upgrade with version",
			reqOrg: org,
			user:   providerUser,
			rackID: rackID,
			queryParams: map[string]string{
				"siteId": site.ID.String(),
			},
			body:           `{"version":"24.11.0"}`,
			mockTaskIDs:    []*rlav1.UUID{{Id: uuid.NewString()}},
			expectedStatus: http.StatusOK,
		},
		{
			name:   "success - firmware upgrade without version",
			reqOrg: org,
			user:   providerUser,
			rackID: rackID,
			queryParams: map[string]string{
				"siteId": site.ID.String(),
			},
			body:           `{}`,
			mockTaskIDs:    []*rlav1.UUID{{Id: uuid.NewString()}},
			expectedStatus: http.StatusOK,
		},
		{
			name:   "failure - missing siteId",
			reqOrg: org,
			user:   providerUser,
			rackID: rackID,
			queryParams: map[string]string{
				// no siteId
			},
			body:           `{}`,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:   "failure - invalid siteId",
			reqOrg: org,
			user:   providerUser,
			rackID: rackID,
			queryParams: map[string]string{
				"siteId": uuid.New().String(),
			},
			body:           `{}`,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:   "failure - tenant access denied",
			reqOrg: org,
			user:   tenantUser,
			rackID: rackID,
			queryParams: map[string]string{
				"siteId": site.ID.String(),
			},
			body:           `{}`,
			expectedStatus: http.StatusForbidden,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockTemporalClient := &tmocks.Client{}
			mockWorkflowRun := &tmocks.WorkflowRun{}
			mockWorkflowRun.On("GetID").Return("test-workflow-id")
			mockWorkflowRun.Mock.On("Get", mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
				resp := args.Get(1).(*rlav1.SubmitTaskResponse)
				if tt.mockTaskIDs != nil {
					resp.TaskIds = tt.mockTaskIDs
				}
			}).Return(nil)
			mockTemporalClient.Mock.On("ExecuteWorkflow", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(mockWorkflowRun, nil)
			scp.IDClientMap[site.ID.String()] = mockTemporalClient

			q := url.Values{}
			for k, v := range tt.queryParams {
				q.Set(k, v)
			}
			path := fmt.Sprintf("/v2/org/%s/carbide/rack/%s/firmware?%s", tt.reqOrg, tt.rackID, q.Encode())

			req := httptest.NewRequest(http.MethodPatch, path, strings.NewReader(tt.body))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()

			ec := e.NewContext(req, rec)
			ec.SetParamNames("orgName", "id")
			ec.SetParamValues(tt.reqOrg, tt.rackID)
			ec.Set("user", tt.user)

			ctx = context.WithValue(ctx, otelecho.TracerKey, tracer)
			ec.SetRequest(ec.Request().WithContext(ctx))

			err := handler.Handle(ec)

			if tt.expectedStatus != rec.Code {
				t.Errorf("FirmwareUpgradeRackHandler.Handle() status = %v, want %v, response: %v, err: %v", rec.Code, tt.expectedStatus, rec.Body.String(), err)
			}

			require.Equal(t, tt.expectedStatus, rec.Code)
			if tt.expectedStatus != http.StatusOK {
				return
			}

			var apiResp model.APIFirmwareUpgradeResponse
			err = json.Unmarshal(rec.Body.Bytes(), &apiResp)
			assert.NoError(t, err)
			assert.NotEmpty(t, apiResp.TaskIDs)
		})
	}
}

func TestFirmwareUpgradeRackBatchHandler_Handle(t *testing.T) {
	e := echo.New()
	dbSession := testRackInitDB(t)
	defer dbSession.Close()

	cfg := common.GetTestConfig()
	tcfg, _ := cfg.GetTemporalConfig()
	scp := sc.NewClientPool(tcfg)

	org := "test-org"
	_, site, _ := testRackSetupTestData(t, dbSession, org)

	providerUser := testRackBuildUser(t, dbSession, "provider-user-fw-rack-batch", org, []string{"FORGE_PROVIDER_ADMIN"})
	tenantUser := testRackBuildUser(t, dbSession, "tenant-user-fw-rack-batch", org, []string{"FORGE_TENANT_ADMIN"})

	handler := NewFirmwareUpgradeRackBatchHandler(dbSession, nil, scp, cfg)

	tracer := oteltrace.NewNoopTracerProvider().Tracer("test")
	ctx := context.Background()

	tests := []struct {
		name           string
		reqOrg         string
		user           *cdbm.User
		queryParams    map[string]string
		body           string
		mockTaskIDs    []*rlav1.UUID
		expectedStatus int
	}{
		{
			name:   "success - firmware upgrade all racks (no filter)",
			reqOrg: org,
			user:   providerUser,
			queryParams: map[string]string{
				"siteId": site.ID.String(),
			},
			body:           `{}`,
			mockTaskIDs:    []*rlav1.UUID{{Id: uuid.NewString()}, {Id: uuid.NewString()}},
			expectedStatus: http.StatusOK,
		},
		{
			name:   "success - firmware upgrade with name filter and version",
			reqOrg: org,
			user:   providerUser,
			queryParams: map[string]string{
				"siteId": site.ID.String(),
				"name":   "rack-1",
			},
			body:           `{"version":"24.11.0"}`,
			mockTaskIDs:    []*rlav1.UUID{{Id: uuid.NewString()}},
			expectedStatus: http.StatusOK,
		},
		{
			name:   "failure - missing siteId",
			reqOrg: org,
			user:   providerUser,
			queryParams: map[string]string{
				// no siteId
			},
			body:           `{}`,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:   "failure - tenant access denied",
			reqOrg: org,
			user:   tenantUser,
			queryParams: map[string]string{
				"siteId": site.ID.String(),
			},
			body:           `{}`,
			expectedStatus: http.StatusForbidden,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockTemporalClient := &tmocks.Client{}
			mockWorkflowRun := &tmocks.WorkflowRun{}
			mockWorkflowRun.On("GetID").Return("test-workflow-id")
			mockWorkflowRun.Mock.On("Get", mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
				resp := args.Get(1).(*rlav1.SubmitTaskResponse)
				if tt.mockTaskIDs != nil {
					resp.TaskIds = tt.mockTaskIDs
				}
			}).Return(nil)
			mockTemporalClient.Mock.On("ExecuteWorkflow", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(mockWorkflowRun, nil)
			scp.IDClientMap[site.ID.String()] = mockTemporalClient

			q := url.Values{}
			for k, v := range tt.queryParams {
				q.Set(k, v)
			}
			path := fmt.Sprintf("/v2/org/%s/carbide/rack/firmware?%s", tt.reqOrg, q.Encode())

			req := httptest.NewRequest(http.MethodPatch, path, strings.NewReader(tt.body))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()

			ec := e.NewContext(req, rec)
			ec.SetParamNames("orgName")
			ec.SetParamValues(tt.reqOrg)
			ec.Set("user", tt.user)

			ctx = context.WithValue(ctx, otelecho.TracerKey, tracer)
			ec.SetRequest(ec.Request().WithContext(ctx))

			err := handler.Handle(ec)

			if tt.expectedStatus != rec.Code {
				t.Errorf("FirmwareUpgradeRackBatchHandler.Handle() status = %v, want %v, response: %v, err: %v", rec.Code, tt.expectedStatus, rec.Body.String(), err)
			}

			require.Equal(t, tt.expectedStatus, rec.Code)
			if tt.expectedStatus != http.StatusOK {
				return
			}

			var apiResp model.APIFirmwareUpgradeResponse
			err = json.Unmarshal(rec.Body.Bytes(), &apiResp)
			assert.NoError(t, err)
			assert.NotEmpty(t, apiResp.TaskIDs)
		})
	}
}

func TestFirmwareUpgradeTrayHandler_Handle(t *testing.T) {
	e := echo.New()
	dbSession := testTrayInitDB(t)
	defer dbSession.Close()

	cfg := common.GetTestConfig()
	tcfg, _ := cfg.GetTemporalConfig()
	scp := sc.NewClientPool(tcfg)

	org := "test-org"
	_, site, _ := testTraySetupTestData(t, dbSession, org)

	providerUser := testTrayBuildUser(t, dbSession, "provider-user-fw-tray", org, []string{"FORGE_PROVIDER_ADMIN"})
	tenantUser := testTrayBuildUser(t, dbSession, "tenant-user-fw-tray", org, []string{"FORGE_TENANT_ADMIN"})

	handler := NewFirmwareUpgradeTrayHandler(dbSession, nil, scp, cfg)

	trayID := uuid.New().String()

	tracer := oteltrace.NewNoopTracerProvider().Tracer("test")
	ctx := context.Background()

	tests := []struct {
		name           string
		reqOrg         string
		user           *cdbm.User
		trayID         string
		queryParams    map[string]string
		body           string
		mockTaskIDs    []*rlav1.UUID
		expectedStatus int
	}{
		{
			name:   "success - firmware upgrade with version",
			reqOrg: org,
			user:   providerUser,
			trayID: trayID,
			queryParams: map[string]string{
				"siteId": site.ID.String(),
			},
			body:           `{"version":"24.11.0"}`,
			mockTaskIDs:    []*rlav1.UUID{{Id: uuid.NewString()}},
			expectedStatus: http.StatusOK,
		},
		{
			name:   "success - firmware upgrade without version",
			reqOrg: org,
			user:   providerUser,
			trayID: trayID,
			queryParams: map[string]string{
				"siteId": site.ID.String(),
			},
			body:           `{}`,
			mockTaskIDs:    []*rlav1.UUID{{Id: uuid.NewString()}},
			expectedStatus: http.StatusOK,
		},
		{
			name:   "failure - missing siteId",
			reqOrg: org,
			user:   providerUser,
			trayID: trayID,
			queryParams: map[string]string{
				// no siteId
			},
			body:           `{}`,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:   "failure - invalid tray ID (not UUID)",
			reqOrg: org,
			user:   providerUser,
			trayID: "not-a-uuid",
			queryParams: map[string]string{
				"siteId": site.ID.String(),
			},
			body:           `{}`,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:   "failure - tenant access denied",
			reqOrg: org,
			user:   tenantUser,
			trayID: trayID,
			queryParams: map[string]string{
				"siteId": site.ID.String(),
			},
			body:           `{}`,
			expectedStatus: http.StatusForbidden,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockTemporalClient := &tmocks.Client{}
			mockWorkflowRun := &tmocks.WorkflowRun{}
			mockWorkflowRun.On("GetID").Return("test-workflow-id")
			mockWorkflowRun.Mock.On("Get", mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
				resp := args.Get(1).(*rlav1.SubmitTaskResponse)
				if tt.mockTaskIDs != nil {
					resp.TaskIds = tt.mockTaskIDs
				}
			}).Return(nil)
			mockTemporalClient.Mock.On("ExecuteWorkflow", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(mockWorkflowRun, nil)
			scp.IDClientMap[site.ID.String()] = mockTemporalClient

			q := url.Values{}
			for k, v := range tt.queryParams {
				q.Set(k, v)
			}
			path := fmt.Sprintf("/v2/org/%s/carbide/tray/%s/firmware?%s", tt.reqOrg, tt.trayID, q.Encode())

			req := httptest.NewRequest(http.MethodPatch, path, strings.NewReader(tt.body))
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
				t.Errorf("FirmwareUpgradeTrayHandler.Handle() status = %v, want %v, response: %v, err: %v", rec.Code, tt.expectedStatus, rec.Body.String(), err)
			}

			require.Equal(t, tt.expectedStatus, rec.Code)
			if tt.expectedStatus != http.StatusOK {
				return
			}

			var apiResp model.APIFirmwareUpgradeResponse
			err = json.Unmarshal(rec.Body.Bytes(), &apiResp)
			assert.NoError(t, err)
			assert.NotEmpty(t, apiResp.TaskIDs)
		})
	}
}

func TestFirmwareUpgradeTrayBatchHandler_Handle(t *testing.T) {
	e := echo.New()
	dbSession := testTrayInitDB(t)
	defer dbSession.Close()

	cfg := common.GetTestConfig()
	tcfg, _ := cfg.GetTemporalConfig()
	scp := sc.NewClientPool(tcfg)

	org := "test-org"
	_, site, _ := testTraySetupTestData(t, dbSession, org)

	providerUser := testTrayBuildUser(t, dbSession, "provider-user-fw-tray-batch", org, []string{"FORGE_PROVIDER_ADMIN"})
	tenantUser := testTrayBuildUser(t, dbSession, "tenant-user-fw-tray-batch", org, []string{"FORGE_TENANT_ADMIN"})

	handler := NewFirmwareUpgradeTrayBatchHandler(dbSession, nil, scp, cfg)

	tracer := oteltrace.NewNoopTracerProvider().Tracer("test")
	ctx := context.Background()

	tests := []struct {
		name           string
		reqOrg         string
		user           *cdbm.User
		queryParams    map[string]string
		body           string
		mockTaskIDs    []*rlav1.UUID
		expectedStatus int
	}{
		{
			name:   "success - firmware upgrade all trays (no filter)",
			reqOrg: org,
			user:   providerUser,
			queryParams: map[string]string{
				"siteId": site.ID.String(),
			},
			body:           `{}`,
			mockTaskIDs:    []*rlav1.UUID{{Id: uuid.NewString()}, {Id: uuid.NewString()}},
			expectedStatus: http.StatusOK,
		},
		{
			name:   "success - firmware upgrade with rackId filter and version",
			reqOrg: org,
			user:   providerUser,
			queryParams: map[string]string{
				"siteId": site.ID.String(),
				"rackId": uuid.NewString(),
			},
			body:           `{"version":"24.11.0"}`,
			mockTaskIDs:    []*rlav1.UUID{{Id: uuid.NewString()}},
			expectedStatus: http.StatusOK,
		},
		{
			name:   "failure - missing siteId",
			reqOrg: org,
			user:   providerUser,
			queryParams: map[string]string{
				// no siteId
			},
			body:           `{}`,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:   "failure - tenant access denied",
			reqOrg: org,
			user:   tenantUser,
			queryParams: map[string]string{
				"siteId": site.ID.String(),
			},
			body:           `{}`,
			expectedStatus: http.StatusForbidden,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockTemporalClient := &tmocks.Client{}
			mockWorkflowRun := &tmocks.WorkflowRun{}
			mockWorkflowRun.On("GetID").Return("test-workflow-id")
			mockWorkflowRun.Mock.On("Get", mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
				resp := args.Get(1).(*rlav1.SubmitTaskResponse)
				if tt.mockTaskIDs != nil {
					resp.TaskIds = tt.mockTaskIDs
				}
			}).Return(nil)
			mockTemporalClient.Mock.On("ExecuteWorkflow", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(mockWorkflowRun, nil)
			scp.IDClientMap[site.ID.String()] = mockTemporalClient

			q := url.Values{}
			for k, v := range tt.queryParams {
				q.Set(k, v)
			}
			path := fmt.Sprintf("/v2/org/%s/carbide/tray/firmware?%s", tt.reqOrg, q.Encode())

			req := httptest.NewRequest(http.MethodPatch, path, strings.NewReader(tt.body))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()

			ec := e.NewContext(req, rec)
			ec.SetParamNames("orgName")
			ec.SetParamValues(tt.reqOrg)
			ec.Set("user", tt.user)

			ctx = context.WithValue(ctx, otelecho.TracerKey, tracer)
			ec.SetRequest(ec.Request().WithContext(ctx))

			err := handler.Handle(ec)

			if tt.expectedStatus != rec.Code {
				t.Errorf("FirmwareUpgradeTrayBatchHandler.Handle() status = %v, want %v, response: %v, err: %v", rec.Code, tt.expectedStatus, rec.Body.String(), err)
			}

			require.Equal(t, tt.expectedStatus, rec.Code)
			if tt.expectedStatus != http.StatusOK {
				return
			}

			var apiResp model.APIFirmwareUpgradeResponse
			err = json.Unmarshal(rec.Body.Bytes(), &apiResp)
			assert.NoError(t, err)
			assert.NotEmpty(t, apiResp.TaskIDs)
		})
	}
}
