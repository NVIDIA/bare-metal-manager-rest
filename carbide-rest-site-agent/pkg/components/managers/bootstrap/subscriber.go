// SPDX-FileCopyrightText: Copyright (c) 2021-2023 NVIDIA CORPORATION & AFFILIATES. All rights reserved.
// SPDX-License-Identifier: LicenseRef-NvidiaProprietary
//
// NVIDIA CORPORATION, its affiliates and licensors retain all intellectual
// property and proprietary rights in and to this material, related
// documentation and any modifications thereto. Any use, reproduction,
// disclosure or distribution of this material and related documentation
// without an express license agreement from NVIDIA CORPORATION or
// its affiliates is strictly prohibited.

package bootstrap

import (
	"go.temporal.io/sdk/activity"
	"go.temporal.io/sdk/workflow"
)

// RegisterSubscriber registers the Bootstrap workflows and activities with the Temporal client
func (api *BoostrapAPI) RegisterSubscriber() error {
	// Initialize logger
	logger := ManagerAccess.Data.EB.Log

	// Only master pod should watch for the OTP rotation workflow
	if !ManagerAccess.Conf.EB.IsMasterPod {
		return nil
	}

	// Register the workflows
	wflowRegisterOptions := workflow.RegisterOptions{
		Name: "RotateTemporalCertAccessOTP",
	}
	ManagerAccess.Data.EB.Managers.Workflow.Temporal.Worker.RegisterWorkflowWithOptions(api.RotateTemporalCertAccessOTP, wflowRegisterOptions)
	logger.Info().Msg("Bootstrap: successfully registered the ReceiveAndProcessOTP workflow")

	// Register the activities
	otpHandler := NewOTPHandler(ManagerAccess.Data.EB.Managers.Bootstrap.Secret)

	activityRegisterOptions := activity.RegisterOptions{
		Name: "ReceiveAndSaveOTP",
	}
	ManagerAccess.Data.EB.Managers.Workflow.Temporal.Worker.RegisterActivityWithOptions(otpHandler.ReceiveAndSaveOTP, activityRegisterOptions)
	logger.Info().Msg("Bootstrap: successfully registered the ReceiveAndSaveOTP activity")

	return nil
}
