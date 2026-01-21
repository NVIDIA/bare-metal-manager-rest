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
	"github.com/nvidia/carbide-rest/cert-manager/pkg/types"
)

// CertificateIssuer is an alias for types.CertificateIssuer for backward compatibility
type CertificateIssuer = types.CertificateIssuer

// CertificateRequest is an alias for types.CertificateRequest for backward compatibility
type CertificateRequest = types.CertificateRequest

// CertificateResponse is an alias for types.CertificateResponse for backward compatibility
type CertificateResponse = types.CertificateResponse
