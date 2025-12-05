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

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestVaultConfigController_getSecret(t *testing.T) {
	vcc, err := NewController("http://localhost", "http://localhost",
		"testdata/secrets")
	require.NoError(t, err)
	require.NotNil(t, vcc)

	s, err := vcc.getSecret("token", "vault-token")
	assert.NoError(t, err)
	assert.Contains(t, s, "774d5e5a7792a7a0f077291e272b4b0eb42b7ad2c18d5c7f4ad29ebb4c46be77")

	s, err = vcc.getSecret("token", "does-not-exist.json")
	assert.Error(t, err)
	assert.Empty(t, s)
}
