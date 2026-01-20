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
	"context"
	"crypto/sha256"
	"encoding/hex"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	mapset "github.com/deckarep/golang-set/v2"
	"github.com/go-jose/go-jose/v4"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/nvidia/carbide-rest/auth/pkg/core"
	cdbm "github.com/nvidia/carbide-rest/db/pkg/db/model"
	"github.com/pkg/errors"
	"github.com/spf13/cast"
)

// Constants for JWKS configuration
const (
	minUpdateInterval = 10 * time.Second
)

var (
	// ErrJWKSURLEmpty is returned when JWKS URL is empty
	ErrJWKSURLEmpty = errors.New("JWKS URL is empty")
	// ErrJWKSNotInitialized is returned when JWKS has not been initialized
	ErrJWKSNotInitialized = errors.New("JWKS not initialized - call UpdateAllJWKS first")
	// ErrEmptyKeySet is returned when JWKS key set is empty
	ErrEmptyKeySet = errors.New("JWKS key set is empty")
	// ErrNoValidKeys is returned when JWKS contains no valid keys
	ErrNoValidKeys = errors.New("JWKS contains no valid keys")
	// ErrJWKSUpdateInProgress is returned when a JWKS update is already in progress
	ErrJWKSUpdateInProgress = errors.New("JWKS update already in progress")

	// ErrInvalidAudience is returned when token audience does not match (401)
	ErrInvalidAudience = errors.New("token audience does not match issuer configuration")
	// ErrInvalidConfiguration is returned when no claim mapping is configured (401)
	ErrInvalidConfiguration = errors.New("no claim mapping configured for requested organization")
	// ErrInvalidScope is returned when token scopes do not match (403)
	ErrInvalidScope = errors.New("token scopes do not match required scopes for issuer")
	// ErrNoClaimRoles is returned when no roles found in token claims (401)
	ErrNoClaimRoles = errors.New("no roles found in token claims for organization")
	// ErrReservedOrgName is returned when token claims a reserved organization name (403)
	ErrReservedOrgName = errors.New("token claims a reserved organization name")
	// ErrInvalidRole is returned when role is not in allowed roles set
	ErrInvalidRole = errors.New("role is not in allowed roles set")

	// ServiceAccountRoles are the default roles assigned to service accounts
	ServiceAccountRoles = []string{"FORGE_PROVIDER_ADMIN", "FORGE_TENANT_ADMIN"}

	// AllowedRoles is the set of valid roles that can be assigned to users.
	// Both static roles in config and dynamic roles from claims must be from this set.
	AllowedRoles = map[string]bool{
		"FORGE_TENANT_ADMIN":   true,
		"FORGE_PROVIDER_ADMIN": true,
	}

	isServiceAccountContextKey = AuthContextKey("isServiceAccount")
	scopeClaims                = []string{"scope", "scopes", "scp"}
)

// AuthContextKey is a custom type for context keys to avoid collisions
type AuthContextKey string

// SetIsServiceAccountInContext stores whether the request is from a service account
func SetIsServiceAccountInContext(c echo.Context, isServiceAccount bool) {
	ctx := context.WithValue(c.Request().Context(), isServiceAccountContextKey, isServiceAccount)
	c.SetRequest(c.Request().WithContext(ctx))
}

// GetIsServiceAccountFromContext returns whether the request is from a service account
func GetIsServiceAccountFromContext(c echo.Context) bool {
	v := c.Request().Context().Value(isServiceAccountContextKey)
	if v == nil {
		return false
	}
	b, ok := v.(bool)
	return ok && b
}

// computeIssuerPrefix returns SHA256(issuerURL)[0:10] for namespacing subject claims.
func computeIssuerPrefix(issuerURL string) string {
	hash := sha256.Sum256([]byte(issuerURL))
	return hex.EncodeToString(hash[:])[:10]
}

// ClaimMapping defines how to map JWT claims to organization data.
// Dynamic mode: set OrgAttribute to extract org from token claims.
// Static mode: set OrgName for a fixed organization.
type ClaimMapping struct {
	// OrgAttribute: JWT claim path to extract org name (e.g., "org", "data.org"). Makes this a dynamic mapping.
	OrgAttribute string `mapstructure:"orgAttribute"`
	// OrgDisplayAttribute: JWT claim path for org display name (dynamic mappings only)
	OrgDisplayAttribute string `mapstructure:"orgDisplayAttribute"`

	// OrgName: fixed organization name (static mapping). Used when OrgAttribute is empty.
	OrgName string `mapstructure:"orgName"`
	// OrgDisplayName: display name for static org mappings
	OrgDisplayName string `mapstructure:"orgDisplayName"`

	// RolesAttribute: JWT claim path to extract roles (e.g., "roles", "data.roles"). Takes precedence over Roles.
	RolesAttribute string `mapstructure:"rolesAttribute"`
	// Roles: static role list. Used when RolesAttribute is empty and IsServiceAccount is false.
	Roles []string `mapstructure:"roles"`

	// IsServiceAccount: if true, assigns admin roles (FORGE_PROVIDER_ADMIN, FORGE_TENANT_ADMIN). Ignores RolesAttribute/Roles.
	IsServiceAccount bool `mapstructure:"isServiceAccount"`
}

// IsOrgDynamic returns true if this is a valid dynamic org mapping.
// Dynamic mappings require all three attributes:
//   - OrgAttribute: JWT claim path to extract org name (e.g., "org", "data.org")
//   - OrgDisplayAttribute: JWT claim path to extract org display name (e.g., "org_display", "data.orgDisplayName")
//   - RolesAttribute: JWT claim path to extract roles (e.g., "roles", "data.roles")
//
// Service accounts are not allowed with dynamic orgs.
func (cm *ClaimMapping) IsOrgDynamic() bool {
	return cm.OrgAttribute != "" && cm.RolesAttribute != "" && !cm.IsServiceAccount
}

// IsOrgStatic returns true if using a fixed org name (OrgName set).
func (cm *ClaimMapping) IsOrgStatic() bool { return cm.OrgName != "" }

// ValidateMapping validates the claim mapping configuration.
// Valid mapping types:
//   - StaticOrg-StaticRoles: OrgName + Roles
//   - StaticOrg-DynamicRoles: OrgName + RolesAttribute
//   - StaticOrg-ServiceAccount: OrgName + IsServiceAccount
//   - DynamicOrg-DynamicRoles: OrgAttribute (required) + RolesAttribute (required) + OrgDisplayAttribute (optional)
//
// Note: DynamicOrg requires DynamicRoles (rolesAttribute). Static roles and service accounts
// are not allowed with dynamic org because the org is determined at runtime from the token.
func (cm *ClaimMapping) ValidateMapping() bool {
	if cm.IsOrgDynamic() {
		return true // IsOrgDynamic already validates the dynamic mapping requirements
	}
	if cm.OrgAttribute != "" {
		// Has OrgAttribute but doesn't satisfy IsOrgDynamic - invalid dynamic mapping
		return false
	}
	if cm.OrgName == "" {
		return false
	}
	return cm.IsServiceAccount || cm.RolesAttribute != "" || len(cm.Roles) > 0 && validateRoles(cm.Roles)
}

// validateRoles checks that all roles are in the AllowedRoles set
// Returns false immediately upon finding the first invalid role
func validateRoles(roles []string) bool {
	for _, role := range roles {
		if !AllowedRoles[role] {
			return false
		}
	}
	return true
}

// IsValidRole checks if a single role is in the AllowedRoles set
func IsValidRole(role string) bool {
	return AllowedRoles[role]
}

// FilterToAllowedRoles filters a list of roles to only include allowed roles
func FilterToAllowedRoles(roles []string) (allowed []string, err error) {
	for _, role := range roles {
		if AllowedRoles[role] {
			allowed = append(allowed, role)
		}
	}
	if len(allowed) == 0 {
		return nil, ErrInvalidRole
	}
	return allowed, nil
}

// ExtractRolesFromClaims returns roles based on mapping config: service account roles, dynamic extraction, or static roles.
func (cm *ClaimMapping) ExtractRolesFromClaims(claims jwt.MapClaims) ([]string, error) {
	if cm.IsServiceAccount {
		return ServiceAccountRoles, nil
	}
	if cm.RolesAttribute != "" {
		return extractNestedClaim(claims, cm.RolesAttribute)
	}
	return cm.Roles, nil
}

func extractNestedClaim(claims jwt.MapClaims, path string) ([]string, error) {
	var current any = claims

	for _, key := range strings.Split(path, ".") {
		switch m := current.(type) {
		case jwt.MapClaims:
			current = m[key]
		case map[string]any:
			current = m[key]
		default:
			return nil, nil
		}

		if current == nil {
			return nil, nil
		}
	}

	roles, err := InterfaceToStringSlice(current)
	if err != nil {
		return nil, err
	}
	return FilterToAllowedRoles(roles)
}

// extractStringClaim extracts a string from a nested claim path (e.g., "data.org"). Returns "" if not found.
func extractStringClaim(claims jwt.MapClaims, path string) string {
	if path == "" {
		return ""
	}

	var current any = claims

	for _, key := range strings.Split(path, ".") {
		switch m := current.(type) {
		case jwt.MapClaims:
			current = m[key]
		case map[string]any:
			current = m[key]
		default:
			return ""
		}

		if current == nil {
			return ""
		}
	}

	if str, ok := current.(string); ok {
		return str
	}
	return ""
}

// ExtractOrgFromClaims extracts org and display name from claims (dynamic mappings only).
func (cm *ClaimMapping) ExtractOrgFromClaims(claims jwt.MapClaims) (orgName string, displayName string) {
	if !cm.IsOrgDynamic() {
		return "", ""
	}

	rawOrgName := extractStringClaim(claims, cm.OrgAttribute)
	orgName = strings.ToLower(rawOrgName)
	displayName = extractStringClaim(claims, cm.OrgDisplayAttribute)

	// If display name not found, use the original (non-lowercased) org name
	if displayName == "" && rawOrgName != "" {
		displayName = rawOrgName
	}

	return orgName, displayName
}

// InterfaceToStringSlice converts interface{} to []string. Handles space-separated strings, arrays, and slices.
func InterfaceToStringSlice(v any) ([]string, error) {
	if v == nil {
		return nil, nil
	}
	if s, ok := v.(string); ok && strings.ContainsAny(s, " \t\n") {
		return strings.Fields(s), nil
	}
	return cast.ToStringSliceE(v)
}

// JwksConfig holds configuration for a JWKS endpoint and token validation.
type JwksConfig struct {
	Name         string
	IsUpdating   uint32        // atomic flag for concurrent JWKS updates
	sync.RWMutex               // protects JWKS access
	URL          string        // JWKS endpoint URL
	Issuer       string        // expected "iss" claim value
	Origin       string        // token origin type (e.g., "kas-legacy", "kas-ssa", "keycloak", "custom")
	LastUpdated  time.Time     // last JWKS update timestamp
	jwks         *core.JWKS    // cached JWKS keys
	JWKSTimeout  time.Duration // fetch timeout (default: 5s)

	Audiences []string // allowed audience values (token must have at least one)
	Scopes    []string // required scopes (token must have ALL)

	ClaimMappings []ClaimMapping // org/role mapping configuration

	// ServiceAccount enables client credentials flow (Keycloak only).
	// For custom issuers, use ClaimMapping.IsServiceAccount instead.
	ServiceAccount bool

	// ReservedOrgNames prevents dynamic org mappings from claiming statically-configured org names.
	// Populated by carbide-rest-api during initialization.
	ReservedOrgNames map[string]bool

	subjectPrefix string // SHA256(issuer)[0:10] - namespaces subject claims
}

// GetKeyByID is a method that returns a JWK secret by ID with enhanced validation
func (jcfg *JwksConfig) GetKeyByID(id string) (interface{}, error) {
	// Validate input parameters
	if strings.TrimSpace(id) == "" {
		return nil, jwt.ErrInvalidKey
	}

	jcfg.RLock()
	defer jcfg.RUnlock()

	if jcfg.jwks == nil {
		return nil, ErrJWKSNotInitialized
	}

	key, err := jcfg.jwks.GetKeyByID(id)
	if err != nil {
		return nil, errors.Wrap(jwt.ErrInvalidKey, err.Error())
	}

	// Validate key using go-jose's built-in validation
	if !key.Valid() {
		return nil, errors.Wrapf(jose.ErrUnsupportedKeyType, "go-jose validation failed for key %s", id)
	}

	return key.Key, nil
}

// KeyCount returns the number of keys in the JWKS
func (jcfg *JwksConfig) KeyCount() int {
	jcfg.RLock()
	defer jcfg.RUnlock()

	if jcfg.jwks == nil || jcfg.jwks.Set == nil {
		return 0
	}

	return len(jcfg.jwks.Set.Keys)
}

// MatchesIssuer checks if the given issuer exactly matches the configured issuer
func (jcfg *JwksConfig) MatchesIssuer(issuer string) bool {
	if jcfg == nil {
		return false
	}

	jcfg.RLock()
	defer jcfg.RUnlock()

	if jcfg.Issuer == "" {
		return false
	}

	return issuer == jcfg.Issuer
}

// shouldAllowJWKSUpdate checks if we should allow JWKS update based on throttling
func (jcfg *JwksConfig) shouldAllowJWKSUpdate() bool {
	jcfg.RLock()
	defer jcfg.RUnlock()

	// Always allow if we've never updated
	if jcfg.LastUpdated.IsZero() {
		return true
	}

	// Allow if enough time has passed since last update (regardless of success/failure)
	return time.Since(jcfg.LastUpdated) >= minUpdateInterval
}

// UpdateJWKS fetches and validates JWKS from the configured URL. Throttled to minUpdateInterval.
func (jcfg *JwksConfig) UpdateJWKS() error {
	if jcfg.URL == "" {
		return ErrJWKSURLEmpty
	}
	if !jcfg.shouldAllowJWKSUpdate() {
		return nil
	}
	if !atomic.CompareAndSwapUint32(&jcfg.IsUpdating, 0, 1) {
		return ErrJWKSUpdateInProgress
	}
	defer atomic.StoreUint32(&jcfg.IsUpdating, 0)

	jcfg.RLock()
	urlCopy, timeout := jcfg.URL, jcfg.JWKSTimeout
	jcfg.RUnlock()

	jwks, err := core.NewJWKSFromURL(urlCopy, timeout)
	if err != nil {
		return errors.Wrapf(err, "failed to update JWKS from %s", urlCopy)
	}
	if jwks.Set == nil || len(jwks.Set.Keys) == 0 {
		return errors.Wrapf(ErrEmptyKeySet, "from %s", urlCopy)
	}

	validKeyCount := 0
	for _, key := range jwks.Set.Keys {
		if key.Valid() {
			validKeyCount++
		}
	}
	if validKeyCount == 0 {
		return errors.Wrapf(ErrNoValidKeys, "from %s", urlCopy)
	}

	jcfg.Lock()
	defer jcfg.Unlock()
	jcfg.jwks = jwks
	jcfg.LastUpdated = time.Now()
	return nil
}

// GetJWKS returns the enhanced JWKS with go-jose capabilities
func (jcfg *JwksConfig) GetJWKS() *core.JWKS {
	jcfg.RLock()
	defer jcfg.RUnlock()
	return jcfg.jwks
}

// ValidateToken parses token from Authorization header with caller-provided claims and enhanced validation
func (jcfg *JwksConfig) ValidateToken(authHeader string, claims jwt.Claims) (*jwt.Token, error) {
	// Validate input parameters
	if strings.TrimSpace(authHeader) == "" {
		return nil, jwt.ErrTokenMalformed
	}

	if claims == nil {
		return nil, jwt.ErrTokenMalformed
	}

	// Use a comprehensive set of common JWT algorithms instead of restricting to current JWKS
	// This allows tokens with algorithms that might become available after JWKS updates
	allCommonAlgorithms := []string{
		"RS256", "RS384", "RS512", // RSA with SHA
		"PS256", "PS384", "PS512", // RSA-PSS with SHA
		"ES256", "ES384", "ES512", // ECDSA with SHA
		"HS256", "HS384", "HS512", // HMAC with SHA
		"EdDSA", // Ed25519/Ed448
	}

	jwtParser := jwt.NewParser(jwt.WithValidMethods(allCommonAlgorithms))

	token, err := jwtParser.ParseWithClaims(authHeader, claims, jcfg.getPublicKey)
	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, jwt.ErrTokenInvalidClaims
	}

	return token, nil
}

// getPublicKey retrieves the public key from the JWKS for JWT validation
func (jcfg *JwksConfig) getPublicKey(token *jwt.Token) (interface{}, error) {
	if token == nil || token.Header == nil {
		return nil, jwt.ErrTokenMalformed
	}

	algorithm, _ := token.Header["alg"].(string)
	if algorithm == "" {
		return nil, jwt.ErrTokenMalformed
	}

	kid, _ := token.Header["kid"].(string)

	// If kid is present, use existing single-key logic
	if kid != "" {
		key, err := jcfg.getKeyFromJWKS(kid)
		if err != nil {
			// Attempt JWKS update with retry logic for concurrent updates
			if updateErr := jcfg.tryUpdateJWKSWithRetry(); updateErr == nil {
				key, err = jcfg.getKeyFromJWKS(kid)
			}
		}
		return key, err
	}

	// No kid provided - try all candidate keys for the algorithm
	return jcfg.tryMultipleKeysForValidation(token, algorithm)
}

// tryUpdateJWKSWithRetry attempts to update JWKS with retry logic for concurrent updates
func (jcfg *JwksConfig) tryUpdateJWKSWithRetry() error {
	const maxRetries = 5
	const retryDelay = 1 * time.Second

	for attempt := 1; attempt <= maxRetries; attempt++ {
		if attempt == 1 {
			updateErr := jcfg.UpdateJWKS()
			if updateErr == nil {
				return nil
			}
			if !errors.Is(updateErr, ErrJWKSUpdateInProgress) {
				return updateErr
			}
		}

		if attempt < maxRetries {
			time.Sleep(retryDelay)
		}

		if jcfg.GetJWKS() != nil {
			return nil
		}
	}

	return ErrJWKSUpdateInProgress
}

// tryMultipleKeysForValidation tries all candidate keys for algorithm-only validation
func (jcfg *JwksConfig) tryMultipleKeysForValidation(token *jwt.Token, algorithm string) (interface{}, error) {
	// Get all candidate keys from current JWKS
	candidateKeys, err := jcfg.getCandidateKeysWithRetry(algorithm)
	if err != nil {
		return nil, errors.Wrap(jwt.ErrInvalidKey, err.Error())
	}

	// Try to validate token with current candidate keys
	key, err := jcfg.tryValidateWithCandidateKeys(token, candidateKeys)
	if err == nil {
		return key, nil
	}

	// If all current keys failed, try with fresh JWKS update
	return jcfg.tryValidateWithFreshJWKS(token, algorithm, err)
}

// getCandidateKeysWithRetry gets candidate keys, with JWKS update retry if initial attempt fails
func (jcfg *JwksConfig) getCandidateKeysWithRetry(algorithm string) ([]interface{}, error) {
	candidateKeys, err := jcfg.getAllCandidateKeys(algorithm)
	if err != nil {
		// Attempt JWKS update and retry
		if updateErr := jcfg.tryUpdateJWKSWithRetry(); updateErr == nil {
			candidateKeys, err = jcfg.getAllCandidateKeys(algorithm)
		}
	}
	return candidateKeys, err
}

// tryValidateWithCandidateKeys attempts to validate token with provided candidate keys
func (jcfg *JwksConfig) tryValidateWithCandidateKeys(token *jwt.Token, candidateKeys []interface{}) (interface{}, error) {
	// Use the same comprehensive algorithm list as ValidateToken
	allCommonAlgorithms := []string{
		"RS256", "RS384", "RS512", // RSA with SHA
		"PS256", "PS384", "PS512", // RSA-PSS with SHA
		"ES256", "ES384", "ES512", // ECDSA with SHA
		"HS256", "HS384", "HS512", // HMAC with SHA
		"EdDSA", // Ed25519/Ed448
	}

	jwtParser := jwt.NewParser(jwt.WithValidMethods(allCommonAlgorithms))

	var lastErr error
	for _, candidateKey := range candidateKeys {
		keyFunc := func(token *jwt.Token) (interface{}, error) {
			return candidateKey, nil
		}

		_, parseErr := jwtParser.Parse(token.Raw, keyFunc)
		if parseErr == nil {
			return candidateKey, nil
		}
		lastErr = parseErr
	}

	return nil, lastErr
}

// tryValidateWithFreshJWKS attempts validation after updating JWKS with fresh keys
func (jcfg *JwksConfig) tryValidateWithFreshJWKS(token *jwt.Token, algorithm string, previousErr error) (interface{}, error) {
	if updateErr := jcfg.tryUpdateJWKSWithRetry(); updateErr == nil {
		freshCandidateKeys, freshErr := jcfg.getAllCandidateKeys(algorithm)
		if freshErr == nil && len(freshCandidateKeys) > 0 {
			key, err := jcfg.tryValidateWithCandidateKeys(token, freshCandidateKeys)
			if err == nil {
				return key, nil
			}
			previousErr = err // Update error from fresh validation attempt
		}
	}

	return nil, errors.Wrap(jwt.ErrInvalidKey, previousErr.Error())
}

// getAllCandidateKeys retrieves all candidate keys for an algorithm (used when no kid provided)
func (jcfg *JwksConfig) getAllCandidateKeys(algorithm string) ([]interface{}, error) {
	jwks := jcfg.GetJWKS()
	if jwks == nil {
		return nil, ErrJWKSNotInitialized
	}

	if algorithm == "" {
		return nil, jwt.ErrTokenMalformed
	}

	supportedKeys := jwks.GetKeysForAlgorithm(algorithm)
	if len(supportedKeys) == 0 {
		return nil, errors.Wrapf(jose.ErrUnsupportedAlgorithm, "algorithm %s", algorithm)
	}

	// Collect all valid keys, preferring signing keys first
	var signingKeys []interface{}
	var otherKeys []interface{}

	for _, key := range supportedKeys {
		if key.Valid() {
			if key.Use == "" || key.Use == "sig" {
				signingKeys = append(signingKeys, key.Key)
			} else {
				otherKeys = append(otherKeys, key.Key)
			}
		}
	}

	// Return signing keys first, then other keys
	result := append(signingKeys, otherKeys...)
	if len(result) == 0 {
		return nil, errors.Wrapf(jose.ErrUnsupportedAlgorithm, "algorithm %s", algorithm)
	}

	return result, nil
}

// getKeyFromJWKS retrieves a key from JWKS by kid
func (jcfg *JwksConfig) getKeyFromJWKS(kid string) (interface{}, error) {
	jwks := jcfg.GetJWKS()
	if jwks == nil {
		return nil, ErrJWKSNotInitialized
	}

	if kid == "" {
		return nil, errors.Wrapf(jwt.ErrInvalidKeyType, "kid is empty")
	}

	return jcfg.GetKeyByID(kid)
}

// HasClaimMappings returns true if claim mappings are configured.
func (jcfg *JwksConfig) HasClaimMappings() bool { return len(jcfg.ClaimMappings) > 0 }

// GetClaimMappings returns the claim mappings.
func (jcfg *JwksConfig) GetClaimMappings() []ClaimMapping { return jcfg.ClaimMappings }

// GetSubjectPrefix returns the issuer-derived prefix for namespacing subjects.
func (jcfg *JwksConfig) GetSubjectPrefix() string {
	if jcfg.subjectPrefix == "" && jcfg.Issuer != "" {
		jcfg.subjectPrefix = computeIssuerPrefix(jcfg.Issuer)
	}
	return jcfg.subjectPrefix
}

// ValidateAudience checks token has at least one configured audience. Returns nil if none configured.
func (jcfg *JwksConfig) ValidateAudience(claims jwt.MapClaims) error {
	if len(jcfg.Audiences) == 0 {
		return nil
	}
	tokenAudiences, err := claims.GetAudience()
	if err != nil {
		return ErrInvalidAudience
	}
	tokenAudSet := mapset.NewSet([]string(tokenAudiences)...)
	requiredAudSet := mapset.NewSet(jcfg.Audiences...)
	if tokenAudSet.Intersect(requiredAudSet).Cardinality() == 0 {
		return ErrInvalidAudience
	}
	return nil
}

// ValidateScopes checks token has ALL configured scopes. Returns nil if none configured.
func (jcfg *JwksConfig) ValidateScopes(claims jwt.MapClaims) error {
	if len(jcfg.Scopes) == 0 {
		return nil
	}
	tokenScopeSet := extractTokenScopes(claims)
	requiredScopeSet := mapset.NewSet(jcfg.Scopes...)
	if !tokenScopeSet.IsSuperset(requiredScopeSet) {
		return ErrInvalidScope
	}
	return nil
}

// GetOrgDataFromClaim extracts org data for the requested org and all accessible orgs.
// This method validates org access and returns errors if:
//   - ErrReservedOrgName: dynamic org claims a statically-configured org name
//   - ErrInvalidConfiguration: no claim mapping configured for the requested org
//   - ErrNoClaimRoles: no roles found for the requested org
//
// Returns orgData, isServiceAccount, and any error.
func (jcfg *JwksConfig) GetOrgDataFromClaim(claims jwt.MapClaims, reqOrgFromRoute string) (cdbm.OrgData, bool, error) {
	reqOrg := strings.ToLower(reqOrgFromRoute)
	orgData := make(cdbm.OrgData)
	foundReqOrgMapping := false
	isServiceAccount := false

	for _, cm := range jcfg.ClaimMappings {
		var orgName, displayName string
		var roles []string
		var err error

		switch {
		case cm.IsOrgDynamic():
			orgName, displayName = cm.ExtractOrgFromClaims(claims)
			if orgName == "" {
				continue
			}
			orgNameLower := strings.ToLower(orgName)
			if jcfg.ReservedOrgNames != nil && jcfg.ReservedOrgNames[orgNameLower] {
				if orgNameLower == reqOrg {
					return nil, false, ErrReservedOrgName
				}
				continue
			}
		case cm.IsOrgStatic():
			orgName = cm.OrgName
			displayName = cm.OrgDisplayName
		}

		orgNameLower := strings.ToLower(orgName)
		isReqOrg := orgNameLower == reqOrg

		roles, err = cm.ExtractRolesFromClaims(claims)
		if err != nil || len(roles) == 0 {
			if isReqOrg {
				return nil, false, ErrNoClaimRoles
			}
			continue
		}

		org := cdbm.Org{
			Name:        orgNameLower,
			DisplayName: displayName,
			OrgType:     "ENTERPRISE",
			Roles:       roles,
			Teams:       []cdbm.Team{},
		}
		// Set Updated timestamp for the requested org
		if isReqOrg {
			foundReqOrgMapping = true
			isServiceAccount = cm.IsServiceAccount

			now := time.Now().UTC()
			org.Updated = &now
		}
		orgData[orgNameLower] = org
	}

	if !foundReqOrgMapping {
		return nil, false, ErrInvalidConfiguration
	}

	return orgData, isServiceAccount, nil
}

// extractTokenScopes extracts scopes from claims (tries "scope", "scopes", "scp").
func extractTokenScopes(claims jwt.MapClaims) mapset.Set[string] {
	scopeSet := mapset.NewSet[string]()
	var scopeClaimValue interface{}
	for _, key := range scopeClaims {
		if val, exists := claims[key]; exists {
			scopeClaimValue = val
			break
		}
	}
	if scopeClaimValue == nil {
		return scopeSet
	}
	scopes, _ := InterfaceToStringSlice(scopeClaimValue)
	for _, scope := range scopes {
		scopeSet.Add(scope)
	}
	return scopeSet
}

// NewJwksConfig is a function that initializes and returns a configuration object for managing JWKS
func NewJwksConfig(name string, url string, issuer string, origin string, serviceAccount bool, audiences []string, scopes []string) *JwksConfig {
	// Default to custom origin if not specified
	if origin == "" {
		origin = TokenOriginCustom
	}
	return &JwksConfig{
		Name:           name,
		URL:            url,
		Issuer:         issuer,
		Origin:         origin,
		ServiceAccount: serviceAccount,
		Audiences:      audiences,
		Scopes:         scopes,
	}
}
