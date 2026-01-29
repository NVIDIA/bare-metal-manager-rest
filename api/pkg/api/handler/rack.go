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
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"go.opentelemetry.io/otel/attribute"
	tClient "go.temporal.io/sdk/client"

	"github.com/nvidia/carbide-rest/api/internal/config"
	"github.com/nvidia/carbide-rest/api/pkg/api/handler/util/common"
	"github.com/nvidia/carbide-rest/api/pkg/api/model"
	sc "github.com/nvidia/carbide-rest/api/pkg/client/site"
	auth "github.com/nvidia/carbide-rest/auth/pkg/authorization"
	cerr "github.com/nvidia/carbide-rest/common/pkg/util"
	sutil "github.com/nvidia/carbide-rest/common/pkg/util"
	cdb "github.com/nvidia/carbide-rest/db/pkg/db"
	rlav1 "github.com/nvidia/carbide-rest/workflow-schema/rla/protobuf/v1"
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
// @Param withComponents query boolean false "Include rack components in response"
// @Success 200 {object} model.APIRack
// @Router /v2/org/{org}/forge/rack/{id} [get]
func (grh GetRackHandler) Handle(c echo.Context) error {
	ctx := c.Request().Context()
	org := c.Param("orgName")

	logger := log.With().Str("Model", "Rack").Str("Handler", "Get").Str("Org", org).Logger()
	logger.Info().Msg("started API handler")

	// Create tracer span
	newctx, handlerSpan := grh.tracerSpan.CreateChildInContext(ctx, "GetRackHandler", logger)
	if handlerSpan != nil {
		ctx = newctx
		defer handlerSpan.End()
		grh.tracerSpan.SetAttribute(handlerSpan, attribute.String("org", org), logger)
	}

	// Get user and validate org membership
	dbUser, logger, err := common.GetUserAndEnrichLogger(c, logger, grh.tracerSpan, handlerSpan)
	if err != nil {
		return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve current user", nil)
	}

	ok, err := auth.ValidateOrgMembership(dbUser, org)
	if !ok {
		if err != nil {
			logger.Error().Err(err).Msg("error validating org membership for User in request")
		}
		return cerr.NewAPIErrorResponse(c, http.StatusForbidden, fmt.Sprintf("Failed to validate membership for org: %s", org), nil)
	}

	// Get rack ID from URL param
	rackStrID := c.Param("id")
	grh.tracerSpan.SetAttribute(handlerSpan, attribute.String("rack_id", rackStrID), logger)

	// Validate rack ID is a valid UUID
	_, err = uuid.Parse(rackStrID)
	if err != nil {
		return cerr.NewAPIErrorResponse(c, http.StatusBadRequest, "Invalid Rack ID in URL", nil)
	}

	// Get site ID from query param (required)
	siteStrID := c.QueryParam("siteId")
	if siteStrID == "" {
		return cerr.NewAPIErrorResponse(c, http.StatusBadRequest, "siteId query parameter is required", nil)
	}

	siteID, err := uuid.Parse(siteStrID)
	if err != nil {
		return cerr.NewAPIErrorResponse(c, http.StatusBadRequest, "Invalid siteId in query", nil)
	}

	// Check withComponents query param
	withComponents := false
	if wc := c.QueryParam("withComponents"); wc != "" {
		withComponents, _ = strconv.ParseBool(wc)
	}

	// Get the temporal client for the site
	stc, err := grh.scp.GetClientByID(siteID)
	if err != nil {
		logger.Error().Err(err).Msg("failed to retrieve Temporal client for Site")
		return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve client for Site", nil)
	}

	// Build RLA request
	rlaRequest := &rlav1.GetRackInfoByIDRequest{
		Id:             &rlav1.UUID{Id: rackStrID},
		WithComponents: withComponents,
	}

	// Execute workflow
	workflowOptions := tClient.StartWorkflowOptions{
		ID:        fmt.Sprintf("GetRackByID-%s-%d", rackStrID, time.Now().UnixNano()),
		TaskQueue: siteID.String(),
	}

	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	we, err := stc.ExecuteWorkflow(ctx, workflowOptions, "GetRackByID", rlaRequest)
	if err != nil {
		logger.Error().Err(err).Msg("failed to execute GetRackByID workflow")
		return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, "Failed to get rack from RLA", nil)
	}

	// Get workflow result
	var rlaResponse rlav1.GetRackInfoResponse
	err = we.Get(ctx, &rlaResponse)
	if err != nil {
		logger.Error().Err(err).Msg("failed to get result from GetRackByID workflow")
		return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, "Failed to get rack from RLA", nil)
	}

	// Convert to API model
	apiRack := model.NewAPIRack(rlaResponse.GetRack(), withComponents)

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
// @Param withComponents query boolean false "Include rack components in response"
// @Success 200 {object} model.APIRackListResponse
// @Router /v2/org/{org}/forge/rack [get]
func (garh GetAllRackHandler) Handle(c echo.Context) error {
	ctx := c.Request().Context()
	org := c.Param("orgName")

	logger := log.With().Str("Model", "Rack").Str("Handler", "GetAll").Str("Org", org).Logger()
	logger.Info().Msg("started API handler")

	// Create tracer span
	newctx, handlerSpan := garh.tracerSpan.CreateChildInContext(ctx, "GetAllRackHandler", logger)
	if handlerSpan != nil {
		ctx = newctx
		defer handlerSpan.End()
		garh.tracerSpan.SetAttribute(handlerSpan, attribute.String("org", org), logger)
	}

	// Get user and validate org membership
	dbUser, logger, err := common.GetUserAndEnrichLogger(c, logger, garh.tracerSpan, handlerSpan)
	if err != nil {
		return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve current user", nil)
	}

	ok, err := auth.ValidateOrgMembership(dbUser, org)
	if !ok {
		if err != nil {
			logger.Error().Err(err).Msg("error validating org membership for User in request")
		}
		return cerr.NewAPIErrorResponse(c, http.StatusForbidden, fmt.Sprintf("Failed to validate membership for org: %s", org), nil)
	}

	// Get site ID from query param (required)
	siteStrID := c.QueryParam("siteId")
	if siteStrID == "" {
		return cerr.NewAPIErrorResponse(c, http.StatusBadRequest, "siteId query parameter is required", nil)
	}

	siteID, err := uuid.Parse(siteStrID)
	if err != nil {
		return cerr.NewAPIErrorResponse(c, http.StatusBadRequest, "Invalid siteId in query", nil)
	}

	// Check withComponents query param
	withComponents := false
	if wc := c.QueryParam("withComponents"); wc != "" {
		withComponents, _ = strconv.ParseBool(wc)
	}

	// Get the temporal client for the site
	stc, err := garh.scp.GetClientByID(siteID)
	if err != nil {
		logger.Error().Err(err).Msg("failed to retrieve Temporal client for Site")
		return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve client for Site", nil)
	}

	// Build RLA request
	rlaRequest := &rlav1.GetListOfRacksRequest{
		WithComponents: withComponents,
	}

	// Execute workflow
	workflowOptions := tClient.StartWorkflowOptions{
		ID:        fmt.Sprintf("GetListOfRacks-%d", time.Now().UnixNano()),
		TaskQueue: siteID.String(),
	}

	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	we, err := stc.ExecuteWorkflow(ctx, workflowOptions, "GetListOfRacks", rlaRequest)
	if err != nil {
		logger.Error().Err(err).Msg("failed to execute GetListOfRacks workflow")
		return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, "Failed to get racks from RLA", nil)
	}

	// Get workflow result
	var rlaResponse rlav1.GetListOfRacksResponse
	err = we.Get(ctx, &rlaResponse)
	if err != nil {
		logger.Error().Err(err).Msg("failed to get result from GetListOfRacks workflow")
		return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, "Failed to get racks from RLA", nil)
	}

	// Convert to API model
	apiResponse := model.NewAPIRackListResponse(&rlaResponse, withComponents)

	logger.Info().Int32("total", apiResponse.Total).Msg("finishing API handler")

	return c.JSON(http.StatusOK, apiResponse)
}
