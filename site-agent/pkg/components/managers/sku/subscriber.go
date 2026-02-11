/*
 * SPDX-FileCopyrightText: Copyright (c) 2026 NVIDIA CORPORATION & AFFILIATES. All rights reserved.
 * SPDX-License-Identifier: Apache-2.0
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

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
