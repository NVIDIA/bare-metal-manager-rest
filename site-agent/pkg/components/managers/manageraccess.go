// SPDX-FileCopyrightText: Copyright (c) 2021-2023 NVIDIA CORPORATION & AFFILIATES. All rights reserved.
// SPDX-License-Identifier: LicenseRef-NvidiaProprietary
//
// NVIDIA CORPORATION, its affiliates and licensors retain all intellectual
// property and proprietary rights in and to this material, related
// documentation and any modifications thereto. Any use, reproduction,
// disclosure or distribution of this material and related documentation
// without an express license agreement from NVIDIA CORPORATION or
// its affiliates is strictly prohibited.

package managers

import (
	"github.com/nvidia/carbide-rest/site-agent/pkg/components/managers/bootstrap"
	"github.com/nvidia/carbide-rest/site-agent/pkg/components/managers/carbide"
	"github.com/nvidia/carbide-rest/site-agent/pkg/components/managers/dpuextensionservice"
	"github.com/nvidia/carbide-rest/site-agent/pkg/components/managers/expectedmachine"
	"github.com/nvidia/carbide-rest/site-agent/pkg/components/managers/health"
	"github.com/nvidia/carbide-rest/site-agent/pkg/components/managers/infinibandpartition"
	"github.com/nvidia/carbide-rest/site-agent/pkg/components/managers/instance"
	"github.com/nvidia/carbide-rest/site-agent/pkg/components/managers/instancetype"
	"github.com/nvidia/carbide-rest/site-agent/pkg/components/managers/machine"
	"github.com/nvidia/carbide-rest/site-agent/pkg/components/managers/machinevalidation"
	"github.com/nvidia/carbide-rest/site-agent/pkg/components/managers/managerapi"
	"github.com/nvidia/carbide-rest/site-agent/pkg/components/managers/networksecuritygroup"
	"github.com/nvidia/carbide-rest/site-agent/pkg/components/managers/nvlinklogicalpartition"
	"github.com/nvidia/carbide-rest/site-agent/pkg/components/managers/operatingsystem"
	"github.com/nvidia/carbide-rest/site-agent/pkg/components/managers/sku"
	"github.com/nvidia/carbide-rest/site-agent/pkg/components/managers/sshkeygroup"
	"github.com/nvidia/carbide-rest/site-agent/pkg/components/managers/subnet"
	"github.com/nvidia/carbide-rest/site-agent/pkg/components/managers/tenant"
	"github.com/nvidia/carbide-rest/site-agent/pkg/components/managers/vpc"
	"github.com/nvidia/carbide-rest/site-agent/pkg/components/managers/vpcprefix"
	"github.com/nvidia/carbide-rest/site-agent/pkg/components/managers/workflow"
)

// ManagerAccess - access to manager struct
var ManagerAccess *Manager

// Manager - Access to all APIs/data/conf in a single struct
type Manager struct {
	//nolint
	API  *managerapi.ManagerAPI
	Data *managerapi.ManagerData
	Conf *managerapi.ManagerConf
}

// Add all the Managers here
// Each manager has to register a new instance here for acceess

// Bootstrap Add bootstrap manager instance here
func (m *Manager) Bootstrap() *bootstrap.BoostrapAPI {
	return bootstrap.NewBootstrapManager(m.Data.EB, m.API, m.Conf)
}

// Orchestrator - Add orchestrator manager instance here
func (m *Manager) Orchestrator() *workflow.API {
	return workflow.NewWorkflowManager(m.Data.EB, m.API, m.Conf)
}

// VPC - Add vpc manager instance here
func (m *Manager) VPC() *vpc.API {
	return vpc.NewVPCManager(m.Data.EB, m.API, m.Conf)
}

// VpcPrefix - Add vpcprefix manager instance here
func (m *Manager) VpcPrefix() *vpcprefix.API {
	return vpcprefix.NewVpcPrefixManager(m.Data.EB, m.API, m.Conf)
}

// Carbide manager instance here
func (m *Manager) Carbide() *carbide.API {
	return carbide.NewCarbideManager(m.Data.EB, m.API, m.Conf)
}

// Machine - Add Machine manager instance here
func (m *Manager) Machine() *machine.API {
	return machine.NewMachineManager(m.Data.EB, m.API, m.Conf)
}

// Subnet - Add Subnet Manager instance here
func (m *Manager) Subnet() *subnet.API {
	return subnet.NewSubnetManager(m.Data.EB, m.API, m.Conf)
}

// Instance - Add Instance Manager instance here
func (m *Manager) Instance() *instance.API {
	return instance.NewInstanceManager(m.Data.EB, m.API, m.Conf)
}

// Health - Add Health Manager instance here
func (m *Manager) Health() *health.API {
	return health.NewHealthManager(m.Data.EB, m.API, m.Conf)
}

// SSHKeyGroup - Add SSHKeyGroup Manager instance here
func (m *Manager) SSHKeyGroup() *sshkeygroup.API {
	return sshkeygroup.NewSSHKeyGroupManager(m.Data.EB, m.API, m.Conf)
}

// InfiniBandPartition - Add InfiniBandPartition Manager instance here
func (m *Manager) InfiniBandPartition() *infinibandpartition.API {
	return infinibandpartition.NewInfiniBandPartitionManager(m.Data.EB, m.API, m.Conf)
}

// Tenant - Add Tenant Manager instance here
func (m *Manager) Tenant() *tenant.API {
	return tenant.NewTenantManager(m.Data.EB, m.API, m.Conf)
}

// OperatingSystem - Add OperatingSystem Manager instance here
func (m *Manager) OperatingSystem() *operatingsystem.API {
	return operatingsystem.NewOperatingSystemManager(m.Data.EB, m.API, m.Conf)
}

// MachineValidation - Add MachineValidation Manager instance here
func (m *Manager) MachineValidation() *machinevalidation.API {
	return machinevalidation.NewMachineValidationManager(m.Data.EB, m.API, m.Conf)
}

// InstanceType - Add InstanceType Manager instance here
func (m *Manager) InstanceType() *instancetype.API {
	return instancetype.NewInstanceTypeManager(m.Data.EB, m.API, m.Conf)
}

// NetworkSecurityGroup - Add NetworkSecurityGroup Manager instance here
func (m *Manager) NetworkSecurityGroup() *networksecuritygroup.API {
	return networksecuritygroup.NewNetworkSecurityGroupManager(m.Data.EB, m.API, m.Conf)
}

// ExpectedMachine - Add ExpectedMachine Manager instance here
func (m *Manager) ExpectedMachine() *expectedmachine.API {
	return expectedmachine.NewExpectedMachineManager(m.Data.EB, m.API, m.Conf)
}

// SKU - Add SKU Manager instance here
func (m *Manager) SKU() *sku.API {
	return sku.NewSKUManager(m.Data.EB, m.API, m.Conf)
}

// DpuExtensionService - Add DPU Extension Service Manager instance here
func (m *Manager) DpuExtensionService() *dpuextensionservice.API {
	return dpuextensionservice.NewDpuExtensionServiceManager(m.Data.EB, m.API, m.Conf)
}

// NVLinkLogicalPartition - Add NVLinkLogicalPartition Manager instance here
func (m *Manager) NVLinkLogicalPartition() *nvlinklogicalpartition.API {
	return nvlinklogicalpartition.NewNVLinkLogicalPartitionManager(m.Data.EB, m.API, m.Conf)
}
