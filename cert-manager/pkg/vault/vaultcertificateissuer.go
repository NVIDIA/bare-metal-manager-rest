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

	"github.com/nvidia/carbide-rest/cert-manager/pkg/core"
)

// CertificateIssuerOptions defines options for vault
type CertificateIssuerOptions struct {
	core.VaultOptions

	BaseDNS        string
	CertificateTTL string
}

type vaultCertificateIssuer struct {
	CertificateIssuerOptions

	root func(ctx context.Context) Identity
}

// NewVaultCertificateIssuer returns a new issuer
func NewVaultCertificateIssuer(_ context.Context, o CertificateIssuerOptions) CertificateIssuer {
	i := &vaultCertificateIssuer{CertificateIssuerOptions: o}
	i.root = func(ctx context.Context) Identity {
		return NewVaultIdentity(ctx, i.CertificateIssuerOptions.VaultOptions)
	}
	return i
}

// NewCertificate implements the NewCertificate method
func (i *vaultCertificateIssuer) NewCertificate(ctx context.Context, req *CertificateRequest) (string, string, error) {
	sans := req.UniqueName(i.CertificateIssuerOptions.BaseDNS)
	return i.root(ctx).RequestIdentityCertificate(ctx, sans, req.TTL)
}

// RawCertificate implements the RawCertificate method
func (i *vaultCertificateIssuer) RawCertificate(ctx context.Context, sans string, ttl int) (string, string, error) {
	return i.root(ctx).RequestIdentityCertificate(ctx, sans, ttl)
}

// GetCACertificate implements the GetCACertificate method
func (i *vaultCertificateIssuer) GetCACertificate(ctx context.Context) (string, error) {
	return i.root(ctx).RequestCACertificate(ctx)
}

// GetCRL implements the GetCRL method
func (i *vaultCertificateIssuer) GetCRL(ctx context.Context) (string, error) {
	return i.root(ctx).RequestCRL(ctx)
}
