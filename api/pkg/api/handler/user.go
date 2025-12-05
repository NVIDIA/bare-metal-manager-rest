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
	"fmt"
	"net/http"

	"github.com/rs/zerolog/log"
	"go.opentelemetry.io/otel/attribute"

	"github.com/labstack/echo/v4"

	cdb "github.com/nvidia/carbide-rest/db/pkg/db"

	"github.com/nvidia/carbide-rest/api/pkg/api/handler/util/common"
	"github.com/nvidia/carbide-rest/api/pkg/api/model"
	auth "github.com/nvidia/carbide-rest/auth/pkg/authorization"
	cerr "github.com/nvidia/carbide-rest/common/pkg/util"
	sutil "github.com/nvidia/carbide-rest/common/pkg/util"
)

// GetUserHandler is an API Handler to return information about the current user
type GetUserHandler struct {
	dbSession  *cdb.Session
	tracerSpan *sutil.TracerSpan
}

// NewGetUserHandler creates and returns a new handler
func NewGetUserHandler(dbSession *cdb.Session) GetUserHandler {
	return GetUserHandler{
		dbSession:  dbSession,
		tracerSpan: sutil.NewTracerSpan(),
	}
}

// Handle godoc
// @Summary Return information about the current user
// @Description Get basic information about the user making the request
// @Tags user
// @Accept */*
// @Produce json
// @Security ApiKeyAuth
// @Param org path string true "Name of NGC organization"
// @Success 200 {object} model.APIUser
// @Router /v2/org/{org}/carbide/user/current [get]
func (guh GetUserHandler) Handle(c echo.Context) error {
	// Get context
	ctx := c.Request().Context()

	// Get org
	org := c.Param("orgName")

	// Initialize logger
	logger := log.With().Str("Model", "User").Str("Handler", "Get").Str("Org", org).Logger()

	logger.Info().Msg("started API handler")

	// Create a child span and set the attributes for current request
	newctx, handlerSpan := guh.tracerSpan.CreateChildInContext(ctx, "GetUserHandler", logger)
	if handlerSpan != nil {
		// Set newly created span context as a current context
		ctx = newctx

		defer handlerSpan.End()

		guh.tracerSpan.SetAttribute(handlerSpan, attribute.String("org", org), logger)
	}

	dbUser, logger, err := common.GetUserAndEnrichLogger(c, logger, guh.tracerSpan, handlerSpan)
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

	apiUser := model.NewAPIUserFromDBUser(*dbUser)

	logger.Info().Msg("finishing API handler")

	return c.JSON(http.StatusOK, apiUser)
}
