// SPDX-FileCopyrightText: Copyright (c) 2021-2023 NVIDIA CORPORATION & AFFILIATES. All rights reserved.
// SPDX-License-Identifier: LicenseRef-NvidiaProprietary
//
// NVIDIA CORPORATION, its affiliates and licensors retain all intellectual
// property and proprietary rights in and to this material, related
// documentation and any modifications thereto. Any use, reproduction,
// disclosure or distribution of this material and related documentation
// without an express license agreement from NVIDIA CORPORATION or
// its affiliates is strictly prohibited.

package pki

// Native Go PKI Certificate Issuer
//
// This replaces the embedded Vault dependency for certificate generation.
// The issuer loads a CA from files, matching how Vault loaded its CA from secrets.
//
// CA Loading Order:
//   1. Primary path: --ca-cert-file / --ca-key-file (default: /vault/secrets/...)
//   2. Alternate path: --alt-ca-cert-file / --alt-ca-key-file (default: /etc/pki/ca/...)
//   3. Error if no CA found
//
// The CA must be provided via K8s secrets. There is no fallback to self-signed.

import (
	"context"
	"fmt"

	"github.com/nvidia/carbide-rest/cert-manager/pkg/types"
)

// NativeCertificateIssuer implements types.CertificateIssuer using native Go crypto
type NativeCertificateIssuer struct {
	ca      *CA
	baseDNS string
}

// NativeCertificateIssuerOptions defines options for the native issuer
type NativeCertificateIssuerOptions struct {
	BaseDNS        string
	CertificateTTL string
	CACommonName   string
	CAOrganization string
	// CACertFile and CAKeyFile are the primary paths (vault-style paths)
	CACertFile string
	CAKeyFile  string
	// AltCACertFile and AltCAKeyFile are alternate paths for easier migration
	AltCACertFile string
	AltCAKeyFile  string
}

// NewNativeCertificateIssuer creates a new native Go certificate issuer
// Tries to load CA from primary paths first, then alternate paths.
// Returns an error if no CA is found - there is no fallback.
func NewNativeCertificateIssuer(opts NativeCertificateIssuerOptions) (types.CertificateIssuer, error) {
	var ca *CA
	var err error
	var loadErr error

	// Try primary paths (vault-style)
	if opts.CACertFile != "" && opts.CAKeyFile != "" {
		ca, err = LoadCA(opts.CACertFile, opts.CAKeyFile)
		if err == nil {
			fmt.Printf("Loaded CA from primary path: %s\n", opts.CACertFile)
			return &NativeCertificateIssuer{
				ca:      ca,
				baseDNS: opts.BaseDNS,
			}, nil
		}
		loadErr = fmt.Errorf("primary path (%s): %w", opts.CACertFile, err)
	}

	// Try alternate paths
	if opts.AltCACertFile != "" && opts.AltCAKeyFile != "" {
		ca, err = LoadCA(opts.AltCACertFile, opts.AltCAKeyFile)
		if err == nil {
			fmt.Printf("Loaded CA from alternate path: %s\n", opts.AltCACertFile)
			return &NativeCertificateIssuer{
				ca:      ca,
				baseDNS: opts.BaseDNS,
			}, nil
		}
		if loadErr != nil {
			loadErr = fmt.Errorf("%v; alternate path (%s): %w", loadErr, opts.AltCACertFile, err)
		} else {
			loadErr = fmt.Errorf("alternate path (%s): %w", opts.AltCACertFile, err)
		}
	}

	// No CA found - error
	if loadErr != nil {
		return nil, fmt.Errorf("CA certificate required but not found: %w", loadErr)
	}
	return nil, fmt.Errorf("CA certificate required: no paths configured")
}

// NewCertificate implements types.CertificateIssuer
func (i *NativeCertificateIssuer) NewCertificate(ctx context.Context, req *types.CertificateRequest) (string, string, error) {
	sans := req.UniqueName(i.baseDNS)
	ttl := req.TTL
	if ttl == 0 {
		ttl = 24 * 90 // 90 days default
	}
	return i.ca.IssueCertificate(sans, ttl)
}

// RawCertificate implements types.CertificateIssuer
func (i *NativeCertificateIssuer) RawCertificate(ctx context.Context, sans string, ttl int) (string, string, error) {
	return i.ca.IssueCertificate(sans, ttl)
}

// GetCACertificate implements types.CertificateIssuer
func (i *NativeCertificateIssuer) GetCACertificate(ctx context.Context) (string, error) {
	return i.ca.GetCACertificatePEM(), nil
}

// GetCRL implements types.CertificateIssuer
func (i *NativeCertificateIssuer) GetCRL(ctx context.Context) (string, error) {
	return i.ca.GetCRL(), nil
}
