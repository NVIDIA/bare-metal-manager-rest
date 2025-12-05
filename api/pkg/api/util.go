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


package api

import (
	"strings"

	"github.com/labstack/echo/v4"
)

// RequestHandler defines the Echo compatible interface all API route handlers
// should implement
type RequestHandler interface {
	Handle(c echo.Context) error
}

// Route defines the data structure to organize route information that can
// be used to initialize Echo routes
type Route struct {
	Path    string
	Method  string
	Handler RequestHandler
}

// MetricsURLSkipper ignores metrics for certain routes
func MetricsURLSkipper(c echo.Context) bool {
	// Allow v2 API paths to be tracked
	if strings.HasPrefix(c.Path(), "/v2/") {
		return false
	}

	if c.Path() == "/metrics" {
		return false
	}

	return true
}
