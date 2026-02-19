/*
 * SPDX-FileCopyrightText: Copyright (c) 2026 NVIDIA CORPORATION & AFFILIATES. All rights reserved.
 * SPDX-License-Identifier: Apache-2.0
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package activity

import (
	"context"
	"errors"

	"github.com/rs/zerolog/log"

	swe "github.com/nvidia/bare-metal-manager-rest/site-workflow/pkg/error"
	cClient "github.com/nvidia/bare-metal-manager-rest/site-workflow/pkg/grpc/client"
	rlav1 "github.com/nvidia/bare-metal-manager-rest/workflow-schema/rla/protobuf/v1"
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
		err = errors.New("received get rack request without rack ID")
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

	logger.Info().Int32("Total", response.GetTotal()).Msg("Completed activity")

	return response, nil
}

// ValidateRackComponents validates rack components by comparing expected vs actual state via RLA.
// Supports validating a single rack, multiple racks with filters, or all racks in a site.
func (mr *ManageRack) ValidateRackComponents(ctx context.Context, request *rlav1.ValidateComponentsRequest) (*rlav1.ValidateComponentsResponse, error) {
	logger := log.With().Str("Activity", "ValidateRackComponents").Logger()
	logger.Info().Msg("Starting activity")

	var err error

	// Validate request
	switch {
	case request == nil:
		err = errors.New("received empty validate rack components request")
	}

	if err != nil {
		return nil, temporal.NewNonRetryableApplicationError(err.Error(), swe.ErrTypeInvalidRequest, err)
	}

	// Call RLA gRPC endpoint
	rlaClient := mr.RlaAtomicClient.GetClient()
	rla := rlaClient.Rla()

	response, err := rla.ValidateComponents(ctx, request)
	if err != nil {
		logger.Warn().Err(err).Msg("Failed to validate rack components using RLA API")
		return nil, swe.WrapErr(err)
	}

	logger.Info().Int32("TotalDiffs", response.GetTotalDiffs()).Msg("Completed activity")

	return response, nil
}
