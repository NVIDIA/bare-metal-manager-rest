// SPDX-FileCopyrightText: Copyright (c) 2021-2023 NVIDIA CORPORATION & AFFILIATES. All rights reserved.
// SPDX-License-Identifier: LicenseRef-NvidiaProprietary
//
// NVIDIA CORPORATION, its affiliates and licensors retain all intellectual
// property and proprietary rights in and to this material, related
// documentation and any modifications thereto. Any use, reproduction,
// disclosure or distribution of this material and related documentation
// without an express license agreement from NVIDIA CORPORATION or
// its affiliates is strictly prohibited.

package elektratypes

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/nvidia/carbide-rest/site-agent/pkg/conftypes"
	"github.com/nvidia/carbide-rest/site-agent/pkg/datatypes/managertypes"
	"go.uber.org/atomic"
)

// Elektra is the main struct for the Elektra plugin
type Elektra struct {
	// Main structure of Elektra
	// All information is contained in this structure
	Managers *managertypes.Managers
	Conf     *conftypes.Config
	// HealthStatus current health state
	HealthStatus atomic.Uint64
	Version      string
	Log          zerolog.Logger
}

// NewElektraTypes - create new Elektra Type
func NewElektraTypes() *Elektra {
	return &Elektra{
		Version:  "0.0.1",
		Managers: managertypes.NewManagerType(),
		Conf:     conftypes.NewConfType(),
		Log:      zerolog.New(os.Stderr).With().Timestamp().Logger(),
	}
}
