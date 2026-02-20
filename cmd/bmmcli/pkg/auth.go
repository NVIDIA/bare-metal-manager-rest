// SPDX-FileCopyrightText: Copyright (c) 2026 NVIDIA CORPORATION & AFFILIATES. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package bmmcli

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"golang.org/x/term"
	cli "github.com/urfave/cli/v2"
)


// TokenResponse is the OAuth2 token endpoint response.
type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
	TokenType    string `json:"token_type"`
}

// LoginCommand returns the 'login' CLI command.
func LoginCommand() *cli.Command {
	return &cli.Command{
		Name:  "login",
		Usage: "Authenticate with the Carbide API and save the token",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "username",
				Usage: "Username for OIDC password grant",
			},
			&cli.StringFlag{
				Name:  "password",
				Usage: "Password for OIDC password grant (prompted if not provided)",
			},
			&cli.StringFlag{
				Name:    "client-secret",
				Usage:   "Client secret (required for confidential OIDC clients)",
				EnvVars: []string{"BMM_CLIENT_SECRET"},
			},
			&cli.StringFlag{
				Name:    "api-key",
				Usage:   "NGC API key for token exchange",
				EnvVars: []string{"BMM_API_KEY"},
			},
			&cli.StringFlag{
				Name:    "authn-url",
				Usage:   "NGC authentication URL for API key exchange",
				EnvVars: []string{"BMM_AUTHN_URL"},
				Value:   "https://authn.nvidia.com/token",
			},
		},
		Action: func(c *cli.Context) error {
			cfg, _ := LoadConfig()

			apiKey := c.String("api-key")
			if apiKey == "" && HasAPIKeyConfig(cfg) {
				apiKey = cfg.Auth.APIKey.Key
			}
			if apiKey != "" {
				authnURL := c.String("authn-url")
				if authnURL == "https://authn.nvidia.com/token" && cfg.Auth.APIKey != nil && cfg.Auth.APIKey.AuthnURL != "" {
					authnURL = cfg.Auth.APIKey.AuthnURL
				}
				return loginWithAPIKey(cfg, authnURL, apiKey)
			}
			return loginWithOIDCCmd(c, cfg)
		},
	}
}

func loginWithAPIKey(cfg *ConfigFile, authnURL, apiKey string) error {
	req, err := http.NewRequest("GET", authnURL, nil)
	if err != nil {
		return fmt.Errorf("building request: %w", err)
	}
	req.Header.Set("Authorization", "ApiKey "+apiKey)
	req.Header.Set("Accept", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("requesting token from NGC: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("reading response: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("NGC token exchange failed (HTTP %d): %s", resp.StatusCode, string(body))
	}

	token := extractNGCToken(body)
	if token == "" {
		return fmt.Errorf("NGC response did not contain a token")
	}

	if cfg.Auth.APIKey == nil {
		cfg.Auth.APIKey = &ConfigAPIKey{}
	}
	cfg.Auth.APIKey.Token = token
	if err := SaveConfig(cfg); err != nil {
		return fmt.Errorf("saving config: %w", err)
	}
	fmt.Fprintf(os.Stderr, "Login successful (NGC API key). Token saved to %s\n", ConfigPath())
	return nil
}

func extractNGCToken(body []byte) string {
	var resp struct {
		Token       string `json:"token"`
		AccessToken string `json:"access_token"`
	}
	if json.Unmarshal(body, &resp) != nil {
		return ""
	}
	if resp.Token != "" {
		return resp.Token
	}
	return resp.AccessToken
}

func loginWithOIDCCmd(c *cli.Context, cfg *ConfigFile) error {
	// Resolve token URL: flag → config → keycloak-url flag
	tokenURL := c.String("token-url")
	if tokenURL == "" && cfg.Auth.OIDC != nil {
		tokenURL = cfg.Auth.OIDC.TokenURL
	}
	if tokenURL == "" {
		if keycloakURL := c.String("keycloak-url"); keycloakURL != "" {
			realm := c.String("realm")
			tokenURL = fmt.Sprintf("%s/realms/%s/protocol/openid-connect/token",
				strings.TrimRight(keycloakURL, "/"), realm)
		}
	}
	if tokenURL == "" {
		return fmt.Errorf("--token-url or --keycloak-url is required (or set auth.oidc.token_url in config)")
	}

	clientID := c.String("client-id")
	if clientID == "" && cfg.Auth.OIDC != nil {
		clientID = cfg.Auth.OIDC.ClientID
	}

	clientSecret := c.String("client-secret")
	if clientSecret == "" && cfg.Auth.OIDC != nil {
		clientSecret = cfg.Auth.OIDC.ClientSecret
	}

	username := c.String("username")
	if username == "" && cfg.Auth.OIDC != nil {
		username = cfg.Auth.OIDC.Username
	}

	password := c.String("password")
	if password == "" && cfg.Auth.OIDC != nil {
		password = cfg.Auth.OIDC.Password
	}

	var tokenResp *TokenResponse
	var err error

	if username == "" && clientSecret != "" {
		tokenResp, err = clientCredentialsGrant(tokenURL, clientID, clientSecret)
	} else {
		if username == "" {
			fmt.Print("Username: ")
			fmt.Scanln(&username)
		}
		if password == "" {
			fmt.Print("Password: ")
			pw, pwErr := term.ReadPassword(int(os.Stdin.Fd()))
			fmt.Println()
			if pwErr != nil {
				return fmt.Errorf("reading password: %w", pwErr)
			}
			password = string(pw)
		}
		tokenResp, err = passwordGrant(tokenURL, clientID, clientSecret, username, password)
	}
	if err != nil {
		return err
	}

	if cfg.Auth.OIDC == nil {
		cfg.Auth.OIDC = &ConfigOIDC{}
	}
	cfg.Auth.OIDC.Token = tokenResp.AccessToken
	cfg.Auth.OIDC.RefreshToken = tokenResp.RefreshToken
	cfg.Auth.OIDC.ExpiresAt = time.Now().Add(time.Duration(tokenResp.ExpiresIn) * time.Second).Format(time.RFC3339)
	cfg.Auth.OIDC.TokenURL = tokenURL
	cfg.Auth.OIDC.ClientID = clientID
	cfg.Auth.OIDC.ClientSecret = clientSecret

	if err := SaveConfig(cfg); err != nil {
		return fmt.Errorf("saving config: %w", err)
	}
	fmt.Fprintf(os.Stderr, "Login successful. Token saved to %s\n", ConfigPath())
	return nil
}

func passwordGrant(tokenURL, clientID, clientSecret, username, password string) (*TokenResponse, error) {
	data := url.Values{
		"grant_type": {"password"},
		"client_id":  {clientID},
		"username":   {username},
		"password":   {password},
		"scope":      {"openid"},
	}
	if clientSecret != "" {
		data.Set("client_secret", clientSecret)
	}
	return postToken(tokenURL, data)
}

func clientCredentialsGrant(tokenURL, clientID, clientSecret string) (*TokenResponse, error) {
	data := url.Values{
		"grant_type":    {"client_credentials"},
		"client_id":     {clientID},
		"client_secret": {clientSecret},
		"scope":         {"openid"},
	}
	return postToken(tokenURL, data)
}

func refreshTokenGrant(tokenURL, clientID, clientSecret, refreshToken string) (*TokenResponse, error) {
	data := url.Values{
		"grant_type":    {"refresh_token"},
		"client_id":     {clientID},
		"refresh_token": {refreshToken},
	}
	if clientSecret != "" {
		data.Set("client_secret", clientSecret)
	}
	return postToken(tokenURL, data)
}

func postToken(tokenURL string, data url.Values) (*TokenResponse, error) {
	resp, err := http.PostForm(tokenURL, data)
	if err != nil {
		return nil, fmt.Errorf("token request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		var errBody struct {
			Error       string `json:"error"`
			Description string `json:"error_description"`
		}
		json.NewDecoder(resp.Body).Decode(&errBody)
		if errBody.Description != "" {
			return nil, fmt.Errorf("authentication failed: %s", errBody.Description)
		}
		return nil, fmt.Errorf("authentication failed: %s", resp.Status)
	}

	var tokenResp TokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return nil, fmt.Errorf("decoding token response: %w", err)
	}
	return &tokenResp, nil
}

// AutoRefreshToken attempts to refresh the OIDC token if it is near expiry.
// Returns the current token if no refresh is needed or possible.
func AutoRefreshToken(cfg *ConfigFile) (string, error) {
	if cfg.Auth.OIDC == nil {
		return GetAuthToken(cfg), nil
	}

	oidc := cfg.Auth.OIDC
	if oidc.Token == "" {
		return "", nil
	}
	if oidc.ExpiresAt == "" {
		return oidc.Token, nil
	}

	expiresAt, err := time.Parse(time.RFC3339, oidc.ExpiresAt)
	if err != nil {
		return oidc.Token, nil
	}

	// Still valid with 30s headroom.
	if time.Now().Before(expiresAt.Add(-30 * time.Second)) {
		return oidc.Token, nil
	}

	if oidc.RefreshToken == "" || oidc.ClientID == "" || oidc.TokenURL == "" {
		return oidc.Token, nil
	}

	tokenResp, err := refreshTokenGrant(oidc.TokenURL, oidc.ClientID, oidc.ClientSecret, oidc.RefreshToken)
	if err != nil {
		return oidc.Token, nil
	}

	oidc.Token = tokenResp.AccessToken
	oidc.RefreshToken = tokenResp.RefreshToken
	oidc.ExpiresAt = time.Now().Add(time.Duration(tokenResp.ExpiresIn) * time.Second).Format(time.RFC3339)
	SaveConfig(cfg)

	return oidc.Token, nil
}

// InteractiveLoginFn builds a loginFn suitable for the TUI session from a loaded config.
// InteractiveLoginFn returns a LoginFunc for the REPL's 'login' command.
// It always performs a fresh OIDC or API-key grant — never silently reuses an
// existing token. AutoRefreshToken is used at session startup; explicit 'login'
// should always obtain a brand-new token and save it to the config file.
func InteractiveLoginFn(cfg *ConfigFile, cfgPath string) func() (string, error) {
	return func() (string, error) {
		// Reload config so we pick up any manual edits since session start.
		current, err := LoadConfigFromPath(cfgPath)
		if err != nil {
			current = cfg
		}

		if HasOIDCConfig(current) {
			return doInteractiveOIDCLogin(current, cfgPath)
		}

		if HasAPIKeyConfig(current) {
			return doInteractiveAPIKeyLogin(current, cfgPath)
		}

		return "", fmt.Errorf("no auth method configured in %s — add oidc or api_key section", cfgPath)
	}
}

func doInteractiveOIDCLogin(cfg *ConfigFile, cfgPath string) (string, error) {
	oidc := cfg.Auth.OIDC

	username := oidc.Username
	password := oidc.Password

	if username == "" {
		fmt.Print("Username: ")
		fmt.Scanln(&username)
	}
	if password == "" {
		fmt.Print("Password: ")
		pw, err := term.ReadPassword(int(os.Stdin.Fd()))
		fmt.Println()
		if err != nil {
			return "", fmt.Errorf("reading password: %w", err)
		}
		password = string(pw)
	}

	var tokenResp *TokenResponse
	var err error
	if username == "" && oidc.ClientSecret != "" {
		tokenResp, err = clientCredentialsGrant(oidc.TokenURL, oidc.ClientID, oidc.ClientSecret)
	} else {
		tokenResp, err = passwordGrant(oidc.TokenURL, oidc.ClientID, oidc.ClientSecret, username, password)
	}
	if err != nil {
		return "", err
	}

	oidc.Token = tokenResp.AccessToken
	oidc.RefreshToken = tokenResp.RefreshToken
	oidc.ExpiresAt = time.Now().Add(time.Duration(tokenResp.ExpiresIn) * time.Second).Format(time.RFC3339)
	SaveConfigToPath(cfg, cfgPath)

	return tokenResp.AccessToken, nil
}

func doInteractiveAPIKeyLogin(cfg *ConfigFile, cfgPath string) (string, error) {
	apiKey := cfg.Auth.APIKey

	req, err := http.NewRequest("GET", apiKey.AuthnURL, nil)
	if err != nil {
		return "", fmt.Errorf("building request: %w", err)
	}
	req.Header.Set("Authorization", "ApiKey "+apiKey.Key)
	req.Header.Set("Accept", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("requesting token from NGC: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("reading response: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("NGC token exchange failed (HTTP %d): %s", resp.StatusCode, string(body))
	}

	token := extractNGCToken(body)
	if token == "" {
		return "", fmt.Errorf("NGC response did not contain a token")
	}

	cfg.Auth.APIKey.Token = token
	SaveConfigToPath(cfg, cfgPath)
	return token, nil
}
