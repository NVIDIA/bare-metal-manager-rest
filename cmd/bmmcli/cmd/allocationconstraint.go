// SPDX-FileCopyrightText: Copyright (c) 2026 NVIDIA CORPORATION & AFFILIATES. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	"fmt"
	"os"

	"github.com/nvidia/bare-metal-manager-rest/client"
	"github.com/nvidia/bare-metal-manager-rest/cmd/bmmcli/internal/pagination"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var allocationConstraintCmd = &cobra.Command{
	Use:   "allocation-constraint",
	Short: "Allocation constraint operations (nested under allocation)",
}

var allocationConstraintListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all allocation constraints for an allocation",
	RunE:  runAllocationConstraintList,
}

var allocationConstraintGetCmd = &cobra.Command{
	Use:   "get <constraint-id>",
	Short: "Get allocation constraint details",
	Args:  cobra.ExactArgs(1),
	RunE:  runAllocationConstraintGet,
}

var allocationConstraintCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create an allocation constraint",
	RunE:  runAllocationConstraintCreate,
}

var allocationConstraintUpdateCmd = &cobra.Command{
	Use:   "update <constraint-id>",
	Short: "Update an allocation constraint",
	Args:  cobra.ExactArgs(1),
	RunE:  runAllocationConstraintUpdate,
}

var allocationConstraintDeleteCmd = &cobra.Command{
	Use:   "delete <constraint-id>",
	Short: "Delete an allocation constraint",
	Args:  cobra.ExactArgs(1),
	RunE:  runAllocationConstraintDelete,
}

func init() {
	allocationConstraintListCmd.Flags().Bool("json", false, "output raw JSON")
	allocationConstraintListCmd.Flags().String("allocation-id", "", "allocation ID (required)")
	allocationConstraintListCmd.MarkFlagRequired("allocation-id")

	allocationConstraintGetCmd.Flags().String("allocation-id", "", "allocation ID (required)")
	allocationConstraintGetCmd.MarkFlagRequired("allocation-id")

	allocationConstraintCreateCmd.Flags().String("allocation-id", "", "allocation ID (required)")
	allocationConstraintCreateCmd.Flags().String("resource-type-id", "", "resource type ID (required)")
	allocationConstraintCreateCmd.Flags().String("constraint-type", "", "constraint type (required)")
	allocationConstraintCreateCmd.Flags().Int32("constraint-value", 0, "constraint value (required)")
	allocationConstraintCreateCmd.MarkFlagRequired("allocation-id")
	allocationConstraintCreateCmd.MarkFlagRequired("resource-type-id")
	allocationConstraintCreateCmd.MarkFlagRequired("constraint-type")
	allocationConstraintCreateCmd.MarkFlagRequired("constraint-value")

	allocationConstraintUpdateCmd.Flags().String("allocation-id", "", "allocation ID (required)")
	allocationConstraintUpdateCmd.Flags().Int32("constraint-value", 0, "new constraint value (required)")
	allocationConstraintUpdateCmd.MarkFlagRequired("allocation-id")
	allocationConstraintUpdateCmd.MarkFlagRequired("constraint-value")

	allocationConstraintDeleteCmd.Flags().String("allocation-id", "", "allocation ID (required)")
	allocationConstraintDeleteCmd.MarkFlagRequired("allocation-id")

	rootCmd.AddCommand(allocationConstraintCmd)
	allocationConstraintCmd.AddCommand(allocationConstraintListCmd)
	allocationConstraintCmd.AddCommand(allocationConstraintGetCmd)
	allocationConstraintCmd.AddCommand(allocationConstraintCreateCmd)
	allocationConstraintCmd.AddCommand(allocationConstraintUpdateCmd)
	allocationConstraintCmd.AddCommand(allocationConstraintDeleteCmd)
}

func runAllocationConstraintList(cmd *cobra.Command, args []string) error {
	org := viper.GetString("api.org")
	if org == "" {
		return fmt.Errorf("org is required: set api.org in config or pass --org")
	}

	allocID, _ := cmd.Flags().GetString("allocation-id")

	apiClient := newAPIClient()
	ctx, err := apiContext()
	if err != nil {
		return err
	}

	constraints, resp, err := apiClient.AllocationAPI.GetAllAllocationConstraint(ctx, org, allocID).Execute()
	if err != nil {
		if resp != nil {
			body := tryReadBody(resp.Body)
			return fmt.Errorf("listing allocation constraints (HTTP %d): %v\n%s", resp.StatusCode, err, body)
		}
		return fmt.Errorf("listing allocation constraints: %v", err)
	}
	pagination.PrintSummary(cmd.ErrOrStderr(), resp, len(constraints))

	jsonFlag, _ := cmd.Flags().GetBool("json")
	outputFlag, _ := cmd.Root().PersistentFlags().GetString("output")
	switch {
	case jsonFlag || outputFlag == "json":
		return printJSON(os.Stdout, constraints)
	case outputFlag == "yaml":
		return printYAML(os.Stdout, constraints)
	default:
		return printJSON(os.Stdout, constraints)
	}
}

func runAllocationConstraintGet(cmd *cobra.Command, args []string) error {
	org := viper.GetString("api.org")
	if org == "" {
		return fmt.Errorf("org is required: set api.org in config or pass --org")
	}

	allocID, _ := cmd.Flags().GetString("allocation-id")

	apiClient := newAPIClient()
	ctx, err := apiContext()
	if err != nil {
		return err
	}

	constraint, resp, err := apiClient.AllocationAPI.GetAllocationConstraint(ctx, org, allocID, args[0]).Execute()
	if err != nil {
		if resp != nil {
			body := tryReadBody(resp.Body)
			return fmt.Errorf("getting allocation constraint (HTTP %d): %v\n%s", resp.StatusCode, err, body)
		}
		return fmt.Errorf("getting allocation constraint: %v", err)
	}

	outputFlag, _ := cmd.Root().PersistentFlags().GetString("output")
	switch outputFlag {
	case "yaml":
		return printYAML(os.Stdout, constraint)
	default:
		return printJSON(os.Stdout, constraint)
	}
}

func runAllocationConstraintCreate(cmd *cobra.Command, args []string) error {
	org := viper.GetString("api.org")
	if org == "" {
		return fmt.Errorf("org is required: set api.org in config or pass --org")
	}

	allocID, _ := cmd.Flags().GetString("allocation-id")
	resourceTypeID, _ := cmd.Flags().GetString("resource-type-id")
	constraintType, _ := cmd.Flags().GetString("constraint-type")
	constraintValue, _ := cmd.Flags().GetInt32("constraint-value")

	apiClient := newAPIClient()
	ctx, err := apiContext()
	if err != nil {
		return err
	}

	req := client.NewAllocationConstraintCreateRequest(resourceTypeID, constraintType, constraintValue)

	constraint, resp, err := apiClient.AllocationAPI.CreateAllocationConstraint(ctx, org, allocID).AllocationConstraintCreateRequest(*req).Execute()
	if err != nil {
		if resp != nil {
			body := tryReadBody(resp.Body)
			return fmt.Errorf("creating allocation constraint (HTTP %d): %v\n%s", resp.StatusCode, err, body)
		}
		return fmt.Errorf("creating allocation constraint: %v", err)
	}

	fmt.Fprintf(os.Stderr, "Allocation constraint created: %s\n", ptrStr(constraint.Id))
	return printJSON(os.Stdout, constraint)
}

func runAllocationConstraintUpdate(cmd *cobra.Command, args []string) error {
	org := viper.GetString("api.org")
	if org == "" {
		return fmt.Errorf("org is required: set api.org in config or pass --org")
	}

	allocID, _ := cmd.Flags().GetString("allocation-id")
	constraintValue, _ := cmd.Flags().GetInt32("constraint-value")

	apiClient := newAPIClient()
	ctx, err := apiContext()
	if err != nil {
		return err
	}

	req := client.NewAllocationConstraintUpdateRequest(constraintValue)

	constraint, resp, err := apiClient.AllocationAPI.UpdateAllocationConstraint(ctx, org, allocID, args[0]).AllocationConstraintUpdateRequest(*req).Execute()
	if err != nil {
		if resp != nil {
			body := tryReadBody(resp.Body)
			return fmt.Errorf("updating allocation constraint (HTTP %d): %v\n%s", resp.StatusCode, err, body)
		}
		return fmt.Errorf("updating allocation constraint: %v", err)
	}

	fmt.Fprintf(os.Stderr, "Allocation constraint updated: %s\n", ptrStr(constraint.Id))
	return printJSON(os.Stdout, constraint)
}

func runAllocationConstraintDelete(cmd *cobra.Command, args []string) error {
	org := viper.GetString("api.org")
	if org == "" {
		return fmt.Errorf("org is required: set api.org in config or pass --org")
	}

	allocID, _ := cmd.Flags().GetString("allocation-id")

	apiClient := newAPIClient()
	ctx, err := apiContext()
	if err != nil {
		return err
	}

	resp, err := apiClient.AllocationAPI.DeleteAllocationConstraint(ctx, org, allocID, args[0]).Execute()
	if err != nil {
		if resp != nil {
			body := tryReadBody(resp.Body)
			return fmt.Errorf("deleting allocation constraint (HTTP %d): %v\n%s", resp.StatusCode, err, body)
		}
		return fmt.Errorf("deleting allocation constraint: %v", err)
	}

	fmt.Fprintf(os.Stderr, "Allocation constraint deleted: %s\n", args[0])
	return nil
}
