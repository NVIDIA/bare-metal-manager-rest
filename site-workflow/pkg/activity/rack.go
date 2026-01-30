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

	rlav1 "github.com/nvidia/carbide-rest/workflow-schema/rla/protobuf/v1"
	swe "github.com/nvidia/carbide-rest/site-workflow/pkg/error"
	cClient "github.com/nvidia/carbide-rest/site-workflow/pkg/grpc/client"
	"go.temporal.io/sdk/temporal"
)

// ManageRack is an activity wrapper for Rack management via RLA
type ManageRack struct {
	RlaAtomicClient *cClient.RlaAtomicClient
}

// NewManageRack returns a new ManageRack client
func NewManageRack(rlaClient *cClient.RlaAtomicClient) ManageRack {
	return ManageRack{
		RlaAtomicClient: rlaClient,
	}
}

// GetRack retrieves a rack by its UUID from RLA
func (mr *ManageRack) GetRack(ctx context.Context, request *rlav1.GetRackInfoByIDRequest) (*rlav1.GetRackInfoResponse, error) {
	logger := log.With().Str("Activity", "GetRack").Logger()
	logger.Info().Msg("Starting activity")

	var err error

	// Validate request
	switch {
	case request == nil:
		err = errors.New("received empty get rack request")
	case request.Id == nil || request.Id.Id == "":
		err = errors.New("received get rack request missing rack ID")
	}

	if err != nil {
		return nil, temporal.NewNonRetryableApplicationError(err.Error(), swe.ErrTypeInvalidRequest, err)
	}

	// Call RLA gRPC endpoint
	rlaClient := mr.RlaAtomicClient.GetClient()
	rla := rlaClient.Rla()

	response, err := rla.GetRackInfoByID(ctx, request)
	if err != nil {
		logger.Warn().Err(err).Msg("Failed to get rack by ID using RLA API")
		return nil, swe.WrapErr(err)
	}

	logger.Info().Msg("Completed activity")

	return response, nil
}

// GetRacks retrieves a list of racks from RLA with optional filters
func (mr *ManageRack) GetRacks(ctx context.Context, request *rlav1.GetListOfRacksRequest) (*rlav1.GetListOfRacksResponse, error) {
	logger := log.With().Str("Activity", "GetRacks").Logger()
	logger.Info().Msg("Starting activity")

	// Request can be nil or empty for getting all racks
	if request == nil {
		request = &rlav1.GetListOfRacksRequest{}
	}

	// Call RLA gRPC endpoint
	rlaClient := mr.RlaAtomicClient.GetClient()
	rla := rlaClient.Rla()

	response, err := rla.GetListOfRacks(ctx, request)
	if err != nil {
		logger.Warn().Err(err).Msg("Failed to get list of racks using RLA API")
		return nil, swe.WrapErr(err)
	}

	logger.Info().Int32("total", response.GetTotal()).Msg("Completed activity")

	return response, nil
}
