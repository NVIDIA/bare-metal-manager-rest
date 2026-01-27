// SPDX-FileCopyrightText: Copyright (c) 2021-2023 NVIDIA CORPORATION & AFFILIATES. All rights reserved.
// SPDX-License-Identifier: LicenseRef-NvidiaProprietary
//
// NVIDIA CORPORATION, its affiliates and licensors retain all intellectual
// property and proprietary rights in and to this material, related
// documentation and any modifications thereto. Any use, reproduction,
// disclosure or distribution of this material and related documentation
// without an express license agreement from NVIDIA CORPORATION or
// its affiliates is strictly prohibited.

package rla

import (
	"fmt"

	computils "github.com/nvidia/carbide-rest/site-agent/pkg/components/utils"
	"github.com/prometheus/client_golang/prometheus"
)

const (
	// MetricRLAStatus - Metric RLA Status
	MetricRLAStatus = "rla_health_status"
)

// Init - initialize RLA manager
func (rla *API) Init() {
	ManagerAccess.Data.EB.Log.Info().Msg("RLA: Initializing RLA manager")

	prometheus.MustRegister(
		prometheus.NewGaugeFunc(prometheus.GaugeOpts{
			Namespace: "elektra_site_agent",
			Name:      MetricRLAStatus,
			Help:      "RLA gRPC health status",
		},
			func() float64 {
				return float64(ManagerAccess.Data.EB.Managers.RLA.State.HealthStatus.Load())
			}))
	ManagerAccess.Data.EB.Managers.RLA.State.HealthStatus.Store(uint64(computils.CompNotKnown))

	// initialize workflow metrics
	ManagerAccess.Data.EB.Managers.RLA.State.WflowMetrics = newWorkflowMetrics()
}

// Start - Start RLA manager
func (rla *API) Start() {
	ManagerAccess.Data.EB.Log.Info().Msg("RLA: Starting RLA manager")

	// Create the client here
	// Each workflow will check and reinitialize the client if needed
	if err := rla.CreateGrpcClient(); err != nil {
		ManagerAccess.Data.EB.Log.Error().Msgf("RLA: failed to create GRPC client: %v", err)
	}
}

// GetState Machine
func (rla *API) GetState() []string {
	state := ManagerAccess.Data.EB.Managers.RLA.State
	var strs []string
	strs = append(strs, fmt.Sprintln(" GRPC Succeeded:", state.GrpcSucc.Load()))
	strs = append(strs, fmt.Sprintln(" GRPC Failed:", state.GrpcFail.Load()))
	strs = append(strs, fmt.Sprintln(" RLA Status:", computils.CompStatus(state.HealthStatus.Load())))
	strs = append(strs, fmt.Sprintln(" RLA Last Error:", state.Err))

	return strs
}

// GetGRPCClientVersion returns the current version of the GRPC client
func (rla *API) GetGRPCClientVersion() int64 {
	return ManagerAccess.Data.EB.Managers.RLA.Client.Version()
}
