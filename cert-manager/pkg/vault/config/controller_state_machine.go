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
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path"
	"strings"

	vault "github.com/hashicorp/vault/api"

	"github.com/nvidia/carbide-rest/cert-manager/pkg/core"
)

/*
* This file is a collection of state transition function
* Each transition returns the next state or an error.
 */

// Initialize is an fsm action method
func (v *VaultConfigController) Initialize(ctx context.Context) (VaultState, error) {
	log := core.GetLogger(ctx)

	// If the Controller crashes at any point until we persist the secret
	// we enter an unstable state that can't be recovered automatically
	resp, err := v.client.Sys().Init(&vault.InitRequest{
		SecretShares:    1,
		SecretThreshold: 1,
	})
	if err != nil {
		return Unitialized, err
	}

	v.setToken(resp)
	log.Debug("Got Vault Credentials:", resp)

	tokenPath := v.getSecretPath(SecretTokenPath, SecretTokenName)
	f, err := os.OpenFile(tokenPath, os.O_TRUNC|os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return Unitialized, err
	}
	defer f.Close()
	err = json.NewEncoder(f).Encode(resp)
	if err != nil {
		return Unitialized, err
	}

	log.Debugf("Credentials/State persisted to %s", tokenPath)
	return vaultNextState[Unitialized], nil
}

// Unseal is an fsm action method
func (v *VaultConfigController) Unseal(ctx context.Context) (VaultState, error) {
	resp, err := v.client.Sys().Unseal(v.unsealShard())
	if err != nil {
		return Unitialized, err
	}

	if resp.Sealed {
		core.GetLogger(ctx).Error("Unexpected Sealed state after unsealing, trying to Unseal again")
		return Sealed, nil
	}

	return vaultNextState[Sealed], nil
}

// ConfigurePolicy is an fsm action method
func (v *VaultConfigController) ConfigurePolicy(_ context.Context) (VaultState, error) {
	err := v.client.Sys().PutPolicy(PolicyName, RootIdentityPolicy)
	if err != nil {
		return Unitialized, err
	}

	err = v.client.Sys().PutPolicy(AppRolePolicyName, AppRolePolicy)
	if err != nil {
		return Unitialized, err
	}

	err = v.client.Sys().PutPolicy(CertManagerPolicyName, CertManagerPolicy)
	if err != nil {
		return Unitialized, err
	}

	return vaultNextState[PolicyNotConfigured], nil
}

// ConfigureTokenAuth is an fsm action method
func (v *VaultConfigController) ConfigureTokenAuth(_ context.Context) (VaultState, error) {
	tokenReq := &vault.TokenCreateRequest{
		Policies: []string{CertManagerPolicyName},
		Period:   certManagerTokenPeriod,
	}
	sec, err := v.client.AuthToken().Create(tokenReq)
	if err != nil {
		return Unitialized, err
	}

	authToken, err := sec.TokenID()
	if err != nil || authToken == "" {
		return Unitialized, fmt.Errorf("token is empty")
	}

	err = v.saveCertMgrToken(authToken)
	if err != nil {
		return Unitialized, err
	}
	return vaultNextState[TokenAuthNotConfigured], nil
}

// EnablePKI is an fsm action method
func (v *VaultConfigController) EnablePKI(_ context.Context) (VaultState, error) {
	err := v.client.Sys().Mount("pki", &vault.MountInput{
		Type: "pki",
	})
	if err != nil {
		return Unitialized, err
	}

	return vaultNextState[PKINotEnabled], nil
}

// ConfigurePKICACert is an fsm action method
func (v *VaultConfigController) ConfigurePKICACert(ctx context.Context) (VaultState, error) {
	log := core.GetLogger(ctx)
	err := v.client.Sys().TuneMount("pki", vault.MountConfigInput{
		MaxLeaseTTL: string(desiredPKICACertTTL),
	})
	if err != nil {
		log.Errorf("failed to tune the pki mount, error: %v", err)
		return Unitialized, err
	}

	cert, err := v.getSecret(SecretCACertPath, "certificate")
	if err != nil {
		log.Errorf("failed to retrieve the certificate from the filesystem, error: %v", err)
		return Unitialized, err
	}

	privKey, err := v.getSecret(SecretCAPrivateKeyPath, "privatekey")
	if err != nil {
		log.Errorf("failed to retrieve the private key from the filesystem, error: %v", err)
		return Unitialized, err
	}

	bundle := fmt.Sprintf("%s\n%s", cert, privKey)

	config := v.desiredPKICACertConfig(bundle)

	var raw map[string]interface{}
	if err := core.ConvertToInterfaceMap(config, &raw); err != nil {
		return Unitialized, err
	}

	_, err = v.client.Logical().Write("pki/config/ca", raw)
	if err != nil {
		log.Errorf("failed to configure the CA, error: %v", err)
		return Unitialized, err
	}

	return vaultNextState[PKICACertNotConfigured], nil
}

// ConfigurePKIURLs is an fsm action method
func (v *VaultConfigController) ConfigurePKIURLs(_ context.Context) (VaultState, error) {
	config := v.desiredPKIURLsConfig()

	var raw map[string]interface{}
	if err := core.ConvertToInterfaceMap(config, &raw); err != nil {
		return Unitialized, err
	}

	_, err := v.client.Logical().Write("pki/config/urls", raw)
	if err != nil {
		return Unitialized, err
	}

	return vaultNextState[PKIURLsNotConfigured], nil
}

// ConfigurePKIRole is an fsm action method
func (v *VaultConfigController) ConfigurePKIRole(_ context.Context) (VaultState, error) {
	config := v.desiredPKIRoleConfig()

	var raw map[string]interface{}
	if err := core.ConvertToInterfaceMap(config, &raw); err != nil {
		return Unitialized, err
	}

	_, err := v.client.Logical().Write(fmt.Sprintf("pki/roles/%s", PKIRoleName), raw)
	if err != nil {
		return Unitialized, err
	}

	return vaultNextState[PKIRoleNotConfigured], nil
}

func (v *VaultConfigController) getSecret(secretName string, secretKey string) (string, error) {
	secretPath := v.getSecretPath(secretName, secretKey)
	secretValue, err := os.ReadFile(v.getSecretPath(secretName, secretKey))
	if err != nil {
		return "", fmt.Errorf("error retrieving key %s, error: %v", secretPath, err)
	}

	formattedValue := strings.ReplaceAll(string(secretValue), "\\n", "\n")
	return formattedValue, nil
}

func (v *VaultConfigController) getSecretPath(secretName string, secretKey string) string {
	return path.Join(v.secretsRootPath, secretName, secretKey)
}

func (v *VaultConfigController) saveCertMgrToken(token string) error {
	v.certMgrToken = token
	tokenPath := v.getSecretPath(SecretTokenPath, CertManagerTokenName)
	return os.WriteFile(tokenPath, []byte(token), 0644)
}
