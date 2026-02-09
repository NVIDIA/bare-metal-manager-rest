// SPDX-FileCopyrightText: Copyright (c) 2021-2023 NVIDIA CORPORATION & AFFILIATES. All rights reserved.
// SPDX-License-Identifier: LicenseRef-NvidiaProprietary
//
// NVIDIA CORPORATION, its affiliates and licensors retain all intellectual
// property and proprietary rights in and to this material, related
// documentation and any modifications thereto. Any use, reproduction,
// disclosure or distribution of this material and related documentation
// without an express license agreement from NVIDIA CORPORATION or
// its affiliates is strictly prohibited.

package activity

import (
	"context"
	"errors"

	"github.com/rs/zerolog/log"

	swe "github.com/nvidia/carbide-rest/site-workflow/pkg/error"
	cClient "github.com/nvidia/carbide-rest/site-workflow/pkg/grpc/client"
	rlav1 "github.com/nvidia/carbide-rest/workflow-schema/rla/protobuf/v1"
	"go.temporal.io/sdk/temporal"
)

// ManageTray is an activity wrapper for Tray management via RLA
type ManageTray struct {
	RlaAtomicClient *cClient.RlaAtomicClient
}

// NewManageTray returns a new ManageTray client
func NewManageTray(rlaClient *cClient.RlaAtomicClient) ManageTray {
	return ManageTray{
		RlaAtomicClient: rlaClient,
	}
}

// GetTray retrieves a tray by its UUID from RLA
func (mt *ManageTray) GetTray(ctx context.Context, request *rlav1.GetComponentInfoByIDRequest) (*rlav1.GetComponentInfoResponse, error) {
	logger := log.With().Str("Activity", "GetTray").Logger()
	logger.Info().Msg("Starting activity")

	var err error

	// Validate request
	switch {
	case request == nil:
		err = errors.New("received empty get tray request")
	case request.Id == nil || request.Id.Id == "":
		err = errors.New("received get tray request missing tray ID")
	}

	if err != nil {
		return nil, temporal.NewNonRetryableApplicationError(err.Error(), swe.ErrTypeInvalidRequest, err)
	}

	// Call RLA gRPC endpoint
	rlaClient := mt.RlaAtomicClient.GetClient()
	rla := rlaClient.Rla()

	response, err := rla.GetComponentInfoByID(ctx, request)
	if err != nil {
		logger.Warn().Err(err).Msg("Failed to get tray by ID using RLA API")
		return nil, swe.WrapErr(err)
	}

	logger.Info().Msg("Completed activity")

	return response, nil
}

// GetTrays retrieves a list of trays from RLA with optional filters
func (mt *ManageTray) GetTrays(ctx context.Context, request *rlav1.GetComponentsRequest) (*rlav1.GetComponentsResponse, error) {
	logger := log.With().Str("Activity", "GetTrays").Logger()
	logger.Info().Msg("Starting activity")

	// Request can be nil or empty for getting all trays
	if request == nil {
		request = &rlav1.GetComponentsRequest{}
	}

	// Call RLA gRPC endpoint
	rlaClient := mt.RlaAtomicClient.GetClient()
	rla := rlaClient.Rla()

	response, err := rla.GetComponents(ctx, request)
	if err != nil {
		logger.Warn().Err(err).Msg("Failed to get list of trays using RLA API")
		return nil, swe.WrapErr(err)
	}

	logger.Info().Int32("total", response.GetTotal()).Msg("Completed activity")

	return response, nil
}

// GetTaskComponentIDs retrieves component (tray) UUIDs for a task via RLA GetTasksByIDs.
// Returns the component_uuids from the task so the tray API can filter by taskId.
func (mt *ManageTray) GetTaskComponentIDs(ctx context.Context, taskID string) ([]string, error) {
	logger := log.With().Str("Activity", "GetTaskComponentIDs").Str("TaskID", taskID).Logger()
	logger.Info().Msg("Starting activity")

	if taskID == "" {
		return nil, temporal.NewNonRetryableApplicationError("task ID is required", swe.ErrTypeInvalidRequest, errors.New("empty task ID"))
	}

	rlaClient := mt.RlaAtomicClient.GetClient()
	rla := rlaClient.Rla()

	req := &rlav1.GetTasksByIDsRequest{
		TaskIds: []*rlav1.UUID{{Id: taskID}},
	}
	resp, err := rla.GetTasksByIDs(ctx, req)
	if err != nil {
		logger.Warn().Err(err).Msg("Failed to get task by ID using RLA API")
		return nil, swe.WrapErr(err)
	}

	var ids []string
	for _, task := range resp.GetTasks() {
		if task == nil {
			continue
		}
		for _, u := range task.GetComponentUuids() {
			if u != nil && u.GetId() != "" {
				ids = append(ids, u.GetId())
			}
		}
	}

	logger.Info().Int("component_count", len(ids)).Msg("Completed activity")
	return ids, nil
}
