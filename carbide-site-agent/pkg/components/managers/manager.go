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

package managers

import (
	"context"
	"net/http"

	"github.com/NVIDIA/carbide-rest-api/carbide-site-agent/pkg/components/managers/machinevalidation"

	"github.com/NVIDIA/carbide-rest-api/carbide-site-agent/pkg/metadata"

	// "github.com/lightstep/otel-launcher-go/launcher" // Removed due to compatibility issues
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog/log"

	"github.com/NVIDIA/carbide-rest-api/carbide-site-agent/pkg/components/managers/bootstrap"
	"github.com/NVIDIA/carbide-rest-api/carbide-site-agent/pkg/components/managers/carbide"
	"github.com/NVIDIA/carbide-rest-api/carbide-site-agent/pkg/components/managers/expectedmachine"
	"github.com/NVIDIA/carbide-rest-api/carbide-site-agent/pkg/components/managers/health"
	"github.com/NVIDIA/carbide-rest-api/carbide-site-agent/pkg/components/managers/infinibandpartition"
	"github.com/NVIDIA/carbide-rest-api/carbide-site-agent/pkg/components/managers/instance"
	"github.com/NVIDIA/carbide-rest-api/carbide-site-agent/pkg/components/managers/instancetype"
	"github.com/NVIDIA/carbide-rest-api/carbide-site-agent/pkg/components/managers/machine"
	"github.com/NVIDIA/carbide-rest-api/carbide-site-agent/pkg/components/managers/managerapi"
	"github.com/NVIDIA/carbide-rest-api/carbide-site-agent/pkg/components/managers/networksecuritygroup"
	"github.com/NVIDIA/carbide-rest-api/carbide-site-agent/pkg/components/managers/operatingsystem"
	"github.com/NVIDIA/carbide-rest-api/carbide-site-agent/pkg/components/managers/sku"
	"github.com/NVIDIA/carbide-rest-api/carbide-site-agent/pkg/components/managers/sshkeygroup"
	"github.com/NVIDIA/carbide-rest-api/carbide-site-agent/pkg/components/managers/subnet"
	"github.com/NVIDIA/carbide-rest-api/carbide-site-agent/pkg/components/managers/tenant"
	"github.com/NVIDIA/carbide-rest-api/carbide-site-agent/pkg/components/managers/vpc"
	"github.com/NVIDIA/carbide-rest-api/carbide-site-agent/pkg/components/managers/vpcprefix"
	"github.com/NVIDIA/carbide-rest-api/carbide-site-agent/pkg/components/managers/workflow"

	"github.com/NVIDIA/carbide-rest-api/carbide-site-agent/pkg/datatypes/elektratypes"

	computils "github.com/NVIDIA/carbide-rest-api/carbide-site-agent/pkg/components/utils"
)

// NewAPIHandlers - handle new api
func NewAPIHandlers() {
	managerapi.ManagerHdl = managerapi.ManagerAPI{
		// Add all the Managers here
		Orchestrator:         &workflow.API{},
		VPC:                  &vpc.API{},
		VpcPrefix:            &vpcprefix.API{},
		Subnet:               &subnet.API{},
		Instance:             &instance.API{},
		Machine:              &machine.API{},
		Carbide:              &carbide.API{},
		Bootstrap:            &bootstrap.BoostrapAPI{},
		Health:               &health.API{},
		SSHKeyGroup:          &sshkeygroup.API{},
		InfiniBandPartition:  &infinibandpartition.API{},
		Tenant:               &tenant.API{},
		OperatingSystem:      &operatingsystem.API{},
		MachineValidation:    &machinevalidation.API{},
		InstanceType:         &instancetype.API{},
		NetworkSecurityGroup: &networksecuritygroup.API{},
		ExpectedMachine:      &expectedmachine.API{},
		SKU:                  &sku.API{},
	}
}

// NewInstance - new instance with the parent datastruct
func NewInstance(superforge *elektratypes.Elektra) (*Manager, error) {
	NewAPIHandlers()
	ManagerAccess = &Manager{
		Data: &managerapi.ManagerData{
			EB: superforge,
		},
		API: &managerapi.ManagerHdl,
		Conf: &managerapi.ManagerConf{
			EB: superforge.Conf,
		},
	}
	ManagerAccess.NewInstance()
	return ManagerAccess, nil
}

// NewInstance - instantiates all the managers
func (Managers *Manager) NewInstance() {
	// Instantiate all the managers here
	Managers.Orchestrator()
	Managers.VPC()
	Managers.VpcPrefix()
	Managers.Subnet()
	Managers.Instance()
	Managers.Carbide()
	Managers.Machine()
	Managers.Bootstrap()
	Managers.Health()
	Managers.SSHKeyGroup()
	Managers.InfiniBandPartition()
	Managers.Tenant()
	Managers.OperatingSystem()
	Managers.MachineValidation()
	Managers.InstanceType()
	Managers.NetworkSecurityGroup()
	Managers.ExpectedMachine()
	Managers.SKU()
}

// Init - initialize all the mgrs
func (Managers *Manager) Init() {
	ManagerAccess.Data.EB.Log.Info().Msg("Managers: Initializing all the managers")
	// register version metric (build_version, build_date)
	versionGauge := prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "elektra_site_agent",
		Name:      "version",
		Help:      "version of the elektra_site_agent",
	}, []string{"build_version", "build_date"})
	prometheus.MustRegister(versionGauge)
	// set the value once, since it does not change
	versionGauge.WithLabelValues(metadata.Version, metadata.BuildTime).Set(1)
	// register health status metric
	prometheus.MustRegister(
		prometheus.NewCounterFunc(prometheus.CounterOpts{
			Namespace: "elektra_site_agent",
			Name:      "health_status",
			Help:      "health status of the elektra_site_agent",
		},
			func() float64 {
				return float64(ManagerAccess.Data.EB.HealthStatus.Load())
			}))
	ManagerAccess.Data.EB.HealthStatus.Store(uint64(computils.CompUnhealthy))

	Managers.Orchestrator().Init()
	Managers.Carbide().Init()
	Managers.Bootstrap().Init()
	Managers.VPC().Init()
	Managers.VpcPrefix().Init()
	Managers.Subnet().Init()
	Managers.Instance().Init()
	Managers.Health().Init()
	Managers.SSHKeyGroup().Init()
	Managers.InfiniBandPartition().Init()
	Managers.Tenant().Init()
	Managers.OperatingSystem().Init()
	Managers.MachineValidation().Init()
	Managers.InstanceType().Init()
	Managers.NetworkSecurityGroup().Init()
	Managers.ExpectedMachine().Init()
	Managers.SKU().Init()
}

// Start - start the mgrs
func (Managers *Manager) Start() {
	// Start otel exporter
	startOTELDaemon(NewDefaultContext(context.Background()))
	go StartMetricServer()
	StartHTTPServer()
	ManagerAccess.Data.EB.Log.Info().Msg("Managers: Starting all the managers")
	Managers.Carbide().Start()
	Managers.Bootstrap().Start()
	Managers.Orchestrator().Start()
}

// StartMetricServer - Start serving Metric Server
func StartMetricServer() {
	log.Info().Msgf("Beginning to serve on port %v", ManagerAccess.Conf.EB.MetricsPort)
	http.Handle("/metrics", promhttp.Handler())
	port := ":" + ManagerAccess.Conf.EB.MetricsPort
	http.ListenAndServe(port, nil)
}

type otelLog struct{}

// Fatalf - otel Fatalf
func (l otelLog) Fatalf(format string, v ...interface{}) {
	ManagerAccess.Data.EB.Log.Fatal().Msgf(format, v...)
}

// Debugf - otel Debugf
func (l otelLog) Debugf(format string, v ...interface{}) {
	ManagerAccess.Data.EB.Log.Debug().Msgf(format, v...)
}

func startOTELDaemon(ctx context.Context) {
	log := ManagerAccess.Data.EB.Log
	log.Info().Msgf("OTEL daemon startup - lightstep launcher disabled, use standard OTel env configuration")
	// NOTE: Lightstep launcher has been disabled due to compatibility issues with newer OTel collector
	// TODO: Implement standard OpenTelemetry setup without lightstep if needed
}
