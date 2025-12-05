// SPDX-FileCopyrightText: Copyright (c) 2021-2023 NVIDIA CORPORATION & AFFILIATES. All rights reserved.
// SPDX-License-Identifier: LicenseRef-NvidiaProprietary
//
// NVIDIA CORPORATION, its affiliates and licensors retain all intellectual
// property and proprietary rights in and to this material, related
// documentation and any modifications thereto. Any use, reproduction,
// disclosure or distribution of this material and related documentation
// without an express license agreement from NVIDIA CORPORATION or
// its affiliates is strictly prohibited.

package config

// PayloadEncryptionConfig holds configuration for payload encryption
type PayloadEncryptionConfig struct {
	EncryptionKey string
}

// NewPayloadEncryptionConfig initializes and returns a configuration object for payload encryption
func NewPayloadEncryptionConfig(encryptionKey string) *PayloadEncryptionConfig {
	return &PayloadEncryptionConfig{
		EncryptionKey: encryptionKey,
	}
}
