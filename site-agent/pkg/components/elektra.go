// SPDX-FileCopyrightText: Copyright (c) 2021-2023 NVIDIA CORPORATION & AFFILIATES. All rights reserved.
// SPDX-License-Identifier: LicenseRef-NvidiaProprietary
//
// NVIDIA CORPORATION, its affiliates and licensors retain all intellectual
// property and proprietary rights in and to this material, related
// documentation and any modifications thereto. Any use, reproduction,
// disclosure or distribution of this material and related documentation
// without an express license agreement from NVIDIA CORPORATION or
// its affiliates is strictly prohibited.

package elektra

import (
	zlog "github.com/rs/zerolog/log"

	"github.com/nvidia/carbide-rest/site-agent/pkg/components/config"
	"github.com/nvidia/carbide-rest/site-agent/pkg/components/managers"
	"github.com/nvidia/carbide-rest/site-agent/pkg/datatypes/elektratypes"
)

// Interface - Managers' interface
type Interface interface {
	Managers() managers.Manager
}

// Elektra - Managers struct
type Elektra struct {
	manager *managers.Manager
}

// Init - initializes the cluster
func (Cluster *Elektra) Init() (err error) {
	zlog.Info().Msg("Elektra: Initializing Elektra cluster")
	Cluster.Managers().Init()
	return nil
}

// Start () Start the Cluster
func (Cluster *Elektra) Start() (err error) {
	zlog.Info().Msg("Elektra: Starting Elektra cluster")
	Cluster.Managers().Start()
	return nil
}

// Managers () Instantiate the Managers
func (Cluster *Elektra) Managers() *managers.Manager {
	return Cluster.manager
}

// NewElektraAPI - Instantiate new struct
func NewElektraAPI(superElektra *elektratypes.Elektra, utMode bool) (*Elektra, error) {
	zlog.Info().Msg("Elektra: Initializing Config Manager")
	var eb Elektra
	var err error
	// Initialize Global Config
	// Load configuration
	if superElektra != nil {
		// Configuration
		zlog.Info().Msg("Elektra: Initializing Config Manager")
		superElektra.Conf = config.NewElektraConfig(utMode)
		eb.manager, err = managers.NewInstance(superElektra)
		zlog.Info().Interface("config", superElektra.Conf).Msg("Elektra: Config Manager initialized")
	}

	return &eb, err
}
