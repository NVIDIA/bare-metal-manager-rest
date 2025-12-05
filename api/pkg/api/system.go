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
	"net/http"

	apiHandler "github.com/nvidia/carbide-rest/api/pkg/api/handler"
)

// NewSystemAPIRoutes returns API routes that provide system level  functions
func NewSystemAPIRoutes() []Route {
	apiRoutes := []Route{
		// Health check endpoints
		{
			Path:    "/healthz",
			Method:  http.MethodGet,
			Handler: apiHandler.NewHealthCheckHandler(),
		},
		{
			Path:    "/readyz",
			Method:  http.MethodGet,
			Handler: apiHandler.NewHealthCheckHandler(),
		},
	}

	return apiRoutes
}

// IsSystemRoute returns true for a path registered as SystemAPIRoute
func IsSystemRoute(p string) bool {
	routes := NewSystemAPIRoutes()
	for _, r := range routes {
		if r.Path == p {
			return true
		}
	}

	return false
}
