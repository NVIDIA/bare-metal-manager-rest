// SPDX-FileCopyrightText: Copyright (c) 2026 NVIDIA CORPORATION & AFFILIATES. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

// Quick smoke test for the full instance creation chain.
// Usage: go run ./cmd/bmmcli/smoketest
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/nvidia/bare-metal-manager-rest/client"
)

func main() {
	// Login
	fmt.Println("-- Login --")
	token, err := login()
	if err != nil {
		fmt.Fprintf(os.Stderr, "login failed: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("OK")

	// Setup client
	cfg := client.NewConfiguration()
	cfg.Servers = client.ServerConfigurations{{URL: "http://localhost:8388"}}
	c := client.NewAPIClient(cfg)
	ctx := context.WithValue(context.Background(), client.ContextAccessToken, token)
	org := "test-org"

	// List sites
	fmt.Println("\n-- Sites --")
	sites, _, err := c.SiteAPI.GetAllSite(ctx, org).Execute()
	must(err, "list sites")
	if len(sites) == 0 {
		fmt.Println("No sites found")
		os.Exit(1)
	}
	siteID := *sites[0].Id
	fmt.Printf("Site: %s (%s)\n", *sites[0].Name, siteID)

	// List VPCs
	fmt.Println("\n-- VPCs --")
	vpcs, _, err := c.VPCAPI.GetAllVpc(ctx, org).Execute()
	must(err, "list vpcs")
	if len(vpcs) == 0 {
		fmt.Println("No VPCs, skipping")
		os.Exit(0)
	}
	vpcID := *vpcs[0].Id
	fmt.Printf("VPC: %s (%s)\n", *vpcs[0].Name, vpcID)

	// Get tenant
	fmt.Println("\n-- Tenant --")
	tenant, _, err := c.TenantAPI.GetCurrentTenant(ctx, org).Execute()
	must(err, "get tenant")
	tenantID := *tenant.Id
	fmt.Printf("Tenant: %s\n", tenantID)

	// List IP blocks
	fmt.Println("\n-- IP Blocks --")
	blocks, _, err := c.IPBlockAPI.GetAllIpblock(ctx, org).Execute()
	must(err, "list ip blocks")
	if len(blocks) == 0 {
		fmt.Println("No IP blocks found, skipping subnet/instance test")
		os.Exit(0)
	}
	blockID := *blocks[0].Id
	fmt.Printf("IP Block: %s (%s)\n", *blocks[0].Name, blockID)

	// Create or find subnet
	fmt.Println("\n-- Subnet --")
	var subnetID string
	subnets, _, err := c.SubnetAPI.GetAllSubnet(ctx, org).VpcId(vpcID).Execute()
	must(err, "list subnets")
	if len(subnets) > 0 {
		subnetID = *subnets[0].Id
		fmt.Printf("Using existing subnet: %s (%s)\n", *subnets[0].Name, subnetID)
	} else {
		subReq := client.NewSubnetCreateRequest("smoke-test-subnet", vpcID, 24)
		subReq.Ipv4BlockId = &blockID
		subnet, resp, err := c.SubnetAPI.CreateSubnet(ctx, org).SubnetCreateRequest(*subReq).Execute()
		if err != nil {
			body := readBody(resp)
			fmt.Fprintf(os.Stderr, "create subnet failed: %v\n%s\n", err, body)
			os.Exit(1)
		}
		subnetID = *subnet.Id
		fmt.Printf("Subnet created: %s (%s) status=%s\n", *subnet.Name, subnetID, string(*subnet.Status))
	}

	// List instance types
	fmt.Println("\n-- Instance Types --")
	types, resp, err := c.InstanceTypeAPI.GetAllInstanceType(ctx, org).SiteId(siteID).TenantId(tenantID).Execute()
	if err != nil {
		body := readBody(resp)
		fmt.Fprintf(os.Stderr, "list instance types failed: %v\n%s\n", err, body)
		os.Exit(1)
	}
	if len(types) == 0 {
		fmt.Println("No instance types found, skipping instance create")
		os.Exit(0)
	}
	itID := *types[0].Id
	fmt.Printf("Instance Type: %s (%s)\n", *types[0].Name, itID)

	// Create instance (using ipxeScript since no OS images exist in local dev)
	fmt.Println("\n-- Create Instance --")
	iface := client.InterfaceCreateRequest{SubnetId: &subnetID}
	instanceName := fmt.Sprintf("smoke-test-%d", time.Now().Unix())
	instReq := client.NewInstanceCreateRequest(instanceName, tenantID, vpcID, []client.InterfaceCreateRequest{iface})
	instReq.InstanceTypeId = &itID
	instReq.SetIpxeScript("#!ipxe\necho smoke test")

	instance, resp, err := c.InstanceAPI.CreateInstance(ctx, org).InstanceCreateRequest(*instReq).Execute()
	if err != nil {
		body := readBody(resp)
		fmt.Fprintf(os.Stderr, "create instance failed: %v\n%s\n", err, body)
		os.Exit(1)
	}
	fmt.Printf("Instance created: %s (%s) status=%s\n", *instance.Name, *instance.Id, string(*instance.Status))

	fmt.Println("\n-- SMOKE TEST PASSED --")
}

func login() (string, error) {
	formData := url.Values{
		"grant_type":    {"password"},
		"client_id":     {"carbide-api"},
		"client_secret": {"carbide-local-secret"},
		"username":      {"admin@example.com"},
		"password":      {"adminpassword"},
	}
	resp, err := http.Post("http://localhost:8080/realms/carbide-dev/protocol/openid-connect/token",
		"application/x-www-form-urlencoded", strings.NewReader(formData.Encode()))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		return "", fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(body))
	}
	var tok struct {
		AccessToken string `json:"access_token"`
	}
	json.Unmarshal(body, &tok)
	return tok.AccessToken, nil
}

func must(err error, msg string) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s failed: %v\n", msg, err)
		os.Exit(1)
	}
}

func readBody(resp *http.Response) string {
	if resp == nil || resp.Body == nil {
		return ""
	}
	defer resp.Body.Close()
	b, _ := io.ReadAll(resp.Body)
	return string(b)
}
