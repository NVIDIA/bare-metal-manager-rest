package managerapi

import (
	"context"

	wflows "github.com/nvidia/carbide-rest/workflow-schema/schema/site-agent/workflows/v1"
	wflowtypes "github.com/nvidia/carbide-rest/site-agent/pkg/datatypes/managertypes/workflow"
	"go.temporal.io/sdk/workflow"
)

// OrchestratorExpansion - Orchestrator Expansion
type OrchestratorExpansion interface{}

// OrchestratorInterface - interface to Orchestrator
type OrchestratorInterface interface {
	// List all the apis of Orchestrator here
	Init()
	Start()
	GetState() []string
	AddWorkflow(wflow interface{})
	DoWorkflow(ctx workflow.Context, TransactionID *wflows.TransactionID,
		ResourceRequest interface{}, wflowMd wflowtypes.WorkflowMetadata,
		retryOptions *wflows.WorkflowOptions) (actErr error, pubErr error)
	DoActivity(ctx context.Context, ResourceVer uint64, ResourceID string,
		ResourceReq interface{}, wflowMd wflowtypes.WorkflowMetadata) (interface{}, error)

	OrchestratorExpansion
}
