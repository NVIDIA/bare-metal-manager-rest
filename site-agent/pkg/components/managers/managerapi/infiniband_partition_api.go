package managerapi

import (
	wflows "github.com/nvidia/carbide-rest/workflow-schema/schema/site-agent/workflows/v1"
	"go.temporal.io/sdk/workflow"
)

// InfiniBandPartitionExpansion - InfiniBandPartition Expansion
type InfiniBandPartitionExpansion interface{}

// InfiniBandPartitionInterface - interface to InfiniBandPartition
type InfiniBandPartitionInterface interface {
	// List all the apis of InfiniBandPartition here
	Init()
	RegisterSubscriber() error
	RegisterPublisher() error
	RegisterCron() error

	// Cloud Workflow APIs
	CreateInfiniBandPartition(ctx workflow.Context, TransactionID *wflows.TransactionID, ResourceRequest *wflows.CreateInfiniBandPartitionRequest) (err error)
	DeleteInfiniBandPartition(ctx workflow.Context, TransactionID *wflows.TransactionID, ResourceRequest *wflows.DeleteInfiniBandPartitionRequest) (err error)

	// CreateInfiniBandPartition
	// RegisterWorkflows() error
	// RegisterActivities() error
	GetState() []string
	InfiniBandPartitionExpansion
}
