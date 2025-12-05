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
	"crypto/x509"
	"encoding/pem"
	"fmt"

	"github.com/nvidia/carbide-rest/cert-manager/pkg/core"
)

/*
* This file is a collection of function that aimed to discover the current vault state.
* The state is discovered by calling these functions in order.
* If a function returns false, then vault is in the state indicated by the function name.
*
* Note that currently all functions are structured as a succession of tests and should have
* only one "return false", at the end of the function's body.
 */

func (v *VaultConfigController) isUninitialized(_ context.Context) (bool, error) {
	health, err := v.client.Sys().Health()
	if err != nil {
		return true, err
	}

	if !health.Initialized {
		return true, nil
	}

	return false, nil
}

func (v *VaultConfigController) isSealed(_ context.Context) (bool, error) {
	health, err := v.client.Sys().Health()
	if err != nil {
		return true, err
	}

	if health.Sealed {
		return true, nil
	}

	return false, nil
}

func (v *VaultConfigController) isPolicyNotConfigured(_ context.Context) (bool, error) {
	err := v.checkPolicy(PolicyName, RootIdentityPolicy)
	if err != nil {
		return true, err
	}

	err = v.checkPolicy(AppRolePolicyName, AppRolePolicy)
	if err != nil {
		return true, err
	}

	err = v.checkPolicy(CertManagerPolicyName, CertManagerPolicy)
	if err != nil {
		return true, err
	}

	return false, nil
}

func (v *VaultConfigController) isTokenAuthNotConfigured(ctx context.Context) (bool, error) {
	log := core.GetLogger(ctx)

	token, err := v.getSecret(SecretTokenPath, CertManagerTokenName)
	if err != nil {
		return true, err
	}

	vct := v.client.AuthToken()
	_, err = vct.Lookup(token)

	if err != nil {
		return true, err
	}

	log.Infof("cert manager token exists")
	v.certMgrToken = token
	return false, nil
}

func (v *VaultConfigController) isPKINotEnabled(_ context.Context) (bool, error) {
	mounts, err := v.client.Sys().ListMounts()
	if err != nil { // Would happen if the client is not authenticated
		return true, err
	}

	if _, ok := mounts["pki/"]; !ok {
		return true, nil
	}

	return false, nil
}

func (v *VaultConfigController) isPKICACertNotConfigured(ctx context.Context) (bool, error) {
	log := core.GetLogger(ctx)
	resp, err := v.client.Logical().Read("pki/cert/ca")
	if err != nil { // Would happen if the client is not authenticated
		log.Errorf("pki/cert/ca retrieval failed from vault: %v", err)
		return true, err
	}

	if resp == nil || resp.Data == nil {
		log.Errorf("pki/cert/ca retrieval returned no data in the response: %v", err)
		return true, nil
	}

	cert, ok := resp.Data["certificate"]
	if !ok {
		log.Errorf("pki/cert/ca response did not contain a certificate field in its data")
		return true, nil
	}

	expectedCert, err := v.getSecret(SecretCACertPath, "certificate")
	if err != nil {
		log.Errorf("error retrieving certificate from secret %s: %v", SecretCACertPath, err)
		return true, err
	}

	if ceq, err := compareCerts(expectedCert, cert.(string)); !ceq {
		log.Infof("expected and actual certs are not equal error: %v", err)
		return true, err
	}

	return false, nil
}

// compareCerts returns true if the certs are equivalent,
// false otherwise.
func compareCerts(certA, certB string) (bool, error) {
	blockA, _ := pem.Decode([]byte(certA))
	if blockA == nil {
		return false, fmt.Errorf("invalid cert %s", certA)
	}
	ca, err := x509.ParseCertificate(blockA.Bytes)
	if err != nil {
		return false, err
	}

	blockB, _ := pem.Decode([]byte(certB))
	if blockB == nil {
		return false, fmt.Errorf("invalid cert %s", certB)
	}
	cb, err := x509.ParseCertificate(blockB.Bytes)
	if err != nil {
		return false, err
	}

	return ca.Equal(cb), nil
}

func (v *VaultConfigController) isPKIURLsNotConfigured(_ context.Context) (bool, error) {
	resp, err := v.client.Logical().Read("pki/config/urls")
	if err != nil { // Would happen if the client is not authenticated
		return true, err
	}

	if resp == nil {
		return true, nil
	}

	desiredConfig := v.desiredPKIURLsConfig()

	var actualConfig vaultPKIURLsConfig
	err = core.ConvertToStruct(resp.Data, &actualConfig)
	if err != nil {
		return true, nil
	}

	if sliceNotEqual(actualConfig.IssuingCertificates, desiredConfig.IssuingCertificates) {
		return true, nil
	}

	if sliceNotEqual(actualConfig.CRLDistributionPoints, desiredConfig.CRLDistributionPoints) {
		return true, nil
	}

	return false, nil
}

func (v *VaultConfigController) isPKIRoleNotConfigured(_ context.Context) (bool, error) {
	resp, err := v.client.Logical().Read(fmt.Sprintf("pki/roles/%s", PKIRoleName))
	if err != nil { // Would happen if the client is not authenticated
		return true, err
	}

	if resp == nil {
		return true, nil
	}

	desiredConfig := v.desiredPKIRoleConfig()

	var actualConfig vaultPKIRoleConfig
	err = core.ConvertToStruct(resp.Data, &actualConfig)
	if err != nil {
		return true, err
	}

	if actualConfig.AllowAnyName != desiredConfig.AllowAnyName {
		return true, nil
	}

	if actualConfig.MaxTTL != desiredConfig.MaxTTL {
		return true, nil
	}

	if actualConfig.NotBeforeDuration != desiredConfig.NotBeforeDuration {
		return true, nil
	}

	if sliceNotEqual(actualConfig.Organization, desiredConfig.Organization) {
		return true, nil
	}

	return false, nil
}

func (v *VaultConfigController) checkPolicy(name, pol string) error {
	p, err := v.client.Sys().GetPolicy(name)
	if err != nil {
		return err
	}

	if p != pol {
		return fmt.Errorf("policy mismatch")
	}
	return nil
}

func sliceNotEqual(a, b []string) bool {
	if len(a) != len(b) {
		return true
	}

	for i := 0; i < len(a); i++ {
		if a[i] != b[i] {
			return true
		}
	}

	return false
}
