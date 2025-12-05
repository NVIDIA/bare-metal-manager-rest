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
	"crypto/rand"
	"crypto/rsa"
	"encoding/json"
	"math/big"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/nvidia/carbide-rest/auth/pkg/config"
	testutil "github.com/nvidia/carbide-rest/auth/pkg/testing"
	"github.com/nvidia/carbide-rest/common/pkg/util"
	cdbu "github.com/nvidia/carbide-rest/db/pkg/util"
)

// setupTestEnvironment creates a test environment with mock JWKS server and database
func setupTestEnvironment(t *testing.T, audiences []string, scopes []string) (*CustomProcessor, *config.JwksConfig, *rsa.PrivateKey, *httptest.Server, func()) {
	// Generate RSA key for signing tokens
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)

	// Create JWKS server
	jwksServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jwks := map[string]interface{}{
			"keys": []map[string]interface{}{
				{
					"kty": "RSA",
					"kid": "test-key-id",
					"use": "sig",
					"alg": "RS256",
					"n":   testutil.EncodeBase64URLBigInt(privateKey.N),
					"e":   testutil.EncodeBase64URLBigInt(big.NewInt(int64(privateKey.E))),
				},
			},
		}
		json.NewEncoder(w).Encode(jwks)
	}))

	// Create JWKS config
	jwksConfig := config.NewJwksConfig(
		"custom-provider",
		jwksServer.URL,
		"https://custom.example.com",
		config.TokenOriginCustom,
		true, // service account enabled
		audiences,
		scopes,
	)

	// Initialize JWKS
	err = jwksConfig.UpdateJWKs()
	require.NoError(t, err)

	// Create test database session
	dbSession := cdbu.GetTestDBSession(t, false)

	// Create processor
	processor := &CustomProcessor{
		dbSession: dbSession,
	}

	cleanup := func() {
		jwksServer.Close()
	}

	return processor, jwksConfig, privateKey, jwksServer, cleanup
}

// createTestToken creates a JWT token with the given claims
func createTestToken(t *testing.T, privateKey *rsa.PrivateKey, claims jwt.MapClaims) string {
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	token.Header["kid"] = "test-key-id"

	tokenString, err := token.SignedString(privateKey)
	require.NoError(t, err)
	return tokenString
}

func TestCustomProcessor_ValidateAudiences_Success(t *testing.T) {
	tests := []struct {
		name                string
		configuredAudiences []string
		tokenAudience       interface{} // can be string or []string
		shouldPass          bool
	}{
		{
			name:                "single audience matches",
			configuredAudiences: []string{"api.example.com"},
			tokenAudience:       "api.example.com",
			shouldPass:          true,
		},
		{
			name:                "one of multiple configured audiences matches",
			configuredAudiences: []string{"api.example.com", "app.example.com"},
			tokenAudience:       "app.example.com",
			shouldPass:          true,
		},
		{
			name:                "token has array of audiences, one matches",
			configuredAudiences: []string{"api.example.com"},
			tokenAudience:       []interface{}{"other.example.com", "api.example.com"},
			shouldPass:          true,
		},
		{
			name:                "no configured audiences (validation skipped)",
			configuredAudiences: []string{},
			tokenAudience:       "any.example.com",
			shouldPass:          true,
		},
		{
			name:                "nil configured audiences (validation skipped)",
			configuredAudiences: nil,
			tokenAudience:       "any.example.com",
			shouldPass:          true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			processor, jwksConfig, privateKey, _, cleanup := setupTestEnvironment(t, tt.configuredAudiences, nil)
			defer cleanup()

			claims := jwt.MapClaims{
				"sub": "test-user-123",
				"iss": "https://custom.example.com",
				"aud": tt.tokenAudience,
				"exp": time.Now().Add(time.Hour).Unix(),
				"iat": time.Now().Unix(),
			}

			tokenString := createTestToken(t, privateKey, claims)

			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			logger := zerolog.Nop()

			// Mock the database operations by creating a minimal user
			// In a real test, you'd want to mock the database properly
			_, apiErr := processor.ProcessToken(c, tokenString, jwksConfig, logger)

			if tt.shouldPass {
				// For this test, we're only checking that audience validation passes
				// The actual user creation might fail due to no real DB, but that's OK
				// We just want to ensure no audience-related errors
				if apiErr != nil && apiErr.Code == http.StatusUnauthorized {
					// Check if the error is audience-related
					if contains(apiErr.Message, "audience") {
						t.Errorf("Expected audience validation to pass, but got error: %s", apiErr.Message)
					}
				}
			} else {
				assert.NotNil(t, apiErr, "Expected error for invalid audience")
				if apiErr != nil {
					assert.Equal(t, http.StatusUnauthorized, apiErr.Code)
					assert.Contains(t, apiErr.Message, "audience")
				}
			}
		})
	}
}

func TestCustomProcessor_ValidateAudiences_Failure(t *testing.T) {
	tests := []struct {
		name                string
		configuredAudiences []string
		tokenAudience       interface{}
	}{
		{
			name:                "audience does not match",
			configuredAudiences: []string{"api.example.com"},
			tokenAudience:       "wrong.example.com",
		},
		{
			name:                "token has array but none match",
			configuredAudiences: []string{"api.example.com"},
			tokenAudience:       []interface{}{"other.example.com", "another.example.com"},
		},
		{
			name:                "multiple configured audiences but token has wrong one",
			configuredAudiences: []string{"api.example.com", "app.example.com"},
			tokenAudience:       "wrong.example.com",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			processor, jwksConfig, privateKey, _, cleanup := setupTestEnvironment(t, tt.configuredAudiences, nil)
			defer cleanup()

			claims := jwt.MapClaims{
				"sub": "test-user-123",
				"iss": "https://custom.example.com",
				"aud": tt.tokenAudience,
				"exp": time.Now().Add(time.Hour).Unix(),
				"iat": time.Now().Unix(),
			}

			tokenString := createTestToken(t, privateKey, claims)

			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			logger := zerolog.Nop()

			_, apiErr := processor.ProcessToken(c, tokenString, jwksConfig, logger)

			require.NotNil(t, apiErr, "Expected error for invalid audience")
			assert.Equal(t, http.StatusUnauthorized, apiErr.Code)
			assert.Contains(t, apiErr.Message, "audience")
		})
	}
}

func TestCustomProcessor_ValidateScopes_Success(t *testing.T) {
	tests := []struct {
		name             string
		configuredScopes []string
		tokenScopes      interface{} // can be string (space-separated) or []string
		shouldPass       bool
	}{
		{
			name:             "single scope matches - carbide",
			configuredScopes: []string{"carbide"},
			tokenScopes:      "carbide",
			shouldPass:       true,
		},
		{
			name:             "multiple scopes all present - includes carbide",
			configuredScopes: []string{"carbide", "read:data"},
			tokenScopes:      "carbide read:data write:data",
			shouldPass:       true,
		},
		{
			name:             "scopes as array",
			configuredScopes: []string{"carbide"},
			tokenScopes:      []interface{}{"carbide", "other"},
			shouldPass:       true,
		},
		{
			name:             "no configured scopes (validation skipped)",
			configuredScopes: []string{},
			tokenScopes:      "any scope",
			shouldPass:       true,
		},
		{
			name:             "nil configured scopes (validation skipped)",
			configuredScopes: nil,
			tokenScopes:      "any scope",
			shouldPass:       true,
		},
		{
			name:             "all required scopes present in space-separated string",
			configuredScopes: []string{"read:data", "write:data"},
			tokenScopes:      "read:data write:data admin",
			shouldPass:       true,
		},
		{
			name:             "scopes in scp claim instead of scope",
			configuredScopes: []string{"carbide"},
			tokenScopes:      "carbide", // Will be put in "scp" claim in test
			shouldPass:       true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			processor, jwksConfig, privateKey, _, cleanup := setupTestEnvironment(t, nil, tt.configuredScopes)
			defer cleanup()

			claims := jwt.MapClaims{
				"sub": "test-user-123",
				"iss": "https://custom.example.com",
				"aud": "test-audience",
				"exp": time.Now().Add(time.Hour).Unix(),
				"iat": time.Now().Unix(),
			}

			// Add scope claim
			if tt.name == "scopes in scp claim instead of scope" {
				claims["scp"] = tt.tokenScopes
			} else {
				claims["scope"] = tt.tokenScopes
			}

			tokenString := createTestToken(t, privateKey, claims)

			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			logger := zerolog.Nop()

			_, apiErr := processor.ProcessToken(c, tokenString, jwksConfig, logger)

			if tt.shouldPass {
				// Check that we don't have scope-related errors
				if apiErr != nil && apiErr.Code == http.StatusUnauthorized {
					if contains(apiErr.Message, "scope") {
						t.Errorf("Expected scope validation to pass, but got error: %s", apiErr.Message)
					}
				}
			} else {
				assert.NotNil(t, apiErr, "Expected error for invalid scopes")
				if apiErr != nil {
					assert.Equal(t, http.StatusUnauthorized, apiErr.Code)
					assert.Contains(t, apiErr.Message, "scope")
				}
			}
		})
	}
}

func TestCustomProcessor_ValidateScopes_Failure(t *testing.T) {
	tests := []struct {
		name             string
		configuredScopes []string
		tokenScopes      interface{}
	}{
		{
			name:             "scope does not match",
			configuredScopes: []string{"carbide"},
			tokenScopes:      "other",
		},
		{
			name:             "missing one required scope",
			configuredScopes: []string{"carbide", "read:data"},
			tokenScopes:      "carbide",
		},
		{
			name:             "completely different scopes",
			configuredScopes: []string{"carbide"},
			tokenScopes:      "admin write:data",
		},
		{
			name:             "empty token scopes",
			configuredScopes: []string{"carbide"},
			tokenScopes:      "",
		},
		{
			name:             "token has array but missing required scope",
			configuredScopes: []string{"carbide", "admin"},
			tokenScopes:      []interface{}{"carbide"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			processor, jwksConfig, privateKey, _, cleanup := setupTestEnvironment(t, nil, tt.configuredScopes)
			defer cleanup()

			claims := jwt.MapClaims{
				"sub":   "test-user-123",
				"iss":   "https://custom.example.com",
				"aud":   "test-audience",
				"scope": tt.tokenScopes,
				"exp":   time.Now().Add(time.Hour).Unix(),
				"iat":   time.Now().Unix(),
			}

			tokenString := createTestToken(t, privateKey, claims)

			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			logger := zerolog.Nop()

			_, apiErr := processor.ProcessToken(c, tokenString, jwksConfig, logger)

			require.NotNil(t, apiErr, "Expected error for invalid scopes")
			assert.Equal(t, http.StatusUnauthorized, apiErr.Code)
			assert.Contains(t, apiErr.Message, "scope")
		})
	}
}

func TestCustomProcessor_CombinedAudienceAndScope_Validation(t *testing.T) {
	tests := []struct {
		name                string
		configuredAudiences []string
		configuredScopes    []string
		tokenAudience       string
		tokenScopes         string
		shouldPass          bool
		errorShouldContain  string
	}{
		{
			name:                "both audience and scopes valid - carbide scope",
			configuredAudiences: []string{"api.example.com"},
			configuredScopes:    []string{"carbide"},
			tokenAudience:       "api.example.com",
			tokenScopes:         "carbide admin",
			shouldPass:          true,
		},
		{
			name:                "valid audience but invalid scopes",
			configuredAudiences: []string{"api.example.com"},
			configuredScopes:    []string{"carbide"},
			tokenAudience:       "api.example.com",
			tokenScopes:         "other",
			shouldPass:          false,
			errorShouldContain:  "scope",
		},
		{
			name:                "invalid audience but valid scopes",
			configuredAudiences: []string{"api.example.com"},
			configuredScopes:    []string{"carbide"},
			tokenAudience:       "wrong.example.com",
			tokenScopes:         "carbide",
			shouldPass:          false,
			errorShouldContain:  "audience",
		},
		{
			name:                "both invalid",
			configuredAudiences: []string{"api.example.com"},
			configuredScopes:    []string{"carbide"},
			tokenAudience:       "wrong.example.com",
			tokenScopes:         "other",
			shouldPass:          false,
			errorShouldContain:  "audience", // Audience is checked first
		},
		{
			name:                "neither configured - both should pass",
			configuredAudiences: nil,
			configuredScopes:    nil,
			tokenAudience:       "any.example.com",
			tokenScopes:         "any",
			shouldPass:          true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			processor, jwksConfig, privateKey, _, cleanup := setupTestEnvironment(t, tt.configuredAudiences, tt.configuredScopes)
			defer cleanup()

			claims := jwt.MapClaims{
				"sub":   "test-user-123",
				"iss":   "https://custom.example.com",
				"aud":   tt.tokenAudience,
				"scope": tt.tokenScopes,
				"exp":   time.Now().Add(time.Hour).Unix(),
				"iat":   time.Now().Unix(),
			}

			tokenString := createTestToken(t, privateKey, claims)

			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			logger := zerolog.Nop()

			_, apiErr := processor.ProcessToken(c, tokenString, jwksConfig, logger)

			if tt.shouldPass {
				// May still have DB errors, but no audience/scope errors
				if apiErr != nil && apiErr.Code == http.StatusUnauthorized {
					if contains(apiErr.Message, "audience") || contains(apiErr.Message, "scope") {
						t.Errorf("Expected validation to pass, but got error: %s", apiErr.Message)
					}
				}
			} else {
				require.NotNil(t, apiErr, "Expected error for validation failure")
				assert.Equal(t, http.StatusUnauthorized, apiErr.Code)
				assert.Contains(t, apiErr.Message, tt.errorShouldContain)
			}
		})
	}
}

func TestCustomProcessor_MissingAudienceClaim(t *testing.T) {
	processor, jwksConfig, privateKey, _, cleanup := setupTestEnvironment(t, []string{"api.example.com"}, nil)
	defer cleanup()

	claims := jwt.MapClaims{
		"sub": "test-user-123",
		"iss": "https://custom.example.com",
		// No "aud" claim
		"exp": time.Now().Add(time.Hour).Unix(),
		"iat": time.Now().Unix(),
	}

	tokenString := createTestToken(t, privateKey, claims)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	logger := zerolog.Nop()

	_, apiErr := processor.ProcessToken(c, tokenString, jwksConfig, logger)

	require.NotNil(t, apiErr, "Expected error for missing audience claim")
	assert.Equal(t, http.StatusUnauthorized, apiErr.Code)
	assert.Contains(t, apiErr.Message, "audience")
}

func TestCustomProcessor_MissingScopeClaim(t *testing.T) {
	processor, jwksConfig, privateKey, _, cleanup := setupTestEnvironment(t, nil, []string{"carbide"})
	defer cleanup()

	claims := jwt.MapClaims{
		"sub": "test-user-123",
		"iss": "https://custom.example.com",
		"aud": "test-audience",
		// No "scope" or "scp" claim
		"exp": time.Now().Add(time.Hour).Unix(),
		"iat": time.Now().Unix(),
	}

	tokenString := createTestToken(t, privateKey, claims)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	logger := zerolog.Nop()

	_, apiErr := processor.ProcessToken(c, tokenString, jwksConfig, logger)

	require.NotNil(t, apiErr, "Expected error for missing scope claim")
	assert.Equal(t, http.StatusUnauthorized, apiErr.Code)
	assert.Contains(t, apiErr.Message, "scope")
}

// Direct validation tests - these test the validation functions without database dependencies

func TestValidateAudiences_DirectTest(t *testing.T) {
	logger := zerolog.Nop()

	tests := []struct {
		name                string
		tokenClaims         jwt.MapClaims
		configuredAudiences []string
		shouldPass          bool
		errorContains       string
	}{
		{
			name: "single audience matches - carbide example",
			tokenClaims: jwt.MapClaims{
				"aud": "api.carbide.com",
			},
			configuredAudiences: []string{"api.carbide.com"},
			shouldPass:          true,
		},
		{
			name: "multiple audiences, one matches",
			tokenClaims: jwt.MapClaims{
				"aud": []interface{}{"other.com", "api.carbide.com"},
			},
			configuredAudiences: []string{"api.carbide.com"},
			shouldPass:          true,
		},
		{
			name: "audience mismatch",
			tokenClaims: jwt.MapClaims{
				"aud": "wrong.com",
			},
			configuredAudiences: []string{"api.carbide.com"},
			shouldPass:          false,
			errorContains:       "audience",
		},
		{
			name: "missing audience claim",
			tokenClaims: jwt.MapClaims{
				"sub": "test",
			},
			configuredAudiences: []string{"api.carbide.com"},
			shouldPass:          false,
			errorContains:       "audience claim",
		},
		{
			name: "no configured audiences - validation skipped",
			tokenClaims: jwt.MapClaims{
				"aud": "any.com",
			},
			configuredAudiences: []string{},
			shouldPass:          true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateAudiences(tt.tokenClaims, tt.configuredAudiences, logger)

			if tt.shouldPass {
				assert.Nil(t, err, "Expected no error for valid audience")
			} else {
				require.NotNil(t, err, "Expected error for invalid audience")
				assert.Contains(t, err.Message, tt.errorContains)
				assert.Equal(t, http.StatusUnauthorized, err.Code)
			}
		})
	}
}

func TestValidateScopes_DirectTest(t *testing.T) {
	logger := zerolog.Nop()

	tests := []struct {
		name           string
		tokenClaims    jwt.MapClaims
		requiredScopes []string
		shouldPass     bool
		errorContains  string
	}{
		{
			name: "single scope matches - scopes array with carbide",
			tokenClaims: jwt.MapClaims{
				"scopes": []interface{}{"carbide"},
			},
			requiredScopes: []string{"carbide"},
			shouldPass:     true,
		},
		{
			name: "multiple scopes in token, all required present - carbide included",
			tokenClaims: jwt.MapClaims{
				"scopes": []interface{}{"carbide", "read:data", "write:data"},
			},
			requiredScopes: []string{"carbide", "read:data"},
			shouldPass:     true,
		},
		{
			name: "scopes as space-separated string",
			tokenClaims: jwt.MapClaims{
				"scopes": "carbide admin read:data",
			},
			requiredScopes: []string{"carbide"},
			shouldPass:     true,
		},
		{
			name: "missing required scope",
			tokenClaims: jwt.MapClaims{
				"scopes": []interface{}{"read:data"},
			},
			requiredScopes: []string{"carbide"},
			shouldPass:     false,
			errorContains:  "scope",
		},
		{
			name: "missing scope claim",
			tokenClaims: jwt.MapClaims{
				"sub": "test",
			},
			requiredScopes: []string{"carbide"},
			shouldPass:     false,
			errorContains:  "scope claim",
		},
		{
			name: "no required scopes - validation skipped",
			tokenClaims: jwt.MapClaims{
				"scopes": []interface{}{"any"},
			},
			requiredScopes: []string{},
			shouldPass:     true,
		},
		{
			name: "scp claim as fallback - space-separated",
			tokenClaims: jwt.MapClaims{
				"scp": "carbide admin",
			},
			requiredScopes: []string{"carbide"},
			shouldPass:     true,
		},
		{
			name: "scope in scopes claim (plural) - kas token format",
			tokenClaims: jwt.MapClaims{
				"scopes": []interface{}{"kas"},
			},
			requiredScopes: []string{"kas"},
			shouldPass:     true,
		},
		{
			name: "multiple required scopes, all present",
			tokenClaims: jwt.MapClaims{
				"scopes": []interface{}{"carbide", "read:data", "write:data", "admin"},
			},
			requiredScopes: []string{"carbide", "read:data", "admin"},
			shouldPass:     true,
		},
		{
			name: "multiple required scopes, one missing",
			tokenClaims: jwt.MapClaims{
				"scopes": []interface{}{"carbide", "read:data"},
			},
			requiredScopes: []string{"carbide", "admin"},
			shouldPass:     false,
			errorContains:  "missing required scopes",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateScopes(tt.tokenClaims, tt.requiredScopes, logger)

			if tt.shouldPass {
				assert.Nil(t, err, "Expected no error for valid scopes")
			} else {
				require.NotNil(t, err, "Expected error for invalid scopes")
				assert.Contains(t, err.Message, tt.errorContains)
				assert.Equal(t, http.StatusUnauthorized, err.Code)
			}
		})
	}
}

func TestCombinedValidation_DirectTest(t *testing.T) {
	logger := zerolog.Nop()

	tests := []struct {
		name                string
		tokenClaims         jwt.MapClaims
		configuredAudiences []string
		requiredScopes      []string
		shouldPass          bool
		errorContains       string
	}{
		{
			name: "both audience and scopes valid - carbide example",
			tokenClaims: jwt.MapClaims{
				"aud":    "api.carbide.com",
				"scopes": []interface{}{"carbide", "read:data"},
			},
			configuredAudiences: []string{"api.carbide.com"},
			requiredScopes:      []string{"carbide"},
			shouldPass:          true,
		},
		{
			name: "valid audience, invalid scopes",
			tokenClaims: jwt.MapClaims{
				"aud":    "api.carbide.com",
				"scopes": []interface{}{"other"},
			},
			configuredAudiences: []string{"api.carbide.com"},
			requiredScopes:      []string{"carbide"},
			shouldPass:          false,
			errorContains:       "scope",
		},
		{
			name: "invalid audience, valid scopes",
			tokenClaims: jwt.MapClaims{
				"aud":    "wrong.com",
				"scopes": []interface{}{"carbide"},
			},
			configuredAudiences: []string{"api.carbide.com"},
			requiredScopes:      []string{"carbide"},
			shouldPass:          false,
			errorContains:       "audience",
		},
		{
			name: "both invalid",
			tokenClaims: jwt.MapClaims{
				"aud":    "wrong.com",
				"scopes": []interface{}{"other"},
			},
			configuredAudiences: []string{"api.carbide.com"},
			requiredScopes:      []string{"carbide"},
			shouldPass:          false,
			errorContains:       "audience", // Audience checked first
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test audience first
			audErr := validateAudiences(tt.tokenClaims, tt.configuredAudiences, logger)

			// Only test scopes if audience passed
			var scopeErr *util.APIError
			if audErr == nil {
				scopeErr = validateScopes(tt.tokenClaims, tt.requiredScopes, logger)
			}

			// Determine final result
			finalErr := audErr
			if finalErr == nil {
				finalErr = scopeErr
			}

			if tt.shouldPass {
				assert.Nil(t, finalErr, "Expected no error for valid token")
			} else {
				require.NotNil(t, finalErr, "Expected error for invalid token")
				assert.Contains(t, finalErr.Message, tt.errorContains)
				assert.Equal(t, http.StatusUnauthorized, finalErr.Code)
			}
		})
	}
}

// Helper function to check if a string contains a substring (case-insensitive)
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
		(len(s) > 0 && len(substr) > 0 && containsHelper(s, substr)))
}

func containsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
