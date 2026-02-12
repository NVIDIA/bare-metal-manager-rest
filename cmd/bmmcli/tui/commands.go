// SPDX-FileCopyrightText: Copyright (c) 2026 NVIDIA CORPORATION & AFFILIATES. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package tui

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/nvidia/bare-metal-manager-rest/client"
	"github.com/spf13/viper"
)

// LoginFunc is a callback to perform login and return a new token
type LoginFunc func() (string, error)

// Session holds the shared state for an interactive session
type Session struct {
	Client   *client.APIClient
	Ctx      context.Context
	Org      string
	Cache    *Cache
	Resolver *Resolver
	LoginFn  LoginFunc // set by the caller to enable in-session login
}

// RefreshToken updates the session context with a new token
func (s *Session) RefreshToken(token string) {
	s.Ctx = context.WithValue(context.Background(), client.ContextAccessToken, token)
}

// NewSession creates a new interactive session
func NewSession(apiClient *client.APIClient, ctx context.Context, org string) *Session {
	cache := NewCache()
	resolver := NewResolver(cache)

	s := &Session{
		Client:   apiClient,
		Ctx:      ctx,
		Org:      org,
		Cache:    cache,
		Resolver: resolver,
	}

	// Register all fetchers
	s.registerFetchers()

	return s
}

func (s *Session) registerFetchers() {
	s.Resolver.RegisterFetcher("site", s.fetchSites)
	s.Resolver.RegisterFetcher("vpc", s.fetchVPCs)
	s.Resolver.RegisterFetcher("subnet", s.fetchSubnets)
	s.Resolver.RegisterFetcher("instance", s.fetchInstances)
	s.Resolver.RegisterFetcher("instance-type", s.fetchInstanceTypes)
	s.Resolver.RegisterFetcher("operating-system", s.fetchOperatingSystems)
	s.Resolver.RegisterFetcher("ssh-key-group", s.fetchSSHKeyGroups)
	s.Resolver.RegisterFetcher("allocation", s.fetchAllocations)
	s.Resolver.RegisterFetcher("machine", s.fetchMachines)
	s.Resolver.RegisterFetcher("ip-block", s.fetchIPBlocks)
	s.Resolver.RegisterFetcher("network-security-group", s.fetchNetworkSecurityGroups)
}

// Command represents a registered interactive command
type Command struct {
	Name        string
	Description string
	Run         func(s *Session, args []string) error
}

// AllCommands returns all available commands
func AllCommands() []Command {
	return []Command{
		{Name: "site list", Description: "List all sites", Run: cmdSiteList},
		{Name: "vpc list", Description: "List all VPCs", Run: cmdVPCList},
		{Name: "vpc create", Description: "Create a VPC", Run: cmdVPCCreate},
		{Name: "subnet list", Description: "List all subnets", Run: cmdSubnetList},
		{Name: "subnet create", Description: "Create a subnet", Run: cmdSubnetCreate},
		{Name: "instance list", Description: "List all instances", Run: cmdInstanceList},
		{Name: "instance create", Description: "Create an instance (guided)", Run: cmdInstanceCreate},
		{Name: "instance-type list", Description: "List all instance types", Run: cmdInstanceTypeList},
		{Name: "instance-type create", Description: "Create an instance type", Run: cmdInstanceTypeCreate},
		{Name: "operating-system list", Description: "List operating systems", Run: cmdOSList},
		{Name: "ssh-key-group list", Description: "List SSH key groups", Run: cmdSSHKeyGroupList},
		{Name: "allocation list", Description: "List allocations", Run: cmdAllocationList},
		{Name: "machine list", Description: "List machines", Run: cmdMachineList},
		{Name: "ip-block list", Description: "List IP blocks", Run: cmdIPBlockList},
		{Name: "ip-block create", Description: "Create an IP block", Run: cmdIPBlockCreate},
		{Name: "network-security-group list", Description: "List network security groups", Run: cmdNSGList},
		{Name: "login", Description: "Login / refresh auth token", Run: cmdLogin},
		{Name: "help", Description: "Show available commands", Run: cmdHelp},
	}
}

// -- Fetchers --

func (s *Session) fetchSites(ctx context.Context) ([]NamedItem, error) {
	sites, _, err := s.Client.SiteAPI.GetAllSite(ctx, s.Org).Execute()
	if err != nil {
		return nil, err
	}
	items := make([]NamedItem, len(sites))
	for i, site := range sites {
		status := ""
		if site.Status != nil {
			status = string(*site.Status)
		}
		items[i] = NamedItem{
			Name:   ptrStr(site.Name),
			ID:     ptrStr(site.Id),
			Status: status,
			Raw:    site,
		}
	}
	return items, nil
}

func (s *Session) fetchVPCs(ctx context.Context) ([]NamedItem, error) {
	vpcs, _, err := s.Client.VPCAPI.GetAllVpc(ctx, s.Org).Execute()
	if err != nil {
		return nil, err
	}
	items := make([]NamedItem, len(vpcs))
	for i, vpc := range vpcs {
		status := ""
		if vpc.Status != nil {
			status = string(*vpc.Status)
		}
		items[i] = NamedItem{
			Name:   ptrStr(vpc.Name),
			ID:     ptrStr(vpc.Id),
			Status: status,
			Extra:  map[string]string{"siteId": ptrStr(vpc.SiteId)},
			Raw:    vpc,
		}
	}
	return items, nil
}

func (s *Session) fetchSubnets(ctx context.Context) ([]NamedItem, error) {
	subnets, _, err := s.Client.SubnetAPI.GetAllSubnet(ctx, s.Org).Execute()
	if err != nil {
		return nil, err
	}
	items := make([]NamedItem, len(subnets))
	for i, subnet := range subnets {
		status := ""
		if subnet.Status != nil {
			status = string(*subnet.Status)
		}
		items[i] = NamedItem{
			Name:   ptrStr(subnet.Name),
			ID:     ptrStr(subnet.Id),
			Status: status,
			Extra:  map[string]string{"vpcId": ptrStr(subnet.VpcId)},
			Raw:    subnet,
		}
	}
	return items, nil
}

func (s *Session) fetchInstances(ctx context.Context) ([]NamedItem, error) {
	instances, _, err := s.Client.InstanceAPI.GetAllInstance(ctx, s.Org).Execute()
	if err != nil {
		return nil, err
	}
	items := make([]NamedItem, len(instances))
	for i, inst := range instances {
		status := ""
		if inst.Status != nil {
			status = string(*inst.Status)
		}
		items[i] = NamedItem{
			Name:   ptrStr(inst.Name),
			ID:     ptrStr(inst.Id),
			Status: status,
			Extra: map[string]string{
				"vpcId":  ptrStr(inst.VpcId),
				"siteId": ptrStr(inst.SiteId),
			},
			Raw: inst,
		}
	}
	return items, nil
}

func (s *Session) fetchInstanceTypes(ctx context.Context) ([]NamedItem, error) {
	// Instance types require a siteId filter; return empty if no site context
	return nil, fmt.Errorf("siteId is required, use fetchInstanceTypesBySite instead")
}

// getTenantID returns the current tenant ID, caching it for the session
func (s *Session) getTenantID(ctx context.Context) (string, error) {
	cached := s.Cache.LookupByName("_tenant", s.Org)
	if cached != nil {
		return cached.ID, nil
	}
	tenant, _, err := s.Client.TenantAPI.GetCurrentTenant(ctx, s.Org).Execute()
	if err != nil {
		return "", fmt.Errorf("fetching tenant: %w", err)
	}
	tid := ptrStr(tenant.Id)
	s.Cache.Set("_tenant", []NamedItem{{Name: s.Org, ID: tid}})
	return tid, nil
}

// fetchInstanceTypesBySite fetches instance types filtered by site ID
func (s *Session) fetchInstanceTypesBySite(ctx context.Context, siteID string) ([]NamedItem, error) {
	tenantID, err := s.getTenantID(ctx)
	if err != nil {
		return nil, err
	}
	types, _, err := s.Client.InstanceTypeAPI.GetAllInstanceType(ctx, s.Org).SiteId(siteID).TenantId(tenantID).Execute()
	if err != nil {
		return nil, err
	}
	items := make([]NamedItem, len(types))
	for i, it := range types {
		status := ""
		if it.Status != nil {
			status = string(*it.Status)
		}
		siteID := ""
		if it.SiteId != nil {
			siteID = *it.SiteId
		}
		items[i] = NamedItem{
			Name:   ptrStr(it.Name),
			ID:     ptrStr(it.Id),
			Status: status,
			Extra:  map[string]string{"siteId": siteID},
			Raw:    it,
		}
	}
	return items, nil
}

func (s *Session) fetchOperatingSystems(ctx context.Context) ([]NamedItem, error) {
	osList, _, err := s.Client.OperatingSystemAPI.GetAllOperatingSystem(ctx, s.Org).Execute()
	if err != nil {
		return nil, err
	}
	items := make([]NamedItem, len(osList))
	for i, os := range osList {
		status := ""
		if os.Status != nil {
			status = string(*os.Status)
		}
		items[i] = NamedItem{
			Name:   ptrStr(os.Name),
			ID:     ptrStr(os.Id),
			Status: status,
			Raw:    os,
		}
	}
	return items, nil
}

func (s *Session) fetchSSHKeyGroups(ctx context.Context) ([]NamedItem, error) {
	groups, _, err := s.Client.SSHKeyGroupAPI.GetAllSshKeyGroup(ctx, s.Org).Execute()
	if err != nil {
		return nil, err
	}
	items := make([]NamedItem, len(groups))
	for i, g := range groups {
		status := ""
		if g.Status != nil {
			status = string(*g.Status)
		}
		items[i] = NamedItem{
			Name:   ptrStr(g.Name),
			ID:     ptrStr(g.Id),
			Status: status,
			Raw:    g,
		}
	}
	return items, nil
}

func (s *Session) fetchAllocations(ctx context.Context) ([]NamedItem, error) {
	allocs, _, err := s.Client.AllocationAPI.GetAllAllocation(ctx, s.Org).Execute()
	if err != nil {
		return nil, err
	}
	items := make([]NamedItem, len(allocs))
	for i, a := range allocs {
		status := ""
		if a.Status != nil {
			status = string(*a.Status)
		}
		items[i] = NamedItem{
			Name:   ptrStr(a.Name),
			ID:     ptrStr(a.Id),
			Status: status,
			Extra:  map[string]string{"siteId": ptrStr(a.SiteId)},
			Raw:    a,
		}
	}
	return items, nil
}

func (s *Session) fetchMachines(ctx context.Context) ([]NamedItem, error) {
	machines, _, err := s.Client.MachineAPI.GetAllMachine(ctx, s.Org).Execute()
	if err != nil {
		return nil, err
	}
	items := make([]NamedItem, len(machines))
	for i, m := range machines {
		status := ""
		if m.Status != nil {
			status = string(*m.Status)
		}
		name := ptrStr(m.Id)
		items[i] = NamedItem{
			Name:   name,
			ID:     ptrStr(m.Id),
			Status: status,
			Extra:  map[string]string{"siteId": ptrStr(m.SiteId)},
			Raw:    m,
		}
	}
	return items, nil
}

func (s *Session) fetchIPBlocks(ctx context.Context) ([]NamedItem, error) {
	blocks, _, err := s.Client.IPBlockAPI.GetAllIpblock(ctx, s.Org).Execute()
	if err != nil {
		return nil, err
	}
	items := make([]NamedItem, len(blocks))
	for i, b := range blocks {
		status := ""
		if b.Status != nil {
			status = string(*b.Status)
		}
		items[i] = NamedItem{
			Name:   ptrStr(b.Name),
			ID:     ptrStr(b.Id),
			Status: status,
			Extra:  map[string]string{"siteId": ptrStr(b.SiteId)},
			Raw:    b,
		}
	}
	return items, nil
}

func (s *Session) fetchNetworkSecurityGroups(ctx context.Context) ([]NamedItem, error) {
	nsgs, _, err := s.Client.NetworkSecurityGroupAPI.GetAllNetworkSecurityGroup(ctx, s.Org).Execute()
	if err != nil {
		return nil, err
	}
	items := make([]NamedItem, len(nsgs))
	for i, nsg := range nsgs {
		status := ""
		if nsg.Status != nil {
			status = string(*nsg.Status)
		}
		items[i] = NamedItem{
			Name:   ptrStr(nsg.Name),
			ID:     ptrStr(nsg.Id),
			Status: status,
			Raw:    nsg,
		}
	}
	return items, nil
}

// -- Command logging --

// LogCmd prints the equivalent CLI one-liner so users can copy/paste it
func LogCmd(parts ...string) {
	cmd := "bmmcli " + strings.Join(parts, " ")
	fmt.Println(Dim("$ " + cmd))
}

// -- Command handlers --

func cmdSiteList(s *Session, args []string) error {
	LogCmd("site", "list")
	items, err := s.Resolver.Fetch(s.Ctx, "site")
	if err != nil {
		return err
	}
	return printResourceTable(os.Stdout, "NAME", "STATUS", "ID", items)
}

func cmdVPCList(s *Session, args []string) error {
	LogCmd("vpc", "list")
	items, err := s.Resolver.Fetch(s.Ctx, "vpc")
	if err != nil {
		return err
	}
	tw := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
	fmt.Fprintln(tw, "NAME\tSTATUS\tSITE\tID")
	for _, item := range items {
		siteName := s.Resolver.ResolveID("site", item.Extra["siteId"])
		fmt.Fprintf(tw, "%s\t%s\t%s\t%s\n", item.Name, item.Status, siteName, item.ID)
	}
	return tw.Flush()
}

func cmdVPCCreate(s *Session, args []string) error {
	site, err := s.Resolver.Resolve(s.Ctx, "site", "Site")
	if err != nil {
		return err
	}

	name, err := PromptText("VPC name", true)
	if err != nil {
		return err
	}

	desc, err := PromptText("Description (optional)", false)
	if err != nil {
		return err
	}

	// Log the equivalent CLI command
	cmdParts := []string{"vpc", "create", "--name", name, "--site-id", site.ID}
	if desc != "" {
		cmdParts = append(cmdParts, "--description", fmt.Sprintf("%q", desc))
	}
	LogCmd(cmdParts...)

	req := client.NewVpcCreateRequest(name, site.ID)
	if desc != "" {
		req.SetDescription(desc)
	}

	vpc, _, err := s.Client.VPCAPI.CreateVpc(s.Ctx, s.Org).VpcCreateRequest(*req).Execute()
	if err != nil {
		return fmt.Errorf("creating VPC: %w", err)
	}

	s.Cache.Invalidate("vpc")
	fmt.Printf("%s VPC created: %s (%s)\n", Green("OK"), ptrStr(vpc.Name), ptrStr(vpc.Id))
	return nil
}

func cmdSubnetList(s *Session, args []string) error {
	LogCmd("subnet", "list")
	items, err := s.Resolver.Fetch(s.Ctx, "subnet")
	if err != nil {
		return err
	}
	tw := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
	fmt.Fprintln(tw, "NAME\tSTATUS\tVPC\tID")
	for _, item := range items {
		vpcName := s.Resolver.ResolveID("vpc", item.Extra["vpcId"])
		fmt.Fprintf(tw, "%s\t%s\t%s\t%s\n", item.Name, item.Status, vpcName, item.ID)
	}
	return tw.Flush()
}

func cmdSubnetCreate(s *Session, args []string) error {
	vpc, err := s.Resolver.Resolve(s.Ctx, "vpc", "VPC")
	if err != nil {
		return err
	}

	ipBlock, err := s.Resolver.Resolve(s.Ctx, "ip-block", "IPv4 Block")
	if err != nil {
		fmt.Println(Yellow("No IP blocks found. Create one first with: ip-block create"))
		return err
	}

	name, err := PromptText("Subnet name", true)
	if err != nil {
		return err
	}

	prefixLenStr, err := PromptText("Prefix length (e.g. 24)", true)
	if err != nil {
		return err
	}
	var prefixLen int32
	fmt.Sscanf(prefixLenStr, "%d", &prefixLen)
	if prefixLen < 1 || prefixLen > 32 {
		return fmt.Errorf("prefix length must be between 1 and 32")
	}

	LogCmd("subnet", "create", "--name", name, "--vpc-id", vpc.ID, "--ipv4-block-id", ipBlock.ID, "--prefix-length", prefixLenStr)

	req := client.NewSubnetCreateRequest(name, vpc.ID, prefixLen)
	req.Ipv4BlockId = &ipBlock.ID

	subnet, _, err := s.Client.SubnetAPI.CreateSubnet(s.Ctx, s.Org).SubnetCreateRequest(*req).Execute()
	if err != nil {
		return fmt.Errorf("creating subnet: %w", err)
	}

	s.Cache.Invalidate("subnet")
	fmt.Printf("%s Subnet created: %s (%s)\n", Green("OK"), ptrStr(subnet.Name), ptrStr(subnet.Id))
	return nil
}

func cmdInstanceList(s *Session, args []string) error {
	LogCmd("instance", "list")
	items, err := s.Resolver.Fetch(s.Ctx, "instance")
	if err != nil {
		return err
	}
	tw := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
	fmt.Fprintln(tw, "NAME\tSTATUS\tVPC\tSITE\tID")
	for _, item := range items {
		vpcName := s.Resolver.ResolveID("vpc", item.Extra["vpcId"])
		siteName := s.Resolver.ResolveID("site", item.Extra["siteId"])
		fmt.Fprintf(tw, "%s\t%s\t%s\t%s\t%s\n", item.Name, item.Status, vpcName, siteName, item.ID)
	}
	return tw.Flush()
}

func cmdInstanceCreate(s *Session, args []string) error {
	return RunInstanceWizard(s)
}

func cmdInstanceTypeList(s *Session, args []string) error {
	// Instance types require a site; prompt for one
	site, err := s.Resolver.Resolve(s.Ctx, "site", "Site")
	if err != nil {
		return err
	}
	LogCmd("instance-type", "list", "--site-id", site.ID)
	items, err := s.fetchInstanceTypesBySite(s.Ctx, site.ID)
	if err != nil {
		return err
	}
	tw := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
	fmt.Fprintln(tw, "NAME\tSTATUS\tSITE\tID")
	for _, item := range items {
		fmt.Fprintf(tw, "%s\t%s\t%s\t%s\n", item.Name, item.Status, site.Name, item.ID)
	}
	return tw.Flush()
}

func cmdInstanceTypeCreate(s *Session, args []string) error {
	site, err := s.Resolver.Resolve(s.Ctx, "site", "Site")
	if err != nil {
		return err
	}

	name, err := PromptText("Instance type name", true)
	if err != nil {
		return err
	}

	desc, err := PromptText("Description (optional)", false)
	if err != nil {
		return err
	}

	LogCmd("instance-type", "create", "--name", name, "--site-id", site.ID)

	req := client.NewInstanceTypeCreateRequest(name, site.ID)
	if desc != "" {
		req.Description = &desc
	}

	it, _, err := s.Client.InstanceTypeAPI.CreateInstanceType(s.Ctx, s.Org).InstanceTypeCreateRequest(*req).Execute()
	if err != nil {
		return fmt.Errorf("creating instance type: %w", err)
	}

	s.Cache.Invalidate("instance-type")
	fmt.Printf("%s Instance type created: %s (%s)\n", Green("OK"), ptrStr(it.Name), ptrStr(it.Id))
	return nil
}

func cmdOSList(s *Session, args []string) error {
	LogCmd("operating-system", "list")
	items, err := s.Resolver.Fetch(s.Ctx, "operating-system")
	if err != nil {
		return err
	}
	return printResourceTable(os.Stdout, "NAME", "STATUS", "ID", items)
}

func cmdSSHKeyGroupList(s *Session, args []string) error {
	LogCmd("ssh-key-group", "list")
	items, err := s.Resolver.Fetch(s.Ctx, "ssh-key-group")
	if err != nil {
		return err
	}
	return printResourceTable(os.Stdout, "NAME", "STATUS", "ID", items)
}

func cmdAllocationList(s *Session, args []string) error {
	LogCmd("allocation", "list")
	items, err := s.Resolver.Fetch(s.Ctx, "allocation")
	if err != nil {
		return err
	}
	tw := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
	fmt.Fprintln(tw, "NAME\tSTATUS\tSITE\tID")
	for _, item := range items {
		siteName := s.Resolver.ResolveID("site", item.Extra["siteId"])
		fmt.Fprintf(tw, "%s\t%s\t%s\t%s\n", item.Name, item.Status, siteName, item.ID)
	}
	return tw.Flush()
}

func cmdMachineList(s *Session, args []string) error {
	LogCmd("machine", "list")
	items, err := s.Resolver.Fetch(s.Ctx, "machine")
	if err != nil {
		return err
	}
	tw := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
	fmt.Fprintln(tw, "NAME\tSTATUS\tSITE\tID")
	for _, item := range items {
		siteName := s.Resolver.ResolveID("site", item.Extra["siteId"])
		fmt.Fprintf(tw, "%s\t%s\t%s\t%s\n", item.Name, item.Status, siteName, item.ID)
	}
	return tw.Flush()
}

func cmdIPBlockList(s *Session, args []string) error {
	LogCmd("ip-block", "list")
	items, err := s.Resolver.Fetch(s.Ctx, "ip-block")
	if err != nil {
		return err
	}
	tw := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
	fmt.Fprintln(tw, "NAME\tSTATUS\tSITE\tID")
	for _, item := range items {
		siteName := s.Resolver.ResolveID("site", item.Extra["siteId"])
		fmt.Fprintf(tw, "%s\t%s\t%s\t%s\n", item.Name, item.Status, siteName, item.ID)
	}
	return tw.Flush()
}

func cmdIPBlockCreate(s *Session, args []string) error {
	site, err := s.Resolver.Resolve(s.Ctx, "site", "Site")
	if err != nil {
		return err
	}

	name, err := PromptText("IP block name", true)
	if err != nil {
		return err
	}

	prefix, err := PromptText("Prefix (e.g. 10.0.0.0)", true)
	if err != nil {
		return err
	}

	prefixLen, err := PromptText("Prefix length (e.g. 16)", true)
	if err != nil {
		return err
	}

	var pl int32
	fmt.Sscanf(prefixLen, "%d", &pl)
	if pl < 1 || pl > 32 {
		return fmt.Errorf("prefix length must be between 1 and 32")
	}

	LogCmd("ip-block", "create", "--name", name, "--site-id", site.ID, "--prefix", prefix, "--prefix-length", prefixLen)

	req := client.NewIpBlockCreateRequest(name, site.ID, "DatacenterOnly", prefix, pl, "IPv4")

	block, _, err := s.Client.IPBlockAPI.CreateIpblock(s.Ctx, s.Org).IpBlockCreateRequest(*req).Execute()
	if err != nil {
		return fmt.Errorf("creating IP block: %w", err)
	}

	s.Cache.Invalidate("ip-block")
	fmt.Printf("%s IP block created: %s (%s)\n", Green("OK"), ptrStr(block.Name), ptrStr(block.Id))
	return nil
}

func cmdNSGList(s *Session, args []string) error {
	LogCmd("network-security-group", "list")
	items, err := s.Resolver.Fetch(s.Ctx, "network-security-group")
	if err != nil {
		return err
	}
	return printResourceTable(os.Stdout, "NAME", "STATUS", "ID", items)
}

func cmdLogin(s *Session, args []string) error {
	if s.LoginFn == nil {
		return fmt.Errorf("login not available (no OIDC provider configured)")
	}
	fmt.Println("Logging in...")
	token, err := s.LoginFn()
	if err != nil {
		return fmt.Errorf("login failed: %w", err)
	}
	s.RefreshToken(token)
	fmt.Printf("%s Logged in successfully.\n", Green("OK"))
	return nil
}

func cmdHelp(s *Session, args []string) error {
	tw := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
	fmt.Fprintln(tw, "\nCOMMAND\tDESCRIPTION")
	fmt.Fprintln(tw, "-------\t-----------")
	for _, cmd := range AllCommands() {
		fmt.Fprintf(tw, "%s\t%s\n", cmd.Name, cmd.Description)
	}
	fmt.Fprintln(tw, "exit\tExit interactive mode")
	tw.Flush()
	fmt.Println()
	return nil
}

// -- Helpers --

func printResourceTable(w io.Writer, col1, col2, col3 string, items []NamedItem) error {
	tw := tabwriter.NewWriter(w, 0, 0, 3, ' ', 0)
	fmt.Fprintf(tw, "%s\t%s\t%s\n", col1, col2, col3)
	for _, item := range items {
		fmt.Fprintf(tw, "%s\t%s\t%s\n", item.Name, item.Status, item.ID)
	}
	return tw.Flush()
}

func ptrStr(p *string) string {
	if p == nil {
		return ""
	}
	return *p
}

func formatAge(d time.Duration) string {
	if d < 0 {
		d = -d
	}
	switch {
	case d < time.Minute:
		return fmt.Sprintf("%ds", int(d.Seconds()))
	case d < time.Hour:
		return fmt.Sprintf("%dm", int(d.Minutes()))
	case d < 24*time.Hour:
		return fmt.Sprintf("%dh", int(d.Hours()))
	default:
		return fmt.Sprintf("%dd", int(d.Hours()/24))
	}
}

func getOrg() string {
	return viper.GetString("api.org")
}
