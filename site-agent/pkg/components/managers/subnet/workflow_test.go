// SPDX-FileCopyrightText: Copyright (c) 2021-2023 NVIDIA CORPORATION & AFFILIATES. All rights reserved.
// SPDX-License-Identifier: LicenseRef-NvidiaProprietary
//
// NVIDIA CORPORATION, its affiliates and licensors retain all intellectual
// property and proprietary rights in and to this material, related
// documentation and any modifications thereto. Any use, reproduction,
// disclosure or distribution of this material and related documentation
// without an express license agreement from NVIDIA CORPORATION or
// its affiliates is strictly prohibited.

package subnet

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	wflows "github.com/nvidia/carbide-rest/workflow-schema/schema/site-agent/workflows/v1"
)

func Test_subnetWorkflowMetadata_ResponseState(t *testing.T) {
	createdNs := &wflows.NetworkSegment{
		Id: &wflows.NetworkSegmentId{Value: uuid.NewString()},
	}

	type fields struct {
		activity activityType
		response *wflows.SubnetInfo
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
		wantResp *wflows.SubnetInfo
	}{
		{
			// TODO: Remove test for create subnet (deprecated)
			name: "test success response state for create",
			fields: fields{
				activity: activityCreate,
				response: &wflows.SubnetInfo{
					NetworkSegment: createdNs,
				},
			},
			args: args{
				status:       wflows.WorkflowStatus_WORKFLOW_STATUS_SUCCESS,
				objectStatus: wflows.ObjectStatus_OBJECT_STATUS_CREATED,
				statusMsg:    "subnet was successfully created",
			},
			wantResp: &wflows.SubnetInfo{
				NetworkSegment: createdNs,
				Status:         wflows.WorkflowStatus_WORKFLOW_STATUS_SUCCESS,
				ObjectStatus:   wflows.ObjectStatus_OBJECT_STATUS_CREATED,
				StatusMsg:      "subnet was successfully created",
			},
		},
		{
			name: "test failure response state for delete",
			fields: fields{
				activity: activityDelete,
				response: &wflows.SubnetInfo{
					NetworkSegment: createdNs,
				},
			},
			args: args{
				status:    wflows.WorkflowStatus_WORKFLOW_STATUS_FAILURE,
				statusMsg: "subnet deletion failed",
			},
			wantResp: &wflows.SubnetInfo{
				Status:    wflows.WorkflowStatus_WORKFLOW_STATUS_FAILURE,
				StatusMsg: "subnet deletion failed",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &subnetWorkflowMetadata{
				activity: tt.fields.activity,
				response: tt.fields.response,
			}
			w.ResponseState(tt.args.status, tt.args.objectStatus, tt.args.statusMsg)

			if tt.wantResp != nil {
				assert.Equal(t, tt.wantResp.Status, w.response.Status)
				assert.Equal(t, tt.wantResp.ObjectStatus, w.response.ObjectStatus)
				assert.Equal(t, tt.wantResp.StatusMsg, w.response.StatusMsg)
			}
		})
	}
}
