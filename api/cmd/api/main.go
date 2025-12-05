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

package main

import (
	"context"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	tClient "go.temporal.io/sdk/client"

	cdb "github.com/nvidia/carbide-rest/db/pkg/db"

	"github.com/nvidia/carbide-rest/api/internal/config"
	capis "github.com/nvidia/carbide-rest/api/internal/server"

	sc "github.com/nvidia/carbide-rest/api/pkg/client/site"

	// Imports for API doc generation
	_ "github.com/nvidia/carbide-rest/api/pkg/api/model"
)

const (
	// ZerologMessageFieldName specifies the field name for log message
	ZerologMessageFieldName = "msg"
	// ZerologLevelFieldName specifies the field name for log level
	ZerologLevelFieldName = "type"
)

// @title NVIDIA Forge Cloud API
// @version 1.0
// @description Forge Cloud API allows you to manage datacenter resources from Cloud
// @termsOfService https://ngc.nvidia.com/legal/terms

// @contact.name NVIDIA Forge Cloud
// @contact.email forge@nvidia.com

// @license.name Proprietary

// @BasePath /
// @schemes http https

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	// Initialize logger
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.LevelFieldName = ZerologLevelFieldName
	zerolog.MessageFieldName = ZerologMessageFieldName

	cfg := config.NewConfig()
	defer cfg.Close()

	dbConfig := cfg.GetDBConfig()

	// Initialize DB connection
	dbSession, err := cdb.NewSession(dbConfig.Host, dbConfig.Port, dbConfig.Name, dbConfig.User, dbConfig.Password, "")
	if err != nil {
		log.Panic().Err(err).Msg("failed to initialize DB session")
	} else {
		defer dbSession.Close()
	}

	// Initialize Temporal client and namespace client
	// Client objects are expensive so they are only initialized once
	tcfg, err := cfg.GetTemporalConfig()

	if err != nil {
		log.Panic().Err(err).Msg("failed to get Temporal config")
	}

	tc, tnc, err := capis.InitTemporalClients(tcfg, cfg.GetTracingEnabled())

	if err != nil {
		log.Panic().Err(err).Msg("failed to create Temporal clients")
	} else {
		defer tc.Close()
		defer tnc.Close()
	}

	_, err = tc.CheckHealth(context.Background(), &tClient.CheckHealthRequest{})
	if err != nil {
		log.Panic().Err(err).Msg("failed to check Temporal health")
	}

	scp := sc.NewClientPool(tcfg)

	// Initialize API Echo instance
	e := capis.InitAPIServer(cfg, dbSession, tc, tnc, scp)

	mconfig := cfg.GetMetricsConfig()
	if mconfig.Enabled {
		// Initialize Prometheus Echo instance
		ep := capis.InitMetricsServer(e)

		// Start Prometheus server
		log.Info().Msg("starting Metrics server")
		go func() {
			ep.Logger.Fatal(ep.Start(mconfig.GetListenAddr()))
		}()
	}

	// Start main server
	log.Info().Msg("starting API server")
	e.Logger.Fatal(e.Start(":8388"))
}
