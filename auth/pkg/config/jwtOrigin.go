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
	"sync"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
	"github.com/nvidia/carbide-rest/common/pkg/util"
	cdbm "github.com/nvidia/carbide-rest/db/pkg/db/model"
)

const (
	TokenOriginUnknown = iota
	TokenOriginKas
	TokenOriginSsa
	TokenOriginKeycloak
	TokenOriginCustom
	TokenOriginMax
)

// TokenProcessor interface for processing JWT tokens
type TokenProcessor interface {
	ProcessToken(c echo.Context, tokenStr string, jwksConfig *JwksConfig, logger zerolog.Logger) (*cdbm.User, *util.APIError)
}

// JWTOriginConfig holds configuration for JWT origins with multiple JWKS configs and handlers
type JWTOriginConfig struct {
	sync.RWMutex                        // protects concurrent access to configs and handlers maps
	configs      map[string]*JwksConfig // map issuer -> JWKSConfig
	processors   map[int]TokenProcessor // map TokenOrigin -> TokenProcessor
}

// NewJWTOriginConfig initializes and returns a configuration object with empty maps
func NewJWTOriginConfig() *JWTOriginConfig {
	return &JWTOriginConfig{
		configs:    make(map[string]*JwksConfig),
		processors: make(map[int]TokenProcessor),
	}
}

// AddConfig adds a new JWKS config with the specified name, issuer, URL, origin, and serviceAccount flag
func (jc *JWTOriginConfig) AddConfig(name, issuer, url string, origin int, serviceAccount bool, audiences []string, scopes []string) {
	jc.Lock()
	defer jc.Unlock()
	jc.configs[issuer] = NewJwksConfig(name, url, issuer, origin, serviceAccount, audiences, scopes)
}

// AddConfigWithProcessor adds a new JWKS config and processor for the specified origin
func (jc *JWTOriginConfig) AddConfigWithProcessor(name, issuer, url string, origin int, serviceAccount bool, audiences []string, scopes []string, processor TokenProcessor) {
	jc.Lock()
	defer jc.Unlock()
	jc.configs[issuer] = NewJwksConfig(name, url, issuer, origin, serviceAccount, audiences, scopes)
	jc.processors[origin] = processor
}

// SetProcessorForOrigin sets a processor for the specified token origin
func (jc *JWTOriginConfig) SetProcessorForOrigin(origin int, processor TokenProcessor) {
	jc.Lock()
	defer jc.Unlock()
	jc.processors[origin] = processor
}

// GetProcessorByOrigin returns the processor for the specified origin
func (jc *JWTOriginConfig) GetProcessorByOrigin(origin int) TokenProcessor {
	jc.RLock()
	defer jc.RUnlock()
	return jc.processors[origin]
}

// GetProcessorByIssuer finds a processor that exactly matches the given issuer
func (jc *JWTOriginConfig) GetProcessorByIssuer(issuer string) TokenProcessor {
	jc.RLock()
	defer jc.RUnlock()
	config := jc.configs[issuer]
	if config != nil {
		return jc.processors[config.Origin]
	}
	return nil
}

// GetConfig returns the JWKS configuration for the specified issuer
func (jc *JWTOriginConfig) GetConfig(issuer string) *JwksConfig {
	jc.RLock()
	defer jc.RUnlock()
	return jc.configs[issuer]
}

// GetConfigsByOrigin returns all JWKS configurations for the specified origin
func (jc *JWTOriginConfig) GetConfigsByOrigin(origin int) map[string]*JwksConfig {
	jc.RLock()
	defer jc.RUnlock()
	result := make(map[string]*JwksConfig)
	for issuer, config := range jc.configs {
		if config.Origin == origin {
			result[issuer] = config
		}
	}
	return result
}

// GetFirstConfigByOrigin returns the first JWKS configuration with the specified origin
func (jc *JWTOriginConfig) GetFirstConfigByOrigin(origin int) *JwksConfig {
	jc.RLock()
	defer jc.RUnlock()
	for _, config := range jc.configs {
		if config.Origin == origin {
			return config
		}
	}
	return nil
}

// RemoveConfig removes the JWKS configuration for the specified issuer
func (jc *JWTOriginConfig) RemoveConfig(issuer string) {
	jc.Lock()
	defer jc.Unlock()
	delete(jc.configs, issuer)
}

// GetAllConfigs returns all JWKS configurations
func (jc *JWTOriginConfig) GetAllConfigs() map[string]*JwksConfig {
	jc.RLock()
	defer jc.RUnlock()
	return jc.configs
}

// UpdateJWKs updates the JWKs for all configurations in the map
func (jc *JWTOriginConfig) UpdateJWKs() error {
	jc.RLock()
	defer jc.RUnlock()
	for _, config := range jc.configs {
		if config != nil && config.URL != "" {
			if err := config.UpdateJWKs(); err != nil {
				return err
			}
		}
	}
	return nil
}

// GetSsaConfig returns the first SSA configuration
// Deprecated: Use GetFirstConfigByOrigin(TokenOriginSsa) instead
func (jc *JWTOriginConfig) GetSsaConfig() *JwksConfig {
	return jc.GetFirstConfigByOrigin(TokenOriginSsa)
}

// GetKasConfig returns the first KAS configuration
// Deprecated: Use GetFirstConfigByOrigin(TokenOriginKas) instead
func (jc *JWTOriginConfig) GetKasConfig() *JwksConfig {
	return jc.GetFirstConfigByOrigin(TokenOriginKas)
}

// GetKeycloakConfig returns the first Keycloak configuration
// Deprecated: Use GetFirstConfigByOrigin(TokenOriginKeycloak) instead
func (jc *JWTOriginConfig) GetKeycloakConfig() *JwksConfig {
	return jc.GetFirstConfigByOrigin(TokenOriginKeycloak)
}

// GetKeycloakProcessor returns the processor for Keycloak tokens
func (jc *JWTOriginConfig) GetKeycloakProcessor() TokenProcessor {
	jc.RLock()
	defer jc.RUnlock()
	return jc.processors[TokenOriginKeycloak]
}

// GetSsaProcessor returns the processor for SSA tokens
func (jc *JWTOriginConfig) GetSsaProcessor() TokenProcessor {
	jc.RLock()
	defer jc.RUnlock()
	return jc.processors[TokenOriginSsa]
}

// GetKasProcessor returns the processor for KAS tokens
func (jc *JWTOriginConfig) GetKasProcessor() TokenProcessor {
	jc.RLock()
	defer jc.RUnlock()
	return jc.processors[TokenOriginKas]
}

// SetProcessors sets all processors at once for easier initialization
func (jc *JWTOriginConfig) SetProcessors(keycloakProcessor, ssaProcessor, kasProcessor TokenProcessor) {
	jc.Lock()
	defer jc.Unlock()
	jc.processors[TokenOriginKeycloak] = keycloakProcessor
	jc.processors[TokenOriginSsa] = ssaProcessor
	jc.processors[TokenOriginKas] = kasProcessor
}

// IsServiceAccount checks if the given issuer supports service account tokens
func (jc *JWTOriginConfig) IsServiceAccount(issuer string) bool {
	jc.RLock()
	defer jc.RUnlock()
	config := jc.configs[issuer]
	if config != nil {
		return config.ServiceAccount
	}
	return false
}
