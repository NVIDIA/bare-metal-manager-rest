// SPDX-FileCopyrightText: Copyright (c) 2021-2023 NVIDIA CORPORATION & AFFILIATES. All rights reserved.
// SPDX-License-Identifier: LicenseRef-NvidiaProprietary
//
// NVIDIA CORPORATION, its affiliates and licensors retain all intellectual
// property and proprietary rights in and to this material, related
// documentation and any modifications thereto. Any use, reproduction,
// disclosure or distribution of this material and related documentation
// without an express license agreement from NVIDIA CORPORATION or
// its affiliates is strictly prohibited.

package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/nvidia/carbide-rest/site-agent/pkg/metadata"

	"github.com/rs/zerolog/log"
	components "github.com/nvidia/carbide-rest/site-agent/pkg/components"
	"github.com/nvidia/carbide-rest/site-agent/pkg/datatypes/elektratypes"
)

// InitElektra initializes the Elektra site agent framework
func InitElektra() {
	// Initialize Elektra microservice
	log.Info().Msg("Elektra: Initializing Elektra service")

	// TODO: this is for verification that we can get version, will move it to a metric after
	log.Info().Msgf("Elektra: version=%s, time=%s", metadata.Version, metadata.BuildTime)

	// Initialize Elektra Data Structures
	elektraTypes := elektratypes.NewElektraTypes()

	// Initialize Elektra API
	api, initErr := components.NewElektraAPI(elektraTypes, false)
	if initErr != nil {
		log.Error().Err(initErr).Msg("Elektra: Failed to initialize Elektra API")
	} else {
		log.Info().Msg("Elektra: Successfully initialized Elektra API")
	}

	// Initialize Elektra Managers
	api.Init()

	// Start Elektra Managers
	api.Start()
}

func main() {
	InitElektra()
	// sleep
	// Wait forever
	termChan := make(chan os.Signal, 1)
	signal.Notify(termChan, syscall.SIGINT, syscall.SIGTERM)
	select {
	case <-termChan:
		return
	}
}
