// SPDX-FileCopyrightText: Copyright (c) 2026 NVIDIA CORPORATION & AFFILIATES. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package interactive

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"text/tabwriter"
	"github.com/nvidia/bare-metal-manager-rest/sdk/standard"
)

// LoginFunc is a callback to perform login and return a new token.
type LoginFunc func() (string, error)

// Scope holds the current filter context for the interactive session.
// When set, list commands automatically filter by these values.
type Scope struct {
	SiteID   string
	SiteName string
	VpcID    string
	VpcName  string
}

// Session holds the shared state for an interactive session.
type Session struct {
	Client     *standard.APIClient
	Ctx        context.Context
	Org        string
	Token      string    // raw JWT token string for decoding claims
	Scope      Scope     // current filter context
	Cache      *Cache
	Resolver   *Resolver
	LoginFn    LoginFunc // set by the caller to enable in-session login
	ConfigPath string    // path to the config file, used in log output
	Verbose    bool      // when true, log raw HTTP details on errors
}

// PromptString returns the prompt showing org and current scope.
func (s *Session) PromptString() string {
	parts := []string{s.Org}
	if s.Scope.SiteName != "" {
		parts = append(parts, s.Scope.SiteName)
	}
	if s.Scope.VpcName != "" {
		parts = append(parts, s.Scope.VpcName)
	}
	return Cyan("bmm:"+strings.Join(parts, "/")) + "> "
}

// RefreshToken updates the session context with a new token.
func (s *Session) RefreshToken(token string) {
	s.Token = token
	s.Ctx = context.WithValue(context.Background(), standard.ContextAccessToken, token)
}

// NewSession creates a new interactive session.
func NewSession(apiClient *standard.APIClient, ctx context.Context, org string) *Session {
	cache := NewCache()
	resolver := NewResolver(cache)

	s := &Session{
		Client:   apiClient,
		Ctx:      ctx,
		Org:      org,
		Cache:    cache,
		Resolver: resolver,
	}

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
	s.Resolver.RegisterFetcher("audit", s.fetchAudits)
	s.Resolver.RegisterFetcher("ssh-key", s.fetchSSHKeys)
	s.Resolver.RegisterFetcher("sku", s.fetchSKUs)
	s.Resolver.RegisterFetcher("rack", s.fetchRacks)
	s.Resolver.RegisterFetcher("vpc-prefix", s.fetchVPCPrefixes)
	s.Resolver.RegisterFetcher("tenant-account", s.fetchTenantAccounts)
	s.Resolver.RegisterFetcher("expected-machine", s.fetchExpectedMachines)
	s.Resolver.RegisterFetcher("dpu-extension-service", s.fetchDPUExtensionServices)
	s.Resolver.RegisterFetcher("infiniband-partition", s.fetchInfiniBandPartitions)
	s.Resolver.RegisterFetcher("nvlink-logical-partition", s.fetchNVLinkLogicalPartitions)
}

// Command represents a registered interactive command.
type Command struct {
	Name        string
	Description string
	Run         func(s *Session, args []string) error
}

// AllCommands returns all available commands.
func AllCommands() []Command {
	return []Command{
		{Name: "site list", Description: "List all sites", Run: cmdSiteList},
		{Name: "site get", Description: "Get site details", Run: cmdSiteGet},
		{Name: "site delete", Description: "Delete a site", Run: cmdSiteDelete},

		{Name: "vpc list", Description: "List all VPCs", Run: cmdVPCList},
		{Name: "vpc get", Description: "Get VPC details", Run: cmdVPCGet},
		{Name: "vpc create", Description: "Create a VPC", Run: cmdVPCCreate},
		{Name: "vpc delete", Description: "Delete a VPC", Run: cmdVPCDelete},

		{Name: "subnet list", Description: "List all subnets", Run: cmdSubnetList},
		{Name: "subnet get", Description: "Get subnet details", Run: cmdSubnetGet},
		{Name: "subnet create", Description: "Create a subnet", Run: cmdSubnetCreate},
		{Name: "subnet delete", Description: "Delete a subnet", Run: cmdSubnetDelete},

		{Name: "instance list", Description: "List all instances", Run: cmdInstanceList},
		{Name: "instance get", Description: "Get instance details", Run: cmdInstanceGet},
		{Name: "instance create", Description: "Create an instance (guided)", Run: cmdInstanceCreate},
		{Name: "instance update", Description: "Update an instance", Run: cmdInstanceUpdate},
		{Name: "instance delete", Description: "Delete an instance", Run: cmdInstanceDelete},

		{Name: "instance-type list", Description: "List all instance types", Run: cmdInstanceTypeList},
		{Name: "instance-type get", Description: "Get instance type details", Run: cmdInstanceTypeGet},
		{Name: "instance-type create", Description: "Create an instance type", Run: cmdInstanceTypeCreate},
		{Name: "instance-type delete", Description: "Delete an instance type", Run: cmdInstanceTypeDelete},

		{Name: "operating-system list", Description: "List operating systems", Run: cmdOSList},
		{Name: "operating-system get", Description: "Get operating system details", Run: cmdOSGet},
		{Name: "operating-system create", Description: "Create an operating system", Run: cmdOSCreate},
		{Name: "operating-system update", Description: "Update an operating system", Run: cmdOSUpdate},
		{Name: "operating-system delete", Description: "Delete an operating system", Run: cmdOSDelete},

		{Name: "ssh-key-group list", Description: "List SSH key groups", Run: cmdSSHKeyGroupList},
		{Name: "ssh-key-group get", Description: "Get SSH key group details", Run: cmdSSHKeyGroupGet},

		{Name: "allocation list", Description: "List allocations", Run: cmdAllocationList},
		{Name: "allocation get", Description: "Get allocation details", Run: cmdAllocationGet},
		{Name: "allocation delete", Description: "Delete an allocation", Run: cmdAllocationDelete},

		{Name: "machine list", Description: "List machines", Run: cmdMachineList},
		{Name: "machine get", Description: "Get machine details", Run: cmdMachineGet},

		{Name: "ip-block list", Description: "List IP blocks", Run: cmdIPBlockList},
		{Name: "ip-block get", Description: "Get IP block details", Run: cmdIPBlockGet},
		{Name: "ip-block create", Description: "Create an IP block", Run: cmdIPBlockCreate},
		{Name: "ip-block delete", Description: "Delete an IP block", Run: cmdIPBlockDelete},

		{Name: "network-security-group list", Description: "List network security groups", Run: cmdNSGList},
		{Name: "network-security-group get", Description: "Get network security group details", Run: cmdNSGGet},

		{Name: "ssh-key list", Description: "List SSH keys", Run: cmdSSHKeyList},
		{Name: "ssh-key get", Description: "Get SSH key details", Run: cmdSSHKeyGet},

		{Name: "sku list", Description: "List SKUs", Run: cmdSKUList},
		{Name: "sku get", Description: "Get SKU details", Run: cmdSKUGet},

		{Name: "rack list", Description: "List racks", Run: cmdRackList},
		{Name: "rack get", Description: "Get rack details", Run: cmdRackGet},

		{Name: "vpc-prefix list", Description: "List VPC prefixes", Run: cmdVPCPrefixList},
		{Name: "vpc-prefix get", Description: "Get VPC prefix details", Run: cmdVPCPrefixGet},

		{Name: "tenant-account list", Description: "List tenant accounts", Run: cmdTenantAccountList},
		{Name: "tenant-account get", Description: "Get tenant account details", Run: cmdTenantAccountGet},

		{Name: "expected-machine list", Description: "List expected machines", Run: cmdExpectedMachineList},
		{Name: "expected-machine get", Description: "Get expected machine details", Run: cmdExpectedMachineGet},

		{Name: "dpu-extension-service list", Description: "List DPU extension services", Run: cmdDPUExtensionServiceList},
		{Name: "dpu-extension-service get", Description: "Get DPU extension service details", Run: cmdDPUExtensionServiceGet},

		{Name: "infiniband-partition list", Description: "List InfiniBand partitions", Run: cmdInfiniBandPartitionList},
		{Name: "infiniband-partition get", Description: "Get InfiniBand partition details", Run: cmdInfiniBandPartitionGet},

		{Name: "nvlink-logical-partition list", Description: "List NVLink logical partitions", Run: cmdNVLinkLogicalPartitionList},
		{Name: "nvlink-logical-partition get", Description: "Get NVLink logical partition details", Run: cmdNVLinkLogicalPartitionGet},

		{Name: "allocation-constraint list", Description: "List allocation constraints for an allocation", Run: cmdAllocationConstraintList},
		{Name: "allocation-constraint get", Description: "Get allocation constraint details", Run: cmdAllocationConstraintGet},

		{Name: "audit list", Description: "List audit log entries", Run: cmdAuditList},
		{Name: "audit get", Description: "Get audit log entry details", Run: cmdAuditGet},
		{Name: "metadata get", Description: "Get API metadata", Run: cmdMetadataGet},
		{Name: "user current", Description: "Get current user", Run: cmdUserCurrent},
		{Name: "service-account current", Description: "Get current service account", Run: cmdServiceAccountCurrent},
		{Name: "infrastructure-provider current", Description: "Get current infrastructure provider", Run: cmdInfrastructureProviderCurrent},
		{Name: "infrastructure-provider stats", Description: "Get infrastructure provider stats", Run: cmdInfrastructureProviderStats},
		{Name: "tenant current", Description: "Get current tenant", Run: cmdTenantCurrent},
		{Name: "tenant stats", Description: "Get tenant stats", Run: cmdTenantStats},
		{Name: "machine-capability list", Description: "List machine capabilities", Run: cmdMachineCapabilityList},

		{Name: "login", Description: "Login / refresh auth token", Run: cmdLogin},
		{Name: "help", Description: "Show available commands", Run: cmdHelp},
	}
}

// -- Fetchers --

func (s *Session) fetchSites(ctx context.Context) ([]NamedItem, error) {
	sites, _, err := FetchAllPages(s, "fetch-sites", func(pageNumber, pageSize int32) ([]standard.Site, *http.Response, error) {
		return s.Client.SiteAPI.GetAllSite(ctx, s.Org).PageNumber(pageNumber).PageSize(pageSize).Execute()
	})
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
	scopeSiteID := s.Scope.SiteID
	vpcs, _, err := FetchAllPages(s, "fetch-vpcs", func(pageNumber, pageSize int32) ([]standard.VPC, *http.Response, error) {
		req := s.Client.VPCAPI.GetAllVpc(ctx, s.Org).PageNumber(pageNumber).PageSize(pageSize)
		if scopeSiteID != "" {
			req = req.SiteId(scopeSiteID)
		}
		return req.Execute()
	})
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
	scopeSiteID := s.Scope.SiteID
	scopeVpcID := s.Scope.VpcID
	subnets, _, err := FetchAllPages(s, "fetch-subnets", func(pageNumber, pageSize int32) ([]standard.Subnet, *http.Response, error) {
		req := s.Client.SubnetAPI.GetAllSubnet(ctx, s.Org).PageNumber(pageNumber).PageSize(pageSize)
		if scopeSiteID != "" {
			req = req.SiteId(scopeSiteID)
		}
		if scopeVpcID != "" {
			req = req.VpcId(scopeVpcID)
		}
		return req.Execute()
	})
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
	scopeSiteID := s.Scope.SiteID
	scopeVpcID := s.Scope.VpcID
	instances, _, err := FetchAllPages(s, "fetch-instances", func(pageNumber, pageSize int32) ([]standard.Instance, *http.Response, error) {
		req := s.Client.InstanceAPI.GetAllInstance(ctx, s.Org).PageNumber(pageNumber).PageSize(pageSize)
		if scopeSiteID != "" {
			req = req.SiteId(scopeSiteID)
		}
		if scopeVpcID != "" {
			req = req.VpcId(scopeVpcID)
		}
		return req.Execute()
	})
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
	return nil, fmt.Errorf("siteId is required, use fetchInstanceTypesBySite instead")
}

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

func (s *Session) fetchInstanceTypesBySite(ctx context.Context, siteID string) ([]NamedItem, error) {
	tenantID, err := s.getTenantID(ctx)
	if err != nil {
		return nil, err
	}
	types, _, err := FetchAllPages(s, "fetch-instance-types", func(pageNumber, pageSize int32) ([]standard.InstanceType, *http.Response, error) {
		return s.Client.InstanceTypeAPI.GetAllInstanceType(ctx, s.Org).
			SiteId(siteID).
			TenantId(tenantID).
			PageNumber(pageNumber).
			PageSize(pageSize).
			Execute()
	})
	if err != nil {
		return nil, err
	}
	items := make([]NamedItem, len(types))
	for i, it := range types {
		status := ""
		if it.Status != nil {
			status = string(*it.Status)
		}
		itSiteID := ""
		if it.SiteId != nil {
			itSiteID = *it.SiteId
		}
		items[i] = NamedItem{
			Name:   ptrStr(it.Name),
			ID:     ptrStr(it.Id),
			Status: status,
			Extra:  map[string]string{"siteId": itSiteID},
			Raw:    it,
		}
	}
	return items, nil
}

func (s *Session) fetchOperatingSystems(ctx context.Context) ([]NamedItem, error) {
	scopeSiteID := s.Scope.SiteID
	osList, _, err := FetchAllPages(s, "fetch-operating-systems", func(pageNumber, pageSize int32) ([]standard.OperatingSystem, *http.Response, error) {
		req := s.Client.OperatingSystemAPI.GetAllOperatingSystem(ctx, s.Org).PageNumber(pageNumber).PageSize(pageSize)
		if scopeSiteID != "" {
			req = req.SiteId(scopeSiteID)
		}
		return req.Execute()
	})
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
	scopeSiteID := s.Scope.SiteID
	groups, _, err := FetchAllPages(s, "fetch-ssh-key-groups", func(pageNumber, pageSize int32) ([]standard.SshKeyGroup, *http.Response, error) {
		req := s.Client.SSHKeyGroupAPI.GetAllSshKeyGroup(ctx, s.Org).PageNumber(pageNumber).PageSize(pageSize)
		if scopeSiteID != "" {
			req = req.SiteId(scopeSiteID)
		}
		return req.Execute()
	})
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
	scopeSiteID := s.Scope.SiteID
	allocs, _, err := FetchAllPages(s, "fetch-allocations", func(pageNumber, pageSize int32) ([]standard.Allocation, *http.Response, error) {
		req := s.Client.AllocationAPI.GetAllAllocation(ctx, s.Org).PageNumber(pageNumber).PageSize(pageSize)
		if scopeSiteID != "" {
			req = req.SiteId(scopeSiteID)
		}
		return req.Execute()
	})
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
	return s.fetchMachinesWithScope(ctx, s.Scope.SiteID, s.Scope.VpcID)
}

func (s *Session) fetchMachinesWithSiteID(ctx context.Context, siteID string) ([]NamedItem, error) {
	return s.fetchMachinesWithScope(ctx, siteID, s.Scope.VpcID)
}

func (s *Session) fetchMachinesWithScope(ctx context.Context, siteID, vpcID string) ([]NamedItem, error) {
	machines, _, err := FetchAllPages(s, "fetch-machines", func(pageNumber, pageSize int32) ([]standard.Machine, *http.Response, error) {
		req := s.Client.MachineAPI.GetAllMachine(ctx, s.Org).PageNumber(pageNumber).PageSize(pageSize)
		if siteID != "" {
			req = req.SiteId(siteID)
		}
		return req.Execute()
	})
	if err != nil {
		return nil, err
	}

	if vpcID != "" {
		instances, _, err := FetchAllPages(s, "fetch-instances-for-machine-vpc", func(pageNumber, pageSize int32) ([]standard.Instance, *http.Response, error) {
			req := s.Client.InstanceAPI.GetAllInstance(ctx, s.Org).PageNumber(pageNumber).PageSize(pageSize).VpcId(vpcID)
			if siteID != "" {
				req = req.SiteId(siteID)
			}
			return req.Execute()
		})
		if err != nil {
			return nil, err
		}

		machineIDs := make(map[string]struct{}, len(instances))
		for _, inst := range instances {
			if machineID, ok := inst.GetMachineIdOk(); ok && machineID != nil {
				id := strings.TrimSpace(*machineID)
				if id != "" {
					machineIDs[id] = struct{}{}
				}
			}
		}

		filtered := make([]standard.Machine, 0, len(machines))
		for _, m := range machines {
			id := strings.TrimSpace(ptrStr(m.Id))
			if _, ok := machineIDs[id]; ok {
				filtered = append(filtered, m)
			}
		}
		machines = filtered
	}

	items := make([]NamedItem, len(machines))
	for i, m := range machines {
		status := ""
		if m.Status != nil {
			status = string(*m.Status)
		}
		name := machineDisplayName(m)
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

func (s *Session) machineVpcNamesByMachineID(ctx context.Context) (map[string]string, error) {
	_, _ = s.Resolver.Fetch(ctx, "vpc")

	scopeSiteID := s.Scope.SiteID
	scopeVpcID := s.Scope.VpcID
	instances, _, err := FetchAllPages(s, "fetch-instances-for-vpc-names", func(pageNumber, pageSize int32) ([]standard.Instance, *http.Response, error) {
		req := s.Client.InstanceAPI.GetAllInstance(ctx, s.Org).PageNumber(pageNumber).PageSize(pageSize)
		if scopeSiteID != "" {
			req = req.SiteId(scopeSiteID)
		}
		if scopeVpcID != "" {
			req = req.VpcId(scopeVpcID)
		}
		return req.Execute()
	})
	if err != nil {
		return nil, err
	}

	vpcSetByMachineID := make(map[string]map[string]struct{})
	for _, inst := range instances {
		machineID := ""
		if machineIDPtr, ok := inst.GetMachineIdOk(); ok && machineIDPtr != nil {
			machineID = strings.TrimSpace(*machineIDPtr)
		}
		vpcID := strings.TrimSpace(ptrStr(inst.VpcId))
		if machineID == "" || vpcID == "" {
			continue
		}
		if vpcSetByMachineID[machineID] == nil {
			vpcSetByMachineID[machineID] = make(map[string]struct{})
		}
		vpcSetByMachineID[machineID][vpcID] = struct{}{}
	}

	vpcNamesByMachineID := make(map[string]string, len(vpcSetByMachineID))
	for machineID, vpcSet := range vpcSetByMachineID {
		names := make([]string, 0, len(vpcSet))
		for vpcID := range vpcSet {
			name := strings.TrimSpace(s.Resolver.ResolveID("vpc", vpcID))
			if name == "" {
				name = vpcID
			}
			names = append(names, name)
		}
		sort.Strings(names)
		vpcNamesByMachineID[machineID] = strings.Join(names, ",")
	}
	return vpcNamesByMachineID, nil
}

func machineDisplayName(m standard.Machine) string {
	for _, key := range []string{"ServerName", "serverName", "hostname", "hostName"} {
		if v, ok := m.Labels[key]; ok {
			v = strings.TrimSpace(v)
			if v != "" {
				return v
			}
		}
	}

	for _, iface := range m.MachineInterfaces {
		if iface.Hostname != nil {
			v := strings.TrimSpace(*iface.Hostname)
			if v != "" {
				return v
			}
		}
	}

	if m.SerialNumber != nil {
		v := strings.TrimSpace(*m.SerialNumber)
		if v != "" {
			return v
		}
	}
	if m.ControllerMachineId != nil {
		v := strings.TrimSpace(*m.ControllerMachineId)
		if v != "" {
			return v
		}
	}
	id := strings.TrimSpace(ptrStr(m.Id))
	if id != "" {
		return id
	}
	return "<unknown>"
}

func (s *Session) fetchIPBlocks(ctx context.Context) ([]NamedItem, error) {
	scopeSiteID := s.Scope.SiteID
	blocks, _, err := FetchAllPages(s, "fetch-ip-blocks", func(pageNumber, pageSize int32) ([]standard.IpBlock, *http.Response, error) {
		req := s.Client.IPBlockAPI.GetAllIpblock(ctx, s.Org).PageNumber(pageNumber).PageSize(pageSize)
		if scopeSiteID != "" {
			req = req.SiteId(scopeSiteID)
		}
		return req.Execute()
	})
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
	scopeSiteID := s.Scope.SiteID
	nsgs, _, err := FetchAllPages(s, "fetch-nsgs", func(pageNumber, pageSize int32) ([]standard.NetworkSecurityGroup, *http.Response, error) {
		req := s.Client.NetworkSecurityGroupAPI.GetAllNetworkSecurityGroup(ctx, s.Org).PageNumber(pageNumber).PageSize(pageSize)
		if scopeSiteID != "" {
			req = req.SiteId(scopeSiteID)
		}
		return req.Execute()
	})
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

func (s *Session) fetchAudits(ctx context.Context) ([]NamedItem, error) {
	entries, _, err := FetchAllPages(s, "fetch-audits", func(pageNumber, pageSize int32) ([]standard.AuditEntry, *http.Response, error) {
		return s.Client.AuditAPI.GetAllAuditEntry(ctx, s.Org).PageNumber(pageNumber).PageSize(pageSize).Execute()
	})
	if err != nil {
		return nil, err
	}
	items := make([]NamedItem, len(entries))
	for i, entry := range entries {
		method := strings.TrimSpace(ptrStr(entry.Method))
		endpoint := strings.TrimSpace(ptrStr(entry.Endpoint))
		if method == "" {
			method = "AUDIT"
		}
		status := ""
		if entry.StatusCode != nil {
			status = fmt.Sprintf("%d", *entry.StatusCode)
		}
		name := strings.TrimSpace(method + " " + endpoint)
		if name == "" {
			name = ptrStr(entry.Id)
		}
		items[i] = NamedItem{
			Name:   name,
			ID:     ptrStr(entry.Id),
			Status: status,
			Extra:  map[string]string{"method": method, "endpoint": endpoint},
			Raw:    entry,
		}
	}
	return items, nil
}

func (s *Session) fetchSSHKeys(ctx context.Context) ([]NamedItem, error) {
	keys, _, err := FetchAllPages(s, "fetch-ssh-keys", func(pageNumber, pageSize int32) ([]standard.SshKey, *http.Response, error) {
		return s.Client.SSHKeyAPI.GetAllSshKey(ctx, s.Org).PageNumber(pageNumber).PageSize(pageSize).Execute()
	})
	if err != nil {
		return nil, err
	}
	items := make([]NamedItem, len(keys))
	for i, key := range keys {
		items[i] = NamedItem{
			Name:  ptrStr(key.Name),
			ID:    ptrStr(key.Id),
			Extra: map[string]string{"fingerprint": ptrStr(key.Fingerprint)},
			Raw:   key,
		}
	}
	return items, nil
}

func (s *Session) fetchSKUs(ctx context.Context) ([]NamedItem, error) {
	scopeSiteID := s.Scope.SiteID
	skus, _, err := FetchAllPages(s, "fetch-skus", func(pageNumber, pageSize int32) ([]standard.Sku, *http.Response, error) {
		req := s.Client.SKUAPI.GetAllSku(ctx, s.Org).PageNumber(pageNumber).PageSize(pageSize)
		if scopeSiteID != "" {
			req = req.SiteId(scopeSiteID)
		}
		return req.Execute()
	})
	if err != nil {
		return nil, err
	}
	items := make([]NamedItem, len(skus))
	for i, sku := range skus {
		deviceType := ""
		if sku.DeviceType.IsSet() {
			deviceType = *sku.DeviceType.Get()
		}
		name := deviceType
		if strings.TrimSpace(name) == "" {
			name = ptrStr(sku.Id)
		}
		items[i] = NamedItem{
			Name:  name,
			ID:    ptrStr(sku.Id),
			Extra: map[string]string{"siteId": ptrStr(sku.SiteId), "deviceType": deviceType},
			Raw:   sku,
		}
	}
	return items, nil
}

func (s *Session) fetchRacks(ctx context.Context) ([]NamedItem, error) {
	scopeSiteID := s.Scope.SiteID
	racks, _, err := FetchAllPages(s, "fetch-racks", func(pageNumber, pageSize int32) ([]standard.Rack, *http.Response, error) {
		req := s.Client.RackAPI.GetAllRack(ctx, s.Org).PageNumber(pageNumber).PageSize(pageSize)
		if scopeSiteID != "" {
			req = req.SiteId(scopeSiteID)
		}
		return req.Execute()
	})
	if err != nil {
		return nil, err
	}
	items := make([]NamedItem, len(racks))
	for i, rack := range racks {
		items[i] = NamedItem{
			Name: ptrStr(rack.Name),
			ID:   ptrStr(rack.Id),
			Extra: map[string]string{
				"manufacturer": ptrStr(rack.Manufacturer),
				"model":        ptrStr(rack.Model),
			},
			Raw: rack,
		}
	}
	return items, nil
}

func (s *Session) fetchVPCPrefixes(ctx context.Context) ([]NamedItem, error) {
	scopeSiteID := s.Scope.SiteID
	scopeVpcID := s.Scope.VpcID
	prefixes, _, err := FetchAllPages(s, "fetch-vpc-prefixes", func(pageNumber, pageSize int32) ([]standard.VpcPrefix, *http.Response, error) {
		req := s.Client.VPCPrefixAPI.GetAllVpcPrefix(ctx, s.Org).PageNumber(pageNumber).PageSize(pageSize)
		if scopeSiteID != "" {
			req = req.SiteId(scopeSiteID)
		}
		if scopeVpcID != "" {
			req = req.VpcId(scopeVpcID)
		}
		return req.Execute()
	})
	if err != nil {
		return nil, err
	}
	items := make([]NamedItem, len(prefixes))
	for i, prefix := range prefixes {
		status := ""
		if prefix.Status != nil {
			status = string(*prefix.Status)
		}
		items[i] = NamedItem{
			Name:   ptrStr(prefix.Name),
			ID:     ptrStr(prefix.Id),
			Status: status,
			Extra:  map[string]string{"vpcId": ptrStr(prefix.VpcId)},
			Raw:    prefix,
		}
	}
	return items, nil
}

func (s *Session) fetchTenantAccounts(ctx context.Context) ([]NamedItem, error) {
	accounts, _, err := FetchAllPages(s, "fetch-tenant-accounts", func(pageNumber, pageSize int32) ([]standard.TenantAccount, *http.Response, error) {
		return s.Client.TenantAccountAPI.GetAllTenantAccount(ctx, s.Org).PageNumber(pageNumber).PageSize(pageSize).Execute()
	})
	if err != nil {
		return nil, err
	}
	items := make([]NamedItem, len(accounts))
	for i, account := range accounts {
		status := ""
		if account.Status != nil {
			status = string(*account.Status)
		}
		tenantOrg := ""
		if account.TenantOrg.IsSet() {
			tenantOrg = *account.TenantOrg.Get()
		}
		name := strings.TrimSpace(tenantOrg)
		if name == "" {
			name = ptrStr(account.Id)
		}
		items[i] = NamedItem{
			Name:   name,
			ID:     ptrStr(account.Id),
			Status: status,
			Extra:  map[string]string{"infrastructureProviderId": ptrStr(account.InfrastructureProviderId)},
			Raw:    account,
		}
	}
	return items, nil
}

func (s *Session) fetchExpectedMachines(ctx context.Context) ([]NamedItem, error) {
	scopeSiteID := s.Scope.SiteID
	machines, _, err := FetchAllPages(s, "fetch-expected-machines", func(pageNumber, pageSize int32) ([]standard.ExpectedMachine, *http.Response, error) {
		req := s.Client.ExpectedMachineAPI.GetAllExpectedMachine(ctx, s.Org).PageNumber(pageNumber).PageSize(pageSize)
		if scopeSiteID != "" {
			req = req.SiteId(scopeSiteID)
		}
		return req.Execute()
	})
	if err != nil {
		return nil, err
	}
	items := make([]NamedItem, len(machines))
	for i, machine := range machines {
		name := strings.TrimSpace(ptrStr(machine.BmcMacAddress))
		if name == "" {
			name = strings.TrimSpace(ptrStr(machine.ChassisSerialNumber))
		}
		if name == "" {
			name = ptrStr(machine.Id)
		}
		items[i] = NamedItem{
			Name: name,
			ID:   ptrStr(machine.Id),
			Extra: map[string]string{
				"siteId":              ptrStr(machine.SiteId),
				"bmcMacAddress":       ptrStr(machine.BmcMacAddress),
				"chassisSerialNumber": ptrStr(machine.ChassisSerialNumber),
			},
			Raw: machine,
		}
	}
	return items, nil
}

func (s *Session) fetchDPUExtensionServices(ctx context.Context) ([]NamedItem, error) {
	scopeSiteID := s.Scope.SiteID
	services, _, err := FetchAllPages(s, "fetch-dpu-extension-services", func(pageNumber, pageSize int32) ([]standard.DpuExtensionService, *http.Response, error) {
		req := s.Client.DPUExtensionServiceAPI.GetAllDpuExtensionService(ctx, s.Org).PageNumber(pageNumber).PageSize(pageSize)
		if scopeSiteID != "" {
			req = req.SiteId(scopeSiteID)
		}
		return req.Execute()
	})
	if err != nil {
		return nil, err
	}
	items := make([]NamedItem, len(services))
	for i, svc := range services {
		items[i] = NamedItem{
			Name: ptrStr(svc.Name),
			ID:   ptrStr(svc.Id),
			Extra: map[string]string{
				"siteId":      ptrStr(svc.SiteId),
				"serviceType": ptrStr(svc.ServiceType),
			},
			Raw: svc,
		}
	}
	return items, nil
}

func (s *Session) fetchInfiniBandPartitions(ctx context.Context) ([]NamedItem, error) {
	scopeSiteID := s.Scope.SiteID
	partitions, _, err := FetchAllPages(s, "fetch-infiniband-partitions", func(pageNumber, pageSize int32) ([]standard.InfiniBandPartition, *http.Response, error) {
		req := s.Client.InfiniBandPartitionAPI.GetAllInfinibandPartition(ctx, s.Org).PageNumber(pageNumber).PageSize(pageSize)
		if scopeSiteID != "" {
			req = req.SiteId(scopeSiteID)
		}
		return req.Execute()
	})
	if err != nil {
		return nil, err
	}
	items := make([]NamedItem, len(partitions))
	for i, partition := range partitions {
		status := ""
		if partition.Status != nil {
			status = string(*partition.Status)
		}
		items[i] = NamedItem{
			Name:   ptrStr(partition.Name),
			ID:     ptrStr(partition.Id),
			Status: status,
			Extra:  map[string]string{"siteId": ptrStr(partition.SiteId)},
			Raw:    partition,
		}
	}
	return items, nil
}

func (s *Session) fetchNVLinkLogicalPartitions(ctx context.Context) ([]NamedItem, error) {
	scopeSiteID := s.Scope.SiteID
	partitions, _, err := FetchAllPages(s, "fetch-nvlink-partitions", func(pageNumber, pageSize int32) ([]standard.NVLinkLogicalPartition, *http.Response, error) {
		req := s.Client.NVLinkLogicalPartitionAPI.GetAllNvlinkLogicalPartition(ctx, s.Org).PageNumber(pageNumber).PageSize(pageSize)
		if scopeSiteID != "" {
			req = req.SiteId(scopeSiteID)
		}
		return req.Execute()
	})
	if err != nil {
		return nil, err
	}
	items := make([]NamedItem, len(partitions))
	for i, partition := range partitions {
		status := ""
		if partition.Status != nil {
			status = string(*partition.Status)
		}
		items[i] = NamedItem{
			Name:   ptrStr(partition.Name),
			ID:     ptrStr(partition.Id),
			Status: status,
			Extra:  map[string]string{"siteId": ptrStr(partition.SiteId)},
			Raw:    partition,
		}
	}
	return items, nil
}

type vpcUsageCounts struct {
	Instances int
	Machines  int
}

func (s *Session) fetchVpcUsageCounts(ctx context.Context) (map[string]vpcUsageCounts, error) {
	scopeSiteID := s.Scope.SiteID
	instances, _, err := FetchAllPages(s, "fetch-instances-for-vpc-counts", func(pageNumber, pageSize int32) ([]standard.Instance, *http.Response, error) {
		req := s.Client.InstanceAPI.GetAllInstance(ctx, s.Org).PageNumber(pageNumber).PageSize(pageSize)
		if scopeSiteID != "" {
			req = req.SiteId(scopeSiteID)
		}
		return req.Execute()
	})
	if err != nil {
		return nil, err
	}

	usageByVpc := make(map[string]vpcUsageCounts)
	machineSetByVpc := make(map[string]map[string]struct{})
	for _, inst := range instances {
		vpcID := strings.TrimSpace(ptrStr(inst.VpcId))
		if vpcID == "" {
			continue
		}
		usage := usageByVpc[vpcID]
		usage.Instances++
		usageByVpc[vpcID] = usage

		if machineID, ok := inst.GetMachineIdOk(); ok && machineID != nil {
			id := strings.TrimSpace(*machineID)
			if id != "" {
				if machineSetByVpc[vpcID] == nil {
					machineSetByVpc[vpcID] = make(map[string]struct{})
				}
				machineSetByVpc[vpcID][id] = struct{}{}
			}
		}
	}

	for vpcID, set := range machineSetByVpc {
		usage := usageByVpc[vpcID]
		usage.Machines = len(set)
		usageByVpc[vpcID] = usage
	}
	return usageByVpc, nil
}

// -- Command logging --

// LogCmd prints the equivalent CLI one-liner. When a session is available,
// prefer LogScopedCmd(s, ...) so --config is included in the output.
func LogCmd(parts ...string) {
	logCmdWithScope(nil, parts...)
}

func defaultConfigPath() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".bmm", "config.yaml")
}

// LogScopedCmd prints the equivalent CLI command including active scope filters.
func LogScopedCmd(s *Session, parts ...string) {
	logCmdWithScope(s, parts...)
}

func logCmdWithScope(s *Session, parts ...string) {
	cmdParts := []string{"bmmcli"}
	if s != nil && s.ConfigPath != "" && s.ConfigPath != defaultConfigPath() {
		cmdParts = append(cmdParts, "--config", s.ConfigPath)
	}
	cmdParts = append(cmdParts, appendScopeFlags(s, parts)...)
	cmd := strings.Join(cmdParts, " ")
	fmt.Printf("INFO: %s\n", cmd)
}

func appendScopeFlags(s *Session, parts []string) []string {
	out := append([]string(nil), parts...)
	if s == nil || len(parts) < 2 {
		return out
	}

	resource := strings.TrimSpace(parts[0])
	action := strings.TrimSpace(parts[1])
	if action != "list" {
		return out
	}

	scopeSiteID := strings.TrimSpace(s.Scope.SiteID)
	scopeVpcID := strings.TrimSpace(s.Scope.VpcID)

	switch resource {
	case "vpc", "allocation", "ip-block",
		"operating-system", "ssh-key-group", "network-security-group",
		"sku", "rack", "expected-machine", "dpu-extension-service",
		"infiniband-partition", "nvlink-logical-partition", "machine-capability":
		if scopeSiteID != "" && !hasArgFlag(out, "--site-id") {
			out = append(out, "--site-id", scopeSiteID)
		}
	case "subnet", "vpc-prefix":
		if scopeSiteID != "" && !hasArgFlag(out, "--site-id") {
			out = append(out, "--site-id", scopeSiteID)
		}
		if scopeVpcID != "" && !hasArgFlag(out, "--vpc-id") {
			out = append(out, "--vpc-id", scopeVpcID)
		}
	case "instance", "machine":
		if scopeSiteID != "" && !hasArgFlag(out, "--site-id") {
			out = append(out, "--site-id", scopeSiteID)
		}
		if scopeVpcID != "" && !hasArgFlag(out, "--vpc-id") {
			out = append(out, "--vpc-id", scopeVpcID)
		}
	}

	return out
}

func hasArgFlag(parts []string, flag string) bool {
	for _, part := range parts {
		if part == flag {
			return true
		}
	}
	return false
}

// -- Command handlers --

func cmdSiteList(s *Session, args []string) error {
	LogScopedCmd(s, "site", "list")
	items, err := s.Resolver.Fetch(s.Ctx, "site")
	if err != nil {
		return err
	}
	return printResourceTable(os.Stdout, "NAME", "STATUS", "ID", items)
}

func cmdVPCList(s *Session, args []string) error {
	LogScopedCmd(s, "vpc", "list")
	items, err := s.Resolver.Fetch(s.Ctx, "vpc")
	if err != nil {
		return err
	}
	siteNameByID := map[string]string{}
	if sites, siteErr := s.Resolver.Fetch(s.Ctx, "site"); siteErr == nil {
		for _, site := range sites {
			siteNameByID[site.ID] = site.Name
		}
	}
	usageByVpc, usageErr := s.fetchVpcUsageCounts(s.Ctx)
	if usageErr != nil {
		fmt.Printf("%s Could not load VPC usage counts: %v\n", Yellow("Note:"), usageErr)
	}
	PrintSummary(os.Stdout, nil, len(items))
	tw := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
	fmt.Fprintln(tw, "NAME\tSTATUS\tSITE\tINSTANCES\tMACHINES\tID")
	for _, item := range items {
		siteID := item.Extra["siteId"]
		siteName := strings.TrimSpace(siteNameByID[siteID])
		if siteName == "" {
			siteName = siteID
		}
		usage := usageByVpc[item.ID]
		fmt.Fprintf(tw, "%s\t%s\t%s\t%d\t%d\t%s\n", item.Name, item.Status, siteName, usage.Instances, usage.Machines, item.ID)
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

	cmdParts := []string{"vpc", "create", "--name", name, "--site-id", site.ID}
	if desc != "" {
		cmdParts = append(cmdParts, "--description", fmt.Sprintf("%q", desc))
	}
	LogScopedCmd(s, cmdParts...)

	req := standard.NewVpcCreateRequest(name, site.ID)
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
	LogScopedCmd(s, "subnet", "list")
	items, err := s.Resolver.Fetch(s.Ctx, "subnet")
	if err != nil {
		return err
	}
	PrintSummary(os.Stdout, nil, len(items))
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

	LogScopedCmd(s, "subnet", "create", "--name", name, "--vpc-id", vpc.ID, "--ipv4-block-id", ipBlock.ID, "--prefix-length", prefixLenStr)

	req := standard.NewSubnetCreateRequest(name, vpc.ID, prefixLen)
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
	LogScopedCmd(s, "instance", "list")
	items, err := s.Resolver.Fetch(s.Ctx, "instance")
	if err != nil {
		return err
	}
	PrintSummary(os.Stdout, nil, len(items))
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

func cmdInstanceUpdate(s *Session, args []string) error {
	instance, err := s.Resolver.ResolveWithArgs(s.Ctx, "instance", "Instance to update", args)
	if err != nil {
		return err
	}

	name, err := PromptText("New instance name (optional)", false)
	if err != nil {
		return err
	}
	desc, err := PromptText("New description (optional)", false)
	if err != nil {
		return err
	}
	if strings.TrimSpace(name) == "" && strings.TrimSpace(desc) == "" {
		return fmt.Errorf("no update fields provided")
	}

	logArgs := []string{"instance", "update", instance.ID}
	if strings.TrimSpace(name) != "" {
		logArgs = append(logArgs, "--name", name)
	}
	if strings.TrimSpace(desc) != "" {
		logArgs = append(logArgs, "--description", desc)
	}
	LogScopedCmd(s, logArgs...)

	req := standard.NewInstanceUpdateRequest()
	if strings.TrimSpace(name) != "" {
		req.SetName(name)
	}
	if strings.TrimSpace(desc) != "" {
		req.SetDescription(desc)
	}

	updated, _, err := s.Client.InstanceAPI.UpdateInstance(s.Ctx, s.Org, instance.ID).InstanceUpdateRequest(*req).Execute()
	if err != nil {
		return fmt.Errorf("updating instance: %w", err)
	}

	s.Cache.Invalidate("instance")
	fmt.Printf("%s Instance updated: %s (%s)\n", Green("OK"), ptrStr(updated.Name), ptrStr(updated.Id))
	return printDetailJSON(os.Stdout, updated)
}

func cmdInstanceTypeList(s *Session, args []string) error {
	site, err := s.Resolver.Resolve(s.Ctx, "site", "Site")
	if err != nil {
		return err
	}
	LogScopedCmd(s, "instance-type", "list", "--site-id", site.ID)
	items, err := s.fetchInstanceTypesBySite(s.Ctx, site.ID)
	if err != nil {
		return err
	}
	PrintSummary(os.Stdout, nil, len(items))
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

	LogScopedCmd(s, "instance-type", "create", "--name", name, "--site-id", site.ID)

	req := standard.NewInstanceTypeCreateRequest(name, site.ID)
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
	LogScopedCmd(s, "operating-system", "list")
	items, err := s.Resolver.Fetch(s.Ctx, "operating-system")
	if err != nil {
		return err
	}
	return printResourceTable(os.Stdout, "NAME", "STATUS", "ID", items)
}

func cmdOSCreate(s *Session, args []string) error {
	name, err := PromptText("Operating system name", true)
	if err != nil {
		return err
	}
	description, err := PromptText("Description (optional)", false)
	if err != nil {
		return err
	}
	ipxeScript, err := PromptText("iPXE script or URL (optional)", false)
	if err != nil {
		return err
	}
	imageURL, err := PromptText("Image URL (optional)", false)
	if err != nil {
		return err
	}
	if strings.TrimSpace(ipxeScript) != "" && strings.TrimSpace(imageURL) != "" {
		return fmt.Errorf("provide either iPXE script or image URL, not both")
	}

	req := standard.NewOperatingSystemCreateRequest(name)
	tenantID, tenantErr := s.getTenantID(s.Ctx)
	if tenantErr == nil && strings.TrimSpace(tenantID) != "" {
		req.SetTenantId(tenantID)
	} else {
		provider, _, providerErr := s.Client.InfrastructureProviderAPI.GetCurrentInfrastructureProvider(s.Ctx, s.Org).Execute()
		if providerErr != nil {
			if tenantErr != nil {
				return fmt.Errorf("unable to determine owning tenant/provider for operating system creation: tenant error: %v; provider error: %w", tenantErr, providerErr)
			}
			return fmt.Errorf("unable to determine owning provider for operating system creation: %w", providerErr)
		}
		providerID := strings.TrimSpace(ptrStr(provider.Id))
		if providerID == "" {
			return fmt.Errorf("unable to determine owning tenant/provider for operating system creation")
		}
		req.SetInfrastructureProviderId(providerID)
	}
	if strings.TrimSpace(description) != "" {
		req.SetDescription(description)
	}
	if strings.TrimSpace(s.Scope.SiteID) != "" {
		req.SetSiteIds([]string{s.Scope.SiteID})
	}
	if strings.TrimSpace(ipxeScript) != "" {
		req.SetIpxeScript(ipxeScript)
	}
	if strings.TrimSpace(imageURL) != "" {
		req.SetImageUrl(imageURL)
	}

	logArgs := []string{"operating-system", "create", "--name", name}
	if req.HasTenantId() {
		logArgs = append(logArgs, "--tenant-id", req.GetTenantId())
	}
	if req.HasInfrastructureProviderId() {
		logArgs = append(logArgs, "--infrastructure-provider-id", req.GetInfrastructureProviderId())
	}
	if strings.TrimSpace(s.Scope.SiteID) != "" {
		logArgs = append(logArgs, "--site-id", s.Scope.SiteID)
	}
	if strings.TrimSpace(ipxeScript) != "" {
		logArgs = append(logArgs, "--ipxe-script", strconv.Quote(ipxeScript))
	}
	if strings.TrimSpace(imageURL) != "" {
		logArgs = append(logArgs, "--image-url", imageURL)
	}
	LogScopedCmd(s, logArgs...)

	created, _, err := s.Client.OperatingSystemAPI.CreateOperatingSystem(s.Ctx, s.Org).OperatingSystemCreateRequest(*req).Execute()
	if err != nil {
		return fmt.Errorf("creating operating system: %w", err)
	}

	s.Cache.Invalidate("operating-system")
	fmt.Printf("%s Operating system created: %s (%s)\n", Green("OK"), ptrStr(created.Name), ptrStr(created.Id))
	return printDetailJSON(os.Stdout, created)
}

func cmdOSUpdate(s *Session, args []string) error {
	osItem, err := s.Resolver.ResolveWithArgs(s.Ctx, "operating-system", "Operating System to update", args)
	if err != nil {
		return err
	}

	name, err := PromptText("New operating system name (optional)", false)
	if err != nil {
		return err
	}
	description, err := PromptText("New description (optional)", false)
	if err != nil {
		return err
	}
	if strings.TrimSpace(name) == "" && strings.TrimSpace(description) == "" {
		return fmt.Errorf("no update fields provided")
	}

	logArgs := []string{"operating-system", "update", osItem.ID}
	if strings.TrimSpace(name) != "" {
		logArgs = append(logArgs, "--name", name)
	}
	if strings.TrimSpace(description) != "" {
		logArgs = append(logArgs, "--description", description)
	}
	LogScopedCmd(s, logArgs...)

	req := standard.NewOperatingSystemUpdateRequest()
	if strings.TrimSpace(name) != "" {
		req.SetName(name)
	}
	if strings.TrimSpace(description) != "" {
		req.SetDescription(description)
	}

	updated, _, err := s.Client.OperatingSystemAPI.UpdateOperatingSystem(s.Ctx, s.Org, osItem.ID).OperatingSystemUpdateRequest(*req).Execute()
	if err != nil {
		return fmt.Errorf("updating operating system: %w", err)
	}

	s.Cache.Invalidate("operating-system")
	fmt.Printf("%s Operating system updated: %s (%s)\n", Green("OK"), ptrStr(updated.Name), ptrStr(updated.Id))
	return printDetailJSON(os.Stdout, updated)
}

func cmdSSHKeyGroupList(s *Session, args []string) error {
	LogScopedCmd(s, "ssh-key-group", "list")
	items, err := s.Resolver.Fetch(s.Ctx, "ssh-key-group")
	if err != nil {
		return err
	}
	return printResourceTable(os.Stdout, "NAME", "STATUS", "ID", items)
}

func cmdAllocationList(s *Session, args []string) error {
	LogScopedCmd(s, "allocation", "list")
	items, err := s.Resolver.Fetch(s.Ctx, "allocation")
	if err != nil {
		return err
	}
	PrintSummary(os.Stdout, nil, len(items))
	tw := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
	fmt.Fprintln(tw, "NAME\tSTATUS\tSITE\tID")
	for _, item := range items {
		siteName := s.Resolver.ResolveID("site", item.Extra["siteId"])
		fmt.Fprintf(tw, "%s\t%s\t%s\t%s\n", item.Name, item.Status, siteName, item.ID)
	}
	return tw.Flush()
}

func cmdMachineList(s *Session, args []string) error {
	LogScopedCmd(s, "machine", "list")
	if s.Scope.VpcID != "" {
		fmt.Printf("%s Showing machines attached to instances in scoped VPC.\n", Dim("Note:"))
	}

	items, err := s.Resolver.Fetch(s.Ctx, "machine")
	if err != nil {
		// Tenant orgs lack FORGE_PROVIDER_ADMIN and receive 400/403 from the
		// machine endpoint without a site filter. When that happens and no site
		// scope is set, prompt for one and retry with the site-scoped path.
		needsSiteScope := s.Scope.SiteID == "" && (strings.Contains(err.Error(), "400") ||
			strings.Contains(err.Error(), "403") ||
			strings.Contains(err.Error(), "permission denied") ||
			strings.Contains(err.Error(), "Bad Request"))
		if !needsSiteScope {
			return err
		}

		fmt.Printf("%s Machine listing requires a site filter for your org. Select a site.\n", Dim("Note:"))
		site, resolveErr := s.Resolver.Resolve(s.Ctx, "site", "Site")
		if resolveErr != nil {
			return resolveErr
		}
		s.Scope.SiteID = site.ID
		s.Scope.SiteName = site.Name
		s.Cache.InvalidateFiltered()

		items, err = s.fetchMachinesWithSiteID(s.Ctx, site.ID)
		if err != nil {
			return err
		}
		fmt.Printf("%s Applied site scope %s (%s) for machine listing.\n", Dim("Note:"), site.Name, site.ID)
	}

	if s.Scope.VpcID != "" && len(items) == 0 {
		fmt.Printf("%s No machines are currently attached to instances in this VPC.\n", Dim("Note:"))
	}

	vpcNamesByMachineID, vpcErr := s.machineVpcNamesByMachineID(s.Ctx)
	if vpcErr != nil {
		fmt.Printf("%s Could not load machine VPC names: %v\n", Yellow("Note:"), vpcErr)
	}

	PrintSummary(os.Stdout, nil, len(items))
	tw := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
	fmt.Fprintln(tw, "NAME\tSTATUS\tSITE\tVPC\tID")
	for _, item := range items {
		siteName := s.Resolver.ResolveID("site", item.Extra["siteId"])
		vpcNames := strings.TrimSpace(vpcNamesByMachineID[item.ID])
		if vpcNames == "" {
			vpcNames = "-"
		}
		fmt.Fprintf(tw, "%s\t%s\t%s\t%s\t%s\n", item.Name, item.Status, siteName, vpcNames, item.ID)
	}
	return tw.Flush()
}

func cmdIPBlockList(s *Session, args []string) error {
	LogScopedCmd(s, "ip-block", "list")
	items, err := s.Resolver.Fetch(s.Ctx, "ip-block")
	if err != nil {
		return err
	}
	PrintSummary(os.Stdout, nil, len(items))
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

	LogScopedCmd(s, "ip-block", "create", "--name", name, "--site-id", site.ID, "--prefix", prefix, "--prefix-length", prefixLen)

	req := standard.NewIpBlockCreateRequest(name, site.ID, "DatacenterOnly", prefix, pl, "IPv4")

	block, _, err := s.Client.IPBlockAPI.CreateIpblock(s.Ctx, s.Org).IpBlockCreateRequest(*req).Execute()
	if err != nil {
		return fmt.Errorf("creating IP block: %w", err)
	}

	s.Cache.Invalidate("ip-block")
	fmt.Printf("%s IP block created: %s (%s)\n", Green("OK"), ptrStr(block.Name), ptrStr(block.Id))
	return nil
}

func cmdNSGList(s *Session, args []string) error {
	LogScopedCmd(s, "network-security-group", "list")
	items, err := s.Resolver.Fetch(s.Ctx, "network-security-group")
	if err != nil {
		return err
	}
	return printResourceTable(os.Stdout, "NAME", "STATUS", "ID", items)
}

func cmdSSHKeyList(s *Session, args []string) error {
	LogScopedCmd(s, "ssh-key", "list")
	items, err := s.Resolver.Fetch(s.Ctx, "ssh-key")
	if err != nil {
		return err
	}
	PrintSummary(os.Stdout, nil, len(items))
	tw := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
	fmt.Fprintln(tw, "NAME\tFINGERPRINT\tID")
	for _, item := range items {
		fmt.Fprintf(tw, "%s\t%s\t%s\n", item.Name, item.Extra["fingerprint"], item.ID)
	}
	return tw.Flush()
}

func cmdSKUList(s *Session, args []string) error {
	LogScopedCmd(s, "sku", "list")
	items, err := s.Resolver.Fetch(s.Ctx, "sku")
	if err != nil {
		return err
	}
	PrintSummary(os.Stdout, nil, len(items))
	tw := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
	fmt.Fprintln(tw, "DEVICE TYPE\tSITE\tID")
	for _, item := range items {
		fmt.Fprintf(tw, "%s\t%s\t%s\n", item.Extra["deviceType"], item.Extra["siteId"], item.ID)
	}
	return tw.Flush()
}

func cmdRackList(s *Session, args []string) error {
	LogScopedCmd(s, "rack", "list")
	items, err := s.Resolver.Fetch(s.Ctx, "rack")
	if err != nil {
		return err
	}
	PrintSummary(os.Stdout, nil, len(items))
	tw := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
	fmt.Fprintln(tw, "NAME\tMANUFACTURER\tMODEL\tID")
	for _, item := range items {
		fmt.Fprintf(tw, "%s\t%s\t%s\t%s\n", item.Name, item.Extra["manufacturer"], item.Extra["model"], item.ID)
	}
	return tw.Flush()
}

func cmdVPCPrefixList(s *Session, args []string) error {
	LogScopedCmd(s, "vpc-prefix", "list")
	items, err := s.Resolver.Fetch(s.Ctx, "vpc-prefix")
	if err != nil {
		return err
	}
	PrintSummary(os.Stdout, nil, len(items))
	tw := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
	fmt.Fprintln(tw, "NAME\tSTATUS\tVPC\tID")
	for _, item := range items {
		vpcName := s.Resolver.ResolveID("vpc", item.Extra["vpcId"])
		fmt.Fprintf(tw, "%s\t%s\t%s\t%s\n", item.Name, item.Status, vpcName, item.ID)
	}
	return tw.Flush()
}

func cmdTenantAccountList(s *Session, args []string) error {
	LogScopedCmd(s, "tenant-account", "list")
	items, err := s.Resolver.Fetch(s.Ctx, "tenant-account")
	if err != nil {
		return err
	}
	PrintSummary(os.Stdout, nil, len(items))
	tw := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
	fmt.Fprintln(tw, "TENANT ORG\tSTATUS\tINFRA PROVIDER ID\tID")
	for _, item := range items {
		fmt.Fprintf(tw, "%s\t%s\t%s\t%s\n", item.Name, item.Status, item.Extra["infrastructureProviderId"], item.ID)
	}
	return tw.Flush()
}

func cmdExpectedMachineList(s *Session, args []string) error {
	LogScopedCmd(s, "expected-machine", "list")
	items, err := s.Resolver.Fetch(s.Ctx, "expected-machine")
	if err != nil {
		return err
	}
	PrintSummary(os.Stdout, nil, len(items))
	tw := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
	fmt.Fprintln(tw, "SITE ID\tBMC MAC\tCHASSIS SN\tID")
	for _, item := range items {
		fmt.Fprintf(tw, "%s\t%s\t%s\t%s\n", item.Extra["siteId"], item.Extra["bmcMacAddress"], item.Extra["chassisSerialNumber"], item.ID)
	}
	return tw.Flush()
}

func cmdDPUExtensionServiceList(s *Session, args []string) error {
	LogScopedCmd(s, "dpu-extension-service", "list")
	items, err := s.Resolver.Fetch(s.Ctx, "dpu-extension-service")
	if err != nil {
		return err
	}
	PrintSummary(os.Stdout, nil, len(items))
	tw := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
	fmt.Fprintln(tw, "NAME\tTYPE\tSITE ID\tID")
	for _, item := range items {
		fmt.Fprintf(tw, "%s\t%s\t%s\t%s\n", item.Name, item.Extra["serviceType"], item.Extra["siteId"], item.ID)
	}
	return tw.Flush()
}

func cmdInfiniBandPartitionList(s *Session, args []string) error {
	LogScopedCmd(s, "infiniband-partition", "list")
	items, err := s.Resolver.Fetch(s.Ctx, "infiniband-partition")
	if err != nil {
		return err
	}
	PrintSummary(os.Stdout, nil, len(items))
	tw := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
	fmt.Fprintln(tw, "NAME\tSTATUS\tSITE ID\tID")
	for _, item := range items {
		fmt.Fprintf(tw, "%s\t%s\t%s\t%s\n", item.Name, item.Status, item.Extra["siteId"], item.ID)
	}
	return tw.Flush()
}

func cmdNVLinkLogicalPartitionList(s *Session, args []string) error {
	LogScopedCmd(s, "nvlink-logical-partition", "list")
	items, err := s.Resolver.Fetch(s.Ctx, "nvlink-logical-partition")
	if err != nil {
		return err
	}
	PrintSummary(os.Stdout, nil, len(items))
	tw := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
	fmt.Fprintln(tw, "NAME\tSTATUS\tSITE ID\tID")
	for _, item := range items {
		fmt.Fprintf(tw, "%s\t%s\t%s\t%s\n", item.Name, item.Status, item.Extra["siteId"], item.ID)
	}
	return tw.Flush()
}

func cmdAllocationConstraintList(s *Session, args []string) error {
	allocation, err := s.Resolver.ResolveWithArgs(s.Ctx, "allocation", "Allocation", args)
	if err != nil {
		return err
	}
	LogScopedCmd(s, "allocation-constraint", "list", "--allocation-id", allocation.ID)
	constraints, _, err := s.Client.AllocationAPI.GetAllAllocationConstraint(s.Ctx, s.Org, allocation.ID).Execute()
	if err != nil {
		return fmt.Errorf("listing allocation constraints: %w", err)
	}
	return printDetailJSON(os.Stdout, constraints)
}

func cmdAuditList(s *Session, args []string) error {
	LogScopedCmd(s, "audit", "list")
	items, err := s.Resolver.Fetch(s.Ctx, "audit")
	if err != nil {
		return err
	}
	PrintSummary(os.Stdout, nil, len(items))
	tw := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
	fmt.Fprintln(tw, "METHOD\tENDPOINT\tSTATUS CODE\tID")
	for _, item := range items {
		method := strings.TrimSpace(item.Extra["method"])
		if method == "" {
			method = item.Name
		}
		fmt.Fprintf(tw, "%s\t%s\t%s\t%s\n", method, item.Extra["endpoint"], item.Status, item.ID)
	}
	return tw.Flush()
}

func cmdMachineCapabilityList(s *Session, args []string) error {
	siteID := strings.TrimSpace(s.Scope.SiteID)
	siteName := s.Scope.SiteName

	// siteId is required by the API — prompt if not already in scope.
	if siteID == "" {
		site, err := s.Resolver.Resolve(s.Ctx, "site", "Site")
		if err != nil {
			return err
		}
		siteID = site.ID
		siteName = site.Name
	}

	LogScopedCmd(s, "machine-capability", "list", "--site-id", siteID)
	_ = siteName

	req := s.Client.MachineAPI.GetAllMachineCapabilities(s.Ctx, s.Org).SiteId(siteID)
	caps, _, err := req.Execute()
	if err != nil {
		if strings.Contains(err.Error(), "403") || strings.Contains(err.Error(), "Forbidden") {
			return fmt.Errorf("machine-capability list requires Infrastructure Provider access (FORGE_PROVIDER_ADMIN) — your org may not have that role")
		}
		return fmt.Errorf("listing machine capabilities: %w", err)
	}
	return printDetailJSON(os.Stdout, caps)
}

// -- Get/Details handlers --

func cmdSiteGet(s *Session, args []string) error {
	site, err := s.Resolver.ResolveWithArgs(s.Ctx, "site", "Site", args)
	if err != nil {
		return err
	}
	LogScopedCmd(s, "site", "get", site.ID)
	result, _, err := s.Client.SiteAPI.GetSite(s.Ctx, s.Org, site.ID).Execute()
	if err != nil {
		return fmt.Errorf("getting site: %w", err)
	}
	return printDetailJSON(os.Stdout, result)
}

func cmdVPCGet(s *Session, args []string) error {
	vpc, err := s.Resolver.ResolveWithArgs(s.Ctx, "vpc", "VPC", args)
	if err != nil {
		return err
	}
	LogScopedCmd(s, "vpc", "get", vpc.ID)
	result, _, err := s.Client.VPCAPI.GetVpc(s.Ctx, s.Org, vpc.ID).Execute()
	if err != nil {
		return fmt.Errorf("getting VPC: %w", err)
	}
	return printDetailJSON(os.Stdout, result)
}

func cmdSubnetGet(s *Session, args []string) error {
	subnet, err := s.Resolver.ResolveWithArgs(s.Ctx, "subnet", "Subnet", args)
	if err != nil {
		return err
	}
	LogScopedCmd(s, "subnet", "get", subnet.ID)
	result, _, err := s.Client.SubnetAPI.GetSubnet(s.Ctx, s.Org, subnet.ID).Execute()
	if err != nil {
		return fmt.Errorf("getting subnet: %w", err)
	}
	return printDetailJSON(os.Stdout, result)
}

func cmdInstanceGet(s *Session, args []string) error {
	instance, err := s.Resolver.ResolveWithArgs(s.Ctx, "instance", "Instance", args)
	if err != nil {
		return err
	}
	LogScopedCmd(s, "instance", "get", instance.ID)
	result, _, err := s.Client.InstanceAPI.GetInstance(s.Ctx, s.Org, instance.ID).Execute()
	if err != nil {
		return fmt.Errorf("getting instance: %w", err)
	}
	return printDetailJSON(os.Stdout, result)
}

func cmdInstanceTypeGet(s *Session, args []string) error {
	site, err := s.Resolver.Resolve(s.Ctx, "site", "Site")
	if err != nil {
		return err
	}
	items, err := s.fetchInstanceTypesBySite(s.Ctx, site.ID)
	if err != nil {
		return err
	}
	selected, err := s.Resolver.SelectFromItems("Instance Type", items)
	if err != nil {
		return err
	}
	LogScopedCmd(s, "instance-type", "get", selected.ID)
	result, _, err := s.Client.InstanceTypeAPI.GetInstanceType(s.Ctx, s.Org, selected.ID).Execute()
	if err != nil {
		return fmt.Errorf("getting instance type: %w", err)
	}
	return printDetailJSON(os.Stdout, result)
}

func cmdOSGet(s *Session, args []string) error {
	osItem, err := s.Resolver.ResolveWithArgs(s.Ctx, "operating-system", "Operating System", args)
	if err != nil {
		return err
	}
	LogScopedCmd(s, "operating-system", "get", osItem.ID)
	result, _, err := s.Client.OperatingSystemAPI.GetOperatingSystem(s.Ctx, s.Org, osItem.ID).Execute()
	if err != nil {
		return fmt.Errorf("getting operating system: %w", err)
	}
	return printDetailJSON(os.Stdout, result)
}

func cmdSSHKeyGroupGet(s *Session, args []string) error {
	group, err := s.Resolver.ResolveWithArgs(s.Ctx, "ssh-key-group", "SSH Key Group", args)
	if err != nil {
		return err
	}
	LogScopedCmd(s, "ssh-key-group", "get", group.ID)
	result, _, err := s.Client.SSHKeyGroupAPI.GetSshKeyGroup(s.Ctx, s.Org, group.ID).Execute()
	if err != nil {
		return fmt.Errorf("getting SSH key group: %w", err)
	}
	return printDetailJSON(os.Stdout, result)
}

func cmdAllocationGet(s *Session, args []string) error {
	alloc, err := s.Resolver.ResolveWithArgs(s.Ctx, "allocation", "Allocation", args)
	if err != nil {
		return err
	}
	LogScopedCmd(s, "allocation", "get", alloc.ID)
	result, _, err := s.Client.AllocationAPI.GetAllocation(s.Ctx, s.Org, alloc.ID).Execute()
	if err != nil {
		return fmt.Errorf("getting allocation: %w", err)
	}
	return printDetailJSON(os.Stdout, result)
}

func cmdMachineGet(s *Session, args []string) error {
	machine, err := s.Resolver.ResolveWithArgs(s.Ctx, "machine", "Machine", args)
	if err != nil {
		return err
	}
	LogScopedCmd(s, "machine", "get", machine.ID)
	result, _, err := s.Client.MachineAPI.GetMachine(s.Ctx, s.Org, machine.ID).Execute()
	if err != nil {
		return fmt.Errorf("getting machine: %w", err)
	}
	return printDetailJSON(os.Stdout, result)
}

func cmdIPBlockGet(s *Session, args []string) error {
	block, err := s.Resolver.ResolveWithArgs(s.Ctx, "ip-block", "IP Block", args)
	if err != nil {
		return err
	}
	LogScopedCmd(s, "ip-block", "get", block.ID)
	result, _, err := s.Client.IPBlockAPI.GetIpblock(s.Ctx, s.Org, block.ID).Execute()
	if err != nil {
		return fmt.Errorf("getting IP block: %w", err)
	}
	return printDetailJSON(os.Stdout, result)
}

func cmdNSGGet(s *Session, args []string) error {
	nsg, err := s.Resolver.ResolveWithArgs(s.Ctx, "network-security-group", "Network Security Group", args)
	if err != nil {
		return err
	}
	LogScopedCmd(s, "network-security-group", "get", nsg.ID)
	result, _, err := s.Client.NetworkSecurityGroupAPI.GetNetworkSecurityGroup(s.Ctx, s.Org, nsg.ID).Execute()
	if err != nil {
		return fmt.Errorf("getting network security group: %w", err)
	}
	return printDetailJSON(os.Stdout, result)
}

func cmdSSHKeyGet(s *Session, args []string) error {
	key, err := s.Resolver.ResolveWithArgs(s.Ctx, "ssh-key", "SSH Key", args)
	if err != nil {
		return err
	}
	LogScopedCmd(s, "ssh-key", "get", key.ID)
	result, _, err := s.Client.SSHKeyAPI.GetSshKey(s.Ctx, s.Org, key.ID).Execute()
	if err != nil {
		return fmt.Errorf("getting SSH key: %w", err)
	}
	return printDetailJSON(os.Stdout, result)
}

func cmdSKUGet(s *Session, args []string) error {
	sku, err := s.Resolver.ResolveWithArgs(s.Ctx, "sku", "SKU", args)
	if err != nil {
		return err
	}
	LogScopedCmd(s, "sku", "get", sku.ID)
	result, _, err := s.Client.SKUAPI.GetSku(s.Ctx, s.Org, sku.ID).Execute()
	if err != nil {
		return fmt.Errorf("getting SKU: %w", err)
	}
	return printDetailJSON(os.Stdout, result)
}

func cmdRackGet(s *Session, args []string) error {
	rack, err := s.Resolver.ResolveWithArgs(s.Ctx, "rack", "Rack", args)
	if err != nil {
		return err
	}
	LogScopedCmd(s, "rack", "get", rack.ID)
	result, _, err := s.Client.RackAPI.GetRack(s.Ctx, s.Org, rack.ID).Execute()
	if err != nil {
		return fmt.Errorf("getting rack: %w", err)
	}
	return printDetailJSON(os.Stdout, result)
}

func cmdVPCPrefixGet(s *Session, args []string) error {
	prefix, err := s.Resolver.ResolveWithArgs(s.Ctx, "vpc-prefix", "VPC Prefix", args)
	if err != nil {
		return err
	}
	LogScopedCmd(s, "vpc-prefix", "get", prefix.ID)
	result, _, err := s.Client.VPCPrefixAPI.GetVpcPrefix(s.Ctx, s.Org, prefix.ID).Execute()
	if err != nil {
		return fmt.Errorf("getting VPC prefix: %w", err)
	}
	return printDetailJSON(os.Stdout, result)
}

func cmdTenantAccountGet(s *Session, args []string) error {
	account, err := s.Resolver.ResolveWithArgs(s.Ctx, "tenant-account", "Tenant Account", args)
	if err != nil {
		return err
	}
	LogScopedCmd(s, "tenant-account", "get", account.ID)
	result, _, err := s.Client.TenantAccountAPI.GetTenantAccount(s.Ctx, s.Org, account.ID).Execute()
	if err != nil {
		return fmt.Errorf("getting tenant account: %w", err)
	}
	return printDetailJSON(os.Stdout, result)
}

func cmdExpectedMachineGet(s *Session, args []string) error {
	machine, err := s.Resolver.ResolveWithArgs(s.Ctx, "expected-machine", "Expected Machine", args)
	if err != nil {
		return err
	}
	LogScopedCmd(s, "expected-machine", "get", machine.ID)
	result, _, err := s.Client.ExpectedMachineAPI.GetExpectedMachine(s.Ctx, s.Org, machine.ID).Execute()
	if err != nil {
		return fmt.Errorf("getting expected machine: %w", err)
	}
	return printDetailJSON(os.Stdout, result)
}

func cmdDPUExtensionServiceGet(s *Session, args []string) error {
	service, err := s.Resolver.ResolveWithArgs(s.Ctx, "dpu-extension-service", "DPU Extension Service", args)
	if err != nil {
		return err
	}
	LogScopedCmd(s, "dpu-extension-service", "get", service.ID)
	result, _, err := s.Client.DPUExtensionServiceAPI.GetDpuExtensionService(s.Ctx, s.Org, service.ID).Execute()
	if err != nil {
		return fmt.Errorf("getting DPU extension service: %w", err)
	}
	return printDetailJSON(os.Stdout, result)
}

func cmdInfiniBandPartitionGet(s *Session, args []string) error {
	partition, err := s.Resolver.ResolveWithArgs(s.Ctx, "infiniband-partition", "InfiniBand Partition", args)
	if err != nil {
		return err
	}
	LogScopedCmd(s, "infiniband-partition", "get", partition.ID)
	result, _, err := s.Client.InfiniBandPartitionAPI.GetInfinibandPartition(s.Ctx, s.Org, partition.ID).Execute()
	if err != nil {
		return fmt.Errorf("getting InfiniBand partition: %w", err)
	}
	return printDetailJSON(os.Stdout, result)
}

func cmdNVLinkLogicalPartitionGet(s *Session, args []string) error {
	partition, err := s.Resolver.ResolveWithArgs(s.Ctx, "nvlink-logical-partition", "NVLink Logical Partition", args)
	if err != nil {
		return err
	}
	LogScopedCmd(s, "nvlink-logical-partition", "get", partition.ID)
	result, _, err := s.Client.NVLinkLogicalPartitionAPI.GetNvlinkLogicalPartition(s.Ctx, s.Org, partition.ID).Execute()
	if err != nil {
		return fmt.Errorf("getting NVLink logical partition: %w", err)
	}
	return printDetailJSON(os.Stdout, result)
}

func cmdAllocationConstraintGet(s *Session, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("usage: allocation-constraint get <allocation-id-or-name> [constraint-id]")
	}
	allocation, err := s.Resolver.ResolveWithArgs(s.Ctx, "allocation", "Allocation", args[:1])
	if err != nil {
		return err
	}

	constraintID := ""
	if len(args) > 1 {
		constraintID = strings.TrimSpace(args[1])
	}
	if constraintID == "" {
		constraints, _, err := s.Client.AllocationAPI.GetAllAllocationConstraint(s.Ctx, s.Org, allocation.ID).Execute()
		if err != nil {
			return fmt.Errorf("listing allocation constraints: %w", err)
		}
		if len(constraints) == 0 {
			return fmt.Errorf("no allocation constraints found for allocation %s", allocation.ID)
		}
		items := make([]NamedItem, 0, len(constraints))
		for _, constraint := range constraints {
			id := strings.TrimSpace(ptrStr(constraint.Id))
			if id == "" {
				continue
			}
			items = append(items, NamedItem{Name: id, ID: id, Raw: constraint})
		}
		if len(items) == 0 {
			return fmt.Errorf("allocation constraints did not include IDs for allocation %s", allocation.ID)
		}
		selected, err := s.Resolver.SelectFromItems("Allocation Constraint", items)
		if err != nil {
			return err
		}
		constraintID = selected.ID
	}

	LogScopedCmd(s, "allocation-constraint", "get", constraintID, "--allocation-id", allocation.ID)
	result, _, err := s.Client.AllocationAPI.GetAllocationConstraint(s.Ctx, s.Org, allocation.ID, constraintID).Execute()
	if err != nil {
		return fmt.Errorf("getting allocation constraint: %w", err)
	}
	return printDetailJSON(os.Stdout, result)
}

func cmdAuditGet(s *Session, args []string) error {
	entry, err := s.Resolver.ResolveWithArgs(s.Ctx, "audit", "Audit Entry", args)
	if err != nil {
		return err
	}
	LogScopedCmd(s, "audit", "get", entry.ID)
	result, _, err := s.Client.AuditAPI.GetAuditEntry(s.Ctx, s.Org, entry.ID).Execute()
	if err != nil {
		return fmt.Errorf("getting audit entry: %w", err)
	}
	return printDetailJSON(os.Stdout, result)
}

func cmdMetadataGet(s *Session, args []string) error {
	LogScopedCmd(s, "metadata", "get")
	result, _, err := s.Client.MetadataAPI.GetMetadata(s.Ctx, s.Org).Execute()
	if err != nil {
		return fmt.Errorf("getting metadata: %w", err)
	}
	return printDetailJSON(os.Stdout, result)
}

func cmdUserCurrent(s *Session, args []string) error {
	LogScopedCmd(s, "user", "current")
	result, _, err := s.Client.UserAPI.GetUser(s.Ctx, s.Org).Execute()
	if err != nil {
		return fmt.Errorf("getting current user: %w", err)
	}
	return printDetailJSON(os.Stdout, result)
}

func cmdServiceAccountCurrent(s *Session, args []string) error {
	LogScopedCmd(s, "service-account", "current")
	result, _, err := s.Client.ServiceAccountAPI.GetCurrentServiceAccount(s.Ctx, s.Org).Execute()
	if err != nil {
		return fmt.Errorf("getting current service account: %w", err)
	}
	return printDetailJSON(os.Stdout, result)
}

func cmdInfrastructureProviderCurrent(s *Session, args []string) error {
	LogScopedCmd(s, "infrastructure-provider", "current")
	result, _, err := s.Client.InfrastructureProviderAPI.GetCurrentInfrastructureProvider(s.Ctx, s.Org).Execute()
	if err != nil {
		return fmt.Errorf("getting current infrastructure provider: %w", err)
	}
	return printDetailJSON(os.Stdout, result)
}

func cmdInfrastructureProviderStats(s *Session, args []string) error {
	LogScopedCmd(s, "infrastructure-provider", "stats")
	result, _, err := s.Client.InfrastructureProviderAPI.GetCurrentInfrastructureProviderStats(s.Ctx, s.Org).Execute()
	if err != nil {
		return fmt.Errorf("getting infrastructure provider stats: %w", err)
	}
	return printDetailJSON(os.Stdout, result)
}

func cmdTenantCurrent(s *Session, args []string) error {
	LogScopedCmd(s, "tenant", "current")
	result, _, err := s.Client.TenantAPI.GetCurrentTenant(s.Ctx, s.Org).Execute()
	if err != nil {
		return fmt.Errorf("getting current tenant: %w", err)
	}
	return printDetailJSON(os.Stdout, result)
}

func cmdTenantStats(s *Session, args []string) error {
	LogScopedCmd(s, "tenant", "stats")
	result, _, err := s.Client.TenantAPI.GetCurrentTenantStats(s.Ctx, s.Org).Execute()
	if err != nil {
		return fmt.Errorf("getting tenant stats: %w", err)
	}
	return printDetailJSON(os.Stdout, result)
}

// -- Delete handlers --

func cmdSiteDelete(s *Session, args []string) error {
	site, err := s.Resolver.ResolveWithArgs(s.Ctx, "site", "Site to delete", args)
	if err != nil {
		return err
	}
	ok, err := PromptConfirm(fmt.Sprintf("Delete site %s (%s)?", site.Name, site.ID))
	if err != nil || !ok {
		return err
	}
	LogScopedCmd(s, "site", "delete", site.ID)
	_, err = s.Client.SiteAPI.DeleteSite(s.Ctx, s.Org, site.ID).Execute()
	if err != nil {
		return fmt.Errorf("deleting site: %w", err)
	}
	s.Cache.Invalidate("site")
	fmt.Printf("%s Site deleted: %s\n", Green("OK"), site.Name)
	return nil
}

func cmdVPCDelete(s *Session, args []string) error {
	vpc, err := s.Resolver.ResolveWithArgs(s.Ctx, "vpc", "VPC to delete", args)
	if err != nil {
		return err
	}
	ok, err := PromptConfirm(fmt.Sprintf("Delete VPC %s (%s)?", vpc.Name, vpc.ID))
	if err != nil || !ok {
		return err
	}
	LogScopedCmd(s, "vpc", "delete", vpc.ID)
	_, err = s.Client.VPCAPI.DeleteVpc(s.Ctx, s.Org, vpc.ID).Execute()
	if err != nil {
		return fmt.Errorf("deleting VPC: %w", err)
	}
	s.Cache.Invalidate("vpc")
	fmt.Printf("%s VPC deleted: %s\n", Green("OK"), vpc.Name)
	return nil
}

func cmdSubnetDelete(s *Session, args []string) error {
	subnet, err := s.Resolver.ResolveWithArgs(s.Ctx, "subnet", "Subnet to delete", args)
	if err != nil {
		return err
	}
	ok, err := PromptConfirm(fmt.Sprintf("Delete subnet %s (%s)?", subnet.Name, subnet.ID))
	if err != nil || !ok {
		return err
	}
	LogScopedCmd(s, "subnet", "delete", subnet.ID)
	_, err = s.Client.SubnetAPI.DeleteSubnet(s.Ctx, s.Org, subnet.ID).Execute()
	if err != nil {
		return fmt.Errorf("deleting subnet: %w", err)
	}
	s.Cache.Invalidate("subnet")
	fmt.Printf("%s Subnet deleted: %s\n", Green("OK"), subnet.Name)
	return nil
}

func cmdInstanceDelete(s *Session, args []string) error {
	instance, err := s.Resolver.ResolveWithArgs(s.Ctx, "instance", "Instance to delete", args)
	if err != nil {
		return err
	}
	ok, err := PromptConfirm(fmt.Sprintf("Delete instance %s (%s)?", instance.Name, instance.ID))
	if err != nil || !ok {
		return err
	}
	LogScopedCmd(s, "instance", "delete", instance.ID)
	_, err = s.Client.InstanceAPI.DeleteInstance(s.Ctx, s.Org, instance.ID).Execute()
	if err != nil {
		return fmt.Errorf("deleting instance: %w", err)
	}
	s.Cache.Invalidate("instance")
	fmt.Printf("%s Instance deleted: %s\n", Green("OK"), instance.Name)
	return nil
}

func cmdOSDelete(s *Session, args []string) error {
	osItem, err := s.Resolver.ResolveWithArgs(s.Ctx, "operating-system", "Operating System to delete", args)
	if err != nil {
		return err
	}
	ok, err := PromptConfirm(fmt.Sprintf("Delete operating system %s (%s)?", osItem.Name, osItem.ID))
	if err != nil || !ok {
		return err
	}
	LogScopedCmd(s, "operating-system", "delete", osItem.ID)
	_, err = s.Client.OperatingSystemAPI.DeleteOperatingSystem(s.Ctx, s.Org, osItem.ID).Execute()
	if err != nil {
		return fmt.Errorf("deleting operating system: %w", err)
	}
	s.Cache.Invalidate("operating-system")
	fmt.Printf("%s Operating system deleted: %s\n", Green("OK"), osItem.Name)
	return nil
}

func cmdInstanceTypeDelete(s *Session, args []string) error {
	site, err := s.Resolver.Resolve(s.Ctx, "site", "Site")
	if err != nil {
		return err
	}
	items, err := s.fetchInstanceTypesBySite(s.Ctx, site.ID)
	if err != nil {
		return err
	}
	selected, err := s.Resolver.SelectFromItems("Instance Type to delete", items)
	if err != nil {
		return err
	}
	ok, err := PromptConfirm(fmt.Sprintf("Delete instance type %s (%s)?", selected.Name, selected.ID))
	if err != nil || !ok {
		return err
	}
	LogScopedCmd(s, "instance-type", "delete", selected.ID)
	_, err = s.Client.InstanceTypeAPI.DeleteInstanceType(s.Ctx, s.Org, selected.ID).Execute()
	if err != nil {
		return fmt.Errorf("deleting instance type: %w", err)
	}
	s.Cache.Invalidate("instance-type")
	fmt.Printf("%s Instance type deleted: %s\n", Green("OK"), selected.Name)
	return nil
}

func cmdAllocationDelete(s *Session, args []string) error {
	alloc, err := s.Resolver.ResolveWithArgs(s.Ctx, "allocation", "Allocation to delete", args)
	if err != nil {
		return err
	}
	ok, err := PromptConfirm(fmt.Sprintf("Delete allocation %s (%s)?", alloc.Name, alloc.ID))
	if err != nil || !ok {
		return err
	}
	LogScopedCmd(s, "allocation", "delete", alloc.ID)
	_, err = s.Client.AllocationAPI.DeleteAllocation(s.Ctx, s.Org, alloc.ID).Execute()
	if err != nil {
		return fmt.Errorf("deleting allocation: %w", err)
	}
	s.Cache.Invalidate("allocation")
	fmt.Printf("%s Allocation deleted: %s\n", Green("OK"), alloc.Name)
	return nil
}

func cmdIPBlockDelete(s *Session, args []string) error {
	block, err := s.Resolver.ResolveWithArgs(s.Ctx, "ip-block", "IP Block to delete", args)
	if err != nil {
		return err
	}
	ok, err := PromptConfirm(fmt.Sprintf("Delete IP block %s (%s)?", block.Name, block.ID))
	if err != nil || !ok {
		return err
	}
	LogScopedCmd(s, "ip-block", "delete", block.ID)
	_, err = s.Client.IPBlockAPI.DeleteIpblock(s.Ctx, s.Org, block.ID).Execute()
	if err != nil {
		return fmt.Errorf("deleting IP block: %w", err)
	}
	s.Cache.Invalidate("ip-block")
	fmt.Printf("%s IP block deleted: %s\n", Green("OK"), block.Name)
	return nil
}

func cmdLogin(s *Session, args []string) error {
	if s.LoginFn == nil {
		return fmt.Errorf("login not available (no auth configured)")
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
	fmt.Fprintln(tw, "org\tShow current org")
	fmt.Fprintln(tw, "org list\tList available orgs")
	fmt.Fprintln(tw, "org set <name>\tSwitch to a different org")
	fmt.Fprintln(tw, "scope\tShow current scope filters")
	fmt.Fprintln(tw, "scope site [name]\tSet site scope (filters lists)")
	fmt.Fprintln(tw, "scope vpc [name]\tSet VPC scope (filters lists)")
	fmt.Fprintln(tw, "scope clear\tClear all scope filters")
	fmt.Fprintln(tw, "exit\tExit interactive mode")
	tw.Flush()
	fmt.Println()
	return nil
}

// -- Helpers --

func printDetailJSON(w io.Writer, v interface{}) error {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return fmt.Errorf("marshaling output: %v", err)
	}
	fmt.Fprintln(w, string(data))
	return nil
}

func printResourceTable(w io.Writer, col1, col2, col3 string, items []NamedItem) error {
	PrintSummary(w, nil, len(items))
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

