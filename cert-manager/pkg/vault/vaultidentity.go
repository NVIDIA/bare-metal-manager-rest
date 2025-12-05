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
	"fmt"
	"strings"
	"sync"

	vault "github.com/hashicorp/vault/api"
	"github.com/nvidia/carbide-rest/cert-manager/pkg/core"
	"github.com/nvidia/carbide-rest/cert-manager/pkg/vault/config"
)

// Identity defines the vault interface
type Identity interface {
	GetRoleTemplate(ctx context.Context, req *TokenRequest) string
	EnsureAppRole(ctx context.Context, role string, policies []string, secret string) error
	EnsureIdentityRole(ctx context.Context, role, key, ttl, template, clientID string) error
	RequestIdentityToken(ctx context.Context, role string) (string, error)
	RequestIdentityCertificate(ctx context.Context, sans string, ttl int) (string, string, error)
	RequestCACertificate(ctx context.Context) (string, error)
	RequestCRL(ctx context.Context) (string, error)
}

type vaultIdentity struct {
	core.VaultOptions

	newClient func(ctx context.Context, o core.VaultOptions) (vaultClient, error)
	sync.Mutex
	vc vaultClient
}

type vaultClient interface {
	LogicalRead(ctx context.Context, path string) (*vault.Secret, error)
	LogicalWrite(ctx context.Context, path string, data map[string]interface{}) (*vault.Secret, error)
}

// NewVaultIdentity returns a VaultIdentity
func NewVaultIdentity(_ context.Context, o core.VaultOptions) Identity {
	i := &vaultIdentity{VaultOptions: o}
	i.newClient = func(ctx context.Context, o core.VaultOptions) (vaultClient, error) {
		return core.NewVaultClient(ctx, o)
	}
	return i
}

func (i *vaultIdentity) client(ctx context.Context) (vaultClient, error) {
	i.Lock()
	defer i.Unlock()
	if i.vc == nil {
		c, err := i.newClient(ctx, i.VaultOptions)
		if err != nil {
			return nil, err
		}
		i.vc = c
	}
	return i.vc, nil
}

func (i *vaultIdentity) GetRoleTemplate(_ context.Context, req *TokenRequest) string {
	username := req.UniqueName()

	groupsWithQuote := []string{}
	for _, g := range req.Groups {
		groupsWithQuote = append(groupsWithQuote, fmt.Sprintf("%q", g))
	}
	groupsStr := strings.Join(groupsWithQuote, ",")

	switch req.Type {
	case TokenTypeDevice:
		return fmt.Sprintf(`{
  "username": "%s",
  "groups": [%s],
  "nbf": {{time.now}},
  "forge.nvidia.com/token_type": "%s",
  "forge.nvidia.com/device_id": "%s"
}`, username, groupsStr, string(req.Type), req.DeviceID)

	case TokenTypeUser:
		return fmt.Sprintf(`{
  "username": "%s",
  "groups": [%s],
  "nbf": {{time.now}},
  "forge.nvidia.com/token_type": "%s"
}`, username, groupsStr, string(req.Type))

	case TokenTypeApp:
		return fmt.Sprintf(`{
  "username": "%s",
  "groups": [%s],
  "nbf": {{time.now}},
  "forge.nvidia.com/token_type": "%s"
}`, username, groupsStr, string(req.Type))
	}

	return ""
}

func (i *vaultIdentity) EnsureAppRole(ctx context.Context, role string, policies []string, secret string) error {
	c, err := i.client(ctx)
	if err != nil {
		return err
	}

	data := map[string]interface{}{"role_id": role, "policies": policies}
	path := fmt.Sprintf("auth/approle/role/%s", role)
	_, err = c.LogicalWrite(ctx, path, data)
	if err != nil {
		return err
	}

	path = fmt.Sprintf("auth/approle/role/%s/custom-secret-id", role)
	data = map[string]interface{}{"secret_id": secret}
	_, err = c.LogicalWrite(ctx, path, data)
	if err != nil {
		// Ignore "failed to store secret_id: SecretID is already
		// registered" errors, enforce "Ensure" semantics
		if strings.Contains(err.Error(), "is already registered") {
			return nil
		}
		return err
	}
	return nil
}

func (i *vaultIdentity) EnsureIdentityRole(ctx context.Context, role, key, ttl, template, clientID string) error {
	c, err := i.client(ctx)
	if err != nil {
		return err
	}

	path := fmt.Sprintf("identity/oidc/role/%s", role)
	data := map[string]interface{}{
		"key":       key,
		"ttl":       ttl,
		"template":  template,
		"client_id": clientID,
	}
	_, err = c.LogicalWrite(ctx, path, data)
	return err
}

func (i *vaultIdentity) RequestIdentityToken(ctx context.Context, role string) (string, error) {
	c, err := i.client(ctx)
	if err != nil {
		return "", err
	}

	path := fmt.Sprintf("identity/oidc/token/%s", role)
	resp, err := c.LogicalRead(ctx, path)
	if err != nil {
		return "", err
	}

	token, err := getFieldFromSecret("token", path, resp)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (i *vaultIdentity) RequestIdentityCertificate(ctx context.Context, sans string, ttl int) (string, string, error) {
	c, err := i.client(ctx)
	if err != nil {
		return "", "", err
	}

	ttlHours := fmt.Sprintf("%dh", ttl)
	certData := map[string]interface{}{
		"common_name": sans,
		"ttl":         ttlHours,
	}

	path := fmt.Sprintf("pki/issue/%s", config.PKIRoleName)
	resp, err := c.LogicalWrite(ctx, path, certData)
	if err != nil {
		return "", "", err
	}

	cert, err := getFieldFromSecret("certificate", path, resp)
	if err != nil {
		return "", "", err
	}

	privKey, err := getFieldFromSecret("private_key", path, resp)
	if err != nil {
		return "", "", err
	}

	return cert, privKey, nil
}

func (i *vaultIdentity) RequestCACertificate(ctx context.Context) (string, error) {
	c, err := i.client(ctx)
	if err != nil {
		return "", err
	}

	path := "pki/cert/ca"
	resp, err := c.LogicalRead(ctx, path)
	if err != nil {
		return "", err
	}

	cert, err := getFieldFromSecret("certificate", path, resp)
	if err != nil {
		return "", err
	}

	return cert, nil
}

func (i *vaultIdentity) RequestCRL(ctx context.Context) (string, error) {
	c, err := i.client(ctx)
	if err != nil {
		return "", err
	}

	path := "pki/cert/crl"
	resp, err := c.LogicalRead(ctx, path)
	if err != nil {
		return "", err
	}

	cert, err := getFieldFromSecret("certificate", path, resp)
	if err != nil {
		return "", err
	}

	return cert, nil
}

func getFieldFromSecret(secretField, vaultPath string, secret *vault.Secret) (string, error) {
	if secret == nil {
		return "", fmt.Errorf("resp is nil from %s", vaultPath)
	}
	if secret.Data == nil {
		return "", fmt.Errorf("resp.Data is nil from %s", vaultPath)
	}
	value, ok := secret.Data[secretField]
	if !ok {
		return "", fmt.Errorf("resp.Data[\"%s\"] is nil from %s", secretField, vaultPath)
	}
	valueStr, ok := value.(string)
	if !ok {
		return "", fmt.Errorf("resp.Data[\"%s\"] is not a string from %s", secretField, vaultPath)
	}
	return valueStr, nil
}
