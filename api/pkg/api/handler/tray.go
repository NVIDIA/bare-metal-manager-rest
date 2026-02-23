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
	"errors"
	"fmt"
	"maps"
	"net/http"
	"slices"
	"strings"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"go.opentelemetry.io/otel/attribute"
	tClient "go.temporal.io/sdk/client"

	"github.com/nvidia/bare-metal-manager-rest/api/internal/config"
	"github.com/nvidia/bare-metal-manager-rest/api/pkg/api/handler/util/common"
	"github.com/nvidia/bare-metal-manager-rest/api/pkg/api/model"
	"github.com/nvidia/bare-metal-manager-rest/api/pkg/api/pagination"
	sc "github.com/nvidia/bare-metal-manager-rest/api/pkg/client/site"
	auth "github.com/nvidia/bare-metal-manager-rest/auth/pkg/authorization"
	cerr "github.com/nvidia/bare-metal-manager-rest/common/pkg/util"
	sutil "github.com/nvidia/bare-metal-manager-rest/common/pkg/util"
	cdb "github.com/nvidia/bare-metal-manager-rest/db/pkg/db"
	rlav1 "github.com/nvidia/bare-metal-manager-rest/workflow-schema/rla/protobuf/v1"
	"github.com/nvidia/bare-metal-manager-rest/workflow/pkg/queue"
	temporalEnums "go.temporal.io/api/enums/v1"
	tp "go.temporal.io/sdk/temporal"
)

// Allowed query parameters for each tray handler
var (
	getTrayAllowedParams              = []string{"siteId"}
	getAllTrayAllowedParams           = []string{"siteId", "rackId", "rackName", "type", "componentId", "id", "pageNumber", "pageSize", "orderBy"}
	powerControlTrayAllowedParams          = []string{"siteId"}
	powerControlTrayBatchAllowedParams     = []string{"siteId", "rackId", "rackName", "type", "componentId", "id"}
	firmwareUpgradeTrayAllowedParams       = []string{"siteId"}
	firmwareUpgradeTrayBatchAllowedParams  = []string{"siteId", "rackId", "rackName", "type", "componentId", "id"}
)

// ~~~~~ Get Tray Handler ~~~~~ //

// GetTrayHandler is the API Handler for getting a Tray by ID
type GetTrayHandler struct {
	dbSession  *cdb.Session
	tc         tClient.Client
	scp        *sc.ClientPool
	cfg        *config.Config
	tracerSpan *sutil.TracerSpan
}

// NewGetTrayHandler initializes and returns a new handler for getting a Tray
func NewGetTrayHandler(dbSession *cdb.Session, tc tClient.Client, scp *sc.ClientPool, cfg *config.Config) GetTrayHandler {
	return GetTrayHandler{
		dbSession:  dbSession,
		tc:         tc,
		scp:        scp,
		cfg:        cfg,
		tracerSpan: sutil.NewTracerSpan(),
	}
}

// Handle godoc
// @Summary Get a Tray
// @Description Get a Tray by ID
// @Tags tray
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param org path string true "Name of NGC organization"
// @Param id path string true "ID of Tray"
// @Param siteId query string true "ID of the Site"
// @Success 200 {object} model.APITray
// @Router /v2/org/{org}/carbide/tray/{id} [get]
func (gth GetTrayHandler) Handle(c echo.Context) error {
	org, dbUser, ctx, logger, handlerSpan := common.SetupHandler("Tray", "Get", c, gth.tracerSpan)
	if handlerSpan != nil {
		defer handlerSpan.End()
	}

	if apiErr := common.ValidateQueryParams(c.QueryParams(), getTrayAllowedParams); apiErr != nil {
		return cerr.NewAPIErrorResponse(c, apiErr.Code, apiErr.Message, nil)
	}

	// Is DB user missing?
	if dbUser == nil {
		logger.Error().Msg("invalid User object found in request context")
		return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve current user", nil)
	}

	// Validate org membership
	ok, err := auth.ValidateOrgMembership(dbUser, org)
	if !ok {
		if err != nil {
			logger.Error().Err(err).Msg("error validating org membership for User in request")
		} else {
			logger.Warn().Msg("could not validate org membership for user, access denied")
		}
		return cerr.NewAPIErrorResponse(c, http.StatusForbidden, fmt.Sprintf("Failed to validate membership for org: %s", org), nil)
	}

	// Validate role, only Provider Admins are allowed to access Tray data
	ok = auth.ValidateUserRoles(dbUser, org, nil, auth.ProviderAdminRole)
	if !ok {
		logger.Warn().Msg("user does not have Provider Admin role, access denied")
		return cerr.NewAPIErrorResponse(c, http.StatusForbidden, "User does not have Provider Admin role with org", nil)
	}

	// Get Infrastructure Provider for org
	infrastructureProvider, err := common.GetInfrastructureProviderForOrg(ctx, nil, gth.dbSession, org)
	if err != nil {
		logger.Warn().Err(err).Msg("error getting infrastructure provider for org")
		return cerr.NewAPIErrorResponse(c, http.StatusBadRequest, "Failed to retrieve Infrastructure Provider for org", nil)
	}

	// Validate siteId is provided
	siteStrID := c.QueryParam("siteId")
	if siteStrID == "" {
		return cerr.NewAPIErrorResponse(c, http.StatusBadRequest, "siteId query parameter is required", nil)
	}

	// Retrieve the Site from the DB
	site, err := common.GetSiteFromIDString(ctx, nil, siteStrID, gth.dbSession)
	if err != nil {
		if errors.Is(err, cdb.ErrDoesNotExist) {
			return cerr.NewAPIErrorResponse(c, http.StatusBadRequest, "Site specified in request does not exist", nil)
		}
		logger.Error().Err(err).Msg("error retrieving Site from DB")
		return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve Site due to DB error", nil)
	}

	// Verify site belongs to the org's Infrastructure Provider
	if site.InfrastructureProviderID != infrastructureProvider.ID {
		return cerr.NewAPIErrorResponse(c, http.StatusForbidden, "Site specified in request doesn't belong to current org's Provider", nil)
	}

	// Get tray ID from URL param
	trayStrID := c.Param("id")
	if _, err := uuid.Parse(trayStrID); err != nil {
		return cerr.NewAPIErrorResponse(c, http.StatusBadRequest, "Invalid Tray ID in URL", nil)
	}
	gth.tracerSpan.SetAttribute(handlerSpan, attribute.String("tray_id", trayStrID), logger)

	// Get the temporal client for the site
	stc, err := gth.scp.GetClientByID(site.ID)
	if err != nil {
		logger.Error().Err(err).Msg("failed to retrieve Temporal client for Site")
		return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve client for Site", nil)
	}

	// Build RLA request
	rlaRequest := &rlav1.GetComponentInfoByIDRequest{
		Id: &rlav1.UUID{Id: trayStrID},
	}

	// Execute workflow
	workflowOptions := tClient.StartWorkflowOptions{
		ID:                       fmt.Sprintf("tray-get-%s", trayStrID),
		WorkflowExecutionTimeout: common.WorkflowExecutionTimeout,
		TaskQueue:                queue.SiteTaskQueue,
		WorkflowIDReusePolicy:    temporalEnums.WORKFLOW_ID_REUSE_POLICY_ALLOW_DUPLICATE,
	}

	ctx, cancel := context.WithTimeout(ctx, common.WorkflowContextTimeout)
	defer cancel()

	we, err := stc.ExecuteWorkflow(ctx, workflowOptions, "GetTray", rlaRequest)
	if err != nil {
		logger.Error().Err(err).Msg("failed to execute GetTray workflow")
		return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, "Failed to get Tray details", nil)
	}

	// Get workflow result
	var rlaResponse rlav1.GetComponentInfoResponse
	err = we.Get(ctx, &rlaResponse)
	if err != nil {
		var timeoutErr *tp.TimeoutError
		if errors.As(err, &timeoutErr) || err == context.DeadlineExceeded || ctx.Err() != nil {
			return common.TerminateWorkflowOnTimeOut(c, logger, stc, fmt.Sprintf("tray-get-%s", trayStrID), err, "Tray", "GetTray")
		}
		logger.Error().Err(err).Msg("failed to get result from GetTray workflow")
		return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, "Failed to get Tray details", nil)
	}

	// Convert to API model
	apiTray := model.NewAPITray(rlaResponse.GetComponent())
	if apiTray == nil {
		return cerr.NewAPIErrorResponse(c, http.StatusNotFound, "Tray not found", nil)
	}

	logger.Info().Msg("finishing API handler")

	return c.JSON(http.StatusOK, apiTray)
}

// ~~~~~ GetAll Trays Handler ~~~~~ //

// GetAllTrayHandler is the API Handler for getting all Trays
type GetAllTrayHandler struct {
	dbSession  *cdb.Session
	tc         tClient.Client
	scp        *sc.ClientPool
	cfg        *config.Config
	tracerSpan *sutil.TracerSpan
}

// NewGetAllTrayHandler initializes and returns a new handler for getting all Trays
func NewGetAllTrayHandler(dbSession *cdb.Session, tc tClient.Client, scp *sc.ClientPool, cfg *config.Config) GetAllTrayHandler {
	return GetAllTrayHandler{
		dbSession:  dbSession,
		tc:         tc,
		scp:        scp,
		cfg:        cfg,
		tracerSpan: sutil.NewTracerSpan(),
	}
}

// Handle godoc
// @Summary Get all Trays
// @Description Get all Trays with optional filters
// @Tags tray
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param org path string true "Name of NGC organization"
// @Param siteId query string true "ID of the Site"
// @Param rackId query string false "Filter by Rack ID"
// @Param rackName query string false "Filter by Rack name"
// @Param type query string false "Filter by tray type (compute, switch, powershelf)"
// @Param componentId query string false "Filter by component ID (use repeated params for multiple values)"
// @Param id query string false "Filter by tray UUID (use repeated params for multiple values)"
// @Param orderBy query string false "Order by field (e.g. name_ASC, manufacturer_DESC)"
// @Param pageNumber query int false "Page number (1-based)"
// @Param pageSize query int false "Page size"
// @Success 200 {array} model.APITray
// @Router /v2/org/{org}/carbide/tray [get]
func (gath GetAllTrayHandler) Handle(c echo.Context) error {
	org, dbUser, ctx, logger, handlerSpan := common.SetupHandler("Tray", "GetAll", c, gath.tracerSpan)
	if handlerSpan != nil {
		defer handlerSpan.End()
	}

	if apiErr := common.ValidateQueryParams(c.QueryParams(), getAllTrayAllowedParams); apiErr != nil {
		return cerr.NewAPIErrorResponse(c, apiErr.Code, apiErr.Message, nil)
	}

	// Is DB user missing?
	if dbUser == nil {
		logger.Error().Msg("invalid User object found in request context")
		return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve current user", nil)
	}

	// Validate org membership
	ok, err := auth.ValidateOrgMembership(dbUser, org)
	if !ok {
		if err != nil {
			logger.Error().Err(err).Msg("error validating org membership for User in request")
		} else {
			logger.Warn().Msg("could not validate org membership for user, access denied")
		}
		return cerr.NewAPIErrorResponse(c, http.StatusForbidden, fmt.Sprintf("Failed to validate membership for org: %s", org), nil)
	}

	// Validate role, only Provider Admins are allowed to access Tray data
	ok = auth.ValidateUserRoles(dbUser, org, nil, auth.ProviderAdminRole)
	if !ok {
		logger.Warn().Msg("user does not have Provider Admin role, access denied")
		return cerr.NewAPIErrorResponse(c, http.StatusForbidden, "User does not have Provider Admin role with org", nil)
	}

	// Get Infrastructure Provider for org
	infrastructureProvider, err := common.GetInfrastructureProviderForOrg(ctx, nil, gath.dbSession, org)
	if err != nil {
		logger.Warn().Err(err).Msg("error getting infrastructure provider for org")
		return cerr.NewAPIErrorResponse(c, http.StatusBadRequest, "Failed to retrieve Infrastructure Provider for org", nil)
	}

	// Validate siteId is provided
	siteStrID := c.QueryParam("siteId")
	if siteStrID == "" {
		return cerr.NewAPIErrorResponse(c, http.StatusBadRequest, "siteId query parameter is required", nil)
	}

	// Retrieve the Site from the DB
	site, err := common.GetSiteFromIDString(ctx, nil, siteStrID, gath.dbSession)
	if err != nil {
		if errors.Is(err, cdb.ErrDoesNotExist) {
			return cerr.NewAPIErrorResponse(c, http.StatusBadRequest, "Site specified in request does not exist", nil)
		}
		logger.Error().Err(err).Msg("error retrieving Site from DB")
		return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve Site due to DB error", nil)
	}

	// Verify site belongs to the org's Infrastructure Provider
	if site.InfrastructureProviderID != infrastructureProvider.ID {
		return cerr.NewAPIErrorResponse(c, http.StatusForbidden, "Site specified in request doesn't belong to current org's Provider", nil)
	}

	// Build and validate tray request from query params
	apiRequest := model.APITrayGetAllRequest{}
	apiRequest.FromQueryParams(c.QueryParams())
	if verr := apiRequest.Validate(); verr != nil {
		logger.Warn().Err(verr).Msg("invalid tray request parameters")
		return cerr.NewAPIErrorResponse(c, http.StatusBadRequest, "Failed to validate request data", verr)
	}

	// Validate pagination request (orderBy, pageNumber, pageSize)
	pageRequest := pagination.PageRequest{}
	err = c.Bind(&pageRequest)
	if err != nil {
		logger.Warn().Err(err).Msg("error binding pagination request data into API model")
		return cerr.NewAPIErrorResponse(c, http.StatusBadRequest, "Failed to parse request pagination data", nil)
	}
	err = pageRequest.Validate(slices.Collect(maps.Keys(model.TrayOrderByFieldMap)))
	if err != nil {
		logger.Warn().Err(err).Msg("error validating pagination request data")
		return cerr.NewAPIErrorResponse(c, http.StatusBadRequest, "Failed to validate pagination request data", err)
	}

	// Get the temporal client for the site
	stc, err := gath.scp.GetClientByID(site.ID)
	if err != nil {
		logger.Error().Err(err).Msg("failed to retrieve Temporal client for Site")
		return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve client for Site", nil)
	}

	// Build RLA request from validated API request
	rlaRequest := apiRequest.ToProto()

	// Set order and pagination on RLA request
	var orderBy *rlav1.OrderBy
	if pageRequest.OrderBy != nil {
		orderBy = model.GetProtoTrayOrderByFromQueryParam(pageRequest.OrderBy.Field, strings.ToUpper(pageRequest.OrderBy.Order))
	}
	rlaRequest.OrderBy = orderBy
	if pageRequest.Offset != nil && pageRequest.Limit != nil {
		rlaRequest.Pagination = &rlav1.Pagination{
			Offset: int32(*pageRequest.Offset),
			Limit:  int32(*pageRequest.Limit),
		}
	}

	// Execute workflow
	workflowOptions := tClient.StartWorkflowOptions{
		ID:                       fmt.Sprintf("tray-get-all-%s", common.QueryParamHash(c)),
		WorkflowExecutionTimeout: common.WorkflowExecutionTimeout,
		TaskQueue:                queue.SiteTaskQueue,
		WorkflowIDReusePolicy:    temporalEnums.WORKFLOW_ID_REUSE_POLICY_ALLOW_DUPLICATE,
	}

	ctx, cancel := context.WithTimeout(ctx, common.WorkflowContextTimeout)
	defer cancel()

	we, err := stc.ExecuteWorkflow(ctx, workflowOptions, "GetTrays", rlaRequest)
	if err != nil {
		logger.Error().Err(err).Msg("failed to execute GetTrays workflow")
		return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, "Failed to get Trays", nil)
	}

	// Get workflow result
	var rlaResponse rlav1.GetComponentsResponse
	err = we.Get(ctx, &rlaResponse)
	if err != nil {
		var timeoutErr *tp.TimeoutError
		if errors.As(err, &timeoutErr) || err == context.DeadlineExceeded || ctx.Err() != nil {
			return common.TerminateWorkflowOnTimeOut(c, logger, stc, fmt.Sprintf("tray-get-all-%s", common.QueryParamHash(c)), err, "Tray", "GetTrays")
		}
		logger.Error().Err(err).Msg("failed to get result from GetTrays workflow")
		return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, "Failed to get Trays", nil)
	}

	apiTrays := make([]*model.APITray, 0, len(rlaResponse.GetComponents()))
	for _, comp := range rlaResponse.GetComponents() {
		apiTray := model.NewAPITray(comp)
		if apiTray != nil {
			apiTrays = append(apiTrays, apiTray)
		}
	}

	// Set pagination response header
	total := int(rlaResponse.GetTotal())
	pageResponse := pagination.NewPageResponse(*pageRequest.PageNumber, *pageRequest.PageSize, total, pageRequest.OrderByStr)
	pageHeader, err := json.Marshal(pageResponse)
	if err != nil {
		logger.Error().Err(err).Msg("error marshaling pagination response")
		return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, "Failed to create pagination response", nil)
	}
	c.Response().Header().Set(pagination.ResponseHeaderName, string(pageHeader))

	logger.Info().Int("count", len(apiTrays)).Int("Total", total).Msg("finishing API handler")

	return c.JSON(http.StatusOK, apiTrays)
}

// ~~~~~ Power Control Tray Handler ~~~~~ //

// PowerControlTrayHandler is the API Handler for power controlling a single Tray by ID
type PowerControlTrayHandler struct {
	dbSession  *cdb.Session
	tc         tClient.Client
	scp        *sc.ClientPool
	cfg        *config.Config
	tracerSpan *sutil.TracerSpan
}

// NewPowerControlTrayHandler initializes and returns a new handler for power controlling a Tray
func NewPowerControlTrayHandler(dbSession *cdb.Session, tc tClient.Client, scp *sc.ClientPool, cfg *config.Config) PowerControlTrayHandler {
	return PowerControlTrayHandler{
		dbSession:  dbSession,
		tc:         tc,
		scp:        scp,
		cfg:        cfg,
		tracerSpan: sutil.NewTracerSpan(),
	}
}

// Handle godoc
// @Summary Power control a Tray
// @Description Power control a single Tray by ID (on, off, cycle, forceoff, forcecycle)
// @Tags tray
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param org path string true "Name of NGC organization"
// @Param id path string true "ID of Tray"
// @Param siteId query string true "ID of the Site"
// @Param body body model.APIPowerControlRequest true "Power control request"
// @Success 200 {object} model.APIPowerControlResponse
// @Router /v2/org/{org}/carbide/tray/{id}/power [patch]
func (pcth PowerControlTrayHandler) Handle(c echo.Context) error {
	org, dbUser, ctx, logger, handlerSpan := common.SetupHandler("Tray", "PowerControl", c, pcth.tracerSpan)
	if handlerSpan != nil {
		defer handlerSpan.End()
	}

	if apiErr := common.ValidateQueryParams(c.QueryParams(), powerControlTrayAllowedParams); apiErr != nil {
		return cerr.NewAPIErrorResponse(c, apiErr.Code, apiErr.Message, nil)
	}

	// Is DB user missing?
	if dbUser == nil {
		logger.Error().Msg("invalid User object found in request context")
		return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve current user", nil)
	}

	// Validate org membership
	ok, err := auth.ValidateOrgMembership(dbUser, org)
	if !ok {
		if err != nil {
			logger.Error().Err(err).Msg("error validating org membership for User in request")
		} else {
			logger.Warn().Msg("could not validate org membership for user, access denied")
		}
		return cerr.NewAPIErrorResponse(c, http.StatusForbidden, fmt.Sprintf("Failed to validate membership for org: %s", org), nil)
	}

	// Validate role, only Provider Admins are allowed to power control Tray
	ok = auth.ValidateUserRoles(dbUser, org, nil, auth.ProviderAdminRole)
	if !ok {
		logger.Warn().Msg("user does not have Provider Admin role, access denied")
		return cerr.NewAPIErrorResponse(c, http.StatusForbidden, "User does not have Provider Admin role with org", nil)
	}

	// Get Infrastructure Provider for org
	infrastructureProvider, err := common.GetInfrastructureProviderForOrg(ctx, nil, pcth.dbSession, org)
	if err != nil {
		logger.Warn().Err(err).Msg("error getting infrastructure provider for org")
		return cerr.NewAPIErrorResponse(c, http.StatusBadRequest, "Failed to retrieve Infrastructure Provider for org", nil)
	}

	// Validate siteId is provided
	siteStrID := c.QueryParam("siteId")
	if siteStrID == "" {
		return cerr.NewAPIErrorResponse(c, http.StatusBadRequest, "siteId query parameter is required", nil)
	}

	// Retrieve the Site from the DB
	site, err := common.GetSiteFromIDString(ctx, nil, siteStrID, pcth.dbSession)
	if err != nil {
		if errors.Is(err, cdb.ErrDoesNotExist) {
			return cerr.NewAPIErrorResponse(c, http.StatusBadRequest, "Site specified in request does not exist", nil)
		}
		logger.Error().Err(err).Msg("error retrieving Site from DB")
		return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve Site due to DB error", nil)
	}

	// Verify site belongs to the org's Infrastructure Provider
	if site.InfrastructureProviderID != infrastructureProvider.ID {
		return cerr.NewAPIErrorResponse(c, http.StatusForbidden, "Site specified in request doesn't belong to current org's Provider", nil)
	}

	// Get tray ID from URL param
	trayStrID := c.Param("id")
	if _, err := uuid.Parse(trayStrID); err != nil {
		return cerr.NewAPIErrorResponse(c, http.StatusBadRequest, "Invalid Tray ID in URL", nil)
	}
	pcth.tracerSpan.SetAttribute(handlerSpan, attribute.String("tray_id", trayStrID), logger)

	// Parse and validate request body
	apiRequest := model.APIPowerControlRequest{}
	if err := c.Bind(&apiRequest); err != nil {
		return cerr.NewAPIErrorResponse(c, http.StatusBadRequest, "Failed to parse request data", nil)
	}
	if verr := apiRequest.Validate(); verr != nil {
		return cerr.NewAPIErrorResponse(c, http.StatusBadRequest, "Failed to validate request data", verr)
	}

	// Get the temporal client for the site
	stc, err := pcth.scp.GetClientByID(site.ID)
	if err != nil {
		logger.Error().Err(err).Msg("failed to retrieve Temporal client for Site")
		return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve client for Site", nil)
	}

	// Build TargetSpec for single tray by ID
	targetSpec := &rlav1.OperationTargetSpec{
		Targets: &rlav1.OperationTargetSpec_Components{
			Components: &rlav1.ComponentTargets{
				Targets: []*rlav1.ComponentTarget{
					{
						Identifier: &rlav1.ComponentTarget_Id{
							Id: &rlav1.UUID{Id: trayStrID},
						},
					},
				},
			},
		},
	}

	return executePowerControlWorkflow(ctx, c, logger, stc, targetSpec, apiRequest.State,
		fmt.Sprintf("tray-power-%s-%s", apiRequest.State, trayStrID), "Tray")
}

// ~~~~~ Power Control Trays (Batch) Handler ~~~~~ //

// PowerControlTrayBatchHandler is the API Handler for power controlling Trays with optional filters
type PowerControlTrayBatchHandler struct {
	dbSession  *cdb.Session
	tc         tClient.Client
	scp        *sc.ClientPool
	cfg        *config.Config
	tracerSpan *sutil.TracerSpan
}

// NewPowerControlTrayBatchHandler initializes and returns a new handler for batch power controlling Trays
func NewPowerControlTrayBatchHandler(dbSession *cdb.Session, tc tClient.Client, scp *sc.ClientPool, cfg *config.Config) PowerControlTrayBatchHandler {
	return PowerControlTrayBatchHandler{
		dbSession:  dbSession,
		tc:         tc,
		scp:        scp,
		cfg:        cfg,
		tracerSpan: sutil.NewTracerSpan(),
	}
}

// Handle godoc
// @Summary Power control Trays
// @Description Power control Trays with optional filters (on, off, cycle, forceoff, forcecycle). If no filter is specified, targets all trays in the Site.
// @Tags tray
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param org path string true "Name of NGC organization"
// @Param siteId query string true "ID of the Site"
// @Param rackId query string false "Filter by Rack ID"
// @Param rackName query string false "Filter by Rack name"
// @Param type query string false "Filter by tray type (compute, switch, powershelf)"
// @Param componentId query string false "Filter by component ID (use repeated params for multiple values)"
// @Param id query string false "Filter by tray UUID (use repeated params for multiple values)"
// @Param body body model.APIPowerControlRequest true "Power control request"
// @Success 200 {object} model.APIPowerControlResponse
// @Router /v2/org/{org}/carbide/tray/power [patch]
func (pctbh PowerControlTrayBatchHandler) Handle(c echo.Context) error {
	org, dbUser, ctx, logger, handlerSpan := common.SetupHandler("Tray", "PowerControlBatch", c, pctbh.tracerSpan)
	if handlerSpan != nil {
		defer handlerSpan.End()
	}

	if apiErr := common.ValidateQueryParams(c.QueryParams(), powerControlTrayBatchAllowedParams); apiErr != nil {
		return cerr.NewAPIErrorResponse(c, apiErr.Code, apiErr.Message, nil)
	}

	// Is DB user missing?
	if dbUser == nil {
		logger.Error().Msg("invalid User object found in request context")
		return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve current user", nil)
	}

	// Validate org membership
	ok, err := auth.ValidateOrgMembership(dbUser, org)
	if !ok {
		if err != nil {
			logger.Error().Err(err).Msg("error validating org membership for User in request")
		} else {
			logger.Warn().Msg("could not validate org membership for user, access denied")
		}
		return cerr.NewAPIErrorResponse(c, http.StatusForbidden, fmt.Sprintf("Failed to validate membership for org: %s", org), nil)
	}

	// Validate role, only Provider Admins are allowed to power control Tray
	ok = auth.ValidateUserRoles(dbUser, org, nil, auth.ProviderAdminRole)
	if !ok {
		logger.Warn().Msg("user does not have Provider Admin role, access denied")
		return cerr.NewAPIErrorResponse(c, http.StatusForbidden, "User does not have Provider Admin role with org", nil)
	}

	// Get Infrastructure Provider for org
	infrastructureProvider, err := common.GetInfrastructureProviderForOrg(ctx, nil, pctbh.dbSession, org)
	if err != nil {
		logger.Warn().Err(err).Msg("error getting infrastructure provider for org")
		return cerr.NewAPIErrorResponse(c, http.StatusBadRequest, "Failed to retrieve Infrastructure Provider for org", nil)
	}

	// Validate siteId is provided
	siteStrID := c.QueryParam("siteId")
	if siteStrID == "" {
		return cerr.NewAPIErrorResponse(c, http.StatusBadRequest, "siteId query parameter is required", nil)
	}

	// Retrieve the Site from the DB
	site, err := common.GetSiteFromIDString(ctx, nil, siteStrID, pctbh.dbSession)
	if err != nil {
		if errors.Is(err, cdb.ErrDoesNotExist) {
			return cerr.NewAPIErrorResponse(c, http.StatusBadRequest, "Site specified in request does not exist", nil)
		}
		logger.Error().Err(err).Msg("error retrieving Site from DB")
		return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve Site due to DB error", nil)
	}

	// Verify site belongs to the org's Infrastructure Provider
	if site.InfrastructureProviderID != infrastructureProvider.ID {
		return cerr.NewAPIErrorResponse(c, http.StatusForbidden, "Site specified in request doesn't belong to current org's Provider", nil)
	}

	// Build and validate tray filter request from query params
	filterRequest := model.APITrayGetAllRequest{}
	filterRequest.FromQueryParams(c.QueryParams())
	if verr := filterRequest.Validate(); verr != nil {
		logger.Warn().Err(verr).Msg("invalid tray filter parameters")
		return cerr.NewAPIErrorResponse(c, http.StatusBadRequest, "Failed to validate filter parameters", verr)
	}

	// Parse and validate request body
	apiRequest := model.APIPowerControlRequest{}
	if err := c.Bind(&apiRequest); err != nil {
		return cerr.NewAPIErrorResponse(c, http.StatusBadRequest, "Failed to parse request data", nil)
	}
	if verr := apiRequest.Validate(); verr != nil {
		return cerr.NewAPIErrorResponse(c, http.StatusBadRequest, "Failed to validate request data", verr)
	}

	// Get the temporal client for the site
	stc, err := pctbh.scp.GetClientByID(site.ID)
	if err != nil {
		logger.Error().Err(err).Msg("failed to retrieve Temporal client for Site")
		return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve client for Site", nil)
	}

	// Build TargetSpec from tray filters (reuses the same logic as GetAll Trays)
	targetSpec := filterRequest.ToProto().GetTargetSpec()

	return executePowerControlWorkflow(ctx, c, logger, stc, targetSpec, apiRequest.State,
		fmt.Sprintf("tray-power-batch-%s-%s", apiRequest.State, common.QueryParamHash(c)), "Tray")
}

// ~~~~~ Firmware Upgrade Tray Handler ~~~~~ //

// FirmwareUpgradeTrayHandler is the API Handler for upgrading firmware on a single Tray by ID
type FirmwareUpgradeTrayHandler struct {
	dbSession  *cdb.Session
	tc         tClient.Client
	scp        *sc.ClientPool
	cfg        *config.Config
	tracerSpan *sutil.TracerSpan
}

// NewFirmwareUpgradeTrayHandler initializes and returns a new handler for firmware upgrading a Tray
func NewFirmwareUpgradeTrayHandler(dbSession *cdb.Session, tc tClient.Client, scp *sc.ClientPool, cfg *config.Config) FirmwareUpgradeTrayHandler {
	return FirmwareUpgradeTrayHandler{
		dbSession:  dbSession,
		tc:         tc,
		scp:        scp,
		cfg:        cfg,
		tracerSpan: sutil.NewTracerSpan(),
	}
}

// Handle godoc
// @Summary Firmware upgrade a Tray
// @Description Upgrade firmware on a Tray identified by UUID. Version is optional; omit to upgrade to the latest available.
// @Tags tray
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param org path string true "Name of NGC organization"
// @Param id path string true "UUID of the Tray"
// @Param siteId query string true "ID of the Site"
// @Param body body model.APIFirmwareUpgradeRequest true "Firmware upgrade request"
// @Success 200 {object} model.APIFirmwareUpgradeResponse
// @Router /v2/org/{org}/carbide/tray/{id}/firmware [patch]
func (futh FirmwareUpgradeTrayHandler) Handle(c echo.Context) error {
	org, dbUser, ctx, logger, handlerSpan := common.SetupHandler("Tray", "FirmwareUpgrade", c, futh.tracerSpan)
	if handlerSpan != nil {
		defer handlerSpan.End()
	}

	if apiErr := common.ValidateQueryParams(c.QueryParams(), firmwareUpgradeTrayAllowedParams); apiErr != nil {
		return cerr.NewAPIErrorResponse(c, apiErr.Code, apiErr.Message, nil)
	}

	if dbUser == nil {
		logger.Error().Msg("invalid User object found in request context")
		return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve current user", nil)
	}

	ok, err := auth.ValidateOrgMembership(dbUser, org)
	if !ok {
		if err != nil {
			logger.Error().Err(err).Msg("error validating org membership for User in request")
		} else {
			logger.Warn().Msg("could not validate org membership for user, access denied")
		}
		return cerr.NewAPIErrorResponse(c, http.StatusForbidden, fmt.Sprintf("Failed to validate membership for org: %s", org), nil)
	}

	ok = auth.ValidateUserRoles(dbUser, org, nil, auth.ProviderAdminRole)
	if !ok {
		logger.Warn().Msg("user does not have Provider Admin role, access denied")
		return cerr.NewAPIErrorResponse(c, http.StatusForbidden, "User does not have Provider Admin role with org", nil)
	}

	infrastructureProvider, err := common.GetInfrastructureProviderForOrg(ctx, nil, futh.dbSession, org)
	if err != nil {
		logger.Warn().Err(err).Msg("error getting infrastructure provider for org")
		return cerr.NewAPIErrorResponse(c, http.StatusBadRequest, "Failed to retrieve Infrastructure Provider for org", nil)
	}

	siteStrID := c.QueryParam("siteId")
	if siteStrID == "" {
		return cerr.NewAPIErrorResponse(c, http.StatusBadRequest, "siteId query parameter is required", nil)
	}

	site, err := common.GetSiteFromIDString(ctx, nil, siteStrID, futh.dbSession)
	if err != nil {
		if errors.Is(err, cdb.ErrDoesNotExist) {
			return cerr.NewAPIErrorResponse(c, http.StatusBadRequest, "Site specified in request does not exist", nil)
		}
		logger.Error().Err(err).Msg("error retrieving Site from DB")
		return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve Site due to DB error", nil)
	}

	if site.InfrastructureProviderID != infrastructureProvider.ID {
		return cerr.NewAPIErrorResponse(c, http.StatusForbidden, "Site specified in request doesn't belong to current org's Provider", nil)
	}

	trayStrID := c.Param("id")
	if _, err := uuid.Parse(trayStrID); err != nil {
		return cerr.NewAPIErrorResponse(c, http.StatusBadRequest, "Invalid Tray ID in URL", nil)
	}
	futh.tracerSpan.SetAttribute(handlerSpan, attribute.String("tray_id", trayStrID), logger)

	apiRequest := model.APIFirmwareUpgradeRequest{}
	if err := c.Bind(&apiRequest); err != nil {
		return cerr.NewAPIErrorResponse(c, http.StatusBadRequest, "Failed to parse request data", nil)
	}

	stc, err := futh.scp.GetClientByID(site.ID)
	if err != nil {
		logger.Error().Err(err).Msg("failed to retrieve Temporal client for Site")
		return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve client for Site", nil)
	}

	targetSpec := &rlav1.OperationTargetSpec{
		Targets: &rlav1.OperationTargetSpec_Components{
			Components: &rlav1.ComponentTargets{
				Targets: []*rlav1.ComponentTarget{
					{
						Identifier: &rlav1.ComponentTarget_Id{
							Id: &rlav1.UUID{Id: trayStrID},
						},
					},
				},
			},
		},
	}

	return executeFirmwareUpgradeWorkflow(ctx, c, logger, stc, targetSpec, apiRequest.Version,
		fmt.Sprintf("tray-fw-upgrade-%s", trayStrID), "Tray")
}

// ~~~~~ Firmware Upgrade Trays (Batch) Handler ~~~~~ //

// FirmwareUpgradeTrayBatchHandler is the API Handler for firmware upgrading Trays with optional filters
type FirmwareUpgradeTrayBatchHandler struct {
	dbSession  *cdb.Session
	tc         tClient.Client
	scp        *sc.ClientPool
	cfg        *config.Config
	tracerSpan *sutil.TracerSpan
}

// NewFirmwareUpgradeTrayBatchHandler initializes and returns a new handler for batch firmware upgrading Trays
func NewFirmwareUpgradeTrayBatchHandler(dbSession *cdb.Session, tc tClient.Client, scp *sc.ClientPool, cfg *config.Config) FirmwareUpgradeTrayBatchHandler {
	return FirmwareUpgradeTrayBatchHandler{
		dbSession:  dbSession,
		tc:         tc,
		scp:        scp,
		cfg:        cfg,
		tracerSpan: sutil.NewTracerSpan(),
	}
}

// Handle godoc
// @Summary Firmware upgrade Trays
// @Description Upgrade firmware on Trays with optional filters. Version is optional; omit to upgrade to the latest available. If no filter is specified, targets all trays in the Site.
// @Tags tray
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param org path string true "Name of NGC organization"
// @Param siteId query string true "ID of the Site"
// @Param rackId query string false "Filter by Rack ID"
// @Param rackName query string false "Filter by Rack name"
// @Param type query string false "Filter by tray type (compute, switch, powershelf)"
// @Param componentId query string false "Filter by component ID (use repeated params for multiple values)"
// @Param id query string false "Filter by tray UUID (use repeated params for multiple values)"
// @Param body body model.APIFirmwareUpgradeRequest true "Firmware upgrade request"
// @Success 200 {object} model.APIFirmwareUpgradeResponse
// @Router /v2/org/{org}/carbide/tray/firmware [patch]
func (futbh FirmwareUpgradeTrayBatchHandler) Handle(c echo.Context) error {
	org, dbUser, ctx, logger, handlerSpan := common.SetupHandler("Tray", "FirmwareUpgradeBatch", c, futbh.tracerSpan)
	if handlerSpan != nil {
		defer handlerSpan.End()
	}

	if apiErr := common.ValidateQueryParams(c.QueryParams(), firmwareUpgradeTrayBatchAllowedParams); apiErr != nil {
		return cerr.NewAPIErrorResponse(c, apiErr.Code, apiErr.Message, nil)
	}

	if dbUser == nil {
		logger.Error().Msg("invalid User object found in request context")
		return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve current user", nil)
	}

	ok, err := auth.ValidateOrgMembership(dbUser, org)
	if !ok {
		if err != nil {
			logger.Error().Err(err).Msg("error validating org membership for User in request")
		} else {
			logger.Warn().Msg("could not validate org membership for user, access denied")
		}
		return cerr.NewAPIErrorResponse(c, http.StatusForbidden, fmt.Sprintf("Failed to validate membership for org: %s", org), nil)
	}

	ok = auth.ValidateUserRoles(dbUser, org, nil, auth.ProviderAdminRole)
	if !ok {
		logger.Warn().Msg("user does not have Provider Admin role, access denied")
		return cerr.NewAPIErrorResponse(c, http.StatusForbidden, "User does not have Provider Admin role with org", nil)
	}

	infrastructureProvider, err := common.GetInfrastructureProviderForOrg(ctx, nil, futbh.dbSession, org)
	if err != nil {
		logger.Warn().Err(err).Msg("error getting infrastructure provider for org")
		return cerr.NewAPIErrorResponse(c, http.StatusBadRequest, "Failed to retrieve Infrastructure Provider for org", nil)
	}

	siteStrID := c.QueryParam("siteId")
	if siteStrID == "" {
		return cerr.NewAPIErrorResponse(c, http.StatusBadRequest, "siteId query parameter is required", nil)
	}

	site, err := common.GetSiteFromIDString(ctx, nil, siteStrID, futbh.dbSession)
	if err != nil {
		if errors.Is(err, cdb.ErrDoesNotExist) {
			return cerr.NewAPIErrorResponse(c, http.StatusBadRequest, "Site specified in request does not exist", nil)
		}
		logger.Error().Err(err).Msg("error retrieving Site from DB")
		return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve Site due to DB error", nil)
	}

	if site.InfrastructureProviderID != infrastructureProvider.ID {
		return cerr.NewAPIErrorResponse(c, http.StatusForbidden, "Site specified in request doesn't belong to current org's Provider", nil)
	}

	filterRequest := model.APITrayGetAllRequest{}
	filterRequest.FromQueryParams(c.QueryParams())
	if verr := filterRequest.Validate(); verr != nil {
		logger.Warn().Err(verr).Msg("invalid tray filter parameters")
		return cerr.NewAPIErrorResponse(c, http.StatusBadRequest, "Failed to validate filter parameters", verr)
	}

	apiRequest := model.APIFirmwareUpgradeRequest{}
	if err := c.Bind(&apiRequest); err != nil {
		return cerr.NewAPIErrorResponse(c, http.StatusBadRequest, "Failed to parse request data", nil)
	}

	stc, err := futbh.scp.GetClientByID(site.ID)
	if err != nil {
		logger.Error().Err(err).Msg("failed to retrieve Temporal client for Site")
		return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve client for Site", nil)
	}

	targetSpec := filterRequest.ToProto().GetTargetSpec()

	return executeFirmwareUpgradeWorkflow(ctx, c, logger, stc, targetSpec, apiRequest.Version,
		fmt.Sprintf("tray-fw-upgrade-batch-%s", common.QueryParamHash(c)), "Tray")
}
