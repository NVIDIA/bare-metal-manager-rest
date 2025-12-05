// SPDX-FileCopyrightText: Copyright (c) 2021-2023 NVIDIA CORPORATION & AFFILIATES. All rights reserved.
// SPDX-License-Identifier: LicenseRef-NvidiaProprietary
//
// NVIDIA CORPORATION, its affiliates and licensors retain all intellectual
// property and proprietary rights in and to this material, related
// documentation and any modifications thereto. Any use, reproduction,
// disclosure or distribution of this material and related documentation
// without an express license agreement from NVIDIA CORPORATION or
// its affiliates is strictly prohibited.

package sku

// RegisterSubscriber registers the SKU workflows/activities with the Temporal client
// Note: SKU does not have create/update/delete capabilities, so no subscriber workflows are registered
func (api *API) RegisterSubscriber() error {
	// Register the subscribers here
	ManagerAccess.Data.EB.Log.Info().Msg("SKU: Registering the subscribers")
	// Note: SKU is read-only, no CRUD workflows to register
	ManagerAccess.Data.EB.Log.Info().Msg("SKU: No CRUD workflows for SKU (read-only resource)")

	return nil
}
