// SPDX-FileCopyrightText: Copyright (c) 2026 NVIDIA CORPORATION & AFFILIATES. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package interactive

import (
	"fmt"

	"github.com/nvidia/bare-metal-manager-rest/sdk/standard"
)

// RunInstanceWizard walks the user through creating an instance step by step.
func RunInstanceWizard(s *Session) error {
	fmt.Println(Bold("\n-- Instance Create Wizard --\n"))

	site, err := s.Resolver.Resolve(s.Ctx, "site", "Site")
	if err != nil {
		return err
	}

	vpc, err := s.Resolver.ResolveFiltered(s.Ctx, "vpc", "VPC", "siteId", site.ID)
	if err != nil {
		return err
	}

	subnet, err := s.Resolver.ResolveFiltered(s.Ctx, "subnet", "Subnet", "vpcId", vpc.ID)
	if err != nil {
		fmt.Println(Yellow("No subnets found. You may need to create one first with: subnet create"))
		return err
	}

	itItems, err := s.fetchInstanceTypesBySite(s.Ctx, site.ID)
	if err != nil {
		return fmt.Errorf("fetching instance types for site: %w", err)
	}
	instanceType, err := s.Resolver.SelectFromItems("Instance Type", itItems)
	if err != nil {
		return err
	}

	var osID *string
	osList, err := s.Resolver.Fetch(s.Ctx, "operating-system")
	if err == nil && len(osList) > 0 {
		os, err := s.Resolver.Resolve(s.Ctx, "operating-system", "Operating System")
		if err == nil {
			osID = &os.ID
		}
	}

	var sshKeyGroupIDs []string
	sshGroups, err := s.Resolver.Fetch(s.Ctx, "ssh-key-group")
	if err == nil && len(sshGroups) > 0 {
		sshGroup, err := s.Resolver.Resolve(s.Ctx, "ssh-key-group", "SSH Key Group")
		if err == nil {
			sshKeyGroupIDs = []string{sshGroup.ID}
		}
	}

	name, err := PromptText("Instance name", true)
	if err != nil {
		return err
	}

	tenantID, err := s.getTenantID(s.Ctx)
	if err != nil {
		return err
	}

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

	subnetID := subnet.ID
	iface := standard.InterfaceCreateRequest{
		SubnetId: &subnetID,
	}

	req := standard.NewInstanceCreateRequest(name, tenantID, vpc.ID, []standard.InterfaceCreateRequest{iface})
	itID := instanceType.ID
	req.InstanceTypeId = &itID
	if osID != nil {
		req.SetOperatingSystemId(*osID)
	} else {
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
