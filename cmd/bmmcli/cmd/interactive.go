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
	"path/filepath"
	"sort"
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
	selectedConfigPath, err := chooseInteractiveConfigFile()
	if err != nil {
		return err
	}
	if selectedConfigPath != "" {
		if err := reloadConfigWithFile(selectedConfigPath); err != nil {
			return fmt.Errorf("loading selected config %q: %w", selectedConfigPath, err)
		}
	}

	org := viper.GetString("api.org")
	if org == "" {
		return fmt.Errorf("org is required: set api.org in config or pass --org")
	}

	apiClient := newAPIClient()

	token, _ := getAuthToken()

	ctx := context.Background()
	if token != "" {
		ctx = context.WithValue(ctx, client.ContextAccessToken, token)
	}

	session := tui.NewSession(apiClient, ctx, org)
	session.Token = token

	// Wire up login -- auto-detect method
	if hasOIDCConfig() || hasAPIKeyConfig() {
		session.LoginFn = interactiveLogin
	}

	if token == "" {
		fmt.Fprintf(os.Stderr, "%s No auth token found. Type %s to authenticate.\n\n",
			tui.Yellow("Warning:"), tui.Bold("login"))
	}

	return tui.RunREPL(session)
}

func chooseInteractiveConfigFile() (string, error) {
	// Respect explicit --config usage.
	if cfgFile != "" {
		return "", nil
	}
	// Interactive selector requires a terminal.
	if !term.IsTerminal(int(os.Stdin.Fd())) || !term.IsTerminal(int(os.Stdout.Fd())) {
		return "", nil
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("finding home directory: %w", err)
	}

	configDir := filepath.Join(home, ".bmm")
	entries, err := os.ReadDir(configDir)
	if err != nil {
		if os.IsNotExist(err) {
			return "", nil
		}
		return "", fmt.Errorf("reading config directory %q: %w", configDir, err)
	}

	var candidatePaths []string
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
		candidatePaths = append(candidatePaths, filepath.Join(configDir, name))
	}

	if len(candidatePaths) <= 1 {
		return "", nil
	}

	sortConfigCandidates(candidatePaths)

	items := make([]tui.SelectItem, len(candidatePaths))
	for i, path := range candidatePaths {
		items[i] = tui.SelectItem{
			Label: displayPath(path, home),
			ID:    path,
		}
	}

	fmt.Println()
	selected, err := tui.Select("Choose config for this interactive session", items)
	if err != nil {
		return "", err
	}
	if selected == nil || strings.TrimSpace(selected.ID) == "" {
		return "", fmt.Errorf("no config selected")
	}
	fmt.Printf("Using config: %s\n\n", selected.Label)
	return selected.ID, nil
}

func sortConfigCandidates(paths []string) {
	sort.Slice(paths, func(i, j int) bool {
		left := filepath.Base(paths[i])
		right := filepath.Base(paths[j])
		leftDefault := left == "config.yaml" || left == "config.yml"
		rightDefault := right == "config.yaml" || right == "config.yml"
		if leftDefault != rightDefault {
			return leftDefault
		}
		return left < right
	})
}

func displayPath(path, home string) string {
	prefix := home + string(os.PathSeparator)
	if strings.HasPrefix(path, prefix) {
		return "~/" + strings.TrimPrefix(path, prefix)
	}
	return path
}

// interactiveLogin performs login using the configured auth method.
// Prefers OIDC if configured, otherwise uses API key.
func interactiveLogin() (string, error) {
	if hasOIDCConfig() {
		return interactiveOIDCLogin()
	}
	if hasAPIKeyConfig() {
		return interactiveAPIKeyLogin()
	}
	return "", fmt.Errorf("no auth method configured")
}

func interactiveOIDCLogin() (string, error) {
	tokenURL := viper.GetString("auth.oidc.token_url")
	clientID := viper.GetString("auth.oidc.client_id")
	clientSecret := viper.GetString("auth.oidc.client_secret")

	username := viper.GetString("auth.oidc.username")
	password := viper.GetString("auth.oidc.password")

	if username == "" {
		fmt.Print("Username: ")
		fmt.Scanln(&username)
	}

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

	var tokenResp oidcTokenResponse
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return "", fmt.Errorf("parsing token response: %w", err)
	}

	viper.Set("auth.oidc.token", tokenResp.AccessToken)
	if err := saveToken("oidc", tokenResp.AccessToken); err != nil {
		fmt.Fprintf(os.Stderr, "Warning: could not save token to config: %v\n", err)
	}

	return tokenResp.AccessToken, nil
}

func interactiveAPIKeyLogin() (string, error) {
	authnURL := viper.GetString("auth.api_key.authn_url")
	apiKey := viper.GetString("auth.api_key.key")

	req, err := http.NewRequest("GET", authnURL, nil)
	if err != nil {
		return "", fmt.Errorf("building request: %w", err)
	}
	req.Header.Set("Authorization", "ApiKey "+apiKey)
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
		return "", fmt.Errorf("NGC token exchange failed (HTTP %d): %s\nURL: %s", resp.StatusCode, string(body), authnURL)
	}

	token := extractToken(body)
	if token == "" {
		return "", fmt.Errorf("NGC authn response did not contain a token: %s", string(body))
	}

	viper.Set("auth.api_key.token", token)
	if err := saveToken("api_key", token); err != nil {
		fmt.Fprintf(os.Stderr, "Warning: could not save token to config: %v\n", err)
	}

	return token, nil
}
