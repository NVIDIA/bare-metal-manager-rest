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
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
	"go.opentelemetry.io/otel/attribute"
	tClient "go.temporal.io/sdk/client"

	"github.com/nvidia/carbide-rest/api/internal/config"
	"github.com/nvidia/carbide-rest/api/pkg/api/handler/util/common"
	"github.com/nvidia/carbide-rest/api/pkg/api/model"
	"github.com/nvidia/carbide-rest/api/pkg/api/pagination"
	sc "github.com/nvidia/carbide-rest/api/pkg/client/site"
	cerr "github.com/nvidia/carbide-rest/common/pkg/util"
	sutil "github.com/nvidia/carbide-rest/common/pkg/util"
	cdb "github.com/nvidia/carbide-rest/db/pkg/db"
	rlav1 "github.com/nvidia/carbide-rest/workflow-schema/rla/protobuf/v1"
	"github.com/nvidia/carbide-rest/workflow/pkg/queue"
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

	// Is DB user missing?
	if dbUser == nil {
		logger.Error().Msg("invalid User object found in request context")
		return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve current user", nil)
	}

	// ensure our user is a provider or privileged tenant for the org
	infrastructureProvider, tenant, apiError := common.IsProviderOrTenant(ctx, logger, grh.dbSession, org, dbUser, false, true)
	if apiError != nil {
		return cerr.NewAPIErrorResponse(c, apiError.Code, apiError.Message, apiError.Data)
	}

	// Get rack ID from URL param
	rackStrID := c.Param("id")
	grh.tracerSpan.SetAttribute(handlerSpan, attribute.String("rack_id", rackStrID), logger)

	// Get site ID from query param (required)
	siteStrID := c.QueryParam("siteId")
	if siteStrID == "" {
		return cerr.NewAPIErrorResponse(c, http.StatusBadRequest, "siteId query parameter is required", nil)
	}

	// Validate the site
	site, err := common.GetSiteFromIDString(ctx, nil, siteStrID, grh.dbSession)
	if err != nil {
		if errors.Is(err, cdb.ErrDoesNotExist) {
			return cerr.NewAPIErrorResponse(c, http.StatusBadRequest, "Site specified in request does not exist", nil)
		}
		logger.Error().Err(err).Msg("error retrieving Site from DB")
		return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve Site specified in request due to DB error", nil)
	}

	// Validate Provider/Tenant Site access
	hasAccess, apiError := common.ValidateProviderOrTenantSiteAccess(ctx, logger, grh.dbSession, site, infrastructureProvider, tenant)
	if apiError != nil {
		return cerr.NewAPIErrorResponse(c, apiError.Code, apiError.Message, apiError.Data)
	}

	if !hasAccess {
		return cerr.NewAPIErrorResponse(c, http.StatusForbidden, "Current org is not associated with the Site specified in query", nil)
	}

	// Check includeComponents query param (API uses includeComponents, RLA uses WithComponents)
	includeComponents := false
	if ic := c.QueryParam("includeComponents"); ic != "" {
		includeComponents, _ = strconv.ParseBool(ic)
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
		WithComponents: includeComponents,
	}

	// Execute workflow
	workflowOptions := tClient.StartWorkflowOptions{
		ID:                       fmt.Sprintf("GetRack-%s", rackStrID),
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
		logger.Error().Err(err).Msg("failed to get result from GetRack workflow")
		return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, "Failed to get Rack details", nil)
	}

	// Convert to API model
	apiRack := model.NewAPIRack(rlaResponse.GetRack(), includeComponents)

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
// @Param model query string false "Filter by model"
// @Param pageNumber query integer false "Page number of results returned"
// @Param pageSize query integer false "Number of results per page"
// @Param orderBy query string false "Order by field (NAME_ASC, NAME_DESC, MANUFACTURER_ASC, MANUFACTURER_DESC, MODEL_ASC, MODEL_DESC)"
// @Success 200 {array} model.APIRack
// @Router /v2/org/{org}/carbide/rack [get]
func (garh GetAllRackHandler) Handle(c echo.Context) error {
	org, dbUser, ctx, logger, handlerSpan := common.SetupHandler("Rack", "GetAll", c, garh.tracerSpan)
	if handlerSpan != nil {
		defer handlerSpan.End()
	}

	// Is DB user missing?
	if dbUser == nil {
		logger.Error().Msg("invalid User object found in request context")
		return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve current user", nil)
	}

	// ensure our user is a provider or privileged tenant for the org
	infrastructureProvider, tenant, apiError := common.IsProviderOrTenant(ctx, logger, garh.dbSession, org, dbUser, false, true)
	if apiError != nil {
		return cerr.NewAPIErrorResponse(c, apiError.Code, apiError.Message, apiError.Data)
	}

	// Get site ID from query param (required)
	siteStrID := c.QueryParam("siteId")
	if siteStrID == "" {
		if tenant != nil {
			// Tenants must specify a Site ID
			return cerr.NewAPIErrorResponse(c, http.StatusBadRequest, "Site ID must be specified in query when retrieving Racks as a Tenant", nil)
		}
		return cerr.NewAPIErrorResponse(c, http.StatusBadRequest, "siteId query parameter is required", nil)
	}

	// Validate the site
	site, err := common.GetSiteFromIDString(ctx, nil, siteStrID, garh.dbSession)
	if err != nil {
		if errors.Is(err, cdb.ErrDoesNotExist) {
			return cerr.NewAPIErrorResponse(c, http.StatusBadRequest, "Site specified in request does not exist", nil)
		}
		logger.Error().Err(err).Msg("error retrieving Site from DB")
		return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve Site specified in request due to DB error", nil)
	}

	// Validate Provider/Tenant Site access
	hasAccess, apiError := common.ValidateProviderOrTenantSiteAccess(ctx, logger, garh.dbSession, site, infrastructureProvider, tenant)
	if apiError != nil {
		return cerr.NewAPIErrorResponse(c, apiError.Code, apiError.Message, apiError.Data)
	}

	if !hasAccess {
		return cerr.NewAPIErrorResponse(c, http.StatusForbidden, "Current org is not associated with the Site specified in query", nil)
	}

	// Validate pagination request
	pageRequest := pagination.PageRequest{}
	err = c.Bind(&pageRequest)
	if err != nil {
		logger.Warn().Err(err).Msg("error binding pagination request data into API model")
		return cerr.NewAPIErrorResponse(c, http.StatusBadRequest, "Failed to parse request pagination data", nil)
	}

	// Validate pagination attributes
	err = pageRequest.Validate(model.RackOrderByFields)
	if err != nil {
		logger.Warn().Err(err).Msg("error validating pagination request data")
		return cerr.NewAPIErrorResponse(c, http.StatusBadRequest, "Failed to validate pagination request data", err)
	}

	// Check includeComponents query param (API uses includeComponents, RLA uses WithComponents)
	includeComponents := false
	if ic := c.QueryParam("includeComponents"); ic != "" {
		includeComponents, _ = strconv.ParseBool(ic)
	}

	// Build filters from query params
	var filters []*rlav1.Filter
	qParams := c.QueryParams()
	for _, field := range model.RackOrderByFields {
		if value := qParams.Get(field); value != "" {
			if f := model.BuildRackStringFilter(field, value); f != nil {
				filters = append(filters, f)
			}
		}
	}

	// Build OrderBy from pagination
	var orderBy *rlav1.OrderBy
	if pageRequest.OrderBy != nil {
		orderBy = model.BuildRackOrderBy(pageRequest.OrderBy.Field, strings.ToUpper(pageRequest.OrderBy.Order))
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

	// Build RLA request
	rlaRequest := &rlav1.GetListOfRacksRequest{
		Filters:        filters,
		WithComponents: includeComponents,
		Pagination:     paginationProto,
		OrderBy:        orderBy,
	}

	// Execute workflow
	workflowOptions := tClient.StartWorkflowOptions{
		ID:                       "GetRacks",
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
		logger.Error().Err(err).Msg("failed to get result from GetRacks workflow")
		return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, "Failed to get Racks", nil)
	}

	// Convert to API model
	apiRacks := make([]*model.APIRack, 0, len(rlaResponse.GetRacks()))
	for _, rack := range rlaResponse.GetRacks() {
		apiRacks = append(apiRacks, model.NewAPIRack(rack, includeComponents))
	}

	// Create pagination response header
	total := int(rlaResponse.GetTotal())
	pageNumber := 1
	pageSize := pagination.MaxPageSize
	if pageRequest.PageNumber != nil {
		pageNumber = *pageRequest.PageNumber
	}
	if pageRequest.PageSize != nil {
		pageSize = *pageRequest.PageSize
	}
	pageResponse := pagination.NewPageResponse(pageNumber, pageSize, total, pageRequest.OrderByStr)
	pageHeader, err := json.Marshal(pageResponse)
	if err != nil {
		logger.Error().Err(err).Msg("error marshaling pagination response")
		return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, "Failed to create pagination response", nil)
	}
	c.Response().Header().Set(pagination.ResponseHeaderName, string(pageHeader))

	logger.Info().Int("count", len(apiRacks)).Int("total", total).Msg("finishing API handler")

	return c.JSON(http.StatusOK, apiRacks)
}
