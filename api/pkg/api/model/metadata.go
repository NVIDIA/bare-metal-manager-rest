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

import (
	"github.com/nvidia/carbide-rest/api/pkg/metadata"
)

// APIMetadata is a data structure to capture Forge API system information
type APIMetadata struct {
	// Version contains the API version
	Version string `json:"version"`
	// BuildTime contains the time the binary was built
	BuildTime string `json:"buildTime"`
}

// NewAPIMetadata creates and returns a new APISystemInfo object
func NewAPIMetadata() *APIMetadata {
	amd := &APIMetadata{
		Version:   metadata.Version,
		BuildTime: metadata.BuildTime,
	}

	return amd
}
