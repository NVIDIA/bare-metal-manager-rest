// SPDX-FileCopyrightText: Copyright (c) 2021-2023 NVIDIA CORPORATION & AFFILIATES. All rights reserved.
// SPDX-License-Identifier: LicenseRef-NvidiaProprietary
//
// NVIDIA CORPORATION, its affiliates and licensors retain all intellectual
// property and proprietary rights in and to this material, related
// documentation and any modifications thereto. Any use, reproduction,
// disclosure or distribution of this material and related documentation
// without an express license agreement from NVIDIA CORPORATION or
// its affiliates is strictly prohibited.

package health

import (
	"encoding/json"
	"net/http"

	"github.com/rs/zerolog/log"
)

// Check captures the API response for workflow service health check
type Check struct {
	IsHealthy bool    `json:"is_healthy"`
	Error     *string `json:"error"`
}

// StatusHandler is an API handler to return health status of the workflow service
func StatusHandler(w http.ResponseWriter, r *http.Request) {
	check := Check{
		IsHealthy: true,
	}
	bytes, err := json.Marshal(check)
	if err != nil {
		log.Error().Err(err).Msg("error converting health check object into JSON")
		http.Error(w, "failed to construct health check response", http.StatusInternalServerError)
		return
	}
	_, err = w.Write(bytes)
	if err != nil {
		log.Error().Err(err).Msg("failed to return health check response")
	}
}
