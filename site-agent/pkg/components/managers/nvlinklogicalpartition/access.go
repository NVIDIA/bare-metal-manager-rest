// SPDX-FileCopyrightText: Copyright (c) 2021-2023 NVIDIA CORPORATION & AFFILIATES. All rights reserved.
// SPDX-License-Identifier: LicenseRef-NvidiaProprietary
//
// NVIDIA CORPORATION, its affiliates and licensors retain all intellectual
// property and proprietary rights in and to this material, related
// documentation and any modifications thereto. Any use, reproduction,
// disclosure or distribution of this material and related documentation
// without an express license agreement from NVIDIA CORPORATION or
// its affiliates is strictly prohibited.

package nvlinklogicalpartition

import (
	Manager "github.com/nvidia/carbide-rest/site-agent/pkg/components/managers/managerapi"
	"github.com/nvidia/carbide-rest/site-agent/pkg/datatypes/elektratypes"
)

// ManagerAccess - access to all managers
var ManagerAccess *Manager.ManagerAccess

// API - all API interface
type API struct{}

// NewNVLinkLogicalPartitionManager - returns a new instance of NVLinkLogicalPartition manager
func NewNVLinkLogicalPartitionManager(superForge *elektratypes.Elektra, superAPI *Manager.ManagerAPI, superConf *Manager.ManagerConf) *API {
	ManagerAccess = &Manager.ManagerAccess{
		Data: &Manager.ManagerData{
			EB: superForge,
		},
		API:  superAPI,
		Conf: superConf,
	}
	return &API{}
}
