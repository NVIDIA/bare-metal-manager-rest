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
	"errors"
	"os"

	"github.com/gogo/status"
	"github.com/rs/zerolog/log"
	wflows "github.com/nvidia/carbide-rest-api/carbide-rest-api-schema/schema/site-agent/workflows/v1"
	"go.opentelemetry.io/otel"
	"google.golang.org/grpc/codes"
)

var (
	ErrInvalidInstanceName = errors.New("gRPC-lib: Instance - invalid name")
	ErrInvalidInstanceID   = errors.New("gRPC-lib: Instance - invalid instance id")
	ErrInvalidMachineID    = errors.New("gRPC-lib: Instance - invalid machine id")
	ErrInvalidSegmentID    = errors.New("gRPC-lib: Instance - invalid segment id")
	ErrInvalidIPxe         = errors.New("gRPC-lib: Instance - invalid custom ipxe")
)

type InstanceInterface interface {
	// Instance Interfaces
	CreateInstance(ctx context.Context, request *wflows.CreateInstanceRequest) (response *wflows.Instance, err error)
	DeleteInstance(ctx context.Context, request *wflows.DeleteInstanceRequest) (response *wflows.InstanceReleaseResult, err error)
	RebootInstance(ctx context.Context, request *wflows.RebootInstanceRequest) (response *wflows.InstancePowerResult, err error)

	FindInstanceIDs(ctx context.Context, request *wflows.InstanceSearchFilter) (response *wflows.InstanceIdList, err error)
	FindInstancesByIDs(ctx context.Context, request *wflows.InstancesByIdsRequest) (response *wflows.InstanceList, err error)

	// DEPRECATED: use GetAllInstances instead
	GetInstance(ctx context.Context, request *wflows.InstanceSearchQuery) (response *wflows.InstanceList, err error)
	GetAllInstances(ctx context.Context, request *wflows.InstanceSearchFilter, pageSize int) (response *wflows.InstanceList, err error)
}

// DEPRECATED: use GetAllInstances instead
func (instance *compute) GetInstance(ctx context.Context, request *wflows.InstanceSearchQuery) (response *wflows.InstanceList, err error) {
	log.Info().Interface("request", request).Msg("GetInstance: received request")
	ctx, span := otel.Tracer(os.Getenv("LS_SERVICE_NAME")).Start(ctx, "CarbideClient-GetInstance")
	defer span.End()

	response, err = instance.carbide.FindInstances(ctx, request)
	if err != nil {
		log.Error().Err(err).Msg("GetInstance: error")
		return nil, err
	}
	log.Info().Int("InstanceListLen", len(response.Instances)).Msg("GetInstance: received result")
	return response, err
}

func (instance *compute) GetAllInstances(ctx context.Context, request *wflows.InstanceSearchFilter, pageSize int) (response *wflows.InstanceList, err error) {
	log.Info().Interface("request", request).Msg("GetAllInstances: received request")
	ctx, span := otel.Tracer(os.Getenv("LS_SERVICE_NAME")).Start(ctx, "CarbideClient-GetAllInstances")
	defer span.End()

	if request == nil {
		request = &wflows.InstanceSearchFilter{}
	}

	idList, err := instance.carbide.FindInstanceIds(ctx, request)
	if err != nil {
		log.Error().Err(err).Msg("FindInstanceIds: error")
		return nil, err
	}
	response = &wflows.InstanceList{}
	idChunks := SliceToChunks(idList.InstanceIds, pageSize)
	for i, chunk := range idChunks {
		list, err := instance.carbide.FindInstancesByIds(ctx, &wflows.InstancesByIdsRequest{InstanceIds: chunk})
		if err != nil {
			log.Error().Err(err).Msgf("FindInstancesByIds: error on chunk index %d", i)
			return nil, err
		}
		response.Instances = append(response.Instances, list.Instances...)
	}
	log.Info().Int("InstanceListLen", len(idList.InstanceIds)).Msg("GetInstances: received result")
	return response, err
}

func (instance *compute) FindInstanceIDs(ctx context.Context, request *wflows.InstanceSearchFilter) (response *wflows.InstanceIdList, err error) {
	log.Info().Interface("request", request).Msg("FindInstanceIDs: received request")
	ctx, span := otel.Tracer(os.Getenv("LS_SERVICE_NAME")).Start(ctx, "CarbideClient-FindInstanceIDs")
	defer span.End()

	if request == nil {
		request = &wflows.InstanceSearchFilter{}
	}
	response, err = instance.carbide.FindInstanceIds(ctx, request)
	if err != nil {
		log.Error().Err(err).Msg("FindInstanceIds: error")
		return nil, err
	}
	return
}

func (instance *compute) FindInstancesByIDs(ctx context.Context, request *wflows.InstancesByIdsRequest) (response *wflows.InstanceList, err error) {
	log.Info().Interface("request", request).Msg("FindInstancesByIDs: received request")
	ctx, span := otel.Tracer(os.Getenv("LS_SERVICE_NAME")).Start(ctx, "CarbideClient-FindInstancesByIDs")
	defer span.End()

	if request == nil {
		request = &wflows.InstancesByIdsRequest{}
	}
	response, err = instance.carbide.FindInstancesByIds(ctx, request)
	if err != nil {
		log.Error().Err(err).Msgf("FindInstancesByIds: error")
		return nil, err
	}

	return
}

func (instance *compute) CreateInstance(ctx context.Context, request *wflows.CreateInstanceRequest) (response *wflows.Instance, err error) {
	log.Info().Interface("request", request).Msg("CreateInstance: received request")
	ctx, span := otel.Tracer(os.Getenv("LS_SERVICE_NAME")).Start(ctx, "CarbideClient-CreateInstance")
	defer span.End()

	// Validations
	if request.MachineId == nil {
		// Id is mandatory
		log.Err(ErrInvalidMachineID).Msg("CreateInstance: invalid request")
		return response, ErrInvalidMachineID
	}

	// Carbide request
	// Convert Resource Request to the type needed by Site controller
	machineID := wflows.MachineId{}
	machineID.Id = request.MachineId.Id
	carbideRequest := &wflows.InstanceAllocationRequest{
		MachineId: &wflows.MachineId{},
	}
	if request.InstanceId != nil {
		carbideRequest.InstanceId = &wflows.InstanceId{Value: request.InstanceId.Value}
	}
	carbideRequest.MachineId = &machineID
	carbideRequest.Config = &wflows.InstanceConfig{}
	carbideRequest.Config.Tenant = &wflows.TenantConfig{
		TenantOrganizationId: request.TenantOrg,
		UserData:             request.UserData,
		TenantKeysetIds:      request.TenantKeysetIds,
		PhoneHomeEnabled:     request.PhoneHomeEnabled,
	}

	if request.CustomIpxe != nil {
		carbideRequest.Config.Tenant.CustomIpxe = *request.CustomIpxe
	}

	if request.AlwaysBootWithCustomIpxe != nil {
		carbideRequest.Config.Tenant.AlwaysBootWithCustomIpxe = *request.AlwaysBootWithCustomIpxe
	}

	carbideRequest.Config.Network = &wflows.InstanceNetworkConfig{}
	carbideRequest.Config.Network.Interfaces = request.Interfaces

	// InfiniBand Interfaces
	if request.IbInterfaces != nil {
		carbideRequest.Config.Infiniband = &wflows.InstanceInfinibandConfig{}
		carbideRequest.Config.Infiniband.IbInterfaces = request.IbInterfaces
	}

	// Instance labels metadata
	if request.Metadata != nil {
		carbideRequest.Metadata = request.Metadata
	}

	// Lets verify the applicable parameters
	response, err = instance.carbide.AllocateInstance(ctx, carbideRequest)
	log.Info().Interface("request", carbideRequest).Msg("CreateInstance: sent gRPC request")
	return response, err
}

func (instance *compute) DeleteInstance(ctx context.Context, request *wflows.DeleteInstanceRequest) (response *wflows.InstanceReleaseResult, err error) {
	log.Info().Interface("request", request).Msg("DeleteInstance: received request")
	ctx, span := otel.Tracer(os.Getenv("LS_SERVICE_NAME")).Start(ctx, "CarbideClient-DeleteInstance")
	defer span.End()

	// Validations
	if request.InstanceId == nil {
		// Name is mandatory
		log.Err(ErrInvalidInstanceID).Msg("DeleteInstance: invalid request")
		return response, ErrInvalidInstanceID
	}
	carbideRequest := &wflows.InstanceReleaseRequest{}
	if request.InstanceId != nil {
		carbideRequest.Id = &wflows.InstanceId{Value: request.InstanceId.Value}
	}
	response, err = instance.carbide.ReleaseInstance(ctx, carbideRequest)
	if err != nil {
		// If site controller don't have Instance, no need to fail the request
		// Check for grpc error code 'NotFound'
		if status.Code(err) == codes.NotFound {
			err = nil
		}
	}
	log.Info().Interface("request", carbideRequest).Msg("DeleteInstance: sent gRPC request")
	return response, err
}
