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
	"github.com/nvidia/carbide-rest/common/pkg/util"
	cdb "github.com/nvidia/carbide-rest/db/pkg/db"
	cdbm "github.com/nvidia/carbide-rest/db/pkg/db/model"
)

// KeycloakProcessor processes Keycloak JWT tokens
type CustomProcessor struct {
	dbSession *cdb.Session
}

// Ensure KeycloakProcessor implements config.TokenProcessor interface
var _ config.TokenProcessor = (*CustomProcessor)(nil)

// HandleToken processes Keycloak JWT tokens
func (h *CustomProcessor) ProcessToken(c echo.Context, tokenStr string, jwksConfig *config.JwksConfig, logger zerolog.Logger) (*cdbm.User, *util.APIError) {
	// Use map claims to be able to extract custom claims like scopes
	claims := jwt.MapClaims{}

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

	// Extract claims from the token
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || claims == nil {
		logger.Error().Msg("claims are nil after type assertion")
		return nil, util.NewAPIError(http.StatusUnauthorized, "Invalid claims in authorization token", nil)
	}

	// Extract necessary information from claims
	sub, _ := token.Claims.GetSubject()
	if sub == "" {
		return nil, util.NewAPIError(http.StatusUnauthorized, "Invalid authorization token, could not find subject ID in claim", nil)
	}

	// Validate audiences if configured
	if len(jwksConfig.Audiences) > 0 {
		if err := validateAudiences(claims, jwksConfig.Audiences, logger); err != nil {
			return nil, err
		}
	}

	// Validate scopes if configured
	if len(jwksConfig.Scopes) > 0 {
		if err := validateScopes(claims, jwksConfig.Scopes, logger); err != nil {
			return nil, err
		}
	}

	// check if service accounts are enabled, if not return an error
	if !jwksConfig.ServiceAccount {
		logger.Error().Msg("Service account detected but service accounts are not enabled")
		return nil, util.NewAPIError(http.StatusUnauthorized, "Service accounts are not enabled", nil)
	}

	auxId := sub
	// Extract issuer from claims
	issuer, _ := claims.GetIssuer()
	firstName := issuer
	lastName := ""
	email := ""
	orgData := cdbm.OrgData{
		jwksConfig.Name: cdbm.Org{
			Name:        jwksConfig.Name, // Issuer name in config represents the org name
			DisplayName: jwksConfig.Name,
			OrgType:     "Enterprise",
			Roles:       []string{"FORGE_PROVIDER_ADMIN", "FORGE_TENANT_ADMIN"},
			Teams:       []cdbm.Team{},
		}}

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
		if created {
			logger.Info().Str("userid", dbUser.ID.String()).Msg("updating newly created user with profile information")
		} else {
			logger.Info().Str("userid", dbUser.ID.String()).Msg("updating user with new information")
		}

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

// validateAudiences checks if the token's audience claim contains at least one of the configured audiences
func validateAudiences(claims jwt.MapClaims, configuredAudiences []string, logger zerolog.Logger) *util.APIError {
	// If no audiences are configured, skip validation
	if len(configuredAudiences) == 0 {
		return nil
	}

	// Extract audience claim from token
	audClaim, exists := claims["aud"]
	if !exists {
		logger.Error().Msg("Token does not contain audience claim")
		return util.NewAPIError(http.StatusUnauthorized,
			"Token missing required audience claim. Token must contain 'aud' claim", nil)
	}

	var tokenAudiences []string

	// Handle both string and array audience claims
	switch aud := audClaim.(type) {
	case string:
		tokenAudiences = []string{aud}
	case []interface{}:
		for _, a := range aud {
			if audStr, ok := a.(string); ok {
				tokenAudiences = append(tokenAudiences, audStr)
			}
		}
	case []string:
		tokenAudiences = aud
	default:
		logger.Error().Msgf("Unexpected audience claim type: %T", audClaim)
		return util.NewAPIError(http.StatusUnauthorized,
			"Invalid audience claim format. Expected string or array of strings", nil)
	}

	// Check if at least one configured audience matches the token audiences
	for _, configAud := range configuredAudiences {
		for _, tokenAud := range tokenAudiences {
			if configAud == tokenAud {
				return nil
			}
		}
	}

	logger.Error().
		Strs("token_audiences", tokenAudiences).
		Strs("configured_audiences", configuredAudiences).
		Msg("Token audience does not match any configured audiences")
	return util.NewAPIError(http.StatusUnauthorized,
		"Token audience does not match required audiences. Token audiences: ["+strings.Join(tokenAudiences, ", ")+
			"], Required audiences: ["+strings.Join(configuredAudiences, ", ")+"]", nil)
}

// validateScopes checks if the token contains all required scopes
func validateScopes(claims jwt.MapClaims, requiredScopes []string, logger zerolog.Logger) *util.APIError {
	// If no scopes are required, skip validation
	if len(requiredScopes) == 0 {
		return nil
	}

	var tokenScopes []string

	scopeClaim, scopeExists := claims["scope"]
	if !scopeExists {
		scopeClaim, scopeExists = claims["scopes"]
	}
	if !scopeExists {
		scopeClaim, scopeExists = claims["scp"]
	}

	if !scopeExists {
		logger.Error().Msg("Token does not contain scope claim (scope, scopes, or scp)")
		return util.NewAPIError(http.StatusUnauthorized,
			"Token missing required scope claim. Token must contain 'scope', 'scopes', or 'scp' claim", nil)
	}

	// Handle different scope claim formats
	switch scope := scopeClaim.(type) {
	case string:
		// Space-separated string of scopes
		if scope != "" {
			tokenScopes = strings.Fields(scope)
		}
	case []interface{}:
		// Array of scopes
		for _, s := range scope {
			if scopeStr, ok := s.(string); ok {
				tokenScopes = append(tokenScopes, scopeStr)
			}
		}
	case []string:
		tokenScopes = scope
	default:
		logger.Error().Msgf("Unexpected scope claim type: %T", scopeClaim)
		return util.NewAPIError(http.StatusUnauthorized,
			"Invalid scope claim format. Expected string (space-separated) or array of strings", nil)
	}

	// Check if token contains all required scopes
	tokenScopeSet := make(map[string]bool)
	for _, scope := range tokenScopes {
		tokenScopeSet[scope] = true
	}

	missingScopes := []string{}
	for _, requiredScope := range requiredScopes {
		if !tokenScopeSet[requiredScope] {
			missingScopes = append(missingScopes, requiredScope)
		}
	}

	if len(missingScopes) > 0 {
		logger.Error().
			Strs("token_scopes", tokenScopes).
			Strs("required_scopes", requiredScopes).
			Strs("missing_scopes", missingScopes).
			Msg("Token is missing required scopes")
		return util.NewAPIError(http.StatusUnauthorized,
			"Token missing required scopes. Missing: ["+strings.Join(missingScopes, ", ")+
				"], Required: ["+strings.Join(requiredScopes, ", ")+"]", nil)
	}

	return nil
}
