// SPDX-FileCopyrightText: Copyright (c) 2021-2023 NVIDIA CORPORATION & AFFILIATES. All rights reserved.
// SPDX-License-Identifier: LicenseRef-NvidiaProprietary
//
// NVIDIA CORPORATION, its affiliates and licensors retain all intellectual
// property and proprietary rights in and to this material, related
// documentation and any modifications thereto. Any use, reproduction,
// disclosure or distribution of this material and related documentation
// without an express license agreement from NVIDIA CORPORATION or
// its affiliates is strictly prohibited.

package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

const (
	MetricsNamespace = "cloud_workflow"
)

type coreMetrics struct {
	Info *prometheus.GaugeVec
}

// NewCoreMetrics creates a new coreMetrics struct and registers the metrics with the provided registerer
func NewCoreMetrics(reg prometheus.Registerer) *coreMetrics {
	m := &coreMetrics{
		Info: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: MetricsNamespace,
			Name:      "info",
			Help:      "Information about the Cloud/Site worker",
		}, []string{"version", "namespace"}),
	}

	reg.MustRegister(m.Info)

	return m
}
