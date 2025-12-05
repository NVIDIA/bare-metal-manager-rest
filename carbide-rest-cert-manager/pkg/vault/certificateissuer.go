// SPDX-FileCopyrightText: Copyright (c) 2021-2023 NVIDIA CORPORATION & AFFILIATES. All rights reserved.
// SPDX-License-Identifier: LicenseRef-NvidiaProprietary
//
// NVIDIA CORPORATION, its affiliates and licensors retain all intellectual
// property and proprietary rights in and to this material, related
// documentation and any modifications thereto. Any use, reproduction,
// disclosure or distribution of this material and related documentation
// without an express license agreement from NVIDIA CORPORATION or
// its affiliates is strictly prohibited.

// Package vault implements vault wrapper
package vault

import (
	"context"
	"fmt"
	"strings"
)

// CertificateIssuer defines a vault certificate issuer
type CertificateIssuer interface {
	NewCertificate(ctx context.Context, req *CertificateRequest) (string, string, error)
	GetCACertificate(ctx context.Context) (string, error)
	GetCRL(ctx context.Context) (string, error)
	RawCertificate(ctx context.Context, sans string, ttl int) (string, string, error)
}

// CertificateRequest defines a request
type CertificateRequest struct {
	// Name should be unique for a particular CertificateType
	Name string `json:"name,omitempty"`

	// App identifies the specific app
	App string `json:"app,omitempty"`

	TTL int `json:"ttl,omitempty"`
}

// CertificateResponse defines a response
type CertificateResponse struct {
	Key         string `json:"key,omitempty"`
	Certificate string `json:"certificate,omitempty"`
}

// UniqueName returns a sans per node/app combination
func (r *CertificateRequest) UniqueName(baseDNS string) string {
	var sans string

	if r.App == "" {
		sans = r.Name
	} else {
		sans = fmt.Sprintf("%s.%s", r.App, r.Name)
	}

	if strings.HasSuffix(sans, ".com") || strings.HasSuffix(sans, ".local") {
		return sans
	}

	// add baseDNS
	return fmt.Sprintf("%s.%s", sans, baseDNS)
}
