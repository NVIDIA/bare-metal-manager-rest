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

package carbide

import (
	"fmt"

	computils "github.com/NVIDIA/carbide-rest-api/carbide-site-agent/pkg/components/utils"
	"github.com/prometheus/client_golang/prometheus"
)

const (
	// MetricCarbideStatus - Metric Carbide Status
	MetricCarbideStatus = "carbide_health_status"
)

// Init - initialize carbide manager
func (carbide *API) Init() {
	ManagerAccess.Data.EB.Log.Info().Msg("Carbide: Initializing the carbide")

	statusGauge := prometheus.NewGaugeFunc(prometheus.GaugeOpts{
		Namespace: "elektra_site_agent",
		Name:      MetricCarbideStatus,
		Help:      "Carbide gRPC health status",
	},
		func() float64 {
			return float64(ManagerAccess.Data.EB.Managers.Carbide.State.HealthStatus.Load())
		})
	if err := prometheus.Register(statusGauge); err != nil {
		if _, ok := err.(prometheus.AlreadyRegisteredError); !ok {
			ManagerAccess.Data.EB.Log.Warn().Err(err).Msg("Failed to register carbide health status metric")
		}
	}
	ManagerAccess.Data.EB.Managers.Carbide.State.HealthStatus.Store(uint64(computils.CompNotKnown))

	// initialize workflow metrics
	ManagerAccess.Data.EB.Managers.Carbide.State.WflowMetrics = newWorkflowMetrics()
}

// Start - start carbide manager
func (carbide *API) Start() {
	ManagerAccess.Data.EB.Log.Info().Msg("Carbide: Starting the carbide")

	// Create the client here
	// Each workflow will check and reinitialize the client if needed
	if err := carbide.CreateGRPCClient(); err != nil {
		ManagerAccess.Data.EB.Log.Error().Msgf("Carbide: failed to create GRPC client: %v", err)
	}
}

// GetState Machine
func (carbide *API) GetState() []string {
	state := ManagerAccess.Data.EB.Managers.Carbide.State
	var strs []string
	strs = append(strs, fmt.Sprintln(" GRPC Succeeded:", state.GrpcSucc.Load()))
	strs = append(strs, fmt.Sprintln(" GRPC Failed:", state.GrpcFail.Load()))
	strs = append(strs, fmt.Sprintln(" GRPC Status:", computils.CompStatus(state.HealthStatus.Load())))
	strs = append(strs, fmt.Sprintln(" GRPC Last Error:", state.Err))

	return strs
}

// GetGRPCClientVersion returns the current version of the GRPC client
func (carbide *API) GetGRPCClientVersion() int64 {
	return ManagerAccess.Data.EB.Managers.Carbide.Client.Version()
}
