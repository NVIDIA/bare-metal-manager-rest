// SPDX-FileCopyrightText: Copyright (c) 2021-2023 NVIDIA CORPORATION & AFFILIATES. All rights reserved.
// SPDX-License-Identifier: LicenseRef-NvidiaProprietary
//
// NVIDIA CORPORATION, its affiliates and licensors retain all intellectual
// property and proprietary rights in and to this material, related
// documentation and any modifications thereto. Any use, reproduction,
// disclosure or distribution of this material and related documentation
// without an express license agreement from NVIDIA CORPORATION or
// its affiliates is strictly prohibited.

package managerapi

import (
	"github.com/nvidia/carbide-rest/site-agent/pkg/conftypes"
	"github.com/nvidia/carbide-rest/site-agent/pkg/datatypes/elektratypes"
)

// ManagerHdl - local handle to be assigned
var ManagerHdl ManagerAPI

// ManagerAccess - access to all APIs/data/conf
// nolint
type ManagerAccess struct {
	API  *ManagerAPI
	Data *ManagerData
	Conf *ManagerConf
}

// ManagerData - super struct
type ManagerData struct {
	EB *elektratypes.Elektra
}

// ManagerAPI struct to hold all mgr interface
type ManagerAPI struct {
	// Add all the manager interfaces here
	Bootstrap              BootstrapInterface
	VPC                    VPCInterface
	VpcPrefix              VpcPrefixInterface
	Subnet                 SubnetInterface
	Instance               InstanceInterface
	Machine                MachineInterface
	Orchestrator           OrchestratorInterface
	Carbide                CarbideInterface
	Health                 HealthInterface
	SSHKeyGroup            SSHKeyGroupInterface
	InfiniBandPartition    InfiniBandPartitionInterface
	Tenant                 TenantInterface
	OperatingSystem        OperatingSystemInterface
	MachineValidation      MachineValidationInterface
	InstanceType           InstanceTypeInterface
	NetworkSecurityGroup   NetworkSecurityGroupInterface
	ExpectedMachine        ExpectedMachineInterface
	SKU                    SKUInterface
	DpuExtensionService    DpuExtensionServiceInterface
	NVLinkLogicalPartition NVLinkLogicalPartitionInterface
}

// ManagerConf - Conf struct
type ManagerConf struct {
	EB *conftypes.Config
}
