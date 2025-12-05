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
	"context"
	"errors"

	wflows "github.com/nvidia/carbide-rest/workflow-schema/schema/site-agent/workflows/v1"
	"go.temporal.io/sdk/temporal"
)

// CreateSubnetActivity - temporal activity
// TODO: Remove (deprecated)
func (ac *Workflows) CreateSubnetActivity(ctx context.Context, ResourceVer uint64, ResourceID string,
	ResourceReq *wflows.CreateSubnetRequest) (*wflows.SubnetInfo, error) {
	// Initialize logger
	logger := ManagerAccess.Data.EB.Log.With().Str("Activity", "CreateSubnetActivity").Str("ResourceID", ResourceID).Logger()
	logger.Info().Msg("Starting Activity")

	wflowMd := &subnetWorkflowMetadata{activity: activityCreate,
		response: &wflows.SubnetInfo{NetworkSegment: &wflows.NetworkSegment{
			Id:    &wflows.NetworkSegmentId{Value: ResourceReq.SubnetId.Value},
			VpcId: &wflows.VpcId{Value: ResourceReq.VpcId.Value},
			Name:  ResourceReq.Name,
		}},
	}

	response, err := ManagerAccess.API.Orchestrator.DoActivity(ctx, ResourceVer, ResourceID, ResourceReq, wflowMd)
	ResourceResp := wflowMd.response
	if err != nil {
		logger.Error().Err(err).Msg("Error creating subnet on site via Orchestrator")
		return ResourceResp, err
	}

	if resp, ok := response.(*wflows.NetworkSegment); ok && resp != nil {
		ResourceResp.NetworkSegment = resp
		logger.Info().Msg("Successfully completed activity")
		return ResourceResp, nil
	}

	err = errors.New("invalid or empty response received from Site Controller")
	logger.Error().Err(err).Msg("Failed to create subnet, invalid or empty response")
	return nil, err
}

// DeleteSubnetActivity - temporal activity
func (ac *Workflows) DeleteSubnetActivity(ctx context.Context, ResourceVer uint64, ResourceID string,
	ResourceReq string) (*wflows.SubnetInfo, error) {
	// Initialize logger
	logger := ManagerAccess.Data.EB.Log.With().Str("Activity", "DeleteSubnetActivity").Str("ResourceID", ResourceID).Logger()
	logger.Info().Msg("Starting Activity")

	wflowMd := &subnetWorkflowMetadata{
		activity: activityDelete,
		response: &wflows.SubnetInfo{
			NetworkSegment: &wflows.NetworkSegment{
				Id: &wflows.NetworkSegmentId{
					Value: ResourceReq,
				},
			}},
	}
	var err error
	if ResourceReq == "" {
		err = errors.New("invalid or empty network segment ID provided as activity argument")
		logger.Error().Err(err).Msg("Failed to delete subnet, invalid or empty network segment ID")
		wflowMd.response.StatusMsg = err.Error()
		wflowMd.response.Status = wflows.WorkflowStatus_WORKFLOW_STATUS_FAILURE
		return nil, temporal.NewApplicationError(err.Error(), "", &wflowMd.response)
	}

	_, err = ManagerAccess.API.Orchestrator.DoActivity(ctx, ResourceVer, ResourceID, &wflows.DeleteSubnetRequest{NetworkSegmentId: &wflows.UUID{Value: ResourceReq}}, wflowMd)
	if err != nil {
		logger.Error().Err(err).Msg("Error deleting subnet from site via Orchestrator")
		wflowMd.response.StatusMsg = err.Error()
		wflowMd.response.Status = wflows.WorkflowStatus_WORKFLOW_STATUS_FAILURE
		return wflowMd.response, err
	}

	logger.Info().Msg("Successfully completed activity")
	wflowMd.response.Status = wflows.WorkflowStatus_WORKFLOW_STATUS_SUCCESS // Ensure response status is updated accordingly
	return wflowMd.response, nil
}
