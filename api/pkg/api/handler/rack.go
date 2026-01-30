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

	"github.com/labstack/echo/v4"
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
// @Param withComponents query boolean false "Include rack components in response"
// @Success 200 {object} model.APIRack
// @Router /v2/org/{org}/forge/rack/{id} [get]
func (grh GetRackHandler) Handle(c echo.Context) error {
	org, dbUser, ctx, logger, handlerSpan := common.SetupHandler("Rack", "Get", c, grh.tracerSpan)
	if handlerSpan != nil {
		defer handlerSpan.End()
	}

	// Validate org membership
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

	// Get site ID from query param (required)
	siteStrID := c.QueryParam("siteId")
	if siteStrID == "" {
		return cerr.NewAPIErrorResponse(c, http.StatusBadRequest, "siteId query parameter is required", nil)
	}

	// Check that infrastructureProvider exists in org
	ip, err := common.GetInfrastructureProviderForOrg(ctx, nil, grh.dbSession, org)
	if err != nil {
		logger.Warn().Err(err).Msg("error getting infrastructure provider for org")
		return cerr.NewAPIErrorResponse(c, http.StatusBadRequest, "Failed to retrieve Infrastructure Provider for org", nil)
	}

	// Validate the site
	site, err := common.GetSiteFromIDString(ctx, nil, siteStrID, grh.dbSession)
	if err != nil {
		logger.Warn().Err(err).Str("Site ID", siteStrID).Msg("error getting site from request")
		return cerr.NewAPIErrorResponse(c, http.StatusBadRequest, "Error retrieving Site in request", nil)
	}

	// Verify site's infrastructure provider matches org's infrastructure provider
	if site.InfrastructureProviderID != ip.ID {
		return cerr.NewAPIErrorResponse(c, http.StatusForbidden, "Site specified in request doesn't belong to current org's Provider", nil)
	}

	// Check withComponents query param
	withComponents := false
	if wc := c.QueryParam("withComponents"); wc != "" {
		withComponents, _ = strconv.ParseBool(wc)
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
		WithComponents: withComponents,
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
		return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, "Failed to get rack from RLA", nil)
	}

	// Get workflow result
	var rlaResponse rlav1.GetRackInfoResponse
	err = we.Get(ctx, &rlaResponse)
	if err != nil {
		logger.Error().Err(err).Msg("failed to get result from GetRack workflow")
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
// @Success 200 {array} model.APIRack
// @Router /v2/org/{org}/forge/rack [get]
func (garh GetAllRackHandler) Handle(c echo.Context) error {
	org, dbUser, ctx, logger, handlerSpan := common.SetupHandler("Rack", "GetAll", c, garh.tracerSpan)
	if handlerSpan != nil {
		defer handlerSpan.End()
	}

	// Validate org membership
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

	// Check that infrastructureProvider exists in org
	ip, err := common.GetInfrastructureProviderForOrg(ctx, nil, garh.dbSession, org)
	if err != nil {
		logger.Warn().Err(err).Msg("error getting infrastructure provider for org")
		return cerr.NewAPIErrorResponse(c, http.StatusBadRequest, "Failed to retrieve Infrastructure Provider for org", nil)
	}

	// Validate the site
	site, err := common.GetSiteFromIDString(ctx, nil, siteStrID, garh.dbSession)
	if err != nil {
		logger.Warn().Err(err).Str("Site ID", siteStrID).Msg("error getting site from request")
		return cerr.NewAPIErrorResponse(c, http.StatusBadRequest, "Error retrieving Site in request", nil)
	}

	// Verify site's infrastructure provider matches org's infrastructure provider
	if site.InfrastructureProviderID != ip.ID {
		return cerr.NewAPIErrorResponse(c, http.StatusForbidden, "Site specified in request doesn't belong to current org's Provider", nil)
	}

	// Check withComponents query param
	withComponents := false
	if wc := c.QueryParam("withComponents"); wc != "" {
		withComponents, _ = strconv.ParseBool(wc)
	}

	// Get the temporal client for the site
	stc, err := garh.scp.GetClientByID(site.ID)
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
		ID:                       "GetRacks",
		WorkflowExecutionTimeout: common.WorkflowExecutionTimeout,
		TaskQueue:                queue.SiteTaskQueue,
	}

	ctx, cancel := context.WithTimeout(ctx, common.WorkflowContextTimeout)
	defer cancel()

	we, err := stc.ExecuteWorkflow(ctx, workflowOptions, "GetRacks", rlaRequest)
	if err != nil {
		logger.Error().Err(err).Msg("failed to execute GetRacks workflow")
		return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, "Failed to get racks from RLA", nil)
	}

	// Get workflow result
	var rlaResponse rlav1.GetListOfRacksResponse
	err = we.Get(ctx, &rlaResponse)
	if err != nil {
		logger.Error().Err(err).Msg("failed to get result from GetRacks workflow")
		return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, "Failed to get racks from RLA", nil)
	}

	// Convert to API model
	apiRacks := model.NewAPIRacks(&rlaResponse, withComponents)

	logger.Info().Int("count", len(apiRacks)).Msg("finishing API handler")

	return c.JSON(http.StatusOK, apiRacks)
}
