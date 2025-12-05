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
)

// TokenIssuer defines an issuer interface
type TokenIssuer interface {
	NewToken(ctx context.Context, req *TokenRequest) (string, error)
}

// TokenType defines the token type
type TokenType string

// Token type definitions
const (
	TokenTypeDevice TokenType = "device"
	TokenTypeUser   TokenType = "user"
	TokenTypeApp    TokenType = "app"
)

// TokenRequest defines the request
type TokenRequest struct {
	// Type is the token type
	Type TokenType

	// Name should be unique for a particular TokenType
	Name   string
	Groups []string

	// DeviceID is base64 encode sha256 hash of akPub, only applicable
	// if token type is device.
	DeviceID string
}

// UniqueName returns a unique name across different TokenType from a
// TokenIssuer
func (r *TokenRequest) UniqueName() string {
	return fmt.Sprintf("%s-%s", string(r.Type), r.Name)
}
