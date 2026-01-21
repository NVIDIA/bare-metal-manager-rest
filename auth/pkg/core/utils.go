// SPDX-FileCopyrightText: Copyright (c) 2021-2023 NVIDIA CORPORATION & AFFILIATES. All rights reserved.
// SPDX-License-Identifier: LicenseRef-NvidiaProprietary
//
// NVIDIA CORPORATION, its affiliates and licensors retain all intellectual
// property and proprietary rights in and to this material, related
// documentation and any modifications thereto. Any use, reproduction,
// disclosure or distribution of this material and related documentation
// without an express license agreement from NVIDIA CORPORATION or
// its affiliates is strictly prohibited.

package core

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/cast"
)

// =============================================================================
// Constants
// =============================================================================

// ScopeClaims are the standard JWT claim keys used for scopes.
var ScopeClaims = []string{"scope", "scopes", "scp"}

// =============================================================================
// Conversion Functions
// =============================================================================

// InterfaceToStringSlice converts interface{} to []string.
// Supports multiple common formats from various IdPs:
//   - Native array/slice: ["role1", "role2"]
//   - JSON-encoded string array: "[\"role1\", \"role2\"]"
//   - Space-separated: "role1 role2"
//   - Comma-separated: "role1,role2" or "role1, role2"
//   - Semicolon-separated: "role1;role2"
//   - Single value: "role1"
func InterfaceToStringSlice(v any) ([]string, error) {
	if v == nil {
		return nil, nil
	}

	// Handle string values with various formats
	if s, ok := v.(string); ok {
		return parseStringToSlice(s), nil
	}

	// Handle native arrays/slices
	return cast.ToStringSliceE(v)
}

// parseStringToSlice parses a string into a slice using common delimiters.
// Tries formats in order: JSON array, comma-separated, semicolon-separated, space-separated.
func parseStringToSlice(s string) []string {
	trimmed := strings.TrimSpace(s)
	if trimmed == "" {
		return nil
	}

	// Try JSON array format first: ["role1", "role2"]
	if strings.HasPrefix(trimmed, "[") && strings.HasSuffix(trimmed, "]") {
		var jsonArray []string
		if err := json.Unmarshal([]byte(trimmed), &jsonArray); err == nil {
			return trimAndFilter(jsonArray)
		}
		// If JSON parsing fails, fall through to other methods
	}

	// Try comma-separated: "role1,role2" or "role1, role2"
	if strings.Contains(trimmed, ",") {
		parts := strings.Split(trimmed, ",")
		return trimAndFilter(parts)
	}

	// Try semicolon-separated: "role1;role2"
	if strings.Contains(trimmed, ";") {
		parts := strings.Split(trimmed, ";")
		return trimAndFilter(parts)
	}

	// Try space/tab/newline-separated: "role1 role2"
	if strings.ContainsAny(trimmed, " \t\n") {
		return strings.Fields(trimmed)
	}

	// Single value
	return []string{trimmed}
}

// trimAndFilter trims whitespace from each element and removes empty strings.
func trimAndFilter(parts []string) []string {
	result := make([]string, 0, len(parts))
	for _, p := range parts {
		trimmed := strings.TrimSpace(p)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}

// ComputeIssuerPrefix returns SHA256(issuerURL)[0:10] for namespacing subject claims.
func ComputeIssuerPrefix(issuerURL string) string {
	hash := sha256.Sum256([]byte(issuerURL))
	return hex.EncodeToString(hash[:])[:10]
}

// =============================================================================
// Claim Extraction Functions
// =============================================================================

// ExtractClaimValue extracts any value from a nested claim path (e.g., "data.roles").
// Returns nil if the path is empty or the value is not found.
func ExtractClaimValue(claims jwt.MapClaims, path string) any {
	if path == "" {
		return nil
	}

	var current any = claims

	for _, key := range strings.Split(path, ".") {
		switch m := current.(type) {
		case jwt.MapClaims:
			current = m[key]
		case map[string]any:
			current = m[key]
		default:
			return nil
		}

		if current == nil {
			return nil
		}
	}

	return current
}

// ExtractStringClaim extracts a string from a nested claim path (e.g., "data.org").
// Returns empty string if not found or if the value is not a string.
func ExtractStringClaim(claims jwt.MapClaims, path string) string {
	value := ExtractClaimValue(claims, path)
	if str, ok := value.(string); ok {
		return str
	}
	return ""
}

// ExtractTokenScopes extracts scopes from claims (tries "scope", "scopes", "scp").
// Returns a slice of scope strings.
func ExtractTokenScopes(claims jwt.MapClaims) []string {
	var scopeClaimValue any
	for _, key := range ScopeClaims {
		if val, exists := claims[key]; exists {
			scopeClaimValue = val
			break
		}
	}
	if scopeClaimValue == nil {
		return nil
	}
	scopes, _ := InterfaceToStringSlice(scopeClaimValue)
	return scopes
}
