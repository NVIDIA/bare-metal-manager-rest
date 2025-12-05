// SPDX-FileCopyrightText: Copyright (c) 2021-2023 NVIDIA CORPORATION & AFFILIATES. All rights reserved.
// SPDX-License-Identifier: LicenseRef-NvidiaProprietary
//
// NVIDIA CORPORATION, its affiliates and licensors retain all intellectual
// property and proprietary rights in and to this material, related
// documentation and any modifications thereto. Any use, reproduction,
// disclosure or distribution of this material and related documentation
// without an express license agreement from NVIDIA CORPORATION or
// its affiliates is strictly prohibited.

package carbide

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
	carbidetypes "github.com/nvidia/carbide-rest/site-agent/pkg/datatypes/managertypes/carbide"
	"github.com/nvidia/carbide-rest/site-workflow/pkg/grpc/client"
)

const (
	metricsNamespace         = "elektra_site_agent"
	metricCarbideGrpcLatency = "carbide_grpc_client_latency_seconds"
	metricWorkflowLatency    = "workflow_latency_seconds"
)

type grpcClientMetrics struct {
	responseLatency *prometheus.HistogramVec
}

func makeGrpcClientMetrics() client.Metrics {
	metrics := &grpcClientMetrics{
		responseLatency: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Namespace: metricsNamespace,
				Name:      metricCarbideGrpcLatency,
				Help:      "Response latency of each RPC",
				Buckets:   []float64{0.0005, 0.001, 0.005, 0.010, 0.025, 0.050, 0.100, 0.250, 0.500, 1.0, 2.5, 5.0, 10.0},
			},
			[]string{"grpc_method", "grpc_status_code"}),
	}
	prometheus.MustRegister(metrics.responseLatency)
	return metrics
}

func (m *grpcClientMetrics) RecordRpcResponse(method, code string, duration time.Duration) {
	ManagerAccess.Data.EB.Log.Debug().Msgf("method=%s, code=%s, duration=%v", method, code, duration)
	m.responseLatency.WithLabelValues(method, code).Observe(duration.Seconds())
}

type wflowMetrics struct {
	latency *prometheus.HistogramVec
}

func newWorkflowMetrics() carbidetypes.WorkflowMetrics {
	metrics := &wflowMetrics{
		latency: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Namespace: metricsNamespace,
				Name:      metricWorkflowLatency,
				Help:      "Latency of each workflow",
				Buckets:   []float64{0.0005, 0.001, 0.005, 0.010, 0.025, 0.050, 0.100, 0.250, 0.500, 1.0, 2.5, 5.0, 10.0},
			},
			[]string{"activity", "status"}),
	}
	prometheus.MustRegister(metrics.latency)
	return metrics
}

func (m *wflowMetrics) RecordLatency(activity string, status carbidetypes.WorkflowStatus, duration time.Duration) {
	ManagerAccess.Data.EB.Log.Debug().Msgf("activity=%s, status=%s, duration=%v", activity, status, duration)
	m.latency.WithLabelValues(activity, string(status)).Observe(duration.Seconds())
}
