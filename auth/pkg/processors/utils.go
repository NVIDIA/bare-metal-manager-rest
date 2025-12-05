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
	"encoding/base64"
	"net/http"
	"sort"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
	"github.com/nvidia/carbide-rest/common/pkg/util"
	cdbm "github.com/nvidia/carbide-rest/db/pkg/db/model"
)

const (
	// NGC KAS Headers
	// Kas v2 NGC team header
	ngcTeamHeader = "NV-Ngc-Team"
	// Kas v2 NGC roles header
	ngcRolesHeader = "NV-Ngc-User-Roles"
	// Kas v2 NGC org display name header
	ngcOrgDisplayNameHeader = "NV-Ngc-Org-Display-Name"
	// Kas v2 NGC user name header
	ngcUserNameHeader = "NV-Ngc-User-Name"
	// Kas v2 NGC user email header
	ngcUserEmailHeader = "X-Ngc-Email-Id"
)

// GetUpdatedUserFromHeaders extracts user information from headers sent by KAS
// Steps include
// 1. Extract NGC user name and email from headers
// 2. Extract NGC roles from headers
// 3. Extract NGC org display name from headers
// 4. Update user record if necessary
// 5. Return updated user record
// Returns updated user record and API error if any
func GetUpdatedUserFromHeaders(c echo.Context, existingUser cdbm.User, ngcOrgName string, logger zerolog.Logger) (*cdbm.User, *util.APIError) {
	// Update user record if necessary
	isUserUpdated := false
	updatedUser := &cdbm.User{}

	// Extract NGC user name
	ngcUserNameB64 := c.Request().Header.Get(ngcUserNameHeader)
	if ngcUserNameB64 != "" {
		// NGC User Name is base64 encoded, decode it
		decodedBytes, err := base64.StdEncoding.DecodeString(ngcUserNameB64)
		if err != nil {
			logger.Warn().Err(err).Msg("failed to decode NGC user name header, invalid base64 value")
			return nil, util.NewAPIError(http.StatusUnauthorized, "Invalid value in NGC org user name header", nil)
		}
		ngcUserName := string(decodedBytes)
		nameComps := strings.SplitN(ngcUserName, " ", 2)
		if len(nameComps) > 0 && nameComps[0] != "" && (existingUser.FirstName == nil || *existingUser.FirstName != nameComps[0]) {
			updatedUser.FirstName = &nameComps[0]
			isUserUpdated = true
		}
		if len(nameComps) > 1 && nameComps[1] != "" && (existingUser.LastName == nil || *existingUser.LastName != nameComps[1]) {
			updatedUser.LastName = &nameComps[1]
			isUserUpdated = true
		}
	} else {
		logger.Warn().Msg("request received without NGC user name header, first/last name may not be available for user")
	}

	// Extract NGC user email
	ngcUserEmailB64 := c.Request().Header.Get(ngcUserEmailHeader)
	if ngcUserEmailB64 != "" {
		// NGC User Email is base64 encoded, decode it
		decodedBytes, err := base64.StdEncoding.DecodeString(ngcUserEmailB64)
		if err != nil {
			logger.Warn().Err(err).Msg("failed to decode NGC user email header, invalid base64 value")
			return nil, util.NewAPIError(http.StatusUnauthorized, "Invalid value in NGC org user email header", nil)
		}
		ngcUserEmail := string(decodedBytes)
		if existingUser.Email == nil || *existingUser.Email != ngcUserEmail {
			updatedUser.Email = &ngcUserEmail
			isUserUpdated = true
		}
	} else {
		logger.Warn().Msg("request received without NGC user email header, email may not be available for user")
	}

	// Extract NGC roles
	ngcRolesValue := c.Request().Header.Get(ngcRolesHeader)
	if ngcRolesValue == "" {
		logger.Warn().Msg("request received without NGC roles header, access denied")
		return nil, util.NewAPIError(http.StatusUnauthorized, "Request is missing NGC roles header", nil)
	}

	ngcRoles := strings.Split(ngcRolesValue, ",")
	newNgcRoles := []string{}
	// Format roles
	for _, role := range ngcRoles {
		curRole := strings.ReplaceAll(role, "-", "_")
		curRole = strings.ToUpper(curRole)
		newNgcRoles = append(newNgcRoles, curRole)
	}
	sort.Strings(newNgcRoles)

	var OrgData cdbm.OrgData
	if existingUser.OrgData != nil {
		OrgData = existingUser.OrgData
	} else {
		OrgData = cdbm.OrgData{}
	}

	// Extract NGC org display name
	var ngcOrgDisplayName string

	ngcOrgDisplayNameB64 := c.Request().Header.Get(ngcOrgDisplayNameHeader)
	if ngcOrgDisplayNameB64 == "" {
		logger.Warn().Msg("request received without NGC org display name header, access denied")
		return nil, util.NewAPIError(http.StatusUnauthorized, "Request is missing NGC org display name header", nil)
	} else {
		// NGC Org Display Name is base64 encoded, decode it
		decodedBytes, err := base64.StdEncoding.DecodeString(ngcOrgDisplayNameB64)
		if err != nil {
			logger.Error().Err(err).Msg("failed to decode NGC org display name header")
			return nil, util.NewAPIError(http.StatusUnauthorized, "Invalid value in NGC org display name header", nil)
		}
		ngcOrgDisplayName = string(decodedBytes)
	}

	ngcOrg, err := existingUser.OrgData.GetOrgByName(ngcOrgName)
	if err != nil {
		// Org not found, create new
		ngcOrg = &cdbm.Org{
			Name:        ngcOrgName,
			DisplayName: ngcOrgDisplayName,
			Roles:       newNgcRoles,
			Teams:       []cdbm.Team{},
		}
		OrgData[ngcOrgName] = *ngcOrg
		updatedUser.OrgData = OrgData
		isUserUpdated = true
	} else {
		// Check if user has any role changes
		updateRoles := len(ngcOrg.Roles) != len(newNgcRoles)
		if !updateRoles {
			existingRoleMap := map[string]bool{}
			for _, role := range ngcOrg.Roles {
				existingRoleMap[role] = true
			}
			for _, role := range newNgcRoles {
				_, found := existingRoleMap[role]
				if !found {
					updateRoles = true
					break
				}
			}
		}
		if updateRoles {
			ngcOrg.Roles = newNgcRoles
			OrgData[ngcOrgName] = *ngcOrg
			updatedUser.OrgData = OrgData
			isUserUpdated = true
		}

		if ngcOrg.DisplayName != ngcOrgDisplayName {
			ngcOrg.DisplayName = ngcOrgDisplayName
			OrgData[ngcOrgName] = *ngcOrg
			updatedUser.OrgData = OrgData
			isUserUpdated = true
		}
	}

	if isUserUpdated {
		return updatedUser, nil
	} else {
		return nil, nil
	}
}
