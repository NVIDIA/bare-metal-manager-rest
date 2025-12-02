/*
 * SPDX-FileCopyrightText: Copyright (c) 2021-2023 NVIDIA CORPORATION & AFFILIATES. All rights reserved.
 * SPDX-License-Identifier: LicenseRef-NvidiaProprietary
 *
 * NVIDIA CORPORATION, its affiliates and licensors retain all intellectual
 * property and proprietary rights in and to this material, related
 * documentation and any modifications thereto. Any use, reproduction,
 * disclosure or distribution of this material and related documentation
 * without an express license agreement from NVIDIA CORPORATION or
 * its affiliates is strictly prohibited.
 */

package client

import (
	"context"
	"os"

	wflows "github.com/NVIDIA/carbide-rest-api/carbide-rest-api-schema/schema/site-agent/workflows/v1"
	"github.com/rs/zerolog/log"
	"go.opentelemetry.io/otel"
)

func (instance *compute) RebootInstance(ctx context.Context, resourceRequest *wflows.RebootInstanceRequest) (result *wflows.InstancePowerResult, err error) {
	log.Info().Interface("request", resourceRequest).Msg("RebootInstance: received request")
	ctx, span := otel.Tracer(os.Getenv("OTEL_SERVICE_NAME")).Start(ctx, "CarbideClient-RebootInstance")
	defer span.End()

	// Validations
	if resourceRequest.MachineId == nil {
		// Name is mandatory
		log.Err(ErrInvalidMachineID).Msg("RebootInstance: invalid request")
		return result, ErrInvalidMachineID
	}
	carbideRequest := &wflows.InstancePowerRequest{MachineId: &wflows.MachineId{}}

	// Convert Resource Request to the type needed by Site controller
	machineID := carbideRequest.MachineId
	machineID.Id = resourceRequest.MachineId.Id
	carbideRequest.MachineId = machineID
	carbideRequest.BootWithCustomIpxe = resourceRequest.BootWithCustomIpxe
	carbideRequest.Operation = wflows.InstancePowerRequest_POWER_RESET
	carbideRequest.ApplyUpdatesOnReboot = resourceRequest.ApplyUpdatesOnReboot
	grpcResponse, err := instance.carbide.InvokeInstancePower(ctx, carbideRequest)
	log.Info().Interface("request", carbideRequest).Interface("response", grpcResponse).Msg("RebootInstance: sent gRPC request")
	return grpcResponse, err
}
