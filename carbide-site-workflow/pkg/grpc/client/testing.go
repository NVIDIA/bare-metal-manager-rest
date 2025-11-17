/*
 * SPDX-FileCopyrightText: Copyright (c) 2021-2023 NVIDIA CORPORATION & AFFILIATES. All rights reserved.
 * SPDX-License-Identifier: LicenseRef-NvidiaProprietary
 *
 * NVIDIA CORPORATION, its affiliates and licensors retain all intellectual
 * property and proprietary rights in and to this material, related
 * documentation and any modifications thereto. Any use, reproduction,
 * disclosure or distribution of this material and related documentation
 * without an express license agreement from NVIDIA CORPORATION or
 * its affiliates is strictly prohibited.
 */

package client

import (
	"context"
	"math/rand"
	"net"

	"github.com/gogo/status"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/protobuf/types/known/emptypb"

	wflows "github.com/NVIDIA/carbide-rest-api/carbide-rest-api-schema/schema/site-agent/workflows/v1"
)

var runes = []rune("abcdefghijklmnopqrstuvwxyz0123456789")

// Add utlity methods here
// randSeq generates a random sequence of runes
func randSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = runes[rand.Intn(len(runes))]
	}
	return string(b)
}

// incrementMAC takes a hardware address (MAC address) and increments it by one.
// It handles carrying over to the next byte when a byte overflows (reaches 255).
func incrementMAC(mac net.HardwareAddr) {
	// Iterate from the last byte to the first.
	for i := len(mac) - 1; i >= 0; i-- {
		// Increment the current byte.
		mac[i]++
		// If the byte is not 0, it means there was no overflow, so we can stop.
		if mac[i] != 0 {
			break
		}
		// If the byte is 0, it means it overflowed from 255, so we continue to the next
		// byte to handle the "carry-over".
	}
}

// MockForgeClient is a mock implementation of Forge gRPC protobuf Client
type MockForgeClient struct {
	wflows.ForgeClient
}

func (c *MockForgeClient) Version(ctx context.Context, in *wflows.VersionRequest, opts ...grpc.CallOption) (*wflows.BuildInfo, error) {
	out := new(wflows.BuildInfo)
	return out, nil
}

func (c *MockForgeClient) CreateDomain(ctx context.Context, in *wflows.Domain, opts ...grpc.CallOption) (*wflows.Domain, error) {
	out := new(wflows.Domain)
	return out, nil
}

func (c *MockForgeClient) UpdateDomain(ctx context.Context, in *wflows.Domain, opts ...grpc.CallOption) (*wflows.Domain, error) {
	out := new(wflows.Domain)
	return out, nil
}

func (c *MockForgeClient) DeleteDomain(ctx context.Context, in *wflows.DomainDeletion, opts ...grpc.CallOption) (*wflows.DomainDeletionResult, error) {
	out := new(wflows.DomainDeletionResult)
	return out, nil
}

func (c *MockForgeClient) FindDomain(ctx context.Context, in *wflows.DomainSearchQuery, opts ...grpc.CallOption) (*wflows.DomainList, error) {
	out := new(wflows.DomainList)
	return out, nil
}

func (c *MockForgeClient) CreateVpc(ctx context.Context, in *wflows.VpcCreationRequest, opts ...grpc.CallOption) (*wflows.Vpc, error) {
	out := new(wflows.Vpc)
	return out, nil
}

func (c *MockForgeClient) UpdateVpc(ctx context.Context, in *wflows.VpcUpdateRequest, opts ...grpc.CallOption) (*wflows.VpcUpdateResult, error) {
	out := new(wflows.VpcUpdateResult)
	return out, nil
}

func (c *MockForgeClient) UpdateVpcVirtualization(ctx context.Context, in *wflows.VpcUpdateVirtualizationRequest, opts ...grpc.CallOption) (*wflows.VpcUpdateVirtualizationResult, error) {
	out := new(wflows.VpcUpdateVirtualizationResult)
	return out, nil
}

func (c *MockForgeClient) DeleteVpc(ctx context.Context, in *wflows.VpcDeletionRequest, opts ...grpc.CallOption) (*wflows.VpcDeletionResult, error) {
	out := new(wflows.VpcDeletionResult)
	return out, nil
}

func (c *MockForgeClient) FindVpcIds(ctx context.Context, in *wflows.VpcSearchFilter, opts ...grpc.CallOption) (*wflows.VpcIdList, error) {
	if err, ok := ctx.Value("wantError").(error); ok {
		return nil, status.Error(status.Code(err), "failed to retrieve vpc ids")
	}

	out := &wflows.VpcIdList{}

	count, ok := ctx.Value("wantCount").(int)
	if ok {
		for i := 0; i < count; i++ {
			out.VpcIds = append(out.VpcIds, &wflows.UUID{Value: uuid.NewString()})
		}
	}

	return out, nil
}

func (c *MockForgeClient) FindVpcsByIds(ctx context.Context, in *wflows.VpcsByIdsRequest, opts ...grpc.CallOption) (*wflows.VpcList, error) {
	err, ok := ctx.Value("wantError").(error)
	if ok {
		return nil, status.Error(status.Code(err), "failed to retrieve vpcs")
	}

	out := &wflows.VpcList{}
	if in != nil {
		for _, id := range in.VpcIds {
			out.Vpcs = append(out.Vpcs, &wflows.Vpc{
				Id: id,
			})
		}
	}

	return out, nil
}

// DEPRECATED: use FindVpcIds and FindVpcsByIds instead
func (c *MockForgeClient) FindVpcs(ctx context.Context, in *wflows.VpcSearchQuery, opts ...grpc.CallOption) (*wflows.VpcList, error) {
	out := new(wflows.VpcList)
	return out, nil
}

func (c *MockForgeClient) FindNetworkSegmentIds(ctx context.Context, in *wflows.NetworkSegmentSearchFilter, opts ...grpc.CallOption) (*wflows.NetworkSegmentIdList, error) {
	if err, ok := ctx.Value("wantError").(error); ok {
		return nil, status.Error(status.Code(err), "failed to retrieve network segment ids")
	}

	out := &wflows.NetworkSegmentIdList{}

	count, ok := ctx.Value("wantCount").(int)
	if ok {
		for i := 0; i < count; i++ {
			out.NetworkSegmentsIds = append(out.NetworkSegmentsIds, &wflows.UUID{Value: uuid.NewString()})
		}
	}

	return out, nil
}

func (c *MockForgeClient) FindNetworkSegmentsByIds(ctx context.Context, in *wflows.NetworkSegmentsByIdsRequest, opts ...grpc.CallOption) (*wflows.NetworkSegmentList, error) {
	err, ok := ctx.Value("wantError").(error)
	if ok {
		return nil, status.Error(status.Code(err), "failed to retrieve network segments")
	}

	out := &wflows.NetworkSegmentList{}
	if in != nil {
		for _, id := range in.NetworkSegmentsIds {
			out.NetworkSegments = append(out.NetworkSegments, &wflows.NetworkSegment{
				Id: id,
			})
		}
	}

	return out, nil
}

// DEPRECATED: use FindNetworkSegmentIDs and FindNetworkSegmentsByIDs instead
func (c *MockForgeClient) FindNetworkSegments(ctx context.Context, in *wflows.NetworkSegmentQuery, opts ...grpc.CallOption) (*wflows.NetworkSegmentList, error) {
	out := &wflows.NetworkSegmentList{}
	count, ok := ctx.Value("wantCount").(int)
	if ok {
		for i := 0; i < count; i++ {
			out.NetworkSegments = append(out.NetworkSegments, &wflows.NetworkSegment{Id: &wflows.UUID{Value: uuid.NewString()}})
		}
	}
	return out, nil
}

func (c *MockForgeClient) CreateNetworkSegment(ctx context.Context, in *wflows.NetworkSegmentCreationRequest, opts ...grpc.CallOption) (*wflows.NetworkSegment, error) {
	out := new(wflows.NetworkSegment)
	return out, nil
}

func (c *MockForgeClient) DeleteNetworkSegment(ctx context.Context, in *wflows.NetworkSegmentDeletionRequest, opts ...grpc.CallOption) (*wflows.NetworkSegmentDeletionResult, error) {
	out := new(wflows.NetworkSegmentDeletionResult)
	return out, nil
}

func (c *MockForgeClient) NetworkSegmentsForVpc(ctx context.Context, in *wflows.VpcSearchQuery, opts ...grpc.CallOption) (*wflows.NetworkSegmentList, error) {
	out := new(wflows.NetworkSegmentList)
	return out, nil
}

func (c *MockForgeClient) FindIBPartitionIds(ctx context.Context, in *wflows.IBPartitionSearchFilter, opts ...grpc.CallOption) (*wflows.IBPartitionIdList, error) {
	if err, ok := ctx.Value("wantError").(error); ok {
		return nil, status.Error(status.Code(err), "failed to retrieve ib partition ids")
	}

	out := &wflows.IBPartitionIdList{}

	count, ok := ctx.Value("wantCount").(int)
	if ok {
		for i := 0; i < count; i++ {
			out.IbPartitionIds = append(out.IbPartitionIds, &wflows.UUID{Value: uuid.NewString()})
		}
	}

	return out, nil
}

func (c *MockForgeClient) FindIBPartitionsByIds(ctx context.Context, in *wflows.IBPartitionsByIdsRequest, opts ...grpc.CallOption) (*wflows.IBPartitionList, error) {
	err, ok := ctx.Value("wantError").(error)
	if ok {
		return nil, status.Error(status.Code(err), "failed to retrieve ib partitions")
	}

	out := &wflows.IBPartitionList{}
	if in != nil {
		for _, id := range in.IbPartitionIds {
			out.IbPartitions = append(out.IbPartitions, &wflows.IBPartition{
				Id: id,
			})
		}
	}

	return out, nil
}

// DEPRECATED: use FindIBPartitionIds and FindIBPartitionsByIds instead
func (c *MockForgeClient) FindIBPartitions(ctx context.Context, in *wflows.IBPartitionQuery, opts ...grpc.CallOption) (*wflows.IBPartitionList, error) {
	out := &wflows.IBPartitionList{}
	count, ok := ctx.Value("wantCount").(int)
	if ok {
		for i := 0; i < count; i++ {
			out.IbPartitions = append(out.IbPartitions, &wflows.IBPartition{Id: &wflows.UUID{Value: uuid.NewString()}})
		}
	}
	return out, nil
}

func (c *MockForgeClient) CreateIBPartition(ctx context.Context, in *wflows.IBPartitionCreationRequest, opts ...grpc.CallOption) (*wflows.IBPartition, error) {
	out := new(wflows.IBPartition)
	return out, nil
}

func (c *MockForgeClient) DeleteIBPartition(ctx context.Context, in *wflows.IBPartitionDeletionRequest, opts ...grpc.CallOption) (*wflows.IBPartitionDeletionResult, error) {
	out := new(wflows.IBPartitionDeletionResult)
	return out, nil
}

func (c *MockForgeClient) IBPartitionsForTenant(ctx context.Context, in *wflows.TenantSearchQuery, opts ...grpc.CallOption) (*wflows.IBPartitionList, error) {
	out := new(wflows.IBPartitionList)
	return out, nil
}

func (c *MockForgeClient) AllocateInstance(ctx context.Context, in *wflows.InstanceAllocationRequest, opts ...grpc.CallOption) (*wflows.Instance, error) {
	out := new(wflows.Instance)
	return out, nil
}

func (c *MockForgeClient) ReleaseInstance(ctx context.Context, in *wflows.InstanceReleaseRequest, opts ...grpc.CallOption) (*wflows.InstanceReleaseResult, error) {
	err, ok := ctx.Value("wantError").(error)
	if ok {
		if status.Code(err) == codes.NotFound {
			return nil, status.Error(codes.NotFound, "instance not found: ")
		}
	}
	out := new(wflows.InstanceReleaseResult)
	return out, nil
}

func (c *MockForgeClient) FindInstanceIds(ctx context.Context, in *wflows.InstanceSearchFilter, opts ...grpc.CallOption) (*wflows.InstanceIdList, error) {
	if err, ok := ctx.Value("wantError").(error); ok {
		return nil, status.Error(status.Code(err), "failed to retrieve instance ids")
	}

	out := &wflows.InstanceIdList{}

	count, ok := ctx.Value("wantCount").(int)
	if ok {
		for i := 0; i < count; i++ {
			out.InstanceIds = append(out.InstanceIds, &wflows.UUID{Value: uuid.NewString()})
		}
	}

	return out, nil
}

func (c *MockForgeClient) FindInstancesByIds(ctx context.Context, in *wflows.InstancesByIdsRequest, opts ...grpc.CallOption) (*wflows.InstanceList, error) {
	err, ok := ctx.Value("wantError").(error)
	if ok {
		return nil, status.Error(status.Code(err), "failed to retrieve instances")
	}

	out := &wflows.InstanceList{}
	if in != nil {
		for _, id := range in.InstanceIds {
			out.Instances = append(out.Instances, &wflows.Instance{
				Id: id,
			})
		}
	}

	return out, nil
}

// DEPRECATED: use FindInstanceIds and FindInstancesByIds instead
func (c *MockForgeClient) FindInstances(ctx context.Context, in *wflows.InstanceSearchQuery, opts ...grpc.CallOption) (*wflows.InstanceList, error) {
	out := new(wflows.InstanceList)
	return out, nil
}

func (c *MockForgeClient) FindInstanceByMachineID(ctx context.Context, in *wflows.MachineId, opts ...grpc.CallOption) (*wflows.InstanceList, error) {
	out := new(wflows.InstanceList)
	return out, nil
}

func (c *MockForgeClient) GetManagedHostNetworkConfig(ctx context.Context, in *wflows.ManagedHostNetworkConfigRequest, opts ...grpc.CallOption) (*wflows.ManagedHostNetworkConfigResponse, error) {
	out := new(wflows.ManagedHostNetworkConfigResponse)
	return out, nil
}

func (c *MockForgeClient) RecordDpuNetworkStatus(ctx context.Context, in *wflows.DpuNetworkStatus, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	return out, nil
}

func (c *MockForgeClient) DpuAgentUpgradeCheck(ctx context.Context, in *wflows.DpuAgentUpgradeCheckRequest, opts ...grpc.CallOption) (*wflows.DpuAgentUpgradeCheckResponse, error) {
	out := new(wflows.DpuAgentUpgradeCheckResponse)
	return out, nil
}

func (c *MockForgeClient) DpuAgentUpgradePolicyAction(ctx context.Context, in *wflows.DpuAgentUpgradePolicyRequest, opts ...grpc.CallOption) (*wflows.DpuAgentUpgradePolicyResponse, error) {
	out := new(wflows.DpuAgentUpgradePolicyResponse)
	return out, nil
}

func (c *MockForgeClient) LookupRecord(ctx context.Context, in *wflows.DNSMessage_DNSQuestion, opts ...grpc.CallOption) (*wflows.DNSMessage_DNSResponse, error) {
	out := new(wflows.DNSMessage_DNSResponse)
	return out, nil
}

func (c *MockForgeClient) InvokeInstancePower(ctx context.Context, in *wflows.InstancePowerRequest, opts ...grpc.CallOption) (*wflows.InstancePowerResult, error) {
	out := new(wflows.InstancePowerResult)
	return out, nil
}

func (c *MockForgeClient) ForgeAgentControl(ctx context.Context, in *wflows.ForgeAgentControlRequest, opts ...grpc.CallOption) (*wflows.ForgeAgentControlResponse, error) {
	out := new(wflows.ForgeAgentControlResponse)
	return out, nil
}

func (c *MockForgeClient) DiscoverMachine(ctx context.Context, in *wflows.MachineDiscoveryInfo, opts ...grpc.CallOption) (*wflows.MachineDiscoveryResult, error) {
	out := new(wflows.MachineDiscoveryResult)
	return out, nil
}

func (c *MockForgeClient) RenewMachineCertificate(ctx context.Context, in *wflows.MachineCertificateRenewRequest, opts ...grpc.CallOption) (*wflows.MachineCertificateResult, error) {
	out := new(wflows.MachineCertificateResult)
	return out, nil
}

func (c *MockForgeClient) DiscoveryCompleted(ctx context.Context, in *wflows.MachineDiscoveryCompletedRequest, opts ...grpc.CallOption) (*wflows.MachineDiscoveryCompletedResponse, error) {
	out := new(wflows.MachineDiscoveryCompletedResponse)
	return out, nil
}

func (c *MockForgeClient) CleanupMachineCompleted(ctx context.Context, in *wflows.MachineCleanupInfo, opts ...grpc.CallOption) (*wflows.MachineCleanupResult, error) {
	out := new(wflows.MachineCleanupResult)
	return out, nil
}

func (c *MockForgeClient) ReportForgeScoutError(ctx context.Context, in *wflows.ForgeScoutErrorReport, opts ...grpc.CallOption) (*wflows.ForgeScoutErrorReportResult, error) {
	out := new(wflows.ForgeScoutErrorReportResult)
	return out, nil
}

func (c *MockForgeClient) DiscoverDhcp(ctx context.Context, in *wflows.DhcpDiscovery, opts ...grpc.CallOption) (*wflows.DhcpRecord, error) {
	out := new(wflows.DhcpRecord)
	return out, nil
}

func (c *MockForgeClient) GetMachine(ctx context.Context, in *wflows.MachineId, opts ...grpc.CallOption) (*wflows.Machine, error) {
	out := new(wflows.Machine)
	return out, nil
}

func (c *MockForgeClient) FindMachines(ctx context.Context, in *wflows.MachineSearchQuery, opts ...grpc.CallOption) (*wflows.MachineList, error) {
	out := new(wflows.MachineList)
	return out, nil
}

func (c *MockForgeClient) FindInterfaces(ctx context.Context, in *wflows.InterfaceSearchQuery, opts ...grpc.CallOption) (*wflows.InterfaceList, error) {
	out := new(wflows.InterfaceList)
	return out, nil
}

func (c *MockForgeClient) FindIpAddress(ctx context.Context, in *wflows.FindIpAddressRequest, opts ...grpc.CallOption) (*wflows.FindIpAddressResponse, error) {
	out := new(wflows.FindIpAddressResponse)
	return out, nil
}

func (c *MockForgeClient) FindMachineIds(ctx context.Context, in *wflows.MachineSearchConfig, opts ...grpc.CallOption) (*wflows.MachineIdList, error) {
	err, ok := ctx.Value("wantError").(error)
	if ok {
		if status.Code(err) == codes.Internal {
			return nil, status.Error(codes.Internal, "failed to retrieve machine ids")
		}
	}

	out := &wflows.MachineIdList{}

	count, ok := ctx.Value("wantCount").(int)
	if ok {
		for i := 0; i < count; i++ {
			out.MachineIds = append(out.MachineIds, &wflows.MachineId{Id: uuid.NewString()})
		}
	}

	return out, nil
}

func (c *MockForgeClient) FindMachinesByIds(ctx context.Context, in *wflows.MachinesByIdsRequest, opts ...grpc.CallOption) (*wflows.MachineList, error) {
	err, ok := ctx.Value("wantError").(error)
	if ok {
		if status.Code(err) == codes.Internal {
			return nil, status.Error(codes.Internal, "failed to retrieve machines by ids")
		}
	}

	out := &wflows.MachineList{}
	if in != nil {
		for _, id := range in.MachineIds {
			out.Machines = append(out.Machines, &wflows.Machine{
				Id:    id,
				State: "Ready",
			})
		}
	}

	return out, nil
}

func (c *MockForgeClient) GetBMCMetaData(ctx context.Context, in *wflows.BMCMetaDataGetRequest, opts ...grpc.CallOption) (*wflows.BMCMetaDataGetResponse, error) {
	out := new(wflows.BMCMetaDataGetResponse)
	return out, nil
}

func (c *MockForgeClient) UpdateMachineCredentials(ctx context.Context, in *wflows.MachineCredentialsUpdateRequest, opts ...grpc.CallOption) (*wflows.MachineCredentialsUpdateResponse, error) {
	out := new(wflows.MachineCredentialsUpdateResponse)
	return out, nil
}

func (c *MockForgeClient) GetPxeInstructions(ctx context.Context, in *wflows.PxeInstructionRequest, opts ...grpc.CallOption) (*wflows.PxeInstructions, error) {
	out := new(wflows.PxeInstructions)
	return out, nil
}

func (c *MockForgeClient) GetCloudInitInstructions(ctx context.Context, in *wflows.CloudInitInstructionsRequest, opts ...grpc.CallOption) (*wflows.CloudInitInstructions, error) {
	out := new(wflows.CloudInitInstructions)
	return out, nil
}

func (c *MockForgeClient) Echo(ctx context.Context, in *wflows.EchoRequest, opts ...grpc.CallOption) (*wflows.EchoResponse, error) {
	out := new(wflows.EchoResponse)
	return out, nil
}

func (c *MockForgeClient) CreateTenant(ctx context.Context, in *wflows.CreateTenantRequest, opts ...grpc.CallOption) (*wflows.CreateTenantResponse, error) {
	out := new(wflows.CreateTenantResponse)
	return out, nil
}

func (c *MockForgeClient) FindTenant(ctx context.Context, in *wflows.FindTenantRequest, opts ...grpc.CallOption) (*wflows.FindTenantResponse, error) {
	out := new(wflows.FindTenantResponse)
	return out, nil
}

func (c *MockForgeClient) UpdateTenant(ctx context.Context, in *wflows.UpdateTenantRequest, opts ...grpc.CallOption) (*wflows.UpdateTenantResponse, error) {
	out := new(wflows.UpdateTenantResponse)
	return out, nil
}

func (c *MockForgeClient) CreateTenantKeyset(ctx context.Context, in *wflows.CreateTenantKeysetRequest, opts ...grpc.CallOption) (*wflows.CreateTenantKeysetResponse, error) {
	out := new(wflows.CreateTenantKeysetResponse)
	return out, nil
}

func (c *MockForgeClient) FindTenantKeysetIds(ctx context.Context, in *wflows.TenantKeysetSearchFilter, opts ...grpc.CallOption) (*wflows.TenantKeysetIdList, error) {
	if err, ok := ctx.Value("wantError").(error); ok {
		return nil, status.Error(status.Code(err), "failed to retrieve tenant keyset ids")
	}

	out := &wflows.TenantKeysetIdList{}

	count, ok := ctx.Value("wantCount").(int)
	if ok {
		orgID := uuid.NewString()
		for i := 0; i < count; i++ {
			out.KeysetIds = append(out.KeysetIds, &wflows.TenantKeysetIdentifier{OrganizationId: orgID, KeysetId: uuid.NewString()})
		}
	}

	return out, nil
}

func (c *MockForgeClient) FindTenantKeysetsByIds(ctx context.Context, in *wflows.TenantKeysetsByIdsRequest, opts ...grpc.CallOption) (*wflows.TenantKeySetList, error) {
	err, ok := ctx.Value("wantError").(error)
	if ok {
		return nil, status.Error(status.Code(err), "failed to retrieve tenant keysets")
	}

	out := &wflows.TenantKeySetList{}
	if in != nil {
		for _, id := range in.KeysetIds {
			out.Keyset = append(out.Keyset, &wflows.TenantKeyset{
				KeysetIdentifier: &wflows.TenantKeysetIdentifier{
					OrganizationId: id.OrganizationId,
					KeysetId:       id.KeysetId,
				},
			})
		}
	}

	return out, nil
}

// DEPRECATED: use FindTenantKeysetIds and FindTenantKeysetsByIds instead
func (c *MockForgeClient) FindTenantKeyset(ctx context.Context, in *wflows.FindTenantKeysetRequest, opts ...grpc.CallOption) (*wflows.TenantKeySetList, error) {
	out := &wflows.TenantKeySetList{}
	count, ok := ctx.Value("wantCount").(int)
	if ok {
		for i := 0; i < count; i++ {
			out.Keyset = append(out.Keyset, &wflows.TenantKeyset{KeysetIdentifier: &wflows.TenantKeysetIdentifier{KeysetId: uuid.NewString()}})
		}
	}
	return out, nil
}

func (c *MockForgeClient) UpdateTenantKeyset(ctx context.Context, in *wflows.UpdateTenantKeysetRequest, opts ...grpc.CallOption) (*wflows.UpdateTenantKeysetResponse, error) {
	out := new(wflows.UpdateTenantKeysetResponse)
	return out, nil
}

func (c *MockForgeClient) DeleteTenantKeyset(ctx context.Context, in *wflows.DeleteTenantKeysetRequest, opts ...grpc.CallOption) (*wflows.DeleteTenantKeysetResponse, error) {
	out := new(wflows.DeleteTenantKeysetResponse)
	return out, nil
}

func (c *MockForgeClient) ValidateTenantPublicKey(ctx context.Context, in *wflows.ValidateTenantPublicKeyRequest, opts ...grpc.CallOption) (*wflows.ValidateTenantPublicKeyResponse, error) {
	out := new(wflows.ValidateTenantPublicKeyResponse)
	return out, nil
}

func (c *MockForgeClient) GetDpuSSHCredential(ctx context.Context, in *wflows.CredentialRequest, opts ...grpc.CallOption) (*wflows.CredentialResponse, error) {
	out := new(wflows.CredentialResponse)
	return out, nil
}

func (c *MockForgeClient) GetAllManagedHostNetworkStatus(ctx context.Context, in *wflows.ManagedHostNetworkStatusRequest, opts ...grpc.CallOption) (*wflows.ManagedHostNetworkStatusResponse, error) {
	out := new(wflows.ManagedHostNetworkStatusResponse)
	return out, nil
}

func (c *MockForgeClient) GetSiteExplorationReport(ctx context.Context, in *wflows.GetSiteExplorationRequest, opts ...grpc.CallOption) (*wflows.SiteExplorationReport, error) {
	out := new(wflows.SiteExplorationReport)
	return out, nil
}

func (c *MockForgeClient) AdminForceDeleteMachine(ctx context.Context, in *wflows.AdminForceDeleteMachineRequest, opts ...grpc.CallOption) (*wflows.AdminForceDeleteMachineResponse, error) {
	out := new(wflows.AdminForceDeleteMachineResponse)
	return out, nil
}

func (c *MockForgeClient) AdminReboot(ctx context.Context, in *wflows.AdminRebootRequest, opts ...grpc.CallOption) (*wflows.AdminRebootResponse, error) {
	out := new(wflows.AdminRebootResponse)
	return out, nil
}

func (c *MockForgeClient) AdminListResourcePools(ctx context.Context, in *wflows.ListResourcePoolsRequest, opts ...grpc.CallOption) (*wflows.ResourcePools, error) {
	out := new(wflows.ResourcePools)
	return out, nil
}

func (c *MockForgeClient) AdminGrowResourcePool(ctx context.Context, in *wflows.GrowResourcePoolRequest, opts ...grpc.CallOption) (*wflows.GrowResourcePoolResponse, error) {
	out := new(wflows.GrowResourcePoolResponse)
	return out, nil
}

func (c *MockForgeClient) MigrateVpcVni(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*wflows.MigrateVpcVniResponse, error) {
	out := new(wflows.MigrateVpcVniResponse)
	return out, nil
}

func (c *MockForgeClient) SetMaintenance(ctx context.Context, in *wflows.MaintenanceRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	return out, nil
}

func (c *MockForgeClient) TriggerDpuReprovisioning(ctx context.Context, in *wflows.DpuReprovisioningRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	return out, nil
}

func (c *MockForgeClient) ListDpuWaitingForReprovisioning(ctx context.Context, in *wflows.DpuReprovisioningListRequest, opts ...grpc.CallOption) (*wflows.DpuReprovisioningListResponse, error) {
	out := new(wflows.DpuReprovisioningListResponse)
	return out, nil
}

func (c *MockForgeClient) GetMachineBootOverride(ctx context.Context, in *wflows.UUID, opts ...grpc.CallOption) (*wflows.MachineBootOverride, error) {
	out := new(wflows.MachineBootOverride)
	return out, nil
}

func (c *MockForgeClient) SetMachineBootOverride(ctx context.Context, in *wflows.MachineBootOverride, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	return out, nil
}

func (c *MockForgeClient) ClearMachineBootOverride(ctx context.Context, in *wflows.UUID, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	return out, nil
}

func (c *MockForgeClient) GetNetworkTopology(ctx context.Context, in *wflows.NetworkTopologyRequest, opts ...grpc.CallOption) (*wflows.NetworkTopologyData, error) {
	out := new(wflows.NetworkTopologyData)
	return out, nil
}

func (c *MockForgeClient) AdminBmcReset(ctx context.Context, in *wflows.AdminBmcResetRequest, opts ...grpc.CallOption) (*wflows.AdminBmcResetResponse, error) {
	out := new(wflows.AdminBmcResetResponse)
	return out, nil
}

func (c *MockForgeClient) CreateCredential(ctx context.Context, in *wflows.CredentialCreationRequest, opts ...grpc.CallOption) (*wflows.CredentialCreationResult, error) {
	out := new(wflows.CredentialCreationResult)
	return out, nil
}

func (c *MockForgeClient) DeleteCredential(ctx context.Context, in *wflows.CredentialDeletionRequest, opts ...grpc.CallOption) (*wflows.CredentialDeletionResult, error) {
	out := new(wflows.CredentialDeletionResult)
	return out, nil
}

func (c *MockForgeClient) GetRouteServers(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*wflows.RouteServers, error) {
	out := new(wflows.RouteServers)
	return out, nil
}

func (c *MockForgeClient) AddRouteServers(ctx context.Context, in *wflows.RouteServers, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	return out, nil
}

func (c *MockForgeClient) RemoveRouteServers(ctx context.Context, in *wflows.RouteServers, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	return out, nil
}

func (c *MockForgeClient) ReplaceRouteServers(ctx context.Context, in *wflows.RouteServers, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	return out, nil
}

func (c *MockForgeClient) UpdateAgentReportedInventory(ctx context.Context, in *wflows.DpuAgentInventoryReport, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	return out, nil
}

func (c *MockForgeClient) UpdateInstanceConfig(ctx context.Context, in *wflows.InstanceConfigUpdateRequest, opts ...grpc.CallOption) (*wflows.Instance, error) {
	out := new(wflows.Instance)
	return out, nil
}

func (c *MockForgeClient) CreateOsImage(ctx context.Context, in *wflows.OsImageAttributes, opts ...grpc.CallOption) (*wflows.OsImage, error) {
	out := new(wflows.OsImage)
	return out, nil
}

func (c *MockForgeClient) UpdateOsImage(ctx context.Context, in *wflows.OsImageAttributes, opts ...grpc.CallOption) (*wflows.OsImage, error) {
	out := new(wflows.OsImage)
	return out, nil
}

func (c *MockForgeClient) DeleteOsImage(ctx context.Context, in *wflows.DeleteOsImageRequest, opts ...grpc.CallOption) (*wflows.DeleteOsImageResponse, error) {
	out := new(wflows.DeleteOsImageResponse)
	return out, nil
}

func (c *MockForgeClient) GetOsImage(ctx context.Context, in *wflows.UUID, opts ...grpc.CallOption) (*wflows.OsImage, error) {
	out := new(wflows.OsImage)
	return out, nil
}

func (c *MockForgeClient) ListOsImage(ctx context.Context, in *wflows.ListOsImageRequest, opts ...grpc.CallOption) (*wflows.ListOsImageResponse, error) {
	if err, ok := ctx.Value("wantError").(error); ok {
		return nil, status.Error(status.Code(err), "failed to retrieve os image list")
	}

	out := &wflows.ListOsImageResponse{}
	count, ok := ctx.Value("wantCount").(int)
	if ok {
		id := uuid.NewString()
		for i := 0; i < count; i++ {
			out.Images = append(out.Images, &wflows.OsImage{Attributes: &wflows.OsImageAttributes{Id: &wflows.UUID{Value: id}}})
		}
	}
	return out, nil
}

func (c *MockForgeClient) FindTenantOrganizationIds(ctx context.Context, in *wflows.TenantSearchFilter, opts ...grpc.CallOption) (*wflows.TenantOrganizationIdList, error) {
	if err, ok := ctx.Value("wantError").(error); ok {
		return nil, status.Error(status.Code(err), "failed to retrieve Tenant organization ids")
	}

	out := &wflows.TenantOrganizationIdList{}

	count, ok := ctx.Value("wantCount").(int)
	if ok {
		for i := 0; i < count; i++ {
			out.TenantOrganizationIds = append(out.TenantOrganizationIds, randSeq(10))
		}
	}

	return out, nil
}

func (c *MockForgeClient) FindTenantsByOrganizationIds(ctx context.Context, in *wflows.TenantByOrganizationIdsRequest, opts ...grpc.CallOption) (*wflows.TenantList, error) {
	err, ok := ctx.Value("wantError").(error)
	if ok {
		return nil, status.Error(status.Code(err), "failed to retrieve Tenants")
	}

	out := &wflows.TenantList{}
	if in != nil {
		for _, id := range in.OrganizationIds {
			out.Tenants = append(out.Tenants, &wflows.Tenant{
				OrganizationId: id,
			})
		}
	}

	return out, nil
}

func (c *MockForgeClient) CreateInstanceType(ctx context.Context, in *wflows.CreateInstanceTypeRequest, opts ...grpc.CallOption) (*wflows.CreateInstanceTypeResponse, error) {
	out := &wflows.CreateInstanceTypeResponse{InstanceType: &wflows.InstanceType{}}
	return out, nil
}

func (c *MockForgeClient) UpdateInstanceType(ctx context.Context, in *wflows.UpdateInstanceTypeRequest, opts ...grpc.CallOption) (*wflows.UpdateInstanceTypeResponse, error) {
	out := &wflows.UpdateInstanceTypeResponse{InstanceType: &wflows.InstanceType{}}
	return out, nil
}

func (c *MockForgeClient) DeleteInstanceType(ctx context.Context, in *wflows.DeleteInstanceTypeRequest, opts ...grpc.CallOption) (*wflows.DeleteInstanceTypeResponse, error) {
	out := &wflows.DeleteInstanceTypeResponse{}
	return out, nil
}

func (c *MockForgeClient) AssociateMachinesWithInstanceType(ctx context.Context, in *wflows.AssociateMachinesWithInstanceTypeRequest, opts ...grpc.CallOption) (*wflows.AssociateMachinesWithInstanceTypeResponse, error) {
	out := &wflows.AssociateMachinesWithInstanceTypeResponse{}
	return out, nil
}

func (c *MockForgeClient) RemoveMachineInstanceTypeAssociation(ctx context.Context, in *wflows.RemoveMachineInstanceTypeAssociationRequest, opts ...grpc.CallOption) (*wflows.RemoveMachineInstanceTypeAssociationResponse, error) {
	out := &wflows.RemoveMachineInstanceTypeAssociationResponse{}
	return out, nil
}

func (c *MockForgeClient) FindInstanceTypeIds(ctx context.Context, in *wflows.FindInstanceTypeIdsRequest, opts ...grpc.CallOption) (*wflows.FindInstanceTypeIdsResponse, error) {
	if err, ok := ctx.Value("wantError").(error); ok {
		return nil, status.Error(status.Code(err), "failed to retrieve InstanceType ids")
	}

	out := &wflows.FindInstanceTypeIdsResponse{}

	count, ok := ctx.Value("wantCount").(int)
	if ok {
		for i := 0; i < count; i++ {
			out.InstanceTypeIds = append(out.InstanceTypeIds, randSeq(10))
		}
	}

	return out, nil
}

func (c *MockForgeClient) FindInstanceTypesByIds(ctx context.Context, in *wflows.FindInstanceTypesByIdsRequest, opts ...grpc.CallOption) (*wflows.FindInstanceTypesByIdsResponse, error) {
	err, ok := ctx.Value("wantError").(error)
	if ok {
		return nil, status.Error(status.Code(err), "failed to retrieve InstanceTypes")
	}

	out := &wflows.FindInstanceTypesByIdsResponse{}
	if in != nil {
		for _, id := range in.InstanceTypeIds {
			out.InstanceTypes = append(out.InstanceTypes, &wflows.InstanceType{
				Id: id,
			})
		}
	}
	return out, nil
}

func (c *MockForgeClient) CreateVpcPrefix(ctx context.Context, in *wflows.VpcPrefixCreationRequest, opts ...grpc.CallOption) (*wflows.VpcPrefix, error) {
	out := new(wflows.VpcPrefix)
	return out, nil
}

func (c *MockForgeClient) UpdateVpcPrefix(ctx context.Context, in *wflows.VpcPrefixUpdateRequest, opts ...grpc.CallOption) (*wflows.VpcPrefix, error) {
	out := new(wflows.VpcPrefix)
	return out, nil
}

func (c *MockForgeClient) DeleteVpcPrefix(ctx context.Context, in *wflows.VpcPrefixDeletionRequest, opts ...grpc.CallOption) (*wflows.VpcPrefixDeletionResult, error) {
	out := new(wflows.VpcPrefixDeletionResult)
	return out, nil
}

func (c *MockForgeClient) SearchVpcPrefixes(ctx context.Context, in *wflows.VpcPrefixSearchQuery, opts ...grpc.CallOption) (*wflows.VpcPrefixIdList, error) {
	if err, ok := ctx.Value("wantError").(error); ok {
		return nil, status.Error(status.Code(err), "failed to retrieve vpcprefix ids")
	}

	out := &wflows.VpcPrefixIdList{}

	count, ok := ctx.Value("wantCount").(int)
	if ok {
		for i := 0; i < count; i++ {
			out.VpcPrefixIds = append(out.VpcPrefixIds, &wflows.UUID{Value: uuid.NewString()})
		}
	}

	return out, nil
}

func (c *MockForgeClient) GetVpcPrefixes(ctx context.Context, in *wflows.VpcPrefixGetRequest, opts ...grpc.CallOption) (*wflows.VpcPrefixList, error) {
	err, ok := ctx.Value("wantError").(error)
	if ok {
		return nil, status.Error(status.Code(err), "failed to retrieve vpcprefixes")
	}

	out := &wflows.VpcPrefixList{}
	if in != nil {
		for _, id := range in.VpcPrefixIds {
			out.VpcPrefixes = append(out.VpcPrefixes, &wflows.VpcPrefix{
				Id: id,
			})
		}
	}

	return out, nil
}

func (c *MockForgeClient) CreateVpcPeering(ctx context.Context, in *wflows.VpcPeeringCreationRequest, opts ...grpc.CallOption) (*wflows.VpcPeering, error) {
	out := new(wflows.VpcPeering)
	out.Id = &wflows.UUID{Value: uuid.NewString()}
	return out, nil
}

func (c *MockForgeClient) FindVpcPeeringIds(ctx context.Context, in *wflows.VpcPeeringSearchFilter, opts ...grpc.CallOption) (*wflows.VpcPeeringIdList, error) {
	if err, ok := ctx.Value("wantError").(error); ok {
		return nil, status.Error(status.Code(err), "failed to retrieve vpc peering ids")
	}

	out := &wflows.VpcPeeringIdList{}

	count, ok := ctx.Value("WantCount").(int)
	if ok {
		for i := 0; i < count; i++ {
			out.VpcPeeringIds = append(out.VpcPeeringIds, &wflows.UUID{Value: uuid.NewString()})
		}
	}

	return out, nil
}

func (c *MockForgeClient) FindVpcPeeringsByIds(ctx context.Context, in *wflows.VpcPeeringsByIdsRequest, opts ...grpc.CallOption) (*wflows.VpcPeeringList, error) {
	err, ok := ctx.Value("wantError").(error)
	if ok {
		return nil, status.Error(status.Code(err), "failed to retrieve vpc peerings")
	}

	out := &wflows.VpcPeeringList{}
	for _, id := range in.VpcPeeringIds {
		out.VpcPeerings = append(out.VpcPeerings, &wflows.VpcPeering{
			Id:        id,
			VpcId:     &wflows.UUID{Value: uuid.NewString()},
			PeerVpcId: &wflows.UUID{Value: uuid.NewString()},
		})
	}

	return out, nil
}

func (c *MockForgeClient) DeleteVpcPeering(ctx context.Context, in *wflows.VpcPeeringDeletionRequest, opts ...grpc.CallOption) (*wflows.VpcPeeringDeletionResult, error) {
	if err, ok := ctx.Value("wantError").(error); ok {
		return nil, status.Error(status.Code(err), "failed to delete vpc peering")
	}

	return &wflows.VpcPeeringDeletionResult{}, nil
}

func (c *MockForgeClient) PersistValidationResult(ctx context.Context, in *wflows.MachineValidationResultPostRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	return out, nil
}

func (c *MockForgeClient) GetMachineValidationResults(ctx context.Context, in *wflows.MachineValidationGetRequest, opts ...grpc.CallOption) (*wflows.MachineValidationResultList, error) {
	out := new(wflows.MachineValidationResultList)
	return out, nil
}

func (c *MockForgeClient) MachineValidationCompleted(ctx context.Context, in *wflows.MachineValidationCompletedRequest, opts ...grpc.CallOption) (*wflows.MachineValidationCompletedResponse, error) {
	out := new(wflows.MachineValidationCompletedResponse)
	return out, nil
}

func (c *MockForgeClient) MachineSetAutoUpdate(ctx context.Context, in *wflows.MachineSetAutoUpdateRequest, opts ...grpc.CallOption) (*wflows.MachineSetAutoUpdateResponse, error) {
	out := new(wflows.MachineSetAutoUpdateResponse)
	return out, nil
}

func (c *MockForgeClient) GetMachineValidationExternalConfigs(ctx context.Context, in *wflows.GetMachineValidationExternalConfigsRequest, opts ...grpc.CallOption) (*wflows.GetMachineValidationExternalConfigsResponse, error) {
	out := new(wflows.GetMachineValidationExternalConfigsResponse)
	return out, nil
}

func (c *MockForgeClient) AddUpdateMachineValidationExternalConfig(ctx context.Context, in *wflows.AddUpdateMachineValidationExternalConfigRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	return out, nil
}

func (c *MockForgeClient) GetMachineValidationRuns(ctx context.Context, in *wflows.MachineValidationRunListGetRequest, opts ...grpc.CallOption) (*wflows.MachineValidationRunList, error) {
	out := new(wflows.MachineValidationRunList)
	return out, nil
}

func (c *MockForgeClient) RemoveMachineValidationExternalConfig(ctx context.Context, in *wflows.RemoveMachineValidationExternalConfigRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	return out, nil
}

func (c *MockForgeClient) GetMachineValidationTests(ctx context.Context, in *wflows.MachineValidationTestsGetRequest, opts ...grpc.CallOption) (*wflows.MachineValidationTestsGetResponse, error) {
	out := new(wflows.MachineValidationTestsGetResponse)
	return out, nil
}

func (c *MockForgeClient) AddMachineValidationTest(ctx context.Context, in *wflows.MachineValidationTestAddRequest, opts ...grpc.CallOption) (*wflows.MachineValidationTestAddUpdateResponse, error) {
	out := new(wflows.MachineValidationTestAddUpdateResponse)
	id, ok := ctx.Value("wantID").(string)
	if ok {
		out.TestId = id
		out.Version = "version-1"
	}
	return out, nil
}

func (c *MockForgeClient) UpdateMachineValidationTest(ctx context.Context, in *wflows.MachineValidationTestUpdateRequest, opts ...grpc.CallOption) (*wflows.MachineValidationTestAddUpdateResponse, error) {
	out := new(wflows.MachineValidationTestAddUpdateResponse)
	out.TestId = in.TestId
	out.Version = in.Version
	return out, nil
}

func (c *MockForgeClient) MachineValidationTestVerfied(ctx context.Context, in *wflows.MachineValidationTestVerfiedRequest, opts ...grpc.CallOption) (*wflows.MachineValidationTestVerfiedResponse, error) {
	out := new(wflows.MachineValidationTestVerfiedResponse)
	return out, nil
}

func (c *MockForgeClient) MachineValidationTestNextVersion(ctx context.Context, in *wflows.MachineValidationTestNextVersionRequest, opts ...grpc.CallOption) (*wflows.MachineValidationTestNextVersionResponse, error) {
	out := new(wflows.MachineValidationTestNextVersionResponse)
	return out, nil
}

func (c *MockForgeClient) MachineValidationTestEnableDisableTest(ctx context.Context, in *wflows.MachineValidationTestEnableDisableTestRequest, opts ...grpc.CallOption) (*wflows.MachineValidationTestEnableDisableTestResponse, error) {
	out := new(wflows.MachineValidationTestEnableDisableTestResponse)
	return out, nil
}

func (c *MockForgeClient) UpdateMachineValidationRun(ctx context.Context, in *wflows.MachineValidationRunRequest, opts ...grpc.CallOption) (*wflows.MachineValidationRunResponse, error) {
	out := new(wflows.MachineValidationRunResponse)
	return out, nil
}

func (c *MockForgeClient) CreateNetworkSecurityGroup(ctx context.Context, in *wflows.CreateNetworkSecurityGroupRequest, opts ...grpc.CallOption) (*wflows.CreateNetworkSecurityGroupResponse, error) {
	out := &wflows.CreateNetworkSecurityGroupResponse{NetworkSecurityGroup: &wflows.NetworkSecurityGroup{}}
	return out, nil
}

func (c *MockForgeClient) UpdateNetworkSecurityGroup(ctx context.Context, in *wflows.UpdateNetworkSecurityGroupRequest, opts ...grpc.CallOption) (*wflows.UpdateNetworkSecurityGroupResponse, error) {
	out := &wflows.UpdateNetworkSecurityGroupResponse{NetworkSecurityGroup: &wflows.NetworkSecurityGroup{}}
	return out, nil
}

func (c *MockForgeClient) DeleteNetworkSecurityGroup(ctx context.Context, in *wflows.DeleteNetworkSecurityGroupRequest, opts ...grpc.CallOption) (*wflows.DeleteNetworkSecurityGroupResponse, error) {
	out := &wflows.DeleteNetworkSecurityGroupResponse{}
	return out, nil
}

func (c *MockForgeClient) GetNetworkSecurityGroupAttachments(ctx context.Context, in *wflows.GetNetworkSecurityGroupAttachmentsRequest, opts ...grpc.CallOption) (*wflows.GetNetworkSecurityGroupAttachmentsResponse, error) {
	out := &wflows.GetNetworkSecurityGroupAttachmentsResponse{}
	return out, nil
}

func (c *MockForgeClient) GetNetworkSecurityGroupPropagationStatus(ctx context.Context, in *wflows.GetNetworkSecurityGroupPropagationStatusRequest, opts ...grpc.CallOption) (*wflows.GetNetworkSecurityGroupPropagationStatusResponse, error) {
	out := &wflows.GetNetworkSecurityGroupPropagationStatusResponse{}
	return out, nil
}

func (c *MockForgeClient) FindNetworkSecurityGroupIds(ctx context.Context, in *wflows.FindNetworkSecurityGroupIdsRequest, opts ...grpc.CallOption) (*wflows.FindNetworkSecurityGroupIdsResponse, error) {
	if err, ok := ctx.Value("wantError").(error); ok {
		return nil, status.Error(status.Code(err), "failed to retrieve NetworkSecurityGroup ids")
	}

	out := &wflows.FindNetworkSecurityGroupIdsResponse{}

	count, ok := ctx.Value("wantCount").(int)
	if ok {
		for i := 0; i < count; i++ {
			out.NetworkSecurityGroupIds = append(out.NetworkSecurityGroupIds, randSeq(10))
		}
	}

	return out, nil
}

func (c *MockForgeClient) FindNetworkSecurityGroupsByIds(ctx context.Context, in *wflows.FindNetworkSecurityGroupsByIdsRequest, opts ...grpc.CallOption) (*wflows.FindNetworkSecurityGroupsByIdsResponse, error) {
	err, ok := ctx.Value("wantError").(error)
	if ok {
		return nil, status.Error(status.Code(err), "failed to retrieve NetworkSecurityGroups")
	}

	out := &wflows.FindNetworkSecurityGroupsByIdsResponse{}
	if in != nil {
		for _, id := range in.NetworkSecurityGroupIds {
			out.NetworkSecurityGroups = append(out.NetworkSecurityGroups, &wflows.NetworkSecurityGroup{
				Id: id,
			})
		}
	}
	return out, nil
}

func (c *MockForgeClient) UpdateMachineMetadata(ctx context.Context, in *wflows.MachineMetadataUpdateRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	err, ok := ctx.Value("wantError").(error)
	if ok {
		return nil, status.Error(status.Code(err), "failed to update machine metadata")
	}

	out := new(emptypb.Empty)
	return out, nil
}

// Expected Machine mock methods
func (c *MockForgeClient) AddExpectedMachine(ctx context.Context, in *wflows.ExpectedMachine, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	if in.Id == nil || in.Id.Value == "" {
		return nil, status.Error(codes.Internal, "ID not provided for AddExpectedMachine")
	}
	if in.BmcMacAddress == "" {
		return nil, status.Error(codes.Internal, "MAC address not provided for AddExpectedMachine")
	}
	if in.ChassisSerialNumber == "" {
		return nil, status.Error(codes.Internal, "Chassis Serial Number not provided for AddExpectedMachine")
	}
	out := new(emptypb.Empty)
	return out, nil
}

func (c *MockForgeClient) DeleteExpectedMachine(ctx context.Context, in *wflows.ExpectedMachineRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	if in.Id == nil || in.Id.Value == "" {
		return nil, status.Error(codes.Internal, "ID not provided for DeleteExpectedMachine")
	}
	out := new(emptypb.Empty)
	return out, nil
}

func (c *MockForgeClient) UpdateExpectedMachine(ctx context.Context, in *wflows.ExpectedMachine, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	if in.Id == nil || in.Id.Value == "" {
		return nil, status.Error(codes.Internal, "ID not provided for UpdateExpectedMachine")
	}
	if in.BmcMacAddress == "" {
		return nil, status.Error(codes.Internal, "MAC address not provided for UpdateExpectedMachine")
	}
	if in.ChassisSerialNumber == "" {
		return nil, status.Error(codes.Internal, "Chassis Serial Number not provided for UpdateExpectedMachine")
	}
	out := new(emptypb.Empty)
	return out, nil
}

func (c *MockForgeClient) GetAllExpectedMachines(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*wflows.ExpectedMachineList, error) {
	err, ok := ctx.Value("wantError").(error)
	if ok {
		if status.Code(err) == codes.Internal {
			return nil, status.Error(codes.Internal, "failed to retrieve machine ids")
		}
	}

	out := &wflows.ExpectedMachineList{}

	// we generate predictable unique IDs and values
	count, ok := ctx.Value("wantCount").(int)
	if ok {
		mac, _ := net.ParseMAC("02:00:00:00:00:00")
		for i := 0; i < count; i++ {
			// Create a 16-byte array for UUID from MAC address (6 bytes) + padding
			var uuidBytes [16]byte
			copy(uuidBytes[:6], mac)
			emID, _ := uuid.FromBytes(uuidBytes[:])
			out.ExpectedMachines = append(out.ExpectedMachines, &wflows.ExpectedMachine{
				Id:                  &wflows.UUID{Value: emID.String()},
				BmcMacAddress:       mac.String(),
				ChassisSerialNumber: "serial-" + mac.String()})
			incrementMAC(mac)
		}
	}

	return out, nil
}

func (c *MockForgeClient) GetExpectedMachine(ctx context.Context, in *wflows.ExpectedMachineRequest, opts ...grpc.CallOption) (*wflows.ExpectedMachine, error) {
	if in.Id == nil || in.Id.Value == "" {
		return nil, status.Error(codes.Internal, "ID not provided for GetExpectedMachine")
	}
	out := new(wflows.ExpectedMachine)
	return out, nil
}

func (c *MockForgeClient) FindSkusByIds(ctx context.Context, in *wflows.SkusByIdsRequest, opts ...grpc.CallOption) (*wflows.SkuList, error) {
	err, ok := ctx.Value("wantError").(error)
	if ok {
		return nil, status.Error(status.Code(err), "failed to retrieve skus")
	}

	out := &wflows.SkuList{}
	if in != nil {
		for _, id := range in.Ids {
			out.Skus = append(out.Skus, &wflows.Sku{
				Id: id,
			})
		}
	}

	return out, nil
}

func (c *MockForgeClient) GetAllSkuIds(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*wflows.SkuIdList, error) {
	if err, ok := ctx.Value("wantError").(error); ok {
		return nil, status.Error(status.Code(err), "failed to retrieve sku ids")
	}

	out := &wflows.SkuIdList{}

	count, ok := ctx.Value("wantCount").(int)
	if ok {
		for i := 0; i < count; i++ {
			out.Ids = append(out.Ids, uuid.NewString())
		}
	}

	return out, nil
}

// NewMockCarbideClient creates a new mock CarbideClient
func NewMockCarbideClient() *CarbideClient {
	return &CarbideClient{
		carbide: &MockForgeClient{},
	}
}
