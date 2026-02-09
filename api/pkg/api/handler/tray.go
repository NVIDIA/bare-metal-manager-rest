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
	"strings"

	"github.com/labstack/echo/v4"
	"go.opentelemetry.io/otel/attribute"
	tClient "go.temporal.io/sdk/client"

	"github.com/nvidia/carbide-rest/api/internal/config"
	"github.com/nvidia/carbide-rest/api/pkg/api/handler/util/common"
	"github.com/nvidia/carbide-rest/api/pkg/api/model"
	sc "github.com/nvidia/carbide-rest/api/pkg/client/site"
	cerr "github.com/nvidia/carbide-rest/common/pkg/util"
	sutil "github.com/nvidia/carbide-rest/common/pkg/util"
	cdb "github.com/nvidia/carbide-rest/db/pkg/db"
	rlav1 "github.com/nvidia/carbide-rest/workflow-schema/rla/protobuf/v1"
	"github.com/nvidia/carbide-rest/workflow/pkg/queue"
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

	// Validate Provider Admin role, org membership, and site access
	site, apiErr := common.ValidateProviderSiteAccess(ctx, logger, gth.dbSession, org, dbUser, c.QueryParam("siteId"))
	if apiErr != nil {
		return c.JSON(apiErr.Code, apiErr)
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
		ID:                       fmt.Sprintf("GetTray-%s", trayStrID),
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
// @Param slot query int false "Filter by tray slot number"
// @Param type query string false "Filter by tray type (compute, switch, powershelf, torswitch, ums, cdu)"
// @Param index query int false "Filter by index of tray in its type"
// @Param componentId query string false "Filter by component IDs (comma-separated)"
// @Param id query string false "Filter by tray UUIDs (comma-separated)"
// @Param taskId query string false "Filter by task UUID"
// @Success 200 {array} model.APITray
// @Router /v2/org/{org}/carbide/tray [get]
func (gath GetAllTrayHandler) Handle(c echo.Context) error {
	org, dbUser, ctx, logger, handlerSpan := common.SetupHandler("Tray", "GetAll", c, gath.tracerSpan)
	if handlerSpan != nil {
		defer handlerSpan.End()
	}

	// Validate Provider Admin role, org membership, and site access
	site, apiErr := common.ValidateProviderSiteAccess(ctx, logger, gath.dbSession, org, dbUser, c.QueryParam("siteId"))
	if apiErr != nil {
		return c.JSON(apiErr.Code, apiErr)
	}

	// Build filter input from query params
	filter := buildTrayFilterInput(c)

	// Get the temporal client for the site
	stc, err := gath.scp.GetClientByID(site.ID)
	if err != nil {
		logger.Error().Err(err).Msg("failed to retrieve Temporal client for Site")
		return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve client for Site", nil)
	}

	// Build RLA request with filters
	rlaRequest := buildRLARequestFromFilter(filter)

	// Extract taskID to pass into the workflow (resolved inside GetTrays)
	taskID := ""
	if filter.TaskID != nil {
		taskID = *filter.TaskID
	}

	// Execute workflow — task ID resolution is handled inside the GetTrays workflow
	workflowOptions := tClient.StartWorkflowOptions{
		ID:                       "GetTrays",
		WorkflowExecutionTimeout: common.WorkflowExecutionTimeout,
		TaskQueue:                queue.SiteTaskQueue,
	}

	ctx, cancel := context.WithTimeout(ctx, common.WorkflowContextTimeout)
	defer cancel()

	we, err := stc.ExecuteWorkflow(ctx, workflowOptions, "GetTrays", rlaRequest, taskID)
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

	// Convert to API model and apply client-side filters if needed
	apiTrays := model.NewAPITraysWithFilter(&rlaResponse, filter)

	logger.Info().Int("count", len(apiTrays)).Msg("finishing API handler")

	return c.JSON(http.StatusOK, apiTrays)
}

// buildTrayFilterInput builds a TrayFilterInput from query parameters
func buildTrayFilterInput(c echo.Context) *model.TrayFilterInput {
	filter := &model.TrayFilterInput{}
	qParams := c.QueryParams()

	if rackID := c.QueryParam("rackId"); rackID != "" {
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
	if slotStr := c.QueryParam("slot"); slotStr != "" {
		if slot, err := strconv.ParseInt(slotStr, 10, 32); err == nil {
			slotVal := int32(slot)
			filter.Slot = &slotVal
		}
	}
	if trayType := c.QueryParam("type"); trayType != "" {
		filter.Type = &trayType
	}
	if indexStr := c.QueryParam("index"); indexStr != "" {
		if index, err := strconv.ParseInt(indexStr, 10, 32); err == nil {
			indexVal := int32(index)
			filter.Index = &indexVal
		}
	}
	// componentId: support comma-separated list (e.g. componentId=id1,id2,id3) and repeated params
	if componentIDs := qParams["componentId"]; len(componentIDs) > 0 {
		filter.ComponentIDs = common.SplitCommaSeparated(componentIDs)
	}
	// id: support comma-separated list (e.g. id=uuid1,uuid2,uuid3) and repeated params
	if ids := qParams["id"]; len(ids) > 0 {
		filter.IDs = common.SplitCommaSeparated(ids)
	}
	if taskID := c.QueryParam("taskId"); taskID != "" {
		filter.TaskID = &taskID
	}

	return filter
}

// buildRLARequestFromFilter builds an RLA GetComponentsRequest from TrayFilterInput
func buildRLARequestFromFilter(filter *model.TrayFilterInput) *rlav1.GetComponentsRequest {
	rlaRequest := &rlav1.GetComponentsRequest{}

	// Build OperationTargetSpec based on filters
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
			componentType := stringToComponentType(*filter.Type)
			if componentType != rlav1.ComponentType_COMPONENT_TYPE_UNKNOWN {
				rackTarget.ComponentTypes = []rlav1.ComponentType{componentType}
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

	// Handle ID-based filtering via ComponentTargets
	if len(filter.IDs) > 0 {
		componentTargets := make([]*rlav1.ComponentTarget, 0, len(filter.IDs))
		for _, id := range filter.IDs {
			componentTargets = append(componentTargets, &rlav1.ComponentTarget{
				Identifier: &rlav1.ComponentTarget_Id{
					Id: &rlav1.UUID{Id: id},
				},
			})
		}
		rlaRequest.TargetSpec = &rlav1.OperationTargetSpec{
			Targets: &rlav1.OperationTargetSpec_Components{
				Components: &rlav1.ComponentTargets{
					Targets: componentTargets,
				},
			},
		}
	}

	return rlaRequest
}

// stringToComponentType converts a string to a ComponentType enum
func stringToComponentType(s string) rlav1.ComponentType {
	switch strings.ToLower(s) {
	case "compute":
		return rlav1.ComponentType_COMPONENT_TYPE_COMPUTE
	case "switch":
		return rlav1.ComponentType_COMPONENT_TYPE_NVLSWITCH
	case "powershelf":
		return rlav1.ComponentType_COMPONENT_TYPE_POWERSHELF
	case "torswitch":
		return rlav1.ComponentType_COMPONENT_TYPE_TORSWITCH
	case "ums":
		return rlav1.ComponentType_COMPONENT_TYPE_UMS
	case "cdu":
		return rlav1.ComponentType_COMPONENT_TYPE_CDU
	default:
		return rlav1.ComponentType_COMPONENT_TYPE_UNKNOWN
	}
}
