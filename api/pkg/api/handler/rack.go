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

	"github.com/labstack/echo/v4"
	"go.opentelemetry.io/otel/attribute"
	temporalEnums "go.temporal.io/api/enums/v1"
	tClient "go.temporal.io/sdk/client"
	tp "go.temporal.io/sdk/temporal"

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
)

// ~~~~~ Get Rack Handler ~~~~~ //

// GetRackHandler is the API Handler for getting a Rack by ID
type GetRackHandler struct {
	dbSession  *cdb.Session
	tc         tClient.Client
	scp        *sc.ClientPool
	cfg        *config.Config
	tracerSpan *sutil.TracerSpan
}

// NewGetRackHandler initializes and returns a new handler for getting a Rack
func NewGetRackHandler(dbSession *cdb.Session, tc tClient.Client, scp *sc.ClientPool, cfg *config.Config) GetRackHandler {
	return GetRackHandler{
		dbSession:  dbSession,
		tc:         tc,
		scp:        scp,
		cfg:        cfg,
		tracerSpan: sutil.NewTracerSpan(),
	}
}

// Handle godoc
// @Summary Get a Rack
// @Description Get a Rack by ID from RLA
// @Tags rack
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param org path string true "Name of NGC organization"
// @Param id path string true "ID of Rack"
// @Param siteId query string true "ID of the Site"
// @Param includeComponents query boolean false "Include rack components in response"
// @Success 200 {object} model.APIRack
// @Router /v2/org/{org}/carbide/rack/{id} [get]
func (grh GetRackHandler) Handle(c echo.Context) error {
	org, dbUser, ctx, logger, handlerSpan := common.SetupHandler("Rack", "Get", c, grh.tracerSpan)
	if handlerSpan != nil {
		defer handlerSpan.End()
	}

	var apiRequest model.APIRackGetRequest
	if err := common.ValidateKnownQueryParams(c.QueryParams(), apiRequest); err != nil {
		return cerr.NewAPIErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
	}
	if err := c.Bind(&apiRequest); err != nil {
		return cerr.NewAPIErrorResponse(c, http.StatusBadRequest, "Failed to parse request data", nil)
	}
	if err := apiRequest.Validate(); err != nil {
		return cerr.NewAPIErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
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

	// Validate role, only Provider Admins are allowed to access Rack data
	ok = auth.ValidateUserRoles(dbUser, org, nil, auth.ProviderAdminRole)
	if !ok {
		logger.Warn().Msg("user does not have Provider Admin role, access denied")
		return cerr.NewAPIErrorResponse(c, http.StatusForbidden, "User does not have Provider Admin role with org", nil)
	}

	// Get Infrastructure Provider for org
	infrastructureProvider, err := common.GetInfrastructureProviderForOrg(ctx, nil, grh.dbSession, org)
	if err != nil {
		logger.Warn().Err(err).Msg("error getting infrastructure provider for org")
		return cerr.NewAPIErrorResponse(c, http.StatusBadRequest, "Failed to retrieve Infrastructure Provider for org", nil)
	}

	// Get rack ID from URL param
	rackStrID := c.Param("id")
	grh.tracerSpan.SetAttribute(handlerSpan, attribute.String("rack_id", rackStrID), logger)

	// Validate the site
	site, err := common.GetSiteFromIDString(ctx, nil, apiRequest.SiteID, grh.dbSession)
	if err != nil {
		if errors.Is(err, cdb.ErrDoesNotExist) {
			return cerr.NewAPIErrorResponse(c, http.StatusBadRequest, "Site specified in request does not exist", nil)
		}
		logger.Error().Err(err).Msg("error retrieving Site from DB")
		return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve Site specified in request due to DB error", nil)
	}

	// Verify site belongs to the org's Infrastructure Provider
	if site.InfrastructureProviderID != infrastructureProvider.ID {
		return cerr.NewAPIErrorResponse(c, http.StatusForbidden, "Site specified in request doesn't belong to current org's Provider", nil)
	}

	// Get the temporal client for the site
	stc, err := grh.scp.GetClientByID(site.ID)
	if err != nil {
		logger.Error().Err(err).Msg("failed to retrieve Temporal client for Site")
		return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve client for Site", nil)
	}

	// Build RLA request
	rlaRequest := &rlav1.GetRackInfoByIDRequest{
		Id:             &rlav1.UUID{Id: rackStrID},
		WithComponents: apiRequest.IncludeComponents,
	}

	// Execute workflow
	workflowOptions := tClient.StartWorkflowOptions{
		ID:                       fmt.Sprintf("rack-get-%s", rackStrID),
		WorkflowIDReusePolicy:    temporalEnums.WORKFLOW_ID_REUSE_POLICY_ALLOW_DUPLICATE,
		WorkflowIDConflictPolicy: temporalEnums.WORKFLOW_ID_CONFLICT_POLICY_USE_EXISTING,
		WorkflowExecutionTimeout: common.WorkflowExecutionTimeout,
		TaskQueue:                queue.SiteTaskQueue,
	}

	ctx, cancel := context.WithTimeout(ctx, common.WorkflowContextTimeout)
	defer cancel()

	we, err := stc.ExecuteWorkflow(ctx, workflowOptions, "GetRack", rlaRequest)
	if err != nil {
		logger.Error().Err(err).Msg("failed to execute GetRack workflow")
		return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, "Failed to get Rack details", nil)
	}

	// Get workflow result
	var rlaResponse rlav1.GetRackInfoResponse
	err = we.Get(ctx, &rlaResponse)
	if err != nil {
		var timeoutErr *tp.TimeoutError
		if errors.As(err, &timeoutErr) || err == context.DeadlineExceeded || ctx.Err() != nil {
			return common.TerminateWorkflowOnTimeOut(c, logger, stc, fmt.Sprintf("rack-get-%s", rackStrID), err, "Rack", "GetRack")
		}
		logger.Error().Err(err).Msg("failed to get result from GetRack workflow")
		return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, "Failed to get Rack details", nil)
	}

	// Convert to API model
	protoRack := rlaResponse.GetRack()
	apiRack := model.NewAPIRack(protoRack, apiRequest.IncludeComponents)
	if apiRack == nil {
		return cerr.NewAPIErrorResponse(c, http.StatusNotFound, "Rack not found", nil)
	}

	logger.Info().Msg("finishing API handler")

	return c.JSON(http.StatusOK, apiRack)
}

// ~~~~~ GetAll Racks Handler ~~~~~ //

// GetAllRackHandler is the API Handler for getting all Racks
type GetAllRackHandler struct {
	dbSession  *cdb.Session
	tc         tClient.Client
	scp        *sc.ClientPool
	cfg        *config.Config
	tracerSpan *sutil.TracerSpan
}

// NewGetAllRackHandler initializes and returns a new handler for getting all Racks
func NewGetAllRackHandler(dbSession *cdb.Session, tc tClient.Client, scp *sc.ClientPool, cfg *config.Config) GetAllRackHandler {
	return GetAllRackHandler{
		dbSession:  dbSession,
		tc:         tc,
		scp:        scp,
		cfg:        cfg,
		tracerSpan: sutil.NewTracerSpan(),
	}
}

// Handle godoc
// @Summary Get all Racks
// @Description Get all Racks from RLA with optional filters
// @Tags rack
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param org path string true "Name of NGC organization"
// @Param siteId query string true "ID of the Site"
// @Param includeComponents query boolean false "Include rack components in response"
// @Param name query string false "Filter by rack name"
// @Param manufacturer query string false "Filter by manufacturer"
// @Param pageNumber query integer false "Page number of results returned"
// @Param pageSize query integer false "Number of results per page"
// @Param orderBy query string false "Order by field"
// @Success 200 {array} model.APIRack
// @Router /v2/org/{org}/carbide/rack [get]
func (garh GetAllRackHandler) Handle(c echo.Context) error {
	org, dbUser, ctx, logger, handlerSpan := common.SetupHandler("Rack", "GetAll", c, garh.tracerSpan)
	if handlerSpan != nil {
		defer handlerSpan.End()
	}

	var apiRequest model.APIRackGetAllRequest
	if err := common.ValidateKnownQueryParams(c.QueryParams(), apiRequest); err != nil {
		return cerr.NewAPIErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
	}
	if err := c.Bind(&apiRequest); err != nil {
		return cerr.NewAPIErrorResponse(c, http.StatusBadRequest, "Failed to parse request data", nil)
	}
	if err := apiRequest.Validate(); err != nil {
		return cerr.NewAPIErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
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

	// Validate role, only Provider Admins are allowed to access Rack data
	ok = auth.ValidateUserRoles(dbUser, org, nil, auth.ProviderAdminRole)
	if !ok {
		logger.Warn().Msg("user does not have Provider Admin role, access denied")
		return cerr.NewAPIErrorResponse(c, http.StatusForbidden, "User does not have Provider Admin role with org", nil)
	}

	// Get Infrastructure Provider for org
	infrastructureProvider, err := common.GetInfrastructureProviderForOrg(ctx, nil, garh.dbSession, org)
	if err != nil {
		logger.Warn().Err(err).Msg("error getting infrastructure provider for org")
		return cerr.NewAPIErrorResponse(c, http.StatusBadRequest, "Failed to retrieve Infrastructure Provider for org", nil)
	}

	// Validate the site
	site, err := common.GetSiteFromIDString(ctx, nil, apiRequest.SiteID, garh.dbSession)
	if err != nil {
		if errors.Is(err, cdb.ErrDoesNotExist) {
			return cerr.NewAPIErrorResponse(c, http.StatusBadRequest, "Site specified in request does not exist", nil)
		}
		logger.Error().Err(err).Msg("error retrieving Site from DB")
		return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve Site specified in request due to DB error", nil)
	}

	// Verify site belongs to the org's Infrastructure Provider
	if site.InfrastructureProviderID != infrastructureProvider.ID {
		return cerr.NewAPIErrorResponse(c, http.StatusForbidden, "Site specified in request doesn't belong to current org's Provider", nil)
	}

	// Validate pagination request
	pageRequest := pagination.PageRequest{}
	err = c.Bind(&pageRequest)
	if err != nil {
		logger.Warn().Err(err).Msg("error binding pagination request data into API model")
		return cerr.NewAPIErrorResponse(c, http.StatusBadRequest, "Failed to parse request pagination data", nil)
	}

	// Validate pagination attributes
	err = pageRequest.Validate(slices.Collect(maps.Keys(model.RackOrderByFieldMap)))
	if err != nil {
		logger.Warn().Err(err).Msg("error validating pagination request data")
		return cerr.NewAPIErrorResponse(c, http.StatusBadRequest, "Failed to validate pagination request data", err)
	}

	// Build OrderBy from pagination
	var orderBy *rlav1.OrderBy
	if pageRequest.OrderBy != nil {
		orderBy = model.GetProtoRackOrderByFromQueryParam(pageRequest.OrderBy.Field, strings.ToUpper(pageRequest.OrderBy.Order))
	}

	// Build Pagination
	var paginationProto *rlav1.Pagination
	if pageRequest.Offset != nil && pageRequest.Limit != nil {
		paginationProto = &rlav1.Pagination{
			Offset: int32(*pageRequest.Offset),
			Limit:  int32(*pageRequest.Limit),
		}
	}

	// Get the temporal client for the site
	stc, err := garh.scp.GetClientByID(site.ID)
	if err != nil {
		logger.Error().Err(err).Msg("failed to retrieve Temporal client for Site")
		return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve client for Site", nil)
	}

	// Build RLA request from validated params
	rlaRequest := &rlav1.GetListOfRacksRequest{
		Filters:        apiRequest.ToFilters(),
		WithComponents: apiRequest.IncludeComponents,
		Pagination:     paginationProto,
		OrderBy:        orderBy,
	}

	workflowID := fmt.Sprintf("rack-get-all-%s", common.QueryParamHash(apiRequest.QueryValues()))

	// Execute workflow
	workflowOptions := tClient.StartWorkflowOptions{
		ID:                       workflowID,
		WorkflowIDReusePolicy:    temporalEnums.WORKFLOW_ID_REUSE_POLICY_ALLOW_DUPLICATE,
		WorkflowIDConflictPolicy: temporalEnums.WORKFLOW_ID_CONFLICT_POLICY_USE_EXISTING,
		WorkflowExecutionTimeout: common.WorkflowExecutionTimeout,
		TaskQueue:                queue.SiteTaskQueue,
	}

	ctx, cancel := context.WithTimeout(ctx, common.WorkflowContextTimeout)
	defer cancel()

	we, err := stc.ExecuteWorkflow(ctx, workflowOptions, "GetRacks", rlaRequest)
	if err != nil {
		logger.Error().Err(err).Msg("failed to execute GetRacks workflow")
		return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, "Failed to get Racks", nil)
	}

	// Get workflow result
	var rlaResponse rlav1.GetListOfRacksResponse
	err = we.Get(ctx, &rlaResponse)
	if err != nil {
		var timeoutErr *tp.TimeoutError
		if errors.As(err, &timeoutErr) || err == context.DeadlineExceeded || ctx.Err() != nil {
			return common.TerminateWorkflowOnTimeOut(c, logger, stc, workflowID, err, "Rack", "GetRacks")
		}
		logger.Error().Err(err).Msg("failed to get result from GetRacks workflow")
		return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, "Failed to get Racks", nil)
	}

	// Convert to API model
	apiRacks := make([]*model.APIRack, 0, len(rlaResponse.GetRacks()))
	for _, rack := range rlaResponse.GetRacks() {
		apiRacks = append(apiRacks, model.NewAPIRack(rack, apiRequest.IncludeComponents))
	}

	// Create pagination response header
	total := int(rlaResponse.GetTotal())
	pageResponse := pagination.NewPageResponse(*pageRequest.PageNumber, *pageRequest.PageSize, total, pageRequest.OrderByStr)
	pageHeader, err := json.Marshal(pageResponse)
	if err != nil {
		logger.Error().Err(err).Msg("error marshaling pagination response")
		return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, "Failed to create pagination response", nil)
	}
	c.Response().Header().Set(pagination.ResponseHeaderName, string(pageHeader))

	logger.Info().Int("Count", len(apiRacks)).Int("Total", total).Msg("finishing API handler")

	return c.JSON(http.StatusOK, apiRacks)
}

// ~~~~~ Validate Rack Handler ~~~~~ //

// ValidateRackHandler is the API Handler for validating a Rack's components
type ValidateRackHandler struct {
	dbSession  *cdb.Session
	tc         tClient.Client
	scp        *sc.ClientPool
	cfg        *config.Config
	tracerSpan *sutil.TracerSpan
}

// NewValidateRackHandler initializes and returns a new handler for validating a Rack
func NewValidateRackHandler(dbSession *cdb.Session, tc tClient.Client, scp *sc.ClientPool, cfg *config.Config) ValidateRackHandler {
	return ValidateRackHandler{
		dbSession:  dbSession,
		tc:         tc,
		scp:        scp,
		cfg:        cfg,
		tracerSpan: sutil.NewTracerSpan(),
	}
}

// Handle godoc
// @Summary Validate a Rack
// @Description Validate a Rack's components by comparing expected vs actual state via RLA
// @Tags rack
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param org path string true "Name of NGC organization"
// @Param id path string true "ID of the Rack"
// @Param siteId query string true "ID of the Site"
// @Success 200 {object} model.APIRackValidationResult
// @Router /v2/org/{org}/carbide/rack/{id}/validation [get]
func (vrh ValidateRackHandler) Handle(c echo.Context) error {
	org, dbUser, ctx, logger, handlerSpan := common.SetupHandler("Rack", "Validate", c, vrh.tracerSpan)
	if handlerSpan != nil {
		defer handlerSpan.End()
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

	// Validate role, only Provider Admins are allowed to access Rack data
	ok = auth.ValidateUserRoles(dbUser, org, nil, auth.ProviderAdminRole)
	if !ok {
		logger.Warn().Msg("user does not have Provider Admin role, access denied")
		return cerr.NewAPIErrorResponse(c, http.StatusForbidden, "User does not have Provider Admin role with org", nil)
	}

	// Get Infrastructure Provider for org
	infrastructureProvider, err := common.GetInfrastructureProviderForOrg(ctx, nil, vrh.dbSession, org)
	if err != nil {
		logger.Warn().Err(err).Msg("error getting infrastructure provider for org")
		return cerr.NewAPIErrorResponse(c, http.StatusBadRequest, "Failed to retrieve Infrastructure Provider for org", nil)
	}

	// Get rack ID from URL param
	rackStrID := c.Param("id")
	vrh.tracerSpan.SetAttribute(handlerSpan, attribute.String("rack_id", rackStrID), logger)

	// Get site ID from query param (required)
	siteStrID := c.QueryParam("siteId")
	if siteStrID == "" {
		return cerr.NewAPIErrorResponse(c, http.StatusBadRequest, "siteId query parameter is required", nil)
	}

	// Validate the site
	site, err := common.GetSiteFromIDString(ctx, nil, siteStrID, vrh.dbSession)
	if err != nil {
		if errors.Is(err, cdb.ErrDoesNotExist) {
			return cerr.NewAPIErrorResponse(c, http.StatusBadRequest, "Site specified in request does not exist", nil)
		}
		logger.Error().Err(err).Msg("error retrieving Site from DB")
		return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve Site specified in request due to DB error", nil)
	}

	// Verify site belongs to the org's Infrastructure Provider
	if site.InfrastructureProviderID != infrastructureProvider.ID {
		return cerr.NewAPIErrorResponse(c, http.StatusForbidden, "Site specified in request doesn't belong to current org's Provider", nil)
	}

	// Get the temporal client for the site
	stc, err := vrh.scp.GetClientByID(site.ID)
	if err != nil {
		logger.Error().Err(err).Msg("failed to retrieve Temporal client for Site")
		return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve client for Site", nil)
	}

	// Build RLA request - target the specific rack by ID
	rlaRequest := &rlav1.ValidateComponentsRequest{
		TargetSpec: &rlav1.OperationTargetSpec{
			Targets: &rlav1.OperationTargetSpec_Racks{
				Racks: &rlav1.RackTargets{
					Targets: []*rlav1.RackTarget{
						{
							Identifier: &rlav1.RackTarget_Id{
								Id: &rlav1.UUID{Id: rackStrID},
							},
						},
					},
				},
			},
		},
	}

	// Execute workflow
	workflowOptions := tClient.StartWorkflowOptions{
		ID:                       fmt.Sprintf("rack-validate-%s", rackStrID),
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
		return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, "Failed to validate Rack", nil)
	}

	// Get workflow result
	var rlaResponse rlav1.ValidateComponentsResponse
	err = we.Get(ctx, &rlaResponse)
	if err != nil {
		var timeoutErr *tp.TimeoutError
		if errors.As(err, &timeoutErr) || err == context.DeadlineExceeded || ctx.Err() != nil {
			return common.TerminateWorkflowOnTimeOut(c, logger, stc, fmt.Sprintf("rack-validate-%s", rackStrID), err, "Rack", "ValidateRackComponents")
		}
		logger.Error().Err(err).Msg("failed to get result from ValidateComponents workflow")
		return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, "Failed to validate Rack", nil)
	}

	// Convert to API model
	apiResult := model.NewAPIRackValidationResult(&rlaResponse)

	logger.Info().Int32("TotalDiffs", rlaResponse.GetTotalDiffs()).Msg("finishing API handler")

	return c.JSON(http.StatusOK, apiResult)
}

// ~~~~~ Validate Racks Handler ~~~~~ //

// ValidateRacksHandler is the API Handler for validating Racks with optional filters.
// If no filter is specified, validates all racks in the Site.
type ValidateRacksHandler struct {
	dbSession  *cdb.Session
	tc         tClient.Client
	scp        *sc.ClientPool
	cfg        *config.Config
	tracerSpan *sutil.TracerSpan
}

// NewValidateRacksHandler initializes and returns a new handler for validating Racks
func NewValidateRacksHandler(dbSession *cdb.Session, tc tClient.Client, scp *sc.ClientPool, cfg *config.Config) ValidateRacksHandler {
	return ValidateRacksHandler{
		dbSession:  dbSession,
		tc:         tc,
		scp:        scp,
		cfg:        cfg,
		tracerSpan: sutil.NewTracerSpan(),
	}
}

// Handle godoc
// @Summary Validate Racks
// @Description Validate Rack components by comparing expected vs actual state via RLA. If no filter is specified, validates all racks in the Site.
// @Tags rack
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param org path string true "Name of NGC organization"
// @Param siteId query string true "ID of the Site"
// @Param name query string false "Filter racks by name"
// @Param manufacturer query string false "Filter racks by manufacturer"
// @Success 200 {object} model.APIRackValidationResult
// @Router /v2/org/{org}/carbide/rack/validation [get]
func (vrsh ValidateRacksHandler) Handle(c echo.Context) error {
	org, dbUser, ctx, logger, handlerSpan := common.SetupHandler("Rack", "ValidateRacks", c, vrsh.tracerSpan)
	if handlerSpan != nil {
		defer handlerSpan.End()
	}

	var apiRequest model.APIRackValidateAllRequest
	if err := common.ValidateKnownQueryParams(c.QueryParams(), apiRequest); err != nil {
		return cerr.NewAPIErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
	}
	if err := c.Bind(&apiRequest); err != nil {
		return cerr.NewAPIErrorResponse(c, http.StatusBadRequest, "Failed to parse request data", nil)
	}
	if err := apiRequest.Validate(); err != nil {
		return cerr.NewAPIErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
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

	// Validate role, only Provider Admins are allowed to access Rack data
	ok = auth.ValidateUserRoles(dbUser, org, nil, auth.ProviderAdminRole)
	if !ok {
		logger.Warn().Msg("user does not have Provider Admin role, access denied")
		return cerr.NewAPIErrorResponse(c, http.StatusForbidden, "User does not have Provider Admin role with org", nil)
	}

	// Get Infrastructure Provider for org
	infrastructureProvider, err := common.GetInfrastructureProviderForOrg(ctx, nil, vrsh.dbSession, org)
	if err != nil {
		logger.Warn().Err(err).Msg("error getting infrastructure provider for org")
		return cerr.NewAPIErrorResponse(c, http.StatusBadRequest, "Failed to retrieve Infrastructure Provider for org", nil)
	}

	// Validate the site
	site, err := common.GetSiteFromIDString(ctx, nil, apiRequest.SiteID, vrsh.dbSession)
	if err != nil {
		if errors.Is(err, cdb.ErrDoesNotExist) {
			return cerr.NewAPIErrorResponse(c, http.StatusBadRequest, "Site specified in request does not exist", nil)
		}
		logger.Error().Err(err).Msg("error retrieving Site from DB")
		return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve Site specified in request due to DB error", nil)
	}

	// Verify site belongs to the org's Infrastructure Provider
	if site.InfrastructureProviderID != infrastructureProvider.ID {
		return cerr.NewAPIErrorResponse(c, http.StatusForbidden, "Site specified in request doesn't belong to current org's Provider", nil)
	}

	// Get the temporal client for the site
	stc, err := vrsh.scp.GetClientByID(site.ID)
	if err != nil {
		logger.Error().Err(err).Msg("failed to retrieve Temporal client for Site")
		return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve client for Site", nil)
	}

	rlaRequest := &rlav1.ValidateComponentsRequest{
		Filters: apiRequest.ToFilters(),
	}

	workflowID := fmt.Sprintf("rack-validate-all-%s", common.QueryParamHash(apiRequest.QueryValues()))

	// Execute workflow
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
		return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, "Failed to validate Racks", nil)
	}

	// Get workflow result
	var rlaResponse rlav1.ValidateComponentsResponse
	err = we.Get(ctx, &rlaResponse)
	if err != nil {
		var timeoutErr *tp.TimeoutError
		if errors.As(err, &timeoutErr) || err == context.DeadlineExceeded || ctx.Err() != nil {
			return common.TerminateWorkflowOnTimeOut(c, logger, stc, workflowID, err, "Rack", "ValidateRackComponents")
		}
		logger.Error().Err(err).Msg("failed to get result from ValidateComponents workflow")
		return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, "Failed to validate Racks", nil)
	}

	// Convert to API model
	apiResult := model.NewAPIRackValidationResult(&rlaResponse)

	logger.Info().Int32("TotalDiffs", rlaResponse.GetTotalDiffs()).Msg("finishing API handler")

	return c.JSON(http.StatusOK, apiResult)
}

// ~~~~~ Power Control Rack Handler ~~~~~ //

// PowerControlRackHandler is the API Handler for power controlling a single Rack by ID
type PowerControlRackHandler struct {
	dbSession  *cdb.Session
	tc         tClient.Client
	scp        *sc.ClientPool
	cfg        *config.Config
	tracerSpan *sutil.TracerSpan
}

// NewPowerControlRackHandler initializes and returns a new handler for power controlling a Rack
func NewPowerControlRackHandler(dbSession *cdb.Session, tc tClient.Client, scp *sc.ClientPool, cfg *config.Config) PowerControlRackHandler {
	return PowerControlRackHandler{
		dbSession:  dbSession,
		tc:         tc,
		scp:        scp,
		cfg:        cfg,
		tracerSpan: sutil.NewTracerSpan(),
	}
}

// Handle godoc
// @Summary Power control a Rack
// @Description Power control a single Rack by ID (on, off, cycle, forceoff, forcecycle)
// @Tags rack
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param org path string true "Name of NGC organization"
// @Param id path string true "ID of Rack"
// @Param siteId query string true "ID of the Site"
// @Param body body model.APIPowerControlRequest true "Power control request"
// @Success 200 {object} model.APIPowerControlResponse
// @Router /v2/org/{org}/carbide/rack/{id}/power [patch]
func (pcrh PowerControlRackHandler) Handle(c echo.Context) error {
	org, dbUser, ctx, logger, handlerSpan := common.SetupHandler("Rack", "PowerControl", c, pcrh.tracerSpan)
	if handlerSpan != nil {
		defer handlerSpan.End()
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

	// Validate role, only Provider Admins are allowed to power control Rack
	ok = auth.ValidateUserRoles(dbUser, org, nil, auth.ProviderAdminRole)
	if !ok {
		logger.Warn().Msg("user does not have Provider Admin role, access denied")
		return cerr.NewAPIErrorResponse(c, http.StatusForbidden, "User does not have Provider Admin role with org", nil)
	}

	// Get Infrastructure Provider for org
	infrastructureProvider, err := common.GetInfrastructureProviderForOrg(ctx, nil, pcrh.dbSession, org)
	if err != nil {
		logger.Warn().Err(err).Msg("error getting infrastructure provider for org")
		return cerr.NewAPIErrorResponse(c, http.StatusBadRequest, "Failed to retrieve Infrastructure Provider for org", nil)
	}

	// Get rack ID from URL param
	rackStrID := c.Param("id")
	pcrh.tracerSpan.SetAttribute(handlerSpan, attribute.String("rack_id", rackStrID), logger)

	// Get site ID from query param (required)
	siteStrID := c.QueryParam("siteId")
	if siteStrID == "" {
		return cerr.NewAPIErrorResponse(c, http.StatusBadRequest, "siteId query parameter is required", nil)
	}

	// Validate the site
	site, err := common.GetSiteFromIDString(ctx, nil, siteStrID, pcrh.dbSession)
	if err != nil {
		if errors.Is(err, cdb.ErrDoesNotExist) {
			return cerr.NewAPIErrorResponse(c, http.StatusBadRequest, "Site specified in request does not exist", nil)
		}
		logger.Error().Err(err).Msg("error retrieving Site from DB")
		return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve Site specified in request due to DB error", nil)
	}

	// Verify site belongs to the org's Infrastructure Provider
	if site.InfrastructureProviderID != infrastructureProvider.ID {
		return cerr.NewAPIErrorResponse(c, http.StatusForbidden, "Site specified in request doesn't belong to current org's Provider", nil)
	}

	// Parse and validate request body
	apiRequest := model.APIPowerControlRequest{}
	if err := c.Bind(&apiRequest); err != nil {
		return cerr.NewAPIErrorResponse(c, http.StatusBadRequest, "Failed to parse request data", nil)
	}
	verr := apiRequest.Validate()
	if verr != nil {
		logger.Warn().Err(verr).Msg("error validating power control request data")
		return cerr.NewAPIErrorResponse(c, http.StatusBadRequest, "Failed to validate power control request data", verr)
	}

	// Get the temporal client for the site
	stc, err := pcrh.scp.GetClientByID(site.ID)
	if err != nil {
		logger.Error().Err(err).Msg("failed to retrieve Temporal client for Site")
		return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve client for Site", nil)
	}

	// Build TargetSpec for single rack by ID
	targetSpec := &rlav1.OperationTargetSpec{
		Targets: &rlav1.OperationTargetSpec_Racks{
			Racks: &rlav1.RackTargets{
				Targets: []*rlav1.RackTarget{
					{
						Identifier: &rlav1.RackTarget_Id{
							Id: &rlav1.UUID{Id: rackStrID},
						},
					},
				},
			},
		},
	}

	return common.ExecutePowerControlWorkflow(ctx, c, logger, stc, targetSpec, apiRequest.State,
		fmt.Sprintf("rack-power-%s-%s", apiRequest.State, rackStrID), "Rack")
}

// ~~~~~ Power Control Racks (Batch) Handler ~~~~~ //

// PowerControlRackBatchHandler is the API Handler for power controlling Racks with optional filters
type PowerControlRackBatchHandler struct {
	dbSession  *cdb.Session
	tc         tClient.Client
	scp        *sc.ClientPool
	cfg        *config.Config
	tracerSpan *sutil.TracerSpan
}

// NewPowerControlRackBatchHandler initializes and returns a new handler for batch power controlling Racks
func NewPowerControlRackBatchHandler(dbSession *cdb.Session, tc tClient.Client, scp *sc.ClientPool, cfg *config.Config) PowerControlRackBatchHandler {
	return PowerControlRackBatchHandler{
		dbSession:  dbSession,
		tc:         tc,
		scp:        scp,
		cfg:        cfg,
		tracerSpan: sutil.NewTracerSpan(),
	}
}

// Handle godoc
// @Summary Power control Racks
// @Description Power control Racks with optional filters (on, off, cycle, forceoff, forcecycle). If no filter is specified, targets all racks in the Site.
// @Tags rack
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param org path string true "Name of NGC organization"
// @Param siteId query string true "ID of the Site"
// @Param name query string false "Filter racks by name (use repeated params for multiple values)"
// @Param body body model.APIPowerControlRequest true "Power control request"
// @Success 200 {object} model.APIPowerControlResponse
// @Router /v2/org/{org}/carbide/rack/power [patch]
func (pcrbh PowerControlRackBatchHandler) Handle(c echo.Context) error {
	org, dbUser, ctx, logger, handlerSpan := common.SetupHandler("Rack", "PowerControlBatch", c, pcrbh.tracerSpan)
	if handlerSpan != nil {
		defer handlerSpan.End()
	}

	var filterRequest model.APIRackPowerControlBatchRequest
	if err := common.ValidateKnownQueryParams(c.QueryParams(), filterRequest); err != nil {
		return cerr.NewAPIErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
	}
	if err := (&echo.DefaultBinder{}).BindQueryParams(c, &filterRequest); err != nil {
		return cerr.NewAPIErrorResponse(c, http.StatusBadRequest, "Failed to parse request data", nil)
	}
	if verr := filterRequest.Validate(); verr != nil {
		logger.Warn().Err(verr).Msg("error validating power control filter parameters")
		return cerr.NewAPIErrorResponse(c, http.StatusBadRequest, "Failed to validate power control filter parameters", verr)
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

	// Validate role, only Provider Admins are allowed to power control Rack
	ok = auth.ValidateUserRoles(dbUser, org, nil, auth.ProviderAdminRole)
	if !ok {
		logger.Warn().Msg("user does not have Provider Admin role, access denied")
		return cerr.NewAPIErrorResponse(c, http.StatusForbidden, "User does not have Provider Admin role with org", nil)
	}

	// Get Infrastructure Provider for org
	infrastructureProvider, err := common.GetInfrastructureProviderForOrg(ctx, nil, pcrbh.dbSession, org)
	if err != nil {
		logger.Warn().Err(err).Msg("error getting infrastructure provider for org")
		return cerr.NewAPIErrorResponse(c, http.StatusBadRequest, "Failed to retrieve Infrastructure Provider for org", nil)
	}

	// Validate the site
	site, err := common.GetSiteFromIDString(ctx, nil, filterRequest.SiteID, pcrbh.dbSession)
	if err != nil {
		if errors.Is(err, cdb.ErrDoesNotExist) {
			return cerr.NewAPIErrorResponse(c, http.StatusBadRequest, "Site specified in request does not exist", nil)
		}
		logger.Error().Err(err).Msg("error retrieving Site from DB")
		return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve Site specified in request due to DB error", nil)
	}

	// Verify site belongs to the org's Infrastructure Provider
	if site.InfrastructureProviderID != infrastructureProvider.ID {
		return cerr.NewAPIErrorResponse(c, http.StatusForbidden, "Site specified in request doesn't belong to current org's Provider", nil)
	}

	// Parse and validate request body
	apiRequest := model.APIPowerControlRequest{}
	if err := c.Bind(&apiRequest); err != nil {
		return cerr.NewAPIErrorResponse(c, http.StatusBadRequest, "Failed to parse request data", nil)
	}
	verr := apiRequest.Validate()
	if verr != nil {
		logger.Warn().Err(verr).Msg("error validating power control request data")
		return cerr.NewAPIErrorResponse(c, http.StatusBadRequest, "Failed to validate power control request data", verr)
	}

	// Get the temporal client for the site
	stc, err := pcrbh.scp.GetClientByID(site.ID)
	if err != nil {
		logger.Error().Err(err).Msg("failed to retrieve Temporal client for Site")
		return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve client for Site", nil)
	}

	// Build TargetSpec from rack filters
	targetSpec := filterRequest.ToTargetSpec()

	return common.ExecutePowerControlWorkflow(ctx, c, logger, stc, targetSpec, apiRequest.State,
		fmt.Sprintf("rack-power-batch-%s-%s", apiRequest.State, common.QueryParamHash(filterRequest.QueryValues())), "Rack")
}

// ~~~~~ Firmware Upgrade Rack Handler ~~~~~ //

// FirmwareUpgradeRackHandler is the API Handler for upgrading firmware on a single Rack by ID
type FirmwareUpgradeRackHandler struct {
	dbSession  *cdb.Session
	tc         tClient.Client
	scp        *sc.ClientPool
	cfg        *config.Config
	tracerSpan *sutil.TracerSpan
}

// NewFirmwareUpgradeRackHandler initializes and returns a new handler for firmware upgrading a Rack
func NewFirmwareUpgradeRackHandler(dbSession *cdb.Session, tc tClient.Client, scp *sc.ClientPool, cfg *config.Config) FirmwareUpgradeRackHandler {
	return FirmwareUpgradeRackHandler{
		dbSession:  dbSession,
		tc:         tc,
		scp:        scp,
		cfg:        cfg,
		tracerSpan: sutil.NewTracerSpan(),
	}
}

// Handle godoc
// @Summary Firmware upgrade a Rack
// @Description Upgrade firmware on a Rack identified by UUID. Version is optional; omit to upgrade to the latest available.
// @Tags rack
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param org path string true "Name of NGC organization"
// @Param id path string true "UUID of the Rack"
// @Param siteId query string true "ID of the Site"
// @Param body body model.APIFirmwareUpgradeRequest true "Firmware upgrade request"
// @Success 200 {object} model.APIFirmwareUpgradeResponse
// @Router /v2/org/{org}/carbide/rack/{id}/firmware [patch]
func (furh FirmwareUpgradeRackHandler) Handle(c echo.Context) error {
	org, dbUser, ctx, logger, handlerSpan := common.SetupHandler("Rack", "FirmwareUpgrade", c, furh.tracerSpan)
	if handlerSpan != nil {
		defer handlerSpan.End()
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

	// Validate role, only Provider Admins are allowed to firmware upgrade Rack
	ok = auth.ValidateUserRoles(dbUser, org, nil, auth.ProviderAdminRole)
	if !ok {
		logger.Warn().Msg("user does not have Provider Admin role, access denied")
		return cerr.NewAPIErrorResponse(c, http.StatusForbidden, "User does not have Provider Admin role with org", nil)
	}

	// Get Infrastructure Provider for org
	infrastructureProvider, err := common.GetInfrastructureProviderForOrg(ctx, nil, furh.dbSession, org)
	if err != nil {
		logger.Warn().Err(err).Msg("error getting infrastructure provider for org")
		return cerr.NewAPIErrorResponse(c, http.StatusBadRequest, "Failed to retrieve Infrastructure Provider for org", nil)
	}

	// Get rack ID from URL param
	rackStrID := c.Param("id")
	furh.tracerSpan.SetAttribute(handlerSpan, attribute.String("rack_id", rackStrID), logger)

	// Get site ID from query param (required)
	siteStrID := c.QueryParam("siteId")
	if siteStrID == "" {
		return cerr.NewAPIErrorResponse(c, http.StatusBadRequest, "siteId query parameter is required", nil)
	}

	// Validate the site
	site, err := common.GetSiteFromIDString(ctx, nil, siteStrID, furh.dbSession)
	if err != nil {
		if errors.Is(err, cdb.ErrDoesNotExist) {
			return cerr.NewAPIErrorResponse(c, http.StatusBadRequest, "Site specified in request does not exist", nil)
		}
		logger.Error().Err(err).Msg("error retrieving Site from DB")
		return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve Site specified in request due to DB error", nil)
	}

	// Verify site belongs to the org's Infrastructure Provider
	if site.InfrastructureProviderID != infrastructureProvider.ID {
		return cerr.NewAPIErrorResponse(c, http.StatusForbidden, "Site specified in request doesn't belong to current org's Provider", nil)
	}

	// Parse request body
	apiRequest := model.APIFirmwareUpgradeRequest{}
	if err := c.Bind(&apiRequest); err != nil {
		return cerr.NewAPIErrorResponse(c, http.StatusBadRequest, "Failed to parse request data", nil)
	}

	// Get the temporal client for the site
	stc, err := furh.scp.GetClientByID(site.ID)
	if err != nil {
		logger.Error().Err(err).Msg("failed to retrieve Temporal client for Site")
		return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve client for Site", nil)
	}

	targetSpec := &rlav1.OperationTargetSpec{
		Targets: &rlav1.OperationTargetSpec_Racks{
			Racks: &rlav1.RackTargets{
				Targets: []*rlav1.RackTarget{
					{
						Identifier: &rlav1.RackTarget_Id{
							Id: &rlav1.UUID{Id: rackStrID},
						},
					},
				},
			},
		},
	}

	return common.ExecuteFirmwareUpgradeWorkflow(ctx, c, logger, stc, targetSpec, apiRequest.Version,
		fmt.Sprintf("rack-fw-upgrade-%s", rackStrID), "Rack")
}

// ~~~~~ Firmware Upgrade Racks (Batch) Handler ~~~~~ //

// FirmwareUpgradeRackBatchHandler is the API Handler for firmware upgrading Racks with optional filters
type FirmwareUpgradeRackBatchHandler struct {
	dbSession  *cdb.Session
	tc         tClient.Client
	scp        *sc.ClientPool
	cfg        *config.Config
	tracerSpan *sutil.TracerSpan
}

// NewFirmwareUpgradeRackBatchHandler initializes and returns a new handler for batch firmware upgrading Racks
func NewFirmwareUpgradeRackBatchHandler(dbSession *cdb.Session, tc tClient.Client, scp *sc.ClientPool, cfg *config.Config) FirmwareUpgradeRackBatchHandler {
	return FirmwareUpgradeRackBatchHandler{
		dbSession:  dbSession,
		tc:         tc,
		scp:        scp,
		cfg:        cfg,
		tracerSpan: sutil.NewTracerSpan(),
	}
}

// Handle godoc
// @Summary Firmware upgrade Racks
// @Description Upgrade firmware on Racks with optional name filter. Version is optional; omit to upgrade to the latest available. If no filter is specified, targets all racks in the Site.
// @Tags rack
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param org path string true "Name of NGC organization"
// @Param siteId query string true "ID of the Site"
// @Param name query string false "Filter racks by name (use repeated params for multiple values)"
// @Param body body model.APIFirmwareUpgradeRequest true "Firmware upgrade request"
// @Success 200 {object} model.APIFirmwareUpgradeResponse
// @Router /v2/org/{org}/carbide/rack/firmware [patch]
func (furbh FirmwareUpgradeRackBatchHandler) Handle(c echo.Context) error {
	org, dbUser, ctx, logger, handlerSpan := common.SetupHandler("Rack", "FirmwareUpgradeBatch", c, furbh.tracerSpan)
	if handlerSpan != nil {
		defer handlerSpan.End()
	}

	var filterRequest model.APIRackFirmwareUpgradeBatchRequest
	if err := common.ValidateKnownQueryParams(c.QueryParams(), filterRequest); err != nil {
		return cerr.NewAPIErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
	}
	if err := (&echo.DefaultBinder{}).BindQueryParams(c, &filterRequest); err != nil {
		return cerr.NewAPIErrorResponse(c, http.StatusBadRequest, "Failed to parse request data", nil)
	}
	if verr := filterRequest.Validate(); verr != nil {
		logger.Warn().Err(verr).Msg("error validating firmware upgrade filter parameters")
		return cerr.NewAPIErrorResponse(c, http.StatusBadRequest, "Failed to validate firmware upgrade filter parameters", verr)
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

	// Validate role, only Provider Admins are allowed to firmware upgrade Rack
	ok = auth.ValidateUserRoles(dbUser, org, nil, auth.ProviderAdminRole)
	if !ok {
		logger.Warn().Msg("user does not have Provider Admin role, access denied")
		return cerr.NewAPIErrorResponse(c, http.StatusForbidden, "User does not have Provider Admin role with org", nil)
	}

	// Get Infrastructure Provider for org
	infrastructureProvider, err := common.GetInfrastructureProviderForOrg(ctx, nil, furbh.dbSession, org)
	if err != nil {
		logger.Warn().Err(err).Msg("error getting infrastructure provider for org")
		return cerr.NewAPIErrorResponse(c, http.StatusBadRequest, "Failed to retrieve Infrastructure Provider for org", nil)
	}

	// Validate the site
	site, err := common.GetSiteFromIDString(ctx, nil, filterRequest.SiteID, furbh.dbSession)
	if err != nil {
		if errors.Is(err, cdb.ErrDoesNotExist) {
			return cerr.NewAPIErrorResponse(c, http.StatusBadRequest, "Site specified in request does not exist", nil)
		}
		logger.Error().Err(err).Msg("error retrieving Site from DB")
		return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve Site specified in request due to DB error", nil)
	}

	// Verify site belongs to the org's Infrastructure Provider
	if site.InfrastructureProviderID != infrastructureProvider.ID {
		return cerr.NewAPIErrorResponse(c, http.StatusForbidden, "Site specified in request doesn't belong to current org's Provider", nil)
	}

	// Parse request body
	apiRequest := model.APIFirmwareUpgradeRequest{}
	if err := c.Bind(&apiRequest); err != nil {
		return cerr.NewAPIErrorResponse(c, http.StatusBadRequest, "Failed to parse request data", nil)
	}

	// Get the temporal client for the site
	stc, err := furbh.scp.GetClientByID(site.ID)
	if err != nil {
		logger.Error().Err(err).Msg("failed to retrieve Temporal client for Site")
		return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve client for Site", nil)
	}

	// Build TargetSpec from rack filters
	targetSpec := filterRequest.ToTargetSpec()

	return common.ExecuteFirmwareUpgradeWorkflow(ctx, c, logger, stc, targetSpec, apiRequest.Version,
		fmt.Sprintf("rack-fw-upgrade-batch-%s", common.QueryParamHash(filterRequest.QueryValues())), "Rack")
}
