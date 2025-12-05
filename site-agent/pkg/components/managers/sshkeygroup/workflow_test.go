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
	"testing"

	"github.com/stretchr/testify/assert"

	wflows "github.com/nvidia/carbide-rest/workflow-schema/schema/site-agent/workflows/v1"
)

func Test_sshkgWorkflowMetadata_ResponseState(t *testing.T) {
	createdTk := &wflows.TenantKeyset{
		KeysetIdentifier: &wflows.TenantKeysetIdentifier{
			OrganizationId: "test-org",
			KeysetId:       "test-keyset-1",
		},
	}

	type fields struct {
		activity activityType
		response *wflows.SSHKeyGroupInfo
		respList *wflows.GetSSHKeyGroupResponse
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
		wantResp *wflows.SSHKeyGroupInfo
	}{
		{
			name: "test success response state for create",
			fields: fields{
				activity: activityCreate,
				response: &wflows.SSHKeyGroupInfo{
					TenantKeyset: createdTk,
				},
			},
			args: args{
				status:       wflows.WorkflowStatus_WORKFLOW_STATUS_SUCCESS,
				objectStatus: wflows.ObjectStatus_OBJECT_STATUS_CREATED,
				statusMsg:    "ssh key group was successfully created",
			},
			wantResp: &wflows.SSHKeyGroupInfo{
				TenantKeyset: createdTk,
				Status:       wflows.WorkflowStatus_WORKFLOW_STATUS_SUCCESS,
				ObjectStatus: wflows.ObjectStatus_OBJECT_STATUS_CREATED,
				StatusMsg:    "ssh key group was successfully created",
			},
		},
		{
			name: "test failure response state for delete",
			fields: fields{
				activity: activityDelete,
				response: &wflows.SSHKeyGroupInfo{
					TenantKeyset: nil,
				},
			},
			args: args{
				status:    wflows.WorkflowStatus_WORKFLOW_STATUS_FAILURE,
				statusMsg: "ssh key group deletion failed",
			},
			wantResp: &wflows.SSHKeyGroupInfo{
				Status:    wflows.WorkflowStatus_WORKFLOW_STATUS_FAILURE,
				StatusMsg: "ssh key group deletion failed",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			skgm := &sshkgWorkflowMetadata{
				activity: tt.fields.activity,
				response: tt.fields.response,
				respList: tt.fields.respList,
			}
			skgm.ResponseState(tt.args.status, tt.args.objectStatus, tt.args.statusMsg)

			if tt.wantResp != nil {
				assert.Equal(t, tt.wantResp.Status, skgm.response.Status)
				assert.Equal(t, tt.wantResp.ObjectStatus, skgm.response.ObjectStatus)
				assert.Equal(t, tt.wantResp.StatusMsg, skgm.response.StatusMsg)
			}
		})
	}
}
