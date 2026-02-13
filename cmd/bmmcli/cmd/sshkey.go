// SPDX-FileCopyrightText: Copyright (c) 2026 NVIDIA CORPORATION & AFFILIATES. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"text/tabwriter"
	"time"

	"github.com/nvidia/bare-metal-manager-rest/client"
	"github.com/nvidia/bare-metal-manager-rest/cmd/bmmcli/internal/pagination"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var sshKeyCmd = &cobra.Command{
	Use:   "ssh-key",
	Short: "SSH key operations",
}

var sshKeyListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all SSH keys",
	RunE:  runSSHKeyList,
}

var sshKeyGetCmd = &cobra.Command{
	Use:   "get <ssh-key-id>",
	Short: "Get SSH key details",
	Args:  cobra.ExactArgs(1),
	RunE:  runSSHKeyGet,
}

var sshKeyCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create an SSH key",
	RunE:  runSSHKeyCreate,
}

var sshKeyUpdateCmd = &cobra.Command{
	Use:   "update <ssh-key-id>",
	Short: "Update an SSH key",
	Args:  cobra.ExactArgs(1),
	RunE:  runSSHKeyUpdate,
}

var sshKeyDeleteCmd = &cobra.Command{
	Use:   "delete <ssh-key-id>",
	Short: "Delete an SSH key",
	Args:  cobra.ExactArgs(1),
	RunE:  runSSHKeyDelete,
}

func init() {
	sshKeyListCmd.Flags().Bool("json", false, "output raw JSON")

	sshKeyCreateCmd.Flags().String("name", "", "name for the SSH key (required)")
	sshKeyCreateCmd.Flags().String("public-key", "", "public key content (required)")
	sshKeyCreateCmd.MarkFlagRequired("name")
	sshKeyCreateCmd.MarkFlagRequired("public-key")

	sshKeyUpdateCmd.Flags().String("name", "", "new name")

	rootCmd.AddCommand(sshKeyCmd)
	sshKeyCmd.AddCommand(sshKeyListCmd)
	sshKeyCmd.AddCommand(sshKeyGetCmd)
	sshKeyCmd.AddCommand(sshKeyCreateCmd)
	sshKeyCmd.AddCommand(sshKeyUpdateCmd)
	sshKeyCmd.AddCommand(sshKeyDeleteCmd)
}

func runSSHKeyList(cmd *cobra.Command, args []string) error {
	org := viper.GetString("api.org")
	if org == "" {
		return fmt.Errorf("org is required: set api.org in config or pass --org")
	}

	apiClient := newAPIClient()
	ctx, err := apiContext()
	if err != nil {
		return err
	}

	keys, resp, err := pagination.FetchAllPages(func(pageNumber, pageSize int32) ([]client.SshKey, *http.Response, error) {
		return apiClient.SSHKeyAPI.GetAllSshKey(ctx, org).PageNumber(pageNumber).PageSize(pageSize).Execute()
	})
	if err != nil {
		if resp != nil {
			body := tryReadBody(resp.Body)
			return fmt.Errorf("listing SSH keys (HTTP %d): %v\n%s", resp.StatusCode, err, body)
		}
		return fmt.Errorf("listing SSH keys: %v", err)
	}
	pagination.PrintSummary(cmd.ErrOrStderr(), resp, len(keys))

	jsonFlag, _ := cmd.Flags().GetBool("json")
	outputFlag, _ := cmd.Root().PersistentFlags().GetString("output")
	switch {
	case jsonFlag || outputFlag == "json":
		return printJSON(os.Stdout, keys)
	case outputFlag == "yaml":
		return printYAML(os.Stdout, keys)
	default:
		return printSSHKeyTable(os.Stdout, keys)
	}
}

func runSSHKeyGet(cmd *cobra.Command, args []string) error {
	org := viper.GetString("api.org")
	if org == "" {
		return fmt.Errorf("org is required: set api.org in config or pass --org")
	}

	apiClient := newAPIClient()
	ctx, err := apiContext()
	if err != nil {
		return err
	}

	key, resp, err := apiClient.SSHKeyAPI.GetSshKey(ctx, org, args[0]).Execute()
	if err != nil {
		if resp != nil {
			body := tryReadBody(resp.Body)
			return fmt.Errorf("getting SSH key (HTTP %d): %v\n%s", resp.StatusCode, err, body)
		}
		return fmt.Errorf("getting SSH key: %v", err)
	}

	outputFlag, _ := cmd.Root().PersistentFlags().GetString("output")
	switch outputFlag {
	case "yaml":
		return printYAML(os.Stdout, key)
	default:
		return printJSON(os.Stdout, key)
	}
}

func runSSHKeyCreate(cmd *cobra.Command, args []string) error {
	org := viper.GetString("api.org")
	if org == "" {
		return fmt.Errorf("org is required: set api.org in config or pass --org")
	}

	name, _ := cmd.Flags().GetString("name")
	publicKey, _ := cmd.Flags().GetString("public-key")

	apiClient := newAPIClient()
	ctx, err := apiContext()
	if err != nil {
		return err
	}

	req := client.NewSshKeyCreateRequest(name, publicKey)

	key, resp, err := apiClient.SSHKeyAPI.CreateSshKey(ctx, org).SshKeyCreateRequest(*req).Execute()
	if err != nil {
		if resp != nil {
			body := tryReadBody(resp.Body)
			return fmt.Errorf("creating SSH key (HTTP %d): %v\n%s", resp.StatusCode, err, body)
		}
		return fmt.Errorf("creating SSH key: %v", err)
	}

	fmt.Fprintf(os.Stderr, "SSH key created: %s (%s)\n", ptrStr(key.Name), ptrStr(key.Id))
	return printJSON(os.Stdout, key)
}

func runSSHKeyUpdate(cmd *cobra.Command, args []string) error {
	org := viper.GetString("api.org")
	if org == "" {
		return fmt.Errorf("org is required: set api.org in config or pass --org")
	}

	apiClient := newAPIClient()
	ctx, err := apiContext()
	if err != nil {
		return err
	}

	req := client.NewSshKeyUpdateRequest()
	if cmd.Flags().Changed("name") {
		name, _ := cmd.Flags().GetString("name")
		req.SetName(name)
	}

	key, resp, err := apiClient.SSHKeyAPI.UpdateSshKey(ctx, org, args[0]).SshKeyUpdateRequest(*req).Execute()
	if err != nil {
		if resp != nil {
			body := tryReadBody(resp.Body)
			return fmt.Errorf("updating SSH key (HTTP %d): %v\n%s", resp.StatusCode, err, body)
		}
		return fmt.Errorf("updating SSH key: %v", err)
	}

	fmt.Fprintf(os.Stderr, "SSH key updated: %s (%s)\n", ptrStr(key.Name), ptrStr(key.Id))
	return printJSON(os.Stdout, key)
}

func runSSHKeyDelete(cmd *cobra.Command, args []string) error {
	org := viper.GetString("api.org")
	if org == "" {
		return fmt.Errorf("org is required: set api.org in config or pass --org")
	}

	apiClient := newAPIClient()
	ctx, err := apiContext()
	if err != nil {
		return err
	}

	resp, err := apiClient.SSHKeyAPI.DeleteSshKey(ctx, org, args[0]).Execute()
	if err != nil {
		if resp != nil {
			body := tryReadBody(resp.Body)
			return fmt.Errorf("deleting SSH key (HTTP %d): %v\n%s", resp.StatusCode, err, body)
		}
		return fmt.Errorf("deleting SSH key: %v", err)
	}

	fmt.Fprintf(os.Stderr, "SSH key deleted: %s\n", args[0])
	return nil
}

func printSSHKeyTable(w io.Writer, keys []client.SshKey) error {
	tw := tabwriter.NewWriter(w, 0, 0, 3, ' ', 0)
	fmt.Fprintln(tw, "NAME\tFINGERPRINT\tAGE\tID")

	for _, k := range keys {
		name := ptrStr(k.Name)
		id := ptrStr(k.Id)
		fingerprint := ptrStr(k.Fingerprint)
		age := "<unknown>"
		if k.Created != nil {
			age = formatAge(time.Since(*k.Created))
		}

		fmt.Fprintf(tw, "%s\t%s\t%s\t%s\n", name, fingerprint, age, id)
	}

	return tw.Flush()
}
