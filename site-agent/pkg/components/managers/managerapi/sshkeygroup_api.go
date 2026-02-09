package managerapi

import (
	wflows "github.com/nvidia/carbide-rest/workflow-schema/schema/site-agent/workflows/v1"
	"go.temporal.io/sdk/workflow"
)

// SSHKeyGroupExpansion - SSHKeyGroup Expansion
type SSHKeyGroupExpansion interface{}

// SSHKeyGroupInterface - interface to SSHKeyGroup
type SSHKeyGroupInterface interface {
	// List all the apis of SSHKeyGroup here
	Init()
	RegisterSubscriber() error
	RegisterPublisher() error
	RegisterCron() error

	// Cloud Workflow APIs
	CreateSSHKeyGroup(ctx workflow.Context, TransactionID *wflows.TransactionID, ResourceRequest *wflows.CreateSSHKeyGroupRequest) (err error)
	DeleteSSHKeyGroup(ctx workflow.Context, TransactionID *wflows.TransactionID, ResourceRequest *wflows.DeleteSSHKeyGroupRequest) (err error)

	// CRUD SSHKeyGroup APIs
	UpdateSSHKeyGroup(ctx workflow.Context, TransactionID *wflows.TransactionID, ResourceRequest *wflows.UpdateSSHKeyGroupRequest) (err error)
	// GetSSHKeyGroupByID(ctx workflow.Context, ResourceID string, SSHKeyGroupID string) (ResourceResponse *wflows.GetSSHKeyGroupResponse, err error)
	GetSSHKeyGroup(ctx workflow.Context, ResourceID string, ResourceRequest *wflows.GetSSHKeyGroup) (ResourceResponse *wflows.GetSSHKeyGroupResponse, err error)

	// CreateSSHKeyGroup
	// RegisterWorkflows() error
	// RegisterActivities() error
	GetState() []string
	SSHKeyGroupExpansion
}
