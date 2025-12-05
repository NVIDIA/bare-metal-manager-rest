// SPDX-FileCopyrightText: Copyright (c) 2021-2023 NVIDIA CORPORATION & AFFILIATES. All rights reserved.
// SPDX-License-Identifier: LicenseRef-NvidiaProprietary
//
// NVIDIA CORPORATION, its affiliates and licensors retain all intellectual
// property and proprietary rights in and to this material, related
// documentation and any modifications thereto. Any use, reproduction,
// disclosure or distribution of this material and related documentation
// without an express license agreement from NVIDIA CORPORATION or
// its affiliates is strictly prohibited.

package client

import (
	"context"
	"os"

	"github.com/rs/zerolog/log"
	wflows "github.com/nvidia/carbide-rest/workflow-schema/schema/site-agent/workflows/v1"
	"go.opentelemetry.io/otel"
)

// MachineInterface grpc client interface
type MachineInterface interface {
	// Machine Interfaces
	// DEPRECATED: use GetAllMachines instead
	GetMachine(ctx context.Context, request *wflows.MachineSearchQuery) (response *wflows.MachineList, err error)
	GetAllMachines(ctx context.Context, request *wflows.MachineSearchConfig, pageSize int) (response *wflows.MachineList, err error)
	FindMachineIDs(ctx context.Context, request *wflows.MachineSearchConfig) (response *wflows.MachineIdList, err error)
	FindMachinesByIDs(ctx context.Context, request *wflows.MachinesByIdsRequest) (response *wflows.MachineList, err error)
	// CreateMachine() error
	// UpdateMachine() error
	// DeleteMachine() error
}

// DEPRECATED: use GetAllMachines instead
func (machine *compute) GetMachine(ctx context.Context, request *wflows.MachineSearchQuery) (response *wflows.MachineList, err error) {
	log.Info().Interface("request", request).Msg("GetMachine: received request")
	ctx, span := otel.Tracer(os.Getenv("LS_SERVICE_NAME")).Start(ctx, "CarbideClient-GetMachine")
	defer span.End()

	response, err = machine.carbide.FindMachines(ctx, request)
	if err != nil {
		log.Error().Err(err).Msg("GetMachine: error")
		return nil, err
	}

	log.Info().Int("MachineListLen", len(response.Machines)).Msg("GetMachine: received result")
	return response, err
}

func (machine *compute) GetAllMachines(ctx context.Context, request *wflows.MachineSearchConfig, pageSize int) (response *wflows.MachineList, err error) {
	log.Info().Interface("request", request).Msg("GetAllMachines: received request")
	ctx, span := otel.Tracer(os.Getenv("LS_SERVICE_NAME")).Start(ctx, "CarbideClient-GetAllMachines")
	defer span.End()

	if request == nil {
		request = &wflows.MachineSearchConfig{}
	}

	idList, err := machine.carbide.FindMachineIds(ctx, request)
	if err != nil {
		log.Error().Err(err).Msg("FindMachineIds: error")
		return nil, err
	}
	response = &wflows.MachineList{}
	idChunks := SliceToChunks(idList.MachineIds, pageSize)
	for i, chunk := range idChunks {
		list, err := machine.carbide.FindMachinesByIds(ctx, &wflows.MachinesByIdsRequest{MachineIds: chunk})
		if err != nil {
			log.Error().Err(err).Msgf("FindMachinesByIds: error on chunk index %d", i)
			return nil, err
		}
		response.Machines = append(response.Machines, list.Machines...)
	}
	log.Info().Int("MachineListLen", len(idList.MachineIds)).Msg("GetMachines: received result")
	return response, err
}

func (machine *compute) FindMachineIDs(ctx context.Context, request *wflows.MachineSearchConfig) (response *wflows.MachineIdList, err error) {
	log.Info().Interface("request", request).Msg("FindMachineIDs: received request")
	ctx, span := otel.Tracer(os.Getenv("LS_SERVICE_NAME")).Start(ctx, "CarbideClient-FindMachineIDs")
	defer span.End()

	if request == nil {
		request = &wflows.MachineSearchConfig{}
	}

	response, err = machine.carbide.FindMachineIds(ctx, request)
	if err != nil {
		log.Error().Err(err).Msg("FindMachineIds: error")
		return nil, err
	}
	return
}

func (machine *compute) FindMachinesByIDs(ctx context.Context, request *wflows.MachinesByIdsRequest) (response *wflows.MachineList, err error) {
	log.Info().Interface("request", request).Msg("FindMachinesByIDs: received request")
	ctx, span := otel.Tracer(os.Getenv("LS_SERVICE_NAME")).Start(ctx, "CarbideClient-FindMachinesByIDs")
	defer span.End()

	if request == nil {
		request = &wflows.MachinesByIdsRequest{}
	}

	response, err = machine.carbide.FindMachinesByIds(ctx, request)
	if err != nil {
		log.Error().Err(err).Msgf("FindMachinesByIds: error")
		return nil, err
	}
	return
}
