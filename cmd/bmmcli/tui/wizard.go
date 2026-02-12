// SPDX-FileCopyrightText: Copyright (c) 2026 NVIDIA CORPORATION & AFFILIATES. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package tui

import (
	"fmt"

	"github.com/nvidia/bare-metal-manager-rest/client"
)

// RunInstanceWizard walks the user through creating an instance step by step
func RunInstanceWizard(s *Session) error {
	fmt.Println(Bold("\n-- Instance Create Wizard --\n"))

	// Step 1: Select site
	site, err := s.Resolver.Resolve(s.Ctx, "site", "Site")
	if err != nil {
		return err
	}

	// Step 2: Select VPC (filtered by site)
	vpc, err := s.Resolver.ResolveFiltered(s.Ctx, "vpc", "VPC", "siteId", site.ID)
	if err != nil {
		return err
	}

	// Step 3: Select subnet (filtered by VPC) -- needed for the network interface
	subnet, err := s.Resolver.ResolveFiltered(s.Ctx, "subnet", "Subnet", "vpcId", vpc.ID)
	if err != nil {
		fmt.Println(Yellow("No subnets found. You may need to create one first with: subnet create"))
		return err
	}

	// Step 4: Select instance type (fetched by site ID directly from API)
	itItems, err := s.fetchInstanceTypesBySite(s.Ctx, site.ID)
	if err != nil {
		return fmt.Errorf("fetching instance types for site: %w", err)
	}
	instanceType, err := s.Resolver.SelectFromItems("Instance Type", itItems)
	if err != nil {
		return err
	}

	// Step 5: Select OS (optional)
	var osID *string
	osList, err := s.Resolver.Fetch(s.Ctx, "operating-system")
	if err == nil && len(osList) > 0 {
		os, err := s.Resolver.Resolve(s.Ctx, "operating-system", "Operating System")
		if err == nil {
			osID = &os.ID
		}
	}

	// Step 6: Select SSH key group (optional)
	var sshKeyGroupIDs []string
	sshGroups, err := s.Resolver.Fetch(s.Ctx, "ssh-key-group")
	if err == nil && len(sshGroups) > 0 {
		sshGroup, err := s.Resolver.Resolve(s.Ctx, "ssh-key-group", "SSH Key Group")
		if err == nil {
			sshKeyGroupIDs = []string{sshGroup.ID}
		}
	}

	// Step 7: Instance name
	name, err := PromptText("Instance name", true)
	if err != nil {
		return err
	}

	// Step 8: Get tenant ID
	tenantID, err := s.getTenantID(s.Ctx)
	if err != nil {
		return err
	}

	// Summary
	fmt.Println(Bold("\n-- Summary --"))
	fmt.Printf("  Site:          %s\n", site.Name)
	fmt.Printf("  VPC:           %s\n", vpc.Name)
	fmt.Printf("  Subnet:        %s\n", subnet.Name)
	fmt.Printf("  Instance Type: %s\n", instanceType.Name)
	if osID != nil {
		osName := s.Resolver.ResolveID("operating-system", *osID)
		fmt.Printf("  OS:            %s\n", osName)
	} else {
		fmt.Printf("  OS:            %s\n", Dim("(iPXE fallback)"))
	}
	if len(sshKeyGroupIDs) > 0 {
		sshName := s.Resolver.ResolveID("ssh-key-group", sshKeyGroupIDs[0])
		fmt.Printf("  SSH Key Group: %s\n", sshName)
	}
	fmt.Printf("  Name:          %s\n", name)
	fmt.Println()

	ok, err := PromptConfirm("Create instance?")
	if err != nil || !ok {
		fmt.Println("Cancelled.")
		return nil
	}

	// Build network interface
	subnetID := subnet.ID
	iface := client.InterfaceCreateRequest{
		SubnetId: &subnetID,
	}

	// Build request
	req := client.NewInstanceCreateRequest(name, tenantID, vpc.ID, []client.InterfaceCreateRequest{iface})
	itID := instanceType.ID
	req.InstanceTypeId = &itID
	if osID != nil {
		req.SetOperatingSystemId(*osID)
	} else {
		// No OS selected; use a minimal iPXE script as required by the API
		req.SetIpxeScript("#!ipxe\necho No OS configured")
	}
	if len(sshKeyGroupIDs) > 0 {
		req.SetSshKeyGroupIds(sshKeyGroupIDs)
	}

	fmt.Print("Creating instance...")
	instance, _, err := s.Client.InstanceAPI.CreateInstance(s.Ctx, s.Org).InstanceCreateRequest(*req).Execute()
	if err != nil {
		fmt.Println()
		return fmt.Errorf("creating instance: %w", err)
	}

	s.Cache.Invalidate("instance")
	fmt.Printf("\r%s Instance created: %s (%s)\n", Green("OK"), ptrStr(instance.Name), ptrStr(instance.Id))
	return nil
}
