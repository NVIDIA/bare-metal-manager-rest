// SPDX-FileCopyrightText: Copyright (c) 2026 NVIDIA CORPORATION & AFFILIATES. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"syscall"

	"github.com/nvidia/bare-metal-manager-rest/client"
	"github.com/nvidia/bare-metal-manager-rest/cmd/bmmcli/tui"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/term"
)

var interactiveCmd = &cobra.Command{
	Use:     "interactive",
	Aliases: []string{"i"},
	Short:   "Open interactive TUI mode",
	Long:    "Start an interactive session with inline autocomplete, arrow-key menus, name resolution, and guided wizards.",
	RunE:    runInteractive,
}

func init() {
	rootCmd.AddCommand(interactiveCmd)
}

func runInteractive(cmd *cobra.Command, args []string) error {
	org := viper.GetString("api.org")
	if org == "" {
		return fmt.Errorf("org is required: set api.org in config or pass --org")
	}

	// Build API client
	cfg := client.NewConfiguration()
	cfg.Servers = client.ServerConfigurations{
		{
			URL:         viper.GetString("api.base"),
			Description: "Configured server",
		},
	}
	apiClient := client.NewAPIClient(cfg)

	// Try to get existing token, but don't fail if missing
	token, _ := getAuthToken()

	ctx := context.Background()
	if token != "" {
		ctx = context.WithValue(ctx, client.ContextAccessToken, token)
	}

	session := tui.NewSession(apiClient, ctx, org)

	// Wire up the login function for in-session login
	if hasAuthProviderConfig() {
		session.LoginFn = interactiveLogin
	}

	// If no token, hint to the user
	if token == "" {
		fmt.Fprintf(os.Stderr, "%s No auth token found. Type %s to authenticate.\n\n",
			tui.Yellow("Warning:"), tui.Bold("login"))
	}

	return tui.RunREPL(session)
}

// interactiveLogin performs the OIDC login flow and returns the new token
func interactiveLogin() (string, error) {
	tokenURL := viper.GetString("auth.oidc.token_url")
	clientID := viper.GetString("auth.oidc.client_id")
	clientSecret := viper.GetString("auth.oidc.client_secret")

	username := viper.GetString("auth.oidc.username")
	password := viper.GetString("auth.oidc.password")

	// Prompt for username if not in config
	if username == "" {
		fmt.Print("Username: ")
		fmt.Scanln(&username)
	}

	// Prompt for password if not in config
	if password == "" {
		fmt.Print("Password: ")
		passwordBytes, err := term.ReadPassword(int(syscall.Stdin))
		if err != nil {
			return "", fmt.Errorf("reading password: %w", err)
		}
		fmt.Println()
		password = string(passwordBytes)
	}

	formData := url.Values{
		"grant_type":    {"password"},
		"client_id":     {clientID},
		"client_secret": {clientSecret},
		"username":      {username},
		"password":      {password},
	}

	resp, err := http.Post(tokenURL, "application/x-www-form-urlencoded", strings.NewReader(formData.Encode()))
	if err != nil {
		return "", fmt.Errorf("requesting token: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("reading response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("authentication failed (HTTP %d): %s", resp.StatusCode, string(body))
	}

	var tokenResp tokenResponse
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return "", fmt.Errorf("parsing token response: %w", err)
	}

	// Save tokens to config
	viper.Set("auth.token", tokenResp.AccessToken)
	viper.Set("auth.refresh_token", tokenResp.RefreshToken)
	if err := saveConfig(); err != nil {
		fmt.Fprintf(os.Stderr, "Warning: could not save token to config: %v\n", err)
	}

	return tokenResp.AccessToken, nil
}
