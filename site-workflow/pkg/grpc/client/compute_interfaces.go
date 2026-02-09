package client

import wflows "github.com/nvidia/carbide-rest/workflow-schema/schema/site-agent/workflows/v1"

// ComputeGetter is the interface for compute workflows
type ComputeGetter interface {
	Compute() ComputeInterface
}

// ComputeInterface for machine gRPC apis
type ComputeInterface interface {
	MachineInterface
	// Instance Interface
	InstanceInterface
	// SSHKeyGroup Interface
	SSHKeyGroupInterface
	// OperatingSystem Interface
	OperatingSystemInterface
	// Tenant Interface
	TenantInterface
}

type compute struct {
	// carbide client
	carbide wflows.ForgeClient
}

func newCompute(carbide wflows.ForgeClient) *compute {
	return &compute{carbide: carbide}
}
