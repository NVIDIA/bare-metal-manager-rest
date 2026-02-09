package client

import (
	wflows "github.com/nvidia/carbide-rest/workflow-schema/schema/site-agent/workflows/v1"
)

// NetworkGetter is the interface for the network workflows
type NetworkGetter interface {
	Network() NetworkInterface
}

// NetworkInterface is the interface for the network client
type NetworkInterface interface {

	// VPC Interface
	VPCInterface
	// Subnet Interface
	SubnetInterface
	// InfiniBandPartition Interface
	InfiniBandPartitionInterface
}

type network struct {
	// carbide client
	carbide wflows.ForgeClient
}

func newNetwork(carbide wflows.ForgeClient) *network {
	return &network{carbide: carbide}
}
