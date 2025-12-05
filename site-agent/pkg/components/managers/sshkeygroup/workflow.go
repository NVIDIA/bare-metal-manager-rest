// SPDX-FileCopyrightText: Copyright (c) 2021-2023 NVIDIA CORPORATION & AFFILIATES. All rights reserved.
// SPDX-License-Identifier: LicenseRef-NvidiaProprietary
//
// NVIDIA CORPORATION, its affiliates and licensors retain all intellectual
// property and proprietary rights in and to this material, related
// documentation and any modifications thereto. Any use, reproduction,
// disclosure or distribution of this material and related documentation
// without an express license agreement from NVIDIA CORPORATION or
// its affiliates is strictly prohibited.

package sshkeygroup

import (
	"context"
	"fmt"

	"time"

	wflows "github.com/nvidia/carbide-rest/workflow-schema/schema/site-agent/workflows/v1"
	common "github.com/nvidia/carbide-rest/site-agent/pkg/components/common"
	workflowtypes "github.com/nvidia/carbide-rest/site-agent/pkg/datatypes/managertypes/workflow"

	"github.com/nvidia/carbide-rest/site-agent/pkg/conftypes"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/workflow"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type activityType int

const (
	activityCreate activityType = iota
	activityUpdate
	activityDelete
	activityGet
	activityPublish
)

var activityStr = []string{"sshkeygroupCreate", "sshkeygroupUpdate", "sshkeygroupDelete", "sshkeygroupGetByName", "sshkeygroupGetList", "sshkeygroupCollectInventory", "sshkeygroupPublish"}

type sshkgWorkflowMetadata struct {
	activity activityType
	response *wflows.SSHKeyGroupInfo
	respList *wflows.GetSSHKeyGroupResponse
	respInv  *wflows.SSHKeyGroupInventory
}

// Todo: Reconcile with DB interfaces
// ResourceType - Resource Type
func (skgm *sshkgWorkflowMetadata) ResourceType() string {
	return common.ResourceTypeSSHKeyGroup
}

// ActivityType - Activity Type
func (skgm *sshkgWorkflowMetadata) ActivityType() string {
	return activityStr[skgm.activity]
}

// DoDbOP - Do Db OP
func (skgm *sshkgWorkflowMetadata) DoDbOP() (act common.OpType) {
	switch skgm.activity {
	case activityCreate:
		act = common.OpCreate
	case activityUpdate:
		act = common.OpUpdate
	case activityDelete:
		act = common.OpDelete
	case activityGet:
		act = common.OpNone
	}
	return
}

// Todo: Reconcile with gRPC & wflow interfaces
// DoSiteControllerOP - Do Site Controller OP
func (skgm *sshkgWorkflowMetadata) DoSiteControllerOP(ctx context.Context,
	TransactionID *wflows.TransactionID, req interface{}) (interface{}, error) {
	switch skgm.activity {
	case activityCreate:
		return ManagerAccess.Data.EB.Managers.Carbide.GetClient().Compute().CreateSSHKeyGroup(ctx, req.(*wflows.CreateSSHKeyGroupRequest))
	case activityUpdate:
		return ManagerAccess.Data.EB.Managers.Carbide.GetClient().Compute().UpdateSSHKeyGroup(ctx, req.(*wflows.UpdateSSHKeyGroupRequest))
	case activityDelete:
		return ManagerAccess.Data.EB.Managers.Carbide.GetClient().Compute().DeleteSSHKeyGroup(ctx, req.(*wflows.DeleteSSHKeyGroupRequest))
	case activityGet:
		return ManagerAccess.Data.EB.Managers.Carbide.GetClient().Compute().GetSSHKeyGroup(ctx, req.(*wflows.GetSSHKeyGroup))
	default:
		panic(fmt.Sprintf("invalid activity type: %v", skgm.activity))
	}
}

// ActivityInvoke - Activity Invoke
func (skgm *sshkgWorkflowMetadata) ActivityInvoke() (act interface{}) {
	sshkgwflowinstance := &Workflows{}
	switch skgm.activity {
	case activityCreate:
		act = sshkgwflowinstance.CreateSSHKeyGroupActivity
	case activityUpdate:
		act = sshkgwflowinstance.UpdateSSHKeyGroupActivity
	case activityDelete:
		act = sshkgwflowinstance.DeleteSSHKeyGroupActivity
	case activityGet:
		act = sshkgwflowinstance.GetSSHKeyGroupActivity
		// No need for this activity yet
		// case activityGetList:
		// 	act = sshkgwflowinstance.CollectSSHKeyGroupListActivity
	}

	return
}

// ActivityPublish - Get the SSHKeyGroup publish activity
func (skgm *sshkgWorkflowMetadata) ActivityPublish() interface{} {
	sshkgwflowinstance := &Workflows{}
	switch skgm.activity {
	case activityCreate, activityUpdate, activityDelete:
		return sshkgwflowinstance.PublishSSHKeyGroupActivity
	case activityGet:
		return sshkgwflowinstance.PublishSSHKeyGroupActivity
	default:
		panic(fmt.Sprintf("invalid activity type: %v", skgm.activity))
	}
}

// ResponseState - Response State
func (skgm *sshkgWorkflowMetadata) ResponseState(status wflows.WorkflowStatus, objectStatus wflows.ObjectStatus, statusMsg string) {
	switch skgm.activity {
	case activityCreate, activityUpdate, activityDelete:
		skgm.response.StatusMsg = statusMsg
		skgm.response.ObjectStatus = objectStatus
		skgm.response.Status = status
	default:
		panic(fmt.Sprintf("invalid activity type: %v", skgm.activity))
	}

}

// Response object
func (skgm *sshkgWorkflowMetadata) Response() interface{} {
	switch skgm.activity {
	case activityCreate, activityUpdate, activityDelete:
		return skgm.response
	case activityGet:
		return skgm.respList
	default:
		panic(fmt.Sprintf("invalid activity type: %v", skgm.activity))
	}
}

// Statistics - SSHKeyGroup Stats
func (skgm *sshkgWorkflowMetadata) Statistics() *workflowtypes.MgrState {
	// Todo: Add stats here
	//return ManagerAccess.Data.EB.Managers.Workflow.sshkeygroupState
	return &workflowtypes.MgrState{}
}

// Workflows - Temporal registration
type Workflows struct {
	tcPublish   client.Client
	tcSubscribe client.Client
	cfg         *conftypes.Config
}

// NewSSHKeyGroupWorkflows creates an instance for SSHKeyGroupWorkflows
func NewSSHKeyGroupWorkflows(TMPublish client.Client, TMSubscribe client.Client, CurrentCFG *conftypes.Config) Workflows {
	return Workflows{
		tcPublish:   TMPublish,
		tcSubscribe: TMSubscribe,
		cfg:         CurrentCFG,
	}
}

// CreateSSHKeyGroup creates a new SSHKeyGroupWorkflow and publishes result to cloud
func (api *API) CreateSSHKeyGroup(ctx workflow.Context, TransactionID *wflows.TransactionID, ResourceRequest *wflows.CreateSSHKeyGroupRequest) (err error) {
	wflowMetadata := &sshkgWorkflowMetadata{activity: activityCreate}
	return doWflow(ctx, TransactionID, ResourceRequest, wflowMetadata, ResourceRequest.Options)
}

// UpdateSSHKeyGroup updates a SSHKeyGroupWorkflow
func (api *API) UpdateSSHKeyGroup(ctx workflow.Context, TransactionID *wflows.TransactionID, ResourceRequest *wflows.UpdateSSHKeyGroupRequest) (err error) {
	wflowMetadata := &sshkgWorkflowMetadata{activity: activityUpdate}
	return doWflow(ctx, TransactionID, ResourceRequest, wflowMetadata, ResourceRequest.Options)
}

// DeleteSSHKeyGroup deletes a SSHKeyGroup
func (api *API) DeleteSSHKeyGroup(ctx workflow.Context, TransactionID *wflows.TransactionID, ResourceRequest *wflows.DeleteSSHKeyGroupRequest) (err error) {
	wflowMetadata := &sshkgWorkflowMetadata{activity: activityDelete}
	return doWflow(ctx, TransactionID, ResourceRequest, wflowMetadata, ResourceRequest.Options)
}

// GetSSHKeyGroup - retrieve SSHKeyGroups
func (api *API) GetSSHKeyGroup(ctx workflow.Context, ResourceID string, ResourceRequest *wflows.GetSSHKeyGroup) (ResourceResponse *wflows.GetSSHKeyGroupResponse, err error) {
	transaction := &wflows.TransactionID{
		ResourceId: common.ResourceTypeSSHKeyGroup,
		Timestamp: &timestamppb.Timestamp{
			Seconds: time.Now().Unix(),
		},
	}
	resourceReq := &wflows.GetSSHKeyGroup{}
	wflowMetadata := &sshkgWorkflowMetadata{activity: activityGet}
	err = doWflow(ctx, transaction, resourceReq, wflowMetadata, &wflows.WorkflowOptions{})
	return wflowMetadata.respList, err
}

func doWflow(ctx workflow.Context, transactionID *wflows.TransactionID, resourceRequest interface{}, wflowMetadata *sshkgWorkflowMetadata, retryOptions *wflows.WorkflowOptions) error {
	wflowMetadata.response = &wflows.SSHKeyGroupInfo{
		Status:       wflows.WorkflowStatus_WORKFLOW_STATUS_FAILURE,
		TenantKeyset: &wflows.TenantKeyset{},
	}
	wflowMetadata.respList = &wflows.GetSSHKeyGroupResponse{List: &wflows.TenantKeySetList{}}
	wflowMetadata.respInv = &wflows.SSHKeyGroupInventory{
		InventoryStatus: wflows.InventoryStatus_INVENTORY_STATUS_FAILED,
	}

	actErr, pubErr := ManagerAccess.API.Orchestrator.DoWorkflow(ctx, transactionID, resourceRequest, wflowMetadata, retryOptions)
	if actErr != nil {
		return actErr
	}
	return pubErr
}
