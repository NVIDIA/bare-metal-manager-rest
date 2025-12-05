// SPDX-FileCopyrightText: Copyright (c) 2021-2023 NVIDIA CORPORATION & AFFILIATES. All rights reserved.
// SPDX-License-Identifier: LicenseRef-NvidiaProprietary
//
// NVIDIA CORPORATION, its affiliates and licensors retain all intellectual
// property and proprietary rights in and to this material, related
// documentation and any modifications thereto. Any use, reproduction,
// disclosure or distribution of this material and related documentation
// without an express license agreement from NVIDIA CORPORATION or
// its affiliates is strictly prohibited.

package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"go.opentelemetry.io/otel/attribute"
	temporalClient "go.temporal.io/sdk/client"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"

	"github.com/labstack/echo/v4"

	cdb "github.com/nvidia/carbide-rest/db/pkg/db"
	cdbm "github.com/nvidia/carbide-rest/db/pkg/db/model"
	"github.com/nvidia/carbide-rest/db/pkg/db/paginator"

	"github.com/nvidia/carbide-rest/api/internal/config"
	common "github.com/nvidia/carbide-rest/api/pkg/api/handler/util/common"
	"github.com/nvidia/carbide-rest/api/pkg/api/model"
	"github.com/nvidia/carbide-rest/api/pkg/api/pagination"
	auth "github.com/nvidia/carbide-rest/auth/pkg/authorization"
	cerr "github.com/nvidia/carbide-rest/common/pkg/util"
	sutil "github.com/nvidia/carbide-rest/common/pkg/util"
)

// ~~~~~ GetAll NVLinkInterface Handler ~~~~~ //

// GetAllNVLinkInterfaceHandler is the API Handler for retrieving all NVLinkInterfaces for an Ins
// tance
type GetAllNVLinkInterfaceHandler struct {
	dbSession  *cdb.Session
	tc         temporalClient.Client
	cfg        *config.Config
	tracerSpan *sutil.TracerSpan
}

// NewGetAllNVLinkInterfaceHandler initializes and returns a new handler for retrieving all subnets for an Instance
func NewGetAllNVLinkInterfaceHandler(dbSession *cdb.Session, tc temporalClient.Client, cfg *config.Config) GetAllNVLinkInterfaceHandler {
	return GetAllNVLinkInterfaceHandler{
		dbSession:  dbSession,
		tc:         tc,
		cfg:        cfg,
		tracerSpan: sutil.NewTracerSpan(),
	}
}

// Handle godoc
// @Summary Retrieve all NVLinkInterfaces for an Instance
// @Description Retrieve all NVLinkInterfaces for an Instance
// @Tags NVLinkInterface
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param org path string true "Name of NGC organization"
// @Param instanceId path string true "ID of Instance"
// @Param status query string false "Filter by status" e.g. 'Pending', 'Error'"
// @Param includeRelation query string false "Related entities to include in response e.g. 'Instance', 'Subnet'"
// @Param pageNumber query integer false "Page number of results returned"
// @Param pageSize query integer false "Number of results per page"
// @Param orderBy query string false "Order by field"
// @Success 200 {object} model.APINVLinkInterface
// @Router /v2/org/{org}/carbide/instance/{instance_id}/nvlink-interface [get]
func (gaish GetAllNVLinkInterfaceHandler) Handle(c echo.Context) error {
	// Get context
	ctx := c.Request().Context()

	// Get org
	org := c.Param("orgName")

	// Initialize logger
	logger := log.With().Str("Model", "NVLinkInterface").Str("Handler", "GetAll").Str("Org", org).Logger()

	logger.Info().Msg("started API handler")

	// Create a child span and set the attributes for current request
	newctx, handlerSpan := gaish.tracerSpan.CreateChildInContext(ctx, "GetAllNVLinkInterfaceHandler", logger)
	if handlerSpan != nil {
		// Set newly created span context as a current context
		ctx = newctx

		defer handlerSpan.End()

		gaish.tracerSpan.SetAttribute(handlerSpan, attribute.String("org", org), logger)
	}

	dbUser, logger, err := common.GetUserAndEnrichLogger(c, logger, gaish.tracerSpan, handlerSpan)
	if err != nil {
		return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve current user", nil)
	}

	// Validate org
	ok, err := auth.ValidateOrgMembership(dbUser, org)
	if !ok {
		if err != nil {
			logger.Error().Err(err).Msg("error validating org membership for User in request")
		} else {
			logger.Warn().Msg("could not validate org membership for user, access denied")
		}
		return cerr.NewAPIErrorResponse(c, http.StatusForbidden, fmt.Sprintf("Failed to validate membership for org: %s", org), nil)
	}

	// Validate role, only Tenant Admins are allowed to retrieve Instances
	ok = auth.ValidateUserRoles(dbUser, org, nil, auth.TenantAdminRole)
	if !ok {
		logger.Warn().Msg("user does not have Tenant Admin role, access denied")
		return cerr.NewAPIErrorResponse(c, http.StatusForbidden, "User does not have Tenant Admin role with org", nil)
	}

	// Validate pagination request
	pageRequest := pagination.PageRequest{}
	err = c.Bind(&pageRequest)
	if err != nil {
		logger.Warn().Err(err).Msg("error binding pagination request data into API model")
		return cerr.NewAPIErrorResponse(c, http.StatusBadRequest, "Failed to parse request pagination data", nil)
	}

	// Validate pagination request attributes
	err = pageRequest.Validate(cdbm.NVLinkInterfaceOrderByFields)
	if err != nil {
		logger.Warn().Err(err).Msg("error validating pagination request data")
		return cerr.NewAPIErrorResponse(c, http.StatusBadRequest,
			"Failed to validate pagination request data", err)
	}

	// Get and validate includeRelation params
	qParams := c.QueryParams()
	qIncludeRelations, errMsg := common.GetAndValidateQueryRelations(qParams, cdbm.NVLinkInterfaceRelatedEntities)
	if errMsg != "" {
		logger.Warn().Msg(errMsg)
		return cerr.NewAPIErrorResponse(c, http.StatusBadRequest, errMsg, nil)
	}

	// Get status from query param
	var status *string

	statusQuery := c.QueryParam("status")
	if statusQuery != "" {
		gaish.tracerSpan.SetAttribute(handlerSpan, attribute.String("status", statusQuery), logger)
		_, ok := cdbm.NVLinkInterfaceStatusMap[statusQuery]
		if !ok {
			logger.Warn().Msg(fmt.Sprintf("invalid value in status query: %v", statusQuery))
			return cerr.NewAPIErrorResponse(c, http.StatusBadRequest, "Invalid Status value in query", nil)
		}
		status = &statusQuery
	}

	// Get Tenant for this org
	tnDAO := cdbm.NewTenantDAO(gaish.dbSession)

	tenants, err := tnDAO.GetAllByOrg(ctx, nil, org, nil)
	if err != nil {
		logger.Error().Err(err).Msg("error retrieving Tenant for this org")
		return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve Tenant", nil)
	}

	if len(tenants) == 0 {
		return cerr.NewAPIErrorResponse(c, http.StatusForbidden, "Org does not have a Tenant associated", nil)
	}
	tenant := tenants[0]

	// Get Instance ID from URL param
	instanceStrID := c.Param("instanceId")
	instanceID, err := uuid.Parse(instanceStrID)
	if err != nil {
		return cerr.NewAPIErrorResponse(c, http.StatusBadRequest, "Invalid Instance ID in URL", nil)
	}

	// Get Instance
	instanceDAO := cdbm.NewInstanceDAO(gaish.dbSession)

	instance, err := instanceDAO.GetByID(ctx, nil, instanceID, nil)
	if err != nil {
		if err == cdb.ErrDoesNotExist {
			return cerr.NewAPIErrorResponse(c, http.StatusNotFound, "Could not find Instance with specified ID", nil)
		}
		logger.Error().Err(err).Msg("error retrieving Instance from DB")
		return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve Instance", nil)
	}

	// Check if Instance belongs to Tenant
	if instance.TenantID != tenant.ID {
		return cerr.NewAPIErrorResponse(c, http.StatusForbidden, "Instance does not belong to current Tenant", nil)
	}

	// Get the instance subnets record from the db
	nvlIfcDAO := cdbm.NewNVLinkInterfaceDAO(gaish.dbSession)

	filterInput := cdbm.NVLinkInterfaceFilterInput{
		InstanceIDs: []uuid.UUID{instanceID},
	}

	if status != nil {
		filterInput.Statuses = append(filterInput.Statuses, *status)
	}

	pageInput := paginator.PageInput{
		Limit:   pageRequest.Limit,
		Offset:  pageRequest.Offset,
		OrderBy: pageRequest.OrderBy,
	}

	dbNVLinkInterfaces, total, err := nvlIfcDAO.GetAll(ctx, nil, filterInput, pageInput, qIncludeRelations)
	if err != nil {
		logger.Error().Err(err).Msg("error retrieving instance NVLink Interface Details from DB")
		return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve Instance NVLink Interface Details for Instance", nil)
	}

	// Create response
	apiNVLinkInterfaces := []model.APINVLinkInterface{}
	for _, dbnvlifc := range dbNVLinkInterfaces {
		curnvlifc := dbnvlifc
		apiNVLinkInterfaces = append(apiNVLinkInterfaces, *model.NewAPINVLinkInterface(&curnvlifc))
	}

	// Create pagination response header
	pageReponse := pagination.NewPageResponse(*pageRequest.PageNumber, *pageRequest.PageSize, total, pageRequest.OrderByStr)
	pageHeader, err := json.Marshal(pageReponse)
	if err != nil {
		logger.Error().Err(err).Msg("error marshaling pagination response")
		return cerr.NewAPIErrorResponse(c, http.StatusInternalServerError, "Failed to generate pagination response header", nil)
	}
	c.Response().Header().Set(pagination.ResponseHeaderName, string(pageHeader))

	logger.Info().Msg("finishing API handler")

	return c.JSON(http.StatusOK, apiNVLinkInterfaces)
}
