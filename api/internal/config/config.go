/*
 * SPDX-FileCopyrightText: Copyright (c) 2021-2023 NVIDIA CORPORATION & AFFILIATES. All rights reserved.
 * SPDX-License-Identifier: LicenseRef-NvidiaProprietary
 *
 * NVIDIA CORPORATION, its affiliates and licensors retain all intellectual
 * property and proprietary rights in and to this material, related
 * documentation and any modifications thereto. Any use, reproduction,
 * disclosure or distribution of this material and related documentation
 * without an express license agreement from NVIDIA CORPORATION or
 * its affiliates is strictly prohibited.
 */

package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"

	cauth "github.com/nvidia/carbide-rest/auth/pkg/config"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

var (
	_, cur, _, _ = runtime.Caller(0)
	// ProjectRoot describes the folder path of this project
	ProjectRoot = filepath.Join(filepath.Dir(cur), "../..")
)

const (
	// ConfigFilePath specifies the path to the config file, this contains the default path
	ConfigFilePath = "config.file"
	// ConfigFilePathEnv specifies the environment variable name for the config file path
	ConfigFilePathEnv = "CONFIG_FILE_PATH"

	// ConfigEnvDev specifies if the service is running in development mode
	ConfigEnvDev = "env.dev"
	// ConfigEnvDisconnected specifies if the service is running in disconnected mode
	ConfigEnvDisconnected = "env.disconnected"

	// ConfigLogLevel specifies the log level
	ConfigLogLevel = "log.level"

	// ConfigSentryDSN is the DSN for Sentry
	ConfigSentryDSN = "log.sentry.dsn"

	// ConfigDBHost specifies the host of the database
	ConfigDBHost = "db.host"
	// ConfigDBPort specifies the port of the database
	ConfigDBPort = "db.port"
	// ConfigDBName specifies the name of the database
	ConfigDBName = "db.name"
	// ConfigDBUser specifies the user of the database
	ConfigDBUser = "db.user"
	// ConfigDBPasswordPath specifies the path for file containing password of the database
	ConfigDBPasswordPath = "db.passwordPath"
	// ConfigDBPassword specifies the password of the database, not read from config file, rather from file specified by db.passwordPath
	ConfigDBPassword = "db.password"

	// ConfigTemporalHost specifies the hostname for Temporal
	ConfigTemporalHost = "temporal.host"
	// ConfigTemporalPort specifies the port for Temporal
	ConfigTemporalPort = "temporal.port"
	// ConfigTemporalServerName specifies the server name for Temporal
	ConfigTemporalServerName = "temporal.serverName"
	// ConfigTemporalNamespace specifies the namespace for Temporal
	ConfigTemporalNamespace = "temporal.namespace"
	// ConfigTemporalQueue specifies the queue for Temporal
	ConfigTemporalQueue = "temporal.queue"
	// ConfigTemporalTlsEnabled specifies if TLS should be used with Temporal
	ConfigTemporalTlsEnabled = "temporal.tls.enabled"
	// ConfigTemporalCertPath specifies the path for client cert
	ConfigTemporalCertPath = "temporal.tls.certPath"
	// ConfigTemporalKeyPath specifies the path for client key
	ConfigTemporalKeyPath = "temporal.tls.keyPath"
	// ConfigTemporalCaPath specifies the path for server CA cert
	ConfigTemporalCaPath = "temporal.tls.caPath"
	// ConfigTemporalEncryptionKey specifies the data encryption key for Temporal
	ConfigTemporalEncryptionKey = "temporal.encryptionKey"
	// ConfigTemporalEncryptionKeyPath specifies the path for file containing data encryption key for Temporal
	ConfigTemporalEncryptionKeyPath = "temporal.encryptionKeyPath"

	// ConfigSiteManagerEnabled is a feature flag for site manager
	ConfigSiteManagerEnabled = "siteManager.enabled"
	// ConfigSiteManagerEndpoint is the service endpoint for site manager
	ConfigSiteManagerEndpoint = "siteManager.svcEndpoint"

	// ConfigSitePhoneHomeUrl is a phone home url from the site
	ConfigSitePhoneHomeUrl = "site.phoneHomeUrl"

	// ConfigMetricsEnabled is a feature flag for Prometheus metrics
	ConfigMetricsEnabled = "metrics.enabled"
	// ConfigMetricsPort specifies the port for Prometheus metrics
	ConfigMetricsPort = "metrics.port"

	// ConfigTracingEnabled is a feature flag for tracing
	ConfigTracingEnabled = "tracing.enabled"
	// ConfigTracingServiceName is the name of the tracing service
	ConfigTracingServiceName = "tracing.serviceName"

	// ConfigKeycloakEnabled is a feature flag for Keycloak authentication
	ConfigKeycloakEnabled = "keycloak.enabled"
	// ConfigKeycloakBaseURL is the base URL of the Keycloak server
	ConfigKeycloakBaseURL = "keycloak.baseURL"
	// ConfigKeycloakExternalBaseURL is the external base URL of the Keycloak server
	ConfigKeycloakExternalBaseURL = "keycloak.externalBaseURL"
	// ConfigKeycloakRealm is the Keycloak realm name
	ConfigKeycloakRealm = "keycloak.realm"
	// ConfigKeycloakClientID is the Keycloak client ID
	ConfigKeycloakClientID = "keycloak.clientID"
	// ConfigKeycloakClientSecretPath is the path to the file containing the client secret
	ConfigKeycloakClientSecretPath = "keycloak.clientSecretPath"
	// ConfigKeycloakClientSecret is the client secret, read from the secret file
	ConfigKeycloakClientSecret = "keycloak.clientSecret"
	// ConfigKeycloakServiceAccountEnabled is a feature flag for service account support
	ConfigKeycloakServiceAccountEnabled = "keycloak.serviceAccount"
)

// AuthConfig represents a single auth configuration entry
type AuthConfig struct {
	Name           string   `mapstructure:"name"`
	Origin         int      `mapstructure:"origin"`
	URL            string   `mapstructure:"url"`
	Issuer         string   `mapstructure:"issuer"`
	ServiceAccount bool     `mapstructure:"serviceAccount"`
	Audiences      []string `mapstructure:"audiences"`
	Scopes         []string `mapstructure:"scopes"`
}

// Maintain a global config object
var config *Config

// Config represents configurations for the service
type Config struct {
	sync.RWMutex
	v               *viper.Viper
	db              *DBConfig
	temporal        *TemporalConfig
	JwtOriginConfig *cauth.JWTOriginConfig
	SiteConfig      *SiteConfig
	KeycloakConfig  *cauth.KeycloakConfig
}

// NewConfig creates a new config object
func NewConfig() *Config {
	if config != nil {
		return config
	}

	c := Config{
		v: viper.New(),
	}

	// Set defaults
	c.v.SetDefault(ConfigLogLevel, "info")
	c.v.SetDefault(ConfigEnvDev, false)
	c.v.SetDefault(ConfigEnvDisconnected, false)

	// Set config file
	// Check environment variable. If not set, use default
	defaultPath := ProjectRoot + "/config.yaml"
	if os.Getenv(ConfigFilePathEnv) != "" {
		c.v.SetDefault(ConfigFilePath, os.Getenv(ConfigFilePathEnv))
	} else {
		log.Warn().Msg("config file path not set, using default")
		c.v.SetDefault(ConfigFilePath, defaultPath)
	}

	c.v.SetDefault(ConfigDBHost, "localhost")
	c.v.SetDefault(ConfigDBPort, 5432)
	c.v.SetDefault(ConfigTemporalHost, "localhost")
	c.v.SetDefault(ConfigTemporalPort, 7233)
	c.v.SetDefault(ConfigTemporalTlsEnabled, true)

	c.v.SetDefault(ConfigMetricsEnabled, true)
	c.v.SetDefault(ConfigMetricsPort, 9360)

	c.v.SetDefault(ConfigTracingEnabled, false)

	// SiteConfig default phone home url
	c.v.SetDefault(ConfigSitePhoneHomeUrl, "http://localhost")

	// Keycloak needs to be explicitly enabled via config
	c.v.SetDefault(ConfigKeycloakEnabled, false)
	// Service account support is disabled by default
	c.v.SetDefault(ConfigKeycloakServiceAccountEnabled, false)

	c.v.AutomaticEnv()
	c.v.SetConfigFile(c.GetPathToConfig())

	log.Warn().Str("config.file", c.GetPathToConfig()).Msg("loading config")

	err := c.v.ReadInConfig() // Find and read the config file
	if _, ok := err.(*os.PathError); ok {
		log.Warn().Msgf("no config file '%s' not found. Using default values", c.GetPathToConfig())
	} else if err != nil { // Handle other errors that occurred while reading the config file
		log.Panic().Err(err).Msgf("fatal error while reading the config file: %s", err)
	}

	// Set values
	c.setLogLevel()

	if c.GetDBPasswordPath() != "" {
		c.setDBPassword()
	}

	if c.GetKeycloakEnabled() {
		if c.GetKeycloakClientSecretPath() != "" {
			log.Info().Str("keycloak.clientSecretPath", c.GetKeycloakClientSecretPath()).Msg("Keycloak is enabled, reading client secret")
			c.SetKeycloakClientSecret()
		} else {
			log.Panic().Msg("Keycloak is enabled but client secret path is not set")
		}
	} else {
		log.Info().Bool("keycloak.enabled", c.GetKeycloakEnabled()).Msg("Keycloak is disabled")
	}

	if c.GetTemporalEncryptionKey() == "" {
		if c.GetTemporalEncryptionKeyPath() != "" {
			c.SetTemporalEncryptionKey()
		} else {
			log.Panic().Msg("Temporal encryption key path is not set")
		}
	}

	// Validate config
	c.Validate()

	// Watch secret files
	c.WatchSecretFilePaths()

	config = &c

	return config
}

// Validate validates the config
func (c *Config) Validate() {
	// Check configs that are essential but were not set
	if c.GetDBName() == "" {
		log.Panic().Msg("DB name config must be specified")
	}

	if c.GetDBUser() == "" {
		log.Panic().Msg("DB user config must be specified")
	}

	if value := os.Getenv("NO_DB_PASSWORD_OK"); strings.EqualFold(value, "true") {
		log.Info().Msg("connecting to DB without password")
	} else {
		if c.GetDBPassword() == "" && c.GetDBPasswordPath() == "" {
			log.Panic().Msg("DB password or password path config must be specified")
		}
	}

	if c.GetTemporalTlsEnabled() {
		if c.GetTemporalCertPath() == "" {
			log.Panic().Msg("Temporal cert path config must be specified")
		}

		if c.GetTemporalKeyPath() == "" {
			log.Panic().Msg("Temporal key path config must be specified")
		}

		if c.GetTemporalCaPath() == "" {
			log.Panic().Msg("Temporal CA path config must be specified")
		}
	}

	if c.GetTemporalEncryptionKey() == "" {
		log.Error().Msg("temporal encryption key config was not specified, arguments won't be encrypted")
	}

	// SiteConfig validations
	if err := c.ValidateSiteConfig(); err != nil {
		log.Panic().Err(err).Msg("SiteConfig must be specified")
	}

	// Validate that at least one auth method is configured
	authConfigs := c.GetAuthConfigs()

	authConfigNameServiceAccountMap := make(map[string]bool)

	for _, authConfig := range authConfigs {
		if authConfig.Name == "" {
			log.Panic().Msg("Auth configuration name must be specified")
		}

		if authConfig.Origin >= cauth.TokenOriginMax {
			log.Panic().Msgf("Auth configuration origin must be less than %v", cauth.TokenOriginMax)
		}

		if authConfig.URL == "" {
			log.Panic().Msg("Auth configuration URL must be specified")
		}

		if authConfig.Issuer == "" {
			log.Panic().Msg("Auth configuration issuer must be specified")
		}

		if authConfig.ServiceAccount {
			if !c.GetEnvDisconnected() {
				log.Panic().Msg("Service account is only supported in disconnected mode")
			}

			_, exists := authConfigNameServiceAccountMap[authConfig.Name]
			if exists {
				log.Panic().Msgf("Auth configuration name '%s' is already defined with service account enabled", authConfig.Name)
			}
			authConfigNameServiceAccountMap[authConfig.Name] = true
		}
	}

	keycloakEnabled := c.GetKeycloakEnabled()
	// Keycloak validations
	if err := c.ValidateKeycloakConfig(); err != nil {
		log.Panic().Err(err).Msg("Keycloak config must be specified")
	}

	if len(authConfigs) == 0 && !keycloakEnabled {
		log.Panic().Msg("No auth configurations specified and Keycloak is disabled - authentication will not work")
	} else if len(authConfigs) > 0 && keycloakEnabled {
		log.Info().Msg("Both auth configurations and Keycloak are enabled - multiple auth methods will be available")
	} else if len(authConfigs) > 0 {
		log.Info().Msgf("Auth configurations loaded: %d", len(authConfigs))
	} else if keycloakEnabled {
		log.Info().Msg("Keycloak authentication is enabled")
	}
}

// Config setters
// setLogLevel sets the log level
func (c *Config) setLogLevel() {
	logLevel := c.GetLogLevel()
	log.Warn().Str("level", logLevel).Msg("setting log level")
	level, err := zerolog.ParseLevel(logLevel)
	if err != nil {
		log.Err(err).Str("level", logLevel).Msg("failed to parse log level from config")
	}
	zerolog.SetGlobalLevel(level)
}

// setDBPassword sets the password for the database
func (c *Config) setDBPassword() {
	log.Warn().Str("db.passwordPath", c.GetDBPasswordPath()).Msg("setting password by reading from password path")

	passwordBytes, err := os.ReadFile(c.GetDBPasswordPath())
	if err != nil {
		log.Err(err).Str("db.passwordPath", c.GetDBPasswordPath()).Msg("failed to read password from file")
	}
	c.Lock()
	defer c.Unlock()
	c.v.Set(ConfigDBPassword, string(passwordBytes))
}

/* Get sub-configurations */
// GetDBConfig returns the database config
func (c *Config) GetDBConfig() *DBConfig {
	if c.db == nil {
		c.db = NewDBConfig(c.GetDBHost(), c.GetDBPort(), c.GetDBName(), c.GetDBUser(), c.GetDBPassword())
	}
	return c.db
}

// GetTemporalConfig returns the temporal config
func (c *Config) GetTemporalConfig() (*TemporalConfig, error) {
	var err error
	if c.temporal == nil {
		c.temporal, err = NewTemporalConfig(c.GetTemporalHost(), c.GetTemporalPort(), c.GetTemporalServerName(), c.GetTemporalNamespace(), c.GetTemporalQueue(), c.GetTemporalEncryptionKey(), c.GetTemporalTlsEnabled(), c.GetTemporalCertPath(), c.GetTemporalKeyPath(), c.GetTemporalCaPath())
	}
	return c.temporal, err
}

// GetOrInitJWTOrigin returns the JWT origin config with all configured auth providers
func (c *Config) GetOrInitJWTOriginConfig() *cauth.JWTOriginConfig {
	if c.JwtOriginConfig == nil {
		c.JwtOriginConfig = cauth.NewJWTOriginConfig()

		// Add auth configurations from the config file
		authConfigs := c.GetAuthConfigs()
		for _, authConfig := range authConfigs {
			c.JwtOriginConfig.AddConfig(authConfig.Name, authConfig.Issuer, authConfig.URL, authConfig.Origin, authConfig.ServiceAccount, authConfig.Audiences, authConfig.Scopes)
		}

		// Add Keycloak configuration if enabled
		if c.GetKeycloakEnabled() {
			keycloakConfig, err := c.GetOrInitKeycloakConfig()
			if err != nil {
				log.Warn().Err(err).Msg("Failed to get Keycloak config, skipping Keycloak JWT origin")
			} else {
				jwksConfig, err := keycloakConfig.GetJwksConfig()
				if err != nil {
					log.Warn().Err(err).Msg("Failed to get Keycloak JWKS config, skipping Keycloak JWT origin")
				} else {
					c.JwtOriginConfig.AddConfig("keycloak", jwksConfig.Issuer, jwksConfig.URL, cauth.TokenOriginKeycloak, c.GetKeycloakServiceAccountEnabled(), nil, nil)
				}
			}
		}

		// Initialize JWKS data
		if err := c.JwtOriginConfig.UpdateJWKs(); err != nil {
			log.Warn().Err(err).Msg("Failed to update JWKS data, skipping Keycloak JWT origin")
			return nil
		} else {
			log.Info().Msg("Successfully updated JWKS data")
		}
	}

	return c.JwtOriginConfig
}

// GetSiteConfig returns the Site config
func (c *Config) GetSiteConfig() *SiteConfig {
	if c.SiteConfig == nil {
		c.SiteConfig = NewSiteConfig(c.GetEnvDisconnected(), c.GetSitePhoneHomeUrl())
	}

	return c.SiteConfig
}

// GetMetricsConfig returns the Metrics config
func (c *Config) GetMetricsConfig() *MetricsConfig {
	return NewMetricsConfig(c.GetMetricsEnabled(), c.GetMetricsPort())
}

// GetAuthConfigs returns the auth configurations from the config file
func (c *Config) GetAuthConfigs() []AuthConfig {
	var authConfigs []AuthConfig
	if err := c.v.UnmarshalKey("auth", &authConfigs); err != nil {
		log.Warn().Err(err).Msg("Failed to unmarshal auth configurations, using empty list")
		return []AuthConfig{}
	}
	return authConfigs
}

/* Config getters */

// GetPathToConfig returns the path to the config file
func (c *Config) GetPathToConfig() string {
	return c.v.GetString(ConfigFilePath)
}

// GetEnvDev returns if the service is running in development mode
func (c *Config) GetEnvDev() bool {
	s := c.v.GetBool(ConfigEnvDev)
	return s
}

// GetEnvDisconnected returns if the service is running in disconnected mode
func (c *Config) GetEnvDisconnected() bool {
	return c.v.GetBool(ConfigEnvDisconnected)
}

// GetLogLevel returns the log level
func (c *Config) GetLogLevel() string {
	return c.v.GetString(ConfigLogLevel)
}

// GetSentryDSN returns the DSN for Sentry
func (c *Config) GetSentryDSN() string {
	return c.v.GetString(ConfigSentryDSN)
}

// GetDBHost returns the host of the database
func (c *Config) GetDBHost() string {
	return c.v.GetString(ConfigDBHost)
}

// GetDBPort returns the port of the database
func (c *Config) GetDBPort() int {
	return c.v.GetInt(ConfigDBPort)
}

// GetDBName returns the name of the database
func (c *Config) GetDBName() string {
	return c.v.GetString(ConfigDBName)
}

// GetDBUser returns the user of the database
func (c *Config) GetDBUser() string {
	return c.v.GetString(ConfigDBUser)
}

// GetDBPasswordPath returns the path for file containing password of the database
func (c *Config) GetDBPasswordPath() string {
	return c.v.GetString(ConfigDBPasswordPath)
}

// GetDBPassword returns the password of the database
func (c *Config) GetDBPassword() string {
	c.RLock()
	defer c.RUnlock()
	return c.v.GetString(ConfigDBPassword)
}

// GetTemporalHost returns the hostname for Temporal
func (c *Config) GetTemporalHost() string {
	return c.v.GetString(ConfigTemporalHost)
}

// GetTemporalPort returns the port for Temporal
func (c *Config) GetTemporalPort() int {
	return c.v.GetInt(ConfigTemporalPort)
}

// GetTemporalServerName returns the server name for Temporal
func (c *Config) GetTemporalServerName() string {
	return c.v.GetString(ConfigTemporalServerName)
}

// GetTemporalNamespace returns the namespace for Temporal
func (c *Config) GetTemporalNamespace() string {
	return c.v.GetString(ConfigTemporalNamespace)
}

// GetTemporalQueue returns the queue for Temporal
func (c *Config) GetTemporalQueue() string {
	return c.v.GetString(ConfigTemporalQueue)
}

// GetTemporalTlsEnabled returns true if TLS is enabled for Temporal connection
func (c *Config) GetTemporalTlsEnabled() bool {
	return c.v.GetBool(ConfigTemporalTlsEnabled)
}

// GetTemporalCertPath returns the path for client cert for Temporal
func (c *Config) GetTemporalCertPath() string {
	return c.v.GetString(ConfigTemporalCertPath)
}

// SetTemporalCertPath set the path for client cert for Temporal
func (c *Config) SetTemporalCertPath(value string) {
	c.v.Set(ConfigTemporalCertPath, value)
}

// GetTemporalKeyPath returns the path for client key for Temporal
func (c *Config) GetTemporalKeyPath() string {
	return c.v.GetString(ConfigTemporalKeyPath)
}

// SetTemporalKeyPath set the path for client key for Temporal
func (c *Config) SetTemporalKeyPath(value string) {
	c.v.Set(ConfigTemporalKeyPath, value)
}

// GetTemporalCaPath returns the path for ca cert for Temporal
func (c *Config) GetTemporalCaPath() string {
	return c.v.GetString(ConfigTemporalCaPath)
}

// SetTemporalCaPath sets the path for ca for Temporal
func (c *Config) SetTemporalCaPath(value string) {
	c.v.Set(ConfigTemporalCaPath, value)
}

// GetTemporalEncryptionKeyPath returns the path for file containing encryption key for Temporal
func (c *Config) GetTemporalEncryptionKeyPath() string {
	return c.v.GetString(ConfigTemporalEncryptionKeyPath)
}

// SetTemporalEncryptionKeyPath sets the path for file containing encryption key for Temporal
func (c *Config) SetTemporalEncryptionKeyPath(value string) {
	c.v.Set(ConfigTemporalEncryptionKeyPath, value)
}

// GetTemporalEncryptionKey returns the encryption key for Temporal
func (c *Config) GetTemporalEncryptionKey() string {
	return c.v.GetString(ConfigTemporalEncryptionKey)
}

// SetTemporalEncryptionKey sets the encryption key for Temporal
func (c *Config) SetTemporalEncryptionKey() {
	log.Warn().Str("temporal.encryptionKeyPath", c.GetTemporalEncryptionKeyPath()).Msg("setting Temporal encryption key by reading from file")

	encryptionKeyBytes, err := os.ReadFile(c.GetTemporalEncryptionKeyPath())
	if err != nil {
		log.Panic().Err(err).Msgf("failed to read encryption key file: %s", err)
	}
	encryptionKey := strings.TrimSpace(string(encryptionKeyBytes))

	c.v.Set(ConfigTemporalEncryptionKey, encryptionKey)
}

// ValidateSiteConfig validates Site configs
func (c *Config) ValidateSiteConfig() error {
	if c.GetSitePhoneHomeUrl() == "" {
		return fmt.Errorf("invalid Site PhoneHome url")
	}
	return nil
}

// SetSiteManagerEnabled sets the enabled field
func (c *Config) SetSiteManagerEnabled(value bool) {
	c.v.Set(ConfigSiteManagerEnabled, value)
}

// GetSiteManagerEnabled sets the enabled field
func (c *Config) GetSiteManagerEnabled() bool {
	return c.v.GetBool(ConfigSiteManagerEnabled)
}

// SetSiteManagerEndpoint sets the endpoint
func (c *Config) SetSiteManagerEndpoint(value string) {
	c.v.Set(ConfigSiteManagerEndpoint, value)
}

// GetSiteManagerEndpoint gets the path for site manager service endpoint
func (c *Config) GetSiteManagerEndpoint() string {
	return c.v.GetString(ConfigSiteManagerEndpoint)
}

// SetSitePhoneHomeUrl sets the url for PhoneHome
func (c *Config) SetSitePhoneHomeUrl(value string) {
	c.v.Set(ConfigSitePhoneHomeUrl, value)
}

// GetSitePhoneHomeUrl gets the url for PhoneHome
func (c *Config) GetSitePhoneHomeUrl() string {
	return c.v.GetString(ConfigSitePhoneHomeUrl)
}

// GetMetricsEnabled gets the enabled field for Metrics
func (c *Config) GetMetricsEnabled() bool {
	return c.v.GetBool(ConfigMetricsEnabled)
}

// GetZincSearchPort gets the port for Metrics
func (c *Config) GetMetricsPort() int {
	return c.v.GetInt(ConfigMetricsPort)
}

// GetTracingEnabled gets the enabled field for tracing
func (c *Config) GetTracingEnabled() bool {
	return c.v.GetBool(ConfigTracingEnabled)
}

// GetTracingServiceName gets the service name for tracing
func (c *Config) GetTracingServiceName() string {
	return c.v.GetString(ConfigTracingServiceName)
}

// Keycloak configuration methods

// GetKeycloakEnabled gets the enabled field for Keycloak
func (c *Config) GetKeycloakEnabled() bool {
	return c.v.GetBool(ConfigKeycloakEnabled)
}

// SetKeycloakEnabled sets the enabled field for Keycloak
func (c *Config) SetKeycloakEnabled(value bool) {
	c.v.Set(ConfigKeycloakEnabled, value)
}

// GetKeycloakBaseURL gets the base URL for Keycloak
func (c *Config) GetKeycloakBaseURL() string {
	return c.v.GetString(ConfigKeycloakBaseURL)
}

// SetKeycloakBaseURL sets the base URL for Keycloak
func (c *Config) SetKeycloakBaseURL(value string) {
	c.v.Set(ConfigKeycloakBaseURL, value)
}

// GetKeycloakExternalBaseURL gets the external base URL for Keycloak
func (c *Config) GetKeycloakExternalBaseURL() string {
	return c.v.GetString(ConfigKeycloakExternalBaseURL)
}

// SetKeycloakExternalBaseURL sets the external base URL for Keycloak
func (c *Config) SetKeycloakExternalBaseURL(value string) {
	c.v.Set(ConfigKeycloakExternalBaseURL, value)
}

// GetKeycloakRealm gets the realm for Keycloak
func (c *Config) GetKeycloakRealm() string {
	return c.v.GetString(ConfigKeycloakRealm)
}

// SetKeycloakRealm sets the realm for Keycloak
func (c *Config) SetKeycloakRealm(value string) {
	c.v.Set(ConfigKeycloakRealm, value)
}

// GetKeycloakClientID gets the client ID for Keycloak
func (c *Config) GetKeycloakClientID() string {
	return c.v.GetString(ConfigKeycloakClientID)
}

// SetKeycloakClientID sets the client ID for Keycloak
func (c *Config) SetKeycloakClientID(value string) {
	c.v.Set(ConfigKeycloakClientID, value)
}

// GetKeycloakClientSecretPath gets the client secret path for Keycloak
func (c *Config) GetKeycloakClientSecretPath() string {
	return c.v.GetString(ConfigKeycloakClientSecretPath)
}

// SetKeycloakClientSecretPath sets the client secret path for Keycloak
func (c *Config) SetKeycloakClientSecretPath(value string) {
	c.v.Set(ConfigKeycloakClientSecretPath, value)
}

// GetKeycloakClientSecret gets the client secret for Keycloak
func (c *Config) GetKeycloakClientSecret() string {
	return c.v.GetString(ConfigKeycloakClientSecret)
}

// SetKeycloakClientSecret sets the client secret for Keycloak
func (c *Config) SetKeycloakClientSecret() {
	log.Warn().Str("keycloak.clientSecretPath", c.GetKeycloakClientSecretPath()).Msg("setting keycloak client secret by reading from secret path")

	secretBytes, err := os.ReadFile(c.GetKeycloakClientSecretPath())
	if err != nil {
		log.Panic().Err(err).Msgf("failed to read keycloak client secret file: %s", err)
	}
	secret := strings.TrimSpace(string(secretBytes))

	c.v.Set(ConfigKeycloakClientSecret, secret)
}

// GetKeycloakServiceAccountEnabled gets the service account enabled flag for Keycloak
func (c *Config) GetKeycloakServiceAccountEnabled() bool {
	return c.v.GetBool(ConfigKeycloakServiceAccountEnabled)
}

// SetKeycloakServiceAccountEnabled sets the service account enabled flag for Keycloak
func (c *Config) SetKeycloakServiceAccountEnabled(value bool) {
	c.v.Set(ConfigKeycloakServiceAccountEnabled, value)
}

// GetOrInitKeycloakConfig gets the Keycloak configuration
func (c *Config) GetOrInitKeycloakConfig() (*cauth.KeycloakConfig, error) {
	if c.KeycloakConfig == nil {
		if !c.GetKeycloakEnabled() {
			log.Info().Msg("Keycloak is disabled in configuration")
			return nil, errors.New("Keycloak is not enabled")
		}

		log.Info().Str("baseURL", c.GetKeycloakBaseURL()).
			Str("externalBaseURL", c.GetKeycloakExternalBaseURL()).
			Str("realm", c.GetKeycloakRealm()).
			Str("clientID", c.GetKeycloakClientID()).
			Str("clientSecretPath", c.GetKeycloakClientSecretPath()).
			Msg("Creating new Keycloak configuration")

		c.KeycloakConfig = cauth.NewKeycloakConfig(
			c.GetKeycloakBaseURL(),
			c.GetKeycloakExternalBaseURL(),
			c.GetKeycloakClientID(),
			c.GetKeycloakClientSecret(),
			c.GetKeycloakRealm(),
			c.GetKeycloakServiceAccountEnabled(),
		)
		log.Info().Msg("Keycloak configuration created successfully")
	}
	return c.KeycloakConfig, nil
}

// ValidateKeycloakConfig validates the Keycloak configuration
func (c *Config) ValidateKeycloakConfig() error {
	if !c.GetKeycloakEnabled() {
		return nil // Skip validation if Keycloak is not enabled
	}

	if c.GetKeycloakBaseURL() == "" {
		return fmt.Errorf("keycloak base URL must be specified when Keycloak is enabled")
	}

	if c.GetKeycloakRealm() == "" {
		return fmt.Errorf("keycloak realm must be specified when Keycloak is enabled")
	}

	if c.GetKeycloakClientID() == "" {
		return fmt.Errorf("keycloak client ID must be specified when Keycloak is enabled")
	}

	if c.GetKeycloakClientSecret() == "" && c.GetKeycloakClientSecretPath() == "" {
		return fmt.Errorf("keycloak client secret or client secret path must be specified when Keycloak is enabled")
	}

	return nil
}

// WatchSecretFilePaths starts watching secret files for changes.
func (c *Config) WatchSecretFilePaths() {
	log.Info().Str("WatchSecretFilePaths", "").Msg("started watching secret file paths")
	initWG := sync.WaitGroup{}
	initWG.Add(1)
	go func() {
		watcher, err := fsnotify.NewWatcher()
		if err != nil {
			log.Panic().Err(err).Msgf("fatal error while creating watcher for secret files: %s", err)
		}
		defer watcher.Close()

		// Add new secret file and path here

		// db password secret file and path
		dbPassFilePath := filepath.Clean(c.GetDBPasswordPath())
		dbPassFileDir := ""
		if dbPassFilePath != "" {
			dbPassFileDir, _ = filepath.Split(dbPassFilePath)
		}

		// keycloak client secret file and path
		keycloakClientSecretFilePath := filepath.Clean(c.GetKeycloakClientSecretPath())
		keycloakClientSecretFileDir := ""
		if keycloakClientSecretFilePath != "" {
			keycloakClientSecretFileDir, _ = filepath.Split(keycloakClientSecretFilePath)
		}

		eventsWG := sync.WaitGroup{}
		eventsWG.Add(1)
		go func() {
			for {
				select {
				case event, ok := <-watcher.Events:
					if !ok { // 'Events' channel is closed
						eventsWG.Done()
						return
					}
					if !(event.Has(fsnotify.Write) || event.Has(fsnotify.Create)) {
						continue
					}
					eventFilePath := filepath.Clean(event.Name)
					if eventFilePath == "" {
						continue
					}
					eventFileDir, _ := filepath.Split(eventFilePath)
					if eventFileDir == dbPassFileDir {
						c.setDBPassword()
					}
					if eventFileDir == keycloakClientSecretFileDir && c.GetKeycloakEnabled() && keycloakClientSecretFileDir != "" {
						c.SetKeycloakClientSecret()
					}
				case err, ok := <-watcher.Errors:
					if ok { // 'Errors' channel is not closed
						log.Info().Err(err).Msgf("one or more secret file received watcher error: %s", err)
					}
					eventsWG.Done()
					return
				}
			}
		}()

		// Add secret file paths to watcher
		log.Info().Str("WatchSecretFilePaths", dbPassFileDir).Msg("adding db path for watcher")
		watcher.Add(dbPassFileDir)

		if c.GetKeycloakEnabled() {
			log.Info().Str("WatchSecretFilePaths", keycloakClientSecretFileDir).Msg("adding keycloak client secret path for watcher")
			watcher.Add(keycloakClientSecretFileDir)
		}

		initWG.Done()   // done initializing the watch in this go routine, so the parent routine can move on...
		eventsWG.Wait() // now, wait for event loop to end in this go-routine...
	}()
	initWG.Wait() // make sure that the go routine above fully ended before returning
}

// Close stops background tasks
func (c *Config) Close() {
	if c.temporal != nil {
		c.temporal.Close()
	}
}
