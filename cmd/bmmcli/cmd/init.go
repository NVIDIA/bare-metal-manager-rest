// SPDX-FileCopyrightText: Copyright (c) 2026 NVIDIA CORPORATION & AFFILIATES. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Create a sample config file",
	Long:  "Write a sample config file to ~/.bmm/config.yaml (or the path specified by --config). Will not overwrite an existing file unless --force is used.",
	RunE:  runInit,
}

func init() {
	initCmd.Flags().Bool("force", false, "overwrite existing config file")
	rootCmd.AddCommand(initCmd)
}

func runInit(cmd *cobra.Command, args []string) error {
	force, _ := cmd.Flags().GetBool("force")
	cfgPath := configFilePath()

	if !force {
		if _, err := os.Stat(cfgPath); err == nil {
			return fmt.Errorf("config file already exists at %s (use --force to overwrite)", cfgPath)
		}
	}

	sampleCfg := Config{
		API: APIConfig{
			Base: "http://localhost:8388",
			Org:  "test-org",
			Name: "carbide",
		},
		Auth: AuthConfig{
			OIDC: &OIDCAuth{
				TokenURL:     "http://localhost:8080/realms/carbide-dev/protocol/openid-connect/token",
				ClientID:     "carbide-api",
				ClientSecret: "carbide-local-secret",
				Username:     "admin@example.com",
				Password:     "adminpassword",
			},
			APIKey: &APIKeyAuth{
				AuthnURL: "https://authn.nvidia.com/token",
				Key:      "",
			},
		},
	}

	dir := filepath.Dir(cfgPath)
	if err := os.MkdirAll(dir, 0700); err != nil {
		return fmt.Errorf("creating config directory: %w", err)
	}

	data, err := yaml.Marshal(sampleCfg)
	if err != nil {
		return fmt.Errorf("marshaling config: %w", err)
	}

	output := configHeader + string(data)

	if err := os.WriteFile(cfgPath, []byte(output), 0600); err != nil {
		return fmt.Errorf("writing config file: %w", err)
	}

	fmt.Printf("Config written to %s\n", cfgPath)
	return nil
}
