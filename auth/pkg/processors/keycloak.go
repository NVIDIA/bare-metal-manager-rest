// SPDX-FileCopyrightText: Copyright (c) 2021-2023 NVIDIA CORPORATION & AFFILIATES. All rights reserved.
// SPDX-License-Identifier: LicenseRef-NvidiaProprietary
//
// NVIDIA CORPORATION, its affiliates and licensors retain all intellectual
// property and proprietary rights in and to this material, related
// documentation and any modifications thereto. Any use, reproduction,
// disclosure or distribution of this material and related documentation
// without an express license agreement from NVIDIA CORPORATION or
// its affiliates is strictly prohibited.

package processors

import (
	"context"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
	"github.com/nvidia/carbide-rest/auth/pkg/config"
	"github.com/nvidia/carbide-rest/auth/pkg/core/claim"
	"github.com/nvidia/carbide-rest/common/pkg/util"
	cdb "github.com/nvidia/carbide-rest/db/pkg/db"
	cdbm "github.com/nvidia/carbide-rest/db/pkg/db/model"
)

// KeycloakProcessor processes Keycloak JWT tokens
type KeycloakProcessor struct {
	dbSession      *cdb.Session
	keycloakConfig *config.KeycloakConfig
}

// Ensure KeycloakProcessor implements config.TokenProcessor interface
var _ config.TokenProcessor = (*KeycloakProcessor)(nil)

// HandleToken processes Keycloak JWT tokens
func (h *KeycloakProcessor) ProcessToken(c echo.Context, tokenStr string, jwksConfig *config.JwksConfig, logger zerolog.Logger) (*cdbm.User, *util.APIError) {
	claims := &claim.KeycloakClaims{}

	token, err := jwksConfig.ValidateToken(tokenStr, claims)
	if err != nil {
		if strings.Contains(err.Error(), jwt.ErrTokenExpired.Error()) {
			logger.Error().Err(err).Msg("Token expired")
			return nil, util.NewAPIError(http.StatusUnauthorized, "Authorization token in request has expired", nil)
		} else {
			logger.Error().Err(err).Msg("failed to validate JWT token in authorization header")
			return nil, util.NewAPIError(http.StatusUnauthorized, "Invalid authorization token in request", nil)
		}
	}

	// Keycloak token, extract claims from the token
	claims, ok := token.Claims.(*claim.KeycloakClaims)
	if !ok || claims == nil {
		logger.Error().Msg("claims are nil after type assertion")
		return nil, util.NewAPIError(http.StatusUnauthorized, "Invalid claims in authorization token", nil)
	}

	// Extract necessary information from claims
	sub, _ := token.Claims.GetSubject()
	if sub == "" {
		return nil, util.NewAPIError(http.StatusUnauthorized, "Invalid authorization token, could not find subject ID in claim", nil)
	}

	email := claims.GetEmail()
	firstName := claims.FirstName
	lastName := claims.LastName
	auxId := claims.GetOidcId()
	if claims.GetClientId() != "" {
		// indicates service account, check if service accounts are enabled
		if !jwksConfig.ServiceAccount {
			logger.Error().Str("clientID", claims.GetClientId()).Msg("Service account detected but service accounts are not enabled")
			return nil, util.NewAPIError(http.StatusUnauthorized, "Service accounts are not enabled", nil)
		}
		// use sub as auxId
		auxId = sub
		firstName = claims.GetClientId()
	}

	orgData := claims.ToOrgData()

	if len(orgData) == 0 {
		return nil, util.NewAPIError(http.StatusForbidden, "User does not have any roles assigned", nil)
	}

	userDAO := cdbm.NewUserDAO(h.dbSession)
	dbUser, created, err := userDAO.GetOrCreate(context.Background(), nil, cdbm.UserGetOrCreateInput{
		AuxiliaryID: &auxId,
	})
	if err != nil {
		logger.Error().Err(err).Msg("failed to get or create user by oidc_id in DB")
		return nil, util.NewAPIError(http.StatusUnauthorized, "Failed to retrieve or create user record, DB error", nil)
	}

	// If user was created or needs updates, update with latest information
	needsUpdate := created || !dbUser.OrgData.Equal(orgData)
	if needsUpdate {
		// Regular update is sufficient since we're updating by UserID (primary key)
		dbUser, err = userDAO.Update(context.Background(), nil, cdbm.UserUpdateInput{
			UserID:    dbUser.ID,
			Email:     &email,
			FirstName: &firstName,
			LastName:  &lastName,
			OrgData:   orgData,
		})
		if err != nil {
			logger.Error().Err(err).Msg("failed to update user in DB")
			return nil, util.NewAPIError(http.StatusUnauthorized, "Failed to update user record, DB error", nil)
		}
	}

	// Set user in context
	c.Set("user", dbUser)
	return dbUser, nil
}
