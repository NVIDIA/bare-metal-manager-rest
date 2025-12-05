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
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/go-jose/go-jose/v4"
	"github.com/golang-jwt/jwt/v5"
	"github.com/pkg/errors"
	"github.com/nvidia/carbide-rest/auth/pkg/core"
)

var (
	// Custom errors specific to our JWKS management implementation
	// ErrJWKSURLEmpty is raised when JWKS URL is empty
	ErrJWKSURLEmpty = errors.New("JWKS URL is empty")
	// ErrJWKSNotInitialized is raised when JWKS has not been loaded
	ErrJWKSNotInitialized = errors.New("JWKS not initialized - call UpdateJWKs first")
	// ErrEmptyKeySet is raised when JWKS contains no keys
	ErrEmptyKeySet = errors.New("JWKS key set is empty")
	// ErrNoValidKeys is raised when JWKS contains no valid keys
	ErrNoValidKeys = errors.New("JWKS contains no valid keys")
	// ErrInvalidUseParameter is raised when an invalid use parameter is provided
	ErrInvalidUseParameter = errors.New("invalid use parameter")
	// ErrJWKSUpdateInProgress is raised when another thread is currently updating JWKS
	ErrJWKSUpdateInProgress = errors.New("JWKS update already in progress")
)

const (
	minUpdateInterval = 10 * time.Second
)

type JwksConfig struct {
	Name           string
	IsUpdating     uint32 // IsUpdating is used to track if the JWKS is being updated
	sync.RWMutex          // mutex is used to handle concurrent access to the JWKS
	URL            string
	Issuer         string
	Origin         int        // Origin indicates the token origin type
	ServiceAccount bool       // ServiceAccount indicates if this config supports service account tokens
	Audiences      []string   // Audiences is the list of valid audiences for token validation (optional)
	Scopes         []string   // Scopes is the list of required scopes for token validation (optional)
	LastUpdated    time.Time  // Track when we last updated JWKS
	jwks           *core.JWKS // Enhanced JWKS with go-jose capabilities
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

// UpdateJWKs populates/updates the JWKs with enhanced go-jose capabilities and validation
func (jcfg *JwksConfig) UpdateJWKs() error {
	if jcfg.URL == "" {
		return ErrJWKSURLEmpty
	}

	// Only allow first time update and 10 sec past the prev update
	if !jcfg.shouldAllowJWKSUpdate() {
		return nil
	}

	// Store IsUpdating so concurrent ops wouldn't issue fetch and update
	if !atomic.CompareAndSwapUint32(&jcfg.IsUpdating, 0, 1) {
		return ErrJWKSUpdateInProgress
	}
	// Ensure we reset the flag when function exits
	defer atomic.StoreUint32(&jcfg.IsUpdating, 0)

	jcfg.RLock()
	urlCopy := jcfg.URL
	jcfg.RUnlock()

	jwks, err := core.NewJWKSFromURL(urlCopy)
	if err != nil {
		return errors.Wrapf(err, "failed to update JWKS from %s", urlCopy)
	}

	// Validate that we have at least one key
	if jwks.Set == nil || len(jwks.Set.Keys) == 0 {
		return errors.Wrapf(ErrEmptyKeySet, "from %s", urlCopy)
	}

	// Validate that all keys in the set are valid according to go-jose
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
	lastUpdateTime := time.Now()
	jcfg.LastUpdated = lastUpdateTime

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
			updateErr := jcfg.UpdateJWKs()
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

// NewJwksConfig is a function that initializes and returns a configuration object for managing JWKS
func NewJwksConfig(name string, url string, issuer string, origin int, serviceAccount bool, audiences []string, scopes []string) *JwksConfig {
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
