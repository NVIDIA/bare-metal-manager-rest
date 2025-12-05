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
	"time"

	wflows "github.com/nvidia/carbide-rest/workflow-schema/schema/site-agent/workflows/v1"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
)

func (ac *HealthWorkflow) GetHealthActivity() (status wflows.HealthStatus, err error) {
	status = wflows.HealthStatus{
		SiteInventoryCollection:   &wflows.HealthStatusMsg{},
		SiteControllerConnection:  &wflows.HealthStatusMsg{},
		SiteAgentHighAvailability: &wflows.HealthStatusMsg{},
	}

	status.Timestamp = timestamppb.New(time.Now())
	status.SiteInventoryCollection.State = ManagerAccess.Data.EB.Managers.Health.Inventory.State
	status.SiteAgentHighAvailability.State = ManagerAccess.Data.EB.Managers.Health.Availabilty.State
	status.SiteControllerConnection.State = ManagerAccess.Data.EB.Managers.Health.CarbideInterface.State

	return status, err
}
