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

package api

import (
	"net/http"

	tClient "go.temporal.io/sdk/client"

	"github.com/nvidia/carbide-rest/api/internal/config"
	apiHandler "github.com/nvidia/carbide-rest/api/pkg/api/handler"
	cdb "github.com/nvidia/carbide-rest/db/pkg/db"

	sc "github.com/nvidia/carbide-rest/api/pkg/client/site"
)

const (
	// EchoUnmatchedPath is the path to use when no route matches the request
	EchoUnmatchedPath = "/v2/*"
)

// NewV2APIRoutes returns all v2 routes
// v2 routes are auto-prefixed with "/v2"
// We start with v2 version to align with NGC API
func NewV2APIRoutes(dbSession *cdb.Session, tc tClient.Client, tnc tClient.NamespaceClient, scp *sc.ClientPool, cfg *config.Config) []Route {
	apiRoutes := []Route{
		// Metadata endpoint
		{
			Path:    "/org/:orgName/carbide/metadata",
			Method:  http.MethodGet,
			Handler: apiHandler.NewMetadataHandler(),
		},
		// User endpoint
		{
			Path:    "/org/:orgName/carbide/user/current",
			Method:  http.MethodGet,
			Handler: apiHandler.NewGetUserHandler(dbSession),
		},
		// Service Account endpoint
		{
			Path:    "/org/:orgName/carbide/service-account/current",
			Method:  http.MethodGet,
			Handler: apiHandler.NewGetCurrentServiceAccountHandler(dbSession, cfg),
		},
		// Infrastructure Provider endpoints
		{
			Path:    "/org/:orgName/carbide/infrastructure-provider",
			Method:  http.MethodPost,
			Handler: apiHandler.NewCreateInfrastructureProviderHandler(dbSession, tc, cfg),
		},
		{
			Path:    "/org/:orgName/carbide/infrastructure-provider/current",
			Method:  http.MethodGet,
			Handler: apiHandler.NewGetCurrentInfrastructureProviderHandler(dbSession, tc, cfg),
		},
		{
			Path:    "/org/:orgName/carbide/infrastructure-provider/current",
			Method:  http.MethodPatch,
			Handler: apiHandler.NewUpdateCurrentInfrastructureProviderHandler(dbSession, tc, cfg),
		},
		{
			Path:    "/org/:orgName/carbide/infrastructure-provider/current/stats",
			Method:  http.MethodGet,
			Handler: apiHandler.NewGetCurrentInfrastructureProviderStatsHandler(dbSession, tc, cfg),
		},
		// Tenant endpoints
		{
			Path:    "/org/:orgName/carbide/tenant",
			Method:  http.MethodPost,
			Handler: apiHandler.NewCreateTenantHandler(dbSession, tc, cfg),
		},
		{
			Path:    "/org/:orgName/carbide/tenant/current",
			Method:  http.MethodGet,
			Handler: apiHandler.NewGetCurrentTenantHandler(dbSession, tc, cfg),
		},
		{
			Path:    "/org/:orgName/carbide/tenant/current",
			Method:  http.MethodPatch,
			Handler: apiHandler.NewUpdateCurrentTenantHandler(dbSession, tc, cfg),
		},
		{
			Path:    "/org/:orgName/carbide/tenant/current/stats",
			Method:  http.MethodGet,
			Handler: apiHandler.NewGetCurrentTenantStatsHandler(dbSession, tc, cfg),
		},
		// TenantAccount endpoints
		{
			Path:    "/org/:orgName/carbide/tenant/account",
			Method:  http.MethodGet,
			Handler: apiHandler.NewGetAllTenantAccountHandler(dbSession, tc, cfg),
		},
		{
			Path:    "/org/:orgName/carbide/tenant/account/:id",
			Method:  http.MethodGet,
			Handler: apiHandler.NewGetTenantAccountHandler(dbSession, tc, cfg),
		},
		{
			Path:    "/org/:orgName/carbide/tenant/account",
			Method:  http.MethodPost,
			Handler: apiHandler.NewCreateTenantAccountHandler(dbSession, tc, cfg),
		},
		{
			Path:    "/org/:orgName/carbide/tenant/account/:id",
			Method:  http.MethodPatch,
			Handler: apiHandler.NewUpdateTenantAccountHandler(dbSession, tc, cfg),
		},
		{
			Path:    "/org/:orgName/carbide/tenant/account/:id",
			Method:  http.MethodDelete,
			Handler: apiHandler.NewDeleteTenantAccountHandler(dbSession, tc, cfg),
		},
		// Site endpoints
		{
			Path:    "/org/:orgName/carbide/site",
			Method:  http.MethodPost,
			Handler: apiHandler.NewCreateSiteHandler(dbSession, tc, tnc, cfg),
		},
		{
			Path:    "/org/:orgName/carbide/site",
			Method:  http.MethodGet,
			Handler: apiHandler.NewGetAllSiteHandler(dbSession, tc, cfg),
		},
		{
			Path:    "/org/:orgName/carbide/site/:id",
			Method:  http.MethodGet,
			Handler: apiHandler.NewGetSiteHandler(dbSession, tc, cfg),
		},
		{
			Path:    "/org/:orgName/carbide/site/:id",
			Method:  http.MethodPatch,
			Handler: apiHandler.NewUpdateSiteHandler(dbSession, tc, cfg),
		},
		{
			Path:    "/org/:orgName/carbide/site/:id",
			Method:  http.MethodDelete,
			Handler: apiHandler.NewDeleteSiteHandler(dbSession, tc, cfg),
		},
		{
			Path:    "/org/:orgName/carbide/site/:id/status-history",
			Method:  http.MethodGet,
			Handler: apiHandler.NewGetSiteStatusDetailsHandler(dbSession),
		},
		// VPC endpoints
		{
			Path:    "/org/:orgName/carbide/vpc",
			Method:  http.MethodPost,
			Handler: apiHandler.NewCreateVPCHandler(dbSession, tc, scp, cfg),
		},
		{
			Path:    "/org/:orgName/carbide/vpc",
			Method:  http.MethodGet,
			Handler: apiHandler.NewGetAllVPCHandler(dbSession, tc, cfg),
		},
		{
			Path:    "/org/:orgName/carbide/vpc/:id",
			Method:  http.MethodGet,
			Handler: apiHandler.NewGetVPCHandler(dbSession, tc, cfg),
		},
		{
			Path:    "/org/:orgName/carbide/vpc/:id",
			Method:  http.MethodPatch,
			Handler: apiHandler.NewUpdateVPCHandler(dbSession, tc, scp, cfg),
		},
		{
			Path:    "/org/:orgName/carbide/vpc/:id",
			Method:  http.MethodDelete,
			Handler: apiHandler.NewDeleteVPCHandler(dbSession, tc, scp, cfg),
		},
		{
			Path:    "/org/:orgName/carbide/vpc/:id/virtualization",
			Method:  http.MethodPatch,
			Handler: apiHandler.NewUpdateVPCVirtualizationHandler(dbSession, tc, scp, cfg),
		},

		// VpcPrefix endpoints
		{
			Path:    "/org/:orgName/carbide/vpc-prefix",
			Method:  http.MethodPost,
			Handler: apiHandler.NewCreateVpcPrefixHandler(dbSession, tc, scp, cfg),
		},
		{
			Path:    "/org/:orgName/carbide/vpc-prefix",
			Method:  http.MethodGet,
			Handler: apiHandler.NewGetAllVpcPrefixHandler(dbSession, tc, cfg),
		},
		{
			Path:    "/org/:orgName/carbide/vpc-prefix/:id",
			Method:  http.MethodGet,
			Handler: apiHandler.NewGetVpcPrefixHandler(dbSession, tc, cfg),
		},
		{
			Path:    "/org/:orgName/carbide/vpc-prefix/:id",
			Method:  http.MethodPatch,
			Handler: apiHandler.NewUpdateVpcPrefixHandler(dbSession, tc, scp, cfg),
		},
		{
			Path:    "/org/:orgName/carbide/vpc-prefix/:id",
			Method:  http.MethodDelete,
			Handler: apiHandler.NewDeleteVpcPrefixHandler(dbSession, tc, scp, cfg),
		},

		// IPBlock endpoints
		{
			Path:    "/org/:orgName/carbide/ipblock",
			Method:  http.MethodPost,
			Handler: apiHandler.NewCreateIPBlockHandler(dbSession, tc, cfg),
		},
		{
			Path:    "/org/:orgName/carbide/ipblock",
			Method:  http.MethodGet,
			Handler: apiHandler.NewGetAllIPBlockHandler(dbSession, tc, cfg),
		},
		{
			Path:    "/org/:orgName/carbide/ipblock/:id",
			Method:  http.MethodGet,
			Handler: apiHandler.NewGetIPBlockHandler(dbSession, tc, cfg),
		},
		{
			Path:    "/org/:orgName/carbide/ipblock/:id/derived",
			Method:  http.MethodGet,
			Handler: apiHandler.NewGetAllDerivedIPBlockHandler(dbSession, tc, cfg),
		},
		{
			Path:    "/org/:orgName/carbide/ipblock/:id",
			Method:  http.MethodPatch,
			Handler: apiHandler.NewUpdateIPBlockHandler(dbSession, tc, cfg),
		},
		{
			Path:    "/org/:orgName/carbide/ipblock/:id",
			Method:  http.MethodDelete,
			Handler: apiHandler.NewDeleteIPBlockHandler(dbSession, tc, cfg),
		},
		// Instance endpoints
		{
			Path:    "/org/:orgName/carbide/instance",
			Method:  http.MethodPost,
			Handler: apiHandler.NewCreateInstanceHandler(dbSession, tc, scp, cfg),
		},
		{
			Path:    "/org/:orgName/carbide/instance",
			Method:  http.MethodGet,
			Handler: apiHandler.NewGetAllInstanceHandler(dbSession, tc, cfg),
		},
		{
			Path:    "/org/:orgName/carbide/instance/:id",
			Method:  http.MethodGet,
			Handler: apiHandler.NewGetInstanceHandler(dbSession, tc, cfg),
		},
		{
			Path:    "/org/:orgName/carbide/instance/:id",
			Method:  http.MethodPatch,
			Handler: apiHandler.NewUpdateInstanceHandler(dbSession, tc, scp, cfg),
		},
		{
			Path:    "/org/:orgName/carbide/instance/:id",
			Method:  http.MethodDelete,
			Handler: apiHandler.NewDeleteInstanceHandler(dbSession, tc, scp, cfg),
		},
		{
			Path:    "/org/:orgName/carbide/instance/:id/status-history",
			Method:  http.MethodGet,
			Handler: apiHandler.NewGetInstanceStatusDetailsHandler(dbSession),
		},
		// Instance Type endpoints
		{
			Path:    "/org/:orgName/carbide/instance/type",
			Method:  http.MethodPost,
			Handler: apiHandler.NewCreateInstanceTypeHandler(dbSession, tc, scp, cfg),
		},
		{
			Path:    "/org/:orgName/carbide/instance/type",
			Method:  http.MethodGet,
			Handler: apiHandler.NewGetAllInstanceTypeHandler(dbSession, tc, cfg),
		},
		{
			Path:    "/org/:orgName/carbide/instance/type/:id",
			Method:  http.MethodGet,
			Handler: apiHandler.NewGetInstanceTypeHandler(dbSession, tc, cfg),
		},
		{
			Path:    "/org/:orgName/carbide/instance/type/:id",
			Method:  http.MethodPatch,
			Handler: apiHandler.NewUpdateInstanceTypeHandler(dbSession, tc, scp, cfg),
		},
		{
			Path:    "/org/:orgName/carbide/instance/type/:id",
			Method:  http.MethodDelete,
			Handler: apiHandler.NewDeleteInstanceTypeHandler(dbSession, tc, scp, cfg),
		},
		// Interface endpoints
		{
			Path:    "/org/:orgName/carbide/instance/:instanceId/interface",
			Method:  http.MethodGet,
			Handler: apiHandler.NewGetAllInterfaceHandler(dbSession, tc, cfg),
		},
		// NVLinkInterface endpoints
		{
			Path:    "/org/:orgName/carbide/instance/:instanceId/nvlink-interface",
			Method:  http.MethodGet,
			Handler: apiHandler.NewGetAllNVLinkInterfaceHandler(dbSession, tc, cfg),
		},
		// InfiniBandPartition endpoints
		{
			Path:    "/org/:orgName/carbide/infiniband-partition",
			Method:  http.MethodPost,
			Handler: apiHandler.NewCreateInfiniBandPartitionHandler(dbSession, tc, scp, cfg),
		},
		{
			Path:    "/org/:orgName/carbide/infiniband-partition",
			Method:  http.MethodGet,
			Handler: apiHandler.NewGetAllInfiniBandPartitionHandler(dbSession, tc, cfg),
		},
		{
			Path:    "/org/:orgName/carbide/infiniband-partition/:id",
			Method:  http.MethodGet,
			Handler: apiHandler.NewGetInfiniBandPartitionHandler(dbSession, tc, cfg),
		},
		{
			Path:    "/org/:orgName/carbide/infiniband-partition/:id",
			Method:  http.MethodPatch,
			Handler: apiHandler.NewUpdateInfiniBandPartitionHandler(dbSession, tc, cfg),
		},
		{
			Path:    "/org/:orgName/carbide/infiniband-partition/:id",
			Method:  http.MethodDelete,
			Handler: apiHandler.NewDeleteInfiniBandPartitionHandler(dbSession, tc, scp, cfg),
		},
		// NVLinkLogicalPartition endpoints
		{
			Path:    "/org/:orgName/carbide/nvlink-logical-partition",
			Method:  http.MethodPost,
			Handler: apiHandler.NewCreateNVLinkLogicalPartitionHandler(dbSession, tc, scp, cfg),
		},
		{
			Path:    "/org/:orgName/carbide/nvlink-logical-partition",
			Method:  http.MethodGet,
			Handler: apiHandler.NewGetAllNVLinkLogicalPartitionHandler(dbSession, tc, cfg),
		},
		{
			Path:    "/org/:orgName/carbide/nvlink-logical-partition/:id",
			Method:  http.MethodGet,
			Handler: apiHandler.NewGetNVLinkLogicalPartitionHandler(dbSession, tc, cfg),
		},
		{
			Path:    "/org/:orgName/carbide/nvlink-logical-partition/:id",
			Method:  http.MethodPatch,
			Handler: apiHandler.NewUpdateNVLinkLogicalPartitionHandler(dbSession, tc, scp, cfg),
		},
		{
			Path:    "/org/:orgName/carbide/nvlink-logical-partition/:id",
			Method:  http.MethodDelete,
			Handler: apiHandler.NewDeleteNVLinkLogicalPartitionHandler(dbSession, tc, scp, cfg),
		},
		// ExpectedMachine endpoints
		{
			Path:    "/org/:orgName/carbide/expected-machine",
			Method:  http.MethodPost,
			Handler: apiHandler.NewCreateExpectedMachineHandler(dbSession, tc, scp, cfg),
		},
		{
			Path:    "/org/:orgName/carbide/expected-machine",
			Method:  http.MethodGet,
			Handler: apiHandler.NewGetAllExpectedMachineHandler(dbSession, tc, cfg),
		},
		{
			Path:    "/org/:orgName/carbide/expected-machine/:id",
			Method:  http.MethodGet,
			Handler: apiHandler.NewGetExpectedMachineHandler(dbSession, tc, cfg),
		},
		{
			Path:    "/org/:orgName/carbide/expected-machine/:id",
			Method:  http.MethodPatch,
			Handler: apiHandler.NewUpdateExpectedMachineHandler(dbSession, tc, scp, cfg),
		},
		{
			Path:    "/org/:orgName/carbide/expected-machine/:id",
			Method:  http.MethodDelete,
			Handler: apiHandler.NewDeleteExpectedMachineHandler(dbSession, tc, scp, cfg),
		},
		// Machine endpoints
		{
			Path:    "/org/:orgName/carbide/machine",
			Method:  http.MethodGet,
			Handler: apiHandler.NewGetAllMachineHandler(dbSession, tc, cfg),
		},
		{
			Path:    "/org/:orgName/carbide/machine/:id",
			Method:  http.MethodGet,
			Handler: apiHandler.NewGetMachineHandler(dbSession, tc, cfg),
		},
		{
			Path:    "/org/:orgName/carbide/machine/:id",
			Method:  http.MethodPatch,
			Handler: apiHandler.NewUpdateMachineHandler(dbSession, tc, scp, cfg),
		},
		{
			Path:    "/org/:orgName/carbide/machine/:id",
			Method:  http.MethodDelete,
			Handler: apiHandler.NewDeleteMachineHandler(dbSession, tc, cfg),
		},
		{

			Path:    "/org/:orgName/carbide/machine/:id/status-history",
			Method:  http.MethodGet,
			Handler: apiHandler.NewGetMachineStatusDetailsHandler(dbSession),
		},
		// Machine/Instance Type association endpoints
		{
			Path:    "/org/:orgName/carbide/instance/type/:instanceTypeId/machine",
			Method:  http.MethodPost,
			Handler: apiHandler.NewCreateMachineInstanceTypeHandler(dbSession, tc, scp, cfg),
		},
		{
			Path:    "/org/:orgName/carbide/instance/type/:instanceTypeId/machine",
			Method:  http.MethodGet,
			Handler: apiHandler.NewGetAllMachineInstanceTypeHandler(dbSession, tc, cfg),
		},
		{
			Path:    "/org/:orgName/carbide/instance/type/:instanceTypeId/machine/:id",
			Method:  http.MethodDelete,
			Handler: apiHandler.NewDeleteMachineInstanceTypeHandler(dbSession, tc, scp, cfg),
		},
		// Allocation endpoints
		{
			Path:    "/org/:orgName/carbide/allocation",
			Method:  http.MethodPost,
			Handler: apiHandler.NewCreateAllocationHandler(dbSession, tc, scp, cfg),
		},
		{
			Path:    "/org/:orgName/carbide/allocation",
			Method:  http.MethodGet,
			Handler: apiHandler.NewGetAllAllocationHandler(dbSession, tc, cfg),
		},
		{
			Path:    "/org/:orgName/carbide/allocation/:id",
			Method:  http.MethodGet,
			Handler: apiHandler.NewGetAllocationHandler(dbSession, tc, cfg),
		},
		{
			Path:    "/org/:orgName/carbide/allocation/:id",
			Method:  http.MethodPatch,
			Handler: apiHandler.NewUpdateAllocationHandler(dbSession, tc, cfg),
		},
		// AllocationConstraint update endpoint
		{
			Path:    "/org/:orgName/carbide/allocation/:allocationId/constraint/:id",
			Method:  http.MethodPatch,
			Handler: apiHandler.NewUpdateAllocationConstraintHandler(dbSession, tc, cfg),
		},
		{
			Path:    "/org/:orgName/carbide/allocation/:id",
			Method:  http.MethodDelete,
			Handler: apiHandler.NewDeleteAllocationHandler(dbSession, tc, cfg),
		},
		// Subnet endpoints
		{
			Path:    "/org/:orgName/carbide/subnet",
			Method:  http.MethodPost,
			Handler: apiHandler.NewCreateSubnetHandler(dbSession, tc, scp, cfg),
		},
		{
			Path:    "/org/:orgName/carbide/subnet",
			Method:  http.MethodGet,
			Handler: apiHandler.NewGetAllSubnetHandler(dbSession, tc, cfg),
		},
		{
			Path:    "/org/:orgName/carbide/subnet/:id",
			Method:  http.MethodGet,
			Handler: apiHandler.NewGetSubnetHandler(dbSession, tc, cfg),
		},
		{
			Path:    "/org/:orgName/carbide/subnet/:id",
			Method:  http.MethodPatch,
			Handler: apiHandler.NewUpdateSubnetHandler(dbSession, tc, cfg),
		},
		{
			Path:    "/org/:orgName/carbide/subnet/:id",
			Method:  http.MethodDelete,
			Handler: apiHandler.NewDeleteSubnetHandler(dbSession, tc, scp, cfg),
		},
		// OperatingSystem endpoints
		{
			Path:    "/org/:orgName/carbide/operating-system",
			Method:  http.MethodPost,
			Handler: apiHandler.NewCreateOperatingSystemHandler(dbSession, tc, scp, cfg),
		},
		{
			Path:    "/org/:orgName/carbide/operating-system",
			Method:  http.MethodGet,
			Handler: apiHandler.NewGetAllOperatingSystemHandler(dbSession, tc, cfg),
		},
		{
			Path:    "/org/:orgName/carbide/operating-system/:id",
			Method:  http.MethodGet,
			Handler: apiHandler.NewGetOperatingSystemHandler(dbSession, tc, cfg),
		},
		{
			Path:    "/org/:orgName/carbide/operating-system/:id",
			Method:  http.MethodPatch,
			Handler: apiHandler.NewUpdateOperatingSystemHandler(dbSession, tc, scp, cfg),
		},
		{
			Path:    "/org/:orgName/carbide/operating-system/:id",
			Method:  http.MethodDelete,
			Handler: apiHandler.NewDeleteOperatingSystemHandler(dbSession, tc, scp, cfg),
		},
		// NetworkSecurityGroup endpoints
		{
			Path:    "/org/:orgName/carbide/network-security-group",
			Method:  http.MethodPost,
			Handler: apiHandler.NewCreateNetworkSecurityGroupHandler(dbSession, tc, scp, cfg),
		},

		{
			Path:    "/org/:orgName/carbide/network-security-group",
			Method:  http.MethodGet,
			Handler: apiHandler.NewGetAllNetworkSecurityGroupHandler(dbSession, tc, cfg),
		},

		{
			Path:    "/org/:orgName/carbide/network-security-group/:id",
			Method:  http.MethodGet,
			Handler: apiHandler.NewGetNetworkSecurityGroupHandler(dbSession, tc, cfg),
		},

		{
			Path:    "/org/:orgName/carbide/network-security-group/:id",
			Method:  http.MethodPatch,
			Handler: apiHandler.NewUpdateNetworkSecurityGroupHandler(dbSession, tc, scp, cfg),
		},

		{
			Path:    "/org/:orgName/carbide/network-security-group/:id",
			Method:  http.MethodDelete,
			Handler: apiHandler.NewDeleteNetworkSecurityGroupHandler(dbSession, tc, scp, cfg),
		},

		// SSHKey endpoints
		{
			Path:    "/org/:orgName/carbide/sshkey",
			Method:  http.MethodPost,
			Handler: apiHandler.NewCreateSSHKeyHandler(dbSession, tc, cfg),
		},
		{
			Path:    "/org/:orgName/carbide/sshkey",
			Method:  http.MethodGet,
			Handler: apiHandler.NewGetAllSSHKeyHandler(dbSession, tc, cfg),
		},
		{
			Path:    "/org/:orgName/carbide/sshkey/:id",
			Method:  http.MethodGet,
			Handler: apiHandler.NewGetSSHKeyHandler(dbSession, tc, cfg),
		},
		{
			Path:    "/org/:orgName/carbide/sshkey/:id",
			Method:  http.MethodPatch,
			Handler: apiHandler.NewUpdateSSHKeyHandler(dbSession, tc, cfg),
		},
		{
			Path:    "/org/:orgName/carbide/sshkey/:id",
			Method:  http.MethodDelete,
			Handler: apiHandler.NewDeleteSSHKeyHandler(dbSession, tc, cfg),
		},
		// SSHKeyGroup endpoints
		{
			Path:    "/org/:orgName/carbide/sshkeygroup",
			Method:  http.MethodPost,
			Handler: apiHandler.NewCreateSSHKeyGroupHandler(dbSession, tc, cfg),
		},
		{
			Path:    "/org/:orgName/carbide/sshkeygroup",
			Method:  http.MethodGet,
			Handler: apiHandler.NewGetAllSSHKeyGroupHandler(dbSession, tc, cfg),
		},
		{
			Path:    "/org/:orgName/carbide/sshkeygroup/:id",
			Method:  http.MethodGet,
			Handler: apiHandler.NewGetSSHKeyGroupHandler(dbSession, tc, cfg),
		},
		{
			Path:    "/org/:orgName/carbide/sshkeygroup/:id",
			Method:  http.MethodPatch,
			Handler: apiHandler.NewUpdateSSHKeyGroupHandler(dbSession, tc, cfg),
		},
		{
			Path:    "/org/:orgName/carbide/sshkeygroup/:id",
			Method:  http.MethodDelete,
			Handler: apiHandler.NewDeleteSSHKeyGroupHandler(dbSession, tc, cfg),
		},
		// Machine Capability endpoints
		{
			Path:    "/org/:orgName/carbide/machine-capability",
			Method:  http.MethodGet,
			Handler: apiHandler.NewGetAllMachineCapabilityHandler(dbSession),
		},
		// Audit Log endpoints
		{
			Path:    "/org/:orgName/carbide/audit",
			Method:  http.MethodGet,
			Handler: apiHandler.NewGetAllAuditEntryHandler(dbSession),
		},
		{
			Path:    "/org/:orgName/carbide/audit/:id",
			Method:  http.MethodGet,
			Handler: apiHandler.NewGetAuditEntryHandler(dbSession),
		},
		// Machine Validation endpoints
		{
			Path:    "/org/:orgName/carbide/site/:siteID/machine-validation/test",
			Method:  http.MethodPost,
			Handler: apiHandler.NewCreateMachineValidationTestHandler(dbSession, tc, scp, cfg),
		},
		{
			Path:    "/org/:orgName/carbide/site/:siteID/machine-validation/test/:id/version/:version",
			Method:  http.MethodPatch,
			Handler: apiHandler.NewUpdateMachineValidationTestHandler(dbSession, tc, scp, cfg),
		},
		{
			Path:    "/org/:orgName/carbide/site/:siteID/machine-validation/test",
			Method:  http.MethodGet,
			Handler: apiHandler.NewGetAllMachineValidationTestHandler(dbSession, tc, scp, cfg),
		},
		{
			Path:    "/org/:orgName/carbide/site/:siteID/machine-validation/test/:id/version/:version",
			Method:  http.MethodGet,
			Handler: apiHandler.NewGetMachineValidationTestHandler(dbSession, tc, scp, cfg),
		},
		{
			Path:    "/org/:orgName/carbide/site/:siteID/machine-validation/machine/:machineID/results",
			Method:  http.MethodGet,
			Handler: apiHandler.NewGetMachineValidationResultsHandler(dbSession, tc, scp, cfg),
		},
		{
			Path:    "/org/:orgName/carbide/site/:siteID/machine-validation/machine/:machineID/runs",
			Method:  http.MethodGet,
			Handler: apiHandler.NewGetAllMachineValidationRunHandler(dbSession, tc, scp, cfg),
		},
		{
			Path:    "/org/:orgName/carbide/site/:siteID/machine-validation/external-config",
			Method:  http.MethodGet,
			Handler: apiHandler.NewGetAllMachineValidationExternalConfigHandler(dbSession, tc, scp, cfg),
		},
		{
			Path:    "/org/:orgName/carbide/site/:siteID/machine-validation/external-config/:cfgName",
			Method:  http.MethodGet,
			Handler: apiHandler.NewGetMachineValidationExternalConfigHandler(dbSession, tc, scp, cfg),
		},
		{
			Path:    "/org/:orgName/carbide/site/:siteID/machine-validation/external-config",
			Method:  http.MethodPost,
			Handler: apiHandler.NewCreateMachineValidationExternalConfigHandler(dbSession, tc, scp, cfg),
		},
		{
			Path:    "/org/:orgName/carbide/site/:siteID/machine-validation/external-config/:cfgName",
			Method:  http.MethodPatch,
			Handler: apiHandler.NewUpdateMachineValidationExternalConfigHandler(dbSession, tc, scp, cfg),
		},
		{
			Path:    "/org/:orgName/carbide/site/:siteID/machine-validation/external-config/:cfgName",
			Method:  http.MethodDelete,
			Handler: apiHandler.NewDeleteMachineValidationExternalConfigHandler(dbSession, tc, scp, cfg),
		},
		// DPU Extension Service endpoints
		{
			Path:    "/org/:orgName/carbide/dpu-extension-service",
			Method:  http.MethodPost,
			Handler: apiHandler.NewCreateDpuExtensionServiceHandler(dbSession, tc, scp, cfg),
		},
		{
			Path:    "/org/:orgName/carbide/dpu-extension-service",
			Method:  http.MethodGet,
			Handler: apiHandler.NewGetAllDpuExtensionServiceHandler(dbSession, tc, cfg),
		},
		{
			Path:    "/org/:orgName/carbide/dpu-extension-service/:id",
			Method:  http.MethodGet,
			Handler: apiHandler.NewGetDpuExtensionServiceHandler(dbSession, tc, cfg),
		},
		{
			Path:    "/org/:orgName/carbide/dpu-extension-service/:id",
			Method:  http.MethodPatch,
			Handler: apiHandler.NewUpdateDpuExtensionServiceHandler(dbSession, tc, scp, cfg),
		},
		{
			Path:    "/org/:orgName/carbide/dpu-extension-service/:id",
			Method:  http.MethodDelete,
			Handler: apiHandler.NewDeleteDpuExtensionServiceHandler(dbSession, tc, scp, cfg),
		},
		{
			Path:    "/org/:orgName/carbide/dpu-extension-service/:id/version/:version",
			Method:  http.MethodGet,
			Handler: apiHandler.NewGetDpuExtensionServiceVersionHandler(dbSession, tc, scp, cfg),
		},
		{
			Path:    "/org/:orgName/carbide/dpu-extension-service/:id/version/:version",
			Method:  http.MethodDelete,
			Handler: apiHandler.NewDeleteDpuExtensionServiceVersionHandler(dbSession, tc, scp, cfg),
		},
		// SKU endpoints
		{
			Path:    "/org/:orgName/carbide/sku",
			Method:  http.MethodGet,
			Handler: apiHandler.NewGetAllSkuHandler(dbSession, tc, cfg),
		},
		{
			Path:    "/org/:orgName/carbide/sku/:id",
			Method:  http.MethodGet,
			Handler: apiHandler.NewGetSkuHandler(dbSession, tc, cfg),
		},
	}

	return apiRoutes
}
