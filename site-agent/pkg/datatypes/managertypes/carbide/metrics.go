// SPDX-FileCopyrightText: Copyright (c) 2021-2023 NVIDIA CORPORATION & AFFILIATES. All rights reserved.
// SPDX-License-Identifier: LicenseRef-NvidiaProprietary
//
// NVIDIA CORPORATION, its affiliates and licensors retain all intellectual
// property and proprietary rights in and to this material, related
// documentation and any modifications thereto. Any use, reproduction,
// disclosure or distribution of this material and related documentation
// without an express license agreement from NVIDIA CORPORATION or
// its affiliates is strictly prohibited.

package carbidetypes

import "time"

type WorkflowStatus string

const (
	// WorkflowStatusSuccess workflow has completed successfully
	WorkflowStatusSuccess WorkflowStatus = "Success"
	// WorkflowStatusActFailed workflow activity execution has failed
	WorkflowStatusActivityFailed WorkflowStatus = "ActivityFailed"
	// WorkflowStatusPubFailed workflow status publish failed
	WorkflowStatusPublishFailed WorkflowStatus = "PublishFailed"
	// WorkflowStatusActPubFailed both workflow activity execution and status publish failed
	WorkflowStatusActivityPublishFailed WorkflowStatus = "ActivityPublishFailed"
)

// WorkflowMetrics defines interface to be used for workflow metrics
type WorkflowMetrics interface {
	// RecordLatency function to record latency for a workflow
	RecordLatency(activity string, status WorkflowStatus, duration time.Duration)
}
