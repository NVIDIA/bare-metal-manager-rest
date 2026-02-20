// SPDX-FileCopyrightText: Copyright (c) 2026 NVIDIA CORPORATION & AFFILIATES. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package bmmcli

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"golang.org/x/term"
	"gopkg.in/yaml.v3"
)

// ConfigFile mirrors the ~/.bmm/config.yaml structure.
type ConfigFile struct {
	API   ConfigAPI   `yaml:"api"`
	Auth  ConfigAuth  `yaml:"auth"`
	Debug ConfigDebug `yaml:"debug"`
}

type ConfigDebug struct {
	// Verbose enables detailed HTTP logging (status, content-type, raw body).
	// Defaults to true so errors are always self-explaining.
	Verbose *bool `yaml:"verbose,omitempty"`
}

// IsVerbose returns true when verbose debug logging is enabled.
// Defaults to false when the field is unset.
func (d ConfigDebug) IsVerbose() bool {
	return d.Verbose != nil && *d.Verbose
}

type ConfigAPI struct {
	Base string `yaml:"base,omitempty"`
	Org  string `yaml:"org,omitempty"`
	Name string `yaml:"name,omitempty"`
}

type ConfigAuth struct {
	Token  string        `yaml:"token,omitempty"`
	OIDC   *ConfigOIDC   `yaml:"oidc,omitempty"`
	APIKey *ConfigAPIKey `yaml:"api_key,omitempty"`
}

type ConfigOIDC struct {
	TokenURL     string `yaml:"token_url,omitempty"`
	ClientID     string `yaml:"client_id,omitempty"`
	ClientSecret string `yaml:"client_secret,omitempty"`
	Username     string `yaml:"username,omitempty"`
	Password     string `yaml:"password,omitempty"`
	Token        string `yaml:"token,omitempty"`
	RefreshToken string `yaml:"refresh_token,omitempty"`
	ExpiresAt    string `yaml:"expires_at,omitempty"`
}

type ConfigAPIKey struct {
	AuthnURL string `yaml:"authn_url,omitempty"`
	Key      string `yaml:"key,omitempty"`
	Token    string `yaml:"token,omitempty"`
}

// ConfigPath returns the default config file path.
func ConfigPath() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".bmm", "config.yaml")
}

// ConfigDir returns the ~/.bmm directory.
func ConfigDir() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".bmm")
}

// LoadConfig reads the default config file.
func LoadConfig() (*ConfigFile, error) {
	return LoadConfigFromPath(ConfigPath())
}

// LoadConfigFromPath reads a config file at a specific path.
func LoadConfigFromPath(path string) (*ConfigFile, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return &ConfigFile{}, nil
		}
		return nil, err
	}
	var cfg ConfigFile
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parsing config %s: %w", path, err)
	}
	return &cfg, nil
}

// SaveConfig writes the config back to ConfigPath(), preserving unknown keys.
func SaveConfig(cfg *ConfigFile) error {
	return SaveConfigToPath(cfg, ConfigPath())
}

// SaveConfigToPath writes the config to a specific path.
func SaveConfigToPath(cfg *ConfigFile, path string) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0700); err != nil {
		return fmt.Errorf("creating config directory: %w", err)
	}
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("marshaling config: %w", err)
	}
	return os.WriteFile(path, data, 0600)
}

// GetAuthToken returns the best available bearer token from the config.
func GetAuthToken(cfg *ConfigFile) string {
	if cfg.Auth.Token != "" {
		return cfg.Auth.Token
	}
	if cfg.Auth.OIDC != nil && cfg.Auth.OIDC.Token != "" {
		return cfg.Auth.OIDC.Token
	}
	if cfg.Auth.APIKey != nil && cfg.Auth.APIKey.Token != "" {
		return cfg.Auth.APIKey.Token
	}
	return ""
}

// HasOIDCConfig returns true when OIDC credentials are present in the config.
func HasOIDCConfig(cfg *ConfigFile) bool {
	return cfg.Auth.OIDC != nil &&
		cfg.Auth.OIDC.TokenURL != "" &&
		cfg.Auth.OIDC.ClientID != ""
}

// HasAPIKeyConfig returns true when NGC API key settings are present.
func HasAPIKeyConfig(cfg *ConfigFile) bool {
	return cfg.Auth.APIKey != nil &&
		cfg.Auth.APIKey.AuthnURL != "" &&
		cfg.Auth.APIKey.Key != ""
}

// ChooseConfigFile scans ~/.bmm/ for config*.yaml files and shows an interactive
// picker when more than one is found. Returns the path of the selected config,
// or "" to use the default.
func ChooseConfigFile() (string, *ConfigFile, error) {
	// Non-interactive environment: skip selection.
	if !term.IsTerminal(int(os.Stdin.Fd())) || !term.IsTerminal(int(os.Stdout.Fd())) {
		cfg, err := LoadConfig()
		return ConfigPath(), cfg, err
	}

	entries, err := os.ReadDir(ConfigDir())
	if err != nil {
		if os.IsNotExist(err) {
			cfg, err := LoadConfig()
			return ConfigPath(), cfg, err
		}
		return "", nil, fmt.Errorf("reading config directory: %w", err)
	}

	var candidates []string
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		if !strings.HasPrefix(name, "config") {
			continue
		}
		if !strings.HasSuffix(name, ".yaml") && !strings.HasSuffix(name, ".yml") {
			continue
		}
		candidates = append(candidates, filepath.Join(ConfigDir(), name))
	}

	// Only one (or zero) config files — no need to ask.
	if len(candidates) <= 1 {
		cfg, err := LoadConfig()
		return ConfigPath(), cfg, err
	}

	sortConfigCandidates(candidates)

	home, _ := os.UserHomeDir()
	items := make([]configSelectItem, len(candidates))
	for i, path := range candidates {
		items[i] = configSelectItem{
			label: displayConfigPath(path, home),
			path:  path,
		}
	}

	fmt.Println()
	selected, err := selectConfig("Choose a config for this session", items)
	if err != nil {
		return "", nil, err
	}
	fmt.Printf("Using: %s\n\n", selected.label)

	cfg, err := LoadConfigFromPath(selected.path)
	return selected.path, cfg, err
}

type configSelectItem struct {
	label string
	path  string
}

func selectConfig(label string, items []configSelectItem) (*configSelectItem, error) {
	if len(items) == 0 {
		return nil, fmt.Errorf("no config files found")
	}
	if len(items) == 1 {
		return &items[0], nil
	}

	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		return &items[0], nil
	}
	defer term.Restore(int(os.Stdin.Fd()), oldState)

	cursor := 0
	renderConfigSelect(label, items, cursor)

	buf := make([]byte, 3)
	for {
		n, err := os.Stdin.Read(buf)
		if err != nil || n == 0 {
			return nil, fmt.Errorf("reading input: %w", err)
		}

		switch {
		case buf[0] == 3 || buf[0] == 4: // Ctrl+C / Ctrl+D
			clearConfigSelect(len(items))
			return nil, fmt.Errorf("cancelled")
		case buf[0] == '\r' || buf[0] == '\n':
			clearConfigSelect(len(items))
			fmt.Printf("%s %s\r\n", label, items[cursor].label)
			return &items[cursor], nil
		case n >= 3 && buf[0] == 27 && buf[1] == '[' && buf[2] == 'A': // up
			if cursor > 0 {
				cursor--
			}
		case n >= 3 && buf[0] == 27 && buf[1] == '[' && buf[2] == 'B': // down
			if cursor < len(items)-1 {
				cursor++
			}
		}

		clearConfigSelect(len(items))
		renderConfigSelect(label, items, cursor)
	}
}

func renderConfigSelect(label string, items []configSelectItem, cursor int) {
	fmt.Printf("\033[?25l") // hide cursor
	fmt.Printf("\033[1m%s\033[0m \033[2m(arrows to move, enter to select)\033[0m\r\n", label)
	for i, item := range items {
		if i == cursor {
			fmt.Printf("  \033[36m>\033[0m \033[7m %s \033[0m\r\n", item.label)
		} else {
			fmt.Printf("    %s\r\n", item.label)
		}
	}
}

func clearConfigSelect(count int) {
	lines := 1 + count
	fmt.Printf("\033[%dA\033[1G\033[J", lines)
}

func sortConfigCandidates(paths []string) {
	sort.Slice(paths, func(i, j int) bool {
		li := filepath.Base(paths[i])
		lj := filepath.Base(paths[j])
		// config.yaml / config.yml sort first
		if (li == "config.yaml" || li == "config.yml") != (lj == "config.yaml" || lj == "config.yml") {
			return li == "config.yaml" || li == "config.yml"
		}
		return li < lj
	})
}

func displayConfigPath(path, home string) string {
	prefix := home + string(os.PathSeparator)
	if strings.HasPrefix(path, prefix) {
		return "~/" + strings.TrimPrefix(path, prefix)
	}
	return path
}

const sampleConfigContent = `# BMM CLI configuration
#
# API connection:
#   api.base -- server URL
#   api.org  -- organization name used in API paths
#   api.name -- API path segment (default: carbide)
#
# Authentication options (choose one):
#   auth.token      -- direct bearer token (no login required)
#   auth.oidc       -- OIDC password/client-credentials flow
#   auth.api_key    -- NGC API key exchange
#
api:
  base: http://localhost:8388
  org: test-org
  name: carbide

auth:
  # Option 1: Direct bearer token
  # token: eyJhbGciOi...

  # Option 2: OIDC provider (e.g. Keycloak)
  oidc:
    token_url: http://localhost:8080/realms/carbide-dev/protocol/openid-connect/token
    client_id: carbide-api
    client_secret: carbide-local-secret
    username: admin@example.com
    password: adminpassword

  # Option 3: NGC API key
  # api_key:
  #   authn_url: https://authn.nvidia.com/token
  #   key: nvapi-xxxx
`
