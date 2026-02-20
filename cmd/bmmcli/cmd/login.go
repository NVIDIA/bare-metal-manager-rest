// SPDX-FileCopyrightText: Copyright (c) 2026 NVIDIA CORPORATION & AFFILIATES. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"syscall"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/term"
)

// oidcTokenResponse represents an OIDC token endpoint response
type oidcTokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
}

// ngcTokenResponse represents the NGC authn token exchange response
type ngcTokenResponse struct {
	Token string `json:"token"`
}

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Authenticate with the Carbide API",
	Long: `Authenticate and save a session token to config.

Login auto-detects the configured auth method:
  - If auth.token is configured, login is not required.
  - If auth.oidc is configured, uses OIDC password flow (default).
  - If auth.api_key is configured, exchanges the NGC API key.
  - If both are configured, defaults to OIDC. Use --api-key to force NGC.

Examples:
  bmm login                  # auto-detect (prefers OIDC)
  bmm login --api-key        # force NGC API key exchange`,
	RunE: runLogin,
}

func init() {
	loginCmd.Flags().String("username", "", "username for OIDC authentication")
	loginCmd.Flags().String("password", "", "password for OIDC authentication")
	loginCmd.Flags().Bool("api-key", false, "force NGC API key exchange instead of OIDC")
	rootCmd.AddCommand(loginCmd)
}

func runLogin(cmd *cobra.Command, args []string) error {
	forceAPIKey, _ := cmd.Flags().GetBool("api-key")

	if forceAPIKey {
		if !hasAPIKeyConfig() {
			return fmt.Errorf("auth.api_key is not configured in %s\n\n"+
				"Add:\n"+
				"  auth:\n"+
				"    api_key:\n"+
				"      authn_url: https://authn.nvidia.com/token\n"+
				"      key: nvapi-xxxx", configFilePath())
		}
		return loginWithAPIKey()
	}

	if hasDirectTokenConfig() && !hasOIDCConfig() && !hasAPIKeyConfig() {
		fmt.Fprintf(os.Stderr, "Direct token is already configured at auth.token in %s\n", configFilePath())
		return nil
	}

	// Auto-detect: prefer OIDC, fall back to API key
	if hasOIDCConfig() {
		return loginWithOIDC(cmd)
	}
	if hasAPIKeyConfig() {
		return loginWithAPIKey()
	}

	return fmt.Errorf(authGuidance, configFilePath())
}

func loginWithAPIKey() error {
	authnURL := viper.GetString("auth.api_key.authn_url")
	apiKey := viper.GetString("auth.api_key.key")

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
		return fmt.Errorf("NGC token exchange failed (HTTP %d): %s\nURL: %s", resp.StatusCode, string(body), authnURL)
	}

	token := extractToken(body)
	if token == "" {
		return fmt.Errorf("NGC authn response did not contain a token: %s", string(body))
	}

	viper.Set("auth.api_key.token", token)
	if err := saveToken("api_key", token); err != nil {
		return fmt.Errorf("saving token: %w", err)
	}

	fmt.Fprintf(os.Stderr, "Login successful (NGC API key). Token saved to %s\n", configFilePath())
	return nil
}

// extractToken tries multiple JSON response formats to find the token.
func extractToken(body []byte) string {
	// Try {"token": "..."}
	var ngcResp ngcTokenResponse
	if err := json.Unmarshal(body, &ngcResp); err == nil && ngcResp.Token != "" {
		return ngcResp.Token
	}
	// Try {"access_token": "..."} (OIDC-style)
	var oidcResp oidcTokenResponse
	if err := json.Unmarshal(body, &oidcResp); err == nil && oidcResp.AccessToken != "" {
		return oidcResp.AccessToken
	}
	return ""
}

func loginWithOIDC(cmd *cobra.Command) error {
	username, _ := cmd.Flags().GetString("username")
	password, _ := cmd.Flags().GetString("password")

	if username == "" {
		username = viper.GetString("auth.oidc.username")
	}
	if password == "" {
		password = viper.GetString("auth.oidc.password")
	}

	if username == "" {
		fmt.Print("Username: ")
		fmt.Scanln(&username)
	}

	if password == "" {
		fmt.Print("Password: ")
		passwordBytes, err := term.ReadPassword(int(syscall.Stdin))
		if err != nil {
			return fmt.Errorf("reading password: %w", err)
		}
		fmt.Println()
		password = string(passwordBytes)
	}

	tokenURL := viper.GetString("auth.oidc.token_url")
	clientID := viper.GetString("auth.oidc.client_id")
	clientSecret := viper.GetString("auth.oidc.client_secret")

	formData := url.Values{
		"grant_type":    {"password"},
		"client_id":     {clientID},
		"client_secret": {clientSecret},
		"username":      {username},
		"password":      {password},
	}

	resp, err := http.Post(tokenURL, "application/x-www-form-urlencoded", strings.NewReader(formData.Encode()))
	if err != nil {
		return fmt.Errorf("requesting token: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("reading response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("authentication failed (HTTP %d): %s", resp.StatusCode, string(body))
	}

	var tokenResp oidcTokenResponse
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return fmt.Errorf("parsing token response: %w", err)
	}

	viper.Set("auth.oidc.token", tokenResp.AccessToken)
	if err := saveToken("oidc", tokenResp.AccessToken); err != nil {
		return fmt.Errorf("saving token: %w", err)
	}

	fmt.Fprintf(os.Stderr, "Login successful (OIDC). Token saved to %s\n", configFilePath())
	return nil
}
