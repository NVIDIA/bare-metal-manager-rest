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
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/nvidia/carbide-rest/api/pkg/api/model"
)

// HealthCheckHandler is an API handler to return health status of the API server
type HealthCheckHandler struct{}

// NewHealthCheckHandler creates and returns a new handler
func NewHealthCheckHandler() HealthCheckHandler {
	return HealthCheckHandler{}
}

// Handle godoc
// @Summary Returns the health status of API server
// @Description Returns the health status of the API server
// @Tags health
// @Accept */*
// @Produce json
// @Success 200 {object} model.APIHealthCheck
// @Router /healthz [get]
func (hch HealthCheckHandler) Handle(c echo.Context) error {
	ahc := model.NewAPIHealthCheck(true, nil)
	return c.JSON(http.StatusOK, ahc)
}
