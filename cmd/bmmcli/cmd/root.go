// SPDX-FileCopyrightText: Copyright (c) 2026 NVIDIA CORPORATION & AFFILIATES. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

var cfgFile string

// Config represents the CLI configuration stored in ~/.bmm/config.yaml
type Config struct {
	API  APIConfig  `yaml:"api"`
	Auth AuthConfig `yaml:"auth"`
}

// APIConfig holds server connection settings
type APIConfig struct {
	Base string `yaml:"base,omitempty"`
	Org  string `yaml:"org,omitempty"`
	Name string `yaml:"name,omitempty"`
}

// AuthConfig holds authentication configuration.
// Token can be set directly for generic bearer-token auth.
// Method-specific sections can also exchange and store session tokens.
type AuthConfig struct {
	Token  string      `yaml:"token,omitempty"`
	OIDC   *OIDCAuth   `yaml:"oidc,omitempty"`
	APIKey *APIKeyAuth `yaml:"api_key,omitempty"`
}

// OIDCAuth holds OIDC provider settings and the resulting session token.
type OIDCAuth struct {
	TokenURL     string `yaml:"token_url,omitempty"`
	ClientID     string `yaml:"client_id,omitempty"`
	ClientSecret string `yaml:"client_secret,omitempty"`
	Username     string `yaml:"username,omitempty"`
	Password     string `yaml:"password,omitempty"`
	Token        string `yaml:"token,omitempty"`
}

// APIKeyAuth holds NGC API key settings and the resulting session token.
type APIKeyAuth struct {
	AuthnURL string `yaml:"authn_url,omitempty"`
	Key      string `yaml:"key,omitempty"`
	Token    string `yaml:"token,omitempty"`
}

var rootCmd = &cobra.Command{
	Use:   "bmm",
	Short: "Bare Metal Manager CLI",
	Long:  "CLI for interacting with the Carbide REST API.",
}

// Execute runs the root command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is ~/.bmm/config.yaml)")
	rootCmd.PersistentFlags().String("org", "", "override org from config")
	rootCmd.PersistentFlags().String("base", "", "override API base URL")
	rootCmd.PersistentFlags().StringP("output", "o", "", "output format: json, yaml (default from config)")

	configureViper()
}

func initConfig() {
	configureViper()
	if err := configureConfigSource(); err != nil {
		fmt.Fprintln(os.Stderr, "Error finding home directory:", err)
		os.Exit(1)
	}

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			fmt.Fprintln(os.Stderr, "Error reading config:", err)
		}
	}
}

func configureViper() {
	_ = viper.BindPFlag("api.org", rootCmd.PersistentFlags().Lookup("org"))
	_ = viper.BindPFlag("api.base", rootCmd.PersistentFlags().Lookup("base"))
	_ = viper.BindPFlag("output", rootCmd.PersistentFlags().Lookup("output"))

	viper.SetDefault("api.base", "http://localhost:8388")
	viper.SetDefault("api.org", "test-org")
	viper.SetDefault("api.name", "carbide")

	viper.SetEnvPrefix("BMM")
	viper.AutomaticEnv()
}

func configureConfigSource() error {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
		return nil
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	configDir := filepath.Join(home, ".bmm")
	viper.AddConfigPath(configDir)
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	return nil
}

func reloadConfigWithFile(configPath string) error {
	cfgFile = configPath
	viper.Reset()
	configureViper()
	if err := configureConfigSource(); err != nil {
		return fmt.Errorf("finding home directory: %w", err)
	}
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return fmt.Errorf("reading config: %w", err)
		}
	}
	return nil
}

// configFilePath returns the path to the config file being used, or the default path
func configFilePath() string {
	if cfgFile != "" {
		return cfgFile
	}
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".bmm", "config.yaml")
}

const configHeader = `# BMM CLI configuration
#
# API connection:
#   api.base -- server URL
#   api.org  -- organization name used in API paths
#   api.name -- API path segment (default: carbide)
#
# Authentication options:
#   auth.token -- direct bearer token (no "bmm login" required)
#   auth.oidc  -- OIDC password flow via "bmm login"
#   auth.api_key -- NGC API key exchange via "bmm login --api-key"
#
# If both OIDC and API key are configured, OIDC is used by default.
# Use "bmm login --api-key" to force NGC API key flow.
#
# The session token is stored under each method's "token" field
# after a successful login.
#
`

const authGuidance = `No authentication configured. Add one of these to %s:

  Option 1 -- Direct bearer token:
    auth:
      token: eyJhbGciOi...

  Option 2 -- NGC API key:
    auth:
      api_key:
        authn_url: https://authn.nvidia.com/token
        key: nvapi-xxxx

  Option 3 -- OIDC provider (e.g. Keycloak):
    auth:
      oidc:
        token_url: http://localhost:8080/realms/carbide-dev/protocol/openid-connect/token
        client_id: carbide-api
        client_secret: carbide-local-secret
        username: admin@example.com
        password: adminpassword

Then run: bmm login (not needed for Option 1)`

// getAuthToken returns the best available session token.
// Prefers direct auth.token, then OIDC/API key session tokens.
func getAuthToken() (string, error) {
	if t := viper.GetString("auth.token"); t != "" {
		return t, nil
	}
	if t := viper.GetString("auth.oidc.token"); t != "" {
		return t, nil
	}
	if t := viper.GetString("auth.api_key.token"); t != "" {
		return t, nil
	}
	return "", fmt.Errorf(authGuidance, configFilePath())
}

func hasDirectTokenConfig() bool {
	return strings.TrimSpace(viper.GetString("auth.token")) != ""
}

// hasOIDCConfig returns true when OIDC provider settings are configured.
func hasOIDCConfig() bool {
	return viper.IsSet("auth.oidc.token_url") &&
		viper.IsSet("auth.oidc.client_id") &&
		viper.IsSet("auth.oidc.client_secret") &&
		viper.GetString("auth.oidc.token_url") != "" &&
		viper.GetString("auth.oidc.client_id") != "" &&
		viper.GetString("auth.oidc.client_secret") != ""
}

// hasAPIKeyConfig returns true when NGC API key settings are configured.
func hasAPIKeyConfig() bool {
	return viper.IsSet("auth.api_key.authn_url") &&
		viper.IsSet("auth.api_key.key") &&
		viper.GetString("auth.api_key.authn_url") != "" &&
		viper.GetString("auth.api_key.key") != ""
}

// readConfig reads the config file into a Config struct.
func readConfig() Config {
	cfgPath := configFilePath()
	var cfg Config
	if data, err := os.ReadFile(cfgPath); err == nil {
		_ = yaml.Unmarshal(data, &cfg)
	}
	return cfg
}

// saveToken reads the config file, sets the token for the given auth method,
// and writes it back. Only the token field is changed.
func saveToken(method, token string) error {
	cfgPath := configFilePath()

	dir := filepath.Dir(cfgPath)
	if err := os.MkdirAll(dir, 0700); err != nil {
		return fmt.Errorf("creating config directory: %w", err)
	}

	cfg := readConfig()

	switch method {
	case "oidc":
		if cfg.Auth.OIDC == nil {
			cfg.Auth.OIDC = &OIDCAuth{}
		}
		cfg.Auth.OIDC.Token = token
	case "api_key":
		if cfg.Auth.APIKey == nil {
			cfg.Auth.APIKey = &APIKeyAuth{}
		}
		cfg.Auth.APIKey.Token = token
	default:
		return fmt.Errorf("unknown auth method: %s", method)
	}

	data, err := yaml.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("marshaling config: %w", err)
	}

	if err := os.WriteFile(cfgPath, []byte(data), 0600); err != nil {
		return fmt.Errorf("writing config file: %w", err)
	}

	return nil
}
