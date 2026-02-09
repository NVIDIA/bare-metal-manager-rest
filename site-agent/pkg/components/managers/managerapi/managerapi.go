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
	RLA                    RLAInterface
}

// ManagerConf - Conf struct
type ManagerConf struct {
	EB *conftypes.Config
}
