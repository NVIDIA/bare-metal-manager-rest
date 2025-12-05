// SPDX-FileCopyrightText: Copyright (c) 2021-2023 NVIDIA CORPORATION & AFFILIATES. All rights reserved.
// SPDX-License-Identifier: LicenseRef-NvidiaProprietary
//
// NVIDIA CORPORATION, its affiliates and licensors retain all intellectual
// property and proprietary rights in and to this material, related
// documentation and any modifications thereto. Any use, reproduction,
// disclosure or distribution of this material and related documentation
// without an express license agreement from NVIDIA CORPORATION or
// its affiliates is strictly prohibited.

package processors

import (
	"github.com/nvidia/carbide-rest/auth/pkg/config"
	commonConfig "github.com/nvidia/carbide-rest/common/pkg/config"
	cdb "github.com/nvidia/carbide-rest/db/pkg/db"
	temporalClient "go.temporal.io/sdk/client"
)

// NewKeycloakProcessor creates a new Keycloak token processor
func NewKeycloakProcessor(dbSession *cdb.Session, kcfg *config.KeycloakConfig) config.TokenProcessor {
	return &KeycloakProcessor{
		dbSession:      dbSession,
		keycloakConfig: kcfg,
	}
}

// NewSSAProcessor creates a new SSA token processor
func NewSSAProcessor(dbSession *cdb.Session) config.TokenProcessor {
	return &SSAProcessor{
		dbSession: dbSession,
	}
}

// NewKASProcessor creates a new KAS token processor
func NewKASProcessor(dbSession *cdb.Session, tc temporalClient.Client, encCfg *commonConfig.PayloadEncryptionConfig) config.TokenProcessor {
	return &KASProcessor{
		dbSession: dbSession,
		tc:        tc,
		encCfg:    encCfg,
	}
}

// NewCustomProcessor creates a new custom token processor
func NewCustomProcessor(dbSession *cdb.Session) config.TokenProcessor {
	return &CustomProcessor{
		dbSession: dbSession,
	}
}

// InitializeProcessors sets up all token processors in the JWTOriginConfig
func InitializeProcessors(joCfg *config.JWTOriginConfig, dbSession *cdb.Session, tc temporalClient.Client, encCfg *commonConfig.PayloadEncryptionConfig, kcfg *config.KeycloakConfig) {
	for _, origin := range []int{config.TokenOriginKeycloak, config.TokenOriginSsa, config.TokenOriginKas, config.TokenOriginCustom} {
		switch origin {
		case config.TokenOriginKeycloak:
			processor := NewKeycloakProcessor(dbSession, kcfg)
			joCfg.SetProcessorForOrigin(origin, processor)
		case config.TokenOriginSsa:
			processor := NewSSAProcessor(dbSession)
			joCfg.SetProcessorForOrigin(origin, processor)
		case config.TokenOriginKas:
			processor := NewKASProcessor(dbSession, tc, encCfg)
			joCfg.SetProcessorForOrigin(origin, processor)
		case config.TokenOriginCustom:
			processor := NewCustomProcessor(dbSession)
			joCfg.SetProcessorForOrigin(origin, processor)
		}
	}
}
