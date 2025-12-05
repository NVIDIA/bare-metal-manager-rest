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

// MetadataHandler is an API handler to return system information about the API
type MetadataHandler struct{}

// NewMetadataHandler creates and returns a new handler
func NewMetadataHandler() MetadataHandler {
	return MetadataHandler{}
}

// Handle godoc
// @Summary Returns system information about the API
// @Description Returns system information about the API
// @Tags metadata
// @Accept */*
// @Produce json
// @Success 200 {object} model.APIMetadata
// @Router /v2/org/{org}/carbide/metadata [get]
func (mdh MetadataHandler) Handle(c echo.Context) error {
	amd := model.NewAPIMetadata()
	return c.JSON(http.StatusOK, amd)
}
