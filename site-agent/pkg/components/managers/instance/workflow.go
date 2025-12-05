// SPDX-FileCopyrightText: Copyright (c) 2021-2023 NVIDIA CORPORATION & AFFILIATES. All rights reserved.
// SPDX-License-Identifier: LicenseRef-NvidiaProprietary
//
// NVIDIA CORPORATION, its affiliates and licensors retain all intellectual
// property and proprietary rights in and to this material, related
// documentation and any modifications thereto. Any use, reproduction,
// disclosure or distribution of this material and related documentation
// without an express license agreement from NVIDIA CORPORATION or
// its affiliates is strictly prohibited.

package instance

import (
	"context"
	"fmt"

	wflows "github.com/nvidia/carbide-rest/workflow-schema/schema/site-agent/workflows/v1"
	computils "github.com/nvidia/carbide-rest/site-agent/pkg/components/common"
	"github.com/nvidia/carbide-rest/site-agent/pkg/conftypes"
	workflowtypes "github.com/nvidia/carbide-rest/site-agent/pkg/datatypes/managertypes/workflow"
	"go.temporal.io/sdk/client"
	workflow "go.temporal.io/sdk/workflow"
)

type activityType int

const (
	activityCreate activityType = iota
	activityDelete
	activityReboot
	activityPublish
)

// TODO: Remove InstanceCreate (deprecated) and any related references.  We've moved to sync workflow
// TODO: Remove InstanceDelete (deprecated) and any related references once we've fully moved to sync workflow.
var activityStr = []string{"InstanceCreate", "InstanceDelete", "InstanceReboot",
	"InstanceGet", "InstancePublish"}

type instanceWorkflowMetadata struct {
	activity       activityType
	response       *wflows.InstanceInfo
	rebootResponse *wflows.InstanceRebootInfo
}

// ResourceType - Resource Type
func (w *instanceWorkflowMetadata) ResourceType() string {
	return computils.ResourceTypeInstance
}

// ActivityType - Activity Type
func (w *instanceWorkflowMetadata) ActivityType() (s string) {
	return activityStr[w.activity]
}

// DoDbOP - Do Db OP
func (w *instanceWorkflowMetadata) DoDbOP() (act computils.OpType) {
	switch w.activity {
	case activityCreate:
		act = computils.OpCreate
	case activityReboot:
		act = computils.OpUpdate
	case activityDelete:
		act = computils.OpDelete
	}
	return
}

// DoSiteControllerOP - Do Site Controller OP
func (w *instanceWorkflowMetadata) DoSiteControllerOP(ctx context.Context, transactionID *wflows.TransactionID, req interface{}) (interface{}, error) {
	compute := ManagerAccess.Data.EB.Managers.Carbide.GetClient().Compute()
	switch w.activity {
	case activityCreate:
		return compute.CreateInstance(ctx, req.(*wflows.CreateInstanceRequest))
	case activityReboot:
		return compute.RebootInstance(ctx, req.(*wflows.RebootInstanceRequest))
	case activityDelete:
		return compute.DeleteInstance(ctx, req.(*wflows.DeleteInstanceRequest))
	default:
		panic(fmt.Sprintf("invalid activity type: %v", w.activity))
	}
}

// ActivityInvoke - Activity Invoke
func (w *instanceWorkflowMetadata) ActivityInvoke() (act interface{}) {
	instance := &Workflows{}
	switch w.activity {
	case activityCreate:
		act = instance.CreateInstanceActivity
	case activityDelete:
		act = instance.DeleteInstanceActivity
	case activityReboot:
		act = instance.RebootInstanceActivity
	}
	return
}

// ActivityPublish Get the Instance publish activity
func (w *instanceWorkflowMetadata) ActivityPublish() interface{} {
	switch w.activity {
	case activityCreate, activityDelete:
		instance := Workflows{}
		return instance.PublishInstanceActivity
	case activityReboot:
		instance := Workflows{}
		return instance.PublishInstancePowerStatus
	default:
		panic(fmt.Sprintf("invalid activity type: %v", w.activity))
	}
}

// ResponseState - Response State
func (w *instanceWorkflowMetadata) ResponseState(status wflows.WorkflowStatus, objectStatus wflows.ObjectStatus, statusMsg string) {
	switch w.activity {
	case activityCreate, activityDelete:
		w.response.StatusMsg = statusMsg
		w.response.ObjectStatus = objectStatus
		w.response.Status = status
	case activityReboot:
		w.rebootResponse.StatusMsg = statusMsg
		w.rebootResponse.ObjectStatus = objectStatus
		w.rebootResponse.Status = status
	default:
		panic(fmt.Sprintf("invalid activity type: %v", w.activity))
	}
}

// Response workflow
func (w *instanceWorkflowMetadata) Response() interface{} {
	switch w.activity {
	case activityCreate, activityDelete:
		return w.response
	case activityReboot:
		return w.rebootResponse
	default:
		panic(fmt.Sprintf("invalid activity type: %v", w.activity))
	}
}

// Statistics instance
func (w *instanceWorkflowMetadata) Statistics() *workflowtypes.MgrState {
	return ManagerAccess.Data.EB.Managers.Workflow.InstanceState
}

// Workflows Temporal registration
type Workflows struct {
	tcPublish   client.Client
	tcSubscribe client.Client
	cfg         *conftypes.Config
}

// NewInstanceWorkflows creates an instance for InstanceWorkflows
func NewInstanceWorkflows(tmPublish client.Client, tmSubscribe client.Client, currentCFG *conftypes.Config) Workflows {
	return Workflows{
		tcPublish:   tmPublish,
		tcSubscribe: tmSubscribe,
		cfg:         currentCFG,
	}
}

// CreateInstance - temporal create instance workflow
func (instance *API) CreateInstance(ctx workflow.Context, transactionID *wflows.TransactionID, resourceRequest *wflows.CreateInstanceRequest) (err error) {
	wflowMd := &instanceWorkflowMetadata{
		activity: activityCreate,
		response: &wflows.InstanceInfo{Status: wflows.WorkflowStatus_WORKFLOW_STATUS_FAILURE,
			ObjectStatus: wflows.ObjectStatus_OBJECT_STATUS_UNSPECIFIED,
			Instance:     &wflows.Instance{},
		},
	}
	return doWflow(ctx, transactionID, resourceRequest, wflowMd, resourceRequest.Options)
}

// DeleteInstance - temporal delete Instance workflow
func (instance *API) DeleteInstance(ctx workflow.Context, transactionID *wflows.TransactionID, resourceRequest *wflows.DeleteInstanceRequest) error {
	wflowMd := &instanceWorkflowMetadata{
		activity: activityDelete,
		response: &wflows.InstanceInfo{Status: wflows.WorkflowStatus_WORKFLOW_STATUS_FAILURE,
			ObjectStatus: wflows.ObjectStatus_OBJECT_STATUS_UNSPECIFIED,
			Instance:     &wflows.Instance{},
		},
	}
	return doWflow(ctx, transactionID, resourceRequest.InstanceId.Value, wflowMd, resourceRequest.Options)
}

// RebootInstance - temporal reboot Instance workflow
func (instance *API) RebootInstance(ctx workflow.Context, transactionID *wflows.TransactionID, resourceRequest *wflows.RebootInstanceRequest) error {
	wflowMd := &instanceWorkflowMetadata{
		activity: activityReboot,
		rebootResponse: &wflows.InstanceRebootInfo{Status: wflows.WorkflowStatus_WORKFLOW_STATUS_FAILURE,
			ObjectStatus: wflows.ObjectStatus_OBJECT_STATUS_UNSPECIFIED,
			MachineId:    &wflows.MachineId{},
		},
	}
	return doWflow(ctx, transactionID, resourceRequest, wflowMd, resourceRequest.Options)
}

func doWflow(ctx workflow.Context, transactionID *wflows.TransactionID, resourceRequest interface{}, wflowMetadata *instanceWorkflowMetadata, retryOptions *wflows.WorkflowOptions) error {
	wflowMetadata.response = &wflows.InstanceInfo{
		Status:   wflows.WorkflowStatus_WORKFLOW_STATUS_FAILURE,
		Instance: &wflows.Instance{},
	}

	actErr, pubErr := ManagerAccess.API.Orchestrator.DoWorkflow(ctx, transactionID, resourceRequest, wflowMetadata, retryOptions)
	if actErr != nil {
		return actErr
	}
	return pubErr
}
