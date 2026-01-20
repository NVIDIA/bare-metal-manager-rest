// SPDX-FileCopyrightText: Copyright (c) 2021-2023 NVIDIA CORPORATION & AFFILIATES. All rights reserved.
// SPDX-License-Identifier: LicenseRef-NvidiaProprietary
//
// NVIDIA CORPORATION, its affiliates and licensors retain all intellectual
// property and proprietary rights in and to this material, related
// documentation and any modifications thereto. Any use, reproduction,
// disclosure or distribution of this material and related documentation
// without an express license agreement from NVIDIA CORPORATION or
// its affiliates is strictly prohibited.

package config

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/nvidia/carbide-rest/auth/pkg/core"
)

// =============================================================================
// Role Constants and Variables
// =============================================================================

var (
	// ServiceAccountRoles are the default roles assigned to service accounts
	ServiceAccountRoles = []string{"FORGE_PROVIDER_ADMIN", "FORGE_TENANT_ADMIN"}

	// AllowedRoles is the set of valid roles that can be assigned to users.
	// Both static roles in config and dynamic roles from claims must be from this set.
	AllowedRoles = map[string]bool{
		"FORGE_TENANT_ADMIN":   true,
		"FORGE_PROVIDER_ADMIN": true,
	}
)

// =============================================================================
// Role Validation Functions
// =============================================================================

// validateRoles checks that all roles are in the AllowedRoles set.
// Returns false immediately upon finding the first invalid role.
func validateRoles(roles []string) bool {
	for _, role := range roles {
		if !AllowedRoles[role] {
			return false
		}
	}
	return true
}

// IsValidRole checks if a single role is in the AllowedRoles set.
func IsValidRole(role string) bool {
	return AllowedRoles[role]
}

// FilterToAllowedRoles filters a list of roles to only include allowed roles.
// Returns core.ErrInvalidRole if no valid roles remain after filtering.
func FilterToAllowedRoles(roles []string) (allowed []string, err error) {
	for _, role := range roles {
		if AllowedRoles[role] {
			allowed = append(allowed, role)
		}
	}
	if len(allowed) == 0 {
		return nil, core.ErrInvalidRole
	}
	return allowed, nil
}

// =============================================================================
// Role Extraction Functions
// =============================================================================

// ExtractRolesFromClaimPath extracts roles from a nested claim path and filters to allowed roles.
// Returns nil if the path doesn't exist or contains no valid roles.
func ExtractRolesFromClaimPath(claims jwt.MapClaims, path string) ([]string, error) {
	value := core.ExtractClaimValue(claims, path)
	if value == nil {
		return nil, nil
	}

	roles, err := core.InterfaceToStringSlice(value)
	if err != nil {
		return nil, err
	}
	return FilterToAllowedRoles(roles)
}
