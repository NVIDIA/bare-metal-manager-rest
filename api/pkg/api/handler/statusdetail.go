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
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
	"github.com/nvidia/carbide-rest/api/pkg/api/model"
	"github.com/nvidia/carbide-rest/api/pkg/api/pagination"
	cerr "github.com/nvidia/carbide-rest/common/pkg/util"
	cdb "github.com/nvidia/carbide-rest/db/pkg/db"
	cdbm "github.com/nvidia/carbide-rest/db/pkg/db/model"
	"net/http"
)

func handleEntityStatusDetails(ctx context.Context, echoCtx echo.Context, dbSession *cdb.Session, entityID string, logger zerolog.Logger) ([]model.APIStatusDetail, error) {
	// Validate pagination request
	pageRequest := pagination.PageRequest{}
	err := echoCtx.Bind(&pageRequest)
	if err != nil {
		logger.Warn().Err(err).Msg("error binding pagination request data into API model")
		return nil, cerr.NewAPIErrorResponse(echoCtx, http.StatusBadRequest, "Failed to parse request pagination data", nil)
	}

	// Validate request attributes
	err = pageRequest.Validate(cdbm.StatusDetailOrderByFields)
	if err != nil {
		logger.Warn().Err(err).Msg("error validating pagination request data")
		return nil, cerr.NewAPIErrorResponse(echoCtx, http.StatusBadRequest, "Failed to validate pagination request data", err)
	}

	// Get status details
	sdDAO := cdbm.NewStatusDetailDAO(dbSession)
	dbSds, total, serr := sdDAO.GetAllByEntityIDs(ctx, nil, []string{entityID}, pageRequest.Offset, pageRequest.Limit, pageRequest.OrderBy)
	if serr != nil {
		logger.Error().Err(serr).Msg("error retrieving Status Details")
		return nil, cerr.NewAPIErrorResponse(echoCtx, http.StatusInternalServerError, "Failed to retrieve Status Details", nil)
	}

	// Create response
	apiSds := []model.APIStatusDetail{}
	for _, sd := range dbSds {
		apiSds = append(apiSds, model.NewAPIStatusDetail(sd))
	}

	// Create pagination response header
	pageReponse := pagination.NewPageResponse(*pageRequest.PageNumber, *pageRequest.PageSize, total, pageRequest.OrderByStr)
	pageHeader, err := json.Marshal(pageReponse)
	if err != nil {
		logger.Error().Err(err).Msg("error marshaling pagination response")
		return nil, cerr.NewAPIErrorResponse(echoCtx, http.StatusInternalServerError, "Failed to generate pagination response header", nil)
	}

	echoCtx.Response().Header().Set(pagination.ResponseHeaderName, string(pageHeader))

	return apiSds, nil
}
