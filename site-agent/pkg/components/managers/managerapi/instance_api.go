package managerapi

import (
	wflows "github.com/nvidia/carbide-rest/workflow-schema/schema/site-agent/workflows/v1"

	"go.temporal.io/sdk/workflow"
)

// InstanceExpansion - Instance Expansion
type InstanceExpansion interface{}

// InstanceInterface - interface to Instance
type InstanceInterface interface {
	// List all the apis of Instance here
	Init()
	RegisterSubscriber() error
	RegisterPublisher() error
	RegisterCron() error

	// Temporal Workflows - Subscriber
	//Create Instance (deprecated)
	CreateInstance(ctx workflow.Context, TransactionID *wflows.TransactionID, ResourceRequest *wflows.CreateInstanceRequest) (err error)
	//Delete Instance (deprecated)
	DeleteInstance(ctx workflow.Context, TransactionID *wflows.TransactionID, ResourceRequest *wflows.DeleteInstanceRequest) (err error)
	//RebootInstance (deprecated)
	RebootInstance(ctx workflow.Context, TransactionID *wflows.TransactionID, ResourceRequest *wflows.RebootInstanceRequest) (err error)
	GetState() []string

	InstanceExpansion
}
