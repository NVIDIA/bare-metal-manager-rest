// SPDX-FileCopyrightText: Copyright (c) 2021-2023 NVIDIA CORPORATION & AFFILIATES. All rights reserved.
// SPDX-License-Identifier: LicenseRef-NvidiaProprietary
//
// NVIDIA CORPORATION, its affiliates and licensors retain all intellectual
// property and proprietary rights in and to this material, related
// documentation and any modifications thereto. Any use, reproduction,
// disclosure or distribution of this material and related documentation
// without an express license agreement from NVIDIA CORPORATION or
// its affiliates is strictly prohibited.

package workflow

import (
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/rs/zerolog/log"
	computils "github.com/nvidia/carbide-rest/site-agent/pkg/components/utils"
	"gopkg.in/fsnotify.v1"
)

const (
	// MetricTemporalConnAttempted - Metric Temporal Conn Attempted
	MetricTemporalConnAttempted = "temporal_connection_attempted"
	// MetricTemporalConnSucc - Metric Temporal Conn Succ
	MetricTemporalConnSucc = "temporal_connection_succeeded"
	// MetricTemporalConnStatus - Metric Temporal Conn Status
	MetricTemporalConnStatus = "temporal_connection_status"
)

// Init - initialize the workflow orchestrator
func (wflow *API) Init() {
	ManagerAccess.Data.EB.Log.Info().Msg("Workflow: Initializing the workflow")

	prometheus.MustRegister(
		prometheus.NewCounterFunc(prometheus.CounterOpts{
			Namespace: "elektra_site_agent",
			Name:      MetricTemporalConnStatus,
			Help:      "temporal health status of the elektra_site_agent",
		},
			func() float64 {
				return float64(ManagerAccess.Data.EB.Managers.Workflow.State.HealthStatus.Load())
			}))
	ManagerAccess.Data.EB.Managers.Workflow.State.HealthStatus.Store(uint64(computils.CompUnhealthy))

	prometheus.MustRegister(
		prometheus.NewCounterFunc(prometheus.CounterOpts{
			Namespace: "elektra_site_agent",
			Name:      MetricTemporalConnAttempted,
			Help:      "temporal connection attempted of elektra_site_agent",
		},
			func() float64 {
				return float64(ManagerAccess.Data.EB.Managers.Workflow.State.ConnectionAttempted.Load())
			}))

	prometheus.MustRegister(
		prometheus.NewCounterFunc(prometheus.CounterOpts{
			Namespace: "elektra_site_agent",
			Name:      MetricTemporalConnSucc,
			Help:      "temporal connection succeded of elektra_site_agent",
		},
			func() float64 {
				return float64(ManagerAccess.Data.EB.Managers.Workflow.State.ConnectionSucc.Load())
			}))
}

// GetState - handle http request
func (wflow *API) GetState() []string {
	wc := ManagerAccess.Conf.EB.Temporal
	wt := ManagerAccess.Data.EB.Managers.Workflow
	var strs []string
	strs = append(strs, fmt.Sprintln("Temporal Host: ", wc.Host, wc.Port))
	strs = append(strs, fmt.Sprintln("Temporal Connection Attempted: ", wt.State.ConnectionAttempted.Load()))
	strs = append(strs, fmt.Sprintln("Temporal Connection Succeeded: ", wt.State.ConnectionSucc.Load()))
	strs = append(strs, fmt.Sprintln("Temporal Status: ", computils.CompStatus(wt.State.HealthStatus.Load()).String()))
	strs = append(strs, fmt.Sprintln("Temporal Last Error: ", *wt.State.Err))
	strs = append(strs, fmt.Sprintln("Temporal Connection Time: ", wt.State.ConnectionTime))

	return strs
}

// Start the workflow orchestrator
func (wflow *API) Start() {
	ManagerAccess.Data.EB.Log.Info().Msg("Workflow: Starting the workflow")
	Orchestrator()
	wflow.WatchCertFile()
}

// WatchCertFile - Watch Cert File
func (wflow *API) WatchCertFile() {
	fileName, caCertPath := ManagerAccess.Conf.EB.Temporal.GetTemporalCACertFilePath()
	kpFileName, clientCertPath := ManagerAccess.Conf.EB.Temporal.GetTemporalClientCertFilePath()
	path := []string{caCertPath, clientCertPath}
	file := make(map[string]bool)
	file[caCertPath+fileName] = true
	for _, v := range kpFileName {
		file[clientCertPath+v] = true
	}
	go wflow.watchFiles(path, file)
}

// watchFiles watch on the secret files
func (wflow *API) watchFiles(path []string, file map[string]bool) {
	log.Info().Msgf("Workflow: Watching secret %v", path)

	// Create a new watcher.
	w, err := fsnotify.NewWatcher()
	if err != nil {
		log.Panic().Msgf("Workflow: creating a new watcher: %v", err)
	}
	defer w.Close()

	// Watch the directory, not the file itself.
	for _, v := range path {
		err = w.Add(v)
		if err != nil {
			log.Panic().Msgf("Workflow: add file to watcher: %v", err)
		}
	}

	for {
		select {
		// Read from Errors.
		case err, ok := <-w.Errors:
			if !ok { // Channel was closed (i.e. Watcher.Close() was called).
				return
			}
			log.Panic().Msgf("Workflow: watcher closed: %v", err.Error())
		// Read from Events.
		case e, ok := <-w.Events:
			if !ok { // Channel was closed (i.e. Watcher.Close() was called).
				return
			}

			// log.Info().Msgf("Got event %v ", e.String())
			_, ok = file[e.Name]
			if !ok {
				continue
			}
			Orchestrator()
			log.Info().Msgf("Workflow: Back to watching secret %v ", e.String())
		}
	}
}
