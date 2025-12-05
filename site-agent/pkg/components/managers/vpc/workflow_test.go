// SPDX-FileCopyrightText: Copyright (c) 2021-2023 NVIDIA CORPORATION & AFFILIATES. All rights reserved.
// SPDX-License-Identifier: LicenseRef-NvidiaProprietary
//
// NVIDIA CORPORATION, its affiliates and licensors retain all intellectual
// property and proprietary rights in and to this material, related
// documentation and any modifications thereto. Any use, reproduction,
// disclosure or distribution of this material and related documentation
// without an express license agreement from NVIDIA CORPORATION or
// its affiliates is strictly prohibited.

package vpc

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	wflows "github.com/nvidia/carbide-rest/workflow-schema/schema/site-agent/workflows/v1"
)

func Test_vpcWorkflowMetadata_ResponseState(t *testing.T) {
	createdVpc := &wflows.Vpc{
		Id: &wflows.VpcId{Value: uuid.NewString()},
	}

	type fields struct {
		activity activityType
		response *wflows.VPCInfo
		respList *wflows.GetVPCResponse
	}
	type args struct {
		status       wflows.WorkflowStatus
		objectStatus wflows.ObjectStatus
		statusMsg    string
	}
	tests := []struct {
		name     string
		fields   fields
		args     args
		wantResp *wflows.VPCInfo
	}{
		// TODO: Remove test for create (deprecated).  We've moved to sync workflow.
		{
			name: "test success response state for create",
			fields: fields{
				activity: activityCreate,
				response: &wflows.VPCInfo{
					Vpc: createdVpc,
				},
			},
			args: args{
				status:       wflows.WorkflowStatus_WORKFLOW_STATUS_SUCCESS,
				objectStatus: wflows.ObjectStatus_OBJECT_STATUS_CREATED,
				statusMsg:    "vpc was successfully created",
			},
			wantResp: &wflows.VPCInfo{
				Vpc:          createdVpc,
				Status:       wflows.WorkflowStatus_WORKFLOW_STATUS_SUCCESS,
				ObjectStatus: wflows.ObjectStatus_OBJECT_STATUS_CREATED,
				StatusMsg:    "vpc was successfully created",
			},
		},
		{
			name: "test failure response state for delete",
			fields: fields{
				activity: activityDelete,
				response: &wflows.VPCInfo{
					Vpc: nil,
				},
			},
			args: args{
				status:    wflows.WorkflowStatus_WORKFLOW_STATUS_FAILURE,
				statusMsg: "vpc deletion failed",
			},
			wantResp: &wflows.VPCInfo{
				Status:    wflows.WorkflowStatus_WORKFLOW_STATUS_FAILURE,
				StatusMsg: "vpc deletion failed",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := &vpcWorkflowMetadata{
				activity: tt.fields.activity,
				response: tt.fields.response,
				respList: tt.fields.respList,
			}
			v.ResponseState(tt.args.status, tt.args.objectStatus, tt.args.statusMsg)

			if tt.wantResp != nil {
				assert.Equal(t, tt.wantResp.Status, v.response.Status)
				assert.Equal(t, tt.wantResp.ObjectStatus, v.response.ObjectStatus)
				assert.Equal(t, tt.wantResp.StatusMsg, v.response.StatusMsg)
			}
		})
	}
}
