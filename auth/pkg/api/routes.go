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
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"

	"github.com/nvidia/carbide-rest/auth/pkg/authentication"
	"github.com/nvidia/carbide-rest/auth/pkg/config"

	cah "github.com/nvidia/carbide-rest/auth/pkg/api/handler"
)

// AuthRoute represents an authentication route
type AuthRoute struct {
	Method  string
	Path    string
	Handler AuthHandler
}

// AuthHandler represents an authentication handler
type AuthHandler interface {
	Handle(c echo.Context) error
}

// NewAuthRoutes creates new authentication routes using cloud-auth services
// This function provides a complete authentication API that can be used by any service
func NewAuthRoutes(keycloakConfig *config.KeycloakConfig) []AuthRoute {
	if keycloakConfig == nil {
		log.Error().Msg("keycloak config is not initialized, cannot create authentication routes")
		return nil
	}

	// Initialize Keycloak auth service
	keycloakAuth := authentication.NewKeycloakAuthService(keycloakConfig)

	// Create handlers
	loginHandler := cah.NewLoginHandler(keycloakAuth)
	callbackHandler := cah.NewCallbackHandler(keycloakAuth)
	logoutHandler := cah.NewLogoutHandler(keycloakAuth)
	refreshHandler := cah.NewRefreshTokenHandler(keycloakAuth)

	return []AuthRoute{
		{
			Method:  "POST",
			Path:    "/login",
			Handler: loginHandler,
		},
		{
			Method:  "POST",
			Path:    "/callback",
			Handler: callbackHandler,
		},
		{
			Method:  "POST",
			Path:    "/logout",
			Handler: logoutHandler,
		},
		{
			Method:  "POST",
			Path:    "/refresh",
			Handler: refreshHandler,
		},
	}
}
