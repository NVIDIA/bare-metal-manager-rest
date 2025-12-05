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

	"github.com/rs/zerolog/log"
	"go.opentelemetry.io/otel"

	wflows "github.com/nvidia/carbide-rest/workflow-schema/schema/site-agent/workflows/v1"
)

var (
	ErrInvalidTenantRequest        = errors.New("gRPC-lib: Tenant - invalid request")
	ErrInvalidTenantOrganizationID = errors.New("gRPC-lib: Tenant - invalid Organization ID")
	ErrInvalidTenantName           = errors.New("gRPC-lib: Tenant - invalid name")
)

type TenantInterface interface {
	// Tenant Interfaces
	CreateTenant(ctx context.Context, request *wflows.CreateTenantRequest) (response *wflows.CreateTenantResponse, err error)
	UpdateTenant(ctx context.Context, request *wflows.UpdateTenantRequest) (response *wflows.UpdateTenantResponse, err error)

	FindTenantOrganizationIDs(ctx context.Context, request *wflows.TenantSearchFilter) (response *wflows.TenantOrganizationIdList, err error)
	FindTenantsByOrganizationIDs(ctx context.Context, request *wflows.TenantByOrganizationIdsRequest) (response *wflows.TenantList, err error)
}

func (tenant *compute) CreateTenant(ctx context.Context, request *wflows.CreateTenantRequest) (response *wflows.CreateTenantResponse, err error) {
	log.Info().Interface("request", request).Msg("CreateTenant: received request")
	ctx, span := otel.Tracer(os.Getenv("LS_SERVICE_NAME")).Start(ctx, "CarbideClient-CreateTenant")
	defer span.End()

	if request == nil {
		log.Err(ErrInvalidMachineID).Msg("CreateTenant: invalid request")
		return response, ErrInvalidTenantRequest
	}

	if request.OrganizationId == "" {
		log.Err(ErrInvalidMachineID).Msg("CreateTenant: invalid Organization ID")
		return response, ErrInvalidTenantOrganizationID
	}

	if request.Metadata != nil && request.Metadata.Name == "" {
		log.Err(ErrInvalidMachineID).Msg("CreateTenant: invalid Name")
		return response, ErrInvalidTenantName
	}

	response, err = tenant.carbide.CreateTenant(ctx, request)
	if err != nil {
		log.Error().Err(err).Msg("CreateTenant: error")
		return nil, err
	}
	return
}

func (tenant *compute) UpdateTenant(ctx context.Context, request *wflows.UpdateTenantRequest) (response *wflows.UpdateTenantResponse, err error) {
	log.Info().Interface("request", request).Msg("UpdateTenant: received request")
	ctx, span := otel.Tracer(os.Getenv("LS_SERVICE_NAME")).Start(ctx, "CarbideClient-UpdateTenant")
	defer span.End()

	if request == nil {
		log.Err(ErrInvalidMachineID).Msg("CreateTenant: invalid request")
		return response, ErrInvalidTenantRequest
	}

	if request.OrganizationId == "" {
		log.Err(ErrInvalidMachineID).Msg("CreateTenant: invalid Organization ID")
		return response, ErrInvalidTenantOrganizationID
	}

	if request.Metadata != nil && request.Metadata.Name == "" {
		log.Err(ErrInvalidMachineID).Msg("CreateTenant: invalid Name")
		return response, ErrInvalidTenantName
	}

	response, err = tenant.carbide.UpdateTenant(ctx, request)
	if err != nil {
		log.Error().Err(err).Msg("UpdateTenant: error")
		return nil, err
	}
	return
}

func (tenant *compute) FindTenantOrganizationIDs(ctx context.Context, request *wflows.TenantSearchFilter) (response *wflows.TenantOrganizationIdList, err error) {
	log.Info().Interface("request", request).Msg("FindTenantOrganizationIDs: received request")
	ctx, span := otel.Tracer(os.Getenv("LS_SERVICE_NAME")).Start(ctx, "CarbideClient-FindTenantOrganizationIDs")
	defer span.End()

	if request == nil {
		request = &wflows.TenantSearchFilter{}
	}
	response, err = tenant.carbide.FindTenantOrganizationIds(ctx, request)
	if err != nil {
		log.Error().Err(err).Msg("FindTenantOrganizationIds: error")
		return nil, err
	}
	return
}

func (tenant *compute) FindTenantsByOrganizationIDs(ctx context.Context, request *wflows.TenantByOrganizationIdsRequest) (response *wflows.TenantList, err error) {
	log.Info().Interface("request", request).Msg("FindTenantsByOrganizationIDs: received request")
	ctx, span := otel.Tracer(os.Getenv("LS_SERVICE_NAME")).Start(ctx, "CarbideClient-FindTenantsByOrganizationIDs")
	defer span.End()

	if request == nil {
		request = &wflows.TenantByOrganizationIdsRequest{}
	}
	response, err = tenant.carbide.FindTenantsByOrganizationIds(ctx, request)
	if err != nil {
		log.Error().Err(err).Msgf("FindTenantsByOrganizationIds: error")
		return nil, err
	}

	return
}
