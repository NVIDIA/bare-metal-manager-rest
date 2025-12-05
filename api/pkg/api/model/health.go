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


package model

// APIHealthCheck is a data structure to capture Forge API health information
type APIHealthCheck struct {
	// IsHealthy provides a flag to accompany an error status code
	IsHealthy bool `json:"is_healthy"`
	// Error contains an error message in case of health issues
	Error *string `json:"error"`
}

// NewAPIHealthCheck creates and returns a new APIHealthCheck object
func NewAPIHealthCheck(isHealthy bool, errorMessage *string) *APIHealthCheck {
	ahc := &APIHealthCheck{
		IsHealthy: isHealthy,
		Error:     errorMessage,
	}

	return ahc
}
