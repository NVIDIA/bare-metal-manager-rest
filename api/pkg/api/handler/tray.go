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
	cerr "github.com/nvidia/bare-metal-manager-rest/common/pkg/util"
	sutil "github.com/nvidia/bare-metal-manager-rest/common/pkg/util"
	cdb "github.com/nvidia/bare-metal-manager-rest/db/pkg/db"
	rlav1 "github.com/nvidia/bare-metal-manager-rest/workflow-schema/rla/protobuf/v1"
	"github.com/nvidia/bare-metal-manager-rest/workflow/pkg/queue"
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
// @Description Get a Tray by ID from RLA
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

	// Is DB user missing?
	if dbUser == nil {
		logger.Error().Msg("invalid User object found in request context")
		return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve current user", nil)
	}

	// Ensure user is a provider or tenant for the org
	infrastructureProvider, tenant, apiErr := common.IsProviderOrTenant(ctx, logger, gth.dbSession, org, dbUser, true, false)
	if apiErr != nil {
		return c.JSON(apiErr.Code, apiErr)
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

	// Validate provider/tenant site access
	hasAccess, apiErr := ValidateProviderOrTenantSiteAccess(ctx, logger, gth.dbSession, site, infrastructureProvider, tenant)
	if apiErr != nil {
		return c.JSON(apiErr.Code, apiErr)
	}
	if !hasAccess {
		return cerr.NewAPIErrorResponse(c, http.StatusForbidden, "User does not have access to Site", nil)
	}

	// Get tray ID from URL param
	trayStrID := c.Param("id")
	if trayStrID == "" {
		return cerr.NewAPIErrorResponse(c, http.StatusBadRequest, "tray id is required", nil)
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
	}

	ctx, cancel := context.WithTimeout(ctx, common.WorkflowContextTimeout)
	defer cancel()

	we, err := stc.ExecuteWorkflow(ctx, workflowOptions, "GetTray", rlaRequest)
	if err != nil {
		logger.Error().Err(err).Msg("failed to execute GetTray workflow")
		return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, "Failed to get tray from RLA", nil)
	}

	// Get workflow result
	var rlaResponse rlav1.GetComponentInfoResponse
	err = we.Get(ctx, &rlaResponse)
	if err != nil {
		logger.Error().Err(err).Msg("failed to get result from GetTray workflow")
		return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, "Failed to get tray from RLA", nil)
	}

	// Convert to API model
	apiTray := model.NewAPITray(rlaResponse.GetComponent())

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
// @Description Get all Trays from RLA with optional filters
// @Tags tray
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param org path string true "Name of NGC organization"
// @Param siteId query string true "ID of the Site"
// @Param rackId query string false "Filter by Rack ID"
// @Param rackName query string false "Filter by Rack name"
// @Param type query string false "Filter by tray type (compute, switch, powershelf, torswitch, ums, cdu)"
// @Param componentId query string false "Filter by component IDs (comma-separated)"
// @Param id query string false "Filter by tray UUIDs (comma-separated)"
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

	// Is DB user missing?
	if dbUser == nil {
		logger.Error().Msg("invalid User object found in request context")
		return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve current user", nil)
	}

	// Ensure user is a provider or tenant for the org
	infrastructureProvider, tenant, apiErr := common.IsProviderOrTenant(ctx, logger, gath.dbSession, org, dbUser, true, false)
	if apiErr != nil {
		return c.JSON(apiErr.Code, apiErr)
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

	// Validate provider/tenant site access
	hasAccess, apiErr := ValidateProviderOrTenantSiteAccess(ctx, logger, gath.dbSession, site, infrastructureProvider, tenant)
	if apiErr != nil {
		return c.JSON(apiErr.Code, apiErr)
	}
	if !hasAccess {
		return cerr.NewAPIErrorResponse(c, http.StatusForbidden, "User does not have access to Site", nil)
	}

	// Build and validate filter input from query params
	filter, filterErr := buildTrayFilterInput(c)
	if filterErr != nil {
		logger.Warn().Err(filterErr).Msg("invalid tray filter query parameters")
		return cerr.NewAPIErrorResponse(c, http.StatusBadRequest, "Invalid filter parameters", filterErr)
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

	// Build RLA request with filters
	rlaRequest := buildRLARequestFromFilter(filter)

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
		ID:                       fmt.Sprintf("tray-get-all-%s", filter.Hash()),
		WorkflowExecutionTimeout: common.WorkflowExecutionTimeout,
		TaskQueue:                queue.SiteTaskQueue,
	}

	ctx, cancel := context.WithTimeout(ctx, common.WorkflowContextTimeout)
	defer cancel()

	we, err := stc.ExecuteWorkflow(ctx, workflowOptions, "GetTrays", rlaRequest)
	if err != nil {
		logger.Error().Err(err).Msg("failed to execute GetTrays workflow")
		return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, "Failed to get trays from RLA", nil)
	}

	// Get workflow result
	var rlaResponse rlav1.GetComponentsResponse
	err = we.Get(ctx, &rlaResponse)
	if err != nil {
		logger.Error().Err(err).Msg("failed to get result from GetTrays workflow")
		return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, "Failed to get trays from RLA", nil)
	}

	apiTrays := model.NewAPITrays(&rlaResponse)

	// Set pagination response header
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

	logger.Info().Int("count", len(apiTrays)).Int("Total", total).Msg("finishing API handler")

	return c.JSON(http.StatusOK, apiTrays)
}

// buildTrayFilterInput builds and validates a TrayFilterInput from query parameters.
// Returns an error if any parameter has an invalid format (e.g., non-integer slot, invalid UUID).
func buildTrayFilterInput(c echo.Context) (*model.TrayFilterInput, error) {
	filter := &model.TrayFilterInput{}
	qParams := c.QueryParams()

	if rackID := c.QueryParam("rackId"); rackID != "" {
		if _, err := uuid.Parse(rackID); err != nil {
			return nil, fmt.Errorf("rackId must be a valid UUID: %s", rackID)
		}
		filter.RackID = &rackID
	}
	if rackName := c.QueryParam("rackName"); rackName != "" {
		filter.RackName = &rackName
	}
	// Support rackname (lowercase) as alias for rackName per API spec
	if filter.RackName == nil {
		if rackName := c.QueryParam("rackname"); rackName != "" {
			filter.RackName = &rackName
		}
	}
	if trayType := c.QueryParam("type"); trayType != "" {
		if !slices.Contains(model.ValidTrayTypes, trayType) {
			return nil, fmt.Errorf("type must be one of %v: %s", model.ValidTrayTypes, trayType)
		}
		filter.Type = &trayType
	}
	// componentId: support comma-separated list (e.g. componentId=id1,id2,id3) and repeated params
	if componentIDs := qParams["componentId"]; len(componentIDs) > 0 {
		filter.ComponentIDs = common.SplitCommaSeparated(componentIDs)
	}
	// id: support comma-separated list (e.g. id=uuid1,uuid2,uuid3) and repeated params
	if ids := qParams["id"]; len(ids) > 0 {
		parsed := common.SplitCommaSeparated(ids)
		for _, id := range parsed {
			if _, err := uuid.Parse(id); err != nil {
				return nil, fmt.Errorf("id must be a valid UUID: %s", id)
			}
		}
		filter.IDs = parsed
	}
	// componentId requires type to be specified so the correct source system can be queried
	if len(filter.ComponentIDs) > 0 && filter.Type == nil {
		return nil, fmt.Errorf("type is required when componentId is provided")
	}

	return filter, nil
}

// buildRLARequestFromFilter builds an RLA GetComponentsRequest from TrayFilterInput.
// Priority: component-level targeting (IDs, componentIDs+type) > rack-level targeting.
func buildRLARequestFromFilter(filter *model.TrayFilterInput) *rlav1.GetComponentsRequest {
	rlaRequest := &rlav1.GetComponentsRequest{}

	// Check for component-level targeting: UUID-based IDs or componentIDs with type (ExternalRef)
	hasIDs := len(filter.IDs) > 0
	hasComponentIDsWithType := len(filter.ComponentIDs) > 0 && filter.Type != nil

	if hasIDs || hasComponentIDsWithType {
		componentTargets := make([]*rlav1.ComponentTarget, 0, len(filter.IDs)+len(filter.ComponentIDs))

		// Add UUID-based targets
		for _, id := range filter.IDs {
			componentTargets = append(componentTargets, &rlav1.ComponentTarget{
				Identifier: &rlav1.ComponentTarget_Id{
					Id: &rlav1.UUID{Id: id},
				},
			})
		}

		// Add ExternalRef-based targets (componentIDs + type)
		if hasComponentIDsWithType {
			if protoName, ok := model.APIToProtoComponentTypeName[*filter.Type]; ok {
				protoType := rlav1.ComponentType(rlav1.ComponentType_value[protoName])
				for _, cid := range filter.ComponentIDs {
					componentTargets = append(componentTargets, &rlav1.ComponentTarget{
						Identifier: &rlav1.ComponentTarget_External{
							External: &rlav1.ExternalRef{
								Type: protoType,
								Id:   cid,
							},
						},
					})
				}
			}
		}

		rlaRequest.TargetSpec = &rlav1.OperationTargetSpec{
			Targets: &rlav1.OperationTargetSpec_Components{
				Components: &rlav1.ComponentTargets{
					Targets: componentTargets,
				},
			},
		}
		return rlaRequest
	}

	// Rack-level targeting (only when no component-level targeting applies)
	if filter.RackID != nil || filter.RackName != nil || filter.Type != nil {
		rackTarget := &rlav1.RackTarget{}

		if filter.RackID != nil {
			rackTarget.Identifier = &rlav1.RackTarget_Id{
				Id: &rlav1.UUID{Id: *filter.RackID},
			}
		} else if filter.RackName != nil {
			rackTarget.Identifier = &rlav1.RackTarget_Name{
				Name: *filter.RackName,
			}
		}

		// Parse component type filter
		if filter.Type != nil {
			if protoName, ok := model.APIToProtoComponentTypeName[*filter.Type]; ok {
				rackTarget.ComponentTypes = []rlav1.ComponentType{
					rlav1.ComponentType(rlav1.ComponentType_value[protoName]),
				}
			}
		}

		rlaRequest.TargetSpec = &rlav1.OperationTargetSpec{
			Targets: &rlav1.OperationTargetSpec_Racks{
				Racks: &rlav1.RackTargets{
					Targets: []*rlav1.RackTarget{rackTarget},
				},
			},
		}
	}

	return rlaRequest
}
