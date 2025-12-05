// SPDX-FileCopyrightText: Copyright (c) 2021-2023 NVIDIA CORPORATION & AFFILIATES. All rights reserved.
// SPDX-License-Identifier: LicenseRef-NvidiaProprietary
//
// NVIDIA CORPORATION, its affiliates and licensors retain all intellectual
// property and proprietary rights in and to this material, related
// documentation and any modifications thereto. Any use, reproduction,
// disclosure or distribution of this material and related documentation
// without an express license agreement from NVIDIA CORPORATION or
// its affiliates is strictly prohibited.

package vault

import (
	"context"
	"sync"

	vault "github.com/hashicorp/vault/api"
)

// mockVaultClient mocks core.VaultClient for testing
type mockVaultClient struct {
	sync.Mutex

	readPath            []string
	readResult          []*vault.Secret
	readErr             []error
	logicalReadOverride func(ctx context.Context, path string) (*vault.Secret, error)

	writePath            []string
	writeData            []map[string]interface{}
	writeResult          []*vault.Secret
	writeErr             []error
	logicalWriteOverride func(ctx context.Context, path string, data map[string]interface{}) (*vault.Secret, error)
}

func (m *mockVaultClient) LogicalRead(ctx context.Context, path string) (*vault.Secret, error) {
	m.Lock()
	defer m.Unlock()
	if m.logicalReadOverride != nil {
		return m.logicalReadOverride(ctx, path)
	}
	res, err := m.readResult[0], m.readErr[0]
	m.readPath, m.readResult, m.readErr = append(m.readPath, path), m.readResult[1:], m.readErr[1:]
	return res, err
}

func (m *mockVaultClient) LogicalWrite(ctx context.Context, path string, data map[string]interface{}) (*vault.Secret, error) {
	m.Lock()
	defer m.Unlock()
	if m.logicalWriteOverride != nil {
		return m.logicalWriteOverride(ctx, path, data)
	}
	res, err := m.writeResult[0], m.writeErr[0]
	m.writePath, m.writeData = append(m.writePath, path), append(m.writeData, data)
	m.writeResult, m.writeErr = m.writeResult[1:], m.writeErr[1:]
	return res, err
}
