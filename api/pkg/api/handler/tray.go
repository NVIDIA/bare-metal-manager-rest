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
	getTrayAllowedParams      = []string{"siteId"}
	getAllTrayAllowedParams    = []string{"siteId", "rackId", "rackName", "type", "componentId", "id", "pageNumber", "pageSize", "orderBy"}
	validateTrayAllowedParams  = []string{"siteId"}
	validateTraysAllowedParams = []string{"siteId", "rackId", "rackName", "name", "manufacturer", "type", "componentId"}
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

// ~~~~~ Validate Tray Handler ~~~~~ //

// ValidateTrayHandler is the API Handler for validating a single Tray's components
type ValidateTrayHandler struct {
	dbSession  *cdb.Session
	tc         tClient.Client
	scp        *sc.ClientPool
	cfg        *config.Config
	tracerSpan *sutil.TracerSpan
}

// NewValidateTrayHandler initializes and returns a new handler for validating a Tray
func NewValidateTrayHandler(dbSession *cdb.Session, tc tClient.Client, scp *sc.ClientPool, cfg *config.Config) ValidateTrayHandler {
	return ValidateTrayHandler{
		dbSession:  dbSession,
		tc:         tc,
		scp:        scp,
		cfg:        cfg,
		tracerSpan: sutil.NewTracerSpan(),
	}
}

// Handle godoc
// @Summary Validate a Tray
// @Description Validate a Tray by comparing expected vs actual state via RLA
// @Tags tray
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param org path string true "Name of NGC organization"
// @Param id path string true "ID of the Tray"
// @Param siteId query string true "ID of the Site"
// @Success 200 {object} model.APIRackValidationResult
// @Router /v2/org/{org}/forge/tray/{id}/validation [get]
func (vth ValidateTrayHandler) Handle(c echo.Context) error {
	org, dbUser, ctx, logger, handlerSpan := common.SetupHandler("Tray", "Validate", c, vth.tracerSpan)
	if handlerSpan != nil {
		defer handlerSpan.End()
	}

	if apiErr := common.ValidateQueryParams(c.QueryParams(), validateTrayAllowedParams); apiErr != nil {
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

	infrastructureProvider, err := common.GetInfrastructureProviderForOrg(ctx, nil, vth.dbSession, org)
	if err != nil {
		logger.Warn().Err(err).Msg("error getting infrastructure provider for org")
		return cerr.NewAPIErrorResponse(c, http.StatusBadRequest, "Failed to retrieve Infrastructure Provider for org", nil)
	}

	trayStrID := c.Param("id")
	if _, err := uuid.Parse(trayStrID); err != nil {
		return cerr.NewAPIErrorResponse(c, http.StatusBadRequest, "Invalid Tray ID in URL", nil)
	}
	vth.tracerSpan.SetAttribute(handlerSpan, attribute.String("tray_id", trayStrID), logger)

	siteStrID := c.QueryParam("siteId")
	if siteStrID == "" {
		return cerr.NewAPIErrorResponse(c, http.StatusBadRequest, "siteId query parameter is required", nil)
	}

	site, err := common.GetSiteFromIDString(ctx, nil, siteStrID, vth.dbSession)
	if err != nil {
		if errors.Is(err, cdb.ErrDoesNotExist) {
			return cerr.NewAPIErrorResponse(c, http.StatusBadRequest, "Site specified in request does not exist", nil)
		}
		logger.Error().Err(err).Msg("error retrieving Site from DB")
		return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve Site specified in request due to DB error", nil)
	}

	if site.InfrastructureProviderID != infrastructureProvider.ID {
		return cerr.NewAPIErrorResponse(c, http.StatusForbidden, "Site specified in request doesn't belong to current org's Provider", nil)
	}

	stc, err := vth.scp.GetClientByID(site.ID)
	if err != nil {
		logger.Error().Err(err).Msg("failed to retrieve Temporal client for Site")
		return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve client for Site", nil)
	}

	rlaRequest := &rlav1.ValidateComponentsRequest{
		TargetSpec: &rlav1.OperationTargetSpec{
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
		},
	}

	workflowOptions := tClient.StartWorkflowOptions{
		ID:                       fmt.Sprintf("tray-validate-%s", trayStrID),
		WorkflowIDReusePolicy:    temporalEnums.WORKFLOW_ID_REUSE_POLICY_ALLOW_DUPLICATE,
		WorkflowIDConflictPolicy: temporalEnums.WORKFLOW_ID_CONFLICT_POLICY_USE_EXISTING,
		WorkflowExecutionTimeout: common.WorkflowExecutionTimeout,
		TaskQueue:                queue.SiteTaskQueue,
	}

	ctx, cancel := context.WithTimeout(ctx, common.WorkflowContextTimeout)
	defer cancel()

	we, err := stc.ExecuteWorkflow(ctx, workflowOptions, "ValidateRackComponents", rlaRequest)
	if err != nil {
		logger.Error().Err(err).Msg("failed to execute ValidateComponents workflow")
		return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, "Failed to validate Tray", nil)
	}

	var rlaResponse rlav1.ValidateComponentsResponse
	err = we.Get(ctx, &rlaResponse)
	if err != nil {
		var timeoutErr *tp.TimeoutError
		if errors.As(err, &timeoutErr) || err == context.DeadlineExceeded || ctx.Err() != nil {
			return common.TerminateWorkflowOnTimeOut(c, logger, stc, fmt.Sprintf("tray-validate-%s", trayStrID), err, "Tray", "ValidateRackComponents")
		}
		logger.Error().Err(err).Msg("failed to get result from ValidateComponents workflow")
		return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, "Failed to validate Tray", nil)
	}

	apiResult := model.NewAPIRackValidationResult(&rlaResponse)

	logger.Info().Int32("TotalDiffs", rlaResponse.GetTotalDiffs()).Msg("finishing API handler")

	return c.JSON(http.StatusOK, apiResult)
}

// ~~~~~ Validate Trays Handler ~~~~~ //

// ValidateTraysHandler is the API Handler for validating Trays with optional filters.
// If no filter is specified, validates all trays in the Site.
type ValidateTraysHandler struct {
	dbSession  *cdb.Session
	tc         tClient.Client
	scp        *sc.ClientPool
	cfg        *config.Config
	tracerSpan *sutil.TracerSpan
}

// NewValidateTraysHandler initializes and returns a new handler for validating Trays
func NewValidateTraysHandler(dbSession *cdb.Session, tc tClient.Client, scp *sc.ClientPool, cfg *config.Config) ValidateTraysHandler {
	return ValidateTraysHandler{
		dbSession:  dbSession,
		tc:         tc,
		scp:        scp,
		cfg:        cfg,
		tracerSpan: sutil.NewTracerSpan(),
	}
}

// Handle godoc
// @Summary Validate Trays
// @Description Validate Tray components by comparing expected vs actual state via RLA. If no filter is specified, validates all trays in the Site. Use rackId/rackName to scope to a specific rack, and name/manufacturer/type to filter by tray attributes.
// @Tags tray
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param org path string true "Name of NGC organization"
// @Param siteId query string true "ID of the Site"
// @Param rackId query string false "Scope to a specific Rack by ID (mutually exclusive with rackName)"
// @Param rackName query string false "Scope to a specific Rack by name (mutually exclusive with rackId)"
// @Param name query string false "Filter trays by name"
// @Param manufacturer query string false "Filter trays by manufacturer"
// @Param type query string false "Filter trays by type (compute, switch, powershelf)"
// @Param componentId query string false "Filter by external component ID (requires type; mutually exclusive with rackId/rackName; use repeated params for multiple values)"
// @Success 200 {object} model.APIRackValidationResult
// @Router /v2/org/{org}/forge/tray/validation [get]
func (vtsh ValidateTraysHandler) Handle(c echo.Context) error {
	org, dbUser, ctx, logger, handlerSpan := common.SetupHandler("Tray", "ValidateTrays", c, vtsh.tracerSpan)
	if handlerSpan != nil {
		defer handlerSpan.End()
	}

	if apiErr := common.ValidateQueryParams(c.QueryParams(), validateTraysAllowedParams); apiErr != nil {
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

	infrastructureProvider, err := common.GetInfrastructureProviderForOrg(ctx, nil, vtsh.dbSession, org)
	if err != nil {
		logger.Warn().Err(err).Msg("error getting infrastructure provider for org")
		return cerr.NewAPIErrorResponse(c, http.StatusBadRequest, "Failed to retrieve Infrastructure Provider for org", nil)
	}

	siteStrID := c.QueryParam("siteId")
	if siteStrID == "" {
		return cerr.NewAPIErrorResponse(c, http.StatusBadRequest, "siteId query parameter is required", nil)
	}

	site, err := common.GetSiteFromIDString(ctx, nil, siteStrID, vtsh.dbSession)
	if err != nil {
		if errors.Is(err, cdb.ErrDoesNotExist) {
			return cerr.NewAPIErrorResponse(c, http.StatusBadRequest, "Site specified in request does not exist", nil)
		}
		logger.Error().Err(err).Msg("error retrieving Site from DB")
		return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve Site specified in request due to DB error", nil)
	}

	if site.InfrastructureProviderID != infrastructureProvider.ID {
		return cerr.NewAPIErrorResponse(c, http.StatusForbidden, "Site specified in request doesn't belong to current org's Provider", nil)
	}

	stc, err := vtsh.scp.GetClientByID(site.ID)
	if err != nil {
		logger.Error().Err(err).Msg("failed to retrieve Temporal client for Site")
		return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve client for Site", nil)
	}

	rackID := c.QueryParam("rackId")
	rackName := c.QueryParam("rackName")
	componentIDs := c.QueryParams()["componentId"]
	componentType := c.QueryParam("type")

	if rackID != "" && rackName != "" {
		return cerr.NewAPIErrorResponse(c, http.StatusBadRequest, "rackId and rackName are mutually exclusive", nil)
	}

	hasRackScope := rackID != "" || rackName != ""
	if hasRackScope && len(componentIDs) > 0 {
		return cerr.NewAPIErrorResponse(c, http.StatusBadRequest, "rackId/rackName and componentId are mutually exclusive", nil)
	}

	if len(componentIDs) > 0 && componentType == "" {
		return cerr.NewAPIErrorResponse(c, http.StatusBadRequest, "type is required when componentId is provided", nil)
	}

	var targetSpec *rlav1.OperationTargetSpec
	if rackID != "" {
		if _, parseErr := uuid.Parse(rackID); parseErr != nil {
			return cerr.NewAPIErrorResponse(c, http.StatusBadRequest, "Invalid rackId: must be a valid UUID", nil)
		}
		targetSpec = &rlav1.OperationTargetSpec{
			Targets: &rlav1.OperationTargetSpec_Racks{
				Racks: &rlav1.RackTargets{
					Targets: []*rlav1.RackTarget{
						{Identifier: &rlav1.RackTarget_Id{Id: &rlav1.UUID{Id: rackID}}},
					},
				},
			},
		}
	} else if rackName != "" {
		targetSpec = &rlav1.OperationTargetSpec{
			Targets: &rlav1.OperationTargetSpec_Racks{
				Racks: &rlav1.RackTargets{
					Targets: []*rlav1.RackTarget{
						{Identifier: &rlav1.RackTarget_Name{Name: rackName}},
					},
				},
			},
		}
	} else if len(componentIDs) > 0 {
		protoName, ok := model.APIToProtoComponentTypeName[componentType]
		if !ok {
			return cerr.NewAPIErrorResponse(c, http.StatusBadRequest, fmt.Sprintf("Invalid type: must be one of %v", slices.Collect(maps.Keys(model.APIToProtoComponentTypeName))), nil)
		}
		protoType := rlav1.ComponentType(rlav1.ComponentType_value[protoName])
		targets := make([]*rlav1.ComponentTarget, 0, len(componentIDs))
		for _, cid := range componentIDs {
			targets = append(targets, &rlav1.ComponentTarget{
				Identifier: &rlav1.ComponentTarget_External{
					External: &rlav1.ExternalRef{
						Type: protoType,
						Id:   cid,
					},
				},
			})
		}
		targetSpec = &rlav1.OperationTargetSpec{
			Targets: &rlav1.OperationTargetSpec_Components{
				Components: &rlav1.ComponentTargets{
					Targets: targets,
				},
			},
		}
	}

	var filters []*rlav1.Filter
	qParams := c.QueryParams()
	for field := range model.TrayFilterFieldMap {
		if value := qParams.Get(field); value != "" {
			if f := model.GetProtoTrayFilterFromQueryParam(field, value); f != nil {
				filters = append(filters, f)
			}
		}
	}

	rlaRequest := &rlav1.ValidateComponentsRequest{
		TargetSpec: targetSpec,
		Filters:    filters,
	}

	workflowID := fmt.Sprintf("tray-validate-all-%s", common.QueryParamHash(c))

	workflowOptions := tClient.StartWorkflowOptions{
		ID:                       workflowID,
		WorkflowIDReusePolicy:    temporalEnums.WORKFLOW_ID_REUSE_POLICY_ALLOW_DUPLICATE,
		WorkflowIDConflictPolicy: temporalEnums.WORKFLOW_ID_CONFLICT_POLICY_USE_EXISTING,
		WorkflowExecutionTimeout: common.WorkflowExecutionTimeout,
		TaskQueue:                queue.SiteTaskQueue,
	}

	ctx, cancel := context.WithTimeout(ctx, common.WorkflowContextTimeout)
	defer cancel()

	we, err := stc.ExecuteWorkflow(ctx, workflowOptions, "ValidateRackComponents", rlaRequest)
	if err != nil {
		logger.Error().Err(err).Msg("failed to execute ValidateComponents workflow")
		return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, "Failed to validate Trays", nil)
	}

	var rlaResponse rlav1.ValidateComponentsResponse
	err = we.Get(ctx, &rlaResponse)
	if err != nil {
		var timeoutErr *tp.TimeoutError
		if errors.As(err, &timeoutErr) || err == context.DeadlineExceeded || ctx.Err() != nil {
			return common.TerminateWorkflowOnTimeOut(c, logger, stc, workflowID, err, "Tray", "ValidateRackComponents")
		}
		logger.Error().Err(err).Msg("failed to get result from ValidateComponents workflow")
		return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, "Failed to validate Trays", nil)
	}

	apiResult := model.NewAPIRackValidationResult(&rlaResponse)

	logger.Info().Int("FilterCount", len(filters)).Int32("TotalDiffs", rlaResponse.GetTotalDiffs()).Msg("finishing API handler")

	return c.JSON(http.StatusOK, apiResult)
}
